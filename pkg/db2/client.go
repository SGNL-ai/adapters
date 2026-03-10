// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/condexpr"
)

// Client defines the interface for querying a DB2 datasource.
type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

// Rows represents the result set from a SQL query.
type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Err() error
	Columns() ([]string, error)
}

// Request is a request to a DB2 database.
type Request struct {
	// BaseURL is the Base URL of the datasource to query.
	BaseURL string `json:"baseURL"`

	// Username is the user name used to authenticate with the DB2 instance.
	Username string `json:"username"`

	// Password is the password used to authenticate with the DB2 instance.
	Password string `json:"password"`

	// PageSize is the maximum number of objects to return from the entity.
	PageSize int64 `json:"pageSize,omitempty"`

	// EntityConfig contains the entity configuration for the current request.
	EntityConfig framework.EntityConfig

	// A filter to apply to the DB2 request when pulling data for the current entity.
	Filter *condexpr.Condition

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *string `json:"cursor,omitempty"`

	// DB2 database to connect to.
	Database string `json:"database"`

	// Schema name for table queries (optional)
	Schema string `json:"schema,omitempty"`

	// UniqueAttributeExternalID is used to specify the unique ID that should be used when ordering results from
	// the specified table.
	UniqueAttributeExternalID string `json:"uniqueAttributeExternalID"`

	// UniqueKeyColumns contains the columns that comprise the unique key for composite ID generation
	UniqueKeyColumns []string `json:"uniqueKeyColumns,omitempty"`

	// ConfigStruct contains the parsed configuration for SSL and other advanced options
	ConfigStruct interface{} `json:"configStruct,omitempty"`
}

// Response is the response returned from a DB2 query.
type Response struct {
	// Objects is a list of objects returned from the query.
	Objects []map[string]interface{} `json:"objects,omitempty"`

	// NextCursor is the cursor to use for the next page of results.
	NextCursor *string `json:"nextCursor,omitempty"`

	// TotalCount is the total number of objects matching the query.
	TotalCount int64 `json:"totalCount,omitempty"`

	// StatusCode is the HTTP status code returned from the DB2 query.
	StatusCode int `json:"statusCode,omitempty"`
}
