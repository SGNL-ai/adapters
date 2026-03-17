// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"context"
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"go.uber.org/zap"
)

// CompositeIDKey is the row map key used to store a generated composite ID
// for tables with multi-column primary keys. It is set during row processing
// and read back during cursor generation for pagination.
const CompositeIDKey = "composite_id"

// SyntheticIDAttr is the attribute external ID that signals the adapter should
// generate a unique identifier from the table's primary key columns rather
// than reading an existing column value.
const SyntheticIDAttr = "id"

type Datasource struct {
	Client SQLClient
}

// generateCursor creates a cursor string for pagination from the last row.
func generateCursor(lastRow map[string]interface{}, uniqueAttrID string, uniqueKeyColumns []string) string {
	if uniqueAttrID == SyntheticIDAttr && len(uniqueKeyColumns) > 0 {
		// For composite keys, use the composite ID as cursor
		if compositeID, exists := lastRow[CompositeIDKey]; exists {
			return fmt.Sprintf("%v", compositeID)
		}
	} else if uniqueAttrID != "" {
		// For single column keys
		if value, exists := lastRow[uniqueAttrID]; exists {
			return fmt.Sprintf("%v", value)
		}
	}

	// Fallback: if no unique key is available, return empty cursor
	return ""
}

// NewClient returns a Client to query the datasource.
func NewClient(client SQLClient) Client {
	return &Datasource{
		Client: client,
	}
}

// GetPage queries a page of objects from a DB2 datasource.
func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	logger := zaplogger.FromContext(ctx)

	// Validate request fields and SQL identifiers
	if err := request.Validate(); err != nil {
		return nil, err
	}

	// Build connection string (includes SSL setup if configured)
	connString, err := request.BuildConnectionString()
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Error building DB2 connection string: %v.", err),
		}
	}

	// Establish database connection
	_, err = d.Client.Connect(connString)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Error connecting to DB2 database: %v.", err),
		}
	}

	// Get unique key columns for composite ID generation
	uniqueKeyColumns, err := d.ExtractUniqueKeyColumns(ctx, request.EntityConfig.ExternalId)
	if err != nil {
		// When using the synthetic "id" attribute, composite key columns are required
		// for unique ID generation and cursor-based pagination.
		if request.UniqueAttributeExternalID == SyntheticIDAttr {
			return nil, &framework.Error{
				Message: fmt.Sprintf(
					"Error extracting unique key columns for table %s: %v. "+
						"Composite key columns are required when using synthetic 'id' attribute.",
					request.EntityConfig.ExternalId, err),
			}
		}

		// For non-synthetic unique attributes, log and continue without composite ID
		logger.Warn("Could not extract unique key for table",
			zap.String("table", request.EntityConfig.ExternalId),
			zap.Error(err))
	}

	// Pass unique key columns to query construction so they can be included in SELECT
	request.UniqueKeyColumns = uniqueKeyColumns

	// Construct the DB2 query
	query, args, err := ConstructQuery(request)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Error constructing DB2 query: %v.", err),
		}
	}

	// Execute the query
	rows, err := d.Client.Query(ctx, query, args...)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Error executing DB2 query: %v.", err),
		}
	}
	defer rows.Close()

	// Process query results using the row processor
	processor := newQueryResultProcessor(request.EntityConfig.Attributes, uniqueKeyColumns, logger)

	results, procErr := processor.process(rows)
	if procErr != nil {
		return nil, procErr
	}

	// Handle pagination - check if we got more rows than requested
	objects := results.objects

	var nextCursor *string

	if request.PageSize > 0 && len(objects) > int(request.PageSize) {
		// Remove the extra row we fetched
		objects = objects[:request.PageSize]

		// Generate cursor from the last row
		if len(objects) > 0 {
			lastRow := objects[len(objects)-1]
			cursor := generateCursor(lastRow, request.UniqueAttributeExternalID, uniqueKeyColumns)
			nextCursor = &cursor
		}
	}

	return &Response{
		Objects:    objects,
		NextCursor: nextCursor,
		TotalCount: results.totalCount,
		StatusCode: 200,
	}, nil
}
