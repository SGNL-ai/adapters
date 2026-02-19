// Copyright 2025 SGNL.ai, Inc.

// Package validation contains address parsing and validation utilities for adapter
// configurations, ensuring that URLs are well-formed and use supported schemes. It trims
// whitespace, correctly handles URLs without schemes, and returns structured errors for invalid
// configurations.
package validation

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// ParseAndValidateAddress trims whitespace, parses the URL, and validates the scheme.
// Returns the trimmed address and parsed URL, or an error if invalid.
//
// The allowedSchemes parameter specifies which URL schemes are permitted.
// Scheme comparison is case-insensitive per RFC 3986 (url.Parse lowercases schemes).
// To allow URLs without a scheme, include an empty string "" in allowedSchemes.
//
// Addresses without "://" are treated as having no scheme (e.g., "example.com:8080"
// is parsed as host:port, not as scheme:opaque).
func ParseAndValidateAddress(address string, allowedSchemes []string) (string, *url.URL, *framework.Error) {
	trimmed := strings.TrimSpace(address)

	// Determine if the address has a scheme by checking for "://"
	// This prevents url.Parse from misinterpreting "host:port" as "scheme:opaque"
	hasScheme := strings.Contains(trimmed, "://")

	var (
		parsed *url.URL
		err    error
	)

	if hasScheme {
		parsed, err = url.Parse(trimmed)
	} else {
		// Prepend "//" so url.Parse treats it as a host (not scheme:opaque)
		parsed, err = url.Parse("//" + trimmed)
	}

	if err != nil {
		return "", nil, &framework.Error{
			Message: "Invalid URL format.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// Scheme is lowercased by url.Parse per RFC 3986.
	// Only validate scheme if the original address had one.
	if hasScheme && !slices.Contains(allowedSchemes, parsed.Scheme) {
		return "", nil, &framework.Error{
			Message: fmt.Sprintf("Scheme %q is not supported.", parsed.Scheme),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	return trimmed, parsed, nil
}
