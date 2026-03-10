// Copyright 2026 SGNL.ai, Inc.

package victorops_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
	victorops_adapter "github.com/sgnl-ai/adapters/pkg/victorops"
)

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := victorops_adapter.NewAdapter(&victorops_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx          context.Context
		request      *framework.Request[victorops_adapter.Config]
		wantResponse framework.Response
		wantCursor   *pagination.CompositeCursor[int64]
	}{
		"valid_request_incidents_first_page": {
			ctx: context.Background(),
			request: &framework.Request[victorops_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockAPIId,
						Password: mockAPIKey,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "currentPhase",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "startTime",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"incidentNumber": "1",
							"currentPhase":   "RESOLVED",
							"startTime":      time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
						},
						{
							"incidentNumber": "2",
							"currentPhase":   "ACKED",
							"startTime":      time.Date(2024, 1, 16, 11, 0, 0, 0, time.UTC),
						},
					},
					NextCursor: base64.StdEncoding.EncodeToString([]byte(`{"cursor":2}`)),
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](2),
			},
		},
		"valid_request_incidents_last_page": {
			ctx: context.Background(),
			request: &framework.Request[victorops_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockAPIId,
						Password: mockAPIKey,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 10,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"incidentNumber": "1",
						},
					},
				},
			},
			wantCursor: nil,
		},
		"valid_request_users": {
			ctx: context.Background(),
			request: &framework.Request[victorops_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockAPIId,
						Password: mockAPIKey,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.User,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "username",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "firstName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "email",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 100,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"username":  "user1",
							"firstName": "Alice",
							"email":     "alice@example.com",
						},
						{
							"username":  "user2",
							"firstName": "Bob",
							"email":     "bob@example.com",
						},
					},
				},
			},
			wantCursor: nil,
		},
		"valid_request_incidents_with_query_parameters": {
			ctx: context.Background(),
			request: &framework.Request[victorops_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockAPIId,
						Password: mockAPIKey,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "currentPhase",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
				Config: &victorops_adapter.Config{
					QueryParameters: map[string]string{
						"IncidentReport": "currentPhase=RESOLVED",
					},
				},
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"incidentNumber": "1",
							"currentPhase":   "RESOLVED",
						},
						{
							"incidentNumber": "3",
							"currentPhase":   "RESOLVED",
						},
					},
				},
			},
			wantCursor: nil,
		},
		"valid_request_incidents_with_child_entities": {
			ctx: context.Background(),
			request: &framework.Request[victorops_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockAPIId,
						Password: mockAPIKey,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "currentPhase",
							Type:       framework.AttributeTypeString,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "pagedUsers",
							Attributes: []*framework.AttributeConfig{
								{ExternalId: "id", Type: framework.AttributeTypeString, UniqueId: true},
								{ExternalId: "value", Type: framework.AttributeTypeString},
							},
						},
						{
							ExternalId: "pagedTeams",
							Attributes: []*framework.AttributeConfig{
								{ExternalId: "id", Type: framework.AttributeTypeString, UniqueId: true},
								{ExternalId: "value", Type: framework.AttributeTypeString},
							},
						},
					},
				},
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"incidentNumber": "10",
							"currentPhase":   "RESOLVED",
							"pagedUsers": []framework.Object{
								{"id": "10_pagedUsers_alice", "value": "alice"},
								{"id": "10_pagedUsers_bob", "value": "bob"},
							},
							"pagedTeams": []framework.Object{
								{"id": "10_pagedTeams_team-alpha", "value": "team-alpha"},
							},
						},
					},
				},
			},
			wantCursor: nil,
		},
		"invalid_request_missing_auth": {
			request: &framework.Request[victorops_adapter.Config]{
				Address: "example.com",
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
						},
					},
				},
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "VictorOps auth is missing required basic credentials (API ID and API key).",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"unable_to_decode_cursor": {
			ctx: context.Background(),
			request: &framework.Request[victorops_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockAPIId,
						Password: mockAPIKey,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
				Cursor:   "invalid_cursor",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to decode base64 cursor: illegal base64 data at input byte 7.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
		"failed_to_unmarshal_cursor": {
			ctx: context.Background(),
			request: &framework.Request[victorops_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockAPIId,
						Password: mockAPIKey,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
				Cursor:   base64.StdEncoding.EncodeToString([]byte(`{"cursor": "not_a_number"}`)),
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to unmarshal JSON cursor: json: cannot unmarshal string into Go struct field " +
						"CompositeCursor[int64].cursor of type int64.",
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
		"victorops_request_returns_400": {
			ctx: context.Background(),
			request: &framework.Request[victorops_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockAPIId,
						Password: mockAPIKey,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Datasource rejected request, returned status code: 400.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"failed_to_make_get_page_request_invalid_certs": {
			ctx: context.Background(),
			request: &framework.Request[victorops_adapter.Config]{
				Address: "example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockAPIId,
						Password: mockAPIKey,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.User,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "username",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 100,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: `Failed to execute VictorOps request: Get "https://example.com/api-public/v2/user": ` +
						`tls: failed to verify certificate: x509: certificate signed by unknown authority.`,
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"failed_to_make_get_page_request_invalid_host": {
			ctx: context.Background(),
			request: &framework.Request[victorops_adapter.Config]{
				Address: "https:///example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockAPIId,
						Password: mockAPIKey,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.User,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "username",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 100,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: `Failed to execute VictorOps request: Get "https:///example.com/api-public/v2/user": ` +
						`http: no Host in request URL.`,
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
	}

	// sortChildEntities sorts child entity arrays by their "id" field for deterministic comparison.
	// CreateChildEntitiesFromList iterates a map, so order is not guaranteed.
	sortChildEntities := func(objects []framework.Object) {
		for _, obj := range objects {
			for _, value := range obj {
				if childObjects, ok := value.([]framework.Object); ok {
					sort.Slice(childObjects, func(i, j int) bool {
						id1, _ := childObjects[i]["id"].(string)
						id2, _ := childObjects[j]["id"].(string)

						return id1 < id2
					})
				}
			}
		}
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			// Sort child entity arrays for order-independent comparison.
			if gotResponse.Success != nil {
				sortChildEntities(gotResponse.Success.Objects)
			}

			if tt.wantResponse.Success != nil {
				sortChildEntities(tt.wantResponse.Success.Objects)
			}

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if gotResponse.Success != nil && tt.wantCursor != nil {
				var gotCursor pagination.CompositeCursor[int64]

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
		})
	}
}
