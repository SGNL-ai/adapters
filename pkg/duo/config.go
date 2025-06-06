// Copyright 2025 SGNL.ai, Inc.
package duo

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
// Duo Adapter configuration example:
// nolint: godot
/*
{
    "requestTimeoutSeconds": 10,
    "localTimeZoneOffset": 43200,
    "apiVersion": "v1"
}
*/
type Config struct {
	// Common configuration
	*config.CommonConfig

	// APIVersion is the version of the Duo API to use.
	APIVersion string `json:"apiVersion,omitempty"`
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
