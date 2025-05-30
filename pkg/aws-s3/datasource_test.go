// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package awss3_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	s3_adapter "github.com/sgnl-ai/adapters/pkg/aws-s3"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestGetObjectKeyFromRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *s3_adapter.Request
		want    string
	}{
		{
			name: "simple",
			request: &s3_adapter.Request{
				PathPrefix:       "data/internal",
				FileType:         "csv",
				EntityExternalID: "users",
			},
			want: "data/internal/users.csv",
		},
		{
			name: "simple_with_trailing_slash",
			request: &s3_adapter.Request{
				PathPrefix:       "data/internal/",
				FileType:         "csv",
				EntityExternalID: "users",
			},
			want: "data/internal/users.csv",
		},
		{
			name: "empty_prefix",
			request: &s3_adapter.Request{
				PathPrefix:       "",
				FileType:         "csv",
				EntityExternalID: "customers",
			},
			want: "customers.csv",
		},
		{
			name: "root_prefix",
			request: &s3_adapter.Request{
				PathPrefix:       "/",
				FileType:         "csv",
				EntityExternalID: "orders",
			},
			want: "/orders.csv",
		},
		{
			name: "nested_path",
			request: &s3_adapter.Request{
				PathPrefix:       "exports/2024/january",
				FileType:         "csv",
				EntityExternalID: "sales",
			},
			want: "exports/2024/january/sales.csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s3_adapter.GetObjectKeyFromRequest(tt.request)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}

func TestDatasource_GetPage(t *testing.T) {
	tests := map[string]struct {
		request              *s3_adapter.Request
		headObjectStatusCode int
		getObjectStatusCode  int
		expectedResponse     *s3_adapter.Response
		expectedError        *framework.Error
	}{
		"success_small_file_traditional_path": {
			request: &s3_adapter.Request{
				Auth: s3_adapter.Auth{
					AccessKey: "test-access-key",
					SecretKey: "test-secret-key",
					Region:    "us-west-1",
				},
				Bucket:                "test-bucket",
				PathPrefix:            "data",
				FileType:              "csv",
				EntityExternalID:      "customers",
				PageSize:              2,
				RequestTimeoutSeconds: 30,
				AttributeConfig: []*framework.AttributeConfig{
					{
						ExternalId: "Email",
						Type:       framework.AttributeTypeString,
						UniqueId:   true,
					},
					{
						ExternalId: "Score",
						Type:       framework.AttributeTypeDouble,
					},
				},
			},
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  http.StatusOK,
			expectedResponse: &s3_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{
						"City":              "Caitlynmouth",
						"Company":           "Blankenship PLC",
						"Country":           "Sao Tome and Principe",
						"Customer Id":       "e685B8690f9fbce",
						"Email":             "shanehester@campbell.org",
						"First Name":        "Erik",
						"KnownAliases":      []any{map[string]any{"alias": "Shane Hester", "primary": true}, map[string]any{"alias": "Cheyne Hester", "primary": false}},
						"Last Name":         "Little",
						"Phone 1":           "457-542-6899",
						"Phone 2":           "055.415.2664x5425",
						"Score":             1.1,
						"Subscription Date": "2021-12-23",
						"Website":           "https://wagner.com/",
					},
					{
						"City":              "Janetfort",
						"Company":           "Jensen and Sons",
						"Country":           "Palestinian Territory",
						"Customer Id":       "6EDdBA3a2DFA7De",
						"Email":             "kleinluis@vang.com",
						"First Name":        "Yvonne",
						"KnownAliases":      []any{map[string]any{"primary": true, "alias": "Klein Luis"}, map[string]any{"alias": "Cline Luis", "primary": false}},
						"Last Name":         "Shaw",
						"Phone 1":           "9610730173",
						"Phone 2":           "531-482-3000x7085",
						"Score":             2.2,
						"Subscription Date": "2021-01-01",
						"Website":           "https://www.paul.org/",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(3)),
				},
			},
		},
		"success_with_cursor": {
			request: &s3_adapter.Request{
				Auth: s3_adapter.Auth{
					AccessKey: "test-access-key",
					SecretKey: "test-secret-key",
					Region:    "us-west-1",
				},
				Bucket:                "test-bucket",
				PathPrefix:            "data",
				FileType:              "csv",
				EntityExternalID:      "customers",
				PageSize:              2,
				RequestTimeoutSeconds: 30,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(3)),
				},
				AttributeConfig: []*framework.AttributeConfig{
					{
						ExternalId: "Email",
						Type:       framework.AttributeTypeString,
						UniqueId:   true,
					},
					{
						ExternalId: "Score",
						Type:       framework.AttributeTypeDouble,
					},
				},
			},
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  http.StatusOK,
			expectedResponse: &s3_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{
						"City":              "Darlenebury",
						"Company":           "Rose, Deleon and Sanders",
						"Country":           "Albania",
						"Customer Id":       "b9Da13bedEc47de",
						"Email":             "deckerjamie@bartlett.biz",
						"First Name":        "Jeffery",
						"KnownAliases":      `[{"alias": "Decker Jaime", "primary": true}`,
						"Last Name":         "Ibarra",
						"Phone 1":           "(840)539-1797x479",
						"Phone 2":           "209-519-5817",
						"Score":             3.3,
						"Subscription Date": "2020-03-30",
						"Website":           "https://www.morgan-phelps.com/",
					},
					{
						"City":              "Donhaven",
						"Company":           "Kline and Sons",
						"Country":           "Bahrain",
						"Customer Id":       "710D4dA2FAa96B5",
						"Email":             "dochoa@carey-morse.com",
						"First Name":        "James",
						"KnownAliases":      []any{map[string]any{"alias": "Do Choa", "primary": true}},
						"Last Name":         "Walters",
						"Phone 1":           "+1-985-596-1072x3040",
						"Phone 2":           "(528)734-8924x054",
						"Score":             4.4,
						"Subscription Date": "2022-01-18",
						"Website":           "https://brennan.com/",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(5)),
				},
			},
		},
		"success_last_page_no_cursor": {
			request: &s3_adapter.Request{
				Auth: s3_adapter.Auth{
					AccessKey: "test-access-key",
					SecretKey: "test-secret-key",
					Region:    "us-west-1",
				},
				Bucket:                "test-bucket",
				PathPrefix:            "data",
				FileType:              "csv",
				EntityExternalID:      "customers",
				PageSize:              2,
				RequestTimeoutSeconds: 30,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(5)),
				},
				AttributeConfig: []*framework.AttributeConfig{
					{
						ExternalId: "Email",
						Type:       framework.AttributeTypeString,
						UniqueId:   true,
					},
					{
						ExternalId: "Score",
						Type:       framework.AttributeTypeDouble,
					},
				},
			},
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  http.StatusOK,
			expectedResponse: &s3_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{
						"City":              "Mossfort",
						"Company":           "Price, Mason and Doyle",
						"Country":           "Central African Republic",
						"Customer Id":       "3c44ed62d7BfEBC",
						"Email":             "darrylbarber@warren.org",
						"First Name":        "Leslie",
						"KnownAliases":      `[{"alias": "Darryl Barber", "primary": true}`,
						"Last Name":         "Snyder",
						"Phone 1":           "812-016-9904x8231",
						"Phone 2":           "254.631.9380",
						"Score":             5.5,
						"Subscription Date": "2020-01-25",
						"Website":           "http://www.trujillo-sullivan.info/",
					},
				},
				NextCursor: nil, // Last page
			},
		},
		"success_headers_only_file": {
			request: &s3_adapter.Request{
				Auth: s3_adapter.Auth{
					AccessKey: "test-access-key",
					SecretKey: "test-secret-key",
					Region:    "us-west-1",
				},
				Bucket:                "test-bucket",
				PathPrefix:            "data",
				FileType:              "csv",
				EntityExternalID:      "customers",
				PageSize:              2,
				RequestTimeoutSeconds: 30,
				AttributeConfig: []*framework.AttributeConfig{
					{
						ExternalId: "Email",
						Type:       framework.AttributeTypeString,
						UniqueId:   true,
					},
				},
			},
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  headersOnlyCSVFileCode,
			expectedResponse: &s3_adapter.Response{
				StatusCode: 200,
				Objects:    []map[string]any{},
				NextCursor: nil,
			},
		},
		"success_large_file_streaming_path": {
			request: &s3_adapter.Request{
				Auth: s3_adapter.Auth{
					AccessKey: "test-access-key",
					SecretKey: "test-secret-key",
					Region:    "us-west-1",
				},
				Bucket:                "test-bucket",
				PathPrefix:            "data",
				FileType:              "csv",
				EntityExternalID:      "large-customers",
				PageSize:              100,
				RequestTimeoutSeconds: 30,
				AttributeConfig: []*framework.AttributeConfig{
					{
						ExternalId: "Email",
						Type:       framework.AttributeTypeString,
						UniqueId:   true,
					},
					{
						ExternalId: "Score",
						Type:       framework.AttributeTypeDouble,
					},
				},
			},
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  largeCSVFileCode,
			expectedResponse: &s3_adapter.Response{
				StatusCode: 200,
				// We'll validate objects separately since we can't predict exact content
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(101)),
				},
			},
		},
		"error_file_not_found_head_object": {
			request: &s3_adapter.Request{
				Auth: s3_adapter.Auth{
					AccessKey: "test-access-key",
					SecretKey: "test-secret-key",
					Region:    "us-west-1",
				},
				Bucket:                "test-bucket",
				PathPrefix:            "data",
				FileType:              "csv",
				EntityExternalID:      "missing",
				PageSize:              2,
				RequestTimeoutSeconds: 30,
				AttributeConfig: []*framework.AttributeConfig{
					{
						ExternalId: "Email",
						Type:       framework.AttributeTypeString,
						UniqueId:   true,
					},
				},
			},
			headObjectStatusCode: http.StatusNotFound,
			getObjectStatusCode:  http.StatusOK,
			expectedError: &framework.Error{
				Message: "Failed to fetch entity from AWS S3: missing, error: failed to check if the file exists: operation error S3: HeadObject, http response error StatusCode: 404, not found: The specified key does not exist.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"error_permission_denied_head_object": {
			request: &s3_adapter.Request{
				Auth: s3_adapter.Auth{
					AccessKey: "test-access-key",
					SecretKey: "test-secret-key",
					Region:    "us-west-1",
				},
				Bucket:                "test-bucket",
				PathPrefix:            "data",
				FileType:              "csv",
				EntityExternalID:      "forbidden",
				PageSize:              2,
				RequestTimeoutSeconds: 30,
				AttributeConfig: []*framework.AttributeConfig{
					{
						ExternalId: "Email",
						Type:       framework.AttributeTypeString,
						UniqueId:   true,
					},
				},
			},
			headObjectStatusCode: http.StatusForbidden,
			getObjectStatusCode:  http.StatusOK,
			expectedError: &framework.Error{
				Message: "Failed to fetch entity from AWS S3: forbidden, error: failed to check if the file exists: operation error S3: HeadObject, http response error StatusCode: 403, AccessDenied: Access Denied.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"error_empty_csv_file": {
			request: &s3_adapter.Request{
				Auth: s3_adapter.Auth{
					AccessKey: "test-access-key",
					SecretKey: "test-secret-key",
					Region:    "us-west-1",
				},
				Bucket:                "test-bucket",
				PathPrefix:            "data",
				FileType:              "csv",
				EntityExternalID:      "empty",
				PageSize:              2,
				RequestTimeoutSeconds: 30,
				AttributeConfig: []*framework.AttributeConfig{
					{
						ExternalId: "Email",
						Type:       framework.AttributeTypeString,
						UniqueId:   true,
					},
				},
			},
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  emptyCSVFileCode,
			expectedError: &framework.Error{
				Message: "Failed to fetch entity from AWS S3: empty, error: unable to process CSV file data: no data found in the CSV file.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"error_corrupted_csv_file": {
			request: &s3_adapter.Request{
				Auth: s3_adapter.Auth{
					AccessKey: "test-access-key",
					SecretKey: "test-secret-key",
					Region:    "us-west-1",
				},
				Bucket:                "test-bucket",
				PathPrefix:            "data",
				FileType:              "csv",
				EntityExternalID:      "corrupt",
				PageSize:              2,
				RequestTimeoutSeconds: 30,
				AttributeConfig: []*framework.AttributeConfig{
					{
						ExternalId: "Email",
						Type:       framework.AttributeTypeString,
						UniqueId:   true,
					},
				},
			},
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  -200,
			expectedError: &framework.Error{
				Message: "Failed to fetch entity from AWS S3: corrupt, error: unable to process CSV file data: failed to read CSV data: parse error on line 4, column 34: bare \" in non-quoted-field.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"error_permission_denied_get_object": {
			request: &s3_adapter.Request{
				Auth: s3_adapter.Auth{
					AccessKey: "test-access-key",
					SecretKey: "test-secret-key",
					Region:    "us-west-1",
				},
				Bucket:                "test-bucket",
				PathPrefix:            "data",
				FileType:              "csv",
				EntityExternalID:      "forbidden-get",
				PageSize:              2,
				RequestTimeoutSeconds: 30,
				AttributeConfig: []*framework.AttributeConfig{
					{
						ExternalId: "Email",
						Type:       framework.AttributeTypeString,
						UniqueId:   true,
					},
				},
			},
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  http.StatusForbidden,
			expectedError: &framework.Error{
				Message: "Failed to fetch entity from AWS S3: forbidden-get, error: unable to read CSV file: failed to convert response: operation error S3: GetObject, http response error StatusCode: 403, access denied: Access Denied.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			awsConfig := mockS3Config(tt.headObjectStatusCode, tt.getObjectStatusCode)

			datasource, err := s3_adapter.NewClient(http.DefaultClient, awsConfig)
			if err != nil {
				t.Fatalf("Failed to create datasource: %v", err)
			}

			ctx := context.Background()
			response, frameworkErr := datasource.GetPage(ctx, tt.request)

			if tt.expectedError != nil {
				if frameworkErr == nil {
					t.Errorf("Expected error but got none")
				} else {
					if frameworkErr.Message != tt.expectedError.Message {
						t.Errorf("Expected error message '%s', got '%s'", tt.expectedError.Message, frameworkErr.Message)
					}
					if frameworkErr.Code != tt.expectedError.Code {
						t.Errorf("Expected error code %v, got %v", tt.expectedError.Code, frameworkErr.Code)
					}
				}
				if response != nil {
					t.Errorf("Expected nil response on error, got %v", response)
				}
			} else {
				if frameworkErr != nil {
					t.Errorf("Expected no error, got: %v", frameworkErr)
				}
				if response == nil {
					t.Errorf("Expected response, got nil")
				} else {
					if name == "success_large_file_streaming_path" {
						validateLargeFileResponse(t, response, tt.expectedResponse)
					} else {
						if !reflect.DeepEqual(response, tt.expectedResponse) {
							t.Errorf("Response mismatch.\nGot: %+v\nWant: %+v", response, tt.expectedResponse)
						}
					}
				}
			}
		})
	}
}

// since we can't predict the exact content but can validate structure and pagination
func validateLargeFileResponse(t *testing.T, got, want *s3_adapter.Response) {
	if got.StatusCode != want.StatusCode {
		t.Errorf("Expected StatusCode %d, got %d", want.StatusCode, got.StatusCode)
	}

	// Should return exactly PageSize objects (100)
	if len(got.Objects) != 100 {
		t.Errorf("Expected 100 objects for large file test, got %d", len(got.Objects))
	}

	// Should have NextCursor indicating more data
	if got.NextCursor == nil {
		t.Error("Expected NextCursor for large file test, got nil")
	} else if *got.NextCursor.Cursor != 101 {
		t.Errorf("Expected NextCursor value 101, got %d", *got.NextCursor.Cursor)
	}

	// Validate structure of first object
	if len(got.Objects) > 0 {
		firstObj := got.Objects[0]

		// Check that required fields exist
		if _, exists := firstObj["Email"]; !exists {
			t.Error("Expected Email field in large file response object")
		}
		if _, exists := firstObj["Score"]; !exists {
			t.Error("Expected Score field in large file response object")
		}
		if _, exists := firstObj["Customer Id"]; !exists {
			t.Error("Expected Customer Id field in large file response object")
		}

		// Validate first row has expected pattern
		if email, ok := firstObj["Email"].(string); ok {
			if email != "user1@example.com" {
				t.Errorf("Expected first row email 'user1@example.com', got '%s'", email)
			}
		} else {
			t.Error("Email field should be string")
		}

		if score, ok := firstObj["Score"].(float64); ok {
			if score != 0.1 {
				t.Errorf("Expected first row score 0.1, got %f", score)
			}
		} else {
			t.Error("Score field should be float64")
		}
	}

	t.Logf("Large file test: Successfully processed %d objects with streaming", len(got.Objects))
}
