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
