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

func TestBambooHRAdapter_Employee(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/bamboohr/employee")
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
			Address: "api.bamboohr.com/api/gateway.php",
			Id:      "Test",
			Type:    "BambooHR-1.0.0",
			Config:  []byte(`{"apiVersion":"v1","companyDomain":"sgnltestdev","onlyCurrent":true,"boolAttributeMappings":{"true":["yes"],"false":["no"]}}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "Employee",
			ExternalId: "Employee",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					Id:         "bestEmail",
					ExternalId: "bestEmail",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "dateOfBirth",
					ExternalId: "dateOfBirth",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
				{
					Id:         "fullName",
					ExternalId: "fullName1",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "isPhotoUploaded",
					ExternalId: "isPhotoUploaded",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "checkboxBoolField",
					ExternalId: "customcustomBoolField",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					List:       false,
				},
				{
					Id:         "supervisorEId",
					ExternalId: "supervisorEId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					List:       false,
				},
				{
					Id:         "supervisorEmail",
					ExternalId: "supervisorEmail",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "lastChanged",
					ExternalId: "lastChanged",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
					List:       false,
				},
			},
		},
		PageSize: 2,
		Cursor:   "eyJjdXJzb3IiOjJ9",
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
									"string_value": "cagluinda@efficientoffice.com"
								}
							],
							"id": "bestEmail"
						},
						{
							"values": [
								{
									"bool_value": false
								}
							],
							"id": "checkboxBoolField"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "1996-08-27T00:00:00Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "dateOfBirth"
						},
						{
							"values": [
								{
									"string_value": "Christina Agluinda"
								}
							],
							"id": "fullName"
						},
						{
							"values": [
								{
									"int64_value": "6"
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
							"id": "isPhotoUploaded"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-04-12T19:33:50+00:00",
										"timezone_offset": 0
									}
								}
							],
							"id": "lastChanged"
						},
						{
							"values": [
								{
									"int64_value": "9"
								}
							],
							"id": "supervisorEId"
						},
						{
							"values": [
								{
									"string_value": "jcaldwell@efficientoffice.com"
								}
							],
							"id": "supervisorEmail"
						}
					],
					"child_objects": []
				},
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "sanderson@efficientoffice.com"
								}
							],
							"id": "bestEmail"
						},
						{
							"values": [
								{
									"bool_value": false
								}
							],
							"id": "checkboxBoolField"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2000-05-08T00:00:00Z",
										"timezone_offset": 0
									}
								}
							],
							"id": "dateOfBirth"
						},
						{
							"values": [
								{
									"string_value": "Shannon Anderson"
								}
							],
							"id": "fullName"
						},
						{
							"values": [
								{
									"int64_value": "7"
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
							"id": "isPhotoUploaded"
						},
						{
							"values": [
								{
									"datetime_value": {
										"timestamp": "2024-04-12T19:33:50+00:00",
										"timezone_offset": 0
									}
								}
							],
							"id": "lastChanged"
						},
						{
							"values": [
								{
									"int64_value": "9"
								}
							],
							"id": "supervisorEId"
						},
						{
							"values": [
								{
									"string_value": "jcaldwell@efficientoffice.com"
								}
							],
							"id": "supervisorEmail"
						}
					],
					"child_objects": []
				}
			],
			"next_cursor": "eyJjdXJzb3IiOjR9"
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
