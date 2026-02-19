// Copyright 2026 SGNL.ai, Inc.

package jiradatacenter

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
	// Each operation can have a different page size limit, and they may change without notice.
	// https://developer.atlassian.com/server/jira/platform/rest/v10000/intro/#pagination.
	// See the "maxResults" query parameter.
	// Use 1000 as an default estimate.
	MaxPageSize = 1000

	// MaxGroupMemberPageSize is the maximum page size for group member operations.
	// This limit is not explicitly documented, but the Jira Data Center API returns
	// only up to 50 members even if a higher maxResults value is provided.
	// The adapter enforces this limit to align with the API behavior.
	MaxGroupMemberPageSize = 50
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

	if _, _, err := validation.ParseAndValidateAddress(request.Address, []string{"https"}); err != nil {
		return err
	}

	// Jira Data Center supports various authentication methods:
	// https://developer.atlassian.com/server/jira/platform/rest/v10000/intro/#authentication
	//
	// This adapter supports two authentication methods:
	// 1. Personal Access Token (PAT) - should be supplied as request.Auth.HTTPAuthorization
	//    with prefix "Bearer ". Example: "Authorization: Bearer NjQ5OTk0NjE4OTI0OhfQ7+..."
	//    https://confluence.atlassian.com/enterprise/using-personal-access-tokens-1026032365.html
	//
	// 2. Basic HTTP Authentication - should be supplied as request.Auth.Basic
	//    Username and password are base64 encoded and supplied with prefix "Basic "
	//    Example: "Authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ="
	//    https://developer.atlassian.com/server/jira/platform/basic-authentication/
	if request.Auth == nil || request.Auth.Basic == nil && request.Auth.HTTPAuthorization == "" {
		return &framework.Error{
			Message: "Request to Jira is missing Basic Auth or Personal Access Token (PAT) credentials.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// If using HTTPAuthorization, ensure it has the correct prefix
	if request.Auth.HTTPAuthorization != "" &&
		!strings.HasPrefix(request.Auth.HTTPAuthorization, "Bearer ") {
		return &framework.Error{
			Message: `Provided auth token is missing required "Bearer " prefix.`,
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

	// Jira supports ordering (via 'orderBy') only for specific APIs, but the endpoints
	// we're accessing for users, groups, and group members don't support customizable sorting.
	// These APIs return results sorted by fixed fields determined by Jira's implementation.
	// Since we cannot control or guarantee the sort order, we require clients to set
	// request.Ordered to false to acknowledge this limitation. (For Issues, sorting can be achieved through
	// 'ORDER BY' clauses in JQL queries, which is configured through IssuesJQLFilter in the adapter config).
	if request.Ordered {
		return &framework.Error{
			Message: "Jira Ordered property must be false.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// Add specific validation for User entity and group members page size.
	if (request.Entity.ExternalId == UserExternalID ||
		request.Entity.ExternalId == GroupMemberExternalID) &&
		(request.PageSize > MaxGroupMemberPageSize) {
		return &framework.Error{
			Message: fmt.Sprintf("User or group member page size (%d) exceeds allowed maximum (%d).",
				request.PageSize, MaxGroupMemberPageSize),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
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
