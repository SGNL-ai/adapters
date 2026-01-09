// Copyright 2026 SGNL.ai, Inc.
package salesforce

import (
	"context"
	"errors"
	"fmt"

	"github.com/sgnl-ai/adapters/pkg/config"
)

var supportedAPIVersions = map[string]struct{}{
	"52.0": {},
	"53.0": {},
	"54.0": {},
	"55.0": {},
	"56.0": {},
	"57.0": {},
	"58.0": {},
}

// Config is the configuration passed in each GetPage calls to the adapter.
// Adapter configuration example:
// nolint: godot
/*
{
    "requestTimeoutSeconds": 10,
    "apiVersion": "58.0",
    "filters": {
        "User": "isActive=true",
        "Case": "isClosed=false"
    }
}
*/
type Config struct {
	// Common configuration
	*config.CommonConfig

	// APIVersion is the Salesforce API version to use for requests.
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
