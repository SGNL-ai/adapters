// Copyright 2026 SGNL.ai, Inc.

// nolint: lll, goconst
package servicenow_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	servicenow_adapter "github.com/sgnl-ai/adapters/pkg/servicenow"
	"github.com/stretchr/testify/assert"
)

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := servicenow_adapter.NewAdapter(&servicenow_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx          context.Context
		request      *framework.Request[servicenow_adapter.Config]
		wantResponse framework.Response
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "email",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "sys_created_on",
							Type:       framework.AttributeTypeDateTime,
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
							"sys_id":         "9a826bf03710200044e0bfc8bcbe5dd1",
							"sys_created_on": time.Date(2012, 2, 18, 3, 4, 51, 0, time.UTC),
							"email":          "freeman.soula@example.com",
						},
						{
							"sys_id":         "a2826bf03710200044e0bfc8bcbe5ddb",
							"sys_created_on": time.Date(2012, 2, 18, 3, 4, 52, 0, time.UTC),
							"email":          "junior.wadlinger@example.com",
						},
						{
							"sys_id":         "aa826bf03710200044e0bfc8bcbe5ddf",
							"sys_created_on": time.Date(2012, 2, 18, 3, 4, 52, 0, time.UTC),
							"email":          "curt.menedez@example.com",
						},
					},
					NextCursor: "https://localhost/api/now/v2/table/sys_user?sysparm_fields=sys_id,manager,email,sys_created_on,active&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=ORDERBYsys_id&sysparm_offset=3",
				},
			},
		},
		"valid_request_with_http_auth": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{HTTPAuthorization: "Bearer testtoken"},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "email",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "sys_created_on",
							Type:       framework.AttributeTypeDateTime,
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
							"sys_id":         "9a826bf03710200044e0bfc8bcbe5dd1",
							"sys_created_on": time.Date(2012, 2, 18, 3, 4, 51, 0, time.UTC),
							"email":          "freeman.soula@example.com",
						},
						{
							"sys_id":         "a2826bf03710200044e0bfc8bcbe5ddb",
							"sys_created_on": time.Date(2012, 2, 18, 3, 4, 52, 0, time.UTC),
							"email":          "junior.wadlinger@example.com",
						},
						{
							"sys_id":         "aa826bf03710200044e0bfc8bcbe5ddf",
							"sys_created_on": time.Date(2012, 2, 18, 3, 4, 52, 0, time.UTC),
							"email":          "curt.menedez@example.com",
						},
					},
					NextCursor: "https://localhost/api/now/v2/table/sys_user?sysparm_fields=sys_id,manager,email,sys_created_on,active&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=ORDERBYsys_id&sysparm_offset=3",
				},
			},
		},
		"valid_request_no_https_prefix": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "email",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "sys_created_on",
							Type:       framework.AttributeTypeDateTime,
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
							"sys_id":         "9a826bf03710200044e0bfc8bcbe5dd1",
							"sys_created_on": time.Date(2012, 2, 18, 3, 4, 51, 0, time.UTC),
							"email":          "freeman.soula@example.com",
						},
						{
							"sys_id":         "a2826bf03710200044e0bfc8bcbe5ddb",
							"sys_created_on": time.Date(2012, 2, 18, 3, 4, 52, 0, time.UTC),
							"email":          "junior.wadlinger@example.com",
						},
						{
							"sys_id":         "aa826bf03710200044e0bfc8bcbe5ddf",
							"sys_created_on": time.Date(2012, 2, 18, 3, 4, 52, 0, time.UTC),
							"email":          "curt.menedez@example.com",
						},
					},
					NextCursor: "https://localhost/api/now/v2/table/sys_user?sysparm_fields=sys_id,manager,email,sys_created_on,active&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=ORDERBYsys_id&sysparm_offset=3",
				},
			},
		},
		"invalid_request_invalid_api_version": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v1",
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "email",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "sys_created_on",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
					},
				},
				Ordered:  true,
				PageSize: 200,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Servicenow config is invalid: apiVersion is not supported: v1.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"invalid_request_http_prefix": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: "http://" + strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "email",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "sys_created_on",
							Type:       framework.AttributeTypeDateTime,
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
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					Filters: map[string]string{
						"sn_customerservice_case": "active=true",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "sn_customerservice_case",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "$.assigned_to.value",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "sys_created_on",
							Type:       framework.AttributeTypeDateTime,
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
							"sys_id":              "f6911038530360100999ddeeff7case1",
							"sys_created_on":      time.Date(2022, 5, 26, 21, 59, 3, 0, time.UTC),
							"$.assigned_to.value": "9a826bf03710200044e0bfc8bcbe5dd1",
						},
						{
							"sys_id":              "f6911038530360100999ddeeff7case4",
							"sys_created_on":      time.Date(2022, 5, 26, 21, 59, 3, 0, time.UTC),
							"$.assigned_to.value": "a2826bf03710200044e0bfc8bcbe5ddb",
						},
					},
					NextCursor: "https://localhost/api/now/v2/table/sn_customerservice_case?sysparm_fields=sys_id,case,parent,assigned_to,account,description,sys_created_on,active&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=active%3Dtrue%5EORDERBYsys_id&sysparm_offset=4",
				},
			},
		},
		"valid_request_with_cursor": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "email",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "sys_created_on",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
					},
				},
				Ordered:  true,
				PageSize: 200,
				Cursor:   server.URL + "/api/now/v2/table/sys_user?sysparm_fields=sys_id,manager,email,sys_created_on,active&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=ORDERBYsys_id&sysparm_offset=3",
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"sys_id":         "cf1ec0b4530360100999ddeeff7b129f",
							"sys_created_on": time.Date(2012, 2, 18, 3, 4, 51, 0, time.UTC),
							"email":          "john.doe@example.com",
						},
					},
				},
			},
		},
		"invalid_request_invalid_url": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL + "/invalid",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "email",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "sys_created_on",
							Type:       framework.AttributeTypeDateTime,
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
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
				},
				Entity: framework.EntityConfig{
					ExternalId: "sys_user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "email",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "sys_created_on",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
					},
				},
				Ordered:  true,
				PageSize: 200,
				Cursor:   server.URL + "/api/now/v2/table/sys_user?sysparm_fields=sys_id,manager,email,sys_created_on,active&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=ORDERBYsys_id&sysparm_offset=4",
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to convert datasource response objects: attribute sys_created_on cannot be parsed " +
						"into a date-time value: failed to parse date-time value: 2021/01/01 00:00:00.000Z.",
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
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

func TestAdapterGetPageWithAdvancedFilters(t *testing.T) {
	fixedAddress := "127.0.0.1:8443"

	listener, err := net.Listen("tcp", fixedAddress)
	if err != nil {
		panic(err)
	}

	server := &httptest.Server{
		Listener: listener,
		Config:   &http.Server{},
	}
	server.StartTLS()
	defer server.Close()

	server.Config.Handler = TestServerAdvancedFiltersHandler

	adapter := servicenow_adapter.NewAdapter(&servicenow_adapter.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx          context.Context
		request      *framework.Request[servicenow_adapter.Config]
		wantResponse framework.Response
	}{
		"valid_request_implicit_filters": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.Group,
									ScopeEntityFilter: "active=true",
									Members: []servicenow_adapter.MemberFilter{
										{
											MemberEntity:       servicenow_adapter.User,
											MemberEntityFilter: "user.active=true",
										},
									},
								},
							},
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.User,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "email",
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
							"sys_id": "9a826bf03710200044e0bfc8bcbe5dd1",
							"email":  "user2@example.com",
						},
					},
				},
			},
		},
		"valid_request_implicit_filters_scope_entity_only_page_1": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.Group,
									ScopeEntityFilter: "sys_idINc38f00f4530360100999ddeexxgroup1,c38f00f4530360100999ddeexxgroup2",
								},
							},
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.Group,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "defaultAssignee",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  true,
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"sys_id":          "c38f00f4530360100999ddeexxgroup1",
							"defaultAssignee": "Tom",
						},
					},
					NextCursor: "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6eyJjb2xsZWN0aW9uQ3Vyc29yIjoiaHR0cHM6Ly8xMjcuMC4wLjE6ODQ0My9hcGkvbm93L3YyL3RhYmxlL3N5c191c2VyX2dyb3VwP3N5c3Bhcm1fZmllbGRzPXN5c19pZCxkZWZhdWx0X2Fzc2lnbmVlXHUwMDI2c3lzcGFybV9leGNsdWRlX3JlZmVyZW5jZV9saW5rPXRydWVcdTAwMjZzeXNwYXJtX2xpbWl0PTFcdTAwMjZzeXNwYXJtX3F1ZXJ5PXN5c19pZElOYzM4ZjAwZjQ1MzAzNjAxMDA5OTlkZGVleHhncm91cDEsYzM4ZjAwZjQ1MzAzNjAxMDA5OTlkZGVleHhncm91cDJeT1JERVJCWXN5c19pZFx1MDAyNnN5c3Bhcm1fb2Zmc2V0PTEifX19",
				},
			},
		},
		"valid_request_implicit_filters_scope_entity_only_page_2": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.Group,
									ScopeEntityFilter: "sys_idINc38f00f4530360100999ddeexxgroup1,c38f00f4530360100999ddeexxgroup2",
								},
							},
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.Group,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "defaultAssignee",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Cursor:   "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6eyJjb2xsZWN0aW9uQ3Vyc29yIjoiaHR0cHM6Ly8xMjcuMC4wLjE6ODQ0My9hcGkvbm93L3YyL3RhYmxlL3N5c191c2VyX2dyb3VwP3N5c3Bhcm1fZmllbGRzPXN5c19pZCxkZWZhdWx0X2Fzc2lnbmVlXHUwMDI2c3lzcGFybV9leGNsdWRlX3JlZmVyZW5jZV9saW5rPXRydWVcdTAwMjZzeXNwYXJtX2xpbWl0PTFcdTAwMjZzeXNwYXJtX3F1ZXJ5PXN5c19pZElOYzM4ZjAwZjQ1MzAzNjAxMDA5OTlkZGVleHhncm91cDEsYzM4ZjAwZjQ1MzAzNjAxMDA5OTlkZGVleHhncm91cDJeT1JERVJCWXN5c19pZFx1MDAyNnN5c3Bhcm1fb2Zmc2V0PTEifX19",
				Ordered:  true,
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"sys_id":          "c38f00f4530360100999ddeexxgroup2",
							"defaultAssignee": "John",
						},
					},
				},
			},
		},
		"valid_request_implicit_filters_scope_entity_and_member_entity_page_1": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.Group,
									ScopeEntityFilter: "sys_idINc38f00f4530360100999ddeexxgroup1,c38f00f4530360100999ddeexxgroup2",
									Members: []servicenow_adapter.MemberFilter{
										{
											MemberEntity:       servicenow_adapter.User,
											MemberEntityFilter: "user.active=true",
										},
									},
								},
							},
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.User,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "active",
							Type:       framework.AttributeTypeBool,
						},
					},
				},
				Ordered:  true,
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"sys_id": "9a826bf03710200044e0bfc8bcbe5dd1",
							"active": true,
						},
					},
					NextCursor: "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6eyJjdXJzb3IiOiJodHRwczovLzEyNy4wLjAuMTo4NDQzL2FwaS9ub3cvdjIvdGFibGUvc3lzX3VzZXJfZ3JtZW1iZXI/c3lzcGFybV9maWVsZHM9c3lzX2lkLHVzZXIuc3lzX2lkXHUwMDI2c3lzcGFybV9leGNsdWRlX3JlZmVyZW5jZV9saW5rPXRydWVcdTAwMjZzeXNwYXJtX2xpbWl0PTFcdTAwMjZzeXNwYXJtX3F1ZXJ5PWdyb3VwSU5jMzhmMDBmNDUzMDM2MDEwMDk5OWRkZWV4eGdyb3VwMV51c2VyLmFjdGl2ZT10cnVlXk9SREVSQllzeXNfaWRcdTAwMjZzeXNwYXJtX29mZnNldD0xIn19fQ==",
				},
			},
		},
		"valid_request_implicit_filters_scope_entity_and_member_entity_page_2": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.Group,
									ScopeEntityFilter: "sys_idINc38f00f4530360100999ddeexxgroup1,c38f00f4530360100999ddeexxgroup2",
									Members: []servicenow_adapter.MemberFilter{
										{
											MemberEntity:       servicenow_adapter.User,
											MemberEntityFilter: "user.active=true",
										},
									},
								},
							},
						},
					},
				},
				Cursor: "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6eyJjdXJzb3IiOiJodHRwczovLzEyNy4wLjAuMTo4NDQzL2FwaS9ub3cvdjIvdGFibGUvc3lzX3VzZXJfZ3JtZW1iZXI/c3lzcGFybV9maWVsZHM9c3lzX2lkLHVzZXIuc3lzX2lkXHUwMDI2c3lzcGFybV9leGNsdWRlX3JlZmVyZW5jZV9saW5rPXRydWVcdTAwMjZzeXNwYXJtX2xpbWl0PTFcdTAwMjZzeXNwYXJtX3F1ZXJ5PWdyb3VwSU5jMzhmMDBmNDUzMDM2MDEwMDk5OWRkZWV4eGdyb3VwMV51c2VyLmFjdGl2ZT10cnVlXk9SREVSQllzeXNfaWRcdTAwMjZzeXNwYXJtX29mZnNldD0xIn19fQ==",
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.User,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "active",
							Type:       framework.AttributeTypeBool,
						},
					},
				},
				Ordered:  true,
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"sys_id": "a2826bf03710200044e0bfc8bcbe5ddb",
							"active": true,
						},
					},
					NextCursor: "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6eyJjb2xsZWN0aW9uQ3Vyc29yIjoiaHR0cHM6Ly8xMjcuMC4wLjE6ODQ0My9hcGkvbm93L3YyL3RhYmxlL3N5c191c2VyX2dyb3VwP3N5c3Bhcm1fZmllbGRzPXN5c19pZCxkZWZhdWx0X2Fzc2lnbmVlXHUwMDI2c3lzcGFybV9leGNsdWRlX3JlZmVyZW5jZV9saW5rPXRydWVcdTAwMjZzeXNwYXJtX2xpbWl0PTFcdTAwMjZzeXNwYXJtX3F1ZXJ5PXN5c19pZElOYzM4ZjAwZjQ1MzAzNjAxMDA5OTlkZGVleHhncm91cDEsYzM4ZjAwZjQ1MzAzNjAxMDA5OTlkZGVleHhncm91cDJeT1JERVJCWXN5c19pZFx1MDAyNnN5c3Bhcm1fb2Zmc2V0PTEifX19",
				},
			},
		},
		"valid_request_implicit_filters_scope_entity_and_member_entity_page_3": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.Group,
									ScopeEntityFilter: "sys_idINc38f00f4530360100999ddeexxgroup1,c38f00f4530360100999ddeexxgroup2",
									Members: []servicenow_adapter.MemberFilter{
										{
											MemberEntity:       servicenow_adapter.User,
											MemberEntityFilter: "user.active=true",
										},
									},
								},
							},
						},
					},
				},
				Cursor: "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6eyJjb2xsZWN0aW9uQ3Vyc29yIjoiaHR0cHM6Ly8xMjcuMC4wLjE6ODQ0My9hcGkvbm93L3YyL3RhYmxlL3N5c191c2VyX2dyb3VwP3N5c3Bhcm1fZmllbGRzPXN5c19pZCxkZWZhdWx0X2Fzc2lnbmVlXHUwMDI2c3lzcGFybV9leGNsdWRlX3JlZmVyZW5jZV9saW5rPXRydWVcdTAwMjZzeXNwYXJtX2xpbWl0PTFcdTAwMjZzeXNwYXJtX3F1ZXJ5PXN5c19pZElOYzM4ZjAwZjQ1MzAzNjAxMDA5OTlkZGVleHhncm91cDEsYzM4ZjAwZjQ1MzAzNjAxMDA5OTlkZGVleHhncm91cDJeT1JERVJCWXN5c19pZFx1MDAyNnN5c3Bhcm1fb2Zmc2V0PTEifX19",
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.User,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "active",
							Type:       framework.AttributeTypeBool,
						},
					},
				},
				Ordered:  true,
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"sys_id": "aa826bf03710200044e0bfc8bcbe5ddf",
							"active": true,
						},
					},
					NextCursor: "",
				},
			},
		},
		"invalid_request_implicit_filters_scope_entity_is_not_group": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.User, // This must be Group.
									ScopeEntityFilter: "active=true",
									Members: []servicenow_adapter.MemberFilter{
										{
											MemberEntity:       servicenow_adapter.User,
											MemberEntityFilter: "user.active=true",
										},
									},
								},
							},
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.User,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
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
					Message: "sys_user is not a supported scope for the current entity: sys_user.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"invalid_request_implicit_filters_invalid_cursor": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.Group,
									ScopeEntityFilter: "active=true",
								},
							},
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.Group,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Cursor:   "invalid_cursor",
				Ordered:  true,
				PageSize: 200,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to decode base64 cursor: illegal base64 data at input byte 7.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
		"invalid_request_implicit_filters_valid_cursor_but_no_implicit_filter": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.Group,
									ScopeEntityFilter: "active=true",
								},
							},
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.Group,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Cursor:   "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjowfX0=",
				Ordered:  true,
				PageSize: 200,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Implicit filter cursor is unexpectedly nil for entity: sys_user_group.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"valid_request_related_filters": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.Group,
									ScopeEntityFilter: "active=true",
									Members: []servicenow_adapter.MemberFilter{
										{
											MemberEntity:       servicenow_adapter.User,
											MemberEntityFilter: "user.active=true",
											RelatedEntities: []servicenow_adapter.RelatedEntityFilter{
												{
													RelatedEntity:       servicenow_adapter.Case,
													RelatedEntityFilter: "assigned_toIN{$.sys_user.sys_id}",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.Case,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "assigned_to",
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
							"sys_id":      "f6911038530360100999ddeeff7case1",
							"assigned_to": "9a826bf03710200044e0bfc8bcbe5dd1",
						},
					},
				},
			},
		},
		"valid_request_related_filters_scope_entity_and_member_entity_page_1": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.Group,
									ScopeEntityFilter: "sys_idINc38f00f4530360100999ddeexxgroup1,c38f00f4530360100999ddeexxgroup2",
									RelatedEntities: []servicenow_adapter.RelatedEntityFilter{
										{
											RelatedEntity:       servicenow_adapter.Case,
											RelatedEntityFilter: "assignment_groupIN{$.sys_user_group.sys_id}",
										},
									},
								},
							},
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.Case,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "assignment_group",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Ordered:  true,
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"sys_id":           "f6911038530360100999ddeeff7case1",
							"assignment_group": "c38f00f4530360100999ddeexxgroup1",
						},
					},
					NextCursor: "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjowLCJlbnRpdHlDdXJzb3IiOiJodHRwczovLzEyNy4wLjAuMTo4NDQzL2FwaS9ub3cvdjIvdGFibGUvc25fY3VzdG9tZXJzZXJ2aWNlX2Nhc2U/c3lzcGFybV9maWVsZHM9c3lzX2lkLGFzc2lnbm1lbnRfZ3JvdXBcdTAwMjZzeXNwYXJtX2V4Y2x1ZGVfcmVmZXJlbmNlX2xpbms9dHJ1ZVx1MDAyNnN5c3Bhcm1fbGltaXQ9MVx1MDAyNnN5c3Bhcm1fcXVlcnk9YXNzaWdubWVudF9ncm91cElOYzM4ZjAwZjQ1MzAzNjAxMDA5OTlkZGVleHhncm91cDFeT1JERVJCWXN5c19pZFx1MDAyNnN5c3Bhcm1fb2Zmc2V0PTEifX0=",
				},
			},
		},
		"valid_request_related_filters_scope_entity_and_member_entity_page_2": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.Group,
									ScopeEntityFilter: "sys_idINc38f00f4530360100999ddeexxgroup1,c38f00f4530360100999ddeexxgroup2",
									RelatedEntities: []servicenow_adapter.RelatedEntityFilter{
										{
											RelatedEntity:       servicenow_adapter.Case,
											RelatedEntityFilter: "assignment_groupIN{$.sys_user_group.sys_id}",
										},
									},
								},
							},
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.Case,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "assignment_group",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Cursor:   "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjowLCJlbnRpdHlDdXJzb3IiOiJodHRwczovLzEyNy4wLjAuMTo4NDQzL2FwaS9ub3cvdjIvdGFibGUvc25fY3VzdG9tZXJzZXJ2aWNlX2Nhc2U/c3lzcGFybV9maWVsZHM9c3lzX2lkLGFzc2lnbm1lbnRfZ3JvdXBcdTAwMjZzeXNwYXJtX2V4Y2x1ZGVfcmVmZXJlbmNlX2xpbms9dHJ1ZVx1MDAyNnN5c3Bhcm1fbGltaXQ9MVx1MDAyNnN5c3Bhcm1fcXVlcnk9YXNzaWdubWVudF9ncm91cElOYzM4ZjAwZjQ1MzAzNjAxMDA5OTlkZGVleHhncm91cDFeT1JERVJCWXN5c19pZFx1MDAyNnN5c3Bhcm1fb2Zmc2V0PTEifX0=",
				Ordered:  true,
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"sys_id":           "f6911038530360100999ddeeff7case2",
							"assignment_group": "c38f00f4530360100999ddeexxgroup1",
						},
					},
					NextCursor: "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjowLCJyZWxhdGVkRW50aXR5Q3Vyc29yIjp7ImNvbGxlY3Rpb25DdXJzb3IiOiJodHRwczovLzEyNy4wLjAuMTo4NDQzL2FwaS9ub3cvdjIvdGFibGUvc3lzX3VzZXJfZ3JvdXA/c3lzcGFybV9maWVsZHM9c3lzX2lkLGRlZmF1bHRfYXNzaWduZWVcdTAwMjZzeXNwYXJtX2V4Y2x1ZGVfcmVmZXJlbmNlX2xpbms9dHJ1ZVx1MDAyNnN5c3Bhcm1fbGltaXQ9MVx1MDAyNnN5c3Bhcm1fcXVlcnk9c3lzX2lkSU5jMzhmMDBmNDUzMDM2MDEwMDk5OWRkZWV4eGdyb3VwMSxjMzhmMDBmNDUzMDM2MDEwMDk5OWRkZWV4eGdyb3VwMl5PUkRFUkJZc3lzX2lkXHUwMDI2c3lzcGFybV9vZmZzZXQ9MSJ9fX0=",
				},
			},
		},
		"valid_request_related_filters_scope_entity_and_member_entity_page_3": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.Group,
									ScopeEntityFilter: "sys_idINc38f00f4530360100999ddeexxgroup1,c38f00f4530360100999ddeexxgroup2",
									RelatedEntities: []servicenow_adapter.RelatedEntityFilter{
										{
											RelatedEntity:       servicenow_adapter.Case,
											RelatedEntityFilter: "assignment_groupIN{$.sys_user_group.sys_id}",
										},
									},
								},
							},
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.Case,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "assignment_group",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Cursor:   "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjowLCJlbnRpdHlDdXJzb3IiOm51bGwsInJlbGF0ZWRFbnRpdHlDdXJzb3IiOnsiY29sbGVjdGlvbkN1cnNvciI6Imh0dHBzOi8vMTI3LjAuMC4xOjg0NDMvYXBpL25vdy92Mi90YWJsZS9zeXNfdXNlcl9ncm91cD9zeXNwYXJtX2ZpZWxkcz1zeXNfaWQsZGVmYXVsdF9hc3NpZ25lZVx1MDAyNnN5c3Bhcm1fZXhjbHVkZV9yZWZlcmVuY2VfbGluaz10cnVlXHUwMDI2c3lzcGFybV9saW1pdD0xXHUwMDI2c3lzcGFybV9xdWVyeT1zeXNfaWRJTmMzOGYwMGY0NTMwMzYwMTAwOTk5ZGRlZXh4Z3JvdXAxLGMzOGYwMGY0NTMwMzYwMTAwOTk5ZGRlZXh4Z3JvdXAyXk9SREVSQllzeXNfaWRcdTAwMjZzeXNwYXJtX29mZnNldD0xIn19fQ==",
				Ordered:  true,
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"sys_id":           "f6911038530360100999ddeeff7case3",
							"assignment_group": "c38f00f4530360100999ddeexxgroup2",
						},
					},
				},
			},
		},
		"invalid_request_related_filters_invalid_cursor": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.Group,
									ScopeEntityFilter: "active=true",
									RelatedEntities: []servicenow_adapter.RelatedEntityFilter{
										{
											RelatedEntity:       servicenow_adapter.Case,
											RelatedEntityFilter: "assignment_groupIN{$.sys_user_group.sys_id}",
										},
									},
								},
							},
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.Case,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Cursor:   "invalid_cursor",
				Ordered:  true,
				PageSize: 200,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to decode base64 cursor: illegal base64 data at input byte 7.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				},
			},
		},
		"invalid_request_related_filters_valid_cursor_but_no_related_filter": {
			ctx: context.Background(),
			request: &framework.Request[servicenow_adapter.Config]{
				Address: server.URL,
				Auth:    &framework.DatasourceAuthCredentials{Basic: &framework.BasicAuthCredentials{Username: "username", Password: "password"}},
				Config: &servicenow_adapter.Config{
					APIVersion: "v2",
					AdvancedFilters: &servicenow_adapter.AdvancedFilters{
						ScopedObjects: map[string][]servicenow_adapter.EntityFilter{
							servicenow_adapter.Group: {
								{
									ScopeEntity:       servicenow_adapter.Group,
									ScopeEntityFilter: "active=true",
									RelatedEntities: []servicenow_adapter.RelatedEntityFilter{
										{
											RelatedEntity:       servicenow_adapter.Case,
											RelatedEntityFilter: "assignment_groupIN{$.sys_user_group.sys_id}",
										},
									},
								},
							},
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: servicenow_adapter.Case,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "sys_id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				Cursor:   "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjB9fQ==",
				Ordered:  true,
				PageSize: 200,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Related filter cursor is unexpectedly nil for entity: sn_customerservice_case.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			assert.Equal(t, tt.wantResponse, gotResponse)
		})
	}
}
