// Copyright 2026 SGNL.ai, Inc.
package db2

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultSQLClient_ConnectReusesConnection(t *testing.T) {
	sqlClient := NewDefaultSQLClient()
	client, ok := sqlClient.(*defaultSQLClient)
	if !ok {
		t.Fatal("expected *defaultSQLClient type")
	}

	dsn := "HOSTNAME=localhost;DATABASE=TEST1;UID=user;PWD=pass;PORT=50000"

	// Create a *sql.DB handle via sql.Open (succeeds even with mock driver)
	existingDB, err := sql.Open(DB2DriverName, dsn)
	require.NoError(t, err)
	defer existingDB.Close()

	// Simulate an existing connection by setting db and DSN
	client.db = existingDB
	client.dataSourceName = dsn

	// Act - calling Connect with the same DSN should return the existing db
	returnedDB, err := client.Connect(dsn)

	// Assert - should reuse the existing connection
	require.NoError(t, err)
	assert.Same(t, existingDB, returnedDB)
}

func TestDefaultSQLClient_QueryWithoutConnect(t *testing.T) {
	client := NewDefaultSQLClient()

	_, err := client.Query(nil, "SELECT 1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database connection not established")
}
