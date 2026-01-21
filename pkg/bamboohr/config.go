// Copyright 2026 SGNL.ai, Inc.

package bamboohr

import (
	"context"
	"errors"
	"fmt"

	"github.com/sgnl-ai/adapter-framework/web"
	"github.com/sgnl-ai/adapters/pkg/config"
)

var supportedAPIVersions = map[string]struct{}{
	"v1": {},
}

var supportedDateFormats = map[string]web.DateTimeFormatWithTimeZone{
	"yyyy-mm-dd": {Format: "2006-01-02", HasTimeZone: false},
	"mm/dd/yyyy": {Format: "01/02/2006", HasTimeZone: false},
	// "dd/mm/yyyy": UNSUPPORTED
	// "dd mon yyyy": UNSUPPORTED
}

// Config is the configuration passed in each GetPage calls to the adapter.
// BambooHR Adapter configuration example:
// nolint: godot
/*
{
    "apiVersion": "v1",
	"companyDomain": "sgnl",
	"onlyCurrent": true,
	"attributeMappings": {
		"date": "yyyy-mm-dd",
		"bool": {
			"true": ["True", "yes", "1"],
			"false": ["False", "no", "0"]
		}
	}
}
*/

type Config struct {
	// Common configuration
	*config.CommonConfig

	// APIVersion is the version of the BambooHR API to use.
	APIVersion string `json:"apiVersion"`

	// CompanyDomain is the domain of the BambooHR company.
	CompanyDomain string `json:"companyDomain"`

	// OnlyCurrent is a boolean value that determines if only current employees should be returned.
	// This optional field will default to false if not set.
	OnlyCurrent bool `json:"onlyCurrent"`

	// AttributeMappings is a map of attribute types to their BambooHR field values.
	// In BambooHR, the field values are strings, but the values can vary for fields of the same type
	// i.e. "bool" -> "true": ["yes"]
	// This is an optional field.
	AttributeMappings *AttributeMappings `json:"attributeMappings"`
}

type AttributeMappings struct {
	// Date is the date format configuration in the BambooHR Console.
	// This optional field will default to "YYYY-MM-DD" if not set.
	Date         *string                `json:"date"`
	BoolMappings *BoolAttributeMappings `json:"bool"`
}

type BoolAttributeMappings struct {
	True  []string `json:"true"`
	False []string `json:"false"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	switch {
	case c == nil:
		return errors.New("request contains no config")
	case c.APIVersion == "":
		return errors.New("apiVersion is not set")
	case c.CompanyDomain == "":
		return errors.New("companyDomain is not set")
	default:
		if _, found := supportedAPIVersions[c.APIVersion]; !found {
			return fmt.Errorf("apiVersion is not supported: %v", c.APIVersion)
		}

		if c.AttributeMappings == nil {
			c.AttributeMappings = &AttributeMappings{}
		}

		return c.AttributeMappings.Validate()
	}
}

func (m *AttributeMappings) Validate() error {
	if m.Date == nil {
		defaultDate := "yyyy-mm-dd"
		m.Date = &defaultDate
	}

	if _, found := supportedDateFormats[*m.Date]; !found {
		return fmt.Errorf("date format is not supported: %v", *m.Date)
	}

	return m.BoolMappings.Validate()
}

func (m *BoolAttributeMappings) Validate() error {
	if m == nil {
		return nil
	}

	errorList := make([]error, 0)

	switch {
	case m.True == nil || m.False == nil:
		return errors.New("Both attributeMappings.bool.true and attributeMappings.bool.false must be set")
	case len(m.True) == 0 || len(m.False) == 0:
		return errors.New("attributeMappings.bool must contain at least one mapping for true and false")
	default:
		setTrue := make(map[string]struct{})
		setFalse := make(map[string]struct{})

		for _, value := range m.True {
			if _, exists := setTrue[value]; exists {
				errorList = append(errorList, fmt.Errorf("attributeMappings.bool.true has a duplicate value: %s", value))
			}

			setTrue[value] = struct{}{}
		}

		for _, value := range m.False {
			if _, exists := setFalse[value]; exists {
				errorList = append(errorList, fmt.Errorf("attributeMappings.bool.false has a duplicate value: %s", value))
			}

			setFalse[value] = struct{}{}
		}

		for _, value := range m.False {
			if _, exists := setTrue[value]; exists {
				errorList = append(errorList, fmt.Errorf("Identical mapping found for both bool.false and bool.true: %s", value))
			}
		}

		if len(errorList) > 0 {
			return fmt.Errorf("These errors were found in your bool mapping list: %v", errorList)
		}

		return nil
	}
}
