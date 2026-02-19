// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package github_test

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
	github_adapter "github.com/sgnl-ai/adapters/pkg/github"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultOrganizationEntityConfig(),
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                   "MDEyOk9yZ2FuaXphdGlvbjk=",
							"enterpriseId":         "MDEwOkVudGVycHJpc2Ux",
							"databaseId":           int64(9),
							"login":                "ArvindOrg1",
							"viewerIsAMember":      true,
							"viewerCanCreateTeams": true,
							"updatedAt":            time.Date(2024, 2, 2, 23, 20, 22, 0, time.UTC),
							"createdAt":            time.Date(2024, 2, 2, 23, 20, 22, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpCWm5SbGNqRWlPaUpaTTFaNVl6STVlVTl1V1hsUGNFdHhVVmhLTW1GWE5XdFVNMHB1VFZGclBTSXNJa0ZtZEdWeU1pSTZiblZzYkN3aVFXWjBaWEl6SWpwdWRXeHNmUT09In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
				nil,
				nil,
			),
		},
		"valid_request_no_https_prefix": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultOrganizationEntityConfig(),
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                   "MDEyOk9yZ2FuaXphdGlvbjk=",
							"enterpriseId":         "MDEwOkVudGVycHJpc2Ux",
							"databaseId":           int64(9),
							"login":                "ArvindOrg1",
							"viewerIsAMember":      true,
							"viewerCanCreateTeams": true,
							"updatedAt":            time.Date(2024, 2, 2, 23, 20, 22, 0, time.UTC),
							"createdAt":            time.Date(2024, 2, 2, 23, 20, 22, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpCWm5SbGNqRWlPaUpaTTFaNVl6STVlVTl1V1hsUGNFdHhVVmhLTW1GWE5XdFVNMHB1VFZGclBTSXNJa0ZtZEdWeU1pSTZiblZzYkN3aVFXWjBaWEl6SWpwdWRXeHNmUT09In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
				nil,
				nil,
			),
		},
		"invalid_request_http_prefix": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: "http://" + strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
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
					Message: `Scheme "http" is not supported.`,
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"valid_request_with_cursor": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultOrganizationEntityConfig(),
				PageSize: 1,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                   "MDEyOk9yZ2FuaXphdGlvbjU=",
							"enterpriseId":         "MDEwOkVudGVycHJpc2Ux",
							"databaseId":           int64(5),
							"login":                "EnterpriseServerOrg",
							"viewerIsAMember":      true,
							"viewerCanCreateTeams": true,
							"updatedAt":            time.Date(2024, 1, 28, 23, 0, 0, 0, time.UTC),
							"createdAt":            time.Date(2024, 1, 28, 22, 59, 59, 0, time.UTC),
						},
					},
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetOrganizationPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultOrganizationEntityConfig(),
				Ordered:  false,
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                   "MDEyOk9yZ2FuaXphdGlvbjk=",
							"enterpriseId":         "MDEwOkVudGVycHJpc2Ux",
							"databaseId":           int64(9),
							"login":                "ArvindOrg1",
							"viewerIsAMember":      true,
							"viewerCanCreateTeams": true,
							"updatedAt":            time.Date(2024, 2, 2, 23, 20, 22, 0, time.UTC),
							"createdAt":            time.Date(2024, 2, 2, 23, 20, 22, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VV3M5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
				nil,
				nil,
			),
		},
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultOrganizationEntityConfig(),
				Ordered:  false,
				PageSize: 1,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                   "MDEyOk9yZ2FuaXphdGlvbjEw",
							"enterpriseId":         "MDEwOkVudGVycHJpc2Ux",
							"databaseId":           int64(10),
							"login":                "ArvindOrg2",
							"viewerIsAMember":      true,
							"viewerCanCreateTeams": true,
							"updatedAt":            time.Date(2024, 2, 15, 17, 0, 12, 0, time.UTC),
							"createdAt":            time.Date(2024, 2, 15, 17, 0, 12, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5aMjg5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
				nil,
				nil,
			),
		},
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultOrganizationEntityConfig(),
				Ordered:  false,
				PageSize: 1,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                   "MDEyOk9yZ2FuaXphdGlvbjU=",
							"enterpriseId":         "MDEwOkVudGVycHJpc2Ux",
							"databaseId":           int64(5),
							"login":                "EnterpriseServerOrg",
							"viewerIsAMember":      true,
							"viewerCanCreateTeams": true,
							"updatedAt":            time.Date(2024, 1, 28, 23, 0, 0, 0, time.UTC),
							"createdAt":            time.Date(2024, 1, 28, 22, 59, 59, 0, time.UTC),
						},
					},
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetOrganizationUserPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultOrganizationUserEntityConfig(),
				Ordered:  false,
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"uniqueId":  "MDEyOk9yZ2FuaXphdGlvbjU=-MDQ6VXNlcjQ=",
							"orgId":     "MDEyOk9yZ2FuaXphdGlvbjU=",
							"$.node.id": "MDQ6VXNlcjQ=",
							"role":      "ADMIN",
							"$.node.organizationVerifiedDomainEmails": []framework.Object{
								{
									"email": "arvind@sgnldemos.com",
								},
							},
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UlVVaUxDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbTUxYkd4OSIsImNvbGxlY3Rpb25JZCI6IkFydmluZE9yZzEiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKb1lYTk9aWGgwVUdGblpTSTZabUZzYzJVc0ltVnVaRU4xY25OdmNpSTZJbGt6Vm5sak1qbDVUMjVaZVU5d1MzRlJXRW95WVZjMWExUXpTbTVOVVdzOUlpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZRPT0ifQ==",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
				testutil.GenPtr("ArvindOrg1"),
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
			),
		},
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultOrganizationUserEntityConfig(),
				Ordered:  false,
				PageSize: 1,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
				testutil.GenPtr("ArvindOrg1"),
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"uniqueId":  "MDEyOk9yZ2FuaXphdGlvbjU=-MDQ6VXNlcjk=",
							"orgId":     "MDEyOk9yZ2FuaXphdGlvbjU=",
							"$.node.id": "MDQ6VXNlcjk=",
							"role":      "MEMBER",
							"$.node.organizationVerifiedDomainEmails": []framework.Object{
								{
									"email": "isabella@sgnldemos.com",
								},
							},
						},
					},
					NextCursor: "eyJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKb1lYTk9aWGgwVUdGblpTSTZabUZzYzJVc0ltVnVaRU4xY25OdmNpSTZJbGt6Vm5sak1qbDVUMjVaZVU5d1MzRlJXRW95WVZjMWExUXpTbTVOVVdzOUlpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZRPT0ifQ==",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				nil,
				nil,
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
			),
		},
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultOrganizationUserEntityConfig(),
				Ordered:  false,
				PageSize: 1,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				nil,
				nil,
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"uniqueId":  "MDEyOk9yZ2FuaXphdGlvbjEy-MDQ6VXNlcjQ=",
							"orgId":     "MDEyOk9yZ2FuaXphdGlvbjEy",
							"$.node.id": "MDQ6VXNlcjQ=",
							"role":      "ADMIN",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UlVVaUxDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbTUxYkd4OSIsImNvbGxlY3Rpb25JZCI6IkFydmluZE9yZzIiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKb1lYTk9aWGgwVUdGblpTSTZabUZzYzJVc0ltVnVaRU4xY25OdmNpSTZJbGt6Vm5sak1qbDVUMjVaZVU5d1MzRlJXRW95WVZjMWExUXpTbTVOWjI4OUlpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZRPT0ifQ==",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
				testutil.GenPtr("ArvindOrg2"),
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
			),
		},
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultOrganizationUserEntityConfig(),
				Ordered:  false,
				PageSize: 1,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
				testutil.GenPtr("ArvindOrg2"),
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "eyJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKb1lYTk9aWGgwVUdGblpTSTZabUZzYzJVc0ltVnVaRU4xY25OdmNpSTZJbGt6Vm5sak1qbDVUMjVaZVU5d1MzRlJXRW95WVZjMWExUXpTbTVOWjI4OUlpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZRPT0ifQ==",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				nil,
				nil,
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
			),
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetTeamPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultTeamEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                  "MDQ6VGVhbTI=",
							"enterpriseId":        "MDEwOkVudGVycHJpc2Ux",
							"orgId":               "MDEyOk9yZ2FuaXphdGlvbjk=",
							"databaseId":          int64(2),
							"slug":                "secret-team-1",
							"viewerCanAdminister": true,
							"updatedAt":           time.Date(2024, 2, 2, 23, 21, 54, 0, time.UTC),
							"createdAt":           time.Date(2024, 2, 2, 23, 21, 54, 0, time.UTC),
							"$.members.edges": []framework.Object{
								{
									"$.node.id":         "MDQ6VXNlcjY=",
									"role":              "MAINTAINER",
									"$.node.databaseId": int64(6),
									"$.node.email":      "",
									"$.node.login":      "arvind",
									"$.node.isViewer":   true,
									"$.node.updatedAt":  time.Date(2024, 1, 31, 5, 9, 26, 0, time.UTC),
									"$.node.createdAt":  time.Date(2024, 1, 28, 23, 28, 3, 0, time.UTC),
								},
							},
							"$.repositories.edges": []framework.Object{
								{
									"$.node.id":                "MDEwOlJlcG9zaXRvcnk2",
									"permission":               "ADMIN",
									"$.node.name":              "arvindrepo2",
									"$.node.databaseId":        int64(6),
									"$.node.url":               "https://ghe-test-server/ArvindOrg1/arvindrepo2",
									"$.node.allowUpdateBranch": false,
									"$.node.pushedAt":          time.Date(2024, 2, 2, 23, 22, 33, 0, time.UTC),
									"$.node.createdAt":         time.Date(2024, 2, 2, 23, 22, 32, 0, time.UTC),
								},
							},
						},
						{
							"id":                  "MDQ6VGVhbTE=",
							"enterpriseId":        "MDEwOkVudGVycHJpc2Ux",
							"orgId":               "MDEyOk9yZ2FuaXphdGlvbjk=",
							"databaseId":          int64(1),
							"slug":                "team1",
							"viewerCanAdminister": true,
							"updatedAt":           time.Date(2024, 2, 2, 23, 21, 2, 0, time.UTC),
							"createdAt":           time.Date(2024, 2, 2, 23, 21, 2, 0, time.UTC),
							"$.members.edges": []framework.Object{
								{
									"$.node.id":         "MDQ6VXNlcjQ=",
									"role":              "MEMBER",
									"$.node.databaseId": int64(4),
									"$.node.email":      "",
									"$.node.login":      "isabella",
									"$.node.isViewer":   false,
									"$.node.updatedAt":  time.Date(2024, 2, 22, 18, 43, 44, 0, time.UTC),
									"$.node.createdAt":  time.Date(2024, 1, 28, 22, 2, 26, 0, time.UTC),
								},
								{
									"$.node.id":         "MDQ6VXNlcjY=",
									"role":              "MAINTAINER",
									"$.node.databaseId": int64(6),
									"$.node.email":      "",
									"$.node.login":      "arvind",
									"$.node.isViewer":   true,
									"$.node.updatedAt":  time.Date(2024, 1, 31, 5, 9, 26, 0, time.UTC),
									"$.node.createdAt":  time.Date(2024, 1, 28, 23, 28, 3, 0, time.UTC),
								},
							},
							"$.repositories.edges": []framework.Object{
								{
									"$.node.id":                "MDEwOlJlcG9zaXRvcnk1",
									"permission":               "MAINTAIN",
									"$.node.name":              "arvindrepo1",
									"$.node.databaseId":        int64(5),
									"$.node.url":               "https://ghe-test-server/ArvindOrg1/arvindrepo1",
									"$.node.allowUpdateBranch": false,
									"$.node.pushedAt":          time.Date(2024, 2, 2, 23, 22, 20, 0, time.UTC),
									"$.node.createdAt":         time.Date(2024, 2, 2, 23, 22, 20, 0, time.UTC),
								},
								{
									"$.node.id":                "MDEwOlJlcG9zaXRvcnk2",
									"permission":               "WRITE",
									"$.node.name":              "arvindrepo2",
									"$.node.databaseId":        int64(6),
									"$.node.url":               "https://ghe-test-server/ArvindOrg1/arvindrepo2",
									"$.node.allowUpdateBranch": false,
									"$.node.pushedAt":          time.Date(2024, 2, 2, 23, 22, 33, 0, time.UTC),
									"$.node.createdAt":         time.Date(2024, 2, 2, 23, 22, 32, 0, time.UTC),
								},
								{
									"$.node.id":                "MDEwOlJlcG9zaXRvcnk3",
									"permission":               "READ",
									"$.node.name":              "arvindrepo3",
									"$.node.databaseId":        int64(7),
									"$.node.url":               "https://ghe-test-server/ArvindOrg1/arvindrepo3",
									"$.node.allowUpdateBranch": false,
									"$.node.pushedAt":          time.Date(2024, 2, 2, 23, 22, 45, 0, time.UTC),
									"$.node.createdAt":         time.Date(2024, 2, 2, 23, 22, 45, 0, time.UTC),
								},
							},
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VV3M5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
				nil,
				nil,
			),
		},
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultTeamEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5aMjg5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
				nil,
				nil,
			),
		},
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultTeamEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                  "MDQ6VGVhbTM=",
							"enterpriseId":        "MDEwOkVudGVycHJpc2Ux",
							"orgId":               "MDEyOk9yZ2FuaXphdGlvbjU=",
							"databaseId":          int64(3),
							"slug":                "random-team-1",
							"viewerCanAdminister": true,
							"updatedAt":           time.Date(2024, 2, 16, 4, 26, 6, 0, time.UTC),
							"createdAt":           time.Date(2024, 2, 16, 4, 26, 6, 0, time.UTC),
							"$.members.edges": []framework.Object{
								{
									"$.node.id":         "MDQ6VXNlcjY=",
									"role":              "MAINTAINER",
									"$.node.databaseId": int64(6),
									"$.node.email":      "",
									"$.node.login":      "arvind",
									"$.node.isViewer":   true,
									"$.node.updatedAt":  time.Date(2024, 1, 31, 5, 9, 26, 0, time.UTC),
									"$.node.createdAt":  time.Date(2024, 1, 28, 23, 28, 3, 0, time.UTC),
								},
							},
							"$.repositories.edges": []framework.Object{
								{
									"$.node.id":                "MDEwOlJlcG9zaXRvcnkx",
									"permission":               "MAINTAIN",
									"$.node.name":              "enterprise_repo1",
									"$.node.databaseId":        int64(1),
									"$.node.url":               "https://ghe-test-server/EnterpriseServerOrg/enterprise_repo1",
									"$.node.allowUpdateBranch": false,
									"$.node.pushedAt":          time.Date(2024, 2, 2, 23, 17, 27, 0, time.UTC),
									"$.node.createdAt":         time.Date(2024, 2, 2, 23, 17, 26, 0, time.UTC),
								},
							},
						},
					},
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetRepositoryPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultRepositoryEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                "MDEwOlJlcG9zaXRvcnk1",
							"enterpriseId":      "MDEwOkVudGVycHJpc2Ux",
							"orgId":             "MDEyOk9yZ2FuaXphdGlvbjk=",
							"name":              "arvindrepo1",
							"databaseId":        int64(5),
							"allowUpdateBranch": false,
							"pushedAt":          time.Date(2024, 2, 2, 23, 22, 20, 0, time.UTC),
							"createdAt":         time.Date(2024, 2, 2, 23, 22, 20, 0, time.UTC),
							"$.collaborators.edges": []framework.Object{
								{
									"$.node.id":  "MDQ6VXNlcjQ=",
									"permission": "ADMIN",
								},
								{
									"$.node.id":  "MDQ6VXNlcjY=",
									"permission": "MAINTAIN",
								},
							},
						},
						{
							"id":                "MDEwOlJlcG9zaXRvcnk2",
							"enterpriseId":      "MDEwOkVudGVycHJpc2Ux",
							"orgId":             "MDEyOk9yZ2FuaXphdGlvbjk=",
							"name":              "arvindrepo2",
							"databaseId":        int64(6),
							"allowUpdateBranch": false,
							"pushedAt":          time.Date(2024, 2, 2, 23, 22, 33, 0, time.UTC),
							"createdAt":         time.Date(2024, 2, 2, 23, 22, 32, 0, time.UTC),
							"$.collaborators.edges": []framework.Object{
								{
									"$.node.id":  "MDQ6VXNlcjQ=",
									"permission": "ADMIN",
								},
							},
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZSeUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEG")},
				nil,
				nil,
			),
		},
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultRepositoryEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEG")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                "MDEwOlJlcG9zaXRvcnk3",
							"enterpriseId":      "MDEwOkVudGVycHJpc2Ux",
							"orgId":             "MDEyOk9yZ2FuaXphdGlvbjk=",
							"name":              "arvindrepo3",
							"databaseId":        int64(7),
							"allowUpdateBranch": false,
							"pushedAt":          time.Date(2024, 2, 2, 23, 22, 45, 0, time.UTC),
							"createdAt":         time.Date(2024, 2, 2, 23, 22, 45, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VV3M5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
				nil,
				nil,
			),
		},
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultRepositoryEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5aMjg5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
				nil,
				nil,
			),
		},
		"fourth_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultRepositoryEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                "MDEwOlJlcG9zaXRvcnkx",
							"enterpriseId":      "MDEwOkVudGVycHJpc2Ux",
							"orgId":             "MDEyOk9yZ2FuaXphdGlvbjU=",
							"name":              "enterprise_repo1",
							"databaseId":        int64(1),
							"allowUpdateBranch": false,
							"pushedAt":          time.Date(2024, 2, 2, 23, 17, 27, 0, time.UTC),
							"createdAt":         time.Date(2024, 2, 2, 23, 17, 26, 0, time.UTC),
						},
						{
							"id":                "MDEwOlJlcG9zaXRvcnky",
							"enterpriseId":      "MDEwOkVudGVycHJpc2Ux",
							"orgId":             "MDEyOk9yZ2FuaXphdGlvbjU=",
							"name":              "enterprise_repo2",
							"databaseId":        int64(2),
							"allowUpdateBranch": false,
							"pushedAt":          time.Date(2024, 2, 2, 23, 17, 42, 0, time.UTC),
							"createdAt":         time.Date(2024, 2, 2, 23, 17, 41, 0, time.UTC),
							"$.collaborators.edges": []framework.Object{
								{
									"$.node.id":  "MDQ6VXNlcjY=",
									"permission": "MAINTAIN",
								},
							},
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5aMjg5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcDdJbWhoYzA1bGVIUlFZV2RsSWpwbVlXeHpaU3dpWlc1a1EzVnljMjl5SWpvaVdUTldlV015T1hsUGJsbDVUM0JGUXlJc0ltOXlaMkZ1YVhwaGRHbHZiazltWm5ObGRDSTZNQ3dpU1c1dVpYSlFZV2RsU1c1bWJ5STZiblZzYkgxOSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
		},
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultRepositoryEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                "MDEwOlJlcG9zaXRvcnkz",
							"enterpriseId":      "MDEwOkVudGVycHJpc2Ux",
							"orgId":             "MDEyOk9yZ2FuaXphdGlvbjU=",
							"name":              "enterprise_repo3",
							"databaseId":        int64(3),
							"allowUpdateBranch": false,
							"pushedAt":          time.Date(2024, 2, 2, 23, 18, 1, 0, time.UTC),
							"createdAt":         time.Date(2024, 2, 2, 23, 18, 1, 0, time.UTC),
						},
					},
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetRepositoryPageWithOrganizations(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
					Organizations:     []string{"arvindorg1", "arvindorg2"},
				},
				Entity:   *PopulateDefaultRepositoryEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                "MDEwOlJlcG9zaXRvcnk1",
							"name":              "arvindrepo1",
							"databaseId":        int64(5),
							"allowUpdateBranch": false,
							"orgId":             "O_kgDOCPwuWw",
							"pushedAt":          time.Date(2024, 2, 2, 23, 22, 20, 0, time.UTC),
							"createdAt":         time.Date(2024, 2, 2, 23, 22, 20, 0, time.UTC),
							"$.collaborators.edges": []framework.Object{
								{
									"$.node.id":  "MDQ6VXNlcjQ=",
									"permission": "ADMIN",
								},
								{
									"$.node.id":  "MDQ6VXNlcjY=",
									"permission": "MAINTAIN",
								},
							},
						},
						{
							"id":                "MDEwOlJlcG9zaXRvcnk2",
							"name":              "arvindrepo2",
							"databaseId":        int64(6),
							"allowUpdateBranch": false,
							// "pushedAt":          "2024-02-02T23:22:33Z",
							// "createdAt":         "2024-02-02T23:22:32Z",
							"orgId":     "O_kgDOCPwuWw",
							"pushedAt":  time.Date(2024, 2, 2, 23, 22, 33, 0, time.UTC),
							"createdAt": time.Date(2024, 2, 2, 23, 22, 32, 0, time.UTC),
							"$.collaborators.edges": []framework.Object{
								{
									"$.node.id":  "MDQ6VXNlcjQ=",
									"permission": "ADMIN",
								},
							},
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UlVjaUxDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbTUxYkd4OSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpEG")},
				nil,
				nil,
			),
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetUserPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultUserEntityConfig(),
				Ordered:  false,
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":         "MDQ6VXNlcjQ=",
							"databaseId": int64(4),
							"email":      "",
							"login":      "arooxa",
							"isViewer":   true,
							"updatedAt":  time.Date(2024, 3, 8, 4, 18, 47, 0, time.UTC),
							"createdAt":  time.Date(2024, 3, 8, 4, 18, 47, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZSU0lzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
				nil,
				nil,
			),
		},
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultUserEntityConfig(),
				Ordered:  false,
				PageSize: 1,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":         "MDQ6VXNlcjk=",
							"databaseId": int64(9),
							"email":      "",
							"login":      "isabella-sgnl",
							"isViewer":   false,
							"updatedAt":  time.Date(2024, 3, 8, 19, 28, 13, 0, time.UTC),
							"createdAt":  time.Date(2024, 3, 8, 17, 52, 21, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
		},
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultUserEntityConfig(),
				Ordered:  false,
				PageSize: 1,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":         "MDQ6VXNlcjQ=",
							"databaseId": int64(4),
							"email":      "",
							"login":      "arooxa",
							"isViewer":   true,
							"updatedAt":  time.Date(2024, 3, 8, 4, 18, 47, 0, time.UTC),
							"createdAt":  time.Date(2024, 3, 8, 4, 18, 47, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcDdJbWhoYzA1bGVIUlFZV2RsSWpwbVlXeHpaU3dpWlc1a1EzVnljMjl5SWpvaVdUTldlV015T1hsUGJsbDVUM0JGUlNJc0ltOXlaMkZ1YVhwaGRHbHZiazltWm5ObGRDSTZNQ3dpU1c1dVpYSlFZV2RsU1c1bWJ5STZiblZzYkgxOSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU="), testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
				nil,
				nil,
			),
		},
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultUserEntityConfig(),
				Ordered:  false,
				PageSize: 1,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU="), testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetCollaboratorPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultCollaboratorEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":         "MDQ6VXNlcjQ=",
							"databaseId": int64(4),
							"email":      "",
							"login":      "arooxa",
							"isViewer":   true,
							"updatedAt":  time.Date(2024, 3, 8, 4, 18, 47, 0, time.UTC),
							"createdAt":  time.Date(2024, 3, 8, 4, 18, 47, 0, time.UTC),
						},
						{
							"id":         "MDQ6VXNlcjk=",
							"databaseId": int64(9),
							"email":      "",
							"login":      "isabella-sgnl",
							"isViewer":   false,
							"updatedAt":  time.Date(2024, 3, 8, 19, 28, 13, 0, time.UTC),
							"createdAt":  time.Date(2024, 3, 8, 17, 52, 21, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
		},
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultCollaboratorEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":         "MDQ6VXNlcjQ=",
							"databaseId": int64(4),
							"email":      "",
							"login":      "arooxa",
							"isViewer":   true,
							"updatedAt":  time.Date(2024, 3, 8, 4, 18, 47, 0, time.UTC),
							"createdAt":  time.Date(2024, 3, 8, 4, 18, 47, 0, time.UTC),
						},
						{
							"id":         "MDQ6VXNlcjk=",
							"databaseId": int64(9),
							"email":      "",
							"login":      "isabella-sgnl",
							"isViewer":   false,
							"updatedAt":  time.Date(2024, 3, 8, 19, 28, 13, 0, time.UTC),
							"createdAt":  time.Date(2024, 3, 8, 17, 52, 21, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWxrelZubGpNamw1VDI1WmVVOXdSVW9pTENKdmNtZGhibWw2WVhScGIyNVBabVp6WlhRaU9qQXNJa2x1Ym1WeVVHRm5aVWx1Wm04aU9tNTFiR3g5ZlgwPSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEJ")},
				nil,
				nil,
			),
		},
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultCollaboratorEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEJ")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":         "MDQ6VXNlcjEw",
							"databaseId": int64(10),
							"email":      "",
							"login":      "r-rakshith",
							"isViewer":   false,
							"updatedAt":  time.Date(2024, 3, 8, 17, 53, 47, 0, time.UTC),
							"createdAt":  time.Date(2024, 3, 8, 17, 52, 54, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
		},
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultCollaboratorEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetLabelPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		// This fetches labels [1, 8] of repo 1/2 for org 1/2.
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultLabelEntityConfig(),
				Ordered:  false,
				PageSize: 8,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":           "MDU6TGFiZWwx",
							"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
							"name":         "bug",
							"color":        "d73a4a",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 30, 0, time.UTC),
							"isDefault":    true,
						},
						{
							"id":           "MDU6TGFiZWwy",
							"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
							"name":         "documentation",
							"color":        "0075ca",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 30, 0, time.UTC),
							"isDefault":    true,
						},
						{
							"id":           "MDU6TGFiZWwz",
							"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
							"name":         "duplicate",
							"color":        "cfd3d7",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 30, 0, time.UTC),
							"isDefault":    true,
						},
						{
							"id":           "MDU6TGFiZWw0",
							"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
							"name":         "enhancement",
							"color":        "a2eeef",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 30, 0, time.UTC),
							"isDefault":    true,
						},
						{
							"id":           "MDU6TGFiZWw1",
							"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
							"name":         "good first issue",
							"color":        "7057ff",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 30, 0, time.UTC),
							"isDefault":    true,
						},
						{
							"id":           "MDU6TGFiZWw2",
							"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
							"name":         "help wanted",
							"color":        "008672",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 30, 0, time.UTC),
							"isDefault":    true,
						},
						{
							"id":           "MDU6TGFiZWw3",
							"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
							"name":         "invalid",
							"color":        "e4e669",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 30, 0, time.UTC),
							"isDefault":    true,
						},
						{
							"id":           "MDU6TGFiZWw4",
							"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
							"name":         "question",
							"color":        "d876e3",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 30, 0, time.UTC),
							"isDefault":    true,
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSlBRU0lzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5ZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, nil, testutil.GenPtr("OA")},
				nil,
				nil,
			),
		},
		// This fetches label 9 of repo 1/2 for org 1/2.
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultLabelEntityConfig(),
				Ordered:  false,
				PageSize: 8,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, nil, testutil.GenPtr("OA")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":           "MDU6TGFiZWw5",
							"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
							"name":         "wontfix",
							"color":        "ffffff",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 30, 0, time.UTC),
							"isDefault":    true,
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
		},
		// This fetches labels [1, 8] of repo 2/2 for org 1/2.
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultLabelEntityConfig(),
				Ordered:  false,
				PageSize: 8,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":           "MDU6TGFiZWwxMA==",
							"repositoryId": "MDEwOlJlcG9zaXRvcnky",
							"name":         "bug",
							"color":        "d73a4a",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 44, 0, time.UTC),
							"isDefault":    true,
						},
						{
							"id":           "MDU6TGFiZWwxMQ==",
							"repositoryId": "MDEwOlJlcG9zaXRvcnky",
							"name":         "documentation",
							"color":        "0075ca",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 44, 0, time.UTC),
							"isDefault":    true,
						},
						{
							"id":           "MDU6TGFiZWwxMg==",
							"repositoryId": "MDEwOlJlcG9zaXRvcnky",
							"name":         "duplicate",
							"color":        "cfd3d7",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 44, 0, time.UTC),
							"isDefault":    true,
						},
						{
							"id":           "MDU6TGFiZWwxMw==",
							"repositoryId": "MDEwOlJlcG9zaXRvcnky",
							"name":         "enhancement",
							"color":        "a2eeef",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 44, 0, time.UTC),
							"isDefault":    true,
						},
						{
							"id":           "MDU6TGFiZWwxNA==",
							"repositoryId": "MDEwOlJlcG9zaXRvcnky",
							"name":         "good first issue",
							"color":        "7057ff",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 44, 0, time.UTC),
							"isDefault":    true,
						},
						{
							"id":           "MDU6TGFiZWwxNQ==",
							"repositoryId": "MDEwOlJlcG9zaXRvcnky",
							"name":         "help wanted",
							"color":        "008672",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 44, 0, time.UTC),
							"isDefault":    true,
						},
						{
							"id":           "MDU6TGFiZWwxNg==",
							"repositoryId": "MDEwOlJlcG9zaXRvcnky",
							"name":         "invalid",
							"color":        "e4e669",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 44, 0, time.UTC),
							"isDefault":    true,
						},
						{
							"id":           "MDU6TGFiZWwxNw==",
							"repositoryId": "MDEwOlJlcG9zaXRvcnky",
							"name":         "question",
							"color":        "d876e3",
							"createdAt":    time.Date(2024, 3, 8, 18, 51, 44, 0, time.UTC),
							"isDefault":    true,
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWs5Qklpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZYMTkifQ==",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("OA")},
				nil,
				nil,
			),
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetIssueLabelPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		// These is 1 issue on label 1/3 for repo 1/2 for org 1/2.
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueLabelEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":       "MDU6SXNzdWUz",
							"labelId":  "MDU6TGFiZWwx",
							"uniqueId": "MDU6TGFiZWwx-MDU6SXNzdWUz",
							"title":    "issue1",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSk5VU0lzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5ZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, nil, testutil.GenPtr("MQ")},
				nil,
				nil,
			),
		},
		// These are no issues on label 2/3 for repo 1/2 for org 1/2.
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueLabelEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, nil, testutil.GenPtr("MQ")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSk5aeUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5ZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, nil, testutil.GenPtr("Mg")},
				nil,
				nil,
			),
		},
		// There are 2 issues for label 3/3 on repo 1/2 for org 1/2.
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueLabelEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, nil, testutil.GenPtr("Mg")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":       "MDU6SXNzdWUz",
							"title":    "issue1",
							"labelId":  "MDU6TGFiZWwy",
							"uniqueId": "MDU6TGFiZWwy-MDU6SXNzdWUz",
						},
						{
							"id":       "MDU6SXNzdWU0",
							"title":    "issue2",
							"labelId":  "MDU6TGFiZWwy",
							"uniqueId": "MDU6TGFiZWwy-MDU6SXNzdWU0",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
		},
		// There are no labels in repo 2/2 for org 1/2, so there are no issue labels.
		"fourth_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueLabelEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
		},
		// There are no repositories in org 2, so there are no labels or issue labels.
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueLabelEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetPullRequestLabelPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		// PullRequestLabels Page 1: Org 1/2, Repo 1/2, Label [1]/2, PullRequest [1]/1
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestLabelEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":       "MDExOlB1bGxSZXF1ZXN0MQ==",
							"labelId":  "MDU6TGFiZWw0",
							"title":    "Create README.md",
							"uniqueId": "MDU6TGFiZWw0-MDExOlB1bGxSZXF1ZXN0MQ==",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSk9VU0lzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5ZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, nil, testutil.GenPtr("NQ")},
				nil,
				nil,
			),
		},
		// PullRequestLabels Page 2: Org 1/2, Repo 1/2, Label [2]/2, (has no pull requests)
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestLabelEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, nil, testutil.GenPtr("NQ")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
		},
		// PullRequestLabels Page 3: Org 1/2, Repo 2/2, Label [1]/1, PullRequest [1, 2]/2
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestLabelEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":       "MDExOlB1bGxSZXF1ZXsdsd$S0=",
							"title":    "BRANCH4PR",
							"labelId":  "MDU6TGFiZWw1",
							"uniqueId": "MDU6TGFiZWw1-MDExOlB1bGxSZXF1ZXsdsd$S0=",
						},
						{
							"id":       "MDExOlB1bGxSZXFsssd@@",
							"title":    "BRANCH5PR UPDATE README",
							"labelId":  "MDU6TGFiZWw1",
							"uniqueId": "MDU6TGFiZWw1-MDExOlB1bGxSZXFsssd@@",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
		},
		// PullRequestLabels Page 4: Org 2/2, (has no repos)
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestLabelEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetIssuePage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		// These are issues 2/2 for repo 1/2 for org 1/2.
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueEntityConfig(),
				Ordered:  false,
				PageSize: 8,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":             "MDU6SXNzdWUz",
							"title":          "issue1",
							"$.author.login": "arooxa",
							"repositoryId":   "MDEwOlJlcG9zaXRvcnkx",
							"createdAt":      time.Date(2024, 3, 15, 18, 40, 52, 0, time.UTC),
							"isPinned":       false,
						},
						{
							"id":             "MDU6SXNzdWU0",
							"title":          "issue2",
							"$.author.login": "arooxa",
							"repositoryId":   "MDEwOlJlcG9zaXRvcnkx",
							"createdAt":      time.Date(2024, 3, 15, 18, 41, 4, 0, time.UTC),
							"isPinned":       false,
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
		},
		// These are issues 2/2 for repo 2/2 for org 1/2.
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueEntityConfig(),
				Ordered:  false,
				PageSize: 8,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":             "MDU6SXNzdWUy",
							"title":          "issue3",
							"$.author.login": "arooxa",
							"repositoryId":   "MDEwOlJlcG9zaXRvcnky",
							"createdAt":      time.Date(2024, 3, 14, 17, 43, 3, 0, time.UTC),
							"isPinned":       false,
						},
						{
							"id":             "MDU6SXNzdWU1",
							"title":          "issue4",
							"$.author.login": "arooxa",
							"repositoryId":   "MDEwOlJlcG9zaXRvcnky",
							"createdAt":      time.Date(2024, 3, 15, 18, 42, 1, 0, time.UTC),
							"isPinned":       false,
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
		},
		// There are no repositories in org 2, so there are no issues.
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueEntityConfig(),
				Ordered:  false,
				PageSize: 8,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetIssueAssigneePage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		// These are the 2 assignees on issue 1/2 for repo 1/2 for org 1/2.
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueAssigneeEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":       "MDQ6VXNlcjQ=",
							"issueId":  "MDU6SXNzdWUz",
							"uniqueId": "MDU6SXNzdWUz-MDQ6VXNlcjQ=",
							"login":    "arooxa",
						},
						{
							"id":       "MDQ6VXNlcjk=",
							"issueId":  "MDU6SXNzdWUz",
							"uniqueId": "MDU6SXNzdWUz-MDQ6VXNlcjk=",
							"login":    "isabella-sgnl",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSlpNMVo1WXpJNWVVOXVXWGxQY0VWRUlpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZYMTkifQ==",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, nil, testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
		},
		// This is the 1 assignee on issue 2/2 for repo 1/2 for org 1/2.
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueAssigneeEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, nil, testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":       "MDQ6VXNlcjk=",
							"issueId":  "MDU6SXNzdWU0",
							"uniqueId": "MDU6SXNzdWU0-MDQ6VXNlcjk=",
							"login":    "isabella-sgnl",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
		},
		// These are the 2 assignees on issue 1/2 for repo 2/2 for org 1/2.
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueAssigneeEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":       "MDQ6VXNlcjQ=",
							"issueId":  "MDU6SXNzdWUy",
							"uniqueId": "MDU6SXNzdWUy-MDQ6VXNlcjQ=",
							"login":    "arooxa",
						},
						{
							"id":       "MDQ6VXNlcjk=",
							"issueId":  "MDU6SXNzdWUy",
							"uniqueId": "MDU6SXNzdWUy-MDQ6VXNlcjk=",
							"login":    "isabella-sgnl",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWxrelZubGpNamw1VDI1WmVVOXdSVU1pTENKdmNtZGhibWw2WVhScGIyNVBabVp6WlhRaU9qQXNJa2x1Ym1WeVVHRm5aVWx1Wm04aU9tNTFiR3g5ZlgwPSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
		},
		// These are the 2 assignees on issue 2/2 for repo 2/2 for org 1/2.
		"fourth_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueAssigneeEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":       "MDQ6VXNlcjQ=",
							"issueId":  "MDU6SXNzdWU1",
							"uniqueId": "MDU6SXNzdWU1-MDQ6VXNlcjQ=",
							"login":    "arooxa",
						},
						{
							"id":       "MDQ6VXNlcjEw",
							"issueId":  "MDU6SXNzdWU1",
							"uniqueId": "MDU6SXNzdWU1-MDQ6VXNlcjEw",
							"login":    "r-rakshith",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
		},
		// There are no repositories in org 2, so there are no issues or issue assignees.
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueAssigneeEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetIssueParticipantPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		// This is the only participant on issue 1/2 for repo 1/2 for org 1/2.
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueParticipantEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":       "MDQ6VXNlcjQ=",
							"issueId":  "MDU6SXNzdWUz",
							"uniqueId": "MDU6SXNzdWUz-MDQ6VXNlcjQ=",
							"login":    "arooxa",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSlpNMVo1WXpJNWVVOXVXWGxQY0VWRUlpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZYMTkifQ==",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, nil, testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
		},
		// This is the only participant on issue 2/2 for repo 1/2 for org 1/2.
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueParticipantEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, nil, testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":       "MDQ6VXNlcjQ=",
							"issueId":  "MDU6SXNzdWU0",
							"uniqueId": "MDU6SXNzdWU0-MDQ6VXNlcjQ=",
							"login":    "arooxa",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
		},
		// These are the 2 participants on issue 1/2 for repo 2/2 for org 1/2.
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueParticipantEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":       "MDQ6VXNlcjQ=",
							"issueId":  "MDU6SXNzdWUy",
							"uniqueId": "MDU6SXNzdWUy-MDQ6VXNlcjQ=",
							"login":    "arooxa",
						},
						{
							"id":       "MDQ6VXNlcjEw",
							"issueId":  "MDU6SXNzdWUy",
							"uniqueId": "MDU6SXNzdWUy-MDQ6VXNlcjEw",
							"login":    "r-rakshith",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWxrelZubGpNamw1VDI1WmVVOXdSVU1pTENKdmNtZGhibWw2WVhScGIyNVBabVp6WlhRaU9qQXNJa2x1Ym1WeVVHRm5aVWx1Wm04aU9tNTFiR3g5ZlgwPSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
		},
		// These are the 2 participants on issue 2/2 for repo 2/2 for org 1/2.
		"fourth_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueParticipantEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":       "MDQ6VXNlcjQ=",
							"issueId":  "MDU6SXNzdWU1",
							"uniqueId": "MDU6SXNzdWU1-MDQ6VXNlcjQ=",
							"login":    "arooxa",
						},
						{
							"id":       "MDQ6VXNlcjEw",
							"issueId":  "MDU6SXNzdWU1",
							"uniqueId": "MDU6SXNzdWU1-MDQ6VXNlcjEw",
							"login":    "r-rakshith",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
		},
		// There are no repositories in org 2, so there are no issues or issue participants.
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultIssueParticipantEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetPullRequestPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		// PullRequests Page 1: Org 1/2, Repo 1/2, PR 1/1
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                  "MDExOlB1bGxSZXF1ZXN0MQ==",
							"title":               "Create README.md",
							"closed":              false,
							"createdAt":           time.Date(2024, 3, 13, 23, 7, 49, 0, time.UTC),
							"$.author.login":      "arooxa",
							"$.baseRepository.id": "MDEwOlJlcG9zaXRvcnkx",
							"$.headRepository.id": "MDEwOlJlcG9zaXRvcnkx",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
		},
		// PullRequests Page 2: Org 1/2, Repo 2/2, PR [1, 2]/3
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                  "MDExOlB1bGxSZXF1ZXN0Mg==",
							"title":               "[branch4PR] README",
							"closed":              false,
							"createdAt":           time.Date(2024, 3, 15, 18, 43, 27, 0, time.UTC),
							"$.author.login":      "arooxa",
							"$.baseRepository.id": "MDEwOlJlcG9zaXRvcnky",
							"$.headRepository.id": "MDEwOlJlcG9zaXRvcnky",
						},
						{
							"id":                  "MDExOlB1bGxSZXF1ZXN0Mw==",
							"title":               "[branch5PR] README.md",
							"closed":              false,
							"createdAt":           time.Date(2024, 3, 15, 18, 46, 54, 0, time.UTC),
							"$.author.login":      "arooxa",
							"$.baseRepository.id": "MDEwOlJlcG9zaXRvcnky",
							"$.headRepository.id": "MDEwOlJlcG9zaXRvcnky",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWxrelZubGpNamw1VDI1WmVVOXdSVVFpTENKdmNtZGhibWw2WVhScGIyNVBabVp6WlhRaU9qQXNJa2x1Ym1WeVVHRm5aVWx1Wm04aU9tNTFiR3g5ZlgwPSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
		},
		// PullRequests Page 3: Org 1/2, Repo 2/2, PR [3]/3
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                  "MDExOlB1bGxSZXF1ZXN0NA==",
							"title":               "[branch6PR] readMe",
							"closed":              false,
							"createdAt":           time.Date(2024, 3, 15, 22, 40, 43, 0, time.UTC),
							"$.author.login":      "arooxa",
							"$.baseRepository.id": "MDEwOlJlcG9zaXRvcnky",
							"$.headRepository.id": "MDEwOlJlcG9zaXRvcnky",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
		},
		// PullRequests Page 4: Org 2/2 (has no repos)
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetPullRequestChangedFilePage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		// PullRequestChangedFiles Page 1: Org 1/2, Repo 1/2, PR 1/1, Files [1]/1
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestChangedFileEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"path":          "README.md",
							"changeType":    "ADDED",
							"pullRequestId": "MDExOlB1bGxSZXF1ZXN0MQ==",
							"uniqueId":      "MDExOlB1bGxSZXF1ZXN0MQ==-README.md",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
		},
		// PullRequestChangedFiles Page 2: Org 1/2, Repo 2/2, PR 1/3, Files [1]/1
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestChangedFileEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"path":          "random/file.txt",
							"changeType":    "DELETED",
							"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mg==",
							"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mg==-random/file.txt",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWxrelZubGpNamw1VDI1WmVVOXdSVU1pTENKdmNtZGhibWw2WVhScGIyNVBabVp6WlhRaU9qQXNJa2x1Ym1WeVVHRm5aVWx1Wm04aU9tNTFiR3g5ZlgwPSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
		},
		// PullRequestChangedFiles Page 3: Org 1/2, Repo 2/2, PR 2/3, Files [1]/1
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestChangedFileEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"path":          "README.md",
							"changeType":    "ADDED",
							"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mw==",
							"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mw==-README.md",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWxrelZubGpNamw1VDI1WmVVOXdSVVFpTENKdmNtZGhibWw2WVhScGIyNVBabVp6WlhRaU9qQXNJa2x1Ym1WeVVHRm5aVWx1Wm04aU9tNTFiR3g5ZlgwPSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
		},
		// PullRequestChangedFiles Page 4: Org 1/2, Repo 2/2, PR 3/3, Files [1]/1
		"fourth_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestChangedFileEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"path":          "random/file.txt",
							"changeType":    "ADDED",
							"pullRequestId": "MDExOlB1bGxSZXF1ZXN0NA==",
							"uniqueId":      "MDExOlB1bGxSZXF1ZXN0NA==-random/file.txt",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
		},
		// PullRequestChangedFiles Page 5: Org 2/2 (has no repos)
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestChangedFileEntityConfig(),
				Ordered:  false,
				PageSize: 2,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetPullRequestAssigneePage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		// PullRequestAssignee Page 1: Org 1/2, Repo 1/2, PullRequest 1/1, Assignees [1]/1
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestAssigneeEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":            "MDQ6VXNlcjQ=",
							"login":         "arooxa",
							"pullRequestId": "MDExOlB1bGxSZXF1ZXN0MQ==",
							"uniqueId":      "MDExOlB1bGxSZXF1ZXN0MQ==-MDQ6VXNlcjQ=",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
		},
		// PullRequestAssignee Page 2: Org 1/2, Repo 2/2, PullRequest 1/3, Assignees [1, 2]/2
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestAssigneeEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":            "MDQ6VXNlcjQ=",
							"login":         "arooxa",
							"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mg==",
							"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mg==-MDQ6VXNlcjQ=",
						},
						{
							"id":            "MDQ6VXNlcjk=",
							"login":         "isabella-sgnl",
							"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mg==",
							"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mg==-MDQ6VXNlcjk=",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWxrelZubGpNamw1VDI1WmVVOXdSVU1pTENKdmNtZGhibWw2WVhScGIyNVBabVp6WlhRaU9qQXNJa2x1Ym1WeVVHRm5aVWx1Wm04aU9tNTFiR3g5ZlgwPSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
		},
		// PullRequestAssignee Page 3: Org 1/2, Repo 2/2, PullRequest 2/3, Assignees [1, 1]/1
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestAssigneeEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":            "MDQ6VXNlcjQ=",
							"login":         "arooxa",
							"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mw==",
							"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mw==-MDQ6VXNlcjQ=",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWxrelZubGpNamw1VDI1WmVVOXdSVVFpTENKdmNtZGhibWw2WVhScGIyNVBabVp6WlhRaU9qQXNJa2x1Ym1WeVVHRm5aVWx1Wm04aU9tNTFiR3g5ZlgwPSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
		},
		// PullRequestAssignee Page 4: Org 1/2, Repo 2/2, PullRequest 3/3, (has no assignees)
		"fourth_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestAssigneeEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
		},
		// PullRequestAssignee Page 5: Org 2/2 (has no repos)
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestAssigneeEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetPullRequestParticipantPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		// PullRequestParticipant Page 1: Org 1/2, Repo 1/2, PullRequest 1/1, Participants [1]/1
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestParticipantEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":            "MDQ6VXNlcjQ=",
							"login":         "arooxa",
							"pullRequestId": "MDExOlB1bGxSZXF1ZXN0MQ==",
							"uniqueId":      "MDExOlB1bGxSZXF1ZXN0MQ==-MDQ6VXNlcjQ=",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
		},
		// PullRequestParticipant Page 2: Org 1/2, Repo 2/2, PullRequest 1/3, Participants [1, 2]/2
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestParticipantEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":            "MDQ6VXNlcjQ=",
							"login":         "arooxa",
							"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mg==",
							"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mg==-MDQ6VXNlcjQ=",
						},
						{
							"id":            "MDQ6VXNlcjEw",
							"login":         "r-rakshith",
							"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mg==",
							"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mg==-MDQ6VXNlcjEw",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWxrelZubGpNamw1VDI1WmVVOXdSVU1pTENKdmNtZGhibWw2WVhScGIyNVBabVp6WlhRaU9qQXNJa2x1Ym1WeVVHRm5aVWx1Wm04aU9tNTFiR3g5ZlgwPSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
		},
		// PullRequestParticipant Page 3: Org 1/2, Repo 2/2, PullRequest 2/3, Participants [1, 2]/2
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestParticipantEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":            "MDQ6VXNlcjQ=",
							"login":         "arooxa",
							"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mw==",
							"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mw==-MDQ6VXNlcjQ=",
						},
						{
							"id":            "MDQ6VXNlcjEw",
							"login":         "r-rakshith",
							"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mw==",
							"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mw==-MDQ6VXNlcjEw",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWxrelZubGpNamw1VDI1WmVVOXdSVVFpTENKdmNtZGhibWw2WVhScGIyNVBabVp6WlhRaU9qQXNJa2x1Ym1WeVVHRm5aVWx1Wm04aU9tNTFiR3g5ZlgwPSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
		},
		// PullRequestParticipant Page 4: Org 1/2, Repo 2/2, PullRequest 3/3, Participants [1]/1
		"fourth_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestParticipantEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":            "MDQ6VXNlcjQ=",
							"login":         "arooxa",
							"pullRequestId": "MDExOlB1bGxSZXF1ZXN0NA==",
							"uniqueId":      "MDExOlB1bGxSZXF1ZXN0NA==-MDQ6VXNlcjQ=",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
		},
		// PullRequestParticipant Page 5: Org 2/2 (has no repos)
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestParticipantEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetPullRequestCommitPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		// PullRequestCommit Page 1: Org 1/2, Repo 1/2, PullRequest 1/1, Commits [1]/1
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestCommitEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                         "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0MTo0YWNkMDEzNTJkNTZjYTMzMTA1ZmMyMjU4ZDFmMTI4NzZmMzhlZjRh",
							"pullRequestId":              "MDExOlB1bGxSZXF1ZXN0MQ==",
							"$.commit.id":                "MDY6Q29tbWl0MTo0YWNkMDEzNTJkNTZjYTMzMTA1ZmMyMjU4ZDFmMTI4NzZmMzhlZjRh",
							"$.commit.committedDate":     time.Date(2024, 3, 13, 23, 7, 39, 0, time.UTC),
							"$.commit.author.email":      "arvind@sgnl.ai",
							"$.commit.author.user.id":    "MDQ6VXNlcjQ=",
							"$.commit.author.user.login": "arooxa",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
		},
		// PullRequestCommit Page 2: Org 1/2, Repo 2/2, PullRequest 1/3, Commits [1, 3]/3
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestCommitEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                         "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0MjozZjBiMmRiMDM3NmJjYTgwNjM0NDRmNjI4ZWI3ZWI5Y2U4NTk1ZGNj",
							"pullRequestId":              "MDExOlB1bGxSZXF1ZXN0Mg==",
							"$.commit.id":                "MDY6Q29tbWl0MjozZjBiMmRiMDM3NmJjYTgwNjM0NDRmNjI4ZWI3ZWI5Y2U4NTk1ZGNj",
							"$.commit.committedDate":     time.Date(2024, 3, 15, 18, 43, 10, 0, time.UTC),
							"$.commit.author.email":      "arvind@sgnl.ai",
							"$.commit.author.user.id":    "MDQ6VXNlcjQ=",
							"$.commit.author.user.login": "arooxa",
						},
						{
							"id":                         "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0Mjo2MTFlOTU3NGUzODNiNWQ2NmVjNjAwNDMxYTg4ODRkMzc4OGJiMTQx",
							"pullRequestId":              "MDExOlB1bGxSZXF1ZXN0Mg==",
							"$.commit.id":                "MDY6Q29tbWl0Mjo2MTFlOTU3NGUzODNiNWQ2NmVjNjAwNDMxYTg4ODRkMzc4OGJiMTQx",
							"$.commit.committedDate":     time.Date(2024, 3, 16, 21, 18, 12, 0, time.UTC),
							"$.commit.author.email":      "arvind@sgnl.ai",
							"$.commit.author.user.id":    "MDQ6VXNlcjQ=",
							"$.commit.author.user.login": "arooxa",
						},
						{
							"id":                         "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0MjpjMWMzNmQ2ZWQ0M2U4ZmVmMjlhNGExNTc2ZWQxZTYxNGZkMGMzNDFi",
							"pullRequestId":              "MDExOlB1bGxSZXF1ZXN0Mg==",
							"$.commit.id":                "MDY6Q29tbWl0MjpjMWMzNmQ2ZWQ0M2U4ZmVmMjlhNGExNTc2ZWQxZTYxNGZkMGMzNDFi",
							"$.commit.committedDate":     time.Date(2024, 3, 22, 21, 48, 21, 0, time.UTC),
							"$.commit.author.email":      "rakshith@sgnl.ai",
							"$.commit.author.user.id":    "MDQ6VXNlcjEw",
							"$.commit.author.user.login": "r-rakshith",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWxrelZubGpNamw1VDI1WmVVOXdSVU1pTENKdmNtZGhibWw2WVhScGIyNVBabVp6WlhRaU9qQXNJa2x1Ym1WeVVHRm5aVWx1Wm04aU9tNTFiR3g5ZlgwPSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
		},
		// PullRequestCommit Page 3: Org 1/2, Repo 2/2, PullRequest 2/3, Commits [1]/1
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestCommitEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                         "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0Mzo1OTYxZGE3NDk1NmJhNWRiYTQ0YWEyYjQ4Mjc2MzM4MGNkNDhhMWZj",
							"pullRequestId":              "MDExOlB1bGxSZXF1ZXN0Mw==",
							"$.commit.id":                "MDY6Q29tbWl0Mjo1OTYxZGE3NDk1NmJhNWRiYTQ0YWEyYjQ4Mjc2MzM4MGNkNDhhMWZj",
							"$.commit.committedDate":     time.Date(2024, 3, 15, 18, 45, 3, 0, time.UTC),
							"$.commit.author.email":      "arvind@sgnl.ai",
							"$.commit.author.user.id":    "MDQ6VXNlcjQ=",
							"$.commit.author.user.login": "arooxa",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWxrelZubGpNamw1VDI1WmVVOXdSVVFpTENKdmNtZGhibWw2WVhScGIyNVBabVp6WlhRaU9qQXNJa2x1Ym1WeVVHRm5aVWx1Wm04aU9tNTFiR3g5ZlgwPSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
		},
		// PullRequestCommit Page 4: Org 1/2, Repo 2/2, PullRequest 3/3, Commits [1, 2]/2
		"fourth_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestCommitEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                         "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0NDo1YTVlNzJmNWQwZjk0MjVlZTk3NDc4NzMxZTc2MDczYjBmMTYzY2Fi",
							"pullRequestId":              "MDExOlB1bGxSZXF1ZXN0NA==",
							"$.commit.id":                "MDY6Q29tbWl0Mjo1YTVlNzJmNWQwZjk0MjVlZTk3NDc4NzMxZTc2MDczYjBmMTYzY2Fi",
							"$.commit.committedDate":     time.Date(2024, 3, 15, 22, 39, 33, 0, time.UTC),
							"$.commit.author.email":      "arvind@sgnl.ai",
							"$.commit.author.user.id":    "MDQ6VXNlcjQ=",
							"$.commit.author.user.login": "arooxa",
						},
						{
							"id":                         "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0NDpkNjE2NmYwYTlmMmQwMGZlYmFjYzZhYTM3MTAwYWY0YzAxNzBlYzhk",
							"pullRequestId":              "MDExOlB1bGxSZXF1ZXN0NA==",
							"$.commit.id":                "MDY6Q29tbWl0MjpkNjE2NmYwYTlmMmQwMGZlYmFjYzZhYTM3MTAwYWY0YzAxNzBlYzhk",
							"$.commit.committedDate":     time.Date(2024, 3, 15, 22, 44, 24, 0, time.UTC),
							"$.commit.author.email":      "arvind@sgnl.ai",
							"$.commit.author.user.id":    "MDQ6VXNlcjQ=",
							"$.commit.author.user.login": "arooxa",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
		},
		// PullRequestCommit Page 5: Org 2/2 (has no repos)
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestCommitEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetPullRequestReviewPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		// PullRequestReview Page 1: Org 1/2, Repo 1/2, PullRequest 1/1, (has no reviews)
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestReviewEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0=",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
		},
		// PullRequestReview Page 2: Org 1/2, Repo 2/2, PullRequest 1/3, Reviews [1]/1
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestReviewEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"$.author.login":            "r-rakshith",
							"state":                     "APPROVED",
							"id":                        "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3NQ==",
							"pullRequestId":             "MDExOlB1bGxSZXF1ZXN0Mg==",
							"createdAt":                 time.Date(2024, 3, 15, 21, 5, 52, 0, time.UTC),
							"authorCanPushToRepository": true,
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWxrelZubGpNamw1VDI1WmVVOXdSVU1pTENKdmNtZGhibWw2WVhScGIyNVBabVp6WlhRaU9qQXNJa2x1Ym1WeVVHRm5aVWx1Wm04aU9tNTFiR3g5ZlgwPSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
		},
		// PullRequestReview Page 3: Org 1/2, Repo 2/2, PullRequest 2/3, Reviews [1, 2]/2
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestReviewEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"$.author.login":            "r-rakshith",
							"state":                     "APPROVED",
							"pullRequestId":             "MDExOlB1bGxSZXF1ZXN0Mw==",
							"id":                        "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3Ng==",
							"createdAt":                 time.Date(2024, 3, 15, 21, 6, 25, 0, time.UTC),
							"authorCanPushToRepository": true,
						},
						{
							"$.author.login":            "isabella-sgnl",
							"state":                     "APPROVED",
							"pullRequestId":             "MDExOlB1bGxSZXF1ZXN0Mw==",
							"id":                        "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3Mg==",
							"createdAt":                 time.Date(2024, 3, 15, 19, 45, 9, 0, time.UTC),
							"authorCanPushToRepository": true,
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmV5Sm9ZWE5PWlhoMFVHRm5aU0k2Wm1Gc2MyVXNJbVZ1WkVOMWNuTnZjaUk2SWxrelZubGpNamw1VDI1WmVVOXdSVVFpTENKdmNtZGhibWw2WVhScGIyNVBabVp6WlhRaU9qQXNJa2x1Ym1WeVVHRm5aVWx1Wm04aU9tNTFiR3g5ZlgwPSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
		},
		// PullRequestReview Page 4: Org 1/2, Repo 2/2, PullRequest 3/3, Reviews [1]/1
		"fourth_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestReviewEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"$.author.login":            "isabella-sgnl",
							"state":                     "CHANGES_REQUESTED",
							"id":                        "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3OA==",
							"pullRequestId":             "MDExOlB1bGxSZXF1ZXN0NA==",
							"createdAt":                 time.Date(2024, 3, 15, 22, 46, 20, 0, time.UTC),
							"authorCanPushToRepository": true,
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9",
				},
			},
			wantCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
		},
		// PullRequestReview Page 5: Org 2/2 (has no repos)
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultPullRequestReviewEntityConfig(),
				Ordered:  false,
				PageSize: 5,
			},
			inputRequestCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
				nil,
				nil,
			),
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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

				if !ValidateGraphQLCompositeCursor(&gotCursor, tt.wantCursor) {
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetSecretScanningAlertPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := github_adapter.NewAdapter(&github_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[github_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultSecretScanningAlertEntityConfig(),
				Ordered:  false,
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"number":               int64(2),
							"state":                "resolved",
							"secret_type":          "adafruit_io_key",
							"$.repository.node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
							"created_at":           time.Date(2020, 11, 6, 18, 48, 51, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2UuY29tL2FwaS92My9lbnRlcnByaXNlcy9TR05ML3NlY3JldC1zY2FubmluZy9hbGVydHM/cGVyX3BhZ2U9MVx1MDAyNnBhZ2U9MiJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://test-instance.com/api/v3/enterprises/SGNL/secret-scanning/alerts?per_page=1&page=2"),
			},
		},
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[github_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &github_adapter.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Entity:   *PopulateDefaultSecretScanningAlertEntityConfig(),
				Ordered:  false,
				PageSize: 1,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr(server.URL + "/api/v3/enterprises/SGNL/secret-scanning/alerts?per_page=1&page=2"),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"number":               int64(1),
							"state":                "open",
							"secret_type":          "mailchimp_api_key",
							"$.repository.node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
							"created_at":           time.Date(2020, 11, 6, 18, 18, 30, 0, time.UTC),
						},
					},
					NextCursor: "",
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
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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
					t.Errorf("gotNextCursor: %v, wantNextCursor: %v", &gotCursor, tt.wantCursor)
				}
			}
		})
	}
}
