// Copyright 2026 SGNL.ai, Inc.

package ldap

import (
	"context"
	"fmt"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"

	"github.com/sgnl-ai/adapters/pkg/validation"
)

const (
	maxPageSize    = 999
	validLDAPPort  = 389
	validLDAPSPort = 636
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if request.Config != nil {
		if err := request.Config.Validate(ctx); err != nil {
			return &framework.Error{
				Message: fmt.Sprintf("Active Directory config is invalid: %v.", err.Error()),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			}
		}

		// Check if entityConfig is present for requested entity
		// and validate memberOf entityConfig if present
		entityConfig, found := request.Config.EntityConfigMap[request.Entity.ExternalId]
		if !found {
			return &framework.Error{
				Message: fmt.Sprintf("entityConfig is missing in config for requested entity %v.", request.Entity.ExternalId),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			}
		}
		// memberOf Present so we expect filter for memberOf entity
		memberOf := entityConfig.MemberOf
		if memberOf != nil {
			if _, found := request.Config.EntityConfigMap[*memberOf]; !found {
				return &framework.Error{
					Message: fmt.Sprintf("Entity configuration entityConfig.%s is missing for "+
						"entity specified in entityConfig.%s.memberOf.", *memberOf, request.Entity.ExternalId),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}
			}

			if request.Config.EntityConfigMap[request.Entity.ExternalId].MemberOfUniqueIDAttribute == nil {
				return &framework.Error{
					Message: fmt.Sprintf("Entity configuration entityConfig.%s.memberOfUniqueIdAttribute is missing.",
						request.Entity.ExternalId),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}
			}

			if request.Config.EntityConfigMap[request.Entity.ExternalId].MemberUniqueIDAttribute == nil {
				return &framework.Error{
					Message: fmt.Sprintf("Entity configuration entityConfig.%s.memberUniqueIdAttribute is missing.",
						request.Entity.ExternalId),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}
			}

			request.Config.EntityConfigMap[request.Entity.ExternalId].SetOptionalDefaults()
		}
	}

	// Validate address scheme - only ldap:// and ldaps:// are allowed
	trimmedAddress, _, err := validation.ParseAndValidateAddress(request.Address, []string{"ldap", "ldaps"})
	if err != nil {
		return err
	}

	// Check if scheme is present in address, if not
	// set scheme based on certificateChain input
	sanitizedAddress := strings.ToLower(trimmedAddress)

	if !strings.HasPrefix(sanitizedAddress, "ldap://") && !strings.HasPrefix(sanitizedAddress, "ldaps://") {
		if request.Config.CertificateChain != "" {
			request.Address = "ldaps://" + trimmedAddress
		} else {
			request.Address = "ldap://" + trimmedAddress
		}
	}

	if request.Auth == nil || request.Auth.Basic == nil ||
		request.Auth.Basic.Username == "" || request.Auth.Basic.Password == "" {
		return &framework.Error{
			Message: "Provided datasource auth is missing required Active Directory authorization credentials.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// Validate that at least the unique ID attribute for the requested entity
	// is requested.
	uniqueIDAttributeFound := getUniqueIDAttribute(request.Entity.Attributes)

	if uniqueIDAttributeFound == nil {
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

// getUniqueIDAttribute checks if uniqueIDAttribute exists and returns its name.
func getUniqueIDAttribute(attrConfig []*framework.AttributeConfig) *string {
	for _, config := range attrConfig {
		if config.UniqueId {
			return &config.ExternalId
		}
	}

	return nil
}
