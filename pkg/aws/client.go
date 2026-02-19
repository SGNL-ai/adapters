// Copyright 2026 SGNL.ai, Inc.

package aws

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

	// PageSize is the maximum number of objects to return from the entity.
	MaxItems int32

	// EntityExternalID is the external ID of the entity.
	// The external ID should match the API's resource name.
	EntityExternalID string

	// AccountIDRequested is a boolean that indicates whether the account ID is requested.
	AccountIDRequested bool

	// EntityConfig is a map containing the config required for each entity associated with this
	EntityConfig map[string]*EntityConfig

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *pagination.CompositeCursor[string]

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds int

	// Ordered is a boolean that indicates whether the results should be ordered.
	Ordered bool

	// ResourceAccountRoles is a list of roleARNs.
	ResourceAccountRoles []string
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
	NextCursor *pagination.CompositeCursor[string]
}
