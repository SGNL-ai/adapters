// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package workday_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
	"github.com/sgnl-ai/adapters/pkg/workday"
)

func TestConstructEndpoint(t *testing.T) {
	tests := map[string]struct {
		request      *workday.Request
		wantEndpoint string
		wantError    *framework.Error
	}{
		"nil_request": {
			request: nil,
			wantError: &framework.Error{
				Message: "Request is nil.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"nil_composite_cursor": {
			request: &workday.Request{
				BaseURL:        "https://test-instance.workday.com",
				APIVersion:     "v1",
				OrganizationID: "SGNL",
				PageSize:       100,
				EntityConfig:   PopulateDefaultWorkerEntityConfig(),
				Cursor:         nil,
			},
			wantEndpoint: "https://test-instance.workday.com/api/wql/v1/SGNL/data?limit=100&offset=0&query=SELECT+FTE%2C+company%2C+email_Work%2C+employeeID%2C+employeeType%2C+gender%2C+hireDate%2C+jobTitle%2C+managementLevel%2C+positionID%2C+worker%2C+workerActive+FROM+allWorkers",
		},
		"nil_cursor": {
			request: &workday.Request{
				BaseURL:        "https://test-instance.workday.com",
				APIVersion:     "v1",
				OrganizationID: "SGNL",
				PageSize:       100,
				EntityConfig:   PopulateDefaultWorkerEntityConfig(),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: nil,
				},
			},
			wantEndpoint: "https://test-instance.workday.com/api/wql/v1/SGNL/data?limit=100&offset=0&query=SELECT+FTE%2C+company%2C+email_Work%2C+employeeID%2C+employeeType%2C+gender%2C+hireDate%2C+jobTitle%2C+managementLevel%2C+positionID%2C+worker%2C+workerActive+FROM+allWorkers",
		},
		"negative_cursor_offset": {
			request: &workday.Request{
				BaseURL:        "https://test-instance.workday.com",
				APIVersion:     "v1",
				OrganizationID: "SGNL",
				PageSize:       100,
				EntityConfig:   PopulateDefaultWorkerEntityConfig(),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(-1)),
				},
			},
			wantError: &framework.Error{
				Message: "Cursor value must be greater than 0.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"cursor_offset_of_zero": {
			request: &workday.Request{
				BaseURL:        "https://test-instance.workday.com",
				APIVersion:     "v1",
				OrganizationID: "SGNL",
				PageSize:       100,
				EntityConfig:   PopulateDefaultWorkerEntityConfig(),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr(int64(0)),
				},
			},
			wantError: &framework.Error{
				Message: "Cursor value must be greater than 0.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"with_cursor_offset": {
			request: &workday.Request{
				BaseURL:        "https://test-instance.workday.com",
				APIVersion:     "v1",
				OrganizationID: "SGNL",
				PageSize:       50,
				EntityConfig:   PopulateDefaultWorkerEntityConfig(),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](50),
				},
			},
			wantEndpoint: "https://test-instance.workday.com/api/wql/v1/SGNL/data?limit=50&offset=50&query=SELECT+FTE%2C+company%2C+email_Work%2C+employeeID%2C+employeeType%2C+gender%2C+hireDate%2C+jobTitle%2C+managementLevel%2C+positionID%2C+worker%2C+workerActive+FROM+allWorkers",
		},
		"uncommon_json_path_support": {
			request: &workday.Request{
				BaseURL:        "https://test-instance.workday.com",
				APIVersion:     "v1",
				OrganizationID: "SGNL",
				PageSize:       50,
				EntityConfig:   PopulateWorkerEntityConfigWithNoChildren(),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](50),
				},
			},
			wantEndpoint: "https://test-instance.workday.com/api/wql/v1/SGNL/data?limit=50&offset=50&query=SELECT+FTE%2C+company%2C+email_Work%2C+employeeID%2C+employeeType%2C+gender%2C+hireDate%2C+jobTitle%2C+managementLevel%2C+positionID%2C+worker%2C+workerActive+FROM+allWorkers",
		},
		"request_with_ordering_true": {
			request: &workday.Request{
				BaseURL:        "https://test-instance.workday.com",
				APIVersion:     "v1",
				OrganizationID: "SGNL",
				Ordered:        true,
				PageSize:       50,
				EntityConfig:   PopulateDefaultWorkerEntityConfig(),
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](50),
				},
			},
			wantEndpoint: "https://test-instance.workday.com/api/wql/v1/SGNL/data?limit=50&offset=50&query=SELECT+FTE%2C+company%2C+email_Work%2C+employeeID%2C+employeeType%2C+gender%2C+hireDate%2C+jobTitle%2C+managementLevel%2C+positionID%2C+worker%2C+workerActive+FROM+allWorkers+ORDER+BY+worker+ASC",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotEndpoint, gotError := workday.ConstructEndpoint(tt.request)

			if !reflect.DeepEqual(gotError, tt.wantError) {
				t.Errorf("gotError: %v, wantError: %v", gotError, tt.wantError)
			}

			if !reflect.DeepEqual(gotEndpoint, tt.wantEndpoint) {
				t.Errorf("gotEndpoint: %v, wantEndpoint: %v", gotEndpoint, tt.wantEndpoint)
			}
		})
	}
}

func PopulateDefaultWorkerEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: "allWorkers",
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "$.worker.id",
				Type:       framework.AttributeTypeString,
				List:       false,
				UniqueId:   true,
			},
			{
				ExternalId: "$.worker.descriptor",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "employeeID",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "workerActive",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "$.managementLevel.descriptor",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.managementLevel.id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.company.descriptor",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.company.id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.gender.descriptor",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.gender.id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "hireDate",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
			{
				ExternalId: "FTE",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "positionID",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "jobTitle",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
		ChildEntities: []*framework.EntityConfig{
			{
				ExternalId: "email_Work",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "descriptor",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
				},
			},
			{
				ExternalId: "employeeType",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "descriptor",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
				},
			},
		},
	}
}

func PopulateWorkerEntityConfigWithNoChildren() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: "allWorkers",
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "$.worker.id",
				Type:       framework.AttributeTypeString,
				List:       false,
				UniqueId:   true,
			},
			{
				ExternalId: "$.worker.descriptor",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "employeeID",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "workerActive",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "$.managementLevel.descriptor",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.managementLevel.id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.company.descriptor",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.company.id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.gender.descriptor",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.gender.id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "hireDate",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
			{
				ExternalId: "FTE",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "positionID",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "jobTitle",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.email_Work[0].descriptor",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.email_Work[0].id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
		ChildEntities: []*framework.EntityConfig{
			{
				ExternalId: "employeeType",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "descriptor",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
				},
			},
		},
	}
}
