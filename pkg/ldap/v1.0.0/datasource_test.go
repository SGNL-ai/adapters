// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package ldap_test

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	ldap_v3 "github.com/go-ldap/ldap/v3"
	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/pkg/connector"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	ldap "github.com/sgnl-ai/adapters/pkg/ldap/v1.0.0"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	validCertificateChain = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURlekNDQW1PZ0F3SUJBZ0lVTWFQRkozK3JCNGZlZWQ1UEVvaTJqR1JBaXpBd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1NqRUxNQWtHQTFVRUJoTUNWVk14RmpBVUJnTlZCQWdNRFZOaGJpQkdjbUZ1WTJselkyOHhEakFNQmdOVgpCQWNNQlZCaGJXRnlNUk13RVFZRFZRUUtEQXBUUjA1TUlFbHVZeTRnTUI0WERUSXpNRGt4TlRFek1qTXdNRm9YCkRUSTBNRGt4TkRFek1qTXdNRm93U2pFTE1Ba0dBMVVFQmhNQ1ZWTXhGakFVQmdOVkJBZ01EVk5oYmlCR2NtRnUKWTJselkyOHhEakFNQmdOVkJBY01CVkJoYldGeU1STXdFUVlEVlFRS0RBcFRSMDVNSUVsdVl5NGdNSUlCSWpBTgpCZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEwL0FzNUg3L2ZwZXZXa0VnbGhNTUJwVmxFWFBvCnBxaWZyY2ZrNlM5ZEFnTUVCOVZYY2tBTVBRYjRXNnZzNlY0K0VkK2Q5YzRTWnFkRXBpcGxvbkQyUVNkU0RnN2gKWVJKQlpKYmNqK1JqTUZrd0JxZnNPUzVDWEVDNHZJdm1HTkZNaGRGZUZINHdIRzZKcWozQVdXbGZaT2FpVHBragpJbVRlY0NkdGhQbW9UR1B1WnJDK0VFRjJwYk9GdGxXVWRwU0VZcTB2NEJmS0JFZDkrZnJVTWYzYnk0cVBVUWxvCnRPU0JKK0pKQmFqY1pZVU9zWVdMUkdZWnZFakpOMzNNaGdxaHJWVzF4QmxHYTZwN1BQZVZkUEZvTVdYRVZYTDIKVGxJRDVEcEd2QlpUcUNNRUU3TGZscEtYR0JFUGdBVjJwVGFBUlZLUGRxRVZxUXdJREFRQUJvMU13VVRBZEJnTlYKSFE0RUZnUVVTOHRGS0VQeGpUTEZQZVdxSEZjT1lXeVd6bTh3SHdZRFZSMGpCQmd3Rm9BVVMxVEdLRVB4alRMRgpQZVdxSEZjT1lXeVd6bTh3RHdZRFZSMFRBUUgvQkFVd0F3RUIvekFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBCnZGZEZvQ0I3Y2tNQzRxNGFRY1p6QmFNNzk3VkdBK0NuUWN3QUlSYWNCZkdYZFNNbnEyVWhjam9QZFQ5ZWR4dFEKWUNFWlkwWmVVUGI3UVJGUGlOZkZjTVRDNTNhUWN4THZhNTdSemVhYjE5Y0Y5M2ZlYWpUdkpVZHVPTVRDK2ZVYwpBTzVGTkFOOVZXamZaZVVad0JXVkl6QkFIMGE1bGVLZWZXY2pUK0JUMWlYTVBnUGhTUGpkNHBiT3VmZzZNRGhzCnRFWVBvUWRyVEtYVVIxWnR0cEZEVmZYVVlOQTBrYVdvWVZUdVcxRnNFRGJXQkxuTGxnRVBBVlJzTGRjMGRYbVQKWVBDRGRFUVBzZnJ2NWJkRGtQUVBxTnZEQ2ZrVnIrQkxlTWJxN2VyYVBJUGRBVHdtdkxmQVdxWVBxRGxDRGxLbQpBVGxJUGZJVGRxRVBxZz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
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
	ctxWithLogger, _ := testutil.NewContextWithObservableLogger(t.Context())
	resp, err := ds.GetPage(ctxWithLogger, testRequest)

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

	ctx, _ := connector.WithContext(t.Context(), ci)
	ctxWithLogger, _ := testutil.NewContextWithObservableLogger(ctx)

	// Act
	resp, err := ds.GetPage(ctxWithLogger, testRequest)

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

	ctx, _ := connector.WithContext(t.Context(), ci)
	ctxWithLogger, _ := testutil.NewContextWithObservableLogger(ctx)

	// Act
	resp, err := ds.GetPage(ctxWithLogger, testRequest)

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

	ds := ldap.NewClient(client, ldap.NewSessionPool(1*time.Minute, time.Minute))

	ctx, _ := connector.WithContext(t.Context(), testutil.TestConnectorInfo)
	ctxWithLogger, _ := testutil.NewContextWithObservableLogger(ctx)

	// Act
	resp, ferr := ds.GetPage(ctxWithLogger, testRequest)
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

	ds := ldap.NewClient(client, ldap.NewSessionPool(1*time.Minute, time.Minute))

	ctx, _ := connector.WithContext(t.Context(), testutil.TestConnectorInfo)
	ctxWithLogger, _ := testutil.NewContextWithObservableLogger(ctx)

	// Act
	resp, ferr := ds.GetPage(ctxWithLogger, testRequest)
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

	ds := ldap.NewClient(client, ldap.NewSessionPool(1*time.Minute, time.Minute))

	ctx, _ := connector.WithContext(t.Context(), testutil.TestConnectorInfo)
	ctxWithLogger, observedLogs := testutil.NewContextWithObservableLogger(ctx)

	expectedLogs := []map[string]any{
		{
			"level":                             "info",
			"msg":                               "Sending request to datasource",
			fields.FieldRequestEntityExternalID: "Person",
			fields.FieldRequestPageSize:         int64(1),
			fields.FieldBaseURL:                 mockLDAPAddr,
			fields.FieldConnectorID:             "test-connector-id",
			fields.FieldConnectorSourceID:       "",
			fields.FieldConnectorSourceType:     int64(0),
		},
		{
			"level":                             "info",
			"msg":                               "Datasource request completed successfully",
			fields.FieldRequestEntityExternalID: "Person",
			fields.FieldRequestPageSize:         int64(1),
			fields.FieldResponseStatusCode:      int64(200),
			fields.FieldResponseObjectCount:     int64(2),
			fields.FieldResponseNextCursor:      (*pagination.CompositeCursor[string])(nil),
			fields.FieldBaseURL:                 mockLDAPAddr,
			fields.FieldConnectorID:             "test-connector-id",
			fields.FieldConnectorSourceID:       "",
			fields.FieldConnectorSourceType:     int64(0),
		},
	}

	// Act
	resp, ferr := ds.GetPage(ctxWithLogger, testRequest)
	if ferr != nil {
		t.Errorf("expecting nil err %v", ferr)
	}

	// Assert
	if http.StatusOK != resp.StatusCode {
		t.Errorf("failed to match the status code, expected %v, got %v",
			http.StatusOK, resp.StatusCode)
	}

	if diff := cmp.Diff(testResponse, resp); diff != "" {
		t.Errorf("response payload mismatch (-want +got):%s", diff)
	}

	testutil.ValidateLogOutput(t, observedLogs, expectedLogs)
}

func TestGivenRequestWithConnectorContextWhenProxyServiceReturnValidResponseWithErrorShouldReturnErroredResponse(t *testing.T) {
	// Arrange
	testErrResponse := &framework.Error{
		Message: "error message",
		Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
	}
	data, _ := json.Marshal(testErrResponse)
	respData := string(data)

	client, cleanup := testutil.ProxyTestCommonSetup(t, &testutil.TestProxyServer{
		Ci:             &testutil.TestConnectorInfo,
		ResponseErrStr: &respData,
		IsLDAPResponse: true,
	})
	defer cleanup()

	ds := ldap.NewClient(client, ldap.NewSessionPool(1*time.Minute, time.Minute))

	ctx, _ := connector.WithContext(t.Context(), testutil.TestConnectorInfo)
	ctxWithLogger, _ := testutil.NewContextWithObservableLogger(ctx)

	// Act
	_, ferr := ds.GetPage(ctxWithLogger, testRequest)

	// Assert
	if ferr == nil {
		t.Errorf("expecting nil err %v", ferr)
	}

	if api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG != ferr.Code {
		t.Errorf("failed to match the error code, expected %v, got %v",
			api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG, ferr.Code)
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
				ConnectionParams: ldap.ConnectionParams{
					IsLDAPS: false,
				},
			},
		},
		{
			name: "invalid_certificate_chain",
			request: &ldap.Request{
				ConnectionParams: ldap.ConnectionParams{
					IsLDAPS:          true,
					CertificateChain: "invalid-base64",
				},
			},
			expectedError: &framework.Error{
				Message: "Failed to load certificates: illegal base64 data at input byte 7.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		{
			name: "valid_certificate_chain_with_hostname",
			request: &ldap.Request{
				ConnectionParams: ldap.ConnectionParams{
					IsLDAPS:          true,
					Host:             "ldap.example.com",
					CertificateChain: validCertificateChain,
				},
				BaseURL: "ldaps://ldap.example.com",
			},
			expectedConfig: &tls.Config{
				ServerName: "ldap.example.com",
			},
		},
		{
			name: "valid_certificate_chain_with_hostname_and_port",
			request: &ldap.Request{
				ConnectionParams: ldap.ConnectionParams{
					IsLDAPS:          true,
					Host:             "ldap.example.com:636",
					CertificateChain: validCertificateChain,
				},
				BaseURL: "ldaps://ldap.example.com:636",
			},
			expectedConfig: &tls.Config{
				ServerName: "ldap.example.com",
			},
		},
		{
			name: "valid_certificate_chain_with_ipv6",
			request: &ldap.Request{
				ConnectionParams: ldap.ConnectionParams{
					IsLDAPS:          true,
					Host:             "[2001:db8::1]:636",
					CertificateChain: validCertificateChain,
				},
				BaseURL: "ldaps://[2001:db8::1]:636",
			},
			expectedConfig: &tls.Config{
				ServerName: "2001:db8::1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ldap.GetTLSConfig(tt.request)

			if tt.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.expectedError)
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
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectedConfig != nil {
				if config.ServerName != tt.expectedConfig.ServerName {
					t.Errorf("expected server name %v, got %v", tt.expectedConfig.ServerName, config.ServerName)
				}
			}
		})
	}
}

func TestStringAttrValuesToRequestedType_EmptyValuesHandling(t *testing.T) {
	tests := []struct {
		name        string
		attr        *ldap_v3.EntryAttribute
		isList      bool
		attrType    framework.AttributeType
		expectError bool
		errorMsg    string
	}{
		{
			name: "empty_values_string_type",
			attr: &ldap_v3.EntryAttribute{
				Name:   "testAttr",
				Values: []string{},
			},
			isList:      false,
			attrType:    framework.AttributeTypeString,
			expectError: false,
		},
		{
			name: "empty_values_bool_type",
			attr: &ldap_v3.EntryAttribute{
				Name:   "testBool",
				Values: []string{},
			},
			isList:      false,
			attrType:    framework.AttributeTypeBool,
			expectError: false,
		},
		{
			name: "empty_values_double_type",
			attr: &ldap_v3.EntryAttribute{
				Name:   "testDouble",
				Values: []string{},
			},
			isList:      false,
			attrType:    framework.AttributeTypeDouble,
			expectError: false,
		},
		{
			name: "empty_values_int64_type",
			attr: &ldap_v3.EntryAttribute{
				Name:   "testInt",
				Values: []string{},
			},
			isList:      false,
			attrType:    framework.AttributeTypeInt64,
			expectError: false,
		},
		{
			name: "empty_values_duration_type",
			attr: &ldap_v3.EntryAttribute{
				Name:   "testDuration",
				Values: []string{},
			},
			isList:      false,
			attrType:    framework.AttributeTypeDuration,
			expectError: false,
		},
		{
			name: "empty_values_datetime_type",
			attr: &ldap_v3.EntryAttribute{
				Name:   "testDateTime",
				Values: []string{},
			},
			isList:      false,
			attrType:    framework.AttributeTypeDateTime,
			expectError: false,
		},
		{
			name: "empty_values_list_type_returns_empty_slice",
			attr: &ldap_v3.EntryAttribute{
				Name:   "testList",
				Values: []string{},
			},
			isList:      true,
			attrType:    framework.AttributeTypeString,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ldap.StringAttrValuesToRequestedType(tt.attr, tt.isList, tt.attrType)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")

					return
				}

				if err.Message != tt.errorMsg {
					t.Errorf("expected error message '%s', got '%s'", tt.errorMsg, err.Message)
				}

				if err.Code != api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE {
					t.Errorf("expected error code ERROR_CODE_INVALID_ATTRIBUTE_TYPE, got %v", err.Code)
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)

				return
			}

			// For list type with empty values, should return empty slice
			if tt.isList {
				if result == nil {
					t.Errorf("expected empty slice, got nil")
				}

				return
			}

			// For non-list types with empty values, should return empty string
			if result != "" {
				t.Errorf("expected empty string, got %v", result)
			}
		})
	}
}

func TestStringAttrValuesToRequestedType_EmptyByteValuesHandling(t *testing.T) {
	tests := []struct {
		name        string
		attr        *ldap_v3.EntryAttribute
		expectError bool
		errorMsg    string
	}{
		{
			name: "empty_byte_values_objectGUID",
			attr: &ldap_v3.EntryAttribute{
				Name:       "objectGUID",
				Values:     []string{"dummy"}, // Values present but ByteValues empty
				ByteValues: [][]byte{},
			},
			expectError: true,
			errorMsg:    "Missing or nil GUID bytes for attribute: objectGUID",
		},
		{
			name: "nil_byte_values_objectGUID",
			attr: &ldap_v3.EntryAttribute{
				Name:       "objectGUID",
				Values:     []string{"dummy"},
				ByteValues: [][]byte{nil},
			},
			expectError: true,
			errorMsg:    "Missing or nil GUID bytes for attribute: objectGUID",
		},
		{
			name: "empty_byte_values_objectSid",
			attr: &ldap_v3.EntryAttribute{
				Name:       "objectSid",
				Values:     []string{"dummy"},
				ByteValues: [][]byte{},
			},
			expectError: true,
			errorMsg:    "Missing or nil SID bytes for attribute: objectSid",
		},
		{
			name: "nil_byte_values_objectSid",
			attr: &ldap_v3.EntryAttribute{
				Name:       "objectSid",
				Values:     []string{"dummy"},
				ByteValues: [][]byte{nil},
			},
			expectError: true,
			errorMsg:    "Missing or nil SID bytes for attribute: objectSid",
		},
		{
			name: "empty_byte_values_sidHistory",
			attr: &ldap_v3.EntryAttribute{
				Name:       "SIDHistory",
				Values:     []string{"dummy"},
				ByteValues: [][]byte{},
			},
			expectError: true,
			errorMsg:    "Missing or nil SID bytes for attribute: SIDHistory",
		},
		{
			name: "empty_byte_values_creatorSID",
			attr: &ldap_v3.EntryAttribute{
				Name:       "mS-DS-CreatorSID",
				Values:     []string{"dummy"},
				ByteValues: [][]byte{},
			},
			expectError: true,
			errorMsg:    "Missing or nil SID bytes for attribute: mS-DS-CreatorSID",
		},
		{
			name: "empty_byte_values_securityIdentifier",
			attr: &ldap_v3.EntryAttribute{
				Name:       "securityIdentifier",
				Values:     []string{"dummy"},
				ByteValues: [][]byte{},
			},
			expectError: true,
			errorMsg:    "Missing or nil SID bytes for attribute: securityIdentifier",
		},
		{
			name: "empty_byte_values_nTSecurityDescriptor",
			attr: &ldap_v3.EntryAttribute{
				Name:       "nTSecurityDescriptor",
				Values:     []string{"dummy"},
				ByteValues: [][]byte{},
			},
			expectError: true,
			errorMsg:    "Missing or nil security descriptor bytes for attribute: nTSecurityDescriptor",
		},
		{
			name: "nil_byte_values_nTSecurityDescriptor",
			attr: &ldap_v3.EntryAttribute{
				Name:       "nTSecurityDescriptor",
				Values:     []string{"dummy"},
				ByteValues: [][]byte{nil},
			},
			expectError: true,
			errorMsg:    "Missing or nil security descriptor bytes for attribute: nTSecurityDescriptor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ldap.StringAttrValuesToRequestedType(tt.attr, false, framework.AttributeTypeString)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")

					return
				}

				if err.Message != tt.errorMsg {
					t.Errorf("expected error message '%s', got '%s'", tt.errorMsg, err.Message)
				}

				if err.Code != api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE {
					t.Errorf("expected error code ERROR_CODE_INVALID_ATTRIBUTE_TYPE, got %v", err.Code)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if result == nil {
					t.Errorf("expected result, got nil")
				}
			}
		})
	}
}
