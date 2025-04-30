// Copyright 2025 SGNL.ai, Inc.
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	grpc_proxy_v1 "github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1"
)

type SQLRows []SQLRow
type SQLRow map[string]string
type SQLColumnTypes map[string]string

type SQLClient interface {
	IsProxied() bool
	Proxy(ctx context.Context, req *grpc_proxy_v1.ProxyRequestMessage) (*grpc_proxy_v1.Response, error)
	Connect(dataSourceName string) error
	Query(query string, args ...any) (*sql.Rows, error)
}

type defaultSQLClient struct {
	proxy grpc_proxy_v1.ProxyServiceClient
	DB    *sql.DB
}

func NewDefaultSQLClient(client grpc_proxy_v1.ProxyServiceClient) SQLClient {
	return &defaultSQLClient{
		proxy: client,
	}
}

// IsProxied returns true if the client is proxied.
func (c *defaultSQLClient) IsProxied() bool {
	return c.proxy != nil
}

// Connect opens a database connection to the provided datasource.
// The database is safe for concurrent use by multiple goroutines and
// maintains its own pool of idle connections. Thus, the Connect function
// should be called just once.
//
// This must be called before calling Query, else that function call will fail.
func (c *defaultSQLClient) Connect(dataSourceName string) error {
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
func (c *defaultSQLClient) Query(query string, args ...any) (*sql.Rows, error) {
	if c.DB == nil {
		return nil, errors.New("no open datasource connection")
	}

	// Prepare query statement. This is done to protect against the risk of SQL
	// injection attacks, which may be a risk to customer instances if a malicious
	// user gains inadvertent access to the SGNL console.
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

func (c *defaultSQLClient) Proxy(ctx context.Context, req *grpc_proxy_v1.ProxyRequestMessage,
) (*grpc_proxy_v1.Response, error) {
	return c.proxy.ProxyRequest(ctx, req)
}
