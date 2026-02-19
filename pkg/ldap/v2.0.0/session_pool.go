// Copyright 2026 SGNL.ai, Inc.

package ldap

import (
	"crypto/tls"
	"strings"
	"sync"
	"time"

	ldap_v3 "github.com/go-ldap/ldap/v3"
)

type Session struct {
	conn     interface{ Close() error }
	currKey  string    // currently active session key
	newKey   string    // new session key (used during pagination)
	lastUsed time.Time // timestamp of last access to manage TTL-based cleanup
	mu       sync.Mutex
}

// GetOrCreateConn retrieves the existing LDAP connection if it's healthy,
// or creates a new one if it doesn't exist or is unhealthy. It uses the provided
// address, TLS configuration, and bind credentials to establish the connection.
// The session's lastUsed timestamp is updated on each access to facilitate
// TTL-based cleanup in the session pool.
func (s *Session) GetOrCreateConn(
	address string,
	tlsConfig *tls.Config,
	bindDN, bindPassword string,
) (*ldap_v3.Conn, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if c, ok := s.conn.(*ldap_v3.Conn); ok && c != nil {
		// Health check: WhoAmI
		_, err := c.WhoAmI(nil)
		if err == nil {
			s.lastUsed = time.Now()

			return c, nil
		}
		// If WhoAmI fails, close and reset
		c.Close()

		s.conn = nil
	}

	// Dial and bind new connection
	conn, err := ldap_v3.DialURL(
		address,
		ldap_v3.DialWithTLSConfig(tlsConfig),
	)
	if err != nil {
		return nil, err
	}

	if err := conn.Bind(bindDN, bindPassword); err != nil {
		conn.Close()

		return nil, err
	}

	s.conn = conn
	s.lastUsed = time.Now()

	return conn, nil
}

type SessionPool struct {
	mu              sync.Mutex
	pool            map[string]*Session
	ttl             time.Duration
	cleanupInterval time.Duration
	done            chan struct{}
	closeOnce       sync.Once
}

func NewSessionPool(ttl, cleanupInterval time.Duration) *SessionPool {
	sp := &SessionPool{
		pool:            make(map[string]*Session),
		ttl:             ttl,
		cleanupInterval: cleanupInterval,
		done:            make(chan struct{}),
	}
	sp.startCleanupLoop()

	return sp
}

// Close stops the cleanup goroutine and closes all connections in the pool.
// Connections are closed after releasing the lock to avoid blocking other operations.
// Safe to call multiple times - subsequent calls are no-ops.
func (sp *SessionPool) Close() {
	sp.closeOnce.Do(func() {
		close(sp.done)

		var toClose []interface{ Close() error }

		sp.mu.Lock()

		for _, session := range sp.pool {
			if session != nil && session.conn != nil {
				toClose = append(toClose, session.conn)
			}
		}

		sp.pool = make(map[string]*Session)

		sp.mu.Unlock()

		for _, conn := range toClose {
			conn.Close()
		}
	})
}

// Get retrieves a session from the pool by key. If the session exists, it updates
// the lastUsed timestamp. It also manages the promotion of newKey to currKey when
// accessed via newKey, and cleanup of alternate keys to ensure that only one active
// key (currKey) is maintained for each session. Note that sessions accessed via
// non-cookie keys (keys ending with "|") are not cleaned up since they are used
// for non-paged queries and should be reusable.
func (sp *SessionPool) Get(key string) (*Session, bool) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	s, ok := sp.pool[key]
	if ok {
		if s == nil || s.conn == nil {
			delete(sp.pool, key)

			return nil, false
		}

		// Clean up the alternate key that wasn't used.
		// This determines which cookie (old vs new) the caller is using.
		// Exception: Don't clean up when accessed via non-cookie key (ends with "|")
		// because internal non-paged queries shouldn't affect paged query keys.
		if key == s.currKey && s.newKey != "" && !strings.HasSuffix(key, "|") {
			// Accessed via current key (with cookie) - discard new key
			delete(sp.pool, s.newKey)
			s.newKey = ""
		} else if key == s.newKey {
			// Accessed via new key - promote new to current
			delete(sp.pool, s.currKey)
			s.currKey = s.newKey
			s.newKey = ""
		}

		s.mu.Lock()
		s.lastUsed = time.Now()
		s.mu.Unlock()
	}

	return s, ok
}

// Set adds a new session to the pool with the given key. If a session already
// exists for the key, it is replaced and its connection is closed. The session's
// currKey is set to the provided key.
func (sp *SessionPool) Set(key string, session *Session) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if old, ok := sp.pool[key]; ok && old != nil {
		// Clean up old session's newKey entry if it exists
		if old.newKey != "" {
			delete(sp.pool, old.newKey)
		}

		if old.conn != nil {
			old.conn.Close()
		}
	}

	session.currKey = key
	sp.pool[key] = session
}

// Delete removes a session from the pool and closes its connection if it exists.
// It deletes both the current and new keys associated with the session to ensure
// complete cleanup. However, it does not delete sessions that are accessed via
// non-cookie keys (keys ending with "|") since those are used for non-paged queries
// and should be reusable.
func (sp *SessionPool) Delete(key string) {
	// Only delete sessions that have a cookie (paging was in progress).
	// Keys ending with "|" have no cookie - keep these sessions for reuse.
	if strings.HasSuffix(key, "|") {
		return
	}

	sp.mu.Lock()
	defer sp.mu.Unlock()

	if s, ok := sp.pool[key]; ok {
		// Delete both keys
		delete(sp.pool, s.currKey)

		if s.newKey != "" {
			delete(sp.pool, s.newKey)
		}

		if s.conn != nil {
			s.conn.Close()
			s.conn = nil
		}
	}
}

func (sp *SessionPool) startCleanupLoop() {
	go func() {
		ticker := time.NewTicker(sp.cleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-sp.done:
				return
			case <-ticker.C:
				sp.cleanupExpiredSessions()
			}
		}
	}()
}

// cleanupExpiredSessions removes sessions that have exceeded their TTL.
// It deletes all keys associated with each expired session and closes its connection.
// Connections are closed after releasing the lock to avoid blocking other operations.
func (sp *SessionPool) cleanupExpiredSessions() {
	now := time.Now()

	var toClose []interface{ Close() error }

	sp.mu.Lock()

	for key, session := range sp.pool {
		session.mu.Lock()
		expired := now.Sub(session.lastUsed) > sp.ttl
		session.mu.Unlock()

		if !expired {
			continue
		}

		sp.deleteSessionKeys(key, session)

		if session.conn != nil {
			toClose = append(toClose, session.conn)
			session.conn = nil
		}
	}

	sp.mu.Unlock()

	for _, conn := range toClose {
		conn.Close()
	}
}

// deleteSessionKeys removes all keys associated with a session from the pool.
// Caller must hold sp.mu.
func (sp *SessionPool) deleteSessionKeys(key string, session *Session) {
	delete(sp.pool, key)

	if session.currKey != "" && session.currKey != key {
		delete(sp.pool, session.currKey)
	}

	if session.newKey != "" && session.newKey != key {
		delete(sp.pool, session.newKey)
	}
}

// UpdateKey updates a new session key to an existing session using currKey in the
// pool. This is used when a new cookie is received in the middle of pagination
// and we want to associate it with the existing session.
func (sp *SessionPool) UpdateKey(currKey, newKey string) {
	if currKey == newKey {
		return
	}

	sp.mu.Lock()
	defer sp.mu.Unlock()

	s, ok := sp.pool[currKey]
	if ok {
		s.newKey = newKey
		sp.pool[newKey] = s
	}
}

func (sp *SessionPool) SessionCount() int {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	return len(sp.pool)
}

// UniqueSessionCount returns the number of unique sessions in the pool.
// This differs from SessionCount() which returns the number of keys.
// With dual-key design, multiple keys can point to the same session.
func (sp *SessionPool) UniqueSessionCount() int {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	seen := make(map[*Session]struct{})
	for _, s := range sp.pool {
		seen[s] = struct{}{}
	}

	return len(seen)
}

// Example usage:
//   pool := NewSessionPool(10 * time.Minute)
//   key := sessionKey(address, cookie)
//   session, found := pool.Get(key)
//   ...
