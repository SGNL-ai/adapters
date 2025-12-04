// Copyright 2025 SGNL.ai, Inc.
package salesforce

import (
	"context"
	"fmt"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

const (
	uniqueIDAttribute = "Id"

	// Minimum/maximum batchSize for v52.0 through v58.0 is 200/2000.
	// https://developer.salesforce.com/docs/atlas.en-us.244.0.api_rest.meta/api_rest/headers_queryoptions.htm
	minPageSize = 200
	maxPageSize = 2000
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if err := request.Config.Validate(ctx); err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("Salesforce config is invalid: %v.", err.Error()),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if strings.HasPrefix(request.Address, "http://") {
		return &framework.Error{
			Message: "The provided HTTP protocol is not supported.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if request.Auth == nil || request.Auth.HTTPAuthorization == "" {
		return &framework.Error{
			Message: "Provided datasource auth is missing required http authorization credentials.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if !strings.HasPrefix(request.Auth.HTTPAuthorization, "Bearer ") {
		return &framework.Error{
			Message: `Provided auth token is missing required "Bearer " prefix.`,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// Validate that at least the unique ID attribute for the requested entity
	// is requested, and validate relationship depth for all attributes.
	var uniqueIDAttributeFound bool

	for _, attribute := range request.Entity.Attributes {
		externalID := attribute.ExternalId

		if externalID == uniqueIDAttribute {
			uniqueIDAttributeFound = true
		}

		// Validate relationship depth: Salesforce SOQL supports up to 5 levels
		// See: https://developer.salesforce.com/docs/atlas.en-us.soql_sosl.meta/soql_sosl/
		//      sforce_api_calls_soql_relationships_query_limits.htm
		if len(strings.Split(strings.TrimPrefix(externalID, "$."), ".")) > 5 {
			return &framework.Error{
				Message: fmt.Sprintf(
					"Attribute '%s' exceeds the maximum relationship depth of 5 levels. "+
						"Salesforce SOQL supports up to 5 levels of child-to-parent relationship traversal.",
					externalID,
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			}
		}
	}

	if !uniqueIDAttributeFound {
		return &framework.Error{
			Message: "Requested entity attributes are missing unique ID attribute.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// Validate child entities for multi-select picklists.
	// Multi-select picklists in Salesforce are returned as semicolon-separated values (e.g., "value1;value2;value3")
	// and must be represented as child entities with exactly one attribute.
	for _, childEntity := range request.Entity.ChildEntities {
		if len(childEntity.Attributes) != 1 {
			return &framework.Error{
				Message: fmt.Sprintf(
					"Child entity '%s' must have exactly one attribute for multi-select picklist support, but has %d attributes.",
					childEntity.ExternalId, len(childEntity.Attributes),
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			}
		}
	}

	if !request.Ordered {
		return &framework.Error{
			Message: "Ordered must be set to true.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if request.PageSize > maxPageSize || request.PageSize < minPageSize {
		return &framework.Error{
			Message: fmt.Sprintf(
				"Provided page size (%d) does not fall within the allowed range (%d-%d).",
				request.PageSize, minPageSize, maxPageSize,
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	return nil
}
