// Copyright 2025 SGNL.ai, Inc.
package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"

	_ "github.com/go-sql-driver/mysql" // Go MySQL Driver is an implementation of Go's database/sql/driver interface.
)

var validSQLIdentifier = regexp.MustCompile(`^[a-zA-Z0-9$_]*$`)

type Datasource struct {
	Client SQLClient
}

// NewClient returns a Client to query the datasource.
func NewClient(client SQLClient) Client {
	return &Datasource{
		Client: client,
	}
}

func (d *Datasource) GetPage(_ context.Context, request *Request) (*Response, *framework.Error) {
	datasourceName := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		request.Username, request.Password, request.BaseURL, request.Database,
	)

	if err := d.Client.Connect(datasourceName); err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to connect to datasource: %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	var cursor int64

	if request.Cursor != nil {
		cursor = *request.Cursor
	}

	// Perform simple validation on the table name to prevent SQL Ingestion attacks, since we can't use table
	// names in prepared queries which leaves us vulnerable.
	if valid := validSQLIdentifier.MatchString(request.EntityExternalID); !valid {
		return nil, &framework.Error{
			Message: "SQL table name validation failed: unsupported characters found.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if valid := validSQLIdentifier.MatchString(request.UniqueAttributeExternalID); !valid {
		return nil, &framework.Error{
			Message: "SQL unique attribute validation failed: unsupported characters found.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	query := fmt.Sprintf(
		"SELECT *, CAST(%s as CHAR(50)) as str_id FROM %s ORDER BY str_id ASC LIMIT ? OFFSET ?",
		request.UniqueAttributeExternalID,
		request.EntityExternalID,
	)

	args := []any{
		request.PageSize,
		cursor,
	}

	rows, err := d.Client.Query(query, args...)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to query datasource: %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	// Parse the rows to a list of objects.
	objs, frameworkErr := ParseResponse(rows, request)
	if frameworkErr != nil {
		return nil, frameworkErr
	}

	response := &Response{
		StatusCode: http.StatusOK,
		Objects:    objs,
	}

	// If we have less objects than the current PageSize, this is the last page and we should not set a NextCursor.
	if len(objs) == int(request.PageSize) {
		if request.Cursor == nil {
			response.NextCursor = &request.PageSize
		} else {
			nextCursor := (*request.Cursor) + request.PageSize

			response.NextCursor = &nextCursor
		}
	}

	return response, nil
}

// nolint: lll
func ParseResponse(rows *sql.Rows, request *Request) ([]map[string]any, *framework.Error) {
	objects := make([]map[string]any, 0)

	// Get column names present in provided rows.
	cols, err := rows.Columns()
	if err != nil {
		return nil, &framework.Error{
			Message: err.Error(),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	// Get column types present in provided rows.
	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, &framework.Error{
			Message: err.Error(),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	// Create an array with the length of cols. This is used when scanning each row.
	vals := make([]any, len(cols))
	for i := range cols {
		vals[i] = new(sql.RawBytes)
	}

	idx := 0

	// Process each row.
	for rows.Next() {
		err := rows.Scan(vals...)
		if err != nil {
			return nil, &framework.Error{
				Message: err.Error(),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
			}
		}

		for i, v := range vals {
			if len(objects) < idx+1 {
				objects = append(objects, map[string]any{})
			}

			b, ok := v.(*sql.RawBytes)
			if !ok || b == nil {
				return nil, &framework.Error{
					Message: "Failed to cast value to sql.RawBytes",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			str := string(*b)

			var castErr error

			// The unique attribute always needs to be cast as a string (due to ingestion scheduler validation).
			if cols[i] == request.UniqueAttributeExternalID {
				objects[idx][cols[i]] = str
			} else {
				// Otherwise, cast each value based on the type.
				switch types[i].DatabaseTypeName() {
				// The adapter framework expects all numbers to be passed as floats, so parse all
				// numeric types as floats here.
				case "SMALLINT", "MEDIUMINT", "INT", "INTEGER", "BIGINT", "UNSIGNED SMALLINT", "UNSIGNED MEDIUMINT", "UNSIGNED INT", "UNSIGNED INTEGER", "UNSIGNED BIGINT":
					if request.CastIntegersToStrings {
						objects[idx][cols[i]] = str
					} else {
						objects[idx][cols[i]], castErr = strconv.ParseFloat(str, 64)
					}
				case "DECIMAL", "NUMERIC", "FLOAT", "DOUBLE":
					objects[idx][cols[i]], castErr = strconv.ParseFloat(str, 64)
				case "BIT", "TINYINT", "BOOL", "BOOLEAN":
					objects[idx][cols[i]], castErr = strconv.ParseBool(str)
				// Default to casting any other values (VARCHAR, TEXT, NVARCHAR, DATETIME, etc) as strings.
				default:
					objects[idx][cols[i]] = str
				}
			}

			if castErr != nil {
				return nil, &framework.Error{
					Message: castErr.Error(),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
				}
			}
		}

		idx++
	}

	if err := rows.Err(); err != nil {
		return nil, &framework.Error{
			Message: err.Error(),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return objects, nil
}
