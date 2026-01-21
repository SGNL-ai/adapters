// Copyright 2026 SGNL.ai, Inc.

package azuread

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

	// PageSize is the maximum number of objects to return from the entity.
	PageSize int64

	// EntityExternalID is the external ID of the entity.
	// The external ID should match the API's resource name.
	EntityExternalID string

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *pagination.CompositeCursor[string]

	// Skip is the number of objects to skip from the beginning of the entity.
	// This is applicable to PIM entities like RoleAssignmentScheduleRequest and GroupAssignmentScheduleRequest.
	Skip int64

	// Filter contains the optional filter to apply to the current request.
	Filter *string

	// ParentFilter contains the optional filter to apply when retrieving parent objects for a member entity.
	// This will only be set if the current entity has a parent and ApplyFiltersToMembers is enabled.
	ParentFilter *string

	// Attributes contains the list of attributes to request along with the current request.
	Attributes []*framework.AttributeConfig

	// APIVersion the API version to use.
	APIVersion string

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds int

	// UseAdvancedFilters is a flag that indicates whether advanced filters should be used.
	// Set $count=true in the query to enable advanced filters.
	// See https://learn.microsoft.com/en-us/graph/aad-advanced-queries for more information.
	UseAdvancedFilters bool

	// AdvancedFilterMemberExternalID is the external ID of the member entity to apply advanced filters to.
	AdvancedFilterMemberExternalID *string
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
	NextCursor *pagination.CompositeCursor[string]
}
