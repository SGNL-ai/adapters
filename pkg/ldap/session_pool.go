// Copyright 2025 SGNL.ai, Inc.

package ldap

import (
	"sync"
	"time"

	"crypto/tls"

	ldap_v3 "github.com/go-ldap/ldap/v3"
)

type Session struct {
	conn     interface{ Close() error }
	lastUsed time.Time
	mu       sync.Mutex
}

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
	mu   sync.Mutex
	pool map[string]*Session
	ttl  time.Duration
}

func NewSessionPool(ttl time.Duration) *SessionPool {
	sp := &SessionPool{
		pool: make(map[string]*Session),
		ttl:  ttl,
	}
	sp.startCleanupLoop()

	return sp
}

func (sp *SessionPool) Get(key string) (*Session, bool) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	s, ok := sp.pool[key]
	if ok {
		if s == nil || s.conn == nil {
			delete(sp.pool, key)

			return nil, false
		}

		s.lastUsed = time.Now()
	}

	return s, ok
}

func (sp *SessionPool) Set(key string, session *Session) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if old, ok := sp.pool[key]; ok && old != nil && old.conn != nil {
		old.conn.Close()
	}

	sp.pool[key] = session
}

func (sp *SessionPool) Delete(key string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if s, ok := sp.pool[key]; ok {
		if s.conn != nil {
			s.conn.Close()
		}

		delete(sp.pool, key)
	}
}

func (sp *SessionPool) startCleanupLoop() {
	go func() {
		// Use shorter cleanup interval (ttl/4) to check for expired sessions more frequently
		cleanupInterval := sp.ttl / 4
		if cleanupInterval < time.Minute {
			cleanupInterval = time.Minute
		}
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()

			sp.mu.Lock()
			for key, session := range sp.pool {
				if now.Sub(session.lastUsed) > sp.ttl {
					if session.conn != nil {
						session.conn.Close()
					}

					delete(sp.pool, key)
				}
			}
			sp.mu.Unlock()
		}
	}()
}

func (sp *SessionPool) UpdateKey(oldKey, newKey string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	s, ok := sp.pool[oldKey]
	if ok {
		sp.pool[newKey] = s
		delete(sp.pool, oldKey)
	}
}

func (sp *SessionPool) SessionCount() int {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	return len(sp.pool)
}

// Example usage:
//   pool := NewSessionPool(10 * time.Minute)
//   key := sessionKey(address, cookie)
//   session, found := pool.Get(key)
//   ...
