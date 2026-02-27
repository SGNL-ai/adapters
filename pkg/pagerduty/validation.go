// Copyright 2026 SGNL.ai, Inc.

package pagerduty

import (
	"context"
	"fmt"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"

	"github.com/sgnl-ai/adapters/pkg/validation"
)

const (
	// MaxPageSize is the maximum page size allowed in a GetPage request.
	// https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTU4-pagination#classic-pagination.
	MaxPageSize = 100

	// UniqueIDAttribute is the name of the attribute containing the unique ID of
	// each returned object for the requested entity.
	// https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTU1-types#id.
	UniqueIDAttribute = "id"
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if err := request.Config.Validate(ctx); err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("PagerDuty config is invalid: %v.", err.Error()),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// PagerDuty uses HTTP auth via an API token.
	// https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTUx-authentication#api-token-authentication.
	if request.Auth == nil || request.Auth.HTTPAuthorization == "" {
		return &framework.Error{
			Message: "PagerDuty auth is missing required token.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if !strings.HasPrefix(request.Auth.HTTPAuthorization, "Token token=") &&
		!strings.HasPrefix(request.Auth.HTTPAuthorization, "Bearer ") {
		return &framework.Error{
			Message: "PagerDuty auth is missing required \"Token token=\" or \"Bearer \" prefix.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	trimmedAddress, parsed, err := validation.ParseAndValidateAddress(request.Address, []string{"https"})
	if err != nil {
		return err
	}

	// All API calls are made to the same DNS domain name.
	// The authentication token dictates what data to return.
	// https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTUw-rest-api-v2-overview#what-and-where.
	if parsed.Host != "api.pagerduty.com" {
		return &framework.Error{
			Message: "Invalid PagerDuty address. Must be api.pagerduty.com.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// Normalize address with https:// scheme if not provided
	if parsed.Scheme == "" {
		request.Address = "https://" + trimmedAddress
	} else {
		request.Address = trimmedAddress
	}

	// Validate that at least the unique ID attribute for the requested entity
	// is requested.
	var uniqueIDAttributeFound bool

	for _, attribute := range request.Entity.Attributes {
		if attribute.ExternalId == UniqueIDAttribute {
			uniqueIDAttributeFound = true

			break
		}
	}

	if !uniqueIDAttributeFound {
		return &framework.Error{
			Message: fmt.Sprintf(
				"PagerDuty requested entity attributes are missing a unique ID attribute: %s.",
				UniqueIDAttribute,
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// Validate that no child entities are requested.
	if len(request.Entity.ChildEntities) > 0 {
		return &framework.Error{
			Message: "PagerDuty requested entity does not support child entities.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// Although PagerDuty reports all responses are sorted, they are not sorted by the unique ID.
	// Furthermore, not every endpoint supports the `sort_by` query parameter.
	// Therefore, assume responses are unsorted.
	// https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTU3-sorting.
	if request.Ordered {
		return &framework.Error{
			Message: "PagerDuty Ordered property must be false.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if request.PageSize > MaxPageSize {
		return &framework.Error{
			Message: fmt.Sprintf("PagerDuty provided page size (%d) exceeds the maximum (%d).", request.PageSize, MaxPageSize),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	return nil
}
