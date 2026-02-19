// Copyright 2026 SGNL.ai, Inc.

package rootly

import (
	"context"
	"fmt"
	"strings"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	"github.com/sgnl-ai/adapters/pkg/config"
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from datasources.
type Adapter struct {
	RootlyClient Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		RootlyClient: client,
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

	trimmedAddress := strings.TrimSpace(request.Address)
	sanitizedAddress := strings.ToLower(trimmedAddress)

	if !strings.HasPrefix(sanitizedAddress, "https://") {
		request.Address = fmt.Sprintf("https://%s", trimmedAddress)
	}

	if !strings.HasSuffix(trimmedAddress, "/") {
		request.Address = fmt.Sprintf("%s/", request.Address)
	}

	baseURL := fmt.Sprintf("%s%s", request.Address, request.Config.APIVersion)

	var authorizationHeader string
	if request.Auth != nil && request.Auth.HTTPAuthorization != "" {
		authorizationHeader = request.Auth.HTTPAuthorization
	}

	var cursor *string
	if request.Cursor != "" {
		cursor = &request.Cursor
	}

	apiRequest := &Request{
		BaseURL:               baseURL,
		HTTPAuthorization:     authorizationHeader,
		EntityExternalID:      request.Entity.ExternalId,
		PageSize:              request.PageSize,
		Cursor:                cursor,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
		Filter:                a.getFilterForEntity(request),
		Includes:              a.getIncludesForEntity(request),
	}

	response, err := a.RootlyClient.GetPage(ctx, apiRequest)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	// Type conversion for Rootly attributes
	for i, obj := range response.Objects {
		if attrsMap, ok := obj["attributes"].(map[string]any); ok {
			val, exists := attrsMap["sequential_id"]
			if !exists {
				continue
			}

			// Update the sequential_id attribute based upon the type.
			for _, attr := range request.Entity.Attributes {
				// Special case for sequential_id attribute for converting to correct type.
				if strings.HasSuffix(attr.ExternalId, "attributes.sequential_id") {
					switch attr.Type {
					case framework.AttributeTypeInt64, framework.AttributeTypeDouble:
						// Convert to float64 as required by the framework
						if floatVal, err := castToFloat64(val); err == nil {
							attrsMap["sequential_id"] = floatVal
						}
					case framework.AttributeTypeBool:
						if boolVal, err := castToBool(val); err == nil {
							attrsMap["sequential_id"] = boolVal
						}
					case framework.AttributeTypeString,
						framework.AttributeTypeDateTime,
						framework.AttributeTypeDuration:
						if strVal, err := castToString(val); err == nil {
							attrsMap["sequential_id"] = strVal
						}
					}

					break
				}
			}

			obj["attributes"] = attrsMap
		}

		response.Objects[i] = obj
	}

	// Convert JSON objects to framework objects
	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&request.Entity,
		response.Objects,
		web.WithJSONPathAttributeNames(),
		web.WithDateTimeFormats(
			[]web.DateTimeFormatWithTimeZone{
				{Format: time.RFC3339, HasTimeZone: true},
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

	if response.NextCursor != nil {
		page.NextCursor = *response.NextCursor
	}

	return framework.NewGetPageResponseSuccess(page)
}

// getFilterForEntity builds the filter string for the given entity from the config.
func (a *Adapter) getFilterForEntity(request *framework.Request[Config]) string {
	if request.Config == nil || request.Config.Filters == nil {
		return ""
	}

	if filter, exists := request.Config.Filters[request.Entity.ExternalId]; exists {
		return filter
	}

	return ""
}

// getIncludesForEntity builds the includes string for the given entity from the config.
func (a *Adapter) getIncludesForEntity(request *framework.Request[Config]) string {
	if request.Config == nil || request.Config.Includes == nil {
		return ""
	}

	if includes, exists := request.Config.Includes[request.Entity.ExternalId]; exists {
		return includes
	}

	return ""
}
