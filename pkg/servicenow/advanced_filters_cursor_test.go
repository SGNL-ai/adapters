// Copyright 2025 SGNL.ai, Inc.

// nolint: lll
package servicenow_test

import (
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/servicenow"
	"github.com/sgnl-ai/adapters/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMarshalAdvancedFilterCursor(t *testing.T) {
	tests := map[string]struct {
		inputCursor *servicenow.AdvancedFilterCursor
		wantCursor  string
		wantErr     *framework.Error
	}{
		"success_with_implicit_filter_only": {
			inputCursor: &servicenow.AdvancedFilterCursor{
				ImplicitFilterCursor: &servicenow.ImplicitFilterCursor{
					EntityFilterIndex: 0,
					MemberFilterIndex: 0,
					Cursor: &pagination.CompositeCursor[string]{
						Cursor:           testutil.GenPtr("cursor1"),
						CollectionCursor: testutil.GenPtr("collection1"),
					},
				},
			},
			wantCursor: "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6eyJjdXJzb3IiOiJjdXJzb3IxIiwiY29sbGVjdGlvbkN1cnNvciI6ImNvbGxlY3Rpb24xIn19fQ==",
			wantErr:    nil,
		},
		"success_with_related_filter_only": {
			inputCursor: &servicenow.AdvancedFilterCursor{
				RelatedFilterCursor: &servicenow.RelatedFilterCursor{
					EntityIndex:  0,
					EntityCursor: testutil.GenPtr("assigned_toIN{$.sys_user.sys_id}"),
					RelatedEntityCursor: &pagination.CompositeCursor[string]{
						Cursor:           testutil.GenPtr("cursor1"),
						CollectionCursor: testutil.GenPtr("collection1"),
					},
				},
			},
			wantCursor: "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjowLCJlbnRpdHlDdXJzb3IiOiJhc3NpZ25lZF90b0lOeyQuc3lzX3VzZXIuc3lzX2lkfSIsInJlbGF0ZWRFbnRpdHlDdXJzb3IiOnsiY3Vyc29yIjoiY3Vyc29yMSIsImNvbGxlY3Rpb25DdXJzb3IiOiJjb2xsZWN0aW9uMSJ9fX0=",
			wantErr:    nil,
		},
		"success_with_nil_input_cursor": {
			inputCursor: nil,
			wantCursor:  "",
			wantErr:     nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotCursor, gotErr := servicenow.MarshalAdvancedFilterCursor(tt.inputCursor)

			assert.Equal(t, tt.wantCursor, gotCursor)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func TestUnmarshalAdvancedFilterCursor(t *testing.T) {
	tests := map[string]struct {
		inputCursor string
		wantCursor  *servicenow.AdvancedFilterCursor
		wantErr     *framework.Error
	}{
		"success_with_implicit_filter_only": {
			inputCursor: "eyJpbXBsaWNpdEZpbHRlckN1cnNvciI6eyJlbnRpdHlGaWx0ZXJJbmRleCI6MCwibWVtYmVyRmlsdGVySW5kZXgiOjAsImN1cnNvciI6eyJjdXJzb3IiOiJjdXJzb3IxIiwiY29sbGVjdGlvbkN1cnNvciI6ImNvbGxlY3Rpb24xIn19fQ==",
			wantCursor: &servicenow.AdvancedFilterCursor{
				ImplicitFilterCursor: &servicenow.ImplicitFilterCursor{
					EntityFilterIndex: 0,
					MemberFilterIndex: 0,
					Cursor: &pagination.CompositeCursor[string]{
						Cursor:           testutil.GenPtr("cursor1"),
						CollectionCursor: testutil.GenPtr("collection1"),
					},
				},
			},
			wantErr: nil,
		},
		"success_with_related_filter_only": {
			inputCursor: "eyJyZWxhdGVkRmlsdGVyQ3Vyc29yIjp7ImVudGl0eUluZGV4IjowLCJlbnRpdHlDdXJzb3IiOiJhc3NpZ25lZF90b0lOeyQuc3lzX3VzZXIuc3lzX2lkfSIsInJlbGF0ZWRFbnRpdHlDdXJzb3IiOnsiY3Vyc29yIjoiY3Vyc29yMSIsImNvbGxlY3Rpb25DdXJzb3IiOiJjb2xsZWN0aW9uMSJ9fX0=",
			wantCursor: &servicenow.AdvancedFilterCursor{
				RelatedFilterCursor: &servicenow.RelatedFilterCursor{
					EntityIndex:  0,
					EntityCursor: testutil.GenPtr("assigned_toIN{$.sys_user.sys_id}"),
					RelatedEntityCursor: &pagination.CompositeCursor[string]{
						Cursor:           testutil.GenPtr("cursor1"),
						CollectionCursor: testutil.GenPtr("collection1"),
					},
				},
			},
			wantErr: nil,
		},
		"success_with_empty_input_cursor": {
			inputCursor: "",
			wantCursor:  &servicenow.AdvancedFilterCursor{},
			wantErr:     nil,
		},
		"invalid_base64": {
			inputCursor: "-",
			wantCursor:  nil,
			wantErr: &framework.Error{
				Message: "Failed to decode base64 cursor: illegal base64 data at input byte 0.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"invalid_json": {
			inputCursor: "eyJrZXkiOiB2YWx1ZQ==", // `{key: value` is invalid JSON.
			wantCursor:  nil,
			wantErr: &framework.Error{
				Message: `Failed to unmarshal JSON cursor: invalid character 'v' looking for beginning of value.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotCursor, gotErr := servicenow.UnmarshalAdvancedFilterCursor(tt.inputCursor)

			assert.Equal(t, tt.wantCursor, gotCursor)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func TestPopulateNextImplicitFilterCursor(t *testing.T) {
	tests := map[string]struct {
		currentCursor      servicenow.ImplicitFilterCursor
		implicitFilters    []servicenow.EntityFilter
		responseNextCursor *pagination.CompositeCursor[string]
		wantNextCursor     *servicenow.ImplicitFilterCursor
	}{
		"first_page_moving_to_next_entity_member": {
			currentCursor: servicenow.ImplicitFilterCursor{},
			implicitFilters: []servicenow.EntityFilter{
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id = '123'",
					Members: []servicenow.MemberFilter{
						{
							MemberEntity:       "User",
							MemberEntityFilter: "id = '456",
						},
						{
							MemberEntity:       "GroupMember",
							MemberEntityFilter: "id = '456",
						},
					},
				},
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id = '456'",
					Members: []servicenow.MemberFilter{
						{
							MemberEntity:       "User",
							MemberEntityFilter: "id = '456",
						},
					},
				},
			},
			responseNextCursor: nil,
			wantNextCursor: &servicenow.ImplicitFilterCursor{
				EntityFilterIndex: 0,
				MemberFilterIndex: 1,
				Cursor:            nil,
			},
		},
		"first_page_moving_to_next_entity_member_page": {
			currentCursor: servicenow.ImplicitFilterCursor{},
			implicitFilters: []servicenow.EntityFilter{
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id = '123'",
					Members: []servicenow.MemberFilter{
						{
							MemberEntity:       "User",
							MemberEntityFilter: "id = '456",
						},
						{
							MemberEntity:       "GroupMember",
							MemberEntityFilter: "id = '456",
						},
					},
				},
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id = '456'",
					Members: []servicenow.MemberFilter{
						{
							MemberEntity:       "User",
							MemberEntityFilter: "id = '456",
						},
					},
				},
			},
			responseNextCursor: &pagination.CompositeCursor[string]{
				Cursor:           testutil.GenPtr("cursor1"),
				CollectionCursor: testutil.GenPtr("collection1"),
			},
			wantNextCursor: &servicenow.ImplicitFilterCursor{
				EntityFilterIndex: 0,
				MemberFilterIndex: 0, // We stay on the current member page because there is a next page to extract more of the current member.
				Cursor:            &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("cursor1"), CollectionCursor: testutil.GenPtr("collection1")},
			},
		},
		"first_page_moving_to_next_entity": {
			currentCursor: servicenow.ImplicitFilterCursor{},
			implicitFilters: []servicenow.EntityFilter{
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id = '123'",
					Members: []servicenow.MemberFilter{
						{
							MemberEntity:       "User",
							MemberEntityFilter: "id = '456",
						},
					},
				},
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id = '456'",
					Members: []servicenow.MemberFilter{
						{
							MemberEntity:       "User",
							MemberEntityFilter: "id = '456",
						},
					},
				},
			},
			responseNextCursor: nil,
			wantNextCursor: &servicenow.ImplicitFilterCursor{
				EntityFilterIndex: 1,
				MemberFilterIndex: 0,
				Cursor:            nil,
			},
		},
		"first_page_no_more_entities": {
			currentCursor: servicenow.ImplicitFilterCursor{},
			implicitFilters: []servicenow.EntityFilter{
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id = '123'",
					Members: []servicenow.MemberFilter{
						{
							MemberEntity:       "User",
							MemberEntityFilter: "id = '456",
						},
					},
				},
			},
			responseNextCursor: nil,
			wantNextCursor:     nil,
		},
		"middle_page_moving_to_next_entity": {
			currentCursor: servicenow.ImplicitFilterCursor{
				EntityFilterIndex: 1,
				MemberFilterIndex: 0,
			},
			implicitFilters: []servicenow.EntityFilter{
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id = '123'",
					Members: []servicenow.MemberFilter{
						{
							MemberEntity:       "User",
							MemberEntityFilter: "id = '456",
						},
					},
				},
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id = '123'",
					Members: []servicenow.MemberFilter{
						{
							MemberEntity:       "User",
							MemberEntityFilter: "id = '456",
						},
					},
				},
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id = '123'",
					Members: []servicenow.MemberFilter{
						{
							MemberEntity:       "User",
							MemberEntityFilter: "id = '456",
						},
					},
				},
			},
			responseNextCursor: nil,
			wantNextCursor: &servicenow.ImplicitFilterCursor{
				EntityFilterIndex: 2,
				MemberFilterIndex: 0,
				Cursor:            nil,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotCursor := servicenow.PopulateNextImplicitFilterCursor(
				tt.currentCursor,
				tt.implicitFilters,
				tt.responseNextCursor,
			)

			assert.Equal(t, tt.wantNextCursor, gotCursor)
		})
	}
}

func TestPopulateNextRelatedFilterCursor(t *testing.T) {
	tests := map[string]struct {
		currentCursor           servicenow.RelatedFilterCursor
		nextEntityCursor        *string
		relatedEntityNextCursor *pagination.CompositeCursor[string]
		relatedFilters          []servicenow.EntityAndRelatedEntityFilter
		wantNextCursor          *servicenow.RelatedFilterCursor
	}{
		"first_page_moving_to_next_entity": {
			currentCursor:           servicenow.RelatedFilterCursor{},
			nextEntityCursor:        nil,
			relatedEntityNextCursor: nil,
			relatedFilters: []servicenow.EntityAndRelatedEntityFilter{
				{
					Entity:       "Case",
					EntityFilter: "assigned_toIN{$.sys_user.sys_id}",
					RelatedEntity: servicenow.EntityFilter{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id = '456",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "User",
								MemberEntityFilter: "id = '456",
							},
						},
					},
				},
				{
					Entity:       "Incident",
					EntityFilter: "assigned_toIN{$.sys_user.sys_id}",
					RelatedEntity: servicenow.EntityFilter{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id = '456",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "User",
								MemberEntityFilter: "id = '456",
							},
						},
					},
				},
			},
			wantNextCursor: &servicenow.RelatedFilterCursor{
				EntityIndex:         1,
				EntityCursor:        nil,
				RelatedEntityCursor: nil,
			},
		},
		"first_page_moving_to_current_entity_next_page": {
			currentCursor:           servicenow.RelatedFilterCursor{},
			nextEntityCursor:        testutil.GenPtr("2"),
			relatedEntityNextCursor: nil,
			relatedFilters: []servicenow.EntityAndRelatedEntityFilter{
				{
					Entity:       "Case",
					EntityFilter: "assigned_toIN{$.sys_user.sys_id}",
					RelatedEntity: servicenow.EntityFilter{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id = '456",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "User",
								MemberEntityFilter: "id = '456",
							},
						},
					},
				},
				{
					Entity:       "Incident",
					EntityFilter: "assigned_toIN{$.sys_user.sys_id}",
					RelatedEntity: servicenow.EntityFilter{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id = '456",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "User",
								MemberEntityFilter: "id = '456",
							},
						},
					},
				},
			},
			wantNextCursor: &servicenow.RelatedFilterCursor{
				EntityIndex:         0,
				EntityCursor:        testutil.GenPtr("2"),
				RelatedEntityCursor: nil,
			},
		},
		"first_page_moving_to_current_entity_related_entity_next_page": {
			currentCursor:           servicenow.RelatedFilterCursor{},
			nextEntityCursor:        nil,
			relatedEntityNextCursor: &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("cursor1"), CollectionCursor: testutil.GenPtr("collection1")},
			relatedFilters: []servicenow.EntityAndRelatedEntityFilter{
				{
					Entity:       "Case",
					EntityFilter: "assigned_toIN{$.sys_user.sys_id}",
					RelatedEntity: servicenow.EntityFilter{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id = '456",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "User",
								MemberEntityFilter: "id = '456",
							},
						},
					},
				},
				{
					Entity:       "Incident",
					EntityFilter: "assigned_toIN{$.sys_user.sys_id}",
					RelatedEntity: servicenow.EntityFilter{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id = '456",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "User",
								MemberEntityFilter: "id = '456",
							},
						},
					},
				},
			},
			wantNextCursor: &servicenow.RelatedFilterCursor{
				EntityIndex:         0,
				EntityCursor:        nil,
				RelatedEntityCursor: &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("cursor1"), CollectionCursor: testutil.GenPtr("collection1")},
			},
		},
		"first_page_no_more_entities": {
			currentCursor:           servicenow.RelatedFilterCursor{},
			nextEntityCursor:        nil,
			relatedEntityNextCursor: nil,
			relatedFilters: []servicenow.EntityAndRelatedEntityFilter{
				{
					Entity:       "Case",
					EntityFilter: "assigned_toIN{$.sys_user.sys_id}",
					RelatedEntity: servicenow.EntityFilter{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id = '456",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "User",
								MemberEntityFilter: "id = '456",
							},
						},
					},
				},
			},
			wantNextCursor: nil,
		},
		"middle_page_moving_to_related_entity_next_page": {
			currentCursor: servicenow.RelatedFilterCursor{
				EntityIndex:         1,
				RelatedEntityCursor: &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("cursor1"), CollectionCursor: testutil.GenPtr("collection1")},
			},
			nextEntityCursor:        nil,
			relatedEntityNextCursor: &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("cursor2"), CollectionCursor: testutil.GenPtr("collection1")},
			relatedFilters: []servicenow.EntityAndRelatedEntityFilter{
				{
					Entity:       "Case",
					EntityFilter: "assigned_toIN{$.sys_user.sys_id}",
					RelatedEntity: servicenow.EntityFilter{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id = '456",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "User",
								MemberEntityFilter: "id = '456",
							},
						},
					},
				},
			},
			wantNextCursor: &servicenow.RelatedFilterCursor{
				EntityIndex:         1,
				EntityCursor:        nil,
				RelatedEntityCursor: &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("cursor2"), CollectionCursor: testutil.GenPtr("collection1")},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotCursor := servicenow.PopulateNextRelatedFilterCursor(
				tt.currentCursor,
				tt.nextEntityCursor,
				tt.relatedEntityNextCursor,
				tt.relatedFilters,
			)

			assert.Equal(t, tt.wantNextCursor, gotCursor)
		})
	}
}
