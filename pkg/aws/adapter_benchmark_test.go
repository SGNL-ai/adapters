// Copyright 2026 SGNL.ai, Inc.

// nolint: goconst

package aws_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	aws_adapter "github.com/sgnl-ai/adapters/pkg/aws"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// BenchmarkAdapterGetUserPage benchmarks the GetPage method for the Policy entity.
// [SCENARIO]
//
//	Number of Records: 1000
//	PageSize: 1000
//	Artificial delay: 100 Millisecond/record
//
// [BEFORE]
//
//	goos: linux
//	goarch: amd64
//	pkg: github.com/sgnl-ai/adapters/pkg/aws
//	cpu: AMD Ryzen 5 5600U with Radeon Graphics
//	BenchmarkAdapterGetPolicyPage/valid_request-12                 1        101056236826 ns/op
//	PASS
//	ok      github.com/sgnl-ai/adapters/pkg/aws 102.346s
//
// [AFTER]
//
//	goos: linux
//	goarch: amd64
//	pkg: github.com/sgnl-ai/adapters/pkg/aws
//	cpu: AMD Ryzen 5 5600U with Radeon Graphics
//	BenchmarkAdapterGetPolicyPage/valid_request-12                 1        5162916760 ns/op
//	PASS
//	ok      github.com/sgnl-ai/adapters/pkg/aws 5.389s
func BenchmarkAdapterGetPolicyPage(b *testing.B) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(defaultTimeout)*time.Second)
	defer cancel()

	cfg, err := SetupTestConfig(ctx, Mocker)
	if err != nil {
		log.Fatalf("Failed to load aws test config: %v", err)
	}

	client, err := aws_adapter.NewClient(http.DefaultClient, cfg, 20)
	if err != nil {
		b.Errorf("error creating client to query datasource: %v", err)
	}

	adapter := aws_adapter.NewAdapter(client)

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[aws_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[aws_adapter.Config]{
				Auth:   validAuthCredentials,
				Config: validCommonConfig,
				Entity: framework.EntityConfig{
					ExternalId: "Policy",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Arn",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "AccountId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "AttachmentCount",
							Type:       framework.AttributeTypeInt64,
						},
						{
							ExternalId: "PolicyName",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "Path",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "DefaultVersionId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "PolicyId",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "IsAttachable",
							Type:       framework.AttributeTypeBool,
						},
						{
							ExternalId: "CreateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "UpdateDate",
							Type:       framework.AttributeTypeDateTime,
						},
						{
							ExternalId: "PermissionsBoundaryUsageCount",
							Type:       framework.AttributeTypeInt64,
						},
					},
				},
				PageSize: 1000,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: largePolicyObjects(),
				},
			},
		},
	}

	for name, tt := range tests {
		b.Run(name, func(b *testing.B) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					b.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				b.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
			if gotResponse.Success != nil && tt.wantCursor != nil {
				var gotCursor pagination.CompositeCursor[string]

				decodedCursor, err := base64.StdEncoding.DecodeString(gotResponse.Success.NextCursor)
				if err != nil {
					b.Errorf("error decoding cursor: %v", err)
				}

				if err := json.Unmarshal(decodedCursor, &gotCursor); err != nil {
					b.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(&gotCursor, tt.wantCursor) {
					b.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}
		})
	}
}
