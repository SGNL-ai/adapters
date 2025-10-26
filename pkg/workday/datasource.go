// Copyright 2025 SGNL.ai, Inc.
package workday

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"go.uber.org/zap"
)

// Datasource directly implements a Client interface to allow querying an external datasource.
type Datasource struct {
	Client *http.Client
}

// Workday API response format.
type DatasourceResponse struct {
	Total   *int64           `json:"total"`
	Objects []map[string]any `json:"data"`
	Error   *string          `json:"error"`
	Errors  []ErrorInfo      `json:"errors"`
}

type ErrorInfo struct {
	Error    string `json:"error"`
	Field    string `json:"field"`
	Location string `json:"location"`
}

// Entity contains entity specific information, such as the entity's unique ID attribute and the
// endpoint path to query that entity.
type Entity struct {
	// TableName is the name of the entity's table in the datasource.
	TableName string
	// UniqueIDAttrExternalID is the external ID of the entity's uniqueId attribute.
	UniqueIDAttrExternalID string
}

// NewClient returns a Client to query the datasource.
func NewClient(client *http.Client) Client {
	return &Datasource{
		Client: client,
	}
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	logger := zaplogger.FromContext(ctx).With(
		fields.RequestEntityExternalID(request.EntityConfig.ExternalId),
		fields.RequestPageSize(request.PageSize),
	)

	logger.Info("Starting datasource request")

	if validationErr := pagination.ValidateCompositeCursor(
		request.Cursor,
		request.EntityConfig.ExternalId,
		false,
	); validationErr != nil {
		return nil, validationErr
	}

	endpoint, endpointErr := ConstructEndpoint(request)
	if endpointErr != nil {
		return nil, endpointErr
	}

	// Timeout API calls that take longer than the configured timeout.
	apiCtx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(apiCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to create request to datasource: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	req.Header.Add("Authorization", request.Token)

	logger.Info("Sending request to datasource", fields.RequestURL(endpoint))

	res, err := d.Client.Do(req)
	if err != nil {
		logger.Error("Request to datasource failed",
			fields.RequestURL(endpoint),
			fields.SGNLEventTypeError(),
			zap.Error(err),
		)

		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute Workday request: %v.", err),
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

	if res.StatusCode != http.StatusOK {
		logger.Error("Datasource responded with an error",
			fields.RequestURL(endpoint),
			fields.ResponseStatusCode(res.StatusCode),
			fields.ResponseRetryAfterHeader(res.Header.Get("Retry-After")),
			fields.ResponseBody(res.Body),
			fields.SGNLEventTypeError(),
		)

		return response, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to read Workday response body: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	objects, nextCursor, frameworkErr := ParseResponse(body, request, endpoint)
	if frameworkErr != nil {
		return nil, frameworkErr
	}

	response.NextCursor = nextCursor
	response.Objects = objects

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(response.StatusCode),
		fields.ResponseObjectCount(len(response.Objects)),
		fields.ResponseNextCursor(response.NextCursor),
	)

	return response, nil
}

func ParseResponse(body []byte, request *Request, endpoint string) (
	objects []map[string]any,
	nextCursor *pagination.CompositeCursor[int64],
	err *framework.Error,
) {
	var data *DatasourceResponse

	if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil || data == nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if data.Error != nil {
		return nil, nil, ParseError(data, endpoint)
	}

	if data.Total == nil {
		return nil, nil, &framework.Error{
			Message: "Total count is missing in the datasource response.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if data.Objects == nil {
		return nil, nil, &framework.Error{
			Message: "Missing data in the datasource response.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	nextOffset := request.PageSize

	// If the cursor is not nil, increment the next offset by the last offset value.
	if request.Cursor != nil && request.Cursor.Cursor != nil {
		nextOffset += *request.Cursor.Cursor
	}

	if *data.Total > nextOffset {
		nextCursor = &pagination.CompositeCursor[int64]{
			Cursor: &nextOffset,
		}
	}

	return data.Objects, nextCursor, nil
}

func ParseError(data *DatasourceResponse, endpoint string) *framework.Error {
	errorMessages := make([]string, 0, len(data.Errors)+1)

	errorMessages = append(errorMessages, *data.Error)

	for _, err := range data.Errors {
		errorMessages = append(errorMessages, fmt.Sprintf("Error: %s, Field: %s, Location: %s",
			err.Error, err.Field, err.Location))
	}

	return &framework.Error{
		Message: fmt.Sprintf("Failed to query the datasource: %s.\nGot errors: %v.",
			endpoint, strings.Join(errorMessages, "\n")),
		Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
	}
}
