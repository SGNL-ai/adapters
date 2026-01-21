// Copyright 2026 SGNL.ai, Inc.

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

func TestWorkdayAdapter_Worker(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/workday/worker")
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
			Address: "test-instance.workday.com",
			Id:      "Test",
			Type:    "Workday-1.0.0",
			Config:  []byte(`{"apiVersion":"v1","organizationId":"{{OMITTED}}"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Worker",
			ExternalId: "allWorkers",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "workerId",
					ExternalId: "$.worker.id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
					UniqueId:   true,
				},
				{
					Id:         "workerDescriptor",
					ExternalId: "$.worker.descriptor",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "employeeID",
					ExternalId: "employeeID",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "workerActive",
					ExternalId: "workerActive",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "managementLevelDescriptor",
					ExternalId: "$.managementLevel.descriptor",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "managementLevelId",
					ExternalId: "$.managementLevel.id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "companyDescriptor",
					ExternalId: "$.company.descriptor",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "companyId",
					ExternalId: "$.company.id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "genderDescriptor",
					ExternalId: "$.gender.descriptor",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "genderId",
					ExternalId: "$.gender.id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "hireDate",
					ExternalId: "hireDate",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "FTE",
					ExternalId: "FTE",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "positionID",
					ExternalId: "positionID",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "jobTitle",
					ExternalId: "jobTitle",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					Id:         "email_Work",
					ExternalId: "email_Work",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "id",
							ExternalId: "id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "descriptor",
							ExternalId: "descriptor",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
					},
				},
				{
					Id:         "employeeType",
					ExternalId: "employeeType",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "id",
							ExternalId: "id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "descriptor",
							ExternalId: "descriptor",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
					},
				},
			},
		},
		PageSize: 3,
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
							"values": [
								{
									"string_value": "Global Modern Services, Inc. (USA)"
								}
							],
							"id": "companyDescriptor"
						},
						{
							"values": [
								{
									"string_value": "cb550da820584750aae8f807882fa79a"
								}
							],
							"id": "companyId"
						},
						{
							"values": [
								{
									"string_value": "Female"
								}
							],
							"id": "genderDescriptor"
						},
						{
							"values": [
								{
									"string_value": "9cce3bec2d0d420283f76f51b928d885"
								}
							],
							"id": "genderId"
						},
						{
							"values": [
								{
									"string_value": "4 Vice President"
								}
							],
							"id": "managementLevelDescriptor"
						},
						{
							"values": [
								{
									"string_value": "679d4d1ac6da40e19deb7d91e170431d"
								}
							],
							"id": "managementLevelId"
						},
						{
							"values": [
								{
									"string_value": "{{OMITTED}}"
								}
							],
							"id": "workerDescriptor"
						},
						{
							"values": [
								{
									"string_value": "3aa5550b7fe348b98d7b5741afc65534"
								}
							],
							"id": "workerId"
						},
						{
							"values": [
								{
									"string_value": "1"
								}
							],
							"id": "FTE"
						},
						{
							"values": [
								{
									"string_value": "21001"
								}
							],
							"id": "employeeID"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2000-01-01T00:00:00Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "hireDate"
						},
						{
							"values": [
								{
									"string_value": "Vice President, Human Resources"
								}
							],
							"id": "jobTitle"
						},
						{
							"values": [
								{
									"string_value": "P-00004"
								}
							],
							"id": "positionID"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "workerActive"
						}
					],
					"child_objects": [
						{
							"objects": [
								{
									"attributes": [
										{
											"values": [
												{
													"string_value": "{{OMITTED}}@workday.net"
												}
											],
											"id": "descriptor"
										},
										{
											"values": [
												{
													"string_value": "f22d772555044fd29f620ff7ff2b0830"
												}
											],
											"id": "id"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "email_Work"
						},
						{
							"objects": [
								{
									"attributes": [
										{
											"values": [
												{
													"string_value": "Regular"
												}
											],
											"id": "descriptor"
										},
										{
											"values": [
												{
													"string_value": "9459f5e6f1084433b767c7901ec04416"
												}
											],
											"id": "id"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "employeeType"
						}
					]
				},
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "Global Modern Services, Inc. (USA)"
								}
							],
							"id": "companyDescriptor"
						},
						{
							"values": [
								{
									"string_value": "cb550da820584750aae8f807882fa79a"
								}
							],
							"id": "companyId"
						},
						{
							"values": [
								{
									"string_value": "Not Declared"
								}
							],
							"id": "genderDescriptor"
						},
						{
							"values": [
								{
									"string_value": "a14bf6afa9204ff48a8ea353dd71eb22"
								}
							],
							"id": "genderId"
						},
						{
							"values": [
								{
									"string_value": "2 Chief Executive Officer"
								}
							],
							"id": "managementLevelDescriptor"
						},
						{
							"values": [
								{
									"string_value": "3de1f2834f064394a40a40a727fb6c6d"
								}
							],
							"id": "managementLevelId"
						},
						{
							"values": [
								{
									"string_value": "{{OMITTED}}"
								}
							],
							"id": "workerDescriptor"
						},
						{
							"values": [
								{
									"string_value": "0e44c92412d34b01ace61e80a47aaf6d"
								}
							],
							"id": "workerId"
						},
						{
							"values": [
								{
									"string_value": "1"
								}
							],
							"id": "FTE"
						},
						{
							"values": [
								{
									"string_value": "21002"
								}
							],
							"id": "employeeID"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2000-01-01T00:00:00Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "hireDate"
						},
						{
							"values": [
								{
									"string_value": "Chief Executive Officer"
								}
							],
							"id": "jobTitle"
						},
						{
							"values": [
								{
									"string_value": "P-00001"
								}
							],
							"id": "positionID"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "workerActive"
						}
					],
					"child_objects": [
						{
							"objects": [
								{
									"attributes": [
										{
											"values": [
												{
													"string_value": "{{OMITTED}}@workdaySJTest.net"
												}
											],
											"id": "descriptor"
										},
										{
											"values": [
												{
													"string_value": "d7fef59db8e21001de457203a69e0001"
												}
											],
											"id": "id"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "email_Work"
						},
						{
							"objects": [
								{
									"attributes": [
										{
											"values": [
												{
													"string_value": "Regular"
												}
											],
											"id": "descriptor"
										},
										{
											"values": [
												{
													"string_value": "9459f5e6f1084433b767c7901ec04416"
												}
											],
											"id": "id"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "employeeType"
						}
					]
				},
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "Global Modern Services, Inc. (USA)"
								}
							],
							"id": "companyDescriptor"
						},
						{
							"values": [
								{
									"string_value": "cb550da820584750aae8f807882fa79a"
								}
							],
							"id": "companyId"
						},
						{
							"values": [
								{
									"string_value": "Male"
								}
							],
							"id": "genderDescriptor"
						},
						{
							"values": [
								{
									"string_value": "d3afbf8074e549ffb070962128e1105a"
								}
							],
							"id": "genderId"
						},
						{
							"values": [
								{
									"string_value": "3 Executive Vice President"
								}
							],
							"id": "managementLevelDescriptor"
						},
						{
							"values": [
								{
									"string_value": "0ceb3292987b474bbc40c751a1e22c69"
								}
							],
							"id": "managementLevelId"
						},
						{
							"values": [
								{
									"string_value": "{{OMITTED}}"
								}
							],
							"id": "workerDescriptor"
						},
						{
							"values": [
								{
									"string_value": "3895af7993ff4c509cbea2e1817172e0"
								}
							],
							"id": "workerId"
						},
						{
							"values": [
								{
									"string_value": "1"
								}
							],
							"id": "FTE"
						},
						{
							"values": [
								{
									"string_value": "21003"
								}
							],
							"id": "employeeID"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2000-01-01T00:00:00Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "hireDate"
						},
						{
							"values": [
								{
									"string_value": "Chief Information Officer"
								}
							],
							"id": "jobTitle"
						},
						{
							"values": [
								{
									"string_value": "P-00002"
								}
							],
							"id": "positionID"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "workerActive"
						}
					],
					"child_objects": [
						{
							"objects": [
								{
									"attributes": [
										{
											"values": [
												{
													"string_value": "{{OMITTED}}@workday.net"
												}
											],
											"id": "descriptor"
										},
										{
											"values": [
												{
													"string_value": "d7fef59db8e21001dddaa607a7d30001"
												}
											],
											"id": "id"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "email_Work"
						},
						{
							"objects": [
								{
									"attributes": [
										{
											"values": [
												{
													"string_value": "Regular"
												}
											],
											"id": "descriptor"
										},
										{
											"values": [
												{
													"string_value": "9459f5e6f1084433b767c7901ec04416"
												}
											],
											"id": "id"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "employeeType"
						}
					]
				}
			],
			"next_cursor": "eyJjdXJzb3IiOjN9"
		}
	}`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}
