// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package duo_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	duo_adapter "github.com/sgnl-ai/adapters/pkg/duo"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := duo_adapter.NewAdapter(&duo_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[duo_adapter.Config]
		inputRequestCursor interface{}
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[int64]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "group_id",
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
				PageSize: 3,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"group_id": "DGKUKQSTG7ZFDN2N1XID",
							"name":     "group1",
						},
						{
							"group_id": "DGIB125DJLJKYZ9W257F",
							"name":     "group2",
						},
						{
							"group_id": "DG36ABPJ1T3RZDL7ISLC",
							"name":     "group3",
						},
					},
					NextCursor: "eyJjdXJzb3IiOjN9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](3),
			},
		},
		"valid_request_no_https_prefix": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "group_id",
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
				PageSize: 3,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"group_id": "DGKUKQSTG7ZFDN2N1XID",
							"name":     "group1",
						},
						{
							"group_id": "DGIB125DJLJKYZ9W257F",
							"name":     "group2",
						},
						{
							"group_id": "DG36ABPJ1T3RZDL7ISLC",
							"name":     "group3",
						},
					},
					NextCursor: "eyJjdXJzb3IiOjN9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](3),
			},
		},
		"invalid_request_invalid_api_version": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v2",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "group_id",
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
				PageSize: 3,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Duo config is invalid: apiVersion is not supported: v2.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"invalid_request_http_prefix": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: "http://" + strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "group_id",
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
				PageSize: 3,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: `Scheme "http" is not supported.`,
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"malformed_cursor_negative_offset": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "group_id",
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
				PageSize: 3,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](-50),
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Cursor must be greater than 0.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
		"malformed_composite_cursor_string_type": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "group_id",
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
				PageSize: 3,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("BROKEN"),
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to unmarshal JSON cursor: json: cannot unmarshal string into Go struct field CompositeCursor[int64].cursor of type int64.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
		"malformed_cursor_includes_collection_cursor": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "group_id",
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
				PageSize: 3,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor:           testutil.GenPtr[int64](-50),
				CollectionCursor: testutil.GenPtr[int64](10),
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Cursor must not contain CollectionID or CollectionCursor fields for entity Group.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
		"valid_request_with_cursor": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "group_id",
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
				PageSize: 3,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](3),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"group_id": "DG6IHDSWM72IJJNXBA82",
							"name":     "group4",
						},
						{
							"group_id": "DGKQMVO91JT365VY36MU",
							"name":     "group5",
						},
					},
				},
			},
		},
		"invalid_request_invalid_url": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL + "/invalid",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "group_id",
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
				PageSize: 3,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Datasource rejected request, returned status code: 404.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				var encodedCursor string

				var err *framework.Error

				switch v := tt.inputRequestCursor.(type) {
				case *pagination.CompositeCursor[int64]:
					encodedCursor, err = pagination.MarshalCursor(v)
				case *pagination.CompositeCursor[string]:
					encodedCursor, err = pagination.MarshalCursor(v)
				default:
					t.Errorf("Unsupported cursor type: %T", tt.inputRequestCursor)
				}

				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if diff := cmp.Diff(gotResponse, tt.wantResponse); diff != "" {
				t.Errorf("adapter.GetPage() mismatch (-want +got):\n%s", diff)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
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

func TestAdapterGetUserPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := duo_adapter.NewAdapter(&duo_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[duo_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[int64]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[int64]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "user_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "realname",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "created",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
						{
							ExternalId: "is_enrolled",
							Type:       framework.AttributeTypeBool,
							List:       false,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "groups",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "group_id",
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
						{
							ExternalId: "phones",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "phone_id",
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
					},
				},
				PageSize: 4,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"created": time.Date(2024, 1, 23, 20, 17, 36, 0, time.UTC),
							"groups": []framework.Object{
								{"group_id": "DGKUKQSTG7ZFDN2N1XID", "name": "group1"},
							},
							"is_enrolled": true,
							"phones": []framework.Object{
								{"name": "", "phone_id": "DPFL36P8Z8LZANN1FFEZ"},
							},
							"realname": "Test User 1",
							"user_id":  "DUYC8O4O953VBGGKLHAL",
						},
						{
							"created": time.Date(2024, 1, 23, 20, 17, 36, 0, time.UTC),
							"groups": []framework.Object{
								{"group_id": "DGKUKQSTG7ZFDN2N1XID", "name": "group1"},
								{"group_id": "DGIB125DJLJKYZ9W257F", "name": "group2"},
							},
							"is_enrolled": true,
							"phones": []framework.Object{
								{"name": "", "phone_id": "DPFL36P8Z8LZANN1FFEZ"},
							},
							"realname": "Test User 2",
							"user_id":  "DUHUTX7KGB6D15WTD3VY",
						},
						{
							"created": time.Date(2024, 1, 23, 20, 17, 36, 0, time.UTC),
							"groups": []framework.Object{
								{"group_id": "DGKUKQSTG7ZFDN2N1XID", "name": "group1"},
							},
							"is_enrolled": true,
							"phones": []framework.Object{
								{"name": "", "phone_id": "DPFL36P8Z8LZANN1FFEZ"},
							},
							"realname": "Test User 3",
							"user_id":  "DUEL2SL4CWLP04CE71SL",
						},
						{
							"created": time.Date(2024, 1, 23, 20, 17, 36, 0, time.UTC),
							"groups": []framework.Object{
								{"group_id": "DGKUKQSTG7ZFDN2N1XID", "name": "group1"},
								{"group_id": "DGIB125DJLJKYZ9W257F", "name": "group2"},
							},
							"is_enrolled": true,
							"phones": []framework.Object{
								{"name": "", "phone_id": "DPFL36P8Z8LZANN1FFEZ"},
							},
							"realname": "Test User 4",
							"user_id":  "DUB3BH17CE2V7B744RLI",
						},
					},
					NextCursor: "eyJjdXJzb3IiOjR9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](4),
			},
		},
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "user_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "realname",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "created",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
						{
							ExternalId: "is_enrolled",
							Type:       framework.AttributeTypeBool,
							List:       false,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "groups",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "group_id",
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
						{
							ExternalId: "phones",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "phone_id",
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
					},
				},
				PageSize: 4,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](4),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"created": time.Date(2024, 1, 23, 20, 17, 36, 0, time.UTC),
							"groups": []framework.Object{
								{"group_id": "DGKUKQSTG7ZFDN2N1XID", "name": "group1"},
							},
							"is_enrolled": true,
							"phones": []framework.Object{
								{"name": "", "phone_id": "DPFL36P8Z8LZANN1FFEZ"},
							},
							"realname": "Test User 5",
							"user_id":  "DUWC7NXJX7IM9I7J26AT",
						},
						{
							"created": time.Date(2024, 1, 23, 20, 17, 36, 0, time.UTC),
							"groups": []framework.Object{
								{"group_id": "DGKUKQSTG7ZFDN2N1XID", "name": "group1"},
							},
							"is_enrolled": false,
							"realname":    "Test User 6",
							"user_id":     "DUQ87KL4A6OU5VYMWWLT",
						},
						{
							"created": time.Date(2024, 1, 23, 20, 17, 37, 0, time.UTC),
							"groups": []framework.Object{
								{"group_id": "DGKUKQSTG7ZFDN2N1XID", "name": "group1"},
							},
							"is_enrolled": true,
							"phones": []framework.Object{
								{"name": "", "phone_id": "DPFL36P8Z8LZANN1FFEZ"},
							},
							"realname": "Test User 7",
							"user_id":  "DU9SRK429IRM2J7OEDCP",
						},
						{
							"created":     time.Date(2024, 1, 23, 20, 17, 37, 0, time.UTC),
							"is_enrolled": false,
							"realname":    "Test User 8",
							"user_id":     "DU1E2CSOI6I5HEO043WN",
						},
					},
					NextCursor: "eyJjdXJzb3IiOjh9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](8),
			},
		},
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "user_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "realname",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "created",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
						{
							ExternalId: "is_enrolled",
							Type:       framework.AttributeTypeBool,
							List:       false,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "groups",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "group_id",
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
						{
							ExternalId: "phones",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "phone_id",
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
					},
				},
				PageSize: 4,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](8),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"created":     time.Date(2024, 1, 23, 20, 17, 37, 0, time.UTC),
							"is_enrolled": true,
							"phones": []framework.Object{
								{"name": "", "phone_id": "DPFL36P8Z8LZANN1FFEZ"},
								{"name": "", "phone_id": "DPX0H7ZWQLSB735FEHVY"},
							},
							"realname": "Test User 9",
							"user_id":  "DU2T7B5VIC0RSCN1A13W",
						},
						{
							"created":     time.Date(2024, 1, 23, 20, 17, 37, 0, time.UTC),
							"is_enrolled": true,
							"phones": []framework.Object{
								{"name": "", "phone_id": "DPFL36P8Z8LZANN1FFEZ"},
								{"name": "", "phone_id": "DP7MW6K4G1OVMP8DTI08"},
							},
							"realname": "Test User 10",
							"user_id":  "DUG1B8MRABMVKYVCFO8H",
						},
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
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

func TestAdapterGetGroupPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := duo_adapter.NewAdapter(&duo_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[duo_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[int64]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[int64]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "group_id",
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
				PageSize: 3,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"group_id": "DGKUKQSTG7ZFDN2N1XID",
							"name":     "group1",
						},
						{
							"group_id": "DGIB125DJLJKYZ9W257F",
							"name":     "group2",
						},
						{
							"group_id": "DG36ABPJ1T3RZDL7ISLC",
							"name":     "group3",
						},
					},
					NextCursor: "eyJjdXJzb3IiOjN9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](3),
			},
		},
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "group_id",
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
				PageSize: 3,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](3),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"group_id": "DG6IHDSWM72IJJNXBA82",
							"name":     "group4",
						},
						{
							"group_id": "DGKQMVO91JT365VY36MU",
							"name":     "group5",
						},
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
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

func TestAdapterGetPhonePage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := duo_adapter.NewAdapter(&duo_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[duo_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[int64]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[int64]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Phone",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "phone_id",
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
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{"name": "", "phone_id": "DPFL36P8Z8LZANN1FFEZ"},
						{"name": "", "phone_id": "DPX0H7ZWQLSB735FEHVY"},
					},
					NextCursor: "eyJjdXJzb3IiOjJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](2),
			},
		},
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Phone",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "phone_id",
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
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](2),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{"name": "", "phone_id": "DP7MW6K4G1OVMP8DTI08"},
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
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

func TestAdapterGetEndpointPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := duo_adapter.NewAdapter(&duo_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[duo_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[int64]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[int64]
	}{
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[duo_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Endpoint",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "epkey",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "device_identifier",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"device_identifier": "3FA47335-1976-3BED-8286-D3F1ABCDEA13",
							"epkey":             "EP18JX1A10AB102M2T2X",
						},
						{
							"device_identifier": "",
							"epkey":             "EP65MWZWXA10AB1027TQ",
						},
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
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
