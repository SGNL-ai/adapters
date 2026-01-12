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

func TestGitHubAdapter_Organization(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/organization")
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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Organization",
			ExternalId: "Organization",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "enterpriseId",
					ExternalId: "enterpriseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "announcement",
					ExternalId: "announcement",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "announcementExpiresAt",
					ExternalId: "announcementExpiresAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "announcementUserDismissible",
					ExternalId: "announcementUserDismissible",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "anyPinnableItems",
					ExternalId: "anyPinnableItems",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "archivedAt",
					ExternalId: "archivedAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "avatarUrl",
					ExternalId: "avatarUrl",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "databaseId",
					ExternalId: "databaseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					Id:         "description",
					ExternalId: "description",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "descriptionHTML",
					ExternalId: "descriptionHTML",
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
					Id:         "login",
					ExternalId: "login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "ipAllowListEnabledSetting",
					ExternalId: "ipAllowListEnabledSetting",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "isVerified",
					ExternalId: "isVerified",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "location",
					ExternalId: "location",
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
					Id:         "newTeamResourcePath",
					ExternalId: "newTeamResourcePath",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "newTeamUrl",
					ExternalId: "newTeamUrl",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "organizationBillingEmail",
					ExternalId: "organizationBillingEmail",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "pinnedItemsRemaining",
					ExternalId: "pinnedItemsRemaining",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					Id:         "projectsResourcePath",
					ExternalId: "projectsResourcePath",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "projectsUrl",
					ExternalId: "projectsUrl",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "requiresTwoFactorAuthentication",
					ExternalId: "requiresTwoFactorAuthentication",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "resourcePath",
					ExternalId: "resourcePath",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "teamsResourcePath",
					ExternalId: "teamsResourcePath",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "teamsUrl",
					ExternalId: "teamsUrl",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "twitterUsername",
					ExternalId: "twitterUsername",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "updatedAt",
					ExternalId: "updatedAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "url",
					ExternalId: "url",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "viewerCanAdminister",
					ExternalId: "viewerCanAdminister",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "viewerCanChangePinnedItems",
					ExternalId: "viewerCanChangePinnedItems",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "viewerCanCreateProjects",
					ExternalId: "viewerCanCreateProjects",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "viewerCanCreateRepositories",
					ExternalId: "viewerCanCreateRepositories",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "viewerCanCreateTeams",
					ExternalId: "viewerCanCreateTeams",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "viewerIsAMember",
					ExternalId: "viewerIsAMember",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "websiteUrl",
					ExternalId: "websiteUrl",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
			},
		},
		PageSize: 1,
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
									"bool_value": false
								}
							],
							"id": "anyPinnableItems"
						},
						{
							"values": [
								{
									"string_value": ""
								}
							],
							"id": "avatarUrl"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-08T04:18:55Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"int64_value": "5"
								}
							],
							"id": "databaseId"
						},
						{
							"values": [
								{
									"string_value": "<div></div>"
								}
							],
							"id": "descriptionHTML"
						},
						{
							"values": [
								{
									"string_value": "MDEwOkVudGVycHJpc2Ux"
								}
							],
							"id": "enterpriseId"
						},
						{
							"values": [
								{
									"string_value": "MDEyOk9yZ2FuaXphdGlvbjU="
								}
							],
							"id": "id"
						},
						{
							"values": [
								{
									"string_value": "DISABLED"
								}
							],
							"id": "ipAllowListEnabledSetting"
						},
						{
							"values": [
								{
									"bool_value": false
								}
							],
							"id": "isVerified"
						},
						{
							"values": [
								{
									"string_value": "ArvindOrg1"
								}
							],
							"id": "login"
						},
						{
							"values": [
								{
									"string_value": "ArvindOrg1"
								}
							],
							"id": "name"
						},
						{
							"values": [
								{
									"string_value": "/orgs/ArvindOrg1/new-team"
								}
							],
							"id": "newTeamResourcePath"
						},
						{
							"values": [
								{
									"string_value": "https://test-instance.com/orgs/ArvindOrg1/new-team"
								}
							],
							"id": "newTeamUrl"
						},
						{
							"values": [
								{
									"int64_value": "6"
								}
							],
							"id": "pinnedItemsRemaining"
						},
						{
							"values": [
								{
									"string_value": "/orgs/ArvindOrg1/projects"
								}
							],
							"id": "projectsResourcePath"
						},
						{
							"values": [
								{
									"string_value": "https://test-instance.com/orgs/ArvindOrg1/projects"
								}
							],
							"id": "projectsUrl"
						},
						{
							"values": [
								{
									"bool_value": false
								}
							],
							"id": "requiresTwoFactorAuthentication"
						},
						{
							"values": [
								{
									"string_value": "/ArvindOrg1"
								}
							],
							"id": "resourcePath"
						},
						{
							"values": [
								{
									"string_value": "/orgs/ArvindOrg1/teams"
								}
							],
							"id": "teamsResourcePath"
						},
						{
							"values": [
								{
									"string_value": "https://test-instance.com/orgs/ArvindOrg1/teams"
								}
							],
							"id": "teamsUrl"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-08T04:18:55Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "updatedAt"
						},
						{
							"values": [
								{
									"string_value": "https://test-instance.com/ArvindOrg1"
								}
							],
							"id": "url"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "viewerCanAdminister"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "viewerCanChangePinnedItems"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "viewerCanCreateProjects"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "viewerCanCreateRepositories"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "viewerCanCreateTeams"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "viewerIsAMember"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9"
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

func TestGitHubAdapter_OrganizationUser(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/organizationuser")
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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "OrganizationUser",
			ExternalId: "OrganizationUser",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					ExternalId: "uniqueId",
					Id:         "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					ExternalId: "orgId",
					Id:         "orgId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					ExternalId: "$.node.id",
					Id:         "userId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					ExternalId: "role",
					Id:         "role",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					Id:         "UserOVDE",
					ExternalId: "$.node.organizationVerifiedDomainEmails",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "email",
							ExternalId: "email",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
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
									"string_value": "MDQ6VXNlcjQ="
								}
							],
							"id": "userId"
						},
						{
							"values": [
								{
									"string_value": "MDEyOk9yZ2FuaXphdGlvbjU="
								}
							],
							"id": "orgId"
						},
						{
							"values": [
								{
									"string_value": "ADMIN"
								}
							],
							"id": "role"
						},
						{
							"values": [
								{
									"string_value": "MDEyOk9yZ2FuaXphdGlvbjU=-MDQ6VXNlcjQ="
								}
							],
							"id": "id"
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
													"string_value": "arvind@sgnldemos.com"
												}
											],
											"id": "email"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "UserOVDE"
						}
					]
				},
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "MDQ6VXNlcjk="
								}
							],
							"id": "userId"
						},
						{
							"values": [
								{
									"string_value": "MDEyOk9yZ2FuaXphdGlvbjU="
								}
							],
							"id": "orgId"
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
									"string_value": "MDEyOk9yZ2FuaXphdGlvbjU=-MDQ6VXNlcjk="
								}
							],
							"id": "id"
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
													"string_value": "isabella@sgnldemos.com"
												}
											],
											"id": "email"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "UserOVDE"
						}
					]
				}
			],
			"next_cursor": "eyJjb2xsZWN0aW9uQ3Vyc29yIjoiZXlKb1lYTk9aWGgwVUdGblpTSTZabUZzYzJVc0ltVnVaRU4xY25OdmNpSTZJbGt6Vm5sak1qbDVUMjVaZVU5d1MzRlJXRW95WVZjMWExUXpTbTVOVVZVOUlpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZRPT0ifQ=="
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

func TestGitHubAdapter_Repository(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/repository")
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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Repository",
			ExternalId: "Repository",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					ExternalId: "id",
					Id:         "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					ExternalId: "enterpriseId",
					Id:         "enterpriseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					ExternalId: "orgId",
					Id:         "orgId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					ExternalId: "name",
					Id:         "name",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					ExternalId: "databaseId",
					Id:         "databaseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					ExternalId: "allowUpdateBranch",
					Id:         "allowUpdateBranch",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					ExternalId: "pushedAt",
					Id:         "pushedAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					ExternalId: "createdAt",
					Id:         "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					ExternalId: "$.collaborators.edges",
					Id:         "RepositoryCollaborator",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							ExternalId: "$.node.id",
							Id:         "id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							ExternalId: "permission",
							Id:         "permission",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
				},
			},
		},
		PageSize: 1,
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
									"bool_value": false
								}
							],
							"id": "allowUpdateBranch"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-08T18:51:29Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"int64_value": "1"
								}
							],
							"id": "databaseId"
						},
						{
							"values": [
								{
									"string_value": "MDEwOkVudGVycHJpc2Ux"
								}
							],
							"id": "enterpriseId"
						},
						{
							"values": [
								{
									"string_value": "MDEwOlJlcG9zaXRvcnkx"
								}
							],
							"id": "id"
						},
						{
							"values": [
								{
									"string_value": "repo1"
								}
							],
							"id": "name"
						},
						{
							"values": [
								{
									"string_value": "MDEyOk9yZ2FuaXphdGlvbjU="
								}
							],
							"id": "orgId"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-13T23:07:49Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "pushedAt"
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
													"string_value": "MDQ6VXNlcjQ="
												}
											],
											"id": "id"
										},
										{
											"values": [
												{
													"string_value": "ADMIN"
												}
											],
											"id": "permission"
										}
									],
									"child_objects": []
								},
								{
									"attributes": [
										{
											"values": [
												{
													"string_value": "MDQ6VXNlcjk="
												}
											],
											"id": "id"
										},
										{
											"values": [
												{
													"string_value": "READ"
												}
											],
											"id": "permission"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "RepositoryCollaborator"
						}
					]
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0="
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

func TestGitHubAdapter_User(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/user")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "User",
			ExternalId: "User",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "databaseId",
					ExternalId: "databaseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					Id:         "email",
					ExternalId: "email",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "login",
					ExternalId: "login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "isViewer",
					ExternalId: "isViewer",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "updatedAt",
					ExternalId: "updatedAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
		},
		PageSize: 1,
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
									"datetime_value": {
										"timestamp": "2024-03-08T04:18:47Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"int64_value": "4"
								}
							],
							"id": "databaseId"
						},
						{
							"values": [
								{
									"string_value": ""
								}
							],
							"id": "email"
						},
						{
							"values": [
								{
									"string_value": "MDQ6VXNlcjQ="
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
							"id": "isViewer"
						},
						{
							"values": [
								{
									"string_value": "arooxa"
								}
							],
							"id": "login"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-08T04:18:47Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "updatedAt"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZSU0lzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0="
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

func TestGitHubAdapter_Team(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/team")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Team",
			ExternalId: "Team",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "enterpriseId",
					ExternalId: "enterpriseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "orgId",
					ExternalId: "orgId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "databaseId",
					ExternalId: "databaseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					Id:         "slug",
					ExternalId: "slug",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "viewerCanAdminister",
					ExternalId: "viewerCanAdminister",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "updatedAt",
					ExternalId: "updatedAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					Id:         "TeamMember",
					ExternalId: "$.members.edges",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "id",
							ExternalId: "$.node.id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "databaseId",
							ExternalId: "$.node.databaseId",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
							List:       false,
						},
						{
							Id:         "role",
							ExternalId: "role",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "email",
							ExternalId: "$.node.email",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "login",
							ExternalId: "$.node.login",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "isViewer",
							ExternalId: "$.node.isViewer",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
							List:       false,
						},
						{
							Id:         "updatedAt",
							ExternalId: "$.node.updatedAt",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
							List:       false,
						},
						{
							Id:         "createdAt",
							ExternalId: "$.node.createdAt",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
							List:       false,
						},
					},
				},
				{
					Id:         "TeamRepository",
					ExternalId: "$.repositories.edges",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "id",
							ExternalId: "$.node.id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "permission",
							ExternalId: "permission",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "name",
							ExternalId: "$.node.name",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "databaseId",
							ExternalId: "$.node.databaseId",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
							List:       false,
						},
						{
							Id:         "url",
							ExternalId: "$.node.url",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "allowUpdateBranch",
							ExternalId: "$.node.allowUpdateBranch",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
							List:       false,
						},
						{
							Id:         "pushedAt",
							ExternalId: "$.node.pushedAt",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
							List:       false,
						},
						{
							Id:         "createdAt",
							ExternalId: "$.node.createdAt",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
							List:       false,
						},
					},
				},
			},
		},
		PageSize: 1,
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
									"datetime_value": {
										"timestamp": "2024-03-08T18:48:56Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"int64_value": "1"
								}
							],
							"id": "databaseId"
						},
						{
							"values": [
								{
									"string_value": "MDEwOkVudGVycHJpc2Ux"
								}
							],
							"id": "enterpriseId"
						},
						{
							"values": [
								{
									"string_value": "MDQ6VGVhbTE="
								}
							],
							"id": "id"
						},
						{
							"values": [
								{
									"string_value": "MDEyOk9yZ2FuaXphdGlvbjU="
								}
							],
							"id": "orgId"
						},
						{
							"values": [
								{
									"string_value": "team1"
								}
							],
							"id": "slug"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-08T18:48:56Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "updatedAt"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "viewerCanAdminister"
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
													"datetime_value": {
														"timestamp": "2024-03-08T04:18:47Z",
														"timezone_offset": 0
													}
												}
											],
											"id": "createdAt"
										},
										{
											"values": [
												{
													"int64_value": "4"
												}
											],
											"id": "databaseId"
										},
										{
											"values": [
												{
													"string_value": ""
												}
											],
											"id": "email"
										},
										{
											"values": [
												{
													"string_value": "MDQ6VXNlcjQ="
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
											"id": "isViewer"
										},
										{
											"values": [
												{
													"string_value": "arooxa"
												}
											],
											"id": "login"
										},
										{
											"values": [
												{
													"datetime_value": {
														"timestamp": "2024-03-08T04:18:47Z",
														"timezone_offset": 0
													}
												}
											],
											"id": "updatedAt"
										},
										{
											"values": [
												{
													"string_value": "MAINTAINER"
												}
											],
											"id": "role"
										}
									],
									"child_objects": []
								},
								{
									"attributes": [
										{
											"values": [
												{
													"datetime_value": {
														"timestamp": "2024-03-08T17:52:21Z",
														"timezone_offset": 0
													}
												}
											],
											"id": "createdAt"
										},
										{
											"values": [
												{
													"int64_value": "9"
												}
											],
											"id": "databaseId"
										},
										{
											"values": [
												{
													"string_value": ""
												}
											],
											"id": "email"
										},
										{
											"values": [
												{
													"string_value": "MDQ6VXNlcjk="
												}
											],
											"id": "id"
										},
										{
											"values": [
												{
													"bool_value": false
												}
											],
											"id": "isViewer"
										},
										{
											"values": [
												{
													"string_value": "isabella-sgnl"
												}
											],
											"id": "login"
										},
										{
											"values": [
												{
													"datetime_value": {
														"timestamp": "2024-03-08T19:28:13Z",
														"timezone_offset": 0
													}
												}
											],
											"id": "updatedAt"
										},
										{
											"values": [
												{
													"string_value": "MEMBER"
												}
											],
											"id": "role"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "TeamMember"
						},
						{
							"objects": [
								{
									"attributes": [
										{
											"values": [
												{
													"bool_value": false
												}
											],
											"id": "allowUpdateBranch"
										},
										{
											"values": [
												{
													"datetime_value": {
														"timestamp": "2024-03-08T18:51:43Z",
														"timezone_offset": 0
													}
												}
											],
											"id": "createdAt"
										},
										{
											"values": [
												{
													"int64_value": "2"
												}
											],
											"id": "databaseId"
										},
										{
											"values": [
												{
													"string_value": "MDEwOlJlcG9zaXRvcnky"
												}
											],
											"id": "id"
										},
										{
											"values": [
												{
													"string_value": "repo2"
												}
											],
											"id": "name"
										},
										{
											"values": [
												{
													"datetime_value": {
														"timestamp": "2024-03-16T21:18:14Z",
														"timezone_offset": 0
													}
												}
											],
											"id": "pushedAt"
										},
										{
											"values": [
												{
													"string_value": "https://test-instance.com/ArvindOrg1/repo2"
												}
											],
											"id": "url"
										},
										{
											"values": [
												{
													"string_value": "WRITE"
												}
											],
											"id": "permission"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "TeamRepository"
						}
					]
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3UzNGUldFb3lZVmMxYTFRelNtNU5VVlU5SWl3aWIzSm5ZVzVwZW1GMGFXOXVUMlptYzJWMElqb3dMQ0pKYm01bGNsQmhaMlZKYm1adklqcHVkV3hzZlE9PSJ9"
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

func TestGitHubAdapter_Collaborator(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/collaborator")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Collaborator",
			ExternalId: "Collaborator",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "databaseId",
					ExternalId: "databaseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					Id:         "email",
					ExternalId: "email",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "login",
					ExternalId: "login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "isViewer",
					ExternalId: "isViewer",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "updatedAt",
					ExternalId: "updatedAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
		},
		PageSize: 1,
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
									"datetime_value": {
										"timestamp": "2024-03-08T04:18:47Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"int64_value": "4"
								}
							],
							"id": "databaseId"
						},
						{
							"values": [
								{
									"string_value": ""
								}
							],
							"id": "email"
						},
						{
							"values": [
								{
									"string_value": "MDQ6VXNlcjQ="
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
							"id": "isViewer"
						},
						{
							"values": [
								{
									"string_value": "arooxa"
								}
							],
							"id": "login"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-08T04:18:47Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "updatedAt"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSlpNMVo1WXpJNWVVOXVXWGxQY0VWRklpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZYMTkifQ=="
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

func TestGitHubAdapter_Label(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/label")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Label",
			ExternalId: "Label",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
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
					Id:         "color",
					ExternalId: "color",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "isDefault",
					ExternalId: "isDefault",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
		},
		PageSize: 4,
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
									"string_value": "d73a4a"
								}
							],
							"id": "color"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-08T18:51:30Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"string_value": "MDU6TGFiZWwx"
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
							"id": "isDefault"
						},
						{
							"values": [
								{
									"string_value": "bug"
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
									"string_value": "0075ca"
								}
							],
							"id": "color"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-08T18:51:30Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"string_value": "MDU6TGFiZWwy"
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
							"id": "isDefault"
						},
						{
							"values": [
								{
									"string_value": "documentation"
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
									"string_value": "cfd3d7"
								}
							],
							"id": "color"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-08T18:51:30Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"string_value": "MDU6TGFiZWwz"
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
							"id": "isDefault"
						},
						{
							"values": [
								{
									"string_value": "duplicate"
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
									"string_value": "a2eeef"
								}
							],
							"id": "color"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-08T18:51:30Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"string_value": "MDU6TGFiZWw0"
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
							"id": "isDefault"
						},
						{
							"values": [
								{
									"string_value": "enhancement"
								}
							],
							"id": "name"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSk9RU0lzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5ZlE9PSJ9"
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

func TestGitHubAdapter_IssueLabel(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/issuelabel")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "IssueLabel",
			ExternalId: "IssueLabel",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "issueId",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "labelId",
					ExternalId: "labelId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "title",
					ExternalId: "title",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 4,
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
									"string_value": "MDU6SXNzdWUz"
								}
							],
							"id": "issueId"
						},
						{
							"values": [
								{
									"string_value": "MDU6TGFiZWwx"
								}
							],
							"id": "labelId"
						},
						{
							"values": [
								{
									"string_value": "issue1"
								}
							],
							"id": "title"
						},
						{
							"values": [
								{
									"string_value": "MDU6TGFiZWwx-MDU6SXNzdWUz"
								}
							],
							"id": "uniqueId"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSk5VU0lzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5ZlE9PSJ9"
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

func TestGitHubAdapter_PullRequestLabel(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/pullrequestlabel")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PullRequestLabel",
			ExternalId: "PullRequestLabel",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "pullRequestId",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "labelId",
					ExternalId: "labelId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "title",
					ExternalId: "title",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 4,
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
									"string_value": "MDExOlB1bGxSZXF1ZXN0MQ=="
								}
							],
							"id": "pullRequestId"
						},
						{
							"values": [
								{
									"string_value": "MDU6TGFiZWwx"
								}
							],
							"id": "labelId"
						},
						{
							"values": [
								{
									"string_value": "Create README.md"
								}
							],
							"id": "title"
						},
						{
							"values": [
								{
									"string_value": "MDU6TGFiZWwx-MDExOlB1bGxSZXF1ZXN0MQ=="
								}
							],
							"id": "uniqueId"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSk5VU0lzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5ZlE9PSJ9"
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

func TestGitHubAdapter_Issue(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/issue")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Issue",
			ExternalId: "Issue",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "repositoryId",
					ExternalId: "repositoryId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "authorLogin",
					ExternalId: "$.author.login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "title",
					ExternalId: "title",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "isPinned",
					ExternalId: "isPinned",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
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
									"string_value": "arooxa"
								}
							],
							"id": "authorLogin"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-15T18:40:52Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"string_value": "MDU6SXNzdWUz"
								}
							],
							"id": "id"
						},
						{
							"values": [
								{
									"bool_value": false
								}
							],
							"id": "isPinned"
						},
						{
							"values": [
								{
									"string_value": "MDEwOlJlcG9zaXRvcnkx"
								}
							],
							"id": "repositoryId"
						},
						{
							"values": [
								{
									"string_value": "issue1"
								}
							],
							"id": "title"
						}
					],
					"child_objects": []
				},
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "arooxa"
								}
							],
							"id": "authorLogin"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-15T18:41:04Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"string_value": "MDU6SXNzdWU0"
								}
							],
							"id": "id"
						},
						{
							"values": [
								{
									"bool_value": false
								}
							],
							"id": "isPinned"
						},
						{
							"values": [
								{
									"string_value": "MDEwOlJlcG9zaXRvcnkx"
								}
							],
							"id": "repositoryId"
						},
						{
							"values": [
								{
									"string_value": "issue2"
								}
							],
							"id": "title"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0="
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

func TestGitHubAdapter_IssueAssignee(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/issueassignee")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "IssueAssignee",
			ExternalId: "IssueAssignee",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "userId",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "issueId",
					ExternalId: "issueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "login",
					ExternalId: "login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 4,
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
									"string_value": "MDQ6VXNlcjQ="
								}
							],
							"id": "userId"
						},
						{
							"values": [
								{
									"string_value": "MDU6SXNzdWUz"
								}
							],
							"id": "issueId"
						},
						{
							"values": [
								{
									"string_value": "arooxa"
								}
							],
							"id": "login"
						},
						{
							"values": [
								{
									"string_value": "MDU6SXNzdWUz-MDQ6VXNlcjQ="
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
									"string_value": "MDQ6VXNlcjk="
								}
							],
							"id": "userId"
						},
						{
							"values": [
								{
									"string_value": "MDU6SXNzdWUz"
								}
							],
							"id": "issueId"
						},
						{
							"values": [
								{
									"string_value": "isabella-sgnl"
								}
							],
							"id": "login"
						},
						{
							"values": [
								{
									"string_value": "MDU6SXNzdWUz-MDQ6VXNlcjk="
								}
							],
							"id": "uniqueId"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSlpNMVo1WXpJNWVVOXVXWGxQY0VWRUlpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZYMTkifQ=="
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

func TestGitHubAdapter_IssueParticipant(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/issueparticipant")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "IssueParticipant",
			ExternalId: "IssueParticipant",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "userId",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "issueId",
					ExternalId: "issueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "login",
					ExternalId: "login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 4,
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
									"string_value": "MDQ6VXNlcjQ="
								}
							],
							"id": "userId"
						},
						{
							"values": [
								{
									"string_value": "MDU6SXNzdWUz"
								}
							],
							"id": "issueId"
						},
						{
							"values": [
								{
									"string_value": "arooxa"
								}
							],
							"id": "login"
						},
						{
							"values": [
								{
									"string_value": "MDU6SXNzdWUz-MDQ6VXNlcjQ="
								}
							],
							"id": "uniqueId"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSlpNMVo1WXpJNWVVOXVXWGxQY0VWRUlpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZYMTkifQ=="
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

func TestGitHubAdapter_PullRequest(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/pullrequest")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PullRequest",
			ExternalId: "PullRequest",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "title",
					ExternalId: "title",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "authorLogin",
					ExternalId: "$.author.login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "baseRepositoryId",
					ExternalId: "$.baseRepository.id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "headRepositoryId",
					ExternalId: "$.headRepository.id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "closed",
					ExternalId: "closed",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
		},
		PageSize: 1,
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
									"string_value": "arooxa"
								}
							],
							"id": "authorLogin"
						},
						{
							"values": [
								{
									"string_value": "MDEwOlJlcG9zaXRvcnkx"
								}
							],
							"id": "baseRepositoryId"
						},
						{
							"values": [
								{
									"string_value": "MDEwOlJlcG9zaXRvcnkx"
								}
							],
							"id": "headRepositoryId"
						},
						{
							"values": [
								{
									"bool_value": false
								}
							],
							"id": "closed"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-13T23:07:49Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"string_value": "MDExOlB1bGxSZXF1ZXN0MQ=="
								}
							],
							"id": "id"
						},
						{
							"values": [
								{
									"string_value": "Create README.md"
								}
							],
							"id": "title"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSlpNMVo1WXpJNWVVOXVXWGxQY0VWQ0lpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZYMTkifQ=="
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

func TestGitHubAdapter_PullRequestChangedFile(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/pullrequestchangedfile")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PullRequestChangedFile",
			ExternalId: "PullRequestChangedFile",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "pullRequestId",
					ExternalId: "pullRequestId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "path",
					ExternalId: "path",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "changeType",
					ExternalId: "changeType",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
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
									"string_value": "ADDED"
								}
							],
							"id": "changeType"
						},
						{
							"values": [
								{
									"string_value": "README.md"
								}
							],
							"id": "path"
						},
						{
							"values": [
								{
									"string_value": "MDExOlB1bGxSZXF1ZXN0MQ=="
								}
							],
							"id": "pullRequestId"
						},
						{
							"values": [
								{
									"string_value": "MDExOlB1bGxSZXF1ZXN0MQ==-README.md"
								}
							],
							"id": "uniqueId"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQkZRaUlzSW05eVoyRnVhWHBoZEdsdmJrOW1abk5sZENJNk1Dd2lTVzV1WlhKUVlXZGxTVzVtYnlJNmJuVnNiSDE5In0="
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

func TestGitHubAdapter_PullRequestAssignee(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/pullrequestassignee")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PullRequestAssignee",
			ExternalId: "PullRequestAssignee",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "userId",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "pullRequestId",
					ExternalId: "pullRequestId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "login",
					ExternalId: "login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 4,
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
									"string_value": "MDQ6VXNlcjQ="
								}
							],
							"id": "userId"
						},
						{
							"values": [
								{
									"string_value": "arooxa"
								}
							],
							"id": "login"
						},
						{
							"values": [
								{
									"string_value": "MDExOlB1bGxSZXF1ZXN0MQ=="
								}
							],
							"id": "pullRequestId"
						},
						{
							"values": [
								{
									"string_value": "MDExOlB1bGxSZXF1ZXN0MQ==-MDQ6VXNlcjQ="
								}
							],
							"id": "uniqueId"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSlpNMVo1WXpJNWVVOXVXWGxQY0VWQ0lpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZYMTkifQ=="
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

func TestGitHubAdapter_PullRequestParticipant(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/pullrequestparticipant")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PullRequestParticipant",
			ExternalId: "PullRequestParticipant",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "userId",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "pullRequestId",
					ExternalId: "pullRequestId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "login",
					ExternalId: "login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 4,
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
									"string_value": "MDQ6VXNlcjQ="
								}
							],
							"id": "userId"
						},
						{
							"values": [
								{
									"string_value": "arooxa"
								}
							],
							"id": "login"
						},
						{
							"values": [
								{
									"string_value": "MDExOlB1bGxSZXF1ZXN0MQ=="
								}
							],
							"id": "pullRequestId"
						},
						{
							"values": [
								{
									"string_value": "MDExOlB1bGxSZXF1ZXN0MQ==-MDQ6VXNlcjQ="
								}
							],
							"id": "uniqueId"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSlpNMVo1WXpJNWVVOXVXWGxQY0VWQ0lpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZYMTkifQ=="
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

func TestGitHubAdapter_PullRequestCommit(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/pullrequestcommit")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PullRequestCommit",
			ExternalId: "PullRequestCommit",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "pullRequestId",
					ExternalId: "pullRequestId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "commitId",
					ExternalId: "$.commit.id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "committedDate",
					ExternalId: "$.commit.committedDate",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "authorEmail",
					ExternalId: "$.commit.author.email",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "authorId",
					ExternalId: "$.commit.author.user.id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "authorLogin",
					ExternalId: "$.commit.author.user.login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 4,
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
									"string_value": "arvind@sgnl.ai"
								}
							],
							"id": "authorEmail"
						},
						{
							"values": [
								{
									"string_value": "MDQ6VXNlcjQ="
								}
							],
							"id": "authorId"
						},
						{
							"values": [
								{
									"string_value": "arooxa"
								}
							],
							"id": "authorLogin"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-13T23:07:39Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "committedDate"
						},
						{
							"values": [
								{
									"string_value": "MDY6Q29tbWl0MTo0YWNkMDEzNTJkNTZjYTMzMTA1ZmMyMjU4ZDFmMTI4NzZmMzhlZjRh"
								}
							],
							"id": "commitId"
						},
						{
							"values": [
								{
									"string_value": "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0MTo0YWNkMDEzNTJkNTZjYTMzMTA1ZmMyMjU4ZDFmMTI4NzZmMzhlZjRh"
								}
							],
							"id": "id"
						},
						{
							"values": [
								{
									"string_value": "MDExOlB1bGxSZXF1ZXN0MQ=="
								}
							],
							"id": "pullRequestId"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSlpNMVo1WXpJNWVVOXVXWGxQY0VWQ0lpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZYMTkifQ=="
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

func TestGitHubAdapter_PullRequestReview(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/pullrequestreview")

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
			Address: "test-instance.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"enterpriseSlug": "SGNL",
				"isEnterpriseCloud": false,
				"apiVersion": "v3"
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PullRequestReview",
			ExternalId: "PullRequestReview",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "pullRequestId",
					ExternalId: "pullRequestId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "authorLogin",
					ExternalId: "$.author.login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "state",
					ExternalId: "state",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "authorCanPushToRepository",
					ExternalId: "authorCanPushToRepository",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
		},
		PageSize: 4,
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
									"string_value": "r-rakshith"
								}
							],
							"id": "authorLogin"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "authorCanPushToRepository"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-15T21:06:25Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"string_value": "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3Ng=="
								}
							],
							"id": "id"
						},
						{
							"values": [
								{
									"string_value": "MDExOlB1bGxSZXF1ZXN0Mw=="
								}
							],
							"id": "pullRequestId"
						},
						{
							"values": [
								{
									"string_value": "APPROVED"
								}
							],
							"id": "state"
						}
					],
					"child_objects": []
				},
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "isabella-sgnl"
								}
							],
							"id": "authorLogin"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "authorCanPushToRepository"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-03-15T19:45:09Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"string_value": "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3Mg=="
								}
							],
							"id": "id"
						},
						{
							"values": [
								{
									"string_value": "MDExOlB1bGxSZXF1ZXN0Mw=="
								}
							],
							"id": "pullRequestId"
						},
						{
							"values": [
								{
									"string_value": "APPROVED"
								}
							],
							"id": "state"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJanB1ZFd4c0xDSnZjbWRoYm1sNllYUnBiMjVQWm1aelpYUWlPakFzSWtsdWJtVnlVR0ZuWlVsdVptOGlPbnNpYUdGelRtVjRkRkJoWjJVaU9tWmhiSE5sTENKbGJtUkRkWEp6YjNJaU9pSlpNMVo1WXpJNWVVOXVXWGxQY0VWRUlpd2liM0puWVc1cGVtRjBhVzl1VDJabWMyVjBJam93TENKSmJtNWxjbEJoWjJWSmJtWnZJanB1ZFd4c2ZYMTkifQ=="
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

func TestGitHubAdapter_Repository_Given_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/repository_with_organizations")
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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Repository",
			ExternalId: "Repository",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					ExternalId: "id",
					Id:         "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					ExternalId: "enterpriseId",
					Id:         "enterpriseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					ExternalId: "orgId",
					Id:         "orgId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					ExternalId: "name",
					Id:         "name",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					ExternalId: "databaseId",
					Id:         "databaseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					ExternalId: "allowUpdateBranch",
					Id:         "allowUpdateBranch",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					ExternalId: "pushedAt",
					Id:         "pushedAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					ExternalId: "createdAt",
					Id:         "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					ExternalId: "$.collaborators.edges",
					Id:         "RepositoryCollaborator",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							ExternalId: "$.node.id",
							Id:         "id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							ExternalId: "permission",
							Id:         "permission",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
				},
			},
		},
		PageSize: 1,
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
									"bool_value": false
								}
							],
							"id": "allowUpdateBranch"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-11-13T03:10:55Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "createdAt"
						},
						{
							"values": [
								{
									"int64_value": "887644645"
								}
							],
							"id": "databaseId"
						},
						{
							"values": [
								{
									"string_value": "R_kgDONOhh5Q"
								}
							],
							"id": "id"
						},
						{
							"values": [
								{
									"string_value": "repo-1"
								}
							],
							"id": "name"
						},
						{
							"values": [
								{
									"string_value": "O_kgDOCzkBcw"
								}
							],
							"id": "orgId"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-11-13T03:11:44Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "pushedAt"
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
													"string_value": "U_kgDOBwOdDw"
												}
											],
											"id": "id"
										},
										{
											"values": [
												{
													"string_value": "ADMIN"
												}
											],
											"id": "permission"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "RepositoryCollaborator"
						}
					]
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

func TestGitHubAdapter_Team_Given_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/team_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Team",
			ExternalId: "Team",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "enterpriseId",
					ExternalId: "enterpriseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "orgId",
					ExternalId: "orgId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "databaseId",
					ExternalId: "databaseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					Id:         "slug",
					ExternalId: "slug",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "viewerCanAdminister",
					ExternalId: "viewerCanAdminister",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "updatedAt",
					ExternalId: "updatedAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					Id:         "TeamMember",
					ExternalId: "$.members.edges",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "id",
							ExternalId: "$.node.id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "databaseId",
							ExternalId: "$.node.databaseId",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
							List:       false,
						},
						{
							Id:         "role",
							ExternalId: "role",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "email",
							ExternalId: "$.node.email",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "login",
							ExternalId: "$.node.login",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "isViewer",
							ExternalId: "$.node.isViewer",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
							List:       false,
						},
						{
							Id:         "updatedAt",
							ExternalId: "$.node.updatedAt",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
							List:       false,
						},
						{
							Id:         "createdAt",
							ExternalId: "$.node.createdAt",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
							List:       false,
						},
					},
				},
				{
					Id:         "TeamRepository",
					ExternalId: "$.repositories.edges",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "id",
							ExternalId: "$.node.id",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "permission",
							ExternalId: "permission",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "name",
							ExternalId: "$.node.name",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "databaseId",
							ExternalId: "$.node.databaseId",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
							List:       false,
						},
						{
							Id:         "url",
							ExternalId: "$.node.url",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "allowUpdateBranch",
							ExternalId: "$.node.allowUpdateBranch",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
							List:       false,
						},
						{
							Id:         "pushedAt",
							ExternalId: "$.node.pushedAt",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
							List:       false,
						},
						{
							Id:         "createdAt",
							ExternalId: "$.node.createdAt",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
							List:       false,
						},
					},
				},
			},
		},
		PageSize: 1,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"objects":[
					{
						"attributes":[
							{
								"id":"createdAt",
								"values":[
									{
										"datetimeValue":{
											"timestamp":"2024-11-15T00:00:05Z"
										}
									}
								]
							},
							{
								"id":"databaseId",
								"values":[
									{
										"int64Value":"11554411"
									}
								]
							},
							{
								"id":"id",
								"values":[
									{"stringValue":"T_kwDOCzkBc84AsE5r"}
								]
							},
							{
								"id":"orgId",
								"values":[
									{"stringValue":"O_kgDOCzkBcw"}
								]
							},
							{
								"id":"slug",
								"values":[
									{"stringValue":"team-dh-1"}
								]
							},
							{
								"id":"updatedAt",
								"values":[
									{"datetimeValue":{"timestamp":"2024-11-15T00:00:05Z"}}
								]
							},
							{
								"id":"viewerCanAdminister",
								"values":[
									{"boolValue":true}
								]
							}
						],
						"childObjects":[
							{
								"entityId":"TeamMember",
								 "objects":[
								 	{
								 		"attributes":[
											{
												"id":"createdAt",
												"values":[
													{"datetimeValue":{"timestamp":"2022-11-07T18:20:50Z"}}
												]
											},
											{
												"id":"databaseId",
												"values":[
													{"int64Value":"117677327"}
												]
											},
											{
												"id":"email",
												"values":[
													{"stringValue":""}
												]
											},
											{
												"id":"id",
												"values":[
													{"stringValue":"U_kgDOBwOdDw"}
												]
											},
											{
												"id":"isViewer",
												"values":[
													{"boolValue":true}]
												},
												{
													"id":"login",
													"values":[
														{"stringValue":"dhanya-sgnl"}
													]
												},
												{
													"id":"updatedAt",
													"values":[
														{"datetimeValue":{"timestamp":"2024-11-15T01:27:58Z"}}
													]
												},
												{
													"id":"role",
													"values":[
														{"stringValue":"MAINTAINER"}
													]
												}
											]
										},
										{
											"attributes":[
												{
													"id":"createdAt",
													"values":[
														{"datetimeValue":{"timestamp":"2023-02-21T18:25:37Z"}}
													]
												},
												{
													"id":"databaseId",
													"values":[
														{"int64Value":"126013561"}
													]
												},
												{
													"id":"email",
													"values":[
														{"stringValue":""}
													]
												},
												{
													"id":"id",
													"values":[
														{"stringValue":"U_kgDOB4LQeQ"}
													]
												},
												{
													"id":"isViewer",
													"values":[
														{"boolValue":false}
													]
												},
												{
													"id":"login",
													"values":[
														{"stringValue":"isabella-sgnl"}
													]
												},
												{
													"id":"updatedAt",
													"values":[
														{"datetimeValue":{"timestamp":"2024-10-02T22:57:26Z"}}
													]
												},
												{
													"id":"role",
													"values":[
														{"stringValue":"MEMBER"}
													]
												}
											]
										},
										{
											"attributes":[
												{
													"id":"createdAt",
													"values":[
														{"datetimeValue":{"timestamp":"2024-09-23T17:16:43Z"}}
													]
												},
												{
													"id":"databaseId",
													"values":[
														{"int64Value":"182546176"}
													]
												},
												{
													"id":"email",
													"values":[
														{"stringValue":""}
													]
												},
												{
													"id":"id",
													"values":[
														{"stringValue":"U_kgDOCuFvAA"}
													]
												},
												{
													"id":"isViewer",
													"values":[
														{"boolValue":false}
													]
												},
												{
													"id":"login",
														"values":[
															{"stringValue":"leminhtri2805"}
														]
													},
													{
														"id":"updatedAt",
														"values":[
															{"datetimeValue":{"timestamp":"2024-10-07T16:58:04Z"}}
														]
													},
													{
														"id":"role",
														"values":[
															{"stringValue":"MEMBER"}
														]
													}
												]
											}
										]
									}
								]
							}
						],
				"next_cursor": "eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3VFVOeFdGSnNXVmN3ZEZwSFozUk5ZelJCYzBVMWNpSXNJbTl5WjJGdWFYcGhkR2x2Yms5bVpuTmxkQ0k2TUN3aVNXNXVaWEpRWVdkbFNXNW1ieUk2Ym5Wc2JIMD0ifQ=="
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

func TestGitHubAdapter_PullRequestLabel_Given_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/pullrequestlabel_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PullRequestLabel",
			ExternalId: "PullRequestLabel",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "pullRequestId",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "labelId",
					ExternalId: "labelId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "title",
					ExternalId: "title",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 10,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`{
		"success":{
			"objects":[
				{
					"attributes":[
						{
							"id":"pullRequestId",
							"values":[{"stringValue":"PR_kwDONOhh5c6CDkoO"}]
						},
						{
							"id":"labelId",
							"values":[{"stringValue":"LA_kwDONOhh5c8AAAABzVViKA"}]
						},
						{
							"id":"title",
							"values":[{"stringValue":"another commit"}]
						},
						{
							"id":"uniqueId",
							"values":[{"stringValue":"LA_kwDONOhh5c8AAAABzVViKA-PR_kwDONOhh5c6CDkoO"}]
						}
					]
				}
			],
			"nextCursor":"eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3U0U5T1QyaG9OVkU5UFNJc0ltOXlaMkZ1YVhwaGRHbHZiazltWm5ObGRDSTZNQ3dpU1c1dVpYSlFZV2RsU1c1bWJ5STZiblZzYkgwPSJ9"
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

func TestGitHubAdapter_User_With_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/user_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2",
					"Test-0ne"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "User",
			ExternalId: "User",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "databaseId",
					ExternalId: "databaseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					Id:         "email",
					ExternalId: "email",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "login",
					ExternalId: "login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "isViewer",
					ExternalId: "isViewer",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "updatedAt",
					ExternalId: "updatedAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
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

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"objects":[
					{
						"attributes":[
							{"id":"createdAt", "values":[{"datetimeValue":{"timestamp":"2022-11-07T18:20:50Z"}}]},
							{"id":"databaseId", "values":[{"int64Value":"117677327"}]},
							{"id":"email", "values":[{"stringValue":""}]},
							{"id":"id", "values":[{"stringValue":"U_kgDOBwOdDw"}]},
							{"id":"isViewer", "values":[{"boolValue":true}]},
							{"id":"login", "values":[{"stringValue":"dhanya-sgnl"}]},
							{"id":"updatedAt", "values":[{"datetimeValue":{"timestamp":"2024-11-15T16:57:07Z"}}]}
						]
					},
					{
						"attributes":[
							{"id":"createdAt", "values":[{"datetimeValue":{"timestamp":"2023-02-21T18:25:37Z"}}]},
							{"id":"databaseId", "values":[{"int64Value":"126013561"}]},
							{"id":"email", "values":[{"stringValue":""}]},
							{"id":"id", "values":[{"stringValue":"U_kgDOB4LQeQ"}]},
							{"id":"isViewer", "values":[{"boolValue":false}]},
							{"id":"login", "values":[{"stringValue":"isabella-sgnl"}]},
							{"id":"updatedAt", "values":[{"datetimeValue":{"timestamp":"2024-10-02T22:57:26Z"}}]}
						]
					},
					{
						"attributes":[
							{"id":"createdAt", "values":[{"datetimeValue":{"timestamp":"2024-09-23T17:16:43Z"}}]},
							{"id":"databaseId", "values":[{"int64Value":"182546176"}]},
							{"id":"email", "values":[{"stringValue":""}]},
							{"id":"id", "values":[{"stringValue":"U_kgDOCuFvAA"}]},
							{"id":"isViewer", "values":[{"boolValue":false}]},
							{"id":"login", "values":[{"stringValue":"leminhtri2805"}]},
							{"id":"updatedAt", "values":[{"datetimeValue":{"timestamp":"2024-10-07T16:58:04Z"}}]}
						]
					}
				],
				"nextCursor":"eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpveExDSkpibTVsY2xCaFoyVkpibVp2SWpwdWRXeHNmUT09In0="
			}
		}`,
	), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestGitHubAdapter_Collaborator_With_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/collaborator_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2",
					"Test-0ne"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Collaborator",
			ExternalId: "Collaborator",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "databaseId",
					ExternalId: "databaseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					Id:         "email",
					ExternalId: "email",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "login",
					ExternalId: "login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "isViewer",
					ExternalId: "isViewer",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "updatedAt",
					ExternalId: "updatedAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
		},
		PageSize: 5,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"objects":[
					{
						"attributes":[
							{"id":"createdAt","values":[{"datetimeValue":{"timestamp":"2022-11-07T18:20:50Z"}}]},
							{"id":"databaseId","values":[{"int64Value":"117677327"}]},
							{"id":"email","values":[{"stringValue":""}]},
							{"id":"id","values":[{"stringValue":"U_kgDOBwOdDw"}]},
							{"id":"isViewer","values":[{"boolValue":true}]},
							{"id":"login","values":[{"stringValue":"dhanya-sgnl"}]},
							{"id":"updatedAt","values":[{"datetimeValue":{"timestamp":"2024-11-15T16:57:07Z"}}]}
						]
					},
					{
						"attributes":[
							{"id":"createdAt","values":[{"datetimeValue":{"timestamp":"2023-02-21T18:25:37Z"}}]},
							{"id":"databaseId","values":[{"int64Value":"126013561"}]},
							{"id":"email","values":[{"stringValue":""}]},
							{"id":"id","values":[{"stringValue":"U_kgDOB4LQeQ"}]},
							{"id":"isViewer","values":[{"boolValue":false}]},
							{"id":"login","values":[{"stringValue":"isabella-sgnl"}]},
							{"id":"updatedAt","values":[{"datetimeValue":{"timestamp":"2024-10-02T22:57:26Z"}}]}
						]
					},
					{
						"attributes":[
							{"id":"createdAt","values":[{"datetimeValue":{"timestamp":"2024-09-23T17:16:43Z"}}]},
							{"id":"databaseId","values":[{"int64Value":"182546176"}]},
							{"id":"email","values":[{"stringValue":""}]},
							{"id":"id","values":[{"stringValue":"U_kgDOCuFvAA"}]},
							{"id":"isViewer","values":[{"boolValue":false}]},
							{"id":"login","values":[{"stringValue":"leminhtri2805"}]},
							{"id":"updatedAt","values":[{"datetimeValue":{"timestamp":"2024-10-07T16:58:04Z"}}]}
						]
					}
				],
				"nextCursor":"eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3U0U5T1QyaG9OVkU5UFNJc0ltOXlaMkZ1YVhwaGRHbHZiazltWm5ObGRDSTZNQ3dpU1c1dVpYSlFZV2RsU1c1bWJ5STZiblZzYkgwPSJ9"
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

func TestGitHubAdapter_Label_With_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/label_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2",
					"Test-0ne"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Label",
			ExternalId: "Label",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
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
					Id:         "color",
					ExternalId: "color",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "isDefault",
					ExternalId: "isDefault",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
		},
		PageSize: 4,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"objects":[
					{
						"attributes":[
							{"id":"color","values":[{"stringValue":"d73a4a"}]},
							{"id":"createdAt","values":[{"datetimeValue":{"timestamp":"2024-11-13T03:10:57Z"}}]},
							{"id":"id","values":[{"stringValue":"LA_kwDONOhh5c8AAAABzVViKA"}]},
							{"id":"isDefault","values":[{"boolValue":true}]},
							{"id":"name","values":[{"stringValue":"bug"}]}
						]
					}
				],
				"nextCursor":"eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3U0U5T1QyaG9OVkU5UFNJc0ltOXlaMkZ1YVhwaGRHbHZiazltWm5ObGRDSTZNQ3dpU1c1dVpYSlFZV2RsU1c1bWJ5STZiblZzYkgwPSJ9"
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

func TestGitHubAdapter_Issue_With_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/issue_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"sgnl-demo-cloud"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Issue",
			ExternalId: "Issue",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "repositoryId",
					ExternalId: "repositoryId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "authorLogin",
					ExternalId: "$.author.login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "title",
					ExternalId: "title",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "isPinned",
					ExternalId: "isPinned",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
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

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"objects":[
					{
						"attributes":[
							{"id":"authorLogin","values":[{"stringValue":"leminhtri2805"}]},
							{"id":"createdAt","values":[{"datetimeValue":{"timestamp":"2024-11-01T02:22:35Z"}}]},
							{"id":"id","values":[{"stringValue":"I_kwDOLLPods6cpqjR"}]},
							{"id":"isPinned","values":[{"boolValue":false}]},
							{"id":"repositoryId","values":[{"stringValue":"R_kgDOLLPodg"}]},
							{"id":"title","values":[{"stringValue":"Test1"}]}
						]
					}
				],
				"nextCursor":"eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3U0U5TVRGQnZaR2M5UFNJc0ltOXlaMkZ1YVhwaGRHbHZiazltWm5ObGRDSTZNQ3dpU1c1dVpYSlFZV2RsU1c1bWJ5STZiblZzYkgwPSJ9"
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

func TestGitHubAdapter_IssueLabel_With_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/issuelabel_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "IssueLabel",
			ExternalId: "IssueLabel",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "issueId",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "labelId",
					ExternalId: "labelId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "title",
					ExternalId: "title",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 4,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"objects":[
					{
						"attributes":[
							{"id":"issueId", "values":[{"stringValue":"I_kwDONOhh5c6eu8DK"}]},
							{"id":"labelId", "values":[{"stringValue":"LA_kwDONOhh5c8AAAABzVViKA"}]},
							{"id":"title", "values":[{"stringValue":"test issue"}]},
							{"id":"uniqueId", "values":[{"stringValue":"LA_kwDONOhh5c8AAAABzVViKA-I_kwDONOhh5c6eu8DK"}]}
						]
					}
				],
				"nextCursor":"eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3U0U5T1QyaG9OVkU5UFNJc0ltOXlaMkZ1YVhwaGRHbHZiazltWm5ObGRDSTZNQ3dpU1c1dVpYSlFZV2RsU1c1bWJ5STZiblZzYkgwPSJ9"
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

func TestGitHubAdapter_PullRequest_With_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/pullrequest_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PullRequest",
			ExternalId: "PullRequest",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "title",
					ExternalId: "title",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "authorLogin",
					ExternalId: "$.author.login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "baseRepositoryId",
					ExternalId: "$.baseRepository.id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "headRepositoryId",
					ExternalId: "$.headRepository.id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "closed",
					ExternalId: "closed",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
		},
		PageSize: 1,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"objects":[
					{
						"attributes":[
							{"id":"authorLogin", "values":[{"stringValue":"dhanya-sgnl"}]},
							{"id":"baseRepositoryId", "values":[{"stringValue":"R_kgDONOhh5Q"}]},
							{"id":"headRepositoryId", "values":[{"stringValue":"R_kgDONOhh5Q"}]},
							{"id":"closed", "values":[{"boolValue":true}]},
							{"id":"createdAt", "values":[{"datetimeValue":{"timestamp":"2024-11-15T00:22:36Z"}}]},
							{"id":"id", "values":[{"stringValue":"PR_kwDONOhh5c6B_FTB"}]},
							{"id":"title", "values":[{"stringValue":"Added a file."}]}
						]
					}
				],
				"nextCursor":"eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQklUMmRtZUZWM1VUMDlJaXdpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwdWRXeHNmWDA9In0="
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

func TestGitHubAdapter_PullRequestChangedFile_With_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/pullrequestchangedfile_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PullRequestChangedFile",
			ExternalId: "PullRequestChangedFile",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "pullRequestId",
					ExternalId: "pullRequestId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "path",
					ExternalId: "path",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "changeType",
					ExternalId: "changeType",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
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

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"objects":[
					{
						"attributes":[
							{"id":"changeType","values":[{"stringValue":"ADDED"}]},
							{"id":"path","values":[{"stringValue":"hello.go"}]},
							{"id":"pullRequestId","values":[{"stringValue":"PR_kwDONOhh5c6B_FTB"}]},
							{"id":"uniqueId","values":[{"stringValue":"PR_kwDONOhh5c6B_FTB-hello.go"}]}
						]
					}
				],
				"nextCursor":"eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQklUMmRtZUZWM1VUMDlJaXdpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwdWRXeHNmWDA9In0="
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

func TestGitHubAdapter_PullRequestAssignee_With_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/pullrequestassignee_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PullRequestAssignee",
			ExternalId: "PullRequestAssignee",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "userId",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "pullRequestId",
					ExternalId: "pullRequestId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "login",
					ExternalId: "login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 4,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"objects":[
					{
						"attributes":[
							{"id":"userId","values":[{"stringValue":"U_kgDOBwOdDw"}]},
							{"id":"login","values":[{"stringValue":"dhanya-sgnl"}]},
							{"id":"pullRequestId","values":[{"stringValue":"PR_kwDONOhh5c6B_FTB"}]},
							{"id":"uniqueId","values":[{"stringValue":"PR_kwDONOhh5c6B_FTB-U_kgDOBwOdDw"}]}
						]
					}
				],
				"nextCursor":"eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQklUMmRtZUZWM1VUMDlJaXdpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwdWRXeHNmWDA9In0="
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

func TestGitHubAdapter_PullRequestParticipant_With_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/pullrequestparticipant_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PullRequestParticipant",
			ExternalId: "PullRequestParticipant",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "userId",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "pullRequestId",
					ExternalId: "pullRequestId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "login",
					ExternalId: "login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 4,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"objects":[
					{
						"attributes":[
							{"id":"userId","values":[{"stringValue":"U_kgDOBwOdDw"}]},
							{"id":"login","values":[{"stringValue":"dhanya-sgnl"}]},
							{"id":"pullRequestId","values":[{"stringValue":"PR_kwDONOhh5c6B_FTB"}]},
							{"id":"uniqueId","values":[{"stringValue":"PR_kwDONOhh5c6B_FTB-U_kgDOBwOdDw"}]}
						]
					}
				],
				"nextCursor":"eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQklUMmRtZUZWM1VUMDlJaXdpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwdWRXeHNmWDA9In0="
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

func TestGitHubAdapter_PullRequestCommit_With_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/pullrequestcommit_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PullRequestCommit",
			ExternalId: "PullRequestCommit",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "pullRequestId",
					ExternalId: "pullRequestId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "commitId",
					ExternalId: "$.commit.id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "committedDate",
					ExternalId: "$.commit.committedDate",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "authorEmail",
					ExternalId: "$.commit.author.email",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "authorId",
					ExternalId: "$.commit.author.user.id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "authorLogin",
					ExternalId: "$.commit.author.user.login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 4,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"objects":[
					{
						"attributes":[
							{"id":"authorEmail", "values":[{"stringValue":"117677327+dhanya-sgnl@users.noreply.github.com"}]},
							{"id":"authorId", "values":[{"stringValue":"U_kgDOBwOdDw"}]},
							{"id":"authorLogin", "values":[{"stringValue":"dhanya-sgnl"}]},
							{"id":"committedDate", "values":[{"datetimeValue":{"timestamp":"2024-11-15T00:22:22Z"}}]},
							{"id":"commitId", "values":[{"stringValue":"C_kwDONOhh5doAKGRlOWQ2YzBiZDczZGVhNmEzY2EzYjk5ZWRjOTNhMjViN2M0ZTAzNjk"}]},
							{"id":"id", "values":[{"stringValue":"PURC_lADONOhh5c6B_FTB2gAoZGU5ZDZjMGJkNzNkZWE2YTNjYTNiOTllZGM5M2EyNWI3YzRlMDM2OQ"}]},
							{"id":"pullRequestId", "values":[{"stringValue":"PR_kwDONOhh5c6B_FTB"}]}
						]
					}
				],
				"nextCursor":"eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQklUMmRtZUZWM1VUMDlJaXdpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwdWRXeHNmWDA9In0="
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

func TestGitHubAdapter_PullRequestReview_With_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/pullrequestreview_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "PullRequestReview",
			ExternalId: "PullRequestReview",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "pullRequestId",
					ExternalId: "pullRequestId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "authorLogin",
					ExternalId: "$.author.login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "state",
					ExternalId: "state",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "authorCanPushToRepository",
					ExternalId: "authorCanPushToRepository",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
		},
		PageSize: 4,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"nextCursor":"eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNmJuVnNiQ3dpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwN0ltaGhjMDVsZUhSUVlXZGxJanBtWVd4elpTd2laVzVrUTNWeWMyOXlJam9pV1ROV2VXTXlPWGxQYmxsNVQzQklUMmRtZUZWM1VUMDlJaXdpYjNKbllXNXBlbUYwYVc5dVQyWm1jMlYwSWpvd0xDSkpibTVsY2xCaFoyVkpibVp2SWpwdWRXeHNmWDA9In0="
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

func TestGitHubAdapter_IssueAssignee_With_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/issueassignee_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"my-comp-1"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "IssueAssignee",
			ExternalId: "IssueAssignee",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "userId",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "issueId",
					ExternalId: "issueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "login",
					ExternalId: "login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
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
		"success":{
			"objects":[
				{"attributes":
					[
						{"id":"userId", "values":[{"stringValue":"U_kgDOBwOdDw"}]},
						{"id":"issueId", "values":[{"stringValue":"I_kwDONOhZIc6evVu-"}]},
						{"id":"login", "values":[{"stringValue":"dhanya-sgnl"}]},
						{"id":"uniqueId", "values":[{"stringValue":"I_kwDONOhZIc6evVu--U_kgDOBwOdDw"}]}
					]
				},
				{"attributes":
					[
						{"id":"userId", "values":[{"stringValue":"U_kgDOCuFvAA"}]},
						{"id":"issueId", "values":[{"stringValue":"I_kwDONOhZIc6evVu-"}]},
						{"id":"login", "values":[{"stringValue":"leminhtri2805"}]},
						{"id":"uniqueId", "values":[{"stringValue":"I_kwDONOhZIc6evVu--U_kgDOCuFvAA"}]}
					]
				}
			]
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

func TestGitHubAdapter_IssueParticipant_With_Organizations(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/issueparticipant_with_organizations")

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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"my-comp-1"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "IssueParticipant",
			ExternalId: "IssueParticipant",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "uniqueId",
					ExternalId: "uniqueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "userId",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "issueId",
					ExternalId: "issueId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "login",
					ExternalId: "login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
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

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"objects":[
					{
						"attributes":[
							{"id":"userId","values":[{"stringValue":"U_kgDOBwOdDw"}]},
							{"id":"issueId","values":[{"stringValue":"I_kwDONOhZIc6evVu-"}]},
							{"id":"login","values":[{"stringValue":"dhanya-sgnl"}]},
							{"id":"uniqueId","values":[{"stringValue":"I_kwDONOhZIc6evVu--U_kgDOBwOdDw"}]}
						]
					},
					{
						"attributes":[
							{"id":"userId","values":[{"stringValue":"U_kgDOCuFvAA"}]},
							{"id":"issueId","values":[{"stringValue":"I_kwDONOhZIc6evVu-"}]},
							{"id":"login","values":[{"stringValue":"leminhtri2805"}]},
							{"id":"uniqueId","values":[{"stringValue":"I_kwDONOhZIc6evVu--U_kgDOCuFvAA"}]}
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

func TestGitHubAdapter_Organization_When_Organization_given(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/organization_with_organizations")
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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Organization",
			ExternalId: "Organization",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "enterpriseId",
					ExternalId: "enterpriseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "announcement",
					ExternalId: "announcement",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "announcementExpiresAt",
					ExternalId: "announcementExpiresAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "announcementUserDismissible",
					ExternalId: "announcementUserDismissible",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "anyPinnableItems",
					ExternalId: "anyPinnableItems",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "archivedAt",
					ExternalId: "archivedAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "avatarUrl",
					ExternalId: "avatarUrl",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "createdAt",
					ExternalId: "createdAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "databaseId",
					ExternalId: "databaseId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					Id:         "description",
					ExternalId: "description",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "descriptionHTML",
					ExternalId: "descriptionHTML",
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
					Id:         "login",
					ExternalId: "login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "ipAllowListEnabledSetting",
					ExternalId: "ipAllowListEnabledSetting",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "isVerified",
					ExternalId: "isVerified",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "location",
					ExternalId: "location",
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
					Id:         "newTeamResourcePath",
					ExternalId: "newTeamResourcePath",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "newTeamUrl",
					ExternalId: "newTeamUrl",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "organizationBillingEmail",
					ExternalId: "organizationBillingEmail",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "pinnedItemsRemaining",
					ExternalId: "pinnedItemsRemaining",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					Id:         "projectsResourcePath",
					ExternalId: "projectsResourcePath",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "projectsUrl",
					ExternalId: "projectsUrl",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "requiresTwoFactorAuthentication",
					ExternalId: "requiresTwoFactorAuthentication",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "resourcePath",
					ExternalId: "resourcePath",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "teamsResourcePath",
					ExternalId: "teamsResourcePath",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "teamsUrl",
					ExternalId: "teamsUrl",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "twitterUsername",
					ExternalId: "twitterUsername",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "updatedAt",
					ExternalId: "updatedAt",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "url",
					ExternalId: "url",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "viewerCanAdminister",
					ExternalId: "viewerCanAdminister",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "viewerCanChangePinnedItems",
					ExternalId: "viewerCanChangePinnedItems",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "viewerCanCreateProjects",
					ExternalId: "viewerCanCreateProjects",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "viewerCanCreateRepositories",
					ExternalId: "viewerCanCreateRepositories",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "viewerCanCreateTeams",
					ExternalId: "viewerCanCreateTeams",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "viewerIsAMember",
					ExternalId: "viewerIsAMember",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "websiteUrl",
					ExternalId: "websiteUrl",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
			},
		},
		PageSize: 1,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"objects":[
					{
						"attributes":[
							{"id":"anyPinnableItems","values":[{"boolValue":true}]},
							{"id":"avatarUrl","values":[{"stringValue":""}]},
							{"id":"createdAt","values":[{"datetimeValue":{"timestamp":"2024-11-13T03:10:21Z"}}]},
							{"id":"databaseId","values":[{"int64Value":"188285299"}]},
							{"id":"descriptionHTML","values":[{"stringValue":"<div></div>"}]},
							{"id":"id","values":[{"stringValue":"O_kgDOCzkBcw"}]},
							{"id":"ipAllowListEnabledSetting","values":[{"stringValue":"DISABLED"}]},
							{"id":"isVerified","values":[{"boolValue":false}]},
							{"id":"login","values":[{"stringValue":"dh-test-org-2"}]},
							{"id":"name","values":[{"stringValue":"dh-test-org-2"}]},
							{"id":"newTeamResourcePath","values":[{"stringValue":"/orgs/dh-test-org-2/new-team"}]},
							{"id":"newTeamUrl","values":[{"stringValue":"https://github.com/orgs/dh-test-org-2/new-team"}]},
							{"id":"organizationBillingEmail","values":[{"stringValue":"dhanya@sgnl.ai"}]},
							{"id":"pinnedItemsRemaining","values":[{"int64Value":"6"}]},
							{"id":"projectsResourcePath","values":[{"stringValue":"/orgs/dh-test-org-2/projects"}]},
							{"id":"projectsUrl","values":[{"stringValue":"https://github.com/orgs/dh-test-org-2/projects"}]},
							{"id":"requiresTwoFactorAuthentication","values":[{"boolValue":false}]},
							{"id":"resourcePath","values":[{"stringValue":"/dh-test-org-2"}]},
							{"id":"teamsResourcePath","values":[{"stringValue":"/orgs/dh-test-org-2/teams"}]},
							{"id":"teamsUrl","values":[{"stringValue":"https://github.com/orgs/dh-test-org-2/teams"}]},
							{"id":"updatedAt","values":[{"datetimeValue":{"timestamp":"2024-11-15T00:43:08Z"}}]},
							{"id":"url","values":[{"stringValue":"https://github.com/dh-test-org-2"}]},
							{"id":"viewerCanAdminister","values":[{"boolValue":true}]},
							{"id":"viewerCanChangePinnedItems","values":[{"boolValue":true}]},
							{"id":"viewerCanCreateProjects","values":[{"boolValue":true}]},
							{"id":"viewerCanCreateRepositories","values":[{"boolValue":true}]},
							{"id":"viewerCanCreateTeams","values":[{"boolValue":true}]},
							{"id":"viewerIsAMember","values":[{"boolValue":true}]}
						]
					}
				],
				"nextCursor":""
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

func TestGitHubAdapter_OrganizationUser_When_Organization_given(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/github/organizationuser_with_organizations")
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
			Address: "api.github.com",
			Id:      "Test",
			Type:    "GitHub-1.0.0",
			Config: []byte(`
			{
				"organizations": [
					"dh-test-org-2"
				],
				"isEnterpriseCloud": true,
				"apiVersion": "v3",
				"requestTimeoutSeconds": 300
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "OrganizationUser",
			ExternalId: "OrganizationUser",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					ExternalId: "uniqueId",
					Id:         "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					ExternalId: "orgId",
					Id:         "orgId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					ExternalId: "$.node.id",
					Id:         "userId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					ExternalId: "role",
					Id:         "role",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					Id:         "UserOVDE",
					ExternalId: "$.node.organizationVerifiedDomainEmails",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "email",
							ExternalId: "email",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
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

	err = protojson.Unmarshal([]byte(`
		{
			"success":{
				"objects":[
					{
						"attributes":[
							{"id":"userId", "values":[{"stringValue":"MDQ6VXNlcjM5MTM0NDM0"}]},
							{"id":"orgId", "values":[{"stringValue":"O_kgDOCzkBcw"}]},
							{"id":"role", "values":[{"stringValue":"MEMBER"}]},
							{"id":"id", "values":[{"stringValue":"O_kgDOCzkBcw-MDQ6VXNlcjM5MTM0NDM0"}]}
						]
					},
					{
						"attributes":[
							{"id":"userId", "values":[{"stringValue":"U_kgDOBwOdDw"}]},
							{"id":"orgId", "values":[{"stringValue":"O_kgDOCzkBcw"}]},
							{"id":"role", "values":[{"stringValue":"ADMIN"}]},
							{"id":"id", "values":[{"stringValue":"O_kgDOCzkBcw-U_kgDOBwOdDw"}]}
						]
					}
				],
				"nextCursor":"eyJjdXJzb3IiOiJleUpvWVhOT1pYaDBVR0ZuWlNJNlptRnNjMlVzSW1WdVpFTjFjbk52Y2lJNklsa3pWbmxqTWpsNVQyNVplVTl3U0U5Q2QwOWtSSGM5UFNJc0ltOXlaMkZ1YVhwaGRHbHZiazltWm5ObGRDSTZNQ3dpU1c1dVpYSlFZV2RsU1c1bWJ5STZiblZzYkgwPSJ9"
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
