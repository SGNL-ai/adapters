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

	// Test that connecting without a running DB2 server gives an error.
	// With mock driver: "mock DB2 driver: cannot connect..."
	// With real driver: IBM CLI driver TCP connection error
	_, err := client.Connect("HOSTNAME=localhost;DATABASE=testdb;UID=test;PWD=test;PORT=50000;PROTOCOL=TCPIP")

	assert.Error(t, err)
}
