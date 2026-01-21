// Copyright 2026 SGNL.ai, Inc.

package pagerduty

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Client is a client that allows querying the PagerDuty datasource which contains JSON objects.
type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

// Request is a request to PagerDuty.
type Request struct {
	// BaseURL is the base URL for PagerDuty. Should always be "https://api.pagerduty.com".
	BaseURL string

	// AdditionalQueryParameters are any additional query parameters to be added to the request.
	// For example, {"users": {"query": "..."}}. The top level map key must match the entity external ID
	// for the query parameters to be added.
	AdditionalQueryParameters map[string]map[string][]string

	// Token is the API token to authenticate a request. For example, "Token token=y_NbAkKc66ryYTWUXYEu".
	Token string

	// PageSize is the maximum number of objects to return per page from the API call.
	// This is used as the "limit" parameter in the PagerDuty API.
	PageSize int64

	// EntityExternalID is the external ID of the entity.
	// The external ID should match the API's resource name, e.g. "users", "teams", "schedules", etc.,
	// with the only exception being "members" for team members.
	EntityExternalID string

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *pagination.CompositeCursor[int64]

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds int
}

// Response is a parsed response returned from PagerDuty.
type Response struct {
	// StatusCode is an HTTP status code.
	StatusCode int

	// RetryAfterHeader is the Retry-After response HTTP header, if set.
	RetryAfterHeader string

	// Objects is the list of parsed entity objects returned from PagerDuty.
	// May be empty.
	Objects []map[string]any

	// NextCursor is the cursor that identifies the first object of the next page.
	// nil if this is the last page in this full sync.
	NextCursor *pagination.CompositeCursor[int64]
}
