// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			got, err := tt.request.BuildConnectionString()
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
