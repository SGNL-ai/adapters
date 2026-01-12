// Copyright 2025 SGNL.ai, Inc.
package validation

import (
	"context"
	"net"
	"testing"
)

func TestValidateExternalURL(t *testing.T) {
	defaultValidator := NewDefaultSSRFValidator()

	tests := map[string]struct {
		url     string
		wantErr bool
	}{
		"valid_https":                  {"https://login.microsoftonline.com/.well-known/openid-configuration", false},
		"valid_http":                   {"http://example.com/path", false},
		"valid_subdomain":              {"https://www.google.com", false},
		"valid_with_port":              {"https://example.com:8443/api", false},
		"invalid_aws_metadata_ipv4":    {"http://169.254.169.254/latest/meta-data/", true},
		"invalid_localhost":            {"http://localhost:8080", true},
		"invalid_loopback_127.0.0.1":   {"http://127.0.0.1:6379", true},
		"invalid_ipv6_loopback":        {"http://[::1]:8080", true},
		"invalid_private_ip_10.x":      {"http://10.0.0.1/admin", true},
		"invalid_private_ip_192.168":   {"http://192.168.1.1", true},
		"invalid_private_ip_172.16":    {"http://172.16.0.1", true},
		"invalid_file_scheme":          {"file:///etc/passwd", true},
		"invalid_gopher_scheme":        {"gopher://internal:70", true},
		"invalid_ftp_scheme":           {"ftp://example.com", true},
		"invalid_empty_url":            {"", true},
		"invalid_no_hostname":          {"http://", true},
		"invalid_domain_starting_dash": {"http://-example.com", true},
		"invalid_domain_ending_dash":   {"http://example-.com", true},
		"invalid_double_dots":          {"http://example..com", true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := defaultValidator.ValidateExternalURL(context.Background(), tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateExternalURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidDomainName(t *testing.T) {
	tests := map[string]struct {
		hostname string
		want     bool
	}{
		"valid_simple_domain":        {"example.com", true},
		"valid_subdomain":            {"api.example.com", true},
		"valid_multiple_subdomains":  {"api.v2.example.com", true},
		"valid_with_numbers":         {"example123.com", true},
		"valid_with_dash":            {"my-example.com", true},
		"valid_idn_prefix":           {"xn--domain.com", true},
		"invalid_starting_with_dash": {"-example.com", false},
		"invalid_ending_with_dash":   {"example-.com", false},
		"invalid_double_dots":        {"example..com", false},
		"invalid_starting_with_dot":  {".example.com", false},
		"invalid_ending_with_dot":    {"example.com.", false},
		"invalid_only_tld":           {"com", false},
		"invalid_empty":              {"", false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := isValidDomainName(tt.hostname)

			if got != tt.want {
				t.Errorf("isValidDomainName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsLocalhost(t *testing.T) {
	tests := map[string]struct {
		hostname string
		want     bool
	}{
		"valid_localhost":             {"localhost", true},
		"valid_localhost_uppercase":   {"Localhost", true},
		"valid_127.0.0.1":             {"127.0.0.1", true},
		"valid_127.1":                 {"127.1", true},
		"valid_::1":                   {"::1", true},
		"valid_[::1]":                 {"[::1]", true},
		"valid_0.0.0.0":               {"0.0.0.0", true},
		"valid_localhost.localdomain": {"localhost.localdomain", true},
		"invalid_example.com":         {"example.com", false},
		"invalid_10.0.0.1":            {"10.0.0.1", false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := isLocalhost(tt.hostname)

			if got != tt.want {
				t.Errorf("isLocalhost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPrivateIP(t *testing.T) {
	tests := map[string]struct {
		ip   string
		want bool
	}{
		"private_10.0.0.1":        {"10.0.0.1", true},
		"private_172.16.0.1":      {"172.16.0.1", true},
		"private_192.168.1.1":     {"192.168.1.1", true},
		"private_127.0.0.1":       {"127.0.0.1", true},
		"private_169.254.169.254": {"169.254.169.254", true},
		"public_8.8.8.8":          {"8.8.8.8", false},
		"public_1.1.1.1":          {"1.1.1.1", false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ip := parseIPHelper(t, tt.ip)
			got := isPrivateIP(ip)

			if got != tt.want {
				t.Errorf("isPrivateIP(%s) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func parseIPHelper(t *testing.T, ipStr string) net.IP {
	t.Helper()
	ip := net.ParseIP(ipStr)

	if ip == nil {
		t.Fatalf("Failed to parse IP: %s", ipStr)
	}

	return ip
}
