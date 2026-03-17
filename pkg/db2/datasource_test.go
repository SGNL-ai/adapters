// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/sgnl-ai/adapters/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestDatasource_NewClient(t *testing.T) {
	mockClient := NewMockSQLClient()
	client := NewClient(mockClient)
	assert.NotNil(t, client)
}

func TestTestConnection_GivenScanError_WhenIteratingRows_ThenLogsWarningAndContinues(t *testing.T) {
	// Arrange - create an observable logger to verify log output
	ctx, observedLogs := testutil.NewContextWithObservableLogger(context.Background())

	scanCount := 0
	scanErr := fmt.Errorf("scan failed: column type mismatch")

	mockClient := &MockSQLClient{
		ConnectFunc: func(_ string) (*sql.DB, error) { return nil, nil },
		QueryFunc: func(_ context.Context, _ string, _ ...interface{}) (Rows, error) {
			mockRows := &MockRows{
				Data: []map[string]interface{}{
					{"SCHEMANAME": "DB2INST1", "TABNAME": "BAD_ROW"},
					{"SCHEMANAME": "DB2INST1", "TABNAME": "GOOD_TABLE"},
				},
				columnsFunc: func() ([]string, error) {
					return []string{"SCHEMANAME", "TABNAME"}, nil
				},
			}
			mockRows.scanFunc = func(dest ...interface{}) error {
				scanCount++
				if scanCount == 1 {
					return scanErr
				}
				// Second row succeeds
				if len(dest) >= 2 {
					if ptr, ok := dest[0].(*string); ok {
						*ptr = "DB2INST1"
					}
					if ptr, ok := dest[1].(*string); ok {
						*ptr = "GOOD_TABLE"
					}
				}

				return nil
			}

			return mockRows, nil
		},
	}

	ds := &Datasource{Client: mockClient}
	request := &Request{
		BaseURL:  "localhost",
		Database: "TESTDB",
		Username: "user",
		Password: "pass",
	}

	// Act
	response, err := ds.TestConnection(ctx, request)

	// Assert - function should succeed and return only the good row
	require.Nil(t, err)
	require.NotNil(t, response)
	assert.Len(t, response.Objects, 1)
	assert.Equal(t, "GOOD_TABLE", response.Objects[0]["table"])

	// Assert - scan error should have been logged as a warning
	warnLogs := observedLogs.FilterLevelExact(zapcore.WarnLevel)
	require.Equal(t, 1, warnLogs.Len(), "expected exactly one warn log entry for scan error")
	logEntry := warnLogs.All()[0]
	assert.Contains(t, logEntry.Message, "scan")
}
