package awsidentitycenter

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Client is a client that allows querying the datasource which contains JSON objects.
type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

type Auth struct {
	// AccessKey is the access key to authenticate with AWS.
	AccessKey string

	// SecretKey is the secret key to authenticate with AWS.
	SecretKey string

	// Region is the AWS region to query.
	Region string
}

// Request is a request to the datasource.
type Request struct {
	Auth

	// IdentityStoreID is the AWS Identity Store identifier.
	IdentityStoreID string

	// InstanceARN is the AWS Identity Center instance ARN.
	InstanceARN string

	// MaxResults is the maximum number of objects to return from the entity.
	MaxResults int32

	// EntityExternalID is the external ID of the entity.
	EntityExternalID string

	// Cursor identifies the first object of the page to return.
	Cursor *pagination.CompositeCursor[string]

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	RequestTimeoutSeconds int
}

// Response is a response returned by the datasource.
type Response struct {
	// StatusCode is an HTTP status code.
	StatusCode int

	// RetryAfterHeader is the Retry-After response HTTP header, if set.
	RetryAfterHeader string

	// Objects is the list of items returned by the datasource.
	Objects []map[string]any

	// NextCursor is the cursor that identifies the first object of the next page.
	NextCursor *pagination.CompositeCursor[string]
}
