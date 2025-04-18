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

func TestJiraDataCenterAdapter_User(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/jira-datacenter/user")
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
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Address: "test-instance.jiradc.ai",
			Id:      "JiraDatacenter",
			Type:    "JiraDatacenter-1.0.0",
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "JiraUser",
			ExternalId: "User",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "key",
					ExternalId: "key",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "emailAddress",
					ExternalId: "emailAddress",
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
								"values": [
									{
										"string_value": "user1@example.com"
									}
								],
								"id": "emailAddress"
							},
							{
								"values": [
									{
										"string_value": "JIRAUSER10000"
									}
								],
								"id": "key"
							}
						],
						"child_objects": []
					},
					{
						"attributes": [
							{
								"values": [
									{
										"string_value": "jiraadmin@sgnl.ai"
									}
								],
								"id": "emailAddress"
							},
							{
								"values": [
									{
										"string_value": "JIRAUSER10105"
									}
								],
								"id": "key"
							}
						],
						"child_objects": []
					}
				],
				"next_cursor": "eyJjb2xsZWN0aW9uSWQiOiJqaXJhLWFkbWluaXN0cmF0b3JzIiwiY29sbGVjdGlvbkN1cnNvciI6MX0="
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

func TestJiraDataCenterAdapter_Group(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/jira-datacenter/group")
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
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Address: "test-instance.jiradc.ai",
			Id:      "JiraDatacenter",
			Type:    "JiraDatacenter-1.0.0",
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "JiraGroup",
			ExternalId: "Group",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
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
								"values": [
									{
										"string_value": "jira-administrators"
									}
								],
								"id": "name"
							}
						],
						"child_objects": []
					},
					{
						"attributes": [
							{
								"values": [
									{
										"string_value": "jira-servicedesk-users"
									}
								],
								"id": "name"
							}
						],
						"child_objects": []
					}
				],
				"next_cursor": ""
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

func TestJiraDataCenterAdapter_GroupMember(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/jira-datacenter/groupMember")
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
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Address: "test-instance.jiradc.ai",
			Id:      "JiraDatacenter",
			Type:    "JiraDatacenter-1.0.0",
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "JiraGroupMember",
			ExternalId: "GroupMember",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "groupId",
					ExternalId: "groupId",
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
								"values": [
									{
										"string_value": "jira-administrators"
									}
								],
								"id": "groupId"
							},
							{
								"values": [
									{
										"string_value": "jira-administrators-JIRAUSER10000"
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
										"string_value": "jira-administrators"
									}
								],
								"id": "groupId"
							},
							{
								"values": [
									{
										"string_value": "jira-administrators-JIRAUSER10105"
									}
								],
								"id": "id"
							}
						],
						"child_objects": []
					}
				],
				"next_cursor": "eyJjb2xsZWN0aW9uSWQiOiJqaXJhLWFkbWluaXN0cmF0b3JzIiwiY29sbGVjdGlvbkN1cnNvciI6MX0="
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

func TestJiraDataCenterAdapter_Issue(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/jira-datacenter/issue")
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
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Address: "test-instance.jiradc.ai",
			Id:      "JiraDatacenter",
			Type:    "JiraDatacenter-1.0.0",
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "JiraIssue",
			ExternalId: "Issue",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "key",
					ExternalId: "key",
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
								"values": [
									{
										"string_value": "12140"
									}
								],
								"id": "id"
							},
							{
								"values": [
									{
										"string_value": "SUP-1"
									}
								],
								"id": "key"
							}
						],
						"child_objects": []
					},
					{
						"attributes": [
							{
								"values": [
									{
										"string_value": "12139"
									}
								],
								"id": "id"
							},
							{
								"values": [
									{
										"string_value": "DTM1-2081"
									}
								],
								"id": "key"
							}
						],
						"child_objects": []
					},
					{
						"attributes": [
							{
								"values": [
									{
										"string_value": "12138"
									}
								],
								"id": "id"
							},
							{
								"values": [
									{
										"string_value": "DTM1-2080"
									}
								],
								"id": "key"
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
