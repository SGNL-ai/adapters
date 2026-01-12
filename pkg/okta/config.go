// Copyright 2025 SGNL.ai, Inc.
package okta

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
    "localTimeZoneOffset": 43200,
    "apiVersion": "v1",
    "filters": {
        "User": "status eq \"ACTIVE\"",
        "Group": "type eq \"OKTA_GROUP\"",
        "Application": "status eq \"ACTIVE\""
    },
	"search": {
        "User": "profile.department eq \"Engineering\""
    }
}
*/
type Config struct {
	*config.CommonConfig

	APIVersion string            `json:"apiVersion,omitempty"`
	Filters    map[string]string `json:"filters,omitempty"`
	Search     map[string]string `json:"search,omitempty"`
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
