// Copyright 2025 SGNL.ai, Inc.
package servicenow

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
				BaseURL:          "https://test-instance.service-now.com",
				APIVersion:       "v2",
				EntityExternalID: "sys_user",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "sys_id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "!url_encoding!",
						Type:       framework.AttributeTypeString,
					},
				},
				PageSize: 100,
			},
			wantEndpoint: "https://test-instance.service-now.com/api/now/v2/table/sys_user" +
				"?sysparm_fields=sys_id,%21url_encoding%21&sysparm_exclude_reference_link=true" +
				"&sysparm_limit=100&sysparm_query=ORDERBYsys_id",
		},
		"simple_with_filter": {
			request: &Request{
				BaseURL:          "https://test-instance.service-now.com",
				APIVersion:       "v2",
				EntityExternalID: "sys_user",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "sys_id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "!url_encoding!",
						Type:       framework.AttributeTypeString,
					},
				},
				PageSize: 100,
				Filter:   testutil.GenPtr("nameLIKETest"),
			},
			wantEndpoint: "https://test-instance.service-now.com/api/now/v2/table/sys_user?sysparm_fields=sys_id," +
				"%21url_encoding%21&sysparm_exclude_reference_link=true&sysparm_limit=100" +
				"&sysparm_query=nameLIKETest%5EORDERBYsys_id",
		},
		"nil_request": {
			request:      nil,
			wantEndpoint: "",
		},
		"simple_with_cursor": {
			request: &Request{
				BaseURL:          "https://test-instance.service-now.com",
				APIVersion:       "v2",
				EntityExternalID: "sys_user",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "sys_id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "!url_encoding!",
						Type:       framework.AttributeTypeString,
					},
				},
				PageSize: 100,
				Cursor: testutil.GenPtr("https://test-instance.service-now.com/api/now/v2/table/customer_account" +
					"?sysparm_fields=sys_id,number,account_parent,parent,sys_created_on,primary" +
					"&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=ORDERBYsys_id" +
					"&sysparm_offset=4"),
			},
			wantEndpoint: "https://test-instance.service-now.com/api/now/v2/table/customer_account" +
				"?sysparm_fields=sys_id,number,account_parent,parent,sys_created_on,primary" +
				"&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=ORDERBYsys_id&sysparm_offset=4",
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
