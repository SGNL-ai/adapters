// Copyright 2026 SGNL.ai, Inc.
package github

import (
	"context"
	"errors"
	"fmt"

	"github.com/sgnl-ai/adapters/pkg/config"
)

var supportedAPIVersions = map[string]bool{
	"v3": true,
}

// Config is the configuration passed in each GetPage calls to the adapter.
// Adapter configuration example:
// nolint: godot
/*
{
	"enterpriseSlug": "SGNL_ENTERPRISE",
	"organizations": [
		"sgnl-demos",
		"wholesalechips"
	],
	"isEnterpriseCloud": true,
	"apiVersion": "v3"
}
*/
type Config struct {
	// Common configuration
	*config.CommonConfig

	// EnterpriseSlug is the enterprise slug to query. This is the top level entity for every Github query.
	EnterpriseSlug *string `json:"enterpriseSlug,omitempty"`

	// Organizations is the list of organizations to query. Either this field or EnterpriseSlug must be set (but not both).
	Organizations []string `json:"organizations,omitempty"`

	// isEnterpriseCloud is a boolean that indicates whether the deployment is GitHub Enterprise Cloud.
	// This is used to determine the base URL to use.
	// If true, the deployment type is Enterprise Cloud. If false, the deployment type is Enterprise Server.
	IsEnterpriseCloud bool `json:"isEnterpriseCloud"`

	// APIVersion is the version of the GitHub API to use.
	// This is only used when constructing REST endpoints.
	APIVersion *string `json:"apiVersion"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context, isRestAPI bool) error {
	switch {
	case c == nil:
		return errors.New("request contains no config")
	case c.EnterpriseSlug != nil && *c.EnterpriseSlug == "" && len(c.Organizations) == 0:
		return errors.New("enterpriseSlug must be specified")
	case c.EnterpriseSlug == nil && len(c.Organizations) == 0:
		return errors.New("either enterpriseSlug or organizations must be specified")
	case c.EnterpriseSlug != nil && len(c.Organizations) > 0:
		return errors.New("only one of enterpriseSlug or organizations must be specified, not both")
	case c.APIVersion == nil && isRestAPI:
		return errors.New("apiVersion is not set for an entity that is retrieve through the GitHub REST API")
	case isRestAPI && !supportedAPIVersions[*c.APIVersion]:
		return fmt.Errorf("apiVersion is not supported: %s", *c.APIVersion)
	}

	return nil
}
