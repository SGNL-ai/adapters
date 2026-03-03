// Copyright 2026 SGNL.ai, Inc.

// nolint: goconst
package identitynow_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	identitynow_adapter "github.com/sgnl-ai/adapters/pkg/identitynow"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[identitynow_adapter.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_request_empty_entity_config": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "IdentityNow config is invalid: request contains no entityConfig.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_config": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "IdentityNow config is invalid: request contains no config.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_entity_not_configured": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "UNCONFIGURED_ENTITY",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Entity with external ID UNCONFIGURED_ENTITY must be present in the Adapter Config for this SoR.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_page_size_too_large": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Ordered:  false,
				PageSize: 500,
			},
			wantErr: &framework.Error{
				Message: "PageSize must be less than or equal to 250.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_invalid_api_version": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "UNSUPPORTED_API_VERSION",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "IdentityNow config is invalid: apiVersion UNSUPPORTED_API_VERSION is not supported.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_api_version_not_set": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "IdentityNow config is invalid: apiVersion is not set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_invalid_entity_api_version": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
							APIVersion:        testutil.GenPtr[string]("UNSUPPORTED_API_VERSION"),
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "IdentityNow config is invalid: apiVersion UNSUPPORTED_API_VERSION for entity accounts is not supported.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_ordered_true": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Ordered must be set to false.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"valid_request_https_prefix": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "https://sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_request_http_prefix": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "http://sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: `Scheme "http" is not supported.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_auth_missing_bearer_prefix": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "NOT_BEARER token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: `Provided auth token is missing required "Bearer " prefix.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_auth": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required http authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_http_auth": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Auth:    &framework.DatasourceAuthCredentials{},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required http authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_unique_attribute": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "WRONG_ID",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Requested entity attributes are missing unique ID attribute.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_request_missing_entity_config": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "IdentityNow config is invalid: entityConfig for entity accounts cannot be empty.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_unique_id_in_entity_config": {
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "sgnl-dev-tenant.api.identitynow.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							APIVersion: testutil.GenPtr[string]("v3"),
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "IdentityNow config is invalid: uniqueIDAttribute for entity accounts cannot be empty.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
	}

	adapter := &identitynow_adapter.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
