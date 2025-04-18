// Copyright 2025 SGNL.ai, Inc.
package mysql

import (
	"database/sql"
	"errors"
	"time"
)

type SQLRows []SQLRow
type SQLRow map[string]string
type SQLColumnTypes map[string]string

type SQLClient interface {
	Connect(dataSourceName string) error
	Query(query string, args ...any) (*sql.Rows, error)
}

type DefaultSQLClient struct {
	DB *sql.DB
}

func NewDefaultSQLClient() *DefaultSQLClient {
	return &DefaultSQLClient{}
}

// Connect opens a database connection to the provided datasource. The database is safe for concurrent use by multiple
// goroutines and maintains its own pool of idle connections. Thus, the Connect function should be called just once.
//
// This must be called before calling Query, else that function call will fail.
func (c *DefaultSQLClient) Connect(dataSourceName string) error {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(1 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	c.DB = db

	return nil
}

// Query prepares a statement and queries a connected database with the provided query.
//
// Returns an error if the query fails or if there is no currently open database connection.
func (c *DefaultSQLClient) Query(query string, args ...any) (*sql.Rows, error) {
	if c.DB == nil {
		return nil, errors.New("no open datasource connection")
	}

	// Prepare query statement. This is done to protect against the risk of SQL injection attacks,
	// which may be a risk to customer instances if a malicious user gains inadvertent access to the
	// SGNL console.
	stmt, err := c.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query with provided arguments.
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
