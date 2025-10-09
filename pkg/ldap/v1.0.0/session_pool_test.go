// Copyright 2025 SGNL.ai, Inc.

// session_pool_mock_test.go
//
// Mock-based tests for session pool and adapter concurrency logic.
// These do not require a real LDAP server and are fast, deterministic, and safe for CI.

package ldap_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapter-framework/pkg/connector"
	ldap_adapter "github.com/sgnl-ai/adapters/pkg/ldap/v1.0.0"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

const (
	bindDNAdmin       = "cn=admin,dc=example,dc=org"
	bindPasswordAdmin = "admin"
	baseDN            = "dc=example,dc=org"
	entityPerson      = "Person"
)

type MockLDAPClient struct {
	mu      sync.Mutex
	calls   int
	cursors []string
}

func (m *MockLDAPClient) Request(_ context.Context, req *ldap_adapter.Request,
) (*ldap_adapter.Response, *framework.Error) {
	m.mu.Lock()
	m.calls++

	if req.Cursor != nil && req.Cursor.Cursor != nil {
		m.cursors = append(m.cursors, *req.Cursor.Cursor)
	}
	m.mu.Unlock()

	dummyCookie := base64.StdEncoding.EncodeToString([]byte("dummy-cookie"))

	return &ldap_adapter.Response{
		StatusCode: 200,
		NextCursor: &pagination.CompositeCursor[string]{Cursor: &dummyCookie},
	}, nil
}

func (m *MockLDAPClient) IsProxied() bool {
	return false
}

func (m *MockLDAPClient) ProxyRequest(_ context.Context, _ *connector.ConnectorInfo, _ *ldap_adapter.Request,
) (*ldap_adapter.Response, *framework.Error) {
	return nil, nil
}

// Helper to create a test LDAP request.
func newTestRequest(baseURL string, cursor *pagination.CompositeCursor[string]) *ldap_adapter.Request {
	return &ldap_adapter.Request{
		ConnectionParams: ldap_adapter.ConnectionParams{
			BindDN:       bindDNAdmin,
			BindPassword: bindPasswordAdmin,
			BaseDN:       baseDN,
		},
		BaseURL:          baseURL,
		PageSize:         1,
		EntityExternalID: entityPerson,
		EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
			entityPerson: {Query: "(objectClass=inetOrgPerson)"},
		},
		Attributes: []*framework.AttributeConfig{},
		Cursor:     cursor,
	}
}

func TestSessionPool_ConcurrentRequests_MockLDAP(t *testing.T) {
	cases := []struct {
		name       string
		withCursor bool
	}{
		{"no cursor", false},
		{"unique cursors", true},
	}

	n := 10

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &MockLDAPClient{}
			client := &ldap_adapter.Datasource{Client: mockClient}
			wg := sync.WaitGroup{}
			uniqueCursors := make(map[string]struct{})

			for i := 0; i < n; i++ {
				var cursor *pagination.CompositeCursor[string]

				if tc.withCursor {
					cursorVal := fmt.Sprintf("cookie-%d", i)
					b64Cursor := base64.StdEncoding.EncodeToString([]byte(cursorVal))
					cursor = &pagination.CompositeCursor[string]{Cursor: &b64Cursor}
					uniqueCursors[b64Cursor] = struct{}{}
				}

				wg.Add(1)

				go func(cursor *pagination.CompositeCursor[string]) {
					defer wg.Done()

					request := newTestRequest("ldap://mock", cursor)

					_, err := client.Client.Request(context.Background(), request)
					if err != nil {
						t.Errorf("mock concurrent request failed: %v", err)
					}
				}(cursor)
			}

			wg.Wait()

			if mockClient.calls != n {
				t.Errorf("expected %d calls, got %d", n, mockClient.calls)
			}

			if tc.withCursor {
				// Assert all unique cursors were received
				received := make(map[string]struct{})
				for _, c := range mockClient.cursors {
					received[c] = struct{}{}
				}

				for c := range uniqueCursors {
					if _, ok := received[c]; !ok {
						t.Errorf("expected cursor %s to be received by mock client", c)
					}
				}
			}
		})
	}
}
