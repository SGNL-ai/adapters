// Copyright 2026 SGNL.ai, Inc.

package victorops

import (
	"context"
	"fmt"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	"github.com/sgnl-ai/adapters/pkg/commonutil"
	"github.com/sgnl-ai/adapters/pkg/config"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from VictorOps (Splunk On-Call) datasources.
type Adapter struct {
	VictorOpsClient Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		VictorOpsClient: client,
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

// RequestPageFromDatasource requests a page of objects from the VictorOps datasource.
func (a *Adapter) RequestPageFromDatasource(
	ctx context.Context,
	request *framework.Request[Config],
) framework.Response {
	var commonConfig *config.CommonConfig
	if request.Config != nil {
		commonConfig = request.Config.CommonConfig
	}

	commonConfig = config.SetMissingCommonConfigDefaults(commonConfig)

	victoropsReq := &Request{
		BaseURL:               request.Address,
		APIId:                 request.Auth.Basic.Username,
		APIKey:                request.Auth.Basic.Password,
		PageSize:              request.PageSize,
		EntityExternalID:      request.Entity.ExternalId,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
	}

	if request.Config != nil && request.Config.QueryParameters != nil {
		victoropsReq.QueryParameters = request.Config.QueryParameters[request.Entity.ExternalId]
	}

	// Unmarshal the current cursor.
	cursor, err := pagination.UnmarshalCursor[int64](request.Cursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	victoropsReq.Cursor = cursor

	res, err := a.VictorOpsClient.GetPage(ctx, victoropsReq)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	// An adapter error message is generated if the response status code is not
	// successful (i.e. if not statusCode >= 200 && statusCode < 300).
	adapterErr := web.HTTPError(res.StatusCode, res.RetryAfterHeader)
	if adapterErr != nil {
		return framework.NewGetPageResponseError(adapterErr)
	}

	// Transform string array fields (e.g., pagedUsers, pagedTeams) into child entity objects.
	// VictorOps returns these as JSON arrays of strings; this converts them into arrays of
	// {id, value} objects so they can participate in SGNL relationships.
	objectsToConvert := res.Objects
	if len(request.Entity.ChildEntities) > 0 {
		objectsToConvert = commonutil.CreateChildEntitiesFromStringArray(
			res.Objects,
			&request.Entity,
		)
	}

	// The raw JSON objects from the response must be parsed and converted into framework.Objects.
	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&request.Entity,
		objectsToConvert,

		web.WithJSONPathAttributeNames(),

		web.WithDateTimeFormats(
			[]web.DateTimeFormatWithTimeZone{
				{Format: time.RFC3339, HasTimeZone: true},
				{Format: "2006-01-02T15:04:05Z", HasTimeZone: true},
				{Format: time.DateOnly, HasTimeZone: false},
			}...,
		),
		web.WithLocalTimeZoneOffset(commonConfig.LocalTimeZoneOffset),
	)
	if parserErr != nil {
		return framework.NewGetPageResponseError(
			&framework.Error{
				Message: fmt.Sprintf("Failed to convert VictorOps response objects: %v.", parserErr),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		)
	}

	// Marshal the next cursor.
	nextCursor, err := pagination.MarshalCursor(res.NextCursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	return framework.NewGetPageResponseSuccess(&framework.Page{
		Objects:    parsedObjects,
		NextCursor: nextCursor,
	})
}
