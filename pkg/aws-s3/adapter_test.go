// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package awss3_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	s3_adapter "github.com/sgnl-ai/adapters/pkg/aws-s3"
	"github.com/sgnl-ai/adapters/pkg/testutil"

	"github.com/sgnl-ai/adapters/pkg/pagination"
)

func TestAdapterGetPage(t *testing.T) {
	tests := map[string]struct {
		ctx                  context.Context
		request              *framework.Request[s3_adapter.Config]
		wantResponse         framework.Response
		wantCursor           *pagination.CompositeCursor[int64]
		headObjectStatusCode int
		getObjectStatusCode  int
	}{
		"success_HeadObject_200_GetObject_200_first_page": {
			headObjectStatusCode: 200,
			getObjectStatusCode:  200,
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Score",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "Subscription Date",
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
							"Email":             "shanehester@campbell.org",
							"Score":             int64(1),
							"Subscription Date": time.Date(2021, 12, 23, 0, 0, 0, 0, time.UTC),
						},
						{
							"Email":             "kleinluis@vang.com",
							"Score":             int64(2),
							"Subscription Date": time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOjY1NX0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(655)),
			},
		},
		"success_HeadObject_200_GetObject_200_middle_page": {
			headObjectStatusCode: 200,
			getObjectStatusCode:  200,
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Score",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "Subscription Date",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOjY1NX0=",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Email":             "deckerjamie@bartlett.biz",
							"Score":             int64(3),
							"Subscription Date": time.Date(2020, 3, 30, 0, 0, 0, 0, time.UTC),
						},
						{
							"Email":             "dochoa@carey-morse.com",
							"Score":             int64(4),
							"Subscription Date": time.Date(2022, 1, 18, 0, 0, 0, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOjEwOTV9",
				},
			},
		},
		"success_HeadObject_200_GetObject_200_last_page": {
			headObjectStatusCode: 200,
			getObjectStatusCode:  200,
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Score",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "Subscription Date",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOjEwOTV9",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Email":             "darrylbarber@warren.org",
							"Score":             int64(5),
							"Subscription Date": time.Date(2020, 1, 25, 0, 0, 0, 0, time.UTC),
						},
					},
				},
			},
		},
		"success_headers_only_csv_file_HeadObject_200_GetObject_801": {
			headObjectStatusCode: 200,
			getObjectStatusCode:  headersOnlyCSVFileCode, // Custom status code to simplify testing
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Score",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "Subscription Date",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: nil,
				},
			},
		},
		// Check if a number in the CSV can be ingested as a string type based on entity configuration
		"success_read_numbers_as_strings_HeadObject_200_GetObject_200": {
			headObjectStatusCode: 200,
			getObjectStatusCode:  200,
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Score",
							Type:       framework.AttributeTypeString, //  This is explicitly set to string
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Email": "shanehester@campbell.org",
							"Score": "1.1",
						},
						{
							"Email": "kleinluis@vang.com",
							"Score": "2.2",
						},
					},
					NextCursor: "eyJjdXJzb3IiOjY1NX0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(655)),
			},
		},
		"success_read_child_objects_HeadObject_200_GetObject_200": {
			headObjectStatusCode: 200,
			getObjectStatusCode:  200,
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "KnownAliases",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "alias",
									Type:       framework.AttributeTypeString,
								},
								{
									ExternalId: "primary",
									Type:       framework.AttributeTypeBool,
								},
							},
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Email": "shanehester@campbell.org",
							`KnownAliases`: []framework.Object{
								{
									"alias":   string("Shane Hester"),
									"primary": bool(true),
								},
								{
									"alias":   string("Cheyne Hester"),
									"primary": bool(false),
								},
							},
						},
						{
							"Email": "kleinluis@vang.com",
							`KnownAliases`: []framework.Object{
								{
									"alias":   string("Klein Luis"),
									"primary": bool(true),
								},
								{
									"alias":   string("Cline Luis"),
									"primary": bool(false),
								},
							},
						},
					},
					NextCursor: "eyJjdXJzb3IiOjY1NX0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(655)),
			},
		},
		// Check if a number in the CSV can be ingested as a double type based on entity configuration
		"success_read_numbers_as_double_HeadObject_200_GetObject_200": {
			headObjectStatusCode: 200,
			getObjectStatusCode:  200,
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Score",
							Type:       framework.AttributeTypeDouble, //  This is explicitly set to double
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Email": "shanehester@campbell.org",
							"Score": float64(1.1),
						},
						{
							"Email": "kleinluis@vang.com",
							"Score": float64(2.2),
						},
					},
					NextCursor: "eyJjdXJzb3IiOjY1NX0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(655)),
			},
		},
		"error_empty_csv_file_HeadObject_200_GetObject_800": {
			headObjectStatusCode: 200,
			getObjectStatusCode:  emptyCSVFileCode, // Custom status code to simplify testing
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Score",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "Subscription Date",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					// nolint: lll
					Message: "Unable to parse CSV file headers: CSV header error: empty or missing",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"error_corrupt_csv_HeadObject_200_GetObject_200": {
			headObjectStatusCode: 200,
			getObjectStatusCode:  -200, // This status code returns a corrupt CSV
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Score",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "Subscription Date",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 5,
				Cursor:   "",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					// nolint: lll
					Message: "Failed to fetch entity from AWS S3: customers, error: CSV file format is invalid or corrupted: parse error on line 1, column 34: bare \" in non-quoted-field.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"error_object_moved_HeadObject_301_GetObject_200": {
			headObjectStatusCode: 301,
			getObjectStatusCode:  200,
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Score",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "Subscription Date",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOjY1NX0=",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					// nolint: lll
					Message: "Failed to fetch entity from AWS S3: customers, error: failed to convert response: operation error S3: HeadObject, http response error StatusCode: 301, permanent redirect: The bucket you are attempting to access must be addressed using the specified endpoint.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"error_no_permission_to_HeadObject_HeadObject_403_GetObject_200": {
			headObjectStatusCode: 403,
			getObjectStatusCode:  200,
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Score",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "Subscription Date",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOjY1NX0=",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					// nolint: lll
					Message: "Failed to fetch entity from AWS S3: customers, error: failed to convert response: operation error S3: HeadObject, http response error StatusCode: 403, AccessDenied: Access Denied.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"error_object_not_found_for_HeadObject_HeadObject_404_GetObject_200": {
			headObjectStatusCode: 404,
			getObjectStatusCode:  200,
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Score",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "Subscription Date",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOjY1NX0=",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					// nolint: lll
					Message: "Failed to fetch entity from AWS S3: customers, error: failed to convert response: operation error S3: HeadObject, http response error StatusCode: 404, not found: The specified key does not exist.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"error_object_moved_before_GetObject_HeadObject_200_GetObject_301": {
			headObjectStatusCode: 200,
			getObjectStatusCode:  301,
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Score",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "Subscription Date",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOjY1NX0=",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					// nolint: lll
					Message: "Failed to fetch entity from AWS S3: customers, error: failed to convert response: operation error S3: GetObject, http response error StatusCode: 301, permanent redirect: The bucket you are attempting to access must be addressed using the specified endpoint.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"error_no_permission_to_GetObject_HeadObject_200_GetObject_403": {
			headObjectStatusCode: 200,
			getObjectStatusCode:  403,
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Score",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "Subscription Date",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOjY1NX0=",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					// nolint: lll
					Message: "Failed to fetch entity from AWS S3: customers, error: failed to convert response: operation error S3: GetObject, http response error StatusCode: 403, access denied: Access Denied.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"error_object_deleted_after_HeadObject_and_before_GetObject_HeadObject_200_GetObject_403": {
			headObjectStatusCode: 200,
			getObjectStatusCode:  404,
			ctx:                  context.Background(),
			request: &framework.Request[s3_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "customers",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Email",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Score",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "Subscription Date",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOjY1NX0=",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					// nolint: lll
					Message: "Failed to fetch entity from AWS S3: customers, error: failed to convert response: operation error S3: GetObject, http response error StatusCode: 404, no such key: The specified key does not exist.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup mock middleware to mimic responses from the SDK
			cfg := mockS3Config(tt.headObjectStatusCode, tt.getObjectStatusCode)

			client, err := s3_adapter.NewClient(http.DefaultClient, cfg, MaxCSVRowSizeBytes, MaxBytesToProcessPerPage)
			if err != nil {
				t.Errorf("error creating client to query datasource: %v", err)
			}

			adapter := s3_adapter.NewAdapter(client)

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
