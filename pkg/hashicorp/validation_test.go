// Copyright 2025 SGNL.ai, Inc.

// nolint: lll
// nolint:forcetypeassert
package hashicorp_test

import (
	"context"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	hashicorp_adapter "github.com/sgnl-ai/adapters/pkg/hashicorp"
	"github.com/sgnl-ai/adapters/pkg/mock"
	"github.com/sgnl-ai/adapters/pkg/validation"
	"github.com/stretchr/testify/assert"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request            *framework.Request[hashicorp_adapter.Config]
		inputSSRFValidator validation.SSRFValidator
		wantErr            *framework.Error
	}{
		"nil_request": {
			request: nil,
			wantErr: &framework.Error{
				Message: "Request is nil",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"nil_config": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Config: nil,
			},
			wantErr: &framework.Error{
				Message: "Request config is nil",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"http_protocol": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Address: "http://example.com",
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "test-auth-method-id",
				},
			},
			wantErr: &framework.Error{
				Message: "The provided HTTP protocol is not supported.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"missing_auth": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "test-auth-method-id",
				},
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required http authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"nil_attributes": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "test-auth-method-id",
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "test",
						Password: "test",
					},
				},
				Entity: framework.EntityConfig{
					Attributes: nil,
				},
			},
			wantErr: &framework.Error{
				Message: "Request entity attributes is nil",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"missing_id_attribute": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "test-auth-method-id",
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "test",
						Password: "test",
					},
				},
				Entity: framework.EntityConfig{
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "Requested entity attributes are missing unique ID attribute.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_page_size": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "test-auth-method-id",
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "test",
						Password: "test",
					},
				},
				Entity: framework.EntityConfig{
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 5,
			},
			wantErr: &framework.Error{
				Message: "Provided page size (5) does not fall within the allowed range (10-10000).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"valid_request": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "test-auth-method-id",
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "test",
						Password: "test",
					},
				},
				Entity: framework.EntityConfig{
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 100,
			},
			wantErr: nil,
		},
		"valid_request_with_cursor": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "test-auth-method-id",
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "test",
						Password: "test",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "hosts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 100,
			},
			wantErr: nil,
		},
		"invalid_cursor": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "test-auth-method-id",
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "test",
						Password: "test",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "hosts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 100,
				Cursor:   "invalid-base64",
			},
			wantErr: &framework.Error{
				Message: "Failed to decode base64 cursor: illegal base64 data at input byte 7.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"invalid_localhost_address": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "test-auth-method-id",
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "test",
						Password: "test",
					},
				},
				Entity: framework.EntityConfig{
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Address:  "https://localhost:8080",
				PageSize: 100,
			},
			inputSSRFValidator: validation.NewDefaultSSRFValidator(),
			wantErr: &framework.Error{
				Message: `Address URL validation failed: localhost URLs are not allowed: "https://localhost:8080".`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_private_ip_address": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "test-auth-method-id",
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "test",
						Password: "test",
					},
				},
				Entity: framework.EntityConfig{
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Address:  "https://192.168.1.1",
				PageSize: 100,
			},
			inputSSRFValidator: validation.NewDefaultSSRFValidator(),
			wantErr: &framework.Error{
				Message: `Address URL validation failed: private IP addresses are not allowed for "https://192.168.1.1": 192.168.1.1.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_aws_metadata_address": {
			request: &framework.Request[hashicorp_adapter.Config]{
				Config: &hashicorp_adapter.Config{
					AuthMethodID: "test-auth-method-id",
				},
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "test",
						Password: "test",
					},
				},
				Entity: framework.EntityConfig{
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Address:  "https://169.254.169.254",
				PageSize: 100,
			},
			inputSSRFValidator: validation.NewDefaultSSRFValidator(),
			wantErr: &framework.Error{
				Message: `Address URL validation failed: private IP addresses are not allowed for "https://169.254.169.254": 169.254.169.254.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testSSRFValidator := mock.NewNoOpSSRFValidator()
			if tt.inputSSRFValidator != nil {
				testSSRFValidator = tt.inputSSRFValidator
			}

			adapter := &hashicorp_adapter.Adapter{
				HashicorpClient: nil,
				SSRFValidator:   testSSRFValidator,
			}
			gotErr := adapter.ValidateGetPageRequest(context.Background(), tt.request)

			if tt.wantErr == nil {
				assert.Nil(t, gotErr)

				return
			}

			if gotErr == nil && tt.wantErr != nil {
				t.Fatalf("expected error %v, got nil", tt.wantErr)
			}

			assert.Equal(t, tt.wantErr.Message, gotErr.Message)
			assert.Equal(t, tt.wantErr.Code, gotErr.Code)
		})
	}
}
