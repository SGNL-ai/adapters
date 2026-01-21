// Copyright 2026 SGNL.ai, Inc.

package okta

import (
	"context"
	"fmt"
	"strings"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	"github.com/sgnl-ai/adapters/pkg/config"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from datasources.
type Adapter struct {
	OktaClient Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		OktaClient: client,
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

	if !strings.HasPrefix(request.Address, "https://") {
		request.Address = "https://" + request.Address
	}

	oktaReq := &Request{
		BaseURL:               request.Address,
		Token:                 request.Auth.HTTPAuthorization,
		PageSize:              request.PageSize,
		EntityExternalID:      request.Entity.ExternalId,
		APIVersion:            request.Config.APIVersion,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
		Filter:                request.Config.Filters[request.Entity.ExternalId],
		Search:                request.Config.Search[request.Entity.ExternalId],
	}

	// Unmarshal the current cursor.
	cursor, err := pagination.UnmarshalCursor[string](request.Cursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	oktaReq.Cursor = cursor

	resp, err := a.OktaClient.GetPage(ctx, oktaReq)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	// An adapter error message is generated if the response status code is not
	// successful (i.e. if not statusCode >= 200 && statusCode < 300).
	if adapterErr := web.HTTPError(resp.StatusCode, resp.RetryAfterHeader); adapterErr != nil {
		return framework.NewGetPageResponseError(adapterErr)
	}

	// The raw JSON objects from the response must be parsed and converted into framework.Objects.
	// Nested attributes are flattened and delimited by the delimiter specified.
	// DateTime values are parsed using the specified DateTimeFormatWithTimeZone.
	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&request.Entity,
		resp.Objects,

		// TODO [sc-14078]: Remove support for complex attribute names.
		web.WithComplexAttributeNameDelimiter("__"),

		web.WithJSONPathAttributeNames(),

		// The Okta developer API specifies that all Date objects are returned in ISO 8601 format:
		// https://developer.okta.com/docs/reference/core-okta-api/#media-types
		web.WithDateTimeFormats(
			[]web.DateTimeFormatWithTimeZone{
				// While the API technically specifies ISO 8601, RFC 3339 is a profile (subset) of ISO 8601 and it
				// appears datetimes in API response are RFC 3339 compliant, so we'll be using the RFC 3339 predefined
				// layout since golang does not have built-in support for ISO 8601. However, this cannot be guaranteed
				// so additional formats should be added here as necessary.
				// https://datatracker.ietf.org/doc/html/rfc3339
				{Format: time.RFC3339, HasTimeZone: true},
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
