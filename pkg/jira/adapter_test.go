// Copyright 2025 SGNL.ai, Inc.
package jira_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/config"
	jira_adapter "github.com/sgnl-ai/adapters/pkg/jira"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := jira_adapter.NewAdapter(&jira_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx          context.Context
		request      *framework.Request[jira_adapter.Config]
		wantResponse framework.Response
		wantCursor   *pagination.CompositeCursor[int64]
	}{
		// More validation tests are done in validation_test.go.
		// The following two tests simply check if the adapter validates the incoming request.
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[jira_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Group,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "groupId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "createdAt",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
					},
				},
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"groupId":   "group1",
							"createdAt": time.Date(2023, 9, 29, 0, 0, 0, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOjF9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(1)),
			},
		},
		"invalid_request_missing_auth": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "example.com",
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Issue,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Jira auth is missing required basic credentials.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"issues_filter_applied": {
			ctx: context.Background(),
			request: &framework.Request[jira_adapter.Config]{
				Address: server.URL,
				Config: &jira_adapter.Config{
					IssuesJQLFilter: testutil.GenPtr("project='SGNL'"),
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Issue,
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
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id": "99",
						},
					},
				},
			},
			wantCursor: nil,
		},
		// This test ensures the filter should only be applied if the entity is an Issue.
		// There is no endpoint defined for Groups with a filter, so if this filter was being incorrectly
		// applied, the test would return a 404 endpoint not defined.
		"issues_filter_not_applied": {
			ctx: context.Background(),
			request: &framework.Request[jira_adapter.Config]{
				Address: server.URL,
				Config: &jira_adapter.Config{
					IssuesJQLFilter: testutil.GenPtr("project='SGNL'"),
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Group,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "groupId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"groupId": "group1",
						},
					},
					NextCursor: "eyJjdXJzb3IiOjF9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(1)),
			},
		},
		// If the query filter is applied correctly, we should only see 2 objects returned, since we
		// have defined only two in the mock server.
		"ql_query_filter_applied": {
			ctx: context.Background(),
			request: &framework.Request[jira_adapter.Config]{
				Address: server.URL,
				Config: &jira_adapter.Config{
					ObjectsQLQuery: testutil.GenPtr("objectType = Customer"),
					AssetBaseURL:   testutil.GenPtr(server.URL + "/assets"),
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(5),
					},
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Object,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "globalId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 10,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"globalId": "1",
						},
						{
							"globalId": "2",
						},
					},
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiIxIiwiY29sbGVjdGlvbkN1cnNvciI6MX0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				CollectionID:     testutil.GenPtr("1"),
				CollectionCursor: testutil.GenPtr[int64](1),
			},
		},
		// When the base URL is not specified, it should default to "https://api.atlassian.com/jsm/assets".
		// This should result in an error due to invalid certs.
		"default_asset_base_url_applied": {
			ctx: context.Background(),
			request: &framework.Request[jira_adapter.Config]{
				Address: server.URL,
				Config: &jira_adapter.Config{
					ObjectsQLQuery: testutil.GenPtr("objectType = Customer"),
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Object,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "globalId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 10,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: `Failed to execute Jira request: Post "https://api.atlassian.com/jsm/assets/` +
						`workspace/1/v1/object/aql?includeAttributes=true&startAt=0&maxResults=10": ` +
						`tls: failed to verify certificate: x509: certificate signed by unknown authority.`,
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
			wantCursor: nil,
		},
		"unable_to_decode_cursor": {
			ctx: context.Background(),
			request: &framework.Request[jira_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Group,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "groupId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 1,
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
			request: &framework.Request[jira_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Group,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "groupId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 1,
				// {"cursor":[]}
				Cursor: "eyJjdXJzb3IiOiBbXX0=",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to unmarshal JSON cursor: json: cannot unmarshal array into Go struct field " +
						"CompositeCursor[int64].cursor of type int64.",
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
		// "failed_to_parse_objects": {
		// 	ctx: context.Background(),
		// 	request: &framework.Request[jira_adapter.Config]{
		// 		Address: server.URL,
		// 		Auth: &framework.DatasourceAuthCredentials{
		// 			Basic: &framework.BasicAuthCredentials{
		// 				Username: mockUsername,
		// 				Password: mockPassword,
		// 			},
		// 		},
		// 		Entity: framework.EntityConfig{
		// 			ExternalId: jira_adapter.Group,
		// 			Attributes: []*framework.AttributeConfig{
		// 				{
		// 					ExternalId: "groupId",
		// 					Type:       framework.AttributeTypeDateTime,
		// 					List:       false,
		// 				},
		// 			},
		// 		},
		// 		PageSize: 1,
		// 		// {"cursor":102}
		// 		Cursor: "eyJjdXJzb3IiOjEwMn0=",
		// 	},
		// 	wantResponse: framework.Response{
		// 		Error: &framework.Error{
		// 			Message: "Failed to convert Jira response objects: attribute groupId cannot be parsed into a " +
		// 				"date-time value: failed to parse date-time value: 2005/07/06.",
		// 			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		// 		},
		// 	},
		// },
		// This test ensures that if the Jira SoR returns a non successful status code, we return an
		// appropriate error.
		"jira_request_returns_400": {
			ctx: context.Background(),
			request: &framework.Request[jira_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Group,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "groupId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 1,
				// {"cursor":103}
				Cursor: "eyJjdXJzb3IiOjEwM30=",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Datasource rejected request, returned status code: 400.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		// If a request to datasource.GetPage fails, we should return an appropriate error
		// wrapped in framework.NewGetPageResponseError.
		// In this case, we make a request to retrieve GroupMembers but the cursor does not contain
		// a group ID.
		"failed_to_make_get_page_request": {
			ctx: context.Background(),
			request: &framework.Request[jira_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.GroupMember,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 1,
				// {"cursor":1,"collectionCursor":1}
				Cursor: "eyJjdXJzb3IiOjEsImNvbGxlY3Rpb25DdXJzb3IiOjF9",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Cursor does not have CollectionID set for entity GroupMember.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
		// This test case uses a random URL instead of the test server's URL, so we should expect an invalid cert.
		// This test case also verifies that the "https://" prefix is added to the address if it's not present
		// which is evident by the error message.
		"failed_to_make_get_page_request_invalid_certs": {
			ctx: context.Background(),
			request: &framework.Request[jira_adapter.Config]{
				Address: "example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Group,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "groupId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: `Failed to execute Jira request: Get "https://example.com/rest/api/3/group/bulk` +
						`?startAt=0&maxResults=1": tls: failed to verify certificate: ` +
						`x509: certificate signed by unknown authority.`,
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"failed_to_make_get_page_request_invalid_host": {
			ctx: context.Background(),
			request: &framework.Request[jira_adapter.Config]{
				// Deliberately add the extra "/".
				// This test case also indirectly verifies that the "https://" prefix is NOT added
				// if the prefix already exists.
				Address: "https:///example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Group,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "groupId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: `Failed to execute Jira request: Get "https:///example.com/rest/api/3/group/bulk` +
						`?startAt=0&maxResults=1": http: no Host in request URL.`,
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
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
