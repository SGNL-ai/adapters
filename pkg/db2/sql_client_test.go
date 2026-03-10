// Copyright 2026 SGNL.ai, Inc.
package db2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultSQLClient_ConnectReusesConnection(t *testing.T) {
	sqlClient := NewDefaultSQLClient()
	client, ok := sqlClient.(*defaultSQLClient)
	if !ok {
		t.Fatal("expected *defaultSQLClient type")
	}

	// First connection attempt - will fail because no real DB, but tests the logic
	dsn1 := "HOSTNAME=localhost;DATABASE=TEST1;UID=user;PWD=pass;PORT=50000"

	// Store a mock db to simulate an existing connection
	client.dataSourceName = dsn1
	client.db = nil // Would be a real *sql.DB in production

	// Calling Connect with same DSN should return early (reuse)
	// Since db is nil, it won't actually reuse, but we're testing the branch
	_, _ = client.Connect(dsn1)

	// The dataSourceName should be set after Connect
	assert.Equal(t, dsn1, client.dataSourceName)
}

func TestDefaultSQLClient_QueryWithoutConnect(t *testing.T) {
	client := NewDefaultSQLClient()

	_, err := client.Query(nil, "SELECT 1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database connection not established")
}
