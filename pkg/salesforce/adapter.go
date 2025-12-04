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

	// Build the list of attributes to query from Salesforce.
	// This includes both regular attributes and multi-select picklist fields (child entities).
	queryAttributes := make(
		[]*framework.AttributeConfig,
		0,
		len(request.Entity.Attributes)+len(request.Entity.ChildEntities),
	)
	queryAttributes = append(queryAttributes, request.Entity.Attributes...)

	for _, childEntity := range request.Entity.ChildEntities {
		queryAttributes = append(queryAttributes, &framework.AttributeConfig{
			ExternalId: childEntity.ExternalId,
			Type:       framework.AttributeTypeString,
		})
	}

	salesforceReq := &Request{
		BaseURL:               request.Address,
		Token:                 request.Auth.HTTPAuthorization,
		PageSize:              request.PageSize,
		EntityExternalID:      request.Entity.ExternalId,
		APIVersion:            request.Config.APIVersion,
		Attributes:            queryAttributes,
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

	// Create a temporary entity config without child entities for framework processing.
	entityForParsing := framework.EntityConfig{
		ExternalId:    request.Entity.ExternalId,
		Attributes:    queryAttributes,
		ChildEntities: nil,
	}

	// The raw JSON objects from the response must be parsed and converted into framework.Objects.
	// DateTime values are parsed using the specified DateTimeFormatWithTimeZone.
	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&entityForParsing, // Use entity without child entities
		resp.Objects,      // Use original objects without transformation

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

	// Transform multi-select picklist fields (semicolon-separated strings) into child entity arrays.
	if len(request.Entity.ChildEntities) > 0 {
		for _, obj := range parsedObjects {
			for _, childEntity := range request.Entity.ChildEntities {
				value, exists := obj[childEntity.ExternalId]

				// Only transform if the value is a non-empty string
				if exists && value != nil {
					if strValue, ok := value.(string); ok && strValue != "" {
						// Split by semicolon and create array of child objects
						values := strings.Split(strValue, ";")
						childObjects := make([]framework.Object, 0, len(values))

						// Get the attribute name from the child entity config (should be exactly one)
						attributeName := childEntity.Attributes[0].ExternalId

						for _, val := range values {
							if val != "" {
								childObjects = append(childObjects, framework.Object{
									attributeName: val,
								})
							}
						}

						// Replace the semicolon-separated string with the array of objects
						obj[childEntity.ExternalId] = childObjects
						continue
					}
				}

				// For nil, empty string, or non-existent fields, set to empty array
				obj[childEntity.ExternalId] = []framework.Object{}
			}
		}
	}

	page := &framework.Page{
		Objects: parsedObjects,
	}

	if resp.NextCursor != nil {
		page.NextCursor = *resp.NextCursor
	}

	return framework.NewGetPageResponseSuccess(page)
}
