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
	if err := request.Config.Validate(ctx); err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("Active Directory config is invalid: %v.", err.Error()),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if _, _, err := validation.ParseAndValidateAddress(request.Address, []string{"ldap", "ldaps"}); err != nil {
		return err
	}

	// set scheme based on certificateChain input
	if request.Config.CertificateChain != "" {
		request.Address = "ldaps://" + request.Address
	} else {
		request.Address = "ldap://" + request.Address
	}

	if request.Auth == nil || request.Auth.Basic == nil ||
		request.Auth.Basic.Username == "" || request.Auth.Basic.Password == "" {
		return &framework.Error{
			Message: "Provided datasource auth is missing required Active Directory authorization credentials.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

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

		if !strings.Contains(request.Config.EntityConfigMap[request.Entity.ExternalId].Query, "{{CollectionId}}") {
			return &framework.Error{
				Message: fmt.Sprintf("{{CollectionId}} is missing in entityConfig.%s.query for Entity configuration.",
					request.Entity.ExternalId),
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
