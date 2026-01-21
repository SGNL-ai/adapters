// Copyright 2026 SGNL.ai, Inc.

// nolint: lll, goconst
package azuread_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/azuread"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestConstructEndpoint(t *testing.T) {
	tests := map[string]struct {
		request      *azuread.Request
		wantEndpoint string
		wantError    *framework.Error
	}{
		"nil_request": {
			request:      nil,
			wantEndpoint: "",
		},
		"invalid_entity": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1",
				EntityExternalID: "invalid",
				PageSize:         100,
				Token:            "SSWS testtoken",
			},
			wantError: &framework.Error{
				Message: "Provided entity external ID is invalid.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"complex_attributes": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "User",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "manager__id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "manager__displayName",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/users?$select=id&$expand=manager($select=id,displayName)&$top=100",
		},
		"invalid_complex_attributes": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "User",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "manager__",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantError: &framework.Error{
				Message: "Provided entity attribute list contains the following unsupported attribute: manager__.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"unsupported_parent_complex_attributes": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "User",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "owner__id",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			// Don't expand unsupported parent attributes
			wantEndpoint: "https://graph.microsoft.com/v1.0/users?$select=id,owner&$top=100",
		},
		"users_simple_no_attrs": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "User",
				PageSize:         100,
				Token:            "SSWS testtoken",
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/users?$select=id&$top=100",
		},
		"users_complex": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "User",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Filter:           testutil.GenPtr("displayName ne null"),
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "displayName",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "userPrincipalName",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/users?$select=id,displayName,userPrincipalName&$top=100&$filter=displayName+ne+null",
		},
		"users_cursor": {
			request: &azuread.Request{
				BaseURL:          "https://sgnl-dev.azureadpreview.com",
				APIVersion:       "v1",
				EntityExternalID: "users",
				PageSize:         100,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/users?$top=1&$skiptoken=RFNwdAIAAQAAACI6QWRhbXNATTM2NXgyMTQzNTUub25taWNyb3NvZnQuY29tKVVzZXJfNmU3Yjc2OGUtMDdlMi00ODEwLTg0NTktNDg1Zjg0ZjhmMjA0uQAAAAAAAAAAAAA"),
				},
				Token: "SSWS testtoken",
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/users?$top=1&$skiptoken=RFNwdAIAAQAAACI6QWRhbXNATTM2NXgyMTQzNTUub25taWNyb3NvZnQuY29tKVVzZXJfNmU3Yjc2OGUtMDdlMi00ODEwLTg0NTktNDg1Zjg0ZjhmMjA0uQAAAAAAAAAAAAA",
		},
		"groups_simple_no_attrs": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "Group",
				PageSize:         100,
				Token:            "SSWS testtoken",
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/groups?$select=id&$top=100",
		},
		"groups_complex": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "Group",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Filter:           testutil.GenPtr("displayName ne null"),
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "displayName",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "mail",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/groups?$select=id,displayName,mail&$top=100&$filter=displayName+ne+null",
		},
		"devices_simple_no_attrs": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "Device",
				PageSize:         100,
				Token:            "SSWS testtoken",
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/devices?$select=id&$top=100",
		},
		"devices_complex": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "Device",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Filter:           testutil.GenPtr("displayName ne null"),
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "displayName",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "deviceId",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/devices?$select=id,displayName,deviceId&$top=100&$filter=displayName+ne+null",
		},
		"applications_simple_no_attrs": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "Application",
				PageSize:         100,
				Token:            "SSWS testtoken",
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/applications?$select=id&$top=100",
		},
		"applications_complex": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "Application",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Filter:           testutil.GenPtr("displayName ne null"),
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "displayName",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "description",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/applications?$select=id,displayName,description&$top=100&$filter=displayName+ne+null",
		},
		"group_member_missing_collection_id": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "GroupMember",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Cursor:           &pagination.CompositeCursor[string]{},
			},
			wantError: &framework.Error{
				Message: "Unable to construct group member endpoint without valid cursor.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"group_member_simple_no_attrs": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "GroupMember",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID: testutil.GenPtr("1"),
				},
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/groups/1/members?$select=id&$top=100",
		},
		"group_member_complex": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "GroupMember",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Filter:           testutil.GenPtr("displayName ne null"),
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "displayName",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "userPrincipalName",
						Type:       framework.AttributeTypeString,
					},
				},
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID: testutil.GenPtr("1"),
				},
			},
			// Only `id` is requested for group members
			wantEndpoint: "https://graph.microsoft.com/v1.0/groups/1/members?$select=id&$top=100&$filter=displayName+ne+null",
		},
		"json_path_simple": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "User",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.manager.id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.manager.displayName",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/users?$select=id&$expand=manager($select=id,displayName)&$top=100",
		},
		"json_path_no_attribute": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "User",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantError: &framework.Error{
				Message: `Unable to extract any attributes from JSON path expression in provided attribute external id: "$.".`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"json_path_too_many_attributes": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "User",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.manager.address.state",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantError: &framework.Error{
				Message: "Too many attributes extracted from JSON path expression in provided attribute external id. Found: 3. Maximum supported: 2.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"json_path_single_attribute": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "User",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$..['displayName']",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/users?$select=id,displayName&$top=100",
		},
		"invalid_json_path": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "User",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$manager.id",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantError: &framework.Error{
				Message: `Provided entity attribute external id contains unsupported JSON path expression: "$manager.id".`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"json_path_unsupported_parent_attribute": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "User",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.invalid.id",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantError: &framework.Error{
				Message: `Unsupported parent attribute provided for the current entity type: "invalid".`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"roles_simple": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "Role",
				PageSize:         999,
				Token:            "SSWS testtoken",
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/directoryRoles?$select=id",
		},
		"roles_with_attrs": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "Role",
				PageSize:         999,
				Token:            "SSWS testtoken",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "displayName",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "description",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/directoryRoles?$select=id,displayName,description",
		},
		"roles_with_cursor": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "Role",
				PageSize:         3,
				Token:            "SSWS testtoken",
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJjdXJzb3IiOiIzIn0="),
				},
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "displayName",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/directoryRoles?$select=id,displayName",
		},
		"role_member_missing_collection_id": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "RoleMember",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Cursor:           &pagination.CompositeCursor[string]{},
			},
			wantError: &framework.Error{
				Message: "Unable to construct role member endpoint without valid cursor.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"role_member_simple_no_attrs": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "RoleMember",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID: testutil.GenPtr("1"),
				},
			},
			wantEndpoint: "https://graph.microsoft.com/v1.0/users/1/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=100",
		},
		"role_member_complex": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: "RoleMember",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "displayName",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "description",
						Type:       framework.AttributeTypeString,
					},
				},
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID: testutil.GenPtr("1"),
				},
			},
			// Only `id` is requested for role members
			wantEndpoint: "https://graph.microsoft.com/v1.0/users/1/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=100",
		},
		"role_assignment_schedule_request_complex": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: azuread.RoleAssignmentScheduleRequest,
				PageSize:         100,
				Token:            "SSWS testtoken",
				Filter:           testutil.GenPtr("status eq 'PendingApproval'"),
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "action",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: `https://graph.microsoft.com/v1.0/roleManagement/directory/roleAssignmentScheduleRequests?$select=id,action&$top=100&$skip=0&$filter=status+eq+%27PendingApproval%27`,
		},
		"role_assignment_schedule_request_complex_with_cursor": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: azuread.RoleAssignmentScheduleRequest,
				PageSize:         100,
				Skip:             100,
				Token:            "SSWS testtoken",
				Filter:           testutil.GenPtr("status eq 'PendingApproval'"),
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "action",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: `https://graph.microsoft.com/v1.0/roleManagement/directory/roleAssignmentScheduleRequests?$select=id,action&$top=100&$skip=100&$filter=status+eq+%27PendingApproval%27`,
		},
		"group_assignment_schedule_request_complex": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: azuread.GroupAssignmentScheduleRequest,
				PageSize:         100,
				Token:            "SSWS testtoken",
				Filter:           testutil.GenPtr("status eq 'PendingApproval'"),
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "status",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: `https://graph.microsoft.com/v1.0/identityGovernance/privilegedAccess/group/assignmentScheduleRequests?$select=id,status&$top=100&$skip=0&$filter=status+eq+%27PendingApproval%27`,
		},
		"group_assignment_schedule_request_complex_with_cursor": {
			request: &azuread.Request{
				BaseURL:          "https://graph.microsoft.com",
				APIVersion:       "v1.0",
				EntityExternalID: azuread.GroupAssignmentScheduleRequest,
				PageSize:         100,
				Skip:             100,
				Token:            "SSWS testtoken",
				Filter:           testutil.GenPtr("status eq 'PendingApproval'"),
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "status",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: `https://graph.microsoft.com/v1.0/identityGovernance/privilegedAccess/group/assignmentScheduleRequests?$select=id,status&$top=100&$skip=100&$filter=status+eq+%27PendingApproval%27`,
		},
		"group_members_with_advanced_filters": { // expect `$count=true` in the endpoint
			request: &azuread.Request{
				BaseURL:            "https://graph.microsoft.com",
				APIVersion:         "v1.0",
				EntityExternalID:   azuread.GroupMember,
				PageSize:           100,
				Token:              "SSWS testtoken",
				Filter:             testutil.GenPtr("status eq 'PendingApproval'"),
				ParentFilter:       testutil.GenPtr("id eq '123'"),
				UseAdvancedFilters: true,
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "status",
						Type:       framework.AttributeTypeString,
					},
				},
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID: testutil.GenPtr("1"),
				},
			},
			wantEndpoint: `https://graph.microsoft.com/v1.0/groups/1/members?$select=id&$top=100&$filter=status+eq+%27PendingApproval%27&$count=true`,
		},
		"groups_with_advanced_filters": { // expect `$count=true` in the endpoint
			request: &azuread.Request{
				BaseURL:            "https://graph.microsoft.com",
				APIVersion:         "v1.0",
				EntityExternalID:   azuread.Group,
				PageSize:           100,
				Token:              "SSWS testtoken",
				Filter:             testutil.GenPtr("startsWith(displayName, 'Infra')"),
				UseAdvancedFilters: true,
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "displayName",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantEndpoint: `https://graph.microsoft.com/v1.0/groups?$select=id,displayName&$top=100&$filter=startsWith%28displayName%2C+%27Infra%27%29&$count=true`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotEndpoint, gotError := azuread.ConstructEndpoint(tt.request)

			if !reflect.DeepEqual(gotError, tt.wantError) {
				t.Errorf("gotError: %v, wantError: %v", gotError, tt.wantError)
			}

			if !reflect.DeepEqual(gotEndpoint, tt.wantEndpoint) {
				t.Errorf("gotEndpoint: %v, wantEndpoint: %v", gotEndpoint, tt.wantEndpoint)
			}
		})
	}
}
