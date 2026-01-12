// Copyright 2025 SGNL.ai, Inc.
package servicenow_test

import (
	"reflect"
	"testing"

	"github.com/sgnl-ai/adapters/pkg/servicenow"
)

func TestExtractImplicitFilters(t *testing.T) {
	tests := map[string]struct {
		advancedFilters servicenow.AdvancedFilters
		wantFilters     map[string][]servicenow.EntityFilter
	}{
		"group_filter_with_single_group_with_user_and_group_members": {
			advancedFilters: servicenow.AdvancedFilters{
				ScopedObjects: map[string][]servicenow.EntityFilter{
					servicenow.Group: {
						{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=California",
							Members: []servicenow.MemberFilter{
								{
									MemberEntity:       "sys_user",
									MemberEntityFilter: "user.active=true",
								},
								{
									MemberEntity:       "sys_user_group",
									MemberEntityFilter: "active=true",
								},
							},
						},
					},
				},
			},
			wantFilters: map[string][]servicenow.EntityFilter{
				servicenow.Group: {
					{
						ScopeEntity:       "sys_user_group",
						ScopeEntityFilter: "name=California",
					},
					{
						ScopeEntity:       "sys_user_group",
						ScopeEntityFilter: "name=California",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "sys_user_group",
								MemberEntityFilter: "active=true",
							},
						},
					},
				},
				servicenow.User: {
					{
						ScopeEntity:       "sys_user_group",
						ScopeEntityFilter: "name=California",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "sys_user",
								MemberEntityFilter: "user.active=true",
							},
						},
					},
				},
				servicenow.GroupMember: {
					{
						ScopeEntity:       "sys_user_group",
						ScopeEntityFilter: "name=California",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "sys_user",
								MemberEntityFilter: "user.active=true",
							},
						},
					},
				},
			},
		},
		"group_filter_with_multiple_group_with_user_and_group_members": {
			advancedFilters: servicenow.AdvancedFilters{
				ScopedObjects: map[string][]servicenow.EntityFilter{
					servicenow.Group: {
						{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=California",
							Members: []servicenow.MemberFilter{
								{
									MemberEntity:       "sys_user",
									MemberEntityFilter: "user.active=true",
								},
							},
						},
						{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=Texas",
							Members: []servicenow.MemberFilter{
								{
									MemberEntity:       "sys_user",
									MemberEntityFilter: "user.active=true",
								},
							},
						},
					},
				},
			},
			wantFilters: map[string][]servicenow.EntityFilter{
				servicenow.Group: {
					{
						ScopeEntity:       "sys_user_group",
						ScopeEntityFilter: "name=California",
					},
					{
						ScopeEntity:       "sys_user_group",
						ScopeEntityFilter: "name=Texas",
					},
				},
				servicenow.User: {
					{
						ScopeEntity:       "sys_user_group",
						ScopeEntityFilter: "name=California",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "sys_user",
								MemberEntityFilter: "user.active=true",
							},
						},
					},
					{
						ScopeEntity:       "sys_user_group",
						ScopeEntityFilter: "name=Texas",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "sys_user",
								MemberEntityFilter: "user.active=true",
							},
						},
					},
				},
				servicenow.GroupMember: {
					{
						ScopeEntity:       "sys_user_group",
						ScopeEntityFilter: "name=California",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "sys_user",
								MemberEntityFilter: "user.active=true",
							},
						},
					},
					{
						ScopeEntity:       "sys_user_group",
						ScopeEntityFilter: "name=Texas",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "sys_user",
								MemberEntityFilter: "user.active=true",
							},
						},
					},
				},
			},
		},
		"group_filter_with_related_filters_are_ignored": {
			advancedFilters: servicenow.AdvancedFilters{
				ScopedObjects: map[string][]servicenow.EntityFilter{
					servicenow.Group: {
						{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=California",
							Members: []servicenow.MemberFilter{
								{
									MemberEntity:       "sys_user",
									MemberEntityFilter: "user.active=true",
									RelatedEntities: []servicenow.RelatedEntityFilter{
										{
											RelatedEntity:       "change_task",
											RelatedEntityFilter: "assigned_toIN{$.sys_user.sys_id}",
										},
									},
								},
							},
							RelatedEntities: []servicenow.RelatedEntityFilter{
								{
									RelatedEntity:       "change_task",
									RelatedEntityFilter: "assigned_toIN{$.sys_user.sys_id}",
								},
							},
						},
					},
				},
			},
			wantFilters: map[string][]servicenow.EntityFilter{
				servicenow.Group: {
					{
						ScopeEntity:       "sys_user_group",
						ScopeEntityFilter: "name=California",
					},
				},
				servicenow.User: {
					{
						ScopeEntity:       "sys_user_group",
						ScopeEntityFilter: "name=California",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "sys_user",
								MemberEntityFilter: "user.active=true",
							},
						},
					},
				},
				servicenow.GroupMember: {
					{
						ScopeEntity:       "sys_user_group",
						ScopeEntityFilter: "name=California",
						Members: []servicenow.MemberFilter{
							{
								MemberEntity:       "sys_user",
								MemberEntityFilter: "user.active=true",
							},
						},
					},
				},
			},
		},
		"non_group_filter_no_implicit_filters_extracted": {
			advancedFilters: servicenow.AdvancedFilters{
				ScopedObjects: map[string][]servicenow.EntityFilter{
					servicenow.User: {
						{
							ScopeEntity:       "sys_user",
							ScopeEntityFilter: "sys_id=1234",
						},
					},
				},
			},
			wantFilters: map[string][]servicenow.EntityFilter{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotFilters := servicenow.ExtractImplicitFilters(tt.advancedFilters)
			if !reflect.DeepEqual(gotFilters, tt.wantFilters) {
				t.Errorf("gotFilters: %v, wantFilters: %v", gotFilters, tt.wantFilters)
			}
		})
	}
}

func TestExtractRelatedFilters(t *testing.T) {
	tests := map[string]struct {
		advancedFilters servicenow.AdvancedFilters
		wantFilters     map[string][]servicenow.EntityAndRelatedEntityFilter
	}{
		"group_filter_with_single_group_with_user_and_related_members": {
			advancedFilters: servicenow.AdvancedFilters{
				ScopedObjects: map[string][]servicenow.EntityFilter{
					servicenow.Group: {
						{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=California",
							Members: []servicenow.MemberFilter{
								{
									MemberEntity:       "sys_user",
									MemberEntityFilter: "user.active=true",
									RelatedEntities: []servicenow.RelatedEntityFilter{
										{
											RelatedEntity:       "change_task",
											RelatedEntityFilter: "assigned_toIN{$.sys_user.sys_id}",
										},
									},
								},
							},
							RelatedEntities: []servicenow.RelatedEntityFilter{
								{
									RelatedEntity:       "change_task",
									RelatedEntityFilter: "assignment_groupIN{$.sys_user_group.sys_id}",
								},
							},
						},
					},
				},
			},
			wantFilters: map[string][]servicenow.EntityAndRelatedEntityFilter{
				servicenow.ChangeTask: {
					{
						Entity:       "change_task",
						EntityFilter: "assigned_toIN{$.sys_user.sys_id}",
						RelatedEntity: servicenow.EntityFilter{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=California",
							Members: []servicenow.MemberFilter{
								{
									MemberEntity:       "sys_user",
									MemberEntityFilter: "user.active=true",
								},
							},
						},
					},
					{
						Entity:       "change_task",
						EntityFilter: "assignment_groupIN{$.sys_user_group.sys_id}",
						RelatedEntity: servicenow.EntityFilter{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=California",
						},
					},
				},
			},
		},
		"group_filter_with_multiple_groups_with_user_and_related_members": {
			advancedFilters: servicenow.AdvancedFilters{
				ScopedObjects: map[string][]servicenow.EntityFilter{
					servicenow.Group: {
						{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=California",
							Members: []servicenow.MemberFilter{
								{
									MemberEntity:       "sys_user",
									MemberEntityFilter: "user.active=true",
									RelatedEntities: []servicenow.RelatedEntityFilter{
										{
											RelatedEntity:       "change_task",
											RelatedEntityFilter: "assigned_toIN{$.sys_user.sys_id}",
										},
										{
											RelatedEntity:       "change_request",
											RelatedEntityFilter: "assigned_toIN{$.sys_user.sys_id}",
										},
									},
								},
							},
							RelatedEntities: []servicenow.RelatedEntityFilter{
								{
									RelatedEntity:       "change_task",
									RelatedEntityFilter: "assignment_groupIN{$.sys_user_group.sys_id}",
								},
								{
									RelatedEntity:       "change_request",
									RelatedEntityFilter: "assignment_groupIN{$.sys_user_group.sys_id}",
								},
							},
						},
						{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=Texas",
							RelatedEntities: []servicenow.RelatedEntityFilter{
								{
									RelatedEntity:       "change_request",
									RelatedEntityFilter: "assignment_groupIN{$.sys_user_group.sys_id}",
								},
							},
						},
					},
				},
			},
			wantFilters: map[string][]servicenow.EntityAndRelatedEntityFilter{
				servicenow.ChangeTask: {
					{
						Entity:       "change_task",
						EntityFilter: "assigned_toIN{$.sys_user.sys_id}",
						RelatedEntity: servicenow.EntityFilter{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=California",
							Members: []servicenow.MemberFilter{
								{
									MemberEntity:       "sys_user",
									MemberEntityFilter: "user.active=true",
								},
							},
						},
					},
					{
						Entity:       "change_task",
						EntityFilter: "assignment_groupIN{$.sys_user_group.sys_id}",
						RelatedEntity: servicenow.EntityFilter{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=California",
						},
					},
				},
				servicenow.ChangeRequest: {
					{
						Entity:       "change_request",
						EntityFilter: "assigned_toIN{$.sys_user.sys_id}",
						RelatedEntity: servicenow.EntityFilter{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=California",
							Members: []servicenow.MemberFilter{
								{
									MemberEntity:       "sys_user",
									MemberEntityFilter: "user.active=true",
								},
							},
						},
					},
					{
						Entity:       "change_request",
						EntityFilter: "assignment_groupIN{$.sys_user_group.sys_id}",
						RelatedEntity: servicenow.EntityFilter{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=California",
						},
					},
					{
						Entity:       "change_request",
						EntityFilter: "assignment_groupIN{$.sys_user_group.sys_id}",
						RelatedEntity: servicenow.EntityFilter{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=Texas",
						},
					},
				},
			},
		},
		"group_filter_with_no_related_entities": {
			advancedFilters: servicenow.AdvancedFilters{
				ScopedObjects: map[string][]servicenow.EntityFilter{
					servicenow.Group: {
						{
							ScopeEntity:       "sys_user_group",
							ScopeEntityFilter: "name=California",
							Members: []servicenow.MemberFilter{
								{
									MemberEntity:       "sys_user",
									MemberEntityFilter: "user.active=true",
								},
							},
						},
					},
				},
			},
			wantFilters: map[string][]servicenow.EntityAndRelatedEntityFilter{},
		},
		"non_group_filter_no_related_filters_extracted": {
			advancedFilters: servicenow.AdvancedFilters{
				ScopedObjects: map[string][]servicenow.EntityFilter{
					servicenow.User: {
						{
							ScopeEntity:       "sys_user",
							ScopeEntityFilter: "sys_id=1234",
						},
					},
				},
			},
			wantFilters: map[string][]servicenow.EntityAndRelatedEntityFilter{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotFilters := servicenow.ExtractRelatedFilters(tt.advancedFilters)
			if !reflect.DeepEqual(gotFilters, tt.wantFilters) {
				t.Errorf("gotFilters: %v, wantFilters: %v", gotFilters, tt.wantFilters)
			}
		})
	}
}

func TestExtractEntityAndAttributeFromJSONPathString(t *testing.T) {
	tests := map[string]struct {
		input      string
		wantEntity string
		wantAttr   string
	}{
		"single_json_path": {
			input:      "assigned_toIN{$.sys_user.sys_id}",
			wantEntity: "sys_user",
			wantAttr:   "sys_id",
		},
		"multiple_json_paths_expect_only_first_to_be_parsed": {
			input:      "assigned_toIN{$.sys_user.sys_id}^priorityIN{$.priority}",
			wantEntity: "sys_user",
			wantAttr:   "sys_id",
		},
		"no_json_path_expect_empty_strings": {
			input:      "assigned_toIN1234",
			wantEntity: "",
			wantAttr:   "",
		},
		"input_empty_string_expect_empty_strings": {
			input:      "",
			wantEntity: "",
			wantAttr:   "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			entity, attr := servicenow.ExtractEntityAndAttributeFromJSONPathString(tt.input)

			if !reflect.DeepEqual(entity, tt.wantEntity) {
				t.Errorf("entity: %v, wantEntity: %v", entity, tt.wantEntity)
			}

			if !reflect.DeepEqual(attr, tt.wantAttr) {
				t.Errorf("attr: %v, wantAttr: %v", attr, tt.wantAttr)
			}
		})
	}
}

func TestReplaceEntityAndAttributeInString(t *testing.T) {
	tests := map[string]struct {
		input              string
		newString          string
		wantReplacedString string
	}{
		"single_json_path": {
			input:              "assigned_toIN{$.sys_user.sys_id}",
			newString:          "1234,5678",
			wantReplacedString: "assigned_toIN1234,5678",
		},
		"no_json_path_in_string_expect_no_change": {
			input:              "assigned_toIN1234",
			newString:          "1234,5678",
			wantReplacedString: "assigned_toIN1234",
		},
		"empty_string": {
			input:              "",
			newString:          "1234,5678",
			wantReplacedString: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			replacedString := servicenow.ReplaceEntityAndAttributeInString(tt.input, tt.newString)

			if !reflect.DeepEqual(replacedString, tt.wantReplacedString) {
				t.Errorf("replacedString: %v, wantReplacedString: %v", replacedString, tt.wantReplacedString)
			}
		})
	}
}
