// Copyright 2025 SGNL.ai, Inc.
package crowdstrike

// We ingest CrowdStrike data using both GraphQL and REST APIs.
// This file contains the functions and structs that are used to interact with the CrowdStrike REST APIs.
// The REST API is two-level. A list endpoint is used to list entity IDs and
// a get endpoint is used to get detailed metadata of a specific entity.
// If a REST API has only one endpoint, it is considered as a get endpoint.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	"go.uber.org/zap"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

type ErrorItem struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PaginationInfo struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

type MetaFields struct {
	PaginationInfo PaginationInfo `json:"pagination"`
}

type ListResourceResponse struct {
	Meta      MetaFields  `json:"meta"`
	Resources []string    `json:"resources"`
	Errors    []ErrorItem `json:"errors"`
}

// "devices/queries/devices-scroll/v1" endpoint has a string offset.
type ScrollPaginationInfo struct {
	Offset string `json:"offset"`
	Limit  int    `json:"limit"`
	Total  int    `json:"total"`
}

type ScrollMetaFields struct {
	PaginationInfo ScrollPaginationInfo `json:"pagination"`
}

type ListScrollResourceResponse struct {
	Meta      ScrollMetaFields `json:"meta"`
	Resources []string         `json:"resources"`
	Errors    []ErrorItem      `json:"errors"`
}

type DetailedResourceRequestBody struct {
	Identifiers []string `json:"ids"`
}

type DetailedResourceResponse struct {
	Meta      MetaFields       `json:"meta"`
	Resources []map[string]any `json:"resources"`
	Errors    []ErrorItem      `json:"errors"`
}

type AlertsRequestBody struct {
	Limit  int     `json:"limit"`
	After  *string `json:"after,omitempty"`
	Filter *string `json:"filter,omitempty"`
	Sort   *string `json:"sort,omitempty"`
}

type AlertsResponse struct {
	Meta      AlertsMeta       `json:"meta"`
	Resources []map[string]any `json:"resources"`
	Errors    []ErrorItem      `json:"errors"`
}

type AlertsMeta struct {
	Pagination AlertsPagination `json:"pagination"`
}

type AlertsPagination struct {
	After string `json:"after"`
	Total int    `json:"total"`
}

// generateRequestBodyBytes creates the request body bytes based on the entity type and request parameters.
func generateRequestBodyBytes(request *Request, resourceIDs []string) ([]byte, error) {
	if request.EntityExternalID == Alerts {
		var after *string
		if request.RESTCursor != nil && request.RESTCursor.Cursor != nil {
			after = request.RESTCursor.Cursor
		}

		reqBody := &AlertsRequestBody{
			Limit:  int(request.PageSize),
			After:  after,
			Filter: request.Filter,
		}

		return json.Marshal(reqBody)
	}

	reqBody := &DetailedResourceRequestBody{
		Identifiers: resourceIDs,
	}

	return json.Marshal(reqBody)
}

func (d *Datasource) getRESTPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	logger := zaplogger.FromContext(ctx).With(
		fields.RequestEntityExternalID(request.EntityExternalID),
		fields.RequestPageSize(request.PageSize),
	)

	logger.Info("Starting datasource request")

	// Fetch all resourceIDs before fetching the detailed information, if applicable.
	// For Alerts we have a combined API which doesn't require to fetch resource IDs separately.
	resourceIDs, nextCursor, httpResp, listErr := d.getResourceIDs(ctx, request)
	if listErr != nil {
		return nil, listErr
	}

	if httpResp != nil && httpResp.StatusCode != http.StatusOK {
		return &Response{
			StatusCode:       httpResp.StatusCode,
			RetryAfterHeader: httpResp.Header.Get("Retry-After"),
		}, nil
	}

	// 1 beyond the last page. See why in `parseListScrollResponse`
	if request.EntityExternalID == Device && len(resourceIDs) == 0 && nextCursor == nil {
		return &Response{
			StatusCode: httpResp.StatusCode,
		}, nil
	}

	// Use resourceIDs to fetch detailed information, if applicable.
	bodyBytes, marshalErr := generateRequestBodyBytes(request, resourceIDs)
	if marshalErr != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to marshal the request body: %v.", marshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	endpointInfo := EntityExternalIDToEndpoint[request.EntityExternalID]

	url, urlErr := ConstructRESTEndpoint(request, endpointInfo.GetEndpoint)
	if urlErr != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to construct the endpoint: %v.", urlErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Timeout API calls that take longer than the configured timeout.
	apiCtx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(apiCtx, http.MethodPost, *url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to create request to datasource: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	req.Header.Add("Authorization", request.Token)
	req.Header.Set("Content-Type", "application/json")

	logger.Info("Sending HTTP request to datasource", fields.URL(*url))

	res, err := d.Client.Do(req)
	if err != nil {
		logger.Error("HTTP request to datasource failed", fields.URL(*url), fields.SGNLEventTypeError(), zap.Error(err))

		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute CrowdStrike request: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds),
		)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to read CrowdStrike response: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	response := &Response{
		StatusCode:       res.StatusCode,
		RetryAfterHeader: res.Header.Get("Retry-After"),
	}

	if res.StatusCode != http.StatusOK {
		logger.Error("Datasource request failed",
			fields.ResponseStatusCode(res.StatusCode),
			fields.ResponseRetryAfterHeader(response.RetryAfterHeader),
			fields.ResponseBody(body),
			fields.SGNLEventTypeError(),
		)

		return response, nil
	}

	var (
		objects      []map[string]any
		frameworkErr *framework.Error
	)

	if request.EntityExternalID == Alerts {
		objects, nextCursor, frameworkErr = parseAlertsResponse(body)
	} else {
		objects, frameworkErr = parseDetailedResponse(body)
	}

	if frameworkErr != nil {
		return nil, frameworkErr
	}

	response.NextRESTCursor = nextCursor
	response.Objects = objects

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(response.StatusCode),
		fields.ResponseObjectCount(len(response.Objects)),
		fields.ResponseNextCursor(response.NextRESTCursor),
	)

	return response, nil
}

func (d *Datasource) getResourceIDs(ctx context.Context, request *Request) (
	[]string,
	*pagination.CompositeCursor[string],
	*http.Response,
	*framework.Error,
) {
	endpointInfo := EntityExternalIDToEndpoint[request.EntityExternalID]

	if request.EntityExternalID == Alerts && endpointInfo.ListEndpoint == "" {
		return []string{}, nil, nil, nil
	}

	url, urlErr := ConstructRESTEndpoint(request, endpointInfo.ListEndpoint)
	if urlErr != nil {
		return nil, nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to construct the endpoint: %v.", urlErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, *url, nil)
	if err != nil {
		return nil, nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to create request to datasource: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Timeout API calls that take longer than the configured timeout.
	apiCtx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	req = req.WithContext(apiCtx)

	req.Header.Add("Authorization", request.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := d.Client.Do(req)
	if err != nil {
		return nil, nil, nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute CrowdStrike request: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds),
		)
	}

	statusCode := res.StatusCode

	if statusCode != http.StatusOK {
		return nil, nil, res, nil
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, res, &framework.Error{
			Message: fmt.Sprintf("Failed to read CrowdStrike response: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	var (
		resourceIDs []string
		nextCursor  *pagination.CompositeCursor[string]
		listErr     *framework.Error
	)

	if ValidRESTEntityExternalIDs[request.EntityExternalID].UseIntCursor {
		resourceIDs, nextCursor, listErr = parseListResponse(body, request)
		if listErr != nil {
			return nil, nil, res, listErr
		}
	} else {
		resourceIDs, nextCursor, listErr = parseListScrollResponse(body, request)
		if listErr != nil {
			return nil, nil, res, listErr
		}
	}

	if len(resourceIDs) == 0 {
		return nil, nil, res, nil
	}

	return resourceIDs, nextCursor, res, nil
}

// parseListResponse parses the list response from the CrowdStrike REST API.
// The metadata contains pagination information, which is used to determine the next cursor.
func parseListResponse(body []byte, request *Request) (
	objects []string,
	nextCursor *pagination.CompositeCursor[string],
	err *framework.Error,
) {
	var data *ListResourceResponse

	if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil || data == nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the resource IDs in datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if len(data.Errors) != 0 {
		return nil, nil, ParseError(data.Errors)
	}

	if data.Resources == nil {
		return nil, nil, &framework.Error{
			Message: "Missing resource IDs in the datasource response.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	nextOffset := request.PageSize
	total := int64(data.Meta.PaginationInfo.Total)

	// If the cursor is not nil, increment the next offset by the last offset value.
	if request.RESTCursor != nil && request.RESTCursor.Cursor != nil {
		prevOffset, parseErr := strconv.ParseInt(*request.RESTCursor.Cursor, 10, 64)
		if parseErr != nil {
			return nil, nil, &framework.Error{
				Message: fmt.Sprintf("Expected a numeric cursor for entity: %s.", request.EntityExternalID),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_AUTHENTICATION_FAILED,
			}
		}

		nextOffset += prevOffset
	}

	if total > nextOffset {
		offsetStr := strconv.Itoa(int(nextOffset))

		nextCursor = &pagination.CompositeCursor[string]{Cursor: &offsetStr}
	}

	return data.Resources, nextCursor, nil
}

/*
parseListScrollResponse parses the list response from the CrowdStrike REST API that use string offsets.
The metadata contains pagination information, which is used to determine the next cursor.

The /devices/queries/devices-scroll/v1 behaves weirdly. Even if there's no items left, it still returns an offset.
One querying using the offset we get a response like this:

	{
		"meta": {
			"query_time": 0.025342538,
			"pagination": {
				"total": 0,
				"offset": ""
			},
			"powered_by": "device-api",
			"trace_id": "ff9ff514-549c-49cf-a22f-b5d981f162c6"
		},
		"resources": [],
		"errors": []
	}

The caller needs to handle this edge case i.e. getRESTPage.
*/
func parseListScrollResponse(body []byte, _ *Request) (
	objects []string,
	nextCursor *pagination.CompositeCursor[string],
	err *framework.Error,
) {
	var data *ListScrollResourceResponse

	if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil || data == nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the resource IDs in datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if len(data.Errors) != 0 {
		return nil, nil, ParseError(data.Errors)
	}

	if data.Resources == nil {
		return nil, nil, &framework.Error{
			Message: "Missing resource IDs in the datasource response.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	if data.Meta.PaginationInfo.Offset != "" {
		return data.Resources, &pagination.CompositeCursor[string]{
			Cursor: &data.Meta.PaginationInfo.Offset,
		}, nil
	}

	return data.Resources, nil, nil
}

// parseDetailedResponse parses the detailed response from the CrowdStrike REST API.
// The metadata does not contain any pagination information, unlike the list endpoint.
func parseDetailedResponse(body []byte) ([]map[string]any, *framework.Error) {
	var data *DetailedResourceResponse

	if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil || data == nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the detailed resources in datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if len(data.Errors) != 0 {
		return nil, ParseError(data.Errors)
	}

	if data.Resources == nil {
		return nil, &framework.Error{
			Message: "Missing  detailed resources in the datasource response.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return data.Resources, nil
}

// parseAlertsResponse parses the response from the alerts API endpoint.
func parseAlertsResponse(body []byte) (
	objects []map[string]any,
	nextCursor *pagination.CompositeCursor[string],
	err *framework.Error,
) {
	var data *AlertsResponse

	if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil || data == nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the alerts response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if len(data.Errors) != 0 {
		return nil, nil, ParseError(data.Errors)
	}

	if data.Resources == nil {
		return nil, nil, &framework.Error{
			Message: "Missing resources in the alerts response.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Set next cursor if there's an "after" token for pagination
	if data.Meta.Pagination.After != "" {
		// The Combined Alerts API returns a base64-encoded cursor directly
		// Use it as-is to maintain consistency with the CrowdStrike API format
		nextCursor = &pagination.CompositeCursor[string]{
			Cursor: &data.Meta.Pagination.After,
		}
	}

	return data.Resources, nextCursor, nil
}

func ParseError(errors []ErrorItem) *framework.Error {
	errorMessages := make([]string, 0, len(errors)+1)

	for _, err := range errors {
		errorMessages = append(errorMessages, fmt.Sprintf("Code: %d, Message: %s",
			err.Code, err.Message))
	}

	return &framework.Error{
		Message: fmt.Sprintf(
			"Failed to query the datasource.\nGot errors: %v.",
			strings.Join(errorMessages, "\n"),
		),
		Code: api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
	}
}
