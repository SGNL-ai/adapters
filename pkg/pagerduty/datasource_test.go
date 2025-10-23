// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst
package pagerduty_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagerduty"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

// Define the endpoints and responses for the mock PagerDuty server.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.URL.RequestURI() {
	// User endpoints
	case "/users?offset=0&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"users": [{"id": "user1"}], "more": true}`))
	case "/users?offset=1&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"users": [{"id": "user2"}], "more": true}`))
	case "/users?offset=2&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"users": [{"id": "user3"}], "more": false}`))
	case "/users?offset=0&limit=1&team_ids%5B%5D=PJ5B6SN&team_ids%5B%5D=1234&query=user4":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"users": [{"id": "user4"}], "more": false}`))
	// Add duplicate endpoint because query parameters are not guaranteed to be in the same order
	// due to Go map iteration.
	case "/users?offset=0&limit=1&query=user4&team_ids%5B%5D=PJ5B6SN&team_ids%5B%5D=1234":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"users": [{"id": "user4"}], "more": false}`))

	// Team endpoints
	case "/teams?offset=0&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"teams": [{"id": "team1"}], "more": true}`))
	case "/teams?offset=1&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"teams": [{"id": "team2"}], "more": false}`))

	// Member endpoints
	case "/teams/team1/members?offset=0&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"members": [{"user": {"id": "user1"}}], "more": true}`))
	case "/teams/team1/members?offset=1&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"members": [{"user": {"id": "user2"}}], "more": false}`))
	case "/teams/team2/members?offset=0&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"members": [{"user": {"id": "user1"}}], "more": true}`))
	case "/teams/team2/members?offset=1&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"members": [{"user": {"id": "user3"}}], "more": false}`))

	case "/oncalls?offset=0&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"oncalls": [{"user": {"id": "user1"}, "escalation_policy": {"id": "policy1"}}], "more": false}`))
	case "/oncalls?offset=99&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"oncalls": [{"user": {"id": "user2"}, "escalation_policy": {"id": "policy2"},` +
			`"start": "2015-03-06T15:28:51-05:00", "end": "2015-03-07T15:28:51-05:00"}], "more": false}`))
	case "/oncalls?offset=100&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"oncalls": [{"user": {"id": "user2"}, "escalation_policy": "NOT_A_MAP"}], "more": false}`))
	case "/oncalls?offset=101&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"oncalls": [{"user": {"id": "user2"}, "escalation_policy": {"id": 1234}}], "more": false}`))
	case "/oncalls?offset=102&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"oncalls": [{"user": "NOT_A_MAP", "escalation_policy": {"id": "policy1"}}], "more": false}`))
	case "/oncalls?offset=103&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"oncalls": [{"user": {"id": 1234}, "escalation_policy": {"id": "policy1"}}], "more": false}`))
	case "/oncalls?offset=104&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"oncalls": [{"user": {"id": "user2"}, "escalation_policy": {"id": "policy2"},` +
			`"start": 1234, "end": "2015-03-07T15:28:51-05:00"}], "more": false}`))
	case "/oncalls?offset=105&limit=1":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"oncalls": [{"user": {"id": "user2"}, "escalation_policy": {"id": "policy2"},` +
			`"start": "2015-03-06T15:28:51-05:00", "end": 1234}], "more": false}`))

	// HTML error response for testing non-JSON error bodies
	case "/html_error?offset=0&limit=1":
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`<html><body><h1>500 Internal Server Error</h1></body></html>`))

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":{"message":"Not Found","code":404}}`))
	}
})

func TestParseResponse(t *testing.T) {
	tests := map[string]struct {
		body             []byte
		entityExternalID string
		pageSize         int64
		cursor           int64
		wantObjects      []map[string]interface{}
		wantNextCursor   *int64
		wantErr          *framework.Error
	}{
		"single_page_no_next_cursor": {
			body:             []byte(`{"users": [{"name": "user1"}, {"name": "user2"}], "more": false}`),
			entityExternalID: pagerduty.Users,
			pageSize:         10,
			cursor:           0,
			wantObjects: []map[string]interface{}{
				{"name": "user1"},
				{"name": "user2"},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"first_page_with_next_cursor": {
			body:             []byte(`{"users": [{"name": "user1"}, {"name": "user2"}], "more": true}`),
			entityExternalID: pagerduty.Users,
			pageSize:         2,
			cursor:           0,
			wantObjects: []map[string]interface{}{
				{"name": "user1"},
				{"name": "user2"},
			},
			wantNextCursor: testutil.GenPtr[int64](2), // This page contains index 0 and 1, so next page starts at index 2.
			wantErr:        nil,
		},
		// If the "more" field is missing, we compute the next cursor regardless.
		"missing_more_field": {
			body:             []byte(`{"users": [{"name": "user1"}, {"name": "user2"}]}`),
			entityExternalID: pagerduty.Users,
			pageSize:         2,
			cursor:           0,
			wantObjects: []map[string]interface{}{
				{"name": "user1"},
				{"name": "user2"},
			},
			wantNextCursor: testutil.GenPtr[int64](2),
			wantErr:        nil,
		},
		"more_field_not_bool": {
			body:             []byte(`{"users": [{"name": "user1"}, {"name": "user2"}], "more": "NOT_A_BOOL"}`),
			entityExternalID: pagerduty.Users,
			pageSize:         2,
			cursor:           0,
			wantObjects:      nil,
			wantNextCursor:   nil,
			wantErr: &framework.Error{
				Message: "Field more exists in PagerDuty response but field value is not a bool: string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"data_not_list": {
			// "users" field should be a list of objects.
			body:             []byte(`{"users": {"name": "user2"}, "more": true}`),
			entityExternalID: pagerduty.Users,
			pageSize:         2,
			cursor:           0,
			wantObjects:      nil,
			wantNextCursor:   nil,
			wantErr: &framework.Error{
				Message: `Entity users field exists in PagerDuty response ` +
					`but field value is not a list of objects: map[string]interface {}.`,
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"missing_entity_external_id_field": {
			// The "users" field should exist.
			body:             []byte(`{"not_users": {"name": "user2"}, "more": true}`),
			entityExternalID: pagerduty.Users,
			pageSize:         2,
			cursor:           0,
			wantObjects:      nil,
			wantNextCursor:   nil,
			wantErr: &framework.Error{
				Message: "Field missing in PagerDuty response: users.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"unmarshal_error": {
			body:             []byte(`{"not_users": {"name": "user2"`),
			entityExternalID: pagerduty.Users,
			pageSize:         2,
			cursor:           0,
			wantObjects:      nil,
			wantNextCursor:   nil,
			wantErr: &framework.Error{
				Message: "Failed to unmarshal PagerDuty response: unexpected end of JSON input.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"object_not_map": {
			body:             []byte(`{"users": [{"name": "user1"}, "STRING"], "more": true}`),
			entityExternalID: pagerduty.Users,
			pageSize:         2,
			cursor:           0,
			wantObjects:      nil,
			wantNextCursor:   nil,
			wantErr: &framework.Error{
				Message: "An object in Entity: users could not be parsed. Expected: map[string]any. Got: string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := pagerduty.ParseResponse(tt.body, tt.entityExternalID, tt.pageSize, tt.cursor)

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

func TestGetPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	pagerdutyClient := pagerduty.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context      context.Context
		request      *pagerduty.Request
		wantRes      *pagerduty.Response
		wantErr      *framework.Error
		expectedLogs []map[string]any
	}{
		"first_page": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.Users,
				PageSize:              1,
			},
			wantRes: &pagerduty.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{"id": "user1"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](1),
				},
			},
			wantErr: nil,
			expectedLogs: []map[string]any{
				{
					"level":                             "info",
					"msg":                               "Starting datasource request",
					fields.FieldRequestEntityExternalID: pagerduty.Users,
					fields.FieldRequestPageSize:         int64(1),
				},
				{
					"level":                             "info",
					"msg":                               "Sending HTTP request to datasource",
					fields.FieldRequestEntityExternalID: pagerduty.Users,
					fields.FieldRequestPageSize:         int64(1),
					fields.FieldURL:                     server.URL + "/users?offset=0&limit=1",
				},
				{
					"level":                             "info",
					"msg":                               "Datasource request completed successfully",
					fields.FieldRequestEntityExternalID: pagerduty.Users,
					fields.FieldRequestPageSize:         int64(1),
					fields.FieldResponseStatusCode:      int64(200),
					fields.FieldResponseObjectCount:     int64(1),
					fields.FieldResponseNextCursor: map[string]any{
						"cursor": int64(1),
					},
				},
			},
		},
		"middle_page": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.Users,
				PageSize:              1,
				Cursor:                &pagination.CompositeCursor[int64]{Cursor: testutil.GenPtr[int64](1)},
			},
			wantRes: &pagerduty.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{"id": "user2"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](2),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.Users,
				PageSize:              1,
				Cursor:                &pagination.CompositeCursor[int64]{Cursor: testutil.GenPtr[int64](2)},
			},
			wantRes: &pagerduty.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{"id": "user3"},
				},
			},
			wantErr: nil,
		},
		"http_not_found_error": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      "invalid_entity",
				PageSize:              1,
			},
			wantRes: &pagerduty.Response{
				StatusCode:       http.StatusNotFound,
				RetryAfterHeader: "",
			},
			wantErr: nil,
			expectedLogs: []map[string]any{
				{
					"level":                             "info",
					"msg":                               "Starting datasource request",
					fields.FieldRequestEntityExternalID: "invalid_entity",
					fields.FieldRequestPageSize:         int64(1),
				},
				{
					"level":                             "info",
					"msg":                               "Sending HTTP request to datasource",
					fields.FieldRequestEntityExternalID: "invalid_entity",
					fields.FieldRequestPageSize:         int64(1),
					fields.FieldURL:                     server.URL + "/invalid_entity?offset=0&limit=1",
				},
				{
					"level":                              "error",
					"msg":                                "Datasource request failed",
					fields.FieldRequestEntityExternalID:  "invalid_entity",
					fields.FieldRequestPageSize:          int64(1),
					fields.FieldResponseStatusCode:       int64(404),
					fields.FieldResponseRetryAfterHeader: "",
					fields.FieldResponseBody: map[string]any{
						"error": map[string]any{
							"message": "Not Found",
							"code":    float64(404),
						},
					},
					fields.FieldSGNLEventType: fields.SgnlEventTypeErrorValue,
				},
			},
		},
		"http_internal_server_error_with_html_body": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      "html_error",
				PageSize:              1,
			},
			wantRes: &pagerduty.Response{
				StatusCode:       http.StatusInternalServerError,
				RetryAfterHeader: "",
			},
			wantErr: nil,
			expectedLogs: []map[string]any{
				{
					"level":                             "info",
					"msg":                               "Starting datasource request",
					fields.FieldRequestEntityExternalID: "html_error",
					fields.FieldRequestPageSize:         int64(1),
				},
				{
					"level":                             "info",
					"msg":                               "Sending HTTP request to datasource",
					fields.FieldRequestEntityExternalID: "html_error",
					fields.FieldRequestPageSize:         int64(1),
					fields.FieldURL:                     server.URL + "/html_error?offset=0&limit=1",
				},
				{
					"level":                              "error",
					"msg":                                "Datasource request failed",
					fields.FieldRequestEntityExternalID:  "html_error",
					fields.FieldRequestPageSize:          int64(1),
					fields.FieldResponseStatusCode:       int64(500),
					fields.FieldResponseRetryAfterHeader: "",
					fields.FieldResponseBody:             "<html><body><h1>500 Internal Server Error</h1></body></html>",
					fields.FieldSGNLEventType:            fields.SgnlEventTypeErrorValue,
				},
			},
		},
		"first_member_page_first_group": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.Members,
				PageSize:              1,
			},
			wantRes: &pagerduty.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{"id": "team1-user1", "userId": "user1", "teamId": "team1"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr[int64](1),
					CollectionID:     testutil.GenPtr("team1"),
					CollectionCursor: testutil.GenPtr[int64](1),
				},
			},
			wantErr: nil,
		},
		"last_member_page_first_group": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.Members,
				PageSize:              1,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr[int64](1),
					CollectionID:     testutil.GenPtr("team1"),
					CollectionCursor: testutil.GenPtr[int64](1),
				},
			},
			wantRes: &pagerduty.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{"id": "team1-user2", "userId": "user2", "teamId": "team1"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					CollectionID:     testutil.GenPtr("team1"),
					CollectionCursor: testutil.GenPtr[int64](1),
				},
			},
			wantErr: nil,
		},
		"first_member_page_last_group": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.Members,
				PageSize:              1,
				// This cursor should match the cursor returned in the last test case.
				Cursor: &pagination.CompositeCursor[int64]{
					CollectionID:     testutil.GenPtr("team1"),
					CollectionCursor: testutil.GenPtr[int64](1),
				},
			},
			wantRes: &pagerduty.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{"id": "team2-user1", "userId": "user1", "teamId": "team2"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor:       testutil.GenPtr[int64](1),
					CollectionID: testutil.GenPtr("team2"),
				},
			},
			wantErr: nil,
		},
		"last_member_page_last_group": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.Members,
				PageSize:              1,
				// This cursor should match the cursor returned in the last test case.
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:       testutil.GenPtr[int64](1),
					CollectionID: testutil.GenPtr("team2"),
				},
			},
			wantRes: &pagerduty.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{"id": "team2-user3", "userId": "user3", "teamId": "team2"},
				},
			},
			wantErr: nil,
		},
		"with_additional_query_parameters": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.Users,
				AdditionalQueryParameters: map[string]map[string][]string{
					// We define our test server handler to return a single user with id "user4" when
					// these query parameters are passed in.
					"users": {
						"team_ids[]": {"PJ5B6SN", "1234"},
						"query":      {"user4"},
					},
				},
				PageSize: 1,
			},
			wantRes: &pagerduty.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{"id": "user4"},
				},
			},
			wantErr: nil,
		},
		"oncall_objects_have_unique_id_created": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.OnCalls,
				PageSize:              1,
			},
			wantRes: &pagerduty.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					0: {
						"escalation_policy": map[string]interface{}{
							"id": "policy1",
						},
						"user": map[string]interface{}{
							"id": "user1",
						},
						// The format of the unique "id" is {policyID}-{userID}-{start}-{end}.
						// Since start and end don't exist in the response, they're converted to empty strings.
						"id": "policy1-user1--",
					},
				},
			},
			wantErr: nil,
		},
		"oncall_objects_have_unique_id_created_with_dates": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.OnCalls,
				PageSize:              1,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](99),
				},
			},
			wantRes: &pagerduty.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					0: {
						"escalation_policy": map[string]interface{}{
							"id": "policy2",
						},
						"user": map[string]interface{}{
							"id": "user2",
						},
						// The format of the unique "id" is {policyID}-{userID}-{start}-{end}.
						// Since start and end don't exist in the response, they're converted to empty strings.
						"id":    "policy2-user2-2015-03-06T15:28:51-05:00-2015-03-07T15:28:51-05:00",
						"start": "2015-03-06T15:28:51-05:00",
						"end":   "2015-03-07T15:28:51-05:00",
					},
				},
			},
			wantErr: nil,
		},
		"oncall_object_escalation_policy_field_not_a_map": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.OnCalls,
				PageSize:              1,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](100),
				},
			},
			wantRes: nil,
			wantErr: &framework.Error{
				Message: "Failed to parse a PagerDuty OnCall object's escalation_policy field as map[string]any.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"oncall_object_escalation_policy_field_id_attribute_not_string": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.OnCalls,
				PageSize:              1,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](101),
				},
			},
			wantRes: nil,
			wantErr: &framework.Error{
				Message: "Failed to parse a field in a PagerDuty OnCall object's escalation_policy object as string: id.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"oncall_object_user_field_not_map": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.OnCalls,
				PageSize:              1,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](102),
				},
			},
			wantRes: nil,
			wantErr: &framework.Error{
				Message: "Failed to parse a PagerDuty OnCall object's user field as map[string]any.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"oncall_object_user_field_id_attribute_not_string": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.OnCalls,
				PageSize:              1,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](103),
				},
			},
			wantRes: nil,
			wantErr: &framework.Error{
				Message: "Failed to parse a field in a PagerDuty OnCall object's user object as string: id.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"oncall_object_start_field_not_date_string": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.OnCalls,
				PageSize:              1,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](104),
				},
			},
			wantRes: nil,
			wantErr: &framework.Error{
				Message: "Failed to parse a PagerDuty OnCall object's start field as string: 1234.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"oncall_object_end_field_not_date_string": {
			context: context.Background(),
			request: &pagerduty.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Token token=1234",
				EntityExternalID:      pagerduty.OnCalls,
				PageSize:              1,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](105),
				},
			},
			wantRes: nil,
			wantErr: &framework.Error{
				Message: "Failed to parse a PagerDuty OnCall object's end field as string: 1234.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctxWithLogger, observedLogs := testutil.NewContextWithObservableLogger(tt.context)

			gotRes, gotErr := pagerdutyClient.GetPage(ctxWithLogger, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}

			testutil.ValidateLogOutput(t, observedLogs, tt.expectedLogs)
		})
	}
}
