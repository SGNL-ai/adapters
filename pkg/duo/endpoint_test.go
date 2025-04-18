// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package duo_test

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/duo"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestConstructEndpoint(t *testing.T) {
	tests := map[string]struct {
		request          *duo.Request
		wantEndpointInfo *duo.EndpointInfo
		wantError        *framework.Error
	}{
		"nil_request": {
			request: nil,
			wantError: &framework.Error{
				Message: "Request is nil.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"negative_cursor_offset": {
			request: &duo.Request{
				BaseURL:          "https://api-xxxxxxxx.duosecurity.com",
				APIVersion:       "v1",
				EntityExternalID: "User",
				PageSize:         100,
				IntegrationKey:   "testkey",
				Secret:           "testsecret",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](-1),
				},
			},
			wantError: &framework.Error{
				Message: "Cursor must be greater than 0.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"zero_cursor_offset": {
			request: &duo.Request{
				BaseURL:          "https://api-xxxxxxxx.duosecurity.com",
				APIVersion:       "v1",
				EntityExternalID: "User",
				PageSize:         100,
				IntegrationKey:   "testkey",
				Secret:           "testsecret",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](0),
				},
			},
			wantError: &framework.Error{
				Message: "Cursor must be greater than 0.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"valid_endpoint_without_cursor": {
			request: &duo.Request{
				BaseURL:          "https://api-xxxxxxxx.duosecurity.com",
				APIVersion:       "v1",
				EntityExternalID: "User",
				PageSize:         100,
				IntegrationKey:   "testkey",
				Secret:           "testsecret",
			},
			wantEndpointInfo: &duo.EndpointInfo{
				URL: "https://api-xxxxxxxx.duosecurity.com/admin/v1/users?limit=100&offset=0",
			},
		},
		"valid_endpoint_with_cursor": {
			request: &duo.Request{
				BaseURL:          "https://api-xxxxxxxx.duosecurity.com",
				APIVersion:       "v1",
				EntityExternalID: "User",
				PageSize:         100,
				IntegrationKey:   "testkey",
				Secret:           "testsecret",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](30),
				},
			},
			wantEndpointInfo: &duo.EndpointInfo{
				URL: "https://api-xxxxxxxx.duosecurity.com/admin/v1/users?limit=100&offset=30",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotEndpointInfo, gotError := duo.ConstructEndpoint(tt.request)

			if gotEndpointInfo != nil {
				parsedURL, err := url.Parse(gotEndpointInfo.URL)
				if err != nil {
					t.Errorf("Error parsing URL: %v", err)
				}

				baseURL := parsedURL.Host
				path := parsedURL.Path
				params := parsedURL.RawQuery

				tt.wantEndpointInfo.Date = gotEndpointInfo.Date
				tt.wantEndpointInfo.Auth = ConfigureAuth(tt.request, baseURL, gotEndpointInfo.Date, path, params)
			}

			if !reflect.DeepEqual(gotError, tt.wantError) {
				t.Errorf("gotError: %v, wantError: %v", gotError, tt.wantError)
			}

			if !reflect.DeepEqual(gotEndpointInfo, tt.wantEndpointInfo) {
				t.Errorf("gotEndpoint: %v, wantEndpoint: %v", gotEndpointInfo, tt.wantEndpointInfo)
			}
		})
	}
}

// ConfigureAuth configures the auth and date headers for the request per the Duo Admin API standards
// Example Request Signature: https://duo.com/docs/adminapi#authentication
// Tue, 21 Aug 2012 17:29:18 -0000
// GET
// api-xxxxxxxx.duosecurity.com
// /admin/v1/users
// &limit=1&offset=0.
func ConfigureAuth(request *duo.Request, baseURL, date, path, params string) string {
	hmac := hmac.New(sha1.New, []byte(request.Secret))
	hmac.Write([]byte(fmt.Sprintf("%s\nGET\n%s\n%s\n%s", date, baseURL, path, params)))

	token := fmt.Sprintf("%s:%s", request.IntegrationKey, hex.EncodeToString(hmac.Sum(nil)))
	auth := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(token)))

	return auth
}
