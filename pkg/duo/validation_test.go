// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst
package duo_test

import (
	"context"
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	duo_adapter "github.com/sgnl-ai/adapters/pkg/duo"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[duo_adapter.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[duo_adapter.Config]{
				Address: "api-xxxxxxxx.duosecurity.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "user_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_request_empty_config": {
			request: &framework.Request[duo_adapter.Config]{
				Address: "api-xxxxxxxx.duosecurity.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "user_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config:   &duo_adapter.Config{},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Duo config is invalid: apiVersion is not set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_config": {
			request: &framework.Request[duo_adapter.Config]{
				Address: "api-xxxxxxxx.duosecurity.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "user_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "displayName",
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
				Message: "Duo config is invalid: request contains no config.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_api_version": {
			request: &framework.Request[duo_adapter.Config]{
				Config:  &duo_adapter.Config{APIVersion: "v1.1"},
				Address: "api-xxxxxxxx.duosecurity.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					}},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "user_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "displayName",
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
				Cursor:   "",
			},
			wantErr: &framework.Error{
				Message: "Duo config is invalid: apiVersion is not supported: v1.1.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_ordered_true": {
			request: &framework.Request[duo_adapter.Config]{
				Address: "api-xxxxxxxx.duosecurity.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "user_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &duo_adapter.Config{
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
			request: &framework.Request[duo_adapter.Config]{
				Address: "https://api-xxxxxxxx.duosecurity.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "user_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_http_prefix": {
			request: &framework.Request[duo_adapter.Config]{
				Address: "http://api-xxxxxxxx.duosecurity.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "user_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "The provided HTTP protocol is not supported.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_auth": {
			request: &framework.Request[duo_adapter.Config]{
				Address: "api-xxxxxxxx.duosecurity.com",
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "user_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required basic authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_basic_auth": {
			request: &framework.Request[duo_adapter.Config]{
				Address: "api-xxxxxxxx.duosecurity.com",
				Auth:    &framework.DatasourceAuthCredentials{},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "user_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required basic authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_basic_auth_details": {
			request: &framework.Request[duo_adapter.Config]{
				Address: "api-xxxxxxxx.duosecurity.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "",
						Password: "",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "user_id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required basic authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_entity_type": {
			request: &framework.Request[duo_adapter.Config]{
				Address: "api-xxxxxxxx.duosecurity.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "invalid",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &duo_adapter.Config{
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
			request: &framework.Request[duo_adapter.Config]{
				Address: "api-xxxxxxxx.duosecurity.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &duo_adapter.Config{
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
		"valid_child_entities": {
			request: &framework.Request[duo_adapter.Config]{
				Address: "api-xxxxxxxx.duosecurity.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test Integration Key",
						Password: "Test Secret",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "user_id",
							Type:       framework.AttributeTypeString,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "groups",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "group_id",
									Type:       framework.AttributeTypeString,
								},
							},
						},
					},
				},
				Config: &duo_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
		},
	}

	adapter := duo_adapter.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(context.TODO(), tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
