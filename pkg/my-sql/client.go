// Copyright 2025 SGNL.ai, Inc.
package mysql

import (
	"context"
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

// Request is a request to a MySQL database.
type Request struct {
	// BaseURL is the Base URL of the datasource to query.
	BaseURL string `json:"baseURL"`

	// Username is the user name used to authenticate with the MySQL instance.
	Username string `json:"username"`

	// Password is the password used to authenticate with the MySQL instance.
	Password string `json:"password"`

	// PageSize is the maximum number of objects to return from the entity.
	PageSize int64 `json:"pageSize,omitempty"`

	// EntityConfig contains the entity configuration for the current request.
	EntityConfig framework.EntityConfig

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *int64 `json:"cursor,omitempty"`

	// MySQL database to connect to.
	Database string `json:"database"`

	// UniqueAttributeExternalID is used to specify the unique ID that should be used when ordering results from
	// the specified table.
	UniqueAttributeExternalID string `json:"uniqueAttributeExternalID"`
}

// DatasourceName to connect to.
func (r *Request) DatasourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s",
		r.Username, r.Password, r.BaseURL, r.Database,
	)
}

// SimpleSQLValidation performs simple validation on the table name to prevent SQL Ingestion attacks,
// since we can't use table names in prepared queries which leaves us vulnerable.
func (r *Request) SimpleSQLValidation() *framework.Error {
	if valid := validSQLTableName.MatchString(r.EntityExternalID); !valid {
		return &framework.Error{
			Message: "SQL table name validation failed: unsupported characters found, or its len is < 1 or > 128.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	return nil
}

// Response is a response returned by the datasource.
type Response struct {
	// StatusCode is an HTTP status code.
	StatusCode int `json:"statusCode"`

	// RetryAfterHeader is the Retry-After response HTTP header, if set.
	RetryAfterHeader string `json:"retryAfterHeader"`

	// Objects is the list of objects returned from the datasource.
	// May be empty.
	Objects []map[string]any `json:"objects"`

	// NextCursor is the cursor that identifies the first object of the next page.
	// nil if this is the last page in this full sync.
	NextCursor *int64 `json:"nextCursor,omitempty"`
}
