// Copyright 2026 SGNL.ai, Inc.

package victorops_test

import (
	"context"
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	victorops_adapter "github.com/sgnl-ai/adapters/pkg/victorops"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[victorops_adapter.Config]
		wantErr *framework.Error
	}{
		"valid_incident": {
			request: &framework.Request[victorops_adapter.Config]{
				Address: "https://api.victorops.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "api-id",
						Password: "api-key",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
						},
					},
				},
			},
			wantErr: nil,
		},
		"valid_user": {
			request: &framework.Request[victorops_adapter.Config]{
				Address: "https://api.victorops.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "api-id",
						Password: "api-key",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.User,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "username",
						},
					},
				},
			},
			wantErr: nil,
		},
		"valid_with_query_parameters": {
			request: &framework.Request[victorops_adapter.Config]{
				Address: "https://api.victorops.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "api-id",
						Password: "api-key",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
						},
					},
				},
				Config: &victorops_adapter.Config{
					QueryParameters: map[string]string{
						"IncidentReport": "currentPhase=RESOLVED&startedAfter=2024-01-01T00:00Z",
					},
				},
			},
			wantErr: nil,
		},
		"invalid_query_parameters_empty_value": {
			request: &framework.Request[victorops_adapter.Config]{
				Address: "https://api.victorops.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "api-id",
						Password: "api-key",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
						},
					},
				},
				Config: &victorops_adapter.Config{
					QueryParameters: map[string]string{
						"IncidentReport": "",
					},
				},
			},
			wantErr: &framework.Error{
				Message: "VictorOps config is invalid: queryParameters[IncidentReport] cannot be an empty string.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_query_parameters_reserved_key": {
			request: &framework.Request[victorops_adapter.Config]{
				Address: "https://api.victorops.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "api-id",
						Password: "api-key",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
						},
					},
				},
				Config: &victorops_adapter.Config{
					QueryParameters: map[string]string{
						"IncidentReport": "offset=5",
					},
				},
			},
			wantErr: &framework.Error{
				Message: "VictorOps config is invalid: queryParameters[IncidentReport] " +
					`contains reserved parameter "offset" which is managed by the adapter.`,
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"valid_with_empty_config": {
			request: &framework.Request[victorops_adapter.Config]{
				Address: "https://api.victorops.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "api-id",
						Password: "api-key",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
						},
					},
				},
				Config: &victorops_adapter.Config{},
			},
			wantErr: nil,
		},
		"missing_auth": {
			request: &framework.Request[victorops_adapter.Config]{
				Address: "https://api.victorops.com",
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "VictorOps auth is missing required basic credentials (API ID and API key).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"missing_api_id": {
			request: &framework.Request[victorops_adapter.Config]{
				Address: "https://api.victorops.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "",
						Password: "api-key",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "VictorOps API ID (basic auth username) must not be empty.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"missing_api_key": {
			request: &framework.Request[victorops_adapter.Config]{
				Address: "https://api.victorops.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "api-id",
						Password: "",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "VictorOps API key (basic auth password) must not be empty.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_external_id": {
			request: &framework.Request[victorops_adapter.Config]{
				Address: "https://api.victorops.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "api-id",
						Password: "api-key",
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
				Message: "VictorOps entity external ID is invalid.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"missing_unique_id": {
			request: &framework.Request[victorops_adapter.Config]{
				Address: "https://api.victorops.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "api-id",
						Password: "api-key",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "NOT_THE_UNIQUE_ID",
						},
					},
				},
			},
			wantErr: &framework.Error{
				Message: "VictorOps requested entity attributes are missing a unique ID attribute.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"ordered_must_be_false": {
			request: &framework.Request[victorops_adapter.Config]{
				Address: "https://api.victorops.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "api-id",
						Password: "api-key",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
						},
					},
				},
				Ordered: true,
			},
			wantErr: &framework.Error{
				Message: "VictorOps Ordered property must be false.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"page_size_too_big": {
			request: &framework.Request[victorops_adapter.Config]{
				Address: "https://api.victorops.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "api-id",
						Password: "api-key",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: victorops_adapter.IncidentReport,
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "incidentNumber",
						},
					},
				},
				PageSize: 9999,
			},
			wantErr: &framework.Error{
				Message: "VictorOps provided page size (9999) exceeds the maximum (100).",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
	}

	adapter := &victorops_adapter.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(context.TODO(), tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
