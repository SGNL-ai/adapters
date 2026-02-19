// Copyright 2026 SGNL.ai, Inc.

package googleworkspace

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

const (
	OrderByAscending = "ASCENDING"
)

// ConstructEndpoint constructs and returns the endpoint to query the datasource.
func ConstructEndpoint(request *Request) (string, *framework.Error) {
	if request == nil {
		return "", &framework.Error{
			Message: "Request is nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// Example URLs:
	// [User]: https://admin.googleapis.com/admin/directory/v1/users?domain=sgnldemos.com&maxResults=500
	// [Group]: https://admin.googleapis.com/admin/directory/v1/groups?domain=sgnldemos.com&maxResults=500
	// [Member]: https://admin.googleapis.com/admin/directory/v1/groups/0300/members?domain=sgnldemos.com&maxResults=2

	var sb strings.Builder

	sb.Grow(len(request.BaseURL) + len(ValidEntityExternalIDs[request.EntityExternalID].Path) + 5)
	sb.WriteString(request.BaseURL)

	if request.EntityExternalID == Member {
		if request.Cursor == nil || request.Cursor.CollectionID == nil {
			return "", &framework.Error{
				Message: "Collection ID is nil for Member entity, unable to form request URI.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		sb.WriteString(fmt.Sprintf(ValidEntityExternalIDs[request.EntityExternalID].Path,
			request.APIVersion, *request.Cursor.CollectionID))
	} else {
		sb.WriteString(fmt.Sprintf(ValidEntityExternalIDs[request.EntityExternalID].Path, request.APIVersion))
	}

	params := url.Values{}
	params.Add("maxResults", strconv.FormatInt(request.PageSize, 10))

	if request.Cursor != nil && request.Cursor.Cursor != nil {
		params.Add("pageToken", *request.Cursor.Cursor)
	}

	if request.Customer != nil {
		params.Add("customer", *request.Customer)
	}

	if request.Domain != nil {
		params.Add("domain", *request.Domain)
	}

	switch request.EntityExternalID {
	case User:
		AddUserParams(&params, request)
	case Group:
		AddGroupParams(&params, request)
	case Member:
		AddMemberParams(&params, request)
	default:
		return "", &framework.Error{
			Message: fmt.Sprintf("Entity ID %v is not supported.", request.EntityExternalID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	paramString := params.Encode()

	sb.Grow(len(paramString) + 1)
	sb.WriteString("?")
	sb.WriteString(paramString)

	return sb.String(), nil
}

func AddUserParams(params *url.Values, request *Request) {
	if request.UserFilters != nil {
		if request.UserFilters.Query != nil {
			params.Add("query", *request.UserFilters.Query)
		}

		params.Add("showDeleted", strconv.FormatBool(request.UserFilters.ShowDeleted))
	}

	if request.Ordered && ValidEntityExternalIDs[request.EntityExternalID].OrderByAttribute != "" {
		params.Add("orderBy", ValidEntityExternalIDs[request.EntityExternalID].OrderByAttribute)
		params.Add("sortOrder", OrderByAscending)
	}
}

func AddGroupParams(params *url.Values, request *Request) {
	if request.GroupFilters != nil && request.GroupFilters.Query != nil {
		params.Add("query", *request.GroupFilters.Query)
	}

	if request.Ordered && ValidEntityExternalIDs[request.EntityExternalID].OrderByAttribute != "" {
		params.Add("orderBy", ValidEntityExternalIDs[request.EntityExternalID].OrderByAttribute)
		params.Add("sortOrder", OrderByAscending)
	}
}

func AddMemberParams(params *url.Values, request *Request) {
	if request.MemberFilters != nil {
		if request.MemberFilters.Roles != nil {
			params.Add("roles", *request.MemberFilters.Roles)
		}

		params.Add("includeDerivedMembership", strconv.FormatBool(request.MemberFilters.IncludeDerivedMembership))
	}
}
