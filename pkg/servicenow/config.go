// Copyright 2026 SGNL.ai, Inc.

package servicenow

import (
	"context"
	"errors"
	"fmt"

	"github.com/sgnl-ai/adapters/pkg/config"
)

var supportedAPIVersions = map[string]struct{}{
	"v2": {},
}

// Config is the configuration passed in each GetPage calls to the adapter.
// Adapter configuration example:
// nolint: godot
/*
{
    "requestTimeoutSeconds": 10,
    "localTimeZoneOffset": 43200,
    "apiVersion": "v2",
    "filters": {
        "incident": "active=true^priority=1"
    }
}
*/
type Config struct {
	// Common configuration
	*config.CommonConfig

	// APIVersion is the Servicenow API version to use for requests.
	APIVersion string `json:"apiVersion,omitempty"`

	// Filters contains a map of filters for each entity associated with this
	// datasource. The key is the entity's external_name, and the value is the filter string.
	Filters map[string]string `json:"filters,omitempty"`

	// Optional advanced filters to apply to the request.
	// See advanced_filters.go for more information.
	AdvancedFilters *AdvancedFilters `json:"advancedFilters,omitempty"`

	// CustomURLPath is an optional custom URL path to use instead of the default /api/now path.
	// If not specified, the default "/api/now" path will be used.
	CustomURLPath string `json:"customURLPath,omitempty"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	if c == nil {
		return errors.New("request contains no config")
	}

	// Only validate apiVersion if it's supplied
	if c.APIVersion != "" {
		if _, found := supportedAPIVersions[c.APIVersion]; !found {
			return fmt.Errorf("apiVersion is not supported: %v", c.APIVersion)
		}
	}

	return nil
}
