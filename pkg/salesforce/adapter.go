// Copyright 2025 SGNL.ai, Inc.
package salesforce

import (
	"context"
	"fmt"
	"strings"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	"github.com/sgnl-ai/adapters/pkg/config"
	processCel "github.com/sgnl-ai/adapters/pkg/process-cel"
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from datasources.
type Adapter struct {
	SalesforceClient Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		SalesforceClient: client,
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

	salesforceReq := &Request{
		BaseURL:               request.Address,
		Token:                 request.Auth.HTTPAuthorization,
		PageSize:              request.PageSize,
		EntityExternalID:      request.Entity.ExternalId,
		APIVersion:            request.Config.APIVersion,
		Attributes:            request.Entity.Attributes,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
	}

	if request.Config.Filters != nil {
		if curFilter, ok := request.Config.Filters[request.Entity.ExternalId]; ok {
			salesforceReq.Filter = &curFilter
		}
	}

	salesforceReq.Cursor = nil

	if request.Cursor != "" {
		salesforceReq.Cursor = &request.Cursor
	}

	resp, err := a.SalesforceClient.GetPage(ctx, salesforceReq)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	// An adapter error message is generated if the response status code is not
	// successful (i.e. if not statusCode >= 200 && statusCode < 300).
	if adapterErr := web.HTTPError(resp.StatusCode, resp.RetryAfterHeader); adapterErr != nil {
		return framework.NewGetPageResponseError(adapterErr)
	}

	// Process CEL expressions before converting to framework objects

	// Process CEL attributes
	celAttrs := make([]processCel.AttributeConfig, len(request.Entity.Attributes))
	for i, attr := range request.Entity.Attributes {
		celAttrs[i] = &attrWrapper{attr}
	}

	if err := processCel.ProcessCELAttributes(celAttrs, resp.Objects); err != nil {
		return framework.NewGetPageResponseError(
			&framework.Error{
				Message: fmt.Sprintf("Failed to process CEL attributes: %v", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		)
	}

	// Process CEL child entities
	celChildEntities := make([]processCel.ChildEntityConfig, len(request.Entity.ChildEntities))
	for i, childEntity := range request.Entity.ChildEntities {
		celChildEntities[i] = &childEntityWrapper{childEntity}
	}

	if err := processCel.ProcessCELChildEntities(celChildEntities, resp.Objects); err != nil {
		return framework.NewGetPageResponseError(
			&framework.Error{
				Message: fmt.Sprintf("Failed to process CEL child entities: %v", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		)
	}

	// The raw JSON objects from the response must be parsed and converted into framework.Objects.
	// DateTime values are parsed using the specified DateTimeFormatWithTimeZone.
	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&request.Entity,
		resp.Objects,

		web.WithJSONPathAttributeNames(),

		// The below formats are explicitly stated as allowed for v53.0 through v58.0, but there is no documentation
		// for v52.0 and below.
		// https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/intro_valid_date_formats.htm
		web.WithDateTimeFormats(
			[]web.DateTimeFormatWithTimeZone{
				{Format: "2006-01-02T15:04:05.999Z0700", HasTimeZone: true},
				{Format: time.DateOnly, HasTimeZone: false},
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

	page := &framework.Page{
		Objects: parsedObjects,
	}

	if resp.NextCursor != nil {
		page.NextCursor = *resp.NextCursor
	}

	return framework.NewGetPageResponseSuccess(page)
}

// attrWrapper wraps framework.AttributeConfig to implement processCel.AttributeConfig interface
type attrWrapper struct {
	*framework.AttributeConfig
}

func (a *attrWrapper) GetExternalId() string {
	return a.ExternalId
}

// childEntityWrapper wraps framework.EntityConfig to implement processCel.ChildEntityConfig interface
type childEntityWrapper struct {
	*framework.EntityConfig
}

func (c *childEntityWrapper) GetExternalId() string {
	return c.ExternalId
}
