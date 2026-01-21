// Copyright 2026 SGNL.ai, Inc.

// nolint: goconst
package azuread_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/azuread"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[azuread.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
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
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_request_empty_config": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
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
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config:   &azuread.Config{},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Azure AD config is invalid: apiVersion is not set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_config": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
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
				Message: "Azure AD config is invalid: request contains no config.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_api_version": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
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
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Azure AD config is invalid: apiVersion is not supported: v1.1.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_ordered_true": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
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
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
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
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
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
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_http_prefix": {
			request: &framework.Request[azuread.Config]{
				Address: "http://graph.microsoft.com",
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
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
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
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Entity: framework.EntityConfig{
					ExternalId: "User",
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
				Config: &azuread.Config{
					APIVersion: "v1.0",
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
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth:    &framework.DatasourceAuthCredentials{},
				Entity: framework.EntityConfig{
					ExternalId: "User",
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
				Config: &azuread.Config{
					APIVersion: "v1.0",
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
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
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
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
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
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
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
				Config: &azuread.Config{
					APIVersion: "v1.0",
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
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
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
				Config: &azuread.Config{
					APIVersion: "v1.0",
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
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
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
							ExternalId: "displayName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastLogin",
							Type:       framework.AttributeTypeString,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "emails",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "email",
									Type:       framework.AttributeTypeString,
								},
							},
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
				},
				Ordered:  false,
				PageSize: 250,
			},
		},
		"invalid_page_size_too_big": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
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
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
				},
				Ordered:  false,
				PageSize: 1000,
			},
			wantErr: &framework.Error{
				Message: "Provided page size (1000) exceeds the maximum allowed (999).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"valid_advanced_filters": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
					AdvancedFilters: &azuread.AdvancedFilters{
						ScopedObjects: map[string][]azuread.EntityFilter{
							azuread.GroupMember: {
								{
									ScopeEntity:       "Group",
									ScopeEntityFilter: "id in ('7df6bf7d-7b09-4399-9aed-0e345d1ea7b2')",
									Members: []azuread.MemberFilter{
										{
											MemberEntity:       "User",
											MemberEntityFilter: "department eq 'Architecture and Authentication'",
										},
										{
											MemberEntity: "Group",
										},
									},
								},
								{
									ScopeEntity:       "Group",
									ScopeEntityFilter: "id eq '94e98d95-8e04-47ee-b2d4-5a1bd96bf0a9'",
									Members: []azuread.MemberFilter{
										{
											MemberEntity: "User",
										},
									},
								},
								{
									ScopeEntity:       "Group",
									ScopeEntityFilter: "id eq '821d9845-3482-41c8-9a0e-d93f4d576731'",
									Members: []azuread.MemberFilter{
										{
											MemberEntity: "User",
										},
										{
											MemberEntity: "Group",
										},
									},
								},
							},
						},
					},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: nil,
		},
		"advanced_filters_valid_if_non_conflicting_regular_filter_exists": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
					Filters: map[string]string{
						// non-conflicting entity external ID `users`
						"users": "department eq 'Architecture and Authentication'",
					},
					AdvancedFilters: &azuread.AdvancedFilters{
						ScopedObjects: map[string][]azuread.EntityFilter{
							azuread.GroupMember: {
								{
									ScopeEntity:       "Group",
									ScopeEntityFilter: "id in ('7df6bf7d-7b09-4399-9aed-0e345d1ea7b2')",
									Members: []azuread.MemberFilter{
										{
											MemberEntity:       "User",
											MemberEntityFilter: "department eq 'Architecture and Authentication'",
										},
										{
											MemberEntity: "Group",
										},
									},
								},
								{
									ScopeEntity:       "Group",
									ScopeEntityFilter: "id eq '94e98d95-8e04-47ee-b2d4-5a1bd96bf0a9'",
									Members: []azuread.MemberFilter{
										{
											MemberEntity: "User",
										},
									},
								},
								{
									ScopeEntity:       "Group",
									ScopeEntityFilter: "id eq '821d9845-3482-41c8-9a0e-d93f4d576731'",
									Members: []azuread.MemberFilter{
										{
											MemberEntity: "User",
										},
										{
											MemberEntity: "Group",
										},
									},
								},
							},
						},
					},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: nil,
		},
		"advanced_filters_invalid_if_entityExternalId_not_supported": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: azuread.Role, // Valid external ID but advanced filters not supported
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
					AdvancedFilters: &azuread.AdvancedFilters{
						ScopedObjects: map[string][]azuread.EntityFilter{
							azuread.Role: { // Valid external ID but advanced filters not supported
								{
									ScopeEntity:       "Group",
									ScopeEntityFilter: "id in ('7df6bf7d-7b09-4399-9aed-0e345d1ea7b2')",
									Members: []azuread.MemberFilter{
										{
											MemberEntity: "User",
										},
									},
								},
							},
						},
					},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "Advanced Filters on advancedFilters.getObjectsByScope.Role is not supported.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"advanced_filters_invalid_without_entityExternalId": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
					AdvancedFilters: &azuread.AdvancedFilters{
						ScopedObjects: map[string][]azuread.EntityFilter{
							// GroupMember from the above example is missing
						},
					},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "advancedFilters.getObjectsByScope cannot be empty.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"advanced_filters_invalid_if_filters_is_empty": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
					AdvancedFilters: &azuread.AdvancedFilters{
						ScopedObjects: map[string][]azuread.EntityFilter{
							azuread.GroupMember: {}, // Must contain at least one filter
						},
					},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: `advancedFilters.getObjectsByScope.GroupMember must have at least one filter defined.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"advanced_filters_invalid_if_scopedEntity_is_empty": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
					AdvancedFilters: &azuread.AdvancedFilters{
						ScopedObjects: map[string][]azuread.EntityFilter{
							azuread.GroupMember: {
								{
									ScopeEntity:       "", // Want non-empty value
									ScopeEntityFilter: "id in ('7df6bf7d-7b09-4399-9aed-0e345d1ea7b2')",
									Members: []azuread.MemberFilter{
										{
											MemberEntity:       "User",
											MemberEntityFilter: "department eq 'Architecture and Authentication'",
										},
										{
											MemberEntity: "Group",
										},
									},
								},
							},
						},
					},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: `advancedFilters.getObjectsByScope.GroupMember.[0].scopeEntity cannot be empty.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"advanced_filters_invalid_if_scopedEntity_is_not_supported": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
					AdvancedFilters: &azuread.AdvancedFilters{
						ScopedObjects: map[string][]azuread.EntityFilter{
							azuread.GroupMember: {
								{
									ScopeEntity:       "User", // User is not a supported scope entity under `GroupMember`.
									ScopeEntityFilter: "id in ('7df6bf7d-7b09-4399-9aed-0e345d1ea7b2')",
									Members: []azuread.MemberFilter{
										{
											MemberEntity:       "User",
											MemberEntityFilter: "department eq 'Architecture and Authentication'",
										},
										{
											MemberEntity: "Group",
										},
									},
								},
							},
						},
					},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: `advancedFilters.getObjectsByScope.GroupMember.[0].scopeEntity is not supported.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"advanced_filters_invalid_if_scopedEntityFilter_is_empty": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
					AdvancedFilters: &azuread.AdvancedFilters{
						ScopedObjects: map[string][]azuread.EntityFilter{
							azuread.GroupMember: {
								{
									ScopeEntity:       "Group",
									ScopeEntityFilter: "", // Want non-empty value
									Members: []azuread.MemberFilter{
										{
											MemberEntity:       "User",
											MemberEntityFilter: "department eq 'Architecture and Authentication'",
										},
										{
											MemberEntity: "Group",
										},
									},
								},
							},
						},
					},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "advancedFilters.getObjectsByScope.GroupMember.[0].scopeEntityFilter cannot be empty.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"advanced_filters_invalid_if_scoped_entity_members_is_empty": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
					AdvancedFilters: &azuread.AdvancedFilters{
						ScopedObjects: map[string][]azuread.EntityFilter{
							azuread.GroupMember: {
								{
									ScopeEntity:       "Group",
									ScopeEntityFilter: "id eq '94e98d95-8e04-47ee-b2d4-5a1bd96bf0a9'",
									Members:           []azuread.MemberFilter{}, // Want non-empty value
								},
							},
						},
					},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "advancedFilters.getObjectsByScope.GroupMember.[0].members cannot be empty.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"advanced_filters_invalid_if_memberEntity_of_scoped_entity_members_is_empty": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
					AdvancedFilters: &azuread.AdvancedFilters{
						ScopedObjects: map[string][]azuread.EntityFilter{
							azuread.GroupMember: {
								{
									ScopeEntity:       "Group",
									ScopeEntityFilter: "id eq '94e98d95-8e04-47ee-b2d4-5a1bd96bf0a9'",
									Members: []azuread.MemberFilter{
										{
											// MemberEntity:       "", // want non-empty value
											MemberEntityFilter: "department eq 'Architecture and Authentication'",
										},
									},
								},
							},
						},
					},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "advancedFilters.getObjectsByScope.GroupMember.[0].members[0].memberEntity cannot be empty.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"advanced_filters_invalid_if_memberEntity_has_no_endpoint_generation_support": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
					AdvancedFilters: &azuread.AdvancedFilters{
						ScopedObjects: map[string][]azuread.EntityFilter{
							azuread.GroupMember: {
								{
									ScopeEntity:       "Group",
									ScopeEntityFilter: "id eq '94e98d95-8e04-47ee-b2d4-5a1bd96bf0a9'",
									Members: []azuread.MemberFilter{
										{
											MemberEntity:       azuread.Device, // This entity does not have suffix registered for endpoint generation
											MemberEntityFilter: "department eq 'Architecture and Authentication'",
										},
									},
								},
							},
						},
					},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "advancedFilters.getObjectsByScope.GroupMember.[0].members[0].memberEntity is not supported.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"advanced_filters_invalid_if_regular_filters_exist": {
			request: &framework.Request[azuread.Config]{
				Address: "https://graph.microsoft.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &azuread.Config{
					APIVersion: "v1.0",
					Filters: map[string]string{
						azuread.GroupMember: "department eq 'IT'",
					},
					AdvancedFilters: &azuread.AdvancedFilters{
						ScopedObjects: map[string][]azuread.EntityFilter{
							azuread.GroupMember: {
								{
									ScopeEntity:       "Group",
									ScopeEntityFilter: "id eq '94e98d95-8e04-47ee-b2d4-5a1bd96bf0a9'",
									Members: []azuread.MemberFilter{
										{
											// MemberEntity:       "", // want non-empty value
											MemberEntityFilter: "department eq 'Architecture and Authentication'",
										},
									},
								},
							},
						},
					},
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: `Only one of advancedFilters.getObjectsByScope.GroupMember OR ` +
					`filters.GroupMember is allowed.`,
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
	}

	adapter := &azuread.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
