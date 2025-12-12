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
		// Check for test scenario query parameter
		scenario := r.URL.Query().Get("filter[scenario]")

		// Scenario: incidents with included relationships
		if scenario == "with_included" {
			response := rootly_adapter.DatasourceResponse{
				Data: []map[string]any{
					{
						"id":   "incident-with-included",
						"type": "incidents",
						"attributes": map[string]any{
							"title":      "Incident with Relationships",
							"status":     "open",
							"slug":       "incident-with-relationships",
							"created_at": "2025-08-03T20:29:27.326-07:00",
							"severity":   "sev1",
							"priority":   "p1",
						},
					},
				},
				Included: []map[string]any{
					{
						"type": "form_field_values",
						"attributes": map[string]any{
							"incident_id":   "incident-with-included",
							"form_field_id": "field-users",
							"selected_users": []any{
								map[string]any{"id": "user-100", "name": "Alice Smith", "email": "alice@example.com"},
								map[string]any{"id": "user-101", "name": "Bob Jones", "email": "bob@example.com"},
							},
						},
					},
					{
						"type": "form_field_values",
						"attributes": map[string]any{
							"incident_id":   "incident-with-included",
							"form_field_id": "field-services",
							"selected_services": []any{
								map[string]any{"id": "service-1", "name": "API Service", "slug": "api-service"},
								map[string]any{"id": "service-2", "name": "Database Service", "slug": "db-service"},
							},
						},
					},
					{
						"type": "form_field_values",
						"attributes": map[string]any{
							"incident_id":   "incident-with-included",
							"form_field_id": "field-groups",
							"selected_groups": []any{
								map[string]any{"id": "group-1", "name": "On-Call Team", "slug": "on-call"},
							},
						},
					},
					{
						"type": "form_field_values",
						"attributes": map[string]any{
							"incident_id":   "incident-with-included",
							"form_field_id": "field-functionalities",
							"selected_functionalities": []any{
								map[string]any{"id": "func-1", "name": "Authentication", "slug": "auth"},
								map[string]any{"id": "func-2", "name": "Payment Processing", "slug": "payment"},
							},
						},
					},
					{
						"type": "form_field_values",
						"attributes": map[string]any{
							"incident_id":   "incident-with-included",
							"form_field_id": "ad5366c5-2680-4656-913b-331433284941",
							"value":         "high",
						},
					},
					{
						"type": "incident_form_field_selections",
						"attributes": map[string]any{
							"incident_id":   "incident-with-included",
							"form_field_id": "uuid-object-format-test",
							"selected_users": map[string]any{
								"id":    nil,
								"value": "Test-Value-Object-Format",
							},
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

		// Scenario: incidents with empty included
		if scenario == "no_included" {
			response := rootly_adapter.DatasourceResponse{
				Data: []map[string]any{
					{
						"id":   "incident-no-included",
						"type": "incidents",
						"attributes": map[string]any{
							"title":      "Incident without Relationships",
							"status":     "resolved",
							"created_at": "2025-08-03T20:29:27.326-07:00",
						},
					},
				},
				Included: []map[string]any{},
			}
			response.Meta.Page = 1
			response.Meta.Pages = 1
			response.Meta.TotalCount = 1

			w.Header().Set("Content-Type", "application/vnd.api+json")
			json.NewEncoder(w).Encode(response)

			return
		}

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
		"incidents_with_included_relationships": {
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
						"incidents": "scenario=with_included",
					},
					Includes: map[string]string{
						"incidents": "form_field_values",
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
							ExternalId: "$.attributes.severity",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.attributes.priority",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: `$.all_selected_users[?(@.field_id=="field-users")].id`,
							Type:       framework.AttributeTypeString,
							List:       true,
						},
					},
				},
				PageSize: 100,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                    "incident-with-included",
							"$.attributes.title":    "Incident with Relationships",
							"$.attributes.status":   "open",
							"$.attributes.severity": "sev1",
							"$.attributes.priority": "p1",
							// Test complex JSONPath filter expression for enriched data
							`$.all_selected_users[?(@.field_id=="field-users")].id`: []string{"user-100", "user-101"},
						},
					},
					NextCursor: "",
				},
			},
		},
		"query_impact_value_from_included": {
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
						"incidents": "scenario=with_included",
					},
					Includes: map[string]string{
						"incidents": "form_field_values",
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
							// Query using nested path (original structure)
							ExternalId: `$.included[?(@.attributes.form_field_id=="ad5366c5-2680-4656-913b-331433284941")].attributes.value`,
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							// Query using flattened form_field_id
							ExternalId: `$.included[?(@.form_field_id=="ad5366c5-2680-4656-913b-331433284941")].attributes.value`,
							Type:       framework.AttributeTypeString,
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
							"id":                 "incident-with-included",
							"$.attributes.title": "Incident with Relationships",
							// Both JSONPath queries should return the same value
							`$.included[?(@.attributes.form_field_id=="ad5366c5-2680-4656-913b-331433284941")].attributes.value`: "high",
							`$.included[?(@.form_field_id=="ad5366c5-2680-4656-913b-331433284941")].attributes.value`:            "high",
						},
					},
					NextCursor: "",
				},
			},
		},
		"incidents_with_no_included": {
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
						"incidents": "scenario=no_included",
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
				PageSize: 50,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                  "incident-no-included",
							"$.attributes.title":  "Incident without Relationships",
							"$.attributes.status": "resolved",
						},
					},
					NextCursor: "",
				},
			},
		},
		"query_object_format_selected_users_by_uuid": {
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
						"incidents": "scenario=with_included",
					},
					Includes: map[string]string{
						"incidents": "form_field_values",
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
							// Query object-format selected_users by UUID field_id
							ExternalId: `$.all_selected_users[?(@.field_id=="uuid-object-format-test")].value`,
							Type:       framework.AttributeTypeString,
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
							"id": "incident-with-included",
							`$.all_selected_users[?(@.field_id=="uuid-object-format-test")].value`: "Test-Value-Object-Format",
						},
					},
					NextCursor: "",
				},
			},
		},
		"query_array_format_selected_groups_by_field_id": {
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
						"incidents": "scenario=with_included",
					},
					Includes: map[string]string{
						"incidents": "form_field_values",
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
							// Query array-format selected_groups by field_id (returns list of IDs)
							ExternalId: `$.all_selected_groups[?(@.field_id=="field-groups")].id`,
							Type:       framework.AttributeTypeString,
							List:       true,
						},
						{
							// Query array-format selected_groups by field_id (returns list of names)
							ExternalId: `$.all_selected_groups[?(@.field_id=="field-groups")].name`,
							Type:       framework.AttributeTypeString,
							List:       true,
						},
					},
				},
				PageSize: 100,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id": "incident-with-included",
							`$.all_selected_groups[?(@.field_id=="field-groups")].id`:   []string{"group-1"},
							`$.all_selected_groups[?(@.field_id=="field-groups")].name`: []string{"On-Call Team"},
						},
					},
					NextCursor: "",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			// Test setup is done in the test case struct

			// Act
			response := adapter.GetPage(tt.ctx, tt.request)

			// Assert
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
