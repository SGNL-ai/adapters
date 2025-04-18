// Copyright 2025 SGNL.ai, Inc.
package smoketests

import (
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
