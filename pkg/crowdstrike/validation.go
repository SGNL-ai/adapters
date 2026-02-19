// Copyright 2026 SGNL.ai, Inc.

package crowdstrike

import (
	"context"
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"

	"github.com/sgnl-ai/adapters/pkg/validation"
)

const (
	// MaxPageSize is the maximum page size allowed in a GetPage request.
	MaxPageSize = 1000
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if err := request.Config.Validate(ctx); err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("CrowdStrike config is invalid: %v.", err.Error()),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if _, _, err := validation.ParseAndValidateAddress(request.Address, []string{"https"}); err != nil {
		return err
	}

	if request.Auth == nil || request.Auth.HTTPAuthorization == "" {
		return &framework.Error{
			Message: "Provided datasource auth is missing required credentials.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// An entity has to be supported to be fetched by one of REST or GraphQL.
	graphQLEntityInfo, isGraphQLEntity := ValidGraphQLEntityExternalIDs[request.Entity.ExternalId]
	_, isRESTEntity := ValidRESTEntityExternalIDs[request.Entity.ExternalId]

	if !isGraphQLEntity && !isRESTEntity {
		return &framework.Error{
			Message: "Provided entity external ID is invalid.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if isGraphQLEntity && isRESTEntity {
		return &framework.Error{
			Message: "Provided entity external ID is misconfigured.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	var (
		uniqueIDAttribute *string
	)

	for _, attribute := range request.Entity.Attributes {
		if attribute.UniqueId {
			uniqueIDAttribute = &attribute.ExternalId

			break
		}
	}

	if uniqueIDAttribute == nil {
		return &framework.Error{
			Message: "Requested entity attributes are missing unique ID attribute.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// GraphQL APIs mandate a specific unique ID attribute.
	if isGraphQLEntity &&
		*uniqueIDAttribute != graphQLEntityInfo.UniqueIDAttrExternalID {
		return &framework.Error{
			Message: "Expected unique ID attribute: " + graphQLEntityInfo.UniqueIDAttrExternalID,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if request.Ordered {
		return &framework.Error{
			Message: "Ordered must be set to false.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if request.PageSize > MaxPageSize {
		return &framework.Error{
			Message: fmt.Sprintf("Provided page size (%d) exceeds the maximum allowed (%d).", request.PageSize, MaxPageSize),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	return nil
}
