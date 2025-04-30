// Copyright 2025 SGNL.ai, Inc.
package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"math"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
	grpc_proxy_v1 "github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1"
)

type MockSQLClient struct {
	DB   *sql.DB
	Mock sqlmock.Sqlmock
}

func NewMockSQLClient() *MockSQLClient {
	return &MockSQLClient{}
}

const (
	TestDatasourceForConnectFailure = "test.connect.failure"
)

func (c *MockSQLClient) Connect(datasourceName string) error {
	if strings.Contains(datasourceName, TestDatasourceForConnectFailure) {
		return errors.New("failed to connect to mock sql service")
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		return err
	}

	c.Mock = mock
	c.DB = db

	return nil
}

// nolint: lll
func (c *MockSQLClient) Query(query string, args ...any) (*sql.Rows, error) {
	if len(args) != 2 {
		return nil, errors.New("mock sql client called with unsupported number of args")
	}

	pageSize, ok := args[0].(int64)
	if !ok {
		return nil, errors.New("mock sql client called with invalid arg[1], unable to cast `pageSize` to int64")
	}

	cursor, ok := args[1].(int64)
	if !ok {
		return nil, errors.New("mock sql client called with invalid arg[2], unable to cast `cursor` to int64")
	}

	if query != "SELECT *, CAST(id as CHAR(50)) as str_id FROM users ORDER BY str_id ASC LIMIT ? OFFSET ?" {
		return nil, errors.New("mock sql client called with unsupported query")
	}

	switch {
	// First page of users.
	case pageSize == 5 && cursor == 0:
		columns := []*sqlmock.Column{
			sqlmock.NewColumn("id").OfType("VARCHAR", ""),
			sqlmock.NewColumn("name").OfType("VARCHAR", ""),
			sqlmock.NewColumn("active").OfType("BOOL", ""),
			sqlmock.NewColumn("employee_number").OfType("INT", ""),
			sqlmock.NewColumn("risk_score").OfType("FLOAT", ""),
			sqlmock.NewColumn("last_modified").OfType("DATETIME", ""),
		}

		mockRows := sqlmock.NewRowsWithColumnDefinition(columns...).
			AddRow("a20bab52-52e3-46c2-bd6a-2ad1512f713f", "Ernesto Gregg", true, 1, 1.0, "2025-02-12T22:38:00+00:00").
			AddRow("d35c298e-d343-4ad8-ac35-f7c5d9d47cb9", "Eleanor Watts", true, 2, 1.562, "2025-02-12T22:38:00+00:00").
			AddRow("62c74831-be4a-4cad-88fa-4e02640269d2", "Chris Griffin", true, 3, 4.23, "2025-02-12T22:38:00+00:00").
			AddRow("65b8fa65-25c5-4682-997f-ca86923e59e4", "Casey Manning", false, 4, 10, "2025-02-12T22:38:00+00:00").
			AddRow("6598acf9-cccc-48c9-ab9b-754bbe9ad146", "Helen Gray", true, 5, 3.25, "2025-02-12T22:38:00+00:00")

		values := []driver.Value{}

		for _, arg := range args {
			values = append(values, driver.Value(arg))
		}

		c.Mock.ExpectQuery(`SELECT \*, CAST\(id as CHAR\(50\)\) as str_id FROM users ORDER BY str_id ASC LIMIT \? OFFSET \?`).
			WithArgs(values...).
			WillReturnRows(mockRows)

	// Second (middle) page of users. Tests providing BOOLs as TINYINT.
	case pageSize == 5 && cursor == 5:
		columns := []*sqlmock.Column{
			sqlmock.NewColumn("id").OfType("VARCHAR", ""),
			sqlmock.NewColumn("name").OfType("VARCHAR", ""),
			sqlmock.NewColumn("active").OfType("TINYINT", ""),
			sqlmock.NewColumn("employee_number").OfType("INT", ""),
			sqlmock.NewColumn("risk_score").OfType("FLOAT", ""),
			sqlmock.NewColumn("last_modified").OfType("DATETIME", ""),
		}

		mockRows := sqlmock.NewRowsWithColumnDefinition(columns...).
			AddRow("7390f7fc-0145-4691-9f55-b5c783369db9", "Martha Pollard", 1, 6, 1.0, "2025-02-12T22:38:00+00:00").
			AddRow("745cf6d6-55c8-4863-9bf6-1b1a80ff1515", "Roxanne Dixon", 1, 7, 1.0, "2025-02-12T22:38:00+00:00").
			AddRow("776b45f0-a2e3-4424-8ef7-84f3052bebc7", "Verna Ferrell", 1, 8, 1.0, "2025-02-12T22:38:00+00:00").
			AddRow("8b9643f9-25b4-458a-ad4f-81e61d106a57", "Adrian Carey", 1, 9, 1.0, "2025-02-12T22:38:00+00:00").
			AddRow("88ff7d742-fb3c-4103-af4b-fcd4315bae66", "Joshua Martinez", 1, 10, 1.0, "2025-02-12T22:38:00+00:00")

		values := []driver.Value{}

		for _, arg := range args {
			values = append(values, driver.Value(arg))
		}

		c.Mock.ExpectQuery(`SELECT \*, CAST\(id as CHAR\(50\)\) as str_id FROM users ORDER BY str_id ASC LIMIT \? OFFSET \?`).
			WithArgs(values...).
			WillReturnRows(mockRows)

	// Third (last) page of users.
	case pageSize == 5 && cursor == 10:
		columns := []*sqlmock.Column{
			sqlmock.NewColumn("id").OfType("VARCHAR", ""),
			sqlmock.NewColumn("name").OfType("VARCHAR", ""),
			sqlmock.NewColumn("active").OfType("BOOL", ""),
			sqlmock.NewColumn("employee_number").OfType("INT", ""),
			sqlmock.NewColumn("risk_score").OfType("FLOAT", ""),
			sqlmock.NewColumn("last_modified").OfType("DATETIME", ""),
		}

		mockRows := sqlmock.NewRowsWithColumnDefinition(columns...).
			AddRow("9cf5a596-0df2-4510-a403-9b514fd500b8", "Erica Meadows", true, 6, 1.0, "2025-02-12T22:38:00+00:00").
			AddRow("987053f0-c06c-48ee-9c99-81f3a96af639", "Carole Crawford", true, 7, 1.0, "2025-02-12T22:38:00+00:00")

		values := []driver.Value{}

		for _, arg := range args {
			values = append(values, driver.Value(arg))
		}

		c.Mock.ExpectQuery(`SELECT \*, CAST\(id as CHAR\(50\)\) as str_id FROM users ORDER BY str_id ASC LIMIT \? OFFSET \?`).
			WithArgs(values...).
			WillReturnRows(mockRows)

	// Test: Failed to query datasource
	case pageSize == 1 && cursor == 101:
		return nil, errors.New("failed to query mock sql service")

	// Test: Edge case with large values
	case pageSize == 5 && cursor == 202:
		columns := []*sqlmock.Column{
			sqlmock.NewColumn("id").OfType("VARCHAR", ""),
			sqlmock.NewColumn("name").OfType("VARCHAR", ""),
			sqlmock.NewColumn("active").OfType("BOOL", ""),
			sqlmock.NewColumn("employee_number").OfType("BIGINT", ""),
			sqlmock.NewColumn("risk_score").OfType("FLOAT", ""),
			sqlmock.NewColumn("last_modified").OfType("DATETIME", ""),
		}

		mockRows := sqlmock.NewRowsWithColumnDefinition(columns...).
			// TODO [sc-42217]: Allow providing values directly as Ints.
			// Any Int ingested with more than 16 digits will lose precision (according to IEEE 754) since we cast
			// all Ints to Float64 before storing them. This is why we're using `1<<53-1` instead of `1<<63-1` since we
			// will lose precision using the MaxInt64 const.
			AddRow("9cf5a596-0df2-4510-a403-9b514fd500b8", "Erica Meadows", true, 1<<53-1, math.MaxFloat64, "2025-02-12T22:38:00+00:00").
			AddRow("dfaf01cc-85b7-4e2e-b2d7-608d1f1904fe", "Eleanor Watts", true, math.MinInt64, -math.MaxFloat64, "2025-02-12T22:38:00+00:00")

		values := []driver.Value{}

		for _, arg := range args {
			values = append(values, driver.Value(arg))
		}

		c.Mock.ExpectQuery(`SELECT \*, CAST\(id as CHAR\(50\)\) as str_id FROM users ORDER BY str_id ASC LIMIT \? OFFSET \?`).
			WithArgs(values...).
			WillReturnRows(mockRows)

	// Test: Edge case with empty values
	case pageSize == 5 && cursor == 203:
		columns := []*sqlmock.Column{
			sqlmock.NewColumn("id").OfType("VARCHAR", ""),
			sqlmock.NewColumn("name").OfType("VARCHAR", ""),
			sqlmock.NewColumn("active").OfType("BOOL", ""),
			sqlmock.NewColumn("employee_number").OfType("BIGINT", ""),
			sqlmock.NewColumn("risk_score").OfType("FLOAT", ""),
			sqlmock.NewColumn("last_modified").OfType("DATETIME", ""),
		}

		mockRows := sqlmock.NewRowsWithColumnDefinition(columns...).
			AddRow("9cf5a596-0df2-4510-a403-9b514fd500b8", "", "", "", "", "").
			AddRow("a20bab52-52e3-46c2-bd6a-2ad1512f713f", "Ernesto Gregg", true, 1, 1.0, "2025-02-12T22:38:00+00:00")

		values := []driver.Value{}

		for _, arg := range args {
			values = append(values, driver.Value(arg))
		}

		c.Mock.ExpectQuery(`SELECT \*, CAST\(id as CHAR\(50\)\) as str_id FROM users ORDER BY str_id ASC LIMIT \? OFFSET \?`).
			WithArgs(values...).
			WillReturnRows(mockRows)

	default:
		return nil, errors.New("mock sql client called with unsupported args")
	}

	// This is not a prepared query so this may be vulnerable to SQL injection attacks, but this is only used
	// for tests so this is not a concern.
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	if err := c.Mock.ExpectationsWereMet(); err != nil {
		return nil, err
	}

	return rows, nil
}

func (c *MockSQLClient) Proxy(_ context.Context, _ *grpc_proxy_v1.ProxyRequestMessage,
) (*grpc_proxy_v1.Response, error) {
	return nil, nil
}
