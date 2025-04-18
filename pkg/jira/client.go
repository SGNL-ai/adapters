// Copyright 2025 SGNL.ai, Inc.
package jira

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Client is a client that allows querying the Jira datasource which contains JSON objects.
type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

// Request is a request to Jira.
type Request struct {
	// BaseURL is the base URL of the Jira instance. For example, "https://{domain}.atlassian.net".
	BaseURL string

	// Username is the user name used to authenticate with Jira using basic auth.
	Username string

	// Password is user's Jira API token used to authenticate with Jira using basic auth.
	Password string

	// PageSize is the maximum number of objects to return per page from the API call.
	// This is used as the "maxResults" parameter in the Jira API.
	PageSize int64

	// EntityExternalID is the external ID of the entity. If it's not a valid external ID, then
	// no request will be made.
	EntityExternalID string

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *pagination.CompositeCursor[int64]

	// IssuesJQLFilter is a JQL filter to apply to the request.
	// This is only used when EntityExternalID = "Issue".
	IssuesJQLFilter *string

	// ObjectsQLQuery is a AQL query to apply to the request.
	// This is only used when EntityExternalID = "Object".
	ObjectsQLQuery *string

	// AssetBaseURL is the base URL to retrieve Asset Objects.
	// This is only used when EntityExternalID = "Object".
	AssetBaseURL *string

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds int
}

// Response is a parsed response returned from Jira.
type Response struct {
	// StatusCode is an HTTP status code.
	StatusCode int

	// RetryAfterHeader is the Retry-After response HTTP header, if set.
	RetryAfterHeader string

	// Objects is the list of parsed entity objects returned from Jira.
	// May be empty.
	Objects []map[string]any

	// NextCursor is the cursor that identifies the first object of the next page.
	// nil if this is the last page in this full sync.
	NextCursor *pagination.CompositeCursor[int64]
}
