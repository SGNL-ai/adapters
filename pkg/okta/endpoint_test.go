// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package okta_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/okta"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestConstructEndpoint(t *testing.T) {
	tests := map[string]struct {
		request      *okta.Request
		wantEndpoint string
		wantError    *framework.Error
	}{
		"users_simple": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "User",
				PageSize:         100,
				Token:            "SSWS testtoken",
			},
			wantEndpoint: "https://test-instance.oktapreview.com/api/v1/users?limit=100",
		},
		"users_simple_filter": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "User",
				PageSize:         100,
				Filter:           "status eq \"ACTIVE\"",
				Token:            "SSWS testtoken",
			},
			wantEndpoint: "https://test-instance.oktapreview.com/api/v1/users?filter=status+eq+%22ACTIVE%22&limit=100",
		},
		"users_simple_search": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "User",
				PageSize:         100,
				Search:           "profile.state eq \"CO\"",
				Token:            "SSWS testtoken",
			},
			wantEndpoint: "https://test-instance.oktapreview.com/api/v1/users?search=profile.state+eq+%22CO%22&limit=100",
		},
		"users_search_and_filter": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "User",
				PageSize:         100,
				Filter:           "status eq \"ACTIVE\"",
				Search:           "profile.state eq \"CO\"",
				Token:            "SSWS testtoken",
			},
			wantEndpoint: "https://test-instance.oktapreview.com/api/v1/users?filter=status+eq+%22ACTIVE%22&search=profile.state+eq+%22CO%22&limit=100",
		},
		"users_cursor": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "User",
				PageSize:         100,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/users?after=00g3zvuhepAwReSDo1d7&limit=100"),
				},
				Token: "SSWS testtoken",
			},
			wantEndpoint: "https://test-instance.oktapreview.com/api/v1/users?after=00g3zvuhepAwReSDo1d7&limit=100",
		},
		"users_cursor_filter": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "User",
				PageSize:         100,
				Filter:           "status eq \"ACTIVE\"",
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/users?filter=status+eq+%22ACTIVE%22&after=00g3zvuhepAwReSDo1d7&limit=100"),
				},
				Token: "SSWS testtoken",
			},
			wantEndpoint: "https://test-instance.oktapreview.com/api/v1/users?filter=status+eq+%22ACTIVE%22&after=00g3zvuhepAwReSDo1d7&limit=100",
		},
		"users_cursor_search": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "User",
				PageSize:         100,
				Search:           "profile.state eq \"CO\"",
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/users?search=profile.state+eq+%22CO%22&after=00g3zvuhepAwReSDo1d7&limit=100"),
				},
				Token: "SSWS testtoken",
			},
			wantEndpoint: "https://test-instance.oktapreview.com/api/v1/users?search=profile.state+eq+%22CO%22&after=00g3zvuhepAwReSDo1d7&limit=100",
		},
		"groups_simple": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "Group",
				PageSize:         100,
				Token:            "SSWS testtoken",
			},
			wantEndpoint: "https://test-instance.oktapreview.com/api/v1/groups?filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22&limit=100",
		},
		"groups_cursor": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "Group",
				PageSize:         100,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/groups?after=00g3zvuhepAwReSDo1d7&limit=100"),
				},
				Token: "SSWS testtoken",
			},
			wantEndpoint: "https://test-instance.oktapreview.com/api/v1/groups?after=00g3zvuhepAwReSDo1d7&limit=100",
		},
		"invalid_entity": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
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
		"nil_request": {
			request:      nil,
			wantEndpoint: "",
		},
		"invalid_user_filter": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "User",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Filter:           "id eq ",
			},
			wantError: &framework.Error{
				Message: "Provided filter is invalid.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_user_search": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "User",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Search:           "id eq ",
			},
			wantError: &framework.Error{
				Message: "Provided search syntax is invalid.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		// First page of the first group of a group member sync. The endpoint should be set to the first page of
		// members for the first group.
		"group_members_group_page_1_of_2_member_page_1_of_2": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "GroupMember",
				PageSize:         100,
				Token:            "SSWS testtoken",
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID:     testutil.GenPtr("00g1emaKYZTWRYYRRTSK"),
					CollectionCursor: testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/groups?after=00g1emaKYZTWRYYRRTSK&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22"),
				},
			},
			wantEndpoint: "https://test-instance.oktapreview.com/api/v1/groups/00g1emaKYZTWRYYRRTSK/users?limit=100",
		},
		// Second page of the first group of a group member sync. The endpoint should be set to the second page of
		// members for the first group, based on the cursor.
		"group_members_group_page_1_of_2_member_page_2_of_2": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "GroupMember",
				PageSize:         100,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:           testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/groups/00g1emaKYZTWRYYRRTSK/users?after=00ub0oNGTSWTBKOLGLNR&limit=100"),
					CollectionID:     testutil.GenPtr("00g1emaKYZTWRYYRRTSK"),
					CollectionCursor: testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/groups?after=00g1emaKYZTWRYYRRTSK&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22"),
				},
				Token: "SSWS testtoken",
			},
			wantEndpoint: "https://test-instance.oktapreview.com/api/v1/groups/00g1emaKYZTWRYYRRTSK/users?after=00ub0oNGTSWTBKOLGLNR&limit=100",
		},
		// First page of the second group of a group member sync. The endpoint should be set to the first page of
		// members for the second group, e.g. the first Group at the CollectionCursor (GroupCursor) specified.
		"group_members_group_page_2_of_2_member_page_1_of_2": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "GroupMember",
				PageSize:         100,
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID: testutil.GenPtr("00garwpuyxHaWOkdV0g4"),
				},
				Token: "SSWS testtoken",
			},
			wantEndpoint: "https://test-instance.oktapreview.com/api/v1/groups/00garwpuyxHaWOkdV0g4/users?limit=100",
		},
		// Second page of the second group of a group member sync. The endpoint should be set to the second page of
		// members for the second group, based on the cursor.
		"group_members_group_page_2_of_2_member_page_2_of_2": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "GroupMember",
				PageSize:         100,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:       testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/groups/00garwpuyxHaWOkdV0g4/users?after=00ub0oNGTSWTBKOLGLNR&limit=100"),
					CollectionID: testutil.GenPtr("00garwpuyxHaWOkdV0g4"),
				},
				Token: "SSWS testtoken",
			},
			wantEndpoint: "https://test-instance.oktapreview.com/api/v1/groups/00garwpuyxHaWOkdV0g4/users?after=00ub0oNGTSWTBKOLGLNR&limit=100",
		},
		// Invalid group member cursors are not evaluated at this point.
		"group_member_invalid_member_cursor": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "GroupMember",
				PageSize:         100,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:       testutil.GenPtr("https://test-instance.oktapreview.com/invalid"),
					CollectionID: testutil.GenPtr("00garwpuyxHaWOkdV0g4"),
				},
				Token: "SSWS testtoken",
			},
			wantEndpoint: "https://test-instance.oktapreview.com/invalid",
		},
		"group_member_missing_collection_id": {
			request: &okta.Request{
				BaseURL:          "https://test-instance.oktapreview.com",
				APIVersion:       "v1",
				EntityExternalID: "GroupMember",
				PageSize:         100,
				Cursor:           &pagination.CompositeCursor[string]{},
			},
			wantError: &framework.Error{
				Message: "Unable to construct group member endpoint without valid cursor.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotEndpoint, gotError := okta.ConstructEndpoint(tt.request)

			if !reflect.DeepEqual(gotError, tt.wantError) {
				t.Errorf("gotError: %v, wantError: %v", gotError, tt.wantError)
			}

			if !reflect.DeepEqual(gotEndpoint, tt.wantEndpoint) {
				t.Errorf("gotEndpoint: %v, wantEndpoint: %v", gotEndpoint, tt.wantEndpoint)
			}
		})
	}
}
