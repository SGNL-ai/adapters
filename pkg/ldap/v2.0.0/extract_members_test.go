// Copyright 2026 SGNL.ai, Inc.

// Unit tests for extractMembersDNFromGroup covering plain attributes,
// range attributes, offset handling, and edge cases.

package ldap

import (
	"fmt"
	"testing"
)

func TestExtractMembersDNFromGroup(t *testing.T) {
	// Helper to create member DN strings as []any for use in test data.
	makeMemberList := func(count int) []any {
		members := make([]any, count)
		for i := range count {
			members[i] = fmt.Sprintf("CN=User%d,OU=Users,DC=example,DC=org", i)
		}

		return members
	}

	// Helper to create expected DN strings.
	makeExpectedDNS := func(start, end int) []string {
		dns := make([]string, end-start)
		for i := start; i < end; i++ {
			dns[i-start] = fmt.Sprintf("CN=User%d,OU=Users,DC=example,DC=org", i)
		}

		return dns
	}

	tests := []struct {
		name           string
		memberObjs     map[string]any
		offset         int
		count          int
		expectedDNS    []string
		expectedAction MemberExtractionAction
	}{
		// Plain attribute: no offset cases.
		{
			name: "plain_attr_all_members_requested",
			memberObjs: map[string]any{
				"member": makeMemberList(5),
			},
			offset:         0,
			count:          5,
			expectedDNS:    makeExpectedDNS(0, 5),
			expectedAction: MemberExtractionActionDone,
		},
		{
			name: "plain_attr_fewer_members_than_count",
			memberObjs: map[string]any{
				"member": makeMemberList(3),
			},
			offset:         0,
			count:          10,
			expectedDNS:    makeExpectedDNS(0, 3),
			expectedAction: MemberExtractionActionDone,
		},
		{
			name: "plain_attr_partial_extraction",
			memberObjs: map[string]any{
				"member": makeMemberList(8),
			},
			offset:         0,
			count:          3,
			expectedDNS:    makeExpectedDNS(0, 3),
			expectedAction: MemberExtractionActionContinueRegular,
		},

		// Plain attribute: with offset (the bug fix scenario).
		{
			name: "plain_attr_offset_skips_first_members",
			memberObjs: map[string]any{
				"member": makeMemberList(8),
			},
			offset:         3,
			count:          3,
			expectedDNS:    makeExpectedDNS(3, 6),
			expectedAction: MemberExtractionActionContinueRegular,
		},
		{
			name: "plain_attr_offset_extracts_remaining_members",
			memberObjs: map[string]any{
				"member": makeMemberList(8),
			},
			offset:         6,
			count:          3,
			expectedDNS:    makeExpectedDNS(6, 8),
			expectedAction: MemberExtractionActionDone,
		},
		{
			name: "plain_attr_offset_exactly_at_end",
			memberObjs: map[string]any{
				"member": makeMemberList(5),
			},
			offset:         5,
			count:          3,
			expectedDNS:    []string{},
			expectedAction: MemberExtractionActionDone,
		},
		{
			name: "plain_attr_offset_past_end",
			memberObjs: map[string]any{
				"member": makeMemberList(5),
			},
			offset:         10,
			count:          3,
			expectedDNS:    []string{},
			expectedAction: MemberExtractionActionDone,
		},
		{
			name: "plain_attr_offset_exact_remaining_count",
			memberObjs: map[string]any{
				"member": makeMemberList(6),
			},
			offset:         3,
			count:          3,
			expectedDNS:    makeExpectedDNS(3, 6),
			expectedAction: MemberExtractionActionDone,
		},

		// Range attribute: offset should NOT be applied (server handles range).
		{
			name: "range_attr_partial_range",
			memberObjs: map[string]any{
				"member;range=0-1499": makeMemberList(5),
			},
			offset:         0,
			count:          3,
			expectedDNS:    makeExpectedDNS(0, 3),
			expectedAction: MemberExtractionActionContinueRange,
		},
		{
			name: "range_attr_all_extracted",
			memberObjs: map[string]any{
				"member;range=0-1499": makeMemberList(5),
			},
			offset:         0,
			count:          5,
			expectedDNS:    makeExpectedDNS(0, 5),
			expectedAction: MemberExtractionActionContinueRange,
		},
		{
			name: "range_attr_final_page_with_star",
			memberObjs: map[string]any{
				"member;range=1500-*": makeMemberList(3),
			},
			offset:         0,
			count:          3,
			expectedDNS:    makeExpectedDNS(0, 3),
			expectedAction: MemberExtractionActionDone,
		},
		{
			name: "range_attr_final_page_partial_count",
			memberObjs: map[string]any{
				"member;range=1500-*": makeMemberList(3),
			},
			offset:         0,
			count:          5,
			expectedDNS:    makeExpectedDNS(0, 3),
			expectedAction: MemberExtractionActionDone,
		},
		{
			name: "range_attr_ignores_offset",
			memberObjs: map[string]any{
				"member;range=500-999": makeMemberList(5),
			},
			offset:         3,
			count:          5,
			expectedDNS:    makeExpectedDNS(0, 5),
			expectedAction: MemberExtractionActionContinueRange,
		},

		// Edge cases.
		{
			name:           "empty_member_objs",
			memberObjs:     map[string]any{},
			offset:         0,
			count:          10,
			expectedDNS:    []string{},
			expectedAction: MemberExtractionActionDone,
		},
		{
			name: "non_list_attribute_skipped",
			memberObjs: map[string]any{
				"distinguishedName": "CN=Group1,DC=example,DC=org",
				"member":            makeMemberList(3),
			},
			offset:         0,
			count:          3,
			expectedDNS:    makeExpectedDNS(0, 3),
			expectedAction: MemberExtractionActionDone,
		},
		{
			name: "non_string_members_skipped",
			memberObjs: map[string]any{
				"member": []any{123, 456},
			},
			offset:         0,
			count:          2,
			expectedDNS:    []string{},
			expectedAction: MemberExtractionActionDone,
		},
		{
			name: "empty_member_list",
			memberObjs: map[string]any{
				"member": []any{},
			},
			offset:         0,
			count:          5,
			expectedDNS:    []string{},
			expectedAction: MemberExtractionActionDone,
		},
		{
			name: "zero_count",
			memberObjs: map[string]any{
				"member": makeMemberList(5),
			},
			offset:         0,
			count:          0,
			expectedDNS:    []string{},
			expectedAction: MemberExtractionActionDone,
		},
		{
			name: "single_member_no_offset",
			memberObjs: map[string]any{
				"member": makeMemberList(1),
			},
			offset:         0,
			count:          1,
			expectedDNS:    makeExpectedDNS(0, 1),
			expectedAction: MemberExtractionActionDone,
		},
		{
			name: "single_member_offset_skips_it",
			memberObjs: map[string]any{
				"member": makeMemberList(1),
			},
			offset:         1,
			count:          1,
			expectedDNS:    []string{},
			expectedAction: MemberExtractionActionDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dns, action := extractMembersDNFromGroup(tt.memberObjs, tt.offset, tt.count)

			if action != tt.expectedAction {
				t.Errorf("expected action %q, got %q", tt.expectedAction, action)
			}

			if len(dns) != len(tt.expectedDNS) {
				t.Fatalf("expected %d DNs, got %d: %v", len(tt.expectedDNS), len(dns), dns)
			}

			for i, dn := range dns {
				if dn != tt.expectedDNS[i] {
					t.Errorf("DN[%d]: expected %q, got %q", i, tt.expectedDNS[i], dn)
				}
			}
		})
	}
}
