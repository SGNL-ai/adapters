// Copyright 2026 SGNL.ai, Inc.

// nolint: lll, goconst
package workday_test

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
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
	"github.com/sgnl-ai/adapters/pkg/workday"
)

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := workday.NewAdapter(&workday.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[workday.Config]
		inputRequestCursor interface{}
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[int64]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[workday.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "SGNL",
				},
				Entity:   *PopulateDefaultWorkerEntityConfig(),
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"$.worker.descriptor":          "user1",
							"$.worker.id":                  "3aa5550b7fe348b98d7b5741afc65534",
							"employeeID":                   "21001",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "4 Vice President",
							"$.managementLevel.id":         "679d4d1ac6da40e19deb7d91e170431d",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Female",
							"$.gender.id":          "9cce3bec2d0d420283f76f51b928d885",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00004",
							"jobTitle":             "Vice President, Human Resources",
						},
						{
							"$.worker.descriptor": "user2",
							"$.worker.id":         "0e44c92412d34b01ace61e80a47aaf6d",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user2@workdaySJTest.net",
									"id":         "d7fef59db8e21001de457203a69e0001",
								},
							},
							"employeeID":                   "21002",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "2 Chief Executive Officer",
							"$.managementLevel.id":         "3de1f2834f064394a40a40a727fb6c6d",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Not Declared",
							"$.gender.id":          "a14bf6afa9204ff48a8ea353dd71eb22",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00001",
							"jobTitle":             "Chief Executive Officer",
						},
						{
							"$.worker.descriptor": "user3",
							"$.worker.id":         "3895af7993ff4c509cbea2e1817172e0",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user3@workday.net",
									"id":         "d7fef59db8e21001dddaa607a7d30001",
								},
							},
							"employeeID":                   "21003",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "3 Executive Vice President",
							"$.managementLevel.id":         "0ceb3292987b474bbc40c751a1e22c69",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Male",
							"$.gender.id":          "d3afbf8074e549ffb070962128e1105a",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00002",
							"jobTitle":             "Chief Information Officer",
						},
						{
							"$.worker.descriptor": "user4",
							"$.worker.id":         "3bf7df19491f4d039fd54decdd84e05c",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user4@workday.net",
									"id":         "2eab98c6070f4a609adf9ce702bfa9c3",
								},
							},
							"employeeID":                   "21004",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "3 Executive Vice President",
							"$.managementLevel.id":         "0ceb3292987b474bbc40c751a1e22c69",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Male",
							"$.gender.id":          "d3afbf8074e549ffb070962128e1105a",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00005",
							"jobTitle":             "Chief Operating Officer",
						},
						{
							"$.worker.descriptor": "user5",
							"$.worker.id":         "26c439a5deed4a7dbab76709e0d2d2ca",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user5@workday.net",
									"id":         "3aff08c6468b45998638dbbaeaaf4ab8",
								},
							},
							"employeeID":                   "21005",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "5 Director",
							"$.managementLevel.id":         "0b778018b3b44ca3959e498041865645",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Female",
							"$.gender.id":          "9cce3bec2d0d420283f76f51b928d885",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00124",
							"jobTitle":             "Director, Field Marketing",
						},
					},
					NextCursor: "eyJjdXJzb3IiOjV9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](5),
			},
		},
		"invalid_request_invalid_api_version": {
			ctx: context.Background(),
			request: &framework.Request[workday.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &workday.Config{
					APIVersion:     "v2",
					OrganizationID: "SGNL",
				},
				Entity:   *PopulateDefaultWorkerEntityConfig(),
				PageSize: 3,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Workday config is invalid: apiVersion is not supported: v2.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"invalid_request_http_prefix": {
			ctx: context.Background(),
			request: &framework.Request[workday.Config]{
				Address: "http://" + strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "SGNL",
				},
				Entity:   *PopulateDefaultWorkerEntityConfig(),
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
			request: &framework.Request[workday.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "SGNL",
				},
				Entity:   *PopulateDefaultWorkerEntityConfig(),
				PageSize: 3,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](-50),
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Cursor value must be greater than 0.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"malformed_composite_cursor_string_type": {
			ctx: context.Background(),
			request: &framework.Request[workday.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "SGNL",
				},
				Entity:   *PopulateDefaultWorkerEntityConfig(),
				PageSize: 3,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr[string]("BROKEN"),
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
			request: &framework.Request[workday.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "SGNL",
				},
				Entity:   *PopulateDefaultWorkerEntityConfig(),
				PageSize: 3,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor:           testutil.GenPtr[int64](50),
				CollectionCursor: testutil.GenPtr[int64](10),
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Cursor must not contain CollectionID or CollectionCursor fields for entity allWorkers.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
		"invalid_request_invalid_url": {
			ctx: context.Background(),
			request: &framework.Request[workday.Config]{
				Address: server.URL + "/invalid",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "SGNL",
				},
				Entity:   *PopulateDefaultWorkerEntityConfig(),
				PageSize: 3,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Datasource rejected request, returned status code: 404.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"valid_special_json_path_requests": {
			ctx: context.Background(),
			request: &framework.Request[workday.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "SGNL",
				},
				Entity:   *PopulateWorkerEntityConfigWithNoChildren(),
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"$.worker.descriptor":          "user1",
							"$.worker.id":                  "3aa5550b7fe348b98d7b5741afc65534",
							"employeeID":                   "21001",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "4 Vice President",
							"$.managementLevel.id":         "679d4d1ac6da40e19deb7d91e170431d",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Female",
							"$.gender.id":          "9cce3bec2d0d420283f76f51b928d885",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00004",
							"jobTitle":             "Vice President, Human Resources",
						},
						{
							"$.worker.descriptor":          "user2",
							"$.worker.id":                  "0e44c92412d34b01ace61e80a47aaf6d",
							"$.email_Work[0].descriptor":   "user2@workdaySJTest.net",
							"$.email_Work[0].id":           "d7fef59db8e21001de457203a69e0001",
							"employeeID":                   "21002",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "2 Chief Executive Officer",
							"$.managementLevel.id":         "3de1f2834f064394a40a40a727fb6c6d",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Not Declared",
							"$.gender.id":          "a14bf6afa9204ff48a8ea353dd71eb22",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00001",
							"jobTitle":             "Chief Executive Officer",
						},
						{
							"$.worker.descriptor":          "user3",
							"$.worker.id":                  "3895af7993ff4c509cbea2e1817172e0",
							"$.email_Work[0].descriptor":   "user3@workday.net",
							"$.email_Work[0].id":           "d7fef59db8e21001dddaa607a7d30001",
							"employeeID":                   "21003",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "3 Executive Vice President",
							"$.managementLevel.id":         "0ceb3292987b474bbc40c751a1e22c69",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Male",
							"$.gender.id":          "d3afbf8074e549ffb070962128e1105a",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00002",
							"jobTitle":             "Chief Information Officer",
						},
						{
							"$.worker.descriptor":          "user4",
							"$.worker.id":                  "3bf7df19491f4d039fd54decdd84e05c",
							"$.email_Work[0].descriptor":   "user4@workday.net",
							"$.email_Work[0].id":           "2eab98c6070f4a609adf9ce702bfa9c3",
							"employeeID":                   "21004",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "3 Executive Vice President",
							"$.managementLevel.id":         "0ceb3292987b474bbc40c751a1e22c69",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Male",
							"$.gender.id":          "d3afbf8074e549ffb070962128e1105a",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00005",
							"jobTitle":             "Chief Operating Officer",
						},
						{
							"$.worker.descriptor":          "user5",
							"$.worker.id":                  "26c439a5deed4a7dbab76709e0d2d2ca",
							"$.email_Work[0].descriptor":   "user5@workday.net",
							"$.email_Work[0].id":           "3aff08c6468b45998638dbbaeaaf4ab8",
							"employeeID":                   "21005",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "5 Director",
							"$.managementLevel.id":         "0b778018b3b44ca3959e498041865645",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Female",
							"$.gender.id":          "9cce3bec2d0d420283f76f51b928d885",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00124",
							"jobTitle":             "Director, Field Marketing",
						},
					},
					NextCursor: "eyJjdXJzb3IiOjV9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](5),
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

func TestAdapterGetWorkerPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := workday.NewAdapter(&workday.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[workday.Config]
		inputRequestCursor *pagination.CompositeCursor[int64]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[int64]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[workday.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "SGNL",
				},
				Entity:   *PopulateDefaultWorkerEntityConfig(),
				PageSize: 5,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"$.worker.descriptor":          "user1",
							"$.worker.id":                  "3aa5550b7fe348b98d7b5741afc65534",
							"employeeID":                   "21001",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "4 Vice President",
							"$.managementLevel.id":         "679d4d1ac6da40e19deb7d91e170431d",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Female",
							"$.gender.id":          "9cce3bec2d0d420283f76f51b928d885",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00004",
							"jobTitle":             "Vice President, Human Resources",
						},
						{
							"$.worker.descriptor": "user2",
							"$.worker.id":         "0e44c92412d34b01ace61e80a47aaf6d",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user2@workdaySJTest.net",
									"id":         "d7fef59db8e21001de457203a69e0001",
								},
							},
							"employeeID":                   "21002",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "2 Chief Executive Officer",
							"$.managementLevel.id":         "3de1f2834f064394a40a40a727fb6c6d",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Not Declared",
							"$.gender.id":          "a14bf6afa9204ff48a8ea353dd71eb22",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00001",
							"jobTitle":             "Chief Executive Officer",
						},
						{
							"$.worker.descriptor": "user3",
							"$.worker.id":         "3895af7993ff4c509cbea2e1817172e0",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user3@workday.net",
									"id":         "d7fef59db8e21001dddaa607a7d30001",
								},
							},
							"employeeID":                   "21003",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "3 Executive Vice President",
							"$.managementLevel.id":         "0ceb3292987b474bbc40c751a1e22c69",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Male",
							"$.gender.id":          "d3afbf8074e549ffb070962128e1105a",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00002",
							"jobTitle":             "Chief Information Officer",
						},
						{
							"$.worker.descriptor": "user4",
							"$.worker.id":         "3bf7df19491f4d039fd54decdd84e05c",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user4@workday.net",
									"id":         "2eab98c6070f4a609adf9ce702bfa9c3",
								},
							},
							"employeeID":                   "21004",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "3 Executive Vice President",
							"$.managementLevel.id":         "0ceb3292987b474bbc40c751a1e22c69",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Male",
							"$.gender.id":          "d3afbf8074e549ffb070962128e1105a",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00005",
							"jobTitle":             "Chief Operating Officer",
						},
						{
							"$.worker.descriptor": "user5",
							"$.worker.id":         "26c439a5deed4a7dbab76709e0d2d2ca",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user5@workday.net",
									"id":         "3aff08c6468b45998638dbbaeaaf4ab8",
								},
							},
							"employeeID":                   "21005",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "5 Director",
							"$.managementLevel.id":         "0b778018b3b44ca3959e498041865645",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Female",
							"$.gender.id":          "9cce3bec2d0d420283f76f51b928d885",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00124",
							"jobTitle":             "Director, Field Marketing",
						},
					},
					NextCursor: "eyJjdXJzb3IiOjV9",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](5),
			},
		},
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[workday.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "SGNL",
				},
				Entity:   *PopulateDefaultWorkerEntityConfig(),
				PageSize: 5,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](5),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"$.worker.descriptor": "user6",
							"$.worker.id":         "cc7fb31eecd544e9ae8e03653c63bfab",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user6@workday.net",
									"id":         "d7fef59db8e21001de09700cef810002",
								},
							},
							"employeeID":                   "21006",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "5 Director",
							"$.managementLevel.id":         "0b778018b3b44ca3959e498041865645",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Female",
							"$.gender.id":          "9cce3bec2d0d420283f76f51b928d885",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00011",
							"jobTitle":             "Director, Employee Benefits",
						},
						{
							"$.worker.descriptor": "user7 (Terminated)",
							"$.worker.id":         "3a37558d68944bf394fad59ff267f4a1",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user7@workday.net",
									"id":         "4c4aa6815de541bfb24cf6144a0550cc",
								},
							},
							"employeeID":          "21007",
							"workerActive":        false,
							"$.gender.descriptor": "Female",
							"$.gender.id":         "9cce3bec2d0d420283f76f51b928d885",
							"hireDate":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                 "0",
						},
						{
							"$.worker.descriptor": "user8",
							"$.worker.id":         "3bcc416214054db6911612ef25d51e9f",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user8@workday.net",
									"id":         "1d53eb9c5247461781f6a415bf94ad49",
								},
							},
							"employeeID":                   "21008",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "5 Director",
							"$.managementLevel.id":         "0b778018b3b44ca3959e498041865645",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Not Declared",
							"$.gender.id":          "a14bf6afa9204ff48a8ea353dd71eb22",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00010",
							"jobTitle":             "Director, Payroll Operations",
						},
						{
							"$.worker.descriptor": "user9",
							"$.worker.id":         "d66d21e0b1c949b2b1a3decd2fad1375",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user9@workday.net",
									"id":         "8261477c74b748a2b03482bb9cdb7287",
								},
							},
							"employeeID":                   "21009",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "5 Director",
							"$.managementLevel.id":         "0b778018b3b44ca3959e498041865645",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Male",
							"$.gender.id":          "d3afbf8074e549ffb070962128e1105a",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00009",
							"jobTitle":             "Director, Workforce Planning",
						},
						{
							"$.worker.descriptor": "user10",
							"$.worker.id":         "50ef79568a9b463a9c5fc431e074125b",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user10@workday.net",
									"id":         "60355b860cae4f7ea300e51594b8e610",
								},
							},
							"employeeID":                   "21012",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "8 Individual Contributor",
							"$.managementLevel.id":         "7a379eea3b0c4a10a2b50663b2bd15e4",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Male",
							"$.gender.id":          "d3afbf8074e549ffb070962128e1105a",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00508",
							"jobTitle":             "Staff Payroll Specialist",
						},
					},
					NextCursor: "eyJjdXJzb3IiOjEwfQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](10),
			},
		},
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[workday.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "SGNL",
				},
				Entity:   *PopulateDefaultWorkerEntityConfig(),
				PageSize: 5,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](10),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"$.worker.descriptor": "user11 (On Leave)",
							"$.worker.id":         "cf9f717959444023b9bc9226a2556661",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user11@workday.net",
									"id":         "d80ef4c876e04e2fadffca124b944ce4",
								},
							},
							"employeeID":                   "21010",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "8 Individual Contributor",
							"$.managementLevel.id":         "7a379eea3b0c4a10a2b50663b2bd15e4",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Female",
							"$.gender.id":          "9cce3bec2d0d420283f76f51b928d885",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00503",
							"jobTitle":             "Senior Benefits Analyst",
						},
						{
							"$.worker.descriptor": "user12",
							"$.worker.id":         "f21231394b71433c8f75f6fe78264f33",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user12@workday.net",
									"id":         "68a02b2bff3a48afbfc4bd7c89750ee1",
								},
							},
							"employeeID":                   "21014",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "8 Individual Contributor",
							"$.managementLevel.id":         "7a379eea3b0c4a10a2b50663b2bd15e4",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Male",
							"$.gender.id":          "d3afbf8074e549ffb070962128e1105a",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00515",
							"jobTitle":             "Staff Recruiter",
						},
						{
							"$.worker.descriptor": "user13",
							"$.worker.id":         "0a46063523fd469f96d4e81ed4d17812",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user13@workday.net",
									"id":         "6a1ae07ebe754bc19fb624d345fc6a68",
								},
							},
							"employeeID":                   "21011",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "8 Individual Contributor",
							"$.managementLevel.id":         "7a379eea3b0c4a10a2b50663b2bd15e4",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Not Declared",
							"$.gender.id":          "a14bf6afa9204ff48a8ea353dd71eb22",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00509",
							"jobTitle":             "Staff Payroll Specialist",
						},
						{
							"$.worker.descriptor": "user14",
							"$.worker.id":         "cb625aa152344212970023a793f2c2ac",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user14@workday.net",
									"id":         "e26999c7731641b8a1c0f678aae7d385",
								},
							},
							"employeeID":                   "21013",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "5 Director",
							"$.managementLevel.id":         "0b778018b3b44ca3959e498041865645",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Male",
							"$.gender.id":          "d3afbf8074e549ffb070962128e1105a",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00512",
							"jobTitle":             "Director, Payroll Operations",
						},
						{
							"$.worker.descriptor": "user15",
							"$.worker.id":         "2014150640fa42ebbafb6ab936b08073",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user15@workday.net",
									"id":         "06d86d5ac21343c5ac866179d320d27e",
								},
							},
							"employeeID":                   "21015",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "8 Individual Contributor",
							"$.managementLevel.id":         "7a379eea3b0c4a10a2b50663b2bd15e4",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Male",
							"$.gender.id":          "d3afbf8074e549ffb070962128e1105a",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00517",
							"jobTitle":             "Senior Workforce Analyst",
						},
					},
					NextCursor: "eyJjdXJzb3IiOjE1fQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](15),
			},
		},
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[workday.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Config: &workday.Config{
					APIVersion:     "v1",
					OrganizationID: "SGNL",
				},
				Entity:   *PopulateDefaultWorkerEntityConfig(),
				PageSize: 5,
			},
			inputRequestCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](15),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"$.worker.descriptor": "user16",
							"$.worker.id":         "16d87047a76a47b399b4a677058d629f",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user16@workday.net",
									"id":         "8e8aaddd60814dc693b89a938c192cba",
								},
							},
							"employeeID":                   "21016",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "8 Individual Contributor",
							"$.managementLevel.id":         "7a379eea3b0c4a10a2b50663b2bd15e4",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services, Inc. (USA)",
							"$.company.id":         "cb550da820584750aae8f807882fa79a",
							"$.gender.descriptor":  "Male",
							"$.gender.id":          "d3afbf8074e549ffb070962128e1105a",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00502",
							"jobTitle":             "Senior Benefits Analyst",
						},
						{
							"$.worker.descriptor": "user17",
							"$.worker.id":         "1cf028c6f4484c248e8d7d573d7b8845",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user17@workday.net",
									"id":         "87dd6909b97e4fc6b2d77769ee1503ac",
								},
							},
							"employeeID":                   "21017",
							"workerActive":                 true,
							"$.managementLevel.descriptor": "5 Director",
							"$.managementLevel.id":         "0b778018b3b44ca3959e498041865645",
							"employeeType": []framework.Object{
								map[string]any{
									"descriptor": "Regular",
									"id":         "9459f5e6f1084433b767c7901ec04416",
								},
							},
							"$.company.descriptor": "Global Modern Services S.p.A (Italy)",
							"$.company.id":         "e4859d59e6094f52a8f2e865cca82cef",
							"$.gender.descriptor":  "Female",
							"$.gender.id":          "9cce3bec2d0d420283f76f51b928d885",
							"hireDate":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                  "1",
							"positionID":           "P-00013",
							"jobTitle":             "Director, Accounting",
						},
						{
							"$.worker.descriptor": "user18 (Terminated)",
							"$.worker.id":         "f2c673e5b73245889be3581d53187731",
							"email_Work": []framework.Object{
								map[string]any{
									"descriptor": "user18@workday.net",
									"id":         "db497b4fba714fd6b457cc56b821c604",
								},
							},
							"employeeID":          "21018",
							"workerActive":        false,
							"$.gender.descriptor": "Male",
							"$.gender.id":         "d3afbf8074e549ffb070962128e1105a",
							"hireDate":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							"FTE":                 "0",
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
