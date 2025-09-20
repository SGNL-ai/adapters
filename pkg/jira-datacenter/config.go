// Copyright 2025 SGNL.ai, Inc.
package jiradatacenter

import (
	"context"
	"errors"
	"fmt"

	"github.com/sgnl-ai/adapters/pkg/config"
)

// Config is the optional configuration passed in each GetPage calls to the
// adapter.
// Adapter configuration example:
// nolint: godot
/*
{
    "requestTimeoutSeconds": 10,
    "localTimeZoneOffset": 43200,
    "apiVersion": "latest",
    "groups": [
        "jira-administrators",
        "jira-users"
    ],
    "issuesJqlFilter": "project=SGNL OR project=MVP"
}
*/

type Config struct {
	// Common configuration
	*config.CommonConfig

	// APIVersion specifies which API version of JIRA Datacenter to use.
	// Must be either "2" or "latest".
	// Version "1" is not supported as Jira API returns 404 for this version.
	// If not specified, the adapter will use the "latest" version.
	APIVersion string `json:"apiVersion,omitempty"`

	// Groups is a list of group names to filter on. For each group in this list,
	// only matching groups will be synchronized, and users are synced based on membership
	// in the listed groups. If this field is empty or nil, the adapter will sync
	// all available groups and their members.
	Groups []string `json:"groups,omitempty"`

	// IssuesJQLFilter is the JQL filter to use when querying for issues.
	// e.g. "project=SGNL OR project=MVP".
	// https://developer.atlassian.com/server/jira/platform/rest/v10005/api-group-search/#api-api-2-search-get.
	// If the JQL is invalid, Jira will return a 400.
	// An invalid JQL does not necessarily mean a syntax error, but also for example
	// if a project does not exist, e.g. project=INVALID.
	// Therefore, it's up to the client to ensure the JQL is valid.
	IssuesJQLFilter *string `json:"issuesJqlFilter,omitempty"`

	// IncludeInactiveUsers determines whether inactive users are included in the results
	// when querying for User or GroupMember entities. If not specified or set to false,
	// only active users are returned.
	IncludeInactiveUsers *bool `json:"includeInactiveUsers,omitempty"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	// Validate APIVersion if provided
	if c.APIVersion != "" {
		// Only "2" and "latest" are supported versions for Jira Data Center API
		// Version "1" is not supported as Jira API returns 404 for this version
		if c.APIVersion != "2" && c.APIVersion != "latest" {
			return fmt.Errorf("apiVersion must be either '2' or 'latest', got '%s'", c.APIVersion)
		}
	}

	// Validate Groups configuration if provided
	if len(c.Groups) > 0 {
		if len(c.Groups) > 255 {
			return errors.New("too many groups specified; maximum allowed is 255")
		}

		groupMap := make(map[string]bool, len(c.Groups))

		for i, group := range c.Groups {
			if group == "" {
				return fmt.Errorf("group at index '%d' cannot be an empty string", i)
			}

			if len(group) > 255 {
				return fmt.Errorf("group name at index '%d' exceeds the 255 character limit", i)
			}

			// Check for duplicate group names
			if groupMap[group] {
				return fmt.Errorf("duplicate group name '%s' found", group)
			}

			groupMap[group] = true
		}
	}

	// The IssuesJQLFilter is optional so only validate if it's set.
	if c.IssuesJQLFilter != nil {
		if *c.IssuesJQLFilter == "" {
			return errors.New("issuesJqlFilter cannot be an empty string")
		}

		// Jira docs don't specify a max length; use a conservative estimate.
		if len(*c.IssuesJQLFilter) > 1024 {
			return errors.New("issuesJqlFilter exceeds the 1024 character limit")
		}
	}

	return nil
}
