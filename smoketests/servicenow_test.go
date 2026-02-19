// Copyright 2026 SGNL.ai, Inc.

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

func TestServicenowAdapter_Case(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/servicenow/case")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "test-instance-username",
						Password: "{{OMITTED}}",
					},
				},
			},
			Address: "test-instance.service-now.com",
			Id:      "ServiceNow",
			Type:    "ServiceNow-1.0.1",
			Config:  []byte(`{"apiVersion":"v2"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "ServicenowCase",
			ExternalId: "sn_customerservice_case",
			Ordered:    true,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "sys_id",
					ExternalId: "sys_id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "number",
					ExternalId: "number",
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
								"id": "number",
								"values": [
									{
										"string_value": "CSE0014483"
									}
								]
							},
							{
								"id": "sys_id",
								"values": [
									{
										"string_value": "0b6027408737c910793a97983cbb3514"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "number",
								"values": [
									{
										"string_value": "CSE0019923"
									}
								]
							},
							{
								"id": "sys_id",
								"values": [
									{
										"string_value": "0c702b8487774910b8137448dabb3564"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "number",
								"values": [
									{
										"string_value": "CSE0010140"
									}
								]
							},
							{
								"id": "sys_id",
								"values": [
									{
										"string_value": "0e70a78487774910b8137448dabb3520"
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

func TestServicenowAdapter_User(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/servicenow/user")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "test-instance-username",
						Password: "{{OMITTED}}",
					},
				},
			},
			Address: "test-instance.service-now.com",
			Id:      "ServiceNow",
			Type:    "ServiceNow-1.0.1",
			Config:  []byte(`{"apiVersion":"v2"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "ServicenowUser",
			ExternalId: "sys_user",
			Ordered:    true,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "sys_id",
					ExternalId: "sys_id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "email",
					ExternalId: "email",
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
								"id": "email",
								"values": [
									{
										"string_value": "jennifer.zamora@wholesalechips.co"
									}
								]
							},
							{
								"id": "sys_id",
								"values": [
									{
										"string_value": "0010af4487774910b8137448dabb351c"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "email",
								"values": [
									{
										"string_value": "noah.villanueva@wholesalechips.co"
									}
								]
							},
							{
								"id": "sys_id",
								"values": [
									{
										"string_value": "0110eb008737c910793a97983cbb3507"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "email",
								"values": [
									{
										"string_value": "romain@sgnl.ai"
									}
								]
							},
							{
								"id": "sys_id",
								"values": [
									{
										"string_value": "016a6dd987634110b8137448dabb3513"
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

var advancedFiltersConfig = `
{
  "requestTimeoutSeconds": 10,
  "apiVersion": "v2",
  "advancedFilters": {
    "getObjectsByScope": {
      "sys_user_group": [
        {
          "scopeEntity": "sys_user_group",
          "scopeEntityFilter": "sys_id=0270c251c3200200be647bfaa2d3aea6",
          "members": [
            {
              "memberEntity": "sys_user",
              "memberEntityFilter": "user.active=true",
              "relatedEntities": [
                {
                  "relatedEntity": "change_task",
                  "relatedEntityFilter": "assigned_toIN{$.sys_user.sys_id}"
                }
              ]
            },
            {
              "memberEntity": "sys_user",
              "memberEntityFilter": "user.first_name=Richard",
              "relatedEntities": [
                {
                  "relatedEntity": "change_task",
                  "relatedEntityFilter": "assigned_toIN{$.sys_user.sys_id}"
                }
              ]
            }
          ],
          "relatedEntities": [
            {
              "relatedEntity": "change_task",
              "relatedEntityFilter": "assignment_groupIN{$.sys_group.sys_id}"
            }
          ]
        },
        {
          "scopeEntity": "sys_user_group",
          "scopeEntityFilter": "sys_id=019ad92ec7230010393d265c95c260dd",
          "members": [
            {
              "memberEntity": "sys_user",
              "memberEntityFilter": "user.active=true",
              "relatedEntities": [
                {
                  "relatedEntity": "change_task",
                  "relatedEntityFilter": "assigned_toIN{$.sys_user.sys_id}"
                }
              ]
            }
          ]
        }
      ]
    }
  }
}
`

// nolint: lll
func TestServiceNowAdapter_User_AdancedFilters(t *testing.T) {
	fixtureFile := "fixtures/servicenow/advanced_filters/user"

	pageSize1Pages := []struct {
		wantResp string
	}{
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Damon"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "b2459a268727d110b8137448dabb350a"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6eyJjdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uuc2VydmljZS1ub3cuY29tL2FwaS9ub3cvdjIvdGFibGUvc3lzX3VzZXJfZ3JtZW1iZXI/c3lzcGFybV9maWVsZHM9c3lzX2lkLGZpcnN0X25hbWUsdXNlci5zeXNfaWQsdXNlci5maXJzdF9uYW1lXHUwMDI2c3lzcGFybV9leGNsdWRlX3JlZmVyZW5jZV9saW5rPXRydWVcdTAwMjZzeXNwYXJtX2xpbWl0PTFcdTAwMjZzeXNwYXJtX3F1ZXJ5PWdyb3VwSU4wMjcwYzI1MWMzMjAwMjAwYmU2NDdiZmFhMmQzYWVhNiU1RXVzZXIuYWN0aXZlJTNEdHJ1ZSU1RU9SREVSQllzeXNfaWRcdTAwMjZzeXNwYXJtX29mZnNldD0xIn19fQ=="
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Joe"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "53eae4e2c33992102177bc159901317f"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6eyJjdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uuc2VydmljZS1ub3cuY29tL2FwaS9ub3cvdjIvdGFibGUvc3lzX3VzZXJfZ3JtZW1iZXI/c3lzcGFybV9maWVsZHM9c3lzX2lkLGZpcnN0X25hbWUsdXNlci5zeXNfaWQsdXNlci5maXJzdF9uYW1lXHUwMDI2c3lzcGFybV9leGNsdWRlX3JlZmVyZW5jZV9saW5rPXRydWVcdTAwMjZzeXNwYXJtX2xpbWl0PTFcdTAwMjZzeXNwYXJtX3F1ZXJ5PWdyb3VwSU4wMjcwYzI1MWMzMjAwMjAwYmU2NDdiZmFhMmQzYWVhNiU1RXVzZXIuYWN0aXZlJTNEdHJ1ZSU1RU9SREVSQllzeXNfaWRcdTAwMjZzeXNwYXJtX29mZnNldD0yIn19fQ=="
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Alejandro"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "5b00ab008737c910793a97983cbb354a"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6eyJjdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uuc2VydmljZS1ub3cuY29tL2FwaS9ub3cvdjIvdGFibGUvc3lzX3VzZXJfZ3JtZW1iZXI/c3lzcGFybV9maWVsZHM9c3lzX2lkLGZpcnN0X25hbWUsdXNlci5zeXNfaWQsdXNlci5maXJzdF9uYW1lXHUwMDI2c3lzcGFybV9leGNsdWRlX3JlZmVyZW5jZV9saW5rPXRydWVcdTAwMjZzeXNwYXJtX2xpbWl0PTFcdTAwMjZzeXNwYXJtX3F1ZXJ5PWdyb3VwSU4wMjcwYzI1MWMzMjAwMjAwYmU2NDdiZmFhMmQzYWVhNiU1RXVzZXIuYWN0aXZlJTNEdHJ1ZSU1RU9SREVSQllzeXNfaWRcdTAwMjZzeXNwYXJtX29mZnNldD0zIn19fQ=="
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Linda"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "ab00ab008737c910793a97983cbb354d"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6eyJjdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uuc2VydmljZS1ub3cuY29tL2FwaS9ub3cvdjIvdGFibGUvc3lzX3VzZXJfZ3JtZW1iZXI/c3lzcGFybV9maWVsZHM9c3lzX2lkLGZpcnN0X25hbWUsdXNlci5zeXNfaWQsdXNlci5maXJzdF9uYW1lXHUwMDI2c3lzcGFybV9leGNsdWRlX3JlZmVyZW5jZV9saW5rPXRydWVcdTAwMjZzeXNwYXJtX2xpbWl0PTFcdTAwMjZzeXNwYXJtX3F1ZXJ5PWdyb3VwSU4wMjcwYzI1MWMzMjAwMjAwYmU2NDdiZmFhMmQzYWVhNiU1RXVzZXIuYWN0aXZlJTNEdHJ1ZSU1RU9SREVSQllzeXNfaWRcdTAwMjZzeXNwYXJtX29mZnNldD00In19fQ=="
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Lisa"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "1f00ab008737c910793a97983cbb3544"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6eyJjdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uuc2VydmljZS1ub3cuY29tL2FwaS9ub3cvdjIvdGFibGUvc3lzX3VzZXJfZ3JtZW1iZXI/c3lzcGFybV9maWVsZHM9c3lzX2lkLGZpcnN0X25hbWUsdXNlci5zeXNfaWQsdXNlci5maXJzdF9uYW1lXHUwMDI2c3lzcGFybV9leGNsdWRlX3JlZmVyZW5jZV9saW5rPXRydWVcdTAwMjZzeXNwYXJtX2xpbWl0PTFcdTAwMjZzeXNwYXJtX3F1ZXJ5PWdyb3VwSU4wMjcwYzI1MWMzMjAwMjAwYmU2NDdiZmFhMmQzYWVhNiU1RXVzZXIuYWN0aXZlJTNEdHJ1ZSU1RU9SREVSQllzeXNfaWRcdTAwMjZzeXNwYXJtX29mZnNldD01In19fQ=="
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Joe"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "51104a4987e035d0793a97983cbb35c7"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjEsImN1cnNvciI6bnVsbH19"
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [],
					"next_cursor": "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MSwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6bnVsbH19"
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Aldo"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "aa4564cd93500e1089e6f7d86cba101c"
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

	pageSize10Pages := []struct {
		wantResp string
	}{
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Damon"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "b2459a268727d110b8137448dabb350a"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						},
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Joe"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "53eae4e2c33992102177bc159901317f"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						},
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Alejandro"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "5b00ab008737c910793a97983cbb354a"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						},
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Linda"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "ab00ab008737c910793a97983cbb354d"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						},
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Lisa"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "1f00ab008737c910793a97983cbb3544"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						},
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Joe"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "51104a4987e035d0793a97983cbb35c7"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjEsImN1cnNvciI6bnVsbH19"
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [],
					"next_cursor": "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MSwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6bnVsbH19"
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Aldo"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "aa4564cd93500e1089e6f7d86cba101c"
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

	pagesMap := map[int64][]struct {
		wantResp string
	}{
		1:  pageSize1Pages,
		10: pageSize10Pages,
	}

	// Test with different page sizes.
	for pageSize, pages := range pagesMap {
		for pageNumber, page := range pages {
			fixtureName := fmt.Sprintf("%s_pagesize_%d_page_%d", fixtureFile, pageSize, pageNumber+1)

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
				if err := protojson.Unmarshal([]byte(pages[pageNumber-1].wantResp), previousResp); err != nil {
					t.Fatal(err)
				}

				cursor = previousResp.GetSuccess().NextCursor
			}

			gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
				Datasource: &adapter_api_v1.DatasourceConfig{
					Auth: &adapter_api_v1.DatasourceAuthCredentials{
						AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
							Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
								Username: "test-instance-username",
								Password: "{{OMITTED}}",
							},
						},
					},
					Address: "test-instance.service-now.com",
					Id:      "ServiceNow",
					Type:    "ServiceNow-1.0.1",
					Config:  []byte(advancedFiltersConfig),
				},
				Entity: &adapter_api_v1.EntityConfig{
					Id:         "User",
					ExternalId: "sys_user",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6",
							ExternalId: "sys_id",
							UniqueId:   true,
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "ce0954e6-cc24-485d-b095-d5f0ee39c68c",
							ExternalId: "first_name",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
				},
				PageSize: pageSize,
				Cursor:   cursor,
			})
			if err != nil {
				t.Fatal(err)
			}

			wantResp := new(adapter_api_v1.GetPageResponse)

			if err := protojson.Unmarshal([]byte(page.wantResp), wantResp); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
				t.Fatal(diff)
			}
		}
	}
}

// nolint: lll
func TestServiceNowAdapter_Group_AdvancedFilters(t *testing.T) {
	fixtureFile := "fixtures/servicenow/advanced_filters/group"

	pageSize1Pages := []struct {
		wantResp string
	}{
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Customer Service Support"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "0270c251c3200200be647bfaa2d3aea6"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MSwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6bnVsbH19"
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Analytics Settings Managers"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "019ad92ec7230010393d265c95c260dd"
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

	pageSize10Pages := []struct {
		wantResp string
	}{
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Customer Service Support"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "0270c251c3200200be647bfaa2d3aea6"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MSwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6bnVsbH19"
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "Analytics Settings Managers"
										}
									],
									"id": "ce0954e6-cc24-485d-b095-d5f0ee39c68c"
								},
								{
									"values": [
										{
											"string_value": "019ad92ec7230010393d265c95c260dd"
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

	pagesMap := map[int64][]struct {
		wantResp string
	}{
		1:  pageSize1Pages,
		10: pageSize10Pages,
	}

	// Test with different page sizes.
	for pageSize, pages := range pagesMap {
		for pageNumber, page := range pages {
			fixtureName := fmt.Sprintf("%s_pagesize_%d_page_%d", fixtureFile, pageSize, pageNumber+1)

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
				if err := protojson.Unmarshal([]byte(pages[pageNumber-1].wantResp), previousResp); err != nil {
					t.Fatal(err)
				}

				cursor = previousResp.GetSuccess().NextCursor
			}

			gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
				Datasource: &adapter_api_v1.DatasourceConfig{
					Auth: &adapter_api_v1.DatasourceAuthCredentials{
						AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
							Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
								Username: "test-instance-username",
								Password: "{{OMITTED}}",
							},
						},
					},
					Address: "test-instance.service-now.com",
					Id:      "ServiceNow",
					Type:    "ServiceNow-1.0.1",
					Config:  []byte(advancedFiltersConfig),
				},
				Entity: &adapter_api_v1.EntityConfig{
					Id:         "Group",
					ExternalId: "sys_user_group",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6",
							ExternalId: "sys_id",
							UniqueId:   true,
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "ce0954e6-cc24-485d-b095-d5f0ee39c68c",
							ExternalId: "name",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
				},
				PageSize: pageSize,
				Cursor:   cursor,
			})
			if err != nil {
				t.Fatal(err)
			}

			wantResp := new(adapter_api_v1.GetPageResponse)

			if err := protojson.Unmarshal([]byte(page.wantResp), wantResp); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
				t.Fatal(diff)
			}
		}
	}
}

// nolint: lll
func TestServiceNowAdapter_ChangeTask_AdvancedFilters(t *testing.T) {
	fixtureFile := "fixtures/servicenow/advanced_filters/changetask"

	pageSize1Pages := []struct {
		wantResp string
	}{
		{
			wantResp: `
			{
				"success": {
					"objects": [],
					"next_cursor": "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjowLCJyZWxhdGVkRW50aXR5Q3Vyc29yIjp7ImN1cnNvciI6Imh0dHBzOi8vdGVzdC1pbnN0YW5jZS5zZXJ2aWNlLW5vdy5jb20vYXBpL25vdy92Mi90YWJsZS9zeXNfdXNlcl9ncm1lbWJlcj9zeXNwYXJtX2ZpZWxkcz1zeXNfaWQsdXNlci5zeXNfaWRcdTAwMjZzeXNwYXJtX2V4Y2x1ZGVfcmVmZXJlbmNlX2xpbms9dHJ1ZVx1MDAyNnN5c3Bhcm1fbGltaXQ9MVx1MDAyNnN5c3Bhcm1fcXVlcnk9Z3JvdXBJTjAyNzBjMjUxYzMyMDAyMDBiZTY0N2JmYWEyZDNhZWE2JTVFdXNlci5hY3RpdmUlM0R0cnVlJTVFT1JERVJCWXN5c19pZFx1MDAyNnN5c3Bhcm1fb2Zmc2V0PTEifX19"
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [],
					"next_cursor": "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjowLCJyZWxhdGVkRW50aXR5Q3Vyc29yIjp7ImN1cnNvciI6Imh0dHBzOi8vdGVzdC1pbnN0YW5jZS5zZXJ2aWNlLW5vdy5jb20vYXBpL25vdy92Mi90YWJsZS9zeXNfdXNlcl9ncm1lbWJlcj9zeXNwYXJtX2ZpZWxkcz1zeXNfaWQsdXNlci5zeXNfaWRcdTAwMjZzeXNwYXJtX2V4Y2x1ZGVfcmVmZXJlbmNlX2xpbms9dHJ1ZVx1MDAyNnN5c3Bhcm1fbGltaXQ9MVx1MDAyNnN5c3Bhcm1fcXVlcnk9Z3JvdXBJTjAyNzBjMjUxYzMyMDAyMDBiZTY0N2JmYWEyZDNhZWE2JTVFdXNlci5hY3RpdmUlM0R0cnVlJTVFT1JERVJCWXN5c19pZFx1MDAyNnN5c3Bhcm1fb2Zmc2V0PTIifX19"
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "5b00ab008737c910793a97983cbb354a"
										}
									],
									"id": "8cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								},
								{
									"values": [
										{
											"string_value": "4d025d2cc37dd6103e3f7fd9d001315e"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjowLCJlbnRpdHlDdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uuc2VydmljZS1ub3cuY29tL2FwaS9ub3cvdjIvdGFibGUvY2hhbmdlX3Rhc2s/c3lzcGFybV9maWVsZHM9c3lzX2lkLGFzc2lnbmVkX3RvXHUwMDI2c3lzcGFybV9leGNsdWRlX3JlZmVyZW5jZV9saW5rPXRydWVcdTAwMjZzeXNwYXJtX2xpbWl0PTFcdTAwMjZzeXNwYXJtX3F1ZXJ5PWFzc2lnbmVkX3RvSU41YjAwYWIwMDg3MzdjOTEwNzkzYTk3OTgzY2JiMzU0YSU1RU9SREVSQllzeXNfaWRcdTAwMjZzeXNwYXJtX29mZnNldD0xIiwicmVsYXRlZEVudGl0eUN1cnNvciI6eyJjdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uuc2VydmljZS1ub3cuY29tL2FwaS9ub3cvdjIvdGFibGUvc3lzX3VzZXJfZ3JtZW1iZXI/c3lzcGFybV9maWVsZHM9c3lzX2lkLHVzZXIuc3lzX2lkXHUwMDI2c3lzcGFybV9leGNsdWRlX3JlZmVyZW5jZV9saW5rPXRydWVcdTAwMjZzeXNwYXJtX2xpbWl0PTFcdTAwMjZzeXNwYXJtX3F1ZXJ5PWdyb3VwSU4wMjcwYzI1MWMzMjAwMjAwYmU2NDdiZmFhMmQzYWVhNiU1RXVzZXIuYWN0aXZlJTNEdHJ1ZSU1RU9SREVSQllzeXNfaWRcdTAwMjZzeXNwYXJtX29mZnNldD0yIn19fQ=="
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "5b00ab008737c910793a97983cbb354a"
										}
									],
									"id": "8cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								},
								{
									"values": [
										{
											"string_value": "e4bc44a0c339d6103e3f7fd9d00131ab"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjowLCJyZWxhdGVkRW50aXR5Q3Vyc29yIjp7ImN1cnNvciI6Imh0dHBzOi8vdGVzdC1pbnN0YW5jZS5zZXJ2aWNlLW5vdy5jb20vYXBpL25vdy92Mi90YWJsZS9zeXNfdXNlcl9ncm1lbWJlcj9zeXNwYXJtX2ZpZWxkcz1zeXNfaWQsdXNlci5zeXNfaWRcdTAwMjZzeXNwYXJtX2V4Y2x1ZGVfcmVmZXJlbmNlX2xpbms9dHJ1ZVx1MDAyNnN5c3Bhcm1fbGltaXQ9MVx1MDAyNnN5c3Bhcm1fcXVlcnk9Z3JvdXBJTjAyNzBjMjUxYzMyMDAyMDBiZTY0N2JmYWEyZDNhZWE2JTVFdXNlci5hY3RpdmUlM0R0cnVlJTVFT1JERVJCWXN5c19pZFx1MDAyNnN5c3Bhcm1fb2Zmc2V0PTMifX19"
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [],
					"next_cursor": "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjowLCJyZWxhdGVkRW50aXR5Q3Vyc29yIjp7ImN1cnNvciI6Imh0dHBzOi8vdGVzdC1pbnN0YW5jZS5zZXJ2aWNlLW5vdy5jb20vYXBpL25vdy92Mi90YWJsZS9zeXNfdXNlcl9ncm1lbWJlcj9zeXNwYXJtX2ZpZWxkcz1zeXNfaWQsdXNlci5zeXNfaWRcdTAwMjZzeXNwYXJtX2V4Y2x1ZGVfcmVmZXJlbmNlX2xpbms9dHJ1ZVx1MDAyNnN5c3Bhcm1fbGltaXQ9MVx1MDAyNnN5c3Bhcm1fcXVlcnk9Z3JvdXBJTjAyNzBjMjUxYzMyMDAyMDBiZTY0N2JmYWEyZDNhZWE2JTVFdXNlci5hY3RpdmUlM0R0cnVlJTVFT1JERVJCWXN5c19pZFx1MDAyNnN5c3Bhcm1fb2Zmc2V0PTQifX19"
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [],
					"next_cursor": "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjowLCJyZWxhdGVkRW50aXR5Q3Vyc29yIjp7ImN1cnNvciI6Imh0dHBzOi8vdGVzdC1pbnN0YW5jZS5zZXJ2aWNlLW5vdy5jb20vYXBpL25vdy92Mi90YWJsZS9zeXNfdXNlcl9ncm1lbWJlcj9zeXNwYXJtX2ZpZWxkcz1zeXNfaWQsdXNlci5zeXNfaWRcdTAwMjZzeXNwYXJtX2V4Y2x1ZGVfcmVmZXJlbmNlX2xpbms9dHJ1ZVx1MDAyNnN5c3Bhcm1fbGltaXQ9MVx1MDAyNnN5c3Bhcm1fcXVlcnk9Z3JvdXBJTjAyNzBjMjUxYzMyMDAyMDBiZTY0N2JmYWEyZDNhZWE2JTVFdXNlci5hY3RpdmUlM0R0cnVlJTVFT1JERVJCWXN5c19pZFx1MDAyNnN5c3Bhcm1fb2Zmc2V0PTUifX19"
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [],
					"next_cursor": "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjoxfX0="
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [],
					"next_cursor": "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjoyfX0="
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [],
					"next_cursor": "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjozfX0="
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [],
					"next_cursor": ""
				}
			}`,
		},
	}

	pageSize10Pages := []struct {
		wantResp string
	}{
		{
			wantResp: `
			{
				"success": {
					"objects": [
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "5b00ab008737c910793a97983cbb354a"
										}
									],
									"id": "8cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								},
								{
									"values": [
										{
											"string_value": "4d025d2cc37dd6103e3f7fd9d001315e"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						},
						{
							"attributes": [
								{
									"values": [
										{
											"string_value": "5b00ab008737c910793a97983cbb354a"
										}
									],
									"id": "8cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								},
								{
									"values": [
										{
											"string_value": "e4bc44a0c339d6103e3f7fd9d00131ab"
										}
									],
									"id": "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6"
								}
							],
							"child_objects": []
						}
					],
					"next_cursor": "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjoxfX0="
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [],
					"next_cursor": "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjoyfX0="
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [],
					"next_cursor": "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjozfX0="
				}
			}`,
		},
		{
			wantResp: `
			{
				"success": {
					"objects": [],
					"next_cursor": ""
				}
			}`,
		},
	}

	pagesMap := map[int64][]struct {
		wantResp string
	}{
		1:  pageSize1Pages,
		10: pageSize10Pages,
	}

	// Test with different page sizes.
	for pageSize, pages := range pagesMap {
		for pageNumber, page := range pages {
			fixtureName := fmt.Sprintf("%s_pagesize_%d_page_%d", fixtureFile, pageSize, pageNumber+1)

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
				if err := protojson.Unmarshal([]byte(pages[pageNumber-1].wantResp), previousResp); err != nil {
					t.Fatal(err)
				}

				cursor = previousResp.GetSuccess().NextCursor
			}

			gotResp, err := adapterClient.GetPage(ctx, &adapter_api_v1.GetPageRequest{
				Datasource: &adapter_api_v1.DatasourceConfig{
					Auth: &adapter_api_v1.DatasourceAuthCredentials{
						AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
							Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
								Username: "test-instance-username",
								Password: "{{OMITTED}}",
							},
						},
					},
					Address: "test-instance.service-now.com",
					Id:      "ServiceNow",
					Type:    "ServiceNow-1.0.1",
					Config:  []byte(advancedFiltersConfig),
				},
				Entity: &adapter_api_v1.EntityConfig{
					Id:         "ChangeTask",
					ExternalId: "change_task",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "7cb589ff-b07c-4190-aa40-ce6aafa5c1d6",
							ExternalId: "sys_id",
							UniqueId:   true,
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "8cb589ff-b07c-4190-aa40-ce6aafa5c1d6",
							ExternalId: "assigned_to",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
				},
				PageSize: pageSize,
				Cursor:   cursor,
			})
			if err != nil {
				t.Fatal(err)
			}

			wantResp := new(adapter_api_v1.GetPageResponse)

			if err := protojson.Unmarshal([]byte(page.wantResp), wantResp); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
				t.Fatal(diff)
			}
		}
	}
}
