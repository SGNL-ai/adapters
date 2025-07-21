// Copyright 2025 SGNL.ai, Inc.
package okta

import (
	"net/url"
	"strconv"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// ConstructEndpoint constructs and returns the endpoint to query the datasource.
func ConstructEndpoint(request *Request) (string, *framework.Error) {
	if request == nil {
		return "", nil
	}

	var endpoint string

	var filter string

	var search string

	// [Users / Groups] This is the cursor to the next page of objects.
	// [GroupMembers] This is the cursor to the next page of Members.
	if request.Cursor != nil && request.Cursor.Cursor != nil {
		endpoint = *request.Cursor.Cursor
	}

	if endpoint == "" {
		formattedPageSize := strconv.FormatInt(request.PageSize, 10)

		var sb strings.Builder

		// URL Format:
		// [Users]		baseURL + "/api/" + apiVersion + "/users?limit=" + pageSize
		// [Filtered Users]	baseURL + "/api/" + apiVersion + "/users?filter="
		// 					+ `status eq \"ACTIVE\"` + "&limit=" + pageSize
		// [Groups]		baseURL + "/api/" + apiVersion + "/groups?filter="
		//                  + `type eq "OKTA_GROUP" or type eq "APP_GROUP"` + "&limit=" + pageSize
		// [GroupMembers] 	baseURL + "/api/" + apiVersion + "/groups/" + groupId + "/users?limit=" + pageSize
		sb.Grow(len(request.BaseURL) + len(request.APIVersion) + len(formattedPageSize) + 12)

		sb.WriteString(request.BaseURL)
		sb.WriteString("/api/")
		sb.WriteString(request.APIVersion)
		sb.WriteString("/")

		if request.Filter != "" {
			// Okta uses double quotes in the filter string, so we need to handle
			// escaping them in config. We need to replace them then encode the filter.
			filter = url.QueryEscape(strings.ReplaceAll(request.Filter, `\"`, `"`))
			// The minimum length of a valid Okta filter is 7 characters
			// given the shortest valid filter is in the form of "id eq x".
			if len(filter) < 7 {
				return "", &framework.Error{
					Message: "Provided filter is invalid.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}
			}
		}

		if request.Search != "" {
			// Okta uses double quotes in the search string, so we need to handle
			// escaping them in config. We need to replace them then encode the search.
			search = url.QueryEscape(strings.ReplaceAll(request.Search, `\"`, `"`))
			// The minimum length of a valid Okta search is 7 characters
			// given the shortest valid search is in the form of "id eq x".
			if len(search) < 7 {
				return "", &framework.Error{
					Message: "Provided search syntax is invalid.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}
			}
		}

		switch request.EntityExternalID {
		case Users:
			// Build the base URL first
			sb.WriteString("users?")

			// Create a slice to hold parameter strings
			paramParts := []string{}

			if filter != "" {
				paramParts = append(paramParts, "filter="+filter)
			}

			if search != "" {
				paramParts = append(paramParts, "search="+search)
			}

			// Join all parameters with &
			if len(paramParts) > 0 {
				sb.WriteString(strings.Join(paramParts, "&"))
				sb.WriteString("&")
			}

		case Groups:
			// Build the base URL first
			sb.WriteString("groups?")

			// Create a slice to hold parameter strings
			paramParts := []string{}

			if filter == "" && search == "" {
				// Some Groups are not useful to ingest into SGNL, automatically filtering.
				filter = url.QueryEscape(`type eq "OKTA_GROUP" or type eq "APP_GROUP"`)
			}

			if filter != "" {
				paramParts = append(paramParts, "filter="+filter)
			}

			if search != "" {
				paramParts = append(paramParts, "search="+search)
			}

			// Join all parameters with &
			if len(paramParts) > 0 {
				sb.WriteString(strings.Join(paramParts, "&"))
				sb.WriteString("&")
			}
		case GroupMembers:
			if request.Cursor == nil || request.Cursor.CollectionID == nil {
				return "", &framework.Error{
					Message: "Unable to construct group member endpoint without valid cursor.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			sb.Grow(len(*request.Cursor.CollectionID) + 14)

			sb.WriteString("groups/")
			sb.WriteString(*request.Cursor.CollectionID)
			sb.WriteString("/users?")
		default:
			return "", &framework.Error{
				Message: "Provided entity external ID is invalid.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			}
		}

		sb.WriteString("limit=")
		sb.WriteString(formattedPageSize)

		endpoint = sb.String()
	}

	return endpoint, nil
}
