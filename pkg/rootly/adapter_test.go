// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package rootly_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/config"
	rootly_adapter "github.com/sgnl-ai/adapters/pkg/rootly"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

// testServerHandler is a comprehensive mock HTTP server handler for testing.
func testServerHandler(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	if auth := r.Header.Get("Authorization"); auth != "Bearer testtoken" {
		w.WriteHeader(http.StatusUnauthorized)

		errorResponse := rootly_adapter.DatasourceErrorResponse{
			Errors: []struct {
				Title  string `json:"title"`
				Detail string `json:"detail"`
				Status string `json:"status"`
			}{
				{
					Title:  "Unauthorized",
					Detail: "Invalid authentication credentials",
					Status: "401",
				},
			},
		}
		json.NewEncoder(w).Encode(errorResponse)

		return
	}

	// Handle incidents endpoint
	if r.URL.Path == "/v1/incidents" {
		// Check for filters
		if r.URL.Query().Get("filter[status]") == "started" {
			response := rootly_adapter.DatasourceResponse{
				Data: []map[string]any{
					{
						"id":   "63e2d211-0132-4a12-86f8-d4afe9e666da",
						"type": "incidents",
						"attributes": map[string]any{
							"title":      "Filtered Test Incident",
							"status":     "started",
							"slug":       "filtered-test-incident",
							"created_at": "2025-08-03T20:29:27.326-07:00",
						},
					},
				},
			}
			response.Meta.Page = 1
			response.Meta.Pages = 1
			response.Meta.TotalCount = 1

			w.Header().Set("Content-Type", "application/vnd.api+json")
			json.NewEncoder(w).Encode(response)

			return
		}

		// Default incidents response
		response := rootly_adapter.DatasourceResponse{
			Data: []map[string]any{
				{
					"id":   "incident-1",
					"type": "incidents",
					"attributes": map[string]any{
						"title":      "Test Incident",
						"status":     "open",
						"slug":       "test-incident",
						"created_at": "2025-08-03T20:29:27.326-07:00",
						"updated_at": "2025-08-03T20:29:27.755-07:00",
					},
				},
				{
					"id":   "incident-2",
					"type": "incidents",
					"attributes": map[string]any{
						"title":      "Second Incident",
						"status":     "resolved",
						"slug":       "second-incident",
						"created_at": "2025-08-02T15:20:10.123-07:00",
						"updated_at": "2025-08-03T10:15:30.456-07:00",
					},
				},
			},
		}

		// Handle pagination
		pageNum := r.URL.Query().Get("page[number]")
		if pageNum == "2" {
			response = rootly_adapter.DatasourceResponse{
				Data: []map[string]any{
					{
						"id":   "incident-3",
						"type": "incidents",
						"attributes": map[string]any{
							"title":      "Third Incident",
							"status":     "mitigated",
							"slug":       "third-incident",
							"created_at": "2025-08-01T12:30:45.789-07:00",
						},
					},
				},
			}
			response.Meta.Page = 2
			response.Meta.Pages = 2
			response.Meta.TotalCount = 3
		} else {
			response.Meta.Page = 1
			response.Meta.Pages = 2
			response.Meta.TotalCount = 3
		}

		w.Header().Set("Content-Type", "application/vnd.api+json")
		json.NewEncoder(w).Encode(response)

		return
	}

	// Handle users endpoint
	if r.URL.Path == "/v1/users" {
		response := rootly_adapter.DatasourceResponse{
			Data: []map[string]any{
				{
					"id":   "user-1",
					"type": "users",
					"attributes": map[string]any{
						"name":       "Test User",
						"email":      "test@example.com",
						"created_at": "2025-07-31T16:21:50.168-07:00",
						"updated_at": "2025-07-31T16:21:57.338-07:00",
					},
				},
			},
		}
		response.Meta.Page = 1
		response.Meta.Pages = 1
		response.Meta.TotalCount = 1

		w.Header().Set("Content-Type", "application/vnd.api+json")
		json.NewEncoder(w).Encode(response)

		return
	}

	// Handle teams endpoint
	if r.URL.Path == "/v1/teams" {
		response := rootly_adapter.DatasourceResponse{
			Data: []map[string]any{
				{
					"id":   "team-1",
					"type": "teams",
					"attributes": map[string]any{
						"name":        "Engineering Team",
						"description": "Core engineering team",
						"created_at":  "2025-07-31T16:21:47.120-07:00",
					},
				},
			},
		}
		response.Meta.Page = 1
		response.Meta.Pages = 1
		response.Meta.TotalCount = 1

		w.Header().Set("Content-Type", "application/vnd.api+json")
		json.NewEncoder(w).Encode(response)

		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(testServerHandler))
	defer server.Close()

	adapter := rootly_adapter.NewAdapter(&rootly_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx          context.Context
		request      *framework.Request[rootly_adapter.Config]
		wantResponse framework.Response
		wantError    bool
	}{
		"valid_incidents_request": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(30),
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.attributes.title",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.attributes.status",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.attributes.created_at",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
					},
				},
				PageSize: 100,
				Cursor:   "",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                      "incident-1",
							"$.attributes.title":      "Test Incident",
							"$.attributes.status":     "open",
							"$.attributes.created_at": time.Date(2025, 8, 4, 3, 29, 27, 326000000, time.FixedZone("PDT", -7*3600)),
						},
						{
							"id":                      "incident-2",
							"$.attributes.title":      "Second Incident",
							"$.attributes.status":     "resolved",
							"$.attributes.created_at": time.Date(2025, 8, 2, 22, 20, 10, 123000000, time.FixedZone("PDT", -7*3600)),
						},
					},
					NextCursor: "2",
				},
			},
		},
		"valid_users_request": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(30),
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
							ExternalId: "$.attributes.name",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.attributes.email",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 50,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                 "user-1",
							"$.attributes.name":  "Test User",
							"$.attributes.email": "test@example.com",
						},
					},
					NextCursor: "",
				},
			},
		},
		"valid_request_no_https_prefix": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(30),
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "teams",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.attributes.name",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 25,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                "team-1",
							"$.attributes.name": "Engineering Team",
						},
					},
					NextCursor: "",
				},
			},
		},
		"request_with_pagination": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(30),
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.attributes.title",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 10,
				Cursor:   "2",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                 "incident-3",
							"$.attributes.title": "Third Incident",
						},
					},
					NextCursor: "",
				},
			},
		},
		"request_with_filters": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(30),
					},
					Filters: map[string]string{
						"incidents": "status=started",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.attributes.title",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.attributes.status",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 20,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                  "63e2d211-0132-4a12-86f8-d4afe9e666da",
							"$.attributes.title":  "Filtered Test Incident",
							"$.attributes.status": "started",
						},
					},
					NextCursor: "",
				},
			},
		},
		"unauthorized_request": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer invalidtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(30),
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 10,
			},
			wantError: true,
		},
		"invalid_config": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "2", // Unsupported version
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 10,
			},
			wantError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			response := adapter.GetPage(tt.ctx, tt.request)

			if tt.wantError {
				if response.Error == nil {
					t.Errorf("GetPage() expected error, got success response")
				}

				return
			}

			if response.Error != nil {
				t.Errorf("GetPage() unexpected error: %v", response.Error.Message)

				return
			}

			if response.Success == nil {
				t.Errorf("GetPage() expected success response, got nil")

				return
			}

			// For simple validation tests, just check that we got a response
			if tt.wantResponse.Success == nil {
				return
			}

			// Check objects count
			if len(response.Success.Objects) != len(tt.wantResponse.Success.Objects) {
				t.Errorf("GetPage() objects count = %d, want %d",
					len(response.Success.Objects), len(tt.wantResponse.Success.Objects))

				return
			}

			// Check cursor
			if response.Success.NextCursor != tt.wantResponse.Success.NextCursor {
				t.Errorf("GetPage() nextCursor = %v, want %v",
					response.Success.NextCursor, tt.wantResponse.Success.NextCursor)
			}

			// For detailed validation, check first object
			if len(response.Success.Objects) > 0 && len(tt.wantResponse.Success.Objects) > 0 {
				actualObj := response.Success.Objects[0]
				expectedObj := tt.wantResponse.Success.Objects[0]

				for key, expectedValue := range expectedObj {
					if actualValue, exists := actualObj[key]; !exists {
						t.Errorf("GetPage() object missing key %s", key)
					} else if key != "$.attributes.created_at" && !reflect.DeepEqual(actualValue, expectedValue) {
						t.Errorf("GetPage() object[%s] = %v, want %v", key, actualValue, expectedValue)
					}
				}
			}
		})
	}
}
