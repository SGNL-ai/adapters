// Copyright 2026 SGNL.ai, Inc.

package mysql

import (
	"context"
	"errors"

	"github.com/sgnl-ai/adapters/pkg/condexpr"
	"github.com/sgnl-ai/adapters/pkg/config"
)

// Config is the configuration passed in each GetPage calls to the adapter.
//
// Adapter configuration example:
// nolint: godot
/*
{
	"requestTimeoutSeconds": 10,
	"localTimeZoneOffset": 43200,
	"database": "sgnl",
	"filters": {
		"users": {
			"or": [
				{
					"and": [
						{
							"field": "age",
							"op": ">",
							"value": 18
						},
						{
							"field": "country",
							"op": "=",
							"value": "USA"
						}
					]
				},
				{
					"field": "verified",
					"op": "=",
					"value": true
				}
			]
		},
		"groups": {
			"field": "country",
			"op": "IN",
			"value": ["active", "inactive"]
		}
	}
}
*/
type Config struct {
	*config.CommonConfig

	// MySQL database to connect to.
	Database string `json:"database,omitempty"`

	Filters map[string]condexpr.Condition `json:"filters,omitempty"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	switch {
	case c == nil:
		return errors.New("request contains no config")
	case c.Database == "":
		return errors.New("database is not set")
	default:
		return nil
	}
}
