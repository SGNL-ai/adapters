// Copyright 2025 SGNL.ai, Inc.
package okta

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Client is a client that allows querying the Okta datasource which contains JSON objects.
type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

// Request is a request to Okta.
type Request struct {
	// BaseURL is the Base URL of the datasource to query.
	BaseURL string

	// Token is the API token to authenticate a request. Okta supports both API tokens and OAuth2 Client Credential
	// auth, so this may either be in the form "SSWS XXXX" (for API tokens) or "Bearer eyJhbG[...]1LQ" (for OAuth2).
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

	// APIVersion the API version to use.
	APIVersion string

	// Filter is the Okta Filter syntax to apply to requests for Users and/or Groups
	Filter string

	// Search is the Okta Search syntax to apply to requests for Users and/or Groups
	Search string

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
