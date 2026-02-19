// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst
package workday_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/workday"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[workday.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[workday.Config]{
				Address: "test-instance.workday.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Worker",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "$.worker.id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "employeeNumber",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "testorgid",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_request_empty_config": {
			request: &framework.Request[workday.Config]{
				Address: "test-instance.workday.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Worker",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "$.worker.id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "employeeNumber",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config:   &workday.Config{},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Workday config is invalid: apiVersion is not set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_config_apiversion_missing": {
			request: &framework.Request[workday.Config]{
				Address: "test-instance.workday.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Worker",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "$.worker.id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "employeeNumber",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &workday.Config{
					OrganizationID: "testorgid",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Workday config is invalid: apiVersion is not set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_config_organizationid_missing": {
			request: &framework.Request[workday.Config]{
				Address: "test-instance.workday.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Worker",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "$.worker.id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "employeeNumber",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &workday.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Workday config is invalid: organizationId is not set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_config": {
			request: &framework.Request[workday.Config]{
				Address: "test-instance.workday.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Worker",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "$.worker.id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "employeeNumber",
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
				Message: "Workday config is invalid: request contains no config.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_api_version": {
			request: &framework.Request[workday.Config]{
				Address: "test-instance.workday.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Worker",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "$.worker.id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "employeeNumber",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &workday.Config{
					APIVersion:     "v1.1",
					OrganizationID: "testorgid",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Workday config is invalid: apiVersion is not supported: v1.1.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"valid_https_prefix": {
			request: &framework.Request[workday.Config]{
				Address: "https://test-instance.workday.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Worker",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "$.worker.id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "employeeNumber",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "testorgid",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_http_prefix": {
			request: &framework.Request[workday.Config]{
				Address: "http://test-instance.workday.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Worker",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "$.worker.id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "employeeNumber",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "testorgid",
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
			request: &framework.Request[workday.Config]{
				Address: "https://test-instance.workday.com",
				Entity: framework.EntityConfig{
					ExternalId: "Worker",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "$.worker.id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "employeeNumber",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "testorgid",
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
			request: &framework.Request[workday.Config]{
				Address: "https://test-instance.workday.com",
				Auth:    &framework.DatasourceAuthCredentials{},
				Entity: framework.EntityConfig{
					ExternalId: "Worker",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "$.worker.id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "employeeNumber",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "testorgid",
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
			request: &framework.Request[workday.Config]{
				Address: "https://test-instance.workday.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Worker",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "$.worker.id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "employeeNumber",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "testorgid",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: `Provided auth token is missing required "Bearer " prefix.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_page_size_too_big": {
			request: &framework.Request[workday.Config]{
				Address: "https://test-instance.workday.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Worker",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "$.worker.id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "employeeNumber",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "testorgid",
				},
				Ordered:  false,
				PageSize: 1001,
			},
			wantErr: &framework.Error{
				Message: "Provided page size (1001) exceeds the maximum allowed (1000).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
	}

	adapter := &workday.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
