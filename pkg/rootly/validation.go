// Copyright 2025 SGNL.ai, Inc.
package rootly

import (
	"context"
	"fmt"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

const (
	uniqueIDAttribute = "id"
	maxPageSize       = 1000
	minPageSize       = 1
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if request.Config == nil {
		return &framework.Error{
			Message: "request contains no config",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if err := request.Config.Validate(ctx); err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("Rootly config is invalid: %v.", err.Error()),
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
