// Copyright 2026 SGNL.ai, Inc.

package victorops

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
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"go.uber.org/zap"
)

const (
	IncidentReport string = "IncidentReport"
	User           string = "User"
)

// ValidEntityExternalIDs is a map of valid external IDs of entities that can be queried.
// Incidents doc:
//
//	https://portal.victorops.com/api-docs/#!/Reporting/get_api_reporting_v2_incidents
//
// Users doc:
//
//	https://portal.victorops.com/api-docs/#!/Users/get_api_public_v2_user
var ValidEntityExternalIDs = map[string]Entity{
	IncidentReport: {
		uniqueIDAttrExternalID: "incidentNumber",
		endpoint:               "/api-reporting/v2/incidents",
	},
	User: {
		uniqueIDAttrExternalID: "username",
		endpoint:               "/api-public/v2/user",
	},
}

// Entity contains entity specific information, such as the entity's unique ID attribute and the
// endpoint to query that entity.
type Entity struct {
	// uniqueIDAttrExternalID is the external ID of the entity's uniqueId attribute.
	uniqueIDAttrExternalID string

	// endpoint is the API endpoint to query the entity.
	endpoint string
}

// Datasource implements the VictorOps Client interface to allow querying the VictorOps datasource.
type Datasource struct {
	Client *http.Client
}

// NewClient instantiates and returns a new VictorOps Client used to query the VictorOps datasource.
func NewClient(client *http.Client) Client {
	return &Datasource{
		Client: client,
	}
}

// GetPage makes a request to the VictorOps datasource to get a page of JSON objects.
func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	logger := zaplogger.FromContext(ctx).With(
		fields.RequestEntityExternalID(request.EntityExternalID),
		fields.RequestPageSize(request.PageSize),
	)

	logger.Info("Starting datasource request")

	entity := ValidEntityExternalIDs[request.EntityExternalID]

	url := constructURL(request, entity)

	// Timeout API calls that take longer than the configured timeout.
	apiCtx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(apiCtx, http.MethodGet, url, nil)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Address in datasource config is an invalid URL: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// VictorOps uses two custom headers for authentication instead of standard Basic auth.
	// The API ID and API key are sourced from request.Auth.Basic.Username and Password respectively.
	// See ValidateGetPageRequest for documentation of this mapping.
	req.Header.Add("X-VO-Api-Id", request.APIId)
	req.Header.Add("X-VO-Api-Key", request.APIKey)
	req.Header.Add("Accept", "application/json")

	logger.Info("Sending request to datasource", fields.RequestURL(url))

	res, err := d.Client.Do(req)
	if err != nil {
		logger.Error("Request to datasource failed",
			fields.RequestURL(url),
			fields.SGNLEventTypeError(),
			zap.Error(err),
		)

		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute VictorOps request: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds),
		)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to read VictorOps response: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	response := &Response{
		StatusCode:       res.StatusCode,
		RetryAfterHeader: res.Header.Get("Retry-After"),
	}

	if res.StatusCode != http.StatusOK {
		logger.Error("Datasource responded with an error",
			fields.RequestURL(url),
			fields.ResponseStatusCode(response.StatusCode),
			fields.ResponseRetryAfterHeader(response.RetryAfterHeader),
			fields.ResponseBody(body),
			fields.SGNLEventTypeError(),
		)

		return response, nil
	}

	var objects []map[string]any

	var nextCursor *pagination.CompositeCursor[int64]

	switch request.EntityExternalID {
	case IncidentReport:
		objects, nextCursor, err = parseIncidentsResponse(body)
	case User:
		objects, nextCursor, err = parseUsersResponse(body, request.PageSize, request.Cursor)
	}

	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to parse VictorOps %s response: %v.", request.EntityExternalID, err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	response.Objects = objects
	response.NextCursor = nextCursor

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(response.StatusCode),
		fields.ResponseObjectCount(len(response.Objects)),
		fields.ResponseNextCursor(response.NextCursor),
	)

	return response, nil
}

// constructURL constructs the VictorOps URL for the given request and entity.
func constructURL(request *Request, entity Entity) string {
	url := request.BaseURL + entity.endpoint

	switch request.EntityExternalID {
	case IncidentReport:
		var offset int64
		if request.Cursor != nil && request.Cursor.Cursor != nil {
			offset = *request.Cursor.Cursor
		}

		url += "?offset=" + strconv.FormatInt(offset, 10) +
			"&limit=" + strconv.FormatInt(request.PageSize, 10)

		if request.QueryParameters != "" {
			url += "&" + request.QueryParameters
		}
	default:
		// Users endpoint has no pagination parameters,
		// but may still have user-configured query parameters.
		if request.QueryParameters != "" {
			url += "?" + request.QueryParameters
		}
	}

	return url
}

// incidentsResponse represents the top-level response from the VictorOps Incidents API.
type incidentsResponse struct {
	Incidents []map[string]any `json:"incidents"`
	Total     int64            `json:"total"`
	Offset    int64            `json:"offset"`
	Limit     int64            `json:"limit"`
}

// parseIncidentsResponse parses the VictorOps Incidents API response and computes the next cursor.
func parseIncidentsResponse(
	body []byte,
) ([]map[string]any, *pagination.CompositeCursor[int64], error) {
	var data incidentsResponse

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal incidents response: %w", err)
	}

	// Use the response's offset and the actual number of incidents returned to compute the next cursor.
	// This is safer than using the request's pageSize, as the API may return fewer items than requested.
	receivedCount := int64(len(data.Incidents))

	if receivedCount == 0 {
		return data.Incidents, nil, nil
	}

	nextOffset := data.Offset + receivedCount

	var nextCursor *pagination.CompositeCursor[int64]

	if nextOffset < data.Total {
		nextCursor = &pagination.CompositeCursor[int64]{
			Cursor: &nextOffset,
		}
	}

	return data.Incidents, nextCursor, nil
}

// usersResponse represents the top-level response from the VictorOps Users API.
type usersResponse struct {
	Users []map[string]any `json:"users"`
}

// parseUsersResponse parses the VictorOps Users API response.
// The Users API returns all users in a single response with no server-side pagination,
// so client-side pagination is applied by slicing the full results based on pageSize and cursor.
func parseUsersResponse(
	body []byte, pageSize int64, cursor *pagination.CompositeCursor[int64],
) ([]map[string]any, *pagination.CompositeCursor[int64], error) {
	var data usersResponse

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal users response: %w", err)
	}

	var offset int64
	if cursor != nil && cursor.Cursor != nil {
		offset = *cursor.Cursor
	}

	// If offset is beyond available data, return empty results with no next cursor.
	if offset >= int64(len(data.Users)) {
		return []map[string]any{}, nil, nil
	}

	end := offset + pageSize
	if end > int64(len(data.Users)) {
		end = int64(len(data.Users))
	}

	var nextCursor *pagination.CompositeCursor[int64]

	if end < int64(len(data.Users)) {
		nextCursor = &pagination.CompositeCursor[int64]{
			Cursor: &end,
		}
	}

	return data.Users[offset:end], nextCursor, nil
}
