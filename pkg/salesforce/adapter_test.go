// Copyright 2025 SGNL.ai, Inc.
package salesforce_test

import (
	"context"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	salesforce_adapter "github.com/sgnl-ai/adapters/pkg/salesforce"
)

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := salesforce_adapter.NewAdapter(&salesforce_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx          context.Context
		request      *framework.Request[salesforce_adapter.Config]
		wantResponse framework.Response
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[salesforce_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Case",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "CaseNumber",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "Status",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  true,
				PageSize: 200,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Id":         "500Hu000020yLuHIAU",
							"CaseNumber": "00001026",
							"Status":     "Closed",
						},
						{
							"Id":         "500Hu000020yLuMIAU",
							"CaseNumber": "00001027",
							"Status":     "New",
						},
					},
					NextCursor: "/services/data/v58.0/query/0r8Hu1lKCluUiC9IMK-200",
				},
			},
		},
		"valid_request_no_https_prefix": {
			ctx: context.Background(),
			request: &framework.Request[salesforce_adapter.Config]{
				Address: strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Case",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "CaseNumber",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "Status",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  true,
				PageSize: 200,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Id":         "500Hu000020yLuHIAU",
							"CaseNumber": "00001026",
							"Status":     "Closed",
						},
						{
							"Id":         "500Hu000020yLuMIAU",
							"CaseNumber": "00001027",
							"Status":     "New",
						},
					},
					NextCursor: "/services/data/v58.0/query/0r8Hu1lKCluUiC9IMK-200",
				},
			},
		},
		"invalid_request_invalid_api_version": {
			ctx: context.Background(),
			request: &framework.Request[salesforce_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "51.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Case",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "CaseNumber",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "Status",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  true,
				PageSize: 200,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Salesforce config is invalid: apiVersion is not supported: 51.0.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"invalid_request_http_prefix": {
			ctx: context.Background(),
			request: &framework.Request[salesforce_adapter.Config]{
				Address: "http://" + strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Case",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "CaseNumber",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "Status",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  true,
				PageSize: 200,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "The provided HTTP protocol is not supported.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"valid_request_with_filter": {
			ctx: context.Background(),
			request: &framework.Request[salesforce_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
					Filters: map[string]string{
						"Case": "Status = 'Closed'",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Case",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "CaseNumber",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "Status",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  true,
				PageSize: 200,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Id":         "500Hu000020yLuHIAU",
							"CaseNumber": "00001026",
							"Status":     "Closed",
						},
					},
				},
			},
		},
		"valid_request_with_cursor": {
			ctx: context.Background(),
			request: &framework.Request[salesforce_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Case",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "CaseNumber",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "Status",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  true,
				PageSize: 200,
				Cursor:   "/services/data/v58.0/query/0r8Hu1lKCluUiC9IMK-200",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Id":         "500Hu000020yLyEIAU",
							"CaseNumber": "00001031",
							"Status":     "New",
						},
						{
							"Id":         "500Hu000020yLyKIAU",
							"CaseNumber": "00001051",
							"Status":     "New",
						},
					},
				},
			},
		},
		"invalid_request_invalid_url": {
			ctx: context.Background(),
			request: &framework.Request[salesforce_adapter.Config]{
				Address: server.URL + "/invalid",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Case",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "CaseNumber",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "Status",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  true,
				PageSize: 200,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Datasource rejected request, returned status code: 404.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"parser_error_invalid_datetime_format": {
			ctx: context.Background(),
			request: &framework.Request[salesforce_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Case",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "CreatedAt",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
					},
				},
				Ordered:  true,
				PageSize: 200,
				Cursor:   "/services/data/v58.0/query/0r8Hu1lKClCJd892jd-200",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to convert datasource response objects: attribute CreatedAt cannot be parsed " +
						"into a date-time value: failed to parse date-time value: 2021/01/01 00:00:00.000Z.",
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"valid_request_with_multi_select_picklist": {
			ctx: context.Background(),
			request: &framework.Request[salesforce_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &salesforce_adapter.Config{
					APIVersion: "58.0",
				},
				Entity: framework.EntityConfig{
					ExternalId: "Contact",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "Interests__c",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "value",
									Type:       framework.AttributeTypeString,
								},
							},
						},
					},
				},
				Ordered:  true,
				PageSize: 200,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"Id": "003Hu000020yLuHIAU",
							"Interests__c": []framework.Object{
								{"value": "Sports"},
								{"value": "Music"},
								{"value": "Reading"},
							},
						},
						{
							"Id": "003Hu000020yLuMIAU",
							"Interests__c": []framework.Object{
								{"value": "Technology"},
							},
						},
						{
							"Id":           "003Hu000020yLuPIAU",
							"Interests__c": []framework.Object{},
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
