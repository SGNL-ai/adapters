// Copyright 2025 SGNL.ai, Inc.

package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/pkg/connector"
	grpc_proxy_v1 "github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockSQLClient struct {
	connectErr   error
	queryErr     error
	proxyErr     error
	proxyRespErr string
	proxyResp    string
	mockConnect  func(string) error
	mockQuery    func(string, ...any) (*sql.Rows, error)
	mockProxy    func(context.Context, *grpc_proxy_v1.ProxyRequestMessage,
	) (*grpc_proxy_v1.Response, error)
}

func (c *mockSQLClient) IsProxied() bool {
	if c.mockProxy != nil || c.proxyResp != "" {
		return true
	}

	return false
}

func (c *mockSQLClient) Connect(name string) error {
	if c.connectErr != nil {
		return c.connectErr
	}

	if c.mockConnect != nil {
		return c.mockConnect(name)
	}

	return nil
}

func (c *mockSQLClient) Query(query string, args ...any) (*sql.Rows, error) {
	if c.mockQuery != nil {
		return c.mockQuery(query, args)
	}

	if c.queryErr != nil {
		return nil, c.queryErr
	}

	return &sql.Rows{}, nil
}

func (c *mockSQLClient) Proxy(ctx context.Context, req *grpc_proxy_v1.ProxyRequestMessage,
) (*grpc_proxy_v1.Response, error) {
	if c.mockProxy != nil {
		return c.mockProxy(ctx, req)
	}

	if c.proxyErr != nil {
		return nil, c.proxyErr
	}

	return &grpc_proxy_v1.Response{
		ResponseType: &grpc_proxy_v1.Response_SqlQueryResponse{
			SqlQueryResponse: &grpc_proxy_v1.SQLQueryResponse{
				Response: c.proxyResp,
				Error:    c.proxyRespErr,
			},
		},
	}, nil
}

var (
	sqlColumns = []*sqlmock.Column{
		sqlmock.NewColumn("id").OfType("VARCHAR", ""),
		sqlmock.NewColumn("name").OfType("VARCHAR", ""),
	}
	sqlRows = sqlmock.NewRowsWithColumnDefinition(sqlColumns...).
		AddRow("1", "test-name-1").
		AddRow("2", "test-name-2")
)

func TestGivenRequestWithoutConnectorCtxWhenGetPageRequestedThenSQLResponseStatusIsOk(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	mockQuery := func(query string, _ ...any) (*sql.Rows, error) {
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT *, CAST(`id` AS CHAR(50)) AS `str_id` FROM `users` ORDER BY `str_id` ASC LIMIT ?"),
		).WillReturnRows(sqlRows)

		return db.Query(query)
	}
	ds := Datasource{
		Client: &mockSQLClient{
			mockQuery: mockQuery,
		},
	}

	request := &Request{
		EntityConfig: framework.EntityConfig{
			ExternalId: "users",
			Attributes: []*framework.AttributeConfig{
				{
					ExternalId: "id",
					Type:       framework.AttributeTypeString,
				},
				{
					ExternalId: "name",
					Type:       framework.AttributeTypeString,
				},
			},
		},
		UniqueAttributeExternalID: "id",
		PageSize:                  100,
		Username:                  "testuser",
		Password:                  "testpass",
		BaseURL:                   "localhost:3306",
		Database:                  "testdb",
	}

	expectedLogs := []map[string]any{
		{
			"level":                             "info",
			"msg":                               "Starting datasource request",
			fields.FieldRequestEntityExternalID: "users",
			fields.FieldRequestPageSize:         int64(100),
			fields.FieldBaseURL:                 "localhost:3306",
			fields.FieldDatabase:                "testdb",
		},
		{
			"level":                             "info",
			"msg":                               "Sending request to datasource",
			fields.FieldRequestEntityExternalID: "users",
			fields.FieldRequestPageSize:         int64(100),
			fields.FieldBaseURL:                 "localhost:3306",
			fields.FieldDatabase:                "testdb",
		},
		{
			"level":                             "info",
			"msg":                               "Datasource request completed successfully",
			fields.FieldRequestEntityExternalID: "users",
			fields.FieldRequestPageSize:         int64(100),
			fields.FieldResponseStatusCode:      int64(200),
			fields.FieldResponseObjectCount:     int64(2),
			fields.FieldResponseNextCursor:      nil,
			fields.FieldBaseURL:                 "localhost:3306",
			fields.FieldDatabase:                "testdb",
		},
	}

	// Act
	ctxWithLogger, observedLogs := testutil.NewContextWithObservableLogger(t.Context())

	resp, err := ds.GetPage(ctxWithLogger, request)

	// Assert
	if err != nil {
		t.Errorf("Error when requesting GetPage() for a datasource, %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %v, got %v http status code", http.StatusOK, resp.StatusCode)
	}

	testutil.ValidateLogOutput(t, observedLogs, expectedLogs)
}

func TestGivenRequestWithConnectorCtxAndWithoutProxyWhenGetPageRequestedThenSQLResponseStatusIsOk(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	mockQuery := func(query string, _ ...any) (*sql.Rows, error) {
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT *, CAST(`id` AS CHAR(50)) AS `str_id` FROM `users` ORDER BY `str_id` ASC"),
		).WillReturnRows(sqlRows)

		return db.Query(query)
	}
	ds := Datasource{
		Client: &mockSQLClient{
			mockQuery: mockQuery,
		},
	}

	ctx, _ := connector.WithContext(context.Background(), connector.ConnectorInfo{})

	// Act
	resp, err := ds.GetPage(ctx, &Request{
		EntityConfig: framework.EntityConfig{
			ExternalId: "users",
		},
		UniqueAttributeExternalID: "id",
	})

	// Assert
	if err != nil {
		t.Errorf("Error when requesting GetPage() for a datasource, %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %v, got %v http status code", http.StatusOK, resp.StatusCode)
	}
}

func TestGivenRequestWithConnectorCtxWhenGetPageRequestedThenSQLResponseStatusIsOk(t *testing.T) {
	// Arrange
	ds := Datasource{
		Client: &mockSQLClient{
			proxyResp: `{"statusCode": 200}`,
		},
	}

	ctx, _ := connector.WithContext(context.Background(), connector.ConnectorInfo{})

	// Act
	resp, err := ds.GetPage(ctx, &Request{})

	// Assert
	if err != nil {
		t.Errorf("Error when requesting GetPage() for a datasource, %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %v, got %v http status code", http.StatusOK, resp.StatusCode)
	}
}

func TestGivenRequestWithConnectorCtxOnGrpcErrorThenGetPageReturnsResponseWithCorrectHttpCode(t *testing.T) {
	// Arrange
	client, cleanup := testutil.ProxyTestCommonSetup(t, &testutil.TestProxyServer{
		Ci:      &testutil.TestConnectorInfo,
		GrpcErr: status.Errorf(codes.Unavailable, "aborted request"),
	})
	defer cleanup()

	ds := Datasource{
		Client: &defaultSQLClient{
			proxy: client,
		},
	}

	// Act
	ctx, _ := connector.WithContext(context.Background(), testutil.TestConnectorInfo)

	resp, err := ds.GetPage(ctx, &Request{})
	if err != nil {
		t.Errorf("Error when requesting GetPage() for a datasource, %v", err)
	}

	// Assert
	if customerror.GRPCStatusCodeToHTTP[codes.Unavailable] != resp.StatusCode {
		t.Errorf("failed to match the error code, expected %v, got %v",
			customerror.GRPCStatusCodeToHTTP[codes.Unavailable], resp.StatusCode)
	}
}

// nolint // verbose test function name
func TestGivenRequestWithConnectorContextWhenProxyServiceReturnEmptyResponseThenGetPageReturnsInternalError(t *testing.T) {
	// Arrange
	emptyResponse := ""

	client, cleanup := testutil.ProxyTestCommonSetup(t, &testutil.TestProxyServer{
		Ci:            &testutil.TestConnectorInfo,
		Response:      &emptyResponse,
		IsSQLResponse: true,
	})
	defer cleanup()

	ds := Datasource{
		Client: &defaultSQLClient{
			proxy: client,
		},
	}

	// Act
	ctx, _ := connector.WithContext(context.Background(), testutil.TestConnectorInfo)

	resp, err := ds.GetPage(ctx, &Request{})
	if resp != nil {
		t.Errorf("expecting nil response %v", resp.StatusCode)
	}

	// Assert
	if api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL != err.Code {
		t.Errorf("failed to match the error code, expected %v, got %v",
			api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL, err.Code)
	}
}

// nolint // verbose function name
func TestGivenRequestWithConnectorContextWhenProxyServiceReturnValidResponseThenGetPageReturnsHttpOkWithCorrectResponse(t *testing.T) {
	// Arrange
	testResponse := &Response{
		StatusCode: http.StatusOK,
		Objects:    []map[string]any{{"a": "b"}, {"c": "d"}},
	}
	data, _ := json.Marshal(testResponse)
	respData := string(data)

	client, cleanup := testutil.ProxyTestCommonSetup(t, &testutil.TestProxyServer{
		Ci:            &testutil.TestConnectorInfo,
		Response:      &respData,
		IsSQLResponse: true,
	})
	defer cleanup()

	ds := Datasource{
		Client: &defaultSQLClient{
			proxy: client,
		},
	}

	// Act
	ctx, _ := connector.WithContext(context.Background(), testutil.TestConnectorInfo)

	resp, err := ds.GetPage(ctx, &Request{})
	if err != nil {
		t.Errorf("expecting nil err %v", err)
	}

	// Assert
	if http.StatusOK != resp.StatusCode {
		t.Errorf("failed to match the error code, expected %v, got %v",
			http.StatusOK, resp.StatusCode)
	}

	if diff := cmp.Diff(testResponse, resp); diff != "" {
		t.Errorf("response payload mismatch (-want +got):%s", diff)
	}
}

// TestValidateProxyResponse tests the pure validateProxyResponse function.
func TestValidateProxyResponse(t *testing.T) {
	testCases := []struct {
		name              string
		proxyResp         *grpc_proxy_v1.Response
		err               error
		expectError       bool
		expectMessagePart string
		expectErrorCode   *api_adapter_v1.ErrorCode
		expectStatusCode  *int // For gRPC status code cases
	}{
		{
			name:        "success_valid_response",
			proxyResp:   createValidProxyResponse(t, 10),
			err:         nil,
			expectError: false,
		},
		{
			name:        "grpc_error_unavailable",
			proxyResp:   nil,
			err:         status.Error(codes.Unavailable, "service unavailable"),
			expectError: false, // Returns response with status code, not framework error
			expectStatusCode: func() *int {
				s := 503

				return &s
			}(),
		},
		{
			name:              "grpc_error_non_status_error",
			proxyResp:         nil,
			err:               errors.New("connection refused"),
			expectError:       true,
			expectMessagePart: "Error querying SQL server",
		},
		{
			name:              "nil_proxy_resp",
			proxyResp:         nil,
			err:               nil,
			expectError:       true,
			expectMessagePart: "nil response",
		},
		{
			name: "response_error_field_set",
			proxyResp: &grpc_proxy_v1.Response{
				Error: "Database connection timeout",
			},
			err:               nil,
			expectError:       true,
			expectMessagePart: "Error received from proxy",
		},
		{
			name: "nil_sql_query_response",
			proxyResp: &grpc_proxy_v1.Response{
				ResponseType: nil,
			},
			err:               nil,
			expectError:       true,
			expectMessagePart: "nil SqlQueryResponse",
		},
		{
			name: "sql_query_response_error_field_with_valid_framework_error",
			proxyResp: createProxyResponseWithError(t, &framework.Error{
				Message: "Table does not exist",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			}),
			err:               nil,
			expectError:       true,
			expectMessagePart: "Table does not exist",
			expectErrorCode: func() *api_adapter_v1.ErrorCode {
				c := api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG

				return &c
			}(),
		},
		{
			name: "empty_response_string",
			proxyResp: &grpc_proxy_v1.Response{
				ResponseType: &grpc_proxy_v1.Response_SqlQueryResponse{
					SqlQueryResponse: &grpc_proxy_v1.SQLQueryResponse{
						Response: "",
					},
				},
			},
			err:               nil,
			expectError:       true,
			expectMessagePart: "empty response",
		},
		{
			name: "invalid_json_in_response",
			proxyResp: &grpc_proxy_v1.Response{
				ResponseType: &grpc_proxy_v1.Response_SqlQueryResponse{
					SqlQueryResponse: &grpc_proxy_v1.SQLQueryResponse{
						Response: `{"statusCode":200,"objects":[{"id":"1"`, // Truncated JSON
					},
				},
			},
			err:               nil,
			expectError:       true,
			expectMessagePart: "unexpected end of JSON input",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, frameworkErr := validateProxyResponse(tc.proxyResp, tc.err)

			if tc.expectError {
				assert.Nil(t, response, "Expected nil response")
				assert.NotNil(t, frameworkErr, "Expected framework error")

				if frameworkErr != nil {
					t.Logf("   Error message: %s", frameworkErr.Message)

					assert.Contains(t, frameworkErr.Message, tc.expectMessagePart,
						"Error message should contain expected text")

					if tc.expectErrorCode != nil {
						assert.Equal(t, *tc.expectErrorCode, frameworkErr.Code,
							"Error code should match")
						t.Logf("   Error code: %v", frameworkErr.Code)
					}
				}
			} else if tc.expectStatusCode != nil {
				// gRPC status code case - returns response with status code, not error
				assert.NotNil(t, response, "Expected response with status code")
				assert.Nil(t, frameworkErr, "Expected no framework error")
				assert.Equal(t, *tc.expectStatusCode, response.StatusCode,
					"Status code should match expected value")
				t.Logf("   Status code: %d", response.StatusCode)
			} else {
				// Success case
				assert.NotNil(t, response, "Expected valid response")
				assert.Nil(t, frameworkErr, "Expected no error")
				assert.Greater(t, len(response.Objects), 0, "Expected objects in response")
				t.Logf("Success: Got %d objects", len(response.Objects))
			}
		})
	}
}

// Helper function to create a valid proxy response with N objects.
func createValidProxyResponse(t *testing.T, numObjects int) *grpc_proxy_v1.Response {
	resp := Response{
		StatusCode: 200,
		Objects:    make([]map[string]any, numObjects),
	}

	for i := 0; i < numObjects; i++ {
		resp.Objects[i] = map[string]any{
			"id":   i + 1,
			"name": "test",
		}
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	return &grpc_proxy_v1.Response{
		ResponseType: &grpc_proxy_v1.Response_SqlQueryResponse{
			SqlQueryResponse: &grpc_proxy_v1.SQLQueryResponse{
				Response: string(data),
			},
		},
	}
}

// Helper function to create proxy response with SQLQueryResponse.Error field.
func createProxyResponseWithError(t *testing.T, frameworkErr *framework.Error) *grpc_proxy_v1.Response {
	errorJSON, err := json.Marshal(frameworkErr)
	require.NoError(t, err)

	return &grpc_proxy_v1.Response{
		ResponseType: &grpc_proxy_v1.Response_SqlQueryResponse{
			SqlQueryResponse: &grpc_proxy_v1.SQLQueryResponse{
				Error: string(errorJSON),
			},
		},
	}
}
