// Copyright 2026 SGNL.ai, Inc.

// nolint: lll, goconst
package azuread_test

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
	azuread_adapter "github.com/sgnl-ai/adapters/pkg/azuread"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := azuread_adapter.NewAdapter(&azuread_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[azuread_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[azuread_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
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
							"id": "6e7b768e-07e2-4810-8459-485f84f8f204",
						},
						{
							"id": "87d349ed-44d7-43e1-9a83-5f2406dee5bd",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20vdjEuMC91c2Vycz8kc2VsZWN0PWlkXHUwMDI2JHRvcD0yXHUwMDI2JHNraXB0b2tlbj1SRk53ZEFJQUFRQUFBQ002UVdSbGJHVldRRTB6TmpWNE1qRTBNelUxTG05dWJXbGpjbTl6YjJaMExtTnZiU2xWYzJWeVh6ZzNaRE0wT1dWa0xUUTBaRGN0TkRObE1TMDVZVGd6TFRWbU1qUXdObVJsWlRWaVpMa0FBQUFBQUFBQUFBQUEifQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/users?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA"),
			},
		},
		"valid_request_no_https_prefix": {
			ctx: context.Background(),
			request: &framework.Request[azuread_adapter.Config]{
				Address: strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
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
							"id": "6e7b768e-07e2-4810-8459-485f84f8f204",
						},
						{
							"id": "87d349ed-44d7-43e1-9a83-5f2406dee5bd",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20vdjEuMC91c2Vycz8kc2VsZWN0PWlkXHUwMDI2JHRvcD0yXHUwMDI2JHNraXB0b2tlbj1SRk53ZEFJQUFRQUFBQ002UVdSbGJHVldRRTB6TmpWNE1qRTBNelUxTG05dWJXbGpjbTl6YjJaMExtTnZiU2xWYzJWeVh6ZzNaRE0wT1dWa0xUUTBaRGN0TkRObE1TMDVZVGd6TFRWbU1qUXdObVJsWlRWaVpMa0FBQUFBQUFBQUFBQUEifQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/users?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA"),
			},
		},
		"invalid_request_invalid_api_version": {
			ctx: context.Background(),
			request: &framework.Request[azuread_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Azure AD config is invalid: apiVersion is not supported: v1.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"invalid_request_http_prefix": {
			ctx: context.Background(),
			request: &framework.Request[azuread_adapter.Config]{
				Address: "http://" + strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
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
			request: &framework.Request[azuread_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr(server.URL + "/v1.0/users?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA"),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id": "5bde3e51-d13b-4db1-9948-fe4b109d11a7",
						},
						{
							"id": "4782e723-f4f4-4af3-a76e-25e3bab0d896",
						},
					},
				},
			},
		},
		"invalid_request_invalid_url": {
			ctx: context.Background(),
			request: &framework.Request[azuread_adapter.Config]{
				Address: server.URL + "/invalid",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
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
	adapter := azuread_adapter.NewAdapter(&azuread_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[azuread_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[azuread_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
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
							"id": "6e7b768e-07e2-4810-8459-485f84f8f204",
						},
						{
							"id": "87d349ed-44d7-43e1-9a83-5f2406dee5bd",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20vdjEuMC91c2Vycz8kc2VsZWN0PWlkXHUwMDI2JHRvcD0yXHUwMDI2JHNraXB0b2tlbj1SRk53ZEFJQUFRQUFBQ002UVdSbGJHVldRRTB6TmpWNE1qRTBNelUxTG05dWJXbGpjbTl6YjJaMExtTnZiU2xWYzJWeVh6ZzNaRE0wT1dWa0xUUTBaRGN0TkRObE1TMDVZVGd6TFRWbU1qUXdObVJsWlRWaVpMa0FBQUFBQUFBQUFBQUEifQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/users?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA"),
			},
		},
		"invalid_request_standard_filter_applied_with_advanced_filter": {
			ctx: context.Background(),
			request: &framework.Request[azuread_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1.0",
					Filters: map[string]string{
						azuread_adapter.User: "displayName eq 'Test'",
					},
					AdvancedFilters: &azuread_adapter.AdvancedFilters{
						ScopedObjects: map[string][]azuread_adapter.EntityFilter{
							azuread_adapter.GroupMember: {
								{
									ScopeEntity:       azuread_adapter.Group,
									ScopeEntityFilter: "displayName eq 'Test'",
									Members: []azuread_adapter.MemberFilter{
										{
											MemberEntity:       azuread_adapter.User,
											MemberEntityFilter: "displayName eq 'Test'",
										},
									},
								},
							},
						},
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
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Implicit filters generated for entity `User` are not allowed with standard filters.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
			wantCursor: nil,
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
	adapter := azuread_adapter.NewAdapter(&azuread_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[azuread_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[azuread_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
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
							"id": "02bd9fd6-8f93-4758-87c3-1fb73740a315",
						},
						{
							"id": "06f62f70-9827-4e6e-93ef-8e0f2d9b7b23",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20vdjEuMC9ncm91cHM/JHNlbGVjdD1pZFx1MDAyNiR0b3A9Mlx1MDAyNiRza2lwdG9rZW49UkZOd2RBSUFBUUFBQUNwSGNtOTFjRjh3Tm1ZMk1tWTNNQzA1T0RJM0xUUmxObVV0T1RObFppMDRaVEJtTW1RNVlqZGlNak1xUjNKdmRYQmZNRFptTmpKbU56QXRPVGd5TnkwMFpUWmxMVGt6WldZdE9HVXdaakprT1dJM1lqSXpBQUFBQUFBQUFBQUFBQUEifQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/groups?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACpHcm91cF8wNmY2MmY3MC05ODI3LTRlNmUtOTNlZi04ZTBmMmQ5YjdiMjMqR3JvdXBfMDZmNjJmNzAtOTgyNy00ZTZlLTkzZWYtOGUwZjJkOWI3YjIzAAAAAAAAAAAAAAA"),
			},
		},
		"invalid_request_standard_filter_applied_with_advanced_filter": {
			ctx: context.Background(),
			request: &framework.Request[azuread_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1.0",
					Filters: map[string]string{
						azuread_adapter.Group: "displayName eq 'Test'",
					},
					AdvancedFilters: &azuread_adapter.AdvancedFilters{
						ScopedObjects: map[string][]azuread_adapter.EntityFilter{
							azuread_adapter.GroupMember: {
								{
									ScopeEntity:       azuread_adapter.Group,
									ScopeEntityFilter: "displayName eq 'Test'",
									Members: []azuread_adapter.MemberFilter{
										{
											MemberEntity:       azuread_adapter.User,
											MemberEntityFilter: "displayName eq 'Test'",
										},
									},
								},
							},
						},
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
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Implicit filters generated for entity `Group` are not allowed with standard filters.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
			wantCursor: nil,
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

func TestAdapterGetApplicationPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := azuread_adapter.NewAdapter(&azuread_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[azuread_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[azuread_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Application",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
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
							"id": "acc848e9-e8ec-4feb-a521-8d58b5482e09",
						},
						{
							"id": "cfa98ac0-a32c-4b4c-a78b-94c9912ed7b2",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20vdjEuMC9hcHBsaWNhdGlvbnM/JHNlbGVjdD1pZFx1MDAyNiR0b3A9Mlx1MDAyNiRza2lwdG9rZW49UkZOd2RBSUFBUUFBQURCQmNIQnNhV05oZEdsdmJsOWpabUU1T0dGak1DMWhNekpqTFRSaU5HTXRZVGM0WWkwNU5HTTVPVEV5WldRM1lqSXdRWEJ3YkdsallYUnBiMjVmWTJaaE9UaGhZekF0WVRNeVl5MDBZalJqTFdFM09HSXRPVFJqT1RreE1tVmtOMkl5QUFBQUFBQUFBQUFBQUFBIn0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/applications?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAADBBcHBsaWNhdGlvbl9jZmE5OGFjMC1hMzJjLTRiNGMtYTc4Yi05NGM5OTEyZWQ3YjIwQXBwbGljYXRpb25fY2ZhOThhYzAtYTMyYy00YjRjLWE3OGItOTRjOTkxMmVkN2IyAAAAAAAAAAAAAAA"),
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

func TestAdapterGetDevicePage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := azuread_adapter.NewAdapter(&azuread_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[azuread_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[azuread_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Device",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
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
							"id": "0357837b-ca6e-402d-9429-9e54dd51d97a",
						},
						{
							"id": "4d1ed9a4-519e-421b-b9f6-158991feff5b",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20vdjEuMC9kZXZpY2VzPyRzZWxlY3Q9aWRcdTAwMjYkdG9wPTJcdTAwMjYkc2tpcHRva2VuPVJGTndkQUlBQVFBQUFEQkJjSEJzYVdOaGRHbHZibDlqWm1FNU9HRmpNQzFoTXpKakxUUmlOR010WVRjNFlpMDVOR001T1RFeVpXUTNZakl3UVhCd2JHbGpZWFJwYjI1ZlkyWmhPVGhoWXpBdFlUTXlZeTAwWWpSakxXRTNPR0l0T1RSak9Ua3hNbVZrTjJJeUFBQUFBQUFBQUFBQUFBQSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/devices?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAADBBcHBsaWNhdGlvbl9jZmE5OGFjMC1hMzJjLTRiNGMtYTc4Yi05NGM5OTEyZWQ3YjIwQXBwbGljYXRpb25fY2ZhOThhYzAtYTMyYy00YjRjLWE3OGItOTRjOTkxMmVkN2IyAAAAAAAAAAAAAAA"),
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
	adapter := azuread_adapter.NewAdapter(&azuread_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[azuread_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[azuread_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
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
							"id": "6e7b768e-07e2-4810-8459-485f84f8f204-02bd9fd6-8f93-4758-87c3-1fb73740a315",
						},
						{
							"id": "87d349ed-44d7-43e1-9a83-5f2406dee5bd-02bd9fd6-8f93-4758-87c3-1fb73740a315",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20vdjEuMC9ncm91cHMvMDJiZDlmZDYtOGY5My00NzU4LTg3YzMtMWZiNzM3NDBhMzE1L21lbWJlcnM/JHNlbGVjdD1pZFx1MDAyNiR0b3A9Mlx1MDAyNiRza2lwdG9rZW49UkZOd2RBSUFBUUFBQUNNNlFXUmxiR1ZXUUUwek5qVjRNakUwTXpVMUxtOXViV2xqY205emIyWjBMbU52YlNsVmMyVnlYemczWkRNME9XVmtMVFEwWkRjdE5ETmxNUzA1WVRnekxUVm1NalF3Tm1SbFpUVmlaTGtBQUFBQUFBQUFBQUFBIiwiY29sbGVjdGlvbklkIjoiMDJiZDlmZDYtOGY5My00NzU4LTg3YzMtMWZiNzM3NDBhMzE1IiwiY29sbGVjdGlvbkN1cnNvciI6Imh0dHBzOi8vZ3JhcGgubWljcm9zb2Z0LmNvbS92MS4wL2dyb3Vwcz8kc2VsZWN0PWlkXHUwMDI2JHRvcD0xXHUwMDI2JHNraXB0b2tlbj1SRk53ZEFJQUFRQUFBQ3BIY205MWNGOHdObVkyTW1ZM01DMDVPREkzTFRSbE5tVXRPVE5sWmkwNFpUQm1NbVE1WWpkaU1qTXFSM0p2ZFhCZk1EWm1OakptTnpBdE9UZ3lOeTAwWlRabExUa3paV1l0T0dVd1pqSmtPV0kzWWpJekFBQUFBQUFBQUFBQUFBQSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor:           testutil.GenPtr("https://graph.microsoft.com/v1.0/groups/02bd9fd6-8f93-4758-87c3-1fb73740a315/members?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA"),
				CollectionID:     testutil.GenPtr("02bd9fd6-8f93-4758-87c3-1fb73740a315"),
				CollectionCursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/groups?$select=id&$top=1&$skiptoken=RFNwdAIAAQAAACpHcm91cF8wNmY2MmY3MC05ODI3LTRlNmUtOTNlZi04ZTBmMmQ5YjdiMjMqR3JvdXBfMDZmNjJmNzAtOTgyNy00ZTZlLTkzZWYtOGUwZjJkOWI3YjIzAAAAAAAAAAAAAAA"),
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

func TestAdapterGetRolePage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := azuread_adapter.NewAdapter(&azuread_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[azuread_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[azuread_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Role",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "description",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "roleTemplateId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "deletedDateTime",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":              "0fea7f0d-dea1-4028-8ce8-a686ec639d75",
							"deletedDateTime": time.Date(2024, 2, 2, 23, 21, 2, 0, time.UTC),
							"description":     "Can read basic directory information. Commonly used to grant directory read access to applications and guests.",
							"displayName":     "Directory Readers",
							"roleTemplateId":  "88d8e3e3-c189-46e8-94e1-9b9898b8876b",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiIxIn0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("1"),
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

			if diff := cmp.Diff(gotResponse.Success.Objects, tt.wantResponse.Success.Objects); diff != "" {
				t.Errorf("Differences found: (-got +want)\n%s", diff)
			}

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

func TestAdapterGetRoleMemberPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := azuread_adapter.NewAdapter(&azuread_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[azuread_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[azuread_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &azuread_adapter.Config{
					APIVersion: "v1.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "RoleMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "memberId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "roleId",
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
							"id":       "0fea7f0d-dea1-458d-9099-69fcc2e3cd42-65bb46a4-7d3j-9302-8a21-4d90f7a0efdb",
							"roleId":   "0fea7f0d-dea1-458d-9099-69fcc2e3cd42",
							"memberId": "65bb46a4-7d3j-9302-8a21-4d90f7a0efdb",
						},
						{
							"id":       "795326a8-6eef-410e-9604-649ca68e1241-65bb46a4-7d3j-9302-8a21-4d90f7a0efdb",
							"roleId":   "795326a8-6eef-410e-9604-649ca68e1241",
							"memberId": "65bb46a4-7d3j-9302-8a21-4d90f7a0efdb",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20vdjEuMC91c2Vycy82NWJiNDZhNC03ZDNqLTkzMDItOGEyMS00ZDkwZjdhMGVmZGIvdHJhbnNpdGl2ZU1lbWJlck9mL21pY3Jvc29mdC5ncmFwaC5kaXJlY3RvcnlSb2xlPyRzZWxlY3Q9aWRcdTAwMjYkdG9wPTJcdTAwMjYkc2tpcHRva2VuPU5FWFRMSU5LX1RPS0VOX1BMQUNFSE9MREVSXzQiLCJjb2xsZWN0aW9uSWQiOiI2NWJiNDZhNC03ZDNqLTkzMDItOGEyMS00ZDkwZjdhMGVmZGIiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiaHR0cHM6Ly9ncmFwaC5taWNyb3NvZnQuY29tL3YxLjAvdXNlcnM/JHNlbGVjdD1pZFx1MDAyNiR0b3A9MVx1MDAyNiRza2lwdG9rZW49TkVYVExJTktfVE9LRU5fUExBQ0VIT0xERVJfMSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor:           testutil.GenPtr("https://graph.microsoft.com/v1.0/users/65bb46a4-7d3j-9302-8a21-4d90f7a0efdb/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=2&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_4"),
				CollectionID:     testutil.GenPtr("65bb46a4-7d3j-9302-8a21-4d90f7a0efdb"),
				CollectionCursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/users?$select=id&$top=1&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_1"),
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
