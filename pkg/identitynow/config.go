// Copyright 2026 SGNL.ai, Inc.
package identitynow

import (
	"context"
	"errors"
	"fmt"

	"github.com/sgnl-ai/adapters/pkg/config"
)

var supportedAPIVersions = map[string]struct{}{
	"v3":   {},
	"beta": {},
}

// Config is the configuration passed in each GetPage calls to the adapter.
// Adapter configuration example:
// nolint: godot
/*
{
	"requestTimeoutSeconds": 10,
	"localTimeZoneOffset": 43200,
	"apiVersion": "v3",
	"entityConfig": {
		"accounts": {
			"uniqueIDAttribute": "id",
			"filter": "identityId eq \"1700926aca594aa2861d9dbd24ca64b9\"",
			"apiVersion": "v3"
		}
	}
}
*/
type Config struct {
	*config.CommonConfig

	// APIVersion is the default API version to use for a request. Required.
	APIVersion string `json:"apiVersion,omitempty"`

	// EntityConfig is a map of configs for each entity associated with this datasource.
	// The key is the entity's external ID, and the value is a map of config values.
	EntityConfig map[string]EntityConfig `json:"entityConfig,omitempty"`
}

type EntityConfig struct {
	// UniqueIDAttribute is the name of the attribute containing the unique ID of
	// each returned object for the requested entity. Required.
	UniqueIDAttribute string `json:"uniqueIDAttribute"`
	// Filter is the filter to apply to the entity request. Optional.
	Filter *string `json:"filter,omitempty"`
	// APIVersion is the API version to use for the entity request. Optional.
	// If set, overrides the default Config.APIVersion.
	APIVersion *string `json:"apiVersion,omitempty"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	switch {
	case c == nil:
		return errors.New("request contains no config")
	case c.APIVersion == "":
		return errors.New("apiVersion is not set")
	case c.EntityConfig == nil:
		return errors.New("request contains no entityConfig")
	default:
		if _, ok := supportedAPIVersions[c.APIVersion]; !ok {
			return fmt.Errorf("apiVersion %s is not supported", c.APIVersion)
		}

		// Loop through each key in the entity config and validate:
		// 1) The entity config is not empty.
		// 2) The uniqueIDAttribute is not empty.
		// 3) The filter is not empty if set.
		// 4) The apiVersion is supported if set.
		for entityExternalID, entityConfig := range c.EntityConfig {
			// TODO [sc-19213]: Collect errors first and then return them to avoid client having to
			// iteratively debug.
			if entityConfig == (EntityConfig{}) {
				return fmt.Errorf("entityConfig for entity %v cannot be empty", entityExternalID)
			}

			if entityConfig.UniqueIDAttribute == "" {
				return fmt.Errorf("uniqueIDAttribute for entity %s cannot be empty", entityExternalID)
			}

			if entityConfig.APIVersion != nil {
				if _, ok := supportedAPIVersions[*entityConfig.APIVersion]; !ok {
					return fmt.Errorf("apiVersion %s for entity %s is not supported", *entityConfig.APIVersion, entityExternalID)
				}
			}
		}

		return nil
	}
}
