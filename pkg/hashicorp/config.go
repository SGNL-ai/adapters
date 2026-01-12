// Copyright 2025 SGNL.ai, Inc.
package hashicorp

import (
	"context"
	"errors"

	"github.com/sgnl-ai/adapters/pkg/config"
)

// Config is the configuration passed in each GetPage calls to the adapter.
// Adapter configuration example:
// nolint: godot
/*
{
    "requestTimeoutSeconds": 10,
	"authMethodId": "test-auth-method-id",
	"entityConfig": {
		"hosts": {
			"scope_id": "global",
			"filter": ""
		}
	}
}
*/
type EntityConfig struct {
	ScopeID string `json:"scopeId,omitempty"`

	Filter string `json:"filter,omitempty"`
}

type Config struct {
	// Common configuration
	*config.CommonConfig

	AuthMethodID string `json:"authMethodId,omitempty"`

	EntityConfig map[string]EntityConfig `json:"entityConfig,omitempty"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	switch {
	case c == nil:
		return errors.New("request contains no config")
	case c.AuthMethodID == "":
		return errors.New("Auth method ID is not set in the configuration")
	default:
		return nil
	}
}
