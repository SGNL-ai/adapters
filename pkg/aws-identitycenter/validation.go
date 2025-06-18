package awsidentitycenter

import (
	"context"
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

const (
	// The AWS Identity Center APIs support returning up to 100 items per
	// request. Reference:
	// https://docs.aws.amazon.com/singlesignon/latest/IdentityStoreAPIReference/API_ListUsers.html
	maxPageSize = 100
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if err := request.Config.Validate(ctx); err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("AWS Identity Center config is invalid: %v.", err.Error()),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if request.Auth == nil || request.Auth.Basic == nil ||
		request.Auth.Basic.Username == "" || request.Auth.Basic.Password == "" {
		return &framework.Error{
			Message: "Provided datasource auth is missing required AWS authorization credentials.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if _, ok := validEntities[request.Entity.ExternalId]; !ok {
		return &framework.Error{
			Message: "Provided entity external ID is invalid.",
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
			Message: fmt.Sprintf("Provided page size (%d) exceeds the maximum allowed (%d).", request.PageSize, maxPageSize),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	return nil
}
