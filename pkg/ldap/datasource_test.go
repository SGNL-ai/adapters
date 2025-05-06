// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package ldap_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/pkg/connector"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	ldap "github.com/sgnl-ai/adapters/pkg/ldap"
	"github.com/sgnl-ai/adapters/pkg/testutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"crypto/tls"
)

func TestBytesToOctetString(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected string
	}{
		{[]byte("hello"), "aGVsbG8="},
		{[]byte("world"), "d29ybGQ="},
		{[]byte(""), ""},
	}

	for _, tc := range testCases {
		actual := ldap.BytesToOctetString(tc.input)
		expected := base64.StdEncoding.EncodeToString(tc.input)

		if *actual != expected {
			t.Errorf("Expected %s, but got %s", expected, *actual)
		}
	}
}

func TestOctetStringToBytes(t *testing.T) {
	testCases := []struct {
		input       string
		expected    []byte
		expectError bool
	}{
		{"aGVsbG8=", []byte("hello"), false},
		{"d29ybGQ=", []byte("world"), false},
		{"", []byte(""), false},
		{"invalidBase64", nil, true},
	}

	for _, tc := range testCases {
		actual, err := ldap.OctetStringToBytes(tc.input)
		if tc.expectError {
			if err == nil {
				t.Errorf("Expected error, but got nil")
			}
		} else {
			if err != nil {
				t.Errorf("Expected no error, but got %v", err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, actual)
			}
		}
	}
}

type testProxyClient struct {
	ci             *connector.ConnectorInfo
	proxiedRequest bool
}

func (c *testProxyClient) ProxyRequest(_ context.Context, ci *connector.ConnectorInfo, _ *ldap.Request) (*ldap.Response, *framework.Error) {
	// Ensure connector info matches the configured value.
	if c.ci != nil {
		if diff := cmp.Diff(c.ci, ci); diff != "" {
			return nil, &framework.Error{
				Message: fmt.Sprintf("connector info mismatch (-want +got):%s", diff),
			}
		}
	}

	return &ldap.Response{
		StatusCode: http.StatusOK,
	}, nil
}

func (c *testProxyClient) Request(_ context.Context, _ *ldap.Request) (*ldap.Response, *framework.Error) {
	if c.proxiedRequest {
		return nil, &framework.Error{}
	}

	return &ldap.Response{
		StatusCode: http.StatusOK,
	}, nil
}

func (c *testProxyClient) IsProxied() bool {
	return c.proxiedRequest
}

var (
	testRequest = &ldap.Request{
		BaseURL:          mockLDAPAddr,
		PageSize:         1,
		EntityExternalID: "Person",
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "dn",
				Type:       framework.AttributeTypeString,
				List:       false,
				UniqueId:   true,
			},
			{
				ExternalId: "objectGUID",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
		ConnectionParams: ldap.ConnectionParams{
			BindDN:       "cn=user,dc=example,dc=org",
			BindPassword: "asdasd",
			BaseDN:       "dc=example,dc=org",
		},
		EntityConfigMap: map[string]*ldap.EntityConfig{
			"Person": {
				Query: "(&(objectClass=person))",
			},
		},
	}
)

func TestGivenRequestWithoutConnectorContextWhenGetPageRequestedThenLdapResponseStatusIsOk(t *testing.T) {
	// Arrange
	ds := ldap.Datasource{
		Client: &testProxyClient{
			proxiedRequest: false,
		},
	}

	// Act
	resp, err := ds.GetPage(context.Background(), testRequest)

	// Assert
	if err != nil {
		t.Errorf("Error when requesting GetPage() for a datasource, %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %v, got %v http status code", http.StatusOK, resp.StatusCode)
	}
}

func TestGivenRequestWithConnectorAndWithoutProxyContextWhenGetPageRequestedThenLdapResponseStatusIsOk(t *testing.T) {
	// Arrange
	ci := connector.ConnectorInfo{
		ID:       "test-connector-id",
		ClientID: "test-client-id",
		TenantID: "test-tenant-id",
	}

	ds := ldap.Datasource{
		Client: &testProxyClient{
			ci:             &ci,
			proxiedRequest: false,
		},
	}

	ctx, _ := connector.WithContext(context.Background(), ci)

	// Act
	resp, err := ds.GetPage(ctx, testRequest)

	// Assert
	if err != nil {
		t.Errorf("Error when requesting GetPage() for a datasource, %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %v, got %v http status code", http.StatusOK, resp.StatusCode)
	}
}

func TestGivenRequestWithConnectorContextWhenGetPageRequestedThenLdapResponseStatusIsOk(t *testing.T) {
	// Arrange
	ci := connector.ConnectorInfo{
		ID:       "test-connector-id",
		ClientID: "test-client-id",
		TenantID: "test-tenant-id",
	}

	ds := ldap.Datasource{
		Client: &testProxyClient{
			ci:             &ci,
			proxiedRequest: true,
		},
	}

	ctx, _ := connector.WithContext(context.Background(), ci)

	// Act
	resp, err := ds.GetPage(ctx, testRequest)

	// Assert
	if err != nil {
		t.Errorf("Error when requesting GetPage() for a datasource, %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %v, got %v http status code", http.StatusOK, resp.StatusCode)
	}
}

func TestGivenRequestWithConnectorContextWhenProxyServiceConnectionFailsWithGrpcErrThenGetPageReturnsResponseWithCorrectHttpErrCode(t *testing.T) {
	// Arrange
	client, cleanup := testutil.ProxyTestCommonSetup(t, &testutil.TestProxyServer{
		Ci:             &testutil.TestConnectorInfo,
		GrpcErr:        status.Errorf(codes.Unavailable, "aborted request"),
		IsLDAPResponse: true,
	})
	defer cleanup()

	ds := ldap.NewClient(client)

	ctx, _ := connector.WithContext(context.Background(), testutil.TestConnectorInfo)

	// Act
	resp, ferr := ds.GetPage(ctx, testRequest)
	if ferr != nil {
		t.Errorf("expecting error code in response, %v", ferr)
	}

	// Assert
	if customerror.GRPCStatusCodeToHTTP[codes.Unavailable] != resp.StatusCode {
		t.Errorf("failed to match the error code, expected %v, got %v",
			customerror.GRPCStatusCodeToHTTP[codes.Unavailable], resp.StatusCode)
	}
}

func TestGivenRequestWithConnectorContextWhenProxyServiceReturnEmptyResponseThenGetPageReturnsInternalError(t *testing.T) {
	// Arrange
	emptyResponse := ""

	client, cleanup := testutil.ProxyTestCommonSetup(t, &testutil.TestProxyServer{
		Ci:             &testutil.TestConnectorInfo,
		Response:       &emptyResponse,
		IsLDAPResponse: true,
	})
	defer cleanup()

	ds := ldap.NewClient(client)

	ctx, _ := connector.WithContext(context.Background(), testutil.TestConnectorInfo)

	// Act
	resp, ferr := ds.GetPage(ctx, testRequest)
	if resp != nil {
		t.Errorf("expecting nil response %v", resp.StatusCode)
	}

	// Assert
	if api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL != ferr.Code {
		t.Errorf("failed to match the error code, expected %v, got %v",
			api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL, ferr.Code)
	}
}

func TestGivenRequestWithConnectorContextWhenProxyServiceReturnValidResponseThenGetPageReturnsHttpOkWithCorrectResponse(t *testing.T) {
	// Arrange
	testResponse := &ldap.Response{
		StatusCode: http.StatusOK,
		Objects:    []map[string]any{{"a": "b"}, {"c": "d"}},
	}
	data, _ := json.Marshal(testResponse)
	respData := string(data)

	client, cleanup := testutil.ProxyTestCommonSetup(t, &testutil.TestProxyServer{
		Ci:             &testutil.TestConnectorInfo,
		Response:       &respData,
		IsLDAPResponse: true,
	})
	defer cleanup()

	ds := ldap.NewClient(client)

	ctx, _ := connector.WithContext(context.Background(), testutil.TestConnectorInfo)

	// Act
	resp, ferr := ds.GetPage(ctx, testRequest)
	if ferr != nil {
		t.Errorf("expecting nil err %v", ferr)
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

func TestGetTLSConfig(t *testing.T) {
	tests := []struct {
		name           string
		request        *ldap.Request
		expectedError  *framework.Error
		expectedConfig *tls.Config
	}{
		{
			name: "non_ldaps_connection",
			request: &ldap.Request{
				IsLDAPS: false,
			},
			expectedConfig: &tls.Config{},
		},
		{
			name: "invalid_certificate_chain",
			request: &ldap.Request{
				IsLDAPS:          true,
				CertificateChain: "invalid-base64",
			},
			expectedError: &framework.Error{
				Message: "Failed to load certificates: illegal base64 data at input byte 7.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		{
			name: "valid_certificate_chain_with_hostname",
			request: &ldap.Request{
				IsLDAPS:          true,
				Host:             "ldap.example.com",
				CertificateChain: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNzRENDQWhtZ0F3SUJBZ0lKQUx3enJKRUlCT2FlTUEwR0NTcUdTSWIzRFFFQkJRVUFNRVV4Q3pBSkJnTlYKQkFZVEFrRlZNUk13RVFZRFZRUUlFd3BUYjIxbExWTjBZWFJsTVNFd0h3WURWUVFLRXhoSmJuUmxjbTVsZENCWAphV1JuYVhSeklGQjBlU0JNZEdRd0hoY05NVEV3T1RNd01UVXlOak0yV2hjTk1qRXdPVEkzTVRVeU5qTTJXakJGCk1Rc3dDUVlEVlFRR0V3SkJWVEVUTUJFR0ExVUVDQk1LVTI5dFpTMVRkR0YwWlRFaE1COEdBMVVFQ2hNWVNXNTAKWlhKdVpYUWdWMmxrWjJsMGN5QlFkSGtnVEhSa01JR2ZNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0R05BRENCaVFLQgpnUUM4OENrd3J1OVZSMnAyS0oxV1F5cWVzTHpyOTV0YU5iaGtZZnNkMGo4VGwwTUdZNWgrZGN6Q2FNUXowWVkzCnhIWHVVNXlBUVFUWmppa3MrRDNLQTNjeCtpS0RmMnAxcTc3b1h4UWN4NUNrclhCV1R hWDJvcVZ0SG0zYVgyM0IKQUlPUkd1UGswMGI0clQzY2xkN1ZoY0VGbXpSTmJ5STBFcUxNQXhJd2NlVUtTUUlEQVFBQm80R25NSUdrTUIwRwpBMVVkRGdRV0JCU0dtT2R2U1hLWGNsaWM1VU9LUFczNUpMTUVFakIxQmdOVkhTTUViakJzZ0JTR21PZHZTWEtYCmNsaWM1VU9LUFczNUpMTUVFcUZKcEVjd1JURUxNQWtHQTFVRUJoTUNRVlV4RXpBUkJnTlZCQWdUQ2xOdmJXVXQKVTNSaGRHVXhJVEFmQmdOVkJBb1RHRWx1ZEdWeWJtVjBJRmRwWkdkcGRITWdVSFI1SUV4MFpJSUpBTHd6ckpFSQpCT2FlTUF3R0ExVWRFd1FGTUFNQkFmOHdEUVlKS29aSWh2Y05BUUVGQlFBRGdZRUFjUGZXbjQ5cGdBWDU0amk1ClNpVVBGRk5DdVFHU1NUSGgySStUTXJzMUcxTWIzYTBYMWRWNUNOTFJ5WHl1VnhzcWhpTS9IMnZlRm5UejJRNFUKd2RZL2tQeEUxOUF1d2N6OUF2Q2t3N29sMUxJbExmSnZCemp6T2pFcFpKTnRrWFR4OFJPU29vTnJEZUpsM0h5TgpjY2lTNWhmODBYeklGcXdoemFWUzlnbWl5TTg9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0=",
			},
			expectedConfig: &tls.Config{
				ServerName: "ldap.example.com",
			},
		},
		{
			name: "valid_certificate_chain_with_hostname_and_port",
			request: &ldap.Request{
				IsLDAPS:          true,
				Host:             "ldap.example.com:636",
				CertificateChain: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNzRENDQWhtZ0F3SUJBZ0lKQUx3enJKRUlCT2FlTUEwR0NTcUdTSWIzRFFFQkJRVUFNRVV4Q3pBSkJnTlYKQkFZVEFrRlZNUk13RVFZRFZRUUlFd3BUYjIxbExWTjBZWFJsTVNFd0h3WURWUVFLRXhoSmJuUmxjbTVsZENCWAphV1JuYVhSeklGQjBlU0JNZEdRd0hoY05NVEV3T1RNd01UVXlOak0yV2hjTk1qRXdPVEkzTVRVeU5qTTJXakJGCk1Rc3dDUVlEVlFRR0V3SkJWVEVUTUJFR0ExVUVDQk1LVTI5dFpTMVRkR0YwWlRFaE1COEdBMVVFQ2hNWVNXNTAKWlhKdVpYUWdWMmxrWjJsMGN5QlFkSGtnVEhSa01JR2ZNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0R05BRENCaVFLQgpnUUM4OENrd3J1OVZSMnAyS0oxV1F5cWVzTHpyOTV0YU5iaGtZZnNkMGo4VGwwTUdZNWgrZGN6Q2FNUXowWVkzCnhIWHVVNXlBUVFUWmppa3MrRDNLQTNjeCtpS0RmMnAxcTc3b1h4UWN4NUNrclhCV1R hWDJvcVZ0SG0zYVgyM0IKQUlPUkd1UGswMGI0clQzY2xkN1ZoY0VGbXpSTmJ5STBFcUxNQXhJd2NlVUtTUUlEQVFBQm80R25NSUdrTUIwRwpBMVVkRGdRV0JCU0dtT2R2U1hLWGNsaWM1VU9LUFczNUpMTUVFakIxQmdOVkhTTUViakJzZ0JTR21PZHZTWEtYCmNsaWM1VU9LUFczNUpMTUVFcUZKcEVjd1JURUxNQWtHQTFVRUJoTUNRVlV4RXpBUkJnTlZCQWdUQ2xOdmJXVXQKVTNSaGRHVXhJVEFmQmdOVkJBb1RHRWx1ZEdWeWJtVjBJRmRwWkdkcGRITWdVSFI1SUV4MFpJSUpBTHd6ckpFSQpCT2FlTUF3R0ExVWRFd1FGTUFNQkFmOHdEUVlKS29aSWh2Y05BUUVGQlFBRGdZRUFjUGZXbjQ5cGdBWDU0amk1ClNpVVBGRk5DdVFHU1NUSGgySStUTXJzMUcxTWIzYTBYMWRWNUNOTFJ5WHl1VnhzcWhpTS9IMnZlRm5UejJRNFUKd2RZL2tQeEUxOUF1d2N6OUF2Q2t3N29sMUxJbExmSnZCemp6T2pFcFpKTnRrWFR4OFJPU29vTnJEZUpsM0h5TgpjY2lTNWhmODBYeklGcXdoemFWUzlnbWl5TTg9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0=",
			},
			expectedConfig: &tls.Config{
				ServerName: "ldap.example.com:636",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ldap.GetTLSConfig(tt.request)
			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
					return
				}
				if err.Message != tt.expectedError.Message {
					t.Errorf("expected error message %v, got %v", tt.expectedError.Message, err.Message)
				}
				if err.Code != tt.expectedError.Code {
					t.Errorf("expected error code %v, got %v", tt.expectedError.Code, err.Code)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if config.ServerName != tt.expectedConfig.ServerName {
				t.Errorf("expected server name %v, got %v", tt.expectedConfig.ServerName, config.ServerName)
			}
		})
	}
}
