// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package bamboohr_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	adapter_api_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/bamboohr"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestParseResponse(t *testing.T) {
	tests := map[string]struct {
		request        *bamboohr.Request
		body           []byte
		wantObjects    []map[string]any
		wantNextCursor *pagination.CompositeCursor[int64]
		wantErr        *framework.Error
	}{
		"page_size_greater_than_num_objects": {
			request: &bamboohr.Request{
				PageSize:          100,
				EntityConfig:      PopulateDefaultEmployeeEntityConfig(),
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "6",
						"bestEmail": "cagluinda@efficientoffice.com",
						"dateOfBirth": "1996-08-27",
						"fullName1": "Christina Agluinda",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "7",
						"bestEmail": "sanderson@efficientoffice.com",
						"dateOfBirth": "2000-05-08",
						"fullName1": "Shannon Anderson",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					}
				]
			}`),
			wantObjects: []map[string]any{
				{
					"id":                    float64(4),
					"bestEmail":             "cabbott@efficientoffice.com",
					"dateOfBirth":           "1996-09-02",
					"fullName1":             "Charlotte Abbott",
					"isPhotoUploaded":       true,
					"customcustomBoolField": false,
					"supervisorEId":         float64(9),
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
				{
					"id":                    float64(5),
					"bestEmail":             "aadams@efficientoffice.com",
					"dateOfBirth":           "1983-06-30",
					"fullName1":             "Ashley Adams",
					"isPhotoUploaded":       true,
					"customcustomBoolField": false,
					"supervisorEId":         float64(9),
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
				{
					"id":                    float64(6),
					"bestEmail":             "cagluinda@efficientoffice.com",
					"dateOfBirth":           "1996-08-27",
					"fullName1":             "Christina Agluinda",
					"isPhotoUploaded":       true,
					"customcustomBoolField": false,
					"supervisorEId":         float64(9),
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
				{
					"id":                    float64(7),
					"bestEmail":             "sanderson@efficientoffice.com",
					"dateOfBirth":           "2000-05-08",
					"fullName1":             "Shannon Anderson",
					"isPhotoUploaded":       true,
					"customcustomBoolField": false,
					"supervisorEId":         float64(9),
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
			},
			wantNextCursor: nil,
		},
		"page_size_greater_than_num_objects_with_cursor": {
			request: &bamboohr.Request{
				PageSize: 100,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(2)),
				},
				EntityConfig:      PopulateDefaultEmployeeEntityConfig(),
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "6",
						"bestEmail": "cagluinda@efficientoffice.com",
						"dateOfBirth": "1996-08-27",
						"fullName1": "Christina Agluinda",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "7",
						"bestEmail": "sanderson@efficientoffice.com",
						"dateOfBirth": "2000-05-08",
						"fullName1": "Shannon Anderson",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					}
				]
			}`),
			wantObjects: []map[string]any{
				{
					"id":                    float64(6),
					"bestEmail":             "cagluinda@efficientoffice.com",
					"dateOfBirth":           "1996-08-27",
					"fullName1":             "Christina Agluinda",
					"isPhotoUploaded":       true,
					"customcustomBoolField": false,
					"supervisorEId":         float64(9),
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
				{
					"id":                    float64(7),
					"bestEmail":             "sanderson@efficientoffice.com",
					"dateOfBirth":           "2000-05-08",
					"fullName1":             "Shannon Anderson",
					"isPhotoUploaded":       true,
					"customcustomBoolField": false,
					"supervisorEId":         float64(9),
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
			},
			wantNextCursor: nil,
		},
		"page_size_less_than_num_objects_with_cursor": {
			request: &bamboohr.Request{
				PageSize: 1,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(2)),
				},
				EntityConfig:      PopulateDefaultEmployeeEntityConfig(),
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "6",
						"bestEmail": "cagluinda@efficientoffice.com",
						"dateOfBirth": "1996-08-27",
						"fullName1": "Christina Agluinda",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "7",
						"bestEmail": "sanderson@efficientoffice.com",
						"dateOfBirth": "2000-05-08",
						"fullName1": "Shannon Anderson",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					}
				]
			}`),
			wantObjects: []map[string]any{
				{
					"id":                    float64(6),
					"bestEmail":             "cagluinda@efficientoffice.com",
					"dateOfBirth":           "1996-08-27",
					"fullName1":             "Christina Agluinda",
					"isPhotoUploaded":       true,
					"customcustomBoolField": false,
					"supervisorEId":         float64(9),
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
			},
			wantNextCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(3)),
			},
		},
		"page_size_less_than_num_objects_without_cursor": {
			request: &bamboohr.Request{
				PageSize:          2,
				EntityConfig:      PopulateDefaultEmployeeEntityConfig(),
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "6",
						"bestEmail": "cagluinda@efficientoffice.com",
						"dateOfBirth": "1996-08-27",
						"fullName1": "Christina Agluinda",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "7",
						"bestEmail": "sanderson@efficientoffice.com",
						"dateOfBirth": "2000-05-08",
						"fullName1": "Shannon Anderson",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					}
				]
			}`),
			wantObjects: []map[string]any{
				{
					"id":                    float64(4),
					"bestEmail":             "cabbott@efficientoffice.com",
					"dateOfBirth":           "1996-09-02",
					"fullName1":             "Charlotte Abbott",
					"isPhotoUploaded":       true,
					"customcustomBoolField": false,
					"supervisorEId":         float64(9),
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
				{
					"id":                    float64(5),
					"bestEmail":             "aadams@efficientoffice.com",
					"dateOfBirth":           "1983-06-30",
					"fullName1":             "Ashley Adams",
					"isPhotoUploaded":       true,
					"customcustomBoolField": false,
					"supervisorEId":         float64(9),
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
			},
			wantNextCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(2)),
			},
		},
		"invalid_requested_int_type_conversion": {
			request: &bamboohr.Request{
				PageSize: 100,
				EntityConfig: &framework.EntityConfig{
					ExternalId: "employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "supervisorEmail",
							Type:       framework.AttributeTypeInt64,
							List:       false,
						},
					},
				},
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "6",
						"bestEmail": "cagluinda@efficientoffice.com",
						"dateOfBirth": "1996-08-27",
						"fullName1": "Christina Agluinda",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "7",
						"bestEmail": "sanderson@efficientoffice.com",
						"dateOfBirth": "2000-05-08",
						"fullName1": "Shannon Anderson",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					}
				]
			}`),
			wantErr: &framework.Error{
				Message: "Failed to parse attribute: supervisorEmail, as type Float with value: jcaldwell@efficientoffice.com.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
			},
			wantNextCursor: nil,
		},
		"invalid_requested_bool_type_conversion": {
			request: &bamboohr.Request{
				PageSize: 100,
				EntityConfig: &framework.EntityConfig{
					ExternalId: "employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "supervisorEmail",
							Type:       framework.AttributeTypeBool,
							List:       false,
						},
					},
				},
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "6",
						"bestEmail": "cagluinda@efficientoffice.com",
						"dateOfBirth": "1996-08-27",
						"fullName1": "Christina Agluinda",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "7",
						"bestEmail": "sanderson@efficientoffice.com",
						"dateOfBirth": "2000-05-08",
						"fullName1": "Shannon Anderson",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					}
				]
			}`),
			wantErr: &framework.Error{
				Message: "Failed to parse attribute: supervisorEmail, as type Bool with value: jcaldwell@efficientoffice.com.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
			},
			wantNextCursor: nil,
		},
		"null_value_conversions": {
			request: &bamboohr.Request{
				PageSize:          100,
				EntityConfig:      PopulateDefaultEmployeeEntityConfig(),
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "0000-00-00",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "Null",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "0000-00-00T00:00:00+00:00"
					},
					{
						"id": "6",
						"bestEmail": "cagluinda@efficientoffice.com",
						"dateOfBirth": "null",
						"fullName1": "",
						"isPhotoUploaded": "null",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					}
				]
			}`),
			wantObjects: []map[string]any{
				{
					"id":                    float64(4),
					"bestEmail":             "cabbott@efficientoffice.com",
					"dateOfBirth":           nil,
					"fullName1":             "Charlotte Abbott",
					"isPhotoUploaded":       true,
					"customcustomBoolField": false,
					"supervisorEId":         float64(9),
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
				{
					"id":                    float64(5),
					"bestEmail":             "aadams@efficientoffice.com",
					"dateOfBirth":           nil,
					"fullName1":             "Ashley Adams",
					"isPhotoUploaded":       true,
					"customcustomBoolField": false,
					"supervisorEId":         nil,
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           nil,
				},
				{
					"id":                    float64(6),
					"bestEmail":             "cagluinda@efficientoffice.com",
					"dateOfBirth":           nil,
					"fullName1":             nil,
					"isPhotoUploaded":       nil,
					"customcustomBoolField": false,
					"supervisorEId":         float64(9),
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
			},
			wantNextCursor: nil,
		},
		"null_value_conversions_with_list": {
			request: &bamboohr.Request{
				PageSize: 100,
				AttributeMappings: &bamboohr.AttributeMappings{
					BoolMappings: &bamboohr.BoolAttributeMappings{
						True:  []string{"yes"},
						False: []string{"no"},
					},
				},
				EntityConfig: &framework.EntityConfig{
					ExternalId: "employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "customcustomBoolField",
							Type:       framework.AttributeTypeBool,
							List:       false,
						},
						{
							ExternalId: "list1",
							Type:       framework.AttributeTypeBool,
							List:       true,
						},
						{
							ExternalId: "list2",
							Type:       framework.AttributeTypeInt64,
							List:       true,
						},
					},
				},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"list1": ["yes", "", "null", "no"],
						"list2": ["1", "", "null", "2"]
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"list2": ["33", "null", "2"]
					}
				]
			}`),
			wantObjects: []map[string]any{
				{
					"id":                    "4",
					"bestEmail":             "cabbott@efficientoffice.com",
					"dateOfBirth":           "1996-09-02",
					"fullName1":             "Charlotte Abbott",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": true,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
					"list1":                 []any{true, nil, nil, false},
					"list2":                 []any{float64(1), nil, nil, float64(2)},
				},
				{
					"id":                    "5",
					"bestEmail":             "aadams@efficientoffice.com",
					"dateOfBirth":           "1983-06-30",
					"fullName1":             "Ashley Adams",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": true,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
					"list2":                 []any{float64(33), nil, float64(2)},
				},
			},
			wantNextCursor: nil,
		},
		"direct_bool_type_conversion_without_mapper": {
			request: &bamboohr.Request{
				PageSize: 100,
				EntityConfig: &framework.EntityConfig{
					ExternalId: "employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "customcustomBoolField",
							Type:       framework.AttributeTypeBool,
							List:       false,
						},
					},
				},
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "6",
						"bestEmail": "cagluinda@efficientoffice.com",
						"dateOfBirth": "1996-08-27",
						"fullName1": "Christina Agluinda",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "7",
						"bestEmail": "sanderson@efficientoffice.com",
						"dateOfBirth": "2000-05-08",
						"fullName1": "Shannon Anderson",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					}
				]
			}`),
			wantObjects: []map[string]any{
				{
					"id":                    "4",
					"bestEmail":             "cabbott@efficientoffice.com",
					"dateOfBirth":           "1996-09-02",
					"fullName1":             "Charlotte Abbott",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": false,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
				{
					"id":                    "5",
					"bestEmail":             "aadams@efficientoffice.com",
					"dateOfBirth":           "1983-06-30",
					"fullName1":             "Ashley Adams",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": false,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
				{
					"id":                    "6",
					"bestEmail":             "cagluinda@efficientoffice.com",
					"dateOfBirth":           "1996-08-27",
					"fullName1":             "Christina Agluinda",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": false,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
				{
					"id":                    "7",
					"bestEmail":             "sanderson@efficientoffice.com",
					"dateOfBirth":           "2000-05-08",
					"fullName1":             "Shannon Anderson",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": false,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
			},
			wantNextCursor: nil,
		},
		"bool_type_conversion_with_mapper": {
			request: &bamboohr.Request{
				PageSize: 100,
				AttributeMappings: &bamboohr.AttributeMappings{
					BoolMappings: &bamboohr.BoolAttributeMappings{
						True:  []string{"yes"},
						False: []string{"no"},
					},
				},
				EntityConfig: &framework.EntityConfig{
					ExternalId: "employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "customcustomBoolField",
							Type:       framework.AttributeTypeBool,
							List:       false,
						},
					},
				},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "6",
						"bestEmail": "cagluinda@efficientoffice.com",
						"dateOfBirth": "1996-08-27",
						"fullName1": "Christina Agluinda",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "no",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "7",
						"bestEmail": "sanderson@efficientoffice.com",
						"dateOfBirth": "2000-05-08",
						"fullName1": "Shannon Anderson",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "no",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					}
				]
			}`),
			wantObjects: []map[string]any{
				{
					"id":                    "4",
					"bestEmail":             "cabbott@efficientoffice.com",
					"dateOfBirth":           "1996-09-02",
					"fullName1":             "Charlotte Abbott",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": true,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
				{
					"id":                    "5",
					"bestEmail":             "aadams@efficientoffice.com",
					"dateOfBirth":           "1983-06-30",
					"fullName1":             "Ashley Adams",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": true,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
				{
					"id":                    "6",
					"bestEmail":             "cagluinda@efficientoffice.com",
					"dateOfBirth":           "1996-08-27",
					"fullName1":             "Christina Agluinda",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": false,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
				{
					"id":                    "7",
					"bestEmail":             "sanderson@efficientoffice.com",
					"dateOfBirth":           "2000-05-08",
					"fullName1":             "Shannon Anderson",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": false,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
				},
			},
			wantNextCursor: nil,
		},
		"int_list_conversion": {
			request: &bamboohr.Request{
				PageSize: 100,
				AttributeMappings: &bamboohr.AttributeMappings{
					BoolMappings: &bamboohr.BoolAttributeMappings{
						True:  []string{"yes"},
						False: []string{"no"},
					},
				},
				EntityConfig: &framework.EntityConfig{
					ExternalId: "employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "customcustomBoolField",
							Type:       framework.AttributeTypeBool,
							List:       false,
						},
						{
							ExternalId: "intList",
							Type:       framework.AttributeTypeInt64,
							List:       true,
						},
					},
				},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"intList": ["1", "2", "3", null]
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"intList": []
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"intList": null
					}
				]
			}`),
			wantObjects: []map[string]any{
				{
					"id":                    "4",
					"bestEmail":             "cabbott@efficientoffice.com",
					"dateOfBirth":           "1996-09-02",
					"fullName1":             "Charlotte Abbott",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": true,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
					"intList":               []any{float64(1), float64(2), float64(3), nil},
				},
				{
					"id":                    "5",
					"bestEmail":             "aadams@efficientoffice.com",
					"dateOfBirth":           "1983-06-30",
					"fullName1":             "Ashley Adams",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": true,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
					"intList":               []any{},
				},
				{
					"id":                    "5",
					"bestEmail":             "aadams@efficientoffice.com",
					"dateOfBirth":           "1983-06-30",
					"fullName1":             "Ashley Adams",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": true,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
					"intList":               nil,
				},
			},
			wantNextCursor: nil,
		},
		"float_list_conversion": {
			request: &bamboohr.Request{
				PageSize: 100,
				AttributeMappings: &bamboohr.AttributeMappings{
					BoolMappings: &bamboohr.BoolAttributeMappings{
						True:  []string{"yes"},
						False: []string{"no"},
					},
				},
				EntityConfig: &framework.EntityConfig{
					ExternalId: "employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "customcustomBoolField",
							Type:       framework.AttributeTypeBool,
							List:       false,
						},
						{
							ExternalId: "intList",
							Type:       framework.AttributeTypeDouble,
							List:       true,
						},
					},
				},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"intList": ["1", "2.74", "3.2", null]
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"intList": []
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"intList": null
					}
				]
			}`),
			wantObjects: []map[string]any{
				{
					"id":                    "4",
					"bestEmail":             "cabbott@efficientoffice.com",
					"dateOfBirth":           "1996-09-02",
					"fullName1":             "Charlotte Abbott",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": true,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
					"intList":               []any{float64(1), float64(2.74), float64(3.2), nil},
				},
				{
					"id":                    "5",
					"bestEmail":             "aadams@efficientoffice.com",
					"dateOfBirth":           "1983-06-30",
					"fullName1":             "Ashley Adams",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": true,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
					"intList":               []any{},
				},
				{
					"id":                    "5",
					"bestEmail":             "aadams@efficientoffice.com",
					"dateOfBirth":           "1983-06-30",
					"fullName1":             "Ashley Adams",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": true,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
					"intList":               nil,
				},
			},
			wantNextCursor: nil,
		},
		"bool_list_conversion_with_mapper": {
			request: &bamboohr.Request{
				AttributeMappings: &bamboohr.AttributeMappings{
					BoolMappings: &bamboohr.BoolAttributeMappings{
						True:  []string{"yes"},
						False: []string{"no"},
					},
				},
				PageSize: 3,
				EntityConfig: &framework.EntityConfig{
					ExternalId: "employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "customcustomBoolField",
							Type:       framework.AttributeTypeBool,
							List:       false,
						},
						{
							ExternalId: "boolList",
							Type:       framework.AttributeTypeBool,
							List:       true,
						},
					},
				},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"boolList": ["1", "true", "no", null]
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"boolList": []
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"boolList": null
					}
				]
			}`),
			wantObjects: []map[string]any{
				{
					"id":                    "4",
					"bestEmail":             "cabbott@efficientoffice.com",
					"dateOfBirth":           "1996-09-02",
					"fullName1":             "Charlotte Abbott",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": true,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
					"boolList":              []any{true, true, false, nil},
				},
				{
					"id":                    "5",
					"bestEmail":             "aadams@efficientoffice.com",
					"dateOfBirth":           "1983-06-30",
					"fullName1":             "Ashley Adams",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": true,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
					"boolList":              []any{},
				},
				{
					"id":                    "5",
					"bestEmail":             "aadams@efficientoffice.com",
					"dateOfBirth":           "1983-06-30",
					"fullName1":             "Ashley Adams",
					"isPhotoUploaded":       "true",
					"customcustomBoolField": true,
					"supervisorEId":         "9",
					"supervisorEmail":       "jcaldwell@efficientoffice.com",
					"lastChanged":           "2024-04-12T19:33:50+00:00",
					"boolList":              nil,
				},
			},
			wantNextCursor: nil,
		},
		"invalid_bool_list_conversion_without_mapper": {
			request: &bamboohr.Request{
				PageSize: 100,
				EntityConfig: &framework.EntityConfig{
					ExternalId: "employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "boolList",
							Type:       framework.AttributeTypeBool,
							List:       true,
						},
					},
				},
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"boolList": ["1", "true", "no"]
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"boolList": []
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"boolList": null
					}
				]
			}`),
			wantErr: &framework.Error{
				Message: "Failed to parse attribute: boolList, as type Bool with value: no.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
			},
			wantNextCursor: nil,
		},
		"invalid_int_list_conversion": {
			request: &bamboohr.Request{
				PageSize: 100,
				EntityConfig: &framework.EntityConfig{
					ExternalId: "employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "boolList",
							Type:       framework.AttributeTypeInt64,
							List:       true,
						},
					},
				},
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"boolList": ["1", "true", "no"]
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"boolList": []
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "yes",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00",
						"boolList": null
					}
				]
			}`),
			wantErr: &framework.Error{
				Message: "Failed to parse attribute: boolList, as type Float with value: true.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
			},
			wantNextCursor: nil,
		},
		"invalid_requested_float_type_conversion": {
			request: &bamboohr.Request{
				PageSize: 100,
				EntityConfig: &framework.EntityConfig{
					ExternalId: "employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "supervisorEmail",
							Type:       framework.AttributeTypeDouble,
							List:       false,
						},
					},
				},
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "isPhotoUploaded",
						"type": "bool",
						"name": "Is employee photo uploaded"
					},
					{
						"id": "customcustomBoolField",
						"type": "checkbox",
						"name": "customBoolField"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": [
					{
						"id": "4",
						"bestEmail": "cabbott@efficientoffice.com",
						"dateOfBirth": "1996-09-02",
						"fullName1": "Charlotte Abbott",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "5",
						"bestEmail": "aadams@efficientoffice.com",
						"dateOfBirth": "1983-06-30",
						"fullName1": "Ashley Adams",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "6",
						"bestEmail": "cagluinda@efficientoffice.com",
						"dateOfBirth": "1996-08-27",
						"fullName1": "Christina Agluinda",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					},
					{
						"id": "7",
						"bestEmail": "sanderson@efficientoffice.com",
						"dateOfBirth": "2000-05-08",
						"fullName1": "Shannon Anderson",
						"isPhotoUploaded": "true",
						"customcustomBoolField": "0",
						"supervisorEId": "9",
						"supervisorEmail": "jcaldwell@efficientoffice.com",
						"lastChanged": "2024-04-12T19:33:50+00:00"
					}
				]
			}`),
			wantErr: &framework.Error{
				Message: "Failed to parse attribute: supervisorEmail, as type Float with value: jcaldwell@efficientoffice.com.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
			},
			wantNextCursor: nil,
		},
		"empty_employees": {
			request: &bamboohr.Request{
				PageSize:          5,
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "employeeNumber",
						"type": "employee_number",
						"name": "Employee #"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": []
			}`),
			wantObjects:    []map[string]any{},
			wantNextCursor: nil,
		},
		"empty_employees_with_cursor": {
			request: &bamboohr.Request{
				PageSize: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(3)),
				},
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "employeeNumber",
						"type": "employee_number",
						"name": "Employee #"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": []
			}`),
			wantErr: &framework.Error{
				Message: "The cursor value: 3, is out of range for number of objects: 0",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"malformed_employees_field": {
			request: &bamboohr.Request{
				PageSize: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(3)),
				},
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "employeeNumber",
						"type": "employee_number",
						"name": "Employee #"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				],
				"employees": {}
			}`),
			wantErr: &framework.Error{
				Message: "Failed to unmarshal the datasource response: json: cannot unmarshal object into Go struct field DatasourceResponse.employees of type []map[string]interface {}.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"missing_employees_field": {
			request: &bamboohr.Request{
				PageSize: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(3)),
				},
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			body: []byte(`{
				"title": "Report",
				"fields": [
					{
						"id": "id",
						"type": "int",
						"name": "EEID"
					},
					{
						"id": "bestEmail",
						"type": "email",
						"name": "Email"
					},
					{
						"id": "dateOfBirth",
						"type": "date",
						"name": "Birth Date"
					},
					{
						"id": "employeeNumber",
						"type": "employee_number",
						"name": "Employee #"
					},
					{
						"id": "fullName1",
						"type": "text",
						"name": "First Name Last Name"
					},
					{
						"id": "supervisorEId",
						"type": "text",
						"name": "Supervisor EID"
					},
					{
						"id": "supervisorEmail",
						"type": "email",
						"name": "Manager's email"
					},
					{
						"id": "lastChanged",
						"type": "timestamp",
						"name": "Last changed"
					}
				]
			}`),
			wantErr: &framework.Error{
				Message: "Failed to parse response body: employees field is missing.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"invalid_object_structure": {
			body: []byte(`{
				[
					{
						"id": "MDEw"
					},
					{
						"id": "MDEw"
					}
				]
			}`),
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: invalid character '[' looking for beginning of object key string.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := bamboohr.ParseResponse(tt.body, tt.request)

			if diff := cmp.Diff(tt.wantObjects, gotObjects); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
			}

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !reflect.DeepEqual(gotNextCursor, tt.wantNextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotNextCursor, tt.wantNextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetEmployeePage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	bamboohrClient := bamboohr.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context      context.Context
		request      *bamboohr.Request
		wantRes      *bamboohr.Response
		wantErr      *framework.Error
		expectedLogs []map[string]any
	}{
		"first_page": {
			context: context.Background(),
			request: &bamboohr.Request{
				BaseURL:               server.URL,
				APIKey:                "apiKey123",
				BasicAuthPassword:     "randomString",
				PageSize:              10,
				APIVersion:            "v1",
				CompanyDomain:         "sgnltestdev",
				OnlyCurrent:           true,
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultEmployeeEntityConfig(),
				AttributeMappings:     &bamboohr.AttributeMappings{},
			},
			expectedLogs: []map[string]any{
				{
					"level":                             "info",
					"msg":                               "Starting datasource request",
					fields.FieldRequestEntityExternalID: "Employee",
					fields.FieldRequestPageSize:         int64(10),
				},
				{
					"level":                             "info",
					"msg":                               "Sending request to datasource",
					fields.FieldRequestEntityExternalID: "Employee",
					fields.FieldRequestPageSize:         int64(10),
					fields.FieldRequestURL:              server.URL + "/sgnltestdev/v1/reports/custom?format=JSON&onlyCurrent=true",
				},
				{
					"level":                             "info",
					"msg":                               "Datasource request completed successfully",
					fields.FieldRequestEntityExternalID: "Employee",
					fields.FieldRequestPageSize:         int64(10),
					fields.FieldResponseStatusCode:      int64(200),
					fields.FieldResponseObjectCount:     int64(10),
					fields.FieldResponseNextCursor: map[string]any{
						"cursor": int64(10),
					},
				},
			},
			wantRes: &bamboohr.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                    float64(4),
						"bestEmail":             "cabbott@efficientoffice.com",
						"dateOfBirth":           "1996-09-02",
						"fullName1":             "Charlotte Abbott",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(9),
						"supervisorEmail":       "jcaldwell@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:50+00:00",
					},
					{
						"id":                    float64(5),
						"bestEmail":             "aadams@efficientoffice.com",
						"dateOfBirth":           "1983-06-30",
						"fullName1":             "Ashley Adams",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(9),
						"supervisorEmail":       "jcaldwell@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:50+00:00",
					},
					{
						"id":                    float64(6),
						"bestEmail":             "cagluinda@efficientoffice.com",
						"dateOfBirth":           "1996-08-27",
						"fullName1":             "Christina Agluinda",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(9),
						"supervisorEmail":       "jcaldwell@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:50+00:00",
					},
					{
						"id":                    float64(7),
						"bestEmail":             "sanderson@efficientoffice.com",
						"dateOfBirth":           nil,
						"fullName1":             "Shannon Anderson",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(9),
						"supervisorEmail":       "jcaldwell@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:50+00:00",
					},
					{
						"id":                    float64(8),
						"bestEmail":             "arvind@sgnl.ai",
						"dateOfBirth":           nil,
						"fullName1":             "Arvind",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         nil,
						"supervisorEmail":       nil,
						"lastChanged":           "2024-04-12T19:33:50+00:00",
					},
					{
						"id":                    float64(9),
						"bestEmail":             "jcaldwell@efficientoffice.com",
						"dateOfBirth":           "1975-01-26",
						"fullName1":             "Jennifer Caldwell",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(8),
						"supervisorEmail":       "arvind@sgnl.ai",
						"lastChanged":           nil,
					},
					{
						"id":                    float64(10),
						"bestEmail":             "rsaito@efficientoffice.com",
						"dateOfBirth":           "1968-12-28",
						"fullName1":             "Ryota Saito",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(8),
						"supervisorEmail":       "arvind@sgnl.ai",
						"lastChanged":           nil,
					},
					{
						"id":                    float64(11),
						"bestEmail":             "dvance@efficientoffice.com",
						"dateOfBirth":           "1978-08-23",
						"fullName1":             "Daniel Vance",
						"isPhotoUploaded":       nil,
						"customcustomBoolField": false,
						"supervisorEId":         float64(8),
						"supervisorEmail":       "arvind@sgnl.ai",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(12),
						"bestEmail":             "easture@efficientoffice.com",
						"dateOfBirth":           "1990-07-01",
						"fullName1":             "Eric Asture",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         nil,
						"supervisorEmail":       "arvind@sgnl.ai",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(13),
						"bestEmail":             "cbarnet@efficientoffice.com",
						"dateOfBirth":           "1987-06-16",
						"fullName1":             "Cheryl Barnet",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         nil,
						"supervisorEmail":       "arvind@sgnl.ai",
						"lastChanged":           "2024-04-12T19:33:50+00:00",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(10)),
				},
			},
			wantErr: nil,
		},
		"second_page": {
			context: context.Background(),
			request: &bamboohr.Request{
				BaseURL:               server.URL,
				APIKey:                "apiKey123",
				BasicAuthPassword:     "randomString",
				PageSize:              10,
				APIVersion:            "v1",
				CompanyDomain:         "sgnltestdev",
				OnlyCurrent:           true,
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultEmployeeEntityConfig(),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(10)),
				},
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			wantRes: &bamboohr.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                    float64(14),
						"bestEmail":             "mandev@efficientoffice.com",
						"dateOfBirth":           "1987-06-05",
						"fullName1":             "Maja Andev",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(8),
						"supervisorEmail":       "arvind@sgnl.ai",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(15),
						"bestEmail":             "twalsh@efficientoffice.com",
						"dateOfBirth":           "1981-03-18",
						"fullName1":             "Trent Walsh",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(8),
						"supervisorEmail":       "arvind@sgnl.ai",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(16),
						"bestEmail":             "jbryan@efficientoffice.com",
						"dateOfBirth":           "1970-12-07",
						"fullName1":             "Jake Bryan",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(8),
						"supervisorEmail":       "arvind@sgnl.ai",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(17),
						"bestEmail":             "dchou@efficientoffice.com",
						"dateOfBirth":           "1987-05-08",
						"fullName1":             "Dorothy Chou",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(8),
						"supervisorEmail":       "arvind@sgnl.ai",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(18),
						"bestEmail":             "javier@efficientoffice.com",
						"dateOfBirth":           "1996-08-28",
						"fullName1":             "Javier Cruz",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(15),
						"supervisorEmail":       "twalsh@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(19),
						"bestEmail":             "shelly@efficientoffice.com",
						"dateOfBirth":           "1993-06-01",
						"fullName1":             "Shelly Cluff",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(10),
						"supervisorEmail":       "rsaito@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(22),
						"bestEmail":             "dillon@efficientoffice.com",
						"dateOfBirth":           "1972-06-06",
						"fullName1":             "Dillon (Remote) Park",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(105),
						"supervisorEmail":       "norma@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(23),
						"bestEmail":             "darlene@efficientoffice.com",
						"dateOfBirth":           "1975-09-16",
						"fullName1":             "Darlene Handley",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(15),
						"supervisorEmail":       "twalsh@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(24),
						"bestEmail":             "zack@efficientoffice.com",
						"dateOfBirth":           "2000-08-02",
						"fullName1":             "Zack Miller",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(11),
						"supervisorEmail":       "dvance@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(25),
						"bestEmail":             "philip@efficientoffice.com",
						"dateOfBirth":           "1975-11-26",
						"fullName1":             "Philip Wagener",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(12),
						"supervisorEmail":       "easture@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(20)),
				},
			},
			wantErr: nil,
		},
		"third_page": {
			context: context.Background(),
			request: &bamboohr.Request{
				BaseURL:               server.URL,
				APIKey:                "apiKey123",
				BasicAuthPassword:     "randomString",
				PageSize:              10,
				APIVersion:            "v1",
				CompanyDomain:         "sgnltestdev",
				OnlyCurrent:           true,
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultEmployeeEntityConfig(),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(20)),
				},
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			wantRes: &bamboohr.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                    float64(26),
						"bestEmail":             "agranger@efficientoffice.com",
						"dateOfBirth":           "1998-11-26",
						"fullName1":             "Amy Granger",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(13),
						"supervisorEmail":       "cbarnet@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(27),
						"bestEmail":             "debra@efficientoffice.com",
						"dateOfBirth":           "1966-10-18",
						"fullName1":             "Debra Tuescher",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(49),
						"supervisorEmail":       "robert@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(28),
						"bestEmail":             "andy@efficientoffice.com",
						"dateOfBirth":           "2001-02-25",
						"fullName1":             "Andy Graves",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(26),
						"supervisorEmail":       "agranger@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(29),
						"bestEmail":             "catherine@efficientoffice.com",
						"dateOfBirth":           "1993-12-18",
						"fullName1":             "Catherine Jones",
						"isPhotoUploaded":       false,
						"customcustomBoolField": false,
						"supervisorEId":         float64(4),
						"supervisorEmail":       "cabbott@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:47+00:00",
					},
					{
						"id":                    float64(30),
						"bestEmail":             "corey@efficientoffice.com",
						"dateOfBirth":           "1995-05-01",
						"fullName1":             "Corey Ross",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(15),
						"supervisorEmail":       "twalsh@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(31),
						"bestEmail":             "sally@efficientoffice.com",
						"dateOfBirth":           "1984-05-31",
						"fullName1":             "Sally Harmon",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(16),
						"supervisorEmail":       "jbryan@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(32),
						"bestEmail":             "carly@efficientoffice.com",
						"dateOfBirth":           "1984-07-02",
						"fullName1":             "Carly Seymour",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(52),
						"supervisorEmail":       "nate@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
					{
						"id":                    float64(33),
						"bestEmail":             "erin@efficientoffice.com",
						"dateOfBirth":           "1993-02-26",
						"fullName1":             "Erin Farr",
						"isPhotoUploaded":       false,
						"customcustomBoolField": false,
						"supervisorEId":         float64(16),
						"supervisorEmail":       "jbryan@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:47+00:00",
					},
					{
						"id":                    float64(34),
						"bestEmail":             "emily@efficientoffice.com",
						"dateOfBirth":           "1994-04-30",
						"fullName1":             "Emily Gomez",
						"isPhotoUploaded":       false,
						"customcustomBoolField": false,
						"supervisorEId":         float64(36),
						"supervisorEmail":       "melany@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:47+00:00",
					},
					{
						"id":                    float64(35),
						"bestEmail":             "aaron@efficientoffice.com",
						"dateOfBirth":           "1998-08-16",
						"fullName1":             "Aaron Eckerly",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(36),
						"supervisorEmail":       "melany@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:48+00:00",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(30)),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &bamboohr.Request{
				BaseURL:               server.URL,
				APIKey:                "apiKey123",
				BasicAuthPassword:     "randomString",
				PageSize:              10,
				APIVersion:            "v1",
				CompanyDomain:         "sgnltestdev",
				OnlyCurrent:           true,
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultEmployeeEntityConfig(),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(30)),
				},
				AttributeMappings: &bamboohr.AttributeMappings{},
			},
			wantRes: &bamboohr.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                    float64(36),
						"bestEmail":             "melany@efficientoffice.com",
						"dateOfBirth":           "1986-11-25",
						"fullName1":             "Melany Olsen",
						"isPhotoUploaded":       false,
						"customcustomBoolField": false,
						"supervisorEId":         float64(13),
						"supervisorEmail":       "cbarnet@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:47+00:00",
					},
					{
						"id":                    float64(37),
						"bestEmail":             "whitney@efficientoffice.com",
						"dateOfBirth":           "1992-12-30",
						"fullName1":             "Whitney Webster",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(5),
						"supervisorEmail":       "aadams@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:47+00:00",
					},
					{
						"id":                    float64(38),
						"bestEmail":             "marrissa@efficientoffice.com",
						"dateOfBirth":           "1995-01-31",
						"fullName1":             "Marrissa Mellon",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(69),
						"supervisorEmail":       "karin@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:47+00:00",
					},
					{
						"id":                    float64(39),
						"bestEmail":             "paige@efficientoffice.com",
						"dateOfBirth":           "1993-02-02",
						"fullName1":             "Paige Rasmussen",
						"isPhotoUploaded":       true,
						"customcustomBoolField": false,
						"supervisorEId":         float64(57),
						"supervisorEmail":       "liam@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:47+00:00",
					},
					{
						"id":                    float64(40),
						"bestEmail":             "kelli@efficientoffice.com",
						"dateOfBirth":           "1988-03-01",
						"fullName1":             "Kelli Crandle",
						"isPhotoUploaded":       false,
						"customcustomBoolField": false,
						"supervisorEId":         float64(49),
						"supervisorEmail":       "robert@efficientoffice.com",
						"lastChanged":           "2024-04-12T19:33:47+00:00",
					},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctxWithLogger, observedLogs := testutil.NewContextWithObservableLogger(t.Context())

			gotRes, gotErr := bamboohrClient.GetPage(ctxWithLogger, tt.request)

			if diff := cmp.Diff(gotRes.Objects, tt.wantRes.Objects); diff != "" {
				t.Errorf("Differences found: (-got +want)\n%s", diff)
			}

			if !reflect.DeepEqual(gotRes.Objects, tt.wantRes.Objects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotRes.Objects, tt.wantRes.Objects)
			}

			if !reflect.DeepEqual(gotRes.NextCursor, tt.wantRes.NextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotRes.NextCursor, tt.wantRes.NextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}

			testutil.ValidateLogOutput(t, observedLogs, tt.expectedLogs)
		})
	}
}
