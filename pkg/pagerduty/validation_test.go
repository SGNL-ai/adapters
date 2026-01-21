// Copyright 2026 SGNL.ai, Inc.

package pagerduty_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	pagerduty_adapter "github.com/sgnl-ai/adapters/pkg/pagerduty"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[pagerduty_adapter.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[pagerduty_adapter.Config]{
				Address: "api.pagerduty.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Token token=y_NbAkKc66ryYTWUXYEu",
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
			},
			wantErr: nil,
		},
		"valid_request_empty_config": {
			request: &framework.Request[pagerduty_adapter.Config]{
				Address: "api.pagerduty.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Token token=y_NbAkKc66ryYTWUXYEu",
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &pagerduty_adapter.Config{},
			},
			wantErr: nil,
		},
		"valid_request_valid_config": {
			request: &framework.Request[pagerduty_adapter.Config]{
				Address: "api.pagerduty.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Token token=y_NbAkKc66ryYTWUXYEu",
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &pagerduty_adapter.Config{
					AdditionalQueryParameters: map[string]map[string]any{
						"users": {
							"include[]": []any{"1234", "5678"},
							"query":     "Random Name",
						},
						"services": {
							"team_ids[]": []any{"1234", "5678"},
							"time_zone":  "UTC",
						},
					},
				},
			},
			wantErr: nil,
		},
		"invalid_value_in_additionalQueryParameters": {
			request: &framework.Request[pagerduty_adapter.Config]{
				Address: "api.pagerduty.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Token token=y_NbAkKc66ryYTWUXYEu",
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &pagerduty_adapter.Config{
					AdditionalQueryParameters: map[string]map[string]any{
						"users": {
							"include[]": 10, // Can only be string or []string.
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "PagerDuty config is invalid: additionalQueryParameters[users][include[]] " +
					"is neither a string nor a list of strings.",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"empty_string_param_value_in_additionalQueryParameters": {
			request: &framework.Request[pagerduty_adapter.Config]{
				Address: "api.pagerduty.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Token token=y_NbAkKc66ryYTWUXYEu",
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &pagerduty_adapter.Config{
					AdditionalQueryParameters: map[string]map[string]any{
						"users": {
							"include[]": "",
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "PagerDuty config is invalid: additionalQueryParameters[users][include[]] is an empty string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"empty_list_param_value_in_additionalQueryParameters": {
			request: &framework.Request[pagerduty_adapter.Config]{
				Address: "api.pagerduty.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Token token=y_NbAkKc66ryYTWUXYEu",
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &pagerduty_adapter.Config{
					AdditionalQueryParameters: map[string]map[string]any{
						"users": {
							"include[]": []any{},
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "PagerDuty config is invalid: additionalQueryParameters[users][include[]] is an empty list.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_param_value_in_additionalQueryParameters": {
			request: &framework.Request[pagerduty_adapter.Config]{
				Address: "api.pagerduty.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Token token=y_NbAkKc66ryYTWUXYEu",
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Config: &pagerduty_adapter.Config{
					AdditionalQueryParameters: map[string]map[string]any{
						"users": {
							"include[]": []any{1, "1", true},
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "PagerDuty config is invalid: additionalQueryParameters[users][include[]][0] is not a string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"missing_httpauth": {
			request: &framework.Request[pagerduty_adapter.Config]{
				Address: "api.pagerduty.com",
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "PagerDuty auth is missing required token.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"httpauth_malformed": {
			request: &framework.Request[pagerduty_adapter.Config]{
				Address: "api.pagerduty.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Not prefixed with Token",
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: `PagerDuty auth is missing required "Token token=" or "Bearer " prefix.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"address_invalid": {
			request: &framework.Request[pagerduty_adapter.Config]{
				Address: "not api.pagerduty.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Token token=y_NbAkKc66ryYTWUXYEu",
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "Invalid PagerDuty address. Must be api.pagerduty.com.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"missing_unique_id": {
			request: &framework.Request[pagerduty_adapter.Config]{
				Address: "api.pagerduty.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Token token=y_NbAkKc66ryYTWUXYEu",
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "NOT_THE_UNIQUE_ID",
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "PagerDuty requested entity attributes are missing a unique ID attribute: id.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"no_child_entities_allowed": {
			request: &framework.Request[pagerduty_adapter.Config]{
				Address: "api.pagerduty.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Token token=y_NbAkKc66ryYTWUXYEu",
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
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
			wantErr: &framework.Error{
				Message: "PagerDuty requested entity does not support child entities.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"ordered_must_be_false": {
			request: &framework.Request[pagerduty_adapter.Config]{
				Address: "api.pagerduty.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Token token=y_NbAkKc66ryYTWUXYEu",
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				Ordered: true,
			},
			wantErr: &framework.Error{
				Message: "PagerDuty Ordered property must be false.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"page_size_too_big": {
			request: &framework.Request[pagerduty_adapter.Config]{
				Address: "api.pagerduty.com",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Token token=y_NbAkKc66ryYTWUXYEu",
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
						},
					},
				},
				PageSize: 9999,
			},
			wantErr: &framework.Error{
				Message: "PagerDuty provided page size (9999) exceeds the maximum (100).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
	}

	adapter := &pagerduty_adapter.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
