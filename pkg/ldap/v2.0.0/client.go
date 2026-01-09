// Copyright 2026 SGNL.ai, Inc.
package ldap

import (
	"context"

	framework "github.com/sgnl-ai/adapter-framework"
	grpc_proxy_v1 "github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Client is a client that allows querying the datasource which contains JSON objects.
type Client interface {
	GetPage(ctx context.Context, request *Request) (*Response, *framework.Error)
}

type ConnectionParams struct {
	// Host is the Hostname of the datasource to query.
	Host string `json:"host"`

	// BaseDN is the Base DN of the datasource to query.
	BaseDN string `json:"baseDN"`

	// BindDN is the Bind DN of the datasource to query.
	BindDN string `json:"bindDN"`

	// BindPassword is the password of the datasource to query.
	BindPassword string `json:"bindPassword"`

	// IsLDAPS flag to check if connection is secured
	IsLDAPS bool `json:"isLDAPS"`

	// CertificateChain contains certificate chain to use for ldaps connection
	CertificateChain string `json:"certificateChain,omitempty"`
}

// Request is a request to the datasource.
type Request struct {
	// ConnectionParams contains LDAP specific params
	ConnectionParams `json:"connectionParams"`

	// BaseURL is the Base URL of the datasource to query.
	BaseURL string `json:"baseURL"`

	// PageSize is the maximum number of objects to return from the entity.
	PageSize int64 `json:"pageSize"`

	// EntityExternalID is the external ID of the entity.
	// The external ID should match the API's resource name.
	EntityExternalID string `json:"entityExternalID"`

	// Cursor identifies the first object of the page to return, as returned by
	// the last request for the entity.
	// nil in the request for the first page.
	Cursor *pagination.CompositeCursor[string] `json:"cursor,omitempty"`

	// UniqueIDAttribute is a attribute which can be used to uniquely identify the Entity.
	// This is specific to ldap server implementation
	UniqueIDAttribute string `json:"uniqueIDAttribute"`

	// EntityConfigMap is an map containing the config required for each entity associated with this
	// datasource. The key is the entity's external_name and value is EntityConfig.
	EntityConfigMap map[string]*EntityConfig `json:"entityConfigMap,omitempty"`

	// Attributes contains the list of attributes to request along with the current request.
	Attributes []*framework.AttributeConfig `json:"attributes,omitempty"`

	// RequestTimeoutSeconds is the timeout duration for requests made to datasources.
	// This should be set to the number of seconds to wait before timing out.
	RequestTimeoutSeconds int `json:"requestTimeoutSeconds"`
}

// Response is a response returned by the datasource.
type Response struct {
	// TODO: Update the comment once we support LDAP status with adapter-framework
	// StatusCode is an HTTP status code.
	StatusCode int `json:"statusCode"`

	// RetryAfterHeader is the Retry-After response HTTP header, if set.
	RetryAfterHeader string `json:"retryAfterHeader"`

	// Objects is the list of
	// May be empty.
	Objects []map[string]any `json:"objects,omitempty"`

	// NextCursor is the cursor that identifies the first object of the next page.
	// nil if this is the last page in this full sync.
	NextCursor *pagination.CompositeCursor[string] `json:"nextCursor"`
}

// NewClient returns a Client to query the datasource.
func NewClient(proxy grpc_proxy_v1.ProxyServiceClient, pool *SessionPool) Client {
	return &Datasource{
		Client: &ldapClient{
			proxyClient: proxy,
			sessionPool: pool,
		},
	}
}
