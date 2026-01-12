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

func TestOktaAdapter_User(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/okta/user")
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
					HttpAuthorization: "SSWS {{OMITTED}}",
				},
			},
			Address: "test-instance.okta.com",
			Id:      "Okta",
			Type:    "Okta-1.0.1",
			Config:  []byte(`{"apiVersion":"v1"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "OktaUser",
			ExternalId: "User",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "status",
					ExternalId: "status",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "firstName",
					ExternalId: "$.profile.firstName",
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
								"id": "firstName",
								"values": [
									{
										"string_value": "Micheal"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "00uabig334Y6bdCWr5d7"
									}
								]
							},
							{
								"id": "status",
								"values": [
									{
										"string_value": "STAGED"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "firstName",
								"values": [
									{
										"string_value": "John"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "00ucwpv0gzdQNHACm5d7"
									}
								]
							},
							{
								"id": "status",
								"values": [
									{
										"string_value": "STAGED"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "firstName",
								"values": [
									{
										"string_value": "Pam"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "00uabigku0fl1RxTT5d7"
									}
								]
							},
							{
								"id": "status",
								"values": [
									{
										"string_value": "PROVISIONED"
									}
								]
							}
						]
					}
				],
				"nextCursor": "eyJjdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uub2t0YS5jb20vYXBpL3YxL3VzZXJzP2FmdGVyPTJcdTAwMjZsaW1pdD0zIn0="
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

func TestOktaAdapter_Group(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/okta/group")
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
					HttpAuthorization: "SSWS {{OMITTED}}",
				},
			},
			Address: "test-instance.okta.com",
			Id:      "Okta",
			Type:    "Okta-1.0.1",
			Config:  []byte(`{"apiVersion":"v1"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "OktaGroups",
			ExternalId: "Group",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "type",
					ExternalId: "type",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "name",
					ExternalId: "$.profile.name",
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
								"id": "name",
								"values": [
									{
										"string_value": "Admins"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "00gcwqjuvsA6o7BIQ5d7"
									}
								]
							},
							{
								"id": "type",
								"values": [
									{
										"string_value": "OKTA_GROUP"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "name",
								"values": [
									{
										"string_value": "Developers"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "00gcwh6k5xmXRfTos5d7"
									}
								]
							},
							{
								"id": "type",
								"values": [
									{
										"string_value": "OKTA_GROUP"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "name",
								"values": [
									{
										"string_value": "East Coast"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "00gcwqhha1g99dU6q5d7"
									}
								]
							},
							{
								"id": "type",
								"values": [
									{
										"string_value": "OKTA_GROUP"
									}
								]
							}
						]
					}
				],
				"nextCursor": "eyJjdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uub2t0YS5jb20vYXBpL3YxL2dyb3Vwcz9hZnRlcj0wMGdjd3FoaGExZzk5ZFU2cTVkN1x1MDAyNmxpbWl0PTNcdTAwMjZmaWx0ZXI9dHlwZStlcSslMjJPS1RBX0dST1VQJTIyK29yK3R5cGUrZXErJTIyQVBQX0dST1VQJTIyIn0="
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

func TestOktaAdapter_GroupMember_Empty(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/okta/groupMemberEmpty")
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
					HttpAuthorization: "SSWS {{OMITTED}}",
				},
			},
			Address: "test-instance.okta.com",
			Id:      "Okta",
			Type:    "Okta-1.0.1",
			Config:  []byte(`{"apiVersion":"v1"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "OktaGroupMembers",
			ExternalId: "GroupMember",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "userId",
					ExternalId: "userId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "groupId",
					ExternalId: "groupId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "firstName",
					ExternalId: "$.profile.firstName",
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
				"objects": [],
				"nextCursor": "eyJjb2xsZWN0aW9uSWQiOiIwMGdjd3FqdXZzQTZvN0JJUTVkNyIsImNvbGxlY3Rpb25DdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uub2t0YS5jb20vYXBpL3YxL2dyb3Vwcz9hZnRlcj0wMGdjd3FqdXZzQTZvN0JJUTVkN1x1MDAyNmxpbWl0PTFcdTAwMjZmaWx0ZXI9dHlwZStlcSslMjJPS1RBX0dST1VQJTIyK29yK3R5cGUrZXErJTIyQVBQX0dST1VQJTIyIn0="
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

func TestOktaAdapter_GroupMember(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/okta/groupMember")
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
					HttpAuthorization: "SSWS {{OMITTED}}",
				},
			},
			Address: "test-instance.okta.com",
			Id:      "Okta",
			Type:    "Okta-1.0.1",
			Config:  []byte(`{"apiVersion":"v1"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "OktaGroupMembers",
			ExternalId: "GroupMember",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "userId",
					ExternalId: "userId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "groupId",
					ExternalId: "groupId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "firstName",
					ExternalId: "$.profile.firstName",
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
								"id": "firstName",
								"values": [
									{
										"string_value": "Nick"
									}
								]
							},
							{
								"id": "groupId",
								"values": [
									{
										"string_value": "00gcwh6k5xmXRfTos5d7"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "00u9guv7o2OwqGqSI5d7-00gcwh6k5xmXRfTos5d7"
									}
								]
							},
							{
								"id": "userId",
								"values": [
									{
										"string_value": "00u9guv7o2OwqGqSI5d7"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "firstName",
								"values": [
									{
										"string_value": "Micheal"
									}
								]
							},
							{
								"id": "groupId",
								"values": [
									{
										"string_value": "00gcwh6k5xmXRfTos5d7"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "00uabig334Y6bdCWr5d7-00gcwh6k5xmXRfTos5d7"
									}
								]
							},
							{
								"id": "userId",
								"values": [
									{
										"string_value": "00uabig334Y6bdCWr5d7"
									}
								]
							}
						]
					},
					{
						"attributes": [
							{
								"id": "firstName",
								"values": [
									{
										"string_value": "Pam"
									}
								]
							},
							{
								"id": "groupId",
								"values": [
									{
										"string_value": "00gcwh6k5xmXRfTos5d7"
									}
								]
							},
							{
								"id": "id",
								"values": [
									{
										"string_value": "00uabigku0fl1RxTT5d7-00gcwh6k5xmXRfTos5d7"
									}
								]
							},
							{
								"id": "userId",
								"values": [
									{
										"string_value": "00uabigku0fl1RxTT5d7"
									}
								]
							}
						]
					}
				],
				"nextCursor": "eyJjb2xsZWN0aW9uSWQiOiIwMGdjd2g2azV4bVhSZlRvczVkNyIsImNvbGxlY3Rpb25DdXJzb3IiOiJodHRwczovL3Rlc3QtaW5zdGFuY2Uub2t0YS5jb20vYXBpL3YxL2dyb3Vwcz9hZnRlcj0wMGdjd2g2azV4bVhSZlRvczVkN1x1MDAyNmxpbWl0PTFcdTAwMjZmaWx0ZXI9dHlwZStlcSslMjJPS1RBX0dST1VQJTIyK29yK3R5cGUrZXErJTIyQVBQX0dST1VQJTIyIn0="
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

func TestOktaAdapter_Applications(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/okta/application")
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
					HttpAuthorization: "SSWS {{OMITTED}}",
				},
			},
			Address: "test-instance.okta.com",
			Id:      "Okta",
			Type:    "Okta-1.0.1",
			Config:  []byte(`{"apiVersion":"v1"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "OktaApplication",
			ExternalId: "Application",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "name",
					ExternalId: "name",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "label",
					ExternalId: "label",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "status",
					ExternalId: "status",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "signOnMode",
					ExternalId: "signOnMode",
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
										"string_value": "0oav0szjt4RXG5wFN697"
									}
								]
							},
							{
								"id": "label",
								"values": [
									{
										"string_value": "Okta Admin Console"
									}
								]
							},
							{
								"id": "name",
								"values": [
									{
										"string_value": "saasure"
									}
								]
							},
							{
								"id": "signOnMode",
								"values": [
									{
										"string_value": "OPENID_CONNECT"
									}
								]
							},
							{
								"id": "status",
								"values": [
									{
										"string_value": "ACTIVE"
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
										"string_value": "0oav0t9spdHM3sWaO697"
									}
								]
							},
							{
								"id": "label",
								"values": [
									{
										"string_value": "Okta Dashboard"
									}
								]
							},
							{
								"id": "name",
								"values": [
									{
										"string_value": "okta_enduser"
									}
								]
							},
							{
								"id": "signOnMode",
								"values": [
									{
										"string_value": "OPENID_CONNECT"
									}
								]
							},
							{
								"id": "status",
								"values": [
									{
										"string_value": "ACTIVE"
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
										"string_value": "0oav0t9srlTfo2iV0697"
									}
								]
							},
							{
								"id": "label",
								"values": [
									{
										"string_value": "Okta Browser Plugin"
									}
								]
							},
							{
								"id": "name",
								"values": [
									{
										"string_value": "okta_browser_plugin"
									}
								]
							},
							{
								"id": "signOnMode",
								"values": [
									{
										"string_value": "OPENID_CONNECT"
									}
								]
							},
							{
								"id": "status",
								"values": [
									{
										"string_value": "ACTIVE"
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
