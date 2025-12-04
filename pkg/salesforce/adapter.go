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

	// Add child entity fields to the query attributes list.
	// Multi-select picklists are stored as semicolon-separated strings in Salesforce,
	// so we need to query them as regular fields.
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

	// Store the transformed multi-select picklist data for later
	// We'll process child entities AFTER the framework converts regular attributes
	var multiSelectPicklistData map[int]map[string][]framework.Object
	if len(request.Entity.ChildEntities) > 0 {
		multiSelectPicklistData = make(map[int]map[string][]framework.Object)
		transformedObjects := transformMultiSelectPicklists(resp.Objects, request.Entity.ChildEntities)

		// Extract child entity data before framework processing
		for i, obj := range transformedObjects {
			multiSelectPicklistData[i] = make(map[string][]framework.Object)

			for _, childEntity := range request.Entity.ChildEntities {
				if childData, exists := obj[childEntity.ExternalId]; exists {
					if childArray, ok := childData.([]map[string]any); ok {
						// Convert to framework.Object array
						frameworkArray := make([]framework.Object, len(childArray))
						for j, item := range childArray {
							frameworkArray[j] = framework.Object(item)
						}

						multiSelectPicklistData[i][childEntity.ExternalId] = frameworkArray
					}
				}
			}
		}
	}

	// The raw JSON objects from the response must be parsed and converted into framework.Objects.
	// DateTime values are parsed using the specified DateTimeFormatWithTimeZone.
	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&request.Entity, // Use original entity without child entity placeholders
		resp.Objects,    // Use original objects without transformation

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

	// Add multi-select picklist child entities to the parsed objects
	if len(multiSelectPicklistData) > 0 {
		for i, obj := range parsedObjects {
			if childData, exists := multiSelectPicklistData[i]; exists {
				for fieldName, fieldValue := range childData {
					obj[fieldName] = fieldValue
				}
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

// transformMultiSelectPicklists transforms multi-select picklist fields from semicolon-separated strings
// into arrays of objects for child entity processing.
//
// In Salesforce, multi-select picklists are returned as semicolon-separated values (e.g., "value1;value2;value3").
// To support these as child entities in the framework, we need to:
// 1. Split the semicolon-separated string into individual values
// 2. Create an array of objects, where each object contains the single attribute specified in the child entity config
//
// Example transformation:
// Input: {"Id": "123", "Interests__c": "Sports;Music;Reading"}
// With child entity ExternalId="Interests__c" and attribute ExternalId="value"
// Output: {"Id": "123", "Interests__c": [{"value": "Sports"}, {"value": "Music"}, {"value": "Reading"}]}.
func transformMultiSelectPicklists(objects []map[string]any, childEntities []*framework.EntityConfig) []map[string]any {
	if len(childEntities) == 0 {
		return objects
	}

	// Build a map of field names that need transformation
	multiSelectFields := make(map[string]string) // maps field name -> attribute name

	for _, childEntity := range childEntities {
		if len(childEntity.Attributes) == 1 {
			multiSelectFields[childEntity.ExternalId] = childEntity.Attributes[0].ExternalId
		}
	}

	if len(multiSelectFields) == 0 {
		return objects
	}

	// Transform each object
	transformedObjects := make([]map[string]any, len(objects))
	for i, obj := range objects {
		transformedObj := make(map[string]any, len(obj))

		// Copy all fields
		for key, value := range obj {
			transformedObj[key] = value
		}

		// Transform multi-select picklist fields
		for fieldName, attributeName := range multiSelectFields {
			if value, exists := obj[fieldName]; exists {
				// Only transform if the value is a non-empty string
				if strValue, ok := value.(string); ok && strValue != "" {
					// Split by semicolon and create array of objects
					values := strings.Split(strValue, ";")
					childObjects := make([]map[string]any, 0, len(values))

					for _, val := range values {
						trimmedVal := strings.TrimSpace(val)
						if trimmedVal != "" {
							childObjects = append(childObjects, map[string]any{
								attributeName: trimmedVal,
							})
						}
					}

					// Replace the semicolon-separated string with the array of objects
					if len(childObjects) > 0 {
						transformedObj[fieldName] = childObjects
					} else {
						// If no valid values, set to empty array
						transformedObj[fieldName] = []map[string]any{}
					}
				} else if value == nil {
					// If the value is nil, set to empty array
					transformedObj[fieldName] = []map[string]any{}
				}
			}
		}

		transformedObjects[i] = transformedObj
	}

	return transformedObjects
}
