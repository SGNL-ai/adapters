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
