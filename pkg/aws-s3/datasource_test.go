// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package awss3_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	s3_adapter "github.com/sgnl-ai/adapters/pkg/aws-s3"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

const (
	validCSVDataHeaderLength = 121
	validCSVDataRow1Length   = 274
	validCSVDataRow2Length   = 260
	validCSVDataRow3Length   = 232
	validCSVDataRow4Length   = 208
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
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetObjectKeyFromRequest() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDatasource_GetPage(t *testing.T) {
	cursorPage1Next := int64(validCSVDataHeaderLength + validCSVDataRow1Length + validCSVDataRow2Length)
	cursorPage2Start := cursorPage1Next
	cursorPage2Next := cursorPage2Start + validCSVDataRow3Length + validCSVDataRow4Length
	cursorPage3Start := cursorPage2Next

	tests := map[string]struct {
		request              *s3_adapter.Request
		headObjectStatusCode int
		getObjectStatusCode  int
		expectedResponse     *s3_adapter.Response
		expectedError        *framework.Error
	}{
		"success_small_file_traditional_path": {
			request: &s3_adapter.Request{
				Auth:   s3_adapter.Auth{AccessKey: "test-access-key", SecretKey: "test-secret-key", Region: "us-west-1"},
				Bucket: "test-bucket", PathPrefix: "data", FileType: "csv", EntityExternalID: "customers",
				PageSize: 2, RequestTimeoutSeconds: 30,
				AttributeConfig: []*framework.AttributeConfig{{ExternalId: "Email", Type: framework.AttributeTypeString,
					UniqueId: true}, {ExternalId: "Score", Type: framework.AttributeTypeDouble}},
			},
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  http.StatusOK,
			expectedResponse: &s3_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"City": "Caitlynmouth", "Company": "Blankenship PLC", "Country": "Sao Tome and Principe",
						"Customer Id": "e685B8690f9fbce", "Email": "shanehester@campbell.org", "First Name": "Erik",
						"KnownAliases": []any{map[string]any{"alias": "Shane Hester", "primary": true},
							map[string]any{"alias": "Cheyne Hester", "primary": false}}, "Last Name": "Little",
						"Phone 1": "457-542-6899", "Phone 2": "055.415.2664x5425", "Score": 1.1,
						"Subscription Date": "2021-12-23", "Website": "https://wagner.com/"},
					{"City": "Janetfort", "Company": "Jensen and Sons", "Country": "Palestinian Territory",
						"Customer Id": "6EDdBA3a2DFA7De", "Email": "kleinluis@vang.com", "First Name": "Yvonne",
						"KnownAliases": []any{map[string]any{"primary": true, "alias": "Klein Luis"},
							map[string]any{"alias": "Cline Luis", "primary": false}}, "Last Name": "Shaw",
						"Phone 1": "9610730173", "Phone 2": "531-482-3000x7085", "Score": 2.2,
						"Subscription Date": "2021-01-01", "Website": "https://www.paul.org/"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{Cursor: testutil.GenPtr(cursorPage1Next)},
			},
		},
		"success_with_cursor": {
			request: &s3_adapter.Request{
				Auth:   s3_adapter.Auth{AccessKey: "test-access-key", SecretKey: "test-secret-key", Region: "us-west-1"},
				Bucket: "test-bucket", PathPrefix: "data", FileType: "csv", EntityExternalID: "customers",
				PageSize: 2, RequestTimeoutSeconds: 30,
				Cursor: &pagination.CompositeCursor[int64]{Cursor: testutil.GenPtr(cursorPage2Start)},
				AttributeConfig: []*framework.AttributeConfig{{ExternalId: "Email", Type: framework.AttributeTypeString,
					UniqueId: true}, {ExternalId: "Score", Type: framework.AttributeTypeDouble}},
			},
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  http.StatusOK,
			expectedResponse: &s3_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"City": "Darlenebury", "Company": "Rose, Deleon and Sanders", "Country": "Albania",
						"Customer Id": "b9Da13bedEc47de", "Email": "deckerjamie@bartlett.biz", "First Name": "Jeffery",
						"KnownAliases": "[{\"alias\": \"Decker Jaime\", \"primary\": true}", "Last Name": "Ibarra",
						"Phone 1": "(840)539-1797x479", "Phone 2": "209-519-5817", "Score": 3.3,
						"Subscription Date": "2020-03-30", "Website": "https://www.morgan-phelps.com/"},
					{"City": "Donhaven", "Company": "Kline and Sons", "Country": "Bahrain",
						"Customer Id": "710D4dA2FAa96B5", "Email": "dochoa@carey-morse.com", "First Name": "James",
						"KnownAliases": []any{map[string]any{"alias": "Do Choa", "primary": true}}, "Last Name": "Walters",
						"Phone 1": "+1-985-596-1072x3040", "Phone 2": "(528)734-8924x054", "Score": 4.4,
						"Subscription Date": "2022-01-18", "Website": "https://brennan.com/"},
				},
				NextCursor: &pagination.CompositeCursor[int64]{Cursor: testutil.GenPtr(cursorPage2Next)},
			},
		},
		"success_last_page_no_cursor": {
			request: &s3_adapter.Request{
				Auth:   s3_adapter.Auth{AccessKey: "test-access-key", SecretKey: "test-secret-key", Region: "us-west-1"},
				Bucket: "test-bucket", PathPrefix: "data", FileType: "csv", EntityExternalID: "customers",
				PageSize: 2, RequestTimeoutSeconds: 30,
				Cursor: &pagination.CompositeCursor[int64]{Cursor: testutil.GenPtr(cursorPage3Start)},
				AttributeConfig: []*framework.AttributeConfig{{ExternalId: "Email", Type: framework.AttributeTypeString,
					UniqueId: true}, {ExternalId: "Score", Type: framework.AttributeTypeDouble}},
			},
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  http.StatusOK,
			expectedResponse: &s3_adapter.Response{
				StatusCode: 200,
				Objects: []map[string]any{
					{"City": "Mossfort", "Company": "Price, Mason and Doyle", "Country": "Central African Republic",
						"Customer Id": "3c44ed62d7BfEBC", "Email": "darrylbarber@warren.org", "First Name": "Leslie",
						"KnownAliases": "[{\"alias\": \"Darryl Barber\", \"primary\": true}", "Last Name": "Snyder",
						"Phone 1": "812-016-9904x8231", "Phone 2": "254.631.9380", "Score": 5.5,
						"Subscription Date": "2020-01-25", "Website": "http://www.trujillo-sullivan.info/"},
				},
				NextCursor: nil,
			},
		},
		"success_headers_only_file": {
			request: &s3_adapter.Request{
				Auth:   s3_adapter.Auth{AccessKey: "test-access-key", SecretKey: "test-secret-key", Region: "us-west-1"},
				Bucket: "test-bucket", PathPrefix: "data", FileType: "csv", EntityExternalID: "customers_headers_only",
				PageSize: 2, RequestTimeoutSeconds: 30,
				AttributeConfig: []*framework.AttributeConfig{{ExternalId: "Email", Type: framework.AttributeTypeString,
					UniqueId: true}},
			},
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  headersOnlyCSVFileCode,
			expectedResponse:     &s3_adapter.Response{StatusCode: 200, Objects: []map[string]any{}, NextCursor: nil},
		},
		"success_large_file_streaming_path": {
			request: &s3_adapter.Request{
				Auth:   s3_adapter.Auth{AccessKey: "test-access-key", SecretKey: "test-secret-key", Region: "us-west-1"},
				Bucket: "test-bucket", PathPrefix: "data", FileType: "csv", EntityExternalID: "large-customers",
				PageSize: 100, RequestTimeoutSeconds: 30,
				AttributeConfig: []*framework.AttributeConfig{{ExternalId: "Email", Type: framework.AttributeTypeString,
					UniqueId: true}, {ExternalId: "Score", Type: framework.AttributeTypeDouble}},
			},
			headObjectStatusCode: largeFileHeaderIndicatorCode, // Use indicator for mock to return large ContentLength
			getObjectStatusCode:  largeCSVFileCode,             // Use indicator for mock to serve large data
			expectedResponse:     &s3_adapter.Response{StatusCode: 200},
		},
		"error_file_not_found_head_object": {
			request: &s3_adapter.Request{
				Auth:   s3_adapter.Auth{AccessKey: "test-access-key", SecretKey: "test-secret-key", Region: "us-west-1"},
				Bucket: "test-bucket", PathPrefix: "data", FileType: "csv", EntityExternalID: "missing",
				PageSize: 2, RequestTimeoutSeconds: 30,
			},
			headObjectStatusCode: http.StatusNotFound,
			getObjectStatusCode:  http.StatusOK,
			expectedError: &framework.Error{ // Adjusted to match datasource.go's direct wrapping of S3 error
				Message: "Failed to fetch entity from AWS S3: missing, error: failed to convert response: " +
					"operation error S3: HeadObject, http response error StatusCode: 404, not found: " +
					"The specified key does not exist.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"error_permission_denied_head_object": {
			request: &s3_adapter.Request{
				Auth:   s3_adapter.Auth{AccessKey: "test-access-key", SecretKey: "test-secret-key", Region: "us-west-1"},
				Bucket: "test-bucket", PathPrefix: "data", FileType: "csv", EntityExternalID: "forbidden",
				PageSize: 2, RequestTimeoutSeconds: 30,
			},
			headObjectStatusCode: http.StatusForbidden,
			getObjectStatusCode:  http.StatusOK,
			expectedError: &framework.Error{ // Adjusted
				Message: "Failed to fetch entity from AWS S3: forbidden, error: failed to convert response: operation error S3: " +
					"HeadObject, http response error StatusCode: 403, AccessDenied: Access Denied.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"error_empty_csv_file_header_parse_fail": {
			request: &s3_adapter.Request{
				Auth:   s3_adapter.Auth{AccessKey: "test-access-key", SecretKey: "test-secret-key", Region: "us-west-1"},
				Bucket: "test-bucket", PathPrefix: "data", FileType: "csv", EntityExternalID: "empty_csv_content",
				PageSize: 2, RequestTimeoutSeconds: 30,
			},
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  emptyCSVFileCode,
			expectedError: &framework.Error{
				Message: "Unable to parse CSV file headers: CSV header is empty or missing",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"error_corrupted_csv_file": {
			request: &s3_adapter.Request{
				Auth:   s3_adapter.Auth{AccessKey: "test-access-key", SecretKey: "test-secret-key", Region: "us-west-1"},
				Bucket: "test-bucket", PathPrefix: "data", FileType: "csv", EntityExternalID: "corrupt",
				PageSize: 5, RequestTimeoutSeconds: 30,
			},
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  -200, // CorruptCSVData
			expectedError: &framework.Error{ // Adjusted to match datasource.go's direct wrapping
				Message: "Failed to fetch entity from AWS S3: corrupt, error: CSV file format is invalid or corrupted: " +
					"parse error on line 1, column 34: bare \" in non-quoted-field.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			currentMockS3Config := mockS3Config(tt.headObjectStatusCode, tt.getObjectStatusCode)

			datasource, err := s3_adapter.NewClient(http.DefaultClient, currentMockS3Config)
			if err != nil {
				t.Fatalf("Failed to create datasource: %v", err)
			}

			ctx := context.Background()
			response, frameworkErr := datasource.GetPage(ctx, tt.request)

			if tt.expectedError != nil {
				validateErrorCase(t, frameworkErr, response, tt.expectedError)
			} else {
				validateSuccessCase(t, frameworkErr, response, tt.expectedResponse, name, int(tt.request.PageSize))
			}
		})
	}
}

func validateErrorCase(t *testing.T, frameworkErr *framework.Error,
	response *s3_adapter.Response, expectedError *framework.Error) {
	if frameworkErr == nil {
		t.Errorf("Expected error but got none. Expected Message: %s", expectedError.Message)

		return
	}

	if !strings.Contains(frameworkErr.Message, expectedError.Message) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedError.Message, frameworkErr.Message)
	}

	if frameworkErr.Code != expectedError.Code {
		t.Errorf("Expected error code %v, got %v", expectedError.Code, frameworkErr.Code)
	}

	if response != nil {
		t.Errorf("Expected nil response on error, got %v", response)
	}
}

func validateSuccessCase(t *testing.T, frameworkErr *framework.Error,
	response, expectedResponse *s3_adapter.Response, name string, requestPageSize int) { // Added requestPageSize
	if frameworkErr != nil {
		t.Errorf("Expected no error, got: %+v", frameworkErr)

		return
	}

	if response == nil {
		t.Errorf("Expected response, got nil")

		return
	}

	if name == "success_large_file_streaming_path" {
		// For the large file test, use the specialized validator.
		// Pass the expected StatusCode from tt.expectedResponse and PageSize from tt.request.
		validateLargeFileResponse(t, response, expectedResponse.StatusCode, requestPageSize)
	} else {
		// DeepEqual for other standard success cases
		if diff := cmp.Diff(expectedResponse, response); diff != "" {
			t.Errorf("Response mismatch (-want +got):\n%s", diff)
		}
	}
}

func validateLargeFileResponse(t *testing.T, got *s3_adapter.Response, expectedStatusCode int, expectedNumObjects int) {
	if got.StatusCode != expectedStatusCode {
		t.Errorf("Expected StatusCode %d, got %d", expectedStatusCode, got.StatusCode)
	}

	if len(got.Objects) != expectedNumObjects {
		t.Errorf("Expected %d objects for this page of large file test, got %d", expectedNumObjects, len(got.Objects))
	}

	// Check if a NextCursor SHOULD exist and is populated
	// For the first page of a large file, we expect a NextCursor.
	if got.NextCursor == nil {
		t.Errorf("Expected NextCursor to be populated (not nil) for this large file page, but got nil")
	} else {
		// NextCursor itself is not nil, now check its inner Cursor field
		if got.NextCursor.Cursor == nil {
			t.Errorf("Expected NextCursor.Cursor to point to a value (not be nil), but it was nil")
		} else {
			// Successfully found a non-nil cursor value
			t.Logf("NextCursor.Cursor is populated and has a value: %d", *got.NextCursor.Cursor)
		}
	}

	if len(got.Objects) > 0 {
		validateFirstObjectOfLargeFile(t, got.Objects[0]) // Assuming this helper is still relevant
	}
}

func validateFirstObjectOfLargeFile(t *testing.T, firstObj map[string]any) {
	expectedEmail := "user1@example.com"
	expectedScore := 0.1

	if email, ok := firstObj["Email"].(string); ok {
		if email != expectedEmail {
			t.Errorf("Expected first row email '%s', got '%s'", expectedEmail, email)
		}
	} else {
		t.Errorf("Email field should be string, got %T", firstObj["Email"])
	}

	if score, ok := firstObj["Score"].(float64); ok {
		if score != expectedScore {
			t.Errorf("Expected first row score %f, got %f", expectedScore, score)
		}
	} else {
		t.Errorf("Score field should be float64, got %T", firstObj["Score"])
	}

	if customerID, ok := firstObj["Customer Id"].(string); ok {
		if customerID != "ID0000001" {
			t.Errorf("Expected first row Customer Id 'ID0000001', got '%s'", customerID)
		}
	} else {
		t.Errorf("Customer Id field should be string, got %T", firstObj["Customer Id"])
	}
}
