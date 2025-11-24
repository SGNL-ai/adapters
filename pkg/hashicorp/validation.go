// Copyright 2025 SGNL.ai, Inc.
package hashicorp

import (
	"context"
	"fmt"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

const (
	uniqueIDAttribute = "id"

	minPageSize = 10
	maxPageSize = 10000
)

var (
	CursorsWithParentEntity = map[string]struct{}{
		EntityTypeAccounts:            {},
		EntityTypeCredentials:         {},
		EntityTypeCredentialLibraries: {},
		EntityTypeHosts:               {},
		EntityTypeHostSets:            {},
	}
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if request == nil {
		return &framework.Error{
			Message: "Request is nil",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if request.Config == nil {
		return &framework.Error{
			Message: "Request config is nil",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if err := request.Config.Validate(ctx); err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("HashiCorp config is invalid: %v.", err.Error()),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if strings.HasPrefix(request.Address, "http://") {
		return &framework.Error{
			Message: "The provided HTTP protocol is not supported.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if a.SSRFValidator != nil {
		if err := a.SSRFValidator.ValidateExternalURL(ctx, request.Address); err != nil {
			return &framework.Error{
				Message: fmt.Sprintf("Address URL validation failed: %v", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			}
		}
	}

	if request.Auth == nil ||
		request.Auth.Basic == nil ||
		request.Auth.Basic.Username == "" || request.Auth.Basic.Password == "" {
		return &framework.Error{
			Message: "Provided datasource auth is missing required http authorization credentials.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if request.Entity.Attributes == nil {
		return &framework.Error{
			Message: "Request entity attributes is nil",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// Validate that at least the unique ID attribute for the requested entity
	// is requested.
	var uniqueIDAttributeFound bool

	for _, attribute := range request.Entity.Attributes {
		if attribute == nil {
			continue
		}

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

	cursor, err := pagination.UnmarshalCursor[string](request.Cursor)
	if err != nil {
		return err
	}

	if cursor != nil {
		_, hasParentEntity := CursorsWithParentEntity[request.Entity.ExternalId]
		validationErr := pagination.ValidateCompositeCursor(
			cursor,
			request.Entity.ExternalId,
			// If the cursor contains a CollectionID, the entity has a parent collection.
			hasParentEntity,
		)

		if validationErr != nil {
			return validationErr
		}
	}

	return nil
}
