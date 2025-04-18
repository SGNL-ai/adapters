// Copyright 2025 SGNL.ai, Inc.
package ldap

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Client is a client that allows querying the datasource which contains JSON objects.
type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

type ConnectionParams struct {
	// Host is the Hostname of the datasource to query.
	Host string

	// BaseDN is the Base DN of the datasource to query.
	BaseDN string

	// BindDN is the Bind DN of the datasource to query.
	BindDN string

	// BindPassword is the password of the datasource to query.
	BindPassword string

	// IsLDAPS flag to check if connection is secured
	IsLDAPS bool

	// CertificateChain contains certificate chain to use for ldaps connection
	CertificateChain string
}

// Request is a request to the datasource.
type Request struct {
	// ConnectionParams contains LDAP specific params
	ConnectionParams

	// BaseURL is the Base URL of the datasource to query.
	BaseURL string

	// PageSize is the maximum number of objects to return from the entity.
	PageSize int64

	// EntityExternalID is the external ID of the entity.
	// The external ID should match the API's resource name.
	EntityExternalID string

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *pagination.CompositeCursor[string]

	// UniqueIDAttribute is a attribute which can be used to uniquely identify the Entity.
	// This is specific to ldap server implementation
	UniqueIDAttribute string

	// EntityConfigMap is an map containing the config required for each entity associated with this
	// datasource. The key is the entity's external_name and value is EntityConfig.
	EntityConfigMap map[string]*EntityConfig

	// Attributes contains the list of attributes to request along with the current request.
	Attributes []*framework.AttributeConfig

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds int
}

// Response is a response returned by the datasource.
type Response struct {
	// TODO: Update the comment once we support LDAP status with adapter-framework
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
