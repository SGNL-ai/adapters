// Copyright 2026 SGNL.ai, Inc.

package ldap

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// SetFilters configures the LDAP search filters based on the inputs received in entityConfig.
func SetFilters(request *Request) (string, *framework.Error) {
	query := request.EntityConfigMap[request.EntityExternalID].Query

	_, err := ldap.CompileFilter(query)
	if err != nil {
		return "", &framework.Error{
			Message: fmt.Sprintf("entityConfig.%s.query is not a valid LDAP filter.",
				request.EntityExternalID),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	return query, nil
}
