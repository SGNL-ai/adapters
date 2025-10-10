// Copyright 2025 SGNL.ai, Inc.

// nolint: lll
package smoketests

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	adapter_api_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/smoketests/common"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestAzureADAdapter_User(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/azuread/user")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "graph.microsoft.com",
			Id:      "AzureAD",
			Type:    "AzureAD-1.0.1",
			Config:  []byte(`{"apiVersion":"v1.0","filters":{"User":"startswith(displayName,'N')"}}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "AzureADUser",
			ExternalId: "User",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "displayName",
					ExternalId: "displayName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "manager__id",
					ExternalId: "manager__id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 3,
		Cursor:   "",
	})
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success": {
				"objects": [
					{
						"attributes": [
							{
								"id": "displayName",
								"values": [
									{
										"string_value": "Nancy Barr"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "displayName",
								"values": [
									{
										"string_value": "Nancy Barton"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "displayName",
								"values": [
									{
										"string_value": "Nancy Bean"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "kkkkkkkk-llll-mmmm-nnnn-oooooooooooo"
									}
								]
							},
							{
								"id": "manager__id",
								"values": [
									{
										"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
									}
								]
							}
						]
					}
				],
				"nextCursor": "eyJjdXJzb3IiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20vdjEuMC91c2Vycz8kc2VsZWN0PWlkJTJjZGlzcGxheU5hbWVcdTAwMjYkZXhwYW5kPW1hbmFnZXIoJTI0c2VsZWN0JTNkaWQpXHUwMDI2JHRvcD0zXHUwMDI2JGZpbHRlcj1zdGFydHN3aXRoKGRpc3BsYXlOYW1lJTJjJTI3TiUyNylcdTAwMjYkc2tpcHRva2VuPU5FWFRMSU5LX1RPS0VOX1BMQUNFSE9MREVSIn0="
			}
		}
	`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAzureADAdapter_Group(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/azuread/group")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "graph.microsoft.com",
			Id:      "AzureAD",
			Type:    "AzureAD-1.0.1",
			Config:  []byte(`{"apiVersion":"v1.0"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Example Group 1",
			ExternalId: "Group",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "displayName",
					ExternalId: "displayName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "createdDateTime",
					ExternalId: "createdDateTime",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
				},
				{
					Id:         "groupTypes",
					ExternalId: "groupTypes",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       true,
				},
			},
		},
		PageSize: 3,
		Cursor:   "",
	})
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success": {
				"objects": [
					{
						"attributes": [
							{
								"id": "createdDateTime",
								"values": [
									{
										"datetime_value": {
											"timestamp": "2023-02-21T18:00:53Z",
											"timezone_offset": 0
										}
									}
								]
							},
							{
								"id": "displayName",
								"values": [
									{
										"string_value": "Example Group 1"
									}
								]
							},
							{
								"id": "groupTypes",
								"values": [
									{
										"string_value": "Unified"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "createdDateTime",
								"values": [
									{
										"datetime_value": {
											"timestamp": "2022-12-09T23:32:10Z",
											"timezone_offset": 0
										}
									}
								]
							},
							{
								"id": "displayName",
								"values": [
									{
										"string_value": "Example Group 2"
									}
								]
							},
							{
								"id": "groupTypes",
								"values": []
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "createdDateTime",
								"values": [
									{
										"datetime_value": {
											"timestamp": "2023-02-16T17:46:16Z",
											"timezone_offset": 0
										}
									}
								]
							},
							{
								"id": "displayName",
								"values": [
									{
										"string_value": "Example Group 3"
									}
								]
							},
							{
								"id": "groupTypes",
								"values": []
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "kkkkkkkk-llll-mmmm-nnnn-oooooooooooo"
									}
								]
							}
						]
					}
				],
				"nextCursor": "eyJjdXJzb3IiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20vdjEuMC9ncm91cHM/JHNlbGVjdD1pZCUyY2Rpc3BsYXlOYW1lJTJjY3JlYXRlZERhdGVUaW1lJTJjZ3JvdXBUeXBlc1x1MDAyNiR0b3A9M1x1MDAyNiRza2lwdG9rZW49TkVYVExJTktfVE9LRU5fUExBQ0VIT0xERVIifQ=="
			}
		}
	`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAzureADAdapter_Role(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/azuread/role")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	tests := map[string]struct {
		request  *adapter_api_v1.GetPageRequest
		response []byte
	}{
		"first_page": {
			request: &adapter_api_v1.GetPageRequest{
				Datasource: &adapter_api_v1.DatasourceConfig{
					Auth: &adapter_api_v1.DatasourceAuthCredentials{
						AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer {{OMITTED}}",
						},
					},
					Address: "graph.microsoft.com",
					Id:      "AzureAD",
					Type:    "AzureAD-1.0.1",
					Config:  []byte(`{"apiVersion":"v1.0"}`),
				},
				Entity: &adapter_api_v1.EntityConfig{
					Id:         "AzureADRole",
					ExternalId: "Role",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "id",
							ExternalId: "id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "displayName",
							ExternalId: "displayName",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "description",
							ExternalId: "description",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "roleTemplateId",
							ExternalId: "roleTemplateId",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
					},
				},
				PageSize: 2,
				Cursor:   "",
			},
			response: []byte(`
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Can read basic directory information. Commonly used to grant directory read access to applications and guests."
										}
									],
									"id": "description"
								},
								{
									"values": [
										{
											"string_value": "Directory Readers"
										}
									],
									"id": "displayName"
								},
								{
									"values": [
										{
											"string_value": "0fea7f0d-dea1-458d-9099-69fcc2e3cd42"
										}
									],
									"id": "id"
								},
								{
									"values": [
										{
											"string_value": "88d8e3e3-8f55-4a1e-953a-9b9898b8876b"
										}
									],
									"id": "roleTemplateId"
								}
							],
							"child_objects": []
						},
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Can manage all aspects of Microsoft Entra ID and Microsoft services that use Microsoft Entra identities."
										}
									],
									"id": "description"
								},
								{
									"values": [
										{
											"string_value": "Global Administrator"
										}
									],
									"id": "displayName"
								},
								{
									"values": [
										{
											"string_value": "18eacdf7-8db3-4028-8ce8-a686ec639d75"
										}
									],
									"id": "id"
								},
								{
									"values": [
										{
											"string_value": "62e90394-69f5-4237-9190-012177145e10"
										}
									],
									"id": "roleTemplateId"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJjdXJzb3IiOiIyIn0="
				}
			}`),
		},

		"second_page": {
			request: &adapter_api_v1.GetPageRequest{
				Datasource: &adapter_api_v1.DatasourceConfig{
					Auth: &adapter_api_v1.DatasourceAuthCredentials{
						AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer {{OMITTED}}",
						},
					},
					Address: "graph.microsoft.com",
					Id:      "AzureAD",
					Type:    "AzureAD-1.0.1",
					Config:  []byte(`{"apiVersion":"v1.0"}`),
				},
				Entity: &adapter_api_v1.EntityConfig{
					Id:         "AzureADRole",
					ExternalId: "Role",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "id",
							ExternalId: "id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "displayName",
							ExternalId: "displayName",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "description",
							ExternalId: "description",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "roleTemplateId",
							ExternalId: "roleTemplateId",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOiIyIn0=",
			},
			response: []byte(`
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Can create application registrations independent of the 'Users can register applications' setting."
										}
									],
									"id": "description"
								},
								{
									"values": [
										{
											"string_value": "Application Developer"
										}
									],
									"id": "displayName"
								},
								{
									"values": [
										{
											"string_value": "321fd63c-c37c-4a77-bf46-ee0acd84476e"
										}
									],
									"id": "id"
								},
								{
									"values": [
										{
											"string_value": "cf1c38e5-3621-4004-a7cb-879624dced7c"
										}
									],
									"id": "roleTemplateId"
								}
							],
							"child_objects": []
						},
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Can manage role assignments in Microsoft Entra ID, and all aspects of Privileged Identity Management."
										}
									],
									"id": "description"
								},
								{
									"values": [
										{
											"string_value": "Privileged Role Administrator"
										}
									],
									"id": "displayName"
								},
								{
									"values": [
										{
											"string_value": "33a4c989-c3ff-4597-81c4-81e0a93ffb6e"
										}
									],
									"id": "id"
								},
								{
									"values": [
										{
											"string_value": "e8611ab8-c189-46e8-94e1-60213ab1f814"
										}
									],
									"id": "roleTemplateId"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJjdXJzb3IiOiI0In0="
				}
			}`),
		},

		"last_page": {
			request: &adapter_api_v1.GetPageRequest{
				Datasource: &adapter_api_v1.DatasourceConfig{
					Auth: &adapter_api_v1.DatasourceAuthCredentials{
						AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer {{OMITTED}}",
						},
					},
					Address: "graph.microsoft.com",
					Id:      "AzureAD",
					Type:    "AzureAD-1.0.1",
					Config:  []byte(`{"apiVersion":"v1.0"}`),
				},
				Entity: &adapter_api_v1.EntityConfig{
					Id:         "AzureADRole",
					ExternalId: "Role",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "id",
							ExternalId: "id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "displayName",
							ExternalId: "displayName",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "description",
							ExternalId: "description",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "roleTemplateId",
							ExternalId: "roleTemplateId",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
					},
				},
				PageSize: 2,
				Cursor:   "eyJjdXJzb3IiOiI4In0=",
			},
			response: []byte(`
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Members of this role can create/manage groups, create/manage groups settings like naming and expiration policies, and view groups activity and audit reports."
										}
									],
									"id": "description"
								},
								{
									"values": [
										{
											"string_value": "Groups Administrator"
										}
									],
									"id": "displayName"
								},
								{
									"values": [
										{
											"string_value": "fb96a81c-6147-4cfe-b7fe-c63c2725e7c9"
										}
									],
									"id": "id"
								},
								{
									"values": [
										{
											"string_value": "fdd7a751-b60b-444a-984c-02652fe8fa1c"
										}
									],
									"id": "roleTemplateId"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": ""
				}
			}`),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResp, err := adapterClient.GetPage(ctx, tt.request)
			if err != nil {
				t.Fatal(err)
			}

			wantResp := new(adapter_api_v1.GetPageResponse)

			err = protojson.Unmarshal(tt.response, wantResp)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}

	close(stop)
}

func TestAzureADAdapter_Application(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/azuread/application")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "graph.microsoft.com",
			Id:      "AzureAD",
			Type:    "AzureAD-1.0.1",
			Config:  []byte(`{"apiVersion":"v1.0"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "AzureADApplication",
			ExternalId: "Application",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "displayName",
					ExternalId: "displayName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "createdDateTime",
					ExternalId: "createdDateTime",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
				},
			},
		},
		PageSize: 3,
		Cursor:   "",
	})
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success": {
				"objects": [
					{
						"attributes": [
							{
								"id": "createdDateTime",
								"values": [
									{
										"datetime_value": {
											"timestamp": "2023-05-05T00:01:52Z",
											"timezone_offset": 0
										}
									}
								]
							},
							{
								"id": "displayName",
								"values": [
									{
										"string_value": "Example App 1"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "createdDateTime",
								"values": [
									{
										"datetime_value": {
											"timestamp": "2023-03-03T01:44:45Z",
											"timezone_offset": 0
										}
									}
								]
							},
							{
								"id": "displayName",
								"values": [
									{
										"string_value": "example-test-sbx"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "createdDateTime",
								"values": [
									{
										"datetime_value": {
											"timestamp": "2022-12-09T23:34:32Z",
											"timezone_offset": 0
										}
									}
								]
							},
							{
								"id": "displayName",
								"values": [
									{
										"string_value": "example-test-stg Test Instance"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "kkkkkkkk-llll-mmmm-nnnn-oooooooooooo"
									}
								]
							}
						]
					}
				],
				"nextCursor": "eyJjdXJzb3IiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20vdjEuMC9hcHBsaWNhdGlvbnM/JHNlbGVjdD1pZCUyY2Rpc3BsYXlOYW1lJTJjY3JlYXRlZERhdGVUaW1lXHUwMDI2JHRvcD0zXHUwMDI2JHNraXB0b2tlbj1ORVhUTElOS19UT0tFTl9QTEFDRUhPTERFUiJ9"
			}
		}
	`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAzureADAdapter_Device(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/azuread/device")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "graph.microsoft.com",
			Id:      "AzureAD",
			Type:    "AzureAD-1.0.1",
			Config:  []byte(`{"apiVersion":"v1.0"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "AzureADDevice",
			ExternalId: "Device",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "physicalIds",
					ExternalId: "physicalIds",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       true,
				},
				{
					Id:         "deviceId",
					ExternalId: "deviceId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 3,
		Cursor:   "",
	})
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success": {
				"objects": [
					{
						"attributes": [
							{
								"id": "deviceId",
								"values": [
									{
										"string_value": "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
									}
								]
							},
							{
								"id": "physicalIds",
								"values": [
									{
										"string_value": "[USER-GID]:kkkkkkkk-llll-mmmm-nnnn-oooooooooooo:1111111111111111"
									},
									{
										"string_value": "[GID]:g:1111111111111111"
									},
									{
										"string_value": "[USER-HWID]:kkkkkkkk-llll-mmmm-nnnn-oooooooooooo:2222222222222222"
									},
									{
										"string_value": "[HWID]:h:2222222222222222"
									}
								]
							}
						]
					}
				]
			}
		}
	`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAzureADAdapter_GroupMember(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/azuread/groupMember")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "graph.microsoft.com",
			Id:      "AzureAD",
			Type:    "AzureAD-1.0.1",
			Config:  []byte(`{"apiVersion":"v1.0"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Example Group 1Member",
			ExternalId: "GroupMember",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "memberId",
					ExternalId: "memberId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "groupId",
					ExternalId: "groupId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 3,
		Cursor:   "",
	})
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success": {
				"objects": [
					{
						"attributes": [
							{
								"id": "groupId",
								"values": [
									{
										"string_value": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt-aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
									}
								]
							},
							{
								"id": "memberId",
								"values": [
									{
										"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "groupId",
								"values": [
									{
										"string_value": "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "zzzzzzzz-yyyy-xxxx-wwww-vvvvvvvvvvvv-ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
									}
								]
							},
							{
								"id": "memberId",
								"values": [
									{
										"string_value": "zzzzzzzz-yyyy-xxxx-wwww-vvvvvvvvvvvv"
									}
								]
							}
						]
					}
				]
			}
		}
	`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAzureADAdapter_RoleMember(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/azuread/roleMember")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	tests := map[string]struct {
		request  *adapter_api_v1.GetPageRequest
		response []byte
	}{
		"request_with_no_members": {
			request: &adapter_api_v1.GetPageRequest{
				Datasource: &adapter_api_v1.DatasourceConfig{
					Auth: &adapter_api_v1.DatasourceAuthCredentials{
						AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer {{OMITTED}}",
						},
					},
					Address: "graph.microsoft.com",
					Id:      "AzureAD",
					Type:    "AzureAD-1.0.1",
					Config:  []byte(`{"apiVersion":"v1.0"}`),
				},
				Entity: &adapter_api_v1.EntityConfig{
					Id:         "AzureADRoleMember",
					ExternalId: "RoleMember",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "id",
							ExternalId: "id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "memberId",
							ExternalId: "memberId",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "roleId",
							ExternalId: "roleId",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
				},
				PageSize: 2,
				Cursor:   "",
			},
			response: []byte(`
			{
				"success": {
					"objects": [],
					"next_cursor": ""
				}
			}
			`),
		},

		// Note: The API response may contain fewer records than the specified page size,
		// especially when filtering records using /microsoft.graph.directoryRole.
		// The API might provide fewer records as it excludes entries other than directory roles from the page.
		// In this case, response will contain a nextLink to access the next set of records.
		"request_with_members": {
			request: &adapter_api_v1.GetPageRequest{
				Datasource: &adapter_api_v1.DatasourceConfig{
					Auth: &adapter_api_v1.DatasourceAuthCredentials{
						AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer {{OMITTED}}",
						},
					},
					Address: "graph.microsoft.com",
					Id:      "AzureAD",
					Type:    "AzureAD-1.0.1",
					Config:  []byte(`{"apiVersion":"v1.0"}`),
				},
				Entity: &adapter_api_v1.EntityConfig{
					Id:         "AzureADRoleMember",
					ExternalId: "RoleMember",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "id",
							ExternalId: "id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "memberId",
							ExternalId: "memberId",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "roleId",
							ExternalId: "roleId",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
				},
				PageSize: 10,
				Cursor:   "eyJjb2xsZWN0aW9uSWQiOiIyMDFkMzFjMC02NTNkLTQzYTYtYWRmMC1hZWU4OWE3OWM4MDUiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiaHR0cHM6Ly9ncmFwaC5taWNyb3NvZnQuY29tL3YxLjAvdXNlcnM/JHNlbGVjdD1pZFx1MDAyNiR0b3A9MVx1MDAyNiRza2lwdG9rZW49TkVYVExJTktfVE9LRU5fUExBQ0VIT0xERVJfMiJ9",
			},
			response: []byte(`
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "0fea7f0d-dea1-458d-9099-69fcc2e3cd42-uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"
										}
									],
									"id": "id"
								},
								{
									"values": [
										{
											"string_value": "uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"
										}
									],
									"id": "memberId"
								},
								{
									"values": [
										{
											"string_value": "0fea7f0d-dea1-458d-9099-69fcc2e3cd42"
										}
									],
									"id": "roleId"
								}
							],
							"child_objects": []
						},
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "45a8712b-a814-4b9f-8a1e-4b31714efa12-uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"
										}
									],
									"id": "id"
								},
								{
									"values": [
										{
											"string_value": "uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"
										}
									],
									"id": "memberId"
								},
								{
									"values": [
										{
											"string_value": "45a8712b-a814-4b9f-8a1e-4b31714efa12"
										}
									],
									"id": "roleId"
								}
							],
							"child_objects": []
						},
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "540b4b34-c25b-437d-8eee-329463952334-uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"
										}
									],
									"id": "id"
								},
								{
									"values": [
										{
											"string_value": "uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"
										}
									],
									"id": "memberId"
								},
								{
									"values": [
										{
											"string_value": "540b4b34-c25b-437d-8eee-329463952334"
										}
									],
									"id": "roleId"
								}
							],
							"child_objects": []
						},
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "d973db57-eb50-4356-959e-f1ce19a22b98-uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"
										}
									],
									"id": "id"
								},
								{
									"values": [
										{
											"string_value": "uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"
										}
									],
									"id": "memberId"
								},
								{
									"values": [
										{
											"string_value": "d973db57-eb50-4356-959e-f1ce19a22b98"
										}
									],
									"id": "roleId"
								}
							],
							"child_objects": []
						},
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "fc6c3c82-669c-4e24-b089-2a2847a43d14-uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"
										}
									],
									"id": "id"
								},
								{
									"values": [
										{
											"string_value": "uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"
										}
									],
									"id": "memberId"
								},
								{
									"values": [
										{
											"string_value": "fc6c3c82-669c-4e24-b089-2a2847a43d14"
										}
									],
									"id": "roleId"
								}
							],
							"child_objects": []
						},
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "4231e0d9-0f7f-47ff-9af6-473ab9356eda-uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"
										}
									],
									"id": "id"
								},
								{
									"values": [
										{
											"string_value": "uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"
										}
									],
									"id": "memberId"
								},
								{
									"values": [
										{
											"string_value": "4231e0d9-0f7f-47ff-9af6-473ab9356eda"
										}
									],
									"id": "roleId"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": ""
				}
			}
			`),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResp, err := adapterClient.GetPage(ctx, tt.request)
			if err != nil {
				t.Fatal(err)
			}

			wantResp := new(adapter_api_v1.GetPageResponse)

			err = protojson.Unmarshal(tt.response, wantResp)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}

	close(stop)
}

func TestAzureADAdapter_RoleAssignment(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/azuread/roleAssignment")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "graph.microsoft.com",
			Id:      "AzureAD",
			Type:    "AzureAD-1.0.1",
			Config:  []byte(`{"apiVersion":"v1.0"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "RoleAssignment",
			ExternalId: "RoleAssignment",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "principalId",
					ExternalId: "principalId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "roleDefinitionId",
					ExternalId: "roleDefinitionId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 3,
		Cursor:   "",
	})
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success": {
				"objects": [
					{
						"attributes": [
							{
								"id": "id",
								"values": [
									{
										"string_value": "role-assignment-id-1"
									}
								]
							},
							{
								"id": "principalId",
								"values": [
									{
										"string_value": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
									}
								]
							},
							{
								"id": "roleDefinitionId",
								"values": [
									{
										"string_value": "62e90394-69f5-4237-9190-012177145e10"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "id",
								"values": [
									{
										"string_value": "role-assignment-id-2"
									}
								]
							},
							{
								"id": "principalId",
								"values": [
									{
										"string_value": "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
									}
								]
							},
							{
								"id": "roleDefinitionId",
								"values": [
									{
										"string_value": "62e90394-69f5-4237-9190-012177145e10"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "id",
								"values": [
									{
										"string_value": "role-assignment-id-3"
									}
								]
							},
							{
								"id": "principalId",
								"values": [
									{
										"string_value": "kkkkkkkk-llll-mmmm-nnnn-oooooooooooo"
									}
								]
							},
							{
								"id": "roleDefinitionId",
								"values": [
									{
										"string_value": "62e90394-69f5-4237-9190-012177145e10"
									}
								]
							}
						]
					}
				]
			}
		}
	`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(wantResp, gotResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

// Page size 3. First page.
func TestAzureADAdapter_RoleAssignmentScheduleRequests_page_1(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/azuread/roleAssignmentScheduleRequests_page_1")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "graph.microsoft.com",
			Id:      "AzureAD",
			Type:    "AzureAD-1.0.1",
			Config:  []byte(`{"apiVersion":"v1.0"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "RoleAssignmentScheduleRequest",
			ExternalId: "RoleAssignmentScheduleRequest",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "principalId",
					ExternalId: "principalId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "roleDefinitionId",
					ExternalId: "roleDefinitionId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "action",
					ExternalId: "action",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "justification",
					ExternalId: "justification",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 3,
		Cursor:   "",
	})
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success": {
				"objects": [
					{
						"attributes": [
							{
								"id": "action",
								"values": [
									{
										"string_value": "adminRemove"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "44a37286-00ec-4315-bf8f-89ddcc0ea41a"
									}
								]
							},
							{
								"id": "principalId",
								"values": [
									{
										"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
									}
								]
							},
							{
								"id": "roleDefinitionId",
								"values": [
									{
										"string_value": "9b895d92-2cd3-44c7-9d02-a6ac2d5ea5c3"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "action",
								"values": [
									{
										"string_value": "selfDeactivate"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "264332a6-3951-4616-bee3-9ab3a8d188c1"
									}
								]
							},
							{
								"id": "justification",
								"values": [
									{
										"string_value": "Deactivation request"
									}
								]
							},
							{
								"id": "principalId",
								"values": [
									{
										"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
									}
								]
							},
							{
								"id": "roleDefinitionId",
								"values": [
									{
										"string_value": "cf1c38e5-3621-4004-a7cb-879624dced7c"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "action",
								"values": [
									{
										"string_value": "selfDeactivate"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "32dc0cdb-d927-428d-b567-72e46e0da6c9"
									}
								]
							},
							{
								"id": "justification",
								"values": [
									{
										"string_value": "Deactivation request"
									}
								]
							},
							{
								"id": "principalId",
								"values": [
									{
										"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
									}
								]
							},
							{
								"id": "roleDefinitionId",
								"values": [
									{
										"string_value": "cf1c38e5-3621-4004-a7cb-879624dced7c"
									}
								]
							}
						]
					}
				],
				"nextCursor": "eyJjdXJzb3IiOiIzIn0="
			}
		}
	`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(wantResp, gotResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

// Page size 3. 2nd page.
func TestAzureADAdapter_RoleAssignmentScheduleRequests_page_2(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/azuread/roleAssignmentScheduleRequests_page_2")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "graph.microsoft.com",
			Id:      "AzureAD",
			Type:    "AzureAD-1.0.1",
			Config:  []byte(`{"apiVersion":"v1.0"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "RoleAssignmentScheduleRequest",
			ExternalId: "RoleAssignmentScheduleRequest",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "principalId",
					ExternalId: "principalId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "roleDefinitionId",
					ExternalId: "roleDefinitionId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "action",
					ExternalId: "action",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "justification",
					ExternalId: "justification",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 3,
		Cursor:   "eyJjdXJzb3IiOiIzIn0=",
	})
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success": {
				"objects": [
					{
						"attributes": [
							{
								"id": "action",
								"values": [
									{
										"string_value": "selfDeactivate"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "42f1f2ac-6a85-41bc-8dea-47e041036c07"
									}
								]
							},
							{
								"id": "justification",
								"values": [
									{
										"string_value": "Deactivation request"
									}
								]
							},
							{
								"id": "principalId",
								"values": [
									{
										"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
									}
								]
							},
							{
								"id": "roleDefinitionId",
								"values": [
									{
										"string_value": "cf1c38e5-3621-4004-a7cb-879624dced7c"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "action",
								"values": [
									{
										"string_value": "selfDeactivate"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "82a9d99b-6d13-4dfc-b173-958ef0ee0525"
									}
								]
							},
							{
								"id": "justification",
								"values": [
									{
										"string_value": "Deactivation request"
									}
								]
							},
							{
								"id": "principalId",
								"values": [
									{
										"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
									}
								]
							},
							{
								"id": "roleDefinitionId",
								"values": [
									{
										"string_value": "cf1c38e5-3621-4004-a7cb-879624dced7c"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "action",
								"values": [
									{
										"string_value": "selfDeactivate"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "8651181b-e433-4a53-8839-0bc85fe4a3a2"
									}
								]
							},
							{
								"id": "justification",
								"values": [
									{
										"string_value": "Deactivation request"
									}
								]
							},
							{
								"id": "principalId",
								"values": [
									{
										"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
									}
								]
							},
							{
								"id": "roleDefinitionId",
								"values": [
									{
										"string_value": "cf1c38e5-3621-4004-a7cb-879624dced7c"
									}
								]
							}
						]
					}
				],
				"nextCursor": "eyJjdXJzb3IiOiI2In0="
			}
		}
	`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(wantResp, gotResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

// This smoketest contains made-up data (as per spec) as the SoR does not have any data for this entity.
func TestAzureADAdapter_GroupAssignmentScheduleRequests(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/azuread/groupAssignmentScheduleRequests")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "graph.microsoft.com",
			Id:      "AzureAD",
			Type:    "AzureAD-1.0.1",
			Config:  []byte(`{"apiVersion":"v1.0"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "GroupAssignmentScheduleRequest",
			ExternalId: "GroupAssignmentScheduleRequest",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "principalId",
					ExternalId: "principalId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "groupId",
					ExternalId: "groupId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "action",
					ExternalId: "action",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 3,
		Cursor:   "",
	})
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success": {
				"objects": [
					{
						"attributes": [
							{
								"id": "action",
								"values": [
									{
										"string_value": "adminAssign"
									}
								]
							},
							{
								"id": "groupId",
								"values": [
									{
										"string_value": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "example-request-1"
									}
								]
							},
							{
								"id": "principalId",
								"values": [
									{
										"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "action",
								"values": [
									{
										"string_value": "adminAssign"
									}
								]
							},
							{
								"id": "groupId",
								"values": [
									{
										"string_value": "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "example-request-2"
									}
								]
							},
							{
								"id": "principalId",
								"values": [
									{
										"string_value": "uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "action",
								"values": [
									{
										"string_value": "adminAssign"
									}
								]
							},
							{
								"id": "groupId",
								"values": [
									{
										"string_value": "kkkkkkkk-llll-mmmm-nnnn-oooooooooooo"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "example-request-3"
									}
								]
							},
							{
								"id": "principalId",
								"values": [
									{
										"string_value": "zzzzzzzz-aaaa-bbbb-cccc-dddddddddddd"
									}
								]
							}
						]
					}
				],
				"nextCursor": "eyJjdXJzb3IiOiIzIn0="
			}
		}
	`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(wantResp, gotResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

// This tests advanced filters.
// The baseline for reference is a groupID: aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee which has 3 members.
// These 3 members are further filtered down for department="IT/IS" using advanced filters.
/*
GET https://graph.microsoft.com/v1.0/groups/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/members/microsoft.graph.user?$count=true&$select=id,displayName,department

{
    "@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users(id,displayName,department)",
    "@odata.count": 3,
    "value": [
        {
            "id": "pppppppp-qqqq-rrrr-ssss-tttttttttttt",
            "displayName": "Bacong, Alejandro",
            "department": "IT/IS"
        },
        {
            "id": "459ba88b-2e43-43d1-9e2e-b8c4987292a1",
            "displayName": "Joe",
            "department": null
        },
        {
            "id": "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj",
            "displayName": "Marc",
            "department": null
        }
    ]
}
This is the only page of data.
*/
func TestAzureADAdapter_GroupMember_With_Advanced_Filter(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/azuread/groupMember_with_advanced_filter")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "graph.microsoft.com",
			Id:      "AzureAD",
			Type:    "AzureAD-1.0.1",
			// Get users of IT/IS dept for a specific group ID.
			Config: []byte(`{
				"apiVersion": "v1.0",
				"advancedFilters": {
					"getObjectsByScope": {
						"GroupMember": [
							{
								"scopeEntity": "Group",
								"scopeEntityFilter": "id eq 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee'",
								"members": [
									{
										"memberEntity": "User",
										"memberEntityFilter": "department eq 'IT/IS'"
									}
								]
							}
						]
					}
				}
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Example Group 1Member",
			ExternalId: "GroupMember",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "memberId",
					ExternalId: "memberId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "groupId",
					ExternalId: "groupId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 4,
		Cursor:   "",
	})
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	// Expecting 1 user in the response as the API returns this
	/*
		GET https://graph.microsoft.com/v1.0/groups/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/members/microsoft.graph.user?$filter=department eq 'IT/IS'&$count=true&$select=id

		{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users(id)",
			"@odata.count": 1,
			"value": [
				{
					"id": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
				}
			]
		}
	*/
	err = protojson.Unmarshal([]byte(`
		{
			"success": {
				"objects": [
					{
						"attributes": [
							{
								"id": "groupId",
								"values": [
									{
										"string_value": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt-aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
									}
								]
							},
							{
								"id": "memberId",
								"values": [
									{
										"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
									}
								]
							}
						]
					}
				],
				"nextCursor": ""
			}
		}
	`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

// This tests advanced filters.
// The baseline for reference is a tenant with about 14 groups.
// These groups are further filtered down using advanced filters i.e. startsWith(displayName, 'Test')
/*
GET https://graph.microsoft.com/v1.0/groups?$select=id,displayName&$top=999&$filter=startswith(displayName, 'Test')&$count=true

{
    "@odata.context": "https://graph.microsoft.com/v1.0/$metadata#groups(id,displayName)",
    "value": [
        {
            "id": "pppppppp-qqqq-rrrr-ssss-tttttttttttt",
            "displayName": "Test Group 3"
        },
        {
            "id": "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj",
            "displayName": "Test Group 6"
        },
        {
            "id": "kkkkkkkk-llll-mmmm-nnnn-oooooooooooo",
            "displayName": "Test Security Group"
        },
        {
            "id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
            "displayName": "Test Security Group 2"
        }
    ]
}
This is the only page of data.
*/

const AdvancedFiltersWithImplicitFilteringConfig = `
{
  "requestTimeoutSeconds": 10,
  "apiVersion": "v1.0",
  "advancedFilters": {
    "getObjectsByScope": {
      "GroupMember": [
        {
          "scopeEntity": "Group",
          "scopeEntityFilter": "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
          "members": [
            {
              "memberEntity": "User",
              "memberEntityFilter": "department eq 'engineering'"
            }
          ]
        },
        {
          "scopeEntity": "Group",
          "scopeEntityFilter": "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
          "members": [
            {
              "memberEntity": "User",
              "memberEntityFilter": "department eq 'product'"
            },
            {
              "memberEntity": "Group",
              "memberEntityFilter": "startswith(displayName, 'California')"
            }
          ]
        },
        {
          "scopeEntity": "Group",
          "scopeEntityFilter": "startswith(displayName, 'California')",
          "members": [
            {
              "memberEntity": "User"
            }
          ]
        }
      ]
    }
  }
}
`

func TestAzureADAdapter_GroupMember_With_Advanced_Filters_And_Implicit_Filtering(t *testing.T) {
	fixtureFile := "fixtures/azuread/groupMember_with_advanced_filters_and_implicit_filtering"

	tests := []struct {
		wantResp string
	}{
		{
			wantResp: `{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
										}
									],
									"id": "af2477f0-e6e6-4d4a-b11e-82e2dab7b148"
								},
								{
									"values": [
										{
											"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt-aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
										}
									],
									"id": "8647328f-e9dc-4603-b015-fac4c21fbaf5"
								},
								{
									"values": [
										{
											"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
										}
									],
									"id": "33a7e7fd-8b1d-4bb0-9aab-d7772435a83d"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJlbnRpdHlGaWx0ZXJJbmRleCI6MSwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6bnVsbH0="
				}
			}`,
		},
		{
			wantResp: `{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
										}
									],
									"id": "af2477f0-e6e6-4d4a-b11e-82e2dab7b148"
								},
								{
									"values": [
										{
											"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt-aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
										}
									],
									"id": "8647328f-e9dc-4603-b015-fac4c21fbaf5"
								},
								{
									"values": [
										{
											"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
										}
									],
									"id": "33a7e7fd-8b1d-4bb0-9aab-d7772435a83d"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJlbnRpdHlGaWx0ZXJJbmRleCI6MSwibWVtYmVyRmlsdGVySW5kZXgiOjEsImN1cnNvciI6bnVsbH0="
				}
			}`,
		},
		{
			wantResp: `{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
										}
									],
									"id": "af2477f0-e6e6-4d4a-b11e-82e2dab7b148"
								},
								{
									"values": [
										{
											"string_value": "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj-aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
										}
									],
									"id": "8647328f-e9dc-4603-b015-fac4c21fbaf5"
								},
								{
									"values": [
										{
											"string_value": "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
										}
									],
									"id": "33a7e7fd-8b1d-4bb0-9aab-d7772435a83d"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJlbnRpdHlGaWx0ZXJJbmRleCI6MiwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6bnVsbH0="
				}
			}`,
		},
		{
			wantResp: `{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
										}
									],
									"id": "af2477f0-e6e6-4d4a-b11e-82e2dab7b148"
								},
								{
									"values": [
										{
											"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt-ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
										}
									],
									"id": "8647328f-e9dc-4603-b015-fac4c21fbaf5"
								},
								{
									"values": [
										{
											"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
										}
									],
									"id": "33a7e7fd-8b1d-4bb0-9aab-d7772435a83d"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": ""
				}
			}`,
		},
	}

	for pageNumber, test := range tests {
		fixtureName := fmt.Sprintf("%s_page_%d", fixtureFile, pageNumber+1)

		httpClient, recorder := common.StartRecorder(t, fixtureName)
		defer recorder.Stop()

		port := common.AvailableTestPort(t)

		stop := make(chan struct{})

		// Start Adapter Server
		go func() {
			stop = common.StartAdapterServer(t, httpClient, port)
		}()

		defer close(stop)

		time.Sleep(10 * time.Millisecond)

		adapterClient, conn := common.GetNewAdapterClient(t, port)
		defer conn.Close()

		ctx, cancelCtx := common.GetAdapterCtx()
		defer cancelCtx()

		var cursor string
		// Extract the cursor from the previous page.
		if pageNumber != 0 {
			previousResp := new(adapter_api_v1.GetPageResponse)
			if err := protojson.Unmarshal([]byte(tests[pageNumber-1].wantResp), previousResp); err != nil {
				t.Fatal(err)
			}

			cursor = previousResp.GetSuccess().NextCursor
		}

		gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
			Datasource: &adapter_api_v1.DatasourceConfig{
				Auth: &adapter_api_v1.DatasourceAuthCredentials{
					AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
						HttpAuthorization: "Bearer {{OMITTED}}",
					},
				},
				Address: "graph.microsoft.com",
				Id:      "AzureAD",
				Type:    "AzureAD-1.0.1",
				Config:  []byte(AdvancedFiltersWithImplicitFilteringConfig),
			},
			Entity: &adapter_api_v1.EntityConfig{
				Id:         "Example Group 1Member",
				ExternalId: "GroupMember",
				Ordered:    false,
				Attributes: []*adapter_api_v1.AttributeConfig{
					{
						Id:         "8647328f-e9dc-4603-b015-fac4c21fbaf5",
						ExternalId: "id",
						Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
					{
						Id:         "33a7e7fd-8b1d-4bb0-9aab-d7772435a83d",
						ExternalId: "memberId",
						Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
					{
						Id:         "af2477f0-e6e6-4d4a-b11e-82e2dab7b148",
						ExternalId: "groupId",
						Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			PageSize: 10,
			Cursor:   cursor,
		})
		if err != nil {
			t.Fatal(err)
		}

		wantResp := new(adapter_api_v1.GetPageResponse)

		if err := protojson.Unmarshal([]byte(test.wantResp), wantResp); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
			t.Fatal(diff)
		}
	}
}

func TestAzureADAdapter_User_With_Advanced_Filters_And_Implicit_Filtering(t *testing.T) {
	fixtureFile := "fixtures/azuread/user_with_advanced_filters_and_implicit_filtering"

	tests := []struct {
		wantResp string
	}{
		{
			wantResp: `{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Richard"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJlbnRpdHlGaWx0ZXJJbmRleCI6MSwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6bnVsbH0="
				}
			}`,
		},
		{
			wantResp: `{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Marc"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJlbnRpdHlGaWx0ZXJJbmRleCI6MiwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6bnVsbH0="
				}
			}`,
		},
		{
			wantResp: `{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Scott"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "pppppppp-qqqq-rrrr-ssss-tttttttttttt"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": ""
				}
			}`,
		},
	}

	for pageNumber, test := range tests {
		fixtureName := fmt.Sprintf("%s_page_%d", fixtureFile, pageNumber+1)

		httpClient, recorder := common.StartRecorder(t, fixtureName)
		defer recorder.Stop()

		port := common.AvailableTestPort(t)

		stop := make(chan struct{})

		// Start Adapter Server
		go func() {
			stop = common.StartAdapterServer(t, httpClient, port)
		}()

		defer close(stop)

		time.Sleep(10 * time.Millisecond)

		adapterClient, conn := common.GetNewAdapterClient(t, port)
		defer conn.Close()

		ctx, cancelCtx := common.GetAdapterCtx()
		defer cancelCtx()

		var cursor string
		// Extract the cursor from the previous page.
		if pageNumber != 0 {
			previousResp := new(adapter_api_v1.GetPageResponse)
			if err := protojson.Unmarshal([]byte(tests[pageNumber-1].wantResp), previousResp); err != nil {
				t.Fatal(err)
			}

			cursor = previousResp.GetSuccess().NextCursor
		}

		gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
			Datasource: &adapter_api_v1.DatasourceConfig{
				Auth: &adapter_api_v1.DatasourceAuthCredentials{
					AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
						HttpAuthorization: "Bearer {{OMITTED}}",
					},
				},
				Address: "graph.microsoft.com",
				Id:      "AzureAD",
				Type:    "AzureAD-1.0.1",
				Config:  []byte(AdvancedFiltersWithImplicitFilteringConfig),
			},
			Entity: &adapter_api_v1.EntityConfig{
				Id:         "AzureADUser",
				ExternalId: "User",
				Ordered:    false,
				Attributes: []*adapter_api_v1.AttributeConfig{
					{
						Id:         "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6",
						ExternalId: "id",
						UniqueId:   true,
						Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
					{
						Id:         "ce0954e6-cc24-485d-b095-d5f0ee39c68c",
						ExternalId: "displayName",
						Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			PageSize: 10,
			Cursor:   cursor,
		})
		if err != nil {
			t.Fatal(err)
		}

		wantResp := new(adapter_api_v1.GetPageResponse)

		if err := protojson.Unmarshal([]byte(test.wantResp), wantResp); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
			t.Fatal(diff)
		}
	}
}

func TestAzureADAdapter_Group_With_Advanced_Filters_And_Implicit_Filtering(t *testing.T) {
	fixtureFile := "fixtures/azuread/group_with_advanced_filters_and_implicit_filtering"

	tests := []struct {
		wantResp string
	}{
		{
			wantResp: `{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Canada"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJlbnRpdHlGaWx0ZXJJbmRleCI6MSwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6bnVsbH0="
				}
			}`,
		},
		{
			wantResp: `{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "United States"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJlbnRpdHlGaWx0ZXJJbmRleCI6MiwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6bnVsbH0="
				}
			}`,
		},
		{
			wantResp: `{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "California"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJlbnRpdHlGaWx0ZXJJbmRleCI6MywibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6bnVsbH0="
				}
			}`,
		},
		// This duplicate Group node is intentional and is due to the nature of the implicit filters being generated.
		// A duplicate Group node across different pages should not affect ingestion correctness.
		{
			wantResp: `{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "California"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": ""
				}
			}`,
		},
	}

	for pageNumber, test := range tests {
		fixtureName := fmt.Sprintf("%s_page_%d", fixtureFile, pageNumber+1)

		httpClient, recorder := common.StartRecorder(t, fixtureName)
		defer recorder.Stop()

		port := common.AvailableTestPort(t)

		stop := make(chan struct{})

		// Start Adapter Server
		go func() {
			stop = common.StartAdapterServer(t, httpClient, port)
		}()

		defer close(stop)

		time.Sleep(10 * time.Millisecond)

		adapterClient, conn := common.GetNewAdapterClient(t, port)
		defer conn.Close()

		ctx, cancelCtx := common.GetAdapterCtx()
		defer cancelCtx()

		var cursor string
		// Extract the cursor from the previous page.
		if pageNumber != 0 {
			previousResp := new(adapter_api_v1.GetPageResponse)
			if err := protojson.Unmarshal([]byte(tests[pageNumber-1].wantResp), previousResp); err != nil {
				t.Fatal(err)
			}

			cursor = previousResp.GetSuccess().NextCursor
		}

		gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
			Datasource: &adapter_api_v1.DatasourceConfig{
				Auth: &adapter_api_v1.DatasourceAuthCredentials{
					AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
						HttpAuthorization: "Bearer {{OMITTED}}",
					},
				},
				Address: "graph.microsoft.com",
				Id:      "AzureAD",
				Type:    "AzureAD-1.0.1",
				Config:  []byte(AdvancedFiltersWithImplicitFilteringConfig),
			},
			Entity: &adapter_api_v1.EntityConfig{
				Id:         "Example Group 1",
				ExternalId: "Group",
				Ordered:    false,
				Attributes: []*adapter_api_v1.AttributeConfig{
					{
						Id:         "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6",
						ExternalId: "id",
						UniqueId:   true,
						Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
					{
						Id:         "ce0954e6-cc24-485d-b095-d5f0ee39c68c",
						ExternalId: "displayName",
						Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			PageSize: 10,
			Cursor:   cursor,
		})
		if err != nil {
			t.Fatal(err)
		}

		wantResp := new(adapter_api_v1.GetPageResponse)

		if err := protojson.Unmarshal([]byte(test.wantResp), wantResp); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
			t.Fatal(diff)
		}
	}
}
