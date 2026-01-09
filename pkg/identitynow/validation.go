// Copyright 2026 SGNL.ai, Inc.
package identitynow

import (
	"context"
	"fmt"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

const (
	// MaxPageSize is the maximum page size allowed in a GetPage request.
	// The IdentityNow documentation specifies a page size limit of 250.
	// https://developer.sailpoint.com/idn/api/standard-collection-parameters/#paginating-results.
	MaxPageSize = 250
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if err := request.Config.Validate(ctx); err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("IdentityNow config is invalid: %v.", err.Error()),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	sanitizedAddress := strings.TrimSpace(strings.ToLower(request.Address))
	if strings.HasPrefix(sanitizedAddress, "http://") {
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

	if request.PageSize > MaxPageSize {
		return &framework.Error{
			Message: fmt.Sprintf("PageSize must be less than or equal to %d.", MaxPageSize),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if !strings.HasPrefix(request.Auth.HTTPAuthorization, "Bearer ") {
		return &framework.Error{
			Message: `Provided auth token is missing required "Bearer " prefix.`,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	entityConfig, found := request.Config.EntityConfig[request.Entity.ExternalId]
	if !found {
		return &framework.Error{
			Message: fmt.Sprintf(
				"Entity with external ID %s must be present in the Adapter Config for this SoR.",
				request.Entity.ExternalId),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// Validate that at least the unique ID attribute for the requested entity
	// is requested.
	var uniqueIDAttributeFound bool

	for _, attribute := range request.Entity.Attributes {
		if attribute.ExternalId == entityConfig.UniqueIDAttribute {
			uniqueIDAttributeFound = true

			break
		}
	}

	if !uniqueIDAttributeFound {
		return &framework.Error{
			Message: "Requested entity attributes are missing unique ID attribute.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// Sorting is supported in the IdentityNow API, however supported sort fields are not always the same.
	// Some resources support sorting by "id" while others do not. Therefore, for consistency,
	// enforce that the results must be unordered.
	// e.g. Accounts ("id" supported): https://developer.sailpoint.com/idn/api/v3/list-accounts.
	// e.g. Roles ("id" not supported): https://developer.sailpoint.com/idn/api/v3/list-roles.
	if request.Ordered {
		return &framework.Error{
			Message: "Ordered must be set to false.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	return nil
}
