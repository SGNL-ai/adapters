// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package github_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	github "github.com/sgnl-ai/adapters/pkg/github"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[github.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("testenterpriseslug"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: nil,
		},
		"invalid_request_empty_config": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config:   &github.Config{},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "GitHub config is invalid: either enterpriseSlug or organizations must be specified.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_config": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "GitHub config is invalid: request contains no config.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"valid_config_empty_deployment_type": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug: testutil.GenPtr("SGNL"),
					APIVersion:     testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: nil,
		},
		"invalid_config_empty_api_version_for_REST": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "SecretScanningAlert",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "number",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        nil,
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "GitHub config is invalid: apiVersion is not set for an entity that is retrieve through the GitHub REST API.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_config_invalid_api_version_for_REST": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "SecretScanningAlert",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "number",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v5"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "GitHub config is invalid: apiVersion is not supported: v5.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_config_empty_api_version_for_GraphQL": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        nil,
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: nil,
		},
		"invalid_config_invalid_api_version_for_GraphQL": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("SGNL"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v5"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: nil,
		},
		"invalid_ordered_true": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("testenterpriseid"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  true,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "Ordered must be set to false.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"valid_https_prefix": {
			request: &framework.Request[github.Config]{
				Address: "https://ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("testenterpriseid"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: nil,
		},
		"invalid_http_prefix": {
			request: &framework.Request[github.Config]{
				Address: "http://ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("testenterpriseid"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "The provided HTTP protocol is not supported.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_auth": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("testenterpriseid"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required http authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_http_auth": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth:    &framework.DatasourceAuthCredentials{},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("testenterpriseid"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required http authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_bearer_prefix": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("testenterpriseid"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: `Provided auth token is missing required "Bearer " prefix.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_entity_type": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "invalid",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("testenterpriseid"),
					IsEnterpriseCloud: true,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "Provided entity external ID is invalid.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_missing_unique_attribute": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "login",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("testenterpriseid"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "Requested entity attributes are missing unique ID attribute.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_orguser_missing_one_of_two_required_attributes": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "OrganizationUser",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "uniqueId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "orgId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "login",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("testenterpriseid"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "Requested entity attributes are missing required attribute: $.node.id",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_page_size_too_big": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("testenterpriseid"),
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 1000,
			},
			wantErr: &framework.Error{
				Message: "Provided page size (1000) exceeds the maximum allowed (100).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"valid_organizations_list_request": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					Organizations:     []string{"org1", "org2"},
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: nil,
		},
		"invalid_both_organizations_and_enterpriseSlug": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					EnterpriseSlug:    testutil.GenPtr("testenterpriseslug"),
					Organizations:     []string{"org1", "org2"},
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "GitHub config is invalid: only one of enterpriseSlug or organizations must be specified, not both.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_neither_organizations_and_enterpriseSlug_are_present": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "GitHub config is invalid: either enterpriseSlug or organizations must be specified.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_empty_enterpriseSlug_present": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					IsEnterpriseCloud: false,
					EnterpriseSlug:    testutil.GenPtr(""),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "GitHub config is invalid: enterpriseSlug must be specified.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_organizations_and_empty_enterpriseSlug_are_present": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
					EnterpriseSlug:    testutil.GenPtr(""),
					Organizations:     []string{"org1", "org2"},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "GitHub config is invalid: only one of enterpriseSlug or organizations must be specified, not both.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_organizations_has_empty_value": {
			request: &framework.Request[github.Config]{
				Address: "ghe-test-server/api/graphql",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Organization",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &github.Config{
					Organizations:     []string{"org1", "", "org3"},
					IsEnterpriseCloud: false,
					APIVersion:        testutil.GenPtr("v3"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "organizations[1] cannot be an empty string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
	}

	adapter := &github.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
