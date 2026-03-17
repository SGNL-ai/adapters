// Copyright 2026 SGNL.ai, Inc.
package db2

import (
	"context"
	"database/sql"
)

// MockRows is a mock implementation of the Rows interface for testing.
type MockRows struct {
	Data        []map[string]interface{}
	index       int
	closed      bool
	scanFunc    func(dest ...interface{}) error
	errFunc     func() error
	columnsFunc func() ([]string, error)
}

// Next advances to the next row.
func (m *MockRows) Next() bool {
	if m.closed || m.index >= len(m.Data) {
		return false
	}

	m.index++

	return m.index <= len(m.Data)
}

// Scan copies column values into dest.
func (m *MockRows) Scan(dest ...interface{}) error {
	if m.scanFunc != nil {
		return m.scanFunc(dest...)
	}
	// Default implementation - just return nil for testing
	return nil
}

// Close closes the rows iterator.
func (m *MockRows) Close() error {
	m.closed = true

	return nil
}

// Err returns any error encountered during iteration.
func (m *MockRows) Err() error {
	if m.errFunc != nil {
		return m.errFunc()
	}

	return nil
}

// Columns returns the column names.
func (m *MockRows) Columns() ([]string, error) {
	if m.columnsFunc != nil {
		return m.columnsFunc()
	}
	// Default mock column names for testing
	return []string{"MANDT", "EBELN", "EBELP", "total_remaining_rows"}, nil
}

// MockSQLClient is a mock implementation of SQLClient for testing.
type MockSQLClient struct {
	ConnectFunc func(dataSourceName string) (*sql.DB, error)
	QueryFunc   func(ctx context.Context, query string, args ...interface{}) (Rows, error)
}

// NewMockSQLClient creates a new MockSQLClient instance.
func NewMockSQLClient() SQLClient {
	return &MockSQLClient{
		ConnectFunc: func(_ string) (*sql.DB, error) {
			return nil, nil
		},
		QueryFunc: func(_ context.Context, _ string, _ ...interface{}) (Rows, error) {
			return &MockRows{}, nil
		},
	}
}

// Connect opens a database connection.
func (m *MockSQLClient) Connect(dataSourceName string) (*sql.DB, error) {
	if m.ConnectFunc != nil {
		return m.ConnectFunc(dataSourceName)
	}

	return nil, nil
}

// Query executes a database query.
func (m *MockSQLClient) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, query, args...)
	}

	return &MockRows{}, nil
}
