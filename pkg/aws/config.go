// Copyright 2026 SGNL.ai, Inc.

package aws

import (
	"context"
	"errors"

	"github.com/sgnl-ai/adapters/pkg/config"
)

type ResourceAccount struct {
	// RoleARN is the ARN of the role to assume in the account.
	RoleARN string `json:"roleARN"`
}

// EntityConfig enables filtering of entities.
type EntityConfig struct {
	// PathPrefix is the path prefix to filter the entities.
	PathPrefix *string `json:"pathPrefix,omitempty"`
}

// Config is the configuration passed in each GetPage calls to the adapter.
// AWS Adapter configuration example:
// nolint: godot
/*
{
  "resourceAccountRoles": [
    "arn:aws:iam::888111444333:role/Cross-Account-Assume-Admin",
    "arn:aws:iam::111111111111:role/Cross-Account-Assume-Admin"
  ],
  "region": "us-west-2",
  "requestTimeoutSeconds": 120
}
*/
type Config struct {
	// Common configuration
	*config.CommonConfig

	// Region is the AWS region to query.
	Region string `json:"region"`

	// EntityConfig is a map containing the config required for each entity associated with this
	EntityConfig map[string]*EntityConfig `json:"entityConfig,omitempty"`

	// ResourceAccountRoles is a list of roleARNs.
	ResourceAccountRoles []string `json:"resourceAccountRoles,omitempty"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	switch {
	case c == nil:
		return errors.New("The request contains an empty configuration")
	case c.Region == "":
		return errors.New("The AWS Region is not set in the configuration")
	default:
		return nil
	}
}
