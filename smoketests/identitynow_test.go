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

func TestIdentityNowAdapter_Account(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/identitynow/account")
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
			Address: "test-instance.api.identitynow-demo.com",
			Id:      "IdentityNow",
			Type:    "IdentityNow-1.0.0",
			Config:  []byte(`{"apiVersion": "v3", "entityConfig": {"accounts": {"uniqueIDAttribute": "id"}}}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Account",
			ExternalId: "accounts",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "name",
					ExternalId: "$.identity.name",
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
								"values": [
									{
										"string_value": "cyril.behen"
									}
								],
								"id": "name"
							},
							{
								"values": [
									{
										"string_value": "1a1bb825eb7e4f76b72fbecb27699b31"
									}
								],
								"id": "id"
							}
						],
						"child_objects": []
					},
					{
						"attributes": [
							{
								"values": [
									{
										"string_value": "James.Grogg"
									}
								],
								"id": "name"
							},
							{
								"values": [
									{
										"string_value": "ba699287e60b4014bcc4319f30e9b59e"
									}
								],
								"id": "id"
							}
						],
						"child_objects": []
					},
					{
						"attributes": [
							{
								"values": [
									{
										"string_value": "Beth.Kalb"
									}
								],
								"id": "name"
							},
							{
								"values": [
									{
										"string_value": "08f557fac4f441e59cbe3b15c1989ffe"
									}
								],
								"id": "id"
							}
						],
						"child_objects": []
					}
				],
				"next_cursor": "eyJjdXJzb3IiOjN9"
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

func TestIdentityNowAdapter_Entitlement(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/identitynow/entitlement")
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
			Address: "test-instance.api.identitynow-demo.com",
			Id:      "IdentityNow",
			Type:    "IdentityNow-1.0.0",
			Config:  []byte(`{"apiVersion": "v3", "entityConfig": {"entitlements": {"uniqueIDAttribute": "id", "apiVersion": "beta"}}}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Entitlement",
			ExternalId: "entitlements",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "displayName",
					ExternalId: "$.attributes.displayName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 2,
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
								"values": [
									{
										"string_value": "Basic purchaser [on] Windows Store for Business"
									}
								],
								"id": "displayName"
							},
							{
								"values": [
									{
										"string_value": "ENTITLEMENT_ID_456"
									}
								],
								"id": "id"
							}
						],
						"child_objects": []
					},
					{
						"attributes": [
							{
								"values": [
									{
										"string_value": "default access [on] TrustedPublishersProxyService"
									}
								],
								"id": "displayName"
							},
							{
								"values": [
									{
										"string_value": "00218206fe614e7da637f528accdf15e"
									}
								],
								"id": "id"
							}
						],
						"child_objects": []
					}
				],
				"next_cursor": "eyJjdXJzb3IiOjJ9"
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

func TestIdentityNowAdapter_AccountEntitlement(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/identitynow/account_entitlement")
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
			Address: "test-instance.api.identitynow-demo.com",
			Id:      "IdentityNow",
			Type:    "IdentityNow-1.0.0",
			Config:  []byte(`{"apiVersion": "v3", "entityConfig": {"accountEntitlements": {"uniqueIDAttribute": "id", "apiVersion": "beta"}}}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "AccountEntitlement",
			ExternalId: "accountEntitlements",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "accountId",
					ExternalId: "accountId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "entitlementId",
					ExternalId: "entitlementId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 250,
		// {"collectionId":"08f557fac4f441e59cbe3b15c1989ffe", "collectionCursor":1}.
		Cursor: "eyJjb2xsZWN0aW9uSWQiOiIwOGY1NTdmYWM0ZjQ0MWU1OWNiZTNiMTVjMTk4OWZmZSIsICJjb2xsZWN0aW9uQ3Vyc29yIjoxfQ==",
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
								"values": [
									{
										"string_value": "ba699287e60b4014bcc4319f30e9b59e"
									}
								],
								"id": "accountId"
							},
							{
								"values": [
									{
										"string_value": "8a8816d25d704e07a47bf3946347502d"
									}
								],
								"id": "entitlementId"
							},
							{
								"values": [
									{
										"string_value": "ba699287e60b4014bcc4319f30e9b59e-8a8816d25d704e07a47bf3946347502d"
									}
								],
								"id": "id"
							}
						],
						"child_objects": []
					},
					{
						"attributes": [
							{
								"values": [
									{
										"string_value": "ba699287e60b4014bcc4319f30e9b59e"
									}
								],
								"id": "accountId"
							},
							{
								"values": [
									{
										"string_value": "a02ec36a78ee45498e9401a7ec6c1325"
									}
								],
								"id": "entitlementId"
							},
							{
								"values": [
									{
										"string_value": "ba699287e60b4014bcc4319f30e9b59e-a02ec36a78ee45498e9401a7ec6c1325"
									}
								],
								"id": "id"
							}
						],
						"child_objects": []
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
