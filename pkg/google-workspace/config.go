// Copyright 2026 SGNL.ai, Inc.

package googleworkspace

import (
	"context"
	"errors"
	"fmt"

	"github.com/sgnl-ai/adapters/pkg/config"
)

var supportedAPIVersions = map[string]struct{}{
	"v1": {},
}

var supportedRoles = map[string]struct{}{
	"OWNER":   {},
	"MANAGER": {},
	"MEMBER":  {},
}

// Config is the configuration passed in each GetPage calls to the adapter.
// Google Workspace Adapter configuration example:
// nolint: godot
/*
{
    "requestTimeoutSeconds": 10,
    "localTimeZoneOffset": 43200,
    "apiVersion": "v1",
	"domain": "sgnldemos.com",
	"filters": {
		"user": {
			"showDeleted": true,
			"query": ""
		},
		"group": {
			"query": ""
		},
		"member": {
			"includeDerivedMembership": false,
			"roles": "MEMBER"
		}
	}
}
*/
type Config struct {
	// Common configuration
	*config.CommonConfig

	// APIVersion is the version of the Google Workspace API to use.
	APIVersion string `json:"apiVersion"`

	// Note: Only one of Customer or Domain should be set.
	// The unique ID for the customer's Google Workspace account. In case of a multi-domain account,
	// to fetch all entities for a customer, use this field instead of domain.
	Customer *string `json:"customer"`

	// Note: Only one of Customer or Domain should be set.
	// Use this field to get entities from only one domain.
	Domain *string `json:"domain"`

	Filters Filters `json:"filters"`
}

type Filters struct {
	UserFilters   *UserFilters   `json:"user"`
	GroupFilters  *GroupFilters  `json:"group"`
	MemberFilters *MemberFilters `json:"member"`
}

type UserFilters struct {
	// If set to true, retrieves the list of deleted users. (Default: false)
	ShowDeleted bool `json:"showDeleted"`

	// Query is a string for searching entity fields. For more information:
	// Constructing user queries: https://developers.google.com/admin-sdk/directory/v1/guides/search-users
	Query *string `json:"query"`
}

type GroupFilters struct {
	// Query is a string for searching entity fields. For more information:
	// Constructing group queries: https://developers.google.com/admin-sdk/directory/v1/guides/search-groups
	Query *string `json:"query"`
}

type MemberFilters struct {
	// This determines whether to list indirect memberships. Default: false.
	IncludeDerivedMembership bool `json:"includeDerivedMembership"`

	// The roles query parameter allows you to retrieve group members by role.
	Roles *string `json:"roles"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	switch {
	case c == nil:
		return errors.New("request contains no config")
	case c.APIVersion == "":
		return errors.New("apiVersion is not set")
	case c.Customer == nil && c.Domain == nil:
		return errors.New("customer or domain must be set")
	case c.Customer != nil && c.Domain != nil:
		return errors.New("One of customer or domain must be set")
	default:
		if _, found := supportedAPIVersions[c.APIVersion]; !found {
			return fmt.Errorf("apiVersion is not supported: %v", c.APIVersion)
		}

		if c.Filters.MemberFilters != nil && c.Filters.MemberFilters.Roles != nil {
			if _, found := supportedRoles[*c.Filters.MemberFilters.Roles]; !found {
				return fmt.Errorf("filters.member.roles is set to an unsupported value: %v", *c.Filters.MemberFilters.Roles)
			}
		}

		return nil
	}
}
