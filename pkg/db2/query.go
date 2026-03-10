// Copyright 2026 SGNL.ai, Inc.
package db2

import (
	"errors"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/condexpr"
	condexprsql "github.com/sgnl-ai/adapters/pkg/condexpr/sql"
)

// TotalRemainingRowsColumn is the column name for the COUNT window function result
// that provides the total count of rows matching the query conditions.
const TotalRemainingRowsColumn = "total_remaining_rows"

// quoteIdentifier quotes DB2 identifiers (table names, column names) that contain special characters.
// DB2 supports delimited identifiers wrapped in double quotes for names containing /, -, or spaces.
// Any embedded double quotes in the identifier are escaped by doubling them (SQL standard).
func quoteIdentifier(name string) string {
	hasSpecialChars := strings.Contains(name, "/") || strings.Contains(name, "-") ||
		strings.Contains(name, " ") || strings.Contains(name, `"`)
	if hasSpecialChars {
		// Escape embedded double quotes by doubling them per SQL standard
		escaped := strings.ReplaceAll(name, `"`, `""`)

		return fmt.Sprintf(`"%s"`, escaped)
	}

	return name
}

// buildSelectColumns constructs the list of columns for the SELECT clause.
// It handles:
// - Regular attribute columns (quoted if containing special characters)
// - Composite key columns (added when "id" attribute is requested)
// - Deduplication of columns
// Returns the column list and whether composite key columns are needed.
func buildSelectColumns(attributes []*framework.AttributeConfig, uniqueKeyColumns []string) ([]string, bool) {
	if len(attributes) == 0 {
		return []string{"*"}, false
	}

	var columns []string

	columnsSet := make(map[string]bool)
	needsCompositeKeyColumns := false

	for _, attr := range attributes {
		if attr.ExternalId == "id" {
			needsCompositeKeyColumns = true
		} else {
			quotedColumn := quoteIdentifier(attr.ExternalId)
			columns = append(columns, quotedColumn)
			columnsSet[attr.ExternalId] = true
		}
	}

	if needsCompositeKeyColumns && len(uniqueKeyColumns) > 0 {
		for _, keyCol := range uniqueKeyColumns {
			if !columnsSet[keyCol] {
				quotedKeyCol := quoteIdentifier(keyCol)
				columns = append(columns, quotedKeyCol)
				columnsSet[keyCol] = true
			}
		}
	}

	return columns, needsCompositeKeyColumns
}

// buildPaginationCastExpression builds the str_id expression for cursor-based pagination.
// For single column keys, it casts the column to VARCHAR.
// For composite keys, it concatenates all key columns with '|' separator.
func buildPaginationCastExpression(uniqueAttrID string, uniqueKeyColumns []string) string {
	if uniqueAttrID == "" {
		return ""
	}

	if uniqueAttrID != "id" {
		quotedAttr := quoteIdentifier(uniqueAttrID)

		return fmt.Sprintf("CAST(%s AS VARCHAR(50)) AS str_id", quotedAttr)
	}

	if len(uniqueKeyColumns) > 0 {
		castParts := make([]string, len(uniqueKeyColumns))
		for i, col := range uniqueKeyColumns {
			quotedCol := quoteIdentifier(col)
			castParts[i] = fmt.Sprintf("CAST(%s AS VARCHAR(50))", quotedCol)
		}

		return fmt.Sprintf("(%s) AS str_id", strings.Join(castParts, " || '|' || "))
	}

	return ""
}

// buildOrderByClause constructs the ORDER BY clause for consistent pagination.
// For composite keys, it orders by all key columns cast to VARCHAR.
// For single column keys, it orders by that column cast to VARCHAR.
func buildOrderByClause(uniqueAttrID string, uniqueKeyColumns []string) string {
	if uniqueAttrID == "id" && len(uniqueKeyColumns) > 0 {
		orderColumns := make([]string, len(uniqueKeyColumns))
		for i, col := range uniqueKeyColumns {
			quotedCol := quoteIdentifier(col)
			orderColumns[i] = fmt.Sprintf("CAST(%s AS VARCHAR(50))", quotedCol)
		}

		return fmt.Sprintf("ORDER BY %s", strings.Join(orderColumns, ", "))
	}

	if uniqueAttrID != "" {
		quotedAttr := quoteIdentifier(uniqueAttrID)

		return fmt.Sprintf("ORDER BY CAST(%s AS VARCHAR(50))", quotedAttr)
	}

	return ""
}

// buildCursorCondition constructs the WHERE condition for cursor-based pagination.
// For composite keys, it builds a tuple comparison: (col1, col2) > (?, ?).
// For single column keys, it builds: col > ?.
// Returns the condition string and the arguments to bind.
func buildCursorCondition(cursor *string, uniqueAttrID string, uniqueKeyColumns []string) (string, []any) {
	if cursor == nil || *cursor == "" {
		return "", nil
	}

	if uniqueAttrID == "id" && len(uniqueKeyColumns) > 0 {
		cursorParts := strings.Split(*cursor, "|")
		if len(cursorParts) != len(uniqueKeyColumns) {
			return "", nil
		}

		colList := make([]string, len(uniqueKeyColumns))
		for i, col := range uniqueKeyColumns {
			quotedCol := quoteIdentifier(col)
			colList[i] = fmt.Sprintf("CAST(%s AS VARCHAR(50))", quotedCol)
		}

		args := make([]any, len(cursorParts))
		for i, part := range cursorParts {
			args[i] = part
		}

		condition := fmt.Sprintf("(%s) > (%s)",
			strings.Join(colList, ", "),
			strings.Repeat("?, ", len(cursorParts)-1)+"?")

		return condition, args
	}

	if uniqueAttrID != "" {
		quotedAttr := quoteIdentifier(uniqueAttrID)

		return fmt.Sprintf("CAST(%s AS VARCHAR(50)) > ?", quotedAttr), []any{*cursor}
	}

	return "", nil
}

// buildFilterCondition converts the filter expression to SQL WHERE condition.
// Returns the condition string and the arguments to bind.
func buildFilterCondition(filter *condexpr.Condition) (string, []any, error) {
	if filter == nil {
		return "", nil, nil
	}

	builder := condexprsql.NewConditionBuilder()

	whereExpr, err := builder.Build(*filter)
	if err != nil {
		return "", nil, err
	}

	tempDS := goqu.Select().Where(whereExpr).Prepared(true)

	filterSQL, filterArgs, err := tempDS.ToSQL()
	if err != nil {
		return "", nil, fmt.Errorf("failed to convert filter to SQL: %w", err)
	}

	filterSQL = strings.TrimPrefix(filterSQL, "SELECT * WHERE ")

	return filterSQL, filterArgs, nil
}

func ConstructQuery(request *Request) (string, []any, error) {
	if request == nil {
		return "", nil, errors.New("nil request provided")
	}

	// Build column list from requested attributes and unique key columns
	columns, _ := buildSelectColumns(request.EntityConfig.Attributes, request.UniqueKeyColumns)

	columnList := strings.Join(columns, ", ")
	if len(columns) == 0 {
		columnList = "1"
	}

	// Add pagination cast expression if needed
	paginationExpr := buildPaginationCastExpression(
		request.UniqueAttributeExternalID, request.UniqueKeyColumns)
	if paginationExpr != "" {
		columnList = fmt.Sprintf("%s, %s", columnList, paginationExpr)
	}

	selectClause := fmt.Sprintf("SELECT %s, COUNT(*) OVER() AS %s", columnList, TotalRemainingRowsColumn)

	// FROM clause with optional schema handling
	tableName := quoteIdentifier(request.EntityConfig.ExternalId)
	if request.Schema != "" {
		tableName = fmt.Sprintf("%s.%s", quoteIdentifier(request.Schema), tableName)
	}

	fromClause := fmt.Sprintf("FROM %s", tableName)

	// Build WHERE conditions
	var whereConditions []string

	args := []any{}

	// Add cursor condition
	cursorCond, cursorArgs := buildCursorCondition(
		request.Cursor, request.UniqueAttributeExternalID, request.UniqueKeyColumns)
	if cursorCond != "" {
		whereConditions = append(whereConditions, cursorCond)
		args = append(args, cursorArgs...)
	}

	// Add filter conditions
	if request.Filter != nil {
		filterSQL, filterArgs, err := buildFilterCondition(request.Filter)
		if err != nil {
			return "", nil, err
		}

		if filterSQL != "" {
			whereConditions = append(whereConditions, filterSQL)
			args = append(args, filterArgs...)
		}
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// ORDER BY and LIMIT clauses
	orderByClause := buildOrderByClause(request.UniqueAttributeExternalID, request.UniqueKeyColumns)

	limitClause := ""
	if request.PageSize > 0 {
		limitClause = fmt.Sprintf("FETCH FIRST %d ROWS ONLY", request.PageSize+1)
	}

	// Combine all clauses
	query := strings.Join([]string{
		selectClause,
		fromClause,
		whereClause,
		orderByClause,
		limitClause,
	}, " ")

	// Clean up extra whitespace
	query = strings.Join(strings.Fields(query), " ")

	return query, args, nil
}
