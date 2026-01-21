// Copyright 2026 SGNL.ai, Inc.

package ldap

import (
	"context"
	"errors"

	"github.com/sgnl-ai/adapters/pkg/config"
)

const (
	defaultMemberAttribute   = "member"
	defaultDistinguishedName = "distinguishedName"
)

// Config is the configuration passed in each GetPage calls to the adapter.
// Adapter configuration example:
// nolint: godot
/*
{
    "baseDN": "dc=org,dc=example,dc=io",
    "certificateChain": "....",
    "entityConfig": {
        "User": {
            "query": "(&(objectCategory=user)(objectClass=user)(distinguishedName=*))",
        },
        "Group": {
            "query": "(&(objectCategory=group)(objectClass=group)(distinguishedName=*))",
        },
        "GroupMember": {
            "memberOf": "Group",
            "query": "(&(objectClass=group)({{CollectionAttribute}}={{CollectionId}}))", // Users who are members
            "collectionAttribute": "distinguishedName",
			"memberUniqueIdAttribute": "memberDistinguishedName",
			"memberOfUniqueIdAttribute": "groupDistinguishedName",
			"memberAttribute": defaultMemberAttribute,
			"memberOfGroupBatchSize": 10,
        }
    }
}
*/

// EntityConfig holds attributes which are used to create LDAP search filter.
type EntityConfig struct {
	Query                     string  `json:"query"`
	CollectionAttribute       *string `json:"collectionAttribute,omitempty"`
	MemberUniqueIDAttribute   *string `json:"memberUniqueIdAttribute,omitempty"`
	MemberOfUniqueIDAttribute *string `json:"memberOfUniqueIdAttribute,omitempty"`
	MemberOf                  *string `json:"memberOf,omitempty"`
	MemberAttribute           *string `json:"memberAttribute,omitempty"`
	MemberOfGroupBatchSize    int64   `json:"memberOfGroupBatchSize,omitempty"`
}

type Config struct {
	// Common configuration
	*config.CommonConfig

	BaseDN string `json:"baseDN"`

	// CertificateChain is a base64 encoded Certificates
	CertificateChain string `json:"certificateChain,omitempty"`

	// EntityConfigMap is an map containing the config required for each entity associated with this
	// datasource. The key is the entity's external_name and value is EntityConfig.
	EntityConfigMap map[string]*EntityConfig `json:"entityConfig"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	if c.EntityConfigMap == nil {
		c.EntityConfigMap = DefaultEntityConfig()
	}

	switch {
	case c == nil:
		return errors.New("request contains no config")
	case c.BaseDN == "":
		return errors.New("baseDN is not set")
	default:
		return nil
	}
}

// defaultEntityConfig: if entityConfig is nil, defaulting to values that pull data from ActiveDirectory.
func DefaultEntityConfig() map[string]*EntityConfig {
	entityConfig := map[string]*EntityConfig{
		"User": {
			Query: "(&(objectCategory=user)(objectClass=user)(distinguishedName=*))",
		},
		"Group": {
			Query: "(&(objectCategory=group)(objectClass=group)(distinguishedName=*))",
		},
		"Computer": {
			Query: "(&(objectCategory=computer)(name=*))",
		},
		"GroupMember": {
			Query: "(&(objectClass=group)({{CollectionAttribute}}={{CollectionId}}))",
			MemberOf: func() *string {
				s := "Group"

				return &s
			}(),
			CollectionAttribute: func() *string {
				s := defaultDistinguishedName

				return &s
			}(),
			MemberUniqueIDAttribute: func() *string {
				s := defaultDistinguishedName

				return &s
			}(),
			MemberOfUniqueIDAttribute: func() *string {
				s := defaultDistinguishedName

				return &s
			}(),
			MemberOfGroupBatchSize: 10,
			MemberAttribute: func() *string {
				s := defaultMemberAttribute

				return &s
			}(),
		},
	}

	return entityConfig
}

// SetOptionalDefaults sets the default values for optional fields in EntityConfig.
func (e *EntityConfig) SetOptionalDefaults() {
	if e.MemberOf != nil {
		if e.MemberOfGroupBatchSize == 0 {
			e.MemberOfGroupBatchSize = 10
		}

		if e.MemberAttribute == nil {
			s := defaultMemberAttribute
			e.MemberAttribute = &s
		}

		if e.Query == "" {
			e.Query = "(&(objectClass=group)({{CollectionAttribute}}={{CollectionId}}))"
		}

		if e.CollectionAttribute == nil {
			s := defaultDistinguishedName
			e.CollectionAttribute = &s
		}

		if e.MemberUniqueIDAttribute == nil {
			s := defaultDistinguishedName
			e.MemberUniqueIDAttribute = &s
		}

		if e.MemberOfUniqueIDAttribute == nil {
			s := defaultDistinguishedName
			e.MemberOfUniqueIDAttribute = &s
		}
	}
}
