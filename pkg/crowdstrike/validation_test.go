// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package crowdstrike_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	crowdstrike_adapter "github.com/sgnl-ai/adapters/pkg/crowdstrike"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[crowdstrike_adapter.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Auth:    validAuthCredentials,
				Address: validAddress,
				Entity: framework.EntityConfig{
					ExternalId: "user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "entityId",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "primaryDisplayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "archived",
							Type:       framework.AttributeTypeBool,
						},
						{
							ExternalId: "secondaryDisplayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "emailAddresses",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config:   validCommonConfig,
				Ordered:  false,
				PageSize: 2,
			},
			wantErr: nil,
		},
		"invalid_request_empty_config": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Auth:    validAuthCredentials,
				Address: validAddress,
				Entity: framework.EntityConfig{
					ExternalId: "user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "entityId",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "primaryDisplayName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "CrowdStrike config is invalid: The request contains an empty configuration.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_api_version": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Auth:    validAuthCredentials,
				Address: validAddress,
				Entity: framework.EntityConfig{
					ExternalId: "user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "entityId",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "primaryDisplayName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config:   &crowdstrike_adapter.Config{},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "CrowdStrike config is invalid: apiVersion is not set in the configuration.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_ordered_true": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Auth:    validAuthCredentials,
				Address: validAddress,
				Entity: framework.EntityConfig{
					ExternalId: "user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "entityId",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "primaryDisplayName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &crowdstrike_adapter.Config{
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
		"invalid_missing_auth": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Entity: framework.EntityConfig{
					ExternalId: "user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "UserName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_entity_type": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Auth:    validAuthCredentials,
				Address: validAddress,
				Entity: framework.EntityConfig{
					ExternalId: "Invalid",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "entityId",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "primaryDisplayName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &crowdstrike_adapter.Config{
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
			request: &framework.Request[crowdstrike_adapter.Config]{
				Auth:    validAuthCredentials,
				Address: validAddress,
				Entity: framework.EntityConfig{
					ExternalId: "user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "secondaryDisplayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryDisplayName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &crowdstrike_adapter.Config{
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
		"invalid_page_size": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Auth:    validAuthCredentials,
				Address: validAddress,
				Entity: framework.EntityConfig{
					ExternalId: "user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "entityId",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "primaryDisplayName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
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

	adapter := &crowdstrike_adapter.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
