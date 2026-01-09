// Copyright 2026 SGNL.ai, Inc.
package servicenow

import (
	"regexp"
)

// ServiceNow advanced filters are similar in syntax to Azure AD advanced filters
// but they have different implementations and semantics.
// Example of ServiceNow advanced filter with new syntax commented:
//
//	{
//		"requestTimeoutSeconds": 10,
//		"apiVersion": "v2",
//		"advancedFilters": {
//		  "getObjectsByScope": {
//			"sys_user_group": [
//			  {
//				"scopeEntity": "sys_user_group",
//				"scopeEntityFilter": "sys_id=0270c251c3200200be647bfaa2d3aea6",
//				"members": [
//				  {
//					"memberEntity": "sys_user",
//					"memberEntityFilter": "user.active=true",
//					"relatedEntities": [ <---------- NEW
//					  {
//						"relatedEntity": "change_task",
//						"relatedEntityFilter": "assigned_toIN{$.sys_user.sys_id}" <--- NEW SYNTAX
//					  }
//					]
//				  }
//				],
//				"relatedEntities": [ <------------- NEW
//				  {
//					"relatedEntity": "change_task",
//					"relatedEntityFilter": "assignment_groupIN{$.sys_group.sys_id}" <--- NEW SYNTAX
//				  }
//				]
//			  }
//			]
//		  }
//		}
//	}
//
// There are two new components in the ServiceNow advanced filter syntax:
// 1) The `relatedEntities` array.
// 2) The JSONPath syntax in the `relatedEntityFilter` key.
// Related entities are related to the `memberEntity` key via the relatedEntityFilter.
// In the example above, we are filtering `change_task` entities by the `assigned_to` field of the `sys_user` entity
// which are `user.active=true` in the `sys_user_group` entity with the `sys_id` of `0270c251c3200200be647bfaa2d3aea6`.
// Advanced filters are a **client facing** concept that we internally parse into two separate filters:
// Implicit Filters and Related Filters.
// Each of those filters are explained in their respective functions below.

// This regex extracts the two components in a JSONPath string, e.g. `{$.sys_user.sys_id}` --> sys_user, sys_id.
const JSONPathTemplateStringRegex = `\{\$\.([a-zA-Z_]+)\.([a-zA-Z_]+)\}`

var (
	SupportedScopeEntity      = Group
	SupportedImplicitEntities = []string{User, Group, GroupMember}
	SupportedRelatedEntities  = []string{Case, Incident, ChangeRequest, ChangeTask}

	JSONPathStringRegex = regexp.MustCompile(JSONPathTemplateStringRegex)
)

type AdvancedFilters struct {
	ScopedObjects map[string][]EntityFilter `json:"getObjectsByScope,omitempty"`
}

type EntityFilter struct {
	ScopeEntity       string                `json:"scopeEntity"`
	ScopeEntityFilter string                `json:"scopeEntityFilter"`
	Members           []MemberFilter        `json:"members,omitempty"`
	RelatedEntities   []RelatedEntityFilter `json:"relatedEntities,omitempty"`
}

type MemberFilter struct {
	MemberEntity       string                `json:"memberEntity"`
	MemberEntityFilter string                `json:"memberEntityFilter,omitempty"`
	RelatedEntities    []RelatedEntityFilter `json:"relatedEntities,omitempty"`
}

type RelatedEntityFilter struct {
	RelatedEntity       string `json:"relatedEntity"`
	RelatedEntityFilter string `json:"relatedEntityFilter,omitempty"`
}

type EntityAndRelatedEntityFilter struct {
	Entity        string       `json:"entity"`
	EntityFilter  string       `json:"entityFilter"`
	RelatedEntity EntityFilter `json:"relatedEntity"`
}

// ExtractImplicitFilters extracts the implicit entity filters from the advanced filters.
// This is best illustrated with an example.
// Given the following advanced filter:
//
//	{
//		"requestTimeoutSeconds": 10,
//		"apiVersion": "v2",
//		"advancedFilters": {
//		  "getObjectsByScope": {
//			"sys_user_group": [
//			  {
//				"scopeEntity": "sys_user_group",
//				"scopeEntityFilter": "sys_id=0270c251c3200200be647bfaa2d3aea6",
//				"members": [
//				  {
//					"memberEntity": "sys_user",
//					"memberEntityFilter": "user.active=true"
//				  }
//				]
//			  }
//			]
//		  }
//		}
//	}
//
// From a client's perspective, this means:
// "Sync all groups with ID `0270c251c3200200be647bfaa2d3aea6` and all users in those groups that are active."
// Therefore, the adapter must retrieve 1 group node and all active user nodes that are members of that group.
// We generate the following implicit filters:
//
//	{
//	    "sys_user": [
//	      {
//	        "scopeEntity": "sys_user_group",
//	        "scopeEntityFilter": "sys_id=0270c251c3200200be647bfaa2d3aea6",
//	        "members": [
//	          {
//	            "memberEntity": "sys_user",
//	            "memberEntityFilter": "user.active=true"
//	          }
//	        ]
//	      }
//	    ],
//	    "sys_user_grmember": [
//	      {
//	        "scopeEntity": "sys_user_group",
//	        "scopeEntityFilter": "sys_id=0270c251c3200200be647bfaa2d3aea6",
//	        "members": [
//	          {
//	            "memberEntity": "sys_user",
//	            "memberEntityFilter": "user.active=true"
//	          }
//	        ]
//	      }
//	    ],
//	    "sys_user_group": [
//	      {
//	        "scopeEntity": "sys_user_group",
//	        "scopeEntityFilter": "sys_id=0270c251c3200200be647bfaa2d3aea6",
//	        "members": null
//	      }
//	    ]
//	}
//
// A `sys_user_grmember` is filter is also generated since it makes logical sense.
// TODO [sc-36825]: This is a copy paste of the logic in the azuread package with small modifications.
func ExtractImplicitFilters(advancedFilters AdvancedFilters) map[string][]EntityFilter {
	implicitFilters := make(map[string][]EntityFilter, len(SupportedImplicitEntities))

	filters := advancedFilters.ScopedObjects[SupportedScopeEntity]

	for _, filter := range filters {
		userFilters, groupFilters, groupMemberFilters := []MemberFilter{}, []MemberFilter{}, []MemberFilter{}

		for _, memberFilter := range filter.Members {
			// Strip any RelatedEntityFilters from the MemberFilter since we don't need that.
			memberFilter := MemberFilter{
				MemberEntity:       memberFilter.MemberEntity,
				MemberEntityFilter: memberFilter.MemberEntityFilter,
			}

			switch memberFilter.MemberEntity {
			case User:
				userFilters = append(userFilters, memberFilter)
				groupMemberFilters = append(groupMemberFilters, memberFilter)
			case Group:
				groupFilters = append(groupFilters, memberFilter)
			}
		}

		if len(userFilters) > 0 {
			implicitFilters[User] = append(implicitFilters[User], EntityFilter{
				ScopeEntity:       filter.ScopeEntity,
				ScopeEntityFilter: filter.ScopeEntityFilter,
				Members:           userFilters,
			})

			implicitFilters[GroupMember] = append(implicitFilters[GroupMember], EntityFilter{
				ScopeEntity:       filter.ScopeEntity,
				ScopeEntityFilter: filter.ScopeEntityFilter,
				Members:           groupMemberFilters,
			})
		}

		parentGroupFilter := EntityFilter{
			ScopeEntity:       filter.ScopeEntity,
			ScopeEntityFilter: filter.ScopeEntityFilter,
		}

		if parentGroupFilter.ScopeEntity == Group {
			implicitFilters[Group] = append(implicitFilters[Group], parentGroupFilter)
		}

		if len(groupFilters) > 0 {
			implicitFilters[Group] = append(implicitFilters[Group], EntityFilter{
				ScopeEntity:       filter.ScopeEntity,
				ScopeEntityFilter: filter.ScopeEntityFilter,
				Members:           groupFilters,
			})
		}
	}

	return implicitFilters
}

// ExtractRelatedFilters extracts the related entity filters from the advanced filters.
// This is best illustrated with an example.
// Given the following advanced filter:
//
//	{
//		"requestTimeoutSeconds": 10,
//		"apiVersion": "v2",
//		"advancedFilters": {
//		  "getObjectsByScope": {
//			"sys_user_group": [
//			  {
//				"scopeEntity": "sys_user_group",
//				"scopeEntityFilter": "sys_id=0270c251c3200200be647bfaa2d3aea6",
//				"members": [
//				  {
//					"memberEntity": "sys_user",
//					"memberEntityFilter": "user.active=true",
//					"relatedEntities": [
//					  {
//						"relatedEntity": "change_task",
//						"relatedEntityFilter": "assigned_toIN{$.sys_user.sys_id}"
//					  }
//					]
//				  }
//				]
//			  }
//			]
//		  }
//		}
//	}
//
// From a client's perspective, this means:
// "Sync all groups with ID `0270c251c3200200be647bfaa2d3aea6` and all users in those groups that are active.
// For all users in that group,
// sync all change tasks where the `assigned_to` field is equal to the `sys_id` of the user."
// We generate the following related filters:
//
//	{
//	    "change_task": [
//	      {
//	        "entity": "change_task",
//	        "entityFilter": "assigned_toIN{$.sys_user.sys_id}",
//	        "relatedEntity": {
//	          "scopeEntity": "sys_user_group",
//	          "scopeEntityFilter": "sys_id=0270c251c3200200be647bfaa2d3aea6",
//	          "members": [
//	            {
//	              "memberEntity": "sys_user",
//	              "memberEntityFilter": "user.active=true"
//	            }
//	          ]
//	        }
//	      }
//	    ]
//	}
//
// The related entity describes how to retrieve the related entity. In this case, it describes how to
// retrieve the users related to the change tasks.
// If the change tasks were related to groups instead, the `members` key would not be present.
// There can only ever be one member in the `members` array after we extract.
func ExtractRelatedFilters(advancedFilters AdvancedFilters) map[string][]EntityAndRelatedEntityFilter {
	relatedEntities := make(map[string][]EntityAndRelatedEntityFilter)

	advancedFiltersByEntity := advancedFilters.ScopedObjects

	for _, entityFilters := range advancedFiltersByEntity {
		for _, filter := range entityFilters {
			// Extract relatedEntites from member filters.
			for _, memberFilter := range filter.Members {
				for _, relatedEntityFilter := range memberFilter.RelatedEntities {
					switch relatedEntityFilter.RelatedEntity {
					case Case, Incident, ChangeRequest, ChangeTask:
						relatedEntityFilterToAdd := EntityAndRelatedEntityFilter{}
						relatedEntityFilterToAdd.Entity = relatedEntityFilter.RelatedEntity
						relatedEntityFilterToAdd.EntityFilter = relatedEntityFilter.RelatedEntityFilter

						relatedEntityFilterToAdd.RelatedEntity.ScopeEntity = filter.ScopeEntity
						relatedEntityFilterToAdd.RelatedEntity.ScopeEntityFilter = filter.ScopeEntityFilter

						relatedEntityFilterToAdd.RelatedEntity.Members = []MemberFilter{{
							MemberEntity:       memberFilter.MemberEntity,
							MemberEntityFilter: memberFilter.MemberEntityFilter,
						}}

						relatedEntities[relatedEntityFilter.RelatedEntity] = append(
							relatedEntities[relatedEntityFilter.RelatedEntity],
							relatedEntityFilterToAdd,
						)
					}
				}
			}

			// Extract relatedEntities from top level entity filters.
			for _, relatedEntityFilter := range filter.RelatedEntities {
				switch relatedEntityFilter.RelatedEntity {
				case Case, Incident, ChangeRequest, ChangeTask:
					relatedEntityFilterToAdd := EntityAndRelatedEntityFilter{}
					relatedEntityFilterToAdd.Entity = relatedEntityFilter.RelatedEntity
					relatedEntityFilterToAdd.EntityFilter = relatedEntityFilter.RelatedEntityFilter

					relatedEntityFilterToAdd.RelatedEntity.ScopeEntity = filter.ScopeEntity
					relatedEntityFilterToAdd.RelatedEntity.ScopeEntityFilter = filter.ScopeEntityFilter

					relatedEntities[relatedEntityFilter.RelatedEntity] = append(
						relatedEntities[relatedEntityFilter.RelatedEntity],
						relatedEntityFilterToAdd,
					)
				}
			}
		}
	}

	return relatedEntities
}

// ExtractEntityAndAttributeFromJSONPathString extracts the entity and attribute components
// from a string containing JSONPath.
// For example, if `s = "assigned_toIN{$.sys_user.sys_id}"`, then we extract `sys_user`, `sys_id`
// from `{$.sys_user.sys_id}`.
func ExtractEntityAndAttributeFromJSONPathString(jsonPath string) (string, string) {
	matches := JSONPathStringRegex.FindStringSubmatch(jsonPath)

	// match[0] is the full match, match[1] and match[2] are the groups.
	if len(matches) == 3 {
		return matches[1], matches[2]
	}

	return "", ""
}

// ReplaceEntityAndAttributeInString replaces the JSONPath string with a new string.
// For example, if `s = "assigned_toIN{$.sys_user.sys_id}"`, then
// we replace `{$.sys_user.sys_id}` with the input string.
func ReplaceEntityAndAttributeInString(jsonPath, newString string) string {
	return JSONPathStringRegex.ReplaceAllString(jsonPath, newString)
}
