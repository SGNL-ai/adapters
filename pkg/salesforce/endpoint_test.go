// Copyright 2025 SGNL.ai, Inc.
package salesforce

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestConstructEndpoint(t *testing.T) {
	tests := map[string]struct {
		request      *Request
		wantEndpoint string
	}{
		"simple": {
			request: &Request{
				BaseURL:          "https://test.salesforce.com",
				APIVersion:       "52.0",
				EntityExternalID: "Account",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "Id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "Name",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: "https://test.salesforce.com/services/data/v52.0/query?q=SELECT+Id,Name+FROM+Account+ORDER+BY+Id+ASC",
		},
		"simple_with_filter": {
			request: &Request{
				BaseURL:          "https://test.salesforce.com",
				APIVersion:       "52.0",
				EntityExternalID: "Account",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "Id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "Name",
						Type:       framework.AttributeTypeString,
					},
				},
				Filter: testutil.GenPtr("Name LIKE 'Sample%'"),
			},
			wantEndpoint: "https://test.salesforce.com/services/data/v52.0/query?q=SELECT+Id,Name+FROM+Account+" +
				"WHERE+Name+LIKE+%27Sample%25%27+ORDER+BY+Id+ASC",
		},
		"nil_request": {
			request:      nil,
			wantEndpoint: "",
		},
		"simple_with_cursor": {
			request: &Request{
				BaseURL:          "https://test.salesforce.com",
				APIVersion:       "52.0",
				EntityExternalID: "Account",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "Id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "Name",
						Type:       framework.AttributeTypeString,
					},
				},
				Cursor: testutil.GenPtr("/services/data/v58.0/query/0r8Hu1lKCluUiC9IMK-200"),
			},
			wantEndpoint: "https://test.salesforce.com/services/data/v58.0/query/0r8Hu1lKCluUiC9IMK-200",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotEndpoint := ConstructEndpoint(tt.request)

			if !reflect.DeepEqual(gotEndpoint, tt.wantEndpoint) {
				t.Errorf("gotEndpoint: %v, wantEndpoint: %v", gotEndpoint, tt.wantEndpoint)
			}
		})
	}
}

func TestExtractFieldName(t *testing.T) {
	tests := map[string]struct {
		attributeName string
		wantFieldName string
	}{
		"simple_field": {
			attributeName: "Name",
			wantFieldName: "Name",
		},
		"id_field": {
			attributeName: "Id",
			wantFieldName: "Id",
		},
		"custom_field": {
			attributeName: "CustomField__c",
			wantFieldName: "CustomField__c",
		},
		"jsonpath_custom_field": {
			attributeName: "$.CustomField__c",
			wantFieldName: "CustomField__c",
		},
		"jsonpath_another_custom_field": {
			attributeName: "$.AnotherCustom__c",
			wantFieldName: "AnotherCustom__c",
		},
		"jsonpath_relationship": {
			attributeName: "$.Account.Name",
			wantFieldName: "Account.Name",
		},
		"jsonpath_deep_relationship": {
			attributeName: "$.Owner.Manager.Name",
			wantFieldName: "Owner.Manager.Name",
		},
		"jsonpath_array_wildcard": {
			attributeName: "$.Emails[*].value",
			wantFieldName: "Emails",
		},
		"jsonpath_array_index": {
			attributeName: "$.Contacts[0].Email",
			wantFieldName: "Contacts",
		},
		"jsonpath_array_filter": {
			attributeName: "$.Emails[?(@.primary==true)].value",
			wantFieldName: "Emails",
		},
		"jsonpath_simple": {
			attributeName: "$.Name",
			wantFieldName: "Name",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotFieldName := extractFieldName(tt.attributeName)

			if gotFieldName != tt.wantFieldName {
				t.Errorf("extractFieldName(%q) = %q, want %q", tt.attributeName, gotFieldName, tt.wantFieldName)
			}
		})
	}
}

func TestRemoveArraySyntax(t *testing.T) {
	tests := map[string]struct {
		fieldName     string
		wantFieldName string
	}{
		"no_array": {
			fieldName:     "CustomField__c",
			wantFieldName: "CustomField__c",
		},
		"array_wildcard": {
			fieldName:     "Emails[*]",
			wantFieldName: "Emails",
		},
		"array_index": {
			fieldName:     "Contacts[0]",
			wantFieldName: "Contacts",
		},
		"array_filter": {
			fieldName:     "Emails[?(@.primary==true)]",
			wantFieldName: "Emails",
		},
		"simple_field": {
			fieldName:     "Name",
			wantFieldName: "Name",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotFieldName := removeArraySyntax(tt.fieldName)

			if gotFieldName != tt.wantFieldName {
				t.Errorf("removeArraySyntax(%q) = %q, want %q", tt.fieldName, gotFieldName, tt.wantFieldName)
			}
		})
	}
}

func TestConstructEndpointWithJSONPath(t *testing.T) {
	tests := map[string]struct {
		request      *Request
		wantEndpoint string
	}{
		"jsonpath_custom_fields": {
			request: &Request{
				BaseURL:          "https://test.salesforce.com",
				APIVersion:       "58.0",
				EntityExternalID: "CustomObject",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "Id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.CustomField__c",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.AnotherCustom__c",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: "https://test.salesforce.com/services/data/v58.0/query?q=SELECT+Id,CustomField__c," +
				"AnotherCustom__c+FROM+CustomObject+ORDER+BY+Id+ASC",
		},
		"jsonpath_relationships": {
			request: &Request{
				BaseURL:          "https://test.salesforce.com",
				APIVersion:       "58.0",
				EntityExternalID: "Contact",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "Id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.Account.Name",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.Owner.Email",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: "https://test.salesforce.com/services/data/v58.0/query?q=SELECT+Id,Account.Name," +
				"Owner.Email+FROM+Contact+ORDER+BY+Id+ASC",
		},
		"jsonpath_array_fields": {
			request: &Request{
				BaseURL:          "https://test.salesforce.com",
				APIVersion:       "58.0",
				EntityExternalID: "Contact",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "Id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.Emails[*].value",
						Type:       framework.AttributeTypeString,
						List:       true,
					},
				},
			},
			wantEndpoint: "https://test.salesforce.com/services/data/v58.0/query?q=SELECT+Id,Emails+" +
				"FROM+Contact+ORDER+BY+Id+ASC",
		},
		"mixed_syntax": {
			request: &Request{
				BaseURL:          "https://test.salesforce.com",
				APIVersion:       "58.0",
				EntityExternalID: "Contact",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "Id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "Name",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.CustomField__c",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.Account.Name",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: "https://test.salesforce.com/services/data/v58.0/query?q=SELECT+Id,Name," +
				"CustomField__c,Account.Name+FROM+Contact+ORDER+BY+Id+ASC",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotEndpoint := ConstructEndpoint(tt.request)

			if gotEndpoint != tt.wantEndpoint {
				t.Errorf("ConstructEndpoint() = %q, want %q", gotEndpoint, tt.wantEndpoint)
			}
		})
	}
}
