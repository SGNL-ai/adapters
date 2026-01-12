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

func TestDuoAdapter_User(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/duo/user")
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
			Address: "test-instance.duosecurity.com",
			Id:      "Test",
			Type:    "Duo-1.0.0",
			Config:  []byte(`{"apiVersion":"v1"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "User",
			ExternalId: "User",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "user_id",
					ExternalId: "user_id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "realname",
					ExternalId: "realname",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "created",
					ExternalId: "created",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "is_enrolled",
					ExternalId: "is_enrolled",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "alias1",
					ExternalId: "alias1",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "alias2",
					ExternalId: "alias2",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "alias3",
					ExternalId: "alias3",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "alias4",
					ExternalId: "alias4",
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
					Id:         "firstname",
					ExternalId: "firstname",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "last_directory_sync",
					ExternalId: "last_directory_sync",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "last_login",
					ExternalId: "last_login",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "lastname",
					ExternalId: "lastname",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "lockout_reason",
					ExternalId: "lockout_reason",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "notes",
					ExternalId: "notes",
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
					Id:         "username",
					ExternalId: "username",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					Id:         "Group",
					ExternalId: "groups",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "group_id",
							ExternalId: "group_id",
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
							Id:         "desc",
							ExternalId: "desc",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "mobile_otp_enabled",
							ExternalId: "mobile_otp_enabled",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
							List:       false,
						},
						{
							Id:         "push_enabled",
							ExternalId: "push_enabled",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
							List:       false,
						},
						{
							Id:         "sms_enabled",
							ExternalId: "sms_enabled",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
							List:       false,
						},
						{
							Id:         "status",
							ExternalId: "status",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "voice_enabled",
							ExternalId: "voice_enabled",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
							List:       false,
						},
					},
				},
				{
					Id:         "Phone",
					ExternalId: "phones",
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "phone_id",
							ExternalId: "phone_id",
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
							Id:         "activated",
							ExternalId: "activated",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
							List:       false,
						},
						{
							Id:         "capabilities",
							ExternalId: "capabilities",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       true,
						},
						{
							Id:         "encrypted",
							ExternalId: "encrypted",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "extension",
							ExternalId: "extension",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "fingerprint",
							ExternalId: "fingerprint",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "last_seen",
							ExternalId: "last_seen",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
							List:       false,
						},
						{
							Id:         "model",
							ExternalId: "model",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "number",
							ExternalId: "number",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "platform",
							ExternalId: "platform",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "postdelay",
							ExternalId: "postdelay",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "predelay",
							ExternalId: "predelay",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "screenlock",
							ExternalId: "screenlock",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							List:       false,
						},
						{
							Id:         "sms_passcodes_sent",
							ExternalId: "sms_passcodes_sent",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
							List:       false,
						},
						{
							Id:         "tampered",
							ExternalId: "tampered",
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
		"success": {
			"objects": [
				{
					"attributes": [
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-01-23T20:17:36Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "created"
						},
						{
							"values": [
								{
									"string_value": "user1@example.com"
								}
							],
							"id": "email"
						},
						{
							"values": [
								{
									"bool_value": true
								}
							],
							"id": "is_enrolled"
						},
						{
							"values": [
								{
									"string_value": ""
								}
							],
							"id": "notes"
						},
						{
							"values": [
								{
									"string_value": "Test User 1"
								}
							],
							"id": "realname"
						},
						{
							"values": [
								{
									"string_value": "active"
								}
							],
							"id": "status"
						},
						{
							"values": [
								{
									"string_value": "DUYC8O4O953VBGGKLHAL"
								}
							],
							"id": "user_id"
						},
						{
							"values": [
								{
									"string_value": "user1"
								}
							],
							"id": "username"
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
													"string_value": "random group 1"
												}
											],
											"id": "desc"
										},
										{
											"values": [
												{
													"string_value": "DGKUKQSTG7ZFDN2N1XID"
												}
											],
											"id": "group_id"
										},
										{
											"values": [
												{
													"bool_value": false
												}
											],
											"id": "mobile_otp_enabled"
										},
										{
											"values": [
												{
													"string_value": "group1"
												}
											],
											"id": "name"
										},
										{
											"values": [
												{
													"bool_value": false
												}
											],
											"id": "push_enabled"
										},
										{
											"values": [
												{
													"bool_value": false
												}
											],
											"id": "sms_enabled"
										},
										{
											"values": [
												{
													"string_value": "Active"
												}
											],
											"id": "status"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "Group"
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
											"id": "activated"
										},
										{
											"values": [
												{
													"string_value": "sms"
												}
											],
											"id": "capabilities"
										},
										{
											"values": [
												{
													"string_value": "Encrypted"
												}
											],
											"id": "encrypted"
										},
										{
											"values": [
												{
													"string_value": ""
												}
											],
											"id": "extension"
										},
										{
											"values": [
												{
													"string_value": "Unknown"
												}
											],
											"id": "fingerprint"
										},
										{
											"values": [
												{
													"string_value": "Unknown"
												}
											],
											"id": "model"
										},
										{
											"values": [
												{
													"string_value": ""
												}
											],
											"id": "name"
										},
										{
											"values": [
												{
													"string_value": "+11111111111"
												}
											],
											"id": "number"
										},
										{
											"values": [
												{
													"string_value": "DPFL36P8Z8LZANN1FFEZ"
												}
											],
											"id": "phone_id"
										},
										{
											"values": [
												{
													"string_value": "Apple iOS"
												}
											],
											"id": "platform"
										},
										{
											"values": [
												{
													"string_value": "Unknown"
												}
											],
											"id": "screenlock"
										},
										{
											"values": [
												{
													"bool_value": false
												}
											],
											"id": "sms_passcodes_sent"
										},
										{
											"values": [
												{
													"string_value": "Unknown"
												}
											],
											"id": "tampered"
										},
										{
											"values": [
												{
													"string_value": "Mobile"
												}
											],
											"id": "type"
										}
									],
									"child_objects": []
								}
							],
							"entity_id": "Phone"
						}
					]
				}
			],
			"next_cursor": "eyJjdXJzb3IiOjF9"
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

func TestDuoAdapter_Group(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/duo/group")
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
			Address: "test-instance.duosecurity.com",
			Id:      "Test",
			Type:    "Duo-1.0.0",
			Config:  []byte(`{"apiVersion":"v1"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Group",
			ExternalId: "Group",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "group_id",
					ExternalId: "group_id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "name",
					ExternalId: "name",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
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
		"success": {
			"objects": [
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "DGKUKQSTG7ZFDN2N1XID"
								}
							],
							"id": "group_id"
						},
						{
							"values": [
								{
									"string_value": "group1"
								}
							],
							"id": "name"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOjF9"
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

func TestDuoAdapter_Phone(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/duo/phone")
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
			Address: "test-instance.duosecurity.com",
			Id:      "Test",
			Type:    "Duo-1.0.0",
			Config:  []byte(`{"apiVersion":"v1"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Phone",
			ExternalId: "Phone",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "phone_id",
					ExternalId: "phone_id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "name",
					ExternalId: "name",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
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
		"success": {
			"objects": [
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": ""
								}
							],
							"id": "name"
						},
						{
							"values": [
								{
									"string_value": "DPFL36P8Z8LZANN1FFEZ"
								}
							],
							"id": "phone_id"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOjF9"
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
