// Copyright 2026 SGNL.ai, Inc.

// nolint: goconst

package aws_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	aws_adapter "github.com/sgnl-ai/adapters/pkg/aws"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestAdapterGetPage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(defaultTimeout)*time.Second)
	defer cancel()

	cfg, err := SetupTestConfig(ctx, Mocker)
	if err != nil {
		log.Fatalf("Failed to load aws test config: %v", err)
	}

	adapter, err := ProvideAWSTestClient(cfg)
	if err != nil {
		log.Fatalf("Failed to load aws test client: %v", err)
	}

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[aws_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn": "arn:aws:iam::000000000000:user/user1",
						},
						{
							"Arn": "arn:aws:iam::000000000000:user/user2",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiIyIn0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("2"),
			},
		},
		"valid_request_with_cursor": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
					},
				},
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("2"),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn": "arn:aws:iam::000000000000:user/user3",
						},
						{
							"Arn": "arn:aws:iam::000000000000:user/user4",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiI0In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("4"),
			},
		},
		"valid_request_with_cursor_last_page": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
					},
				},
				PageSize: 4,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("4"),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn": "arn:aws:iam::000000000000:user/user5",
						},
						{
							"Arn": "arn:aws:iam::000000000000:user/user6",
						},
					},
				},
			},
		},
		"valid_empty_entity_config": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Config: &aws_adapter.Config{
					Region: "us-west-2",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn": "arn:aws:iam::000000000000:user/user1",
						},
						{
							"Arn": "arn:aws:iam::000000000000:user/user2",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiIyIn0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("2"),
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
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(defaultTimeout)*time.Second)
	defer cancel()

	cfg, err := SetupTestConfig(ctx, Mocker)
	if err != nil {
		log.Fatalf("Failed to load aws test config: %v", err)
	}

	adapter, err := ProvideAWSTestClient(cfg)
	if err != nil {
		log.Fatalf("Failed to load aws test client: %v", err)
	}

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[aws_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
		expectedLogs       []map[string]any
	}{
		"valid_request_page_1": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "PasswordLastUsed",
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
							"Arn":              "arn:aws:iam::000000000000:user/user1",
							"AccountId":        "000000000000",
							"UserName":         "user1",
							"Path":             "/",
							"UserId":           "user1",
							"CreateDate":       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							"PasswordLastUsed": time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":              "arn:aws:iam::000000000000:user/user2",
							"AccountId":        "000000000000",
							"UserName":         "user2",
							"Path":             "/",
							"UserId":           "user2",
							"CreateDate":       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							"PasswordLastUsed": time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiIyIn0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("2"),
			},
			expectedLogs: []map[string]any{
				{
					"level":                             "info",
					"msg":                               "Starting datasource request",
					fields.FieldRequestEntityExternalID: "User",
					fields.FieldRequestPageSize:         int64(2),
				},
				{
					"level":                             "info",
					"msg":                               "Datasource request completed successfully",
					fields.FieldRequestEntityExternalID: "User",
					fields.FieldRequestPageSize:         int64(2),
					fields.FieldResponseStatusCode:      int64(200),
					fields.FieldResponseObjectCount:     int64(2),
					fields.FieldResponseNextCursor: map[string]any{
						"cursor": "2",
					},
				},
			},
		},
		"valid_request_page_2": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "PasswordLastUsed",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOiIyIn0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":              "arn:aws:iam::000000000000:user/user3",
							"AccountId":        "000000000000",
							"UserName":         "user3",
							"Path":             "/",
							"UserId":           "user3",
							"CreateDate":       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							"PasswordLastUsed": time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":              "arn:aws:iam::000000000000:user/user4",
							"AccountId":        "000000000000",
							"UserName":         "user4",
							"Path":             "/",
							"UserId":           "user4",
							"CreateDate":       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							"PasswordLastUsed": time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiI0In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("4"),
			},
		},
		"valid_request_sdk_error": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Config: &aws_adapter.Config{
					Region: "us-west-2",
					EntityConfig: map[string]*aws_adapter.EntityConfig{
						"User": {
							PathPrefix: func() *string {
								s := internalError

								return &s

							}(),
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Unable to fetch AWS entity: User, error: Failed to list entities: " +
						"operation error IAM: ListUsers, InternalFailure.",
					Code: v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"valid_request_not_found_with_filter": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Config: &aws_adapter.Config{
					Region: "us-west-2",
					EntityConfig: map[string]*aws_adapter.EntityConfig{
						"User": {
							PathPrefix: func() *string {
								s := "/not-found"

								return &s

							}(),
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{},
			},
		},
		"valid_request_with_multiple_accounts_fetch_first_page_first_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "PasswordLastUsed",
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
							"Arn":              "arn:aws:iam::000000000000:user/user1",
							"AccountId":        "000000000000",
							"UserName":         "user1",
							"Path":             "/",
							"UserId":           "user1",
							"CreateDate":       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							"PasswordLastUsed": time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":              "arn:aws:iam::000000000000:user/user2",
							"AccountId":        "000000000000",
							"UserName":         "user2",
							"Path":             "/",
							"UserId":           "user2",
							"CreateDate":       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							"PasswordLastUsed": time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					// Base64 encoded value of {"cursor":"eyJPZmZzZXQiOjAsIk5leHRNYXJrZXIiOiIyIn0="}
					NextCursor: "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakFzSWs1bGVIUk5ZWEpyWlhJaU9pSXlJbjA9In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				// Base64 encoded value of {"Offset":0,"NextMarker":"2"}
				Cursor: testutil.GenPtr("eyJPZmZzZXQiOjAsIk5leHRNYXJrZXIiOiIyIn0="),
			},
		},
		"valid_request_with_multiple_accounts_fetch_second_page_first_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "PasswordLastUsed",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakFzSWs1bGVIUk5ZWEpyWlhJaU9pSXlJbjA9In0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":              "arn:aws:iam::000000000000:user/user3",
							"AccountId":        "000000000000",
							"UserName":         "user3",
							"Path":             "/",
							"UserId":           "user3",
							"CreateDate":       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							"PasswordLastUsed": time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":              "arn:aws:iam::000000000000:user/user4",
							"AccountId":        "000000000000",
							"UserName":         "user4",
							"Path":             "/",
							"UserId":           "user4",
							"CreateDate":       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							"PasswordLastUsed": time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					// Base64 encoded value of {"cursor":"eyJPZmZzZXQiOjAsIk5leHRNYXJrZXIiOiI0In0="}
					NextCursor: "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakFzSWs1bGVIUk5ZWEpyWlhJaU9pSTBJbjA9In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				// Base64 encoded value of {"Offset":0,"NextMarker":"4"}
				Cursor: testutil.GenPtr("eyJPZmZzZXQiOjAsIk5leHRNYXJrZXIiOiI0In0="),
			},
		},
		"valid_request_with_multiple_accounts_fetch_last_page_first_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "PasswordLastUsed",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakFzSWs1bGVIUk5ZWEpyWlhJaU9pSTBJbjA9In0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":              "arn:aws:iam::000000000000:user/user5",
							"AccountId":        "000000000000",
							"UserName":         "user5",
							"Path":             "/",
							"UserId":           "user5",
							"CreateDate":       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							"PasswordLastUsed": time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":              "arn:aws:iam::000000000000:user/user6",
							"AccountId":        "000000000000",
							"UserName":         "user6",
							"Path":             "/",
							"UserId":           "user6",
							"CreateDate":       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							"PasswordLastUsed": time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					// Base64 encoded value of {"cursor":"eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOm51bGx9"}
					NextCursor: "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakVzSWs1bGVIUk5ZWEpyWlhJaU9tNTFiR3g5In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				// Base64 encoded value of {"Offset":1,"NextMarker":null}
				Cursor: testutil.GenPtr("eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOm51bGx9"),
			},
		},
		"valid_request_with_multiple_accounts_fetch_first_page_second_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "PasswordLastUsed",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakVzSWs1bGVIUk5ZWEpyWlhJaU9tNTFiR3g5In0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":              "arn:aws:iam::000000000000:user/user1",
							"AccountId":        "000000000000",
							"UserName":         "user1",
							"Path":             "/",
							"UserId":           "user1",
							"CreateDate":       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							"PasswordLastUsed": time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":              "arn:aws:iam::000000000000:user/user2",
							"AccountId":        "000000000000",
							"UserName":         "user2",
							"Path":             "/",
							"UserId":           "user2",
							"CreateDate":       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							"PasswordLastUsed": time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					// Base64 encoded value of {"cursor":"eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOiIyIn0="}
					NextCursor: "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakVzSWs1bGVIUk5ZWEpyWlhJaU9pSXlJbjA9In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				// Base64 encoded value of {"Offset":1,"NextMarker":"2"}
				Cursor: testutil.GenPtr("eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOiIyIn0="),
			},
		},
		"valid_request_with_multiple_accounts_fetch_last_page_second_account_no_more_accounts": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "PasswordLastUsed",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakVzSWs1bGVIUk5ZWEpyWlhJaU9pSTBJbjA9In0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":              "arn:aws:iam::000000000000:user/user5",
							"AccountId":        "000000000000",
							"UserName":         "user5",
							"Path":             "/",
							"UserId":           "user5",
							"CreateDate":       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							"PasswordLastUsed": time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":              "arn:aws:iam::000000000000:user/user6",
							"AccountId":        "000000000000",
							"UserName":         "user6",
							"Path":             "/",
							"UserId":           "user6",
							"CreateDate":       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							"PasswordLastUsed": time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
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

			ctxWithLogger, observedLogs := testutil.NewContextWithObservableLogger(tt.ctx)

			gotResponse := adapter.GetPage(ctxWithLogger, tt.request)

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

			testutil.ValidateLogOutput(t, observedLogs, tt.expectedLogs)
		})
	}
}

func TestAdapterGetGroupPage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(defaultTimeout)*time.Second)
	defer cancel()

	cfg, err := SetupTestConfig(ctx, Mocker)
	if err != nil {
		log.Fatalf("Failed to load aws test config: %v", err)
	}

	adapter, err := ProvideAWSTestClient(cfg)
	if err != nil {
		log.Fatalf("Failed to load aws test client: %v", err)
	}

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[aws_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
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
							"Arn":        "arn:aws:iam::000000000000:group/Group1",
							"AccountId":  "000000000000",
							"GroupName":  "Group1",
							"Path":       "/",
							"GroupId":    "Group1",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":        "arn:aws:iam::000000000000:group/Group2",
							"AccountId":  "000000000000",
							"GroupName":  "Group2",
							"Path":       "/",
							"GroupId":    "Group2",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiIyIn0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("2"),
			},
		},
		"valid_request_sdk_error": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Config: &aws_adapter.Config{
					Region: "us-west-2",
					EntityConfig: map[string]*aws_adapter.EntityConfig{
						"Group": {
							PathPrefix: func() *string {
								s := internalError

								return &s

							}(),
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Unable to fetch AWS entity: Group, error: Failed to list entities: " +
						"operation error IAM: ListGroups, InternalFailure.",
					Code: v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"valid_request_fetch_first_page_first_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
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
							"Arn":        "arn:aws:iam::000000000000:group/Group1",
							"AccountId":  "000000000000",
							"GroupName":  "Group1",
							"Path":       "/",
							"GroupId":    "Group1",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":        "arn:aws:iam::000000000000:group/Group2",
							"AccountId":  "000000000000",
							"GroupName":  "Group2",
							"Path":       "/",
							"GroupId":    "Group2",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					// Base64 encoded value for `{"cursor":"eyJPZmZzZXQiOjAsIk5leHRNYXJrZXIiOiIyIn0="}`
					NextCursor: "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakFzSWs1bGVIUk5ZWEpyWlhJaU9pSXlJbjA9In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				// Base64 encoded value of {"Offset":0,"NextMarker":"2"}
				Cursor: testutil.GenPtr("eyJPZmZzZXQiOjAsIk5leHRNYXJrZXIiOiIyIn0="),
			},
		},
		"valid_request_fetch_last_page_first_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakFzSWs1bGVIUk5ZWEpyWlhJaU9pSXlJbjA9In0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":        "arn:aws:iam::000000000000:group/Group3",
							"AccountId":  "000000000000",
							"GroupName":  "Group3",
							"Path":       "/",
							"GroupId":    "Group3",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":        "arn:aws:iam::000000000000:group/Group4",
							"AccountId":  "000000000000",
							"GroupName":  "Group4",
							"Path":       "/",
							"GroupId":    "Group4",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					// Base64 encoded value for `{"cursor":"eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOm51bGx9"}`
					NextCursor: "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakVzSWs1bGVIUk5ZWEpyWlhJaU9tNTFiR3g5In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				// Base64 encoded value of {"Offset":1,"NextMarker":null}
				Cursor: testutil.GenPtr("eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOm51bGx9"),
			},
		},
		"valid_request_fetch_last_page_second_account_no_more_accounts_left": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakVzSWs1bGVIUk5ZWEpyWlhJaU9pSXlJbjA9In0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":        "arn:aws:iam::000000000000:group/Group3",
							"AccountId":  "000000000000",
							"GroupName":  "Group3",
							"Path":       "/",
							"GroupId":    "Group3",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":        "arn:aws:iam::000000000000:group/Group4",
							"AccountId":  "000000000000",
							"GroupName":  "Group4",
							"Path":       "/",
							"GroupId":    "Group4",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
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
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(defaultTimeout)*time.Second)
	defer cancel()

	cfg, err := SetupTestConfig(ctx, Mocker)
	if err != nil {
		log.Fatalf("Failed to load aws test config: %v", err)
	}

	adapter, err := ProvideAWSTestClient(cfg)
	if err != nil {
		log.Fatalf("Failed to load aws test client: %v", err)
	}

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[aws_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "Role",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "RoleName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "RoleId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
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
							"Arn":        "arn:aws:iam::000000000000:role/sso.amazonaws.com/AWSReservedSSO_de25js739eef1832",
							"AccountId":  "000000000000",
							"RoleName":   "AWSReservedSSO_de25js739eef1832",
							"Path":       "/sso.amazonaws.com",
							"RoleId":     "AROAXXXXXXXXXXXXXXXX2",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":        "arn:aws:iam::000000000000:role/sso.amazonaws.com/AWSReservedSSO_abcdef1234567890",
							"AccountId":  "000000000000",
							"RoleName":   "AWSReservedSSO_abcdef1234567890",
							"Path":       "/sso.amazonaws.com",
							"RoleId":     "AROA3C7OBZZZCD433N4DQ",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiIyIn0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("2"),
			},
		},
		"valid_request_sdk_error": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Config: &aws_adapter.Config{
					Region: "us-west-2",
					EntityConfig: map[string]*aws_adapter.EntityConfig{
						"Role": {
							PathPrefix: func() *string {
								s := internalError

								return &s

							}(),
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Role",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "RoleName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "RoleId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Unable to fetch AWS entity: Role, error: Failed to list entities: " +
						"operation error IAM: ListRoles, InternalFailure.",
					Code: v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"valid_request_multiple_accounts_fetch_first_page_first_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "Role",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "RoleName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "RoleId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
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
							"Arn":        "arn:aws:iam::000000000000:role/sso.amazonaws.com/AWSReservedSSO_de25js739eef1832",
							"AccountId":  "000000000000",
							"RoleName":   "AWSReservedSSO_de25js739eef1832",
							"Path":       "/sso.amazonaws.com",
							"RoleId":     "AROAXXXXXXXXXXXXXXXX2",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":        "arn:aws:iam::000000000000:role/sso.amazonaws.com/AWSReservedSSO_abcdef1234567890",
							"AccountId":  "000000000000",
							"RoleName":   "AWSReservedSSO_abcdef1234567890",
							"Path":       "/sso.amazonaws.com",
							"RoleId":     "AROA3C7OBZZZCD433N4DQ",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					// Base64 encoded value for `{"cursor":"eyJPZmZzZXQiOjAsIk5leHRNYXJrZXIiOiIyIn0="}`
					NextCursor: "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakFzSWs1bGVIUk5ZWEpyWlhJaU9pSXlJbjA9In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				// Base64 encoded value of {"Offset":0,"NextMarker":"2"}
				Cursor: testutil.GenPtr("eyJPZmZzZXQiOjAsIk5leHRNYXJrZXIiOiIyIn0="),
			},
		},
		"valid_request_multiple_accounts_fetch_last_page_first_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "Role",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "RoleName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "RoleId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakFzSWs1bGVIUk5ZWEpyWlhJaU9pSXlJbjA9In0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":        "arn:aws:iam::000000000000:role/sso.amazonaws.com/role_3",
							"AccountId":  "000000000000",
							"RoleName":   "role_3",
							"Path":       "/sso.amazonaws.com",
							"RoleId":     "AROA3C7OBYYYCD433N4DQ",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					// Base64 encoded value for `{"cursor":"eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOm51bGx9"}`
					NextCursor: "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakVzSWs1bGVIUk5ZWEpyWlhJaU9tNTFiR3g5In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				// Base64 encoded value of {"Offset":1,"NextMarker":null}
				Cursor: testutil.GenPtr("eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOm51bGx9"),
			},
		},
		"valid_request_multiple_accounts_fetch_last_page_second_account_and_no_more_accounts": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "Role",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "RoleName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "RoleId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakVzSWs1bGVIUk5ZWEpyWlhJaU9pSXlJbjA9In0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":        "arn:aws:iam::000000000000:role/sso.amazonaws.com/role_3",
							"AccountId":  "000000000000",
							"RoleName":   "role_3",
							"Path":       "/sso.amazonaws.com",
							"RoleId":     "AROA3C7OBYYYCD433N4DQ",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
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

func TestAdapterGetPolicyPage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(defaultTimeout)*time.Second)
	defer cancel()

	cfg, err := SetupTestConfig(ctx, Mocker)
	if err != nil {
		log.Fatalf("Failed to load aws test config: %v", err)
	}

	adapter, err := ProvideAWSTestClient(cfg)
	if err != nil {
		log.Fatalf("Failed to load aws test client: %v", err)
	}

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[aws_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "Policy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "AttachmentCount",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "PolicyName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "DefaultVersionId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "PolicyId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "IsAttachable",
							Type:       framework.AttributeTypeBool,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "UpdateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "PermissionsBoundaryUsageCount",
							Type:       framework.AttributeTypeInt64,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":                           "arn:aws:iam::000000000000:policy/ExampleEngPolicy",
							"AttachmentCount":               int64(1),
							"IsAttachable":                  true,
							"CreateDate":                    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"UpdateDate":                    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"PermissionsBoundaryUsageCount": int64(0),
							"AccountId":                     "000000000000",
							"PolicyName":                    "ExampleEngPolicy",
							"PolicyId":                      "ANPA3C7OBZZZCD411N4DQ",
							"DefaultVersionId":              "v1",
							"Path":                          "/",
						},
						{
							"Arn":                           "arn:aws:iam::000000000000:policy/Policy2",
							"AccountId":                     "000000000000",
							"PolicyName":                    "Policy2",
							"PolicyId":                      "ANPA3C7OBZZZCD433N4DQ",
							"AttachmentCount":               int64(1),
							"DefaultVersionId":              "v1",
							"CreateDate":                    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"IsAttachable":                  true,
							"Path":                          "/",
							"UpdateDate":                    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"PermissionsBoundaryUsageCount": int64(0),
						},
					},
				},
			},
		},
		"valid_request_sdk_error": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Config: &aws_adapter.Config{
					Region: "us-west-2",
					EntityConfig: map[string]*aws_adapter.EntityConfig{
						"Policy": {
							PathPrefix: func() *string {
								s := internalError

								return &s

							}(),
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Policy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "AttachmentCount",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "PolicyName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "DefaultVersionId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "PolicyId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "IsAttachable",
							Type:       framework.AttributeTypeBool,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "UpdateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "PermissionsBoundaryUsageCount",
							Type:       framework.AttributeTypeInt64,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Unable to fetch AWS entity: Policy, error: Failed to list entities: " +
						"operation error IAM: ListPolicies, InternalFailure.",
					Code: v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"valid_request_config_has_resource_accounts_fetch_page_from_first_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "Policy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "AttachmentCount",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "PolicyName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "DefaultVersionId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "PolicyId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "IsAttachable",
							Type:       framework.AttributeTypeBool,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "UpdateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "PermissionsBoundaryUsageCount",
							Type:       framework.AttributeTypeInt64,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":                           "arn:aws:iam::000000000000:policy/ExampleEngPolicy",
							"AttachmentCount":               int64(1),
							"IsAttachable":                  true,
							"CreateDate":                    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"UpdateDate":                    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"PermissionsBoundaryUsageCount": int64(0),
							"AccountId":                     "000000000000",
							"PolicyName":                    "ExampleEngPolicy",
							"PolicyId":                      "ANPA3C7OBZZZCD411N4DQ",
							"DefaultVersionId":              "v1",
							"Path":                          "/",
						},
						{
							"Arn":                           "arn:aws:iam::000000000000:policy/Policy2",
							"AccountId":                     "000000000000",
							"PolicyName":                    "Policy2",
							"PolicyId":                      "ANPA3C7OBZZZCD433N4DQ",
							"AttachmentCount":               int64(1),
							"DefaultVersionId":              "v1",
							"CreateDate":                    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"IsAttachable":                  true,
							"Path":                          "/",
							"UpdateDate":                    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"PermissionsBoundaryUsageCount": int64(0),
						},
					},
					// Base64 encoded value for `{"cursor":"eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOm51bGx9"}`
					NextCursor: "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakVzSWs1bGVIUk5ZWEpyWlhJaU9tNTFiR3g5In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				// Base64 encoded value of {"Offset":1,"NextMarker":null}
				Cursor: testutil.GenPtr("eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOm51bGx9"),
			},
		},
		"valid_request_config_has_resource_accounts_fetch_from_2nd_account_no_more_accounts_after_that": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "Policy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "AttachmentCount",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "PolicyName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "DefaultVersionId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "PolicyId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "IsAttachable",
							Type:       framework.AttributeTypeBool,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "UpdateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "PermissionsBoundaryUsageCount",
							Type:       framework.AttributeTypeInt64,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakVzSWs1bGVIUk5ZWEpyWlhJaU9tNTFiR3g5In0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":                           "arn:aws:iam::000000000000:policy/ExampleEngPolicy",
							"AttachmentCount":               int64(1),
							"IsAttachable":                  true,
							"CreateDate":                    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"UpdateDate":                    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"PermissionsBoundaryUsageCount": int64(0),
							"AccountId":                     "000000000000",
							"PolicyName":                    "ExampleEngPolicy",
							"PolicyId":                      "ANPA3C7OBZZZCD411N4DQ",
							"DefaultVersionId":              "v1",
							"Path":                          "/",
						},
						{
							"Arn":                           "arn:aws:iam::000000000000:policy/Policy2",
							"AccountId":                     "000000000000",
							"PolicyName":                    "Policy2",
							"PolicyId":                      "ANPA3C7OBZZZCD433N4DQ",
							"AttachmentCount":               int64(1),
							"DefaultVersionId":              "v1",
							"CreateDate":                    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"IsAttachable":                  true,
							"Path":                          "/",
							"UpdateDate":                    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"PermissionsBoundaryUsageCount": int64(0),
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

func TestAdapterGetSAMLProvidersPage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(defaultTimeout)*time.Second)
	defer cancel()

	cfg, err := SetupTestConfig(ctx, Mocker)
	if err != nil {
		log.Fatalf("Failed to load aws test config: %v", err)
	}

	adapter, err := ProvideAWSTestClient(cfg)
	if err != nil {
		log.Fatalf("Failed to load aws test client: %v", err)
	}

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[aws_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request_page_1": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "IdentityProvider",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "ValidUntil",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":        "arn:aws:iam::123456789012:saml-provider/Provider1",
							"AccountId":  "123456789012",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"ValidUntil": time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOiIxIn0=",
				},
			},
		},
		"valid_request_page_2": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "IdentityProvider",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "ValidUntil",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 1,
				Cursor:   "eyJjdXJzb3IiOiIxIn0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":        "arn:aws:iam::123456789012:saml-provider/Provider2",
							"AccountId":  "123456789012",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"ValidUntil": time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
						},
					},
				},
			},
		},
		"valid_request_all": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "IdentityProvider",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "ValidUntil",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 10,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":        "arn:aws:iam::123456789012:saml-provider/Provider1",
							"AccountId":  "123456789012",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"ValidUntil": time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":        "arn:aws:iam::123456789012:saml-provider/Provider2",
							"AccountId":  "123456789012",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"ValidUntil": time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
						},
					},
				},
			},
		},
		"invalid_request_with_filter": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Config: &aws_adapter.Config{
					Region: "us-west-2",
					EntityConfig: map[string]*aws_adapter.EntityConfig{
						"IdentityProvider": {
							PathPrefix: testutil.GenPtr("/some-filter"),
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "IdentityProvider",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "ValidUntil",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Entity IdentityProvider does not supports filtering.",
					Code:    v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				},
			},
		},
		"valid_request_with_multiple_accounts_fetch_all_from_first_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "IdentityProvider",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "ValidUntil",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 10,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":        "arn:aws:iam::123456789012:saml-provider/Provider1",
							"AccountId":  "123456789012",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"ValidUntil": time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":        "arn:aws:iam::123456789012:saml-provider/Provider2",
							"AccountId":  "123456789012",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"ValidUntil": time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					// Base64 encoded value for `{"cursor":"eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOm51bGx9"}`
					NextCursor: "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakVzSWs1bGVIUk5ZWEpyWlhJaU9tNTFiR3g5In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				// Base64 encoded value of {"Offset":1,"NextMarker":null}
				Cursor: testutil.GenPtr("eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOm51bGx9"),
			},
		},
		"valid_request_with_multiple_accounts_fetch_all_from_second_account_no_more_accounts_left": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "IdentityProvider",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "ValidUntil",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 10,
				Cursor:   "eyJjdXJzb3IiOiJleUpQWm1aelpYUWlPakVzSWs1bGVIUk5ZWEpyWlhJaU9tNTFiR3g5In0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Arn":        "arn:aws:iam::123456789012:saml-provider/Provider1",
							"AccountId":  "123456789012",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"ValidUntil": time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							"Arn":        "arn:aws:iam::123456789012:saml-provider/Provider2",
							"AccountId":  "123456789012",
							"CreateDate": time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
							"ValidUntil": time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
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

func TestAdapterGetSAMLProvidersEmptyPage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(defaultTimeout)*time.Second)
	defer cancel()

	cfg, err := SetupTestConfig(ctx, EmptyMocker)
	if err != nil {
		log.Fatalf("Failed to load aws test config: %v", err)
	}

	adapter, err := ProvideAWSTestClient(cfg)
	if err != nil {
		log.Fatalf("Failed to load aws test client: %v", err)
	}

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[aws_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request_empty_page": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "IdentityProvider",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "ValidUntil",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{},
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

func TestAdapterGetGroupPolicyPage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(defaultTimeout)*time.Second)
	defer cancel()

	cfg, err := SetupTestConfig(ctx, Mocker)
	if err != nil {
		log.Fatalf("Failed to load aws test config: %v", err)
	}

	adapter, err := ProvideAWSTestClient(cfg)
	if err != nil {
		log.Fatalf("Failed to load aws test client: %v", err)
	}

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[aws_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_group_1_policy_1_of_2": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "GroupPolicy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "PolicyArn",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "arn:aws:iam::000000000000:policy/ExampleEngPolicy-Group1",
							"PolicyArn": "arn:aws:iam::000000000000:policy/ExampleEngPolicy",
							"GroupName": "Group1",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiIxIiwiY29sbGVjdGlvbklkIjoiR3JvdXAxIiwiY29sbGVjdGlvbkN1cnNvciI6IjEifQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID:     testutil.GenPtr("Group1"),
				CollectionCursor: testutil.GenPtr("1"),
				Cursor:           testutil.GenPtr("1"),
			},
		},
		"valid_group_1_policy_2_of_2": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "GroupPolicy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "PolicyArn",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 1,
				Cursor:   "eyJjdXJzb3IiOiIxIiwiY29sbGVjdGlvbklkIjoiR3JvdXAxIiwiY29sbGVjdGlvbkN1cnNvciI6IjEifQ==",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "arn:aws:iam::000000000000:policy/Policy2-Group1",
							"PolicyArn": "arn:aws:iam::000000000000:policy/Policy2",
							"GroupName": "Group1",
						},
					},
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDEiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiMSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID:     testutil.GenPtr("Group1"),
				CollectionCursor: testutil.GenPtr("1"),
			},
		},
		"valid_group_2_policy_0_of_0": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "GroupPolicy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "PolicyArn",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 1,
				Cursor:   "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDEiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiMSJ9",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDIiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiMiJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID:     testutil.GenPtr("Group2"),
				CollectionCursor: testutil.GenPtr("2"),
			},
		},
		"valid_group_4_policy_1_of_1": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "GroupPolicy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "PolicyArn",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 1,
				Cursor:   "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDMiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiMyJ9",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "arn:aws:iam::000000000000:policy/ExampleEngPolicy-Group4",
							"PolicyArn": "arn:aws:iam::000000000000:policy/ExampleEngPolicy",
							"GroupName": "Group4",
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

func TestAdapterGetUserPolicyPage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(defaultTimeout)*time.Second)
	defer cancel()

	cfg, err := SetupTestConfig(ctx, Mocker)
	if err != nil {
		log.Fatalf("Failed to load aws test config: %v", err)
	}

	adapter, err := ProvideAWSTestClient(cfg)
	if err != nil {
		log.Fatalf("Failed to load aws test client: %v", err)
	}

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[aws_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_user_1_policy_1_of_1": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "UserPolicy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "PolicyArn",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "arn:aws:iam::000000000000:policy/Policy104-user1",
							"PolicyArn": "arn:aws:iam::000000000000:policy/Policy104",
							"UserName":  "user1",
						},
					},
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJ1c2VyMSIsImNvbGxlY3Rpb25DdXJzb3IiOiIxIn0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID:     testutil.GenPtr("user1"),
				CollectionCursor: testutil.GenPtr("1"),
			},
		},
		"valid_user_2_policy_0_of_0": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "UserPolicy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "PolicyArn",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 1,
				Cursor:   "eyJjb2xsZWN0aW9uSWQiOiJ1c2VyMSIsImNvbGxlY3Rpb25DdXJzb3IiOiIxIn0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJ1c2VyMiIsImNvbGxlY3Rpb25DdXJzb3IiOiIyIn0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID:     testutil.GenPtr("user2"),
				CollectionCursor: testutil.GenPtr("2"),
			},
		},

		"valid_user_3_policy_0_of_0": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "UserPolicy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "PolicyArn",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 1,
				Cursor:   "eyJjb2xsZWN0aW9uSWQiOiJ1c2VyMiIsImNvbGxlY3Rpb25DdXJzb3IiOiIyIn0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJ1c2VyMyIsImNvbGxlY3Rpb25DdXJzb3IiOiIzIn0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID:     testutil.GenPtr("user3"),
				CollectionCursor: testutil.GenPtr("3"),
			},
		},
		"valid_user_6_policy_1_of_1": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "UserPolicy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "PolicyArn",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 1,
				Cursor:   "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDUiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiNSJ9",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "arn:aws:iam::000000000000:policy/Policy105-user6",
							"PolicyArn": "arn:aws:iam::000000000000:policy/Policy105",
							"UserName":  "user6",
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

func TestAdapterGetRolePolicyPage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(defaultTimeout)*time.Second)
	defer cancel()

	cfg, err := SetupTestConfig(ctx, Mocker)
	if err != nil {
		log.Fatalf("Failed to load aws test config: %v", err)
	}

	adapter, err := ProvideAWSTestClient(cfg)
	if err != nil {
		log.Fatalf("Failed to load aws test client: %v", err)
	}

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[aws_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_role_1_policy_2_of_2": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "RolePolicy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "PolicyArn",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "RoleName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "arn:aws:iam::000000000000:policy/Policy106-AWSReservedSSO_de25js739eef1832",
							"PolicyArn": "arn:aws:iam::000000000000:policy/Policy106",
							"RoleName":  "AWSReservedSSO_de25js739eef1832",
						},
						{
							"id":        "arn:aws:iam::000000000000:policy/Policy102-AWSReservedSSO_de25js739eef1832",
							"PolicyArn": "arn:aws:iam::000000000000:policy/Policy102",
							"RoleName":  "AWSReservedSSO_de25js739eef1832",
						},
					},
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJBV1NSZXNlcnZlZFNTT19kZTI1anM3MzllZWYxODMyIiwiY29sbGVjdGlvbkN1cnNvciI6IjEifQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID:     testutil.GenPtr("AWSReservedSSO_de25js739eef1832"),
				CollectionCursor: testutil.GenPtr("1"),
			},
		},
		"valid_role_2_policy_1_of_1": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "RolePolicy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "PolicyArn",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "RoleName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjb2xsZWN0aW9uSWQiOiJBV1NSZXNlcnZlZFNTT19kZTI1anM3MzllZWYxODMyIiwiY29sbGVjdGlvbkN1cnNvciI6IjEifQ==",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "arn:aws:iam::000000000000:policy/Policy107-AWSReservedSSO_abcdef1234567890",
							"PolicyArn": "arn:aws:iam::000000000000:policy/Policy107",
							"RoleName":  "AWSReservedSSO_abcdef1234567890",
						},
					},
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJBV1NSZXNlcnZlZFNTT19hYmNkZWYxMjM0NTY3ODkwIiwiY29sbGVjdGlvbkN1cnNvciI6IjIifQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID:     testutil.GenPtr("AWSReservedSSO_abcdef1234567890"),
				CollectionCursor: testutil.GenPtr("2"),
			},
		},
		"valid_role_3_policy_0_of_0": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "RolePolicy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "PolicyArn",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "RoleName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjb2xsZWN0aW9uSWQiOiJBV1NSZXNlcnZlZFNTT19hYmNkZWYxMjM0NTY3ODkwIiwiY29sbGVjdGlvbkN1cnNvciI6IjIifQ==",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{},
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

func TestAdapterGetGroupMembersPage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(defaultTimeout)*time.Second)
	defer cancel()

	cfg, err := SetupTestConfig(ctx, Mocker)
	if err != nil {
		log.Fatalf("Failed to load aws test config: %v", err)
	}

	adapter, err := ProvideAWSTestClient(cfg)
	if err != nil {
		log.Fatalf("Failed to load aws test client: %v", err)
	}

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[aws_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_group_1_user_2_of_2": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "user1-Group1",
							"UserId":    "user1",
							"GroupName": "Group1",
						},
						{
							"id":        "user2-Group1",
							"UserId":    "user2",
							"GroupName": "Group1",
						},
					},
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDEiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiMSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID:     testutil.GenPtr("Group1"),
				CollectionCursor: testutil.GenPtr("1"),
			},
		},
		"valid_group_2_user_1_of_1": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDEiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiMSJ9",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "user1-Group2",
							"UserId":    "user1",
							"GroupName": "Group2",
						},
						{
							"id":        "user2-Group2",
							"UserId":    "user2",
							"GroupName": "Group2",
						},
					},
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDIiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiMiJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID:     testutil.GenPtr("Group2"),
				CollectionCursor: testutil.GenPtr("2"),
			},
		},
		"valid_group_3_user_0_of_0": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDIiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiMiJ9",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDMiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiMyJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID:     testutil.GenPtr("Group3"),
				CollectionCursor: testutil.GenPtr("3"),
			},
		},
		"valid_group_4_user_0_of_0": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDMiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiMyJ9",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					NextCursor: "",
				},
			},
		},
		"valid_group_1_user_2_of_2_with_multiple_accounts_result_for_first_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "user1-Group1",
							"UserId":    "user1",
							"GroupName": "Group1",
						},
						{
							"id":        "user2-Group1",
							"UserId":    "user2",
							"GroupName": "Group1",
						},
					},
					// Base 64 encoded value of: {"collectionId":"Group1","collectionCursor":"eyJPZmZzZXQiOjAsIk5leHRNYXJrZXIiOiIxIn0="}
					// nolint: lll
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDEiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKUFptWnpaWFFpT2pBc0lrNWxlSFJOWVhKclpYSWlPaUl4SW4wPSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID: testutil.GenPtr("Group1"),
				// Base 64 encoded value of: {"Offset":0,"NextMarker":"1"}
				CollectionCursor: testutil.GenPtr("eyJPZmZzZXQiOjAsIk5leHRNYXJrZXIiOiIxIn0="),
			},
		},
		"valid_group_2_user_1_of_1_with_multiple_accounts_result_for_first_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
				// Base64 encoded value of: {"collectionId":"Group1","collectionCursor":"eyJPZmZzZXQiOjAsIk5leHRNYXJrZXIiOiIxIn0="}
				// nolint: lll
				Cursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDEiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKUFptWnpaWFFpT2pBc0lrNWxlSFJOWVhKclpYSWlPaUl4SW4wPSJ9",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "user1-Group2",
							"UserId":    "user1",
							"GroupName": "Group2",
						},
						{
							"id":        "user2-Group2",
							"UserId":    "user2",
							"GroupName": "Group2",
						},
					},
					// Base 64 value of: {"collectionId":"Group2","collectionCursor":"eyJPZmZzZXQiOjAsIk5leHRNYXJrZXIiOiIyIn0="}
					// nolint: lll
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDIiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKUFptWnpaWFFpT2pBc0lrNWxlSFJOWVhKclpYSWlPaUl5SW4wPSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID: testutil.GenPtr("Group2"),
				// Base 64 value of: {"Offset":0,"NextMarker":"2"}
				CollectionCursor: testutil.GenPtr("eyJPZmZzZXQiOjAsIk5leHRNYXJrZXIiOiIyIn0="),
			},
		},
		"valid_group_3_user_0_of_0_with_multiple_accounts_result_for_first_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
				// nolint:lll
				Cursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDIiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKUFptWnpaWFFpT2pBc0lrNWxlSFJOWVhKclpYSWlPaUl5SW4wPSJ9",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					// nolint:lll
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDMiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKUFptWnpaWFFpT2pBc0lrNWxlSFJOWVhKclpYSWlPaUl6SW4wPSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID: testutil.GenPtr("Group3"),
				// Base 64 value of: {"Offset":0,"NextMarker":"3"}
				CollectionCursor: testutil.GenPtr("eyJPZmZzZXQiOjAsIk5leHRNYXJrZXIiOiIzIn0="),
			},
		},
		"valid_group_4_user_0_of_0_with_multiple_accounts_result_for_first_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
				// nolint:lll
				Cursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDMiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKUFptWnpaWFFpT2pBc0lrNWxlSFJOWVhKclpYSWlPaUl6SW4wPSJ9",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					// nolint:lll
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDQiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKUFptWnpaWFFpT2pFc0lrNWxlSFJOWVhKclpYSWlPbTUxYkd4OSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID: testutil.GenPtr("Group4"),
				// Base 64 value of: {"Offset":1,"NextMarker":null}
				CollectionCursor: testutil.GenPtr("eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOm51bGx9"),
			},
		},
		"valid_group_1_user_2_of_2_with_multiple_accounts_result_for_second_account": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
				// nolint:lll
				Cursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDQiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKUFptWnpaWFFpT2pFc0lrNWxlSFJOWVhKclpYSWlPbTUxYkd4OSJ9",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "user1-Group1",
							"UserId":    "user1",
							"GroupName": "Group1",
						},
						{
							"id":        "user2-Group1",
							"UserId":    "user2",
							"GroupName": "Group1",
						},
					},
					// Base 64 encoded value of: {"collectionId":"Group1","collectionCursor":"eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOiIxIn0="}
					// nolint: lll
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDEiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKUFptWnpaWFFpT2pFc0lrNWxlSFJOWVhKclpYSWlPaUl4SW4wPSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				CollectionID: testutil.GenPtr("Group1"),
				// Base 64 encoded value of: {"Offset":0,"NextMarker":"1"}
				CollectionCursor: testutil.GenPtr("eyJPZmZzZXQiOjEsIk5leHRNYXJrZXIiOiIxIn0="),
			},
		},
		"valid_group_4_user_0_of_0_with_multiple_accounts_result_for_second_account_no_more_accounts_left": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfigWithAccounts,
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "UserId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "GroupName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 2,
				// nolint:lll
				Cursor: "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDMiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKUFptWnpaWFFpT2pFc0lrNWxlSFJOWVhKclpYSWlPaUl6SW4wPSJ9",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{},
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
