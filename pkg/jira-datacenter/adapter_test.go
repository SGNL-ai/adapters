// Copyright 2026 SGNL.ai, Inc.
package jiradatacenter_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	jiradatacenter_adapter "github.com/sgnl-ai/adapters/pkg/jira-datacenter"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := jiradatacenter_adapter.NewAdapter(&jiradatacenter_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx          context.Context
		request      *framework.Request[jiradatacenter_adapter.Config]
		wantResponse framework.Response
		wantCursor   *pagination.CompositeCursor[int64]
	}{
		// More validation tests are done in validation_test.go.
		// The following two tests simply check if the adapter validates the incoming request.
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.GroupExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
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
							"name": "group1",
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
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "example.com",
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Request to Jira is missing Basic Auth or Personal Access Token (PAT) credentials.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"issues_filter_applied": {
			ctx: context.Background(),
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: server.URL,
				Config: &jiradatacenter_adapter.Config{
					IssuesJQLFilter: testutil.GenPtr("project='SGNL'"),
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
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
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: server.URL,
				Config: &jiradatacenter_adapter.Config{
					IssuesJQLFilter: testutil.GenPtr("project='SGNL'"),
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.GroupExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
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
							"name": "group1",
						},
					},
					NextCursor: "eyJjdXJzb3IiOjF9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(1)),
			},
		},
		"unable_to_decode_cursor": {
			ctx: context.Background(),
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.GroupExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
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
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.GroupExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
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
		"failed_to_parse_objects": {
			ctx: context.Background(),
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Config: &jiradatacenter_adapter.Config{
					IssuesJQLFilter: testutil.GenPtr("project='BAD_DATE_FORMAT'"),
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
					},
				},
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to convert Jira response objects: attribute id cannot be parsed into a " +
						"date-time value: failed to parse date-time value: 2005/07/06.",
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		// This test ensures that if the Jira SoR returns a non successful status code, we return an
		// appropriate error.
		"jira_request_returns_400": {
			ctx: context.Background(),
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 1,
				Config: &jiradatacenter_adapter.Config{
					IssuesJQLFilter: testutil.GenPtr("project='NONEXISTENT_PROJECT'"),
				},
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
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.GroupMemberExternalID,
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
		// This test case uses a non-existent local address to ensure connection failure.
		// This test case also verifies that the "https://" prefix is added to the address if it's not present
		// which is evident by the error message.
		"failed_to_make_get_page_request_connection_refused": {
			ctx: context.Background(),
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "localhost:1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.GroupExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: `Failed to execute Jira request: Get "https://localhost:1/rest/api/latest/groups/picker": ` +
						`dial tcp [::1]:1: connect: connection refused.`,
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"failed_to_make_get_page_request_invalid_host": {
			ctx: context.Background(),
			request: &framework.Request[jiradatacenter_adapter.Config]{
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
					ExternalId: jiradatacenter_adapter.GroupExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: `Failed to execute Jira request: Get "https:///example.com/rest/api/latest/groups/picker": ` +
						`http: no Host in request URL.`,
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"issues_with_child_entities": {
			ctx: context.Background(),
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: server.URL,
				Config: &jiradatacenter_adapter.Config{
					IssuesJQLFilter: testutil.GenPtr("project='CHILD_ENTITIES_PRESENT'"),
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: mockUsername,
						Password: mockPassword,
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "key",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.fields.summary[0].description",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "$.fields.issuetype",
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
								{
									ExternalId: "description",
									Type:       framework.AttributeTypeString,
									List:       false,
								},
							},
						},
						{
							ExternalId: "$.fields.assignee",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "accountId",
									Type:       framework.AttributeTypeString,
									List:       false,
								},
								{
									ExternalId: "displayName",
									Type:       framework.AttributeTypeString,
									List:       false,
								},
								{
									ExternalId: "emailAddress",
									Type:       framework.AttributeTypeString,
									List:       false,
								},
							},
						},
						{
							ExternalId: "$.fields.non_existent",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "accountId",
									Type:       framework.AttributeTypeString,
									List:       false,
								},
							},
						},
					},
				},
				PageSize: 10,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                              "ISSUE-100",
							"key":                             "TEST-100",
							"$.fields.summary[0].description": "Issue with child entities",
							"$.fields.assignee": []framework.Object{
								{
									"accountId":    "user123",
									"displayName":  "John Doe",
									"emailAddress": "john.doe@example.com",
								},
							},
							"$.fields.issuetype": []framework.Object{
								{
									"description": "A bug issue type",
									"id":          "1",
									"name":        "Bug",
								},
							},
						},
					},
					NextCursor: "",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				// Print detailed error information
				if gotResponse.Error != nil && tt.wantResponse.Error != nil {
					t.Logf("Got error message: %q", gotResponse.Error.Message)
					t.Logf("Want error message: %q", tt.wantResponse.Error.Message)
					t.Logf("Got error code: %v", gotResponse.Error.Code)
					t.Logf("Want error code: %v", tt.wantResponse.Error.Code)
				}

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
