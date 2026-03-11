// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BuildConnectionString constructs the DB2 connection string from request parameters.
// It handles hostname/port parsing, credentials, and optional SSL configuration.
func (r *Request) BuildConnectionString() (string, error) {
	host, port := parseHostPort(r.BaseURL)

	connectionString := fmt.Sprintf(ConnectionStringFormat,
		host, r.Database, r.Username, r.Password, port)

	// Add SSL parameters if certificate is provided
	if r.ConfigStruct != nil {
		if config, ok := r.ConfigStruct.(*Config); ok && config.CertificateChain != "" {
			certPath, err := setupSSLCertificate(config.CertificateChain)
			if err != nil {
				return "", fmt.Errorf("SSL certificate setup failed: %w", err)
			}

			connectionString += fmt.Sprintf(SSLConnectionSuffix, certPath)
		}
	}

	return connectionString, nil
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

	return host, port
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
