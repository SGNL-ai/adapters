// Copyright 2026 SGNL.ai, Inc.

package salesforce_test

import (
	"context"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	salesforce_adapter "github.com/sgnl-ai/adapters/pkg/salesforce"
)

// sortChildEntitiesByID sorts child entity arrays by their "id" field to enable order-independent comparison.
// This is needed because map iteration order is non-deterministic in Go, but the test needs consistent results.
func sortChildEntitiesByID(objects []framework.Object) {
	for _, obj := range objects {
		for key, value := range obj {
			// Check if this field is a child entity array
			if childArray, ok := value.([]framework.Object); ok {
				sort.Slice(childArray, func(i, j int) bool {
					id1, _ := childArray[i]["id"].(string)
					id2, _ := childArray[j]["id"].(string)

					return id1 < id2
				})
				// Update the object with sorted array
				obj[key] = childArray
			}
		}
	}
}

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
					Message: `Scheme "http" is not supported.`,
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
		"valid_request_with_custom_fields": {
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
					ExternalId: "CustomObject",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "CustomField__c",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "AnotherCustom__c",
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
							"Id":               "a00Hu000000AbCDEF",
							"CustomField__c":   "CustomValue1",
							"AnotherCustom__c": "CustomValue2",
						},
					},
				},
			},
		},
		"valid_request_with_custom_fields_jsonpath": {
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
					ExternalId: "CustomObject",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.CustomField__c",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.AnotherCustom__c",
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
							"Id":                 "a00Hu000000AbCDEF",
							"$.CustomField__c":   "CustomValue1",
							"$.AnotherCustom__c": "CustomValue2",
						},
					},
				},
			},
		},
		"valid_request_with_jsonpath_relationship": {
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
						{
							ExternalId: "$.Name",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.Account.Name",
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
							"Id":             "003Hu000000AbCDEF",
							"$.Name":         "John Doe",
							"$.Account.Name": "Acme Corporation",
						},
					},
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
					ExternalId: "Account",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "Interests__c",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "id",
									Type:       framework.AttributeTypeString,
									UniqueId:   true,
								},
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
								{"id": "003Hu000020yLuHIAU_Interests__c_sports", "value": "Sports"},
								{"id": "003Hu000020yLuHIAU_Interests__c_music", "value": "Music"},
								{"id": "003Hu000020yLuHIAU_Interests__c_reading", "value": "Reading"},
							},
						},
						{
							"Id": "003Hu000020yLuMIAU",
							"Interests__c": []framework.Object{
								{"id": "003Hu000020yLuMIAU_Interests__c_technology", "value": "Technology"},
							},
						},
						{
							"Id": "003Hu000020yLuPIAU",
						},
					},
				},
			},
		},
		"valid_request_with_list_and_complex_object_child_entities": {
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
					ExternalId: "Account",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "Interests__c",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "id",
									Type:       framework.AttributeTypeString,
									UniqueId:   true,
								},
								{
									ExternalId: "value",
									Type:       framework.AttributeTypeString,
								},
							},
						},
						{
							ExternalId: "Tags__c",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "Name",
									Type:       framework.AttributeTypeString,
									UniqueId:   true,
								},
								{
									ExternalId: "Priority",
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
							"Id": "001Hu000020yLuXYZ",
							"Interests__c": []framework.Object{
								{"id": "001Hu000020yLuXYZ_Interests__c_sports", "value": "Sports"},
								{"id": "001Hu000020yLuXYZ_Interests__c_music", "value": "Music"},
							},
							"Tags__c": []framework.Object{
								{"Name": "VIP", "Priority": "High"},
								{"Name": "Region", "Priority": "Medium"},
							},
						},
						{
							"Id": "001Hu000020yLuABC",
							"Interests__c": []framework.Object{
								{"id": "001Hu000020yLuABC_Interests__c_technology", "value": "Technology"},
							},
							"Tags__c": []framework.Object{
								{"Name": "Status", "Priority": "Low"},
							},
						},
					},
				},
			},
		},
		"valid_request_with_missing_picklist_field": {
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
					ExternalId: "Account",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "Id",
							Type:       framework.AttributeTypeString,
							UniqueId:   true,
						},
						{
							ExternalId: "Name",
							Type:       framework.AttributeTypeString,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "Locations__c",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "id",
									Type:       framework.AttributeTypeString,
									UniqueId:   true,
								},
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
							"Id":   "001Hu000020yLuJKL",
							"Name": "Account Without Locations",
						},
						{
							"Id":   "001Hu000020yLuMNO",
							"Name": "Another Account",
							"Locations__c": []framework.Object{
								{"id": "001Hu000020yLuMNO_Locations__c_seattle", "value": "Seattle"},
								{"id": "001Hu000020yLuMNO_Locations__c_portland", "value": "Portland"},
							},
						},
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			// For multi-select picklist tests, sort child entities before comparison
			if name == "valid_request_with_multi_select_picklist" ||
				name == "valid_request_with_list_and_complex_object_child_entities" ||
				name == "valid_request_with_missing_picklist_field" {
				sortChildEntitiesByID(gotResponse.Success.Objects)
				sortChildEntitiesByID(tt.wantResponse.Success.Objects)
			}

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
