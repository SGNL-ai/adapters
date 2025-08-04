// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst, revive
package rootly_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	rootly_adapter "github.com/sgnl-ai/adapters/pkg/rootly"
)

func TestDatasourceGetPage(t *testing.T) {
	tests := map[string]struct {
		serverHandler    http.HandlerFunc
		request          *rootly_adapter.Request
		expectedResponse *rootly_adapter.Response
		expectedError    *framework.Error
	}{
		"successful_incidents_request": {
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "Bearer testtoken" {
					w.WriteHeader(http.StatusUnauthorized)

					return
				}
				if r.URL.Path != "/v1/incidents" {
					w.WriteHeader(http.StatusNotFound)

					return
				}

				response := rootly_adapter.DatasourceResponse{
					Data: []map[string]any{
						{
							"id":   "63e2d211-0132-4a12-86f8-d4afe9e666da",
							"type": "incidents",
							"attributes": map[string]any{
								"title":  "Test Incident",
								"status": "started",
								"slug":   "test-incident",
							},
						},
					},
				}
				response.Meta.Page = 1
				response.Meta.Pages = 1
				response.Meta.TotalCount = 1

				w.Header().Set("Content-Type", "application/vnd.api+json")
				json.NewEncoder(w).Encode(response)
			},
			request: &rootly_adapter.Request{
				HTTPAuthorization:     "Bearer testtoken",
				EntityExternalName:    "incidents",
				PageSize:              10,
				RequestTimeoutSeconds: 30,
			},
			expectedResponse: &rootly_adapter.Response{
				Objects: []map[string]any{
					{
						"id":   "63e2d211-0132-4a12-86f8-d4afe9e666da",
						"type": "incidents",
						"attributes": map[string]any{
							"title":  "Test Incident",
							"status": "started",
							"slug":   "test-incident",
						},
					},
				},
				NextCursor: nil,
			},
		},
		"successful_users_request": {
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "Bearer testtoken" {
					w.WriteHeader(http.StatusUnauthorized)

					return
				}
				if r.URL.Path != "/v1/users" {
					w.WriteHeader(http.StatusNotFound)

					return
				}

				response := rootly_adapter.DatasourceResponse{
					Data: []map[string]any{
						{
							"id":   "116641",
							"type": "users",
							"attributes": map[string]any{
								"name":  "Test User",
								"email": "test@example.com",
							},
						},
					},
				}
				response.Meta.Page = 1
				response.Meta.Pages = 2
				response.Meta.TotalCount = 15

				w.Header().Set("Content-Type", "application/vnd.api+json")
				json.NewEncoder(w).Encode(response)
			},
			request: &rootly_adapter.Request{
				HTTPAuthorization:     "Bearer testtoken",
				EntityExternalName:    "users",
				PageSize:              10,
				RequestTimeoutSeconds: 30,
			},
			expectedResponse: &rootly_adapter.Response{
				Objects: []map[string]any{
					{
						"id":   "116641",
						"type": "users",
						"attributes": map[string]any{
							"name":  "Test User",
							"email": "test@example.com",
						},
					},
				},
				NextCursor: strPtr("2"),
			},
		},
		"pagination_second_page": {
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "Bearer testtoken" {
					w.WriteHeader(http.StatusUnauthorized)

					return
				}
				if r.URL.Query().Get("page[number]") != "2" {
					w.WriteHeader(http.StatusBadRequest)

					return
				}

				response := rootly_adapter.DatasourceResponse{
					Data: []map[string]any{
						{
							"id":   "116642",
							"type": "users",
							"attributes": map[string]any{
								"name":  "Second User",
								"email": "second@example.com",
							},
						},
					},
				}
				response.Meta.Page = 2
				response.Meta.Pages = 2
				response.Meta.TotalCount = 15

				w.Header().Set("Content-Type", "application/vnd.api+json")
				json.NewEncoder(w).Encode(response)
			},
			request: &rootly_adapter.Request{
				HTTPAuthorization:     "Bearer testtoken",
				EntityExternalName:    "users",
				PageSize:              10,
				Cursor:                strPtr("2"),
				RequestTimeoutSeconds: 30,
			},
			expectedResponse: &rootly_adapter.Response{
				Objects: []map[string]any{
					{
						"id":   "116642",
						"type": "users",
						"attributes": map[string]any{
							"name":  "Second User",
							"email": "second@example.com",
						},
					},
				},
				NextCursor: nil,
			},
		},
		"request_with_filters": {
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "Bearer testtoken" {
					w.WriteHeader(http.StatusUnauthorized)

					return
				}

				// Check that filters are properly formatted
				if r.URL.Query().Get("filter[status]") != "started" {
					w.WriteHeader(http.StatusBadRequest)

					return
				}
				if r.URL.Query().Get("filter[severity]") != "high" {
					w.WriteHeader(http.StatusBadRequest)

					return
				}

				response := rootly_adapter.DatasourceResponse{
					Data: []map[string]any{
						{
							"id":   "filtered-incident",
							"type": "incidents",
							"attributes": map[string]any{
								"title":    "Filtered Incident",
								"status":   "started",
								"severity": "high",
							},
						},
					},
				}
				response.Meta.Page = 1
				response.Meta.Pages = 1
				response.Meta.TotalCount = 1

				w.Header().Set("Content-Type", "application/vnd.api+json")
				json.NewEncoder(w).Encode(response)
			},
			request: &rootly_adapter.Request{
				HTTPAuthorization:     "Bearer testtoken",
				EntityExternalName:    "incidents",
				PageSize:              10,
				Filter:                "status=started&severity=high",
				RequestTimeoutSeconds: 30,
			},
			expectedResponse: &rootly_adapter.Response{
				Objects: []map[string]any{
					{
						"id":   "filtered-incident",
						"type": "incidents",
						"attributes": map[string]any{
							"title":    "Filtered Incident",
							"status":   "started",
							"severity": "high",
						},
					},
				},
				NextCursor: nil,
			},
		},
		"unauthorized_request": {
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
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
			},
			request: &rootly_adapter.Request{
				HTTPAuthorization:     "Bearer invalidtoken",
				EntityExternalName:    "incidents",
				PageSize:              10,
				RequestTimeoutSeconds: 30,
			},
			expectedError: &framework.Error{
				Message: "HTTP 401: Invalid authentication credentials",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"not_found_entity": {
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "Bearer testtoken" {
					w.WriteHeader(http.StatusUnauthorized)

					return
				}
				w.WriteHeader(http.StatusNotFound)
				errorResponse := rootly_adapter.DatasourceErrorResponse{
					Errors: []struct {
						Title  string `json:"title"`
						Detail string `json:"detail"`
						Status string `json:"status"`
					}{
						{
							Title:  "Not Found",
							Detail: "The requested resource was not found",
							Status: "404",
						},
					},
				}
				json.NewEncoder(w).Encode(errorResponse)
			},
			request: &rootly_adapter.Request{
				HTTPAuthorization:     "Bearer testtoken",
				EntityExternalName:    "nonexistent",
				PageSize:              10,
				RequestTimeoutSeconds: 30,
			},
			expectedError: &framework.Error{
				Message: "HTTP 404: The requested resource was not found",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"malformed_json_response": {
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "Bearer testtoken" {
					w.WriteHeader(http.StatusUnauthorized)

					return
				}
				w.Header().Set("Content-Type", "application/vnd.api+json")
				w.WriteHeader(http.StatusOK)           // Add explicit 200 status
				w.Write([]byte(`{"invalid": "json"}`)) // Missing required fields
			},
			request: &rootly_adapter.Request{
				HTTPAuthorization:     "Bearer testtoken",
				EntityExternalName:    "incidents",
				PageSize:              10,
				RequestTimeoutSeconds: 30,
			},
			expectedError: &framework.Error{
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				Message: "Invalid response format: missing required data field. Body: {\"invalid\": \"json\"}.",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewTLSServer(tt.serverHandler)
			defer server.Close()

			datasource := &rootly_adapter.Datasource{
				Client: server.Client(),
			}

			tt.request.BaseURL = server.URL + "/v1"

			response, err := datasource.GetPage(context.Background(), tt.request)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("GetPage() expected error, got nil")

					return
				}

				if err.Code != tt.expectedError.Code {
					t.Errorf("GetPage() error code = %v, want %v", err.Code, tt.expectedError.Code)
				}

				if tt.expectedError.Message != "" && err.Message != tt.expectedError.Message {
					t.Errorf("GetPage() error message = %v, want %v", err.Message, tt.expectedError.Message)
				}

				return
			}

			if err != nil {
				t.Errorf("GetPage() unexpected error: %v", err)

				return
			}

			if response == nil {
				t.Errorf("GetPage() response is nil")

				return
			}

			if !reflect.DeepEqual(response.Objects, tt.expectedResponse.Objects) {
				t.Errorf("GetPage() objects = %v, want %v", response.Objects, tt.expectedResponse.Objects)
			}

			if !reflect.DeepEqual(response.NextCursor, tt.expectedResponse.NextCursor) {
				t.Errorf("GetPage() nextCursor = %v, want %v", response.NextCursor, tt.expectedResponse.NextCursor)
			}
		})
	}
}

func TestDatasourceConstructEndpoint(t *testing.T) {
	tests := map[string]struct {
		request     *rootly_adapter.Request
		expectedURL string
	}{
		"basic_request": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           10,
			},
			expectedURL: "https://api.rootly.com/v1/incidents?page%5Bnumber%5D=1&page%5Bsize%5D=10",
		},
		"request_with_cursor": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "users",
				PageSize:           25,
				Cursor:             strPtr("3"),
			},
			expectedURL: "https://api.rootly.com/v1/users?page%5Bnumber%5D=3&page%5Bsize%5D=25",
		},
		"request_with_filters": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           50,
				Filter:             "status=open&severity=high",
			},
			expectedURL: "https://api.rootly.com/v1/incidents?filter%5Bseverity%5D=high&filter%5Bstatus%5D=open&page%5Bnumber%5D=1&page%5Bsize%5D=50",
		},
		"request_with_complex_filters": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           100,
				Filter:             "status=started,mitigated&severity=major,minor",
			},
			expectedURL: "https://api.rootly.com/v1/incidents?filter%5Bseverity%5D=major%2Cminor&filter%5Bstatus%5D=started%2Cmitigated&page%5Bnumber%5D=1&page%5Bsize%5D=100",
		},
		"request_no_page_size": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "teams",
				PageSize:           0,
			},
			expectedURL: "https://api.rootly.com/v1/teams?page%5Bnumber%5D=1",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actualURL := rootly_adapter.ConstructEndpoint(tt.request)
			if actualURL != tt.expectedURL {
				t.Errorf("ConstructEndpoint() = %v, want %v", actualURL, tt.expectedURL)
			}
		})
	}
}

// Helper function to create string pointers.
func strPtr(s string) *string {
	return &s
}
