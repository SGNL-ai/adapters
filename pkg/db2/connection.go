// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"go.uber.org/zap"
)

// BuildConnectionString constructs the DB2 connection string from request parameters.
// It handles hostname/port parsing, credentials, and optional SSL configuration.
func (r *Request) BuildConnectionString(ctx context.Context) (string, error) {
	host, port := parseHostPort(r.BaseURL)

	connectionString := fmt.Sprintf(ConnectionStringFormat,
		host, r.Database, r.Username, r.Password, port)

	// Add SSL parameters and connection properties if config is provided
	if r.ConfigStruct != nil {
		if cfg, ok := r.ConfigStruct.(*Config); ok {
			if cfg.CertificateChain != "" {
				certPath, err := setupSSLCertificate(cfg.CertificateChain)
				if err != nil {
					return "", fmt.Errorf("SSL certificate setup failed: %w", err)
				}

				connectionString += fmt.Sprintf(SSLConnectionSuffix, certPath)
			}

			suffix, err := buildConnectionPropertiesSuffix(ctx, cfg.ConnectionProperties)
			if err != nil {
				return "", err
			}

			connectionString += suffix
		}
	}

	return connectionString, nil
}

// allowedConnectionProperties defines the set of DB2 CLI driver keywords
// that may be passed via Config.ConnectionProperties. Only keywords in this
// map are accepted; all others are rejected to prevent misuse of dangerous
// driver options (tracing, diagnostics, file paths, etc.).
// Keys are stored in lowercase because DB2 CLI keywords are case-insensitive.
// Reference: https://www.ibm.com/docs/en/db2/11.5.x?topic=odbc-cliodbc-configuration-keywords
var allowedConnectionProperties = map[string]bool{
	// Security & Authentication
	"authentication":    true,
	"securitymechanism": true,
	"krbplugin":         true,
	"pwdplugin":         true,
	"clientencalg":      true,
	"targetprincipal":   true,

	// Timeouts
	"connecttimeout":       true,
	"querytimeoutinterval": true,
	"receivetimeout":       true,
	"locktimeout":          true,

	// Transaction & Connection
	"autocommit":         true,
	"txnisolation":       true,
	"readonlyconnection": true,
	"currentschema":      true,
	"currentpackageset":  true,

	// TLS
	"tlsversion":                    true,
	"sslclientkeystash":             true,
	"sslclientkeystoredb":           true,
	"sslclientkeystoredbpassword":   true,
	"sslclientlabel":                true,
	"sslclienthostnamevalidation":   true,
}

// buildConnectionPropertiesSuffix builds a connection string suffix from
// allowed key-value properties. Keys are validated against the allow-list,
// sorted for deterministic output, and checked for injection characters.
func buildConnectionPropertiesSuffix(
	ctx context.Context,
	properties map[string]string,
) (string, error) {
	if len(properties) == 0 {
		return "", nil
	}

	logger := zaplogger.FromContext(ctx)
	keys := make([]string, 0, len(properties))
	seen := make(map[string]string, len(properties))

	for k := range properties {
		lower := strings.ToLower(k)

		if !allowedConnectionProperties[lower] {
			return "", fmt.Errorf(
				"unsupported connection property %q", k)
		}

		if prev, exists := seen[lower]; exists {
			logger.Warn("Skipping duplicate connection property",
				zap.String("skipped_key", k),
				zap.String("kept_key", prev),
			)

			continue
		}

		seen[lower] = k
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var builder strings.Builder

	for _, key := range keys {
		value := properties[key]

		if err := validatePropertyValue(key, value); err != nil {
			return "", err
		}

		builder.WriteString(";")
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(value)
	}

	return builder.String(), nil
}

// validatePropertyValue checks that a connection property value is non-empty
// printable ASCII and does not contain semicolons.
func validatePropertyValue(key, value string) error {
	if value == "" {
		return fmt.Errorf(
			"invalid connection property value for %q: "+
				"must not be empty", key)
	}

	for _, b := range []byte(value) {
		if b < PrintableASCIIMin || b > PrintableASCIIMax || b == ';' {
			return fmt.Errorf(
				"invalid connection property value for %q: "+
					"must be printable ASCII without semicolons",
				key)
		}
	}

	return nil
}

// parseHostPort extracts hostname and port from a URL string.
// Accepts "hostname" or "hostname:port" format.
// Returns default DB2 port if not specified.
func parseHostPort(baseURL string) (host, port string) {
	host = baseURL
	port = DefaultDB2Port

	if colonIndex := strings.LastIndex(baseURL, ":"); colonIndex != -1 {
		host = baseURL[:colonIndex]
		port = baseURL[colonIndex+1:]
	}

	return
}

// setupSSLCertificate creates a certificate file from the base64-encoded certificate chain
// and returns the path to use in the DB2 connection string.
// Uses a deterministic filename based on the certificate content hash to avoid
// creating duplicate temp files on repeated calls with the same certificate.
func setupSSLCertificate(certificateChain string) (string, error) {
	if certificateChain == "" {
		return "", nil
	}

	certPEM, err := base64.StdEncoding.DecodeString(certificateChain)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 certificate chain: %w", err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		return "", fmt.Errorf("failed to parse certificate chain: no valid PEM data found")
	}

	if _, err = x509.ParseCertificate(block.Bytes); err != nil {
		return "", fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Use a deterministic filename based on cert content hash to prevent
	// unbounded temp file growth from repeated calls with the same cert.
	hash := sha256.Sum256(certPEM)
	certFileName := fmt.Sprintf("db2-cert-%s.pem", hex.EncodeToString(hash[:8]))
	certPath := filepath.Join(os.TempDir(), certFileName)

	// If the file already exists with correct content, reuse it
	if _, err := os.Stat(certPath); err == nil {
		return certPath, nil
	}

	certFile, err := os.Create(certPath)
	if err != nil {
		return "", fmt.Errorf("failed to create certificate file: %w", err)
	}

	if _, err := certFile.Write(certPEM); err != nil {
		certFile.Close()
		os.Remove(certPath)

		return "", fmt.Errorf("failed to write certificate to file: %w", err)
	}

	if err := certFile.Close(); err != nil {
		os.Remove(certPath)

		return "", fmt.Errorf("failed to close certificate file: %w", err)
	}

	return certPath, nil
}
