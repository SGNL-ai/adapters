// Copyright 2026 SGNL.ai, Inc.

package jiradatacenter

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

	// AuthorizationHeader is the Authorization header sent to the Jira SoR.
	AuthorizationHeader string

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

	// Groups is a list of group names to filter on. For each group in this list,
	// only matching groups will be synchronized, and users are synced based on membership
	// in the listed groups. If this field is empty or nil, the adapter will sync
	// all available groups and their members.
	Groups []string

	// IssuesJQLFilter is a JQL filter to apply to the request.
	// This is only used when EntityExternalID = "Issue".
	IssuesJQLFilter *string

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds int

	// APIVersion specifies which API version of JIRA Datacenter to use
	// In many cases, this should be set to 'latest'
	APIVersion string

	// IncludeInactiveUsers determines whether inactive users are included in the results
	// when querying for User or GroupMember entities.
	IncludeInactiveUsers *bool

	// GroupsMaxResults is the maximum number of groups to return per page from the groups/picker API.
	// This is only used when EntityExternalID = "Group".
	// If not specified, the Jira API default behavior is used.
	GroupsMaxResults *int64

	// Attributes contains the list of attributes to request along with the current request.
	// This is used to limit the fields returned in the API response.
	Attributes []*framework.AttributeConfig

	// Entity is the configuration of the entity to get data from.
	Entity *framework.EntityConfig `json:"entityConfig"`
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
