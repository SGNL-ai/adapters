// Copyright 2025 SGNL.ai, Inc.
package jira

import (
	"context"
	"errors"
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
    "localTimeZoneOffset": 43200,
    "issuesJqlFilter": "project=SGNL OR project=MVP",
    "objectsQlQuery": "objectType = Customer",
    "assetBaseUrl": "https://api.atlassian.com/jsm/assets"
}
*/
type Config struct {
	// Common configuration
	*config.CommonConfig

	// IssuesJQLFilter is the JQL filter to use when querying for issues.
	// e.g. "project=SGNL OR project=MVP".
	// https://support.atlassian.com/jira-software-cloud/docs/what-is-advanced-search-in-jira-cloud/.
	// If the JQL is invalid, Jira will return a 400.
	// An invalid JQL does not necessarily mean a syntax error, but also for example
	// if a project does not exist, e.g. project=INVALID.
	// Therefore, it's up to the client to ensure the JQL is valid.
	IssuesJQLFilter *string `json:"issuesJqlFilter,omitempty"`

	// ObjectsQLQuery is the AQL query to use when querying for custom Objects.
	// e.g. "qlQuery="objectType = Customer".
	// https://developer.atlassian.com/cloud/assets/rest/api-group-object/#api-object-aql-post.
	// It is up to the client to ensure the ObjectsQLQuery is valid.
	// This field is only used for the Object entity.
	ObjectsQLQuery *string `json:"objectsQlQuery,omitempty"`

	// AssetBaseURL is the base URL to use when querying for Objects.
	// e.g. "https://api.atlassian.com/jsm/assets".
	// If not, specified it defaults to "https://api.atlassian.com/jsm/assets".
	// This field is only used for the Object entity.
	AssetBaseURL *string `json:"assetBaseUrl,omitempty"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
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

	if c.ObjectsQLQuery != nil {
		if *c.ObjectsQLQuery == "" {
			return errors.New("objectsQlQuery cannot be an empty string")
		}

		if len(*c.ObjectsQLQuery) > 1024 {
			return errors.New("objectsQlQuery exceeds the 1024 character limit")
		}
	}

	if c.AssetBaseURL != nil {
		if _, err := url.ParseRequestURI(*c.AssetBaseURL); err != nil {
			return fmt.Errorf("assetBaseUrl is not a valid URL: %w", err)
		}
	}

	return nil
}
