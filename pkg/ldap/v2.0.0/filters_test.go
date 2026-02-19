// Copyright 2026 SGNL.ai, Inc.

// nolint: lll, goconst

package ldap_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	ldap_adapter "github.com/sgnl-ai/adapters/pkg/ldap/v2.0.0"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestSetFilters(t *testing.T) {
	tests := []struct {
		name            string
		request         *ldap_adapter.Request
		expectedFilters string
		expectedError   *framework.Error
	}{
		// user
		{
			name: "user_filter",
			request: &ldap_adapter.Request{
				EntityExternalID: "User",
				EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
					"User": {
						Query: "(&(objectCategory=user)(objectClass=user)(distinguishedName=*))",
					},
				},
			},
			expectedFilters: "(&(objectCategory=user)(objectClass=user)(distinguishedName=*))",
		},
		{
			name: "invalud_user_filter",
			request: &ldap_adapter.Request{
				EntityExternalID: "User",
				EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
					"User": {
						Query: "(&(objectCategory=user)(objectClass=user)(distinguishedName=*",
					},
				},
			},
			expectedError: &framework.Error{
				Message: "entityConfig.User.query is not a valid LDAP filter.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		// person
		{
			name: "person_without_filter",
			request: &ldap_adapter.Request{
				EntityExternalID: "Person",
				EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
					"Person": {
						Query: "(&(objectClass=person))",
					},
				},
			},
			expectedFilters: "(&(objectClass=person))",
		},
		// group
		{
			name: "group_filter",
			request: &ldap_adapter.Request{
				EntityExternalID: "Group",
				EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
					"Group": {
						Query: "(&(objectCategory=group)(objectClass=group)(distinguishedName=*)" +
							"(instanceType=4)(sAMAccountName=Blocked Sign-In))",
					},
				},
			},
			expectedFilters: "(&(objectCategory=group)(objectClass=group)(distinguishedName=*)" +
				"(instanceType=4)(sAMAccountName=Blocked Sign-In))",
		},
		// group member
		{
			name: "groupMember_filter",
			request: &ldap_adapter.Request{
				EntityExternalID: "GroupMember",
				EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
					"Group": {
						Query: "(&(objectCategory=group)(objectClass=group)(distinguishedName=*))",
					},
					"GroupMember": {
						MemberOf: testutil.GenPtr("Group"),
						// CollectionID already set by previous step
						Query:                     "(&(memberOf=cn=Administrator,ou=Groups,dc=example,dc=org)(logonCount=0))",
						MemberUniqueIDAttribute:   testutil.GenPtr("memberDistinguishedName"),
						MemberOfUniqueIDAttribute: testutil.GenPtr("groupDistinguishedName"),
					},
				},
			},
			expectedFilters: "(&(memberOf=cn=Administrator,ou=Groups,dc=example,dc=org)(logonCount=0))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFilters, gotErr := ldap_adapter.SetFilters(tt.request)
			if gotFilters != tt.expectedFilters {
				t.Errorf("Expected filters %s, but got %s", tt.expectedFilters, gotFilters)
			}

			if !reflect.DeepEqual(gotErr, tt.expectedError) {
				t.Errorf("Expected filtering error %s, but got %s", tt.expectedError, gotErr)
			}
		})
	}
}
