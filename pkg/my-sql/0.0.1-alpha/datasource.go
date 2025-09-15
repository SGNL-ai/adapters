// Copyright 2025 SGNL.ai, Inc.

package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/pkg/connector"
	grpc_proxy_v1 "github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"google.golang.org/grpc/status"

	_ "github.com/go-sql-driver/mysql" // Go MySQL Driver is an implementation of Go's database/sql/driver interface.
)

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
func (d *Datasource) ProxyRequest(ctx context.Context, request *Request, ci *connector.ConnectorInfo,
) (*Response, *framework.Error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to proxy sql request, %v.", err),
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
			code := customerror.GRPCErrStatusToHTTPStatusCode(st, err)
			response.StatusCode = code

			return response, nil
		}

		return nil, &framework.Error{
			Message: fmt.Sprintf("Error querying SQL server: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	resp := proxyResp.GetSqlQueryResponse()
	if resp == nil {
		return nil, &framework.Error{
			Message: "Error received nil response from the proxy.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if err = json.Unmarshal([]byte(resp.Response), response); err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Error unmarshalling SQL response from the proxy: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return response, nil
}

// Request function to directly connect to the SQL datasource and execute a query
// to fetch data.
func (d *Datasource) Request(_ context.Context, request *Request) (*Response, *framework.Error) {
	if err := request.SimpleSQLValidation(); err != nil {
		return nil, err
	}

	if err := d.Client.Connect(request.DatasourceName()); err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to connect to datasource: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	query, args, err := ConstructQuery(request)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to construct query: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	rows, err := d.Client.Query(query, args...)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to query datasource: %v.", err),
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
func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	// Make sure if the connector context is set and client can proxy the request.
	if d.Client.IsProxied() {
		if ci, ok := connector.FromContext(ctx); ok {
			return d.ProxyRequest(ctx, request, &ci)
		}
	}

	return d.Request(ctx, request)
}

// ParseResponse for parsing the SQL query response.
func ParseResponse(rows *sql.Rows, request *Request) ([]map[string]any, *framework.Error) {
	objects := make([]map[string]any, 0)

	// Get column names present in provided rows.
	cols, err := rows.Columns()
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to retrieve column names: %s.", err.Error()),
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
				Message: fmt.Sprintf("Failed to scan current row: %s.", err.Error()),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
			}
		}

		for i, v := range vals {
			columnName := cols[i]

			var attribute *framework.AttributeConfig

			for _, curAttribute := range request.EntityConfig.Attributes {
				if curAttribute != nil && curAttribute.ExternalId == columnName {
					attribute = curAttribute
				}
			}

			// Skipping current attribute since this wasn't requested.
			if attribute == nil {
				continue
			}

			if len(objects) < idx+1 {
				objects = append(objects, map[string]any{})
			}

			b, ok := v.(*sql.RawBytes)
			if !ok || b == nil {
				return nil, &framework.Error{
					Message: "Failed to cast value to sql.RawBytes.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			// If no data is returned for the current value, skip and don't return a value for this attribute.
			if len(*b) == 0 {
				continue
			}

			str := string(*b)

			var castErr error

			// Attempt to cast each attribute based on the requested type.
			switch attribute.Type {
			case framework.AttributeTypeBool:
				objects[idx][columnName], castErr = strconv.ParseBool(str)
			// The adapter framework expects all numbers to be passed as floats, so parse all
			// numeric types as floats here.
			case framework.AttributeTypeDouble, framework.AttributeTypeInt64:
				objects[idx][columnName], castErr = strconv.ParseFloat(str, 64)
			case framework.AttributeTypeString, framework.AttributeTypeDuration, framework.AttributeTypeDateTime:
				objects[idx][columnName] = str
			default:
				return nil, &framework.Error{
					Message: "Unsupported attribute type provided.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
				}
			}

			if castErr != nil {
				return nil, &framework.Error{
					Message: fmt.Sprintf("Failed to parse attribute: (%s) %v.", columnName, castErr),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
				}
			}
		}

		idx++
	}

	if err := rows.Err(); err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to process rows: %s.", err.Error()),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return objects, nil
}
