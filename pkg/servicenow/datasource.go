// Copyright 2026 SGNL.ai, Inc.

package servicenow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"github.com/sgnl-ai/adapters/pkg/extractor"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"go.uber.org/zap"
)

const (
	User          = "sys_user"
	Group         = "sys_user_group"
	GroupMember   = "sys_user_grmember"
	Case          = "sn_customerservice_case"
	Incident      = "incident"
	ChangeRequest = "change_request"
	ChangeTask    = "change_task"
)

// Datasource directly implements a Client interface to allow querying an external datasource.
type Datasource struct {
	Client *http.Client
}

type DatasourceResponse struct {
	Result []map[string]any `json:"result"`
}

type DatasourceErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Detail  string `json:"detail"`
	} `json:"error,omitempty"`
	Status string `json:"status,omitempty"`
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

	req.Header.Add("Authorization", request.AuthorizationHeader)

	logger.Info("Sending request to datasource", fields.RequestURL(endpoint))

	res, err := d.Client.Do(req)
	if err != nil {
		logger.Error("Request to datasource failed",
			fields.RequestURL(endpoint),
			fields.SGNLEventTypeError(),
			zap.Error(err),
		)

		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute request: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds),
		)
	}

	defer res.Body.Close()

	response := &Response{
		StatusCode:       res.StatusCode,
		RetryAfterHeader: res.Header.Get("Retry-After"),
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to read response (%d): %v.", res.StatusCode, err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Edge case: If the `sysparm_query` parameter is too large and the page size is too small,
	// ServiceNow will return a 400 Bad Request with a message "Pagination not supported" and the reason.
	// We need to surface this error to the user.
	if res.StatusCode != http.StatusOK {
		logger.Error("Datasource responded with an error",
			fields.RequestURL(endpoint),
			fields.ResponseStatusCode(res.StatusCode),
			fields.ResponseRetryAfterHeader(res.Header.Get("Retry-After")),
			fields.ResponseBody(io.NopCloser(bytes.NewReader(body))),
			fields.SGNLEventTypeError(),
		)

		var errorResponse DatasourceErrorResponse

		if err := json.Unmarshal(body, &errorResponse); err == nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf(
					"Failed to get page from datasource: %d. Message: `%s`. Details: `%s`.",
					response.StatusCode,
					errorResponse.Error.Message,
					errorResponse.Error.Detail,
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
			}
		}

		// If we can't parse the error messsage, just return the status code.
		return response, nil
	}

	objects, frameworkErr := ParseResponse(body)
	if frameworkErr != nil {
		return nil, frameworkErr
	}

	response.Objects = objects

	if cursor := extractor.ValueFromList(res.Header.Values("Link"), "https://", ">;rel=\"next\""); cursor != "" {
		response.NextCursor = &cursor
	}

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(response.StatusCode),
		fields.ResponseObjectCount(len(response.Objects)),
		fields.ResponseNextCursor(response.NextCursor),
	)

	return response, nil
}

func ParseResponse(body []byte) ([]map[string]any, *framework.Error) {
	var data *DatasourceResponse

	if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return data.Result, nil
}
