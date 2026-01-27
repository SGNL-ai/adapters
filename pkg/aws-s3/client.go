// Copyright 2026 SGNL.ai, Inc.

package awss3

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// Client is a client that allows querying the datasource which contains JSON objects.
type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

type Auth struct {
	// AccessKey is the access key to authenticate with the AWS.
	AccessKey string

	// SecretKey is the secret key to authenticate with the AWS.
	SecretKey string

	// Region is the AWS region to query.
	Region string
}

// S3Cursor contains pagination state for S3 CSV files.
// Headers are cached to avoid re-fetching on subsequent pages.
type S3Cursor struct {
	// Cursor is the byte position offset in the file.
	Cursor *int64 `json:"cursor,omitempty"`

	// Headers contains the parsed CSV headers from the first page.
	// Cached to avoid re-fetching headers on subsequent pages.
	Headers []string `json:"headers,omitempty"`
}

// UnmarshalS3Cursor unmarshals the cursor from a base64 encoded JSON string.
// Returns nil cursor if the input is empty.
func UnmarshalS3Cursor(cursor string) (*S3Cursor, *framework.Error) {
	if cursor == "" {
		return nil, nil
	}

	cursorBytes, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to decode base64 cursor: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	unmarshaledCursor := &S3Cursor{}

	unmarshalErr := json.Unmarshal(cursorBytes, unmarshaledCursor)
	if unmarshalErr != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal JSON cursor: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	return unmarshaledCursor, nil
}

// MarshalS3Cursor marshals the cursor into a base64 encoded JSON string.
func MarshalS3Cursor(cursor *S3Cursor) (string, *framework.Error) {
	if cursor == nil {
		return "", nil
	}

	nextCursorBytes, marshalErr := json.Marshal(cursor)
	if marshalErr != nil {
		return "", &framework.Error{
			Message: fmt.Sprintf("Failed to marshal cursor into JSON: %v.", marshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return base64.StdEncoding.EncodeToString(nextCursorBytes), nil
}

// Request is a request to the datasource.
type Request struct {
	Auth

	// Bucket is the AWS S3 bucket containing the files with entity data.
	Bucket string

	// PathPrefix is the prefix of the path containing the files with entity data.
	PathPrefix string

	// FileType is the extension of the files containing the entity data.
	FileType string

	// PageSize is the maximum number of objects to return from the entity.
	PageSize int64

	// EntityExternalID is the external ID of the entity.
	// The external ID should match the file name.
	EntityExternalID string

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *S3Cursor

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds int

	// Ordered is a boolean that indicates whether the results should be ordered.
	Ordered bool

	// AttributeConfig is the list of attributes requested by the datasource.
	AttributeConfig []*framework.AttributeConfig
}

// Response is a response returned by the datasource.
type Response struct {
	// StatusCode is an HTTP status code.
	StatusCode int

	// RetryAfterHeader is the Retry-After response HTTP header, if set.
	RetryAfterHeader string

	// Objects is the list of items returned by the datasource.
	// May be empty.
	Objects []map[string]any

	// NextCursor is the cursor that identifies the first object of the next page.
	// nil if this is the last page in this full sync.
	NextCursor *S3Cursor
}
