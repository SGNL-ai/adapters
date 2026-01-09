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

func TestAWSAdapter_User(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/user")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"region":"us-west-2"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "user",
			ExternalId: "User",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Arn",
					ExternalId: "Arn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "AccountId",
					ExternalId: "AccountId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "Path",
					ExternalId: "Path",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "UserId",
					ExternalId: "UserId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "UserName",
					ExternalId: "UserName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 2,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
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
									"string_value": "000000000000"
								}
							],
							"id": "AccountId"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::000000000000:user/sampleuser1"
								}
							],
							"id": "Arn"
						},
						{
							"values": [
								{
									"string_value": "/"
								}
							],
							"id": "Path"
						},
						{
							"values": [
								{
									"string_value": "AIDAXXXXXXXXXXXXXXXX1"
								}
							],
							"id": "UserId"
						},
						{
							"values": [
								{
									"string_value": "sampleuser1"
								}
							],
							"id": "UserName"
						}
					],
					"child_objects": []
				},
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "000000000000"
								}
							],
							"id": "AccountId"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::000000000000:user/sgnl-test"
								}
							],
							"id": "Arn"
						},
						{
							"values": [
								{
									"string_value": "/"
								}
							],
							"id": "Path"
						},
						{
							"values": [
								{
									"string_value": "AIDAXXXXXXXXXXXXXXXX2"
								}
							],
							"id": "UserId"
						},
						{
							"values": [
								{
									"string_value": "sgnl-test"
								}
							],
							"id": "UserName"
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

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapter_Group(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/group")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"region":"us-west-2"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "group",
			ExternalId: "Group",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Arn",
					ExternalId: "Arn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "AccountId",
					ExternalId: "AccountId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "Path",
					ExternalId: "Path",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "GroupId",
					ExternalId: "GroupId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "GroupName",
					ExternalId: "GroupName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 2,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
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
									"string_value": "000000000000"
								}
							],
							"id": "AccountId"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::000000000000:group/Group1"
								}
							],
							"id": "Arn"
						},
						{
							"values": [
								{
									"string_value": "AGPAXXXXXXXXXXXXXXXX1"
								}
							],
							"id": "GroupId"
						},
						{
							"values": [
								{
									"string_value": "Group1"
								}
							],
							"id": "GroupName"
						},
						{
							"values": [
								{
									"string_value": "/"
								}
							],
							"id": "Path"
						}
					],
					"child_objects": []
				},
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "000000000000"
								}
							],
							"id": "AccountId"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::000000000000:group/sgnl-group-test"
								}
							],
							"id": "Arn"
						},
						{
							"values": [
								{
									"string_value": "AGPAXXXXXXXXXXXXXXXX2"
								}
							],
							"id": "GroupId"
						},
						{
							"values": [
								{
									"string_value": "sgnl-group-test"
								}
							],
							"id": "GroupName"
						},
						{
							"values": [
								{
									"string_value": "/"
								}
							],
							"id": "Path"
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

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapter_Role(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/role")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"region":"us-west-2"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "role",
			ExternalId: "Role",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Arn",
					ExternalId: "Arn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "AccountId",
					ExternalId: "AccountId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "AssumeRolePolicyDocument",
					ExternalId: "AssumeRolePolicyDocument",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "Description",
					ExternalId: "Description",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "MaxSessionDuration",
					ExternalId: "MaxSessionDuration",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
				},
				{
					Id:         "Path",
					ExternalId: "Path",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "PermissionsBoundary",
					ExternalId: "PermissionsBoundary",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "RoleId",
					ExternalId: "RoleId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "RoleName",
					ExternalId: "RoleName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 2,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
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
									"string_value": "000000000000"
								}
							],
							"id": "AccountId"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::000000000000:role/aws-reserved/sso.amazonaws.com/AWSReservedSSO_AdministratorAccess_1234567890abcdef"
								}
							],
							"id": "Arn"
						},
						{
							"values": [
								{
									"string_value": "%7B%22Version%22%3A%222012-10-17%22%2C%22Statement%22%3A%5B%7B%22Effect%22%3A%22Allow%22%2C%22Principal%22%3A%7B%22Federated%22%3A%22arn%3Aaws%3Aiam%3A%3A000000000000%3Asaml-provider%2FAWSSSO_123456abcdef_DO_NOT_DELETE%22%7D%2C%22Action%22%3A%5B%22sts%3AAssumeRoleWithSAML%22%2C%22sts%3ATagSession%22%5D%2C%22Condition%22%3A%7B%22StringEquals%22%3A%7B%22SAML%3Aaud%22%3A%22https%3A%2F%2Fsignin.aws.amazon.com%2Fsaml%22%7D%7D%7D%5D%7D"
								}
							],
							"id": "AssumeRolePolicyDocument"
						},
						{
							"values": [
								{
									"int64_value": "43200"
								}
							],
							"id": "MaxSessionDuration"
						},
						{
							"values": [
								{
									"string_value": "/aws-reserved/sso.amazonaws.com/"
								}
							],
							"id": "Path"
						},
						{
							"values": [
								{
									"string_value": "AROAXXXXXXXXXXXXXXXX1"
								}
							],
							"id": "RoleId"
						},
						{
							"values": [
								{
									"string_value": "AWSReservedSSO_AdministratorAccess_1234567890abcdef"
								}
							],
							"id": "RoleName"
						}
					],
					"child_objects": []
				},
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "000000000000"
								}
							],
							"id": "AccountId"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::000000000000:role/aws-reserved/sso.amazonaws.com/AWSReservedSSO_Okta_Transform_Permission_Set_abcdef1234567890"
								}
							],
							"id": "Arn"
						},
						{
							"values": [
								{
									"string_value": "%7B%22Version%22%3A%222012-10-17%22%2C%22Statement%22%3A%5B%7B%22Effect%22%3A%22Allow%22%2C%22Principal%22%3A%7B%22Federated%22%3A%22arn%3Aaws%3Aiam%3A%3A000000000000%3Asaml-provider%2FAWSSSO_123456abcdef_DO_NOT_DELETE%22%7D%2C%22Action%22%3A%5B%22sts%3AAssumeRoleWithSAML%22%2C%22sts%3ATagSession%22%5D%2C%22Condition%22%3A%7B%22StringEquals%22%3A%7B%22SAML%3Aaud%22%3A%22https%3A%2F%2Fsignin.aws.amazon.com%2Fsaml%22%7D%7D%7D%5D%7D"
								}
							],
							"id": "AssumeRolePolicyDocument"
						},
						{
							"values": [
								{
									"int64_value": "43200"
								}
							],
							"id": "MaxSessionDuration"
						},
						{
							"values": [
								{
									"string_value": "/aws-reserved/sso.amazonaws.com/"
								}
							],
							"id": "Path"
						},
						{
							"values": [
								{
									"string_value": "AROAXXXXXXXXXXXXXXXX2"
								}
							],
							"id": "RoleId"
						},
						{
							"values": [
								{
									"string_value": "AWSReservedSSO_Okta_Transform_Permission_Set_abcdef1234567890"
								}
							],
							"id": "RoleName"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJBR1RQRmdGRGlBUTlRbjZ1TzNVbnR0ZXMvZWNZVU5GbDVQNWcxUS8zQlN2b0g5bTQzVk5RWktJQjRuMFZhaSt2azBEazc2aGJQSlhReFJSdCtmaUcyVHhVcU5iZnUvL1VMMWhyVE9tblJhbG5nRGdabENMMmwrNENtTTNxdThodTVhbzJEOUo0TURUTlJtamdKdHNtNlhVdUhQbmJiYXFFVENrVkJ6cnpXdTM1N0hvR3FJNzRPVnRNMHc2ZGdtcWQ3WW9pWDhmcTJpUmgifQ=="
		}
	}
	`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapter_Policy(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/policy")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"region":"us-west-2"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "policy",
			ExternalId: "Policy",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Arn",
					ExternalId: "Arn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "AccountId",
					ExternalId: "AccountId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "DefaultVersionId",
					ExternalId: "DefaultVersionId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "Description",
					ExternalId: "Description",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "AttachmentCount",
					ExternalId: "AttachmentCount",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
				},
				{
					Id:         "Path",
					ExternalId: "Path",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "IsAttachable",
					ExternalId: "IsAttachable",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "PolicyId",
					ExternalId: "PolicyId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "PermissionsBoundaryUsageCount",
					ExternalId: "PermissionsBoundaryUsageCount",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					Id:         "PolicyName",
					ExternalId: "PolicyName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 2,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
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
										"string_value": "000000000000"
									}
								],
								"id": "AccountId"
							},
							{
								"values": [
									{
										"string_value": "arn:aws:iam::000000000000:policy/ExamplePolicy1"
									}
								],
								"id": "Arn"
							},
							{
								"values": [
									{
										"int64_value": "1"
									}
								],
								"id": "AttachmentCount"
							},
							{
								"values": [
									{
										"string_value": "v6"
									}
								],
								"id": "DefaultVersionId"
							},
							{
								"values": [
									{
										"bool_value": true
									}
								],
								"id": "IsAttachable"
							},
							{
								"values": [
									{
										"string_value": "/"
									}
								],
								"id": "Path"
							},
							{
								"values": [
									{
										"int64_value": "0"
									}
								],
								"id": "PermissionsBoundaryUsageCount"
							},
							{
								"values": [
									{
										"string_value": "ANPAXXXXXXXXXXXXXXXX1"
									}
								],
								"id": "PolicyId"
							},
							{
								"values": [
									{
										"string_value": "ExamplePolicy1"
									}
								],
								"id": "PolicyName"
							}
						],
						"child_objects": []
					},
					{
						"attributes": [
							{
								"values": [
									{
										"string_value": "aws"
									}
								],
								"id": "AccountId"
							},
							{
								"values": [
									{
										"string_value": "arn:aws:iam::aws:policy/AdministratorAccess"
									}
								],
								"id": "Arn"
							},
							{
								"values": [
									{
										"int64_value": "1"
									}
								],
								"id": "AttachmentCount"
							},
							{
								"values": [
									{
										"string_value": "v1"
									}
								],
								"id": "DefaultVersionId"
							},
							{
								"values": [
									{
										"string_value": "Provides full access to AWS services and resources."
									}
								],
								"id": "Description"
							},
							{
								"values": [
									{
										"bool_value": true
									}
								],
								"id": "IsAttachable"
							},
							{
								"values": [
									{
										"string_value": "/"
									}
								],
								"id": "Path"
							},
							{
								"values": [
									{
										"int64_value": "0"
									}
								],
								"id": "PermissionsBoundaryUsageCount"
							},
							{
								"values": [
									{
										"string_value": "ANPAXXXXXXXXXXXXXXXX2"
									}
								],
								"id": "PolicyId"
							},
							{
								"values": [
									{
										"string_value": "AdministratorAccess"
									}
								],
								"id": "PolicyName"
							}
						],
						"child_objects": []
					}
				],
				"next_cursor": "eyJjdXJzb3IiOiJBRGJPNFBGUGRFZDhtRENoVzRMcE1ndWlpMFFUelFjM2JBblNwdGJ4OWw1cWM1a1hvU0QzN2s1QjMzelRvTU9hNjdJNFNzbC9tUUpHV1BhbFFDenZmOGk5dXZMRTZCREorZk9EZzkyVFlmNExHMUVRQ2U1QWlSNDlrdHk3Vm9vVXV5RT0ifQ=="
			}
		}
		`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapter_IdentityProvider(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/identityprovider")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"region":"us-west-2"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "identityProvider",
			ExternalId: "IdentityProvider",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Arn",
					ExternalId: "Arn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "AccountId",
					ExternalId: "AccountId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "CreateDate",
					ExternalId: "CreateDate",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "ValidUntil",
					ExternalId: "ValidUntil",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 2,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
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
										"string_value": "000000000000"
									}
								],
								"id": "AccountId"
							},
							{
								"values": [
									{
										"string_value": "arn:aws:iam::000000000000:saml-provider/AWSSSO_123456abcdef_DO_NOT_DELETE"
									}
								],
								"id": "Arn"
							},
							{
								"values": [
									{
										"string_value": "2023-12-13T19:51:02Z"
									}
								],
								"id": "CreateDate"
							},
							{
								"values": [
									{
										"string_value": "2123-12-13T19:51:02Z"
									}
								],
								"id": "ValidUntil"
							}
						],
						"child_objects": []
					},
					{
						"attributes": [
							{
								"values": [
									{
										"string_value": "000000000000"
									}
								],
								"id": "AccountId"
							},
							{
								"values": [
									{
										"string_value": "arn:aws:iam::000000000000:saml-provider/DemoProvider"
									}
								],
								"id": "Arn"
							},
							{
								"values": [
									{
										"string_value": "2024-05-10T17:28:41Z"
									}
								],
								"id": "CreateDate"
							},
							{
								"values": [
									{
										"string_value": "2124-05-10T17:28:41Z"
									}
								],
								"id": "ValidUntil"
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

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapter_GroupPolicy(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/grouppolicy")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"region":"us-west-2"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "groupPolicy",
			ExternalId: "GroupPolicy",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "PolicyArn",
					ExternalId: "PolicyArn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "GroupName",
					ExternalId: "GroupName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 2,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
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
									"string_value": "Group1"
								}
							],
							"id": "GroupName"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::000000000000:policy/ExamplePolicy1"
								}
							],
							"id": "PolicyArn"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::000000000000:policy/ExamplePolicy1-Group1"
								}
							],
							"id": "id"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDEiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiQUZTM2g2OXUzaTlFV0hjcThOUmhTdmVHMGtEYzZFNjZnYnZBKzJNRm03WWQwYXFHeDhSRjJ6VkU2dkZYYnV3Skp4MTRVdjJkYzhiMnFyWW4xQ3I1NTF2NW5jT2NrK1Z2dmJyVXY2ST0ifQ=="
		}
	}
	`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapter_RolePolicy(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/rolepolicy")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"region":"us-west-2"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "rolePolicy",
			ExternalId: "RolePolicy",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "PolicyArn",
					ExternalId: "PolicyArn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "RoleName",
					ExternalId: "RoleName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 1,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
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
									"string_value": "arn:aws:iam::aws:policy/AdministratorAccess"
								}
							],
							"id": "PolicyArn"
						},
						{
							"values": [
								{
									"string_value": "AWSReservedSSO_AdministratorAccess_1234567890abcdef"
								}
							],
							"id": "RoleName"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::aws:policy/AdministratorAccess-AWSReservedSSO_AdministratorAccess_1234567890abcdef"
								}
							],
							"id": "id"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjb2xsZWN0aW9uSWQiOiJBV1NSZXNlcnZlZFNTT19BZG1pbmlzdHJhdG9yQWNjZXNzXzM2MjMyZGM4MDZlYjZjMjEiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiQUZHWE96S2c3N1pRK0RlMTgwUlpWY1h0bTZ0dWlBSUVRNjRSZzlmTGJvUmFVczZDS29SckJmZHZDSGNNNDRYalFiVDJqaGZVRkxjYTQ1d1FiTmM3S3Y3WGlxMnhiTUlpYmVTZHBsUXViQkMrTTVKVWNsWXBtQ1pvQ2Yxb2ZGMHFITXNmRnNSV0VYMnBUTnBkOHY0WWpCYnFlb3p0eWp3SldoZ0xySU5zN1libElwWG9YQzcwN3JDV2ZyNnRGcUE9In0="
		}
	}
	`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapter_UserPolicy(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/userpolicy")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"region":"us-west-2"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "userPolicy",
			ExternalId: "UserPolicy",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "PolicyArn",
					ExternalId: "PolicyArn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "UserName",
					ExternalId: "UserName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 1,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
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
									"string_value": "arn:aws:iam::aws:policy/ReadOnlyAccess"
								}
							],
							"id": "PolicyArn"
						},
						{
							"values": [
								{
									"string_value": "sampleuser1"
								}
							],
							"id": "UserName"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::aws:policy/ReadOnlyAccess-sampleuser1"
								}
							],
							"id": "id"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjb2xsZWN0aW9uSWQiOiJzYW1wbGV1c2VyMSIsImNvbGxlY3Rpb25DdXJzb3IiOiJBR1RBdGp4MTJ4Vm9QM21mV3dvMDNpWDd5Rm5ZMy9FeTV0Uk5iRHBpc3RnWEVyelVLNEx4MjRzUVVsVGVRcThoZXpZL29QOFcyam5RcUhEM0p0Tit1VlBLODNPcVp4bnQwcitTSFVqN2J3Y0IifQ=="
		}
	}
	`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapter_GroupMember(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/groupmember")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"region":"us-west-2"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "groupMember",
			ExternalId: "GroupMember",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "UserId",
					ExternalId: "UserId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "GroupName",
					ExternalId: "GroupName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 2,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
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
									"string_value": "Group1"
								}
							],
							"id": "GroupName"
						},
						{
							"values": [
								{
									"string_value": "AIDAXXXXXXXXXXXXXXXX1"
								}
							],
							"id": "UserId"
						},
						{
							"values": [
								{
									"string_value": "AIDAXXXXXXXXXXXXXXXX1-Group1"
								}
							],
							"id": "id"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjb2xsZWN0aW9uSWQiOiJHcm91cDEiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiQUU0dlVQNGF1bU12aW9POWdjVmtmdXNmUUZDbjdlaVFhSm1XM3hPb0x4UVRCdDRYOSsvRjdSZU44Q0pSUythMEFUTHNNcDBJK3pxVjVqdDdxT1VhbGh3TFRIOGZDRHBOTUt5SDF6QT0ifQ=="
		}
	}
	`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapterWithMultipleAccounts_User(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/user_multi_account")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"resourceAccounts":["arn:aws:iam::82202838614:role/Cross-Account-Assume-Admin"],"region": "us-west-2","requestTimeoutSeconds": 300}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "user",
			ExternalId: "User",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Arn",
					ExternalId: "Arn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "AccountId",
					ExternalId: "AccountId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "Path",
					ExternalId: "Path",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "UserId",
					ExternalId: "UserId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "UserName",
					ExternalId: "UserName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 2,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
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
									"string_value": "000000000000"
								}
							],
							"id": "AccountId"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::000000000000:user/user1@example.com"
								}
							],
							"id": "Arn"
						},
						{
							"values": [
								{
									"string_value": "/"
								}
							],
							"id": "Path"
						},
						{
							"values": [
								{
									"string_value": "AIDAXXXXXXXXXXXXXXXX1"
								}
							],
							"id": "UserId"
						},
						{
							"values": [
								{
									"string_value": "user1@example.com"
								}
							],
							"id": "UserName"
						}
					],
					"child_objects": []
				},
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "000000000000"
								}
							],
							"id": "AccountId"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::000000000000:user/user2@example.com"
								}
							],
							"id": "Arn"
						},
						{
							"values": [
								{
									"string_value": "/"
								}
							],
							"id": "Path"
						},
						{
							"values": [
								{
									"string_value": "AIDAXXXXXXXXXXXXXXXX2"
								}
							],
							"id": "UserId"
						},
						{
							"values": [
								{
									"string_value": "user2@example.com"
								}
							],
							"id": "UserName"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJBRGNEd0t4RFJhbHNlSlN6eC9mVFRTMW9BNjVNVlREOHpUOTd5K2FHMG5pWm9XUEhkUWIraFJKRWhtT3dLY3FieW4valRvclVxVnNnNm9lNzYwakU0ZXRCaFZWQ1NmenQ3NTdVUTk0ODJGN1UifQ=="
		}
	}
	`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapterWithMultipleAccounts_Group(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/group_multi_account")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"resourceAccounts":["arn:aws:iam::82202838614:role/Cross-Account-Assume-Admin"],"region": "us-west-2","requestTimeoutSeconds": 300}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "group",
			ExternalId: "Group",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Arn",
					ExternalId: "Arn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "AccountId",
					ExternalId: "AccountId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "Path",
					ExternalId: "Path",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "GroupId",
					ExternalId: "GroupId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "GroupName",
					ExternalId: "GroupName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 2,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
	{
		"success": {
			"objects": [
			],
			"next_cursor": ""
		}
	}
	`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapterWithMultipleAccounts_Role(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/role_multi_account")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"resourceAccounts":["arn:aws:iam::82202838614:role/Cross-Account-Assume-Admin"],"region": "us-west-2","requestTimeoutSeconds": 300}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "role",
			ExternalId: "Role",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Arn",
					ExternalId: "Arn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "AccountId",
					ExternalId: "AccountId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "AssumeRolePolicyDocument",
					ExternalId: "AssumeRolePolicyDocument",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "MaxSessionDuration",
					ExternalId: "MaxSessionDuration",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
				},
				{
					Id:         "Path",
					ExternalId: "Path",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "PermissionsBoundary",
					ExternalId: "PermissionsBoundary",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "RoleId",
					ExternalId: "RoleId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "RoleName",
					ExternalId: "RoleName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 2,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
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
									"string_value": "000000000000"
								}
							],
							"id": "AccountId"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::000000000000:role/AdminAccessRole"
								}
							],
							"id": "Arn"
						},
						{
							"values": [
								{
									"string_value": "%7B%22Version%22%3A%222012-10-17%22%2C%22Statement%22%3A%5B%7B%22Effect%22%3A%22Allow%22%2C%22Principal%22%3A%7B%22Federated%22%3A%22arn%3Aaws%3Aiam%3A%3A000000000000%3Asaml-provider%2FCOMPANY_A-COMPANY_B%22%7D%2C%22Action%22%3A%22sts%3AAssumeRoleWithSAML%22%2C%22Condition%22%3A%7B%22StringEquals%22%3A%7B%22SAML%3Aaud%22%3A%22https%3A%2F%2Fsignin.aws.amazon.com%2Fsaml%22%7D%7D%7D%2C%7B%22Effect%22%3A%22Allow%22%2C%22Principal%22%3A%7B%22Federated%22%3A%22arn%3Aaws%3Aiam%3A%3A000000000000%3Asaml-provider%2Fsgnl-demo.oktapreview.com%22%7D%2C%22Action%22%3A%22sts%3AAssumeRoleWithSAML%22%2C%22Condition%22%3A%7B%22StringEquals%22%3A%7B%22SAML%3Aaud%22%3A%22https%3A%2F%2Fsignin.aws.amazon.com%2Fsaml%22%7D%7D%7D%5D%7D"
								}
							],
							"id": "AssumeRolePolicyDocument"
						},
						{
							"values": [
								{
									"int64_value": "3600"
								}
							],
							"id": "MaxSessionDuration"
						},
						{
							"values": [
								{
									"string_value": "/"
								}
							],
							"id": "Path"
						},
						{
							"values": [
								{
									"string_value": "AROAXXXXXXXXXXXXXXXX1"
								}
							],
							"id": "RoleId"
						},
						{
							"values": [
								{
									"string_value": "AdminAccessRole"
								}
							],
							"id": "RoleName"
						}
					],
					"child_objects": []
				},
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "000000000000"
								}
							],
							"id": "AccountId"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::000000000000:role/aws-service-role/support.amazonaws.com/AWSServiceRoleForSupport"
								}
							],
							"id": "Arn"
						},
						{
							"values": [
								{
									"string_value": "%7B%22Version%22%3A%222012-10-17%22%2C%22Statement%22%3A%5B%7B%22Effect%22%3A%22Allow%22%2C%22Principal%22%3A%7B%22Service%22%3A%22support.amazonaws.com%22%7D%2C%22Action%22%3A%22sts%3AAssumeRole%22%7D%5D%7D"
								}
							],
							"id": "AssumeRolePolicyDocument"
						},
						{
							"values": [
								{
									"int64_value": "3600"
								}
							],
							"id": "MaxSessionDuration"
						},
						{
							"values": [
								{
									"string_value": "/aws-service-role/support.amazonaws.com/"
								}
							],
							"id": "Path"
						},
						{
							"values": [
								{
									"string_value": "AROAXXXXXXXXXXXXXXXX2"
								}
							],
							"id": "RoleId"
						},
						{
							"values": [
								{
									"string_value": "AWSServiceRoleForSupport"
								}
							],
							"id": "RoleName"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJBRHJxZWFnRDVWZnA1Z0JTUDZTT1A5ZkdNUVNYbk05L1FlaXdORWJ2QmtCLzZ4ZjlzU3pUbDdlWmpCTy83ckZUM3Q2dDdTTGptdmgxWFVTT2R1b3VITkhaLzN6RHRNN3FGWGpnbVhaenp3cTdUTGVlZ0dST3BUcm9zQkYwZXVBZm4xcCt4YTV0cEpEYzhtemJGSVFIeVFUT2dyOD0ifQ=="
		}
	}
	`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapterWithMultipleAccounts_Policy(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/policy_multi_account")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"resourceAccounts":["arn:aws:iam::82202838614:role/Cross-Account-Assume-Admin"],"region": "us-west-2","requestTimeoutSeconds": 300}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "policy",
			ExternalId: "Policy",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "Arn",
					ExternalId: "Arn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "AccountId",
					ExternalId: "AccountId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "DefaultVersionId",
					ExternalId: "DefaultVersionId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "Description",
					ExternalId: "Description",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "AttachmentCount",
					ExternalId: "AttachmentCount",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
				},
				{
					Id:         "Path",
					ExternalId: "Path",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "IsAttachable",
					ExternalId: "IsAttachable",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "PolicyId",
					ExternalId: "PolicyId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "PermissionsBoundaryUsageCount",
					ExternalId: "PermissionsBoundaryUsageCount",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					Id:         "PolicyName",
					ExternalId: "PolicyName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 2,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
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
										"string_value": "000000000000"
									}
								],
								"id": "AccountId"
							},
							{
								"values": [
									{
										"string_value": "arn:aws:iam::000000000000:policy/ExampleS3AccessPolicy"
									}
								],
								"id": "Arn"
							},
							{
								"values": [
									{
										"int64_value": "1"
									}
								],
								"id": "AttachmentCount"
							},
							{
								"values": [
									{
										"string_value": "v3"
									}
								],
								"id": "DefaultVersionId"
							},
							{
								"values": [
									{
										"string_value": "Allow access to buckets tagged with CustomerID principal tags."
									}
								],
								"id": "Description"
							},
							{
								"values": [
									{
										"bool_value": true
									}
								],
								"id": "IsAttachable"
							},
							{
								"values": [
									{
										"string_value": "/"
									}
								],
								"id": "Path"
							},
							{
								"values": [
									{
										"int64_value": "0"
									}
								],
								"id": "PermissionsBoundaryUsageCount"
							},
							{
								"values": [
									{
										"string_value": "ANPAXXXXXXXXXXXXXXXX1"
									}
								],
								"id": "PolicyId"
							},
							{
								"values": [
									{
										"string_value": "ExampleS3AccessPolicy"
									}
								],
								"id": "PolicyName"
							}
						],
						"child_objects": []
					},
					{
						"attributes": [
							{
								"values": [
									{
										"string_value": "aws"
									}
								],
								"id": "AccountId"
							},
							{
								"values": [
									{
										"string_value": "arn:aws:iam::aws:policy/AdministratorAccess"
									}
								],
								"id": "Arn"
							},
							{
								"values": [
									{
										"int64_value": "8"
									}
								],
								"id": "AttachmentCount"
							},
							{
								"values": [
									{
										"string_value": "v1"
									}
								],
								"id": "DefaultVersionId"
							},
							{
								"values": [
									{
										"string_value": "Provides full access to AWS services and resources."
									}
								],
								"id": "Description"
							},
							{
								"values": [
									{
										"bool_value": true
									}
								],
								"id": "IsAttachable"
							},
							{
								"values": [
									{
										"string_value": "/"
									}
								],
								"id": "Path"
							},
							{
								"values": [
									{
										"int64_value": "0"
									}
								],
								"id": "PermissionsBoundaryUsageCount"
							},
							{
								"values": [
									{
										"string_value": "ANPAXXXXXXXXXXXXXXXX2"
									}
								],
								"id": "PolicyId"
							},
							{
								"values": [
									{
										"string_value": "AdministratorAccess"
									}
								],
								"id": "PolicyName"
							}
						],
						"child_objects": []
					}
				],
				"next_cursor": "eyJjdXJzb3IiOiJBRkgvUWhDanJ0Y2xiVy9UUUxacmRaWUxJVmYxMjRpeGFUVk81blR5U3BSN3dLaTRGRFlqTGdPR3pIeEVLaTE3d3Q5ZFJobi9pMlhRZnVrK2l4NUkrR3BwWXl3MGNOOTBtbUxOaHA2Q2RYbjc4MmRoYjV5ZXpOd3VrQ3JPL3E3VHM1TT0ifQ=="
			}
		}
		`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapterWithMultipleAccounts_GroupPolicy(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/grouppolicy_multi_account")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"resourceAccounts":["arn:aws:iam::82202838614:role/Cross-Account-Assume-Admin"],"region": "us-west-2","requestTimeoutSeconds": 300}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "groupPolicy",
			ExternalId: "GroupPolicy",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "PolicyArn",
					ExternalId: "PolicyArn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "GroupName",
					ExternalId: "GroupName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 2,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
	{
		"success": {
			"objects": [],
			"next_cursor": ""
		}
	}
	`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapterWithMultipleAccounts_RolePolicy(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/rolepolicy_multi_account")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"resourceAccounts":["arn:aws:iam::82202838614:role/Cross-Account-Assume-Admin"],"region": "us-west-2","requestTimeoutSeconds": 300}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "rolePolicy",
			ExternalId: "RolePolicy",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "PolicyArn",
					ExternalId: "PolicyArn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "RoleName",
					ExternalId: "RoleName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 1,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
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
									"string_value": "arn:aws:iam::aws:policy/AdministratorAccess"
								}
							],
							"id": "PolicyArn"
						},
						{
							"values": [
								{
									"string_value": "AdminAccessRole"
								}
							],
							"id": "RoleName"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::aws:policy/AdministratorAccess-AdminAccessRole"
								}
							],
							"id": "id"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjb2xsZWN0aW9uSWQiOiJBZG1pbkFjY2Vzc1JvbGUiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiQUV2b3ltWXhYMnJUYkdoc0pBdUlMekZWUk54NW8yR0ZiQWNpZTdROVdicnVHUTVwc1ZoRHQ0MkQ0c1dXUjVxYlVGOW1SWkw1UVpMSnphRk1yUnZTQWlYTTR6MjNXajNGY3IrN0diZXgzbThyNnlSbXlLRzdxRHZZNjdINkRWRTZzSklTaXhUZFRadTFnckk9In0="
		}
	}
	`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapterWithMultipleAccounts_UserPolicy(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/userpolicy_multi_account")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"resourceAccounts":["arn:aws:iam::82202838614:role/Cross-Account-Assume-Admin"],"region": "us-west-2","requestTimeoutSeconds": 300}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "userPolicy",
			ExternalId: "UserPolicy",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "PolicyArn",
					ExternalId: "PolicyArn",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "UserName",
					ExternalId: "UserName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 1,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
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
									"string_value": "arn:aws:iam::aws:policy/AdministratorAccess"
								}
							],
							"id": "PolicyArn"
						},
						{
							"values": [
								{
									"string_value": "user1@example.com"
								}
							],
							"id": "UserName"
						},
						{
							"values": [
								{
									"string_value": "arn:aws:iam::aws:policy/AdministratorAccess-user1@example.com"
								}
							],
							"id": "id"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOiJBR0lMNkxQM2lBSU9hczZlSW5XRkdwS3RqQVVGSWxuRE9LQ1lCU2JPOW54WXhxYnZicSt2MWowU1hXSTY3WXcrcXc3K09KZEtJZUc3WFhha2xCRStmK0FtZEt1dTA2S0daOU9DQy9wVklmT0lHeDN0ZklPNE1KN0FuVFhMQUQyZmlxR1FsWFI4Zmc9PSIsImNvbGxlY3Rpb25JZCI6ImRhbW9uQHNnbmwuYWkiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiQUU3NzdLUlZSQ0F6aGY3ZmJlcWllSmMrOC9weXVyMU1lakhEZlhEeDJkTWFqRitXN0tRNnl5NHNXek13a3BVam5BdlVCRGNaVlMvZFhxdS9QZWR2YXFqa296em5FdDhqTm1Vb21QVXloTmtmIn0="
		}
	}
	`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestAWSAdapterWithMultipleAccounts_GroupMember(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws/groupmember_multi_accounts")
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
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Id:     "Test",
			Type:   "AWS-1.0.0",
			Config: []byte(`{"resourceAccounts":["arn:aws:iam::82202838614:role/Cross-Account-Assume-Admin"],"region": "us-west-2","requestTimeoutSeconds": 300}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "groupMember",
			ExternalId: "GroupMember",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					UniqueId:   true,
				},
				{
					Id:         "UserId",
					ExternalId: "UserId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "GroupName",
					ExternalId: "GroupName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
		},
		PageSize: 2,
	}

	gotResp, err := adapterClient.GetPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`
	{
		"success": {
			"objects": [],
			"next_cursor": ""
		}
	}
	`), wantResp)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp.GetSuccess().Objects, wantResp.GetSuccess().Objects, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}
