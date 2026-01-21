// Copyright 2026 SGNL.ai, Inc.

package github

import (
	"context"
	"fmt"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// The maximum API page size for the GitHub GraphQL APIs are 100.
// https://docs.github.com/en/graphql/guides/using-pagination-in-the-graphql-api#about-pagination
const (
	maxPageSize = 100
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if _, found := ValidEntityExternalIDs[request.Entity.ExternalId]; !found {
		return &framework.Error{
			Message: "Provided entity external ID is invalid.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if err := request.Config.Validate(ctx, ValidEntityExternalIDs[request.Entity.ExternalId].isRestAPI); err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("GitHub config is invalid: %v.", err.Error()),
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

	var uniqueIDAttributeFound bool

	for _, attribute := range request.Entity.Attributes {
		if attribute.ExternalId == ValidEntityExternalIDs[request.Entity.ExternalId].UniqueExternalIDAttribute {
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

	requiredAttributes := ValidEntityExternalIDs[request.Entity.ExternalId].RequiredAttributes
	if requiredAttributes != nil {
		requiredAttributesFound := make(map[string]bool)

		for _, attribute := range request.Entity.Attributes {
			for _, requiredID := range requiredAttributes {
				if attribute.ExternalId == requiredID {
					requiredAttributesFound[requiredID] = true
				}
			}
		}

		for _, requiredID := range requiredAttributes {
			if !requiredAttributesFound[requiredID] {
				return &framework.Error{
					Message: fmt.Sprintf("Requested entity attributes are missing required attribute: %s", requiredID),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}
			}
		}
	}

	// [sc-22432] TODO: Add Ordering Support
	if request.Ordered {
		return &framework.Error{
			Message: "Ordered must be set to false.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if request.PageSize > maxPageSize {
		return &framework.Error{
			Message: fmt.Sprintf(
				"Provided page size (%d) exceeds the maximum allowed (%d).", request.PageSize, maxPageSize,
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	if len(request.Config.Organizations) > 0 {
		for idx, organization := range request.Config.Organizations {
			if organization == "" {
				return &framework.Error{
					Message: fmt.Sprintf("organizations[%d] cannot be an empty string.", idx),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}
			}
		}
	}

	return nil
}
