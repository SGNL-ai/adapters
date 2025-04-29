// Copyright 2025 SGNL.ai, Inc.
package mysql_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/config"
	mysql "github.com/sgnl-ai/adapters/pkg/my-sql"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestValidationGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[mysql.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[mysql.Config]{
				Address: "sgnl.testaddress.us-east-1.rds.amazonaws.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastModified",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &mysql.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 100,
			},
			wantErr: nil,
		},
		"invalid_request_empty_config": {
			request: &framework.Request[mysql.Config]{
				Address: "sgnl.testaddress.us-east-1.rds.amazonaws.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastModified",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config:   &mysql.Config{},
				Ordered:  true,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "MySQL config is invalid: database is not set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_nil_config": {
			request: &framework.Request[mysql.Config]{
				Address: "sgnl.testaddress.us-east-1.rds.amazonaws.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastModified",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config:   nil,
				Ordered:  true,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "MySQL config is invalid: request contains no config.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_auth": {
			request: &framework.Request[mysql.Config]{
				Address: "sgnl.testaddress.us-east-1.rds.amazonaws.com",
				Auth:    &framework.DatasourceAuthCredentials{},
				Entity: framework.EntityConfig{
					ExternalId: "user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastModified",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &mysql.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required basic authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_auth_password": {
			request: &framework.Request[mysql.Config]{
				Address: "sgnl.testaddress.us-east-1.rds.amazonaws.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastModified",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &mysql.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required basic authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_configured_child_entities": {
			request: &framework.Request[mysql.Config]{
				Address: "sgnl.testaddress.us-east-1.rds.amazonaws.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastModified",
							Type:       framework.AttributeTypeDateTime,
						},
					},
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "phone_number",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "id",
									Type:       framework.AttributeTypeString,
								},
								{
									ExternalId: "number",
									Type:       framework.AttributeTypeString,
								},
							},
						},
					},
				},
				Config: &mysql.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  true,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "Requested entity does not support child entities.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_ordered_false": {
			request: &framework.Request[mysql.Config]{
				Address: "sgnl.testaddress.us-east-1.rds.amazonaws.com",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testusername",
						Password: "testpassword",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "user",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastModified",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &mysql.Config{
					CommonConfig: &config.CommonConfig{
						RequestTimeoutSeconds: testutil.GenPtr(10),
						LocalTimeZoneOffset:   -18000, // UTC−05:00 (EST)
					},
					Database: "sgnl",
				},
				Ordered:  false,
				PageSize: 100,
			},
			wantErr: &framework.Error{
				Message: "Ordered must be set to true.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
	}

	adapter := mysql.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
