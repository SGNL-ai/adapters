// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst, lll
package googleworkspace_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	googleworkspace "github.com/sgnl-ai/adapters/pkg/google-workspace"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[googleworkspace.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_request_empty_config": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config:   &googleworkspace.Config{},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Google Workspace adapter config is invalid: apiVersion is not set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_config": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Google Workspace adapter config is invalid: request contains no config.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_api_version": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1.1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Google Workspace adapter config is invalid: apiVersion is not supported: v1.1.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"missing_api_version": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					Domain: testutil.GenPtr("sgnldemos.com"),
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Google Workspace adapter config is invalid: apiVersion is not set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"missing_domain_and_customer": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Google Workspace adapter config is invalid: customer or domain must be set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"valid_only_customer_set": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Customer:   testutil.GenPtr("customer"),
				},
				Ordered:  false,
				PageSize: 250,
			},
		},
		"valid_only_one_filter_set": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Customer:   testutil.GenPtr("customer"),
					Filters: googleworkspace.Filters{
						UserFilters: &googleworkspace.UserFilters{
							ShowDeleted: true,
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
		},
		"invalid_members_filter_roles_enum": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Member",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Customer:   testutil.GenPtr("customer"),
					Filters: googleworkspace.Filters{
						MemberFilters: &googleworkspace.MemberFilters{
							Roles: testutil.GenPtr("INVALID"),
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Google Workspace adapter config is invalid: filters.member.roles is set to an unsupported value: INVALID.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"valid_https_prefix": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "https://admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_http_prefix": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "http://admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
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
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
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
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth:    &framework.DatasourceAuthCredentials{},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
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
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
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
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: `Provided auth token is missing required "Bearer " prefix.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_entity_type": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "INVALID",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided entity external ID is invalid.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_missing_unique_attribute_user": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Requested entity User is missing the required unique ID attribute: id",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_missing_unique_attribute_member": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Member",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "groupId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Requested entity Member is missing the required unique ID attribute: uniqueId",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"valid_child_entities": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "emails",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "address",
									Type:       framework.AttributeTypeString,
									List:       false,
								},
								{
									ExternalId: "type",
									Type:       framework.AttributeTypeString,
									List:       false,
								},
								{
									ExternalId: "primary",
									Type:       framework.AttributeTypeBool,
									List:       false,
								},
							},
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Ordered:  false,
				PageSize: 250,
			},
		},
		"invalid_page_size_too_big": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "User",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Ordered:  false,
				PageSize: 1000,
			},
			wantErr: &framework.Error{
				Message: "Requested page size, 1000, exceeds the maximum allowed value of 500 for entity: User.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"invalid_page_size_too_big_member_entity": {
			request: &framework.Request[googleworkspace.Config]{
				Address: "admin.googleapis.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Member",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "uniqueId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "groupId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "primaryEmail",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Ordered:  false,
				PageSize: 1500,
			},
			wantErr: &framework.Error{
				Message: "Requested page size, 1500, exceeds the maximum allowed value of 1000 for entity: Member.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
	}

	adapter := &googleworkspace.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
