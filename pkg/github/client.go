// Copyright 2026 SGNL.ai, Inc.
package github

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

	// EnterpriseSlug is the slug attribute belonging to a GitHub Enterprise entity.
	// Either this field or Organizations must be set (but not both).
	EnterpriseSlug *string

	// Organizations is the list of organizations to query. Either this field or EnterpriseSlug must be set (but not both).
	Organizations []string

	// APIVersion is the version of the GitHub API that is used for constructing REST endpoints.
	APIVersion *string

	// IsEnterpriseCloud is a boolean that indicates whether the deployment is GitHub Enterprise Cloud.
	IsEnterpriseCloud bool

	// Attributes contains the list of attributes to request along with the current request.
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

	// Objects is the list of
	// May be empty.
	Objects []map[string]any

	// NextCursor is the cursor that identifies the first object of the next page.
	// nil if this is the last page in this full sync.
	NextCursor *pagination.CompositeCursor[string]
}
