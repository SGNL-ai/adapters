// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package okta_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/config"
	okta_adapter "github.com/sgnl-ai/adapters/pkg/okta"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := okta_adapter.NewAdapter(&okta_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[okta_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[okta_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(5),
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                  "00ub0oNGTSWTBKOLGLNR",
							"lastLogin":           time.Date(2013, 6, 24, 17, 39, 19, 0, time.UTC),
							"$.profile.firstName": "Isaac",
						},
						{
							"id":                  "00ub0oNGTSWTBKOCNDJI",
							"lastLogin":           time.Date(2013, 6, 24, 16, 43, 12, 0, time.UTC),
							"$.profile.firstName": "John",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uub2t0YXByZXZpZXcuY29tL2FwaS92MS91c2Vycz9hZnRlcj0xMDB1NjV4dHAzMk5vdkhvUHgxZDdcdTAwMjZsaW1pdD0yIn0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2"),
			},
		},
		"valid_request_no_https_prefix": {
			ctx: context.Background(),
			request: &framework.Request[okta_adapter.Config]{
				Address: strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(5),
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                  "00ub0oNGTSWTBKOLGLNR",
							"lastLogin":           time.Date(2013, 6, 24, 17, 39, 19, 0, time.UTC),
							"$.profile.firstName": "Isaac",
						},
						{
							"id":                  "00ub0oNGTSWTBKOCNDJI",
							"lastLogin":           time.Date(2013, 6, 24, 16, 43, 12, 0, time.UTC),
							"$.profile.firstName": "John",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uub2t0YXByZXZpZXcuY29tL2FwaS92MS91c2Vycz9hZnRlcj0xMDB1NjV4dHAzMk5vdkhvUHgxZDdcdTAwMjZsaW1pdD0yIn0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2"),
			},
		},
		"invalid_request_invalid_api_version": {
			ctx: context.Background(),
			request: &framework.Request[okta_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1.1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Okta config is invalid: apiVersion is not supported: v1.1.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"invalid_request_http_prefix": {
			ctx: context.Background(),
			request: &framework.Request[okta_adapter.Config]{
				Address: "http://" + strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(5),
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "The provided HTTP protocol is not supported.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"valid_request_with_cursor": {
			ctx: context.Background(),
			request: &framework.Request[okta_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(5),
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr(server.URL + "/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2"),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                  "00ub0oNGTSWTBKOMSUFE",
							"lastLogin":           time.Date(2013, 6, 24, 19, 14, 58, 0, time.UTC),
							"$.profile.firstName": "Brooke",
						},
					},
				},
			},
		},
		"invalid_request_invalid_url": {
			ctx: context.Background(),
			request: &framework.Request[okta_adapter.Config]{
				Address: server.URL + "/invalid",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(5),
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
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
		})
	}
}

func TestAdapterGetUserPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := okta_adapter.NewAdapter(&okta_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[okta_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[okta_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(5),
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                  "00ub0oNGTSWTBKOLGLNR",
							"lastLogin":           time.Date(2013, 6, 24, 17, 39, 19, 0, time.UTC),
							"$.profile.firstName": "Isaac",
						},
						{
							"id":                  "00ub0oNGTSWTBKOCNDJI",
							"lastLogin":           time.Date(2013, 6, 24, 16, 43, 12, 0, time.UTC),
							"$.profile.firstName": "John",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uub2t0YXByZXZpZXcuY29tL2FwaS92MS91c2Vycz9hZnRlcj0xMDB1NjV4dHAzMk5vdkhvUHgxZDdcdTAwMjZsaW1pdD0yIn0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2"),
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
		})
	}
}

func TestAdapterGetGroupPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := okta_adapter.NewAdapter(&okta_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[okta_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[okta_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(5),
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "type",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.profile.name",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":             "00g1emaKYZTWRYYRRTSK",
							"type":           "OKTA_GROUP",
							"$.profile.name": "West Coast Users",
						},
						{
							"id":             "00garwpuyxHaWOkdV0g4",
							"type":           "APP_GROUP",
							"$.profile.name": "Engineering Users",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uub2t0YXByZXZpZXcuY29tL2FwaS92MS9ncm91cHM/YWZ0ZXI9MDBnM3p2dWhlcEF3UmVTRG8xZDdcdTAwMjZsaW1pdD0yXHUwMDI2ZmlsdGVyPXR5cGUrZXErJTIyT0tUQV9HUk9VUCUyMitvcit0eXBlK2VxKyUyMkFQUF9HUk9VUCUyMiJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/groups?after=00g3zvuhepAwReSDo1d7&limit=2&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22"),
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
		})
	}
}

func TestAdapterGetGroupMemberPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := okta_adapter.NewAdapter(&okta_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[okta_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"simple": {
			ctx: context.Background(),
			request: &framework.Request[okta_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(5),
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "userId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "groupId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":      "00ub0oNGTSWTBKOLGLNR-00g1emaKYZTWRYYRRTSK",
							"userId":  "00ub0oNGTSWTBKOLGLNR",
							"groupId": "00g1emaKYZTWRYYRRTSK",
						},
						{
							"id":      "00ub0oNGTSWTBKOCNDJI-00g1emaKYZTWRYYRRTSK",
							"userId":  "00ub0oNGTSWTBKOCNDJI",
							"groupId": "00g1emaKYZTWRYYRRTSK",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uub2t0YXByZXZpZXcuY29tL2FwaS92MS9ncm91cHMvMDBnMWVtYUtZWlRXUllZUlJUU0svdXNlcnM/YWZ0ZXI9MDB1YjBvTkdUU1dUQktPQ05ESklcdTAwMjZsaW1pdD0yIiwiY29sbGVjdGlvbklkIjoiMDBnMWVtYUtZWlRXUllZUlJUU0siLCJjb2xsZWN0aW9uQ3Vyc29yIjoiaHR0cHM6Ly90ZXN0LWluc3RhbmNlLm9rdGFwcmV2aWV3LmNvbS9hcGkvdjEvZ3JvdXBzP2FmdGVyPTAwZzFlbWFLWVpUV1JZWVJSVFNLXHUwMDI2bGltaXQ9MVx1MDAyNmZpbHRlcj10eXBlK2VxKyUyMk9LVEFfR1JPVVAlMjIrb3IrdHlwZStlcSslMjJBUFBfR1JPVVAlMjIifQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor:           testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/groups/00g1emaKYZTWRYYRRTSK/users?after=00ub0oNGTSWTBKOCNDJI&limit=2"),
				CollectionID:     testutil.GenPtr("00g1emaKYZTWRYYRRTSK"),
				CollectionCursor: testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/groups?after=00g1emaKYZTWRYYRRTSK&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22"),
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
		})
	}
}
