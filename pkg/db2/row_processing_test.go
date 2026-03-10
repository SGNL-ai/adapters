// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"context"
	"database/sql"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// createTestLogger creates a no-op logger for testing.
func createTestLogger() *zap.Logger {
	return zap.NewNop()
}

func TestMockRows_ScanWorks(t *testing.T) {
	// Verify our mock rows actually work before using them in GetPage tests
	data := []map[string]interface{}{
		{"col1": "value1", "col2": int64(42)},
	}

	mockRows := createMockRowsWithData(data)

	// Test Columns
	cols, err := mockRows.Columns()
	require.NoError(t, err)
	t.Logf("Columns: %v", cols)

	// Test Next
	hasNext := mockRows.Next()
	assert.True(t, hasNext, "should have first row")
	t.Logf("After Next: index=%d", mockRows.index)

	// Test Scan
	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	err = mockRows.Scan(valuePtrs...)
	require.NoError(t, err)

	t.Logf("Scanned values: %v", values)

	// Verify values were populated
	foundValue := false
	for _, v := range values {
		if v == "value1" || v == int64(42) {
			foundValue = true
		}
	}
	assert.True(t, foundValue, "should have scanned data")

	// Test Next again - should be no more rows
	hasNext = mockRows.Next()
	assert.False(t, hasNext, "should not have more rows")
}

func TestGetPage_Pagination(t *testing.T) {
	tests := []struct {
		name          string
		rowData       []map[string]interface{}
		pageSize      int64
		wantObjectLen int
		wantHasCursor bool
		wantCursor    string
	}{
		{
			name: "returns_cursor_when_more_rows_exist",
			rowData: []map[string]interface{}{
				{"user_id": "1", "name": "Alice", TotalRemainingRowsColumn: int64(3)},
				{"user_id": "2", "name": "Bob", TotalRemainingRowsColumn: int64(3)},
				{"user_id": "3", "name": "Charlie", TotalRemainingRowsColumn: int64(3)},
			},
			pageSize:      2,
			wantObjectLen: 2,
			wantHasCursor: true,
			wantCursor:    "2",
		},
		{
			name: "no_cursor_when_all_rows_returned",
			rowData: []map[string]interface{}{
				{"user_id": "1", "name": "Alice", TotalRemainingRowsColumn: int64(2)},
				{"user_id": "2", "name": "Bob", TotalRemainingRowsColumn: int64(2)},
			},
			pageSize:      2,
			wantObjectLen: 2,
			wantHasCursor: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := createMockClientForGetPage(nil, tt.rowData)
			ds := &Datasource{Client: mockClient}
			request := &Request{
				Username: "testuser",
				Password: "testpass",
				BaseURL:  "localhost",
				Database: "TESTDB",
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{ExternalId: "user_id", Type: framework.AttributeTypeString, UniqueId: true},
						{ExternalId: "name", Type: framework.AttributeTypeString},
					},
				},
				UniqueAttributeExternalID: "user_id",
				PageSize:                  tt.pageSize,
			}

			response, err := ds.GetPage(context.Background(), request)

			require.Nil(t, err)
			require.Len(t, response.Objects, tt.wantObjectLen)

			if tt.wantHasCursor {
				require.NotNil(t, response.NextCursor, "should have next cursor")
				assert.Equal(t, tt.wantCursor, *response.NextCursor)
			} else {
				assert.Nil(t, response.NextCursor, "should not have next cursor")
			}
		})
	}
}

func TestValueToString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string_value",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "byte_slice_value",
			input:    []byte("world"),
			expected: "world",
		},
		{
			name:     "int_value",
			input:    42,
			expected: "42",
		},
		{
			name:     "int32_value",
			input:    int32(100),
			expected: "100",
		},
		{
			name:     "int64_value",
			input:    int64(999),
			expected: "999",
		},
		{
			name:     "float32_value",
			input:    float32(3.14),
			expected: "3.140000",
		},
		{
			name:     "float64_value",
			input:    float64(2.718),
			expected: "2.718000",
		},
		{
			name:     "bool_true",
			input:    true,
			expected: "true",
		},
		{
			name:     "bool_false",
			input:    false,
			expected: "false",
		},
		{
			name:     "custom_type_uses_sprintf",
			input:    struct{ Name string }{Name: "test"},
			expected: "{test}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := valueToString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCastAttributeValue(t *testing.T) {
	logger := createTestLogger()
	processor := &queryResultProcessor{logger: logger}

	tests := []struct {
		name          string
		value         interface{}
		attrType      framework.AttributeType
		attrName      string
		expected      interface{}
		expectError   bool
		errorContains string
	}{
		{
			name:     "bool_true_string",
			value:    "true",
			attrType: framework.AttributeTypeBool,
			attrName: "is_active",
			expected: true,
		},
		{
			name:     "bool_false_string",
			value:    "false",
			attrType: framework.AttributeTypeBool,
			attrName: "is_active",
			expected: false,
		},
		{
			name:     "bool_1_string",
			value:    "1",
			attrType: framework.AttributeTypeBool,
			attrName: "is_active",
			expected: true,
		},
		{
			name:     "bool_0_string",
			value:    "0",
			attrType: framework.AttributeTypeBool,
			attrName: "is_active",
			expected: false,
		},
		{
			name:     "bool_invalid_falls_back_to_string",
			value:    "notabool",
			attrType: framework.AttributeTypeBool,
			attrName: "is_active",
			expected: "notabool",
		},
		{
			name:     "int64_returns_float64",
			value:    "42",
			attrType: framework.AttributeTypeInt64,
			attrName: "count",
			expected: float64(42),
		},
		{
			name:     "double_parses_decimal",
			value:    "3.14159",
			attrType: framework.AttributeTypeDouble,
			attrName: "rate",
			expected: float64(3.14159),
		},
		{
			name:     "double_invalid_falls_back_to_string",
			value:    "not_a_number",
			attrType: framework.AttributeTypeDouble,
			attrName: "rate",
			expected: "not_a_number",
		},
		{
			name:     "string_type_returns_string",
			value:    "hello world",
			attrType: framework.AttributeTypeString,
			attrName: "name",
			expected: "hello world",
		},
		{
			name:     "duration_type_returns_string",
			value:    "PT1H30M",
			attrType: framework.AttributeTypeDuration,
			attrName: "duration",
			expected: "PT1H30M",
		},
		{
			name:     "datetime_type_returns_string",
			value:    "2024-01-15T10:30:00Z",
			attrType: framework.AttributeTypeDateTime,
			attrName: "created_at",
			expected: "2024-01-15T10:30:00Z",
		},
		{
			name:          "unsupported_type_returns_error",
			value:         "value",
			attrType:      framework.AttributeType(999),
			attrName:      "unknown",
			expectError:   true,
			errorContains: "Unsupported attribute type",
		},
		{
			name:     "byte_array_converted_to_string",
			value:    []byte("bytes"),
			attrType: framework.AttributeTypeString,
			attrName: "data",
			expected: "bytes",
		},
		{
			name:     "int64_value_converted",
			value:    int64(100),
			attrType: framework.AttributeTypeInt64,
			attrName: "count",
			expected: float64(100),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.castAttributeValue(tt.value, tt.attrType, tt.attrName)

			if tt.expectError {
				require.NotNil(t, err)
				assert.Contains(t, err.Message, tt.errorContains)
			} else {
				require.Nil(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBuildObject(t *testing.T) {
	logger := createTestLogger()

	t.Run("filters_to_requested_attributes_only", func(t *testing.T) {
		processor := newQueryResultProcessor(
			[]*framework.AttributeConfig{
				{ExternalId: "name", Type: framework.AttributeTypeString},
				{ExternalId: "email", Type: framework.AttributeTypeString},
			},
			nil,
			logger,
		)

		allColumns := map[string]interface{}{
			"name":   "John",
			"email":  "john@example.com",
			"secret": "should_not_appear",
			"salary": 50000,
		}

		obj, err := processor.buildObject(allColumns)

		require.Nil(t, err)
		assert.Equal(t, "John", obj["name"])
		assert.Equal(t, "john@example.com", obj["email"])
		_, secretExists := obj["secret"]
		assert.False(t, secretExists)
		_, salaryExists := obj["salary"]
		assert.False(t, salaryExists)
	})

	t.Run("handles_nil_values", func(t *testing.T) {
		processor := newQueryResultProcessor(
			[]*framework.AttributeConfig{
				{ExternalId: "name", Type: framework.AttributeTypeString},
			},
			nil,
			logger,
		)

		allColumns := map[string]interface{}{
			"name": nil,
		}

		obj, err := processor.buildObject(allColumns)

		require.Nil(t, err)
		val, exists := obj["name"]
		assert.True(t, exists)
		assert.Nil(t, val)
	})

	t.Run("skips_missing_columns", func(t *testing.T) {
		processor := newQueryResultProcessor(
			[]*framework.AttributeConfig{
				{ExternalId: "name", Type: framework.AttributeTypeString},
				{ExternalId: "missing_col", Type: framework.AttributeTypeString},
			},
			nil,
			logger,
		)

		allColumns := map[string]interface{}{
			"name": "John",
		}

		obj, err := processor.buildObject(allColumns)

		require.Nil(t, err)
		assert.Equal(t, "John", obj["name"])
		_, exists := obj["missing_col"]
		assert.False(t, exists)
	})

	t.Run("builds_composite_id_when_id_requested", func(t *testing.T) {
		processor := newQueryResultProcessor(
			[]*framework.AttributeConfig{
				{ExternalId: "id", Type: framework.AttributeTypeString},
				{ExternalId: "name", Type: framework.AttributeTypeString},
			},
			[]string{"MANDT", "EBELN"},
			logger,
		)

		allColumns := map[string]interface{}{
			"MANDT": "100",
			"EBELN": "4500001234",
			"name":  "Test",
		}

		obj, err := processor.buildObject(allColumns)

		require.Nil(t, err)
		assert.Equal(t, "100|4500001234", obj["id"])
		assert.Equal(t, "100|4500001234", obj["composite_id"])
		assert.Equal(t, "Test", obj["name"])
	})

	t.Run("no_composite_id_without_unique_key_columns", func(t *testing.T) {
		processor := newQueryResultProcessor(
			[]*framework.AttributeConfig{
				{ExternalId: "id", Type: framework.AttributeTypeString},
			},
			nil,
			logger,
		)

		allColumns := map[string]interface{}{
			"id": "123",
		}

		obj, err := processor.buildObject(allColumns)

		require.Nil(t, err)
		_, idExists := obj["id"]
		assert.False(t, idExists, "id should not be set when no unique key columns")
	})
}

func TestScanRowToMap(t *testing.T) {
	logger := createTestLogger()
	processor := newQueryResultProcessor(nil, nil, logger)

	t.Run("scans_all_columns", func(t *testing.T) {
		data := []map[string]interface{}{
			{"col1": "value1", "col2": int64(42), "col3": true},
		}
		mockRows := createMockRowsWithData(data)
		mockRows.Next()

		columns, _ := mockRows.Columns()
		result, err := processor.scanRowToMap(mockRows, columns)

		require.NoError(t, err)
		assert.Contains(t, result, "col1")
		assert.Contains(t, result, "col2")
		assert.Contains(t, result, "col3")
	})

	t.Run("converts_byte_arrays_to_strings", func(t *testing.T) {
		data := []map[string]interface{}{
			{"data": []byte("byte_content")},
		}
		mockRows := createMockRowsWithData(data)
		mockRows.Next()

		columns, _ := mockRows.Columns()
		result, err := processor.scanRowToMap(mockRows, columns)

		require.NoError(t, err)
		assert.Equal(t, "byte_content", result["data"])
	})

	t.Run("preserves_nil_values", func(t *testing.T) {
		data := []map[string]interface{}{
			{"nullable_col": nil},
		}
		mockRows := createMockRowsWithData(data)
		mockRows.Next()

		columns, _ := mockRows.Columns()
		result, err := processor.scanRowToMap(mockRows, columns)

		require.NoError(t, err)
		val, exists := result["nullable_col"]
		assert.True(t, exists)
		assert.Nil(t, val)
	})
}

func TestProcessorProcess(t *testing.T) {
	logger := createTestLogger()

	t.Run("processes_multiple_rows", func(t *testing.T) {
		processor := newQueryResultProcessor(
			[]*framework.AttributeConfig{
				{ExternalId: "name", Type: framework.AttributeTypeString},
			},
			nil,
			logger,
		)

		data := []map[string]interface{}{
			{"name": "Alice", TotalRemainingRowsColumn: int64(3)},
			{"name": "Bob", TotalRemainingRowsColumn: int64(3)},
			{"name": "Charlie", TotalRemainingRowsColumn: int64(3)},
		}
		mockRows := createMockRowsWithData(data)

		results, err := processor.process(mockRows)

		require.Nil(t, err)
		require.Len(t, results.objects, 3)
		assert.Equal(t, "Alice", results.objects[0]["name"])
		assert.Equal(t, "Bob", results.objects[1]["name"])
		assert.Equal(t, "Charlie", results.objects[2]["name"])
	})

	t.Run("extracts_total_count", func(t *testing.T) {
		processor := newQueryResultProcessor(
			[]*framework.AttributeConfig{
				{ExternalId: "name", Type: framework.AttributeTypeString},
			},
			nil,
			logger,
		)

		data := []map[string]interface{}{
			{"name": "Test", TotalRemainingRowsColumn: int64(500)},
		}
		mockRows := createMockRowsWithData(data)

		results, err := processor.process(mockRows)

		require.Nil(t, err)
		assert.Equal(t, int64(500), results.totalCount)
	})

	t.Run("handles_empty_result_set", func(t *testing.T) {
		processor := newQueryResultProcessor(
			[]*framework.AttributeConfig{
				{ExternalId: "name", Type: framework.AttributeTypeString},
			},
			nil,
			logger,
		)

		mockRows := createMockRowsWithData([]map[string]interface{}{})

		results, err := processor.process(mockRows)

		require.Nil(t, err)
		assert.Empty(t, results.objects)
		assert.Equal(t, int64(0), results.totalCount)
	})

	t.Run("returns_error_on_columns_failure", func(t *testing.T) {
		processor := newQueryResultProcessor(nil, nil, logger)

		mockRows := &MockRows{
			columnsFunc: func() ([]string, error) {
				return nil, assert.AnError
			},
		}

		results, err := processor.process(mockRows)

		require.NotNil(t, err)
		assert.Nil(t, results)
		assert.Contains(t, err.Message, "Error getting column information")
	})

	t.Run("returns_error_on_rows_iteration_error", func(t *testing.T) {
		processor := newQueryResultProcessor(
			[]*framework.AttributeConfig{
				{ExternalId: "name", Type: framework.AttributeTypeString},
			},
			nil,
			logger,
		)

		data := []map[string]interface{}{
			{"name": "Test"},
		}
		mockRows := createMockRowsWithData(data)
		mockRows.errFunc = func() error {
			return assert.AnError
		}

		results, err := processor.process(mockRows)

		require.NotNil(t, err)
		assert.Nil(t, results)
		assert.Contains(t, err.Message, "Error reading DB2 query results")
	})

	t.Run("continues_on_scan_error", func(t *testing.T) {
		processor := newQueryResultProcessor(
			[]*framework.AttributeConfig{
				{ExternalId: "name", Type: framework.AttributeTypeString},
			},
			nil,
			logger,
		)

		// Create mock that fails on first row but succeeds on second
		scanCount := 0
		mockRows := &MockRows{
			Data: []map[string]interface{}{
				{"name": "Bad"},
				{"name": "Good"},
			},
			columnsFunc: func() ([]string, error) {
				return []string{"name"}, nil
			},
			scanFunc: func(dest ...interface{}) error {
				scanCount++
				if scanCount == 1 {
					return assert.AnError
				}
				if ptr, ok := dest[0].(*interface{}); ok {
					*ptr = "Good"
				}

				return nil
			},
		}

		results, err := processor.process(mockRows)

		require.Nil(t, err)
		assert.Len(t, results.objects, 1)
		assert.Equal(t, "Good", results.objects[0]["name"])
	})
}

func TestGenerateCursor(t *testing.T) {
	tests := []struct {
		name             string
		lastRow          map[string]interface{}
		uniqueAttrID     string
		uniqueKeyColumns []string
		expected         string
	}{
		{
			name:             "single_column_cursor",
			lastRow:          map[string]interface{}{"user_id": "123"},
			uniqueAttrID:     "user_id",
			uniqueKeyColumns: []string{},
			expected:         "123",
		},
		{
			name:             "composite_key_cursor",
			lastRow:          map[string]interface{}{"composite_id": "100|4500001234|10"},
			uniqueAttrID:     "id",
			uniqueKeyColumns: []string{"MANDT", "EBELN", "EBELP"},
			expected:         "100|4500001234|10",
		},
		{
			name:             "missing_attribute_returns_empty",
			lastRow:          map[string]interface{}{"other": "value"},
			uniqueAttrID:     "missing",
			uniqueKeyColumns: []string{},
			expected:         "",
		},
		{
			name:             "empty_unique_attr_returns_empty",
			lastRow:          map[string]interface{}{"id": "123"},
			uniqueAttrID:     "",
			uniqueKeyColumns: []string{},
			expected:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateCursor(tt.lastRow, tt.uniqueAttrID, tt.uniqueKeyColumns)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// createMockClientForGetPage creates a MockSQLClient that handles both constraint and data queries.
// constraintData can be nil if no constraints should be returned.
func createMockClientForGetPage(constraintData, rowData []map[string]interface{}) *MockSQLClient {
	callCount := 0

	return &MockSQLClient{
		ConnectFunc: func(_ string) (*sql.DB, error) { return nil, nil },
		QueryFunc: func(_ context.Context, _ string, _ ...interface{}) (Rows, error) {
			callCount++
			// First query is for constraints (SYSCAT.TABCONST)
			if callCount == 1 {
				if constraintData == nil {
					constraintData = []map[string]interface{}{}
				}

				return createMockRowsForConstraints(constraintData), nil
			}

			// Second query is for actual data
			return createMockRowsWithData(rowData), nil
		},
	}
}

// setStringFromMap extracts a string value from a map and sets it to the destination pointer.
func setStringFromMap(row map[string]interface{}, key string, dest interface{}) {
	ptr, ok := dest.(*string)
	if !ok {
		return
	}

	v, ok := row[key].(string)
	if ok {
		*ptr = v
	}
}

// setIntFromMap extracts an int value from a map and sets it to the destination pointer.
func setIntFromMap(row map[string]interface{}, key string, dest interface{}) {
	ptr, ok := dest.(*int)
	if !ok {
		return
	}

	v, ok := row[key].(int)
	if ok {
		*ptr = v
	}
}

// createMockRowsForConstraints creates mock rows for constraint queries.
func createMockRowsForConstraints(data []map[string]interface{}) *MockRows {
	columns := []string{"CONSTNAME", "TABNAME", "COLNAME", "COLSEQ"}
	mockRows := &MockRows{
		Data: data,
		columnsFunc: func() ([]string, error) {
			return columns, nil
		},
	}
	mockRows.scanFunc = func(dest ...interface{}) error {
		if mockRows.index == 0 || mockRows.index > len(data) {
			return nil
		}

		row := data[mockRows.index-1]
		if len(dest) >= 4 {
			setStringFromMap(row, "CONSTNAME", dest[0])
			setStringFromMap(row, "TABNAME", dest[1])
			setStringFromMap(row, "COLNAME", dest[2])
			setIntFromMap(row, "COLSEQ", dest[3])
		}

		return nil
	}

	return mockRows
}

// createMockRowsWithData creates a MockRows with proper scan function for generic data.
func createMockRowsWithData(data []map[string]interface{}) *MockRows {
	// Build consistent column list from first row
	var columns []string
	if len(data) > 0 {
		// Use sorted keys for consistent column order
		for col := range data[0] {
			columns = append(columns, col)
		}
		// Sort for consistency
		sortStrings(columns)
	}

	mockRows := &MockRows{
		Data: data,
		columnsFunc: func() ([]string, error) {
			return columns, nil
		},
	}

	mockRows.scanFunc = func(dest ...interface{}) error {
		if mockRows.index == 0 || mockRows.index > len(data) {
			return nil
		}

		row := data[mockRows.index-1]
		for i, col := range columns {
			if i < len(dest) {
				// dest[i] is a *interface{} pointing to values[i]
				if ptr, ok := dest[i].(*interface{}); ok {
					*ptr = row[col]
				}
			}
		}

		return nil
	}

	return mockRows
}

// sortStrings sorts a slice of strings in place.
func sortStrings(s []string) {
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] > s[j] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}
