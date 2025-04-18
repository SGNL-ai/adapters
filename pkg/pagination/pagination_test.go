// Copyright 2025 SGNL.ai, Inc.

// nolint: lll
package pagination_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestUnmarshalCursorInt64(t *testing.T) {
	tests := map[string]struct {
		inputCursor         string
		wantCompositeCursor *pagination.CompositeCursor[int64]
		wantErr             *framework.Error
	}{
		"empty_cursor": {
			inputCursor:         "",
			wantCompositeCursor: nil,
			wantErr:             nil,
		},
		"valid_cursor": {
			inputCursor: "eyJjdXJzb3IiOjF9", // {"cursor":1}.
			wantCompositeCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(1)),
			},
			wantErr: nil,
		},
		"valid_member_cursor": {
			// {"cursor":1,"collectionId":"1","collectionCursor":1}.
			inputCursor: "eyJjdXJzb3IiOjEsImNvbGxlY3Rpb25JZCI6IjEiLCJjb2xsZWN0aW9uQ3Vyc29yIjoxfQ==",
			wantCompositeCursor: &pagination.CompositeCursor[int64]{
				Cursor:           testutil.GenPtr(int64(1)),
				CollectionID:     testutil.GenPtr("1"),
				CollectionCursor: testutil.GenPtr(int64(1)),
			},
			wantErr: nil,
		},
		"invalid_b64_cursor": {
			inputCursor:         "NOT_B64_ENCODED",
			wantCompositeCursor: nil,
			wantErr: &framework.Error{
				Message: "Failed to decode base64 cursor: illegal base64 data at input byte 3.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"invalid_cursor_struct": {
			inputCursor:         "ImZvbyI=", // "foo" b64 encoded.
			wantCompositeCursor: nil,
			wantErr: &framework.Error{
				Message: "Failed to unmarshal JSON cursor: json: cannot unmarshal string into Go " +
					"value of type pagination.CompositeCursor[int64].",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotCompositeCursor, gotErr := pagination.UnmarshalCursor[int64](tt.inputCursor)

			if !reflect.DeepEqual(gotCompositeCursor, tt.wantCompositeCursor) {
				t.Errorf("gotCompositeCursor: %v, wantCompositeCursor: %v", gotCompositeCursor, tt.wantCompositeCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestUnmarshalCursorString(t *testing.T) {
	tests := map[string]struct {
		inputCursor         string
		wantCompositeCursor *pagination.CompositeCursor[string]
		wantErr             *framework.Error
	}{
		"empty_cursor": {
			inputCursor:         "",
			wantCompositeCursor: nil,
			wantErr:             nil,
		},
		"valid_cursor": {
			// {"cursor":"http://localhost/api/users"}.
			inputCursor: "eyJjdXJzb3IiOiJodHRwOi8vbG9jYWxob3N0L2FwaS91c2VycyJ9",
			wantCompositeCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("http://localhost/api/users"),
			},
			wantErr: nil,
		},
		"valid_member_cursor": {
			// {"cursor":"http://localhost/api/users","collectionId":"1","collectionCursor":"http://localhost/api/groups"}.
			inputCursor: "eyJjdXJzb3IiOiJodHRwOi8vbG9jYWxob3N0L2FwaS91c2VycyIsImNvbGxlY3Rpb25JZCI6IjEiLCJjb2xsZWN0aW9uQ3Vyc29yIjoiaHR0cDovL2xvY2FsaG9zdC9hcGkvZ3JvdXBzIn0=",
			wantCompositeCursor: &pagination.CompositeCursor[string]{
				Cursor:           testutil.GenPtr("http://localhost/api/users"),
				CollectionID:     testutil.GenPtr("1"),
				CollectionCursor: testutil.GenPtr("http://localhost/api/groups"),
			},
			wantErr: nil,
		},
		"invalid_b64_cursor": {
			inputCursor:         "NOT_B64_ENCODED",
			wantCompositeCursor: nil,
			wantErr: &framework.Error{
				Message: "Failed to decode base64 cursor: illegal base64 data at input byte 3.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"invalid_cursor_struct": {
			inputCursor:         "ImZvbyI=", // "foo" b64 encoded.
			wantCompositeCursor: nil,
			wantErr: &framework.Error{
				Message: "Failed to unmarshal JSON cursor: json: cannot unmarshal string into Go " +
					"value of type pagination.CompositeCursor[string].",
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotCompositeCursor, gotErr := pagination.UnmarshalCursor[string](tt.inputCursor)

			if !reflect.DeepEqual(gotCompositeCursor, tt.wantCompositeCursor) {
				t.Errorf("gotCompositeCursor: %v, wantCompositeCursor: %v", gotCompositeCursor, tt.wantCompositeCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestMarshalCursor(t *testing.T) {
	tests := map[string]struct {
		inputCompositeCursor *pagination.CompositeCursor[int64]
		wantB64Cursor        string
		wantErr              *framework.Error
	}{
		"nil_composite_cursor": {
			inputCompositeCursor: nil,
			wantB64Cursor:        "",
			wantErr:              nil,
		},
		"empty_composite_cursor": {
			inputCompositeCursor: &pagination.CompositeCursor[int64]{},
			wantB64Cursor:        "e30=",
			wantErr:              nil,
		},
		"valid_cursor": {
			inputCompositeCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(1)),
			},
			wantB64Cursor: "eyJjdXJzb3IiOjF9", // {"cursor":1}.
			wantErr:       nil,
		},
		"valid_member_cursor": {
			inputCompositeCursor: &pagination.CompositeCursor[int64]{
				Cursor:           testutil.GenPtr(int64(1)),
				CollectionID:     testutil.GenPtr("1"),
				CollectionCursor: testutil.GenPtr(int64(1)),
			},
			// {"cursor":1,"collectionId":"1","collectionCursor":1}.
			wantB64Cursor: "eyJjdXJzb3IiOjEsImNvbGxlY3Rpb25JZCI6IjEiLCJjb2xsZWN0aW9uQ3Vyc29yIjoxfQ==",
			wantErr:       nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotB64Cursor, gotErr := pagination.MarshalCursor(tt.inputCompositeCursor)

			if !reflect.DeepEqual(gotB64Cursor, tt.wantB64Cursor) {
				t.Errorf("gotB64Cursor: %v, wantB64Cursor: %v", gotB64Cursor, tt.wantB64Cursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestValidateCompositeCursor(t *testing.T) {
	tests := map[string]struct {
		inputCompositeCursor  *pagination.CompositeCursor[int64]
		inputEntityExternalID string
		inputIsMemberEntity   bool
		wantErr               *framework.Error
	}{
		"nil_cursor": {
			inputCompositeCursor:  nil,
			inputEntityExternalID: "Member",
			inputIsMemberEntity:   true,
			wantErr:               nil,
		},
		"valid_member_cursor": {
			inputCompositeCursor: &pagination.CompositeCursor[int64]{
				Cursor:           testutil.GenPtr(int64(1)),
				CollectionID:     testutil.GenPtr("1"),
				CollectionCursor: testutil.GenPtr(int64(1)),
			},
			inputEntityExternalID: "Member",
			inputIsMemberEntity:   true,
			wantErr:               nil,
		},
		"valid_non_member_cursor": {
			inputCompositeCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(1)),
			},
			inputEntityExternalID: "User",
			inputIsMemberEntity:   false,
			wantErr:               nil,
		},
		"group_id_missing": {
			inputCompositeCursor: &pagination.CompositeCursor[int64]{
				Cursor:           testutil.GenPtr(int64(1)),
				CollectionCursor: testutil.GenPtr(int64(1)),
			},
			inputEntityExternalID: "Member",
			inputIsMemberEntity:   true,
			wantErr: &framework.Error{
				Message: "Cursor does not have CollectionID set for entity Member.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"group_cursor_field_should_not_be_present": {
			inputCompositeCursor: &pagination.CompositeCursor[int64]{
				Cursor:           testutil.GenPtr(int64(1)),
				CollectionCursor: testutil.GenPtr(int64(1)),
			},
			inputEntityExternalID: "User",
			inputIsMemberEntity:   false,
			wantErr: &framework.Error{
				Message: "Cursor must not contain CollectionID or CollectionCursor fields for entity User.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"group_id_field_should_not_be_present": {
			inputCompositeCursor: &pagination.CompositeCursor[int64]{
				Cursor:       testutil.GenPtr(int64(1)),
				CollectionID: testutil.GenPtr("1"),
			},
			inputEntityExternalID: "User",
			inputIsMemberEntity:   false,
			wantErr: &framework.Error{
				Message: "Cursor must not contain CollectionID or CollectionCursor fields for entity User.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := pagination.ValidateCompositeCursor(
				tt.inputCompositeCursor, tt.inputEntityExternalID, tt.inputIsMemberEntity,
			)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetNextCursorFromPageSize(t *testing.T) {
	tests := map[string]struct {
		inputObjectsInPage int
		inputPageSize      int64
		currentCursor      int64
		wantNextCursor     *int64
	}{
		"last_page": {
			inputObjectsInPage: 1,
			inputPageSize:      10,
			currentCursor:      5,
			wantNextCursor:     nil,
		},
		"unit_page_size": {
			inputObjectsInPage: 1,
			inputPageSize:      1,
			currentCursor:      0,
			// Not possible to know this is the last page if objectsInPage == pageSize,
			// so we specify a next cursor.
			wantNextCursor: testutil.GenPtr(int64(1)),
		},
		"large_page_size": {
			inputObjectsInPage: 1000,
			inputPageSize:      1000,
			currentCursor:      0,
			wantNextCursor:     testutil.GenPtr(int64(1000)),
		},
		"non_zero_start_cursor": {
			inputObjectsInPage: 1,
			inputPageSize:      1,
			currentCursor:      5,
			wantNextCursor:     testutil.GenPtr(int64(6)),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotNextCursor := pagination.GetNextCursorFromPageSize(tt.inputObjectsInPage, tt.inputPageSize, tt.currentCursor)

			if !reflect.DeepEqual(gotNextCursor, tt.wantNextCursor) {
				t.Errorf("gotNextCursor: %v, wantB64Cursor: %v", gotNextCursor, tt.wantNextCursor)
			}
		})
	}
}

func TestGetNextCursorFromLinkHeader(t *testing.T) {
	tests := map[string]struct {
		inputLinkHeader []string
		wantNextCursor  *pagination.CompositeCursor[string]
	}{
		"nil_link_header": {
			inputLinkHeader: nil,
			wantNextCursor:  nil,
		},
		"empty_link_header": {
			inputLinkHeader: []string{},
			wantNextCursor:  nil,
		},
		"valid_link_header": {
			inputLinkHeader: []string{"<https://localhost/api/users?page=2>; rel=\"next\""},
			wantNextCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://localhost/api/users?page=2"),
			},
		},
		"multiple_links_header": {
			inputLinkHeader: []string{"<https://localhost/api/users?page=2>; rel=\"next\"; <https://localhost/api/users?page=1>; rel=\"prev\""},
			wantNextCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://localhost/api/users?page=2"),
			},
		},
		"missing_next_link_header": {
			inputLinkHeader: []string{"<https://localhost/api/users?page=2>; rel=\"first\"; <https://localhost/api/users?page=1>; rel=\"prev\""},
			wantNextCursor:  nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotNextCursor := pagination.GetNextCursorFromLinkHeader(tt.inputLinkHeader)

			if !reflect.DeepEqual(gotNextCursor, tt.wantNextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotNextCursor, tt.wantNextCursor)
			}
		})
	}
}

func TestPaginateObjectsInt64(t *testing.T) {
	tests := map[string]struct {
		objects        []map[string]any
		pageSize       int64
		inputCursor    *pagination.CompositeCursor[int64]
		wantObjects    []map[string]any
		wantNextCursor *int64
		wantErr        *framework.Error
	}{
		"empty_objects": {
			objects:        []map[string]any{},
			pageSize:       3,
			inputCursor:    nil,
			wantObjects:    []map[string]any{},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"empty_objects_with_cursor": {
			objects:  []map[string]any{},
			pageSize: 3,
			inputCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(1)),
			},
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "The cursor value: 1, is out of range for number of objects: 0",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"valid_objects_without_cursor": {
			objects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			pageSize:    3,
			inputCursor: nil,
			wantObjects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
			},
			wantNextCursor: testutil.GenPtr(int64(3)),
			wantErr:        nil,
		},
		"valid_objects_with_cursor": {
			objects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			pageSize: 3,
			inputCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(2)),
			},
			wantObjects: []map[string]any{
				{"id": 3},
				{"id": 4},
				{"id": 5},
			},
			wantNextCursor: testutil.GenPtr(int64(5)),
			wantErr:        nil,
		},
		"page_size_greater_than_num_objects": {
			objects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			pageSize:    10,
			inputCursor: nil,
			wantObjects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"page_size_greater_than_num_objects_with_cursor": {
			objects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			pageSize: 10,
			inputCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr(int64(2)),
			},
			wantObjects: []map[string]any{
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"negative_cursor_value": {
			objects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			pageSize:       10,
			inputCursor:    &pagination.CompositeCursor[int64]{Cursor: testutil.GenPtr(int64(-1))},
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "The cursor value: -1, is out of range for number of objects: 6",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"cursor_greater_than_num_objects": {
			objects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			pageSize:       10,
			inputCursor:    &pagination.CompositeCursor[int64]{Cursor: testutil.GenPtr(int64(15))},
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "The cursor value: 15, is out of range for number of objects: 6",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := pagination.PaginateObjects(tt.objects, tt.pageSize, tt.inputCursor)

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !reflect.DeepEqual(gotNextCursor, tt.wantNextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotNextCursor, tt.wantNextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestPaginateObjectsString(t *testing.T) {
	tests := map[string]struct {
		objects        []map[string]any
		pageSize       int64
		inputCursor    *pagination.CompositeCursor[string]
		wantObjects    []map[string]any
		wantNextCursor *string
		wantErr        *framework.Error
	}{
		"empty_objects": {
			objects:        []map[string]any{},
			pageSize:       3,
			inputCursor:    nil,
			wantObjects:    []map[string]any{},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"empty_objects_with_cursor": {
			objects:  []map[string]any{},
			pageSize: 3,
			inputCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("1"),
			},
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "The cursor value: 1, is out of range for number of objects: 0",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"valid_objects_without_cursor": {
			objects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			pageSize:    3,
			inputCursor: nil,
			wantObjects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
			},
			wantNextCursor: testutil.GenPtr("3"),
			wantErr:        nil,
		},
		"valid_objects_with_cursor": {
			objects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			pageSize: 3,
			inputCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("2"),
			},
			wantObjects: []map[string]any{
				{"id": 3},
				{"id": 4},
				{"id": 5},
			},
			wantNextCursor: testutil.GenPtr("5"),
			wantErr:        nil,
		},
		"page_size_greater_than_num_objects": {
			objects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			pageSize:    10,
			inputCursor: nil,
			wantObjects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"page_size_greater_than_num_objects_with_cursor": {
			objects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			pageSize: 10,
			inputCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("2"),
			},
			wantObjects: []map[string]any{
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"negative_cursor_value": {
			objects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			pageSize:       10,
			inputCursor:    &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("-1")},
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "The cursor value: -1, is out of range for number of objects: 6",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"invalid_cursor_value": {
			objects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			pageSize:       10,
			inputCursor:    &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("random")},
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "unable to parse cursor: want valid number, got {random}",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"cursor_greater_than_num_objects": {
			objects: []map[string]any{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
				{"id": 5},
				{"id": 6},
			},
			pageSize:       10,
			inputCursor:    &pagination.CompositeCursor[string]{Cursor: testutil.GenPtr("15")},
			wantObjects:    nil,
			wantNextCursor: nil,
			wantErr: &framework.Error{
				Message: "The cursor value: 15, is out of range for number of objects: 6",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		// Tests below this are from AzureAD for the Role entity.
		"single_page": {
			objects: []map[string]any{
				{
					"id":              "0fea7f0d-dea1-4028-8ce8-a686ec639d75",
					"deletedDateTime": nil,
					"description":     "Can read basic directory information. Commonly used to grant directory read access to applications and guests.",
					"displayName":     "Directory Readers",
					"roleTemplateId":  "88d8e3e3-c189-46e8-94e1-9b9898b8876b",
				},
				{
					"id":              "18eacdf7-8db3-458d-9099-69fcc2e3cd42",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of Microsoft Entra ID and Microsoft services that use Microsoft Entra identities.",
					"displayName":     "Global Administrator",
					"roleTemplateId":  "62e90394-3621-4004-a7cb-012177145e10",
				},
				{
					"id":              "33a4c989-c3ff-4a77-bf46-ee0acd84476e",
					"deletedDateTime": nil,
					"description":     "Can create application registrations independent of the 'Users can register applications' setting.",
					"displayName":     "Application Developer",
					"roleTemplateId":  "cf1c38e5-69f5-4237-9190-879624dced7c",
				},
				{
					"id":              "321fd63c-c37c-4597-81c4-81e0a93ffb6e",
					"deletedDateTime": nil,
					"description":     "Can manage role assignments in Microsoft Entra ID, and all aspects of Privileged Identity Management.",
					"displayName":     "Privileged Role Administrator",
					"roleTemplateId":  "e8611ab8-8f55-4a1e-953a-60213ab1f814",
				},
				{
					"id":              "8ed2c3eb-1709-444f-bea2-4fb8e2038d6e",
					"deletedDateTime": nil,
					"description":     "Read custom security attribute keys and values for supported Microsoft Entra objects.",
					"displayName":     "Attribute Assignment Reader",
					"roleTemplateId":  "ffd52fa5-98dc-465c-991d-fc073eb59f8f",
				},
				{
					"id":              "e8d9279e-6883-4add-96e8-5f7c8df5637f",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of users and groups, including resetting passwords for limited admins.",
					"displayName":     "User Administrator",
					"roleTemplateId":  "fe930be7-5e62-47db-91af-98c3a49a38b1",
				},
				{
					"id":              "fb96a81c-6147-4cfe-b7fe-c63c2725e7c9",
					"deletedDateTime": nil,
					"description":     "Members of this role can create/manage groups, create/manage groups settings like naming and expiration policies, and view groups activity and audit reports.",
					"displayName":     "Groups Administrator",
					"roleTemplateId":  "fdd7a751-b60b-444a-984c-02652fe8fa1c",
				},
			},
			pageSize: 100,
			wantObjects: []map[string]any{
				{
					"id":              "0fea7f0d-dea1-4028-8ce8-a686ec639d75",
					"deletedDateTime": nil,
					"description":     "Can read basic directory information. Commonly used to grant directory read access to applications and guests.",
					"displayName":     "Directory Readers",
					"roleTemplateId":  "88d8e3e3-c189-46e8-94e1-9b9898b8876b",
				},
				{
					"id":              "18eacdf7-8db3-458d-9099-69fcc2e3cd42",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of Microsoft Entra ID and Microsoft services that use Microsoft Entra identities.",
					"displayName":     "Global Administrator",
					"roleTemplateId":  "62e90394-3621-4004-a7cb-012177145e10",
				},
				{
					"id":              "33a4c989-c3ff-4a77-bf46-ee0acd84476e",
					"deletedDateTime": nil,
					"description":     "Can create application registrations independent of the 'Users can register applications' setting.",
					"displayName":     "Application Developer",
					"roleTemplateId":  "cf1c38e5-69f5-4237-9190-879624dced7c",
				},
				{
					"id":              "321fd63c-c37c-4597-81c4-81e0a93ffb6e",
					"deletedDateTime": nil,
					"description":     "Can manage role assignments in Microsoft Entra ID, and all aspects of Privileged Identity Management.",
					"displayName":     "Privileged Role Administrator",
					"roleTemplateId":  "e8611ab8-8f55-4a1e-953a-60213ab1f814",
				},
				{
					"id":              "8ed2c3eb-1709-444f-bea2-4fb8e2038d6e",
					"deletedDateTime": nil,
					"description":     "Read custom security attribute keys and values for supported Microsoft Entra objects.",
					"displayName":     "Attribute Assignment Reader",
					"roleTemplateId":  "ffd52fa5-98dc-465c-991d-fc073eb59f8f",
				},
				{
					"id":              "e8d9279e-6883-4add-96e8-5f7c8df5637f",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of users and groups, including resetting passwords for limited admins.",
					"displayName":     "User Administrator",
					"roleTemplateId":  "fe930be7-5e62-47db-91af-98c3a49a38b1",
				},
				{
					"id":              "fb96a81c-6147-4cfe-b7fe-c63c2725e7c9",
					"deletedDateTime": nil,
					"description":     "Members of this role can create/manage groups, create/manage groups settings like naming and expiration policies, and view groups activity and audit reports.",
					"displayName":     "Groups Administrator",
					"roleTemplateId":  "fdd7a751-b60b-444a-984c-02652fe8fa1c",
				},
			},
			wantErr:        nil,
			wantNextCursor: nil,
		},
		"first_page": {
			objects: []map[string]any{
				{
					"id":              "0fea7f0d-dea1-4028-8ce8-a686ec639d75",
					"deletedDateTime": nil,
					"description":     "Can read basic directory information. Commonly used to grant directory read access to applications and guests.",
					"displayName":     "Directory Readers",
					"roleTemplateId":  "88d8e3e3-c189-46e8-94e1-9b9898b8876b",
				},
				{
					"id":              "18eacdf7-8db3-458d-9099-69fcc2e3cd42",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of Microsoft Entra ID and Microsoft services that use Microsoft Entra identities.",
					"displayName":     "Global Administrator",
					"roleTemplateId":  "62e90394-3621-4004-a7cb-012177145e10",
				},
				{
					"id":              "33a4c989-c3ff-4a77-bf46-ee0acd84476e",
					"deletedDateTime": nil,
					"description":     "Can create application registrations independent of the 'Users can register applications' setting.",
					"displayName":     "Application Developer",
					"roleTemplateId":  "cf1c38e5-69f5-4237-9190-879624dced7c",
				},
				{
					"id":              "321fd63c-c37c-4597-81c4-81e0a93ffb6e",
					"deletedDateTime": nil,
					"description":     "Can manage role assignments in Microsoft Entra ID, and all aspects of Privileged Identity Management.",
					"displayName":     "Privileged Role Administrator",
					"roleTemplateId":  "e8611ab8-8f55-4a1e-953a-60213ab1f814",
				},
				{
					"id":              "8ed2c3eb-1709-444f-bea2-4fb8e2038d6e",
					"deletedDateTime": nil,
					"description":     "Read custom security attribute keys and values for supported Microsoft Entra objects.",
					"displayName":     "Attribute Assignment Reader",
					"roleTemplateId":  "ffd52fa5-98dc-465c-991d-fc073eb59f8f",
				},
				{
					"id":              "e8d9279e-6883-4add-96e8-5f7c8df5637f",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of users and groups, including resetting passwords for limited admins.",
					"displayName":     "User Administrator",
					"roleTemplateId":  "fe930be7-5e62-47db-91af-98c3a49a38b1",
				},
				{
					"id":              "fb96a81c-6147-4cfe-b7fe-c63c2725e7c9",
					"deletedDateTime": nil,
					"description":     "Members of this role can create/manage groups, create/manage groups settings like naming and expiration policies, and view groups activity and audit reports.",
					"displayName":     "Groups Administrator",
					"roleTemplateId":  "fdd7a751-b60b-444a-984c-02652fe8fa1c",
				},
			},
			pageSize: 1,
			wantObjects: []map[string]any{
				{
					"id":              "0fea7f0d-dea1-4028-8ce8-a686ec639d75",
					"deletedDateTime": nil,
					"description":     "Can read basic directory information. Commonly used to grant directory read access to applications and guests.",
					"displayName":     "Directory Readers",
					"roleTemplateId":  "88d8e3e3-c189-46e8-94e1-9b9898b8876b",
				},
			},
			wantNextCursor: testutil.GenPtr("1"),
			wantErr:        nil,
		},
		"second_page": {
			objects: []map[string]any{
				{
					"id":              "0fea7f0d-dea1-4028-8ce8-a686ec639d75",
					"deletedDateTime": nil,
					"description":     "Can read basic directory information. Commonly used to grant directory read access to applications and guests.",
					"displayName":     "Directory Readers",
					"roleTemplateId":  "88d8e3e3-c189-46e8-94e1-9b9898b8876b",
				},
				{
					"id":              "18eacdf7-8db3-458d-9099-69fcc2e3cd42",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of Microsoft Entra ID and Microsoft services that use Microsoft Entra identities.",
					"displayName":     "Global Administrator",
					"roleTemplateId":  "62e90394-3621-4004-a7cb-012177145e10",
				},
				{
					"id":              "33a4c989-c3ff-4a77-bf46-ee0acd84476e",
					"deletedDateTime": nil,
					"description":     "Can create application registrations independent of the 'Users can register applications' setting.",
					"displayName":     "Application Developer",
					"roleTemplateId":  "cf1c38e5-69f5-4237-9190-879624dced7c",
				},
				{
					"id":              "321fd63c-c37c-4597-81c4-81e0a93ffb6e",
					"deletedDateTime": nil,
					"description":     "Can manage role assignments in Microsoft Entra ID, and all aspects of Privileged Identity Management.",
					"displayName":     "Privileged Role Administrator",
					"roleTemplateId":  "e8611ab8-8f55-4a1e-953a-60213ab1f814",
				},
				{
					"id":              "8ed2c3eb-1709-444f-bea2-4fb8e2038d6e",
					"deletedDateTime": nil,
					"description":     "Read custom security attribute keys and values for supported Microsoft Entra objects.",
					"displayName":     "Attribute Assignment Reader",
					"roleTemplateId":  "ffd52fa5-98dc-465c-991d-fc073eb59f8f",
				},
				{
					"id":              "e8d9279e-6883-4add-96e8-5f7c8df5637f",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of users and groups, including resetting passwords for limited admins.",
					"displayName":     "User Administrator",
					"roleTemplateId":  "fe930be7-5e62-47db-91af-98c3a49a38b1",
				},
				{
					"id":              "fb96a81c-6147-4cfe-b7fe-c63c2725e7c9",
					"deletedDateTime": nil,
					"description":     "Members of this role can create/manage groups, create/manage groups settings like naming and expiration policies, and view groups activity and audit reports.",
					"displayName":     "Groups Administrator",
					"roleTemplateId":  "fdd7a751-b60b-444a-984c-02652fe8fa1c",
				},
			},
			pageSize: 1,
			inputCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("1"),
			},
			wantObjects: []map[string]any{
				{
					"id":              "18eacdf7-8db3-458d-9099-69fcc2e3cd42",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of Microsoft Entra ID and Microsoft services that use Microsoft Entra identities.",
					"displayName":     "Global Administrator",
					"roleTemplateId":  "62e90394-3621-4004-a7cb-012177145e10",
				},
			},
			wantNextCursor: testutil.GenPtr("2"),
			wantErr:        nil,
		},
		"last_page": {
			objects: []map[string]any{
				{
					"id":              "0fea7f0d-dea1-4028-8ce8-a686ec639d75",
					"deletedDateTime": nil,
					"description":     "Can read basic directory information. Commonly used to grant directory read access to applications and guests.",
					"displayName":     "Directory Readers",
					"roleTemplateId":  "88d8e3e3-c189-46e8-94e1-9b9898b8876b",
				},
				{
					"id":              "18eacdf7-8db3-458d-9099-69fcc2e3cd42",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of Microsoft Entra ID and Microsoft services that use Microsoft Entra identities.",
					"displayName":     "Global Administrator",
					"roleTemplateId":  "62e90394-3621-4004-a7cb-012177145e10",
				},
				{
					"id":              "33a4c989-c3ff-4a77-bf46-ee0acd84476e",
					"deletedDateTime": nil,
					"description":     "Can create application registrations independent of the 'Users can register applications' setting.",
					"displayName":     "Application Developer",
					"roleTemplateId":  "cf1c38e5-69f5-4237-9190-879624dced7c",
				},
				{
					"id":              "321fd63c-c37c-4597-81c4-81e0a93ffb6e",
					"deletedDateTime": nil,
					"description":     "Can manage role assignments in Microsoft Entra ID, and all aspects of Privileged Identity Management.",
					"displayName":     "Privileged Role Administrator",
					"roleTemplateId":  "e8611ab8-8f55-4a1e-953a-60213ab1f814",
				},
				{
					"id":              "8ed2c3eb-1709-444f-bea2-4fb8e2038d6e",
					"deletedDateTime": nil,
					"description":     "Read custom security attribute keys and values for supported Microsoft Entra objects.",
					"displayName":     "Attribute Assignment Reader",
					"roleTemplateId":  "ffd52fa5-98dc-465c-991d-fc073eb59f8f",
				},
				{
					"id":              "e8d9279e-6883-4add-96e8-5f7c8df5637f",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of users and groups, including resetting passwords for limited admins.",
					"displayName":     "User Administrator",
					"roleTemplateId":  "fe930be7-5e62-47db-91af-98c3a49a38b1",
				},
				{
					"id":              "fb96a81c-6147-4cfe-b7fe-c63c2725e7c9",
					"deletedDateTime": nil,
					"description":     "Members of this role can create/manage groups, create/manage groups settings like naming and expiration policies, and view groups activity and audit reports.",
					"displayName":     "Groups Administrator",
					"roleTemplateId":  "fdd7a751-b60b-444a-984c-02652fe8fa1c",
				},
			},
			pageSize: 1,
			inputCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("6"),
			},
			wantObjects: []map[string]any{
				{
					"id":              "fb96a81c-6147-4cfe-b7fe-c63c2725e7c9",
					"deletedDateTime": nil,
					"description":     "Members of this role can create/manage groups, create/manage groups settings like naming and expiration policies, and view groups activity and audit reports.",
					"displayName":     "Groups Administrator",
					"roleTemplateId":  "fdd7a751-b60b-444a-984c-02652fe8fa1c",
				},
			},
			wantErr:        nil,
			wantNextCursor: nil,
		},
		"bigger_page": {
			objects: []map[string]any{
				{
					"id":              "0fea7f0d-dea1-4028-8ce8-a686ec639d75",
					"deletedDateTime": nil,
					"description":     "Can read basic directory information. Commonly used to grant directory read access to applications and guests.",
					"displayName":     "Directory Readers",
					"roleTemplateId":  "88d8e3e3-c189-46e8-94e1-9b9898b8876b",
				},
				{
					"id":              "18eacdf7-8db3-458d-9099-69fcc2e3cd42",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of Microsoft Entra ID and Microsoft services that use Microsoft Entra identities.",
					"displayName":     "Global Administrator",
					"roleTemplateId":  "62e90394-3621-4004-a7cb-012177145e10",
				},
				{
					"id":              "33a4c989-c3ff-4a77-bf46-ee0acd84476e",
					"deletedDateTime": nil,
					"description":     "Can create application registrations independent of the 'Users can register applications' setting.",
					"displayName":     "Application Developer",
					"roleTemplateId":  "cf1c38e5-69f5-4237-9190-879624dced7c",
				},
				{
					"id":              "321fd63c-c37c-4597-81c4-81e0a93ffb6e",
					"deletedDateTime": nil,
					"description":     "Can manage role assignments in Microsoft Entra ID, and all aspects of Privileged Identity Management.",
					"displayName":     "Privileged Role Administrator",
					"roleTemplateId":  "e8611ab8-8f55-4a1e-953a-60213ab1f814",
				},
				{
					"id":              "8ed2c3eb-1709-444f-bea2-4fb8e2038d6e",
					"deletedDateTime": nil,
					"description":     "Read custom security attribute keys and values for supported Microsoft Entra objects.",
					"displayName":     "Attribute Assignment Reader",
					"roleTemplateId":  "ffd52fa5-98dc-465c-991d-fc073eb59f8f",
				},
				{
					"id":              "e8d9279e-6883-4add-96e8-5f7c8df5637f",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of users and groups, including resetting passwords for limited admins.",
					"displayName":     "User Administrator",
					"roleTemplateId":  "fe930be7-5e62-47db-91af-98c3a49a38b1",
				},
				{
					"id":              "fb96a81c-6147-4cfe-b7fe-c63c2725e7c9",
					"deletedDateTime": nil,
					"description":     "Members of this role can create/manage groups, create/manage groups settings like naming and expiration policies, and view groups activity and audit reports.",
					"displayName":     "Groups Administrator",
					"roleTemplateId":  "fdd7a751-b60b-444a-984c-02652fe8fa1c",
				},
			},
			pageSize: 2,
			inputCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("5"),
			},
			wantObjects: []map[string]any{
				{
					"id":              "e8d9279e-6883-4add-96e8-5f7c8df5637f",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of users and groups, including resetting passwords for limited admins.",
					"displayName":     "User Administrator",
					"roleTemplateId":  "fe930be7-5e62-47db-91af-98c3a49a38b1",
				},
				{
					"id":              "fb96a81c-6147-4cfe-b7fe-c63c2725e7c9",
					"deletedDateTime": nil,
					"description":     "Members of this role can create/manage groups, create/manage groups settings like naming and expiration policies, and view groups activity and audit reports.",
					"displayName":     "Groups Administrator",
					"roleTemplateId":  "fdd7a751-b60b-444a-984c-02652fe8fa1c",
				},
			},
			wantNextCursor: nil,
			wantErr:        nil,
		},
		"invalid_cursor": {
			objects: []map[string]any{
				{
					"id":              "0fea7f0d-dea1-4028-8ce8-a686ec639d75",
					"deletedDateTime": nil,
					"description":     "Can read basic directory information. Commonly used to grant directory read access to applications and guests.",
					"displayName":     "Directory Readers",
					"roleTemplateId":  "88d8e3e3-c189-46e8-94e1-9b9898b8876b",
				},
				{
					"id":              "18eacdf7-8db3-458d-9099-69fcc2e3cd42",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of Microsoft Entra ID and Microsoft services that use Microsoft Entra identities.",
					"displayName":     "Global Administrator",
					"roleTemplateId":  "62e90394-3621-4004-a7cb-012177145e10",
				},
				{
					"id":              "33a4c989-c3ff-4a77-bf46-ee0acd84476e",
					"deletedDateTime": nil,
					"description":     "Can create application registrations independent of the 'Users can register applications' setting.",
					"displayName":     "Application Developer",
					"roleTemplateId":  "cf1c38e5-69f5-4237-9190-879624dced7c",
				},
				{
					"id":              "321fd63c-c37c-4597-81c4-81e0a93ffb6e",
					"deletedDateTime": nil,
					"description":     "Can manage role assignments in Microsoft Entra ID, and all aspects of Privileged Identity Management.",
					"displayName":     "Privileged Role Administrator",
					"roleTemplateId":  "e8611ab8-8f55-4a1e-953a-60213ab1f814",
				},
				{
					"id":              "8ed2c3eb-1709-444f-bea2-4fb8e2038d6e",
					"deletedDateTime": nil,
					"description":     "Read custom security attribute keys and values for supported Microsoft Entra objects.",
					"displayName":     "Attribute Assignment Reader",
					"roleTemplateId":  "ffd52fa5-98dc-465c-991d-fc073eb59f8f",
				},
				{
					"id":              "e8d9279e-6883-4add-96e8-5f7c8df5637f",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of users and groups, including resetting passwords for limited admins.",
					"displayName":     "User Administrator",
					"roleTemplateId":  "fe930be7-5e62-47db-91af-98c3a49a38b1",
				},
				{
					"id":              "fb96a81c-6147-4cfe-b7fe-c63c2725e7c9",
					"deletedDateTime": nil,
					"description":     "Members of this role can create/manage groups, create/manage groups settings like naming and expiration policies, and view groups activity and audit reports.",
					"displayName":     "Groups Administrator",
					"roleTemplateId":  "fdd7a751-b60b-444a-984c-02652fe8fa1c",
				},
			},
			pageSize: 1,
			inputCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("h7fy7"),
			},
			wantObjects: nil,
			wantErr: &framework.Error{
				Message: "unable to parse cursor: want valid number, got {h7fy7}",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
			wantNextCursor: nil,
		},
		"invalid_cursor_out_of_range": {
			objects: []map[string]any{
				{
					"id":              "0fea7f0d-dea1-4028-8ce8-a686ec639d75",
					"deletedDateTime": nil,
					"description":     "Can read basic directory information. Commonly used to grant directory read access to applications and guests.",
					"displayName":     "Directory Readers",
					"roleTemplateId":  "88d8e3e3-c189-46e8-94e1-9b9898b8876b",
				},
				{
					"id":              "18eacdf7-8db3-458d-9099-69fcc2e3cd42",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of Microsoft Entra ID and Microsoft services that use Microsoft Entra identities.",
					"displayName":     "Global Administrator",
					"roleTemplateId":  "62e90394-3621-4004-a7cb-012177145e10",
				},
				{
					"id":              "33a4c989-c3ff-4a77-bf46-ee0acd84476e",
					"deletedDateTime": nil,
					"description":     "Can create application registrations independent of the 'Users can register applications' setting.",
					"displayName":     "Application Developer",
					"roleTemplateId":  "cf1c38e5-69f5-4237-9190-879624dced7c",
				},
				{
					"id":              "321fd63c-c37c-4597-81c4-81e0a93ffb6e",
					"deletedDateTime": nil,
					"description":     "Can manage role assignments in Microsoft Entra ID, and all aspects of Privileged Identity Management.",
					"displayName":     "Privileged Role Administrator",
					"roleTemplateId":  "e8611ab8-8f55-4a1e-953a-60213ab1f814",
				},
				{
					"id":              "8ed2c3eb-1709-444f-bea2-4fb8e2038d6e",
					"deletedDateTime": nil,
					"description":     "Read custom security attribute keys and values for supported Microsoft Entra objects.",
					"displayName":     "Attribute Assignment Reader",
					"roleTemplateId":  "ffd52fa5-98dc-465c-991d-fc073eb59f8f",
				},
				{
					"id":              "e8d9279e-6883-4add-96e8-5f7c8df5637f",
					"deletedDateTime": nil,
					"description":     "Can manage all aspects of users and groups, including resetting passwords for limited admins.",
					"displayName":     "User Administrator",
					"roleTemplateId":  "fe930be7-5e62-47db-91af-98c3a49a38b1",
				},
				{
					"id":              "fb96a81c-6147-4cfe-b7fe-c63c2725e7c9",
					"deletedDateTime": nil,
					"description":     "Members of this role can create/manage groups, create/manage groups settings like naming and expiration policies, and view groups activity and audit reports.",
					"displayName":     "Groups Administrator",
					"roleTemplateId":  "fdd7a751-b60b-444a-984c-02652fe8fa1c",
				},
			},
			pageSize: 1,
			inputCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("9"),
			},
			wantObjects: nil,
			wantErr: &framework.Error{
				Message: "The cursor value: 9, is out of range for number of objects: 7",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
			wantNextCursor: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := pagination.PaginateObjects(tt.objects, tt.pageSize, tt.inputCursor)

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !reflect.DeepEqual(gotNextCursor, tt.wantNextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotNextCursor, tt.wantNextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
