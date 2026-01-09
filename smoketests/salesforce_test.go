// Copyright 2026 SGNL.ai, Inc.
package smoketests

import (
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	adapter_api_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/smoketests/common"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestSalesforceAdapter_Case(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/salesforce/case")
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
			Address: "test-instance.my.salesforce.com",
			Id:      "Salesforce",
			Type:    "Salesforce-1.0.1",
			Config:  []byte(`{"apiVersion":"58.0"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "SalesforceCase",
			ExternalId: "Case",
			Ordered:    true,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Id",
					ExternalId: "Id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "CaseNumber",
					ExternalId: "CaseNumber",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
			ChildEntities: nil,
		},
		PageSize: 200,
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
								"id": "CaseNumber",
								"values": [
									{
										"string_value": "00001001"
									}
								]
							},
							{
								"id": "Id",
								"values": [
									{
										"string_value": "5008Z00001xKVNCQA4"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "CaseNumber",
								"values": [
									{
										"string_value": "00001000"
									}
								]
							},
							{
								"id": "Id",
								"values": [
									{
										"string_value": "5008Z00001xKVNDQA4"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "CaseNumber",
								"values": [
									{
										"string_value": "00001002"
									}
								]
							},
							{
								"id": "Id",
								"values": [
									{
										"string_value": "5008Z00001xKVNEQA4"
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

func TestSalesforceAdapter_User(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/salesforce/user")
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
			Address: "test-instance.my.salesforce.com",
			Id:      "Salesforce",
			Type:    "Salesforce-1.0.1",
			Config:  []byte(`{"apiVersion":"58.0"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "SalesforceUser",
			ExternalId: "User",
			Ordered:    true,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Id",
					ExternalId: "Id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "Name",
					ExternalId: "Name",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
			ChildEntities: nil,
		},
		PageSize: 200,
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
								"id": "Id",
								"values": [
									{
										"string_value": "0058Z000008hrrgQAA"
									}
								]
							},
							{
								"id": "Name",
								"values": [
									{
										"string_value": "SGNL Integration"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "Id",
								"values": [
									{
										"string_value": "0058Z000008opHuQAI"
									}
								]
							},
							{
								"id": "Name",
								"values": [
									{
										"string_value": "Nick"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "Id",
								"values": [
									{
										"string_value": "0058Z000008szVuQAI"
									}
								]
							},
							{
								"id": "Name",
								"values": [
									{
										"string_value": "Integration User"
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

func TestSalesforceAdapter_CustomFields(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/salesforce/account_customfield")
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
			Address: "test-instance.my.salesforce.com",
			Id:      "Salesforce",
			Type:    "Salesforce-1.0.1",
			Config:  []byte(`{"apiVersion":"58.0"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "SalesforceAccount",
			ExternalId: "Account",
			Ordered:    true,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Id",
					ExternalId: "Id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "Name",
					ExternalId: "$.Name",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "Active",
					ExternalId: "$.Active__c",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
			ChildEntities: nil,
		},
		PageSize: 200,
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
								"id": "Active",
								"values": [
									{
										"string_value": "Yes"
									}
								]
							},
							{
								"id": "Name",
								"values": [
									{
										"string_value": "Edge Communications"
									}
								]
							},
							{
								"id": "Id",
								"values": [
									{
										"string_value": "001gL00000WW0iPQAT"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "Name",
								"values": [
									{
										"string_value": "Burlington Textiles Corp of America"
									}
								]
							},
							{
								"id": "Id",
								"values": [
									{
										"string_value": "001gL00000WW0iQQAT"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "Active",
								"values": [
									{
										"string_value": "Yes"
									}
								]
							},
							{
								"id": "Name",
								"values": [
									{
										"string_value": "Pyramid Construction Inc."
									}
								]
							},
							{
								"id": "Id",
								"values": [
									{
										"string_value": "001gL00000WW0iRQAT"
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

func TestSalesforceAdapter_MultiLevelRelationship(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/salesforce/account_relationship")
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
			Address: "test-instance.my.salesforce.com",
			Id:      "Salesforce",
			Type:    "Salesforce-1.0.1",
			Config:  []byte(`{"apiVersion":"58.0"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "SalesforceAccount",
			ExternalId: "Account",
			Ordered:    true,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Id",
					ExternalId: "Id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "Name",
					ExternalId: "$.Name",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "OwnerName",
					ExternalId: "$.Owner.Name",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
			ChildEntities: nil,
		},
		PageSize: 200,
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
								"id": "Name",
								"values": [
									{
										"string_value": "Edge Communications"
									}
								]
							},
							{
								"id": "OwnerName",
								"values": [
									{
										"string_value": "OrgFarm EPIC"
									}
								]
							},
							{
								"id": "Id",
								"values": [
									{
										"string_value": "001gL00000WW0iPQAT"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "Name",
								"values": [
									{
										"string_value": "Burlington Textiles Corp of America"
									}
								]
							},
							{
								"id": "OwnerName",
								"values": [
									{
										"string_value": "OrgFarm EPIC"
									}
								]
							},
							{
								"id": "Id",
								"values": [
									{
										"string_value": "001gL00000WW0iQQAT"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "Name",
								"values": [
									{
										"string_value": "Pyramid Construction Inc."
									}
								]
							},
							{
								"id": "OwnerName",
								"values": [
									{
										"string_value": "OrgFarm EPIC"
									}
								]
							},
							{
								"id": "Id",
								"values": [
									{
										"string_value": "001gL00000WW0iRQAT"
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

func TestSalesforceAdapter_FiveLevelRelationship(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/salesforce/user_fivelevelrelationship")
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
			Address: "test-instance.my.salesforce.com",
			Id:      "Salesforce",
			Type:    "Salesforce-1.0.1",
			Config:  []byte(`{"apiVersion":"58.0"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "SalesforceUser",
			ExternalId: "User",
			Ordered:    true,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Id",
					ExternalId: "Id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "Name",
					ExternalId: "$.Name",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "ManagerManagerManagerManagerName",
					ExternalId: "$.Manager.Manager.Manager.Manager.Name",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
			ChildEntities: nil,
		},
		PageSize: 200,
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
								"id": "ManagerManagerManagerManagerName",
								"values": [
									{
										"string_value": "TopLevel Manager5"
									}
								]
							},
							{
								"id": "Name",
								"values": [
									{
										"string_value": "Level Manager1"
									}
								]
							},
							{
								"id": "Id",
								"values": [
									{
										"string_value": "005gL00000BHeFdQAL"
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

// sortChildObjectsByID sorts child objects within a GetPageResponse by their "id" attribute
// to enable order-independent comparison in tests.
func sortChildObjectsByID(resp *adapter_api_v1.GetPageResponse) {
	if resp == nil || resp.GetSuccess() == nil {
		return
	}

	for _, obj := range resp.GetSuccess().GetObjects() {
		for _, childEntity := range obj.GetChildObjects() {
			objects := childEntity.GetObjects()
			sort.Slice(objects, func(i, j int) bool {
				var id1, id2 string
				for _, attr := range objects[i].GetAttributes() {
					if attr.GetId() == "id" && len(attr.GetValues()) > 0 {
						id1 = attr.GetValues()[0].GetStringValue()

						break
					}
				}
				for _, attr := range objects[j].GetAttributes() {
					if attr.GetId() == "id" && len(attr.GetValues()) > 0 {
						id2 = attr.GetValues()[0].GetStringValue()

						break
					}
				}

				return id1 < id2
			})
		}
	}
}

func TestSalesforceAdapter_MultiSelectPicklist(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/salesforce/account_multiselectpicklist")
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
			Address: "test-instance.my.salesforce.com",
			Id:      "Salesforce",
			Type:    "Salesforce-1.0.1",
			Config:  []byte(`{"apiVersion":"58.0"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "SalesforceAccount",
			ExternalId: "Account",
			Ordered:    true,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Id",
					ExternalId: "Id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "Name",
					ExternalId: "Name",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					Id:         "Locations",
					ExternalId: "Locations__c",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "id",
							ExternalId: "id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							UniqueId:   true,
						},
						{
							Id:         "value",
							ExternalId: "value",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
				},
			},
		},
		PageSize: 200,
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
								"id": "Id",
								"values": [
									{
										"string_value": "001gL00000XgJxhQAF"
									}
								]
							},
							{
								"id": "Name",
								"values": [
									{
										"string_value": "TestMultiFieldPicklist"
									}
								]
							}
						],
						"child_objects": [
							{
								"entity_id": "Locations",
								"objects": [
									{
										"attributes": [
											{
												"id": "id",
												"values": [
													{
														"string_value": "001gL00000XgJxhQAF_Locations__c_chicago"
													}
												]
											},
											{
												"id": "value",
												"values": [
													{
														"string_value": "Chicago"
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
														"string_value": "001gL00000XgJxhQAF_Locations__c_new-york"
													}
												]
											},
											{
												"id": "value",
												"values": [
													{
														"string_value": "New York"
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
														"string_value": "001gL00000XgJxhQAF_Locations__c_seattle"
													}
												]
											},
											{
												"id": "value",
												"values": [
													{
														"string_value": "Seattle"
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
								"id": "Id",
								"values": [
									{
										"string_value": "001gL00000XjjndQAB"
									}
								]
							},
							{
								"id": "Name",
								"values": [
									{
										"string_value": "TestMultiFieldPicklist-1"
									}
								]
							}
						],
						"child_objects": [
							{
								"entity_id": "Locations",
								"objects": [
									{
										"attributes": [
											{
												"id": "id",
												"values": [
													{
														"string_value": "001gL00000XjjndQAB_Locations__c_chicago"
													}
												]
											},
											{
												"id": "value",
												"values": [
													{
														"string_value": "Chicago"
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
								"id": "Id",
								"values": [
									{
										"string_value": "001gL00000WW0iRQAT"
									}
								]
							},
							{
								"id": "Name",
								"values": [
									{
										"string_value": "Pyramid Construction Inc."
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

	// Sort child objects by ID for order-independent comparison
	sortChildObjectsByID(gotResp)
	sortChildObjectsByID(wantResp)

	if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}
