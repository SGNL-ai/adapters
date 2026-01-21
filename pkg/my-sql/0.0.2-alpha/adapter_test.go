// Copyright 2026 SGNL.ai, Inc.

package mysql_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/condexpr"
	"github.com/sgnl-ai/adapters/pkg/config"
	mysql_0_0_2_alpha "github.com/sgnl-ai/adapters/pkg/my-sql/0.0.2-alpha"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestAdapterGetPage(t *testing.T) {
	adapter := mysql_0_0_2_alpha.NewAdapter(mysql_0_0_2_alpha.NewClient(mysql_0_0_2_alpha.NewMockSQLClient()))

	tests := map[string]struct {
		request            *framework.Request[mysql_0_0_2_alpha.Config]
		inputRequestCursor interface{}
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request_first_page": {
			request: &framework.Request[mysql_0_0_2_alpha.Config]{
				Address: "127.0.0.1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "missing_field",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "active",
							Type:       framework.AttributeTypeBool,
						},
						{
							ExternalId: "employee_number",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "risk_score",
							Type:       framework.AttributeTypeDouble,
						},
						{
							ExternalId: "last_modified",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &mysql_0_0_2_alpha.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"active":          true,
							"employee_number": int64(1),
							"id":              "a20bab52-52e3-46c2-bd6a-2ad1512f713f",
							"last_modified":   time.Date(2025, 2, 12, 22, 38, 00, 00, time.UTC),
							"name":            "Ernesto Gregg",
							"risk_score":      float64(1),
						},
						{
							"active":          true,
							"employee_number": int64(2),
							"id":              "d35c298e-d343-4ad8-ac35-f7c5d9d47cb9",
							"last_modified":   time.Date(2025, 2, 12, 22, 38, 00, 00, time.UTC),
							"name":            "Eleanor Watts",
							"risk_score":      float64(1.562),
						},
						{
							"active":          true,
							"employee_number": int64(3),
							"id":              "62c74831-be4a-4cad-88fa-4e02640269d2",
							"last_modified":   time.Date(2025, 2, 12, 22, 38, 00, 00, time.UTC),
							"name":            "Chris Griffin",
							"risk_score":      float64(4.23),
						},
						{
							"active":          false,
							"employee_number": int64(4),
							"id":              "65b8fa65-25c5-4682-997f-ca86923e59e4",
							"last_modified":   time.Date(2025, 2, 12, 22, 38, 00, 00, time.UTC),
							"name":            "Casey Manning",
							"risk_score":      float64(10),
						},
						{
							"active":          true,
							"employee_number": int64(5),
							"id":              "6598acf9-cccc-48c9-ab9b-754bbe9ad146",
							"last_modified":   time.Date(2025, 2, 12, 22, 38, 00, 00, time.UTC),
							"name":            "Helen Gray",
							"risk_score":      float64(3.25),
						},
					},
					NextCursor: "6598acf9-cccc-48c9-ab9b-754bbe9ad146",
				},
			},
		},
		"valid_request_second_middle_page": {
			request: &framework.Request[mysql_0_0_2_alpha.Config]{
				Address: "127.0.0.1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "missing_field",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &mysql_0_0_2_alpha.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 5,
				Cursor:   "5",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":   "7390f7fc-0145-4691-9f55-b5c783369db9",
							"name": "Martha Pollard",
						},
						{
							"id":   "745cf6d6-55c8-4863-9bf6-1b1a80ff1515",
							"name": "Roxanne Dixon",
						},
						{
							"id":   "776b45f0-a2e3-4424-8ef7-84f3052bebc7",
							"name": "Verna Ferrell",
						},
						{
							"id":   "8b9643f9-25b4-458a-ad4f-81e61d106a57",
							"name": "Adrian Carey",
						},
						{
							"id":   "88ff7d742-fb3c-4103-af4b-fcd4315bae66",
							"name": "Joshua Martinez",
						},
					},
					NextCursor: "88ff7d742-fb3c-4103-af4b-fcd4315bae66",
				},
			},
		},
		"valid_request_third_second_last_page": {
			request: &framework.Request[mysql_0_0_2_alpha.Config]{
				Address: "127.0.0.1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "missing_field",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &mysql_0_0_2_alpha.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 5,
				Cursor:   "10",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":   "9cf5a596-0df2-4510-a403-9b514fd500b8",
							"name": "Erica Meadows",
						},
						{
							"id":   "987053f0-c06c-48ee-9c99-81f3a96af639",
							"name": "Carole Crawford",
						},
					},
					NextCursor: "987053f0-c06c-48ee-9c99-81f3a96af639",
				},
			},
		},
		"valid_request_fourth_last_page": {
			request: &framework.Request[mysql_0_0_2_alpha.Config]{
				Address: "127.0.0.1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "missing_field",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &mysql_0_0_2_alpha.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 5,
				Cursor:   "12",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{},
			},
		},
		"failed_to_connect_to_datasource": {
			request: &framework.Request[mysql_0_0_2_alpha.Config]{
				Address: "127.0.0.1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "missing_field",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &mysql_0_0_2_alpha.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: mysql_0_0_2_alpha.TestDatasourceForConnectFailure,
				},
				Ordered:  true,
				PageSize: 5,
				Cursor:   "10",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to connect to datasource: failed to connect to mock sql service.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
				},
			},
		},
		"failed_to_query_datasource": {
			request: &framework.Request[mysql_0_0_2_alpha.Config]{
				Address: "127.0.0.1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "missing_field",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &mysql_0_0_2_alpha.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered: true,
				// Hardcoded PageSize + Cursor in Mock SQL to prompt query failure.
				PageSize: 1,
				Cursor:   "101",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to query datasource: failed to query mock sql service.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
				},
			},
		},
		"failed_to_cast_to_type": {
			request: &framework.Request[mysql_0_0_2_alpha.Config]{
				Address: "127.0.0.1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						// Returned as an string, unable to cast to int
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeInt64,
						},
					},
				},
				Config: &mysql_0_0_2_alpha.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 5,
				Cursor:   "10",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: `Failed to parse attribute: (name) strconv.ParseFloat: parsing "Erica Meadows": invalid syntax.`,
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
				},
			},
		},
		"max_int_float": {
			request: &framework.Request[mysql_0_0_2_alpha.Config]{
				Address: "127.0.0.1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "missing_field",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "active",
							Type:       framework.AttributeTypeBool,
						},
						{
							ExternalId: "employee_number",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "risk_score",
							Type:       framework.AttributeTypeDouble,
						},
						{
							ExternalId: "last_modified",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &mysql_0_0_2_alpha.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 5,
				Cursor:   "202",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"active": true,
							// See `mock_sql_client.go` for context on why we're using `1<<53-1` instead of `MaxInt64`
							"employee_number": int64(1<<53 - 1),
							"id":              "9cf5a596-0df2-4510-a403-9b514fd500b8",
							"last_modified":   time.Date(2025, 2, 12, 22, 38, 00, 00, time.UTC),
							"name":            "Erica Meadows",
							"risk_score":      math.MaxFloat64,
						},
						{
							"active":          true,
							"employee_number": int64(math.MinInt64),
							"id":              "dfaf01cc-85b7-4e2e-b2d7-608d1f1904fe",
							"last_modified":   time.Date(2025, 2, 12, 22, 38, 00, 00, time.UTC),
							"name":            "Eleanor Watts",
							"risk_score":      -math.MaxFloat64,
						},
					},
					NextCursor: "dfaf01cc-85b7-4e2e-b2d7-608d1f1904fe",
				},
			},
		},
		"null_values": {
			request: &framework.Request[mysql_0_0_2_alpha.Config]{
				Address: "127.0.0.1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "missing_field",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "active",
							Type:       framework.AttributeTypeBool,
						},
						{
							ExternalId: "employee_number",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "risk_score",
							Type:       framework.AttributeTypeDouble,
						},
						{
							ExternalId: "last_modified",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &mysql_0_0_2_alpha.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 5,
				Cursor:   "203",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id": "9cf5a596-0df2-4510-a403-9b514fd500b8",
						},
						{
							"id": "8f678b7c-2571-45fe-ba01-a6cad31b02de",
						},
						{
							"active":          true,
							"employee_number": int64(1),
							"id":              "a20bab52-52e3-46c2-bd6a-2ad1512f713f",
							"last_modified":   time.Date(2025, 2, 12, 22, 38, 00, 00, time.UTC),
							"name":            "Ernesto Gregg",
							"risk_score":      float64(1),
						},
					},
					NextCursor: "a20bab52-52e3-46c2-bd6a-2ad1512f713f",
				},
			},
		},
		"validation_prevent_sql_injection_table": {
			request: &framework.Request[mysql_0_0_2_alpha.Config]{
				Address: "127.0.0.1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "DROP sampletable;-- ",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
					},
				},
				Config: &mysql_0_0_2_alpha.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 5,
				Cursor:   "10",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "SQL table name validation failed: unsupported characters found or length is not in range 1-128.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				},
			},
		},
		"validation_prevent_sql_injection_unique_attribute": {
			request: &framework.Request[mysql_0_0_2_alpha.Config]{
				Address: "127.0.0.1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "DROP sampletable;-- ",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
					},
				},
				Config: &mysql_0_0_2_alpha.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 5,
				Cursor:   "10",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "SQL unique attribute validation failed: unsupported characters found or length is not in range 1-128.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				},
			},
		},
		"valid_request_first_page_cast_multiple_fields_to_string": {
			request: &framework.Request[mysql_0_0_2_alpha.Config]{
				Address: "127.0.0.1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "missing_field",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "active",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "employee_number",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "risk_score",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "last_modified",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &mysql_0_0_2_alpha.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"active":          "true",
							"employee_number": "1",
							"id":              "a20bab52-52e3-46c2-bd6a-2ad1512f713f",
							"last_modified":   "2025-02-12T22:38:00+00:00",
							"name":            "Ernesto Gregg",
							"risk_score":      "1",
						},
						{
							"active":          "true",
							"employee_number": "2",
							"id":              "d35c298e-d343-4ad8-ac35-f7c5d9d47cb9",
							"last_modified":   "2025-02-12T22:38:00+00:00",
							"name":            "Eleanor Watts",
							"risk_score":      "1.562",
						},
						{
							"active":          "true",
							"employee_number": "3",
							"id":              "62c74831-be4a-4cad-88fa-4e02640269d2",
							"last_modified":   "2025-02-12T22:38:00+00:00",
							"name":            "Chris Griffin",
							"risk_score":      "4.23",
						},
						{
							"active":          "false",
							"employee_number": "4",
							"id":              "65b8fa65-25c5-4682-997f-ca86923e59e4",
							"last_modified":   "2025-02-12T22:38:00+00:00",
							"name":            "Casey Manning",
							"risk_score":      "10",
						},
						{
							"active":          "true",
							"employee_number": "5",
							"id":              "6598acf9-cccc-48c9-ab9b-754bbe9ad146",
							"last_modified":   "2025-02-12T22:38:00+00:00",
							"name":            "Helen Gray",
							"risk_score":      "3.25",
						},
					},
					NextCursor: "6598acf9-cccc-48c9-ab9b-754bbe9ad146",
				},
			},
		},
		"valid_request_fields_not_in_db": {
			request: &framework.Request[mysql_0_0_2_alpha.Config]{
				Address: "127.0.0.1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "invalid_field_not_present",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				Config: &mysql_0_0_2_alpha.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id": "a20bab52-52e3-46c2-bd6a-2ad1512f713f",
						},
						{
							"id": "d35c298e-d343-4ad8-ac35-f7c5d9d47cb9",
						},
						{
							"id": "62c74831-be4a-4cad-88fa-4e02640269d2",
						},
						{
							"id": "65b8fa65-25c5-4682-997f-ca86923e59e4",
						},
						{
							"id": "6598acf9-cccc-48c9-ab9b-754bbe9ad146",
						},
					},
					NextCursor: "6598acf9-cccc-48c9-ab9b-754bbe9ad146",
				},
			},
		},
		"valid_request_first_page_filtered": {
			request: &framework.Request[mysql_0_0_2_alpha.Config]{
				Address: "127.0.0.1",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "missing_field",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "active",
							Type:       framework.AttributeTypeBool,
						},
						{
							ExternalId: "employee_number",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "risk_score",
							Type:       framework.AttributeTypeDouble,
						},
						{
							ExternalId: "last_modified",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &mysql_0_0_2_alpha.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
					Filters: map[string]condexpr.Condition{
						"users": {
							And: []condexpr.Condition{
								{
									Field:    "active",
									Operator: "=",
									Value:    true,
								},
								{
									Field:    "risk_score",
									Operator: ">",
									Value:    2.0,
								},
							},
						},
					},
				},
				Ordered:  true,
				PageSize: 5,
				Cursor:   "204",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"active":          true,
							"employee_number": int64(3),
							"id":              "62c74831-be4a-4cad-88fa-4e02640269d2",
							"last_modified":   time.Date(2025, 2, 12, 22, 38, 00, 00, time.UTC),
							"name":            "Chris Griffin",
							"risk_score":      float64(4.23),
						},
						{
							"active":          true,
							"employee_number": int64(5),
							"id":              "6598acf9-cccc-48c9-ab9b-754bbe9ad146",
							"last_modified":   time.Date(2025, 2, 12, 22, 38, 00, 00, time.UTC),
							"name":            "Helen Gray",
							"risk_score":      float64(3.25),
						},
					},
					NextCursor: "6598acf9-cccc-48c9-ab9b-754bbe9ad146",
				},
			},
		},
		// TODO: Test with missing unique ID. Current mock doesn't cover this.
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

			gotResponse := adapter.GetPage(context.Background(), tt.request)

			if diff := cmp.Diff(tt.wantResponse, gotResponse); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
			if gotResponse.Success != nil && tt.wantCursor != nil {
				var gotCursor pagination.CompositeCursor[string]

				decodedCursor, err := base64.StdEncoding.DecodeString(gotResponse.Success.NextCursor)
				if err != nil {
					t.Errorf("error decoding cursor: %v", err)
				}

				if err := json.Unmarshal(decodedCursor, &gotCursor); err != nil {
					t.Errorf("error unmarshaling cursor: %v", err)
				}

				if !reflect.DeepEqual(&gotCursor, tt.wantCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}
		})
	}
}
