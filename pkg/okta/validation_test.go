// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst
package okta_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	okta_adapter "github.com/sgnl-ai/adapters/pkg/okta"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[okta_adapter.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[okta_adapter.Config]{
				Address: "test-instance.oktapreview.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "SSWS testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_request_empty_config": {
			request: &framework.Request[okta_adapter.Config]{
				Address: "test-instance.oktapreview.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "SSWS testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config:   &okta_adapter.Config{},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Okta config is invalid: apiVersion is not set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_config": {
			request: &framework.Request[okta_adapter.Config]{
				Address: "test-instance.oktapreview.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "SSWS testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Okta config is invalid: request contains no config.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_api_version": {
			request: &framework.Request[okta_adapter.Config]{
				Address: "test-instance.oktapreview.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "SSWS testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1.1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Okta config is invalid: apiVersion is not supported: v1.1.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_ordered_true": {
			request: &framework.Request[okta_adapter.Config]{
				Address: "test-instance.oktapreview.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "SSWS testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Ordered must be set to false.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"valid_https_prefix": {
			request: &framework.Request[okta_adapter.Config]{
				Address: "https://test-instance.oktapreview.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "SSWS testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_http_prefix": {
			request: &framework.Request[okta_adapter.Config]{
				Address: "http://test-instance.oktapreview.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "SSWS testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: `Scheme "http" is not supported.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_auth": {
			request: &framework.Request[okta_adapter.Config]{
				Address: "test-instance.oktapreview.com",
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required http authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_http_auth": {
			request: &framework.Request[okta_adapter.Config]{
				Address: "test-instance.oktapreview.com",
				Auth:    &framework.DatasourceAuthCredentials{},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required http authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_prefix": {
			request: &framework.Request[okta_adapter.Config]{
				Address: "test-instance.oktapreview.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: `Provided auth token is missing required "Bearer " or "SSWS " prefix.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_entity_type": {
			request: &framework.Request[okta_adapter.Config]{
				Address: "test-instance.oktapreview.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "invalid",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided entity external ID is invalid.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_missing_unique_attribute": {
			request: &framework.Request[okta_adapter.Config]{
				Address: "test-instance.oktapreview.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "SSWS testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Requested entity attributes are missing unique ID attribute.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_child_entities": {
			request: &framework.Request[okta_adapter.Config]{
				Address: "test-instance.oktapreview.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "SSWS testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "$.profile.firstName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "emails",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "email",
									Type:       framework.AttributeTypeString,
								},
							},
						},
					},
				},
				Config: &okta_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Requested entity does not support child entities.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
	}

	adapter := &okta_adapter.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
