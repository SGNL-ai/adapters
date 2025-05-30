// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package awss3_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	s3_adapter "github.com/sgnl-ai/adapters/pkg/aws-s3"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestCSVHeaders(t *testing.T) {
	tests := map[string]struct {
		input         *[]byte
		expectedError bool
		errorContains string
		expected      []string
	}{
		"valid_headers": {
			input:    testutil.GenPtr([]byte("name,age,city\nJohn,25,NYC")),
			expected: []string{"name", "age", "city"},
		},
		"headers_with_spaces": {
			input:    testutil.GenPtr([]byte("First Name, Last Name, Email Address\nJohn,Doe,john@example.com")),
			expected: []string{"First Name", " Last Name", " Email Address"},
		},
		"single_header": {
			input:    testutil.GenPtr([]byte("email\ntest@example.com")),
			expected: []string{"email"},
		},
		"empty_input": {
			input:         testutil.GenPtr([]byte("")),
			expectedError: true,
			errorContains: "CSV header is empty or missing",
		},
		"nil_input": {
			input:         nil,
			expectedError: true,
			errorContains: "CSV header is empty or missing",
		},
		"invalid_csv_format": {
			input:         testutil.GenPtr([]byte("\"unclosed quote field")),
			expectedError: true,
			errorContains: "CSV file format is invalid or corrupted",
		},
		"headers_only": {
			input:    testutil.GenPtr([]byte("Score,Customer Id,First Name")),
			expected: []string{"Score", "Customer Id", "First Name"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := s3_adapter.CSVHeaders(tt.input)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !cmp.Equal(err.Error(), tt.errorContains) &&
					!contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}

				if result != nil {
					t.Errorf("Expected nil result on error, got: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}

				if diff := cmp.Diff(result, tt.expected); diff != "" {
					t.Errorf("Headers mismatch: %s", diff)
				}
			}
		})
	}
}

func TestCSVBytesToPage(t *testing.T) {
	tests := map[string]struct {
		data            *[]byte
		start           int64
		pageSize        int64
		attrConfig      []*framework.AttributeConfig
		expectedObjects []map[string]any
		expectedHasNext bool
		expectedError   bool
		errorContains   string
	}{
		"success_basic_csv": {
			data:     testutil.GenPtr([]byte("name,age\nJohn,25\nJane,30\nBob,35")),
			start:    1,
			pageSize: 2,
			attrConfig: []*framework.AttributeConfig{
				{
					ExternalId: "name",
					Type:       framework.AttributeTypeString,
				},
				{
					ExternalId: "age",
					Type:       framework.AttributeTypeInt64,
				},
			},
			expectedObjects: []map[string]any{
				{"name": "John", "age": float64(25)},
				{"name": "Jane", "age": float64(30)},
			},
			expectedHasNext: true,
		},
		"success_last_page": {
			data:     testutil.GenPtr([]byte("name,age\nJohn,25\nJane,30")),
			start:    1,
			pageSize: 3,
			attrConfig: []*framework.AttributeConfig{
				{
					ExternalId: "name",
					Type:       framework.AttributeTypeString,
				},
			},
			expectedObjects: []map[string]any{
				{"name": "John", "age": "25"},
				{"name": "Jane", "age": "30"},
			},
			expectedHasNext: false,
		},
		"success_with_json_field": {
			data: testutil.GenPtr([]byte(`name,aliases
John,"[{""alias"": ""Johnny"", ""primary"": true}]"`)),
			start:    1,
			pageSize: 1,
			attrConfig: []*framework.AttributeConfig{
				{
					ExternalId: "name",
					Type:       framework.AttributeTypeString,
				},
			},
			expectedObjects: []map[string]any{
				{
					"name": "John",
					"aliases": []any{
						map[string]any{"alias": "Johnny", "primary": true},
					},
				},
			},
			expectedHasNext: false,
		},
		"success_numeric_conversion": {
			data:     testutil.GenPtr([]byte("name,score,rating\nJohn,85.5,4")),
			start:    1,
			pageSize: 1,
			attrConfig: []*framework.AttributeConfig{
				{
					ExternalId: "score",
					Type:       framework.AttributeTypeDouble,
				},
				{
					ExternalId: "rating",
					Type:       framework.AttributeTypeInt64,
				},
			},
			expectedObjects: []map[string]any{
				{"name": "John", "score": 85.5, "rating": float64(4)},
			},
			expectedHasNext: false,
		},
		"success_headers_only": {
			data:     testutil.GenPtr([]byte("name,age")),
			start:    1,
			pageSize: 2,
			attrConfig: []*framework.AttributeConfig{
				{
					ExternalId: "name",
					Type:       framework.AttributeTypeString,
				},
			},
			expectedObjects: []map[string]any{},
			expectedHasNext: false,
		},
		"error_empty_data": {
			data:          testutil.GenPtr([]byte("")),
			start:         1,
			pageSize:      2,
			attrConfig:    []*framework.AttributeConfig{},
			expectedError: true,
			errorContains: "no data found in the CSV file",
		},
		"error_invalid_csv": {
			data:          testutil.GenPtr([]byte("name,\"age\nJohn,25")),
			start:         1,
			pageSize:      2,
			attrConfig:    []*framework.AttributeConfig{},
			expectedError: true,
			errorContains: "failed to read CSV data",
		},
		"error_invalid_number": {
			data:     testutil.GenPtr([]byte("name,age\nJohn,not_a_number")),
			start:    1,
			pageSize: 1,
			attrConfig: []*framework.AttributeConfig{
				{
					ExternalId: "age",
					Type:       framework.AttributeTypeInt64,
				},
			},
			expectedError: true,
			errorContains: "failed to convert the value",
		},
		"error_invalid_json": {
			data: testutil.GenPtr([]byte(`name,data
John,"[{invalid json}]"`)),
			start:    1,
			pageSize: 1,
			attrConfig: []*framework.AttributeConfig{
				{
					ExternalId: "name",
					Type:       framework.AttributeTypeString,
				},
			},
			expectedError: true,
			errorContains: "failed to unmarshal the value",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			objects, hasNext, err := s3_adapter.CSVBytesToPage(tt.data, tt.start, tt.pageSize, tt.attrConfig)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}

				if diff := cmp.Diff(objects, tt.expectedObjects); diff != "" {
					t.Errorf("Objects mismatch: %s", diff)
				}

				if hasNext != tt.expectedHasNext {
					t.Errorf("HasNext mismatch: got %v, want %v", hasNext, tt.expectedHasNext)
				}
			}
		})
	}
}

func TestStreamingCSVToPage(t *testing.T) {
	tests := map[string]struct {
		bucket              string
		key                 string
		fileSize            int64
		headers             []string
		start               int64
		pageSize            int64
		attrConfig          []*framework.AttributeConfig
		getObjectStatusCode int
		expectedObjects     []map[string]any
		expectedHasNext     bool
		expectedError       bool
		errorContains       string
	}{
		"success_streaming_basic": {
			bucket:              "test-bucket",
			key:                 "basic-file.csv",
			fileSize:            int64(len(validCSVData)),
			headers:             []string{"Email", "Score"},
			start:               1,
			pageSize:            2,
			getObjectStatusCode: http.StatusOK,
			attrConfig: []*framework.AttributeConfig{
				{
					ExternalId: "Email",
					Type:       framework.AttributeTypeString,
				},
				{
					ExternalId: "Score",
					Type:       framework.AttributeTypeString,
				},
			},
			expectedObjects: []map[string]any{
				{"Email": "1.1", "Score": "e685B8690f9fbce"},
				{"Email": "2.2", "Score": "6EDdBA3a2DFA7De"},
			},
			expectedHasNext: true,
		},
		"success_streaming_with_cursor": {
			bucket:              "test-bucket",
			key:                 "basic-file.csv",
			fileSize:            int64(len(validCSVData)),
			headers:             []string{"Email", "Score"},
			start:               3,
			pageSize:            2,
			getObjectStatusCode: http.StatusOK,
			attrConfig: []*framework.AttributeConfig{
				{
					ExternalId: "Email",
					Type:       framework.AttributeTypeString,
				},
				{
					ExternalId: "Score",
					Type:       framework.AttributeTypeString,
				},
			},
			expectedObjects: []map[string]any{
				{"Email": "3.3", "Score": "b9Da13bedEc47de"},
				{"Email": "4.4", "Score": "710D4dA2FAa96B5"},
			},
			expectedHasNext: true,
		},
		"error_s3_read_failure": {
			bucket:              "test-bucket",
			key:                 "missing-file.csv",
			fileSize:            1000,
			headers:             []string{"Email"},
			start:               1,
			pageSize:            2,
			getObjectStatusCode: http.StatusNotFound,
			attrConfig: []*framework.AttributeConfig{
				{
					ExternalId: "Email",
					Type:       framework.AttributeTypeString,
				},
			},
			expectedError: true,
			errorContains: "unable to read CSV file data",
		},
		"error_permission_denied": {
			bucket:              "test-bucket",
			key:                 "forbidden.csv",
			fileSize:            1000,
			headers:             []string{"Email"},
			start:               1,
			pageSize:            2,
			getObjectStatusCode: http.StatusForbidden,
			attrConfig: []*framework.AttributeConfig{
				{
					ExternalId: "Email",
					Type:       framework.AttributeTypeString,
				},
			},
			expectedError: true,
			errorContains: "unable to read CSV file data",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			awsConfig := mockS3Config(http.StatusOK, tt.getObjectStatusCode)
			handler := &s3_adapter.S3Handler{Client: s3.NewFromConfig(*awsConfig)}

			ctx := context.Background()
			objects, hasNext, err := s3_adapter.StreamingCSVToPage(
				ctx, handler, tt.bucket, tt.key, tt.fileSize,
				tt.headers, tt.start, tt.pageSize, tt.attrConfig,
			)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}

				if diff := cmp.Diff(objects, tt.expectedObjects); diff != "" {
					t.Errorf("Objects mismatch: %s", diff)
				}

				if hasNext != tt.expectedHasNext {
					t.Errorf("HasNext mismatch: got %v, want %v", hasNext, tt.expectedHasNext)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}

				return false
			}())))
}
