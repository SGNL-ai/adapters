// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package bamboohr_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/bamboohr"
)

func TestConstructEndpoint(t *testing.T) {
	tests := map[string]struct {
		request          *bamboohr.Request
		wantEndpointInfo *bamboohr.EndpointInfo
		wantError        *framework.Error
	}{
		"nil_request": {
			request: nil,
			wantError: &framework.Error{
				Message: "Request is nil.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"only_current_true": {
			request: &bamboohr.Request{
				BaseURL:       "https://test-instance.bamboohr.com",
				APIVersion:    "v1",
				CompanyDomain: "SGNL",
				OnlyCurrent:   true,
				PageSize:      100,
				EntityConfig:  PopulateDefaultEmployeeEntityConfig(),
			},
			wantEndpointInfo: &bamboohr.EndpointInfo{
				URL:  "https://test-instance.bamboohr.com/SGNL/v1/reports/custom?format=JSON&onlyCurrent=true",
				Body: "{\"fields\":[\"id\",\"bestEmail\",\"dateOfBirth\",\"fullName1\",\"isPhotoUploaded\",\"customcustomBoolField\",\"supervisorEId\",\"supervisorEmail\",\"lastChanged\"]}",
			},
		},
		"only_current_false": {
			request: &bamboohr.Request{
				BaseURL:       "https://test-instance.bamboohr.com",
				APIVersion:    "v1",
				CompanyDomain: "SGNL",
				OnlyCurrent:   false,
				PageSize:      100,
				EntityConfig:  PopulateDefaultEmployeeEntityConfig(),
			},
			wantEndpointInfo: &bamboohr.EndpointInfo{
				URL:  "https://test-instance.bamboohr.com/SGNL/v1/reports/custom?format=JSON&onlyCurrent=false",
				Body: "{\"fields\":[\"id\",\"bestEmail\",\"dateOfBirth\",\"fullName1\",\"isPhotoUploaded\",\"customcustomBoolField\",\"supervisorEId\",\"supervisorEmail\",\"lastChanged\"]}",
			},
		},
		"random_attributes": {
			request: &bamboohr.Request{
				BaseURL:       "https://test-instance.bamboohr.com",
				APIVersion:    "v1",
				CompanyDomain: "SGNL",
				OnlyCurrent:   false,
				PageSize:      100,
				EntityConfig: &framework.EntityConfig{
					ExternalId: "Employee",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeInt64,
							List:       false,
						},
						{
							ExternalId: "bestEmail",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
			},
			wantEndpointInfo: &bamboohr.EndpointInfo{
				URL:  "https://test-instance.bamboohr.com/SGNL/v1/reports/custom?format=JSON&onlyCurrent=false",
				Body: "{\"fields\":[\"id\",\"bestEmail\"]}",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotEndpointInfo, gotError := bamboohr.ConstructEndpoint(tt.request)

			if !reflect.DeepEqual(gotError, tt.wantError) {
				t.Errorf("gotError: %v, wantError: %v", gotError, tt.wantError)
			}

			if diff := cmp.Diff(gotEndpointInfo, tt.wantEndpointInfo); diff != "" {
				t.Errorf("gotEndpoint: %v, wantEndpoint: %v", gotEndpointInfo, tt.wantEndpointInfo)
			}
		})
	}
}

func PopulateDefaultEmployeeEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: bamboohr.Employee,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeInt64,
				List:       false,
			},
			{
				ExternalId: "bestEmail",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "dateOfBirth",
				Type:       framework.AttributeTypeDateTime, // BambooHR Date type (YYYY-MM-DD)
				List:       false,
			},
			{
				ExternalId: "fullName1",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "isPhotoUploaded",
				Type:       framework.AttributeTypeBool, // true or false
				List:       false,
			},
			{
				ExternalId: "customcustomBoolField",
				Type:       framework.AttributeTypeBool, // 0 or 1
				List:       false,
			},
			{
				ExternalId: "supervisorEId",
				Type:       framework.AttributeTypeInt64,
				List:       false,
			},
			{
				ExternalId: "supervisorEmail",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "lastChanged",
				Type:       framework.AttributeTypeDateTime, // BambooHR timestamp type (RFC3339)
				List:       false,
			},
		},
	}
}
