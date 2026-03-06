// Copyright 2026 SGNL.ai, Inc.

package victorops

import (
	"context"
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"

	"github.com/sgnl-ai/adapters/pkg/validation"
)

const (
	// MaxPageSize is the maximum page size allowed in a GetPage request.
	// The VictorOps Reporting Incidents API supports a maximum limit of 100.
	// The Users API returns all users in a single response, so this limit applies to IncidentReport only,
	// but we enforce it uniformly for consistency.
	MaxPageSize = 100
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if request.Config != nil {
		if err := request.Config.Validate(ctx); err != nil {
			return &framework.Error{
				Message: fmt.Sprintf("VictorOps config is invalid: %v.", err.Error()),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			}
		}
	}

	trimmedAddress, parsed, err := validation.ParseAndValidateAddress(request.Address, []string{"https"})
	if err != nil {
		return err
	}

	// Normalize address with https:// scheme if not provided.
	if parsed.Scheme == "" {
		request.Address = "https://" + trimmedAddress
	} else {
		request.Address = trimmedAddress
	}

	// VictorOps uses two custom headers for authentication: X-VO-Api-Id and X-VO-Api-Key.
	// Since the SGNL framework does not have a native auth type for dual custom headers,
	// we overload HTTP Basic authentication:
	//   - Username field → X-VO-Api-Id header value
	//   - Password field → X-VO-Api-Key header value
	// This is NOT standard Basic auth — the values are sent as separate custom headers,
	// not as a base64-encoded Authorization header.
	if request.Auth == nil || request.Auth.Basic == nil {
		return &framework.Error{
			Message: "VictorOps auth is missing required basic credentials (API ID and API key).",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if request.Auth.Basic.Username == "" {
		return &framework.Error{
			Message: "VictorOps API ID (basic auth username) must not be empty.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if request.Auth.Basic.Password == "" {
		return &framework.Error{
			Message: "VictorOps API key (basic auth password) must not be empty.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if _, found := ValidEntityExternalIDs[request.Entity.ExternalId]; !found {
		return &framework.Error{
			Message: "VictorOps entity external ID is invalid.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// Validate that at least the unique ID attribute for the requested entity
	// is requested.
	var uniqueIDAttributeFound bool

	for _, attribute := range request.Entity.Attributes {
		if attribute.ExternalId == ValidEntityExternalIDs[request.Entity.ExternalId].uniqueIDAttrExternalID {
			uniqueIDAttributeFound = true

			break
		}
	}

	if !uniqueIDAttributeFound {
		return &framework.Error{
			Message: "VictorOps requested entity attributes are missing a unique ID attribute.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// VictorOps does not guarantee ordered responses.
	if request.Ordered {
		return &framework.Error{
			Message: "VictorOps Ordered property must be false.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if request.PageSize > MaxPageSize {
		return &framework.Error{
			Message: fmt.Sprintf("VictorOps provided page size (%d) exceeds the maximum (%d).", request.PageSize, MaxPageSize),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	return nil
}
