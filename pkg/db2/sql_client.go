// Copyright 2026 SGNL.ai, Inc.
package db2

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"
)

type SQLRows []SQLRow
type SQLRow map[string]string
type SQLColumnTypes map[string]string

type SQLClient interface {
	Connect(dataSourceName string) (*sql.DB, error)
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
}

// sqlRowsWrapper wraps sql.Rows to implement our Rows interface.
type sqlRowsWrapper struct {
	*sql.Rows
}

type defaultSQLClient struct {
	db             *sql.DB
	dataSourceName string
	mu             sync.Mutex
}

// NewDefaultSQLClient creates a new SQLClient instance.
// It is used to connect to a DB2 database and execute queries.
func NewDefaultSQLClient() SQLClient {
	return &defaultSQLClient{}
}

// Connect opens a database connection to the provided datasource.
// The database is safe for concurrent use by multiple goroutines and
// maintains its own pool of idle connections. Connections are reused
// if the datasource name matches; otherwise the old connection is closed.
//
// This must be called before calling Query, else that function call will fail.
func (c *defaultSQLClient) Connect(dataSourceName string) (*sql.DB, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Reuse existing connection if DSN matches
	if c.db != nil && c.dataSourceName == dataSourceName {
		return c.db, nil
	}

	// Close existing connection if DSN changed
	if c.db != nil {
		c.db.Close()
		c.db = nil
		c.dataSourceName = ""
	}

	db, err := sql.Open(DB2DriverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings for long-lived connection reuse
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		db.Close()

		return nil, err
	}

	// Store the connection and data source name for reuse
	c.db = db
	c.dataSourceName = dataSourceName

	return db, nil
}

// Query executes a query against the database and returns the result.
func (c *defaultSQLClient) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	// Ensure we have a database connection
	if c.db == nil {
		return nil, errors.New("database connection not established - call Connect() first")
	}

	// Execute the query
	sqlRows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	// Wrap the sql.Rows to implement our Rows interface
	return &sqlRowsWrapper{Rows: sqlRows}, nil
}
