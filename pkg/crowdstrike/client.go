// Copyright 2026 SGNL.ai, Inc.
package crowdstrike

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Client is a client that allows querying the datasource which
// contains JSON objects.
type Client interface {
	// GetPage returns a page of JSON objects from the datasource for the
	// requested entity.
	// Returns a (possibly empty) list of JSON objects, each object being
	// unmarshaled into a map by Golang's JSON unmarshaler.
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

// Request is a request to the datasource.
type Request struct {
	// BaseURL is the Base URL of the datasource to query.
	BaseURL string

	// Token is the Authorization token to use to authentication with the datasource.
	Token string

	// PageSize is the maximum number of objects to return from the entity.
	PageSize int64

	// EntityExternalID is the external ID of the entity.
	// The external ID should match the API's resource name.
	EntityExternalID string

	// A Falcon Query Language filter applicable for REST API based entities.
	// See more at https://falconpy.io/Usage/Falcon-Query-Language.html#operators
	Filter *string

	// GraphQLCursor identifies the first object of the page to return, as returned by
	// the last request for the entity. This field is used to paginate entities from GraphQL APIs.
	// Optional. If not set, return the first page for this entity.
	GraphQLCursor *pagination.CompositeCursor[string]

	// RESTCursor identifies the first object of the page to return, as returned by
	// the last request for the entity. This field is used to paginate entities from REST APIs.
	// Optional. If not set, return the first page for this entity.
	RESTCursor *pagination.CompositeCursor[string]

	// EntityConfig contains entity metadata and a list of attributes to request along with the current request.
	EntityConfig *framework.EntityConfig

	// Ordered is a boolean that indicates whether the results should be ordered.
	Ordered bool

	Config *Config

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds int
}

type Response struct {
	// StatusCode is an HTTP status code.
	StatusCode int

	// RetryAfterHeader is the Retry-After response HTTP header, if set.
	RetryAfterHeader string

	// Objects is the list of parsed entity objects returned from the datasource.	// May be empty.
	Objects []map[string]any

	// NextGraphQLCursor is the cursor that identifies the first object of the next
	// page for GraphQL APIs.
	// May be empty.
	NextGraphQLCursor *pagination.CompositeCursor[string]

	// NextRESTCursor is the cursor that identifies the first object of the next
	// page for REST APIs.
	// May be empty.
	NextRESTCursor *pagination.CompositeCursor[string]
}
