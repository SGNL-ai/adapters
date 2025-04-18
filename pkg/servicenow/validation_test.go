// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst
package servicenow_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	servicenow_adapter "github.com/sgnl-ai/adapters/pkg/servicenow"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[servicenow_adapter.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[servicenow_adapter.Config]{
				Address: "test-instance.service-now.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_request_empty_config": {
			request: &framework.Request[servicenow_adapter.Config]{
				Address: "test-instance.service-now.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config:   &servicenow_adapter.Config{},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Servicenow config is invalid: apiVersion is not set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_config": {
			request: &framework.Request[servicenow_adapter.Config]{
				Address: "test-instance.service-now.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Servicenow config is invalid: request contains no config.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_api_version": {
			request: &framework.Request[servicenow_adapter.Config]{
				Address: "test-instance.service-now.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Servicenow config is invalid: apiVersion is not supported: v1.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_http_prefix": {
			request: &framework.Request[servicenow_adapter.Config]{
				Address: "http://test-instance.service-now.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "The provided HTTP protocol is not supported.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_page_size_too_large": {
			request: &framework.Request[servicenow_adapter.Config]{
				Address: "test-instance.service-now.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Ordered:  true,
				PageSize: 10001,
			},
			wantErr: &framework.Error{
				Message: "Provided page size (10001) is greater than the allowed maximum (10000).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"invalid_missing_auth": {
			request: &framework.Request[servicenow_adapter.Config]{
				Address: "test-instance.service-now.com",
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "System of Record is missing required authentication credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_basic_auth": {
			request: &framework.Request[servicenow_adapter.Config]{
				Address: "test-instance.service-now.com",
				Auth:    &framework.DatasourceAuthCredentials{},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "System of Record is missing required authentication credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_bearer_prefix": {
			request: &framework.Request[servicenow_adapter.Config]{
				Address: "test-instance.service-now.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "SSWS token", // Expected: `Bearer token`
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: `Provided auth token is missing required "Bearer " prefix.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_username": {
			request: &framework.Request[servicenow_adapter.Config]{
				Address: "test-instance.service-now.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "One of username or password required for basic auth is empty.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_password": {
			request: &framework.Request[servicenow_adapter.Config]{
				Address: "test-instance.service-now.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "One of username or password required for basic auth is empty.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_unique_attribute": {
			request: &framework.Request[servicenow_adapter.Config]{
				Address: "test-instance.service-now.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Requested entity attributes are missing unique ID attribute.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
	}

	adapter := &servicenow_adapter.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
