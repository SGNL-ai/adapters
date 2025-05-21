package awsidentitycenter

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
	Client Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{Client: client}
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

	awsReq := &Request{
		Auth: Auth{
			AccessKey: request.Auth.Basic.Username,
			SecretKey: request.Auth.Basic.Password,
			Region:    request.Config.Region,
		},
		IdentityStoreID:       request.Config.IdentityStoreID,
		InstanceARN:           request.Config.InstanceARN,
		MaxResults:            int32(request.PageSize),
		EntityExternalID:      request.Entity.ExternalId,
		Cursor:                cursor,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
	}

	resp, err := a.Client.GetPage(ctx, awsReq)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	if adapterErr := web.HTTPError(resp.StatusCode, resp.RetryAfterHeader); adapterErr != nil {
		return framework.NewGetPageResponseError(adapterErr)
	}

	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&request.Entity,
		resp.Objects,
		web.WithLocalTimeZoneOffset(commonConfig.LocalTimeZoneOffset),
		web.WithDateTimeFormats(
			[]web.DateTimeFormatWithTimeZone{{Format: time.RFC3339, HasTimeZone: true}}...,
		),
	)
	if parserErr != nil {
		return framework.NewGetPageResponseError(
			&framework.Error{Message: fmt.Sprintf("Failed to convert datasource response objects: %v.", parserErr), Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL},
		)
	}

	nextCursor, err := pagination.MarshalCursor(resp.NextCursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	return framework.NewGetPageResponseSuccess(&framework.Page{Objects: parsedObjects, NextCursor: nextCursor})
}
