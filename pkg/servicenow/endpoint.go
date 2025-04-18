// Copyright 2025 SGNL.ai, Inc.
package servicenow

import (
	"net/url"
	"strconv"
	"strings"
)

// ConstructEndpoint constructs and returns the endpoint to query the datasource.
func ConstructEndpoint(request *Request) string {
	if request == nil {
		return ""
	}

	if request.Cursor != nil {
		return *request.Cursor
	}

	// URL Format:
	// baseURL + "/api/now/" + apiVersion + "/table/" + tableName + "?sysparm_fields=sys_id"
	// 		+ "&sysparm_exclude_reference_link=true&sysparm_limit=" + pageSize
	// 		+ ["&sysparm_query=" + filter + "%5EORDERBYsys_id"] | ["&sysparm_query=ORDERBYsys_id"]

	var sb strings.Builder

	pageSizeStr := strconv.FormatInt(request.PageSize, 10)

	sb.Grow(len(request.BaseURL) + len(request.APIVersion) + len(request.EntityExternalID) + len(pageSizeStr) + 89)

	sb.WriteString(request.BaseURL)
	sb.WriteString("/api/now/")
	sb.WriteString(request.APIVersion)
	sb.WriteString("/table/")
	sb.WriteString(request.EntityExternalID)
	sb.WriteString("?sysparm_fields=sys_id")

	for _, attribute := range request.Attributes {
		// sys_id is added to all requests by default to enable sorting, so don't re-add
		if attribute.ExternalId == "sys_id" {
			continue
		}

		encodedExternalID := url.QueryEscape(attribute.ExternalId)

		sb.Grow(1 + len(encodedExternalID))
		sb.WriteRune(',')
		sb.WriteString(encodedExternalID)
	}

	sb.WriteString("&sysparm_exclude_reference_link=true&sysparm_limit=")
	sb.WriteString(pageSizeStr)

	escapedFilter := ""

	if request.Filter != nil && *request.Filter != "" {
		escapedFilter = url.QueryEscape(*request.Filter) + "%5E"
	}

	sb.Grow(31 + len(escapedFilter))

	sb.WriteString("&sysparm_query=")
	sb.WriteString(escapedFilter)
	sb.WriteString("ORDERBYsys_id")

	return sb.String()
}
