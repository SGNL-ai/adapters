// Copyright 2026 SGNL.ai, Inc.

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
	keyAddressCookie    = "address|cookie"
	keyAddressCookie1   = "address|cookie1"
	keyAddressCookie2   = "address|cookie2"
	keyAddressNoCookie  = "address|"           // Key without cookie (non-paged query)
	keyAddressNoCookie2 = "ldap://server.com|" // Another key without cookie
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

	// Assert: Both keys should exist after UpdateKey (dual-key design)
	// The session can be accessed via either key until one is used
	gotViaNew, okNew := pool.Get(newKey)
	if !okNew || gotViaNew != session {
		t.Fatalf("expected session to be accessible via new key")
	}

	// After accessing via new key, old key should be cleaned up
	if _, ok := pool.Get(oldKey); ok {
		t.Fatalf("expected old key to be gone after accessing via new key")
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

// TestSessionPool_DualKey_BothKeysExistAfterUpdateKey verifies that after
// UpdateKey, both old and new keys point to the same session until one is accessed.
func TestSessionPool_DualKey_BothKeysExistAfterUpdateKey(t *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	oldKey := keyAddressCookie1
	newKey := keyAddressCookie2
	conn := &testConn{}
	session := &Session{conn: conn, lastUsed: time.Now()}
	pool.Set(oldKey, session)

	// Act
	pool.UpdateKey(oldKey, newKey)

	// Assert: Both keys should exist in the pool
	pool.mu.Lock()
	_, oldExists := pool.pool[oldKey]
	_, newExists := pool.pool[newKey]
	pool.mu.Unlock()

	if !oldExists {
		t.Fatalf("expected old key to still exist after UpdateKey")
	}

	if !newExists {
		t.Fatalf("expected new key to exist after UpdateKey")
	}

	// Verify session fields are set correctly
	if session.currKey != oldKey {
		t.Fatalf("expected currKey to be %q, got %q", oldKey, session.currKey)
	}

	if session.newKey != newKey {
		t.Fatalf("expected newKey to be %q, got %q", newKey, session.newKey)
	}
}

// TestSessionPool_DualKey_AccessCleansUpUnusedKey verifies that accessing
// a session via one key cleans up the other key.
func TestSessionPool_DualKey_AccessCleansUpUnusedKey(t *testing.T) {
	tests := []struct {
		name          string
		accessKey     string // Which key to access
		otherKey      string // The key that should be cleaned up
		expectCurrKey string // Expected currKey after access
		expectNewKey  string // Expected newKey after access (should be empty)
	}{
		{
			name:          "access_via_old_key",
			accessKey:     keyAddressCookie1,
			otherKey:      keyAddressCookie2,
			expectCurrKey: keyAddressCookie1,
			expectNewKey:  "",
		},
		{
			name:          "access_via_new_key",
			accessKey:     keyAddressCookie2,
			otherKey:      keyAddressCookie1,
			expectCurrKey: keyAddressCookie2, // newKey promoted to currKey
			expectNewKey:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			pool := NewSessionPool(1*time.Minute, time.Minute)
			conn := &testConn{}
			session := &Session{conn: conn, lastUsed: time.Now()}
			pool.Set(keyAddressCookie1, session)
			pool.UpdateKey(keyAddressCookie1, keyAddressCookie2)

			// Act
			got, ok := pool.Get(tt.accessKey)

			// Assert
			if !ok || got != session {
				t.Fatalf("expected session to be accessible via %s", tt.accessKey)
			}

			// Other key should be cleaned up
			if _, ok := pool.Get(tt.otherKey); ok {
				t.Fatalf("expected %s to be gone after accessing via %s", tt.otherKey, tt.accessKey)
			}

			// Verify session state
			if session.currKey != tt.expectCurrKey {
				t.Fatalf("expected currKey=%q, got %q", tt.expectCurrKey, session.currKey)
			}

			if session.newKey != tt.expectNewKey {
				t.Fatalf("expected newKey=%q, got %q", tt.expectNewKey, session.newKey)
			}
		})
	}
}

// TestSessionPool_Delete_BehaviorByKeyType verifies Delete behavior based on
// whether the key has a cookie or not.
func TestSessionPool_Delete_BehaviorByKeyType(t *testing.T) {
	tests := []struct {
		name             string
		key              string
		expectDeleted    bool
		expectConnClosed bool
	}{
		{
			name:             "empty_cookie_key_not_deleted",
			key:              keyAddressNoCookie,
			expectDeleted:    false,
			expectConnClosed: false,
		},
		{
			name:             "another_empty_cookie_key_not_deleted",
			key:              keyAddressNoCookie2,
			expectDeleted:    false,
			expectConnClosed: false,
		},
		{
			name:             "non_empty_cookie_key_deleted",
			key:              keyAddressCookie,
			expectDeleted:    true,
			expectConnClosed: true,
		},
		{
			name:             "another_non_empty_cookie_key_deleted",
			key:              keyAddressCookie1,
			expectDeleted:    true,
			expectConnClosed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			pool := NewSessionPool(1*time.Minute, time.Minute)
			conn := &testConn{}
			session := &Session{conn: conn, lastUsed: time.Now()}
			pool.Set(tt.key, session)

			// Act
			pool.Delete(tt.key)

			// Assert
			_, exists := pool.Get(tt.key)

			if tt.expectDeleted && exists {
				t.Fatalf("expected session to be deleted for key %q", tt.key)
			}

			if !tt.expectDeleted && !exists {
				t.Fatalf("expected session to NOT be deleted for key %q", tt.key)
			}

			if tt.expectConnClosed && !conn.IsClosed() {
				t.Fatalf("expected connection to be closed for key %q", tt.key)
			}

			if !tt.expectConnClosed && conn.IsClosed() {
				t.Fatalf("expected connection to NOT be closed for key %q", tt.key)
			}
		})
	}
}

// TestSessionPool_Delete_DeletesBothKeysWhenDualKeyExists verifies that
// Delete removes both keys when a session has dual keys.
func TestSessionPool_Delete_DeletesBothKeysWhenDualKeyExists(t *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	oldKey := keyAddressCookie1
	newKey := keyAddressCookie2
	conn := &testConn{}
	session := &Session{conn: conn, lastUsed: time.Now()}
	pool.Set(oldKey, session)
	pool.UpdateKey(oldKey, newKey)

	// Verify both keys exist
	pool.mu.Lock()
	oldExists := pool.pool[oldKey] != nil
	newExists := pool.pool[newKey] != nil
	pool.mu.Unlock()

	if !oldExists || !newExists {
		t.Fatalf("expected both keys to exist before delete")
	}

	// Act: Delete using the new key
	pool.Delete(newKey)

	// Assert: Both keys should be gone
	if _, ok := pool.Get(oldKey); ok {
		t.Fatalf("expected old key to be deleted when deleting via new key")
	}

	if _, ok := pool.Get(newKey); ok {
		t.Fatalf("expected new key to be deleted")
	}

	if !conn.IsClosed() {
		t.Fatalf("expected connection to be closed")
	}
}

// TestSessionPool_SessionReuse_NonPagedQueries verifies that sessions with
// empty cookie keys are reused for subsequent non-paged queries.
func TestSessionPool_SessionReuse_NonPagedQueries(t *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	key := keyAddressNoCookie
	conn := &testConn{}
	session := &Session{conn: conn, lastUsed: time.Now()}
	pool.Set(key, session)

	// Simulate multiple non-paged queries using the same key
	first, ok1 := pool.Get(key)
	if !ok1 {
		t.Fatalf("expected first Get to succeed")
	}

	// Simulate query completion - Delete is called but should be no-op
	pool.Delete(key)

	second, ok2 := pool.Get(key)
	if !ok2 {
		t.Fatalf("expected second Get to succeed after Delete on empty cookie key")
	}

	// Assert
	if first != second {
		t.Fatalf("expected same session to be reused for non-paged queries")
	}

	if conn.IsClosed() {
		t.Fatalf("expected connection to remain open for non-paged queries")
	}
}

// TestSessionPool_TTLExpiry_CleansBothDualKeys verifies that TTL cleanup
// removes both keys when a session with dual keys expires.
func TestSessionPool_TTLExpiry_CleansBothDualKeys(t *testing.T) {
	// Arrange
	pool := NewSessionPool(100*time.Millisecond, 10*time.Millisecond)
	oldKey := keyAddressCookie1
	newKey := keyAddressCookie2
	conn := &testConn{}
	session := &Session{
		conn:     conn,
		lastUsed: time.Now().Add(-11 * time.Minute), // Already expired
	}
	pool.Set(oldKey, session)
	pool.UpdateKey(oldKey, newKey)

	// Act: Wait for cleanup loop to run
	time.Sleep(500 * time.Millisecond)

	// Assert: Both keys should be gone
	if _, ok := pool.Get(oldKey); ok {
		t.Fatalf("expected old key to be expired and deleted")
	}

	if _, ok := pool.Get(newKey); ok {
		t.Fatalf("expected new key to be expired and deleted")
	}

	if !conn.IsClosed() {
		t.Fatalf("expected connection to be closed on TTL expiry")
	}
}

// TestSessionPool_UpdateKey_SameKeyNoOp verifies that UpdateKey with the
// same old and new key is a no-op.
func TestSessionPool_UpdateKey_SameKeyNoOp(t *testing.T) {
	// Arrange
	pool := NewSessionPool(1*time.Minute, time.Minute)
	key := keyAddressCookie
	conn := &testConn{}
	session := &Session{conn: conn, lastUsed: time.Now()}
	pool.Set(key, session)

	// Act
	pool.UpdateKey(key, key)

	// Assert: Session should still be accessible
	got, ok := pool.Get(key)
	if !ok || got != session {
		t.Fatalf("expected session to be accessible after UpdateKey with same key")
	}

	// newKey should not be set
	if session.newKey != "" {
		t.Fatalf("expected newKey to remain empty when UpdateKey uses same key")
	}
}
