// Copyright 2025 SGNL.ai, Inc.
package awss3

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/pagination"
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
	Cursor *pagination.CompositeCursor[int64]

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
	NextCursor *pagination.CompositeCursor[int64]
}
