// Copyright 2026 SGNL.ai, Inc.

// nolint: lll, goconst
package bamboohr_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/bamboohr"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := bamboohr.NewAdapter(&bamboohr.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[bamboohr.Config]
		inputRequestCursor interface{}
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[int64]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[bamboohr.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "apiKey123",
						Password: "randomString",
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnltestdev",
					OnlyCurrent:   true,
				},
				Entity:   *PopulateDefaultEmployeeEntityConfig(),
				PageSize: 3,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                    int64(4),
							"bestEmail":             "cabbott@efficientoffice.com",
							"dateOfBirth":           time.Date(1996, 9, 2, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Charlotte Abbott",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(9),
							"supervisorEmail":       "jcaldwell@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 50, 0, time.UTC),
						},
						{
							"id":                    int64(5),
							"bestEmail":             "aadams@efficientoffice.com",
							"dateOfBirth":           time.Date(1983, 6, 30, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Ashley Adams",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(9),
							"supervisorEmail":       "jcaldwell@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 50, 0, time.UTC),
						},
						{
							"id":                    int64(6),
							"bestEmail":             "cagluinda@efficientoffice.com",
							"dateOfBirth":           time.Date(1996, 8, 27, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Christina Agluinda",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(9),
							"supervisorEmail":       "jcaldwell@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 50, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOjN9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](3),
			},
		},
		"valid_request_no_https_prefix": {
			ctx: context.Background(),
			request: &framework.Request[bamboohr.Config]{
				Address: strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "apiKey123",
						Password: "randomString",
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnltestdev",
					OnlyCurrent:   true,
				},
				Entity:   *PopulateDefaultEmployeeEntityConfig(),
				PageSize: 3,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                    int64(4),
							"bestEmail":             "cabbott@efficientoffice.com",
							"dateOfBirth":           time.Date(1996, 9, 2, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Charlotte Abbott",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(9),
							"supervisorEmail":       "jcaldwell@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 50, 0, time.UTC),
						},
						{
							"id":                    int64(5),
							"bestEmail":             "aadams@efficientoffice.com",
							"dateOfBirth":           time.Date(1983, 6, 30, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Ashley Adams",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(9),
							"supervisorEmail":       "jcaldwell@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 50, 0, time.UTC),
						},
						{
							"id":                    int64(6),
							"bestEmail":             "cagluinda@efficientoffice.com",
							"dateOfBirth":           time.Date(1996, 8, 27, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Christina Agluinda",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(9),
							"supervisorEmail":       "jcaldwell@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 50, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOjN9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](3),
			},
		},
		"invalid_request_invalid_api_version": {
			ctx: context.Background(),
			request: &framework.Request[bamboohr.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "apiKey123",
						Password: "randomString",
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v2",
					CompanyDomain: "sgnltestdev",
				},
				Entity:   *PopulateDefaultEmployeeEntityConfig(),
				PageSize: 3,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "BambooHR config is invalid: apiVersion is not supported: v2.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"invalid_request_http_prefix": {
			ctx: context.Background(),
			request: &framework.Request[bamboohr.Config]{
				Address: "http://" + strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "apiKey123",
						Password: "randomString",
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnltestdev",
				},
				Entity:   *PopulateDefaultEmployeeEntityConfig(),
				PageSize: 3,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: `Scheme "http" is not supported.`,
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"malformed_cursor_negative_offset": {
			ctx: context.Background(),
			request: &framework.Request[bamboohr.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "apiKey123",
						Password: "randomString",
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnltestdev",
				},
				Entity:   *PopulateDefaultEmployeeEntityConfig(),
				PageSize: 3,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](-50),
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Cursor value must be greater than 0.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
		"malformed_composite_cursor_string_type": {
			ctx: context.Background(),
			request: &framework.Request[bamboohr.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "apiKey123",
						Password: "randomString",
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnltestdev",
				},
				Entity:   *PopulateDefaultEmployeeEntityConfig(),
				PageSize: 3,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("BROKEN"),
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to unmarshal JSON cursor: json: cannot unmarshal string into Go struct field CompositeCursor[int64].cursor of type int64.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
		"malformed_cursor_includes_collection_cursor": {
			ctx: context.Background(),
			request: &framework.Request[bamboohr.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "apiKey123",
						Password: "randomString",
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnltestdev",
				},
				Entity:   *PopulateDefaultEmployeeEntityConfig(),
				PageSize: 3,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor:           testutil.GenPtr[int64](50),
				CollectionCursor: testutil.GenPtr[int64](10),
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Cursor must not contain CollectionID or CollectionCursor fields for entity Employee.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
		"valid_request_with_cursor": {
			ctx: context.Background(),
			request: &framework.Request[bamboohr.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "apiKey123",
						Password: "randomString",
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnltestdev",
					OnlyCurrent:   true,
				},
				Entity:   *PopulateDefaultEmployeeEntityConfig(),
				PageSize: 3,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](3),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                    int64(7),
							"bestEmail":             "sanderson@efficientoffice.com",
							"fullName1":             "Shannon Anderson",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(9),
							"supervisorEmail":       "jcaldwell@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 50, 0, time.UTC),
						},
						{
							"id":                    int64(8),
							"bestEmail":             "arvind@sgnl.ai",
							"fullName1":             "Arvind",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 50, 0, time.UTC),
						},
						{
							"id":                    int64(9),
							"bestEmail":             "jcaldwell@efficientoffice.com",
							"dateOfBirth":           time.Date(1975, 1, 26, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Jennifer Caldwell",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(8),
							"supervisorEmail":       "arvind@sgnl.ai",
						},
					},
					NextCursor: "eyJjdXJzb3IiOjZ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](6),
			},
		},
		"invalid_request_invalid_url": {
			ctx: context.Background(),
			request: &framework.Request[bamboohr.Config]{
				Address: server.URL + "/invalid",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "apiKey123",
						Password: "randomString",
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnltestdev",
				},
				Entity:   *PopulateDefaultEmployeeEntityConfig(),
				PageSize: 3,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Datasource rejected request, returned status code: 404.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				var encodedCursor string

				var err *framework.Error

				switch v := tt.inputRequestCursor.(type) {
				case *pagination.CompositeCursor[int64]:
					encodedCursor, err = pagination.MarshalCursor(v)
				case *pagination.CompositeCursor[string]:
					encodedCursor, err = pagination.MarshalCursor(v)
				default:
					t.Errorf("Unsupported cursor type: %T", tt.inputRequestCursor)
				}

				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if diff := cmp.Diff(tt.wantResponse, gotResponse); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
			if gotResponse.Success != nil && tt.wantCursor != nil {
				var gotCursor pagination.CompositeCursor[int64]

				decodedCursor, err := base64.StdEncoding.DecodeString(gotResponse.Success.NextCursor)
				if err != nil {
					t.Errorf("error decoding cursor: %v", err)
				}

				if err := json.Unmarshal(decodedCursor, &gotCursor); err != nil {
					t.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(&gotCursor, tt.wantCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetEmployeePage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := bamboohr.NewAdapter(&bamboohr.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[bamboohr.Config]
		inputRequestCursor *pagination.CompositeCursor[int64]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[int64]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[bamboohr.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "apiKey123",
						Password: "randomString",
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnltestdev",
					OnlyCurrent:   true,
				},
				Entity:   *PopulateDefaultEmployeeEntityConfig(),
				PageSize: 10,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                    int64(4),
							"bestEmail":             "cabbott@efficientoffice.com",
							"dateOfBirth":           time.Date(1996, 9, 2, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Charlotte Abbott",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(9),
							"supervisorEmail":       "jcaldwell@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 50, 0, time.UTC),
						},
						{
							"id":                    int64(5),
							"bestEmail":             "aadams@efficientoffice.com",
							"dateOfBirth":           time.Date(1983, 6, 30, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Ashley Adams",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(9),
							"supervisorEmail":       "jcaldwell@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 50, 0, time.UTC),
						},
						{
							"id":                    int64(6),
							"bestEmail":             "cagluinda@efficientoffice.com",
							"dateOfBirth":           time.Date(1996, 8, 27, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Christina Agluinda",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(9),
							"supervisorEmail":       "jcaldwell@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 50, 0, time.UTC),
						},
						{
							"id":                    int64(7),
							"bestEmail":             "sanderson@efficientoffice.com",
							"fullName1":             "Shannon Anderson",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(9),
							"supervisorEmail":       "jcaldwell@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 50, 0, time.UTC),
						},
						{
							"id":                    int64(8),
							"bestEmail":             "arvind@sgnl.ai",
							"fullName1":             "Arvind",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 50, 0, time.UTC),
						},
						{
							"id":                    int64(9),
							"bestEmail":             "jcaldwell@efficientoffice.com",
							"dateOfBirth":           time.Date(1975, 1, 26, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Jennifer Caldwell",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(8),
							"supervisorEmail":       "arvind@sgnl.ai",
						},
						{
							"id":                    int64(10),
							"bestEmail":             "rsaito@efficientoffice.com",
							"dateOfBirth":           time.Date(1968, 12, 28, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Ryota Saito",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(8),
							"supervisorEmail":       "arvind@sgnl.ai",
						},
						{
							"id":                    int64(11),
							"bestEmail":             "dvance@efficientoffice.com",
							"dateOfBirth":           time.Date(1978, 8, 23, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Daniel Vance",
							"customcustomBoolField": false,
							"supervisorEId":         int64(8),
							"supervisorEmail":       "arvind@sgnl.ai",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(12),
							"bestEmail":             "easture@efficientoffice.com",
							"dateOfBirth":           time.Date(1990, 7, 1, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Eric Asture",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEmail":       "arvind@sgnl.ai",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(13),
							"bestEmail":             "cbarnet@efficientoffice.com",
							"dateOfBirth":           time.Date(1987, 6, 16, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Cheryl Barnet",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEmail":       "arvind@sgnl.ai",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 50, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOjEwfQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](10),
			},
		},
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[bamboohr.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "apiKey123",
						Password: "randomString",
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnltestdev",
					OnlyCurrent:   true,
				},
				Entity:   *PopulateDefaultEmployeeEntityConfig(),
				PageSize: 10,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](10),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                    int64(14),
							"bestEmail":             "mandev@efficientoffice.com",
							"dateOfBirth":           time.Date(1987, 6, 5, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Maja Andev",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(8),
							"supervisorEmail":       "arvind@sgnl.ai",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(15),
							"bestEmail":             "twalsh@efficientoffice.com",
							"dateOfBirth":           time.Date(1981, 3, 18, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Trent Walsh",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(8),
							"supervisorEmail":       "arvind@sgnl.ai",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(16),
							"bestEmail":             "jbryan@efficientoffice.com",
							"dateOfBirth":           time.Date(1970, 12, 7, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Jake Bryan",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(8),
							"supervisorEmail":       "arvind@sgnl.ai",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(17),
							"bestEmail":             "dchou@efficientoffice.com",
							"dateOfBirth":           time.Date(1987, 5, 8, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Dorothy Chou",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(8),
							"supervisorEmail":       "arvind@sgnl.ai",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(18),
							"bestEmail":             "javier@efficientoffice.com",
							"dateOfBirth":           time.Date(1996, 8, 28, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Javier Cruz",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(15),
							"supervisorEmail":       "twalsh@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(19),
							"bestEmail":             "shelly@efficientoffice.com",
							"dateOfBirth":           time.Date(1993, 6, 1, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Shelly Cluff",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(10),
							"supervisorEmail":       "rsaito@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(22),
							"bestEmail":             "dillon@efficientoffice.com",
							"dateOfBirth":           time.Date(1972, 6, 6, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Dillon (Remote) Park",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(105),
							"supervisorEmail":       "norma@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(23),
							"bestEmail":             "darlene@efficientoffice.com",
							"dateOfBirth":           time.Date(1975, 9, 16, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Darlene Handley",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(15),
							"supervisorEmail":       "twalsh@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(24),
							"bestEmail":             "zack@efficientoffice.com",
							"dateOfBirth":           time.Date(2000, 8, 2, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Zack Miller",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(11),
							"supervisorEmail":       "dvance@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(25),
							"bestEmail":             "philip@efficientoffice.com",
							"dateOfBirth":           time.Date(1975, 11, 26, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Philip Wagener",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(12),
							"supervisorEmail":       "easture@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOjIwfQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](20),
			},
		},
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[bamboohr.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "apiKey123",
						Password: "randomString",
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnltestdev",
					OnlyCurrent:   true,
				},
				Entity:   *PopulateDefaultEmployeeEntityConfig(),
				PageSize: 10,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](20),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                    int64(26),
							"bestEmail":             "agranger@efficientoffice.com",
							"dateOfBirth":           time.Date(1998, 11, 26, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Amy Granger",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(13),
							"supervisorEmail":       "cbarnet@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(27),
							"bestEmail":             "debra@efficientoffice.com",
							"dateOfBirth":           time.Date(1966, 10, 18, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Debra Tuescher",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(49),
							"supervisorEmail":       "robert@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(28),
							"bestEmail":             "andy@efficientoffice.com",
							"dateOfBirth":           time.Date(2001, 2, 25, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Andy Graves",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(26),
							"supervisorEmail":       "agranger@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(29),
							"bestEmail":             "catherine@efficientoffice.com",
							"dateOfBirth":           time.Date(1993, 12, 18, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Catherine Jones",
							"isPhotoUploaded":       false,
							"customcustomBoolField": false,
							"supervisorEId":         int64(4),
							"supervisorEmail":       "cabbott@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 47, 0, time.UTC),
						},
						{
							"id":                    int64(30),
							"bestEmail":             "corey@efficientoffice.com",
							"dateOfBirth":           time.Date(1995, 5, 1, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Corey Ross",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(15),
							"supervisorEmail":       "twalsh@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(31),
							"bestEmail":             "sally@efficientoffice.com",
							"dateOfBirth":           time.Date(1984, 5, 31, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Sally Harmon",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(16),
							"supervisorEmail":       "jbryan@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(32),
							"bestEmail":             "carly@efficientoffice.com",
							"dateOfBirth":           time.Date(1984, 7, 2, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Carly Seymour",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(52),
							"supervisorEmail":       "nate@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
						{
							"id":                    int64(33),
							"bestEmail":             "erin@efficientoffice.com",
							"dateOfBirth":           time.Date(1993, 2, 26, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Erin Farr",
							"isPhotoUploaded":       false,
							"customcustomBoolField": false,
							"supervisorEId":         int64(16),
							"supervisorEmail":       "jbryan@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 47, 0, time.UTC),
						},
						{
							"id":                    int64(34),
							"bestEmail":             "emily@efficientoffice.com",
							"dateOfBirth":           time.Date(1994, 4, 30, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Emily Gomez",
							"isPhotoUploaded":       false,
							"customcustomBoolField": false,
							"supervisorEId":         int64(36),
							"supervisorEmail":       "melany@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 47, 0, time.UTC),
						},
						{
							"id":                    int64(35),
							"bestEmail":             "aaron@efficientoffice.com",
							"dateOfBirth":           time.Date(1998, 8, 16, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Aaron Eckerly",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(36),
							"supervisorEmail":       "melany@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 48, 0, time.UTC),
						},
					},
					NextCursor: "eyJjdXJzb3IiOjMwfQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](30),
			},
		},
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[bamboohr.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "apiKey123",
						Password: "randomString",
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnltestdev",
					OnlyCurrent:   true,
				},
				Entity:   *PopulateDefaultEmployeeEntityConfig(),
				PageSize: 10,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](30),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                    int64(36),
							"bestEmail":             "melany@efficientoffice.com",
							"dateOfBirth":           time.Date(1986, 11, 25, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Melany Olsen",
							"isPhotoUploaded":       false,
							"customcustomBoolField": false,
							"supervisorEId":         int64(13),
							"supervisorEmail":       "cbarnet@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 47, 0, time.UTC),
						},
						{
							"id":                    int64(37),
							"bestEmail":             "whitney@efficientoffice.com",
							"dateOfBirth":           time.Date(1992, 12, 30, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Whitney Webster",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(5),
							"supervisorEmail":       "aadams@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 47, 0, time.UTC),
						},
						{
							"id":                    int64(38),
							"bestEmail":             "marrissa@efficientoffice.com",
							"dateOfBirth":           time.Date(1995, 1, 31, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Marrissa Mellon",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(69),
							"supervisorEmail":       "karin@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 47, 0, time.UTC),
						},
						{
							"id":                    int64(39),
							"bestEmail":             "paige@efficientoffice.com",
							"dateOfBirth":           time.Date(1993, 2, 2, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Paige Rasmussen",
							"isPhotoUploaded":       true,
							"customcustomBoolField": false,
							"supervisorEId":         int64(57),
							"supervisorEmail":       "liam@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 47, 0, time.UTC),
						},
						{
							"id":                    int64(40),
							"bestEmail":             "kelli@efficientoffice.com",
							"dateOfBirth":           time.Date(1988, 3, 1, 0, 0, 0, 0, time.UTC),
							"fullName1":             "Kelli Crandle",
							"isPhotoUploaded":       false,
							"customcustomBoolField": false,
							"supervisorEId":         int64(49),
							"supervisorEmail":       "robert@efficientoffice.com",
							"lastChanged":           time.Date(2024, 4, 12, 19, 33, 47, 0, time.UTC),
						},
					},
					NextCursor: "",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if diff := cmp.Diff(tt.wantResponse, gotResponse); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
			if gotResponse.Success != nil && tt.wantCursor != nil {
				var gotCursor pagination.CompositeCursor[int64]

				decodedCursor, err := base64.StdEncoding.DecodeString(gotResponse.Success.NextCursor)
				if err != nil {
					t.Errorf("error decoding cursor: %v", err)
				}

				if err := json.Unmarshal(decodedCursor, &gotCursor); err != nil {
					t.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(&gotCursor, tt.wantCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}
		})
	}
}
