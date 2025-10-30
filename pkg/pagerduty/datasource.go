// Copyright 2025 SGNL.ai, Inc.
package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

const (
	Users   string = "users"
	Teams   string = "teams"
	Members string = "members"
	OnCalls string = "oncalls"
)

// Datasource implements the PagerDuty Client interface to allow querying the PagerDuty datasource.
type Datasource struct {
	Client *http.Client
}

// NewClient instantiates and returns a new PagerDuty Client used to query the PagerDuty datasource.
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

	cursor := request.Cursor

	if cursor == nil || cursor.Cursor == nil {
		var zero int64

		switch request.EntityExternalID {
		case Members:
			var teamCursor *int64
			if cursor != nil {
				teamCursor = cursor.CollectionCursor
			}

			if teamCursor == nil {
				teamCursor = &zero
			}

			// We have no more members to query for the last requested team,
			// or this is a request for the first page.
			// Get the ID of the next team.
			pagerDutyTeamsReq := &Request{
				BaseURL:               request.BaseURL,
				Token:                 request.Token,
				PageSize:              1,
				Cursor:                &pagination.CompositeCursor[int64]{Cursor: teamCursor},
				EntityExternalID:      Teams,
				RequestTimeoutSeconds: request.RequestTimeoutSeconds,
			}

			teamsRes, err := d.GetPage(ctx, pagerDutyTeamsReq)
			if err != nil {
				return nil, err
			}

			// If we fail to get teams, then we can't get members. Terminate and return the error.
			if teamsRes.StatusCode != http.StatusOK {
				return teamsRes, nil
			}

			// There are no more teams. Return an empty last page.
			if len(teamsRes.Objects) == 0 {
				return &Response{
					StatusCode: 202,
				}, nil
			}

			firstTeamIDAsAny, found := teamsRes.Objects[0][UniqueIDAttribute]
			if !found {
				return nil, &framework.Error{
					Message: fmt.Sprintf("PagerDuty team object contains no %s field.", UniqueIDAttribute),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
				}
			}

			teamID, ok := firstTeamIDAsAny.(string)
			if !ok {
				return nil, &framework.Error{
					Message: fmt.Sprintf("Failed to convert PagerDuty team object %s field to string.", UniqueIDAttribute),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
				}
			}

			cursor = &pagination.CompositeCursor[int64]{
				CollectionID: &teamID,
				Cursor:       &zero,
			}

			if teamsRes.NextCursor == nil {
				cursor.CollectionCursor = nil
			} else {
				cursor.CollectionCursor = teamsRes.NextCursor.Cursor
			}

		default:
			// The request is for the first page, initialize the cursor.
			cursor = &pagination.CompositeCursor[int64]{
				Cursor: &zero,
			}
		}
	}

	validationErr := pagination.ValidateCompositeCursor(
		cursor,
		request.EntityExternalID,
		request.EntityExternalID == Members,
	)
	if validationErr != nil {
		return nil, validationErr
	}

	var sb strings.Builder
	// URL format:
	// request.BaseURL + "/" + {resource path} + "?offset=" + request.Cursor + "&limit=" + request.PageSize
	// where len("/") + len("?offset=") + len("&limit=") == 16.
	// {resource path} is computed below.
	cursorStr := strconv.FormatInt(*cursor.Cursor, 10)
	pageSizeStr := strconv.FormatInt(request.PageSize, 10)
	sb.Grow(len(request.BaseURL) + len(cursorStr) + len(pageSizeStr) + 16)
	sb.WriteString(request.BaseURL)
	sb.WriteRune('/')

	if request.EntityExternalID == Members {
		// If we sync team members, the endpoint becomes the following:
		// baseURL/ + teams/:teamID/members + query params
		escapedTeamID := url.PathEscape(*cursor.CollectionID)
		sb.Grow(len(escapedTeamID) + 14)
		sb.WriteString(Teams)
		sb.WriteRune('/')
		sb.WriteString(escapedTeamID)
		sb.WriteRune('/')
		sb.WriteString(Members)
	} else {
		// Otherwise, baseURL/ + :EntityExternalID + query params
		sb.WriteString(request.EntityExternalID)
	}

	sb.WriteString("?offset=")
	sb.WriteString(cursorStr)
	sb.WriteString("&limit=")
	sb.WriteString(pageSizeStr)

	if len(request.AdditionalQueryParameters) != 0 {
		// Extract the additional query parameters for this entity.
		// If it does not exist, entityAdditionalQueryParameters is an empty map and no query parameters
		// are appended.
		entityAdditionalQueryParameters := request.AdditionalQueryParameters[request.EntityExternalID]

		for queryParam, values := range entityAdditionalQueryParameters {
			// e.g. Append &includes[]=foo&includes[]=bar.
			for _, v := range values {
				sb.WriteRune('&')
				sb.WriteString(url.QueryEscape(queryParam))
				sb.WriteRune('=')
				sb.WriteString(url.QueryEscape(v))
			}
		}
	}

	requestURL := sb.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Adapter generated an invalid URL: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Timeout API calls that take longer than the configured timeout.
	apiCtx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	req = req.WithContext(apiCtx)

	req.Header.Add("Authorization", request.Token)

	logger.Info("Sending request to datasource", fields.RequestURL(requestURL))

	res, err := d.Client.Do(req)
	if err != nil {
		logger.Error("Request to datasource failed",
			fields.RequestURL(requestURL),
			fields.SGNLEventTypeError(),
			zap.Error(err),
		)

		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute PagerDuty request: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds),
		)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to read PagerDuty response: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	response := &Response{
		StatusCode:       res.StatusCode,
		RetryAfterHeader: res.Header.Get("Retry-After"),
	}

	if res.StatusCode != http.StatusOK {
		logger.Error("Datasource responded with an error",
			fields.RequestURL(requestURL),
			fields.ResponseStatusCode(response.StatusCode),
			fields.ResponseRetryAfterHeader(response.RetryAfterHeader),
			fields.ResponseBody(body),
			fields.SGNLEventTypeError(),
		)

		return response, nil
	}

	objects, nextCursor, frameworkErr := ParseResponse(body, request.EntityExternalID, request.PageSize, *cursor.Cursor)
	if frameworkErr != nil {
		return nil, frameworkErr
	}

	response.NextCursor = &pagination.CompositeCursor[int64]{
		Cursor: nextCursor,
	}

	switch request.EntityExternalID {
	case Members:
		// We must create a unique ID for Members that is a
		// combination of the teamID and userID, since a user may belong
		// to multiple teams.
		teamMemberObjects := make([]map[string]any, 0, len(objects))

		for _, object := range objects {
			var teamMemberObject = make(map[string]any, 4)

			userObject, ok := object["user"].(map[string]any)
			if !ok {
				return nil, &framework.Error{
					Message: "Failed to parse user field in PagerDuty team members response as map[string]any.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			userID, ok := userObject[UniqueIDAttribute].(string)
			if !ok {
				return nil, &framework.Error{
					Message: fmt.Sprintf(
						"Failed to parse %s field in PagerDuty team members object as string.",
						UniqueIDAttribute,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			teamMemberObject["userId"] = userID
			teamMemberObject["teamId"] = *cursor.CollectionID

			if object["role"] != nil {
				teamMemberObject["role"] = object["role"]
			}

			teamMemberObject[UniqueIDAttribute] = *cursor.CollectionID + "-" + userID

			teamMemberObjects = append(teamMemberObjects, teamMemberObject)
		}

		objects = teamMemberObjects
		response.NextCursor.CollectionID = cursor.CollectionID
		response.NextCursor.CollectionCursor = cursor.CollectionCursor
	case OnCalls:
		// OnCall objects do not have a unique "id" attribute.
		// Instead, we create it using the following fields: {escalation_policy.id}-{user.id}-{start}-{end}.
		// If start and end are null, they are set to empty strings.
		for _, object := range objects {
			escalationPolicyMap, ok := object["escalation_policy"].(map[string]any)
			if !ok {
				return nil, &framework.Error{
					Message: "Failed to parse a PagerDuty OnCall object's escalation_policy field as map[string]any.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			escalationPolicyID, ok := escalationPolicyMap[UniqueIDAttribute].(string)
			if !ok {
				return nil, &framework.Error{
					Message: fmt.Sprintf(
						"Failed to parse a field in a PagerDuty OnCall object's escalation_policy object as string: %s.",
						UniqueIDAttribute,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			userMap, ok := object["user"].(map[string]any)
			if !ok {
				return nil, &framework.Error{
					Message: "Failed to parse a PagerDuty OnCall object's user field as map[string]any.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			userID, ok := userMap[UniqueIDAttribute].(string)
			if !ok {
				return nil, &framework.Error{
					Message: fmt.Sprintf(
						"Failed to parse a field in a PagerDuty OnCall object's user object as string: %s.",
						UniqueIDAttribute,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			startDate, endDate := object["start"], object["end"]

			if startDate == nil {
				startDate = ""
			}

			if endDate == nil {
				endDate = ""
			}

			startDateString, ok := startDate.(string)
			if !ok {
				return nil, &framework.Error{
					Message: fmt.Sprintf(
						"Failed to parse a PagerDuty OnCall object's start field as string: %v.",
						startDate,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			endDateString, ok := endDate.(string)
			if !ok {
				return nil, &framework.Error{
					Message: fmt.Sprintf(
						"Failed to parse a PagerDuty OnCall object's end field as string: %v.",
						endDate,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			object[UniqueIDAttribute] = escalationPolicyID + "-" + userID + "-" + startDateString + "-" + endDateString
		}
	}

	response.Objects = objects

	// We must return a response with nil NextCursor to indicate a full sync has completed.
	// A full sync completes depending on the following:
	// If we aren't syncing team members, then if the computed nextCursor is nil, we have
	// reached the end of a sync.
	// If we are syncing team members, then if the next CollectionCursor is nil, we have reached the
	// end of a sync because we have iterated through all teams.
	// These two conditions can be combined because a CollectionCursor should always be nil if we aren't
	// syncing team members.
	if nextCursor == nil && cursor.CollectionCursor == nil {
		response.NextCursor = nil
	}

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(response.StatusCode),
		fields.ResponseObjectCount(len(response.Objects)),
		fields.ResponseNextCursor(response.NextCursor),
	)

	return response, nil
}

// ParseResponse parses the response body into an array of objects with the cursor to the next page.
// It assumes the response is structured as follows:
// {"{entityExternalID}": []objects, "more": bool, ...}.
func ParseResponse(
	body []byte, entityExternalID string, pageSize int64, cursor int64,
) (objects []map[string]any, nextCursor *int64, err *framework.Error) {
	var data map[string]any

	unmarshalErr := json.Unmarshal(body, &data)
	if unmarshalErr != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal PagerDuty response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	rawData, found := data[entityExternalID]
	if !found {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Field missing in PagerDuty response: %s.", entityExternalID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	rawObjects, ok := rawData.([]any)
	if !ok {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf(
				"Entity %s field exists in PagerDuty response but field value is not a list of objects: %T.",
				entityExternalID,
				rawData,
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	parsedObjects, parserErr := parseObjects(rawObjects, entityExternalID)
	if parserErr != nil {
		return nil, nil, parserErr
	}

	more, found := data["more"]
	if found {
		moreData, ok := more.(bool)
		if !ok {
			return nil, nil, &framework.Error{
				Message: fmt.Sprintf(
					"Field more exists in PagerDuty response but field value is not a bool: %T.",
					more,
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		if !moreData {
			return parsedObjects, nil, nil
		}
	}

	nextCursor = pagination.GetNextCursorFromPageSize(len(parsedObjects), pageSize, cursor)

	return parsedObjects, nextCursor, nil
}

// parseObjects parses []any into []map[string]any. If any object in the slice is not a map[string]any,
// a framework.Error is returned.
func parseObjects(objects []any, entityExternalID string) ([]map[string]any, *framework.Error) {
	parsedObjects := make([]map[string]any, 0, len(objects))

	for _, object := range objects {
		parsedObject, ok := object.(map[string]any)
		if !ok {
			return nil, &framework.Error{
				Message: fmt.Sprintf(
					"An object in Entity: %s could not be parsed. Expected: map[string]any. Got: %T.",
					entityExternalID,
					object,
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		parsedObjects = append(parsedObjects, parsedObject)
	}

	return parsedObjects, nil
}
