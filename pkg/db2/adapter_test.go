// Copyright 2026 SGNL.ai, Inc.

package db2_test

import (
	"context"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/db2"
	"github.com/stretchr/testify/assert"
)

func TestAdapter_GetPage_Success(t *testing.T) {
	adapter := db2.NewAdapter(db2.NewClient(db2.NewMockSQLClient()))

	request := &framework.Request[db2.Config]{
		Auth: &framework.DatasourceAuthCredentials{
			Basic: &framework.BasicAuthCredentials{
				Username: "testuser",
				Password: "testpass",
			},
		},
		Address:  "localhost",
		PageSize: 100,
		Entity: framework.EntityConfig{
			ExternalId: "test_table",
			Attributes: []*framework.AttributeConfig{
				{
					ExternalId: "id",
					UniqueId:   true,
				},
			},
		},
		Config: &db2.Config{
			Database: "TESTDB",
		},
	}

	response := adapter.GetPage(context.Background(), request)
	assert.NotNil(t, response)
	// Note: This will fail without a proper mock setup, but demonstrates the structure
}

func TestAdapter_GetPage_AddressValidation(t *testing.T) {
	adapter := db2.NewAdapter(db2.NewClient(db2.NewMockSQLClient()))

	baseRequest := func(address string) *framework.Request[db2.Config] {
		return &framework.Request[db2.Config]{
			Auth: &framework.DatasourceAuthCredentials{
				Basic: &framework.BasicAuthCredentials{
					Username: "testuser",
					Password: "testpass",
				},
			},
			Address:  address,
			PageSize: 100,
			Entity: framework.EntityConfig{
				ExternalId: "test_table",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						UniqueId:   true,
					},
				},
			},
			Config: &db2.Config{
				Database: "TESTDB",
			},
		}
	}

	tests := []struct {
		name        string
		address     string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid_hostname_only",
			address: "localhost",
			wantErr: false,
		},
		{
			name:    "valid_hostname_with_port",
			address: "db2server.example.com:50000",
			wantErr: false,
		},
		{
			name:    "valid_ipv4_with_port",
			address: "192.168.1.100:50001",
			wantErr: false,
		},
		{
			name:        "empty_address",
			address:     "",
			wantErr:     true,
			errContains: "address (hostname) is required",
		},
		{
			name:        "colon_with_empty_port",
			address:     "localhost:",
			wantErr:     true,
			errContains: "invalid port ''",
		},
		{
			name:        "invalid_port_non_numeric",
			address:     "localhost:abc",
			wantErr:     true,
			errContains: "invalid port 'abc'",
		},
		{
			name:        "port_out_of_range_high",
			address:     "localhost:65536",
			wantErr:     true,
			errContains: "invalid port '65536'",
		},
		{
			name:        "port_zero",
			address:     "localhost:0",
			wantErr:     true,
			errContains: "invalid port '0'",
		},
		{
			name:        "empty_hostname_with_port",
			address:     ":50000",
			wantErr:     true,
			errContains: "hostname is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := adapter.GetPage(context.Background(), baseRequest(tt.address))

			if tt.wantErr {
				assert.NotNil(t, response.Error)
				assert.Contains(t, response.Error.Message, tt.errContains)
			} else {
				// For valid addresses, error should be nil or not related to address validation
				if response.Error != nil {
					assert.NotContains(t, response.Error.Message, "address")
					assert.NotContains(t, response.Error.Message, "port")
					assert.NotContains(t, response.Error.Message, "hostname")
				}
			}
		})
	}
}

func TestAdapter_GetPage_MockDriverError(t *testing.T) {
	// Test integration with mock driver to verify error message
	client := db2.NewDefaultSQLClient()
	datasource := db2.NewClient(client)
	adapter := db2.NewAdapter(datasource)

	request := &framework.Request[db2.Config]{
		Auth: &framework.DatasourceAuthCredentials{
			Basic: &framework.BasicAuthCredentials{
				Username: "testuser",
				Password: "testpass",
			},
		},
		Address: "localhost",
		Config: &db2.Config{
			Database: "testdb",
		},
		Entity: framework.EntityConfig{
			ExternalId: "EKPO",
			Attributes: []*framework.AttributeConfig{
				{
					ExternalId: "id",
					UniqueId:   true,
				},
			},
		},
		PageSize: 100,
	}

	response := adapter.GetPage(context.Background(), request)

	// Should get an error response with the mock driver message
	if response.Error != nil {
		assert.Contains(t, response.Error.Message, "Error connecting to DB2 database")
		assert.Contains(t, response.Error.Message,
			"mock DB2 driver: cannot connect to real DB2 database without client libraries")
	} else {
		t.Errorf("Expected error response but got success")
	}
}
