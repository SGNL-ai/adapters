// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package bamboohr_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/bamboohr"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[bamboohr.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[bamboohr.Config]{
				Address: "api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_request_empty_config": {
			request: &framework.Request[bamboohr.Config]{
				Address: "api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config:   &bamboohr.Config{},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "BambooHR config is invalid: apiVersion is not set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_request_missing_config": {
			request: &framework.Request[bamboohr.Config]{
				Address: "api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "BambooHR config is invalid: request contains no config.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_api_version": {
			request: &framework.Request[bamboohr.Config]{
				Config: &bamboohr.Config{
					APIVersion:    "v1.1",
					CompanyDomain: "sgnl",
				},
				Address: "api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
				Cursor:   "",
			},
			wantErr: &framework.Error{
				Message: "BambooHR config is invalid: apiVersion is not supported: v1.1.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"valid_date_and_onlyCurrent_format": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
					AttributeMappings: &bamboohr.AttributeMappings{
						Date: testutil.GenPtr("yyyy-mm-dd"),
					},
					OnlyCurrent: true,
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_date_format": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
					AttributeMappings: &bamboohr.AttributeMappings{
						Date: testutil.GenPtr("DDMMYYYY"),
					},
					OnlyCurrent: true,
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "BambooHR config is invalid: date format is not supported: DDMMYYYY.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"missing_company_domain": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion: "v1",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "BambooHR config is invalid: companyDomain is not set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"missing_optional_fields": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"attribute_mapping_with_missing_bool_mapping": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:        "v1",
					CompanyDomain:     "sgnl",
					AttributeMappings: &bamboohr.AttributeMappings{},
				},
				Ordered:  false,
				PageSize: 250,
			},
		},
		"bool_mapping_missing_both_true_and_false_keys": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
					AttributeMappings: &bamboohr.AttributeMappings{
						BoolMappings: &bamboohr.BoolAttributeMappings{},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "BambooHR config is invalid: Both attributeMappings.bool.true and attributeMappings.bool.false must be set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_bool_mapping_missing_false_key": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
					AttributeMappings: &bamboohr.AttributeMappings{
						BoolMappings: &bamboohr.BoolAttributeMappings{
							True: []string{"True", "yes", "1"},
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "BambooHR config is invalid: Both attributeMappings.bool.true and attributeMappings.bool.false must be set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_bool_mapping_missing_true_key": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
					AttributeMappings: &bamboohr.AttributeMappings{
						BoolMappings: &bamboohr.BoolAttributeMappings{
							False: []string{"False", "no", "0"},
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "BambooHR config is invalid: Both attributeMappings.bool.true and attributeMappings.bool.false must be set.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"valid_bool_mappings": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
					AttributeMappings: &bamboohr.AttributeMappings{
						BoolMappings: &bamboohr.BoolAttributeMappings{
							True:  []string{"True", "yes", "1"},
							False: []string{"False", "no", "0"},
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
		},
		"invalid_shared_bool_mapping_values": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
					AttributeMappings: &bamboohr.AttributeMappings{
						BoolMappings: &bamboohr.BoolAttributeMappings{
							True:  []string{"True", "yes", "1", "SHARED"},
							False: []string{"False", "no", "0", "SHARED"},
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "BambooHR config is invalid: These errors were found in your bool mapping list: [Identical mapping found for both bool.false and bool.true: SHARED].",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_duplicate_bool_mapping_value_in_true_key": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
					AttributeMappings: &bamboohr.AttributeMappings{
						BoolMappings: &bamboohr.BoolAttributeMappings{
							True:  []string{"True", "DUPLICATE", "DUPLICATE"},
							False: []string{"False", "no", "0"},
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "BambooHR config is invalid: These errors were found in your bool mapping list: [attributeMappings.bool.true has a duplicate value: DUPLICATE].",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_duplicate_bool_mapping_value_in_false_key": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
					AttributeMappings: &bamboohr.AttributeMappings{
						BoolMappings: &bamboohr.BoolAttributeMappings{
							True:  []string{"True", "yes", "1"},
							False: []string{"False", "DUPLICATE", "DUPLICATE"},
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "BambooHR config is invalid: These errors were found in your bool mapping list: [attributeMappings.bool.false has a duplicate value: DUPLICATE].",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_multiple_shared_mappings_and_duplicate_values": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
					AttributeMappings: &bamboohr.AttributeMappings{
						BoolMappings: &bamboohr.BoolAttributeMappings{
							True:  []string{"True", "yes", "1", "yes", "no"},
							False: []string{"False", "no", "0", "DUPLICATE", "DUPLICATE"},
						},
					},
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "BambooHR config is invalid: These errors were found in your bool mapping list: [attributeMappings.bool.true has a duplicate value: yes attributeMappings.bool.false has a duplicate value: DUPLICATE Identical mapping found for both bool.false and bool.true: no].",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_ordered_true": {
			request: &framework.Request[bamboohr.Config]{
				Address: "api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
				},
				Ordered:  true,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Ordered must be set to false.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"valid_https_prefix": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: nil,
		},
		"invalid_http_prefix": {
			request: &framework.Request[bamboohr.Config]{
				Address: "http://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "The provided HTTP protocol is not supported.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_auth": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required basic authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_missing_basic_auth": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth:    &framework.DatasourceAuthCredentials{},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required basic authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_basic_auth_details": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required basic authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_basic_auth_user_details": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required basic authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_basic_auth_password_details": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required basic authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
		"invalid_entity_type": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "INVALID ENTITY",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Provided entity external ID is invalid.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_missing_unique_attribute": {
			request: &framework.Request[bamboohr.Config]{
				Address: "https://api.bamboohr.com/api/gateway.php/SGNL",
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "Test API Key",
						Password: "xxx",
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "fullName1",
							Type:       framework.AttributeTypeString,
						},
						{
							ExternalId: "lastChanged",
							Type:       framework.AttributeTypeDateTime,
						},
					},
				},
				Config: &bamboohr.Config{
					APIVersion:    "v1",
					CompanyDomain: "sgnl",
				},
				Ordered:  false,
				PageSize: 250,
			},
			wantErr: &framework.Error{
				Message: "Requested entity attributes are missing unique ID attribute.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
	}

	adapter := bamboohr.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := adapter.ValidateGetPageRequest(nil, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
