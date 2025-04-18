// Copyright 2025 SGNL.ai, Inc.
package pagerduty

import (
	"context"
	"fmt"

	"github.com/sgnl-ai/adapters/pkg/config"
)

// Config is the optional configuration passed in each GetPage call to the adapter.
// Adapter configuration example:
// nolint: godot
/*
{
    "requestTimeoutSeconds": 10,
    "localTimeZoneOffset": 43200,
    "additionalQueryParameters": {
        "users": {
            "query": "user",
            "include[]": ["contact_methods", "notification_rules", "teams"]
        }
    }
}
*/
type Config struct {
	// Common configuration
	*config.CommonConfig

	// AdditionalQueryParameters is a map of entity external ID to the entity's query parameters
	// that are added to the request.
	// In some cases, invalid values for query parameters will cause the SoR to return a 400.
	// e.g. /users&team_ids[]=INVALID_ID.
	// Therefore, it's up to the client to ensure the query parameter values produce valid results.
	AdditionalQueryParameters map[string]map[string]any `json:"additionalQueryParameters,omitempty"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	switch {
	case c != nil:
		// Validate that each query param value is not empty and is a string or []string.
		for entity, entityQueryParams := range c.AdditionalQueryParameters {
			for queryParam, value := range entityQueryParams {
				switch v := value.(type) {
				case string:
					if v == "" {
						return fmt.Errorf("additionalQueryParameters[%s][%s] is an empty string", entity, queryParam)
					}
				case []any:
					if len(v) == 0 {
						return fmt.Errorf("additionalQueryParameters[%s][%s] is an empty list", entity, queryParam)
					}

					for i, e := range v {
						s, ok := e.(string)

						switch {
						case !ok:
							return fmt.Errorf("additionalQueryParameters[%s][%s][%d] is not a string", entity, queryParam, i)
						case s == "":
							return fmt.Errorf("additionalQueryParameters[%s][%s][%d] is an empty string", entity, queryParam, i)
						}
					}
				default:
					return fmt.Errorf(
						"additionalQueryParameters[%s][%s] is neither a string nor a list of strings",
						entity,
						queryParam,
					)
				}
			}
		}

		return nil
	default:
		return nil
	}
}
