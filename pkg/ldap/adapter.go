// Copyright 2025 SGNL.ai, Inc.

package ldap

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	grpc_proxy_v1 "github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	"github.com/sgnl-ai/adapters/pkg/config"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from datasources.
type Adapter struct {
	ADClient Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client grpc_proxy_v1.ProxyServiceClient) framework.Adapter[Config] {
	return &Adapter{
		ADClient: NewClient(client),
	}
}

// GetPage is called by SGNL's ingestion service to query a page of objects
// from a datasource.
func (a *Adapter) GetPage(ctx context.Context, request *framework.Request[Config]) framework.Response {
	if err := a.ValidateGetPageRequest(ctx, request); err != nil {
		return framework.NewGetPageResponseError(err)
	}

	return a.RequestPageFromDatasource(ctx, request)
}

// RequestPageFromDatasource requests a page of objects from a datasource.
func (a *Adapter) RequestPageFromDatasource(
	ctx context.Context, request *framework.Request[Config],
) framework.Response {
	var commonConfig *config.CommonConfig
	if request.Config != nil {
		commonConfig = request.Config.CommonConfig
	}

	commonConfig = config.SetMissingCommonConfigDefaults(commonConfig)

	// Unmarshal the current cursor.
	cursor, err := pagination.UnmarshalCursor[string](request.Cursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	// At this level request.Address is already validated, skipping parsing error
	url, _ := url.Parse(request.Address)

	var isLDAPS bool
	if strings.HasPrefix(request.Address, "ldaps://") {
		isLDAPS = true
	}

	uniqueIDAttribute := getUniqueIDAttribute(request.Entity.Attributes)

	adReq := &Request{
		BaseURL:          request.Address,
		PageSize:         request.PageSize,
		EntityExternalID: request.Entity.ExternalId,
		Attributes:       request.Entity.Attributes,
		ConnectionParams: ConnectionParams{
			BindDN:           request.Auth.Basic.Username,
			BindPassword:     request.Auth.Basic.Password,
			BaseDN:           request.Config.BaseDN,
			CertificateChain: request.Config.CertificateChain,
			IsLDAPS:          isLDAPS,
			Host:             url.Host,
		},
		UniqueIDAttribute:     *uniqueIDAttribute,
		Cursor:                cursor,
		EntityConfigMap:       request.Config.EntityConfigMap,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
	}

	resp, err := a.ADClient.GetPage(ctx, adReq)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	// An adapter error message is generated if the response status code is not
	// successful (i.e. if not statusCode >= 200 && statusCode < 300).
	if adapterErr := web.HTTPError(resp.StatusCode, resp.RetryAfterHeader); adapterErr != nil {
		return framework.NewGetPageResponseError(adapterErr)
	}

	// The ldap.SearchResults object from the response must be parsed and converted into framework.Objects.
	// Nested attributes are flattened and delimited by the delimiter specified.
	// DateTime values are parsed using the specified DateTimeFormatWithTimeZone.
	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&request.Entity,
		resp.Objects,

		web.WithJSONPathAttributeNames(),
		// TODO: Add tests for all datetime formats defined by LDAP
		web.WithDateTimeFormats(
			[]web.DateTimeFormatWithTimeZone{
				{Format: web.SGNLGeneralizedTime, HasTimeZone: true},
			}...,
		),
		web.WithLocalTimeZoneOffset(commonConfig.LocalTimeZoneOffset),
	)
	if parserErr != nil {
		return framework.NewGetPageResponseError(
			&framework.Error{
				Message: fmt.Sprintf("Failed to convert datasource response objects: %v.", parserErr),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		)
	}

	// Marshal the next cursor.
	nextCursor, err := pagination.MarshalCursor(resp.NextCursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	return framework.NewGetPageResponseSuccess(&framework.Page{
		Objects:    parsedObjects,
		NextCursor: nextCursor,
	})
}
