// Copyright 2025 SGNL.ai, Inc.
package mysql

import (
	"context"
	"errors"

	"github.com/sgnl-ai/adapters/pkg/config"
)

// Config is the configuration passed in each GetPage calls to the adapter.
//
// WARNING: In order to allow flexibility with filtering, we do NOT validate the provided filters. This effectively
// means that any user with access to configure this SoR has access to execute any command on the database due to
// SQL injection.
//
// Reasonable precautions should be made to ensure only authorized users have permission to configure this SoR / the
// provided credentials only have the required privileges assigned (e.g. read access scoped to a specific table).
//
// Adapter configuration example:
// nolint: godot
/*
{
	"requestTimeoutSeconds": 10,
	"localTimeZoneOffset": 43200,
	"database": "sgnl",
	"filters": {
		"users": "(age > 18 AND country = 'USA') OR verified = TRUE",
		"groups": "country IN ('USA', 'Canada')"
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
