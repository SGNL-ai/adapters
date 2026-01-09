// Copyright 2026 SGNL.ai, Inc.
package jira

import (
	"context"
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

const (
	// MaxPageSize is the maximum page size allowed in a GetPage request.
	// Each operation can have a different page size limit, and they may change without notice.
	// https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/#pagination. See the "maxResults" query parameter.
	// Use 1000 as an estimate.
	MaxPageSize = 1000
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if request.Config != nil {
		if err := request.Config.Validate(ctx); err != nil {
			return &framework.Error{
				Message: fmt.Sprintf("Jira config is invalid: %v.", err.Error()),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			}
		}
	}

	// Jira uses Basic Auth for REST API clients:
	//   https://developer.atlassian.com/cloud/jira/platform/security-overview/#scripts-and-other-rest-api-clients.
	// The username is the email address and the password is an API token:
	//   https://developer.atlassian.com/cloud/jira/platform/basic-auth-for-rest-apis/#basic-auth-for-rest-apis.
	// For example, {email}:{api_token}, which is then base64 encoded and supplied as the Authorization header
	// with prefix "Basic ". Like the following: "Authorization: Basic ZnJlZDpmcmVk".
	if request.Auth == nil || request.Auth.Basic == nil {
		return &framework.Error{
			Message: "Jira auth is missing required basic credentials.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if _, found := ValidEntityExternalIDs[request.Entity.ExternalId]; !found {
		return &framework.Error{
			Message: "Jira entity external ID is invalid.",
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
			Message: "Jira requested entity attributes are missing a unique ID attribute.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// Manual testing with Jira users and issues indicates the response is not ordered, despite
	// setting the `orderBy=accountId` (users) and `orderBy=id` (issues) query parameters.
	// Only some resources support ordering: https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/#ordering.
	if request.Ordered {
		return &framework.Error{
			Message: "Jira Ordered property must be false.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if request.PageSize > MaxPageSize {
		return &framework.Error{
			Message: fmt.Sprintf("Jira provided page size (%d) exceeds the maximum (%d).", request.PageSize, MaxPageSize),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	return nil
}
