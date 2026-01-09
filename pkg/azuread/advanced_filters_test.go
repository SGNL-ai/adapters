// Copyright 2026 SGNL.ai, Inc.
package azuread

import (
	"testing"

	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPopulateNextAdvancedFilterCursor(t *testing.T) {
	tests := []struct {
		name                     string
		currAdvancedFilterCursor AdvancedFilterCursor
		advancedFilters          []EntityFilter
		responseNextCursor       *pagination.CompositeCursor[string]
		wantAdvancedFilterCursor *AdvancedFilterCursor
	}{
		{
			name: "next_page_exists_for_current_member_entity",
			currAdvancedFilterCursor: AdvancedFilterCursor{
				EntityFilterIndex: 0,
				MemberFilterIndex: 0,
				Cursor:            nil,
			},
			advancedFilters: []EntityFilter{
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id eq 'group1'",
					Members: []MemberFilter{
						{MemberEntity: "User"},
					},
				},
			},
			responseNextCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("dummyCursor"),
			},
			wantAdvancedFilterCursor: &AdvancedFilterCursor{
				// These fields are not changed
				EntityFilterIndex: 0,
				MemberFilterIndex: 0,

				// Cursor is updated
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("dummyCursor"),
				},
			},
		},
		{
			name: "next_page_does_not_exist__want_next_member_entity",
			currAdvancedFilterCursor: AdvancedFilterCursor{
				EntityFilterIndex: 0,
				MemberFilterIndex: 0,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("dummyCursor"),
				},
			},
			advancedFilters: []EntityFilter{
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id eq 'group1'",
					Members: []MemberFilter{
						{MemberEntity: "User"},
						{MemberEntity: "Group"},
					},
				},
			},
			responseNextCursor: nil,
			wantAdvancedFilterCursor: &AdvancedFilterCursor{
				// This field is not changed
				EntityFilterIndex: 0,

				// This field is incremented
				MemberFilterIndex: 1,

				// Cursor is reset
				Cursor: nil,
			},
		},
		{
			name: "next_page_and_next_member_entity_does_not_exist__want_next_filter",
			currAdvancedFilterCursor: AdvancedFilterCursor{
				EntityFilterIndex: 0,
				MemberFilterIndex: 1,
				Cursor:            nil,
			},
			advancedFilters: []EntityFilter{
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id eq 'group1'",
					Members: []MemberFilter{
						{MemberEntity: "User"},
						{MemberEntity: "Group"},
					},
				},
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id eq 'group2'",
					Members: []MemberFilter{
						{MemberEntity: "User"},
					},
				},
			},
			responseNextCursor: nil,
			wantAdvancedFilterCursor: &AdvancedFilterCursor{
				// This field is incremented
				EntityFilterIndex: 1,

				// These fields are reset
				MemberFilterIndex: 0,
				Cursor:            nil,
			},
		},
		{
			name: "sync_complete",
			currAdvancedFilterCursor: AdvancedFilterCursor{
				EntityFilterIndex: 1,
				MemberFilterIndex: 0,
				Cursor:            nil,
			},
			advancedFilters: []EntityFilter{
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id eq 'group1'",
					Members: []MemberFilter{
						{MemberEntity: "User"},
					},
				},
				{
					ScopeEntity:       "Group",
					ScopeEntityFilter: "id eq 'group2'",
					Members: []MemberFilter{
						{MemberEntity: "User"},
					},
				},
			},
			responseNextCursor:       nil, // indicates that no next page exists for (scopeEntity, memberEntity) pair.
			wantAdvancedFilterCursor: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := populateNextAdvancedFilterCursor(tt.currAdvancedFilterCursor, tt.advancedFilters, tt.responseNextCursor)
			assert.Equal(t, tt.wantAdvancedFilterCursor, result)
		})
	}
}

func TestExtractImplicitFilters(t *testing.T) {
	tests := map[string]struct {
		advancedFilters AdvancedFilters
		wantFilters     map[string][]EntityFilter
	}{
		"group_member_filter_with_only_group_filters": {
			advancedFilters: AdvancedFilters{
				ScopedObjects: map[string][]EntityFilter{
					GroupMember: {
						{
							ScopeEntity:       "Group",
							ScopeEntityFilter: "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
							Members: []MemberFilter{
								{
									MemberEntity:       "Group",
									MemberEntityFilter: "startswith(displayName, 'California')",
								},
							},
						},
					},
				},
			},
			wantFilters: map[string][]EntityFilter{
				Group: {
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
					},
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
						Members: []MemberFilter{
							{
								MemberEntity:       "Group",
								MemberEntityFilter: "startswith(displayName, 'California')",
							},
						},
					},
				},
			},
		},
		"group_member_filter_with_group_and_user_filters": {
			advancedFilters: AdvancedFilters{
				ScopedObjects: map[string][]EntityFilter{
					GroupMember: {
						{
							ScopeEntity:       "Group",
							ScopeEntityFilter: "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
							Members: []MemberFilter{
								{
									MemberEntity:       "Group",
									MemberEntityFilter: "startswith(displayName, 'California')",
								},
								{
									MemberEntity:       "User",
									MemberEntityFilter: "department eq 'engineering'",
								},
								{
									MemberEntity:       "User",
									MemberEntityFilter: "department eq 'product'",
								},
							},
						},
					},
				},
			},
			wantFilters: map[string][]EntityFilter{
				User: {
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
						Members: []MemberFilter{
							{
								MemberEntity:       "User",
								MemberEntityFilter: "department eq 'engineering'",
							},
							{
								MemberEntity:       "User",
								MemberEntityFilter: "department eq 'product'",
							},
						},
					},
				},
				Group: {
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
					},
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
						Members: []MemberFilter{
							{
								MemberEntity:       "Group",
								MemberEntityFilter: "startswith(displayName, 'California')",
							},
						},
					},
				},
			},
		},
		"group_member_filter_with_multiple_group_and_user_filters": {
			advancedFilters: AdvancedFilters{
				ScopedObjects: map[string][]EntityFilter{
					GroupMember: {
						{
							ScopeEntity:       "Group",
							ScopeEntityFilter: "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
							Members: []MemberFilter{
								{
									MemberEntity:       "User",
									MemberEntityFilter: "department eq 'engineering'",
								},
							},
						},
						{
							ScopeEntity:       "Group",
							ScopeEntityFilter: "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
							Members: []MemberFilter{
								{
									MemberEntity:       "Group",
									MemberEntityFilter: "startswith(displayName, 'California')",
								},
								{
									MemberEntity:       "User",
									MemberEntityFilter: "department eq 'engineering'",
								},
								{
									MemberEntity:       "User",
									MemberEntityFilter: "department eq 'product'",
								},
							},
						},
						{
							ScopeEntity:       "Group",
							ScopeEntityFilter: "startswith(displayName, 'California')",
							Members: []MemberFilter{
								{
									MemberEntity: "User",
								},
							},
						},
						{
							ScopeEntity:       "Group",
							ScopeEntityFilter: "id in ('fc765ab6-40bc-4eec-99cb-1fe45369c44e')",
							Members: []MemberFilter{
								{
									MemberEntity:       "User",
									MemberEntityFilter: "department eq 'product'",
								},
							},
						},
						{
							ScopeEntity:       "Group",
							ScopeEntityFilter: "id in ('f7d73513-9c19-4b6b-bbc6-cbcad7881081')",
							Members: []MemberFilter{
								{
									MemberEntity: "User",
								},
							},
						},
					},
				},
			},
			wantFilters: map[string][]EntityFilter{
				User: {
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
						Members: []MemberFilter{
							{
								MemberEntity:       "User",
								MemberEntityFilter: "department eq 'engineering'",
							},
						},
					},
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
						Members: []MemberFilter{
							{
								MemberEntity:       "User",
								MemberEntityFilter: "department eq 'engineering'",
							},
							{
								MemberEntity:       "User",
								MemberEntityFilter: "department eq 'product'",
							},
						},
					},
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "startswith(displayName, 'California')",
						Members: []MemberFilter{
							{
								MemberEntity: "User",
							},
						},
					},
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id in ('fc765ab6-40bc-4eec-99cb-1fe45369c44e')",
						Members: []MemberFilter{
							{
								MemberEntity:       "User",
								MemberEntityFilter: "department eq 'product'",
							},
						},
					},
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id in ('f7d73513-9c19-4b6b-bbc6-cbcad7881081')",
						Members: []MemberFilter{
							{
								MemberEntity: "User",
							},
						},
					},
				},
				Group: {
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
					},
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
					},
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')",
						Members: []MemberFilter{
							{
								MemberEntity:       "Group",
								MemberEntityFilter: "startswith(displayName, 'California')",
							},
						},
					},
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "startswith(displayName, 'California')",
					},
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id in ('fc765ab6-40bc-4eec-99cb-1fe45369c44e')",
					},
					{
						ScopeEntity:       "Group",
						ScopeEntityFilter: "id in ('f7d73513-9c19-4b6b-bbc6-cbcad7881081')",
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotFilters := ExtractImplicitFilters(tt.advancedFilters)
			assert.Equal(t, tt.wantFilters, gotFilters)
		})
	}
}
