// Copyright 2025 SGNL.ai, Inc.
package validation

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
)

// SSRFValidator defines the interface for SSRF validation.
type SSRFValidator interface {
	// ValidateExternalURL checks if a URL is safe for external HTTP requests.
	ValidateExternalURL(ctx context.Context, rawURL string) error
}

// DefaultSSRFValidator implements SSRFValidator with full SSRF protection.
type defaultSSRFValidator struct{}

var _ SSRFValidator = (*defaultSSRFValidator)(nil)

// NewDefaultSSRFValidator creates a new DefaultSSRFValidator.
func NewDefaultSSRFValidator() *defaultSSRFValidator {
	return &defaultSSRFValidator{}
}

var (
	// Private IP ranges per RFC 1918 and other reserved ranges.
	privateIPRanges = []string{
		"10.0.0.0/8",         // RFC 1918 - Private network
		"172.16.0.0/12",      // RFC 1918 - Private network
		"192.168.0.0/16",     // RFC 1918 - Private network
		"127.0.0.0/8",        // RFC 1122 - Loopback
		"169.254.0.0/16",     // RFC 3927 - Link-local (AWS metadata)
		"::1/128",            // RFC 4291 - IPv6 loopback
		"fc00::/7",           // RFC 4193 - IPv6 unique local addresses
		"fe80::/10",          // RFC 4291 - IPv6 link-local
		"0.0.0.0/8",          // RFC 1122 - Current network
		"100.64.0.0/10",      // RFC 6598 - Shared address space
		"192.0.0.0/24",       // RFC 6890 - IETF Protocol Assignments
		"192.0.2.0/24",       // RFC 5737 - Documentation (TEST-NET-1)
		"198.18.0.0/15",      // RFC 2544 - Benchmarking
		"198.51.100.0/24",    // RFC 5737 - Documentation (TEST-NET-2)
		"203.0.113.0/24",     // RFC 5737 - Documentation (TEST-NET-3)
		"224.0.0.0/4",        // RFC 5771 - Multicast
		"240.0.0.0/4",        // RFC 1112 - Reserved
		"255.255.255.255/32", // RFC 919 - Limited broadcast
	}

	parsedPrivateRanges []*net.IPNet

	// Full domain validation: must have at least one dot and valid TLD.
	// Each label must start/end with alphanumeric, max 63 chars per label.
	// This is taken from OWASP.
	domainRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*\.[a-z]{2,}$`)
)

func init() {
	for _, cidr := range privateIPRanges {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err == nil {
			parsedPrivateRanges = append(parsedPrivateRanges, ipNet)
		}
	}
}

// ValidateExternalURL validates that a URL is safe for external HTTP requests.
// It performs the following checks:
// 1. URL format is valid
// 2. Scheme is http or https only
// 3. Hostname is not localhost or variations
// 4. Domain name follows DNS naming conventions
// 5. All resolved IPs are not in private/reserved ranges.
func (v *defaultSSRFValidator) ValidateExternalURL(ctx context.Context, rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format for %q: %w", rawURL, err)
	}

	// Only allow HTTP and HTTPS schemes.
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("unsupported URL scheme for %q: %s (only http/https allowed)", rawURL, scheme)
	}

	hostname := parsedURL.Hostname()
	if hostname == "" {
		return fmt.Errorf("URL must contain a hostname: %q", rawURL)
	}

	// Block localhost variations.
	if isLocalhost(hostname) {
		return fmt.Errorf("localhost URLs are not allowed: %q", rawURL)
	}

	// Check if hostname is already an IP address.
	if ip := net.ParseIP(hostname); ip != nil {
		if isPrivateIP(ip) {
			return fmt.Errorf("private IP addresses are not allowed for %q: %s", rawURL, ip)
		}

		// IP is public, allow it.
		return nil
	}

	// Validate domain name format.
	if !isValidDomainName(hostname) {
		return fmt.Errorf("invalid domain name format for %q: %s", rawURL, hostname)
	}

	resolver := &net.Resolver{}

	ips, err := resolver.LookupIP(ctx, "ip", hostname)
	if err != nil {
		return fmt.Errorf("failed to resolve hostname %s for %q: %w", hostname, rawURL, err)
	}

	if len(ips) == 0 {
		return fmt.Errorf("hostname %s did not resolve to any IP addresses for %q", hostname, rawURL)
	}

	// Check all resolved IPs are not internal IPs to prevent DNS rebinding attacks.
	for _, ip := range ips {
		if isPrivateIP(ip) {
			return fmt.Errorf("hostname %s resolves to private IP address %s which is not allowed for %q", hostname, ip, rawURL)
		}
	}

	return nil
}

// isLocalhost checks if a hostname is a localhost variation.
func isLocalhost(hostname string) bool {
	lower := strings.ToLower(hostname)

	return lower == "localhost" ||
		lower == "127.0.0.1" ||
		lower == "::1" ||
		strings.HasPrefix(lower, "localhost.") ||
		strings.HasPrefix(lower, "127.") ||
		lower == "0.0.0.0" ||
		lower == "[::1]"
}

// isPrivateIP checks if an IP address is in a private or reserved range.
func isPrivateIP(ip net.IP) bool {
	for _, ipNet := range parsedPrivateRanges {
		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}

// isValidDomainName validates that a hostname follows DNS naming conventions.
func isValidDomainName(hostname string) bool {
	lower := strings.ToLower(hostname)

	if len(lower) == 0 || len(lower) > 253 {
		return false
	}

	if strings.Contains(lower, "..") || strings.HasPrefix(lower, ".") || strings.HasSuffix(lower, ".") {
		return false
	}

	// Must contain at least one dot (to have a TLD).
	if !strings.Contains(lower, ".") {
		return false
	}

	// This ensures each label is 1-63 chars, starts/ends with alphanumeric.
	return domainRegex.MatchString(lower)
}
