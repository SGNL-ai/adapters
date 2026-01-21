// Copyright 2026 SGNL.ai, Inc.

package jira_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	jira_adapter "github.com/sgnl-ai/adapters/pkg/jira"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[jira_adapter.Config]
		wantErr *framework.Error
	}{
		"valid": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Issue,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jira_adapter.Config{
					IssuesJQLFilter: testutil.GenPtr("project=SGNL"),
				},
			},
			wantErr: nil,
		},
		"valid_with_empty_config": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Issue,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jira_adapter.Config{},
			},
			wantErr: nil,
		},
		"invalid_config_filter_empty": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Issue,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jira_adapter.Config{
					IssuesJQLFilter: testutil.GenPtr(""),
				},
			},
			wantErr: &framework.Error{
				Message: "Jira config is invalid: issuesJqlFilter cannot be an empty string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_config_filter_too_long": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Issue,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jira_adapter.Config{
					IssuesJQLFilter: testutil.GenPtr(strings.Repeat("a", 1025)),
				},
			},
			wantErr: &framework.Error{
				Message: "Jira config is invalid: issuesJqlFilter exceeds the 1024 character limit.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"valid_ql_query": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Object,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "globalId",
						},
					},
				},
				Config: &jira_adapter.Config{
					ObjectsQLQuery: testutil.GenPtr("objectType = Customer"),
				},
			},
			wantErr: nil,
		},
		"empty_ql_query": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Object,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jira_adapter.Config{
					ObjectsQLQuery: testutil.GenPtr(""),
				},
			},
			wantErr: &framework.Error{
				Message: "Jira config is invalid: objectsQlQuery cannot be an empty string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"ql_query_too_long": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Object,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jira_adapter.Config{
					ObjectsQLQuery: testutil.GenPtr(strings.Repeat("a", 1025)),
				},
			},
			wantErr: &framework.Error{
				Message: "Jira config is invalid: objectsQlQuery exceeds the 1024 character limit.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"valid_asset_base_url": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Object,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "globalId",
						},
					},
				},
				Config: &jira_adapter.Config{
					AssetBaseURL: testutil.GenPtr("https://api.atlassian.com/jsm/assets"),
				},
			},
			wantErr: nil,
		},
		"empty_asset_base_url": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Object,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "globalId",
						},
					},
				},
				Config: &jira_adapter.Config{
					AssetBaseURL: testutil.GenPtr(""),
				},
			},
			wantErr: &framework.Error{
				Message: `Jira config is invalid: assetBaseUrl is not a valid URL: parse "": empty url.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_asset_base_url": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Object,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "globalId",
						},
					},
				},
				Config: &jira_adapter.Config{
					AssetBaseURL: testutil.GenPtr("INVALID_URL"),
				},
			},
			wantErr: &framework.Error{
				Message: `Jira config is invalid: assetBaseUrl is not a valid URL: parse "INVALID_URL": invalid URI for request.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"missing_auth": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Issue,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "Jira auth is missing required basic credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_external_id": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "INVALID_EXTERNAL_ID",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "Jira entity external ID is invalid.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"missing_unique_id": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Issue,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "NOT_THE_UNIQUE_ID",
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "Jira requested entity attributes are missing a unique ID attribute.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"child_entities_allowed": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Issue,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "child",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "id",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		"ordered_must_be_false": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Issue,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Ordered: true,
			},
			wantErr: &framework.Error{
				Message: "Jira Ordered property must be false.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"page_size_too_big": {
			request: &framework.Request[jira_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jira_adapter.Issue,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				PageSize: 9999,
			},
			wantErr: &framework.Error{
				Message: "Jira provided page size (9999) exceeds the maximum (1000).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
	}

	adapter := &jira_adapter.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(context.TODO(), tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
