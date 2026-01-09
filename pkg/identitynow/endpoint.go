// Copyright 2026 SGNL.ai, Inc.
package identitynow

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

	// GetPage should always set the cursor to a non-nil value. If it is nil here, this is a bug within the adapter.
	if request.Cursor == nil || request.Cursor.Cursor == nil {
		return "", &framework.Error{
			Message: "Request cursor is unexpectedly nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	pageSizeAsStr := strconv.FormatInt(request.PageSize, 10)
	offsetAsStr := strconv.FormatInt(*request.Cursor.Cursor, 10)

	var endpoint strings.Builder

	// URL Format:
	// nolint: lll
	// [AccountEntitlements] baseURL + "/" + apiVersion + "/accounts/" + request.Cursor.CollectionID + "/entitlements?limit=" + pageSizeAsStr + "&offset=" + offsetAsStr
	// [any other entity]     baseURL + "/" + apiVersion + "/" + entityExternalID + "?limit=" + pageSizeAsStr + "&offset=" + offsetAsStr
	endpoint.Grow(len(request.BaseURL) + len(request.APIVersion) + len(pageSizeAsStr) + len(offsetAsStr) + 17)
	endpoint.WriteString(request.BaseURL)
	endpoint.WriteString("/")
	endpoint.WriteString(request.APIVersion)
	endpoint.WriteString("/")

	switch request.EntityExternalID {
	case AccountEntitlements:
		if request.Cursor.CollectionID == nil {
			return "", &framework.Error{
				Message: "CollectionId field is unexpectedly nil.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		collectionIDEscaped := url.PathEscape(*request.Cursor.CollectionID) // Escapes spaces to %20.

		endpoint.Grow(len(collectionIDEscaped) + 22)
		endpoint.WriteString("accounts/")
		endpoint.WriteString(collectionIDEscaped)
		endpoint.WriteString("/entitlements")
	default:
		endpoint.WriteString(request.EntityExternalID)
	}

	endpoint.WriteString("?limit=")
	endpoint.WriteString(pageSizeAsStr)
	endpoint.WriteString("&offset=")
	endpoint.WriteString(offsetAsStr)

	if request.Sorters != nil {
		endpoint.WriteString("&sorters=")
		// From the documentation, it doesn't seem like we need to escape the ',' in the sorters.
		endpoint.WriteString(*request.Sorters)
	}

	// Apply any filters.
	// TODO [sc-19213]: We don't apply any filters when querying AccountEntitlements. This is because
	// the filter that is passed in must be used to filter the accounts that we're retrieving entitlements for
	// and NOT the entitlements themselves. This is a limitation that must be addressed in the future.
	if request.Filter != nil && request.EntityExternalID != AccountEntitlements {
		// IdentityNow requires spaces to be encoded as %20 instead of +.
		// https://developer.sailpoint.com/idn/api/standard-collection-parameters/#known-limitations.
		// Golang's url.QueryEscape() encodes spaces as +, so we need to replace them with %20.
		escapedFilter := strings.Replace(url.QueryEscape(*request.Filter), "+", "%20", -1)
		endpoint.Grow(len(escapedFilter) + 9)

		endpoint.WriteString("&filters=")
		endpoint.WriteString(escapedFilter)
	}

	return endpoint.String(), nil
}
