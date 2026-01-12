// Copyright 2025 SGNL.ai, Inc.
package salesforce_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	salesforce_adapter "github.com/sgnl-ai/adapters/pkg/salesforce"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[salesforce_adapter.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[salesforce_adapter.Config]{
				Address: "sgnl-dev.my.salesforce.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
					Filters: map[string]string{
						"User": "Name = 'John'",
					},
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_request_empty_config": {
			request: &framework.Request[salesforce_adapter.Config]{
				Address: "sgnl-dev.my.salesforce.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config:   &salesforce_adapter.Config{},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Salesforce config is invalid: apiVersion is not set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_config": {
			request: &framework.Request[salesforce_adapter.Config]{
				Address: "sgnl-dev.my.salesforce.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Salesforce config is invalid: request contains no config.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_api_version": {
			request: &framework.Request[salesforce_adapter.Config]{
				Address: "sgnl-dev.my.salesforce.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "50.0",
					Filters: map[string]string{
						"User": "Name = 'John'",
					},
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Salesforce config is invalid: apiVersion is not supported: 50.0.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_ordered_false": {
			request: &framework.Request[salesforce_adapter.Config]{
				Address: "sgnl-dev.my.salesforce.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
					Filters: map[string]string{
						"User": "Name = 'John'",
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Ordered must be set to true.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_page_size_too_large": {
			request: &framework.Request[salesforce_adapter.Config]{
				Address: "sgnl-dev.my.salesforce.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
					Filters: map[string]string{
						"User": "Name = 'John'",
					},
				},
				Ordered:  true,
				PageSize: 2001,
			},
			wantErr: &framework.Error{
				Message: "Provided page size (2001) does not fall within the allowed range (200-2000).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"invalid_page_size_too_small": {
			request: &framework.Request[salesforce_adapter.Config]{
				Address: "sgnl-dev.my.salesforce.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
					Filters: map[string]string{
						"User": "Name = 'John'",
					},
				},
				Ordered:  true,
				PageSize: 199,
			},
			wantErr: &framework.Error{
				Message: "Provided page size (199) does not fall within the allowed range (200-2000).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"invalid_missing_auth": {
			request: &framework.Request[salesforce_adapter.Config]{
				Address: "sgnl-dev.my.salesforce.com",
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
					Filters: map[string]string{
						"User": "Name = 'John'",
					},
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required http authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_http_auth": {
			request: &framework.Request[salesforce_adapter.Config]{
				Address: "sgnl-dev.my.salesforce.com",
				Auth:    &framework.DatasourceAuthCredentials{},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
					Filters: map[string]string{
						"User": "Name = 'John'",
					},
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required http authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_unique_attribute": {
			request: &framework.Request[salesforce_adapter.Config]{
				Address: "sgnl-dev.my.salesforce.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
					Filters: map[string]string{
						"User": "Name = 'John'",
					},
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Requested entity attributes are missing unique ID attribute.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_auth_token_missing_bearer_prefix": {
			request: &framework.Request[salesforce_adapter.Config]{
				Address: "sgnl-dev.my.salesforce.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
					Filters: map[string]string{
						"User": "Name = 'John'",
					},
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: `Provided auth token is missing required "Bearer " prefix.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_relationship_depth_6_levels": {
			request: &framework.Request[salesforce_adapter.Config]{
				Address: "sgnl-dev.my.salesforce.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Case",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "$.Account.Parent.Parent.Parent.Parent.Name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Attribute '$.Account.Parent.Parent.Parent.Parent.Name' exceeds the maximum " +
					"relationship depth of 5 levels. Salesforce SOQL supports up to 5 levels of " +
					"child-to-parent relationship traversal.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
	}

	adapter := &salesforce_adapter.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
