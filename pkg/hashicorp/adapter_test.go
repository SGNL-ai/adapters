// Copyright 2025 SGNL.ai, Inc.
//
//nolint:forcetypeassert
package hashicorp_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	hashicorp_adapter "github.com/sgnl-ai/adapters/pkg/hashicorp"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

const (
	mockUsername = "test_user"
	mockPassword = "test_pass"
	mockToken    = "mock_auth_token"
)

var mockAuth = &framework.DatasourceAuthCredentials{
	Basic: &framework.BasicAuthCredentials{
		Username: mockUsername,
		Password: mockPassword,
	},
}

type testHandler struct {
	users []map[string]interface{}
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check bearer token for non-auth endpoints
	if !strings.Contains(r.URL.Path, "/v1/auth-methods/") {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") || strings.TrimPrefix(authHeader, "Bearer ") != mockToken {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}
	}

	// Handle authentication endpoint
	if strings.Contains(r.URL.Path, "/v1/auth-methods/") && strings.Contains(r.URL.Path, ":authenticate") {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"attributes": {
				"token": "` + mockToken + `"
			},
			"command": "authenticate"
		}`))

		return
	}

	// Handle different test cases based on path and query parameters
	switch {
	case strings.Contains(r.URL.Path, "/v1/auth-methods"):
		h.handleAuthMethodsEndpoint(w)
	case strings.Contains(r.URL.Path, "/v1/hosts"):
		h.handleHostsEndpoint(w, r)
	case strings.Contains(r.URL.Path, "/v1/host-catalogs"):
		h.handleHostCatalogsEndpoint(w)
	case strings.Contains(r.URL.Path, "/v1/accounts"):
		h.handleAccountsEndpoint(w)
	case strings.Contains(r.URL.Path, "/v1/users"):
		h.handleUsersEndpoint(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *testHandler) handleAuthMethodsEndpoint(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{
		"items": [{
			"id": "ampw_123",
			"name": "test-auth-method",
			"type": "password",
			"created_time": "2023-09-29T00:00:00Z"
		}],
		"response_type": "complete",
		"list_token": "next_page_token",
		"sort_by": "created_time",
		"sort_dir": "desc",
		"est_item_count": 1
	}`))
}

func (h *testHandler) handleHostsEndpoint(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	if filter != "" && filter != "name eq \"test-host\"" {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{
		"items": [{
			"id": "hst_123",
			"name": "test-host",
			"created_time": "2023-09-29T00:00:00Z"
		}],
		"response_type": "complete",
		"list_token": "next_page_token",
		"sort_by": "created_time",
		"sort_dir": "desc",
		"est_item_count": 1
	}`))
}

func (h *testHandler) handleHostCatalogsEndpoint(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{
		"items": [{
			"id": "hc_123",
			"name": "test-catalog",
			"created_time": "2023-09-29T00:00:00Z"
		}],
		"response_type": "complete",
		"list_token": "next_page_token",
		"sort_by": "created_time",
		"sort_dir": "desc",
		"est_item_count": 1
	}`))
}

func (h *testHandler) handleAccountsEndpoint(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{
		"items": [{
			"id": "acc_123",
			"login_name": "test-account",
			"created_time": "2023-09-29T00:00:00Z"
		}],
		"response_type": "complete",
		"list_token": "next_page_token",
		"sort_by": "created_time",
		"sort_dir": "desc",
		"est_item_count": 1
	}`))
}

func (h *testHandler) handleUsersEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pageSize := 100

	if size := r.URL.Query().Get("page_size"); size != "" {
		if parsedSize, err := strconv.Atoi(size); err == nil {
			pageSize = parsedSize
		}
	}

	cursor := r.URL.Query().Get("list_token")
	startIndex := 0

	if cursor != "" {
		if parsedIndex, err := strconv.Atoi(cursor); err == nil {
			startIndex = parsedIndex
		}
	}

	endIndex := startIndex + pageSize
	if endIndex > len(h.users) {
		endIndex = len(h.users)
	}

	pageUsers := h.users[startIndex:endIndex]

	response := map[string]interface{}{
		"items":          pageUsers,
		"sort_by":        "created_time",
		"sort_dir":       "desc",
		"est_item_count": len(h.users),
	}

	if endIndex < len(h.users) {
		response["list_token"] = strconv.Itoa(endIndex)
		response["response_type"] = "delta"
	} else {
		response["response_type"] = "complete"
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Write(jsonResponse)
}

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(&testHandler{})
	adapter := hashicorp_adapter.NewAdapter(&hashicorp_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx          context.Context
		request      *framework.Request[hashicorp_adapter.Config]
		wantResponse framework.Response
		wantCursor   *pagination.CompositeCursor[string]
		expectedLogs []map[string]any
	}{
		"valid_request_for_hosts": {
			ctx: context.Background(),
			request: &framework.Request[hashicorp_adapter.Config]{
				Address: server.URL,
				Auth:    mockAuth,
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "ampw_123",
					EntityConfig: map[string]hashicorp_adapter.EntityConfig{
						"hosts": {
							ScopeID: "global",
							Filter:  "name eq \"test-host\"",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "hosts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "created_time",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
					},
				},
				PageSize: 100,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":           "hst_123",
							"name":         "test-host",
							"created_time": time.Date(2023, 9, 29, 0, 0, 0, 0, time.UTC),
						},
					},
					NextCursor: "",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("next_page_token"),
			},
		},
		"valid_request_for_accounts": {
			ctx: context.Background(),
			request: &framework.Request[hashicorp_adapter.Config]{
				Address: server.URL,
				Auth:    mockAuth,
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "ampw_123",
					EntityConfig: map[string]hashicorp_adapter.EntityConfig{
						"accounts": {
							ScopeID: "global",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "login_name",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "created_time",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
					},
				},
				PageSize: 100,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"created_time": time.Date(2023, 9, 29, 0, 0, 0, 0, time.UTC),
							"id":           "acc_123",
							"login_name":   "test-account",
						},
					},
					NextCursor: "",
				},
			},
		},
		"valid_request_with_filter": {
			ctx: context.Background(),
			request: &framework.Request[hashicorp_adapter.Config]{
				Address: server.URL,
				Auth:    mockAuth,
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "ampw_123",
					EntityConfig: map[string]hashicorp_adapter.EntityConfig{
						"hosts": {
							Filter: "name eq \"test-host\"",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "hosts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "created_time",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
					},
				},
				PageSize: 100,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":           "hst_123",
							"name":         "test-host",
							"created_time": time.Date(2023, 9, 29, 0, 0, 0, 0, time.UTC),
						},
					},
					NextCursor: "",
				},
			},
		},
		"valid_request_with_filter_matching_no_results": {
			ctx: context.Background(),
			request: &framework.Request[hashicorp_adapter.Config]{
				Address: server.URL,
				Auth:    mockAuth,
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "ampw_123",
					EntityConfig: map[string]hashicorp_adapter.EntityConfig{
						"hosts": {
							Filter: "name eq \"my-host\"",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "hosts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "created_time",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
					},
				},
				PageSize: 100,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
				},
			},
		},
		"invalid_request_missing_auth": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Address: server.URL,
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "ampw_123",
				},
				Entity: framework.EntityConfig{
					ExternalId: "hosts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Provided datasource auth is missing required http authorization credentials.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"invalid_request_missing_id_attribute": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Address: server.URL,
				Auth:    mockAuth,
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "ampw_123",
				},
				Entity: framework.EntityConfig{
					ExternalId: "hosts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Requested entity attributes are missing unique ID attribute.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				},
			},
		},
		"invalid_request_http_protocol": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Address: "http://example.com",
				Auth:    mockAuth,
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "ampw_123",
				},
				Entity: framework.EntityConfig{
					ExternalId: "hosts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "The provided HTTP protocol is not supported.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"invalid_request_page_size": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Address: server.URL,
				Auth:    mockAuth,
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "ampw_123",
				},
				Entity: framework.EntityConfig{
					ExternalId: "hosts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				PageSize: 5, // Less than minimum page size
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Provided page size (5) does not fall within the allowed range (10-10000).",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctxWithLogger, observedLogs := testutil.NewContextWithObservableLogger(tt.ctx)

			gotResponse := adapter.GetPage(ctxWithLogger, tt.request)

			if diff := cmp.Diff(tt.wantResponse, gotResponse); diff != "" {
				t.Errorf("gotResponse: %v, wantResponse: %v. Diff: %v", gotResponse, tt.wantResponse, diff)
			}

			if gotResponse.Success != nil && tt.wantCursor != nil && gotResponse.Success.NextCursor != "" {
				var gotCursor pagination.CompositeCursor[string]

				decodedCursor, err := base64.StdEncoding.DecodeString(gotResponse.Success.NextCursor)
				if err != nil {
					t.Errorf("error decoding cursor: %v", err)
				}

				if err := json.Unmarshal(decodedCursor, &gotCursor); err != nil {
					t.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(&gotCursor, tt.wantCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}

			testutil.ValidateLogOutput(t, observedLogs, tt.expectedLogs)
		})
	}
}

func TestAdapterGetPageWithPagination(t *testing.T) {
	allUsers := make(map[string]struct{})

	handler := &testHandler{
		users: make([]map[string]interface{}, 0),
	}

	for i := 0; i < 220; i++ {
		user := map[string]interface{}{
			"id":           fmt.Sprintf("usr_%d", i),
			"name":         fmt.Sprintf("User %d", i),
			"created_time": time.Date(2023, 9, 29, 0, 0, i, 0, time.UTC).Format(time.RFC3339),
		}
		handler.users = append(handler.users, user)
		allUsers[user["id"].(string)] = struct{}{}
	}

	server := httptest.NewTLSServer(handler)
	adapter := hashicorp_adapter.NewAdapter(&hashicorp_adapter.Datasource{
		Client: server.Client(),
	})

	ctx := context.Background()
	request := &framework.Request[hashicorp_adapter.Config]{
		Address: server.URL,
		Auth:    mockAuth,
		Config: &hashicorp_adapter.Config{
			AuthMethodID: "ampw_123",
			EntityConfig: map[string]hashicorp_adapter.EntityConfig{
				"users": {
					ScopeID: "global",
				},
			},
		},
		Entity: framework.EntityConfig{
			ExternalId: "users",
			Attributes: []*framework.AttributeConfig{
				{
					ExternalId: "id",
					Type:       framework.AttributeTypeString,
					List:       false,
				},
				{
					ExternalId: "name",
					Type:       framework.AttributeTypeString,
					List:       false,
				},
				{
					ExternalId: "created_time",
					Type:       framework.AttributeTypeDateTime,
					List:       false,
				},
			},
		},
		PageSize: 100,
	}

	// Track retrieved users
	retrievedUsers := make(map[string]struct{})

	// Make requests until we get all users
	for {
		response := adapter.GetPage(ctx, request)
		if response.Error != nil {
			t.Fatalf("Unexpected error: %v", response.Error)
		}

		if response.Success == nil {
			t.Fatal("Expected success response")
		}

		// Add retrieved users to our map
		for _, obj := range response.Success.Objects {
			id := obj["id"].(string)
			retrievedUsers[id] = struct{}{}
		}

		if response.Success.NextCursor == "" {
			break
		}

		request.Cursor = response.Success.NextCursor
	}

	// Verify we got all users
	if len(retrievedUsers) != len(allUsers) {
		t.Errorf("Expected %d users, got %d", len(allUsers), len(retrievedUsers))
	}

	// Verify each user was retrieved
	for id := range allUsers {
		if _, exists := retrievedUsers[id]; !exists {
			t.Errorf("User %s was not retrieved", id)
		}
	}
}
