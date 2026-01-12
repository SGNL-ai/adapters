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

func TestPagerdutyAdapter_User(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/pagerduty/user")
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
					HttpAuthorization: "Token token={{OMITTED}}",
				},
			},
			Address: "api.pagerduty.com",
			Id:      "Pagerduty",
			Type:    "PagerDuty-1.0.0",
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PagerdutyUser",
			ExternalId: "users",
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
			},
			ChildEntities: nil,
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
							 			"string_value": "PT9HH3C"
							 		}
								]
							},
							{
								"id": "name",
								"values": [
								   {
										"string_value": "John Doe"
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
							 			"string_value": "P2U4KPO"
							 		}
								]
							},
							{
								"id": "name",
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
							 	"id": "id",
							 	"values": [
									{
							 			"string_value": "PPAWGRP"
							 		}
								]
							},
							{
								"id": "name",
								"values": [
								   {
										"string_value": "Owen Smith"
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

func TestPagerdutyAdapter_Team(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/pagerduty/team")
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
					HttpAuthorization: "Token token={{OMITTED}}",
				},
			},
			Address: "api.pagerduty.com",
			Id:      "Pagerduty",
			Type:    "PagerDuty-1.0.0",
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PagerdutyTeam",
			ExternalId: "teams",
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
			},
			ChildEntities: nil,
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
							 			"string_value": "PHKRHH4"
							 		}
								]
							},
							{
								"id": "name",
								"values": [
								   {
										"string_value": "Admins"
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
							 			"string_value": "PLWCIWU"
							 		}
								]
							},
							{
								"id": "name",
								"values": [
								   {
										"string_value": "Developers"
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
							 			"string_value": "PEO3H4E"
							 		}
								]
							},
							{
								"id": "name",
								"values": [
								   {
										"string_value": "Marketing"
									}
							   ]
						   }
						]
					}
				],
				"nextCursor": "eyJjdXJzb3IiOjN9"
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

func TestPagerdutyAdapter_Members(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/pagerduty/member")
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
					HttpAuthorization: "Token token={{OMITTED}}",
				},
			},
			Address: "api.pagerduty.com",
			Id:      "Pagerduty",
			Type:    "PagerDuty-1.0.0",
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PagerdutyMember",
			ExternalId: "members",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
			ChildEntities: nil,
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
							 			"string_value": "PHKRHH4-PT9HH3C"
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
							 			"string_value": "PHKRHH4-PPAWGRP"
							 		}
								]
							}
						]
					}
				],
				"nextCursor": "eyJjb2xsZWN0aW9uSWQiOiJQSEtSSEg0IiwiY29sbGVjdGlvbkN1cnNvciI6MX0="
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
