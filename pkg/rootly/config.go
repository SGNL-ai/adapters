// Copyright 2025 SGNL.ai, Inc.
package rootly

import (
	"context"
	"errors"
	"fmt"

	"github.com/sgnl-ai/adapters/pkg/config"
)

var supportedAPIVersions = map[string]struct{}{
	"v1": {},
}

// Config is the configuration passed in each GetPage calls to the adapter.
// Adapter configuration example:
// nolint: godot
/*
{
    "requestTimeoutSeconds": 10,
    "apiVersion": "v1",
    "filters": {
        "users": "email=rufus_raynor@hegmann.test",
        "incidents": "status=started&severity=high"
    }
}
*/
type Config struct {
	// Common configuration
	*config.CommonConfig

	// APIVersion is the Rootly API version to use for requests.
	APIVersion string `json:"apiVersion,omitempty"`

	// Filters contains a map of filters for each entity associated with this
	// datasource. The key is the entity's external_name, and the value is the filter string.
	Filters map[string]string `json:"filters,omitempty"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	switch {
	case c == nil:
		return errors.New("request contains no config")
	case c.APIVersion == "":
		return errors.New("apiVersion is not set")
	default:
		if _, found := supportedAPIVersions[c.APIVersion]; !found {
			return fmt.Errorf("apiVersion is not supported: %v", c.APIVersion)
		}

		return nil
	}
}
