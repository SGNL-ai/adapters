// Copyright 2026 SGNL.ai, Inc.
package bamboohr

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
	BambooHRClient Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		BambooHRClient: client,
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

	if cursor != nil && cursor.Cursor != nil && *cursor.Cursor <= 0 {
		return framework.NewGetPageResponseError(
			&framework.Error{
				Message: "Cursor value must be greater than 0.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		)
	}

	trimmedAddress := strings.TrimSpace(request.Address)
	sanitizedAddress := strings.ToLower(trimmedAddress)

	if !strings.HasPrefix(sanitizedAddress, "https://") {
		request.Address = "https://" + trimmedAddress
	}

	bambooReq := &Request{
		BaseURL:               request.Address,
		APIVersion:            request.Config.APIVersion,
		CompanyDomain:         request.Config.CompanyDomain,
		OnlyCurrent:           request.Config.OnlyCurrent,
		APIKey:                request.Auth.Basic.Username,
		BasicAuthPassword:     request.Auth.Basic.Password,
		AttributeMappings:     request.Config.AttributeMappings,
		PageSize:              request.PageSize,
		EntityConfig:          &request.Entity,
		Cursor:                cursor,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
	}

	resp, err := a.BambooHRClient.GetPage(ctx, bambooReq)
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
				// BambooHR has two datetime formats: date: 'yyyy-mm-dd' and timestamp: '####-##-##T##:##:##+##:##'
				// https://documentation.bamboohr.com/docs/field-types
				supportedDateFormats[*request.Config.AttributeMappings.Date],
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
