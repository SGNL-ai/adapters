// Copyright 2026 SGNL.ai, Inc.

// nolint: lll, goconst
package workday_test

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
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
	"github.com/sgnl-ai/adapters/pkg/workday"
)

func TestParseResponse(t *testing.T) {
	tests := map[string]struct {
		request        *workday.Request
		body           []byte
		endpoint       string
		wantObjects    []map[string]any
		wantNextCursor *pagination.CompositeCursor[int64]
		wantErr        *framework.Error
	}{
		"empty_response": {
			request: &workday.Request{
				PageSize: 5,
			},
			body: []byte(`{}`),
			wantErr: &framework.Error{
				Message: "Total count is missing in the datasource response.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"empty_workers": {
			request: &workday.Request{
				PageSize: 5,
			},
			body: []byte(`{
				"total": 0,
				"data": []
			}`),
			wantObjects:    []map[string]any{},
			wantNextCursor: nil,
		},
		"empty_workers_with_cursor": {
			request: &workday.Request{
				PageSize: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(3)),
				},
			},
			body: []byte(`{
				"total": 0,
				"data": []
			}`),
			wantObjects:    []map[string]any{},
			wantNextCursor: nil,
		},
		"malformed_data_field": {
			request: &workday.Request{
				PageSize: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(3)),
				},
			},
			body: []byte(`{
				"total": 0,
				"data": {
					"id": "1"
				}
			}`),
			wantErr: &framework.Error{
				Message: "Failed to unmarshal the datasource response: json: cannot unmarshal object into Go struct field DatasourceResponse.data of type []map[string]interface {}.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"missing_data_field": {
			request: &workday.Request{
				PageSize: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(3)),
				},
			},
			body: []byte(`{
				"total": 100
			}`),
			wantErr: &framework.Error{
				Message: "Missing data in the datasource response.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"missing_total_field": {
			request: &workday.Request{
				PageSize: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(3)),
				},
			},
			body: []byte(`{
				"data": []
			}`),
			wantErr: &framework.Error{
				Message: "Total count is missing in the datasource response.",
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
		"parsing_single_error_message": {
			body: []byte(`{
				"error": "invalid request: WQL error.",
				"errors": [
					{
						"error": "Invalid WQL syntax.",
						"field": "at or near 'ASC'",
						"location": "Invalid ORDER BY clause starting with character: 149"
					}
				]
			}`),
			endpoint: "https://test-instance.workday.com/api/wql/v1/capgeminisrv_dpt4/data?query=SELECT+FTE%2C+company%2C+employeeID%2C+email_Work%2C+gender%2C+hireDate%2C+jobTitle%2C+managementLevel%2C+positionID%2C+worker%2C+orkerActive+FROM+allWorkers+ORDER+BY++ASC&limit=1840&offset=0",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to query the datasource: https://test-instance.workday.com/api/wql/v1/capgeminisrv_dpt4/data?query=SELECT+FTE%2C+company%2C+employeeID%2C+email_Work%2C+gender%2C+hireDate%2C+jobTitle%2C+managementLevel%2C+positionID%2C+worker%2C+orkerActive+FROM+allWorkers+ORDER+BY++ASC&limit=1840&offset=0.\nGot errors: invalid request: WQL error.\nError: Invalid WQL syntax., Field: at or near 'ASC', Location: Invalid ORDER BY clause starting with character: 149.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"parsing_multiple_error_messages": {
			body: []byte(`{
				"error": "invalid request: WQL error.",
				"errors": [
					{
						"error": "Invalid WQL syntax.",
						"field": "at or near 'ASC'",
						"location": "Invalid ORDER BY clause starting with character: 149"
					},
					{
						"error": "Random Error 1",
						"field": "at or near 'DESC'",
						"location": "Invalid SQL KEYWORD"
					}
				]
			}`),
			endpoint: "https://test-instance.workday.com/api/wql/v1/capgeminisrv_dpt4/data?query=SELECT+FTE%2C+company%2C+employeeID%2C+email_Work%2C+gender%2C+hireDate%2C+jobTitle%2C+managementLevel%2C+positionID%2C+worker%2C+orkerActive+FROM+allWorkers+ORDER+BY++ASC&limit=1840&offset=0",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to query the datasource: https://test-instance.workday.com/api/wql/v1/capgeminisrv_dpt4/data?query=SELECT+FTE%2C+company%2C+employeeID%2C+email_Work%2C+gender%2C+hireDate%2C+jobTitle%2C+managementLevel%2C+positionID%2C+worker%2C+orkerActive+FROM+allWorkers+ORDER+BY++ASC&limit=1840&offset=0.\nGot errors: invalid request: WQL error.\nError: Invalid WQL syntax., Field: at or near 'ASC', Location: Invalid ORDER BY clause starting with character: 149\nError: Random Error 1, Field: at or near 'DESC', Location: Invalid SQL KEYWORD.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := workday.ParseResponse(tt.body, tt.request, tt.endpoint)

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

func TestGetWorkerPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	workdayClient := workday.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context      context.Context
		request      *workday.Request
		wantRes      *workday.Response
		wantErr      *framework.Error
		expectedLogs []map[string]any
	}{
		"first_page": {
			context: context.Background(),
			request: &workday.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer testtoken",
				PageSize:              5,
				APIVersion:            "v1",
				OrganizationID:        "SGNL",
				RequestTimeoutSeconds: 20,
				EntityConfig:          PopulateDefaultWorkerEntityConfig(),
			},
			wantRes: &workday.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"worker": map[string]any{
							"descriptor": "user1",
							"id":         "3aa5550b7fe348b98d7b5741afc65534",
						},
						"email_Work":   []any{},
						"employeeID":   "21001",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "4 Vice President",
							"id":         "679d4d1ac6da40e19deb7d91e170431d",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Female",
							"id":         "9cce3bec2d0d420283f76f51b928d885",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00004",
						"jobTitle":   "Vice President, Human Resources",
					},
					{
						"worker": map[string]any{
							"descriptor": "user2",
							"id":         "0e44c92412d34b01ace61e80a47aaf6d",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user2@workdaySJTest.net",
								"id":         "d7fef59db8e21001de457203a69e0001",
							},
						},
						"employeeID":   "21002",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "2 Chief Executive Officer",
							"id":         "3de1f2834f064394a40a40a727fb6c6d",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Not Declared",
							"id":         "a14bf6afa9204ff48a8ea353dd71eb22",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00001",
						"jobTitle":   "Chief Executive Officer",
					},
					{
						"worker": map[string]any{
							"descriptor": "user3",
							"id":         "3895af7993ff4c509cbea2e1817172e0",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user3@workday.net",
								"id":         "d7fef59db8e21001dddaa607a7d30001",
							},
						},
						"employeeID":   "21003",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "3 Executive Vice President",
							"id":         "0ceb3292987b474bbc40c751a1e22c69",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Male",
							"id":         "d3afbf8074e549ffb070962128e1105a",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00002",
						"jobTitle":   "Chief Information Officer",
					},
					{
						"worker": map[string]any{
							"descriptor": "user4",
							"id":         "3bf7df19491f4d039fd54decdd84e05c",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user4@workday.net",
								"id":         "2eab98c6070f4a609adf9ce702bfa9c3",
							},
						},
						"employeeID":   "21004",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "3 Executive Vice President",
							"id":         "0ceb3292987b474bbc40c751a1e22c69",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Male",
							"id":         "d3afbf8074e549ffb070962128e1105a",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00005",
						"jobTitle":   "Chief Operating Officer",
					},
					{
						"worker": map[string]any{
							"descriptor": "user5",
							"id":         "26c439a5deed4a7dbab76709e0d2d2ca",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user5@workday.net",
								"id":         "3aff08c6468b45998638dbbaeaaf4ab8",
							},
						},
						"employeeID":   "21005",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "5 Director",
							"id":         "0b778018b3b44ca3959e498041865645",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Female",
							"id":         "9cce3bec2d0d420283f76f51b928d885",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00124",
						"jobTitle":   "Director, Field Marketing",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(5)),
				},
			},
			wantErr: nil,
			expectedLogs: []map[string]any{
				{
					"level":                             "info",
					"msg":                               "Starting datasource request",
					fields.FieldRequestEntityExternalID: "allWorkers",
					fields.FieldRequestPageSize:         int64(5),
				},
				{
					"level":                             "info",
					"msg":                               "Sending request to datasource",
					fields.FieldRequestEntityExternalID: "allWorkers",
					fields.FieldRequestPageSize:         int64(5),
					fields.FieldRequestURL:              server.URL + "/api/wql/v1/SGNL/data?limit=5&offset=0&query=SELECT+FTE%2C+company%2C+email_Work%2C+employeeID%2C+employeeType%2C+gender%2C+hireDate%2C+jobTitle%2C+managementLevel%2C+positionID%2C+worker%2C+workerActive+FROM+allWorkers",
				},
				{
					"level":                             "info",
					"msg":                               "Datasource request completed successfully",
					fields.FieldRequestEntityExternalID: "allWorkers",
					fields.FieldRequestPageSize:         int64(5),
					fields.FieldResponseStatusCode:      int64(200),
					fields.FieldResponseObjectCount:     int64(5),
					fields.FieldResponseNextCursor: map[string]any{
						"cursor": int64(5),
					},
				},
			},
		},
		"second_page": {
			context: context.Background(),
			request: &workday.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer testtoken",
				PageSize:              5,
				APIVersion:            "v1",
				OrganizationID:        "SGNL",
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultWorkerEntityConfig(),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(5)),
				},
			},
			wantRes: &workday.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"worker": map[string]any{
							"descriptor": "user6",
							"id":         "cc7fb31eecd544e9ae8e03653c63bfab",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user6@workday.net",
								"id":         "d7fef59db8e21001de09700cef810002",
							},
						},
						"employeeID":   "21006",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "5 Director",
							"id":         "0b778018b3b44ca3959e498041865645",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Female",
							"id":         "9cce3bec2d0d420283f76f51b928d885",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00011",
						"jobTitle":   "Director, Employee Benefits",
					},
					{
						"worker": map[string]any{
							"descriptor": "user7 (Terminated)",
							"id":         "3a37558d68944bf394fad59ff267f4a1",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user7@workday.net",
								"id":         "4c4aa6815de541bfb24cf6144a0550cc",
							},
						},
						"employeeID":   "21007",
						"workerActive": false,
						"gender": map[string]any{
							"descriptor": "Female",
							"id":         "9cce3bec2d0d420283f76f51b928d885",
						},
						"hireDate": "2000-01-01",
						"FTE":      "0",
					},
					{
						"worker": map[string]any{
							"descriptor": "user8",
							"id":         "3bcc416214054db6911612ef25d51e9f",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user8@workday.net",
								"id":         "1d53eb9c5247461781f6a415bf94ad49",
							},
						},
						"employeeID":   "21008",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "5 Director",
							"id":         "0b778018b3b44ca3959e498041865645",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Not Declared",
							"id":         "a14bf6afa9204ff48a8ea353dd71eb22",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00010",
						"jobTitle":   "Director, Payroll Operations",
					},
					{
						"worker": map[string]any{
							"descriptor": "user9",
							"id":         "d66d21e0b1c949b2b1a3decd2fad1375",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user9@workday.net",
								"id":         "8261477c74b748a2b03482bb9cdb7287",
							},
						},
						"employeeID":   "21009",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "5 Director",
							"id":         "0b778018b3b44ca3959e498041865645",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Male",
							"id":         "d3afbf8074e549ffb070962128e1105a",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00009",
						"jobTitle":   "Director, Workforce Planning",
					},
					{
						"worker": map[string]any{
							"descriptor": "user10",
							"id":         "50ef79568a9b463a9c5fc431e074125b",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user10@workday.net",
								"id":         "60355b860cae4f7ea300e51594b8e610",
							},
						},
						"employeeID":   "21012",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "8 Individual Contributor",
							"id":         "7a379eea3b0c4a10a2b50663b2bd15e4",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Male",
							"id":         "d3afbf8074e549ffb070962128e1105a",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00508",
						"jobTitle":   "Staff Payroll Specialist",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(10)),
				},
			},
			wantErr: nil,
		},
		"third_page": {
			context: context.Background(),
			request: &workday.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer testtoken",
				PageSize:              5,
				APIVersion:            "v1",
				OrganizationID:        "SGNL",
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultWorkerEntityConfig(),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(10)),
				},
			},
			wantRes: &workday.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"worker": map[string]any{
							"descriptor": "user11 (On Leave)",
							"id":         "cf9f717959444023b9bc9226a2556661",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user11@workday.net",
								"id":         "d80ef4c876e04e2fadffca124b944ce4",
							},
						},
						"employeeID":   "21010",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "8 Individual Contributor",
							"id":         "7a379eea3b0c4a10a2b50663b2bd15e4",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Female",
							"id":         "9cce3bec2d0d420283f76f51b928d885",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00503",
						"jobTitle":   "Senior Benefits Analyst",
					},
					{
						"worker": map[string]any{
							"descriptor": "user12",
							"id":         "f21231394b71433c8f75f6fe78264f33",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user12@workday.net",
								"id":         "68a02b2bff3a48afbfc4bd7c89750ee1",
							},
						},
						"employeeID":   "21014",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "8 Individual Contributor",
							"id":         "7a379eea3b0c4a10a2b50663b2bd15e4",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Male",
							"id":         "d3afbf8074e549ffb070962128e1105a",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00515",
						"jobTitle":   "Staff Recruiter",
					},
					{
						"worker": map[string]any{
							"descriptor": "user13",
							"id":         "0a46063523fd469f96d4e81ed4d17812",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user13@workday.net",
								"id":         "6a1ae07ebe754bc19fb624d345fc6a68",
							},
						},
						"employeeID":   "21011",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "8 Individual Contributor",
							"id":         "7a379eea3b0c4a10a2b50663b2bd15e4",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Not Declared",
							"id":         "a14bf6afa9204ff48a8ea353dd71eb22",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00509",
						"jobTitle":   "Staff Payroll Specialist",
					},
					{
						"worker": map[string]any{
							"descriptor": "user14",
							"id":         "cb625aa152344212970023a793f2c2ac",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user14@workday.net",
								"id":         "e26999c7731641b8a1c0f678aae7d385",
							},
						},
						"employeeID":   "21013",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "5 Director",
							"id":         "0b778018b3b44ca3959e498041865645",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Male",
							"id":         "d3afbf8074e549ffb070962128e1105a",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00512",
						"jobTitle":   "Director, Payroll Operations",
					},
					{
						"worker": map[string]any{
							"descriptor": "user15",
							"id":         "2014150640fa42ebbafb6ab936b08073",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user15@workday.net",
								"id":         "06d86d5ac21343c5ac866179d320d27e",
							},
						},
						"employeeID":   "21015",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "8 Individual Contributor",
							"id":         "7a379eea3b0c4a10a2b50663b2bd15e4",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Male",
							"id":         "d3afbf8074e549ffb070962128e1105a",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00517",
						"jobTitle":   "Senior Workforce Analyst",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(15)),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &workday.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer testtoken",
				PageSize:              5,
				APIVersion:            "v1",
				OrganizationID:        "SGNL",
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultWorkerEntityConfig(),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(15)),
				},
			},
			wantRes: &workday.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"worker": map[string]any{
							"descriptor": "user16",
							"id":         "16d87047a76a47b399b4a677058d629f",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user16@workday.net",
								"id":         "8e8aaddd60814dc693b89a938c192cba",
							},
						},
						"employeeID":   "21016",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "8 Individual Contributor",
							"id":         "7a379eea3b0c4a10a2b50663b2bd15e4",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services, Inc. (USA)",
							"id":         "cb550da820584750aae8f807882fa79a",
						},
						"gender": map[string]any{
							"descriptor": "Male",
							"id":         "d3afbf8074e549ffb070962128e1105a",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00502",
						"jobTitle":   "Senior Benefits Analyst",
					},
					{
						"worker": map[string]any{
							"descriptor": "user17",
							"id":         "1cf028c6f4484c248e8d7d573d7b8845",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user17@workday.net",
								"id":         "87dd6909b97e4fc6b2d77769ee1503ac",
							},
						},
						"employeeID":   "21017",
						"workerActive": true,
						"managementLevel": map[string]any{
							"descriptor": "5 Director",
							"id":         "0b778018b3b44ca3959e498041865645",
						},
						"employeeType": []any{
							map[string]any{
								"descriptor": "Regular",
								"id":         "9459f5e6f1084433b767c7901ec04416",
							},
						},
						"company": map[string]any{
							"descriptor": "Global Modern Services S.p.A (Italy)",
							"id":         "e4859d59e6094f52a8f2e865cca82cef",
						},
						"gender": map[string]any{
							"descriptor": "Female",
							"id":         "9cce3bec2d0d420283f76f51b928d885",
						},
						"hireDate":   "2000-01-01",
						"FTE":        "1",
						"positionID": "P-00013",
						"jobTitle":   "Director, Accounting",
					},
					{
						"worker": map[string]any{
							"descriptor": "user18 (Terminated)",
							"id":         "f2c673e5b73245889be3581d53187731",
						},
						"email_Work": []any{
							map[string]any{
								"descriptor": "user18@workday.net",
								"id":         "db497b4fba714fd6b457cc56b821c604",
							},
						},
						"employeeID":   "21018",
						"workerActive": false,
						"gender": map[string]any{
							"descriptor": "Male",
							"id":         "d3afbf8074e549ffb070962128e1105a",
						},
						"hireDate": "2000-01-01",
						"FTE":      "0",
					},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctxWithLogger, observedLogs := testutil.NewContextWithObservableLogger(tt.context)

			gotRes, gotErr := workdayClient.GetPage(ctxWithLogger, tt.request)

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
