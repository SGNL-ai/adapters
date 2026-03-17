// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockDriverBehavior(t *testing.T) {
	// Test that we can create a SQL client without DB2 libraries
	client := NewDefaultSQLClient()
	assert.NotNil(t, client)

	// Test that connecting with mock driver gives appropriate error
	_, err := client.Connect("HOSTNAME=localhost;DATABASE=testdb;UID=test;PWD=test;PORT=50000;PROTOCOL=TCPIP")

	// Should get an error since we're using the mock driver
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mock DB2 driver: cannot connect to real DB2 database without client libraries")
	// Verify connection string (which may contain credentials) is NOT in the error
	assert.NotContains(t, err.Error(), "Connection string")
}
