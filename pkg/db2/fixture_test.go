// Copyright 2026 SGNL.ai, Inc.

// Test fixture loader for contract testing with recorded DB2 responses.
// Provides helpers to load JSON fixtures and create mock rows for testing.

package db2_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/sgnl-ai/adapters/pkg/condexpr"
	"github.com/sgnl-ai/adapters/pkg/db2"
)

// FixtureRequest captures the request parameters from a recorded fixture.
type FixtureRequest struct {
	Entity     string              `json:"entity"`
	Schema     string              `json:"schema"`
	Database   string              `json:"database"`
	PageSize   int64               `json:"pageSize"`
	Cursor     string              `json:"cursor,omitempty"`
	Filter     *condexpr.Condition `json:"filter,omitempty"`
	Attributes []string            `json:"attributes"`
}

// FixtureResponse captures the response from a recorded fixture.
type FixtureResponse struct {
	Objects    []map[string]interface{} `json:"objects"`
	NextCursor string                   `json:"nextCursor,omitempty"`
	Error      *FixtureError            `json:"error,omitempty"`
}

// FixtureError captures error details from a recorded fixture.
type FixtureError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Fixture represents a recorded test fixture with request and response.
type Fixture struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	RecordedAt  string          `json:"recordedAt"`
	Request     FixtureRequest  `json:"request"`
	Response    FixtureResponse `json:"response"`
}

// LoadFixture loads a single fixture from a JSON file.
func LoadFixture(name string) (*Fixture, error) {
	path := filepath.Join("testdata", "fixtures", name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read fixture %s: %w", name, err)
	}

	var fixture Fixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return nil, fmt.Errorf("failed to parse fixture %s: %w", name, err)
	}

	return &fixture, nil
}

// LoadAllFixtures loads all fixtures from the combined fixtures file.
func LoadAllFixtures() ([]Fixture, error) {
	path := filepath.Join("testdata", "fixtures", "all_fixtures.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read all fixtures: %w", err)
	}

	var fixtures []Fixture
	if err := json.Unmarshal(data, &fixtures); err != nil {
		return nil, fmt.Errorf("failed to parse fixtures: %w", err)
	}

	return fixtures, nil
}

// FixtureMockRows creates a MockRows implementation from fixture response objects.
type FixtureMockRows struct {
	objects []map[string]interface{}
	columns []string
	index   int
	closed  bool
}

// NewFixtureMockRows creates MockRows from fixture data.
func NewFixtureMockRows(objects []map[string]interface{}, columns []string) *FixtureMockRows {
	return &FixtureMockRows{
		objects: objects,
		columns: columns,
		index:   -1,
	}
}

// Next advances to the next row.
func (f *FixtureMockRows) Next() bool {
	if f.closed || f.index >= len(f.objects)-1 {
		return false
	}

	f.index++

	return true
}

// Scan copies column values into dest.
func (f *FixtureMockRows) Scan(dest ...interface{}) error {
	if f.index < 0 || f.index >= len(f.objects) {
		return fmt.Errorf("no current row")
	}

	row := f.objects[f.index]
	for i, col := range f.columns {
		if i >= len(dest) {
			break
		}

		val, ok := row[col]
		if !ok {
			// Column not in row - set to nil
			setNilValue(dest[i])

			continue
		}

		if err := scanValue(val, dest[i]); err != nil {
			return fmt.Errorf("scanning column %s: %w", col, err)
		}
	}

	return nil
}

// Close closes the rows iterator.
func (f *FixtureMockRows) Close() error {
	f.closed = true

	return nil
}

// Err returns any error encountered during iteration.
func (f *FixtureMockRows) Err() error {
	return nil
}

// Columns returns the column names.
func (f *FixtureMockRows) Columns() ([]string, error) {
	return f.columns, nil
}

// scanValue assigns a value from the fixture to a destination pointer.
func scanValue(val interface{}, dest interface{}) error {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}

	destElem := destVal.Elem()

	if val == nil {
		destElem.Set(reflect.Zero(destElem.Type()))

		return nil
	}

	switch d := dest.(type) {
	case *string:
		switch v := val.(type) {
		case string:
			*d = v
		case float64:
			*d = fmt.Sprintf("%v", v)
		default:
			*d = fmt.Sprintf("%v", v)
		}
	case *int64:
		switch v := val.(type) {
		case float64:
			*d = int64(v)
		case int:
			*d = int64(v)
		case int64:
			*d = v
		default:
			return fmt.Errorf("cannot convert %T to int64", val)
		}
	case *float64:
		switch v := val.(type) {
		case float64:
			*d = v
		case int:
			*d = float64(v)
		case int64:
			*d = float64(v)
		default:
			return fmt.Errorf("cannot convert %T to float64", val)
		}
	case *interface{}:
		*d = val
	default:
		// Try reflection-based assignment
		srcVal := reflect.ValueOf(val)
		if srcVal.Type().AssignableTo(destElem.Type()) {
			destElem.Set(srcVal)
		} else if srcVal.Type().ConvertibleTo(destElem.Type()) {
			destElem.Set(srcVal.Convert(destElem.Type()))
		} else {
			return fmt.Errorf("cannot assign %T to %T", val, dest)
		}
	}

	return nil
}

// setNilValue sets a pointer's value to its zero value.
func setNilValue(dest interface{}) {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() == reflect.Ptr {
		destVal.Elem().Set(reflect.Zero(destVal.Elem().Type()))
	}
}

// CreateMockClientFromFixture creates a MockSQLClient that returns fixture data.
func CreateMockClientFromFixture(fixture *Fixture) db2.SQLClient {
	// Build columns from fixture attributes plus total_remaining_rows for pagination
	columns := append([]string{}, fixture.Request.Attributes...)
	columns = append(columns, "total_remaining_rows")

	// Add total_remaining_rows to each object for pagination calculation
	objectsWithCount := make([]map[string]interface{}, len(fixture.Response.Objects))
	for i, obj := range fixture.Response.Objects {
		objectsWithCount[i] = make(map[string]interface{})
		for k, v := range obj {
			objectsWithCount[i][k] = v
		}
		// If there's a next cursor, there are more rows; otherwise 0
		if fixture.Response.NextCursor != "" {
			objectsWithCount[i]["total_remaining_rows"] = int64(1)
		} else {
			objectsWithCount[i]["total_remaining_rows"] = int64(0)
		}
	}

	return &db2.MockSQLClient{
		ConnectFunc: func(_ string) (*sql.DB, error) { return nil, nil },
		QueryFunc: func(_ context.Context, _ string, _ ...interface{}) (db2.Rows, error) {
			return NewFixtureMockRows(objectsWithCount, columns), nil
		},
	}
}
