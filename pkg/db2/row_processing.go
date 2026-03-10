// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"fmt"
	"strconv"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"go.uber.org/zap"
)

// queryResultProcessor handles processing of database query results.
type queryResultProcessor struct {
	attributes       []*framework.AttributeConfig
	uniqueKeyColumns []string
	logger           *zap.Logger
}

// processedResults contains the results of processing query rows.
type processedResults struct {
	objects    []map[string]interface{}
	totalCount int64
}

// newQueryResultProcessor creates a new processor for query results.
func newQueryResultProcessor(
	attributes []*framework.AttributeConfig,
	uniqueKeyColumns []string,
	logger *zap.Logger,
) *queryResultProcessor {
	return &queryResultProcessor{
		attributes:       attributes,
		uniqueKeyColumns: uniqueKeyColumns,
		logger:           logger,
	}
}

// process iterates over rows and builds the result objects.
func (p *queryResultProcessor) process(rows Rows) (*processedResults, *framework.Error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Error getting column information: %v.", err),
		}
	}

	results := &processedResults{
		objects: []map[string]interface{}{},
	}

	for rows.Next() {
		allColumns, scanErr := p.scanRowToMap(rows, columns)
		if scanErr != nil {
			p.logger.Warn("Failed to scan row", zap.Error(scanErr))

			continue
		}

		obj, buildErr := p.buildObject(allColumns)
		if buildErr != nil {
			return nil, buildErr
		}

		// Extract total count (don't include in response object)
		if val, exists := allColumns[TotalRemainingRowsColumn]; exists {
			if count, ok := val.(int64); ok {
				results.totalCount = count
			}
		}

		results.objects = append(results.objects, obj)
	}

	if err := rows.Err(); err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Error reading DB2 query results: %v.", err),
		}
	}

	return results, nil
}

// scanRowToMap scans a single row into a map of column names to values.
func (p *queryResultProcessor) scanRowToMap(rows Rows, columns []string) (map[string]interface{}, error) {
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})

	for i, col := range columns {
		val := values[i]
		if val != nil {
			// Convert byte arrays to strings (common in DB2)
			if b, ok := val.([]byte); ok {
				result[col] = string(b)
			} else {
				result[col] = val
			}
		} else {
			result[col] = nil
		}
	}

	return result, nil
}

// buildObject filters and casts attributes, and adds composite ID if needed.
func (p *queryResultProcessor) buildObject(
	allColumns map[string]interface{},
) (map[string]interface{}, *framework.Error) {
	obj := make(map[string]interface{})
	needsCompositeID := false

	// Build map of requested attributes
	requestedAttrs := make(map[string]framework.AttributeConfig)

	for _, attr := range p.attributes {
		if attr.ExternalId == "id" {
			needsCompositeID = true
		} else {
			requestedAttrs[attr.ExternalId] = *attr
		}
	}

	// Add requested attributes with type casting
	for attrName, attr := range requestedAttrs {
		val, exists := allColumns[attrName]
		if !exists {
			continue
		}

		if val == nil {
			obj[attrName] = nil

			continue
		}

		castedVal, err := p.castAttributeValue(val, attr.Type, attrName)
		if err != nil {
			return nil, err
		}

		obj[attrName] = castedVal
	}

	// Build composite ID if needed
	if needsCompositeID && len(p.uniqueKeyColumns) > 0 {
		keyColumns := make([]UniqueConstraintColumn, len(p.uniqueKeyColumns))
		for i, colName := range p.uniqueKeyColumns {
			keyColumns[i] = UniqueConstraintColumn{
				ColumnName: colName,
				Position:   i + 1,
			}
		}

		compositeID := BuildCompositeID(allColumns, keyColumns, "|")
		obj["id"] = compositeID
		obj["composite_id"] = compositeID // Keep for cursor generation
	}

	return obj, nil
}

// castAttributeValue converts a value to the expected attribute type.
func (p *queryResultProcessor) castAttributeValue(
	val interface{},
	attrType framework.AttributeType,
	attrName string,
) (interface{}, *framework.Error) {
	// Convert value to string first
	str := valueToString(val)

	var (
		result  interface{}
		castErr error
	)

	switch attrType {
	case framework.AttributeTypeBool:
		result, castErr = strconv.ParseBool(str)
	case framework.AttributeTypeDouble, framework.AttributeTypeInt64:
		// The adapter framework expects all numbers to be passed as floats
		result, castErr = strconv.ParseFloat(str, 64)
	case framework.AttributeTypeString, framework.AttributeTypeDuration, framework.AttributeTypeDateTime:
		result = str
	default:
		return nil, &framework.Error{
			Message: "Unsupported attribute type provided.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
		}
	}

	if castErr != nil {
		p.logger.Warn("Failed to cast column",
			zap.String("column", attrName),
			zap.Any("type", attrType),
			zap.Error(castErr))
		// Return raw string value if casting fails
		result = str
	}

	return result, nil
}

// valueToString converts any value to its string representation.
func valueToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
