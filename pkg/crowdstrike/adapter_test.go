// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst, lll
package crowdstrike_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	crowdstrike_adapter "github.com/sgnl-ai/adapters/pkg/crowdstrike"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestAdapterUserGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestGraphQLServerHandler)
	adapter := crowdstrike_adapter.NewAdapter(&crowdstrike_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		request            *framework.Request[crowdstrike_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
		expectedLogs       []map[string]any
	}{
		"first_page": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateUserEntityConfig(),
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"creationTime": time.Date(2024, 5, 15, 15, 29, 10, 0, time.UTC),
							"entityId":     string("095b6929-44b9-4525-a0cc-9ef4552011f3"),
							"inactive":     bool(true),

							// Child Objects
							`$.accounts[?(@.__typename=="ActiveDirectoryAccountDescriptor")]`: []framework.Object{
								{
									"objectGuid":     string("095b6929-44b9-4525-a0cc-9ef4552011f3"),
									"samAccountName": string("Wendolyn.Garber"),
								},
							},
						},
						{
							"creationTime": time.Date(2024, 8, 25, 18, 4, 22, 0, time.UTC),
							"entityId":     string("45dc40e2-7b7b-4f38-9ac7-98f4a35b24e1"),
							"inactive":     bool(true),

							// Child Objects
							`$.accounts[?(@.__typename=="ActiveDirectoryAccountDescriptor")]`: []framework.Object{
								{
									"objectGuid":     string("45dc40e2-7b7b-4f38-9ac7-98f4a35b24e1"),
									"samAccountName": string("sgnl-user"),
								},
							},
						},
					},
					NextCursor: "eyJyaXNrU2NvcmUiOjAuNjQ1NDg3MTMzOTk5OTk5OSwiX2lkIjoiNDVkYzQwZTItN2I3Yi00ZjM4LTlhYzctOThmNGEzNWIyNGUxIn0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("eyJyaXNrU2NvcmUiOjAuNjQ1NDg3MTMzOTk5OTk5OSwiX2lkIjoiNDVkYzQwZTItN2I3Yi00ZjM4LTlhYzctOThmNGEzNWIyNGUxIn0="),
			},
			expectedLogs: []map[string]any{
				{
					"level":                             "info",
					"msg":                               "Starting datasource request",
					fields.FieldRequestEntityExternalID: "user",
					fields.FieldRequestPageSize:         int64(2),
				},
				{
					"level":                             "info",
					"msg":                               "Sending HTTP request to datasource",
					fields.FieldRequestEntityExternalID: "user",
					fields.FieldRequestPageSize:         int64(2),
					fields.FieldURL:                     server.URL + "/identity-protection/combined/graphql/v1",
				},
				{
					"level":                             "info",
					"msg":                               "Datasource request completed successfully",
					fields.FieldRequestEntityExternalID: "user",
					fields.FieldRequestPageSize:         int64(2),
					fields.FieldResponseStatusCode:      int64(200),
					fields.FieldResponseObjectCount:     int64(2),
					fields.FieldResponseNextCursor: map[string]any{
						"cursor": "eyJyaXNrU2NvcmUiOjAuNjQ1NDg3MTMzOTk5OTk5OSwiX2lkIjoiNDVkYzQwZTItN2I3Yi00ZjM4LTlhYzctOThmNGEzNWIyNGUxIn0=",
					},
				},
			},
		},
		"middle_page": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateUserEntityConfig(),
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("eyJyaXNrU2NvcmUiOjAuNjQ1NDg3MTMzOTk5OTk5OSwiX2lkIjoiNDVkYzQwZTItN2I3Yi00ZjM4LTlhYzctOThmNGEzNWIyNGUxIn0="),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"creationTime": time.Date(2024, 5, 15, 15, 16, 27, 0, time.UTC),
							"entityId":     string("c1732de2-853c-4375-a479-17b0afbe114f"),
							"inactive":     bool(false),

							// Child Objects
							`$.accounts[?(@.__typename=="ActiveDirectoryAccountDescriptor")]`: []framework.Object{
								{
									"objectGuid":     string("c1732de2-853c-4375-a479-17b0afbe114f"),
									"samAccountName": string("marc"),
								},
							},
						},
						{
							"creationTime": time.Date(2024, 8, 25, 18, 18, 0, 0, time.UTC),
							"entityId":     string("83a49ef1-17a7-4fa4-b90f-9142dfa49577"),
							"inactive":     bool(false),

							// Child Objects
							`$.accounts[?(@.__typename=="ActiveDirectoryAccountDescriptor")]`: []framework.Object{
								{
									"objectGuid":     string("83a49ef1-17a7-4fa4-b90f-9142dfa49577"),
									"samAccountName": string("sgnl.sor"),
								},
							},
						},
					},
					NextCursor: "eyJyaXNrU2NvcmUiOjAuNjQwNDc5MTcxNzM1MjQ4OSwiX2lkIjoiODNhNDllZjEtMTdhNy00ZmE0LWI5MGYtOTE0MmRmYTQ5NTc3In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("eyJyaXNrU2NvcmUiOjAuNjQwNDc5MTcxNzM1MjQ4OSwiX2lkIjoiODNhNDllZjEtMTdhNy00ZmE0LWI5MGYtOTE0MmRmYTQ5NTc3In0="),
			},
		},
		"last_page": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateUserEntityConfig(),
				PageSize: 2,
			},
			// Request last page of users.
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("eyJyaXNrU2NvcmUiOjAuNjQwNDc5MTcxNzM1MjQ4OSwiX2lkIjoiODNhNDllZjEtMTdhNy00ZmE0LWI5MGYtOTE0MmRmYTQ5NTc3In0="),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"creationTime": time.Date(2024, 5, 23, 15, 8, 11, 0, time.UTC),
							"entityId":     string("6b4c76ba-2493-4a87-bfb3-1ea91985cce5"),
							"inactive":     bool(true),

							// Child Objects
							`$.accounts[?(@.__typename=="ActiveDirectoryAccountDescriptor")]`: []framework.Object{
								{
									"objectGuid":     string("6b4c76ba-2493-4a87-bfb3-1ea91985cce5"),
									"samAccountName": string("alejandro.bacong"),
								},
							},
						},
					},
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

			ctxWithLogger, observedLogs := testutil.NewContextWithObservableLogger(t.Context())

			gotResponse := adapter.GetPage(ctxWithLogger, tt.request)
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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
					t.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(*tt.wantCursor, gotCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}

			testutil.ValidateLogOutput(t, observedLogs, tt.expectedLogs)
		})
	}
}

func TestAdapterEndpointGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestGraphQLServerHandler)
	adapter := crowdstrike_adapter.NewAdapter(&crowdstrike_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		request            *framework.Request[crowdstrike_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateEndpointEntityConfig(),
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"creationTime": time.Date(2024, 5, 29, 21, 30, 17, 0, time.UTC),
							"entityId":     string("89be47c3-f51b-48af-884a-ecb02ed0807a"),
							"inactive":     bool(true),

							// Child Objects
							`$.accounts[?(@.__typename=="ActiveDirectoryAccountDescriptor")]`: []framework.Object{
								{
									"objectGuid":     string("89be47c3-f51b-48af-884a-ecb02ed0807a"),
									"samAccountName": string("ALICE-WIN11$"),
								},
							},
						},
						{
							"creationTime": time.Date(2024, 5, 15, 15, 17, 19, 0, time.UTC),
							"entityId":     string("3c7aebb9-411b-4ee9-b481-e881f29afcc8"),
							"inactive":     bool(false),

							// Child Objects
							`$.accounts[?(@.__typename=="ActiveDirectoryAccountDescriptor")]`: []framework.Object{
								{
									"objectGuid":     string("3c7aebb9-411b-4ee9-b481-e881f29afcc8"),
									"samAccountName": string("mj-dc$"),
								},
							},
						},
					},
					NextCursor: "eyJyaXNrU2NvcmUiOjAuNDU5NDAwMDAwMDAwMDAwMDMsIl9pZCI6IjNjN2FlYmI5LTQxMWItNGVlOS1iNDgxLWU4ODFmMjlhZmNjOCJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("eyJyaXNrU2NvcmUiOjAuNDU5NDAwMDAwMDAwMDAwMDMsIl9pZCI6IjNjN2FlYmI5LTQxMWItNGVlOS1iNDgxLWU4ODFmMjlhZmNjOCJ9"),
			},
		},
		"last_page": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateEndpointEntityConfig(),
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("eyJyaXNrU2NvcmUiOjAuNDU5NDAwMDAwMDAwMDAwMDMsIl9pZCI6IjNjN2FlYmI5LTQxMWItNGVlOS1iNDgxLWU4ODFmMjlhZmNjOCJ9"),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"creationTime": time.Date(2024, 8, 25, 18, 6, 23, 0, time.UTC),
							"entityId":     string("fd1e0f0b-f1e1-4224-8d60-4f297aa91c29"),
							"inactive":     bool(false),

							// Child Objects
							`$.accounts[?(@.__typename=="ActiveDirectoryAccountDescriptor")]`: []framework.Object{
								{
									"objectGuid":     string("fd1e0f0b-f1e1-4224-8d60-4f297aa91c29"),
									"samAccountName": string("SE-Demo-Active-$"),
								},
							},
						},
					},
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

			gotResponse := adapter.GetPage(context.Background(), tt.request)
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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
					t.Errorf("error unmarshalling cursor: %v", err)
				}
			}
		})
	}
}

func TestAdapterIncidentGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestGraphQLServerHandler)
	adapter := crowdstrike_adapter.NewAdapter(&crowdstrike_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		request            *framework.Request[crowdstrike_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateIncidentEntityConfig(),
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"startTime":  time.Date(2024, 9, 23, 13, 0, 21, 0, time.UTC),
							"incidentId": string("INC-16"),
							"severity":   string("INFO"),

							// Child Objects
							`$.compromisedEntities`: []framework.Object{
								{
									"entityId":           string("3c7aebb9-411b-4ee9-b481-e881f29afcc8"),
									"primaryDisplayName": string("mj-dc"),
								},
							},
						},
						{
							"startTime":  time.Date(2024, 9, 20, 1, 49, 27, 0, time.UTC),
							"incidentId": string("INC-15"),
							"severity":   string("INFO"),

							// Child Objects
							`$.compromisedEntities`: []framework.Object{
								{
									"entityId":           string("60ee5bb1-805f-46d2-8f3a-9d7cadc52909"),
									"primaryDisplayName": string("Alice Wu"),
								},
							},
						},
					},
					NextCursor: "eyJlbmRUaW1lIjp7IiRkYXRlIjoiMjAyNC0wOS0yMFQwMTo1NToxMC4yNzRaIn0sInNlcXVlbmNlSWQiOjE1fQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("eyJlbmRUaW1lIjp7IiRkYXRlIjoiMjAyNC0wOS0yMFQwMTo1NToxMC4yNzRaIn0sInNlcXVlbmNlSWQiOjE1fQ=="),
			},
		},
		"middle_page": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateIncidentEntityConfig(),
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("eyJlbmRUaW1lIjp7IiRkYXRlIjoiMjAyNC0wOS0yMFQwMTo1NToxMC4yNzRaIn0sInNlcXVlbmNlSWQiOjE1fQ=="),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"startTime":  time.Date(2024, 9, 20, 1, 36, 36, 0, time.UTC),
							"incidentId": string("INC-14"),
							"severity":   string("INFO"),

							// Child Objects
							`$.compromisedEntities`: []framework.Object{
								{
									"entityId":           string("c1732de2-853c-4375-a479-17b0afbe114f"),
									"primaryDisplayName": string("marc"),
								},
							},
						},
						{
							"startTime":  time.Date(2024, 9, 9, 14, 28, 0, 0, time.UTC),
							"incidentId": string("INC-13"),
							"severity":   string("INFO"),

							// Child Objects
							`$.compromisedEntities`: []framework.Object{
								{
									"entityId":           string("3c7aebb9-411b-4ee9-b481-e881f29afcc8"),
									"primaryDisplayName": string("mj-dc"),
								},
							},
						},
					},
					NextCursor: "eyJlbmRUaW1lIjp7IiRkYXRlIjoiMjAyNC0wOS0wOVQxNDoyODowNC4wMDhaIn0sInNlcXVlbmNlSWQiOjEzfQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("eyJlbmRUaW1lIjp7IiRkYXRlIjoiMjAyNC0wOS0wOVQxNDoyODowNC4wMDhaIn0sInNlcXVlbmNlSWQiOjEzfQ=="),
			},
		},
		"last_page": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateIncidentEntityConfig(),
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("eyJlbmRUaW1lIjp7IiRkYXRlIjoiMjAyNC0wOS0wOVQxNDoyODowNC4wMDhaIn0sInNlcXVlbmNlSWQiOjEzfQ=="),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"startTime":  time.Date(2024, 9, 4, 2, 23, 23, 0, time.UTC),
							"incidentId": string("INC-12"),
							"severity":   string("INFO"),

							// Child Objects
							`$.compromisedEntities`: []framework.Object{
								{
									"entityId":           string("83a49ef1-17a7-4fa4-b90f-9142dfa49577"),
									"primaryDisplayName": string("sgnl sor"),
								},
								{
									"entityId":           string("40ff0c2d-a1d3-3676-a924-7688b73c163a"),
									"primaryDisplayName": string("1.1.1.1"),
								},
							},
						},
					},
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

			gotResponse := adapter.GetPage(context.Background(), tt.request)
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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
					t.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(*tt.wantCursor, gotCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterDetectionGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestRESTServerHandler)
	adapter := crowdstrike_adapter.NewAdapter(&crowdstrike_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		request            *framework.Request[crowdstrike_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateDetectionEntityConfig(),
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"detection_id": string("ldt:9b9b1e4f7512492f95f8039c065a28a9:4298086570"),
							"email_sent":   false,
							"status":       string("new"),
						},
						{
							"detection_id": string("ldt:9b9b1e4f7512492f95f8039c065a28a9:4298709414"),
							"email_sent":   false,
							"status":       string("new"),
						},
					},
					NextCursor: "eyJjdXJzb3IiOjJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr[string]("2"),
			},
		},
		"middle_page": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateDetectionEntityConfig(),
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr[string]("2"),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"detection_id": string("ldt:9b9b1e4f7512492f95f8039c065a28a9:1169567"),
							"email_sent":   false,
							"status":       string("new"),
						},
						{
							"detection_id": string("ldt:9b9b1e4f7512492f95f8039c065a28a9:4295459139"),
							"email_sent":   false,
							"status":       string("new"),
						},
					},
					NextCursor: "eyJjdXJzb3IiOjR9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr[string]("4"),
			},
		},
		"last_page": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateDetectionEntityConfig(),
				PageSize: 2,
			},

			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr[string]("4"),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"detection_id": string("ldt:eca21da34c934e8e95c97a4f7af1d9a5:77310702382"),
							"email_sent":   false,
							"status":       string("new"),
						},
						{
							"detection_id": string("ldt:eca21da34c934e8e95c97a4f7af1d9a5:77309428075"),
							"email_sent":   false,
							"status":       string("new"),
						},
					},
				},
			},
		},
		// Non existent page
		"err_404": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateDetectionEntityConfig(),
				PageSize: 2,
			},

			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr[string]("1000"), // Non existent page
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Datasource rejected request, returned status code: 404.",
					Code:    v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		// Specialized error from CRWD APIs
		"err_specialized": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateDetectionEntityConfig(),
				PageSize: 2,
			},

			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr[string]("999"), // Non existent page
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to query the datasource.\n" +
						"Got errors: Code: 404, Message: 404: Page Not Found.",
					Code: v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
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

			gotResponse := adapter.GetPage(context.Background(), tt.request)
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
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
					t.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(*tt.wantCursor, gotCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterAlertGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestRESTServerHandler)
	adapter := crowdstrike_adapter.NewAdapter(&crowdstrike_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		request            *framework.Request[crowdstrike_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateAlertsEntityConfig(),
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"aggregate_id": string("aggind:c36c42b64ce54b39a32e1d57240704c8:625985642613668398"),
							"composite_id": string("8693deb4bf134cfb8855ee118d9a0243:ind:c36c42b64ce54b39a32e1d57240704c8:625985642593750673-20151-7049"),
							"status":       string("new"),
							"$.files_accessed": []framework.Object{
								{
									"filename": string("cat"),
									"filepath": string("/bin/"),
								},
								{
									"filename": string("zshnW4W3l"),
									"filepath": string("/private/tmp/"),
								},
							},
							"$.files_written": []framework.Object{
								{
									"filename": string("eicar.com"),
									"filepath": string("/Users/joe/Desktop/"),
								},
								{
									"filename": string("zshnW4W3l"),
									"filepath": string("/private/tmp/"),
								},
							},
							"$.mitre_attack": []framework.Object{
								{
									"pattern_id": int64(20151),
								},
							},
						},
						{
							"aggregate_id": string("aggind:5388c592189444ad9e84df071c8f3954:8592364792"),
							"composite_id": string("8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:12119912898-10304-117513744"),
							"status":       string("closed"),
							"$.files_accessed": []framework.Object{
								{
									"filename": string("cat"),
									"filepath": string("/bin/"),
								},
								{
									"filename": string("zshnW4W3l"),
									"filepath": string("/private/tmp/"),
								},
							},
							"$.files_written": []framework.Object{
								{
									"filename": string("eicar.com"),
									"filepath": string("/Users/joe/Desktop/"),
								},
								{
									"filename": string("zshnW4W3l"),
									"filepath": string("/private/tmp/"),
								},
							},
							"$.mitre_attack": []framework.Object{
								{
									"pattern_id": int64(10304),
								},
							},
						},
					},
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("eyJ2ZXJzaW9uIjoidjEiLCJ0b3RhbF9oaXRzIjoyMywidG90YWxfcmVsYXRpb24iOiJlcSIsImNsdXN0ZXJfaWQiOiJ0ZXN0IiwiYWZ0ZXIiOlsxNzQ5NjExMTU3MjIxLCJ0ZXN0aWQ6aW5kOjUzODhjNTkyMTg5NDQ0YWQ5ZTg0ZGYwNzFjOGYzOTU0Ojk3ODI3ODI2MTQtMTAzMDMtMzE4MzE1NjgiXSwidG90YWxfZmV0Y2hlZCI6Mn0="),
			},
		},
		"middle_page": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateAlertsEntityConfig(),
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("eyJ2ZXJzaW9uIjoidjEiLCJ0b3RhbF9oaXRzIjoyMywidG90YWxfcmVsYXRpb24iOiJlcSIsImNsdXN0ZXJfaWQiOiJ0ZXN0IiwiYWZ0ZXIiOlsxNzQ5NjExMTU3MjIxLCJ0ZXN0aWQ6aW5kOjUzODhjNTkyMTg5NDQ0YWQ5ZTg0ZGYwNzFjOGYzOTU0Ojk3ODI3ODI2MTQtMTAzMDMtMzE4MzE1NjgiXSwidG90YWxfZmV0Y2hlZCI6Mn0="),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"aggregate_id": string("aggind:5388c592189444ad9e84df071c8f3954:8592364792"),
							"composite_id": string("8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:10653769300-10304-81908752"),
							"status":       string("closed"),
							"$.mitre_attack": []framework.Object{
								{
									"pattern_id": int64(10304),
								},
							},
						},
						{
							"aggregate_id": string("aggind:5388c592189444ad9e84df071c8f3954:8592364792"),
							"composite_id": string("8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:10557972040-10304-81543952"),
							"status":       string("closed"),
							"$.mitre_attack": []framework.Object{
								{
									"pattern_id": int64(10304),
								},
							},
						},
					},
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("eyJ2ZXJzaW9uIjoidjEiLCJ0b3RhbF9oaXRzIjoyMywidG90YWxfcmVsYXRpb24iOiJlcSIsImNsdXN0ZXJfaWQiOiJ0ZXN0IiwiYWZ0ZXIiOlsxNzQ5NTEyMzQ1Njc4LCJ0ZXN0aWQ6aW5kOmU0NTY3ODkwMTIzNDU2YWI3ODkwY2RlZjEyMzQ1Njc4LTIwMTUzLTcwNTEiXSwidG90YWxfZmV0Y2hlZCI6NH0="),
			},
		},
		"last_page": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateAlertsEntityConfig(),
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("eyJ2ZXJzaW9uIjoidjEiLCJ0b3RhbF9oaXRzIjoyMywidG90YWxfcmVsYXRpb24iOiJlcSIsImNsdXN0ZXJfaWQiOiJ0ZXN0IiwiYWZ0ZXIiOlsxNzQ5NTEyMzQ1Njc4LCJ0ZXN0aWQ6aW5kOmU0NTY3ODkwMTIzNDU2Nzg5MGNkZWYxMjM0NTY3ODkwLTIwMTUzLTcwNTEiXSwidG90YWxfZmV0Y2hlZCI6NH0="),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"aggregate_id": string("aggind:5388c592189444ad9e84df071c8f3954:8591071260"),
							"composite_id": string("8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:10230629714-10303-52340240"),
							"status":       string("closed"),
							"$.mitre_attack": []framework.Object{
								{
									"pattern_id": int64(10303),
								},
							},
						},
						{
							"aggregate_id": string("aggind:5388c592189444ad9e84df071c8f3954:8591071260"),
							"composite_id": string("8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:10208107226-10304-52110864"),
							"status":       string("closed"),
							"$.mitre_attack": []framework.Object{
								{
									"pattern_id": int64(10304),
								},
							},
						},
					},
				},
			},
		},
		// Non existent page
		"err_404": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateAlertsEntityConfig(),
				PageSize: 2,
			},

			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("1000"), // Non existent page
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Datasource rejected request, returned status code: 404.",
					Code:    v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		// Specialized error from CRWD APIs
		"err_specialized": {
			request: &framework.Request[crowdstrike_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &crowdstrike_adapter.Config{
					APIVersion: "v1",
					Archived:   false,
					Enabled:    true,
				},
				Entity:   *PopulateAlertsEntityConfig(),
				PageSize: 2,
			},

			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("999"), // Non existent page
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to query the datasource.\n" +
						"Got errors: Code: 404, Message: 404: Page Not Found.",
					Code: v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
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

			gotResponse := adapter.GetPage(context.Background(), tt.request)

			// Check success response.
			if tt.wantResponse.Error == nil {
				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Fatalf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}

				// Check cursor comparison
				if tt.wantCursor != nil && gotResponse.Success.NextCursor != "" {
					gotCursor, err := pagination.UnmarshalCursor[string](gotResponse.Success.NextCursor)
					if err != nil {
						t.Fatalf("error unmarshalling cursor: %v", err)
					}

					if gotCursor.Cursor == nil || tt.wantCursor.Cursor == nil {
						t.Fatalf("gotCursor or wantCursor is nil: gotCursor: %v, wantCursor: %v", gotCursor, *tt.wantCursor)
					} else if *gotCursor.Cursor != *tt.wantCursor.Cursor {
						t.Fatalf("gotCursor: %s, wantCursor: %s", *gotCursor.Cursor, *tt.wantCursor.Cursor)
					}
				}
			}

			// Check error respose.
			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Fatalf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
			}
		})
	}
}
