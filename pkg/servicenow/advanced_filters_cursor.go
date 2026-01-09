// Copyright 2026 SGNL.ai, Inc.
package servicenow

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

type AdvancedFilterCursor struct {
	ImplicitFilterCursor *ImplicitFilterCursor `json:"implicitFilterCursor,omitempty"`
	RelatedFilterCursor  *RelatedFilterCursor  `json:"relatedFilterCursor,omitempty"`
}

// nolint: lll
// Example of ImplictFilter:
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
//	    ]
//	}
//
// ImplicitFilterCursor.EntityFilterIndex keeps track of items in the "sys_user" array.
// ImplicitFilterCursor.MemberFilterIndex keeps track of members in the "sys_user"[EntityFilterIndex].members array.
// ImplicitFilterCursor.Cursor keeps track of the cursor used for "sys_user"[EntityFilterIndex].members[MemberFilterIndex].
// ImplicitFilterCursor is semantically equivalent to the AdvancedFilterCursor in the azuread package.
// It's been renamed in this package to improve clarity in the context of ServiceNow which
// behaves differently than AzureAD.
type ImplicitFilterCursor struct {
	EntityFilterIndex int                                 `json:"entityFilterIndex"`
	MemberFilterIndex int                                 `json:"memberFilterIndex"`
	Cursor            *pagination.CompositeCursor[string] `json:"cursor"`
}

// Example of a RelatedFilter:
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
// RelatedFilterCursor.EntityIndex keeps track of items in the "change_task" array.
// RelatedFilterCursor.EntityCursor keeps track of the ServiceNow cursor used for "change_task"[EntityIndex].
// RelatedFilterCursor.RelatedEntityCursor keeps track of the cursor used for "change_task"[EntityIndex].relatedEntity.
type RelatedFilterCursor struct {
	EntityIndex         int                                 `json:"entityIndex"`
	EntityCursor        *string                             `json:"entityCursor,omitempty"`
	RelatedEntityCursor *pagination.CompositeCursor[string] `json:"relatedEntityCursor,omitempty"`
}

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

// TODO [sc-36825]: This is a copy paste of the logic in the azuread package.
func PopulateNextImplicitFilterCursor(
	currImplicitFilterCursor ImplicitFilterCursor,
	implicitFilters []EntityFilter,
	responseNextCursor *pagination.CompositeCursor[string],
) *ImplicitFilterCursor {
	// Default value assumes next page exists for this (scopedEntity, memberEntity) pair, continue paginating...
	nextAdvancedFilterCursor := &ImplicitFilterCursor{
		EntityFilterIndex: currImplicitFilterCursor.EntityFilterIndex,
		MemberFilterIndex: currImplicitFilterCursor.MemberFilterIndex,
		Cursor:            responseNextCursor,
	}

	memberFilters := implicitFilters[currImplicitFilterCursor.EntityFilterIndex].Members

	// if there is no next page for this (scopedEntity, memberEntity) pair
	// 1. move to the next memberEntity for that scopedEntity if it exists.
	// 2. if there is no next memberEntity, move to the next scopedEntity if it exists.
	// 3. If no more scopedEntity exists, sync is complete.
	if responseNextCursor == nil {
		nextAdvancedFilterCursor.Cursor = nil // reset cursor for the next (scopedEntity, memberEntity) pair

		// if there is no next memberEntity, move to the next scopedEntity
		nextAdvancedFilterCursor.MemberFilterIndex = currImplicitFilterCursor.MemberFilterIndex + 1
		if nextAdvancedFilterCursor.MemberFilterIndex >= len(memberFilters) {
			nextAdvancedFilterCursor.MemberFilterIndex = 0 // reset memberFilterIndex for a new scopedEntity

			// if there is no next scopedEntity, sync is complete.
			nextAdvancedFilterCursor.EntityFilterIndex = currImplicitFilterCursor.EntityFilterIndex + 1
			if nextAdvancedFilterCursor.EntityFilterIndex >= len(implicitFilters) {
				nextAdvancedFilterCursor = nil
			}
		}
	}

	return nextAdvancedFilterCursor
}

func PopulateNextRelatedFilterCursor(
	currRelatedFilterCursor RelatedFilterCursor,
	nextEntityCursor *string,
	relatedEntityNextCursor *pagination.CompositeCursor[string],
	relatedFilters []EntityAndRelatedEntityFilter,
) *RelatedFilterCursor {
	nextRelatedFilterCursor := &RelatedFilterCursor{
		EntityIndex:         currRelatedFilterCursor.EntityIndex,
		EntityCursor:        currRelatedFilterCursor.EntityCursor,
		RelatedEntityCursor: currRelatedFilterCursor.RelatedEntityCursor,
	}

	// If there are more pages for the current entity in the (entity, relatedEntity) pair, continue.
	if nextEntityCursor != nil {
		nextRelatedFilterCursor.EntityCursor = nextEntityCursor

		return nextRelatedFilterCursor
	}

	// If there are no more entities, but there are more pages for relatedEntity in the (entity, relatedEntity) pair,
	// move to the next relatedEntity page.
	if relatedEntityNextCursor != nil {
		nextRelatedFilterCursor.RelatedEntityCursor = relatedEntityNextCursor
		nextRelatedFilterCursor.EntityCursor = nil

		return nextRelatedFilterCursor
	}

	// If there are no more pages for the current (entity, relatedEntity) pair, move onto the next pair.
	if currRelatedFilterCursor.EntityIndex < len(relatedFilters)-1 {
		nextRelatedFilterCursor.EntityIndex++
		nextRelatedFilterCursor.EntityCursor = nil
		nextRelatedFilterCursor.RelatedEntityCursor = nil

		return nextRelatedFilterCursor
	}

	// Otherwise, we're done syncing.
	return nil
}
