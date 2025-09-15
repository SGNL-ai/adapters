// Copyright 2025 SGNL.ai, Inc.

package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
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
	"github.com/sgnl-ai/adapters/pkg/testutil"
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

	// Act
	resp, err := ds.GetPage(context.Background(), &Request{
		EntityConfig: framework.EntityConfig{
			ExternalId: "users",
		},
		UniqueAttributeExternalID: "id",
		PageSize:                  100,
	})

	// Assert
	if err != nil {
		t.Errorf("Error when requesting GetPage() for a datasource, %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %v, got %v http status code", http.StatusOK, resp.StatusCode)
	}
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
