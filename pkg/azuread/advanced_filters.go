// Copyright 2025 SGNL.ai, Inc.
package azuread

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// See docs: https://learn.microsoft.com/en-us/graph/aad-advanced-queries

// nolint: godot,lll
// Example of AdvancedFilters:
/*
"advancedFilters": {
	"getObjectsByScope": {
		// The externalId is a key
		// The value is a list of filter configurations
		"GroupMember": [
			{
				// To obtain a subset of groupMembers, we scope the entity `Group` and apply a scope-filter.
				// Observe that scopeEntityFilter is extensible because
				// - native AAD filter syntax and not limited to uniqueIds
				// - allows applying same user-filter/group-filter to several groups matching the scope-filter
				"scopeEntity": "Group",
				"scopeEntityFilter": "id in ('7df6bf7d-7b09-4399-9aed-0e345d1ea7b2', 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', '84168a05-7f94-46f5-8d2f-302f6c2ca638')",

				// Fetch both users and groups belonging to the groups matching the scope-filter
				// The set of users and groups are further filtered down using memberEntityFilter
				"members": [
					{
						// All Users belonging to group matching scope-filters will be further filtered down using memberEntityFilter
						"memberEntity": "User",
						"memberEntityFilter": "department eq 'Architecture and Authentication'"
					},
					{
						// All Groups belonging to group matching scope-filters will NOT be further filtered down as memberEntityFilter is absent
						"memberEntity": "Group"
					}
				]
			},
			{
				"scopeEntity": "Group",
				"scopeEntityFilter": "id eq '94e98d95-8e04-47ee-b2d4-5a1bd96bf0a9'",
				// Only fetch users belonging to the groups matching the scope-filter
				"members": [
					{
						"memberEntity": "User"
					}
				]
			},
			{
				"scopeEntity": "Group",
				"scopeEntityFilter": "id eq '821d9845-3482-41c8-9a0e-d93f4d576731'",

				// Fetch both users and groups belonging to the groups matching the scope-filter
				"members": [
					{
						"memberEntity": "User"
					},
					{
						"memberEntity": "Group"
					}
				]
			}
		],
		// This ingests only those group entity nodes starting with `Test` or `SGNL`
		"Group": [
			{
				"scopeEntity": "Group",
				"scopeEntityFilter": "startswith(displayName, 'Test')"
			},
			{
				"scopeEntity": "Group",
				"scopeEntityFilter": "startswith(displayName, 'SGNL')"
			}
		]
	}
}
*/
type AdvancedFilters struct {
	// ScopedObjects is a map of entityExternalID to a list of EntityFilter configurations.
	ScopedObjects map[string][]EntityFilter `json:"getObjectsByScope,omitempty"`
}

type EntityFilter struct {
	ScopeEntity       string         `json:"scopeEntity"`
	ScopeEntityFilter string         `json:"scopeEntityFilter"`
	Members           []MemberFilter `json:"members"`
}

type MemberFilter struct {
	MemberEntity       string `json:"memberEntity"`
	MemberEntityFilter string `json:"memberEntityFilter,omitempty"`
}

/*
AdvancedFilterCursor is a simple wrapper around the existing composite cursor.
The existing composite cursor helps paginate through G groups and fetch U users.

When using advanced filters for an entityExternalID (e.g. `groupMembers“),

 1. Iterate each scopedEntity (i-th index) (e.g. a `group`) and get its filter, a scopeFilter
 2. For that scopedEntity, iterate each memberEntity (j-th index) (e.g. `users`), and get its filter, a memberFilter
 3. For every (scopedEntity, memberEntity) pair,
    3.1 the scopeFilter would filter down the /groups response to a subset of G groups
    3.2 the memberFilter would filter down the /members response to a subset of U users for each group
    3.3 Call GetPage with these filters
    3.3 The pagination within GetPage is captured by the existing composite cursor
    3.4 If composite cursor is not nil, it means there are more pages for the current (scopedEntity, memberEntity) pair.
    3.5 Move to the next (scopedEntity, memberEntity) pair if the composite cursor is nil.
 4. Continue calling GetPage until all (scopedEntity, memberEntity) pairs are processed

This is how the advanced filter cursor wraps the existing composite cursor.
*/
type AdvancedFilterCursor struct {
	EntityFilterIndex int                                 `json:"entityFilterIndex"`
	MemberFilterIndex int                                 `json:"memberFilterIndex"`
	Cursor            *pagination.CompositeCursor[string] `json:"cursor"`
}

var (
	advancedFilterEntityHasMembers = map[string]bool{
		GroupMember: true,
	}

	allowedScopeEntities = map[string]bool{
		Group: true,
	}

	// Valid member entity external ID to endpoint suffix mapping.
	// With advanced filters, a customer can choose to apply filters individually to different types of members.
	// Every top level key here must also be present in `advancedFilterEntityHasMembers`.
	// If no suffix is applicable, set the value to an empty string.
	memberEntityToEndpointSuffix = map[string]map[string]string{
		// A group member can either be a user or a group.
		GroupMember: {
			User:  "/microsoft.graph.user",
			Group: "/microsoft.graph.group",
		},

		// TODO: Extend this for other member entities when required e.g. RoleMember -> Role
	}
)

// MarshalAdvancedFilterCursor marshals the struct and b64 encodes it.
func MarshalAdvancedFilterCursor(cursor *AdvancedFilterCursor) (string, *framework.Error) {
	if cursor == nil {
		return "", nil
	}

	nextCursorBytes, marshalErr := json.Marshal(cursor)
	if marshalErr != nil {
		return "", &framework.Error{
			Message: fmt.Sprintf("Failed to marshal advanced filter cursor into JSON: %v.", marshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return base64.StdEncoding.EncodeToString(nextCursorBytes), nil
}

// UnmarshalAdvancedFilterCursor decodes the b64 encoded string and unmarshals it.
func UnmarshalAdvancedFilterCursor(cursor string) (*AdvancedFilterCursor, *framework.Error) {
	var advancedFilterCursor = &AdvancedFilterCursor{}
	if cursor == "" {
		return advancedFilterCursor, nil
	}

	advancedFilterCursorBytes, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to decode base64 cursor: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	unmarshalErr := json.Unmarshal(advancedFilterCursorBytes, advancedFilterCursor)
	if unmarshalErr != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal JSON cursor: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	return advancedFilterCursor, nil
}

// ValidateAdvancedFilterCursor validates the advanced filter cursor based on the current request configuration.
func validateAdvancedFilterCursor(
	advancedFilterCursor AdvancedFilterCursor,
	advancedFilters []EntityFilter,
	entityExternalID string,
) *framework.Error {
	if advancedFilterCursor.EntityFilterIndex >= len(advancedFilters) {
		return &framework.Error{
			Message: fmt.Sprintf(
				"Invalid filter index for %v: %v.",
				entityExternalID,
				advancedFilterCursor.EntityFilterIndex,
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	if hasMembers := advancedFilterEntityHasMembers[entityExternalID]; !hasMembers {
		return nil
	}

	if advancedFilterCursor.MemberFilterIndex >= len(advancedFilters[advancedFilterCursor.EntityFilterIndex].Members) {
		return &framework.Error{
			Message: fmt.Sprintf(
				"Invalid member filter index for %v: %v.",
				entityExternalID,
				advancedFilterCursor.MemberFilterIndex,
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	return nil
}

// populateNextAdvancedFilterCursor populates the next advanced filter cursor based on
// the current AdvancedFilterCursor, the next compositeCursor in the GetPage response, and
// the advanced filter cursor based on the current request configuration.
func populateNextAdvancedFilterCursor(
	currAdvancedFilterCursor AdvancedFilterCursor,
	advancedFilters []EntityFilter,
	responseNextCursor *pagination.CompositeCursor[string],
) *AdvancedFilterCursor {
	// Default value assumes next page exists for this (scopedEntity, memberEntity) pair, continue paginating...
	nextAdvancedFilterCursor := &AdvancedFilterCursor{
		EntityFilterIndex: currAdvancedFilterCursor.EntityFilterIndex,
		MemberFilterIndex: currAdvancedFilterCursor.MemberFilterIndex,
		Cursor:            responseNextCursor,
	}

	memberFilters := advancedFilters[currAdvancedFilterCursor.EntityFilterIndex].Members

	// if there is no next page for this (scopedEntity, memberEntity) pair
	// 1. move to the next memberEntity for that scopedEntity if it exists.
	// 2. if there is no next memberEntity, move to the next scopedEntity if it exists.
	// 3. If no more scopedEntity exists, sync is complete.
	if responseNextCursor == nil {
		nextAdvancedFilterCursor.Cursor = nil // reset cursor for the next (scopedEntity, memberEntity) pair

		// if there is no next memberEntity, move to the next scopedEntity
		nextAdvancedFilterCursor.MemberFilterIndex = currAdvancedFilterCursor.MemberFilterIndex + 1
		if nextAdvancedFilterCursor.MemberFilterIndex >= len(memberFilters) {
			nextAdvancedFilterCursor.MemberFilterIndex = 0 // reset memberFilterIndex for a new scopedEntity

			// if there is no next scopedEntity, sync is complete.
			nextAdvancedFilterCursor.EntityFilterIndex = currAdvancedFilterCursor.EntityFilterIndex + 1
			if nextAdvancedFilterCursor.EntityFilterIndex >= len(advancedFilters) {
				nextAdvancedFilterCursor = nil
			}
		}
	}

	return nextAdvancedFilterCursor
}

func validateAdvancedFilterConfiguration(request *framework.Request[Config]) *framework.Error {
	if request.Config.AdvancedFilters == nil {
		return nil
	}

	if len(request.Config.AdvancedFilters.ScopedObjects) == 0 {
		return &framework.Error{
			Message: "advancedFilters.getObjectsByScope cannot be empty.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	for entityExternalID, filters := range request.Config.AdvancedFilters.ScopedObjects {
		if _, found := advancedFilterEntityHasMembers[entityExternalID]; !found {
			return &framework.Error{
				Message: fmt.Sprintf(
					"Advanced Filters on advancedFilters.getObjectsByScope.%v is not supported.",
					entityExternalID,
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			}
		}

		if request.Config.Filters != nil {
			if _, exists := request.Config.Filters[entityExternalID]; exists {
				return &framework.Error{
					Message: fmt.Sprintf(
						"Only one of advancedFilters.getObjectsByScope.%v OR filters.%v is allowed.",
						entityExternalID,
						entityExternalID,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				}
			}
		}

		if len(filters) == 0 {
			return &framework.Error{
				Message: fmt.Sprintf(
					"advancedFilters.getObjectsByScope.%v must have at least one filter defined.",
					entityExternalID,
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			}
		}

		for idx, filter := range filters {
			if filter.ScopeEntity == "" {
				return &framework.Error{
					Message: fmt.Sprintf(
						"advancedFilters.getObjectsByScope.%v.[%d].scopeEntity cannot be empty.",
						entityExternalID,
						idx,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				}
			}

			if _, found := allowedScopeEntities[filter.ScopeEntity]; !found {
				return &framework.Error{
					Message: fmt.Sprintf(
						"advancedFilters.getObjectsByScope.%v.[%d].scopeEntity is not supported.",
						entityExternalID,
						idx,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				}
			}

			if filter.ScopeEntityFilter == "" {
				return &framework.Error{
					Message: fmt.Sprintf(
						"advancedFilters.getObjectsByScope.%v.[%d].scopeEntityFilter cannot be empty.",
						entityExternalID,
						idx,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				}
			}

			hasMember, found := advancedFilterEntityHasMembers[entityExternalID]
			if !found {
				return &framework.Error{
					Message: fmt.Sprintf(
						"advancedFilters.getObjectsByScope.%v.[%d].scopeEntity is not supported.",
						entityExternalID,
						idx,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				}
			}

			if !hasMember {
				continue
			}

			if len(filter.Members) == 0 {
				return &framework.Error{
					Message: fmt.Sprintf(
						"advancedFilters.getObjectsByScope.%v.[%d].members cannot be empty.",
						entityExternalID,
						idx,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				}
			}

			for memberIdx, member := range filter.Members {
				if member.MemberEntity == "" {
					return &framework.Error{
						Message: fmt.Sprintf(
							"advancedFilters.getObjectsByScope.%v.[%d].members[%d].memberEntity cannot be empty.",
							entityExternalID,
							idx,
							memberIdx,
						),
						Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
					}
				}

				// A missing entry corresponding to incorrect endpoint generation.
				if _, found := memberEntityToEndpointSuffix[entityExternalID][member.MemberEntity]; !found {
					return &framework.Error{
						Message: fmt.Sprintf(
							"advancedFilters.getObjectsByScope.%v.[%d].members[%d].memberEntity is not supported.",
							entityExternalID,
							idx,
							memberIdx,
						),
						Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
					}
				}
			}
		}
	}

	return nil
}

// ExtractImplicitFilters extracts implicit User/Group filters from the GroupMember's advanced filters.
// A GroupMember advanced filter may look like:
//
// "GroupMember": [
//
//	{
//	  "scopeEntity": "Group",
//	  "scopeEntityFilter": "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')", // Canada Group
//	  "members": [
//		{
//		  "memberEntity": "User",
//		  "memberEntityFilter": "department eq 'engineering'"
//		}
//	  ]
//	},
//	{
//	  "scopeEntity": "Group",
//	  "scopeEntityFilter": "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')", // United States Group
//	  "members": [
//		{
//		  "memberEntity": "User",
//		  "memberEntityFilter": "department eq 'engineering'"
//		},
//		{
//		  "memberEntity": "User",
//		  "memberEntityFilter": "department eq 'product'"
//		},
//		{
//		  "memberEntity": "Group",
//		  "memberEntityFilter": "startswith(displayName, 'California')"
//		}
//	  ]
//	},
//
// ]
// The above filter means:
// 1) Sync all GroupMember user nodes from the Canada group where the department is engineering.
// 2) Sync all GroupMember user nodes from the United States group where the department is engineering or product.
// 3) Sync all GroupMember group nodes from the United States group where the group name starts with California.
// This implies we should only sync:
// 1) User nodes from the Canada group where the department is engineering.
// 2) User nodes from the United States group where the department is engineering or product.
// 3) Group nodes from the United States group where the group name starts with California.
// 4) The Canada and United States groups themselves.
// Thus, this function will extract the following implicit filters (only User filters are shown for brevity)
// "User": [
//
//	{
//	  "scopeEntity": "Group",
//	  "scopeEntityFilter": "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')", // Canada Group
//	  "members": [
//		{
//		  "memberEntity": "User",
//		  "memberEntityFilter": "department eq 'engineering'"
//		}
//	  ]
//	},
//	{
//	  "scopeEntity": "Group",
//	  "scopeEntityFilter": "id in ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')", // United States Group
//	  "members": [
//		{
//		  "memberEntity": "User",
//		  "memberEntityFilter": "department eq 'engineering'"
//		},
//		{
//		  "memberEntity": "User",
//		  "memberEntityFilter": "department eq 'product'"
//		},
//	  ]
//	},
//
// ].
func ExtractImplicitFilters(advancedFilters AdvancedFilters) map[string][]EntityFilter {
	implicitFilters := make(map[string][]EntityFilter, 2)

	// Implicit filters are extracted from GroupMember filters.
	groupMemberFilters := advancedFilters.ScopedObjects[GroupMember]

	for _, groupMemberFilter := range groupMemberFilters {
		// Each GroupMember filter may or may not contain user and group filters.
		// That's what we are attempting to extract.
		userFilters, groupFilters := []MemberFilter{}, []MemberFilter{}

		// We extract the user and group filters from the `members` array.
		// See the function doc for an example of a GroupMember filter.
		for _, memberFilter := range groupMemberFilter.Members {
			switch memberFilter.MemberEntity {
			case User:
				userFilters = append(userFilters, memberFilter)
			case Group:
				groupFilters = append(groupFilters, memberFilter)
			}
		}

		// If we've found user filters, we add them to the implicit filters.
		// This is essentially just copying the example above and putting it under the `Users` key
		// in the implicit filters map.
		if len(userFilters) > 0 {
			implicitFilters[User] = append(implicitFilters[User], EntityFilter{
				ScopeEntity:       groupMemberFilter.ScopeEntity,
				ScopeEntityFilter: groupMemberFilter.ScopeEntityFilter,
				Members:           userFilters,
			})
		}

		// For group filters, we need to create two implicit filters:
		// 1) For the parent group itself.
		// 2) For the child groups of that parent.
		// This ensures we return both of those groups.
		parentGroupFilter := EntityFilter{
			ScopeEntity:       groupMemberFilter.ScopeEntity,
			ScopeEntityFilter: groupMemberFilter.ScopeEntityFilter,
		}

		if parentGroupFilter.ScopeEntity == Group {
			implicitFilters[Group] = append(implicitFilters[Group], parentGroupFilter)
		}

		if len(groupFilters) > 0 {
			implicitFilters[Group] = append(implicitFilters[Group], EntityFilter{
				ScopeEntity:       groupMemberFilter.ScopeEntity,
				ScopeEntityFilter: groupMemberFilter.ScopeEntityFilter,
				Members:           groupFilters,
			})
		}
	}

	return implicitFilters
}
