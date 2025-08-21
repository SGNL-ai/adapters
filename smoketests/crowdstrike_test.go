// Copyright 2025 SGNL.ai, Inc.

// nolint: lll
package smoketests

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	adapter_api_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/smoketests/common"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestCrowdStrikeUser(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/crowdstrike/user")
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

	req := &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "https://api.us-2.crowdstrike.com",
			Id:      "Test",
			Type:    "CrowdStrike-1.0.0",
			Config: []byte(`
			{
				"archived": false,
				"enabled": true,
				"apiVersion": "v1"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "User",
			ExternalId: "user",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "entityId",
					ExternalId: "entityId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
					UniqueId:   true,
				},
				{
					Id:         "inactive",
					ExternalId: "inactive",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "creationTime",
					ExternalId: "creationTime",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					Id:         "ActiveDirectoryAccount",
					ExternalId: `$.accounts[?(@.__typename=="ActiveDirectoryAccountDescriptor")]`,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "objectGuid",
							ExternalId: "objectGuid",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
							UniqueId:   true,
						},
						{
							Id:         "samAccountName",
							ExternalId: "samAccountName",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
					},
				},
			},
		},
		PageSize: 2,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`{
    "success": {
        "objects": [
            {
                "attributes": [
                    {
                        "id": "creationTime",
                        "values": [
                            {
                                "datetimeValue": {
                                    "timestamp": "2024-05-15T15:29:10Z"
                                }
                            }
                        ]
                    },
                    {
                        "id": "entityId",
                        "values": [
                            {
                                "stringValue": "095b6929-44b9-4525-a0cc-9ef4552011f3"
                            }
                        ]
                    },
                    {
                        "id": "inactive",
                        "values": [
                            {
                                "boolValue": true
                            }
                        ]
                    }
                ],
                "childObjects": [
                    {
                        "entityId": "ActiveDirectoryAccount",
                        "objects": [
                            {
                                "attributes": [
                                    {
                                        "id": "objectGuid",
                                        "values": [
                                            {
                                                "stringValue": "095b6929-44b9-4525-a0cc-9ef4552011f3"
                                            }
                                        ]
                                    },
                                    {
                                        "id": "samAccountName",
                                        "values": [
                                            {
                                                "stringValue": "Wendolyn.Garber"
                                            }
                                        ]
                                    }
                                ]
                            }
                        ]
                    }
                ]
            },
            {
                "attributes": [
                    {
                        "id": "creationTime",
                        "values": [
                            {
                                "datetimeValue": {
                                    "timestamp": "2024-08-25T18:18:00Z"
                                }
                            }
                        ]
                    },
                    {
                        "id": "entityId",
                        "values": [
                            {
                                "stringValue": "83a49ef1-17a7-4fa4-b90f-9142dfa49577"
                            }
                        ]
                    },
                    {
                        "id": "inactive",
                        "values": [
                            {
                                "boolValue": true
                            }
                        ]
                    }
                ],
                "childObjects": [
                    {
                        "entityId": "ActiveDirectoryAccount",
                        "objects": [
                            {
                                "attributes": [
                                    {
                                        "id": "objectGuid",
                                        "values": [
                                            {
                                                "stringValue": "83a49ef1-17a7-4fa4-b90f-9142dfa49577"
                                            }
                                        ]
                                    },
                                    {
                                        "id": "samAccountName",
                                        "values": [
                                            {
                                                "stringValue": "sgnl.sor"
                                            }
                                        ]
                                    }
                                ]
                            }
                        ]
                    }
                ]
            }
        ],
        "nextCursor": "eyJjdXJzb3IiOiJleUp5YVhOclUyTnZjbVVpT2pBdU5qUTFOVGswTnpjMU5UQTJNRFE0TkN3aVgybGtJam9pT0ROaE5EbGxaakV0TVRkaE55MDBabUUwTFdJNU1HWXRPVEUwTW1SbVlUUTVOVGMzSW4wPSJ9"
    }
}`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(wantResp, gotResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestCrowdStrikeEndpoint(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/crowdstrike/endpoint")
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

	req := &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "https://api.us-2.crowdstrike.com",
			Id:      "Test",
			Type:    "CrowdStrike-1.0.0",
			Config: []byte(`
			{
				"archived": false,
				"enabled": true,
				"apiVersion": "v1"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Endpoint",
			ExternalId: "endpoint",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "entityId",
					ExternalId: "entityId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
					UniqueId:   true,
				},
				{
					Id:         "inactive",
					ExternalId: "inactive",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "creationTime",
					ExternalId: "creationTime",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					Id:         "ActiveDirectoryAccount",
					ExternalId: `$.accounts[?(@.__typename=="ActiveDirectoryAccountDescriptor")]`,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "objectGuid",
							ExternalId: "objectGuid",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
							UniqueId:   true,
						},
						{
							Id:         "samAccountName",
							ExternalId: "samAccountName",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
					},
				},
			},
		},
		PageSize: 2,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`{
    "success": {
        "objects": [
            {
                "attributes": [
                    {
                        "id": "creationTime",
                        "values": [
                            {
                                "datetimeValue": {
                                    "timestamp": "2024-05-15T15:17:19Z"
                                }
                            }
                        ]
                    },
                    {
                        "id": "entityId",
                        "values": [
                            {
                                "stringValue": "3c7aebb9-411b-4ee9-b481-e881f29afcc8"
                            }
                        ]
                    },
                    {
                        "id": "inactive",
                        "values": [
                            {
                                "boolValue": false
                            }
                        ]
                    }
                ],
                "childObjects": [
                    {
                        "entityId": "ActiveDirectoryAccount",
                        "objects": [
                            {
                                "attributes": [
                                    {
                                        "id": "objectGuid",
                                        "values": [
                                            {
                                                "stringValue": "3c7aebb9-411b-4ee9-b481-e881f29afcc8"
                                            }
                                        ]
                                    },
                                    {
                                        "id": "samAccountName",
                                        "values": [
                                            {
                                                "stringValue": "mj-dc$"
                                            }
                                        ]
                                    }
                                ]
                            }
                        ]
                    }
                ]
            },
            {
                "attributes": [
                    {
                        "id": "creationTime",
                        "values": [
                            {
                                "datetimeValue": {
                                    "timestamp": "2024-05-29T21:30:17Z"
                                }
                            }
                        ]
                    },
                    {
                        "id": "entityId",
                        "values": [
                            {
                                "stringValue": "89be47c3-f51b-48af-884a-ecb02ed0807a"
                            }
                        ]
                    },
                    {
                        "id": "inactive",
                        "values": [
                            {
                                "boolValue": true
                            }
                        ]
                    }
                ],
                "childObjects": [
                    {
                        "entityId": "ActiveDirectoryAccount",
                        "objects": [
                            {
                                "attributes": [
                                    {
                                        "id": "objectGuid",
                                        "values": [
                                            {
                                                "stringValue": "89be47c3-f51b-48af-884a-ecb02ed0807a"
                                            }
                                        ]
                                    },
                                    {
                                        "id": "samAccountName",
                                        "values": [
                                            {
                                                "stringValue": "ALICE-WIN11$"
                                            }
                                        ]
                                    }
                                ]
                            }
                        ]
                    }
                ]
            }
        ],
        "nextCursor": "eyJjdXJzb3IiOiJleUp5YVhOclUyTnZjbVVpT2pBdU5EYzVNaXdpWDJsa0lqb2lPRGxpWlRRM1l6TXRaalV4WWkwME9HRm1MVGc0TkdFdFpXTmlNREpsWkRBNE1EZGhJbjA9In0="
    }
}`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(wantResp, gotResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestCrowdStrikeIncident(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/crowdstrike/incident")
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

	req := &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "https://api.us-2.crowdstrike.com",
			Id:      "Test",
			Type:    "CrowdStrike-1.0.0",
			Config: []byte(`
			{
				"archived": false,
				"enabled": true,
				"apiVersion": "v1"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Incident",
			ExternalId: "incident",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "incidentId",
					ExternalId: "incidentId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
					UniqueId:   true,
				},
				{
					Id:         "severity",
					ExternalId: "severity",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "startTime",
					ExternalId: "startTime",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					Id:         "CompromisedEntities",
					ExternalId: "$.compromisedEntities",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "entityId",
							ExternalId: "entityId",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
							UniqueId:   true,
						},
						{
							Id:         "primaryDisplayName",
							ExternalId: "primaryDisplayName",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
					},
				},
			},
		},
		PageSize: 2,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`{
    "success": {
        "objects": [
            {
                "attributes": [
                    {
                        "id": "incidentId",
                        "values": [
                            {
                                "stringValue": "INC-16"
                            }
                        ]
                    },
                    {
                        "id": "severity",
                        "values": [
                            {
                                "stringValue": "INFO"
                            }
                        ]
                    },
                    {
                        "id": "startTime",
                        "values": [
                            {
                                "datetimeValue": {
                                    "timestamp": "2024-09-23T13:00:21.995Z"
                                }
                            }
                        ]
                    }
                ],
                "childObjects": [
                    {
                        "entityId": "CompromisedEntities",
                        "objects": [
                            {
                                "attributes": [
                                    {
                                        "id": "entityId",
                                        "values": [
                                            {
                                                "stringValue": "3c7aebb9-411b-4ee9-b481-e881f29afcc8"
                                            }
                                        ]
                                    },
                                    {
                                        "id": "primaryDisplayName",
                                        "values": [
                                            {
                                                "stringValue": "mj-dc"
                                            }
                                        ]
                                    }
                                ]
                            }
                        ]
                    }
                ]
            },
            {
                "attributes": [
                    {
                        "id": "incidentId",
                        "values": [
                            {
                                "stringValue": "INC-15"
                            }
                        ]
                    },
                    {
                        "id": "severity",
                        "values": [
                            {
                                "stringValue": "INFO"
                            }
                        ]
                    },
                    {
                        "id": "startTime",
                        "values": [
                            {
                                "datetimeValue": {
                                    "timestamp": "2024-09-20T01:49:27.080Z"
                                }
                            }
                        ]
                    }
                ],
                "childObjects": [
                    {
                        "entityId": "CompromisedEntities",
                        "objects": [
                            {
                                "attributes": [
                                    {
                                        "id": "entityId",
                                        "values": [
                                            {
                                                "stringValue": "60ee5bb1-805f-46d2-8f3a-9d7cadc52909"
                                            }
                                        ]
                                    },
                                    {
                                        "id": "primaryDisplayName",
                                        "values": [
                                            {
                                                "stringValue": "Alice Wu"
                                            }
                                        ]
                                    }
                                ]
                            }
                        ]
                    }
                ]
            }
        ],
        "nextCursor": "eyJjdXJzb3IiOiJleUpsYm1SVWFXMWxJanA3SWlSa1lYUmxJam9pTWpBeU5DMHdPUzB5TUZRd01UbzFOVG94TUM0eU56UmFJbjBzSW5ObGNYVmxibU5sU1dRaU9qRTFmUT09In0="
    }
}`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(wantResp, gotResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestCrowdStrikeEndpointIncident(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/crowdstrike/endpoint-incident")
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

	req := &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "https://api.us-2.crowdstrike.com",
			Id:      "Test",
			Type:    "CrowdStrike-1.0.0",
			Config: []byte(`
			{
				"archived": false,
				"enabled": true,
				"apiVersion": "v1"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "EndpointIncident",
			ExternalId: "endpoint_protection_incident",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "incidentId",
					ExternalId: "incident_id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
					UniqueId:   true,
				},
				{
					Id:         "state",
					ExternalId: "state",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 1,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`{
    "success": {
        "objects": [
            {
                "attributes": [
                    {
                        "id": "incidentId",
                        "values": [
                            {
                                "stringValue": "inc:eca21da34c934e8e95c97a4f7af1d9a5:fede7474a2634f16997504abe3d21974"
                            }
                        ]
                    },
                    {
                        "id": "state",
                        "values": [
                            {
                                "stringValue": "closed"
                            }
                        ]
                    }
                ]
            }
        ],
        "nextCursor": "eyJjdXJzb3IiOiIxIn0="
    }
}`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(wantResp, gotResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestCrowdStrikeDevice(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/crowdstrike/device")
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

	req := &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "https://api.us-2.crowdstrike.com",
			Id:      "Test",
			Type:    "CrowdStrike-1.0.0",
			Config: []byte(`
			{
				"archived": false,
				"enabled": true,
				"apiVersion": "v1"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Device",
			ExternalId: "endpoint_protection_device",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "DeviceId",
					ExternalId: "device_id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
					UniqueId:   true,
				},
				{
					Id:         "DeviceType",
					ExternalId: "product_type_desc",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "ProvisionStatus",
					ExternalId: "provision_status",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "ModifiedAt",
					ExternalId: "modified_timestamp",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
				},
			},
		},
		PageSize: 1,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`{
    "success": {
        "objects": [
            {
                "attributes": [
                    {
                        "id": "DeviceId",
                        "values": [
                            {
                                "stringValue": "9b9b1e4f7512492f95f8039c065a28a9"
                            }
                        ]
                    },
                    {
                        "id": "ModifiedAt",
                        "values": [
                            {
                                "datetimeValue": {
                                    "timestamp": "2025-01-29T23:01:16Z"
                                }
                            }
                        ]
                    },
                    {
                        "id": "DeviceType",
                        "values": [
                            {
                                "stringValue": "Server"
                            }
                        ]
                    },
                    {
                        "id": "ProvisionStatus",
                        "values": [
                            {
                                "stringValue": "Provisioned"
                            }
                        ]
                    }
                ]
            }
        ],
        "nextCursor": "eyJjdXJzb3IiOiJGR2x1WTJ4MVpHVmZZMjl1ZEdWNGRGOTFkV2xrRG5GMVpYSjVWR2hsYmtabGRHTm9BaFp2YmtreFRUVnRhVlEzVTBZMVNteDBSbU4wYzNsbkFBQUFBQ1o4MFNNV2FrNVNkR3B2V2xsVVIyRlhja3N5Um10U1pHRXRkeFozY1dsNFFuRkRTRkZOWVRNMVkycDVSbkJNYmxOUkFBQUFBQ2JXOTlvV1Mzb3RjRGx2VFRsVVVrTk9hMWhtTjFaWlVWTkZRUT09In0="
    }
}`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(wantResp, gotResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestCrowdStrikeDetection(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/crowdstrike/detect")
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

	req := &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "https://api.us-2.crowdstrike.com",
			Id:      "Test",
			Type:    "CrowdStrike-1.0.0",
			Config: []byte(`
			{
				"archived": false,
				"enabled": true,
				"apiVersion": "v1"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Detection",
			ExternalId: "endpoint_protection_detect",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "DetectionId",
					ExternalId: "detection_id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
					UniqueId:   true,
				},
				{
					Id:         "EmailSent",
					ExternalId: "email_sent",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
				},
				{
					Id:         "Status",
					ExternalId: "status",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 2,
		Cursor:   "eyJjdXJzb3IiOiIyIn0=",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`{
    "success": {
        "objects": [
            {
                "attributes": [
                    {
                        "id": "DetectionId",
                        "values": [
                            {
                                "stringValue": "ldt:9b9b1e4f7512492f95f8039c065a28a9:1169567"
                            }
                        ]
                    },
                    {
                        "id": "EmailSent",
                        "values": [
                            {
                                "boolValue": false
                            }
                        ]
                    },
                    {
                        "id": "Status",
                        "values": [
                            {
                                "stringValue": "new"
                            }
                        ]
                    }
                ]
            },
            {
                "attributes": [
                    {
                        "id": "DetectionId",
                        "values": [
                            {
                                "stringValue": "ldt:9b9b1e4f7512492f95f8039c065a28a9:4295459139"
                            }
                        ]
                    },
                    {
                        "id": "EmailSent",
                        "values": [
                            {
                                "boolValue": false
                            }
                        ]
                    },
                    {
                        "id": "Status",
                        "values": [
                            {
                                "stringValue": "new"
                            }
                        ]
                    }
                ]
            }
        ],
        "nextCursor": "eyJjdXJzb3IiOiI0In0="
    }
}`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(wantResp, gotResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestCrowdStrikeCombinedAlert(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/crowdstrike/combinedAlert")
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

	req := &adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer {{OMITTED}}",
				},
			},
			Address: "https://api.us-2.crowdstrike.com",
			Id:      "Test",
			Type:    "CrowdStrike-1.0.0",
			Config: []byte(`
			{
				"archived": false,
				"enabled": true,
				"apiVersion": "v1"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Alert",
			ExternalId: "endpoint_protection_alert",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "CompositeId",
					ExternalId: "composite_id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
					UniqueId:   true,
				},
				{
					Id:         "AggregateId",
					ExternalId: "aggregate_id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "Status",
					ExternalId: "status",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					Id:         "files_accessed",
					ExternalId: "$.files_accessed",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "filename",
							ExternalId: "filename",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
							UniqueId:   false,
						},
						{
							Id:         "filepath",
							ExternalId: "filepath",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
					},
				},
				{
					Id:         "files_written",
					ExternalId: "$.files_written",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "filename",
							ExternalId: "filename",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
							UniqueId:   false,
						},
						{
							Id:         "filepath",
							ExternalId: "filepath",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
					},
				},
				{
					Id:         "mitre_attack",
					ExternalId: "$.mitre_attack",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "pattern_id",
							ExternalId: "pattern_id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
							List:       false,
							UniqueId:   true,
						},
					},
				},
			},
		},
		PageSize: 2,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`{
		"success": {
			"objects": [
			{
				"attributes": [
				{
					"id": "AggregateId",
					"values": [
					{
						"stringValue": "aggind:c36c42b64ce54b39a32e1d57240704c8:625985642613668398"
					}
					]
				},
				{
					"id": "CompositeId",
					"values": [
					{
						"stringValue": "8693deb4bf134cfb8855ee118d9a0243:ind:c36c42b64ce54b39a32e1d57240704c8:625985642593750673-20151-7049"
					}
					]
				},
				{
					"id": "Status",
					"values": [
					{
						"stringValue": "new"
					}
					]
				}
				],
				"childObjects": [
				{
					"entityId": "files_accessed",
					"objects": [
					{
						"attributes": [
						{
							"id": "filename",
							"values": [
							{
								"stringValue": "cat"
							}
							]
						},
						{
							"id": "filepath",
							"values": [
							{
								"stringValue": "/bin/"
							}
							]
						}
						]
					},
					{
						"attributes": [
						{
							"id": "filename",
							"values": [
							{
								"stringValue": "zshnW4W3l"
							}
							]
						},
						{
							"id": "filepath",
							"values": [
							{
								"stringValue": "/private/tmp/"
							}
							]
						}
						]
					}
					]
				},
				{
					"entityId": "files_written",
					"objects": [
					{
						"attributes": [
						{
							"id": "filename",
							"values": [
							{
								"stringValue": "eicar.com"
							}
							]
						},
						{
							"id": "filepath",
							"values": [
							{
								"stringValue": "/Users/joe/Desktop/"
							}
							]
						}
						]
					},
					{
						"attributes": [
						{
							"id": "filename",
							"values": [
							{
								"stringValue": "zshnW4W3l"
							}
							]
						},
						{
							"id": "filepath",
							"values": [
							{
								"stringValue": "/private/tmp/"
							}
							]
						}
						]
					}
					]
				},
				{
					"entityId": "mitre_attack",
					"objects": [
					{
						"attributes": [
						{
							"id": "pattern_id",
							"values": [
							{
								"int64Value": 20151
							}
							]
						}
						]
					}
					]
				}
				]
			},
			{
				"attributes": [
				{
					"id": "AggregateId",
					"values": [
					{
						"stringValue": "aggind:5388c592189444ad9e84df071c8f3954:8592364792"
					}
					]
				},
				{
					"id": "CompositeId",
					"values": [
					{
						"stringValue": "8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:12119912898-10304-117513744"
					}
					]
				},
				{
					"id": "Status",
					"values": [
					{
						"stringValue": "closed"
					}
					]
				}
				],
				"childObjects": [
				{
					"entityId": "files_accessed",
					"objects": [
					{
						"attributes": [
						{
							"id": "filename",
							"values": [
							{
								"stringValue": "cat"
							}
							]
						},
						{
							"id": "filepath",
							"values": [
							{
								"stringValue": "/bin/"
							}
							]
						}
						]
					},
					{
						"attributes": [
						{
							"id": "filename",
							"values": [
							{
								"stringValue": "zshnW4W3l"
							}
							]
						},
						{
							"id": "filepath",
							"values": [
							{
								"stringValue": "/private/tmp/"
							}
							]
						}
						]
					}
					]
				},
				{
					"entityId": "files_written",
					"objects": [
					{
						"attributes": [
						{
							"id": "filename",
							"values": [
							{
								"stringValue": "eicar.com"
							}
							]
						},
						{
							"id": "filepath",
							"values": [
							{
								"stringValue": "/Users/joe/Desktop/"
							}
							]
						}
						]
					},
					{
						"attributes": [
						{
							"id": "filename",
							"values": [
							{
								"stringValue": "zshnW4W3l"
							}
							]
						},
						{
							"id": "filepath",
							"values": [
							{
								"stringValue": "/private/tmp/"
							}
							]
						}
						]
					}
					]
				},
				{
					"entityId": "mitre_attack",
					"objects": [
					{
						"attributes": [
						{
							"id": "pattern_id",
							"values": [
							{
								"int64Value": 10304
							}
							]
						}
						]
					}
					]
				}
				]
			}
			],
			"nextCursor": "eyJjdXJzb3IiOiJleUoyWlhKemFXOXVJam9pZGpFaUxDSjBiM1JoYkY5b2FYUnpJam95TXl3aWRHOTBZV3hmY21Wc1lYUnBiMjRpT2lKbGNTSXNJbU5zZFhOMFpYSmZhV1FpT2lKMFpYTjBJaXdpWVdaMFpYSWlPbHN4TnpRNU5qRXhNVFUzTWpJeExDSjBaWE4wYVdRNmFXNWtPalV6T0Roak5Ua3lNVGc1TkRRMFlXUTVaVGcwWkdZd056RmpPR1l6T1RVME9qazNPREkzT0RJMk1UUXRNVEF6TURNdE16RTRNekUxTmpnaVhTd2lkRzkwWVd4ZlptVjBZMmhsWkNJNk1uMD0ifQ=="
		}
}`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(wantResp, gotResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}
