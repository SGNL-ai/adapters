// Copyright 2026 SGNL.ai, Inc.

package crowdstrike

import (
	"context"
	"errors"

	"github.com/sgnl-ai/adapters/pkg/config"
)

var (
	SupportedAPIVersions = map[string]struct{}{
		"v1": {},
	}
)

// Example Config:
//
//	{
//	   "apiVersion": "v1",
//	   "archived": false,
//	   "enabled": true,
//	   "filters": {
//	       "endpoint_protection_device": "platform:'Windows'"
//	   }
//	}
//
// Config is the optional configuration passed in each GetPage calls to the
// adapter.
type Config struct {
	// Common configuration
	*config.CommonConfig

	APIVersion string            `json:"apiVersion,omitempty"`
	Archived   bool              `json:"archived,omitempty"`
	Enabled    bool              `json:"enabled,omitempty"`
	Filters    map[string]string `json:"filters,omitempty"`
}

// Validate ensures that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	switch {
	case c == nil:
		return errors.New("The request contains an empty configuration")
	case c.APIVersion == "":
		return errors.New("apiVersion is not set in the configuration")
	case c.APIVersion != "":
		if _, found := SupportedAPIVersions[c.APIVersion]; !found {
			return errors.New("apiVersion is not supported")
		}
	default:
		return nil
	}

	return nil
}
