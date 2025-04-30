// Copyright 2025 SGNL.ai, Inc.
package mysql

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
)

type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

// Request is a request to a MySQL database.
type Request struct {
	// BaseURL is the Base URL of the datasource to query.
	BaseURL string

	// Username is the user name used to authenticate with the MySQL instance.
	Username string

	// Password is the password used to authenticate with the MySQL instance.
	Password string

	// PageSize is the maximum number of objects to return from the entity.
	PageSize int64

	// EntityConfig contains the entity configuration for the current request.
	EntityConfig framework.EntityConfig

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *int64

	// MySQL database to connect to.
	Database string

	// UniqueAttributeExternalID is used to specify the unique ID that should be used when ordering results from
	// the specified table.
	UniqueAttributeExternalID string
}

// Response is a response returned by the datasource.
type Response struct {
	// StatusCode is an HTTP status code.
	StatusCode int

	// RetryAfterHeader is the Retry-After response HTTP header, if set.
	RetryAfterHeader string

	// Objects is the list of objects returned from the datasource.
	// May be empty.
	Objects []map[string]any

	// NextCursor is the cursor that identifies the first object of the next page.
	// nil if this is the last page in this full sync.
	NextCursor *int64
}
