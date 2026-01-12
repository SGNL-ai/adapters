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

func TestS3Adapter(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/aws-s3/user")
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
			Type:   "S3-1.0.0",
			Config: []byte(`{"region":"us-east-1","bucket":"example-bucket"}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "User",
			ExternalId: "users",
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "userId",
					ExternalId: "userId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
					UniqueId:   true,
				},
				{
					Id:         "userName",
					ExternalId: "userName",
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

	err = protojson.Unmarshal([]byte(`{
		"success": {
			"objects": [
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "1"
								}
							],
							"id": "userId"
						},
						{
							"values": [
								{
									"string_value": "user_1"
								}
							],
							"id": "userName"
						}
					]
				},
				{
					"attributes": [
						{
							"values": [
								{
									"string_value": "2"
								}
							],
							"id": "userId"
						},
						{
							"values": [
								{
									"string_value": "user_2"
								}
							],
							"id": "userName"
						}
					]
				}
			],
			"next_cursor": "eyJjdXJzb3IiOjM0fQ=="
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
