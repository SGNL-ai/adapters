// Copyright 2025 SGNL.ai, Inc.

package validation

import (
	"testing"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

func TestParseAndValidateAddress(t *testing.T) {
	tests := map[string]struct {
		// Arrange
		address        string
		allowedSchemes []string

		// Assert expectations
		wantTrimmed   string
		wantScheme    string
		wantHost      string
		wantErr       bool
		wantErrCode   api_adapter_v1.ErrorCode
		wantErrSubstr string
	}{
		// Valid cases - HTTPS only (most common use case for adapters)
		"valid_https_simple": {
			address:        "https://example.com",
			allowedSchemes: []string{"https"},
			wantTrimmed:    "https://example.com",
			wantScheme:     "https",
			wantHost:       "example.com",
			wantErr:        false,
		},
		"valid_https_with_path": {
			address:        "https://api.example.com/v1/users",
			allowedSchemes: []string{"https"},
			wantTrimmed:    "https://api.example.com/v1/users",
			wantScheme:     "https",
			wantHost:       "api.example.com",
			wantErr:        false,
		},
		"valid_https_with_port": {
			address:        "https://example.com:8443",
			allowedSchemes: []string{"https"},
			wantTrimmed:    "https://example.com:8443",
			wantScheme:     "https",
			wantHost:       "example.com:8443",
			wantErr:        false,
		},

		// Valid cases - HTTP and HTTPS allowed
		"valid_http_when_allowed": {
			address:        "http://example.com",
			allowedSchemes: []string{"http", "https"},
			wantTrimmed:    "http://example.com",
			wantScheme:     "http",
			wantHost:       "example.com",
			wantErr:        false,
		},
		"valid_https_when_both_allowed": {
			address:        "https://example.com",
			allowedSchemes: []string{"http", "https"},
			wantTrimmed:    "https://example.com",
			wantScheme:     "https",
			wantHost:       "example.com",
			wantErr:        false,
		},

		// Valid cases - LDAP schemes
		"valid_ldaps": {
			address:        "ldaps://ldap.example.com:636",
			allowedSchemes: []string{"ldap", "ldaps"},
			wantTrimmed:    "ldaps://ldap.example.com:636",
			wantScheme:     "ldaps",
			wantHost:       "ldap.example.com:636",
			wantErr:        false,
		},
		"valid_ldap": {
			address:        "ldap://ldap.example.com:389",
			allowedSchemes: []string{"ldap", "ldaps"},
			wantTrimmed:    "ldap://ldap.example.com:389",
			wantScheme:     "ldap",
			wantHost:       "ldap.example.com:389",
			wantErr:        false,
		},

		// Whitespace trimming
		"trim_leading_whitespace": {
			address:        "   https://example.com",
			allowedSchemes: []string{"https"},
			wantTrimmed:    "https://example.com",
			wantScheme:     "https",
			wantHost:       "example.com",
			wantErr:        false,
		},
		"trim_trailing_whitespace": {
			address:        "https://example.com   ",
			allowedSchemes: []string{"https"},
			wantTrimmed:    "https://example.com",
			wantScheme:     "https",
			wantHost:       "example.com",
			wantErr:        false,
		},
		"trim_both_whitespace": {
			address:        "  https://example.com  ",
			allowedSchemes: []string{"https"},
			wantTrimmed:    "https://example.com",
			wantScheme:     "https",
			wantHost:       "example.com",
			wantErr:        false,
		},
		"trim_tab_and_newline": {
			address:        "\t\nhttps://example.com\n\t",
			allowedSchemes: []string{"https"},
			wantTrimmed:    "https://example.com",
			wantScheme:     "https",
			wantHost:       "example.com",
			wantErr:        false,
		},

		// Scheme case insensitivity - url.Parse lowercases schemes per RFC 3986
		"scheme_uppercase_http": {
			address:        "HTTP://example.com",
			allowedSchemes: []string{"http"},
			wantTrimmed:    "HTTP://example.com",
			wantScheme:     "http",
			wantHost:       "example.com",
			wantErr:        false,
		},
		"scheme_uppercase_https": {
			address:        "HTTPS://example.com",
			allowedSchemes: []string{"https"},
			wantTrimmed:    "HTTPS://example.com",
			wantScheme:     "https",
			wantHost:       "example.com",
			wantErr:        false,
		},
		"scheme_mixed_case": {
			address:        "HtTpS://example.com",
			allowedSchemes: []string{"https"},
			wantTrimmed:    "HtTpS://example.com",
			wantScheme:     "https",
			wantHost:       "example.com",
			wantErr:        false,
		},

		// No scheme provided - addresses without "://" are treated as schemeless
		// and parsed with "//" prefix so host is correctly extracted
		"no_scheme_just_hostname": {
			address:        "example.com",
			allowedSchemes: []string{"https"},
			wantTrimmed:    "example.com",
			wantScheme:     "",
			wantHost:       "example.com",
			wantErr:        false,
		},
		"no_scheme_with_port": {
			address:        "example.com:8080",
			allowedSchemes: []string{"https"},
			wantTrimmed:    "example.com:8080",
			wantScheme:     "",
			wantHost:       "example.com:8080",
			wantErr:        false,
		},
		"no_scheme_with_path": {
			address:        "example.com/api/v1",
			allowedSchemes: []string{"https"},
			wantTrimmed:    "example.com/api/v1",
			wantScheme:     "",
			wantHost:       "example.com",
			wantErr:        false,
		},

		// Invalid scheme cases
		"http_not_allowed": {
			address:        "http://example.com",
			allowedSchemes: []string{"https"},
			wantErr:        true,
			wantErrCode:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrSubstr:  `Scheme "http" is not supported`,
		},
		"ftp_not_allowed": {
			address:        "ftp://files.example.com",
			allowedSchemes: []string{"http", "https"},
			wantErr:        true,
			wantErrCode:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrSubstr:  `Scheme "ftp" is not supported`,
		},
		"file_scheme_not_allowed": {
			address:        "file:///etc/passwd",
			allowedSchemes: []string{"https"},
			wantErr:        true,
			wantErrCode:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrSubstr:  `Scheme "file" is not supported`,
		},
		"gopher_scheme_not_allowed": {
			address:        "gopher://evil.com",
			allowedSchemes: []string{"https"},
			wantErr:        true,
			wantErrCode:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrSubstr:  `Scheme "gopher" is not supported`,
		},
		// URLs with "://" but disallowed schemes are properly rejected
		"javascript_scheme_not_allowed": {
			address:        "javascript://alert(1)",
			allowedSchemes: []string{"https"},
			wantErr:        true,
			wantErrCode:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrSubstr:  `Scheme "javascript" is not supported`,
		},
		"data_scheme_not_allowed": {
			address:        "data://text/html",
			allowedSchemes: []string{"https"},
			wantErr:        true,
			wantErrCode:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrSubstr:  `Scheme "data" is not supported`,
		},

		// Empty allowed schemes - no scheme should be allowed except empty
		"scheme_present_but_none_allowed": {
			address:        "https://example.com",
			allowedSchemes: []string{},
			wantErr:        true,
			wantErrCode:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrSubstr:  `Scheme "https" is not supported`,
		},

		// Invalid URL format
		"invalid_url_format_control_chars": {
			address:        "https://example.com/path\x00with\x1fnull",
			allowedSchemes: []string{"https"},
			wantErr:        true,
			wantErrCode:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			wantErrSubstr:  "Invalid URL format",
		},

		// Edge cases - empty and whitespace only
		"empty_address": {
			address:        "",
			allowedSchemes: []string{"https"},
			wantTrimmed:    "",
			wantScheme:     "",
			wantHost:       "",
			wantErr:        false,
		},
		"whitespace_only_address": {
			address:        "   ",
			allowedSchemes: []string{"https"},
			wantTrimmed:    "",
			wantScheme:     "",
			wantHost:       "",
			wantErr:        false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			gotTrimmed, gotParsed, gotErr := ParseAndValidateAddress(tt.address, tt.allowedSchemes)

			// Assert - error cases
			if tt.wantErr {
				if gotErr == nil {
					t.Errorf("ParseAndValidateAddress() expected error, got nil")

					return
				}

				if gotErr.Code != tt.wantErrCode {
					t.Errorf("ParseAndValidateAddress() error code = %v, want %v", gotErr.Code, tt.wantErrCode)
				}

				if tt.wantErrSubstr != "" && !containsSubstring(gotErr.Message, tt.wantErrSubstr) {
					t.Errorf("ParseAndValidateAddress() error message = %q, want substring %q", gotErr.Message, tt.wantErrSubstr)
				}

				return
			}

			// Assert - success cases
			if gotErr != nil {
				t.Errorf("ParseAndValidateAddress() unexpected error = %v", gotErr)

				return
			}

			if gotTrimmed != tt.wantTrimmed {
				t.Errorf("ParseAndValidateAddress() trimmed = %q, want %q", gotTrimmed, tt.wantTrimmed)
			}

			if gotParsed == nil {
				t.Errorf("ParseAndValidateAddress() parsed URL is nil, expected non-nil")

				return
			}

			if gotParsed.Scheme != tt.wantScheme {
				t.Errorf("ParseAndValidateAddress() scheme = %q, want %q", gotParsed.Scheme, tt.wantScheme)
			}

			if gotParsed.Host != tt.wantHost {
				t.Errorf("ParseAndValidateAddress() host = %q, want %q", gotParsed.Host, tt.wantHost)
			}
		})
	}
}

func TestParseAndValidateAddress_CommonAdapterSchemes(t *testing.T) {
	// Test the common use case: HTTPS only (rejecting HTTP)
	// This is the pattern used by ServiceNow, Okta, Workday adapters

	httpsOnlySchemes := []string{"https"}

	tests := map[string]struct {
		// Arrange
		address string

		// Assert expectations
		wantErr       bool
		wantErrSubstr string
	}{
		"https_allowed": {
			address: "https://myinstance.service-now.com",
			wantErr: false,
		},
		"http_rejected": {
			address:       "http://myinstance.service-now.com",
			wantErr:       true,
			wantErrSubstr: `Scheme "http" is not supported`,
		},
		"http_rejected_with_whitespace": {
			address:       "  http://myinstance.service-now.com  ",
			wantErr:       true,
			wantErrSubstr: `Scheme "http" is not supported`,
		},
		"https_with_leading_whitespace": {
			address: "   https://myinstance.service-now.com",
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			_, _, gotErr := ParseAndValidateAddress(tt.address, httpsOnlySchemes)

			// Assert
			if tt.wantErr {
				if gotErr == nil {
					t.Errorf("ParseAndValidateAddress() expected error for address %q", tt.address)

					return
				}

				if tt.wantErrSubstr != "" && !containsSubstring(gotErr.Message, tt.wantErrSubstr) {
					t.Errorf("ParseAndValidateAddress() error = %q, want substring %q", gotErr.Message, tt.wantErrSubstr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("ParseAndValidateAddress() unexpected error = %v", gotErr)
				}
			}
		})
	}
}

// containsSubstring checks if s contains substr (case-sensitive).
func containsSubstring(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
