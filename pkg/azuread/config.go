// Copyright 2025 SGNL.ai, Inc.
package azuread

import (
	"context"
	"errors"
	"fmt"

	"github.com/sgnl-ai/adapters/pkg/config"
)

var supportedAPIVersions = map[string]struct{}{
	"v1.0": {},
}

// Config is the configuration passed in each GetPage calls to the adapter.
// Adapter configuration example:
// nolint: godot
/*
{
    "requestTimeoutSeconds": 10,
    "localTimeZoneOffset": 43200,
    "apiVersion": "v1.0",
    "filters": {
        "users": "displayName ne null"
    },
    "applyFiltersToMembers": true
}
*/
type Config struct {
	// Common configuration
	*config.CommonConfig

	// APIVersion is the version of the Microsoft Graph API to use.
	APIVersion string `json:"apiVersion,omitempty"`

	// Filters contains a map of filters for each entity associated with this
	// datasource. The key is the entity's external_name, and the value is the filter string.
	Filters map[string]string `json:"filters,omitempty"`

	// ApplyFiltersToMembers determines whether filters applied to parent objects should be applied when querying
	// member objects. For example, if true, Group filters be applied when pulling GroupMembers.
	ApplyFiltersToMembers bool `json:"applyFiltersToMembers,omitempty"`

	// Optional advanced filters to apply to the request.
	// See advanced_filters.go for more information.
	AdvancedFilters *AdvancedFilters `json:"advancedFilters,omitempty"`
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
