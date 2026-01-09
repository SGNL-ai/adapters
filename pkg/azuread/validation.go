// Copyright 2026 SGNL.ai, Inc.

package azuread

import (
	"context"
	"fmt"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// The maximum API page size varies, but it appears the typical maximum page size for the Graph API is 999,
// so this is what we'll enforce. Depending on the specific endpoint and attributes requested, the specified
// page size may not be respected if your request configuration supports a lower max page size than 999.
// For example, if you request users with the `signInActivity` attribute, any request with a page size greater
// than 120 will return pages with up to 120 users.
// https://learn.microsoft.com/en-us/graph/paging?tabs=http
// https://learn.microsoft.com/en-us/graph/api/user-list?view=graph-rest-1.0&tabs=http#optional-query-parameters
// https://learn.microsoft.com/en-us/graph/api/group-list?view=graph-rest-1.0&tabs=http#optional-query-parameters
// https://learn.microsoft.com/en-us/graph/api/application-list?view=graph-rest-1.0&tabs=http#optional-query-parameters

const (
	uniqueIDAttribute = "id"
	maxPageSize       = 999
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if err := request.Config.Validate(ctx); err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("Azure AD config is invalid: %v.", err.Error()),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if err := validateAdvancedFilterConfiguration(request); err != nil {
		return err
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

	// Validate that at least the unique ID attribute for the requested entity
	// is requested.
	var uniqueIDAttributeFound bool

	for _, attribute := range request.Entity.Attributes {
		if attribute.ExternalId == uniqueIDAttribute {
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

	return nil
}
