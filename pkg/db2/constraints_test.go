// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getUniqueConstraints(t *testing.T) {
	tests := []struct {
		name          string
		tableName     string
		mockData      []map[string]interface{}
		expectedCount int
		expectError   bool
	}{
		{
			name:      "table with primary key",
			tableName: "USERS",
			mockData: []map[string]interface{}{
				{
					"CONSTNAME": "PK_USERS",
					"TABNAME":   "USERS",
					"COLNAME":   "ID",
					"COLSEQ":    1,
				},
			},
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:      "table with composite key",
			tableName: "EKPO",
			mockData: []map[string]interface{}{
				{
					"CONSTNAME": "PK_EKPO",
					"TABNAME":   "EKPO",
					"COLNAME":   "MANDT",
					"COLSEQ":    1,
				},
				{
					"CONSTNAME": "PK_EKPO",
					"TABNAME":   "EKPO",
					"COLNAME":   "EBELN",
					"COLSEQ":    2,
				},
				{
					"CONSTNAME": "PK_EKPO",
					"TABNAME":   "EKPO",
					"COLNAME":   "EBELP",
					"COLSEQ":    3,
				},
			},
			expectedCount: 1,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockSQLClient{
				QueryFunc: func(_ context.Context, _ string, _ ...interface{}) (Rows, error) {
					mockRows := &MockRows{
						Data: tt.mockData,
					}
					// Set up the scan function to populate the provided destinations
					mockRows.scanFunc = func(dest ...interface{}) error {
						if len(tt.mockData) == 0 || mockRows.index == 0 || mockRows.index > len(tt.mockData) {
							return nil
						}

						row := tt.mockData[mockRows.index-1]
						if len(dest) < 4 {
							return nil
						}

						// Populate constraint name
						if constName, ok := dest[0].(*string); ok {
							if val, ok := row["CONSTNAME"].(string); ok {
								*constName = val
							}
						}

						// Populate table name
						if tabName, ok := dest[1].(*string); ok {
							if val, ok := row["TABNAME"].(string); ok {
								*tabName = val
							}
						}

						// Populate column name
						if colName, ok := dest[2].(*string); ok {
							if val, ok := row["COLNAME"].(string); ok {
								*colName = val
							}
						}

						// Populate column sequence
						if colSeq, ok := dest[3].(*int); ok {
							if val, ok := row["COLSEQ"].(int); ok {
								*colSeq = val
							}
						}

						return nil
					}

					return mockRows, nil
				},
			}

			datasource := &Datasource{Client: mockClient}
			constraints, err := datasource.getUniqueConstraints(context.Background(), tt.tableName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, constraints, tt.expectedCount)

				if len(constraints) > 0 {
					// Verify the constraint structure
					constraint := constraints[0]
					assert.Equal(t, tt.tableName, constraint.tableName)
					assert.NotEmpty(t, constraint.constraintName)
					assert.NotEmpty(t, constraint.columns)

					// For composite key test, verify column order
					if tt.name == "table with composite key" {
						assert.Len(t, constraint.columns, 3)
						assert.Equal(t, "MANDT", constraint.columns[0].ColumnName)
						assert.Equal(t, "EBELN", constraint.columns[1].ColumnName)
						assert.Equal(t, "EBELP", constraint.columns[2].ColumnName)
					}
				}
			}
		})
	}
}

func TestProcessConstraintRows(t *testing.T) {
	tests := []struct {
		name           string
		mockData       []map[string]interface{}
		scanErr        error
		iterErr        error
		expectedCount  int
		expectError    bool
		expectedErrMsg string
		validateResult func(t *testing.T, constraints []uniqueConstraint)
	}{
		{
			name:          "empty_result_returns_empty_slice",
			mockData:      []map[string]interface{}{},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name: "single_column_constraint",
			mockData: []map[string]interface{}{
				{"CONSTNAME": "PK_USERS", "TABNAME": "USERS", "COLNAME": "ID", "COLSEQ": 1},
			},
			expectedCount: 1,
			expectError:   false,
			validateResult: func(t *testing.T, constraints []uniqueConstraint) {
				assert.Equal(t, "PK_USERS", constraints[0].constraintName)
				assert.Equal(t, "USERS", constraints[0].tableName)
				require.Len(t, constraints[0].columns, 1)
				assert.Equal(t, "ID", constraints[0].columns[0].ColumnName)
			},
		},
		{
			name: "multiple_constraints_on_same_table",
			mockData: []map[string]interface{}{
				{"CONSTNAME": "PK_ORDERS", "TABNAME": "ORDERS", "COLNAME": "ORDER_ID", "COLSEQ": 1},
				{"CONSTNAME": "UK_ORDERS_REF", "TABNAME": "ORDERS", "COLNAME": "ORDER_REF", "COLSEQ": 1},
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name: "columns_sorted_by_position_regardless_of_input_order",
			mockData: []map[string]interface{}{
				{"CONSTNAME": "PK_TEST", "TABNAME": "TEST", "COLNAME": "COL_C", "COLSEQ": 3},
				{"CONSTNAME": "PK_TEST", "TABNAME": "TEST", "COLNAME": "COL_A", "COLSEQ": 1},
				{"CONSTNAME": "PK_TEST", "TABNAME": "TEST", "COLNAME": "COL_B", "COLSEQ": 2},
			},
			expectedCount: 1,
			expectError:   false,
			validateResult: func(t *testing.T, constraints []uniqueConstraint) {
				require.Len(t, constraints[0].columns, 3)
				assert.Equal(t, "COL_A", constraints[0].columns[0].ColumnName)
				assert.Equal(t, 1, constraints[0].columns[0].Position)
				assert.Equal(t, "COL_B", constraints[0].columns[1].ColumnName)
				assert.Equal(t, 2, constraints[0].columns[1].Position)
				assert.Equal(t, "COL_C", constraints[0].columns[2].ColumnName)
				assert.Equal(t, 3, constraints[0].columns[2].Position)
			},
		},
		{
			name: "scan_error_returns_wrapped_error",
			mockData: []map[string]interface{}{
				{"CONSTNAME": "PK_TEST", "TABNAME": "TEST", "COLNAME": "ID", "COLSEQ": 1},
			},
			scanErr:        errors.New("scan failed"),
			expectError:    true,
			expectedErrMsg: "failed to scan constraint row",
		},
		{
			name:           "rows_iteration_error_returns_wrapped_error",
			mockData:       []map[string]interface{}{},
			iterErr:        errors.New("iteration error"),
			expectError:    true,
			expectedErrMsg: "error iterating over constraint rows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRows := &MockRows{Data: tt.mockData}
			if tt.scanErr != nil {
				mockRows.scanFunc = func(_ ...interface{}) error {
					return tt.scanErr
				}
			} else {
				mockRows.scanFunc = createConstraintScanFunc(tt.mockData, mockRows)
			}
			if tt.iterErr != nil {
				mockRows.errFunc = func() error {
					return tt.iterErr
				}
			}

			// Act
			constraints, err := processConstraintRows(mockRows)

			// Assert
			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, constraints)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Len(t, constraints, tt.expectedCount)
				if tt.validateResult != nil {
					tt.validateResult(t, constraints)
				}
			}
		})
	}
}

func TestGetUniqueConstraintsQueryExecution(t *testing.T) {
	t.Run("query_error_returns_wrapped_error", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("connection failed")
		mockClient := &MockSQLClient{
			QueryFunc: func(_ context.Context, _ string, _ ...interface{}) (Rows, error) {
				return nil, expectedErr
			},
		}
		datasource := &Datasource{Client: mockClient}

		// Act
		constraints, err := datasource.getUniqueConstraints(context.Background(), "USERS")

		// Assert
		require.Error(t, err)
		assert.Nil(t, constraints)
		assert.Contains(t, err.Error(), "failed to query unique constraints")
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("table_name_is_uppercased_in_query", func(t *testing.T) {
		// Arrange
		var capturedArgs []interface{}
		mockClient := &MockSQLClient{
			QueryFunc: func(_ context.Context, _ string, args ...interface{}) (Rows, error) {
				capturedArgs = args

				return &MockRows{Data: []map[string]interface{}{}}, nil
			},
		}
		datasource := &Datasource{Client: mockClient}

		// Act
		_, _ = datasource.getUniqueConstraints(context.Background(), "lowercase_table")

		// Assert
		require.Len(t, capturedArgs, 1)
		assert.Equal(t, "LOWERCASE_TABLE", capturedArgs[0])
	})
}

func TestGetPrimaryKey(t *testing.T) {
	tests := []struct {
		name               string
		mockData           []map[string]interface{}
		queryErr           error
		expectError        bool
		expectedErrMsg     string
		expectNil          bool
		expectedConstraint string
	}{
		{
			name: "constraint_with_pk_in_name",
			mockData: []map[string]interface{}{
				{"CONSTNAME": "PK_USERS", "TABNAME": "USERS", "COLNAME": "ID", "COLSEQ": 1},
			},
			expectedConstraint: "PK_USERS",
		},
		{
			name: "constraint_with_primary_in_name",
			mockData: []map[string]interface{}{
				{"CONSTNAME": "PRIMARY_KEY_1", "TABNAME": "ORDERS", "COLNAME": "ORDER_ID", "COLSEQ": 1},
			},
			expectedConstraint: "PRIMARY_KEY_1",
		},
		{
			name: "pk_constraint_selected_over_unique",
			mockData: []map[string]interface{}{
				{"CONSTNAME": "UK_EMAIL", "TABNAME": "USERS", "COLNAME": "EMAIL", "COLSEQ": 1},
				{"CONSTNAME": "PK_USERS", "TABNAME": "USERS", "COLNAME": "ID", "COLSEQ": 1},
			},
			expectedConstraint: "PK_USERS",
		},
		{
			name: "fallback_to_first_unique_constraint_when_no_pk_pattern",
			mockData: []map[string]interface{}{
				{"CONSTNAME": "UNIQUE_EMAIL", "TABNAME": "USERS", "COLNAME": "EMAIL", "COLSEQ": 1},
			},
			expectedConstraint: "UNIQUE_EMAIL",
		},
		{
			name:      "no_constraints_returns_nil",
			mockData:  []map[string]interface{}{},
			expectNil: true,
		},
		{
			name:           "query_error_propagates",
			queryErr:       errors.New("connection failed"),
			expectError:    true,
			expectedErrMsg: "failed to query unique constraints",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockClient := &MockSQLClient{
				QueryFunc: func(_ context.Context, _ string, _ ...interface{}) (Rows, error) {
					if tt.queryErr != nil {
						return nil, tt.queryErr
					}

					mockRows := &MockRows{Data: tt.mockData}
					mockRows.scanFunc = createConstraintScanFunc(tt.mockData, mockRows)

					return mockRows, nil
				},
			}
			datasource := &Datasource{Client: mockClient}

			// Act
			constraint, err := datasource.getPrimaryKey(context.Background(), "TEST")

			// Assert
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)

				return
			}

			require.NoError(t, err)
			if tt.expectNil {
				assert.Nil(t, constraint)
			} else {
				require.NotNil(t, constraint)
				assert.Equal(t, tt.expectedConstraint, constraint.constraintName)
			}
		})
	}
}

// createConstraintScanFunc is a helper to create a scan function for constraint tests.
func createConstraintScanFunc(data []map[string]interface{}, mockRows *MockRows) func(dest ...interface{}) error {
	return func(dest ...interface{}) error {
		if len(data) == 0 || mockRows.index == 0 || mockRows.index > len(data) {
			return nil
		}

		row := data[mockRows.index-1]
		if len(dest) < 4 {
			return nil
		}

		if constName, ok := dest[0].(*string); ok {
			if val, ok := row["CONSTNAME"].(string); ok {
				*constName = val
			}
		}
		if tabName, ok := dest[1].(*string); ok {
			if val, ok := row["TABNAME"].(string); ok {
				*tabName = val
			}
		}
		if colName, ok := dest[2].(*string); ok {
			if val, ok := row["COLNAME"].(string); ok {
				*colName = val
			}
		}
		if colSeq, ok := dest[3].(*int); ok {
			if val, ok := row["COLSEQ"].(int); ok {
				*colSeq = val
			}
		}

		return nil
	}
}

func TestBuildCompositeID(t *testing.T) {
	tests := []struct {
		name      string
		row       map[string]interface{}
		columns   []UniqueConstraintColumn
		separator string
		expected  string
	}{
		{
			name: "single column ID",
			row: map[string]interface{}{
				"ID": "123",
			},
			columns: []UniqueConstraintColumn{
				{ColumnName: "ID", Position: 1},
			},
			separator: "|",
			expected:  "123",
		},
		{
			name: "composite ID with three columns",
			row: map[string]interface{}{
				"MANDT": "100",
				"EBELN": "4500001234",
				"EBELP": "10",
			},
			columns: []UniqueConstraintColumn{
				{ColumnName: "MANDT", Position: 1},
				{ColumnName: "EBELN", Position: 2},
				{ColumnName: "EBELP", Position: 3},
			},
			separator: "|",
			expected:  "100|4500001234|10",
		},
		{
			name: "composite ID with different separator",
			row: map[string]interface{}{
				"CLIENT": "100",
				"DOC_NO": "ABC123",
			},
			columns: []UniqueConstraintColumn{
				{ColumnName: "CLIENT", Position: 1},
				{ColumnName: "DOC_NO", Position: 2},
			},
			separator: ":",
			expected:  "100:ABC123",
		},
		{
			name: "missing column value",
			row: map[string]interface{}{
				"MANDT": "100",
				// EBELN missing
				"EBELP": "10",
			},
			columns: []UniqueConstraintColumn{
				{ColumnName: "MANDT", Position: 1},
				{ColumnName: "EBELN", Position: 2},
				{ColumnName: "EBELP", Position: 3},
			},
			separator: "|",
			expected:  "100||10",
		},
		{
			name: "nil column value",
			row: map[string]interface{}{
				"MANDT": "100",
				"EBELN": nil,
				"EBELP": "10",
			},
			columns: []UniqueConstraintColumn{
				{ColumnName: "MANDT", Position: 1},
				{ColumnName: "EBELN", Position: 2},
				{ColumnName: "EBELP", Position: 3},
			},
			separator: "|",
			expected:  "100||10",
		},
		{
			name: "empty_columns_returns_empty_string",
			row: map[string]interface{}{
				"ID": "123",
			},
			columns:   []UniqueConstraintColumn{},
			separator: "|",
			expected:  "",
		},
		{
			name: "integer_values_converted_to_string",
			row: map[string]interface{}{
				"ID":    123,
				"COUNT": 456,
			},
			columns: []UniqueConstraintColumn{
				{ColumnName: "ID", Position: 1},
				{ColumnName: "COUNT", Position: 2},
			},
			separator: "|",
			expected:  "123|456",
		},
		{
			name: "float_values_converted_to_string",
			row: map[string]interface{}{
				"AMOUNT": 123.45,
			},
			columns: []UniqueConstraintColumn{
				{ColumnName: "AMOUNT", Position: 1},
			},
			separator: "|",
			expected:  "123.45",
		},
		{
			name: "boolean_values_converted_to_string",
			row: map[string]interface{}{
				"ACTIVE": true,
				"LOCKED": false,
			},
			columns: []UniqueConstraintColumn{
				{ColumnName: "ACTIVE", Position: 1},
				{ColumnName: "LOCKED", Position: 2},
			},
			separator: "|",
			expected:  "true|false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildCompositeID(tt.row, tt.columns, tt.separator)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractUniqueKeyColumns(t *testing.T) {
	tests := []struct {
		name            string
		tableName       string
		mockConstraints []map[string]interface{}
		expectedColumns []string
		expectError     bool
	}{
		{
			name:      "table with simple primary key",
			tableName: "USERS",
			mockConstraints: []map[string]interface{}{
				{
					"CONSTNAME": "PK_USERS",
					"TABNAME":   "USERS",
					"COLNAME":   "ID",
					"COLSEQ":    1,
				},
			},
			expectedColumns: []string{"ID"},
			expectError:     false,
		},
		{
			name:      "table with composite primary key",
			tableName: "EKPO",
			mockConstraints: []map[string]interface{}{
				{
					"CONSTNAME": "PK_EKPO",
					"TABNAME":   "EKPO",
					"COLNAME":   "MANDT",
					"COLSEQ":    1,
				},
				{
					"CONSTNAME": "PK_EKPO",
					"TABNAME":   "EKPO",
					"COLNAME":   "EBELN",
					"COLSEQ":    2,
				},
				{
					"CONSTNAME": "PK_EKPO",
					"TABNAME":   "EKPO",
					"COLNAME":   "EBELP",
					"COLSEQ":    3,
				},
			},
			expectedColumns: []string{"MANDT", "EBELN", "EBELP"},
			expectError:     false,
		},
		{
			name:            "table with no constraints",
			tableName:       "TEMP_TABLE",
			mockConstraints: []map[string]interface{}{},
			expectedColumns: nil,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockSQLClient{
				QueryFunc: func(_ context.Context, _ string, _ ...interface{}) (Rows, error) {
					mockRows := &MockRows{
						Data: tt.mockConstraints,
					}
					mockRows.scanFunc = func(dest ...interface{}) error {
						if len(tt.mockConstraints) == 0 || mockRows.index == 0 || mockRows.index > len(tt.mockConstraints) {
							return nil
						}

						row := tt.mockConstraints[mockRows.index-1]
						if len(dest) < 4 {
							return nil
						}

						// Populate constraint name
						if constName, ok := dest[0].(*string); ok {
							if val, ok := row["CONSTNAME"].(string); ok {
								*constName = val
							}
						}

						// Populate table name
						if tabName, ok := dest[1].(*string); ok {
							if val, ok := row["TABNAME"].(string); ok {
								*tabName = val
							}
						}

						// Populate column name
						if colName, ok := dest[2].(*string); ok {
							if val, ok := row["COLNAME"].(string); ok {
								*colName = val
							}
						}

						// Populate column sequence
						if colSeq, ok := dest[3].(*int); ok {
							if val, ok := row["COLSEQ"].(int); ok {
								*colSeq = val
							}
						}

						return nil
					}

					return mockRows, nil
				},
			}

			datasource := &Datasource{Client: mockClient}
			columns, err := datasource.ExtractUniqueKeyColumns(context.Background(), tt.tableName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedColumns, columns)
			}
		})
	}
}

func TestExtractUniqueKeyColumnsErrorPropagation(t *testing.T) {
	t.Run("query_error_propagates_from_getPrimaryKey", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("connection failed")
		mockClient := &MockSQLClient{
			QueryFunc: func(_ context.Context, _ string, _ ...interface{}) (Rows, error) {
				return nil, expectedErr
			},
		}
		datasource := &Datasource{Client: mockClient}

		// Act
		columns, err := datasource.ExtractUniqueKeyColumns(context.Background(), "TEST")

		// Assert
		require.Error(t, err)
		assert.Nil(t, columns)
		assert.Contains(t, err.Error(), "failed to query unique constraints")
		assert.ErrorIs(t, err, expectedErr)
	})
}
