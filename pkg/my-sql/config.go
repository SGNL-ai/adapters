// Copyright 2025 SGNL.ai, Inc.
package mysql

import (
	"context"
	"errors"

	"github.com/sgnl-ai/adapters/pkg/config"
)

// Config is the configuration passed in each GetPage calls to the adapter.
// Adapter configuration example:
// nolint: godot
/*
{
"re	questTimeoutSeconds": 10,
	"localTimeZoneOffset": 43200,
	"database": "sgnl",
	"filters": {
		"users": "WHERE active = true",
		"groups": "WHERE x = y"
	}
}
*/
type Config struct {
	*config.CommonConfig

	// MySQL database to connect to.
	Database string `json:"database,omitempty"`

	Filters map[string]string `json:"filters,omitempty"`
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
