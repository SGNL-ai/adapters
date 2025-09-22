// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package ldap_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	ldap "github.com/sgnl-ai/adapters/pkg/ldap"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name                 string
		config               *ldap.Config
		isLdaps              bool
		expectedError        error
		expectedEntityConfig map[string]*ldap.EntityConfig
	}{
		{
			name: "missing_baseDN",
			config: &ldap.Config{
				BaseDN: "",
			},
			expectedError: errors.New("baseDN is not set"),
		},
		{
			name: "valid_adapter_configuration",
			config: &ldap.Config{
				BaseDN:           "example",
				CertificateChain: "certificate",
			},
			expectedError:        nil,
			expectedEntityConfig: ldap.DefaultEntityConfig(),
		},
		{
			name: "valid_adapter_config_with_entityConfig",
			config: &ldap.Config{
				BaseDN:           "example",
				CertificateChain: "certificate",
				EntityConfigMap: map[string]*ldap.EntityConfig{
					"User": {
						Query: "(&(objectCategory=user)(objectClass=user))",
					},
				},
			},
			expectedEntityConfig: map[string]*ldap.EntityConfig{
				"User": {
					Query: "(&(objectCategory=user)(objectClass=user))",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate(context.Background())
			if (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError == nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("Validate() error = %v, expectedError %v", err, tt.expectedError)
			}

			if tt.expectedEntityConfig != nil {
				if ok := reflect.DeepEqual(tt.expectedEntityConfig, tt.config.EntityConfigMap); !ok {
					t.Errorf("Validate() entityConfig mismatched")
				}
			}
		})
	}
}

func TestDefaultMemberAttribute(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "DefaultEntityConfig_uses_defaultMemberAttribute",
			testFunc: func(t *testing.T) {
				config := ldap.DefaultEntityConfig()

				// Verify that GroupMember entity has the default member attribute
				groupMember, exists := config["GroupMember"]
				if !exists {
					t.Fatal("GroupMember entity should exist in default config")
				}
				if groupMember.MemberAttribute == nil {
					t.Fatal("MemberAttribute should not be nil")
				}
				if *groupMember.MemberAttribute != "member" {
					t.Errorf("MemberAttribute should equal 'member', got %s", *groupMember.MemberAttribute)
				}
			},
		},
		{
			name: "EntityConfig_SetOptionalDefaults_for_GroupMember_sets_defaults",
			testFunc: func(t *testing.T) {
				entityConfig := &ldap.EntityConfig{
					MemberOf: func() *string {
						s := "Group"

						return &s
					}(),
				}
				entityConfig.SetOptionalDefaults()

				// After calling SetOptionalDefaults, MemberAttribute should be set to "member"
				if entityConfig.MemberAttribute == nil {
					t.Fatal("MemberAttribute should be set after SetOptionalDefaults")
				}
				if *entityConfig.MemberAttribute != "member" {
					t.Errorf("MemberAttribute should be set to 'member', got %s", *entityConfig.MemberAttribute)
				}
				if entityConfig.MemberOfGroupBatchSize != 10 {
					t.Errorf("MemberOfGroupBatchSize should be set to 10, got %d", entityConfig.MemberOfGroupBatchSize)
				}
				if entityConfig.Query != "(&(objectClass=group)({{CollectionAttribute}}={{CollectionId}}))" {
					t.Errorf("Query should be set to default value, got %s", entityConfig.Query)
				}
				if entityConfig.CollectionAttribute == nil || *entityConfig.CollectionAttribute != "distinguishedName" {
					t.Errorf("CollectionAttribute should be set to 'distinguishedName', got %v", entityConfig.CollectionAttribute)
				}
			},
		},
		{
			name: "EntityConfig_SetOptionalDefaults_preserves_existing_memberAttribute",
			testFunc: func(t *testing.T) {
				customMemberAttr := "uniqueMember"
				entityConfig := &ldap.EntityConfig{
					MemberAttribute: &customMemberAttr,
					MemberOf: func() *string {
						s := "Group"

						return &s
					}(),
				}
				entityConfig.SetOptionalDefaults()

				// SetOptionalDefaults should not override existing MemberAttribute
				if *entityConfig.MemberAttribute != "uniqueMember" {
					t.Errorf("Custom MemberAttribute should be preserved, got %s", *entityConfig.MemberAttribute)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}
