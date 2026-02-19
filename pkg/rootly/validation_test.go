// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package rootly_test

import (
	"context"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/config"
	rootly_adapter "github.com/sgnl-ai/adapters/pkg/rootly"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestValidateGetPageRequest(t *testing.T) {
	adapter := &rootly_adapter.Adapter{}

	tests := map[string]struct {
		ctx         context.Context
		request     *framework.Request[rootly_adapter.Config]
		wantErrCode api_adapter_v1.ErrorCode
		wantErrMsg  string
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: "https://api.rootly.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "v1",
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(30),
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 100,
			},
		},
		"missing_config": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: "https://api.rootly.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: nil,
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 100,
			},
			wantErrCode: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrMsg:  "request contains no config",
		},
		"invalid_api_version": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: "https://api.rootly.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "2",
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 100,
			},
			wantErrCode: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrMsg:  "Rootly config is invalid: apiVersion is not supported: 2.",
		},
		"missing_api_version": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: "https://api.rootly.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "",
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 100,
			},
			wantErrCode: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrMsg:  "Rootly config is invalid: apiVersion is not set.",
		},
		"http_protocol_not_supported": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: "http://api.rootly.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 100,
			},
			wantErrCode: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrMsg:  `Scheme "http" is not supported.`,
		},
		"missing_auth": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: "https://api.rootly.com",
				Auth:    nil,
				Config: &rootly_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 100,
			},
			wantErrCode: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrMsg:  "Provided datasource auth is missing required http authorization credentials.",
		},
		"empty_http_authorization": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: "https://api.rootly.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 100,
			},
			wantErrCode: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrMsg:  "Provided datasource auth is missing required http authorization credentials.",
		},
		"missing_bearer_prefix": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: "https://api.rootly.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 100,
			},
			wantErrCode: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrMsg:  `Provided auth token is missing required "Bearer " prefix.`,
		},
		"missing_unique_id_attribute": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: "https://api.rootly.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "title",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 100,
			},
			wantErrCode: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			wantErrMsg:  "Requested entity attributes are missing unique ID attribute.",
		},
		"page_size_too_large": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: "https://api.rootly.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 1001,
			},
			wantErrCode: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			wantErrMsg:  "Provided page size (1001) does not fall within the allowed range (1-1000).",
		},
		"page_size_too_small": {
			ctx: context.Background(),
			request: &framework.Request[rootly_adapter.Config]{
				Address: "https://api.rootly.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &rootly_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "incidents",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 0,
			},
			wantErrCode: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			wantErrMsg:  "Provided page size (0) does not fall within the allowed range (1-1000).",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := adapter.ValidateGetPageRequest(tt.ctx, tt.request)

			if tt.wantErrCode == api_adapter_v1.ErrorCode_ERROR_CODE_UNSPECIFIED {
				if err != nil {
					t.Errorf("ValidateGetPageRequest() expected no error, got: %v", err)
				}

				return
			}

			if err == nil {
				t.Errorf("ValidateGetPageRequest() expected error with code %v, got no error", tt.wantErrCode)

				return
			}

			if err.Code != tt.wantErrCode {
				t.Errorf("ValidateGetPageRequest() error code = %v, want %v", err.Code, tt.wantErrCode)
			}

			if tt.wantErrMsg != "" && err.Message != tt.wantErrMsg {
				t.Errorf("ValidateGetPageRequest() error message = %v, want %v", err.Message, tt.wantErrMsg)
			}
		})
	}
}
