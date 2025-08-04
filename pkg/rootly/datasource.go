// Copyright 2025 SGNL.ai, Inc.
package rootly

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// Datasource directly implements a Client interface to allow querying an external datasource.
type Datasource struct {
	Client *http.Client
}

type DatasourceResponse struct {
	Data []map[string]any `json:"data"`
	Meta struct {
		Page       int `json:"current_page"`
		Pages      int `json:"total_pages"`
		TotalCount int `json:"total_count"`
	} `json:"meta"`
}

type DatasourceErrorResponse struct {
	Errors []struct {
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Status string `json:"status"`
	} `json:"errors"`
}

// NewClient returns a Client to query the datasource.
func NewClient(client *http.Client) Client {
	return &Datasource{
		Client: client,
	}
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ConstructEndpoint(request), nil)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to create request to datasource: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Timeout API calls that take longer than the configured timeout.
	apiCtx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	req = req.WithContext(apiCtx)

	req.Header.Add("Authorization", request.HTTPAuthorization)
	req.Header.Add("Content-Type", "application/vnd.api+json")

	// Use the client from the datasource instead of the request
	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to query datasource: %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to read response body: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse DatasourceErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf(
					"Failed to parse error response: %v. Status: %d. Body: %s.", err, resp.StatusCode, string(body),
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		errorMsg := fmt.Sprintf("HTTP %d", resp.StatusCode)
		if len(errorResponse.Errors) > 0 {
			errorMsg = fmt.Sprintf("%s: %s", errorMsg, errorResponse.Errors[0].Detail)
		}

		return nil, &framework.Error{
			Message: errorMsg,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	var datasourceResponse DatasourceResponse
	if err := json.Unmarshal(body, &datasourceResponse); err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to parse response: %v. Body: %s.", err, string(body)),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Check if the response has the required data field
	if datasourceResponse.Data == nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Invalid response format: missing required data field. Body: %s.", string(body)),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Determine next cursor based on pagination
	var nextCursor *string

	currentPage := datasourceResponse.Meta.Page
	totalPages := datasourceResponse.Meta.Pages

	if currentPage < totalPages {
		nextPageStr := strconv.Itoa(currentPage + 1)
		nextCursor = &nextPageStr
	}

	return &Response{
		Objects:    datasourceResponse.Data,
		NextCursor: nextCursor,
	}, nil
}
