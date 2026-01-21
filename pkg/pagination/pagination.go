// Copyright 2026 SGNL.ai, Inc.

package pagination

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	"github.com/sgnl-ai/adapters/pkg/extractor"
)

// CompositeCursor is used to store all required information for pagination.
// This struct is marshaled into JSON and then base-64 encoded.
type CompositeCursor[T int64 | string] struct {
	// Cursor is the cursor that identifies the first object of the page to return.
	//
	// For datasources where this value is used to populate the query parameter to offset the results in a page,
	// this should have an `int64` type. e.g. the "startAt" parameter in Jira or the "offset" parameter in PagerDuty.
	//
	// For datasources where this contains a request URL, this should have a `string` type.
	Cursor *T `json:"cursor,omitempty"`

	// "Collection" is a generic term for a type of entity that can contain other entities.
	// For example, the "teams" entity in PagerDuty can contain "users". It qualifies as a "Collection".
	// The "group" entity in Jira can contain "users". It qualifies as a "Collection".
	// A collection member entity ID is defined per SoR. For Jira, it is "GroupMember". For PagerDuty, it is "members".
	// The following fields are to ONLY be used when the entity is a member entity, otherwise they must be nil.
	// It is up to the adapter to enforce this rule.
	// CollectionID is the ID of the collection when querying for members.
	// Only used when the entity is a member entity, otherwise it must be nil.
	CollectionID *string `json:"collectionId,omitempty"`

	// CollectionCursor is the cursor that identifies the first collection object of the page to return.
	// Only used when the entity is a member entity, otherwise it must be nil.
	CollectionCursor *T `json:"collectionCursor,omitempty"`
}

// UnmarshalCursor unmarshals the cursor from a base64 encoded JSON string.
// If unmarshalling fails, an error is returned.
func UnmarshalCursor[T int64 | string](cursor string) (*CompositeCursor[T], *framework.Error) {
	if cursor == "" {
		return nil, nil
	}

	unmarshaledCursor := &CompositeCursor[T]{}

	cursorBytes, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to decode base64 cursor: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	unmarshalErr := json.Unmarshal(cursorBytes, unmarshaledCursor)
	if unmarshalErr != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal JSON cursor: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	return unmarshaledCursor, nil
}

// MarshalCursor marshals the cursor into a base64 encoded JSON string.
// If marshalling fails, an error is returned.
func MarshalCursor[T int64 | string](cursor *CompositeCursor[T]) (string, *framework.Error) {
	if cursor == nil {
		return "", nil
	}

	nextCursorBytes, marshalErr := json.Marshal(cursor)
	if marshalErr != nil {
		return "", &framework.Error{
			Message: fmt.Sprintf("Failed to marshal cursor into JSON: %v.", marshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return base64.StdEncoding.EncodeToString(nextCursorBytes), nil
}

// ValidateCompositeCursor validates that the composite cursor is correctly set.
// The following rules are enforced:
// 1) If the entity is a member entity, then CollectionID must be set.
// 2) If the entity is not a member entity, then CollectionID and CollectionCursor must be nil.
//
// If no cursor is provided, this returns nil.
func ValidateCompositeCursor[T int64 | string](
	cursor *CompositeCursor[T], entityExternalID string, isMemberEntity bool,
) *framework.Error {
	if cursor == nil {
		return nil
	}

	switch isMemberEntity {
	case true:
		if cursor.CollectionID == nil {
			return &framework.Error{
				Message: fmt.Sprintf("Cursor does not have CollectionID set for entity %s.", entityExternalID),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			}
		}

		return nil
	default:
		if cursor.CollectionID != nil || cursor.CollectionCursor != nil {
			return &framework.Error{
				Message: fmt.Sprintf(
					"Cursor must not contain CollectionID or CollectionCursor fields for entity %s.",
					entityExternalID,
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			}
		}

		return nil
	}
}

// GetNextCursorFromPageSize computes the next cursor based on the number of objects in the current page,
// the current cursor, and the page size.
func GetNextCursorFromPageSize(objectsInPage int, pageSize int64, currentCursor int64) (nextCursor *int64) {
	// Cases:
	// 1) If objectsInPage < pageSize, then this is the last page.
	// 2) If objectsInPage == pageSize, then there MAY be more pages. We cannot determine if this is the
	// last page from this information alone, so we must specify a cursor to get the next page.
	if int64(objectsInPage) == pageSize {
		n := currentCursor + pageSize
		nextCursor = &n
	}

	return
}

// The Link Header is a comma separated list of links with the following format:
// <https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=1>; rel="prev",
// <https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=3>; rel="next",
// <https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=3>; rel="last",
// <https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=1>; rel="first"
// We want to retrieve the "next" link, or return a nil cursor if it is missing to indicate the end of the sync.
func GetNextCursorFromLinkHeader(links []string) *CompositeCursor[string] {
	if cursor := extractor.ValueFromList(links, "https://", ">;rel=\"next\""); cursor != "" {
		return &CompositeCursor[string]{
			Cursor: &cursor,
		}
	}

	return nil
}

// UpdateNextCursorFromCollectionAPI populates the provided cursor in place based on the following rules:
//  1. If `cursor.Cursor` is set, we're in the middle of syncing a page of member entity objects for a collection so we
//     should exit early without modifying the cursor.
//  2. If `cursor.Cursor` is not set, we need to determine the current collection ID to use. In this case, we'll
//     use the provided `getPageFunc` and `collectionRequest` to request the next collection object. If
//     `cursor.CollectionCursor` is not set we're at the start of a sync so we'll request the first collection object.
//     Otherwise, we'll use the `cursor.CollectionCursor` to request the correct collection object.
//  3. If there is a collection object returned, set the `cursor.CollectionID`.
//  4. If there is a cursor to the next collection object, we'll save that as `cursor.CollectionCursor`,
//     which will be used to request the next collection object after we process all pages of member
//     objects for the current collection.
//  5. If the current collection page does not return any objects, we're done with the current sync and should
//     set `cursor.CollectionID` and `cursor.CollectionCursor` to nil and return.
//
// The provided `getPageFunc` should make a request to the collection API based on the provided `collectionRequest`
// and parse / return the expected fields. The `collectionRequest` MUST have a page size set to 1, as well as setting
// any additional params required by the function called in `getPageFunc` to request data from the API.
func UpdateNextCursorFromCollectionAPI[T int64 | string, Request any](
	ctx context.Context,
	cursor *CompositeCursor[T],
	getPageFunc func(
		ctx context.Context, request *Request,
	) (
		int, string, []map[string]any, *CompositeCursor[T], *framework.Error,
	),
	collectionRequest *Request,
	uniqueIDAttribute string,
) (bool, *framework.Error) {
	// If `cursor.Cursor` is set, then we're in the middle of syncing a page of the
	// member entity objects for a specific collection, so this function should exit early.
	if cursor != nil && cursor.Cursor != nil {
		return false, nil
	}

	statusCode, retryAfterHeader, objects, nextCursor, err := getPageFunc(ctx, collectionRequest)
	if err != nil {
		return false, err
	}

	// An adapter error message is generated if the response status code from the
	// collection API is not successful (i.e. if not statusCode >= 200 && statusCode < 300).
	if adapterErr := web.HTTPError(statusCode, retryAfterHeader); adapterErr != nil {
		return false, adapterErr
	}

	switch len(objects) {
	case 0:
		cursor.CollectionID = nil
	case 1:
		if groupID, ok := objects[0][uniqueIDAttribute].(string); ok {
			cursor.CollectionID = &groupID
		} else {
			cursor.CollectionID = nil
		}
	default:
		return false, &framework.Error{
			Message: fmt.Sprintf("Too many collection objects returned in response; expected 1, got %d.", len(objects)),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	cursor.CollectionCursor = nil

	// If there was a next cursor returned with the request to get the current page of collection objects, save that
	// as the CollectionCursor.
	if nextCursor != nil && nextCursor.Cursor != nil {
		cursor.CollectionCursor = nextCursor.Cursor
	}

	// If CollectionID is nil, there are no collection objects included in the current page.
	// If CollectionCursor is also nil, this means we've completed the current sync.
	if cursor.CollectionID == nil && cursor.CollectionCursor == nil {
		// Return true to indicate the sync should complete with no data.
		return true, nil
	}

	return false, nil
}

func (c *CompositeCursor[T]) ParseOffsetValue() (int64, *framework.Error) {
	if c == nil || c.Cursor == nil {
		return 0, nil
	}

	switch v := any(c.Cursor).(type) {
	case *int64:
		return *v, nil
	case *string:
		offsetStr := *v

		offsetInt, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			return 0, &framework.Error{
				Message: fmt.Sprintf("Unable to parse cursor: want valid number, got {%v}.", *c.Cursor),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			}
		}

		return offsetInt, nil
	}

	return 0, &framework.Error{
		Message: "Unable to parse offset value from cursor.",
		Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
	}
}

// PaginateObjects implements pagination for APIs that do not support pagination on the server side.
// It returns a requested page from available records along with a cursor for the next page.
func PaginateObjects[T int64 | string](
	objects []map[string]any,
	pageSize int64,
	cursor *CompositeCursor[T],
) (
	[]map[string]any,
	*T,
	*framework.Error,
) {
	startIndex, err := cursor.ParseOffsetValue()
	if err != nil {
		return nil, nil, err
	}

	numObjects := int64(len(objects))

	// If the start index is 0, we won't require it to be less than the number of objects to account for empty pages.
	// Otherwise, the start index must be less than the number of objects to avoid out of range errors.
	if startIndex != 0 && (startIndex >= numObjects || startIndex < 0) {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("The cursor value: %v, is out of range for number of objects: %v", startIndex, numObjects),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	endIndex := startIndex + pageSize

	// If end index exceeds the length of objects, adjust it to the length of objects
	if endIndex > numObjects {
		endIndex = numObjects
	}

	var nextCursor *T

	// Calculate the next cursor if there are objects remaining past the end index.
	if endIndex < numObjects {
		var ok bool

		nextCursor = new(T)
		switch any(nextCursor).(type) {
		case *int64:
			*nextCursor, ok = any(endIndex).(T)
		case *string:
			*nextCursor, ok = any(strconv.FormatInt(endIndex, 10)).(T)
		}

		if !ok {
			return nil, nil, &framework.Error{
				Message: "Unable to convert next index to cursor type.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			}
		}
	}

	return objects[startIndex:endIndex], nextCursor, nil
}
