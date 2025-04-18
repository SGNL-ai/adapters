// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package googleworkspace_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	googleworkspace "github.com/sgnl-ai/adapters/pkg/google-workspace"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestConstructEndpoint(t *testing.T) {
	tests := map[string]struct {
		request      *googleworkspace.Request
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
		"nil_compositecursor_user_entity": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "User",
				Domain:           testutil.GenPtr("sgnldemos.com"),
				Cursor:           nil,
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/users?domain=sgnldemos.com&maxResults=100",
		},
		"nil_cursor_value_user_entity": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "User",
				Domain:           testutil.GenPtr("sgnldemos.com"),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: nil,
				},
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/users?domain=sgnldemos.com&maxResults=100",
		},
		"customer_used_instead_of_domain": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "User",
				Customer:         testutil.GenPtr("sgnldemos.com"),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: nil,
				},
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/users?customer=sgnldemos.com&maxResults=100",
		},
		"user_entity_with_cursor": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "User",
				Customer:         testutil.GenPtr("sgnldemos.com"),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("nextPage"),
				},
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/users?customer=sgnldemos.com&maxResults=100&pageToken=nextPage",
		},
		"user_entity_with_default_filters": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "User",
				Customer:         testutil.GenPtr("sgnldemos.com"),
				Filters: googleworkspace.Filters{
					UserFilters: &googleworkspace.UserFilters{},
				},
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("nextPage"),
				},
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/users?customer=sgnldemos.com&maxResults=100&pageToken=nextPage&showDeleted=false",
		},
		"user_entity_excludes_membered_filters": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "User",
				Customer:         testutil.GenPtr("sgnldemos.com"),
				Filters: googleworkspace.Filters{
					MemberFilters: &googleworkspace.MemberFilters{},
				},
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("nextPage"),
				},
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/users?customer=sgnldemos.com&maxResults=100&pageToken=nextPage",
		},
		"user_entity_with_filters": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "User",
				Customer:         testutil.GenPtr("sgnldemos.com"),
				Filters: googleworkspace.Filters{
					UserFilters: &googleworkspace.UserFilters{
						Query:       testutil.GenPtr("email:*@gmail.com"),
						ShowDeleted: true,
					},
				},
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("nextPage"),
				},
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/users?customer=sgnldemos.com&maxResults=100&pageToken=nextPage&query=email%3A%2A%40gmail.com&showDeleted=true",
		},
		"group_entity_with_default_filters": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "Group",
				Customer:         testutil.GenPtr("sgnldemos.com"),
				Filters: googleworkspace.Filters{
					GroupFilters: &googleworkspace.GroupFilters{},
				},
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("nextPage"),
				},
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/groups?customer=sgnldemos.com&maxResults=100&pageToken=nextPage",
		},
		"group_entity_excludes_user_filters": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "Group",
				Customer:         testutil.GenPtr("sgnldemos.com"),
				Filters: googleworkspace.Filters{
					UserFilters: &googleworkspace.UserFilters{},
				},
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("nextPage"),
				},
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/groups?customer=sgnldemos.com&maxResults=100&pageToken=nextPage",
		},
		"group_entity_with_filters": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "Group",
				Customer:         testutil.GenPtr("sgnldemos.com"),
				Filters: googleworkspace.Filters{
					GroupFilters: &googleworkspace.GroupFilters{
						Query: testutil.GenPtr("email:*@gmail.com"),
					},
				},
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("nextPage"),
				},
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/groups?customer=sgnldemos.com&maxResults=100&pageToken=nextPage&query=email%3A%2A%40gmail.com",
		},
		"nil_compositecursor_member_entity": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "Member",
				Domain:           testutil.GenPtr("sgnldemos.com"),
				Cursor:           nil,
			},
			wantError: &framework.Error{
				Message: "Collection ID is nil for Member entity, unable to form request URI.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"nil_collection_id_value_member_entity": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "Member",
				Domain:           testutil.GenPtr("sgnldemos.com"),
				Cursor:           &pagination.CompositeCursor[string]{},
			},
			wantError: &framework.Error{
				Message: "Collection ID is nil for Member entity, unable to form request URI.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"nil_cursor_value_member_entity": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "Member",
				Domain:           testutil.GenPtr("sgnldemos.com"),
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID: testutil.GenPtr("collectionId"),
				},
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/groups/collectionId/members?domain=sgnldemos.com&maxResults=100",
		},
		"member_entity_with_cursor": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "Member",
				Domain:           testutil.GenPtr("sgnldemos.com"),
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID: testutil.GenPtr("collectionId"),
					Cursor:       testutil.GenPtr("nextPage"),
				},
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/groups/collectionId/members?domain=sgnldemos.com&maxResults=100&pageToken=nextPage",
		},
		"member_entity_excludes_group_filters": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "Member",
				Domain:           testutil.GenPtr("sgnldemos.com"),
				Filters: googleworkspace.Filters{
					GroupFilters: &googleworkspace.GroupFilters{},
				},
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID: testutil.GenPtr("collectionId"),
					Cursor:       testutil.GenPtr("nextPage"),
				},
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/groups/collectionId/members?domain=sgnldemos.com&maxResults=100&pageToken=nextPage",
		},
		"member_entity_with_default_filters": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "Member",
				Domain:           testutil.GenPtr("sgnldemos.com"),
				Filters: googleworkspace.Filters{
					MemberFilters: &googleworkspace.MemberFilters{},
				},
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID: testutil.GenPtr("collectionId"),
					Cursor:       testutil.GenPtr("nextPage"),
				},
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/groups/collectionId/members?domain=sgnldemos.com&includeDerivedMembership=false&maxResults=100&pageToken=nextPage",
		},
		"member_entity_with_set_filters": {
			request: &googleworkspace.Request{
				BaseURL:          "https://admin.googleapis.com",
				APIVersion:       "v1",
				PageSize:         100,
				EntityExternalID: "Member",
				Domain:           testutil.GenPtr("sgnldemos.com"),
				Filters: googleworkspace.Filters{
					MemberFilters: &googleworkspace.MemberFilters{
						IncludeDerivedMembership: true,
						Roles:                    testutil.GenPtr("ADMIN"),
					},
				},
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID: testutil.GenPtr("collectionId"),
					Cursor:       testutil.GenPtr("nextPage"),
				},
			},
			wantEndpoint: "https://admin.googleapis.com/admin/directory/v1/groups/collectionId/members?domain=sgnldemos.com&includeDerivedMembership=true&maxResults=100&pageToken=nextPage&roles=ADMIN",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotEndpoint, gotError := googleworkspace.ConstructEndpoint(tt.request)

			if !reflect.DeepEqual(gotError, tt.wantError) {
				t.Errorf("gotError: %v, wantError: %v", gotError, tt.wantError)
			}

			if !reflect.DeepEqual(gotEndpoint, tt.wantEndpoint) {
				t.Errorf("gotEndpoint: %v, wantEndpoint: %v", gotEndpoint, tt.wantEndpoint)
			}
		})
	}
}

func PopulateDefaultUserEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: "User",
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "primaryEmail",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.name.fullName",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "isAdmin",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "isDelegatedAdmin",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "changePasswordAtNextLogin",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "creationTime",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
			{
				ExternalId: "nonEditableAliases",
				Type:       framework.AttributeTypeString,
				List:       true,
			},
			{
				ExternalId: "customerId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
		ChildEntities: []*framework.EntityConfig{
			{
				ExternalId: "emails",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "address",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
					{
						ExternalId: "type",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
					{
						ExternalId: "primary",
						Type:       framework.AttributeTypeBool,
						List:       false,
					},
				},
			},
		},
	}
}

func PopulateDefaultGroupEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: "Group",
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "kind",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "email",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "name",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "adminCreated",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "nonEditableAliases",
				Type:       framework.AttributeTypeString,
				List:       true,
			},
			{
				ExternalId: "directMembersCount",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
	}
}

func PopulateDefaultMemberEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: "Member",
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "uniqueId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "groupId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "kind",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "email",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "role",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "type",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "status",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
	}
}
