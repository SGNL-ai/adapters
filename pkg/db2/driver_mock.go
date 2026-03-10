//go:build !db2

// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

// mockDriver is a mock implementation of the DB2 driver for development/testing.
// when DB2 client libraries are not available.
type mockDriver struct{}

func init() {
	sql.Register("go_ibm_db", &mockDriver{})
}

func (d *mockDriver) Open(name string) (driver.Conn, error) {
	return nil, fmt.Errorf(
		"mock DB2 driver: cannot connect to real DB2 database without client libraries. Connection string: %s",
		name,
	)
}
