// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package salesforce_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	adapter_api_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/salesforce"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

// Define the endpoints and responses for the mock Salesforce server.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	cases := [][]map[string]string{
		// Page 0
		{
			{
				"Id":         "500Hu000020yLuHIAU",
				"CaseNumber": "00001026",
				"Status":     "Closed",
			},
			{
				"Id":         "500Hu000020yLuMIAU",
				"CaseNumber": "00001027",
				"Status":     "New",
			},
		},
		// Page 1
		{
			{
				"Id":         "500Hu000020yLyEIAU",
				"CaseNumber": "00001031",
				"Status":     "New",
			},
			{
				"Id":         "500Hu000020yLyKIAU",
				"CaseNumber": "00001051",
				"Status":     "New",
			},
		},
	}

	switch r.URL.RequestURI() {
	// Cases Page 1
	case "/services/data/v58.0/query?q=SELECT+Id,CaseNumber,Status+FROM+Case+ORDER+BY+Id+ASC":
		w.Write([]byte(`{
				"totalSize": 2,
				"done": false,
				"nextRecordsUrl": "/services/data/v58.0/query/0r8Hu1lKCluUiC9IMK-200",
				"records": [
					{
						"attributes": {
							"type": "Case",
							"url": "/services/data/v58.0/sobjects/Case/` + cases[0][0]["Id"] + `"
						},
						"Id": "` + cases[0][0]["Id"] + `",
						"CaseNumber": "` + cases[0][0]["CaseNumber"] + `",
						"Status": "` + cases[0][0]["Status"] + `"
					},
					{
						"attributes": {
							"type": "Case",
							"url": "/services/data/v58.0/sobjects/Case/` + cases[0][1]["Id"] + `"
						},
						"Id": "` + cases[0][1]["Id"] + `",
						"CaseNumber": "` + cases[0][1]["CaseNumber"] + `",
						"Status": "` + cases[0][1]["Status"] + `"
					}
				]
			}`))

	// Cases Page 2
	case "/services/data/v58.0/query/0r8Hu1lKCluUiC9IMK-200":
		w.Write([]byte(`{
			"totalSize": 2,
			"done": true,
			"nextRecordsUrl": "",
			"records": [
				{
					"attributes": {
						"type": "Case",
						"url": "/services/data/v58.0/sobjects/Case/` + cases[1][0]["Id"] + `"
					},
					"Id": "` + cases[1][0]["Id"] + `",
					"CaseNumber": "` + cases[1][0]["CaseNumber"] + `",
					"Status": "` + cases[1][0]["Status"] + `"
				},
				{
					"attributes": {
						"type": "Case",
						"url": "/services/data/v58.0/sobjects/Case/` + cases[1][1]["Id"] + `"
					},
					"Id": "` + cases[1][1]["Id"] + `",
					"CaseNumber": "` + cases[1][1]["CaseNumber"] + `",
					"Status": "` + cases[1][1]["Status"] + `"
				}
			]
		}`))

	// Closed Cases Page 1
	case "/services/data/v58.0/query?q=SELECT+Id,CaseNumber,Status+FROM+Case+WHERE+Status+%3D+%27Closed%27+ORDER+BY+Id+ASC":
		w.Write([]byte(`{
				"totalSize": 1,
				"done": true,
				"nextRecordsUrl": "",
				"records": [
					{
						"attributes": {
							"type": "Case",
							"url": "/services/data/v58.0/sobjects/Case/` + cases[0][1]["Id"] + `"
						},
						"Id": "` + cases[0][0]["Id"] + `",
						"CaseNumber": "` + cases[0][0]["CaseNumber"] + `",
						"Status": "` + cases[0][0]["Status"] + `"
					}
				]
			}`))

	// Cases Page 3 - Invalid Datetime Format
	case "/services/data/v58.0/query/0r8Hu1lKClCJd892jd-200":
		w.Write([]byte(`{
			"totalSize": 2,
			"done": true,
			"nextRecordsUrl": "",
			"records": [
				{
					"attributes": {
						"type": "Case",
						"url": "/services/data/v58.0/sobjects/Case/` + cases[0][1]["Id"] + `"
					},
					"Id": "` + cases[1][0]["Id"] + `",
					"CreatedAt": "2021/01/01 00:00:00.000Z",
					"Status": "` + cases[1][0]["Status"] + `"
				}
			]
		}`))

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})

func TestParseResponse(t *testing.T) {
	tests := map[string]struct {
		body             []byte
		entityExternalID string
		cursor           string
		wantObjects      []map[string]interface{}
		wantNextCursor   *string
		wantErr          *framework.Error
	}{
		"single_page_no_next_cursor": {
			body:             []byte(`{"records": [{"Id": "500Hu000020yLuHIAU"}, {"Id": "500Hu000020yLuMIAU"}], "done": true}`),
			entityExternalID: "User",
			wantObjects: []map[string]interface{}{
				{"Id": "500Hu000020yLuHIAU"},
				{"Id": "500Hu000020yLuMIAU"},
			},
		},
		"single_page_next_cursor": {
			body:             []byte(`{"records": [{"Id": "500Hu000020yLuHIAU"}, {"Id": "500Hu000020yLuMIAU"}], "nextRecordsUrl": "/services/data/v58.0/query/0r8Hu1lKCluUiC9IMK-200", "done": true}`),
			entityExternalID: "User",
			wantObjects: []map[string]interface{}{
				{"Id": "500Hu000020yLuHIAU"},
				{"Id": "500Hu000020yLuMIAU"},
			},
			wantNextCursor: testutil.GenPtr("/services/data/v58.0/query/0r8Hu1lKCluUiC9IMK-200"),
		},
		"next_page_cursor_invalid_type": {
			body:             []byte(`{"records": [{"Id": "500Hu000020yLuHIAU"}, {"Id": "500Hu000020yLuMIAU"}], "nextRecordsUrl": 10, "done": true}`),
			entityExternalID: "User",
			wantNextCursor:   nil,
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: json: cannot unmarshal number into Go struct field DatasourceResponse.nextRecordsUrl of type string.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"records_invalid_type": {
			body:             []byte(`{"records": {"Id": "500Hu000020yLuHIAU"}, "done": true}`),
			entityExternalID: "User",
			wantNextCursor:   nil,
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: json: cannot unmarshal object into Go struct field DatasourceResponse.records of type []map[string]interface {}.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"invalid_objects": {
			body:             []byte(`{"records": ["500Hu000020yLuHIAU", "500Hu000020yLuMIAU"], "done": true}`),
			entityExternalID: "User",
			wantNextCursor:   nil,
			wantErr: testutil.GenPtr(framework.Error{
				Message: `Failed to unmarshal the datasource response: json: cannot unmarshal string into Go struct field DatasourceResponse.records of type map[string]interface {}.`,
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"invalid_object_structure": {
			body:             []byte(`{"records": [{"500Hu000020yLuHIAU"}, {"500Hu000020yLuMIAU"}], "done": true}`),
			entityExternalID: "User",
			wantNextCursor:   nil,
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: invalid character '}' after object key.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := salesforce.ParseResponse(tt.body)

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

func TestGetPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	salesforceClient := salesforce.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context      context.Context
		request      *salesforce.Request
		wantRes      *salesforce.Response
		wantErr      *framework.Error
		expectedLogs []map[string]any
	}{
		"first_page": {
			context: context.Background(),
			request: &salesforce.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Bearer testtoken",
				EntityExternalID:      "Case",
				PageSize:              200,
				APIVersion:            "58.0",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "Id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "CaseNumber",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "Status",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			expectedLogs: []map[string]any{
				{
					"level":                             "info",
					"msg":                               "Starting datasource request",
					fields.FieldRequestEntityExternalID: "Case",
					fields.FieldRequestPageSize:         int64(200),
				},
				{
					"level":                             "info",
					"msg":                               "Sending request to datasource",
					fields.FieldRequestEntityExternalID: "Case",
					fields.FieldRequestPageSize:         int64(200),
					fields.FieldRequestURL:              server.URL + "/services/data/v58.0/query?q=SELECT+Id,CaseNumber,Status+FROM+Case+ORDER+BY+Id+ASC",
				},
				{
					"level":                             "info",
					"msg":                               "Datasource request completed successfully",
					fields.FieldRequestEntityExternalID: "Case",
					fields.FieldRequestPageSize:         int64(200),
					fields.FieldResponseStatusCode:      int64(200),
					fields.FieldResponseObjectCount:     int64(2),
					fields.FieldResponseNextCursor:      "/services/data/v58.0/query/0r8Hu1lKCluUiC9IMK-200",
				},
			},
			wantRes: &salesforce.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"Id":         "500Hu000020yLuHIAU",
						"CaseNumber": "00001026",
						"Status":     "Closed",
						"attributes": map[string]any{
							"type": "Case",
							"url":  "/services/data/v58.0/sobjects/Case/500Hu000020yLuHIAU",
						},
					},
					{
						"Id":         "500Hu000020yLuMIAU",
						"CaseNumber": "00001027",
						"Status":     "New",
						"attributes": map[string]any{
							"type": "Case",
							"url":  "/services/data/v58.0/sobjects/Case/500Hu000020yLuMIAU",
						},
					},
				},
				NextCursor: testutil.GenPtr("/services/data/v58.0/query/0r8Hu1lKCluUiC9IMK-200"),
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &salesforce.Request{
				BaseURL:               server.URL,
				RequestTimeoutSeconds: 5,
				Token:                 "Bearer testtoken",
				EntityExternalID:      "Case",
				PageSize:              200,
				APIVersion:            "58.0",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "Id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "CaseNumber",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "Status",
						Type:       framework.AttributeTypeString,
					},
				},
				Cursor: testutil.GenPtr("/services/data/v58.0/query/0r8Hu1lKCluUiC9IMK-200"),
			},
			wantRes: &salesforce.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"Id":         "500Hu000020yLyEIAU",
						"CaseNumber": "00001031",
						"Status":     "New",
						"attributes": map[string]any{
							"type": "Case",
							"url":  "/services/data/v58.0/sobjects/Case/500Hu000020yLyEIAU",
						},
					},
					{
						"Id":         "500Hu000020yLyKIAU",
						"CaseNumber": "00001051",
						"Status":     "New",
						"attributes": map[string]any{
							"type": "Case",
							"url":  "/services/data/v58.0/sobjects/Case/500Hu000020yLyKIAU",
						},
					},
				},
				NextCursor: testutil.GenPtr(""),
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctxWithLogger, observedLogs := testutil.NewContextWithObservableLogger(tt.context)

			gotRes, gotErr := salesforceClient.GetPage(ctxWithLogger, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}

			testutil.ValidateLogOutput(t, observedLogs, tt.expectedLogs)
		})
	}
}
