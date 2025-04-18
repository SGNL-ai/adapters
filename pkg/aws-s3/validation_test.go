// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package awss3_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	s3_adapter "github.com/sgnl-ai/adapters/pkg/aws-s3"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[s3_adapter.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[s3_adapter.Config]{
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
			request: &framework.Request[s3_adapter.Config]{
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
				Message: "AWS config is invalid: the request contains an empty configuration.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_region": {
			request: &framework.Request[s3_adapter.Config]{
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
				Config:   &s3_adapter.Config{},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "AWS config is invalid: the AWS Region is not set in the configuration.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_unsupported_file_type": {
			request: &framework.Request[s3_adapter.Config]{
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
				Config: &s3_adapter.Config{
					Region:   "us-west-2",
					Bucket:   "bucket",
					FileType: testutil.GenPtr("json"),
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "The filetype json in config.fileType is not supported.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_auth": {
			request: &framework.Request[s3_adapter.Config]{
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
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required AWS authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_unique_attribute": {
			request: &framework.Request[s3_adapter.Config]{
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
				Config:   validCommonConfig,
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Requested entity attributes are missing unique ID attribute.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_page_size": {
			request: &framework.Request[s3_adapter.Config]{
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
				PageSize: 1001,
			},
			wantErr: &framework.Error{
				Message: "Provided page size (1001) exceeds the maximum allowed (1000).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
	}

	adapter := &s3_adapter.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
