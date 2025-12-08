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

	objectsToConvert := resp.Objects
	if len(request.Entity.ChildEntities) > 0 {
		objectsToConvert = transformMultiSelectPicklists(resp.Objects, request.Entity.ChildEntities)
	}

	// The raw JSON objects from the response must be parsed and converted into framework.Objects.
	// DateTime values are parsed using the specified DateTimeFormatWithTimeZone.
	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&request.Entity,
		objectsToConvert,

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

// transformMultiSelectPicklists transforms multi-select picklist fields from semicolon-separated strings
// into arrays of objects for child entity processing.
//
// In Salesforce, multi-select picklists are returned as semicolon-separated values (e.g., "value1;value2;value3").
// To support these as child entities in the framework, we need to:
// 1. Split the semicolon-separated string into individual values
// 2. Create an array of objects, where each object contains the single attribute specified in the child entity config
// 3. Use []any type (not []map[string]any) so the framework can properly assert the type
//
// Example transformation:
// Input: {"Id": "123", "Interests__c": "Sports;Music;Reading"}
// With child entity ExternalId="Interests__c" and attribute ExternalId="value"
// Output: {"Id": "123", "Interests__c": []any{map[string]any{"value":"Sports"},
// map[string]any{"value":"Music"}, map[string]any{"value":"Reading"}}}.
func transformMultiSelectPicklists(objects []map[string]any, childEntities []*framework.EntityConfig) []map[string]any {
	if len(childEntities) == 0 {
		return objects
	}

	multiSelectFields := make(map[string]string)

	for _, childEntity := range childEntities {
		if len(childEntity.Attributes) == 1 {
			multiSelectFields[childEntity.ExternalId] = childEntity.Attributes[0].ExternalId
		}
	}

	if len(multiSelectFields) == 0 {
		return objects
	}

	transformedObjects := make([]map[string]any, len(objects))
	for i, obj := range objects {
		transformedObj := make(map[string]any, len(obj))

		for key, value := range obj {
			transformedObj[key] = value
		}

		for fieldName, attributeName := range multiSelectFields {
			value, exists := obj[fieldName]

			if !exists || value == nil {
				transformedObj[fieldName] = []any{}

				continue
			}

			strValue, ok := value.(string)
			if !ok {
				continue
			}

			if strValue == "" {
				transformedObj[fieldName] = []any{}
			} else {
				values := strings.Split(strValue, ";")

				// Create []any because framework requires this type for child entities
				childObjects := make([]any, 0, len(values))

				for _, val := range values {
					trimmedVal := strings.TrimSpace(val)
					if trimmedVal != "" {
						childObjects = append(childObjects, map[string]any{
							attributeName: trimmedVal,
						})
					}
				}

				transformedObj[fieldName] = childObjects
			}
		}

		transformedObjects[i] = transformedObj
	}

	return transformedObjects
}
