// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package awss3_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
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
		errorContains        string // Expects errors from the mock middleware directly or from "failed to convert response"
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
			errorContains:        "not found: The specified key does not exist", // Error from mock middleware
		},
		"permission_denied": {
			bucket:               "test-bucket",
			key:                  "forbidden.csv",
			headObjectStatusCode: http.StatusForbidden,
			expectedError:        true,
			errorContains:        "AccessDenied: Access Denied", // Error from mock middleware
		},
		"moved_permanently": {
			bucket:               "test-bucket",
			key:                  "moved.csv",
			headObjectStatusCode: http.StatusMovedPermanently,
			expectedError:        true,
			errorContains:        "permanent redirect", // Error from mock middleware
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			awsConfig := mockS3Config(tt.headObjectStatusCode, http.StatusOK) // getObjectStatusCode doesn't matter here
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
				// Check response fields if needed, e.g., response.ContentLength
				// The mock for http.StatusOK now returns len(validCSVData)
				expectedLen := int64(len(validCSVData))
				if response.ContentLength == nil || *response.ContentLength != expectedLen {
					t.Errorf("Expected ContentLength to be %d, got: %v", expectedLen, response.ContentLength)
				}
			}
		})
	}
}

func TestS3Handler_GetObjectStream(t *testing.T) {
	tests := map[string]struct {
		bucket              string
		key                 string
		rangeHeader         *string
		getObjectStatusCode int
		expectedError       bool
		errorContains       string
		expectedContent     string // For simple validation of the stream
	}{
		"success_no_range": {
			bucket:              "test-bucket",
			key:                 "test-file.csv",
			rangeHeader:         nil,
			getObjectStatusCode: http.StatusOK, // Mock returns validCSVData
			expectedError:       false,
			expectedContent:     validCSVData,
		},
		"success_empty_string_range": {
			bucket:              "test-bucket",
			key:                 "test-file.csv",
			rangeHeader:         aws.String(""), // s3.go treats this as no range
			getObjectStatusCode: http.StatusOK,
			expectedError:       false,
			expectedContent:     validCSVData,
		},
		"success_with_range": {
			bucket:              "test-bucket",
			key:                 "test-file.csv",
			rangeHeader:         aws.String("bytes=0-100"),
			getObjectStatusCode: http.StatusOK, // Mock doesn't change content based on range, but S3 client would
			expectedError:       false,
			expectedContent:     validCSVData, // Mock will still return full validCSVData
		},
		"error_not_found": {
			bucket:              "test-bucket",
			key:                 "nonexistent.csv",
			rangeHeader:         nil,
			getObjectStatusCode: http.StatusNotFound,
			expectedError:       true,
			errorContains:       "no such key: The specified key does not exist", // Error from mock middleware
		},
		"error_permission_denied": {
			bucket:              "test-bucket",
			key:                 "forbidden.csv",
			rangeHeader:         nil,
			getObjectStatusCode: http.StatusForbidden,
			expectedError:       true,
			errorContains:       "access denied: Access Denied", // Error from mock middleware (note lowercase 'access denied')
		},
		"error_range_not_satisfiable": {
			bucket:              "test-bucket",
			key:                 "test-file.csv",
			rangeHeader:         aws.String("bytes=1000000-1000001"), // A range likely outside the mock data size
			getObjectStatusCode: http.StatusRequestedRangeNotSatisfiable,
			expectedError:       true,
			// This error often comes from the SDK in a way that httpResponseFromError might not parse perfectly,
			// leading to the "failed to convert response" wrapper.
			errorContains: "failed to convert response",
		},
		"error_get_object_nil_response_simulation": { // Theoretical case not easily triggered by mock
			bucket:              "test-bucket",
			key:                 "nil-response-sim.csv",
			getObjectStatusCode: -999, // Custom code to indicate mock should return nil response (would need mock update)
			expectedError:       true,
			errorContains:       "failed to download file range: response is nil",
			// Note: Current mockS3Middleware does not support simulating a nil response with a success status.
			// This test case is conceptual for s3.go's nil check.
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Skip conceptual test if not implemented in mock
			if tt.getObjectStatusCode == -999 && name == "error_get_object_nil_response_simulation" {
				t.Skip("Skipping nil response simulation test as mock does not support it.")
			}

			awsConfig := mockS3Config(http.StatusOK, tt.getObjectStatusCode) // headObjectStatusCode doesn't matter here
			handler := &s3_adapter.S3Handler{Client: s3.NewFromConfig(*awsConfig)}

			ctx := context.Background()
			response, err := handler.GetObjectStream(ctx, tt.bucket, tt.key, tt.rangeHeader)

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
					t.Fatalf("Expected non-nil response")
				}
				if response.Body == nil {
					t.Fatalf("Expected non-nil response.Body")
				}
				defer response.Body.Close()
				bodyBytes, readErr := io.ReadAll(response.Body)
				if readErr != nil {
					t.Fatalf("Failed to read response body: %v", readErr)
				}
				if diff := cmp.Diff(tt.expectedContent, string(bodyBytes)); diff != "" {
					t.Errorf("Content mismatch (-want +got):\n%s", diff)
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
		// Mock setup for simulating nil ContentLength on HeadObjectOutput for a specific test.
		// If nilContentLengthTest is true, the mock should be adjusted (not done by this code).
		nilContentLengthTest bool
		expectedError        bool
		errorContains        string
		expectedSize         int64
	}{
		"successful_size_retrieval": {
			bucket:               "test-bucket",
			key:                  "test.csv",
			headObjectStatusCode: http.StatusOK,
			expectedError:        false,
			expectedSize:         int64(len(validCSVData)),
		},
		"file_not_found": {
			bucket:               "test-bucket",
			key:                  "missing.csv",
			headObjectStatusCode: http.StatusNotFound,
			expectedError:        true,
			errorContains:        "not found: The specified key does not exist", // Error from mock middleware
		},
		"permission_denied": {
			bucket:               "test-bucket",
			key:                  "forbidden.csv",
			headObjectStatusCode: http.StatusForbidden,
			expectedError:        true,
			errorContains:        "AccessDenied: Access Denied", // Error from mock middleware
		},
		"moved_permanently": {
			bucket:               "test-bucket",
			key:                  "moved.csv",
			headObjectStatusCode: http.StatusMovedPermanently,
			expectedError:        true,
			errorContains:        "permanent redirect", // Error from mock middleware
		},
		"error_unable_to_determine_file_size": {
			bucket:               "test-bucket",
			key:                  "nil-content-length.csv",
			headObjectStatusCode: http.StatusOK, // HeadObject succeeds
			nilContentLengthTest: true,          // Indicates this test needs special mock behavior
			expectedError:        true,
			errorContains:        "unable to determine file size",
			// Note: Current mockS3Middleware always returns a ContentLength on 200 OK.
			// This test case requires mockS3Middleware to be modified to return a HeadObjectOutput
			// with a nil ContentLength for this specific key/scenario.
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.nilContentLengthTest {
				// This is a conceptual reminder that the mock would need adjustment for this path.
				// For now, we can't directly trigger this path with the existing mock without changes to it.
				// One way to test would be to make mockS3Config return a specially crafted HeadObjectOutput.
				t.Skip("Skipping test for 'unable to determine file size' as mock needs modification to support nil ContentLength on success.")
			}

			awsConfig := mockS3Config(tt.headObjectStatusCode, http.StatusOK) // getObjectStatusCode doesn't matter
			handler := &s3_adapter.S3Handler{Client: s3.NewFromConfig(*awsConfig)}

			ctx := context.Background()
			size, err := handler.GetFileSize(ctx, tt.bucket, tt.key)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}

				if size != 0 { // As per s3.go, error returns 0 size
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
