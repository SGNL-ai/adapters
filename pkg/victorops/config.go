// Copyright 2026 SGNL.ai, Inc.

package victorops

import (
	"context"
	"fmt"
	"net/url"

	"github.com/sgnl-ai/adapters/pkg/config"
)

// Config is the optional configuration passed in each GetPage calls to the
// adapter.
// Adapter configuration example:
// nolint: godot
/*
{
    "requestTimeoutSeconds": 10,
    "localTimeZoneOffset": 0,
    "queryParameters": {
        "IncidentReport": "currentPhase=RESOLVED&startedAfter=2024-01-01T00:00Z"
    }
}
*/
type Config struct {
	// Common configuration.
	*config.CommonConfig

	// QueryParameters contains a map of query parameter strings for each entity associated
	// with this datasource. The key is the entity's external ID (e.g. "IncidentReport"),
	// and the value is a URL query string (e.g. "currentPhase=RESOLVED&startedAfter=2024-01-01T00:00Z").
	// These parameters are appended to the API request URL for the matching entity.
	// Invalid parameter values may cause the SoR to return a 400.
	// It is up to the client to ensure the query parameter values are valid.
	QueryParameters map[string]string `json:"queryParameters,omitempty"`
}

// Validate validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	// Reserved query parameter keys that are managed by the adapter for pagination.
	reservedKeys := map[string]bool{
		"offset": true,
		"limit":  true,
	}

	for entity, params := range c.QueryParameters {
		if params == "" {
			return fmt.Errorf("queryParameters[%s] cannot be an empty string", entity)
		}

		parsed, err := url.ParseQuery(params)
		if err != nil {
			return fmt.Errorf("queryParameters[%s] is not a valid query string: %w", entity, err)
		}

		for key := range parsed {
			if reservedKeys[key] {
				return fmt.Errorf("queryParameters[%s] contains reserved parameter %q which is managed by the adapter", entity, key)
			}
		}
	}

	return nil
}
