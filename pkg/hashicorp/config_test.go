// Copyright 2026 SGNL.ai, Inc.

package hashicorp_test

import (
	"context"
	"testing"

	hashicorp_adapter "github.com/sgnl-ai/adapters/pkg/hashicorp"
	"github.com/stretchr/testify/assert"
)

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *hashicorp_adapter.Config
		expectedErr string
	}{
		{
			name:        "nil_config",
			config:      nil,
			expectedErr: "request contains no config",
		},
		{
			name: "missing_auth_method_ID",
			config: &hashicorp_adapter.Config{
				AuthMethodID: "",
			},
			expectedErr: "Auth method ID is not set in the configuration",
		},
		{
			name: "valid_config",
			config: &hashicorp_adapter.Config{
				AuthMethodID: "test-auth-method-id",
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate(context.Background())
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedErr)
			}
		})
	}
}
