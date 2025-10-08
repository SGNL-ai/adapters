// Copyright 2025 SGNL.ai, Inc.

// nolint: lll
package smoketests

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	adapter_api_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/smoketests/common"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestJiraAdapter_User(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/jira/user")
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
						Username: "nick@sgnl.ai",
						Password: "{{OMITTED}}",
					},
				},
			},
			Address: "test-instance.atlassian.net",
			Id:      "Jira",
			Type:    "Jira-1.0.0",
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "JiraUser",
			ExternalId: "User",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "accountId",
					ExternalId: "accountId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "displayName",
					ExternalId: "displayName",
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
							 	"id": "accountId",
							 	"values": [
									{
							 			"string_value": "712020:c95ef8cc-fe03-43cc-8f50-c390bcb9499f"
							 		}
								]
							},
							{
								"id": "displayName",
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
							 	"id": "accountId",
							 	"values": [
									{
							 			"string_value": "557058:f58131cb-b67d-43c7-b30d-6b58d40bd077"
							 		}
								]
							},
							{
								"id": "displayName",
								"values": [
								   {
										"string_value": "Automation for Jira"
									}
							   ]
						   }
						]
					},
					{
						"attributes": [
							{
							 	"id": "accountId",
							 	"values": [
									{
							 			"string_value": "5d53f3cbc6b9320d9ea5bdc2"
							 		}
								]
							},
							{
								"id": "displayName",
								"values": [
								   {
										"string_value": "Jira Outlook"
									}
							   ]
						   }
						]
					}
				],
				"nextCursor": "`+base64.StdEncoding.EncodeToString([]byte(`{"cursor":"3"}`))+`"
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

func TestJiraAdapter_Issue(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/jira/issue")
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
						Username: "nick@sgnl.ai",
						Password: "{{OMITTED}}",
					},
				},
			},
			Address: "test-instance.atlassian.net",
			Id:      "Jira",
			Type:    "Jira-1.0.0",
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
							 	"id": "id",
							 	"values": [
									{
							 			"string_value": "10002"
							 		}
								]
							},
							{
								"id": "key",
								"values": [
								   {
										"string_value": "TES-3"
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
							 			"string_value": "10001"
							 		}
								]
							},
							{
								"id": "key",
								"values": [
								   {
										"string_value": "TES-2"
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
							 			"string_value": "10000"
							 		}
								]
							},
							{
								"id": "key",
								"values": [
								   {
										"string_value": "TES-1"
									}
							   ]
						   }
						]
					}
				],
				"nextCursor": "`+base64.StdEncoding.EncodeToString([]byte(`{"cursor":"3"}`))+`"
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

// TODO: Remove this after fully deprecating the legacy Issue endpoint.TODO
func TestJiraAdapter_EnhancedIssue_LastPage(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/jira/enhancedIssueLastPage")
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
						Username: "nick@sgnl.ai",
						Password: "{{OMITTED}}",
					},
				},
			},
			Address: "test-instance.atlassian.net",
			Id:      "Jira",
			Type:    "Jira-1.0.0",
			Config:  []byte(`{"enhancedIssueSearch":true,"issuesJqlFilter":"project=Test1 OR project=10000"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "JiraIssue",
			ExternalId: "EnhancedIssue",
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
							 	"id": "id",
							 	"values": [
									{
							 			"string_value": "10002"
							 		}
								]
							},
							{
								"id": "key",
								"values": [
								   {
										"string_value": "ED-1"
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

// TODO: Remove this after fully deprecating the legacy Issue endpoint.
func TestJiraAdapter_EnhancedIssue_FirstPage(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/jira/enhancedIssueFirstPage")
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
						Username: "nick@sgnl.ai",
						Password: "{{OMITTED}}",
					},
				},
			},
			Address: "test-instance.atlassian.net",
			Id:      "Jira",
			Type:    "Jira-1.0.0",
			Config:  []byte(`{"enhancedIssueSearch":true,"issuesJqlFilter":"project=Test1 OR project=10000"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "JiraIssue",
			ExternalId: "EnhancedIssue",
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
		PageSize: 1,
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
							 			"string_value": "10002"
							 		}
								]
							},
							{
								"id": "key",
								"values": [
								   {
										"string_value": "ED-1"
									}
							   ]
						   }
						]
					}
				],
				"nextCursor": "`+base64.StdEncoding.EncodeToString([]byte(`{"cursor":"example-cursor"}`))+`"
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

func TestJiraAdapter_Group(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/jira/group")
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
						Username: "nick@sgnl.ai",
						Password: "{{OMITTED}}",
					},
				},
			},
			Address: "test-instance.atlassian.net",
			Id:      "Jira",
			Type:    "Jira-1.0.0",
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "JiraGroup",
			ExternalId: "Group",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "groupId",
					ExternalId: "groupId",
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
							 	"id": "groupId",
							 	"values": [
									{
							 			"string_value": "0f1c156b-a88b-4918-ba29-35eec8fed41c"
							 		}
								]
							},
							{
								"id": "name",
								"values": [
								   {
										"string_value": "jira-software-user-access-admins-713b7785-bfe1-44f7-aaed-80dc5ec72616"
									}
							   ]
						   }
						]
					},
					{
						"attributes": [
							{
							 	"id": "groupId",
							 	"values": [
									{
							 			"string_value": "114e2bfb-5281-4bae-8b5c-100999f0580d"
							 		}
								]
							},
							{
								"id": "name",
								"values": [
								   {
										"string_value": "atlassian-addons-admin"
									}
							   ]
						   }
						]
					},
					{
						"attributes": [
							{
							 	"id": "groupId",
							 	"values": [
									{
							 			"string_value": "3bbab507-d2c6-481c-918e-36722af115d3"
							 		}
								]
							},
							{
								"id": "name",
								"values": [
								   {
										"string_value": "jira-admins-713b7785-bfe1-44f7-aaed-80dc5ec72616"
									}
							   ]
						   }
						]
					}
				],
				"nextCursor": "`+base64.StdEncoding.EncodeToString([]byte(`{"cursor":"3"}`))+`"
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

func TestJiraAdapter_GroupMember(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/jira/groupMember")
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
						Username: "nick@sgnl.ai",
						Password: "{{OMITTED}}",
					},
				},
			},
			Address: "test-instance.atlassian.net",
			Id:      "Jira",
			Type:    "Jira-1.0.0",
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
							 	"id": "groupId",
							 	"values": [
									{
							 			"string_value": "0f1c156b-a88b-4918-ba29-35eec8fed41c"
							 		}
								]
							},
							{
								"id": "id",
								"values": [
								   {
										"string_value": "0f1c156b-a88b-4918-ba29-35eec8fed41c-712020:c95ef8cc-fe03-43cc-8f50-c390bcb9499f"
									}
							   ]
						   }
						]
					}
				],
				"nextCursor": "`+base64.StdEncoding.EncodeToString([]byte(`{"collectionId":"0f1c156b-a88b-4918-ba29-35eec8fed41c","collectionCursor":"1"}`))+`"
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

func TestJiraAdapter_Workspace(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/jira/workspace")
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
						Username: "shashank@sgnl.ai",
						Password: "{{OMITTED}}",
					},
				},
			},
			Address: "test-instance.atlassian.net",
			Id:      "Jira",
			Type:    "Jira-1.0.0",
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "JiraWorkspace",
			ExternalId: "Workspace",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "workspaceId",
					ExternalId: "workspaceId",
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
							 	"id": "workspaceId",
							 	"values": [
									{
							 			"string_value": "10c19baf-1ce3-4558-ad80-47bc4494ded7"
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

func TestJiraAdapter_Object(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/jira/object")
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
						Username: "shashank@sgnl.ai",
						Password: "{{OMITTED}}",
					},
				},
			},
			Config:  []byte(`{"objectsQlQuery": "objectType = Customer"}`),
			Address: "test-instance.atlassian.net",
			Id:      "Jira",
			Type:    "Jira-1.0.0",
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "JiraObject",
			ExternalId: "Object",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "globalId",
					ExternalId: "globalId",
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
							 	"id": "globalId",
							 	"values": [
									{
							 			"string_value": "10c19baf-1ce3-4558-ad80-47bc4494ded7:4"
							 		}
								]
							}
						]
					},
					{
						"attributes": [
							{
							 	"id": "globalId",
							 	"values": [
									{
							 			"string_value": "10c19baf-1ce3-4558-ad80-47bc4494ded7:3"
							 		}
								]
							}
						]
					},
					{
						"attributes": [
							{
							 	"id": "globalId",
							 	"values": [
									{
							 			"string_value": "10c19baf-1ce3-4558-ad80-47bc4494ded7:6"
							 		}
								]
							}
						]
					}
				],
				"nextCursor": "`+base64.StdEncoding.EncodeToString([]byte(`{"cursor":"3","collectionId":"10c19baf-1ce3-4558-ad80-47bc4494ded7"}`))+`"
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
