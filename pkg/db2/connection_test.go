// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildConnectionString(t *testing.T) {
	tests := []struct {
		name    string
		request *Request
		want    string
	}{
		{
			name: "basic_connection_string_with_default_port",
			request: &Request{
				BaseURL:  "localhost",
				Database: "TESTDB",
				Username: "testuser",
				Password: "testpass",
			},
			want: "HOSTNAME=localhost;DATABASE=TESTDB;UID=testuser;PWD=testpass;PORT=50000;PROTOCOL=TCPIP",
		},
		{
			name: "connection_with_custom_port",
			request: &Request{
				BaseURL:  "db2server.example.com:60000",
				Database: "SGNL",
				Username: "admin",
				Password: "secret123",
			},
			want: "HOSTNAME=db2server.example.com;DATABASE=SGNL;UID=admin;PWD=secret123;PORT=60000;PROTOCOL=TCPIP",
		},
		{
			name: "connection_without_port_uses_default",
			request: &Request{
				BaseURL:  "db2server.example.com",
				Database: "SGNL",
				Username: "admin",
				Password: "secret123",
			},
			want: "HOSTNAME=db2server.example.com;DATABASE=SGNL;UID=admin;PWD=secret123;PORT=50000;PROTOCOL=TCPIP",
		},
		{
			name: "empty_baseurl_produces_empty_hostname",
			request: &Request{
				BaseURL:  "",
				Database: "TESTDB",
				Username: "testuser",
				Password: "testpass",
			},
			want: "HOSTNAME=;DATABASE=TESTDB;UID=testuser;PWD=testpass;PORT=50000;PROTOCOL=TCPIP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.request.BuildConnectionString(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseHostPort(t *testing.T) {
	tests := []struct {
		name         string
		baseURL      string
		expectedHost string
		expectedPort string
	}{
		{
			name:         "hostname_only_uses_default_port",
			baseURL:      "localhost",
			expectedHost: "localhost",
			expectedPort: DefaultDB2Port,
		},
		{
			name:         "hostname_with_port",
			baseURL:      "db2server.example.com:60000",
			expectedHost: "db2server.example.com",
			expectedPort: "60000",
		},
		{
			name:         "ipv4_address_only",
			baseURL:      "192.168.1.100",
			expectedHost: "192.168.1.100",
			expectedPort: DefaultDB2Port,
		},
		{
			name:         "ipv4_address_with_port",
			baseURL:      "192.168.1.100:50001",
			expectedHost: "192.168.1.100",
			expectedPort: "50001",
		},
		{
			name:         "colon_but_no_port_returns_empty_port",
			baseURL:      "localhost:",
			expectedHost: "localhost",
			expectedPort: "",
		},
		{
			name:         "empty_baseurl",
			baseURL:      "",
			expectedHost: "",
			expectedPort: DefaultDB2Port,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, port := parseHostPort(tt.baseURL)
			assert.Equal(t, tt.expectedHost, host)
			assert.Equal(t, tt.expectedPort, port)
		})
	}
}

// generateTestCertBase64 creates a self-signed certificate and returns it as a base64 string.
func generateTestCertBase64(t *testing.T) string {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	return base64.StdEncoding.EncodeToString(certPEM)
}

func TestBuildConnectionString_GivenConnectionProperties_WhenBuilding_ThenAppendsProperties(t *testing.T) {
	tests := []struct {
		name    string
		request *Request
		want    string
	}{
		{
			name: "single_property_appended",
			request: &Request{
				BaseURL:  "localhost",
				Database: "TESTDB",
				Username: "testuser",
				Password: "testpass",
				ConfigStruct: &Config{
					ConnectionProperties: map[string]string{
						"SecurityMechanism": "9",
					},
				},
			},
			want: "HOSTNAME=localhost;DATABASE=TESTDB;UID=testuser;PWD=testpass;PORT=50000;PROTOCOL=TCPIP;SecurityMechanism=9",
		},
		{
			name: "multiple_properties_appended_in_sorted_order",
			request: &Request{
				BaseURL:  "localhost",
				Database: "TESTDB",
				Username: "testuser",
				Password: "testpass",
				ConfigStruct: &Config{
					ConnectionProperties: map[string]string{
						"SecurityMechanism": "9",
						"ConnectTimeout":    "30",
						"Authentication":    "SERVER",
					},
				},
			},
			want: "HOSTNAME=localhost;DATABASE=TESTDB;UID=testuser;PWD=testpass;" +
				"PORT=50000;PROTOCOL=TCPIP;" +
				"Authentication=SERVER;ConnectTimeout=30;SecurityMechanism=9",
		},
		{
			name: "empty_map_no_extra_properties",
			request: &Request{
				BaseURL:  "localhost",
				Database: "TESTDB",
				Username: "testuser",
				Password: "testpass",
				ConfigStruct: &Config{
					ConnectionProperties: map[string]string{},
				},
			},
			want: "HOSTNAME=localhost;DATABASE=TESTDB;UID=testuser;PWD=testpass;PORT=50000;PROTOCOL=TCPIP",
		},
		{
			name: "nil_map_no_extra_properties",
			request: &Request{
				BaseURL:  "localhost",
				Database: "TESTDB",
				Username: "testuser",
				Password: "testpass",
				ConfigStruct: &Config{
					ConnectionProperties: nil,
				},
			},
			want: "HOSTNAME=localhost;DATABASE=TESTDB;UID=testuser;PWD=testpass;PORT=50000;PROTOCOL=TCPIP",
		},
		{
			name: "nil_config_struct_no_extra_properties",
			request: &Request{
				BaseURL:  "localhost",
				Database: "TESTDB",
				Username: "testuser",
				Password: "testpass",
			},
			want: "HOSTNAME=localhost;DATABASE=TESTDB;UID=testuser;PWD=testpass;PORT=50000;PROTOCOL=TCPIP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			got, err := tt.request.BuildConnectionString(context.Background())

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBuildConnectionString_GivenSSLAndConnectionProperties_WhenBuilding_ThenAppendsBoth(t *testing.T) {
	// Arrange
	certB64 := generateTestCertBase64(t)

	request := &Request{
		BaseURL:  "db2server:50000",
		Database: "SGNL",
		Username: "admin",
		Password: "secret",
		ConfigStruct: &Config{
			CertificateChain: certB64,
			ConnectionProperties: map[string]string{
				"ConnectTimeout": "30",
			},
		},
	}

	// Act
	got, err := request.BuildConnectionString(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Contains(t, got, "HOSTNAME=db2server;DATABASE=SGNL;UID=admin;PWD=secret;PORT=50000;PROTOCOL=TCPIP")
	assert.Contains(t, got, ";SECURITY=SSL;SSLServerCertificate=")
	assert.Contains(t, got, ";ConnectTimeout=30")
}

func TestBuildConnectionString_GivenDuplicateKeysWithDifferentCase_WhenBuilding_ThenSkipsDuplicate(t *testing.T) {
	// Arrange
	request := &Request{
		BaseURL:  "localhost",
		Database: "TESTDB",
		Username: "testuser",
		Password: "testpass",
		ConfigStruct: &Config{
			ConnectionProperties: map[string]string{
				"ConnectTimeout": "30",
				"CONNECTTIMEOUT": "60",
			},
		},
	}

	// Act
	got, err := request.BuildConnectionString(context.Background())

	// Assert — duplicate is silently skipped, only one connecttimeout appears
	assert.NoError(t, err)
	assert.Equal(t, 1, strings.Count(strings.ToLower(got), "connecttimeout"),
		"expected exactly one connecttimeout key in connection string")
}

func TestBuildConnectionString_GivenDisallowedProperty_WhenBuilding_ThenReturnsError(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{
			name: "trace_keyword_rejected",
			key:  "TraceFileName",
		},
		{
			name: "diagnostic_keyword_rejected",
			key:  "DiagPath",
		},
		{
			name: "arbitrary_unknown_keyword_rejected",
			key:  "MadeUpKeyword",
		},
		{
			name: "semicolon_in_key_rejected",
			key:  "Bad;Key",
		},
		{
			name: "empty_key_rejected",
			key:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			request := &Request{
				BaseURL:  "localhost",
				Database: "TESTDB",
				Username: "testuser",
				Password: "testpass",
				ConfigStruct: &Config{
					ConnectionProperties: map[string]string{
						tt.key: "value",
					},
				},
			}

			// Act
			_, err := request.BuildConnectionString(context.Background())

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(),
				"unsupported connection property")
		})
	}
}

func TestBuildConnectionString_GivenInvalidPropertyValue_WhenBuilding_ThenReturnsError(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "semicolon_rejected",
			value: "bad;value",
		},
		{
			name:  "empty_value_rejected",
			value: "",
		},
		{
			name:  "newline_rejected",
			value: "bad\nvalue",
		},
		{
			name:  "tab_rejected",
			value: "bad\tvalue",
		},
		{
			name:  "null_byte_rejected",
			value: "bad\x00value",
		},
		{
			name:  "carriage_return_rejected",
			value: "bad\rvalue",
		},
		{
			name:  "non_ascii_rejected",
			value: "bad\x80value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			request := &Request{
				BaseURL:  "localhost",
				Database: "TESTDB",
				Username: "testuser",
				Password: "testpass",
				ConfigStruct: &Config{
					ConnectionProperties: map[string]string{
						"ConnectTimeout": tt.value,
					},
				},
			}

			// Act
			_, err := request.BuildConnectionString(context.Background())

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(),
				"invalid connection property value")
		})
	}
}

func TestAllowedConnectionProperties_GivenExpectedKeys_WhenChecking_ThenContainsAll(t *testing.T) {
	// Verify the allow-list contains the keywords we expect to support.
	// This catches accidental deletions from the map.
	expectedKeys := []string{
		"authentication",
		"securitymechanism",
		"connecttimeout",
		"txnisolation",
		"tlsversion",
		"sslclientkeystash",
		"receivetimeout",
		"locktimeout",
		"readonlyconnection",
		"currentschema",
		"sslclienthostnamevalidation",
	}

	for _, key := range expectedKeys {
		assert.True(t, allowedConnectionProperties[key],
			"expected %q to be in allowedConnectionProperties", key)
	}
}

func TestSetupSSLCertificate_ReusesExistingFile(t *testing.T) {
	certB64 := generateTestCertBase64(t)

	// First call creates the file
	path1, err := setupSSLCertificate(certB64)
	require.NoError(t, err)
	require.NotEmpty(t, path1)
	defer os.Remove(path1)

	// Second call with same cert should return the same path
	path2, err := setupSSLCertificate(certB64)
	require.NoError(t, err)
	assert.Equal(t, path1, path2)
}

func TestSetupSSLCertificate_EmptyChainReturnsEmpty(t *testing.T) {
	path, err := setupSSLCertificate("")
	assert.NoError(t, err)
	assert.Empty(t, path)
}

func TestSetupSSLCertificate_InvalidBase64ReturnsError(t *testing.T) {
	_, err := setupSSLCertificate("not-valid-base64!!!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode base64")
}

func TestSetupSSLCertificate_InvalidPEMReturnsError(t *testing.T) {
	// Valid base64 but not valid PEM
	notPEM := base64.StdEncoding.EncodeToString([]byte("this is not PEM data"))
	_, err := setupSSLCertificate(notPEM)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no valid PEM data found")
}
