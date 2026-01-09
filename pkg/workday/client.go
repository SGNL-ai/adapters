// Copyright 2026 SGNL.ai, Inc.
package workday

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Client is a client that allows querying the datasource which contains JSON objects.
type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

// Request is a request to the datasource.
type Request struct {
	// BaseURL is the Base URL of the datasource to query.
	BaseURL string

	// Token is the Bearer API token to authenticate a request.
	Token string

	// APIVersion the API version to use.
	APIVersion string

	// OrganizationID is the ID of the organization in Workday.
	OrganizationID string

	// PageSize is the maximum number of objects to return from the entity.
	PageSize int64

	// Ordered is a flag that indicates if the objects should be ordered.
	Ordered bool

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *pagination.CompositeCursor[int64]

	// EntityConfig contains the attributes that will be used to build the wql query.
	EntityConfig *framework.EntityConfig

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds int
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
