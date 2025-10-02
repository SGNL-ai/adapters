// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst
package ldap_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	ldap_adapter "github.com/sgnl-ai/adapters/pkg/ldap"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[ldap_adapter.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[ldap_adapter.Config]{
				Address: mockLDAPSAddr,
				Auth:    validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							ExternalId: "objectGUID",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Config:   validCommonConfig,
				Ordered:  false,
				PageSize: 2,
			},
			wantErr: nil,
		},

		"invalid_ordered_true": {
			request: &framework.Request[ldap_adapter.Config]{
				Address: mockLDAPSAddr,
				Auth:    validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							ExternalId: "objectGUID",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Config:   validCommonConfig,
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Ordered must be set to false.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_missing_auth": {
			request: &framework.Request[ldap_adapter.Config]{
				Address: mockLDAPAddr,
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							ExternalId: "objectGUID",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Config:   validCommonConfig,
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required Active Directory authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_http_auth": {
			request: &framework.Request[ldap_adapter.Config]{
				Address: mockLDAPSAddr,
				Auth:    &framework.DatasourceAuthCredentials{},
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							ExternalId: "objectGUID",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Config:   validCommonConfig,
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required Active Directory authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_http_auth_basic": {
			request: &framework.Request[ldap_adapter.Config]{
				Address: mockLDAPAddr,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							ExternalId: "objectGUID",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Config:   validCommonConfig,
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required Active Directory authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_http_auth_basic_username": {
			request: &framework.Request[ldap_adapter.Config]{
				Address: mockLDAPSAddr,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Password: "asdasd",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "objectGUID",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Config:   validCommonConfig,
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required Active Directory authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_http_auth_basic_password": {
			request: &framework.Request[ldap_adapter.Config]{
				Address: mockLDAPSAddr,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "asdasd",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "objectGUID",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Config:   validCommonConfig,
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required Active Directory authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_unique_attribute": {
			request: &framework.Request[ldap_adapter.Config]{
				Address: mockLDAPSAddr,
				Auth:    validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sAMAccountName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "objectGUID",
							Type:       framework.AttributeTypeString,
							List:       false,
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
		"invalid_missing_entityConfig": {
			request: &framework.Request[ldap_adapter.Config]{
				Address: mockLDAPAddr,
				Auth:    validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sAMAccountName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "objectGUID",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Config: &ldap_adapter.Config{
					BaseDN: validCommonConfig.BaseDN,
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"User": {
							Query: "(&(objectClass=user))",
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "entityConfig is missing in config for requested entity Person.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_page_size_too_big": {
			request: &framework.Request[ldap_adapter.Config]{
				Address: mockLDAPSAddr,
				Auth:    validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
					},
				},
				Config:   validCommonConfig,
				Ordered:  false,
				PageSize: 1000,
			},
			wantErr: &framework.Error{
				Message: "Provided page size (1000) exceeds the maximum allowed (999).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"invalid_entityConfig_missing_memberOf_entity": {
			request: &framework.Request[ldap_adapter.Config]{
				Address: mockLDAPAddr,
				Auth:    validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							ExternalId: "objectGUID",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Config: &ldap_adapter.Config{
					BaseDN: "dc=corp,dc=example,dc=io",
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"GroupMember": {
							MemberOf:                  testutil.GenPtr("Group"),
							Query:                     "(memberof={{CollectionID}})",
							MemberUniqueIDAttribute:   testutil.GenPtr("memberDistingushedName"),
							MemberOfUniqueIDAttribute: testutil.GenPtr("memberOfDistingushedName"),
						},
					},
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Entity configuration entityConfig.Group is missing for " +
					"entity specified in entityConfig.GroupMember.memberOf.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_entityConfig_missing_memberUniqueIDAttribute": {
			request: &framework.Request[ldap_adapter.Config]{
				Address: mockLDAPAddr,
				Auth:    validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							ExternalId: "objectGUID",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Config: &ldap_adapter.Config{
					BaseDN: "dc=corp,dc=example,dc=io",
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"Group": {
							Query: "(objectClass=groupofuniquenames)",
						},
						"GroupMember": {
							MemberOf:                  testutil.GenPtr("Group"),
							Query:                     "(&(memberOf={{CollectionId}}))",
							MemberOfUniqueIDAttribute: testutil.GenPtr("memberOfDistingushedName"),
						},
					},
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Entity configuration entityConfig.GroupMember.memberUniqueIdAttribute is missing.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_entityConfig_missing_memberOfUniqueIDAttribute": {
			request: &framework.Request[ldap_adapter.Config]{
				Address: mockLDAPAddr,
				Auth:    validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							ExternalId: "objectGUID",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Config: &ldap_adapter.Config{
					BaseDN: "dc=corp,dc=example,dc=io",
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"Group": {
							Query: "(objectClass=groupofuniquenames)",
						},
						"GroupMember": {
							MemberOf:                testutil.GenPtr("Group"),
							Query:                   "(&(memberOf={{CollectionId}}))",
							MemberUniqueIDAttribute: testutil.GenPtr("memberDistingushedName"),
						},
					},
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Entity configuration entityConfig.GroupMember.memberOfUniqueIdAttribute is missing.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
	}

	adapter := &ldap_adapter.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
