// Copyright 2026 SGNL.ai, Inc.

// nolint: goconst

package aws_test

import (
	"reflect"
	"strconv"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	aws_adapter "github.com/sgnl-ai/adapters/pkg/aws"
)

func TestValidateGetPageRequest(t *testing.T) {
	resourceAccounts := make([]string, 0)

	for i := 0; i < aws_adapter.MaxResourceAccounts+20; i++ {
		resourceAccounts = append(resourceAccounts, "arn:aws:iam::"+strconv.Itoa(i)+":role/role-name")
	}

	tests := map[string]struct {
		request *framework.Request[aws_adapter.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "User",
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
				Config:   validCommonConfig,
				Ordered:  false,
				PageSize: 2,
			},
			wantErr: nil,
		},
		"invalid_request_empty_config": {
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "User",
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
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "AWS config is invalid: The request contains an empty configuration.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_region": {
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "User",
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
				Config:   &aws_adapter.Config{},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "AWS config is invalid: The AWS Region is not set in the configuration.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_ordered_true": {
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "User",
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
				Config: &aws_adapter.Config{
					Region: "us-west-2",
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
			request: &framework.Request[aws_adapter.Config]{
				Entity: framework.EntityConfig{
					ExternalId: "User",
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
				Config: &aws_adapter.Config{
					Region: "us-west-2",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required AWS authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_entity_type": {
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "Invalid",
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
				Config: &aws_adapter.Config{
					Region: "us-west-2",
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
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "UserName",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &aws_adapter.Config{
					Region: "us-west-2",
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
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "User",
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
				Config: &aws_adapter.Config{
					Region: "us-west-2",
				},
				Ordered:  false,
				PageSize: 1001,
			},
			wantErr: &framework.Error{
				Message: "Provided page size (1001) exceeds the maximum allowed (1000).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"invalid_request_more_than_100_resource_accounts": {
			request: &framework.Request[aws_adapter.Config]{
				Auth: validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "User",
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
				Ordered:  false,
				PageSize: 100,
				Config: &aws_adapter.Config{
					Region:               "us-west-2",
					ResourceAccountRoles: resourceAccounts,
				},
			},
			wantErr: &framework.Error{
				Message: "Provided number of resource accounts (120) exceeds the maximum allowed limit: (100).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
	}

	adapter := &aws_adapter.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
