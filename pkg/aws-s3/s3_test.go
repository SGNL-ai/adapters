// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package awss3_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/go-cmp/cmp"
	s3_adapter "github.com/sgnl-ai/adapters/pkg/aws-s3"
)

func TestS3Handler_FileExists(t *testing.T) {
	tests := map[string]struct {
		bucket               string
		key                  string
		headObjectStatusCode int
		expectedError        bool
		errorContains        string
	}{
		"file_exists_success": {
			bucket:               "test-bucket",
			key:                  "test-file.csv",
			headObjectStatusCode: http.StatusOK,
			expectedError:        false,
		},
		"file_not_found": {
			bucket:               "test-bucket",
			key:                  "nonexistent.csv",
			headObjectStatusCode: http.StatusNotFound,
			expectedError:        true,
			errorContains:        "failed to check if the file exists",
		},
		"permission_denied": {
			bucket:               "test-bucket",
			key:                  "forbidden.csv",
			headObjectStatusCode: http.StatusForbidden,
			expectedError:        true,
			errorContains:        "failed to check if the file exists",
		},
		"moved_permanently": {
			bucket:               "test-bucket",
			key:                  "moved.csv",
			headObjectStatusCode: http.StatusMovedPermanently,
			expectedError:        true,
			errorContains:        "failed to check if the file exists",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			awsConfig := mockS3Config(tt.headObjectStatusCode, http.StatusOK)
			handler := &s3_adapter.S3Handler{Client: s3.NewFromConfig(*awsConfig)}

			ctx := context.Background()
			response, err := handler.FileExists(ctx, tt.bucket, tt.key)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}

				if response != nil {
					t.Errorf("Expected nil response on error, got: %v", response)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}

				if response == nil {
					t.Errorf("Expected response, got nil")
				}
			}
		})
	}
}

func TestS3Handler_GetFile(t *testing.T) {
	tests := map[string]struct {
		bucket               string
		key                  string
		headObjectStatusCode int
		getObjectStatusCode  int
		expectedError        bool
		errorContains        string
		expectedContent      []byte
	}{
		"successful_file_read": {
			bucket:               "test-bucket",
			key:                  "test.csv",
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  http.StatusOK,
			expectedError:        false,
			expectedContent:      []byte(validCSVData),
		},
		"file_not_found_head": {
			bucket:               "test-bucket",
			key:                  "missing.csv",
			headObjectStatusCode: http.StatusNotFound,
			expectedError:        true,
			errorContains:        "failed to check if the file exists",
		},
		"permission_denied_get": {
			bucket:               "test-bucket",
			key:                  "forbidden.csv",
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  http.StatusForbidden,
			expectedError:        true,
			errorContains:        "failed to convert response",
		},
		"get_object_not_found": {
			bucket:               "test-bucket",
			key:                  "missing-get.csv",
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  http.StatusNotFound,
			expectedError:        true,
			errorContains:        "failed to convert response",
		},
		"get_object_moved": {
			bucket:               "test-bucket",
			key:                  "moved.csv",
			headObjectStatusCode: http.StatusOK,
			getObjectStatusCode:  http.StatusMovedPermanently,
			expectedError:        true,
			errorContains:        "failed to convert response",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			awsConfig := mockS3Config(tt.headObjectStatusCode, tt.getObjectStatusCode)
			handler := &s3_adapter.S3Handler{Client: s3.NewFromConfig(*awsConfig)}

			ctx := context.Background()
			result, err := handler.GetFile(ctx, tt.bucket, tt.key)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}

				if result != nil {
					t.Errorf("Expected nil result on error, got: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}

				if result == nil {
					t.Errorf("Expected result, got nil")
				} else if diff := cmp.Diff(*result, tt.expectedContent); diff != "" {
					t.Errorf("Content mismatch: %s", diff)
				}
			}
		})
	}
}

func TestS3Handler_GetFileRange(t *testing.T) {
	tests := map[string]struct {
		bucket              string
		key                 string
		startByte           int64
		endByte             int64
		getObjectStatusCode int
		expectedError       bool
		errorContains       string
		shouldHaveBOM       bool
	}{
		"successful_range_request_from_start": {
			bucket:              "test-bucket",
			key:                 "test.csv",
			startByte:           0,
			endByte:             100,
			getObjectStatusCode: http.StatusOK,
			expectedError:       false,
			shouldHaveBOM:       false,
		},
		"range_request_from_middle": {
			bucket:              "test-bucket",
			key:                 "test.csv",
			startByte:           100,
			endByte:             200,
			getObjectStatusCode: http.StatusOK,
			expectedError:       false,
			shouldHaveBOM:       false,
		},
		"range_not_satisfiable": {
			bucket:              "test-bucket",
			key:                 "test.csv",
			startByte:           0,
			endByte:             100,
			getObjectStatusCode: http.StatusRequestedRangeNotSatisfiable,
			expectedError:       true,
			errorContains:       "failed to convert response",
		},
		"permission_denied": {
			bucket:              "test-bucket",
			key:                 "forbidden.csv",
			startByte:           0,
			endByte:             100,
			getObjectStatusCode: http.StatusForbidden,
			expectedError:       true,
			errorContains:       "failed to convert response",
		},
		"file_not_found": {
			bucket:              "test-bucket",
			key:                 "missing.csv",
			startByte:           0,
			endByte:             100,
			getObjectStatusCode: http.StatusNotFound,
			expectedError:       true,
			errorContains:       "failed to convert response",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			awsConfig := mockS3Config(http.StatusOK, tt.getObjectStatusCode)
			handler := &s3_adapter.S3Handler{Client: s3.NewFromConfig(*awsConfig)}

			ctx := context.Background()
			result, err := handler.GetFileRange(ctx, tt.bucket, tt.key, tt.startByte, tt.endByte)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}

				if result != nil {
					t.Errorf("Expected nil result on error, got: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}

				if result == nil {
					t.Errorf("Expected result, got nil")
				}
			}
		})
	}
}

func TestS3Handler_GetFileSize(t *testing.T) {
	tests := map[string]struct {
		bucket               string
		key                  string
		headObjectStatusCode int
		expectedError        bool
		errorContains        string
		expectedSize         int64
	}{
		"successful_size_retrieval": {
			bucket:               "test-bucket",
			key:                  "test.csv",
			headObjectStatusCode: http.StatusOK,
			expectedError:        false,
			expectedSize:         1234,
		},
		"file_not_found": {
			bucket:               "test-bucket",
			key:                  "missing.csv",
			headObjectStatusCode: http.StatusNotFound,
			expectedError:        true,
			errorContains:        "failed to check if the file exists",
		},
		"permission_denied": {
			bucket:               "test-bucket",
			key:                  "forbidden.csv",
			headObjectStatusCode: http.StatusForbidden,
			expectedError:        true,
			errorContains:        "failed to check if the file exists",
		},
		"moved_permanently": {
			bucket:               "test-bucket",
			key:                  "moved.csv",
			headObjectStatusCode: http.StatusMovedPermanently,
			expectedError:        true,
			errorContains:        "failed to check if the file exists",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			awsConfig := mockS3Config(tt.headObjectStatusCode, http.StatusOK)
			handler := &s3_adapter.S3Handler{Client: s3.NewFromConfig(*awsConfig)}

			ctx := context.Background()
			size, err := handler.GetFileSize(ctx, tt.bucket, tt.key)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}

				if size != 0 {
					t.Errorf("Expected zero size on error, got: %d", size)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}

				if size != tt.expectedSize {
					t.Errorf("Expected size %d, got: %d", tt.expectedSize, size)
				}
			}
		})
	}
}

func TestS3Handler_GetHeaderChunk(t *testing.T) {
	tests := map[string]struct {
		bucket              string
		key                 string
		getObjectStatusCode int
		expectedError       bool
		errorContains       string
	}{
		"successful_header_read": {
			bucket:              "test-bucket",
			key:                 "test.csv",
			getObjectStatusCode: http.StatusOK,
			expectedError:       false,
		},
		"permission_denied": {
			bucket:              "test-bucket",
			key:                 "forbidden.csv",
			getObjectStatusCode: http.StatusForbidden,
			expectedError:       true,
			errorContains:       "failed to convert response",
		},
		"file_not_found": {
			bucket:              "test-bucket",
			key:                 "missing.csv",
			getObjectStatusCode: http.StatusNotFound,
			expectedError:       true,
			errorContains:       "failed to convert response",
		},
		"range_not_satisfiable": {
			bucket:              "test-bucket",
			key:                 "tiny.csv",
			getObjectStatusCode: http.StatusRequestedRangeNotSatisfiable,
			expectedError:       true,
			errorContains:       "failed to convert response",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			awsConfig := mockS3Config(http.StatusOK, tt.getObjectStatusCode)
			handler := &s3_adapter.S3Handler{Client: s3.NewFromConfig(*awsConfig)}

			ctx := context.Background()
			result, err := handler.GetHeaderChunk(ctx, tt.bucket, tt.key)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}

				if result != nil {
					t.Errorf("Expected nil result on error, got: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}

				if result == nil {
					t.Errorf("Expected result, got nil")
				}
			}
		})
	}
}
