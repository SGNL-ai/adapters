// Copyright 2026 SGNL.ai, Inc.

package jira

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	net_url "net/url"
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
	User  string = "User"
	Issue string = "Issue"
	// TODO: Remove this after fully deprecating the legacy Issue endpoint.
	EnhancedIssue string = "EnhancedIssue"
	Group         string = "Group"
	GroupMember   string = "GroupMember"
	Workspace     string = "Workspace"
	Object        string = "Object"

	isLastFieldName     = "isLast"
	nextCursorFieldName = "nextPageToken"
)

var EntityIDToParentCollectionID = map[string]string{
	GroupMember: Group,
	Object:      Workspace,
}

var (
	// ValidEntityExternalIDs is a map of valid external IDs of entities that can be queried.
	// The map value is the Entity struct, which contains the unique ID attribute, the endpoint to query that entity,
	// and a function to parse the response.
	// Users doc:
	//   https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-users/#api-rest-api-3-users-search-get.
	// Issues doc:
	//   https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-search/#api-rest-api-3-search-get.
	// Groups doc:
	//   https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-groups/#api-rest-api-3-group-bulk-get.
	// GroupMembers doc:
	//   https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-groups/#api-rest-api-3-group-member-get.
	// Workspaces doc:
	// nolint:lll
	//   https://developer.atlassian.com/cloud/jira/service-desk/rest/api-group-assets/#api-rest-servicedeskapi-assets-workspace-get.
	// Objects doc:
	// 	 https://developer.atlassian.com/cloud/assets/rest/api-group-object/#api-object-aql-post.
	ValidEntityExternalIDs = map[string]Entity{
		User: {
			uniqueIDAttrExternalID: "accountId",
			endpoint:               "users/search",
			parseResponse:          ParseUsersResponse,
		},
		Issue: {
			uniqueIDAttrExternalID: "id",
			endpoint:               "search",
			parseResponse:          ParseIssuesResponse,
		},
		// EnhancedIssue is a temporary mapping type that is set implicitly based on the
		// EnhancedIssueSearch config setting. This will be removed with that setting eventually and we'll
		// remove this mapping.
		//
		// While it should be possible to specify this EnhancedIssue external entity ID from the UI, this is not
		// recommended (instead set the EnhancedIssueSearch config setting), as this will not be backwards compatible
		// once this is removed.
		//
		// TODO: Remove this after fully deprecating the legacy Issue endpoint.
		EnhancedIssue: {
			uniqueIDAttrExternalID: "id",
			endpoint:               "search/jql",
			parseResponse:          ParseEnhancedIssuesResponse,
		},
		Group: {
			uniqueIDAttrExternalID: "groupId",
			endpoint:               "group/bulk",
			parseResponse:          ParseGroupsResponse,
		},
		GroupMember: {
			// "id" is created by combining the groupId and accountId.
			// It is not returned by the Jira API.
			uniqueIDAttrExternalID: "id",
			endpoint:               "group/member",
			parseResponse:          ParseGroupMembersResponse,
		},
		Workspace: {
			uniqueIDAttrExternalID: "workspaceId",
			endpoint:               "workspace",
			parseResponse:          ParseWorkspacesResponse,
		},
		Object: {
			uniqueIDAttrExternalID: "globalId",
			endpoint:               "object", // Not used.
			parseResponse:          ParseObjectsResponse,
		},
	}
)

type responseParser func(
	body []byte, pageSize int64, cursor string,
) (
	objects []map[string]any, nextCursor *string, err *framework.Error,
)

// Entity contains entity specific information, such as the entity's unique ID attribute and the
// endpoint to query that entity.
type Entity struct {
	// uniqueIDAttrExternalID is the external ID of the entity's uniqueId attribute.
	uniqueIDAttrExternalID string

	// endpoint is the endpoint to query the entity, e.g. "users/search" for users.
	// It does not need to include "/rest/api/3/".
	endpoint string

	// parseResponse is a function that parses the entity response body and returns the objects, the next cursor,
	// and an error.
	parseResponse responseParser
}

// Datasource implements the Jira Client interface to allow querying the Jira datasource.
type Datasource struct {
	Client *http.Client
}

// NewClient instantiates and returns a new Jira Client used to query the Jira datasource.
func NewClient(client *http.Client) Client {
	return &Datasource{
		Client: client,
	}
}

// GetPage makes a request to the Jira datasource to get a page of JSON objects. If a response is received,
// regardless of status code, a Response object is returned with the response body and the status code.
// If the request fails, an appropriate framework.Error is returned.
func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	logger := zaplogger.FromContext(ctx).With(
		fields.RequestEntityExternalID(request.EntityExternalID),
		fields.RequestPageSize(request.PageSize),
	)

	logger.Info("Starting datasource request")

	// ValidateGetPageRequest already checks if the entity exists in the valid entities map.
	entity := ValidEntityExternalIDs[request.EntityExternalID]

	cursor, isEmptyLastPage, cursorErr := d.constructCursor(ctx, request)
	if cursorErr != nil {
		return nil, cursorErr
	}

	if isEmptyLastPage {
		return &Response{
			StatusCode: http.StatusOK,
		}, nil
	}

	validationErr := pagination.ValidateCompositeCursor(
		cursor,
		request.EntityExternalID,
		// Send a bool indicating if the entity is a member of a collection.
		request.EntityExternalID == GroupMember || request.EntityExternalID == Object,
	)
	if validationErr != nil {
		return nil, validationErr
	}

	var req *http.Request

	var err error

	url, err := ConstructURL(request, entity, cursor)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Unable to construct URL: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// Timeout API calls that take longer than the configured timeout.
	apiCtx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	switch request.EntityExternalID {
	case Object:
		// The QL query must be included in the request body of a POST request.
		// Contrary to the IssuesJQL query, which can be included as a query parameter in the URL.
		jsonValue, marshalErr := json.Marshal(map[string]any{
			"qlQuery": request.ObjectsQLQuery,
		})
		if marshalErr != nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf("Failed to marshal Jira request body due to QL query: %v.", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		req, err = http.NewRequestWithContext(apiCtx, http.MethodPost, url, bytes.NewBuffer(jsonValue))
		req.Header.Add("Content-Type", "application/json")
	default:
		req, err = http.NewRequestWithContext(apiCtx, http.MethodGet, url, nil)
	}

	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Address in datasource config is an invalid URL: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	req.Header.Add("Authorization", basicAuth(request.Username, request.Password))

	logger.Info("Sending request to datasource", fields.RequestURL(url))

	res, err := d.Client.Do(req)
	if err != nil {
		logger.Error("Request to datasource failed",
			fields.RequestURL(url),
			fields.SGNLEventTypeError(),
			zap.Error(err),
		)

		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute Jira request: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds),
		)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to read Jira response: %v.", err),
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

	objects, nextCursor, frameworkErr := entity.parseResponse(body, request.PageSize, *cursor.Cursor)
	if frameworkErr != nil {
		return nil, frameworkErr
	}

	response.NextCursor = &pagination.CompositeCursor[string]{
		Cursor: nextCursor,
	}

	switch request.EntityExternalID {
	case GroupMember:
		// We must create a unique ID for GroupMembers that is a
		// combination of the groupId and accountId, since a user may belong
		// to multiple groups.
		userUniqueIDAttrExternalID := ValidEntityExternalIDs[User].uniqueIDAttrExternalID
		groupID := *cursor.CollectionID
		groupMemberUniqueIDAttrExternalID := ValidEntityExternalIDs[GroupMember].uniqueIDAttrExternalID

		for _, object := range objects {
			userID, ok := object[userUniqueIDAttrExternalID].(string)
			if !ok {
				return nil, &framework.Error{
					Message: fmt.Sprintf(
						"Failed to parse %s field in Jira GroupMember response as string.",
						userUniqueIDAttrExternalID,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			object["groupId"] = groupID
			object[groupMemberUniqueIDAttrExternalID] = groupID + "-" + userID
		}

		response.NextCursor.CollectionID = cursor.CollectionID
		response.NextCursor.CollectionCursor = cursor.CollectionCursor
	case Object:
		// The Jira API returns Objects with a globalId, which is already a combination of the workspace + object ID.
		// So the globalId can be used as the unique ID. e.g. globalId: f1668d0c-828c-470c-b7d1-8c4f48cd345a:88.
		response.NextCursor.CollectionID = cursor.CollectionID
		response.NextCursor.CollectionCursor = cursor.CollectionCursor
	}

	response.Objects = objects

	// We must return a response with nil NextCursor to indicate a full sync has completed.
	// A full sync completes depending on the following:
	// If we aren't syncing groupMembers, then if the computed nextCursor is nil, we have
	// reached the end of a sync.
	// If we are syncing groupMembers, then if the next CollectionCursor is nil, we have reached the
	// end of a sync because we have iterated through all groups.
	// These two conditions can be combined because a CollectionCursor should always be nil if we aren't
	// syncing groupMembers.
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

// constructCursor constructs the page's composite cursor. For entities that are
// members of a collection (and therefore require collection IDs in the request) a request
// is made to retrieve the collection ID. If there are no more collections, a bool is returned indicating
// this is the last page and it's empty.
func (d *Datasource) constructCursor(
	ctx context.Context, request *Request,
) (*pagination.CompositeCursor[string], bool, *framework.Error) {
	cursor := request.Cursor

	// nolint:nestif
	if cursor == nil || cursor.Cursor == nil {
		zero := "0"

		if parentCollectionEntityID, hasParentCollection :=
			EntityIDToParentCollectionID[request.EntityExternalID]; hasParentCollection {
			var collectionCursor *string
			if cursor != nil {
				collectionCursor = cursor.CollectionCursor
			}

			if collectionCursor == nil {
				collectionCursor = &zero
			}

			// We have no more members to query for the last requested collection,
			// or this is a request for the first page.
			// Get the ID of the next collection.
			nextCollectionReq := &Request{
				BaseURL:               request.BaseURL,
				Username:              request.Username,
				Password:              request.Password,
				PageSize:              1,
				Cursor:                &pagination.CompositeCursor[string]{Cursor: collectionCursor},
				EntityExternalID:      parentCollectionEntityID,
				RequestTimeoutSeconds: request.RequestTimeoutSeconds,
			}

			nextCollectionRes, err := d.GetPage(ctx, nextCollectionReq)
			if err != nil {
				return nil, false, err
			}

			// There are no more collections. Return a bool indicating this was the last page.
			if len(nextCollectionRes.Objects) == 0 {
				return nil, true, nil
			}

			collectionUniqueID := ValidEntityExternalIDs[parentCollectionEntityID].uniqueIDAttrExternalID

			firstCollectionIDRaw, found := nextCollectionRes.Objects[0][collectionUniqueID]
			if !found {
				return nil, false, &framework.Error{
					Message: fmt.Sprintf("Jira %s object contains no %s field.", parentCollectionEntityID, collectionUniqueID),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
				}
			}

			firstCollectionID, ok := firstCollectionIDRaw.(string)
			if !ok {
				return nil, false, &framework.Error{
					Message: fmt.Sprintf(
						"Failed to convert Jira %s object %s field to string.",
						parentCollectionEntityID,
						collectionUniqueID,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
				}
			}

			cursor = &pagination.CompositeCursor[string]{
				CollectionID: &firstCollectionID,
				Cursor:       &zero,
			}

			if nextCollectionRes.NextCursor == nil {
				cursor.CollectionCursor = nil
			} else {
				cursor.CollectionCursor = nextCollectionRes.NextCursor.Cursor
			}
		} else {
			// The request is for the first page, initialize the cursor.
			cursor = &pagination.CompositeCursor[string]{
				Cursor: &zero,
			}
		}
	}

	return cursor, false, nil
}

func ParseUsersResponse(
	body []byte, pageSize int64, cursor string,
) (objects []map[string]any, nextCursor *string, err *framework.Error) {
	// Users response contains a list of User objects: []User.
	unmarshalErr := json.Unmarshal(body, &objects)
	if unmarshalErr != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal Jira users response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	cursorInt, convErr := strconv.ParseInt(cursor, 10, 64)
	if convErr != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf(
				"Failed to parse cursor as an int (%s): %s.",
				cursor,
				convErr.Error(),
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	nextCursorInt := pagination.GetNextCursorFromPageSize(len(objects), pageSize, cursorInt)

	if nextCursorInt != nil {
		tmp := strconv.FormatInt(*nextCursorInt, 10)
		nextCursor = &tmp
	}

	return objects, nextCursor, nil
}

// TODO: Until we fully deprecate the legacy issue search endpoint, any changes to this func should be copied to
// ParseEnhancedIssuesResponse.
func ParseIssuesResponse(
	body []byte, pageSize int64, cursor string,
) (objects []map[string]any, nextCursor *string, err *framework.Error) {
	// Issues response contains a single object map, with the top level field "issues": {"issues": []Issue, ...}.
	// First unmarshal into single object, then extract the "issues" field.
	var data map[string]any

	unmarshalErr := json.Unmarshal(body, &data)
	if unmarshalErr != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal Jira issues response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	_, found := data["issues"]
	if !found {
		return nil, nil, &framework.Error{
			Message: "Field missing in Jira issues response: issues.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	issues, ok := data["issues"].([]any)
	if !ok {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf(
				"Entity field exists in Jira issues response but field value is not a list of objects: %T.",
				data["issues"],
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	objects, parserErr := parseObjects(Issue, issues)
	if parserErr != nil {
		return nil, nil, parserErr
	}

	cursorInt, convErr := strconv.ParseInt(cursor, 10, 64)
	if convErr != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf(
				"Failed to parse cursor as an int (%s): %s.",
				cursor,
				convErr.Error(),
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	nextCursorInt := pagination.GetNextCursorFromPageSize(len(objects), pageSize, cursorInt)

	if nextCursorInt != nil {
		tmp := strconv.FormatInt(*nextCursorInt, 10)
		nextCursor = &tmp
	}

	return objects, nextCursor, nil
}

func ParseEnhancedIssuesResponse(
	body []byte, _ int64, _ string,
) (objects []map[string]any, nextCursor *string, err *framework.Error) {
	// Issues response contains a single object map, with the top level field "issues": {"issues": []Issue, ...}.
	// First unmarshal into single object, then extract the "issues" field.
	var data map[string]any

	unmarshalErr := json.Unmarshal(body, &data)
	if unmarshalErr != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal Jira issues response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	_, found := data["issues"]
	if !found {
		return nil, nil, &framework.Error{
			Message: "Field missing in Jira issues response: issues.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	issues, ok := data["issues"].([]any)
	if !ok {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf(
				"Entity field exists in Jira issues response but field value is not a list of objects: %T.",
				data["issues"],
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	objects, parserErr := parseObjects(Issue, issues)
	if parserErr != nil {
		return nil, nil, parserErr
	}

	// Optimization:
	// If the response contains the lastPageFieldName field, and it's a bool, and it's true,
	// the current page is the last page and there is no need to compute the next cursor.
	isLast, found := data[isLastFieldName]
	if found {
		isLastPage, ok := isLast.(bool)
		if !ok {
			return nil, nil, &framework.Error{
				Message: fmt.Sprintf(
					"Field %s exists in Jira %s response but field value is not a bool: %T.",
					isLastFieldName,
					Issue,
					isLast,
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		if isLastPage {
			return objects, nil, nil
		}
	}

	foundNextCursor, found := data[nextCursorFieldName]
	if found {
		strNextCursor, ok := foundNextCursor.(string)
		if !ok {
			return nil, nil, &framework.Error{
				Message: fmt.Sprintf(
					"Field %s exists in Jira %s response but field value is not a string: %T.",
					nextCursorFieldName,
					Issue,
					foundNextCursor,
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		nextCursor = &strNextCursor

		return objects, nextCursor, nil
	}

	// If no cursor is found above, treat this as the last page. We don't expect this code path to execute since
	// this would mean we had no nextCursor and lastPage was not set to true.
	return objects, nil, nil
}

func ParseGroupsResponse(
	body []byte, pageSize int64, cursor string,
) (objects []map[string]any, nextCursor *string, err *framework.Error) {
	return parseResponse(body, pageSize, cursor, Group, "isLast")
}

func ParseGroupMembersResponse(
	body []byte, pageSize int64, cursor string,
) (objects []map[string]any, nextCursor *string, err *framework.Error) {
	return parseResponse(body, pageSize, cursor, GroupMember, "isLast")
}

func ParseWorkspacesResponse(
	body []byte, pageSize int64, cursor string,
) (objects []map[string]any, nextCursor *string, err *framework.Error) {
	return parseResponse(body, pageSize, cursor, Workspace, "isLastPage")
}

func ParseObjectsResponse(
	body []byte, pageSize int64, cursor string,
) (objects []map[string]any, nextCursor *string, err *framework.Error) {
	return parseResponse(body, pageSize, cursor, Object, "isLast")
}

// parseResponse parses Jira responses that have the format {"values": []Entity}.
// If the lastPageFieldName field exists, it is used to determine if the current page is the last page.
// If parsing fails, a framework.Error is returned.
func parseResponse(
	body []byte, pageSize int64, cursor string, entityExternalID string, lastPageFieldName string,
) (objects []map[string]any, nextCursor *string, err *framework.Error) {
	var data map[string]any

	unmarshalErr := json.Unmarshal(body, &data)
	if unmarshalErr != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal Jira %s response: %v.", entityExternalID, unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	_, found := data["values"]
	if !found {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Field missing in Jira %s response: values.", entityExternalID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	valuesAsList, ok := data["values"].([]any)
	if !ok {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf(
				"Entity field exists in Jira %s response but field value is not a list of objects: %T.",
				entityExternalID,
				data["values"],
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	objects, parserErr := parseObjects(entityExternalID, valuesAsList)
	if parserErr != nil {
		return nil, nil, parserErr
	}

	// Optimization:
	// If the response contains the lastPageFieldName field, and it's a bool, and it's true,
	// the current page is the last page and there is no need to compute the next cursor.
	isLast, found := data[lastPageFieldName]
	if found {
		isLastPage, ok := isLast.(bool)
		if !ok {
			return nil, nil, &framework.Error{
				Message: fmt.Sprintf(
					"Field %s exists in Jira %s response but field value is not a bool: %T.",
					lastPageFieldName,
					entityExternalID,
					isLast,
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		if isLastPage {
			return objects, nil, nil
		}
	}

	cursorInt, convErr := strconv.ParseInt(cursor, 10, 64)
	if convErr != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf(
				"Failed to parse cursor as an int (%s): %s.",
				cursor,
				convErr.Error(),
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	nextCursorInt := pagination.GetNextCursorFromPageSize(len(objects), pageSize, cursorInt)

	if nextCursorInt != nil {
		tmp := strconv.FormatInt(*nextCursorInt, 10)
		nextCursor = &tmp
	}

	return objects, nextCursor, nil
}

// parseObjects parses []any into []map[string]any. If any object in the slice is not a map[string]any,
// a framework.Error is returned.
func parseObjects(entityExternalID string, objects []any) ([]map[string]any, *framework.Error) {
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

		// Issue entity requires additional parsing of custom fields.
		// TODO: Remove this after properly supporting parsing complex multi-valued attributes
		// into lists based on JSON Paths given in the attribute external IDs.
		if entityExternalID == Issue || entityExternalID == EnhancedIssue {
			parsedIssueObject, err := parseIssueCustomFields(parsedObject)
			if err != nil {
				return nil, err
			}

			parsedObject = parsedIssueObject
		}

		parsedObjects = append(parsedObjects, parsedObject)
	}

	return parsedObjects, nil
}

// parseIssueCustomFields parses multi-valued complex custom fields in an Issue object.
// It takes in the Issue object and returns the same object with the custom fields parsed.
// TODO: Remove this function after implementing a proper solution based on using JSON Paths as external IDs.
func parseIssueCustomFields(object map[string]any) (map[string]any, *framework.Error) {
	if object == nil {
		return nil, nil
	}

	// Custom fields are nested under the "fields" field. If that does not exist, return early.
	objectFields, found := object["fields"]
	if !found {
		return object, nil
	}

	objectFieldsAsMap, ok := objectFields.(map[string]any)
	if !ok {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to parse fields field in Jira %s object as map[string]any.", Issue),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	customFields := []string{"customfield_10069", "customfield_11605"}

	for _, customField := range customFields {
		// These specific custom fields are []objects. For example,
		// "customfield_10069": [
		// {
		// 	"self": "https://test-instance.atlassian.net/rest/api/3/customFieldOption/10153",
		// 	"value": "A",
		// 	"id": "10153"
		// }]
		// Convert the custom field into "customfield_10069": ["A"] instead of an object.
		customFieldValue := objectFieldsAsMap[customField]

		if customFieldValue == nil {
			continue
		}

		customFieldObjs, ok := customFieldValue.([]any)
		if !ok {
			return nil, &framework.Error{
				Message: fmt.Sprintf("Failed to parse %s field in Jira %s object as []any.", customField, Issue),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}
		// This must be left as []any to let adapter framework parse into []string.
		// https://github.com/SGNL-ai/adapter-framework/blob/b1c34af0488a20c96c986b9e8d9a9bd44c5820f0/web/json_value.go#L141
		values := make([]any, 0, len(customFieldObjs))

		for _, obj := range customFieldObjs {
			objAsMap, ok := obj.(map[string]any)
			if !ok {
				return nil, &framework.Error{
					Message: fmt.Sprintf(
						"Failed to parse one of the objects in Jira %s field into map[string]any: %v.",
						customField, obj,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			// We only extract the "value" field.
			value, ok := objAsMap["value"].(string)
			if !ok {
				return nil, &framework.Error{
					Message: fmt.Sprintf("Failed to parse value field in Jira %s object as string: %v.", customField, objAsMap),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			values = append(values, value)
		}

		// Replace the object["fields"][customField] array of objects value with this array of strings.
		objectFieldsAsMap[customField] = values
		object["fields"] = objectFieldsAsMap
	}

	return object, nil
}

// ConstructURL constructs the Jira URL for the given request and entity.
// This URL is used to make a request to the Jira API.
// For URLs which require a resource ID, the resource ID must be non nil (i.e. a non nil cursor.CollectionID)
// otherwise an error is returned.
func ConstructURL(request *Request, entity Entity, cursor *pagination.CompositeCursor[string]) (string, error) {
	var sb strings.Builder

	var cursorStr string

	if cursor.Cursor != nil {
		cursorStr = *cursor.Cursor
	}

	pageSizeStr := strconv.FormatInt(request.PageSize, 10)

	switch request.EntityExternalID {
	// The URL format depends on entity.
	// The Workspace entity is retrieved from the Jira Service Management (JSM) API.
	// The Object entity is retrieved from the Jira Assets API.
	// Other entities are retrieved from the Jira Cloud API.
	case Workspace:
		// request.BaseURL + "/rest/servicedeskapi/assets/" + entity.endpoint + "?"
		// len("/rest/servicedeskapi/assets/") + len("?") == 29.
		sb.Grow(len(request.BaseURL) + len(entity.endpoint) + 29)
		sb.WriteString(request.BaseURL)
		sb.WriteString("/rest/servicedeskapi/assets/")
		sb.WriteString(entity.endpoint)
		sb.WriteRune('?')
	case Object:
		if cursor.CollectionID == nil {
			return "", fmt.Errorf("cursor.CollectionID must not be nil for Object entity")
		}

		if request.AssetBaseURL == nil {
			return "", fmt.Errorf("request.AssetBaseURL must not be nil for Object entity")
		}
		// assetBaseURL + "/workspace/ + workspaceID + "/v1/object/aql?includeAttributes=true?"
		// len("/workspace/") + len("/v1/object/aql?includeAttributes=true?") == 49.
		sb.Grow(len(*request.AssetBaseURL) + len(entity.endpoint) + len(*cursor.CollectionID) + 49)
		sb.WriteString(*request.AssetBaseURL)
		sb.WriteString("/workspace/")
		sb.WriteString(*cursor.CollectionID)
		sb.WriteString("/v1/object/aql?includeAttributes=true&")
	default:
		// request.BaseURL + "/rest/api/3/" + entity.endpoint
		// len("/rest/api/3/") + len("?") == 13.
		sb.Grow(len(request.BaseURL) + len(entity.endpoint) + 13)
		sb.WriteString(request.BaseURL)
		sb.WriteString("/rest/api/3/")
		sb.WriteString(entity.endpoint)
		sb.WriteRune('?')

		switch request.EntityExternalID {
		case GroupMember:
			if cursor.CollectionID == nil {
				return "", fmt.Errorf("cursor.CollectionID must not be nil for GroupMember entity")
			}

			sb.Grow(len(*cursor.CollectionID) + 9)
			sb.WriteString("groupId=")
			sb.WriteString(*cursor.CollectionID)
			sb.WriteRune('&')
		case Issue, EnhancedIssue:
			if request.IssuesJQLFilter != nil {
				escapedFilter := net_url.QueryEscape(*request.IssuesJQLFilter)
				sb.Grow(len(escapedFilter) + 5)
				sb.WriteString("jql=")
				sb.WriteString(escapedFilter)
				sb.WriteRune('&')
			}
		}
	}

	if request.EntityExternalID == EnhancedIssue {
		// len("nextPageToken=") + len("&maxResults=") + len("&fields=*navigable") == 44.
		sb.Grow(len(cursorStr) + len(pageSizeStr) + 44)

		// Don't specify a next page token if the value is 0. This functionality differs from all other endpoints
		// that set `startAt=0` for the first page.
		if cursorStr != "0" {
			sb.WriteString("nextPageToken=")
			sb.WriteString(cursorStr)
			sb.WriteRune('&')
		}

		sb.WriteString("maxResults=")
		sb.WriteString(pageSizeStr)

		// This was the default with the deprecated API endpoint. The new default is just to return IDs, so we're
		// manually specifying this to keep the same functionality.
		sb.WriteString("&fields=*navigable")

		return sb.String(), nil
	}

	// All endpoints share these query parameters, so they are added last.
	// len("startAt=") + len("&maxResults=") == 20.
	sb.Grow(len(cursorStr) + len(pageSizeStr) + 20)
	sb.WriteString("startAt=")
	sb.WriteString(cursorStr)
	sb.WriteString("&maxResults=")
	sb.WriteString(pageSizeStr)

	return sb.String(), nil
}

// basicAuth returns the basic auth header value for the given username and password base64 encoded.
// Cf. https://developer.atlassian.com/cloud/jira/platform/basic-auth-for-rest-apis/#supply-basic-auth-headers.
func basicAuth(username, password string) string {
	auth := username + ":" + password

	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
