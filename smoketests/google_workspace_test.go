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

func TestGoogleWorkspaceAdapter_User(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/google-workspace/user")
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
			Address: "admin.googleapis.com",
			Id:      "Test",
			Type:    "GoogleWorkspace-1.0.0",
			Config:  []byte(`{"apiVersion":"v1","domain":"sgnldemos.com"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "User",
			ExternalId: "User",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "kind",
					ExternalId: "kind",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "etag",
					ExternalId: "etag",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "primaryEmail",
					ExternalId: "primaryEmail",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "givenName",
					ExternalId: "$.name.givenName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "familyName",
					ExternalId: "$.name.familyName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "fullName",
					ExternalId: "$.name.fullName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "isAdmin",
					ExternalId: "isAdmin",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "isDelegatedAdmin",
					ExternalId: "isDelegatedAdmin",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "lastLoginTime",
					ExternalId: "lastLoginTime",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "nonEditableAliases",
					ExternalId: "nonEditableAliases",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       true,
				},
			},
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					Id:         "Email",
					ExternalId: "emails",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "address",
							ExternalId: "address",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "type",
							ExternalId: "type",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "primary",
							ExternalId: "primary",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
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
							"values": [
								{
									"string_value": "user1"
								}
							],
							"id": "familyName"
						},
						{
							"values": [
								{
									"string_value": "user1 user1"
								}
							],
							"id": "fullName"
						},
						{
							"values": [
								{
									"string_value": "user1"
								}
							],
							"id": "givenName"
						},
						{
							"values": [
								{
									"string_value": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/HFAgBzXOF0pogUtDLHSgZQ9lvls\""
								}
							],
							"id": "etag"
						},
						{
							"values": [
								{
									"string_value": "USER987654321"
								}
							],
							"id": "id"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "isAdmin"
						},
						{
							"values": [
								{
									"bool_value": false
								}
							],
							"id": "isDelegatedAdmin"
						},
						{
							"values": [
								{
									"string_value": "admin#directory#user"
								}
							],
							"id": "kind"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-02-02T23:30:53.000Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "lastLoginTime"
						},
						{
							"values": [
								{
									"string_value": "user1@sgnldemos.com.test-google-a.com"
								}
							],
							"id": "nonEditableAliases"
						},
						{
							"values": [
								{
									"string_value": "user1@sgnldemos.com"
								}
							],
							"id": "primaryEmail"
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
													"string_value": "user1@sgnldemos.com"
												}
											],
											"id": "address"
										},
										{
											"values": [
												{
													"bool_value": true
												}
											],
											"id": "primary"
										}
									],
									"child_objects": []
								},
								{
									"attributes": [
										{
											"values": [
												{
													"string_value": "user1@sgnldemos.com.test-google-a.com"
												}
											],
											"id": "address"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "Email"
						}
					]
				},
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "user2"
								}
							],
							"id": "familyName"
						},
						{
							"values": [
								{
									"string_value": "user2 user2"
								}
							],
							"id": "fullName"
						},
						{
							"values": [
								{
									"string_value": "user2"
								}
							],
							"id": "givenName"
						},
						{
							"values": [
								{
									"string_value": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/VaUMmgVxJHqNreIe-pK6XVo-Kyk\""
								}
							],
							"id": "etag"
						},
						{
							"values": [
								{
									"string_value": "102475661842232156723"
								}
							],
							"id": "id"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "isAdmin"
						},
						{
							"values": [
								{
									"bool_value": false
								}
							],
							"id": "isDelegatedAdmin"
						},
						{
							"values": [
								{
									"string_value": "admin#directory#user"
								}
							],
							"id": "kind"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-04-22T15:33:41.000Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "lastLoginTime"
						},
						{
							"values": [
								{
									"string_value": "user2@sgnldemos.com.test-google-a.com"
								}
							],
							"id": "nonEditableAliases"
						},
						{
							"values": [
								{
									"string_value": "user2@sgnldemos.com"
								}
							],
							"id": "primaryEmail"
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
													"string_value": "user2@sgnldemos.com"
												}
											],
											"id": "address"
										},
										{
											"values": [
												{
													"bool_value": true
												}
											],
											"id": "primary"
										}
									],
									"child_objects": []
								},
								{
									"attributes": [
										{
											"values": [
												{
													"string_value": "user2@sgnldemos.com.test-google-a.com"
												}
											],
											"id": "address"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "Email"
						}
					]
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJRMEZGVXpoUlNVSnJVSHBXUVVWUVpHa3hWSFJOVldVMlZtOUdNM0J1VkdreU5XZHVZV2xWUmxGQlJVNVROVVp1Y0cxdVpWQTBSMUoyVDB0aFJsZ3dVQ3RtWkU1UWIxVkdaelp1UjNOSlYxcGFhemxxTmxRMVIxVkpUMHB3YnpSUFpqRllMMHRJTUhjdk9WZHhVbGxXY3pCNlYzbEVhWGhRWW10dWRuZE9lV28yY2tWYVlXTTVhWEpyUzBNNE1XSmFOMU0yWkhaSE5HdzRXRWxJZVhWb1R6QTVVR2sxYlhaamFVZzFUM0l6U1V3ek4zQnhWbWxZYkdKVlVraHpXbE5ZWXpObGVrZEhXUzgzUjNsYVVUaHNVMUZCZEZJeFRXRmpMM3BzVmxKTFJWSkhXa2xrTDA5SWRYWXJNMWsyWTNjMldtWktXVlZ6VDFRNU9IcFNkWEo0TTJwRU56UnNRbU5DY25KamFEZHJTalV3TDB4TFpFWTRTakpVY1RablN6Rk5jeXRvV0dkbWFVUkVXVGR2UzJSWVIybHJSekpUY0RrelVVeEJNRXB0UXpkUVNHSjFhVXd2TTFZeU1UUkVZVXRoUjA5UVR6VmpjQzh6UTJOelUzSmtiRk5wYjJoblEzaFlaSFJNZUZkek5sSlFWbGhsZWxKcE1EaDVNVFVyU0N0NldWQlpXamRZTDJwTGVFOWpTbU5EUlVnNVlrTkJXVWxWV2tJM1VsZDRiM1ptVGxGeVozWnJjVXhxUzBkYWFrNHlSQzlwUW1WVFUyODNXSE12VTJRd2JEVklkRXBtVjJwTVkxbDJTa05HV0VWalJIRXpiSFY0YTBOcWJuSnVjM1paY2sweFZWRllZa2wwZW5GVlYyNTBlbVJMVFQwPSJ9"
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

func TestGoogleWorkspaceAdapter_Group(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/google-workspace/group")
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
			Address: "admin.googleapis.com",
			Id:      "Test",
			Type:    "GoogleWorkspace-1.0.0",
			Config:  []byte(`{"apiVersion":"v1","domain":"sgnldemos.com"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Group",
			ExternalId: "Group",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "kind",
					ExternalId: "kind",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "etag",
					ExternalId: "etag",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "email",
					ExternalId: "email",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "name",
					ExternalId: "name",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "adminCreated",
					ExternalId: "adminCreated",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "nonEditableAliases",
					ExternalId: "nonEditableAliases",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       true,
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
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "adminCreated"
						},
						{
							"values": [
								{
									"string_value": "emptygroup@sgnldemos.com"
								}
							],
							"id": "email"
						},
						{
							"values": [
								{
									"string_value": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/jegUmo9HAPpclSc2UeC_v9oHm2E\""
								}
							],
							"id": "etag"
						},
						{
							"values": [
								{
									"string_value": "01qoc8b13vgdlqb"
								}
							],
							"id": "id"
						},
						{
							"values": [
								{
									"string_value": "admin#directory#group"
								}
							],
							"id": "kind"
						},
						{
							"values": [
								{
									"string_value": "Empty Group"
								}
							],
							"id": "name"
						},
						{
							"values": [
								{
									"string_value": "emptygroup@sgnldemos.com.test-google-a.com"
								}
							],
							"id": "nonEditableAliases"
						}
					],
					"child_objects": []
				},
				{
					"attributes": [
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "adminCreated"
						},
						{
							"values": [
								{
									"string_value": "group2@sgnldemos.com"
								}
							],
							"id": "email"
						},
						{
							"values": [
								{
									"string_value": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/n5phxV9NHrbJaK4QhiAp3pDR21k\""
								}
							],
							"id": "etag"
						},
						{
							"values": [
								{
									"string_value": "048pi1tg0qf1f8g"
								}
							],
							"id": "id"
						},
						{
							"values": [
								{
									"string_value": "admin#directory#group"
								}
							],
							"id": "kind"
						},
						{
							"values": [
								{
									"string_value": "Group2"
								}
							],
							"id": "name"
						},
						{
							"values": [
								{
									"string_value": "group2@sgnldemos.com.test-google-a.com"
								}
							],
							"id": "nonEditableAliases"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJRMmxWZDB4RFNtNWpiVGt4WTBSS1FXTXlaSFZpUjFKc1lsYzVla3h0VG5aaVUwbHpUbnBWTTA1RVdURk5WRkV4VFZSbk1GTkJUbWRvTldKbk5tWTNYMTlmWDE5QlVUMDkifQ=="
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

func TestGoogleWorkspaceAdapter_Member(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/google-workspace/member")
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
			Address: "admin.googleapis.com",
			Id:      "Test",
			Type:    "GoogleWorkspace-1.0.0",
			Config:  []byte(`{"apiVersion":"v1","domain":"sgnldemos.com"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Member",
			ExternalId: "Member",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "kind",
					ExternalId: "kind",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "memberId",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "groupId",
					ExternalId: "groupId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "etag",
					ExternalId: "etag",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "email",
					ExternalId: "email",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "role",
					ExternalId: "role",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "status",
					ExternalId: "status",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "type",
					ExternalId: "type",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 2,
		Cursor:   "eyJjb2xsZWN0aW9uQ3Vyc29yIjoiUTJsdmQweERTbXhpV0VJd1pWZGtlV0l6Vm5kUlNFNXVZbTE0YTFwWE1YWmplVFZxWWpJd2FVeEVSWGRPVkdNeVRWUm5lazFFV1hoTmVteEpRVEpET0cxMVUybENRVDA5In0=",
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
									"string_value": "user1@sgnldemos.com"
								}
							],
							"id": "email"
						},
						{
							"values": [
								{
									"string_value": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/1pVaKo8kmZjl-_gkcKlCYvYBPjA\""
								}
							],
							"id": "etag"
						},
						{
							"values": [
								{
									"string_value": "048pi1tg0qf1f8g"
								}
							],
							"id": "groupId"
						},
						{
							"values": [
								{
									"string_value": "USER987654321"
								}
							],
							"id": "memberId"
						},
						{
							"values": [
								{
									"string_value": "admin#directory#member"
								}
							],
							"id": "kind"
						},
						{
							"values": [
								{
									"string_value": "MEMBER"
								}
							],
							"id": "role"
						},
						{
							"values": [
								{
									"string_value": "ACTIVE"
								}
							],
							"id": "status"
						},
						{
							"values": [
								{
									"string_value": "USER"
								}
							],
							"id": "type"
						},
						{
							"values": [
								{
									"string_value": "048pi1tg0qf1f8g-USER987654321"
								}
							],
							"id": "uniqueId"
						}
					],
					"child_objects": []
				},
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "user2@sgnldemos.com"
								}
							],
							"id": "email"
						},
						{
							"values": [
								{
									"string_value": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/4j9E51Yhqt_d3t65g3OEza89nps\""
								}
							],
							"id": "etag"
						},
						{
							"values": [
								{
									"string_value": "048pi1tg0qf1f8g"
								}
							],
							"id": "groupId"
						},
						{
							"values": [
								{
									"string_value": "102475661842232156723"
								}
							],
							"id": "memberId"
						},
						{
							"values": [
								{
									"string_value": "admin#directory#member"
								}
							],
							"id": "kind"
						},
						{
							"values": [
								{
									"string_value": "OWNER"
								}
							],
							"id": "role"
						},
						{
							"values": [
								{
									"string_value": "ACTIVE"
								}
							],
							"id": "status"
						},
						{
							"values": [
								{
									"string_value": "USER"
								}
							],
							"id": "type"
						},
						{
							"values": [
								{
									"string_value": "048pi1tg0qf1f8g-102475661842232156723"
								}
							],
							"id": "uniqueId"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJDalJKYURoTFNGRnFTWFpOYVdaNlowVlRSVzB4YUdOdFRrRmpNbVIxWWtkU2JHSlhPWHBNYlU1MllsSm5RbGxLZVVwcFRUaEZJaDhLSFFqSXZNaWZ6Z0VTRW0xaGNtTkFjMmR1YkdSbGJXOXpMbU52YlJnQllKeUppTThFIiwiY29sbGVjdGlvbklkIjoiMDQ4cGkxdGcwcWYxZjhnIiwiY29sbGVjdGlvbkN1cnNvciI6IlEybFZkMHhEU201amJUa3hZMFJLUVdNeVpIVmlSMUpzWWxjNWVreHRUblppVTBselRucFZNMDVFV1RGTlZGRXhUVlJuTUZOQlRtZG9OV0puTm1ZM1gxOWZYMTlCVVQwOSJ9"
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
