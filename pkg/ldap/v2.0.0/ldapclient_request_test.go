// Copyright 2026 SGNL.ai, Inc.

package ldap

import (
	"context"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

// ldapClientRequestTestSuite provides integration tests for ldapClient.Request
// using a real OpenLDAP server via testcontainers.
type ldapClientRequestTestSuite struct {
	testutil.CommonSuite
	ldapContainer testcontainers.Container
	ldapHost      string
	ctx           context.Context
	cancel        context.CancelFunc
}

func Test_LdapClientRequestSuite(t *testing.T) {
	testutil.Run(t, new(ldapClientRequestTestSuite))
}

func (s *ldapClientRequestTestSuite) SetupSuite() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), time.Minute*5)

	var ldapPort nat.Port
	s.ldapContainer, ldapPort = s.StartLDAPServer(s.ctx, false)
	s.ldapHost = "ldap://localhost:" + ldapPort.Port()

	time.Sleep(10 * time.Second)
}

func (s *ldapClientRequestTestSuite) TearDownSuite() {
	if s.cancel != nil {
		s.cancel()
	}

	if s.ldapContainer != nil {
		s.ldapContainer.Terminate(s.ctx)
	}
}

func (s *ldapClientRequestTestSuite) newRequest(cursor *pagination.CompositeCursor[string]) *Request {
	return &Request{
		BaseURL: s.ldapHost,
		ConnectionParams: ConnectionParams{
			BaseDN:       "dc=example,dc=org",
			BindDN:       "cn=admin,dc=example,dc=org",
			BindPassword: "admin",
		},
		PageSize:         2, // Small page size to force pagination
		EntityExternalID: "Person",
		EntityConfigMap: map[string]*EntityConfig{
			"Person": {Query: "(objectClass=inetOrgPerson)"},
		},
		Attributes: []*framework.AttributeConfig{
			{ExternalId: "dn", Type: framework.AttributeTypeString, UniqueId: true},
			{ExternalId: "cn", Type: framework.AttributeTypeString},
		},
		Cursor: cursor,
	}
}

// Test_GivenNonPagedRequests_WhenSameAddress_ThenSessionIsReused verifies that
// multiple non-paged requests to the same address reuse the same session.
func (s *ldapClientRequestTestSuite) Test_GivenNonPagedRequests_WhenSameAddress_ThenSessionIsReused() {
	// Arrange
	pool := NewSessionPool(5*time.Minute, time.Minute)
	client := &ldapClient{sessionPool: pool}
	request := s.newRequest(nil)

	// Act - Make first request
	resp1, err1 := client.Request(s.ctx, request)

	// Assert first request succeeded
	s.Require().Nil(err1, "first request should succeed")
	s.Require().NotNil(resp1, "first response should not be nil")
	uniqueSessionsAfterFirst := pool.UniqueSessionCount()
	s.Assert().Equal(1, uniqueSessionsAfterFirst, "should have 1 unique session after first request")

	// Act - Make second request (no cursor, same address)
	resp2, err2 := client.Request(s.ctx, request)

	// Assert second request succeeded and reused session
	s.Require().Nil(err2, "second request should succeed")
	s.Require().NotNil(resp2, "second response should not be nil")
	s.Assert().Equal(1, pool.UniqueSessionCount(), "should still have 1 unique session after reuse")
}

// Test_GivenPagedRequest_WhenNewCookieReceived_ThenSessionKeyIsUpdated verifies that
// when a paged response contains a new cookie, the session key is updated.
func (s *ldapClientRequestTestSuite) Test_GivenPagedRequest_WhenNewCookieReceived_ThenSessionKeyIsUpdated() {
	// Arrange
	pool := NewSessionPool(5*time.Minute, time.Minute)
	client := &ldapClient{sessionPool: pool}
	request := s.newRequest(nil)

	// Act - Make first request (should get a cursor if there are more results)
	resp1, err1 := client.Request(s.ctx, request)

	// Assert
	s.Require().Nil(err1, "first request should succeed")
	s.Require().NotNil(resp1, "first response should not be nil")
	s.Assert().Equal(1, pool.UniqueSessionCount(), "should have 1 unique session")

	// If there's a next cursor, there should be 2 keys (currKey + newKey)
	// pointing to the same session
	if resp1.NextCursor != nil && resp1.NextCursor.Cursor != nil {
		s.T().Log("Response has next cursor - session key was updated")
		s.Assert().Equal(2, pool.SessionCount(), "should have 2 keys (currKey + newKey)")
		s.Assert().Equal(1, pool.UniqueSessionCount(), "but only 1 unique session")
	}
}

// Test_GivenPagedSequence_WhenMultiplePages_ThenSessionIsReusedAcrossPages verifies
// that the same session is reused across multiple pages of results.
func (s *ldapClientRequestTestSuite) Test_GivenPagedSequence_WhenMultiplePages_ThenSessionIsReusedAcrossPages() {
	// Arrange
	pool := NewSessionPool(5*time.Minute, time.Minute)
	client := &ldapClient{sessionPool: pool}

	// Act - Page through all results
	var cursor *pagination.CompositeCursor[string]
	pageCount := 0
	maxPages := 10 // Safety limit

	for pageCount < maxPages {
		request := s.newRequest(cursor)
		resp, err := client.Request(s.ctx, request)

		s.Require().Nil(err, "request for page %d should succeed", pageCount+1)
		s.Require().NotNil(resp, "response for page %d should not be nil", pageCount+1)

		pageCount++

		// Should always have exactly 1 unique session (reused across pages)
		s.Assert().Equal(1, pool.UniqueSessionCount(),
			"should have exactly 1 unique session after page %d", pageCount)

		if resp.NextCursor == nil || resp.NextCursor.Cursor == nil {
			break // No more pages
		}

		cursor = resp.NextCursor
	}

	s.T().Logf("Completed %d pages with session reuse", pageCount)
	s.Assert().Equal(1, pool.UniqueSessionCount(), "should have exactly 1 unique session after all pages")
}

// TestSessionKey verifies the sessionKey function generates correct keys.
func TestSessionKey(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		cookie   []byte
		expected string
	}{
		{
			name:     "no cookie",
			address:  "ldap://server.com",
			cookie:   nil,
			expected: "ldap://server.com|",
		},
		{
			name:     "empty cookie",
			address:  "ldap://server.com",
			cookie:   []byte{},
			expected: "ldap://server.com|",
		},
		{
			name:     "with cookie",
			address:  "ldap://server.com",
			cookie:   []byte{0x01, 0x02, 0x03},
			expected: "ldap://server.com|AQID",
		},
		{
			name:     "ldaps address",
			address:  "ldaps://secure.server.com:636",
			cookie:   []byte("test-cookie"),
			expected: "ldaps://secure.server.com:636|dGVzdC1jb29raWU=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sessionKey(tt.address, tt.cookie)
			require.Equal(t, tt.expected, result)
		})
	}
}

// TestSessionKey_NonCookieVsCookieKey verifies the distinction between
// keys ending with "|" (non-cookie) and keys with data after "|" (cookie).
func TestSessionKey_NonCookieVsCookieKey(t *testing.T) {
	t.Run("NonCookieKeyEndsWithPipe", func(t *testing.T) {
		// Arrange & Act
		key := sessionKey("ldap://server.com", nil)

		// Assert
		assert.Equal(t, "ldap://server.com|", key)
		assert.True(t, len(key) > 0 && key[len(key)-1] == '|',
			"non-cookie key should end with |")
	})

	t.Run("CookieKeyHasDataAfterPipe", func(t *testing.T) {
		// Arrange & Act
		key := sessionKey("ldap://server.com", []byte("cookie-data"))

		// Assert
		assert.Contains(t, key, "|")
		parts := splitKey(key)
		assert.Equal(t, 2, len(parts))
		assert.NotEmpty(t, parts[1], "cookie key should have data after |")
	})
}

// splitKey is a helper to split a session key into address and cookie parts.
func splitKey(key string) []string {
	for i := len(key) - 1; i >= 0; i-- {
		if key[i] == '|' {
			return []string{key[:i], key[i+1:]}
		}
	}

	return []string{key}
}
