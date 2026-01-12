// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package identitynow_test

import (
	"context"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	identitynow_adapter "github.com/sgnl-ai/adapters/pkg/identitynow"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := identitynow_adapter.NewAdapter(&identitynow_adapter.Datasource{
		Client:                    server.Client(),
		AccountCollectionPageSize: 5,
	})

	tests := map[string]struct {
		ctx          context.Context
		request      *framework.Request[identitynow_adapter.Config]
		wantResponse framework.Response
	}{
		"valid_request_first_page_accounts": {
			ctx: context.Background(),
			request: &framework.Request[identitynow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "created",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
						{
							ExternalId: "$.attributes.user_name",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.attributes.roles",
							Type:       framework.AttributeTypeString,
							List:       true,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                     "1a1bb825eb7e4f76b72fbecb27699b31",
							"created":                time.Date(2023, 9, 22, 16, 46, 54, 250000000, time.UTC),
							"$.attributes.user_name": "victor.mcintosh",
						},
						{
							"id":                     "ba699287e60b4014bcc4319f30e9b59e",
							"created":                time.Date(2023, 9, 22, 16, 47, 35, 702000000, time.UTC),
							"$.attributes.user_name": "Cleo.Yoder",
							"$.attributes.roles":     []string{"e098ecf6c0a80165002aaec84d906014"},
						},
					},
					// {"cursor":2}
					NextCursor: "eyJjdXJzb3IiOjJ9",
				},
			},
		},
		"valid_request_last_page_accounts": {
			ctx: context.Background(),
			request: &framework.Request[identitynow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
				// {"cursor":2}
				Cursor: "eyJjdXJzb3IiOjJ9",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id": "28b1e9bf40ab458981067f4e4dc330b3",
						},
					},
				},
			},
		},
		"valid_request_first_page_entitlements": {
			ctx: context.Background(),
			request: &framework.Request[identitynow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3", // Will be overwritten by the entity config.
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"entitlements": {
							UniqueIDAttribute: "id",
							APIVersion:        testutil.GenPtr[string]("beta"),
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "entitlements",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "created",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
						{
							ExternalId: "$.attributes.displayName",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                       "ENTITLEMENT_ID_456",
							"created":                  time.Date(2023, 9, 22, 16, 50, 12, 53000000, time.UTC),
							"$.attributes.displayName": "Basic purchaser [on] Windows Store for Business",
						},
						{
							"id":                       "00218206fe614e7da637f528accdf15e",
							"created":                  time.Date(2023, 9, 22, 16, 50, 54, 856000000, time.UTC),
							"$.attributes.displayName": "default access [on] TrustedPublishersProxyService",
						},
					},
					// {"cursor":2}
					NextCursor: "eyJjdXJzb3IiOjJ9",
				},
			},
		},
		// The majority of these AccountEntitlement pagination tests are covered in datasource_test.go.
		// They're duplicated here to ensure the marshaling/unmarshaling of the composite cursor works.
		"valid_request_first_account_first_page_account_entitlements": {
			ctx: context.Background(),
			request: &framework.Request[identitynow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accountEntitlements": {
							UniqueIDAttribute: "id",
							APIVersion:        testutil.GenPtr[string]("beta"),
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "accountEntitlements",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "accountId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "entitlementId",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 10,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":            "testaccountId1-entitlementId1",
							"accountId":     "testaccountId1",
							"entitlementId": "entitlementId1",
						},
						{
							"id":            "testaccountId1-entitlementId2",
							"accountId":     "testaccountId1",
							"entitlementId": "entitlementId2",
						},
						{
							"id":            "testaccountId2-entitlementId3",
							"accountId":     "testaccountId2",
							"entitlementId": "entitlementId3",
						},
						{
							"id":            "testaccountId2-entitlementId4",
							"accountId":     "testaccountId2",
							"entitlementId": "entitlementId4",
						},
						{
							"id":            "testaccountId2-entitlementId5",
							"accountId":     "testaccountId2",
							"entitlementId": "entitlementId5",
						},
						{
							"id":            "testaccountId3-entitlementId6",
							"accountId":     "testaccountId3",
							"entitlementId": "entitlementId6",
						},
						{
							"id":            "testaccountId3-entitlementId7",
							"accountId":     "testaccountId3",
							"entitlementId": "entitlementId7",
						},
						{
							"id":            "testaccountId3-entitlementId8",
							"accountId":     "testaccountId3",
							"entitlementId": "entitlementId8",
						},
						{
							"id":            "testaccountId4-entitlementId9",
							"accountId":     "testaccountId4",
							"entitlementId": "entitlementId9",
						},
						{
							"id":            "testaccountId4-entitlementId10",
							"accountId":     "testaccountId4",
							"entitlementId": "entitlementId10",
						},
					},
					// {"collectionId":"testaccountId4","collectionCursor":4}
					NextCursor: "eyJjb2xsZWN0aW9uSWQiOiJ0ZXN0YWNjb3VudElkNCIsImNvbGxlY3Rpb25DdXJzb3IiOjR9",
				},
			},
		},
		"valid_request_no_https_prefix": {
			ctx: context.Background(),
			request: &framework.Request[identitynow_adapter.Config]{
				Address: strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id": "1a1bb825eb7e4f76b72fbecb27699b31",
						},
						{
							"id": "ba699287e60b4014bcc4319f30e9b59e",
						},
					},
					// {"cursor":2}
					NextCursor: "eyJjdXJzb3IiOjJ9",
				},
			},
		},
		"invalid_request_invalid_api_version": {
			ctx: context.Background(),
			request: &framework.Request[identitynow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
							APIVersion:        testutil.GenPtr[string]("v1"),
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "IdentityNow config is invalid: apiVersion v1 for entity accounts is not supported.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"invalid_request_http_prefix": {
			ctx: context.Background(),
			request: &framework.Request[identitynow_adapter.Config]{
				Address: "http://" + strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "The provided HTTP protocol is not supported.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"invalid_request_unable_to_parse_cursor": {
			ctx: context.Background(),
			request: &framework.Request[identitynow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
				Cursor:   "NOT_B64_ENCODED_COMPOSITE_CURSOR",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to decode base64 cursor: illegal base64 data at input byte 3.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
		"invalid_request_not_found": {
			ctx: context.Background(),
			request: &framework.Request[identitynow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
				// {"cursor":99}
				Cursor: "eyJjdXJzb3IiOjk5fQ==",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Datasource rejected request, returned status code: 404.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"valid_request_with_filter": {
			ctx: context.Background(),
			request: &framework.Request[identitynow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
							Filter:            testutil.GenPtr[string](`id eq "1a1bb825eb7e4f76b72fbecb27699b31"`),
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id": "1a1bb825eb7e4f76b72fbecb27699b31",
						},
					},
				},
			},
		},
		"valid_request_with_concatenated_attribute": {
			ctx: context.Background(),
			request: &framework.Request[identitynow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "groupsMembership",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "memberOfMembership",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "GroupsMembership",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				// {"cursor":97}
				Cursor:   "eyJjdXJzb3IiOjk3fQ==",
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                 "1a1bb825eb7e4f76b72fbecb27699b31",
							"memberOfMembership": "GROUP1 | GROUP2 | GROUP3",
							"GroupsMembership":   "GROUP1 | GROUP2 | GROUP3",
							"groupsMembership":   "GROUP1 | GROUP2 | GROUP3",
						},
					},
				},
			},
		},
		"invalid_request_attribute_not_string_array": {
			ctx: context.Background(),
			request: &framework.Request[identitynow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer token",
				},
				Config: &identitynow_adapter.Config{
					APIVersion: "v3",
					EntityConfig: map[string]identitynow_adapter.EntityConfig{
						"accounts": {
							UniqueIDAttribute: "id",
							APIVersion:        testutil.GenPtr[string]("v3"),
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "accounts",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "groupsMembership",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "memberOfMembership",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "GroupsMembership",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				// {"cursor":98}
				Cursor:   "eyJjdXJzb3IiOjk4fQ==",
				Ordered:  false,
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":               "1a1bb825eb7e4f76b72fbecb27699b31",
							"GroupsMembership": "GROUP1 | GROUP2 | GROUP3",
						},
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
