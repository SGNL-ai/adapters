// Copyright 2025 SGNL.ai, Inc.
package rootly_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	rootly_adapter "github.com/sgnl-ai/adapters/pkg/rootly"
)

// testServerHandler is a mock HTTP server handler for testing
func testServerHandler(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	if auth := r.Header.Get("Authorization"); auth != "Bearer testtoken" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Mock response for users endpoint
	if r.URL.Path == "/v1/users" {
		response := rootly_adapter.DatasourceResponse{
			Data: []map[string]any{
				{
					"id":         "user-1",
					"type":       "users",
					"attributes": map[string]any{
						"name":  "Test User",
						"email": "test@example.com",
					},
				},
			},
		}
		response.Meta.Pagination.Count = 1
		response.Meta.Pagination.Page = 1
		response.Meta.Pagination.Pages = 1
		response.Meta.Pagination.PerPage = 100
		response.Meta.Pagination.TotalCount = 1

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
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "1.0.0",
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
					},
				},
				PageSize: 100,
				Cursor:   "",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Basic test to ensure no panic occurred
			_ = adapter.GetPage(tt.ctx, tt.request)
		})
	}
}