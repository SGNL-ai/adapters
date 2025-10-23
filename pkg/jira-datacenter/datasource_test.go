// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst
package jiradatacenter_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	jiradatacenter "github.com/sgnl-ai/adapters/pkg/jira-datacenter"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

const (
	mockUsername            = "username"
	mockPassword            = "password"
	mockAuthorizationHeader = "authHeader"
)

type TestSuite struct {
	client jiradatacenter.Client
	server *httptest.Server
}

// Define the endpoints and responses for the mock Jira server.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.URL.RequestURI() {
	// Issue endpoints
	case "/rest/api/latest/search?fields=*navigable&startAt=0&maxResults=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues": [{"id": "1"}, {"id": "2"}]}`))
	case "/rest/api/latest/search?fields=*navigable&startAt=0&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues": [{"id": "1"}]}`))
	case "/rest/api/latest/search?fields=*navigable&startAt=1&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues": [{"id": "2"}]}`))
	case "/rest/api/latest/search?fields=*navigable&startAt=2&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues":[]}`))
	// With a JQL filter.
	case "/rest/api/latest/search?jql=project%3D%27SGNL%27&fields=*navigable&startAt=0&maxResults=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues": [{"id": "99"}]}`))
	// With a JQL filter and specific fields
	case "/rest/api/latest/search?jql=project%3D%27SGNL%27&fields=id&startAt=0&maxResults=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues": [{"id": "99"}]}`))
	// With child entities present in response (detailed field paths)
	case "/rest/api/latest/search?jql=project%3D%27CHILD_ENTITIES_PRESENT%27&fields=" +
		"assignee.accountId%2Cassignee.displayName%2Cassignee.emailAddress%2Cid%2C" +
		"issuetype.description%2Cissuetype.id%2Cissuetype.name%2Ckey%2Csummary&startAt=0&maxResults=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"issues": [
				{
					"id": "ISSUE-100",
					"key": "TEST-100",
					"fields": {
						"summary": [
							{
								"id": "100",
								"description": "Issue with child entities"
							}
						],
						"issuetype": [
							{
								"id": "1",
								"name": "Bug",
								"description": "A bug issue type"
							}
						],
						"assignee": [
							{
								"accountId": "user123",
								"displayName": "John Doe",
								"emailAddress": "john.doe@example.com"
							}
						]
					}
				}
			]
		}`))
	// With child entities present in response (simple field paths)
	case "/rest/api/latest/search?jql=project%3D%27CHILD_ENTITIES_PRESENT%27&fields=" +
		"assignee%2Cid%2Cissuetype%2Ckey%2Cnon_existent%2Csummary&startAt=0&maxResults=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"startAt": 0,
			"maxResults": 10,
			"total": 1,
			"issues": [
				{
					"id": "ISSUE-100",
					"key": "TEST-100",
					"fields": {
						"summary": [
							{
								"id": "100",
								"description": "Issue with child entities"
							},
							{
								"id": "200",
								"description": "Invalid description"
							}
						],
						"issuetype": [
							{
								"id": "1",
								"name": "Bug",
								"description": "A bug issue type"
							}
						],
						"assignee": [
							{
								"accountId": "user123",
								"displayName": "John Doe",
								"emailAddress": "john.doe@example.com"
							}
						]
					}
				}
			]
		}`))
	// With child entities missing from response
	case "/rest/api/latest/search?jql=project%3D%27CHILD_ENTITIES_MISSING%27&fields=" +
		"assignee.accountId%2Cassignee.displayName%2Cid%2Cissuetype.id%2Cissuetype.name%2C" +
		"key%2Csummary&startAt=0&maxResults=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"issues": [
				{
					"id": "ISSUE-101",
					"key": "TEST-101",
					"fields": {
						"summary": [
							{
								"id": "100",
								"description": "Issue with child entities"
							}
						]
					}
				}
			]
		}`))
	// With a JQL filter and child entity fields (legacy test)
	case "/rest/api/latest/search?jql=project%3D%27SGNL%27&fields=" +
		"id%2Cissuetype.%24.fields.assignee%2Cissuetype.%24.fields.customfield_10209%2C" +
		"issuetype.id%2Cissuetype.name%2Ckey%2Csummary%2Cwatchers.%24.values&startAt=0&maxResults=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues": [{"id": "99"}]}`))
	// With an invalid JQL filter (e.g. a project that doesn't exist).
	case "/rest/api/latest/search?jql=project%3D%27INVALID%27&fields=*navigable&startAt=0&maxResults=10":
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"issues": []}`))

	// Simulates a 400 Bad Request response when trying to query a nonexistent project
	case "/rest/api/latest/search?jql=project%3D%27NONEXISTENT_PROJECT%27&fields=*navigable&startAt=0&maxResults=1":
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errorMessages":["The project 'NONEXISTENT_PROJECT' does not exist."]}`))
	// With specific field (for jira_request_returns_400 test)
	case "/rest/api/latest/search?jql=project%3D%27NONEXISTENT_PROJECT%27&fields=id&startAt=0&maxResults=1":
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errorMessages":["The project 'NONEXISTENT_PROJECT' does not exist."]}`))

	// Simulates a response with a date in an unparseable format (2005/07/06)
	case "/rest/api/latest/search?jql=project%3D%27BAD_DATE_FORMAT%27&fields=*navigable&startAt=0&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues": [{"id": "2005/07/06"}]}`))
	// With a JQL filter and specific field (for failed_to_parse_objects test)
	case "/rest/api/latest/search?jql=project%3D%27BAD_DATE_FORMAT%27&fields=id&startAt=0&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues": [{"id": "2005/07/06"}]}`))

	// Test with specific attributes
	case "/rest/api/latest/search?fields=id%2Ckey%2Cstatus%2Csummary&startAt=10&maxResults=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issues": [{"id": "1", "key": "TEST-1", "summary": "Test Issue", "status": "Open"}]}`))

	// Group endpoints
	case "/rest/api/latest/groups/picker":
		w.WriteHeader(http.StatusOK)
		// nolint: lll
		w.Write([]byte(`{"groups": [{"name": "group1"}, {"name": "group2"}, {"name": "group3"}]}`))

	// GroupMember endpoints
	// Group1 has 2 members.
	case "/rest/api/latest/group/member?groupname=group1&startAt=0&maxResults=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"key": "member1"}, {"key": "member2"}], "isLast": true}`))
	case "/rest/api/latest/group/member?groupname=group1&startAt=0&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"key": "member1"}], "isLast": false}`))
	case "/rest/api/latest/group/member?groupname=group1&startAt=1&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"key": "member2"}], "isLast": true}`))
	// Group2 has 0 members.
	case "/rest/api/latest/group/member?groupname=group2&startAt=0&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [], "isLast": true}`))
	// Group3 has 2 members.
	case "/rest/api/latest/group/member?groupname=group3&startAt=0&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"key": "member3"}], "isLast": false}`))
	case "/rest/api/latest/group/member?groupname=group3&startAt=1&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"key": "member4"}], "isLast": true}`))

	// These endpoints define cases where tests should fail, e.g. missing fields, empty, etc.
	// Hence, they start from page 99 to avoid colliding with the above endpoints.
	// Return an empty list of groups.
	// Omit the Group's uniqueId.
	case "/rest/api/failing-version-one/groups/picker":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"groups": [{"NOT_UNIQUE_ID":"group1"}]}`))
	// Make the Group's uniqueId not parsable as a string.
	case "/rest/api/failing-version-two/groups/picker":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"groups": [{"name":10}]}`))

	// Create a group member uniqueId that is not parsable into a string.
	case "/rest/api/latest/group/member?groupname=group1&startAt=99&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"values": [{"key": 4}], "isLast": true}`))
	// Create a group member response structure that is not expected.
	case "/rest/api/latest/group/member?groupname=group1&startAt=100&maxResults=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"UNEXPECTED_FIELD": [{"key": 4}], "isLast": true}`))

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})

func TestParseUsersResponse(t *testing.T) {
	tests := map[string]struct {
		body           []byte
		cursor         int64
		pageSize       int64
		wantObjects    []map[string]interface{}
		wantNextCursor *int64
		wantErr        *framework.Error
	}{
		"user_objects_last_page": {
			// Two user objects in response with page size = 10, so this must be last page.
			body:     []byte(`{"values": [{"key": "user1"}, {"key": "user2"}]}`),
			cursor:   0,
			pageSize: 10,
			wantObjects: []map[string]interface{}{
				{"key": "user1"},
				{"key": "user2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"user_objects_not_last_page": {
			// Two user objects in response with page size = 2, so there is a possibility of next page.
			body:     []byte(`{"values": [{"key": "user1"}, {"key": "user2"}]}`),
			cursor:   0,
			pageSize: 2,
			wantObjects: []map[string]interface{}{
				{"key": "user1"},
				{"key": "user2"},
			},
			wantNextCursor: testutil.GenPtr(int64(2)), // This page contains index 0 and 1, so next page starts at index 2.
			wantErr:        nil,
		},
		"invalid_user_response": {
			// Users response, which is group member response, should return a top level values object, not a list of users.
			body:           []byte(`[{"key": "user1"}, {"key": "user2"}]`),
			cursor:         0,
			pageSize:       2,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Failed to unmarshal Jira User response: json: cannot unmarshal array into " +
					"Go value of type map[string]interface {}.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"response_contains_valid_is_last_field": {
			// The next cursor does not need to be computed if `isLast` field is present.
			body:     []byte(`{"values": [{"key": "user1"}, {"key": "user2"}], "isLast": true}`),
			cursor:   0,
			pageSize: 2,
			wantObjects: []map[string]interface{}{
				{"key": "user1"},
				{"key": "user2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			request := &jiradatacenter.Request{
				PageSize: tt.pageSize,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(tt.cursor),
				},
			}

			userEntity := jiradatacenter.ValidEntityExternalIDs[jiradatacenter.UserExternalID]

			gotObjects, gotNextCursor, gotErr := userEntity.Parse(tt.body, *request)

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
		cursor         int64
		pageSize       int64
		wantObjects    []map[string]interface{}
		wantNextCursor *int64
		wantErr        *framework.Error
	}{
		"issue_objects_last_page": {
			body:     []byte(`{"issues": [{"id": "issue1"}, {"id": "issue2"}]}`),
			cursor:   0,
			pageSize: 10,
			wantObjects: []map[string]interface{}{
				{"id": "issue1"},
				{"id": "issue2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"issue_objects_not_last_page": {
			body:     []byte(`{"issues": [{"id": "issue1"}, {"id": "issue2"}]}`),
			cursor:   0,
			pageSize: 2,
			wantObjects: []map[string]interface{}{
				{"id": "issue1"},
				{"id": "issue2"},
			},
			wantNextCursor: testutil.GenPtr(int64(2)),
			wantErr:        nil,
		},
		"invalid_issue_response": {
			// Issues response should return a single top level object, not a list.
			body:           []byte(`[{"id": "issue1"}, {"id": "issue2"}]`),
			cursor:         0,
			pageSize:       2,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Failed to unmarshal Jira Issue response: json: cannot unmarshal array into " +
					"Go value of type map[string]interface {}.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"invalid_issue_object": {
			// The "issues" value should return []map[string]any, not []any.
			body:           []byte(`{"issues": ["issue1", {"id": "issue2"}]}`),
			cursor:         0,
			pageSize:       2,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "An object in Entity: Issue could not be parsed. Expected: map[string]any. Got: string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"issues_field_does_not_exist": {
			body:           []byte(`{"WRONG_FIELD": [{"id": "issue1"}, {"id": "issue2"}]}`),
			cursor:         0,
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Field missing in Jira issues response: Issue.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"issues_field_exists_but_invalid_format": {
			// The "issues" field value should be a list of issue objects, not a map.
			body:           []byte(`{"issues": {"id": "issue1"}}`),
			cursor:         0,
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Entity field exists in Jira Issue response but field value is not a list of objects: " +
					"map[string]interface {}.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			request := &jiradatacenter.Request{
				PageSize: tt.pageSize,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(tt.cursor),
				},
			}

			issueEntity := jiradatacenter.ValidEntityExternalIDs[jiradatacenter.IssueExternalID]

			gotObjects, gotNextCursor, gotErr := issueEntity.Parse(tt.body, *request)

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
		cursor         int64
		pageSize       int64
		wantObjects    []map[string]interface{}
		wantNextCursor *int64
		wantErr        *framework.Error
	}{
		"group_objects_last_page": {
			body:     []byte(`{"groups": [{"name": "group1"}, {"name": "group2"}]}`),
			cursor:   0,
			pageSize: 10,
			wantObjects: []map[string]interface{}{
				{"name": "group1"},
				{"name": "group2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"group_objects_not_last_page": {
			body:     []byte(`{"groups": [{"name": "group1"}, {"name": "group2"}, {"name": "group3"}]}`),
			cursor:   0,
			pageSize: 2,
			wantObjects: []map[string]interface{}{
				{"name": "group1"},
				{"name": "group2"},
			},
			wantNextCursor: testutil.GenPtr(int64(2)),
			wantErr:        nil,
		},
		"invalid_group_response": {
			// Groups response should return a single top level object, not a list.
			body:           []byte(`[{"name": "group1"}, {"name": "group2"}]`),
			cursor:         0,
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
			body:           []byte(`{"groups": ["group1", {"name": "group2"}]}`),
			cursor:         0,
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
			cursor:         0,
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Field missing in Jira groups response: Group.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"values_field_exists_but_invalid_format": {
			// The "values" field value should be a list of group objects, not a map.
			body:           []byte(`{"groups": {"name": "group1"}}`),
			cursor:         0,
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Entity field exists in Jira Group response but field value is not a list of objects: " +
					"map[string]interface {}.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			request := &jiradatacenter.Request{
				PageSize: tt.pageSize,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(tt.cursor),
				},
			}

			groupEntity := jiradatacenter.ValidEntityExternalIDs[jiradatacenter.GroupExternalID]

			gotObjects, gotNextCursor, gotErr := groupEntity.Parse(tt.body, *request)

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
		cursor         int64
		pageSize       int64
		wantObjects    []map[string]interface{}
		wantNextCursor *int64
		wantErr        *framework.Error
	}{
		"group_member_objects_last_page": {
			body:     []byte(`{"values": [{"key": "groupMember1"}, {"key": "groupMember2"}]}`),
			cursor:   0,
			pageSize: 10,
			wantObjects: []map[string]interface{}{
				{"key": "groupMember1"},
				{"key": "groupMember2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"group_objects_not_last_page": {
			body:     []byte(`{"values": [{"key": "groupMember1"}, {"key": "groupMember2"}]}`),
			cursor:   0,
			pageSize: 2,
			wantObjects: []map[string]interface{}{
				{"key": "groupMember1"},
				{"key": "groupMember2"},
			},
			wantNextCursor: testutil.GenPtr(int64(2)),
			wantErr:        nil,
		},
		"invalid_group_members_response": {
			// Groups response should return a single top level object, not a list.
			body:           []byte(`[{"key": "groupMember1"}, {"key": "groupMember2"}]`),
			cursor:         0,
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
			body:           []byte(`{"values": ["groupMember1", {"key": "groupMember2"}]}`),
			cursor:         0,
			pageSize:       2,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "An object in Entity: GroupMember could not be parsed. Expected: map[string]any. Got: string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"values_field_does_not_exist": {
			body:           []byte(`{"WRONG_FIELD": [{"key": "groupMember1"}, {"key": "groupMember2"}]}`),
			cursor:         0,
			pageSize:       10,
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "Field missing in Jira values response: GroupMember.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"values_field_exists_but_invalid_format": {
			// The "values" field value should be a list of group objects, not a map.
			body:           []byte(`{"values": {"key": "groupMember1"}}`),
			cursor:         0,
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
			body:     []byte(`{"values": [{"key": "groupMember1"}, {"key": "groupMember2"}], "isLast": true}`),
			cursor:   0,
			pageSize: 2,
			wantObjects: []map[string]interface{}{
				{"key": "groupMember1"},
				{"key": "groupMember2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			request := &jiradatacenter.Request{
				PageSize: tt.pageSize,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(tt.cursor),
				},
			}

			groupMemberEntity := jiradatacenter.ValidEntityExternalIDs[jiradatacenter.GroupMemberExternalID]

			gotObjects, gotNextCursor, gotErr := groupMemberEntity.Parse(tt.body, *request)

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
		request *jiradatacenter.Request
		entity  jiradatacenter.Entity
		cursor  *pagination.CompositeCursor[int64]
		wantURL string
		wantErr error
	}{
		"users": {
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jiradatacenter.UserExternalID,
			},
			entity: jiradatacenter.ValidEntityExternalIDs[jiradatacenter.UserExternalID],
			cursor: &pagination.CompositeCursor[int64]{
				Cursor:           testutil.GenPtr[int64](10),
				CollectionID:     testutil.GenPtr("group1"),
				CollectionCursor: testutil.GenPtr[int64](1),
			},
			wantURL: "https://jira.com/rest/api/latest/group/member?groupname=group1&startAt=10&maxResults=10",
		},
		"users_with_inactive_users": {
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jiradatacenter.UserExternalID,
				IncludeInactiveUsers:  testutil.GenPtr(true),
			},
			entity: jiradatacenter.ValidEntityExternalIDs[jiradatacenter.UserExternalID],
			cursor: &pagination.CompositeCursor[int64]{
				Cursor:           testutil.GenPtr[int64](10),
				CollectionID:     testutil.GenPtr("group1"),
				CollectionCursor: testutil.GenPtr[int64](1),
			},
			wantURL: "https://jira.com/rest/api/latest/group/member?groupname=group1" +
				"&includeInactiveUsers=true&startAt=10&maxResults=10",
		},
		"groups": {
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jiradatacenter.GroupExternalID,
			},
			entity: jiradatacenter.ValidEntityExternalIDs[jiradatacenter.GroupExternalID],
			cursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](10),
			},
			wantURL: "https://jira.com/rest/api/latest/groups/picker",
		},
		"groups_with_inactive_users_set_true": {
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jiradatacenter.GroupExternalID,
				IncludeInactiveUsers:  testutil.GenPtr(true),
			},
			entity: jiradatacenter.ValidEntityExternalIDs[jiradatacenter.GroupExternalID],
			cursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](10),
			},
			wantURL: "https://jira.com/rest/api/latest/groups/picker",
		},
		"issues": {
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jiradatacenter.IssueExternalID,
			},
			entity: jiradatacenter.ValidEntityExternalIDs[jiradatacenter.IssueExternalID],
			cursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](10),
			},
			wantURL: "https://jira.com/rest/api/latest/search?fields=*navigable&startAt=10&maxResults=10",
		},
		"issues_with_filter": {
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jiradatacenter.IssueExternalID,
				IssuesJQLFilter:       testutil.GenPtr("project=TEST"),
			},
			entity: jiradatacenter.ValidEntityExternalIDs[jiradatacenter.IssueExternalID],
			cursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](10),
			},
			wantURL: "https://jira.com/rest/api/latest/search?jql=project%3DTEST&fields=*navigable&startAt=10&maxResults=10",
		},
		"issues_with_attributes": {
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jiradatacenter.IssueExternalID,
				Entity: &framework.EntityConfig{
					Attributes: []*framework.AttributeConfig{
						{ExternalId: "id"},
						{ExternalId: "key"},
						{ExternalId: "$.fields.summary"},
						{ExternalId: "$.fields.status.id"},
					},
				},
			},
			entity: jiradatacenter.ValidEntityExternalIDs[jiradatacenter.IssueExternalID],
			cursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](10),
			},
			wantURL: "https://jira.com/rest/api/latest/search?fields=id%2Ckey%2Cstatus%2Csummary&startAt=10&maxResults=10",
		},
		"group_members": {
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jiradatacenter.GroupMemberExternalID,
			},
			entity: jiradatacenter.ValidEntityExternalIDs[jiradatacenter.GroupMemberExternalID],
			cursor: &pagination.CompositeCursor[int64]{
				Cursor:           testutil.GenPtr[int64](10),
				CollectionID:     testutil.GenPtr("1"),
				CollectionCursor: testutil.GenPtr[int64](1),
			},
			wantURL: "https://jira.com/rest/api/latest/group/member?groupname=1&startAt=10&maxResults=10",
		},
		"group_members_with_inactive_users": {
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jiradatacenter.GroupMemberExternalID,
				IncludeInactiveUsers:  testutil.GenPtr(true),
			},
			entity: jiradatacenter.ValidEntityExternalIDs[jiradatacenter.GroupMemberExternalID],
			cursor: &pagination.CompositeCursor[int64]{
				Cursor:           testutil.GenPtr[int64](10),
				CollectionID:     testutil.GenPtr("1"),
				CollectionCursor: testutil.GenPtr[int64](1),
			},
			wantURL: "https://jira.com/rest/api/latest/group/member?groupname=1" +
				"&includeInactiveUsers=true&startAt=10&maxResults=10",
		},
		"group_members_missing_group_id": {
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://jira.com",
				PageSize:              10,
				EntityExternalID:      jiradatacenter.GroupMemberExternalID,
			},
			entity: jiradatacenter.ValidEntityExternalIDs[jiradatacenter.GroupMemberExternalID],
			cursor: &pagination.CompositeCursor[int64]{
				Cursor:           testutil.GenPtr[int64](10),
				CollectionCursor: testutil.GenPtr[int64](1),
			},
			wantErr: errors.New("cursor.CollectionID must not be nil for User entity or GroupMember entity"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			entity := jiradatacenter.ValidEntityExternalIDs[tt.request.EntityExternalID]

			gotURL, gotErr := entity.ConstructURL(tt.request, tt.cursor)

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
		client: jiradatacenter.NewClient(client),
		server: httptest.NewServer(TestServerHandler),
	}
	defer ts.server.Close()
	t.Run("TestGetPageErrors", ts.TestGetPageErrors)
	t.Run("TestGetPageUsers", ts.TestGetPageUsers)
	t.Run("TestGetPageIssues", ts.TestGetPageIssues)
	t.Run("TestGetPageGroups", ts.TestGetPageGroups)
	t.Run("TestGetPageGroupMembers", ts.TestGetPageGroupMembers)
}

func (ts *TestSuite) TestGetPageErrors(t *testing.T) {
	tests := map[string]struct {
		ctx          context.Context
		request      *jiradatacenter.Request
		wantResponse *jiradatacenter.Response
		wantErr      *framework.Error
	}{
		"invalid_url": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "https://{hello}",
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: `Address in datasource config is an invalid URL: parse ` +
					`"https://{hello}/rest/api/latest/?startAt=0&maxResults=0": invalid character "{" in host name.`,
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "BAD_PROTOCOL",
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: `Failed to execute Jira request: Get "BAD_PROTOCOL/rest/api/latest/?startAt=0&maxResults=0": ` +
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
	externalEntityID := jiradatacenter.UserExternalID

	tests := map[string]struct {
		ctx          context.Context
		request      *jiradatacenter.Request
		wantResponse *jiradatacenter.Response
		wantErr      *framework.Error
		expectedLogs []map[string]any
	}{
		"first_page_first_group": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"key": "member1"},
					{"key": "member2"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					// Since the page size = 10, and group1 only has 2 members, we should not have a
					// next cursor. We can sync all group members for group1 in one page.
					Cursor: nil,
					// CollectionID is the ID of the group that we're currently/just finished syncing.
					CollectionID: testutil.GenPtr("group1"),
					// CollectionCursor is the cursor of the NEXT group that we're going to sync.
					CollectionCursor: testutil.GenPtr(int64(1)),
				},
			},
			wantErr: nil,
			expectedLogs: []map[string]any{
				{
					"level":                             "info",
					"msg":                               "Starting datasource request",
					fields.FieldRequestEntityExternalID: "User",
					fields.FieldRequestPageSize:         int64(10),
				},
				{
					"level":                             "info",
					"msg":                               "Starting datasource request",
					fields.FieldRequestEntityExternalID: "Group",
					fields.FieldRequestPageSize:         int64(1),
				},
				{
					"level":                             "info",
					"msg":                               "Sending HTTP request to datasource",
					fields.FieldRequestEntityExternalID: "Group",
					fields.FieldRequestPageSize:         int64(1),
					fields.FieldRequestURL:              ts.server.URL + "/rest/api/latest/groups/picker",
				},
				{
					"level":                             "info",
					"msg":                               "Datasource request completed successfully",
					fields.FieldRequestEntityExternalID: "Group",
					fields.FieldRequestPageSize:         int64(1),
					fields.FieldResponseStatusCode:      int64(200),
					fields.FieldResponseObjectCount:     int64(1),
					fields.FieldResponseNextCursor: map[string]any{
						"cursor": int64(1),
					},
				},
				{
					"level":                             "info",
					"msg":                               "Sending HTTP request to datasource",
					fields.FieldRequestEntityExternalID: "User",
					fields.FieldRequestPageSize:         int64(10),
					// nolint: lll
					fields.FieldRequestURL: ts.server.URL + "/rest/api/latest/group/member?groupname=group1&startAt=0&maxResults=10",
				},
				{
					"level":                             "info",
					"msg":                               "Datasource request completed successfully",
					fields.FieldRequestEntityExternalID: "User",
					fields.FieldRequestPageSize:         int64(10),
					fields.FieldResponseStatusCode:      int64(200),
					fields.FieldResponseObjectCount:     int64(2),
					fields.FieldResponseNextCursor: map[string]any{
						"collectionId":     "group1",
						"collectionCursor": int64(1),
					},
				},
			},
		},
		"first_page_first_group_page_size_1": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"key": "member1"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					// Since the page size = 1, and group1 has 2 members, we should have a
					// next cursor.
					Cursor: testutil.GenPtr(int64(1)),
					// CollectionID is the ID of the group that we're currently/just finished syncing.
					CollectionID: testutil.GenPtr("group1"),
					// CollectionCursor is the cursor of the NEXT group that we're going to sync.
					CollectionCursor: testutil.GenPtr(int64(1)),
				},
			},
			wantErr: nil,
		},
		"second_page_first_group_page_size_1": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr(int64(1)),
					CollectionID:     testutil.GenPtr("group1"),
					CollectionCursor: testutil.GenPtr(int64(1)),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"key": "member2"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					// This is the last group member of this group. There should not be a next cursor.
					Cursor: nil,
					// CollectionID is the ID of the group that we're currently/just finished syncing.
					CollectionID: testutil.GenPtr("group1"),
					// CollectionCursor is the cursor of the NEXT group that we're going to sync.
					CollectionCursor: testutil.GenPtr(int64(1)),
				},
			},
			wantErr: nil,
		},
		"empty_page_second_group_page_size_1": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					CollectionID:     testutil.GenPtr("group1"),
					CollectionCursor: testutil.GenPtr(int64(1)),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects:    []map[string]any{},
				NextCursor: &pagination.CompositeCursor[int64]{
					// Group2 has no members. There should be no next cursor.
					Cursor: nil,
					// CollectionID is the ID of the group that we're currently/just finished syncing.
					CollectionID: testutil.GenPtr("group2"),
					// CollectionCursor is the cursor of the NEXT group that we're going to sync.
					CollectionCursor: testutil.GenPtr(int64(2)),
				},
			},
			wantErr: nil,
		},
		"first_page_third_group_page_size_1": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					CollectionID:     testutil.GenPtr("group2"),
					CollectionCursor: testutil.GenPtr(int64(2)),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"key": "member3"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					// Group3 has two members. We should have a next cursor.
					Cursor: testutil.GenPtr(int64(1)),
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
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr(int64(1)),
					CollectionID:     testutil.GenPtr("group3"),
					CollectionCursor: nil,
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"key": "member4"},
				},
				// Group3 has two members. We've synced all of them. There should be no next cursor.
				// Group3 is also the last group, so CollectionCursor should be nil.
				// Both these conditions imply we have finished syncing all group members.
				// The entire NextCursor should be nil.
				NextCursor: nil,
			},
			wantErr: nil,
		},
		"first_page_first_group_with_config_groups": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
				Groups:                []string{"group3"},
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"key": "member3"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr(int64(1)),
					CollectionID:     testutil.GenPtr("group3"),
					CollectionCursor: nil,
				},
			},
			wantErr: nil,
		},
		// If there are no groups to begin with, there is no need to sync any group members.
		"no_groups": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					CollectionCursor: testutil.GenPtr(int64(99)),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
			},
			wantErr: nil,
		},
		"group_unique_id_missing": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor:                &pagination.CompositeCursor[int64]{},
				EntityExternalID:      externalEntityID,
				APIVersion:            "failing-version-one",
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: "Jira Group object contains no name field.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
			},
		},
		"group_unique_id_not_string": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor:                &pagination.CompositeCursor[int64]{},
				EntityExternalID:      externalEntityID,
				APIVersion:            "failing-version-two",
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: "Failed to convert Jira Group object name field to string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
			},
		},
		// If we're syncing group members, we must have a group id to sync.
		"composite_cursor_missing_group_id": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					// Need to set Cursor here otherwise the code will think we're syncing the first page,
					// which doesn't require a group id.
					Cursor:           testutil.GenPtr(int64(1)),
					CollectionCursor: testutil.GenPtr(int64(1)),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: "Cursor does not have CollectionID set for entity User.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		// On the first page sync of a GroupMember, we make a request for the first group. If that request fails,
		// we should see an error.
		"group_get_page_fails_when_nil_cursor": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "http://localhost:1234",
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor:                nil,
				EntityExternalID:      externalEntityID,
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: `Failed to execute Jira request: Get "http://localhost:1234/rest/api/latest/groups/picker": ` +
					`dial tcp [::1]:1234: connect: connection refused.`,
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		// If the GroupMember request is successful, but the response structure is not what we expect
		// (e.g. missing a field), we should see an error.
		"group_member_response_structure_is_invalid": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr(int64(100)),
					CollectionID:     testutil.GenPtr("group1"),
					CollectionCursor: testutil.GenPtr(int64(1)),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: "Field missing in Jira values response: User.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
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

func (ts *TestSuite) TestGetPageIssues(t *testing.T) {
	externalEntityID := jiradatacenter.IssueExternalID

	tests := map[string]struct {
		ctx          context.Context
		request      *jiradatacenter.Request
		wantResponse *jiradatacenter.Response
		wantErr      *framework.Error
	}{
		"first_page_no_next_cursor": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
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
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "1"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(1)),
				},
			},
			wantErr: nil,
		},
		"second_last_page_last_issue": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				// The last issue occurs on page 2.
				Cursor:           &pagination.CompositeCursor[int64]{Cursor: testutil.GenPtr(int64(1))},
				PageSize:         int64(1),
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "2"},
				},
				// We've synced the last issue but it's not possible for our code to know that.
				// We still need to check if there are more results.
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(2)),
				},
			},
			wantErr: nil,
		},
		"last_page_empty": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				Cursor:                &pagination.CompositeCursor[int64]{Cursor: testutil.GenPtr(int64(2))},
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects:    []map[string]any{},
				// An empty response means we've completed the sync.
				NextCursor: nil,
			},
			wantErr: nil,
		},
		"with_jql_filter": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
				IssuesJQLFilter:       testutil.GenPtr("project='SGNL'"),
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "99"},
				},
			},
		},
		"with_invalid_jql_filter": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
				IssuesJQLFilter:       testutil.GenPtr("project='INVALID'"),
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 400,
			},
		},
		"with_entity_config": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
				IssuesJQLFilter:       testutil.GenPtr("project='SGNL'"),
				Entity: &framework.EntityConfig{
					Attributes: []*framework.AttributeConfig{
						{ExternalId: "id"},
					},
				},
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "99"},
				},
			},
		},
		"with_child_entities_present_in_response": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
				IssuesJQLFilter:       testutil.GenPtr("project='CHILD_ENTITIES_PRESENT'"),
				Entity: &framework.EntityConfig{
					Attributes: []*framework.AttributeConfig{
						{ExternalId: "$.id"},
						{ExternalId: "$.key"},
						{ExternalId: "$.fields.summary"},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "issuetype",
							Attributes: []*framework.AttributeConfig{
								{ExternalId: "id"},
								{ExternalId: "name"},
								{ExternalId: "description"},
							},
						},
						{
							ExternalId: "assignee",
							Attributes: []*framework.AttributeConfig{
								{ExternalId: "accountId"},
								{ExternalId: "displayName"},
								{ExternalId: "emailAddress"},
							},
						},
					},
				},
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{
						"id":  "ISSUE-100",
						"key": "TEST-100",
						"fields": map[string]any{
							"summary": []any{
								map[string]any{
									"id":          "100",
									"description": "Issue with child entities",
								},
							},
							"assignee": []any{
								map[string]any{
									"accountId":    "user123",
									"displayName":  "John Doe",
									"emailAddress": "john.doe@example.com",
								},
							},
							"issuetype": []any{
								map[string]any{
									"description": "A bug issue type",
									"id":          "1",
									"name":        "Bug",
								},
							},
						},
					},
				},
			},
		},
		"with_child_entities_missing_from_response": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
				IssuesJQLFilter:       testutil.GenPtr("project='CHILD_ENTITIES_MISSING'"),
				Entity: &framework.EntityConfig{
					Attributes: []*framework.AttributeConfig{
						{ExternalId: "$.id"},
						{ExternalId: "$.key"},
						{ExternalId: "$.fields.summary"},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "issuetype",
							Attributes: []*framework.AttributeConfig{
								{ExternalId: "id"},
								{ExternalId: "name"},
							},
						},
						{
							ExternalId: "assignee",
							Attributes: []*framework.AttributeConfig{
								{ExternalId: "accountId"},
								{ExternalId: "displayName"},
							},
						},
					},
				},
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{
						"id":  "ISSUE-101",
						"key": "TEST-101",
						"fields": map[string]any{
							"summary": []any{
								map[string]any{
									"id":          "100",
									"description": "Issue with child entities",
								},
							},
							// Note: issuetype and assignee fields are missing from response
						},
					},
				},
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
	externalEntityID := jiradatacenter.GroupExternalID

	tests := map[string]struct {
		ctx          context.Context
		request      *jiradatacenter.Request
		wantResponse *jiradatacenter.Response
		wantErr      *framework.Error
	}{
		"first_page_no_next_cursor": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"name": "group1"},
					{"name": "group2"},
					{"name": "group3"},
				},
			},
			wantErr: nil,
		},
		"first_page_with_next_cursor": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"name": "group1"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(1)),
				},
			},
			wantErr: nil,
		},
		"middle_page_with_next_cursor": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				Cursor:                &pagination.CompositeCursor[int64]{Cursor: testutil.GenPtr(int64(1))},
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"name": "group2"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(2)),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				// The last group occurs on page 3.
				Cursor:           &pagination.CompositeCursor[int64]{Cursor: testutil.GenPtr(int64(2))},
				PageSize:         int64(1),
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"name": "group3"},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
		"no_group_found_from_the_specified_config_group_list": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
				Groups:                []string{"group5", "group6"},
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects:    []map[string]any{},
			},
			wantErr: nil,
		},
		"only_one_group_found_from_the_specified_config_group_list": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
				Groups:                []string{"group3", "group6"},
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"name": "group3"},
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

func (ts *TestSuite) TestGetPageGroupMembers(t *testing.T) {
	externalEntityID := jiradatacenter.GroupMemberExternalID

	tests := map[string]struct {
		ctx          context.Context
		request      *jiradatacenter.Request
		wantResponse *jiradatacenter.Response
		wantErr      *framework.Error
	}{
		"first_page_first_group": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(10),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "group1-member1", "key": "member1", "groupId": "group1"},
					{"id": "group1-member2", "key": "member2", "groupId": "group1"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					// Since the page size = 10, and group1 only has 2 members, we should not have a
					// next cursor. We can sync all group members for group1 in one page.
					Cursor: nil,
					// CollectionID is the ID of the group that we're currently/just finished syncing.
					CollectionID: testutil.GenPtr("group1"),
					// CollectionCursor is the cursor of the NEXT group that we're going to sync.
					CollectionCursor: testutil.GenPtr(int64(1)),
				},
			},
			wantErr: nil,
		},
		"first_page_first_group_page_size_1": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "group1-member1", "key": "member1", "groupId": "group1"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					// Since the page size = 1, and group1 has 2 members, we should have a
					// next cursor.
					Cursor: testutil.GenPtr(int64(1)),
					// CollectionID is the ID of the group that we're currently/just finished syncing.
					CollectionID: testutil.GenPtr("group1"),
					// CollectionCursor is the cursor of the NEXT group that we're going to sync.
					CollectionCursor: testutil.GenPtr(int64(1)),
				},
			},
			wantErr: nil,
		},
		"second_page_first_group_page_size_1": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr(int64(1)),
					CollectionID:     testutil.GenPtr("group1"),
					CollectionCursor: testutil.GenPtr(int64(1)),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "group1-member2", "key": "member2", "groupId": "group1"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					// This is the last group member of this group. There should not be a next cursor.
					Cursor: nil,
					// CollectionID is the ID of the group that we're currently/just finished syncing.
					CollectionID: testutil.GenPtr("group1"),
					// CollectionCursor is the cursor of the NEXT group that we're going to sync.
					CollectionCursor: testutil.GenPtr(int64(1)),
				},
			},
			wantErr: nil,
		},
		"empty_page_second_group_page_size_1": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					CollectionID:     testutil.GenPtr("group1"),
					CollectionCursor: testutil.GenPtr(int64(1)),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects:    []map[string]any{},
				NextCursor: &pagination.CompositeCursor[int64]{
					// Group2 has no members. There should be no next cursor.
					Cursor: nil,
					// CollectionID is the ID of the group that we're currently/just finished syncing.
					CollectionID: testutil.GenPtr("group2"),
					// CollectionCursor is the cursor of the NEXT group that we're going to sync.
					CollectionCursor: testutil.GenPtr(int64(2)),
				},
			},
			wantErr: nil,
		},
		"first_page_third_group_page_size_1": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					CollectionID:     testutil.GenPtr("group2"),
					CollectionCursor: testutil.GenPtr(int64(2)),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "group3-member3", "key": "member3", "groupId": "group3"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					// Group3 has two members. We should have a next cursor.
					Cursor: testutil.GenPtr(int64(1)),
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
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr(int64(1)),
					CollectionID:     testutil.GenPtr("group3"),
					CollectionCursor: nil,
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "group3-member4", "key": "member4", "groupId": "group3"},
				},
				// Group3 has two members. We've synced all of them. There should be no next cursor.
				// Group3 is also the last group, so CollectionCursor should be nil.
				// Both these conditions imply we have finished syncing all group members.
				// The entire NextCursor should be nil.
				NextCursor: nil,
			},
			wantErr: nil,
		},
		"first_page_first_group_with_config_groups": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				EntityExternalID:      externalEntityID,
				Groups:                []string{"group3"},
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"id": "group3-member3", "key": "member3", "groupId": "group3"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr(int64(1)),
					CollectionID:     testutil.GenPtr("group3"),
					CollectionCursor: nil,
				},
			},
			wantErr: nil,
		},
		// If there are no groups to begin with, there is no need to sync any group members.
		"no_groups": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					CollectionCursor: testutil.GenPtr(int64(99)),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: &jiradatacenter.Response{
				StatusCode: 200,
			},
			wantErr: nil,
		},
		"group_unique_id_missing": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor:                &pagination.CompositeCursor[int64]{},
				EntityExternalID:      externalEntityID,
				APIVersion:            "failing-version-one",
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: "Jira Group object contains no name field.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
			},
		},
		"group_unique_id_not_string": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor:                &pagination.CompositeCursor[int64]{},
				EntityExternalID:      externalEntityID,
				APIVersion:            "failing-version-two",
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: "Failed to convert Jira Group object name field to string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
			},
		},
		"group_member_unique_id_not_string": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr(int64(99)),
					CollectionID:     testutil.GenPtr("group1"),
					CollectionCursor: testutil.GenPtr(int64(3)),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: "Failed to parse key field in Jira GroupMember response as string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		// If we're syncing group members, we must have a group id to sync.
		"composite_cursor_missing_group_id": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					// Need to set Cursor here otherwise the code will think we're syncing the first page,
					// which doesn't require a group id.
					Cursor:           testutil.GenPtr(int64(1)),
					CollectionCursor: testutil.GenPtr(int64(1)),
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
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               "http://localhost:1234",
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor:                nil,
				EntityExternalID:      externalEntityID,
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: `Failed to execute Jira request: Get "http://localhost:1234/rest/api/latest/groups/picker": ` +
					`dial tcp [::1]:1234: connect: connection refused.`,
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		// If the GroupMember request is successful, but the response structure is not what we expect
		// (e.g. missing a field), we should see an error.
		"group_member_response_structure_is_invalid": {
			ctx: context.Background(),
			request: &jiradatacenter.Request{
				RequestTimeoutSeconds: 5,
				BaseURL:               ts.server.URL,
				AuthorizationHeader:   mockAuthorizationHeader,
				PageSize:              int64(1),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr(int64(100)),
					CollectionID:     testutil.GenPtr("group1"),
					CollectionCursor: testutil.GenPtr(int64(1)),
				},
				EntityExternalID: externalEntityID,
			},
			wantResponse: nil,
			wantErr: &framework.Error{
				Message: "Field missing in Jira values response: GroupMember.",
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

func TestBuildJiraFieldsParam(t *testing.T) {
	tests := []struct {
		name   string
		entity *framework.EntityConfig
		want   string
	}{
		{
			name: "empty attributes list",
			entity: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{},
			},
			want: "*navigable",
		},
		{
			name: "nil attributes list",
			entity: &framework.EntityConfig{
				Attributes: nil,
			},
			want: "*navigable",
		},
		{
			name: "single field extraction",
			entity: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.fields.summary"},
				},
			},
			want: "summary",
		},
		{
			name: "multiple different fields",
			entity: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.fields.summary"},
					{ExternalId: "$.fields.description"},
				},
			},
			want: "description%2Csummary", // URL encoded "description,summary" (sorted)
		},
		{
			name: "nested field extraction",
			entity: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.fields.issuetype.id"},
				},
			},
			want: "issuetype",
		},
		{
			name: "deduplication of same field",
			entity: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.fields.assignee.key"},
					{ExternalId: "$.fields.assignee.name"},
				},
			},
			want: "assignee",
		},
		{
			name: "root level fields",
			entity: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.id"},
					{ExternalId: "$.key"},
					{ExternalId: "$.self"},
				},
			},
			want: "id%2Ckey%2Cself", // URL encoded "id,key,self" (sorted)
		},
		{
			name: "mix of different field types",
			entity: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.fields.summary"},
					{ExternalId: "$.fields.issuetype.id"},
					{ExternalId: "$.fields.assignee.key"},
					{ExternalId: "$.fields.assignee.name"}, // duplicate assignee
					{ExternalId: "$.id"},
					{ExternalId: "$.key"},
				},
			},
			want: "assignee%2Cid%2Cissuetype%2Ckey%2Csummary", // URL encoded and sorted
		},
		{
			name: "attributes with empty ExternalId",
			entity: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: ""},
					{ExternalId: "$.fields.summary"},
					{ExternalId: ""},
				},
			},
			want: "summary",
		},
		{
			name: "all attributes have empty ExternalId",
			entity: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: ""},
					{ExternalId: ""},
				},
			},
			want: "*navigable",
		},
		{
			name: "custom fields",
			entity: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.fields.customfield_10000"},
					{ExternalId: "$.fields.customfield_10001.value"},
				},
			},
			want: "customfield_10000%2Ccustomfield_10001", // URL encoded and sorted
		},
		{
			name: "field names with special characters needing URL encoding",
			entity: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.fields.summary"},
					{ExternalId: "$.fields.status"},
					{ExternalId: "$.fields.project"},
				},
			},
			want: "project%2Cstatus%2Csummary", // URL encoded "project,status,summary" (sorted)
		},
		{
			name: "non-JSON path field names",
			entity: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "id"},
					{ExternalId: "key"},
					{ExternalId: "summary"},
				},
			},
			want: "id%2Ckey%2Csummary", // URL encoded "id,key,summary" (sorted)
		},
		{
			name: "array indices in field names",
			entity: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.fields.customfield_10209[0].value"},
					{ExternalId: "$.fields.customfield_10210[0]"},
					{ExternalId: "$.fields.assignee[0].key"},
				},
			},
			want: "assignee%2Ccustomfield_10209%2Ccustomfield_10210", // URL encoded and sorted, array indices removed
		},
		{
			name: "mixed array and non-array fields with deduplication",
			entity: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.fields.customfield_10209[0].value"},
					{ExternalId: "$.fields.customfield_10209.id"}, // same field without array index
					{ExternalId: "$.fields.assignee[0].key"},
					{ExternalId: "$.fields.assignee.name"}, // same field without array index
				},
			},
			want: "assignee%2Ccustomfield_10209", // URL encoded and sorted, deduplicated
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := jiradatacenter.BuildJiraFieldsParam(tt.entity)
			if got != tt.want {
				t.Errorf("BuildJiraFieldsParam() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractEntityFieldNames(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		entity *framework.EntityConfig
		want   map[string]struct{}
	}{
		{
			name:   "nil_entity_returns_empty_map",
			prefix: "",
			entity: nil,
			want:   map[string]struct{}{},
		},
		{
			name:   "entity_with_simple_attributes",
			prefix: "",
			entity: &framework.EntityConfig{
				ExternalId: "Issue",
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "id"},
					{ExternalId: "key"},
					{ExternalId: "$.fields.summary"},
				},
			},
			want: map[string]struct{}{
				"id":      {},
				"key":     {},
				"summary": {},
			},
		},
		{
			name:   "entity_with_json_path_attributes",
			prefix: "",
			entity: &framework.EntityConfig{
				ExternalId: "Issue",
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.id"},
					{ExternalId: "$.fields.assignee.key"},
					{ExternalId: "$.fields.issuetype.id"},
					{ExternalId: "$.fields.customfield_10209[0].value"},
				},
			},
			want: map[string]struct{}{
				"id":                {},
				"assignee":          {},
				"issuetype":         {},
				"customfield_10209": {},
			},
		},
		{
			name:   "entity_with_duplicate_field_names",
			prefix: "",
			entity: &framework.EntityConfig{
				ExternalId: "Issue",
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.fields.assignee.key"},
					{ExternalId: "$.fields.assignee.name"},
					{ExternalId: "$.fields.assignee[0].displayName"},
					{ExternalId: "assignee"},
				},
			},
			want: map[string]struct{}{
				"assignee": {},
			},
		},
		{
			name:   "entity_with_child_entities",
			prefix: "",
			entity: &framework.EntityConfig{
				ExternalId: "Issue",
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.id"},
					{ExternalId: "$.fields.summary"},
				},
				ChildEntities: []*framework.EntityConfig{
					{
						ExternalId: "$.fields.issuetype",
						Attributes: []*framework.AttributeConfig{
							{ExternalId: "id"},
							{ExternalId: "name"},
						},
					},
				},
			},
			want: map[string]struct{}{
				"id":        {},
				"summary":   {},
				"issuetype": {},
			},
		},
		{
			name:   "entity_with_multiple_child_entities",
			prefix: "",
			entity: &framework.EntityConfig{
				ExternalId: "Issue",
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.id"},
					{ExternalId: "$.fields.summary"},
				},
				ChildEntities: []*framework.EntityConfig{
					{
						ExternalId: "$.fields.issuetype",
						Attributes: []*framework.AttributeConfig{
							{ExternalId: "id"},
							{ExternalId: "name"},
						},
					},
					{
						ExternalId: "$.fields.assignee",
						Attributes: []*framework.AttributeConfig{
							{ExternalId: "key"},
							{ExternalId: "displayName"},
						},
					},
				},
			},
			want: map[string]struct{}{
				"id":        {},
				"summary":   {},
				"issuetype": {},
				"assignee":  {},
			},
		},
		{
			name:   "entity_with_prefix",
			prefix: "project",
			entity: &framework.EntityConfig{
				ExternalId: "Issue",
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "id"},
					{ExternalId: "$.fields.summary"},
				},
			},
			want: map[string]struct{}{
				"project.id":               {},
				"project.$.fields.summary": {},
			},
		},
		{
			name:   "entity_with_array_indices_and_nested_paths",
			prefix: "",
			entity: &framework.EntityConfig{
				ExternalId: "Issue",
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "$.id"},
					{ExternalId: "$.fields.customfield_10209[0].value"},
					{ExternalId: "$.fields.assignee[0].key"},
					{ExternalId: "$.fields.components[0].name"},
					{ExternalId: "$.fields.fixVersions[0].id"},
				},
			},
			want: map[string]struct{}{
				"id":                {},
				"customfield_10209": {},
				"assignee":          {},
				"components":        {},
				"fixVersions":       {},
			},
		},
		{
			name:   "entity_with_empty_attribute_external_ids",
			prefix: "",
			entity: &framework.EntityConfig{
				ExternalId: "Issue",
				Attributes: []*framework.AttributeConfig{
					{ExternalId: ""},
					{ExternalId: "$.id"},
					{ExternalId: "$.fields.summary"},
				},
			},
			want: map[string]struct{}{
				"id":      {},
				"summary": {},
			},
		},
		{
			name:   "entity_with_no_attributes_only_child_entities",
			prefix: "",
			entity: &framework.EntityConfig{
				ExternalId: "Issue",
				Attributes: []*framework.AttributeConfig{},
				ChildEntities: []*framework.EntityConfig{
					{
						ExternalId: "$.fields.issuetype",
						Attributes: []*framework.AttributeConfig{
							{ExternalId: "id"},
						},
					},
				},
			},
			want: map[string]struct{}{
				"issuetype": {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := jiradatacenter.ExtractEntityFieldNames(tt.prefix, tt.entity)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractEntityFieldNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
