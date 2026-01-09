// Copyright 2026 SGNL.ai, Inc.

package googleworkspace

import (
	"context"
	"fmt"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if err := request.Config.Validate(ctx); err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("Google Workspace adapter config is invalid: %v.", err.Error()),
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

	if !strings.HasPrefix(request.Auth.HTTPAuthorization, "Bearer ") {
		return &framework.Error{
			Message: `Provided auth token is missing required "Bearer " prefix.`,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if _, found := ValidEntityExternalIDs[request.Entity.ExternalId]; !found {
		return &framework.Error{
			Message: "Provided entity external ID is invalid.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	var uniqueIDAttributeFound bool

	for _, attribute := range request.Entity.Attributes {
		if attribute.ExternalId == ValidEntityExternalIDs[request.Entity.ExternalId].UniqueIDAttribute {
			uniqueIDAttributeFound = true

			break
		}
	}

	if !uniqueIDAttributeFound {
		return &framework.Error{
			Message: fmt.Sprintf("Requested entity %s is missing the required unique ID attribute: %s",
				request.Entity.ExternalId, ValidEntityExternalIDs[request.Entity.ExternalId].UniqueIDAttribute),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	requiredAttributes := append([]string(nil), ValidEntityExternalIDs[request.Entity.ExternalId].RequiredAttributes...)
	if requiredAttributes != nil {
		for _, attribute := range request.Entity.Attributes {
			for idx, requiredID := range requiredAttributes {
				if attribute.ExternalId == requiredID {
					requiredAttributes = append(requiredAttributes[:idx], requiredAttributes[idx+1:]...)

					break
				}
			}
		}

		if len(requiredAttributes) > 0 {
			return &framework.Error{
				Message: fmt.Sprintf("Requested entity attributes are missing required attributes: %v", requiredAttributes),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			}
		}
	}

	if request.Ordered && ValidEntityExternalIDs[request.Entity.ExternalId].OrderByAttribute == "" {
		return &framework.Error{
			Message: fmt.Sprintf("Requested entity external ID, %s, does not support ordering.", request.Entity.ExternalId),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if request.PageSize > ValidEntityExternalIDs[request.Entity.ExternalId].MaxPageSize {
		return &framework.Error{
			Message: fmt.Sprintf("Requested page size, %d, exceeds the maximum allowed value of %d for entity: %s.",
				request.PageSize, ValidEntityExternalIDs[request.Entity.ExternalId].MaxPageSize, request.Entity.ExternalId),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	return nil
}
