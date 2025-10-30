// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst
package jira_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/jira"
	jira_adapter "github.com/sgnl-ai/adapters/pkg/jira"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

const (
	mockUsername = "username"
	mockPassword = "password"
)

type TestSuite struct {
	client jira_adapter.Client
	server *httptest.Server
}

// Define the endpoints and responses for the mock Jira server.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.URL.RequestURI() {
	// User endpoints
	case "/rest/api/3/users/search?startAt=0&maxResults=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"accountId": "1"}, {"accountId": "2"}]`))
	case "/rest/api/3/users/search?startAt=0&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"accountId": "1"}]`))
	case "/rest/api/3/users/search?startAt=1&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"accountId": "2"}]`))
	// For users and issues, we always query an extra page to check if there are
	// more results and determine if we've completed a sync.
	// So we define an extra page with no results to indicate a sync is complete.
	case "/rest/api/3/users/search?startAt=2&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))

	// Issue endpoints
	case "/rest/api/3/search?startAt=0&maxResults=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues": [{"id": "1"}, {"id": "2"}]}`))
	case "/rest/api/3/search?startAt=0&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues": [{"id": "1"}]}`))
	case "/rest/api/3/search?startAt=1&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues": [{"id": "2"}]}`))
	case "/rest/api/3/search?startAt=2&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues":[]}`))
	// With a JQL filter.
	case "/rest/api/3/search?jql=project%3D%27SGNL%27&startAt=0&maxResults=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues": [{"id": "99"}]}`))
	// With an invalid JQL filter (e.g. a project that doesn't exist).
	case "/rest/api/3/search?jql=project%3D%27INVALID%27&startAt=0&maxResults=10":
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"issues": []}`))

	// Group endpoints
	// Group endpoints have a convenient `isLast` field, compared to Users and Issues.
	case "/rest/api/3/group/bulk?startAt=0&maxResults=10":
		w.WriteHeader(http.StatusOK)
		// nolint: lll
		w.Write([]byte(`{"values": [{"groupId": "group1", "createdAt": "2023-09-29"}, {"groupId": "group2", "createdAt": "2023-09-29T11:17:42.000Z0700"}], "isLast": true}`))
	case "/rest/api/3/group/bulk?startAt=0&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"groupId": "group1", "createdAt": "2023-09-29"}], "isLast": false}`))
	case "/rest/api/3/group/bulk?startAt=1&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"groupId": "group2", "createdAt": "2023-09-29T11:17:42.000Z0700"}], "isLast": false}`))
	case "/rest/api/3/group/bulk?startAt=2&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"groupId": "group3"}], "isLast": true}`))

	// GroupMember endpoints
	// Group1 has 2 members.
	case "/rest/api/3/group/member?groupId=group1&startAt=0&maxResults=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"accountId": "member1"}, {"accountId": "member2"}], "isLast": true}`))
	case "/rest/api/3/group/member?groupId=group1&startAt=0&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"accountId": "member1"}], "isLast": false}`))
	case "/rest/api/3/group/member?groupId=group1&startAt=1&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"accountId": "member2"}], "isLast": true}`))
	// Group2 has 0 members.
	case "/rest/api/3/group/member?groupId=group2&startAt=0&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [], "isLast": true}`))
	// Group3 has 2 members.
	case "/rest/api/3/group/member?groupId=group3&startAt=0&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"accountId": "member3"}], "isLast": false}`))
	case "/rest/api/3/group/member?groupId=group3&startAt=1&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"accountId": "member4"}], "isLast": true}`))

	// Workspace endpoints
	case "/rest/servicedeskapi/assets/workspace?startAt=0&maxResults=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"workspaceId": "1"}, {"workspaceId": "2"}], "isLastPage": true}`))
	case "/rest/servicedeskapi/assets/workspace?startAt=0&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"workspaceId": "1"}], "isLastPage": false}`))
	case "/rest/servicedeskapi/assets/workspace?startAt=1&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"workspaceId": "2"}], "isLastPage": true}`))

	// Object endpoints
	case "/assets/workspace/1/v1/object/aql?includeAttributes=true&startAt=0&maxResults=10":
		body, _ := io.ReadAll(r.Body)
		switch string(body) {
		case `{"qlQuery":""}`:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"values": [{"globalId": "1"}, {"globalId": "2"}, {"globalId": "3"}], "isLast": true}`))
		case `{"qlQuery":"objectType = Customer"}`:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"values": [{"globalId": "1"}, {"globalId": "2"}], "isLast": true}`))
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"values": []}`))
		}

	// These endpoints define cases where tests should fail, e.g. missing fields, empty, etc.
	// Hence, they start from page 99 to avoid colliding with the above endpoints.
	// Return an empty list of groups.
	case "/rest/api/3/group/bulk?startAt=99&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [], "isLast": true}`))
	// Omit the Group's uniqueId.
	case "/rest/api/3/group/bulk?startAt=100&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"NOT_UNIQUE_ID":"group1"}], "isLast": true}`))
	// Make the Group's uniqueId not parsable as a string.
	case "/rest/api/3/group/bulk?startAt=101&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"groupId":10}], "isLast": true}`))
	// Return unparsable objects.
	case "/rest/api/3/group/bulk?startAt=102&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"groupId": "2005/07/06"}], "isLast": true}`))
	// Return 400 error.
	case "/rest/api/3/group/bulk?startAt=103&maxResults=1":
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"values": [{"groupId": "2005/07/06"}], "isLast": true}`))

	// Create a group member uniqueId that is not parsable into a string.
	case "/rest/api/3/group/member?groupId=group1&startAt=99&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"accountId": 4}], "isLast": true}`))
	// Create a group member response structure that is not expected.
	case "/rest/api/3/group/member?groupId=group1&startAt=100&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"UNEXPECTED_FIELD": [{"accountId": 4}], "isLast": true}`))

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})

func TestParseUsersResponse(t *testing.T) {
	tests := map[string]struct {
		body           []byte
		cursor         string
		pageSize       int64
		wantObjects    []map[string]interface{}
		wantNextCursor *string
		wantErr        *framework.Error
	}{
		"user_objects_last_page": {
			// Two user objects in response with page size = 10, so this must be last page.
			body:     []byte(`[{"name": "user1"}, {"name": "user2"}]`),
			cursor:   "0",
			pageSize: 10,
			wantObjects: []map[string]interface{}{
				{"name": "user1"},
				{"name": "user2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"user_objects_not_last_page": {
			// Two user objects in response with page size = 2, so there is a possibility of next page.
			body:     []byte(`[{"name": "user1"}, {"name": "user2"}]`),
			cursor:   "0",
			pageSize: 2,
			wantObjects: []map[string]interface{}{
				{"name": "user1"},
				{"name": "user2"},
			},
			wantNextCursor: testutil.GenPtr("2"), // This page contains index 0 and 1, so next page starts at index 2.
			wantErr:        nil,
		},
		"invalid_user_response": {
			// Users response should return a list of user objects, not a single top level object.
			body:           []byte(`{"users": [{"name": "user1"}, {"name": "user2"}]}`),
			cursor:         "0",
			pageSize:       2,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Failed to unmarshal Jira users response: json: cannot unmarshal object into " +
					"Go value of type []map[string]interface {}.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := jira_adapter.ParseUsersResponse(tt.body, tt.pageSize, tt.cursor)

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !reflect.DeepEqual(gotNextCursor, tt.wantNextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotNextCursor, tt.wantNextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestParseIssuesResponse(t *testing.T) {
	tests := map[string]struct {
		body           []byte
		cursor         string
		pageSize       int64
		wantObjects    []map[string]interface{}
		wantNextCursor *string
		wantErr        *framework.Error
	}{
		"issue_objects_last_page": {
			body:     []byte(`{"issues": [{"name": "issue1"}, {"name": "issue2"}]}`),
			cursor:   "0",
			pageSize: 10,
			wantObjects: []map[string]interface{}{
				{"name": "issue1"},
				{"name": "issue2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"issue_objects_not_last_page": {
			body:     []byte(`{"issues": [{"name": "issue1"}, {"name": "issue2"}]}`),
			cursor:   "0",
			pageSize: 2,
			wantObjects: []map[string]interface{}{
				{"name": "issue1"},
				{"name": "issue2"},
			},
			wantNextCursor: testutil.GenPtr("2"),
			wantErr:        nil,
		},
		"invalid_issue_response": {
			// Issues response should return a single top level object, not a list.
			body:           []byte(`[{"name": "issue1"}, {"name": "issue2"}]`),
			cursor:         "0",
			pageSize:       2,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Failed to unmarshal Jira issues response: json: cannot unmarshal array into " +
					"Go value of type map[string]interface {}.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"invalid_issue_object": {
			// The "issues" value should return []map[string]any, not []any.
			body:           []byte(`{"issues": ["issue1", {"name": "issue2"}]}`),
			cursor:         "0",
			pageSize:       2,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "An object in Entity: Issue could not be parsed. Expected: map[string]any. Got: string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"issues_field_does_not_exist": {
			body:           []byte(`{"WRONG_FIELD": [{"name": "issue1"}, {"name": "issue2"}]}`),
			cursor:         "0",
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Field missing in Jira issues response: issues.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"issues_field_exists_but_invalid_format": {
			// The "issues" field value should be a list of issue objects, not a map.
			body:           []byte(`{"issues": {"name": "issue1"}}`),
			cursor:         "0",
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Entity field exists in Jira issues response but field value is not a list of objects: " +
					"map[string]interface {}.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"fields_field_malformed": {
			// The "fields" field value should be an object.
			body:           []byte(`{"issues": [{"fields": "NOT_AN_OBJECT"}]}`),
			cursor:         "0",
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Failed to parse fields field in Jira Issue object as map[string]any.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		// Custom fields should be parsed from []objects into []any.
		"custom_field": {
			body: []byte(`{
				"issues": [
					{
						"fields": {
							"customfield_10069": [
								{
									"id": 1,
									"value": "A",
									"self": "1"
								},
								{
									"id": 2,
									"value": "B",
									"self": "2"
								}
							]
						}
					}
				]
			}`),
			cursor:   "0",
			pageSize: 10,
			wantObjects: []map[string]interface{}{
				{
					"fields": map[string]interface{}{
						"customfield_10069": []any{"A", "B"},
					},
				},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"null_custom_field": {
			body: []byte(`{
				"issues": [
					{
						"fields": {
							"customfield_10069": null
						}
					}
				]
			}`),
			cursor:   "0",
			pageSize: 10,
			wantObjects: []map[string]interface{}{
				{
					"fields": map[string]interface{}{
						"customfield_10069": nil,
					},
				},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		// Custom fields should have the "value" field.
		"custom_field_missing_value_field": {
			body: []byte(`{
						"issues": [
							{
								"fields": {
									"customfield_10069": [
										{
											"id": 1,
											"self": "1"
										},
										{
											"id": 2,
											"self": "2"
										}
									]
								}
							}
						]
					}`),
			cursor:         "0",
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Failed to parse value field in Jira customfield_10069 object as string: map[id:1 self:1].",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"custom_field_not_array_of_objects": {
			body: []byte(`{
						"issues": [
							{
								"fields": {
									"customfield_10069": "NOT_AN_ARRAY"
								}
							}
						]
					}`),
			cursor:         "0",
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Failed to parse customfield_10069 field in Jira Issue object as []any.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"custom_field_array_element_not_an_object": {
			body: []byte(`{
						"issues": [
							{
								"fields": {
									"customfield_10069": [
										null,
										{
											"id": 2,
											"self": "2"
										}
									]
								}
							}
						]
					}`),
			cursor:         "0",
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Failed to parse one of the objects in Jira customfield_10069 field into map[string]any: <nil>.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := jira_adapter.ParseIssuesResponse(tt.body, tt.pageSize, tt.cursor)

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !reflect.DeepEqual(gotNextCursor, tt.wantNextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotNextCursor, tt.wantNextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestParseGroupsResponse(t *testing.T) {
	tests := map[string]struct {
		body           []byte
		cursor         string
		pageSize       int64
		wantObjects    []map[string]interface{}
		wantNextCursor *string
		wantErr        *framework.Error
	}{
		"group_objects_last_page": {
			body:     []byte(`{"values": [{"name": "group1"}, {"name": "group2"}]}`),
			cursor:   "0",
			pageSize: 10,
			wantObjects: []map[string]interface{}{
				{"name": "group1"},
				{"name": "group2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"group_objects_not_last_page": {
			body:     []byte(`{"values": [{"name": "group1"}, {"name": "group2"}]}`),
			cursor:   "0",
			pageSize: 2,
			wantObjects: []map[string]interface{}{
				{"name": "group1"},
				{"name": "group2"},
			},
			wantNextCursor: testutil.GenPtr("2"),
			wantErr:        nil,
		},
		"invalid_group_response": {
			// Groups response should return a single top level object, not a list.
			body:           []byte(`[{"name": "group1"}, {"name": "group2"}]`),
			cursor:         "0",
			pageSize:       2,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Failed to unmarshal Jira Group response: json: cannot unmarshal array into " +
					"Go value of type map[string]interface {}.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"invalid_group_object": {
			// The "values" field should return []map[string]any, not []any.
			body:           []byte(`{"values": ["group1", {"name": "group2"}]}`),
			cursor:         "0",
			pageSize:       2,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "An object in Entity: Group could not be parsed. Expected: map[string]any. Got: string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"values_field_does_not_exist": {
			body:           []byte(`{"WRONG_FIELD": [{"name": "group1"}, {"name": "group2"}]}`),
			cursor:         "0",
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Field missing in Jira Group response: values.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"values_field_exists_but_invalid_format": {
			// The "values" field value should be a list of group objects, not a map.
			body:           []byte(`{"values": {"name": "group1"}}`),
			cursor:         "0",
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Entity field exists in Jira Group response but field value is not a list of objects: " +
					"map[string]interface {}.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"response_contains_valid_is_last_field": {
			// The next cursor does not need to be computed if `isLast` field is present.
			body:     []byte(`{"values": [{"name": "group1"}, {"name": "group2"}], "isLast": true}`),
			cursor:   "0",
			pageSize: 10,
			wantObjects: []map[string]interface{}{
				{"name": "group1"},
				{"name": "group2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"response_contains_invalid_is_last_field": {
			// The next cursor is computed if `isLast` field is present, but not a bool.
			body:           []byte(`{"values": [{"name": "group1"}, {"name": "group2"}], "isLast": "NOT A BOOL"}`),
			cursor:         "0",
			pageSize:       2,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Field isLast exists in Jira Group response but field value is not a bool: string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := jira_adapter.ParseGroupsResponse(tt.body, tt.pageSize, tt.cursor)

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !reflect.DeepEqual(gotNextCursor, tt.wantNextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotNextCursor, tt.wantNextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestParseGroupMembersResponse(t *testing.T) {
	tests := map[string]struct {
		body           []byte
		cursor         string
		pageSize       int64
		wantObjects    []map[string]interface{}
		wantNextCursor *string
		wantErr        *framework.Error
	}{
		"group_member_objects_last_page": {
			body:     []byte(`{"values": [{"name": "groupMember1"}, {"name": "groupMember2"}]}`),
			cursor:   "0",
			pageSize: 10,
			wantObjects: []map[string]interface{}{
				{"name": "groupMember1"},
				{"name": "groupMember2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"group_objects_not_last_page": {
			body:     []byte(`{"values": [{"name": "groupMember1"}, {"name": "groupMember2"}]}`),
			cursor:   "0",
			pageSize: 2,
			wantObjects: []map[string]interface{}{
				{"name": "groupMember1"},
				{"name": "groupMember2"},
			},
			wantNextCursor: testutil.GenPtr("2"),
			wantErr:        nil,
		},
		"invalid_group_members_response": {
			// Groups response should return a single top level object, not a list.
			body:           []byte(`[{"name": "groupMember1"}, {"name": "groupMember2"}]`),
			cursor:         "0",
			pageSize:       2,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Failed to unmarshal Jira GroupMember response: json: cannot unmarshal array into " +
					"Go value of type map[string]interface {}.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"invalid_group_member_object": {
			// The "values" field should return []map[string]any, not []any.
			body:           []byte(`{"values": ["groupMember1", {"name": "groupMember2"}]}`),
			cursor:         "0",
			pageSize:       2,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "An object in Entity: GroupMember could not be parsed. Expected: map[string]any. Got: string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"values_field_does_not_exist": {
			body:           []byte(`{"WRONG_FIELD": [{"name": "groupMember1"}, {"name": "groupMember2"}]}`),
			cursor:         "0",
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Field missing in Jira GroupMember response: values.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"values_field_exists_but_invalid_format": {
			// The "values" field value should be a list of group objects, not a map.
			body:           []byte(`{"values": {"name": "groupMember1"}}`),
			cursor:         "0",
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Entity field exists in Jira GroupMember response but field value is not " +
					"a list of objects: map[string]interface {}.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"response_contains_valid_is_last_field": {
			// The next cursor does not need to be computed if `isLast` field is present.
			body:     []byte(`{"values": [{"name": "groupMember1"}, {"name": "groupMember2"}], "isLast": true}`),
			cursor:   "0",
			pageSize: 10,
			wantObjects: []map[string]interface{}{
				{"name": "groupMember1"},
				{"name": "groupMember2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"response_contains_invalid_is_last_field": {
			// An error is thrown if we cannot parse the `isLast` field as a bool.
			body:           []byte(`{"values": [{"name": "group1"}, {"name": "group2"}], "isLast": "NOT A BOOL"}`),
			cursor:         "0",
			pageSize:       2,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Field isLast exists in Jira GroupMember response but field value is not a bool: string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := jira_adapter.ParseGroupMembersResponse(tt.body, tt.pageSize, tt.cursor)

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !reflect.DeepEqual(gotNextCursor, tt.wantNextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotNextCursor, tt.wantNextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestParseWorkspacesResponse(t *testing.T) {
	tests := map[string]struct {
		body           []byte
		cursor         string
		pageSize       int64
		wantObjects    []map[string]interface{}
		wantNextCursor *string
		wantErr        *framework.Error
	}{
		// The logic to parse Workspaces is shared by Groups and GroupMembers, which have already been tested above.
		"single_page": {
			body:     []byte(`{"values": [{"workspaceId": "1"}, {"workspaceId": "2"}]}`),
			cursor:   "0",
			pageSize: 10,
			wantObjects: []map[string]interface{}{
				{"workspaceId": "1"},
				{"workspaceId": "2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"first_page": {
			// The next cursor does not need to be computed if `isLastPage` field is present.
			body:     []byte(`{"values": [{"workspaceId": "1"}, {"workspaceId": "2"}], "isLastPage": false}`),
			cursor:   "0",
			pageSize: 2,
			wantObjects: []map[string]interface{}{
				{"workspaceId": "1"},
				{"workspaceId": "2"},
			},
			wantNextCursor: testutil.GenPtr("2"),
			wantErr:        nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := jira_adapter.ParseWorkspacesResponse(tt.body, tt.pageSize, tt.cursor)

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !reflect.DeepEqual(gotNextCursor, tt.wantNextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotNextCursor, tt.wantNextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestParseObjectsResponse(t *testing.T) {
	tests := map[string]struct {
		body           []byte
		cursor         string
		pageSize       int64
		wantObjects    []map[string]interface{}
		wantNextCursor *string
		wantErr        *framework.Error
	}{
		// The logic to parse Objects is shared by Groups and GroupMembers, which have already been tested above.
		"single_page": {
			body:     []byte(`{"values": [{"globalId": "1"}, {"globalId": "2"}], "isLast": true}`),
			cursor:   "0",
			pageSize: 10,
			wantObjects: []map[string]interface{}{
				{"globalId": "1"},
				{"globalId": "2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"first_page": {
			body:     []byte(`{"values": [{"globalId": "1"}, {"globalId": "2"}], "isLast": false}`),
			cursor:   "0",
			pageSize: 2,
			wantObjects: []map[string]interface{}{
				{"globalId": "1"},
				{"globalId": "2"},
			},
			wantNextCursor: testutil.GenPtr("2"),
			wantErr:        nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := jira_adapter.ParseObjectsResponse(tt.body, tt.pageSize, tt.cursor)

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !reflect.DeepEqual(gotNextCursor, tt.wantNextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotNextCursor, tt.wantNextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestConstructURL(t *testing.T) {
	tests := map[string]struct {
		request *jira_adapter.Request
		entity  jira.Entity
		cursor  *pagination.CompositeCursor[string]
		wantURL string
		wantErr error
	}{
		"users": {
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jira_adapter.User,
			},
			entity: jira.ValidEntityExternalIDs[jira_adapter.User],
			cursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("10"),
			},
			wantURL: "https://jira.com/rest/api/3/users/search?startAt=10&maxResults=10",
		},
		"groups": {
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jira_adapter.Group,
			},
			entity: jira.ValidEntityExternalIDs[jira_adapter.Group],
			cursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("10"),
			},
			wantURL: "https://jira.com/rest/api/3/group/bulk?startAt=10&maxResults=10",
		},
		"issues": {
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jira_adapter.Issue,
			},
			entity: jira.ValidEntityExternalIDs[jira_adapter.Issue],
			cursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("10"),
			},
			wantURL: "https://jira.com/rest/api/3/search?startAt=10&maxResults=10",
		},
		"enhanced_issues": {
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jira_adapter.EnhancedIssue,
			},
			entity: jira.ValidEntityExternalIDs[jira_adapter.EnhancedIssue],
			cursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("10"),
			},
			wantURL: "https://jira.com/rest/api/3/search/jql?nextPageToken=10&maxResults=10&fields=*navigable",
		},
		"issues_with_filter": {
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jira_adapter.Issue,
				IssuesJQLFilter:       testutil.GenPtr("project=TEST"),
			},
			entity: jira.ValidEntityExternalIDs[jira_adapter.Issue],
			cursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("10"),
			},
			wantURL: "https://jira.com/rest/api/3/search?jql=project%3DTEST&startAt=10&maxResults=10",
		},
		"group_members": {
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jira_adapter.GroupMember,
			},
			entity: jira.ValidEntityExternalIDs[jira_adapter.GroupMember],
			cursor: &pagination.CompositeCursor[string]{
				Cursor:           testutil.GenPtr("10"),
				CollectionID:     testutil.GenPtr("1"),
				CollectionCursor: testutil.GenPtr("1"),
			},
			wantURL: "https://jira.com/rest/api/3/group/member?groupId=1&startAt=10&maxResults=10",
		},
		"group_members_missing_group_id": {
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jira_adapter.GroupMember,
			},
			entity: jira.ValidEntityExternalIDs[jira_adapter.GroupMember],
			cursor: &pagination.CompositeCursor[string]{
				Cursor:           testutil.GenPtr("10"),
				CollectionCursor: testutil.GenPtr("1"),
			},
			wantErr: errors.New("cursor.CollectionID must not be nil for GroupMember entity"),
		},
		"workspaces": {
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jira_adapter.Workspace,
			},
			entity: jira.ValidEntityExternalIDs[jira_adapter.Workspace],
			cursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("10"),
			},
			wantURL: "https://jira.com/rest/servicedeskapi/assets/workspace?startAt=10&maxResults=10",
		},
		"objects": {
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jira_adapter.Object,
				AssetBaseURL:          testutil.GenPtr(jira.DefaultAssetBaseURL),
			},
			entity: jira.ValidEntityExternalIDs[jira_adapter.Object],
			cursor: &pagination.CompositeCursor[string]{
				Cursor:           testutil.GenPtr("10"),
				CollectionID:     testutil.GenPtr("1"),
				CollectionCursor: testutil.GenPtr("1"),
			},
			wantURL: `https://api.atlassian.com/jsm/assets/workspace/1/v1/object/aql?` +
				`includeAttributes=true&startAt=10&maxResults=10`,
		},
		"objects_missing_asset_base_url": {
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jira_adapter.Object,
			},
			entity: jira.ValidEntityExternalIDs[jira_adapter.Object],
			cursor: &pagination.CompositeCursor[string]{
				Cursor:           testutil.GenPtr("10"),
				CollectionID:     testutil.GenPtr("1"),
				CollectionCursor: testutil.GenPtr("1"),
			},
			wantErr: errors.New("request.AssetBaseURL must not be nil for Object entity"),
		},
		"objects_missing_workspace_id": {
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jira_adapter.Object,
			},
			entity: jira.ValidEntityExternalIDs[jira_adapter.Object],
			cursor: &pagination.CompositeCursor[string]{
				Cursor:           testutil.GenPtr("10"),
				CollectionCursor: testutil.GenPtr("1"),
			},
			wantErr: errors.New("cursor.CollectionID must not be nil for Object entity"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotURL, gotErr := jira_adapter.ConstructURL(tt.request, tt.entity, tt.cursor)

			if !reflect.DeepEqual(gotURL, tt.wantURL) {
				t.Errorf("gotURL: %v, wantURL: %v", gotURL, tt.wantURL)
			}

			if gotErr != nil && gotErr.Error() != tt.wantErr.Error() {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr.Error(), tt.wantErr.Error())
			}
		})
	}
}

func TestGetPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(30) * time.Second,
	}

	ts := &TestSuite{
		client: jira_adapter.NewClient(client),
		server: httptest.NewServer(TestServerHandler),
	}
	defer ts.server.Close()
	t.Run("TestGetPageErrors", ts.TestGetPageErrors)
	t.Run("TestGetPageUsers", ts.TestGetPageUsers)
	t.Run("TestGetPageIssues", ts.TestGetPageIssues)
	t.Run("TestGetPageGroups", ts.TestGetPageGroups)
	t.Run("TestGetPageGroupMembers", ts.TestGetPageGroupMembers)
	t.Run("TestGetPageWorkspaces", ts.TestGetPageWorkspaces)
	t.Run("TestGetPageObjects", ts.TestGetPageObjects)
}

func (ts *TestSuite) TestGetPageErrors(t *testing.T) {
	tests := map[string]struct {
		ctx          context.Context
		request      *jira_adapter.Request
		wantResponse *jira_adapter.Response
		wantErr      *framework.Error
	}{
		"invalid_url": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://{hello}",
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: `Address in datasource config is an invalid URL: parse ` +
					`"https://{hello}/rest/api/3/?startAt=0&maxResults=0": invalid character "{" in host name.`,
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "BAD_PROTOCOL",
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: `Failed to execute Jira request: Get "BAD_PROTOCOL/rest/api/3/?startAt=0&maxResults=0": ` +
					`unsupported protocol scheme "".`,
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResponse, gotErr := ts.client.GetPage(tt.ctx, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (ts *TestSuite) TestGetPageUsers(t *testing.T) {
	externalEntityID := jira_adapter.User

	tests := map[string]struct {
		ctx          context.Context
		request      *jira_adapter.Request
		wantResponse *jira_adapter.Response
		wantErr      *framework.Error
	}{
		"first_page_no_next_cursor": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"accountId": "1"},
					{"accountId": "2"},
				},
			},
			wantErr: nil,
		},
		"first_page_with_next_cursor": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"accountId": "1"},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("1"),
				},
			},
			wantErr: nil,
		},
		"second_last_page_last_user": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				// The last user occurs on page 2.
				Cursor:           &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("1")},
				PageSize:         int64(1),
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"accountId": "2"},
				},
				// We've synced the last user but it's not possible for our code to know that.
				// We still need to check if there are more results.
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("2"),
				},
			},
			wantErr: nil,
		},
		"last_page_no_users": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				Cursor:                &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("2")},
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects:    []map[string]any{},
				// An empty response means we've completed the sync.
				NextCursor: nil,
			},
			wantErr: nil,
		},
		"with_unnecessary_cursor_fields": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				// CollectionID and CollectionCursor fields should only ever be present if we're syncing GroupMembers.
				// This test ensures we throw an error if we do see these fields present and we're syncing a non GroupMember entity.
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:           testutil.GenPtr("2"),
					CollectionID:     testutil.GenPtr("1"),
					CollectionCursor: testutil.GenPtr("1"),
				},
				PageSize:         int64(1),
				EntityExternalID: externalEntityID,
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: "Cursor must not contain CollectionID or CollectionCursor fields for entity User.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResponse, gotErr := ts.client.GetPage(tt.ctx, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (ts *TestSuite) TestGetPageIssues(t *testing.T) {
	externalEntityID := jira_adapter.Issue

	tests := map[string]struct {
		ctx          context.Context
		request      *jira_adapter.Request
		wantResponse *jira_adapter.Response
		wantErr      *framework.Error
	}{
		"first_page_no_next_cursor": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "1"},
					{"id": "2"},
				},
			},
			wantErr: nil,
		},
		"first_page_with_next_cursor": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "1"},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("1"),
				},
			},
			wantErr: nil,
		},
		"second_last_page_last_issue": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				// The last issue occurs on page 2.
				Cursor:           &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("1")},
				PageSize:         int64(1),
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "2"},
				},
				// We've synced the last issue but it's not possible for our code to know that.
				// We still need to check if there are more results.
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("2"),
				},
			},
			wantErr: nil,
		},
		"last_page_empty": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				Cursor:                &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("2")},
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects:    []map[string]any{},
				// An empty response means we've completed the sync.
				NextCursor: nil,
			},
			wantErr: nil,
		},
		"with_jql_filter": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
				IssuesJQLFilter:       testutil.GenPtr("project='SGNL'"),
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "99"},
				},
			},
		},
		"with_invalid_jql_filter": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
				IssuesJQLFilter:       testutil.GenPtr("project='INVALID'"),
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 400,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResponse, gotErr := ts.client.GetPage(tt.ctx, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (ts *TestSuite) TestGetPageGroups(t *testing.T) {
	externalEntityID := jira_adapter.Group

	tests := map[string]struct {
		ctx          context.Context
		request      *jira_adapter.Request
		wantResponse *jira_adapter.Response
		wantErr      *framework.Error
		expectedLogs []map[string]any
	}{
		"first_page_no_next_cursor": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"groupId": "group1", "createdAt": "2023-09-29"},
					{"groupId": "group2", "createdAt": "2023-09-29T11:17:42.000Z0700"},
				},
			},
			wantErr: nil,
			expectedLogs: []map[string]any{
				{
					"level":                             "info",
					"msg":                               "Starting datasource request",
					fields.FieldRequestEntityExternalID: "Group",
					fields.FieldRequestPageSize:         int64(10),
				},
				{
					"level":                             "info",
					"msg":                               "Sending request to datasource",
					fields.FieldRequestEntityExternalID: "Group",
					fields.FieldRequestPageSize:         int64(10),
					fields.FieldRequestURL:              ts.server.URL + "/rest/api/3/group/bulk?startAt=0&maxResults=10",
				},
				{
					"level":                             "info",
					"msg":                               "Datasource request completed successfully",
					fields.FieldRequestEntityExternalID: "Group",
					fields.FieldRequestPageSize:         int64(10),
					fields.FieldResponseStatusCode:      int64(200),
					fields.FieldResponseObjectCount:     int64(2),
					fields.FieldResponseNextCursor:      (*pagination.CompositeCursor[string])(nil),
				},
			},
		},
		"first_page_with_next_cursor": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"groupId": "group1", "createdAt": "2023-09-29"},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("1"),
				},
			},
			wantErr: nil,
		},
		"middle_page_with_next_cursor": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				Cursor:                &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("1")},
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"groupId": "group2", "createdAt": "2023-09-29T11:17:42.000Z0700"},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("2"),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				// The last group occurs on page 3.
				Cursor:           &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("2")},
				PageSize:         int64(1),
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"groupId": "group3"},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctxWithLogger, observedLogs := testutil.NewContextWithObservableLogger(tt.ctx)

			gotResponse, gotErr := ts.client.GetPage(ctxWithLogger, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}

			testutil.ValidateLogOutput(t, observedLogs, tt.expectedLogs)
		})
	}
}

func (ts *TestSuite) TestGetPageGroupMembers(t *testing.T) {
	externalEntityID := jira_adapter.GroupMember

	tests := map[string]struct {
		ctx          context.Context
		request      *jira_adapter.Request
		wantResponse *jira_adapter.Response
		wantErr      *framework.Error
	}{
		"first_page_first_group": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "group1-member1", "accountId": "member1", "groupId": "group1"},
					{"id": "group1-member2", "accountId": "member2", "groupId": "group1"},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// Since the page size = 10, and group1 only has 2 members, we should not have a
					// next cursor. We can sync all group members for group1 in one page.
					Cursor: nil,
					// CollectionID is the ID of the group that we're currently/just finished syncing.
					CollectionID: testutil.GenPtr("group1"),
					// CollectionCursor is the cursor of the NEXT group that we're going to sync.
					CollectionCursor: testutil.GenPtr("1"),
				},
			},
			wantErr: nil,
		},
		"first_page_first_group_page_size_1": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "group1-member1", "accountId": "member1", "groupId": "group1"},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// Since the page size = 1, and group1 has 2 members, we should have a
					// next cursor.
					Cursor: testutil.GenPtr("1"),
					// CollectionID is the ID of the group that we're currently/just finished syncing.
					CollectionID: testutil.GenPtr("group1"),
					// CollectionCursor is the cursor of the NEXT group that we're going to sync.
					CollectionCursor: testutil.GenPtr("1"),
				},
			},
			wantErr: nil,
		},
		"second_page_first_group_page_size_1": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:           testutil.GenPtr("1"),
					CollectionID:     testutil.GenPtr("group1"),
					CollectionCursor: testutil.GenPtr("1"),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "group1-member2", "accountId": "member2", "groupId": "group1"},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// This is the last group member of this group. There should not be a next cursor.
					Cursor: nil,
					// CollectionID is the ID of the group that we're currently/just finished syncing.
					CollectionID: testutil.GenPtr("group1"),
					// CollectionCursor is the cursor of the NEXT group that we're going to sync.
					CollectionCursor: testutil.GenPtr("1"),
				},
			},
			wantErr: nil,
		},
		"empty_page_second_group_page_size_1": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID:     testutil.GenPtr("group1"),
					CollectionCursor: testutil.GenPtr("1"),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects:    []map[string]any{},
				NextCursor: &pagination.CompositeCursor[string]{
					// Group2 has no members. There should be no next cursor.
					Cursor: nil,
					// CollectionID is the ID of the group that we're currently/just finished syncing.
					CollectionID: testutil.GenPtr("group2"),
					// CollectionCursor is the cursor of the NEXT group that we're going to sync.
					CollectionCursor: testutil.GenPtr("2"),
				},
			},
			wantErr: nil,
		},
		"first_page_third_group_page_size_1": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID:     testutil.GenPtr("group2"),
					CollectionCursor: testutil.GenPtr("2"),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "group3-member3", "accountId": "member3", "groupId": "group3"},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// Group3 has two members. We should have a next cursor.
					Cursor: testutil.GenPtr("1"),
					// CollectionID is the ID of the group that we're currently/just finished syncing.
					CollectionID: testutil.GenPtr("group3"),
					// Group3 is the last group. There should be no next group cursor.
					CollectionCursor: nil,
				},
			},
			wantErr: nil,
		},
		"second_page_third_group_page_size_1": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:           testutil.GenPtr("1"),
					CollectionID:     testutil.GenPtr("group3"),
					CollectionCursor: nil,
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "group3-member4", "accountId": "member4", "groupId": "group3"},
				},
				// Group3 has two members. We've synced all of them. There should be no next cursor.
				// Group3 is also the last group, so CollectionCursor should be nil.
				// Both these conditions imply we have finished syncing all group members.
				// The entire NextCursor should be nil.
				NextCursor: nil,
			},
			wantErr: nil,
		},
		// If there are no groups to begin with, there is no need to sync any group members.
		"no_groups": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[string]{
					CollectionCursor: testutil.GenPtr("99"),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
			},
			wantErr: nil,
		},
		"group_unique_id_missing": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[string]{
					CollectionCursor: testutil.GenPtr("100"),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: "Jira Group object contains no groupId field.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
			},
		},
		"group_unique_id_not_string": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[string]{
					CollectionCursor: testutil.GenPtr("101"),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: "Failed to convert Jira Group object groupId field to string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
			},
		},
		"group_member_unique_id_not_string": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:           testutil.GenPtr("99"),
					CollectionID:     testutil.GenPtr("group1"),
					CollectionCursor: testutil.GenPtr("3"),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: "Failed to parse accountId field in Jira GroupMember response as string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		// If we're syncing group members, we must have a group id to sync.
		"composite_cursor_missing_group_id": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[string]{
					// Need to set Cursor here otherwise the code will think we're syncing the first page,
					// which doesn't require a group id.
					Cursor:           testutil.GenPtr("1"),
					CollectionCursor: testutil.GenPtr("1"),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: "Cursor does not have CollectionID set for entity GroupMember.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		// On the first page sync of a GroupMember, we make a request for the first group. If that request fails,
		// we should see an error.
		"group_get_page_fails_when_nil_cursor": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "http://localhost:1234",
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				Cursor:                nil,
				EntityExternalID:      externalEntityID,
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: `Failed to execute Jira request: Get "http://localhost:1234/rest/api/3/group/bulk` +
					`?startAt=0&maxResults=1": dial tcp [::1]:1234: connect: connection refused.`,
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		// If the GroupMember request is successful, but the response structure is not what we expect
		// (e.g. missing a field), we should see an error.
		"group_member_response_structure_is_invalid": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:           testutil.GenPtr("100"),
					CollectionID:     testutil.GenPtr("group1"),
					CollectionCursor: testutil.GenPtr("1"),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: "Field missing in Jira GroupMember response: values.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResponse, gotErr := ts.client.GetPage(tt.ctx, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (ts *TestSuite) TestGetPageWorkspaces(t *testing.T) {
	externalEntityID := jira_adapter.Workspace

	tests := map[string]struct {
		ctx          context.Context
		request      *jira_adapter.Request
		wantResponse *jira_adapter.Response
		wantErr      *framework.Error
	}{
		// The majority of this logic has already been tested in e.g. TestGetPageGroups, so
		// most duplicate test cases are omitted here.
		"first_page_no_next_cursor": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"workspaceId": "1"},
					{"workspaceId": "2"},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
		"first_page_with_next_cursor": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"workspaceId": "1"},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("1"),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				// The last workspace occurs on page 2.
				Cursor:           &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("1")},
				PageSize:         int64(1),
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"workspaceId": "2"},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResponse, gotErr := ts.client.GetPage(tt.ctx, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (ts *TestSuite) TestGetPageObjects(t *testing.T) {
	externalEntityID := jira_adapter.Object
	assetBaseURL := ts.server.URL + "/assets"

	tests := map[string]struct {
		ctx          context.Context
		request      *jira_adapter.Request
		wantResponse *jira_adapter.Response
		wantErr      *framework.Error
	}{
		"empty_ql_query": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
				ObjectsQLQuery:        testutil.GenPtr(""),
				AssetBaseURL:          &assetBaseURL,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"globalId": "1"},
					{"globalId": "2"},
					{"globalId": "3"},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// "1" is the collection (i.e. Workspace) ID that we just synced. Since the request
					// did not specify the cursor, this was the first workspace ID. Now that we've completed
					// the sync for all assets of this workspace, the next collection cursor should be 1 (up from 0).
					CollectionID:     testutil.GenPtr("1"),
					CollectionCursor: testutil.GenPtr("1"),
				},
			},
			wantErr: nil,
		},
		"non_empty_ql_query": {
			ctx: context.Background(),
			request: &jira_adapter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				Username:              mockUsername,
				Password:              mockPassword,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
				ObjectsQLQuery:        testutil.GenPtr("objectType = Customer"),
				AssetBaseURL:          &assetBaseURL,
			},
			wantResponse: &jira_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"globalId": "1"},
					{"globalId": "2"},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					CollectionID:     testutil.GenPtr("1"),
					CollectionCursor: testutil.GenPtr("1"),
				},
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResponse, gotErr := ts.client.GetPage(tt.ctx, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
