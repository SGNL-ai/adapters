// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"context"
	"errors"

	"github.com/sgnl-ai/adapters/pkg/condexpr"
	"github.com/sgnl-ai/adapters/pkg/config"
)

// DB2 connection constants.
const (
	// DB2DriverName is the driver name used for sql.Open() calls.
	DB2DriverName = "go_ibm_db"

	// DefaultDB2Port is the default port for DB2 connections.
	DefaultDB2Port = "50000"

	// ConnectionStringFormat is the format template for DB2 connection strings.
	// Parameters: hostname, database, username, password, port.
	ConnectionStringFormat = "HOSTNAME=%s;DATABASE=%s;UID=%s;PWD=%s;PORT=%s;PROTOCOL=TCPIP"

	// SSLConnectionSuffix is appended to connection strings when SSL is enabled.
	// Parameter: certificate file path.
	SSLConnectionSuffix = ";SECURITY=SSL;SSLServerCertificate=%s"

	// PrintableASCIIMin is the lowest allowed byte in property values (space).
	PrintableASCIIMin = 0x20

	// PrintableASCIIMax is the highest allowed byte in property values (~).
	PrintableASCIIMax = 0x7E
)

// Config is the configuration passed in each GetPage calls to the adapter.
//
// Adapter configuration example:
// nolint: godot
/*
{
	"requestTimeoutSeconds": 10,
	"localTimeZoneOffset": 43200,
	"database": "sgnl",
	"schema": "MYSCHEMA",
	"certificateChain": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURYVENDQWtXZ0F3SUJBZ0lKQUtvSy9Pdk16b...",
	"filters": {
		"users": {
			"or": [
				{
					"and": [
						{
							"field": "age",
							"op": ">",
							"value": 18
						},
						{
							"field": "country",
							"op": "=",
							"value": "USA"
						}
					]
				},
				{
					"field": "verified",
					"op": "=",
					"value": true
				}
			]
		},
		"groups": {
			"field": "country",
			"op": "IN",
			"value": ["active", "inactive"]
		}
	},
	"connectionProperties": {
		"SecurityMechanism": "9",
		"ConnectTimeout": "30",
		"Authentication": "SERVER"
	}
}
*/
type Config struct {
	// Common configuration
	*config.CommonConfig

	// API Version for configuration compatibility
	APIVersion string `json:"apiVersion,omitempty"`

	// DB2 database to connect to.
	Database string `json:"database,omitempty"`

	// Schema name to use for table queries (optional, defaults to username if not specified)
	Schema string `json:"schema,omitempty"`

	// CertificateChain is a base64-encoded PEM certificate chain to use for SSL connections.
	// When provided, SSL will be enabled and this certificate chain will be used
	// to validate the DB2 server's certificate.
	CertificateChain string `json:"certificateChain,omitempty"`

	Filters map[string]condexpr.Condition `json:"filters,omitempty"`

	// ConnectionProperties is a map of arbitrary key-value pairs appended to the
	// DB2 CLI driver connection string as ;KEY=VALUE. This enables passing driver
	// properties like SecurityMechanism, ConnectTimeout, etc. without needing a
	// dedicated Config field for each one.
	ConnectionProperties map[string]string `json:"connectionProperties,omitempty"`
}

// Validate validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	switch {
	case c == nil:
		return errors.New("request contains no config")
	case c.Database == "":
		return errors.New("database is not set")
	default:
		return nil
	}
}
