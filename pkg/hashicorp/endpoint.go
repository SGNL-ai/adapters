// Copyright 2025 SGNL.ai, Inc.
package hashicorp

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ConstructEndpoint constructs and returns the endpoint to query the datasource.
func ConstructEndpoint(request *Request) string {
	var sb strings.Builder

	// URL Format:
	// baseURL + "/v1/" + resourceType + "?page_size=" + pageSize + ["&filter=" + filter] + ["&list_token=" + cursor]
	// Example: https://boundary.example.com/v1/roles?page_size=100&filter=name%20eq%20"admin"&list_token=abc123
	sb.Grow(len(request.BaseURL) + len(request.EntityExternalID) + 52)

	sb.WriteString(request.BaseURL)
	sb.WriteString(fmt.Sprintf("/%s/", APIVersion))
	sb.WriteString(request.EntityExternalID)

	params := url.Values{}
	params.Add("page_size", strconv.FormatInt(request.PageSize, 10))
	params.Add("recursive", "true")

	if config, ok := request.EntityConfig[request.EntityExternalID]; ok && config.Filter != "" {
		params.Add("filter", config.Filter)
	}

	if request.Cursor != nil && request.Cursor.Cursor != nil && len(*request.Cursor.Cursor) > 0 {
		params.Add("list_token", *request.Cursor.Cursor)
	}

	if config, ok := request.EntityConfig[request.EntityExternalID]; ok && config.ScopeID != "" {
		params.Add("scope_id", config.ScopeID)
	} else {
		params.Add("scope_id", "global")
	}

	for key, value := range request.AdditionalParams {
		params.Add(key, value)
	}

	paramString := params.Encode()
	if paramString != "" {
		sb.WriteString("?")
		sb.WriteString(paramString)
	}

	return sb.String()
}
