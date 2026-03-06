// Copyright 2026 SGNL.ai, Inc.

package victorops

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Client is a client that allows querying the VictorOps (Splunk On-Call) datasource.
type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

// Request is a request to VictorOps.
type Request struct {
	// BaseURL is the base URL of the VictorOps API. For example, "https://api.victorops.com".
	BaseURL string

	// APIId is the VictorOps API ID used for authentication via the X-VO-Api-Id header.
	// This value is sourced from request.Auth.Basic.Username.
	APIId string

	// APIKey is the VictorOps API key used for authentication via the X-VO-Api-Key header.
	// This value is sourced from request.Auth.Basic.Password.
	APIKey string

	// PageSize is the maximum number of objects to return per page from the API call.
	PageSize int64

	// EntityExternalID is the external ID of the entity.
	EntityExternalID string

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *pagination.CompositeCursor[int64]

	// QueryParameters is an optional query string to append to the API request URL.
	// For example, "currentPhase=RESOLVED&startedAfter=2024-01-01T00:00Z".
	QueryParameters string

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	RequestTimeoutSeconds int
}

// Response is a parsed response returned from VictorOps.
type Response struct {
	// StatusCode is an HTTP status code.
	StatusCode int

	// RetryAfterHeader is the Retry-After response HTTP header, if set.
	RetryAfterHeader string

	// Objects is the list of parsed entity objects returned from VictorOps.
	Objects []map[string]any

	// NextCursor is the cursor that identifies the first object of the next page.
	// nil if this is the last page in this full sync.
	NextCursor *pagination.CompositeCursor[int64]
}
