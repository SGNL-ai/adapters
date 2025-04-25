// Copyright 2025 SGNL.ai, Inc.

package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/pkg/connector"
	grpc_proxy_v1 "github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"google.golang.org/grpc/status"

	_ "github.com/go-sql-driver/mysql" // Go MySQL Driver is an implementation of Go's database/sql/driver interface.
)

var validSQLTableName = regexp.MustCompile(`^[a-zA-Z0-9$_]*$`)

type Datasource struct {
	Client SQLClient
}

// NewClient returns a Client to query the datasource.
func NewClient(client SQLClient) Client {
	return &Datasource{
		Client: client,
	}
}

// ProxyRequest sends serialized SQL query request to the on-premises connector.
func (d *Datasource) ProxyRequest(ctx context.Context, request *Request,
	ci *connector.ConnectorInfo) (*Response, *framework.Error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to proxy sql request, %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	proxyRequest := &grpc_proxy_v1.ProxyRequestMessage{
		ConnectorId: ci.ID,
		ClientId:    ci.ClientID,
		TenantId:    ci.TenantID,
		Request: &grpc_proxy_v1.Request{
			RequestType: &grpc_proxy_v1.Request_SqlQueryReq{
				SqlQueryReq: &grpc_proxy_v1.SQLQueryRequest{
					Request: string(data),
				},
			},
		},
	}

	response := &Response{}

	proxyResp, err := d.Client.Proxy(ctx, proxyRequest)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			code := customerror.GRPCStatusCodeToHTTPStatusCode(st, err)
			response.StatusCode = code

			return response, nil
		}

		return nil, &framework.Error{
			Message: fmt.Sprintf("Error querying SQL server: - %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	resp := proxyResp.GetSqlQueryResponse()
	if resp == nil {
		return nil, &framework.Error{
			Message: "Error received nil response from the proxy",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if err = json.Unmarshal([]byte(resp.Response), response); err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Error unmarshalling SQL response from the proxy: %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return response, nil
}

// Request function to directly connect to the SQL datasource and execute a query
// to fetch data.
func (d *Datasource) Request(_ context.Context, request *Request,
) (*Response, *framework.Error) {
	if err := request.SimpleSQLValidation(); err != nil {
		return nil, err
	}

	if err := d.Client.Connect(request.DatasourceName()); err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to connect to datasource: %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	var cursor int64

	if request.Cursor != nil {
		cursor = *request.Cursor
	}

	query := fmt.Sprintf(
		"SELECT *, CAST(? as CHAR(50)) as str_id FROM %s ORDER BY str_id ASC LIMIT ? OFFSET ?",
		request.EntityExternalID,
	)

	args := []any{
		request.UniqueAttributeExternalID,
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

	// If we have less objects than the current PageSize, this is the last page and
	// we should not set a NextCursor.
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

// GetPage for requesting data from a datasource.
func (d *Datasource) GetPage(ctx context.Context, request *Request,
) (*Response, *framework.Error) {
	ci, ok := connector.FromContext(ctx)
	if ok {
		return d.ProxyRequest(ctx, request, &ci)
	}

	return d.Request(ctx, request)
}

// ParseResponse for parsing the SQL query response.
func ParseResponse(rows *sql.Rows, request *Request,
) ([]map[string]any, *framework.Error) {
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

			// The unique attribute always needs to be cast as a string (due to
			// ingestion scheduler validation).
			if cols[i] == request.UniqueAttributeExternalID {
				objects[idx][cols[i]] = str
			} else {
				// Otherwise, cast each value based on the type.
				switch types[i].DatabaseTypeName() {
				// The adapter framework expects all numbers to be passed as floats,
				// so parse all numeric types as floats here.
				// TODO [sc-42217]: Split out the logic to parse ints into a separate
				// case once we add support for providing ints to the action framework.
				case "DECIMAL", "NUMERIC", "FLOAT", "DOUBLE", "SMALLINT", "MEDIUMINT",
					"INT", "INTEGER", "BIGINT", "UNSIGNED SMALLINT", "UNSIGNED MEDIUMINT",
					"UNSIGNED INT", "UNSIGNED INTEGER", "UNSIGNED BIGINT":
					objects[idx][cols[i]], castErr = strconv.ParseFloat(str, 64)
				case "BIT", "TINYINT", "BOOL", "BOOLEAN":
					objects[idx][cols[i]], castErr = strconv.ParseBool(str)
				// Default to casting any other values
				// (VARCHAR, TEXT, NVARCHAR, DATETIME, etc) as strings.
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
