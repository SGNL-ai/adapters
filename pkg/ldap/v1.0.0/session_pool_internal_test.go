// Copyright 2025 SGNL.ai, Inc.

// session_pool_internal_test.go
//
// Pure unit tests for the LDAP session pool logic (set/get/delete, TTL, cleanup, thread safety, etc.).
// These tests use a testConn and require access to unexported fields, so they are in package ldap.

package ldap

import (
	"sync"
	"testing"
	"time"
)

const (
	keyAddressCookie  = "address|cookie"
	keyAddressCookie1 = "address|cookie1"
	keyAddressCookie2 = "address|cookie2"
)

type testConn struct {
	closed bool
	mu     sync.Mutex
}

func (m *testConn) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true

	return nil
}

func (m *testConn) IsClosed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.closed
}

func TestSessionPool_SetGetDelete(t *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	key := keyAddressCookie
	conn := &testConn{}
	session := &Session{conn: conn, lastUsed: time.Now()}

	pool.Set(key, session)

	if got, ok := pool.Get(key); !ok || got != session {
		t.Fatalf("expected to get the same session back")
	}

	// Act
	pool.Delete(key)

	// Assert
	if _, ok := pool.Get(key); ok {
		t.Fatalf("expected session to be deleted")
	}

	if !conn.IsClosed() {
		t.Fatalf("expected connection to be closed on delete")
	}
}

func TestSessionPool_TTLExpiry(t *testing.T) {
	// Arrange
	pool := NewSessionPool(100*time.Millisecond, 10*time.Millisecond)
	key := keyAddressCookie
	conn := &testConn{}
	session := &Session{
		conn:     conn,
		lastUsed: time.Now().Add(-11 * time.Minute),
	}
	pool.Set(key, session)

	// Act
	// Wait a bit longer than the cleanup interval to ensure the goroutine runs
	time.Sleep(1 * time.Second)

	// Assert
	if _, ok := pool.Get(key); ok {
		t.Fatalf("expected session to be expired and deleted")
	}

	if !conn.IsClosed() {
		t.Fatalf("expected connection to be closed on TTL expiry")
	}
}

func TestSessionPool_ThreadSafety(_ *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	key := keyAddressCookie
	wg := sync.WaitGroup{}
	conn := &testConn{}
	session := &Session{conn: conn, lastUsed: time.Now()}

	// Act, Assert
	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			pool.Set(key, session)
			pool.Get(key)
			pool.Delete(key)
			pool.Set(key, session)
		}()
	}

	wg.Wait()
}

func TestSessionPool_GetOnNonexistentKey(t *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	key := "nonexistent|cookie"

	// Act
	session, ok := pool.Get(key)

	// Assert
	if ok || session != nil {
		t.Fatalf("expected Get on nonexistent key to return (nil, false)")
	}
}

func TestSessionPool_SetOverwritesExistingSession(t *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	key := keyAddressCookie

	conn1 := &testConn{}
	session1 := &Session{
		conn:     conn1,
		lastUsed: time.Now(),
	}
	pool.Set(key, session1)

	conn2 := &testConn{}
	session2 := &Session{
		conn:     conn2,
		lastUsed: time.Now(),
	}

	// Act
	pool.Set(key, session2)

	// Assert
	got, ok := pool.Get(key)
	if !ok || got != session2 {
		t.Fatalf("expected to get the new session after overwrite")
	}

	if !conn1.IsClosed() {
		t.Fatalf("expected old connection to be closed on overwrite")
	}
}

func TestSessionPool_DeleteOnNonexistentKey(_ *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	key := "nonexistent|cookie"

	// Act, Assert
	// Should not panic or error; nothing to check.
	pool.Delete(key)
}

func TestSessionPool_UpdateKeyMovesSession(t *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	oldKey := keyAddressCookie1
	newKey := keyAddressCookie2
	conn := &testConn{}
	session := &Session{conn: conn, lastUsed: time.Now()}
	pool.Set(oldKey, session)

	// Act
	pool.UpdateKey(oldKey, newKey)

	// Assert
	if _, ok := pool.Get(oldKey); ok {
		t.Fatalf("expected old key to be gone after UpdateKey")
	}

	got, ok := pool.Get(newKey)
	if !ok || got != session {
		t.Fatalf("expected session to be moved to new key")
	}
}

func TestSessionPool_UpdateKeyOnNonexistentKey(t *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	oldKey := "doesnotexist"
	newKey := "newkey"

	// Act
	pool.UpdateKey(oldKey, newKey)

	// Assert
	if _, ok := pool.Get(newKey); ok {
		t.Fatalf("expected new key to not exist after UpdateKey on nonexistent old key")
	}
}

func TestSessionPool_GetRemovesNilConn(t *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	key := keyAddressCookie
	session := &Session{conn: nil, lastUsed: time.Now()}

	pool.Set(key, session)

	// Act
	got, ok := pool.Get(key)

	// Asser	t
	if ok || got != nil {
		t.Fatalf("expected Get to remove key and return (nil, false) if conn is nil")
	}
}

func TestSessionPool_MultipleKeys(t *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	key1 := keyAddressCookie1
	key2 := keyAddressCookie2
	conn1 := &testConn{}
	conn2 := &testConn{}

	session1 := &Session{
		conn:     conn1,
		lastUsed: time.Now(),
	}
	session2 := &Session{
		conn:     conn2,
		lastUsed: time.Now(),
	}

	pool.Set(key1, session1)
	pool.Set(key2, session2)

	// Act
	got1, ok1 := pool.Get(key1)
	got2, ok2 := pool.Get(key2)

	// Assert
	if !ok1 || got1 != session1 {
		t.Fatalf("expected to get session: %v for key: %v", session1, key1)
	}

	if !ok2 || got2 != session2 {
		t.Fatalf("expected to get session: %v for key: %v", session2, key2)
	}
}

func TestSessionPool_CleanupClosesAllExpired(t *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	key1 := keyAddressCookie1
	key2 := keyAddressCookie2
	conn1 := &testConn{}
	conn2 := &testConn{}
	session1 := &Session{
		conn:     conn1,
		lastUsed: time.Now().Add(-11 * time.Minute),
	}
	session2 := &Session{
		conn:     conn2,
		lastUsed: time.Now().Add(-12 * time.Minute),
	}

	pool.Set(key1, session1)
	pool.Set(key2, session2)

	// Act
	pool.mu.Lock()
	for key, session := range pool.pool {
		if time.Since(session.lastUsed) > pool.ttl {
			if session.conn != nil {
				session.conn.Close()
			}

			delete(pool.pool, key)
		}
	}
	pool.mu.Unlock()

	// Assert
	if _, ok := pool.Get(key1); ok {
		t.Fatalf("expected session: %v to be expired and deleted", session1)
	}

	if _, ok := pool.Get(key2); ok {
		t.Fatalf("expected session: %v to be expired and deleted", session2)
	}

	if !conn1.IsClosed() || !conn2.IsClosed() {
		t.Fatalf("expected both connections to be closed on cleanup")
	}
}

func TestSessionPool_ReuseAcrossPages(t *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	key := keyAddressCookie
	conn := &testConn{}
	session := &Session{
		conn:     conn,
		lastUsed: time.Now(),
	}

	pool.Set(key, session)

	// Act
	first, ok1 := pool.Get(key)
	if !ok1 || first != session {
		t.Fatalf("expected to get the same session for first page")
	}

	second, ok2 := pool.Get(key)
	if !ok2 || second != session {
		t.Fatalf("expected to get the same session for next page")
	}

	// Assert
	if first != second {
		t.Fatalf("expected session to be reused across pages")
	}
}
