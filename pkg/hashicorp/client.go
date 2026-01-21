// Copyright 2026 SGNL.ai, Inc.

package hashicorp

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Client is a client that allows querying the HashiCorp datasource which contains JSON objects.
type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

type Auth struct {
	Username     string
	Password     string
	AuthMethodID string
}

// Request is a request to HashiCorp.
type Request struct {
	// BaseURL is the Base URL of the datasource to query.
	BaseURL string

	// Auth is the authentication information to authenticate a request.
	Auth

	// authToken is the authentication token to use for the request retrieved from the auth-method.
	authToken string

	// PageSize is the maximum number of objects to return from the entity.
	PageSize int64

	// EntityExternalID is the external ID of the entity.
	// The external ID should match the API's resource name.
	EntityExternalID string

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *pagination.CompositeCursor[string]

	// Filter contains the optional filter to apply to the current request.
	Filter *string

	// Attributes contains the list of attributes to request along with the current request.
	Attributes []*framework.AttributeConfig

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds int

	// EntityConfig is the configuration for the each entity.
	EntityConfig map[string]EntityConfig

	AdditionalParams map[string]string
}

// Response is a response returned by the datasource.
type Response struct {
	// StatusCode is an HTTP status code.
	StatusCode int

	// Objects is the list of
	// May be empty.
	Objects []map[string]any

	// NextCursor is the cursor that identifies the first object of the next page.
	// nil if this is the last page in this full sync.
	NextCursor *string

	// RetryAfter is the number of seconds to wait before retrying the request.
	RetryAfterHeader string
}
