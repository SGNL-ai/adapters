// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// UniqueConstraintColumn represents a column in a unique constraint.
type UniqueConstraintColumn struct {
	ColumnName string
	Position   int
}

// uniqueConstraint represents a unique constraint on a table.
type uniqueConstraint struct {
	constraintName string
	tableName      string
	columns        []UniqueConstraintColumn
}

// getUniqueConstraints queries DB2 system tables to find unique constraints for a given table.
func (d *Datasource) getUniqueConstraints(ctx context.Context, tableName string) ([]uniqueConstraint, error) {
	// DB2 system catalog query to find unique constraints
	// This queries SYSCAT.TABCONST (table constraints) and SYSCAT.KEYCOLUSE (key column usage)
	query := `
		SELECT
			tc.CONSTNAME,
			tc.TABNAME,
			kcu.COLNAME,
			kcu.COLSEQ
		FROM SYSCAT.TABCONST tc
		JOIN SYSCAT.KEYCOLUSE kcu ON tc.CONSTNAME = kcu.CONSTNAME
			AND tc.TABSCHEMA = kcu.TABSCHEMA
			AND tc.TABNAME = kcu.TABNAME
		WHERE tc.TABNAME = ?
			AND tc.TABSCHEMA = CURRENT SCHEMA
			AND tc.TYPE IN ('P', 'U')
		ORDER BY tc.CONSTNAME, kcu.COLSEQ
	`

	rows, err := d.Client.Query(ctx, query, strings.ToUpper(tableName))
	if err != nil {
		return nil, fmt.Errorf("failed to query unique constraints: %w", err)
	}
	defer rows.Close()

	return processConstraintRows(rows)
}

// processConstraintRows iterates over constraint query results and builds uniqueConstraint structs.
// It groups columns by constraint name and sorts them by position.
func processConstraintRows(rows Rows) ([]uniqueConstraint, error) {
	constraintMap := make(map[string]*uniqueConstraint)

	for rows.Next() {
		var constraintName, tableNameResult, columnName string

		var position int

		if err := rows.Scan(&constraintName, &tableNameResult, &columnName, &position); err != nil {
			return nil, fmt.Errorf("failed to scan constraint row: %w", err)
		}

		if constraint, exists := constraintMap[constraintName]; exists {
			constraint.columns = append(constraint.columns, UniqueConstraintColumn{
				ColumnName: columnName,
				Position:   position,
			})
		} else {
			constraintMap[constraintName] = &uniqueConstraint{
				constraintName: constraintName,
				tableName:      tableNameResult,
				columns: []UniqueConstraintColumn{
					{
						ColumnName: columnName,
						Position:   position,
					},
				},
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over constraint rows: %w", err)
	}

	// Convert map to slice and sort columns by position
	constraints := make([]uniqueConstraint, 0, len(constraintMap))
	for _, constraint := range constraintMap {
		// Sort columns by their position in the constraint
		sort.Slice(constraint.columns, func(i, j int) bool {
			return constraint.columns[i].Position < constraint.columns[j].Position
		})
		constraints = append(constraints, *constraint)
	}

	return constraints, nil
}

// getPrimaryKey returns the primary key constraint for a table, or nil if none exists.
func (d *Datasource) getPrimaryKey(ctx context.Context, tableName string) (*uniqueConstraint, error) {
	constraints, err := d.getUniqueConstraints(ctx, tableName)
	if err != nil {
		return nil, err
	}

	if len(constraints) == 0 {
		return nil, nil
	}

	// Look for primary key constraint (usually named with a pattern or has type 'P')
	for _, constraint := range constraints {
		// Primary keys in DB2 often have names containing "PK" or are the first constraint
		// We could make this more sophisticated by checking the constraint type
		if strings.Contains(strings.ToUpper(constraint.constraintName), "PK") ||
			strings.Contains(strings.ToUpper(constraint.constraintName), "PRIMARY") {
			return &constraint, nil
		}
	}

	// If no explicit primary key found, return the first unique constraint
	return &constraints[0], nil
}

// BuildCompositeID creates a composite ID by concatenating the values of the key columns.
//
// This is used when a table has a composite primary key (multiple columns) but the entity
// configuration expects a single "id" attribute. The function concatenates the values of
// all key columns using the specified separator to produce a unique string identifier.
//
// Example: For a row with MANDT="100", EBELN="4500001234", EBELP="10" and separator "|",
// the resulting composite ID would be "100|4500001234|10".
//
// The columns slice determines the order of concatenation (sorted by Position).
// Missing or nil values result in empty strings in the composite ID.
func BuildCompositeID(row map[string]interface{}, columns []UniqueConstraintColumn, separator string) string {
	parts := make([]string, 0, len(columns))

	for _, col := range columns {
		value := ""
		if val, exists := row[col.ColumnName]; exists && val != nil {
			value = fmt.Sprintf("%v", val)
		}

		parts = append(parts, value)
	}

	return strings.Join(parts, separator)
}

// ExtractUniqueKeyColumns returns the column names that make up the unique key for a table.
//
// This function queries the DB2 system catalog (SYSCAT.TABCONST and SYSCAT.KEYCOLUSE) to
// discover the primary key or unique constraint columns for the given table. The returned
// columns are used for:
//
//  1. Composite ID Generation: When an entity is configured with a synthetic "id" attribute,
//     the adapter needs to generate a unique identifier from the actual database columns.
//     For tables with composite primary keys (e.g., MANDT|EBELN|EBELP in SAP tables), the
//     unique key columns are concatenated with a separator to form a single ID value.
//
//  2. Cursor-based Pagination: The unique key columns are used to construct ORDER BY clauses
//     and WHERE conditions for keyset pagination. This ensures consistent ordering across
//     pages and allows resuming from the last fetched row.
//
//  3. Query Construction: The columns are added to the SELECT clause to ensure they're
//     available for ID generation even if not explicitly requested in the entity attributes.
//
// The function prioritizes primary keys (constraints with "PK" or "PRIMARY" in the name)
// over unique constraints. If no primary key is found, it falls back to the first unique
// constraint available.
func (d *Datasource) ExtractUniqueKeyColumns(ctx context.Context, tableName string) ([]string, error) {
	primaryKey, err := d.getPrimaryKey(ctx, tableName)
	if err != nil {
		return nil, err
	}

	if primaryKey == nil {
		return nil, fmt.Errorf("no unique key found for table %s", tableName)
	}

	columnNames := make([]string, 0, len(primaryKey.columns))
	for _, col := range primaryKey.columns {
		columnNames = append(columnNames, col.ColumnName)
	}

	return columnNames, nil
}
