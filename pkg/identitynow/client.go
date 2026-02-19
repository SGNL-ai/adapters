// Copyright 2026 SGNL.ai, Inc.

package identitynow

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Client is a client that allows querying the IdentityNow datasource which contains JSON objects.
type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

// Request is a request to IdentityNow.
type Request struct {
	// BaseURL is the Base URL of the datasource to query.
	// For example, https://{tenant}.api.identitynow.com.
	BaseURL string

	// Token is the API token to authenticate a request. IdentityNow supports OAuth2 Client Credential auth
	// so this must be in the form "Bearer eyJhbG[...]1LQ".
	Token string

	// PageSize is the maximum number of objects to return from the entity.
	PageSize int64

	// EntityExternalID is the external ID of the entity.
	// The external ID should match the API's resource name.
	EntityExternalID string

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *pagination.CompositeCursor[int64]

	// APIVersion the API version to use.
	APIVersion string

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds int

	// Filter contains the optional filter to apply to the current request.
	// It's applied using the `filters` query parameter.
	Filter *string

	// Sorters contains the optional sorters param to apply to the current request.
	// It contains a comma separated value of field names by which the results are sorted.
	//
	// Example: sorters=type,-modified
	// Results are sorted primarily by type in ascending order, and secondarily by modified date in descending order.
	Sorters *string
}

// Response is a response returned by the datasource.
type Response struct {
	// StatusCode is an HTTP status code.
	StatusCode int

	// RetryAfterHeader is the Retry-After response HTTP header, if set.
	RetryAfterHeader string

	// Objects is the list of
	// May be empty.
	Objects []map[string]any

	// NextCursor is the cursor that identifies the first object of the next page.
	// nil if this is the last page in this full sync.
	NextCursor *pagination.CompositeCursor[int64]
}
