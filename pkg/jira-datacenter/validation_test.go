// Copyright 2025 SGNL.ai, Inc.
package jiradatacenter_test

import (
	"reflect"
	"strings"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	jiradatacenter_adapter "github.com/sgnl-ai/adapters/pkg/jira-datacenter"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[jiradatacenter_adapter.Config]
		wantErr *framework.Error
	}{
		"valid": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jiradatacenter_adapter.Config{
					IssuesJQLFilter: testutil.GenPtr("project=SGNL"),
					Groups:          []string{"jira-administrators", "jira-users"},
				},
			},
			wantErr: nil,
		},
		"valid_with_empty_config": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jiradatacenter_adapter.Config{},
			},
			wantErr: nil,
		},
		"invalid_config_filter_empty": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jiradatacenter_adapter.Config{
					IssuesJQLFilter: testutil.GenPtr(""),
				},
			},
			wantErr: &framework.Error{
				Message: "Jira config is invalid: issuesJqlFilter cannot be an empty string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_config_empty_group_name": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jiradatacenter_adapter.Config{
					Groups: []string{"jira-administrators", ""},
				},
			},
			wantErr: &framework.Error{
				Message: "Jira config is invalid: group at index '1' cannot be an empty string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_config_filter_too_long": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jiradatacenter_adapter.Config{
					IssuesJQLFilter: testutil.GenPtr(strings.Repeat("a", 1025)),
				},
			},
			wantErr: &framework.Error{
				Message: "Jira config is invalid: issuesJqlFilter exceeds the 1024 character limit.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_config_group_name_too_long": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jiradatacenter_adapter.Config{
					Groups: []string{strings.Repeat("a", 256)},
				},
			},
			wantErr: &framework.Error{
				Message: "Jira config is invalid: group name at index '0' exceeds the 255 character limit.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"missing_auth": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "Request to Jira is missing Basic Auth or Personal Access Token (PAT) credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_external_id": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
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
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
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
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
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
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
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
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
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
		"page_size_too_big_for_users": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.UserExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "key",
						},
					},
				},
				PageSize: 51,
			},
			wantErr: &framework.Error{
				Message: "User or group member page size (51) exceeds allowed maximum (50).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"page_size_too_big_for_group_members": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.GroupMemberExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				PageSize: 51,
			},
			wantErr: &framework.Error{
				Message: "User or group member page size (51) exceeds allowed maximum (50).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"invalid_api_version": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jiradatacenter_adapter.Config{
					APIVersion: "1",
				},
			},
			wantErr: &framework.Error{
				Message: "Jira config is invalid: apiVersion must be either '2' or 'latest', got '1'.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"valid_api_version_2": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jiradatacenter_adapter.Config{
					APIVersion: "2",
				},
			},
			wantErr: nil,
		},
		"valid_api_version_latest": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &jiradatacenter_adapter.Config{
					APIVersion: "latest",
				},
			},
			wantErr: nil,
		},
		"invalid_missing_bearer_prefix": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "testtoken", // missing Bearer prefix
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: `Provided auth token is missing required "Bearer " prefix.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"valid_with_bearer_prefix": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer testtoken",
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.IssueExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
			},
			wantErr: nil,
		},
		"valid_with_groups_max_results": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.GroupExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
						},
					},
				},
				Config: &jiradatacenter_adapter.Config{
					GroupsMaxResults: testutil.GenPtr[int64](50),
				},
			},
			wantErr: nil,
		},
		"invalid_groups_max_results_zero": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.GroupExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
						},
					},
				},
				Config: &jiradatacenter_adapter.Config{
					GroupsMaxResults: testutil.GenPtr[int64](0),
				},
			},
			wantErr: &framework.Error{
				Message: "Jira config is invalid: groupsMaxResults must be greater than 0.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_groups_max_results_too_large": {
			request: &framework.Request[jiradatacenter_adapter.Config]{
				Address: "https://example.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "username",
						Password: "password",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: jiradatacenter_adapter.GroupExternalID,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
						},
					},
				},
				Config: &jiradatacenter_adapter.Config{
					GroupsMaxResults: testutil.GenPtr[int64](1001),
				},
			},
			wantErr: &framework.Error{
				Message: "Jira config is invalid: groupsMaxResults cannot exceed 1000.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
	}

	adapter := &jiradatacenter_adapter.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
