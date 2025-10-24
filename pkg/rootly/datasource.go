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
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"go.uber.org/zap"
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
	logger := zaplogger.FromContext(ctx).With(
		fields.RequestEntityExternalID(request.EntityExternalID),
		fields.RequestPageSize(request.PageSize),
	)

	logger.Info("Starting datasource request")

	endpoint := ConstructEndpoint(request)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
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

	logger.Info("Sending request to datasource", fields.RequestURL(endpoint))

	// Use the client from the datasource instead of the request
	resp, err := d.Client.Do(req)
	if err != nil {
		logger.Error("Request to datasource failed",
			fields.RequestURL(endpoint),
			fields.SGNLEventTypeError(),
			zap.Error(err),
		)

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
		logger.Error("Datasource responded with an error",
			fields.RequestURL(endpoint),
			fields.ResponseStatusCode(resp.StatusCode),
			fields.ResponseRetryAfterHeader(resp.Header.Get("Retry-After")),
			fields.ResponseBody(body),
			fields.SGNLEventTypeError(),
		)

		var errorResponse DatasourceErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf(
					"Failed to parse error response: %v. Status: %d. Body: %s.", err, resp.StatusCode, string(body),
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		errorMsg := fmt.Sprintf("Received Http Error %d: %s", resp.StatusCode, resp.Status)

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

	response := &Response{
		Objects:    datasourceResponse.Data,
		NextCursor: nextCursor,
	}

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(resp.StatusCode),
		fields.ResponseObjectCount(len(response.Objects)),
		fields.ResponseNextCursor(response.NextCursor),
	)

	return response, nil
}
