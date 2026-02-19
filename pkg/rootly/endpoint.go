// Copyright 2026 SGNL.ai, Inc.

package rootly

import (
	"fmt"
	"net/url"
	"strconv"
)

// ConstructEndpoint constructs the endpoint URL for the given request.
func ConstructEndpoint(request *Request) string {
	endpoint := fmt.Sprintf("%s/%s", request.BaseURL, request.EntityExternalID)

	params := url.Values{}

	// Add page size
	if request.PageSize > 0 {
		params.Add("page[size]", strconv.FormatInt(request.PageSize, 10))
	}

	// Add page number if cursor is provided
	if request.Cursor != nil && *request.Cursor != "" {
		params.Add("page[number]", *request.Cursor)
	} else {
		params.Add("page[number]", "1")
	}

	// Add filter if provided
	if request.Filter != "" {
		// Parse user-provided filters and convert to Rootly's filter[key]=value format
		filterParams, err := url.ParseQuery(request.Filter)
		if err == nil {
			for key, values := range filterParams {
				for _, value := range values {
					params.Add(fmt.Sprintf("filter[%s]", key), value)
				}
			}
		}
	}

	// Add includes if provided
	if request.Includes != "" {
		// Rootly expects includes as a comma-separated list
		params.Add("include", request.Includes)
	}

	if len(params) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	}

	return endpoint
}
