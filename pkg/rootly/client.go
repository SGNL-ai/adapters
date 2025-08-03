// Copyright 2025 SGNL.ai, Inc.
package rootly

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
)

// Client is a client that allows querying the Rootly datasource which contains JSON objects.
type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

// Request is a request to Rootly.
type Request struct {
	// BaseURL is the Base URL of the datasource to query.
	BaseURL string

	// HTTPAuthorization is the HTTP authorization header value to authenticate a request.
	// This will be provided in the form "Bearer <token>".
	HTTPAuthorization string

	// PageSize is the maximum number of objects to return from the entity.
	PageSize int64

	// EntityExternalName is the external name of the entity.
	// The external name should match the API's resource name.
	EntityExternalName string

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *string

	// Filter contains the optional filter to apply to the current request.
	Filter string

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds int
}

// Response is a response returned by the datasource.
type Response struct {
	// Objects contains the list of objects returned by the datasource.
	Objects []map[string]any

	// NextCursor identifies the first object of the next page to return.
	// nil if the current page is the last page for the entity.
	NextCursor *string
}