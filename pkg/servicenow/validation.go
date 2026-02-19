// Copyright 2026 SGNL.ai, Inc.

package servicenow

import (
	"context"
	"fmt"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"

	"github.com/sgnl-ai/adapters/pkg/validation"
)

const (
	uniqueIDAttribute = "sys_id"

	// While Servicenow does not have a documented maximum page size, it is recommended in multiple places to enforce
	// a soft maximum page size at around 10,000 results per page and perform pagination for any additional data
	// (since larger page sizes can cause degraded performance). Since we have support for pagination we'll enforce a
	// limit of 10,000 internally to prevent timing out during page request.
	// https://support.servicenow.com/kb?id=kb_article_view&sysparm_article=KB0727636
	maxPageSize = 10_000
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if err := request.Config.Validate(ctx); err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("Servicenow config is invalid: %v.", err.Error()),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	trimmedAddress, parsed, err := validation.ParseAndValidateAddress(request.Address, []string{"https"})
	if err != nil {
		return err
	}

	// Normalize address with https:// scheme if not provided
	if parsed.Scheme == "" {
		request.Address = "https://" + trimmedAddress
	} else {
		request.Address = trimmedAddress
	}

	if request.Auth == nil || (request.Auth.HTTPAuthorization == "" && request.Auth.Basic == nil) {
		return &framework.Error{
			Message: "System of Record is missing required authentication credentials.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if request.Auth.Basic != nil && (request.Auth.Basic.Username == "" || request.Auth.Basic.Password == "") {
		return &framework.Error{
			Message: "One of username or password required for basic auth is empty.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if request.Auth.HTTPAuthorization != "" && !strings.HasPrefix(request.Auth.HTTPAuthorization, "Bearer ") {
		return &framework.Error{
			Message: `Provided auth token is missing required "Bearer " prefix.`,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
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

	if request.PageSize > maxPageSize {
		return &framework.Error{
			Message: fmt.Sprintf(
				"Provided page size (%d) is greater than the allowed maximum (%d).",
				request.PageSize, maxPageSize,
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	return nil
}
