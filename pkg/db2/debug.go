// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"context"
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"go.uber.org/zap"
)

// TestConnection tries to connect and query system tables to help debug connection issues.
func (d *Datasource) TestConnection(ctx context.Context, request *Request) (*Response, *framework.Error) {
	logger := zaplogger.FromContext(ctx)

	// Build DB2 connection string
	connectionString := fmt.Sprintf("HOSTNAME=%s;DATABASE=%s;UID=%s;PWD=%s;PORT=50000;PROTOCOL=TCPIP",
		request.BaseURL, request.Database, request.Username, request.Password)

	// Establish database connection
	_, err := d.Client.Connect(connectionString)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Error connecting to DB2 database: %v.", err),
		}
	}

	// Query to list available tables
	query := `SELECT SCHEMANAME, TABNAME FROM SYSCAT.TABLES ` +
		`WHERE TABNAME LIKE '%EKPO%' OR SCHEMANAME IN (?, 'DB2INST1') ` +
		`ORDER BY SCHEMANAME, TABNAME`

	rows, err := d.Client.Query(ctx, query, request.Username)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Error querying system tables: %v.", err),
		}
	}
	defer rows.Close()

	var objects []map[string]interface{}

	for rows.Next() {
		var schemaName, tableName string
		if err := rows.Scan(&schemaName, &tableName); err != nil {
			logger.Warn("Failed to scan system table row", zap.Error(err))

			continue
		}

		objects = append(objects, map[string]interface{}{
			"schema": schemaName,
			"table":  tableName,
		})
	}

	return &Response{
		Objects:    objects,
		NextCursor: nil,
		TotalCount: int64(len(objects)),
		StatusCode: 200,
	}, nil
}
