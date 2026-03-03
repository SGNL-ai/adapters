// Copyright 2026 SGNL.ai, Inc.

package duo

import (
	"context"
	"fmt"
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
	DuoClient Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		DuoClient: client,
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
	cursor, err := pagination.UnmarshalCursor[int64](request.Cursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	duoReq := &Request{
		BaseURL:               request.Address,
		IntegrationKey:        request.Auth.Basic.Username,
		Secret:                request.Auth.Basic.Password,
		PageSize:              request.PageSize,
		EntityExternalID:      request.Entity.ExternalId,
		APIVersion:            request.Config.APIVersion,
		Cursor:                cursor,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
	}

	resp, err := a.DuoClient.GetPage(ctx, duoReq)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	// An adapter error message is generated if the response status code is not
	// successful (i.e. if not statusCode >= 200 && statusCode < 300).
	if adapterErr := web.HTTPError(resp.StatusCode, resp.RetryAfterHeader); adapterErr != nil {
		return framework.NewGetPageResponseError(adapterErr)
	}

	// The raw JSON objects from the response must be parsed and converted into framework.Objects.
	// DateTime values are parsed using the specified DateTimeFormatWithTimeZone.
	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&request.Entity,
		resp.Objects,
		web.WithJSONPathAttributeNames(),
		web.WithDateTimeFormats(
			[]web.DateTimeFormatWithTimeZone{
				// The only date formats that Duo responses contain are UNIX timestamps and RFC3339.
				// RFC3339: Phone Entity -> "last_seen": "2019-11-18T15:51:13"
				// Unix Timestamp: User Entity -> "created": 1574100000
				{Format: time.RFC3339, HasTimeZone: true},
				{Format: web.SGNLUnixSec, HasTimeZone: false},
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
