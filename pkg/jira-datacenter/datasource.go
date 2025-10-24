// Copyright 2025 SGNL.ai, Inc.
package jiradatacenter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	net_url "net/url"
	"sort"
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
	// JiraFieldsPrefix is the prefix used for Jira field paths in JSONPath expressions.
	JiraFieldsPrefix = "fields"

	// Entity external IDs.
	UserExternalID        string = "User"
	IssueExternalID       string = "Issue"
	GroupExternalID       string = "Group"
	GroupMemberExternalID string = "GroupMember"

	// Field names.
	GroupIDFieldName    string = "groupId"
	GroupMemberIDFormat string = "%s-%s" // Format: groupID-userID
)

var EntityIDToParentCollectionID = map[string]string{
	UserExternalID:        GroupExternalID,
	GroupMemberExternalID: GroupExternalID,
}

var isLastFieldName = "isLast"

var (
	// ValidEntityExternalIDs is a map of valid external IDs of entities that can be queried.
	// The map value is the Entity struct.
	// Users doc:
	//   https://docs.atlassian.com/software/jira/docs/api/REST/9.17.0/#api/2/group-getUsersFromGroup
	// Issues doc:
	//   https://docs.atlassian.com/software/jira/docs/api/REST/9.17.0/#api/2/search-search
	// Groups doc:
	//   https://docs.atlassian.com/software/jira/docs/api/REST/9.17.0/#api/2/groups-findGroups
	// GroupMembers doc:
	//   https://docs.atlassian.com/software/jira/docs/api/REST/9.17.0/#api/2/group-getUsersFromGroup
	ValidEntityExternalIDs = map[string]Entity{
		UserExternalID: {
			externalID:             UserExternalID,
			uniqueIDAttrExternalID: "key",
			endpoint:               "group/member", // User entity is retrieved through group membership.
			objectArrayFieldName:   "values",
			lastPageFieldName:      &isLastFieldName,
		},
		IssueExternalID: {
			externalID:             IssueExternalID,
			uniqueIDAttrExternalID: "id",
			endpoint:               "search",
			objectArrayFieldName:   "issues",
		},
		GroupExternalID: {
			externalID:             GroupExternalID,
			uniqueIDAttrExternalID: "name",
			endpoint:               "groups/picker",
			objectArrayFieldName:   "groups",
		},
		GroupMemberExternalID: {
			// "id" is created by combining the groupId and accountId.
			// It is not returned by the Jira API.
			externalID:             GroupMemberExternalID,
			uniqueIDAttrExternalID: "id",
			endpoint:               "group/member",
			objectArrayFieldName:   "values",
			lastPageFieldName:      &isLastFieldName,
		},
	}
)

// Datasource implements the Jira Client interface to allow querying the Jira datasource.
type Datasource struct {
	Client *http.Client
}

// Entity contains entity specific information, such as the entity's unique ID attribute and the
// endpoint to query that entity.
type Entity struct {
	externalID string

	// uniqueIDAttrExternalID is the external ID of the entity's uniqueId attribute.
	uniqueIDAttrExternalID string

	// endpoint is the endpoint to query the entity, e.g. "users/search" for users.
	// It does not need to include "/rest/api/2/".
	endpoint string

	// objectArrayFieldName is the name of the field in the API response that contains
	// the array of entity objects (e.g., "values" for users, "issues" for issues).
	objectArrayFieldName string

	// lastPageFieldName is the name of the field in the API response that indicates
	// whether the current page is the last page (e.g., "isLast").
	// This field is optional and can be nil.
	lastPageFieldName *string
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

	request.Cursor = cursor

	if isEmptyLastPage {
		return &Response{
			StatusCode: http.StatusOK,
		}, nil
	}

	validationErr := pagination.ValidateCompositeCursor(
		cursor,
		request.EntityExternalID,
		// Send a bool indicating if the entity is a member of a collection.
		request.EntityExternalID == UserExternalID || request.EntityExternalID == GroupMemberExternalID,
	)
	if validationErr != nil {
		return nil, validationErr
	}

	var req *http.Request

	var err error

	url, err := entity.ConstructURL(request, cursor)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Unable to construct URL: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// Timeout API calls that take longer than the configured timeout.
	apiCtx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	req, err = http.NewRequestWithContext(apiCtx, http.MethodGet, url, nil)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Address in datasource config is an invalid URL: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	req.Header.Add("Authorization", request.AuthorizationHeader)

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
			fields.ResponseStatusCode(res.StatusCode),
			fields.ResponseRetryAfterHeader(response.RetryAfterHeader),
			fields.ResponseBody(body),
			fields.SGNLEventTypeError(),
		)

		return response, nil
	}

	objects, nextCursor, frameworkErr := entity.Parse(body, *request)
	if frameworkErr != nil {
		return nil, frameworkErr
	}

	response.NextCursor = &pagination.CompositeCursor[int64]{
		Cursor: nextCursor,
	}

	switch request.EntityExternalID {
	case GroupMemberExternalID:
		// We must create a unique ID for GroupMembers that is a
		// combination of the groupID and userID, since a user may belong
		// to multiple groups.
		userUniqueIDAttrExternalID := ValidEntityExternalIDs[UserExternalID].uniqueIDAttrExternalID
		groupID := *cursor.CollectionID
		groupMemberUniqueIDAttrExternalID := ValidEntityExternalIDs[GroupMemberExternalID].uniqueIDAttrExternalID

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

			object[GroupIDFieldName] = groupID
			object[groupMemberUniqueIDAttrExternalID] = fmt.Sprintf(GroupMemberIDFormat, groupID, userID)
		}

		response.NextCursor.CollectionID = cursor.CollectionID
		response.NextCursor.CollectionCursor = cursor.CollectionCursor
	case UserExternalID:
		// For User entity, we need to preserve the group context
		// since we're fetching users through group membership
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
) (*pagination.CompositeCursor[int64], bool, *framework.Error) {
	if request.Cursor != nil && request.Cursor.Cursor != nil {
		return request.Cursor, false, nil
	}

	var zero int64

	cursor := &pagination.CompositeCursor[int64]{
		Cursor: &zero,
	}

	parentCollectionEntityID, hasParentCollection := EntityIDToParentCollectionID[request.EntityExternalID]
	if !hasParentCollection {
		return cursor, false, nil
	}

	collectionCursor := &zero
	if request.Cursor != nil && request.Cursor.CollectionCursor != nil {
		collectionCursor = request.Cursor.CollectionCursor
	}

	// We have no more members to query for the last requested collection,
	// or this is a request for the first page.
	// Get the ID of the next collection.
	nextCollectionReq := &Request{
		BaseURL:               request.BaseURL,
		AuthorizationHeader:   request.AuthorizationHeader,
		PageSize:              1,
		Cursor:                &pagination.CompositeCursor[int64]{Cursor: collectionCursor},
		Groups:                request.Groups,
		EntityExternalID:      parentCollectionEntityID,
		RequestTimeoutSeconds: request.RequestTimeoutSeconds,
		APIVersion:            request.APIVersion,
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

	cursor.CollectionID = &firstCollectionID

	if nextCollectionRes.NextCursor != nil {
		cursor.CollectionCursor = nextCollectionRes.NextCursor.Cursor
	}

	return cursor, false, nil
}

func (e Entity) Parse(
	body []byte, request Request,
) (objects []map[string]any, nextCursor *int64, err *framework.Error) {
	switch e.externalID {
	case UserExternalID, GroupMemberExternalID, IssueExternalID:
		return e.parseResponse(body, request.PageSize, *request.Cursor.Cursor)
	case GroupExternalID:
		return e.parseGroupsResponse(body, request)
	default:
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Unexpected entity type: %s", e.externalID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}
}

// parseResponse parses Jira responses that have the format {"values": []Entity}.
// If the lastPageFieldName field exists, it is used to determine if the current page is the last page.
// If parsing fails, a framework.Error is returned.
func (e Entity) parseResponse(
	body []byte, pageSize int64, cursor int64,
) (objects []map[string]any, nextCursor *int64, err *framework.Error) {
	var data map[string]any

	unmarshalErr := json.Unmarshal(body, &data)
	if unmarshalErr != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal Jira %s response: %v.", e.externalID, unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	_, found := data[e.objectArrayFieldName]
	if !found {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Field missing in Jira %s response: %s.", e.objectArrayFieldName, e.externalID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	valuesAsList, ok := data[e.objectArrayFieldName].([]any)
	if !ok {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf(
				"Entity field exists in Jira %s response but field value is not a list of objects: %T.",
				e.externalID,
				data[e.objectArrayFieldName],
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	objects, parserErr := e.parseObjects(valuesAsList)
	if parserErr != nil {
		return nil, nil, parserErr
	}

	// Optimization:
	// If the response contains the lastPageFieldName field, and it's a bool, and it's true,
	// the current page is the last page and there is no need to compute the next cursor.
	if e.lastPageFieldName != nil {
		if isLast, ok := data[*e.lastPageFieldName].(bool); isLast && ok {
			return objects, nil, nil
		}
	}

	nextCursor = pagination.GetNextCursorFromPageSize(len(objects), pageSize, cursor)

	return objects, nextCursor, nil
}

// parseObjects parses []any into []map[string]any. If any object in the slice is not a map[string]any,
// a framework.Error is returned.
func (e Entity) parseObjects(objects []any) ([]map[string]any, *framework.Error) {
	parsedObjects := make([]map[string]any, 0, len(objects))

	for _, object := range objects {
		parsedObject, ok := object.(map[string]any)
		if !ok {
			return nil, &framework.Error{
				Message: fmt.Sprintf(
					"An object in Entity: %s could not be parsed. Expected: map[string]any. Got: %T.",
					e.externalID,
					object,
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		parsedObjects = append(parsedObjects, parsedObject)
	}

	return parsedObjects, nil
}

func (e Entity) parseGroupsResponse(
	body []byte, request Request,
) (objects []map[string]any, nextCursor *int64, err *framework.Error) {
	objects, _, frameworkErr := e.parseResponse(
		body,
		request.PageSize,
		*request.Cursor.Cursor,
	)
	if frameworkErr != nil {
		return nil, nil, frameworkErr
	}

	objects, nextCursor = filterAndPaginateGroups(objects, request.Groups, request.PageSize, *request.Cursor.Cursor)

	return objects, nextCursor, nil
}

// ConstructURL constructs the Jira URL for the given request and entity.
// This URL is used to make a request to the Jira API.
// For URLs which require a resource ID, the resource ID must be non nil (i.e. a non nil cursor.CollectionID)
// otherwise an error is returned.
func (e Entity) ConstructURL(request *Request, cursor *pagination.CompositeCursor[int64]) (string, error) {
	var sb strings.Builder

	cursorStr := strconv.FormatInt(*cursor.Cursor, 10)
	pageSizeStr := strconv.FormatInt(request.PageSize, 10)

	apiVersion := request.APIVersion
	if apiVersion == "" {
		apiVersion = "latest"
	}

	// Preallocate for the base URL: BaseURL + "/rest/api/" + apiVersion + "/"
	baseLen := len(request.BaseURL) + len("/rest/api/") + len(apiVersion) + 1
	sb.Grow(baseLen)

	// Build the base URL path for all endpoints
	sb.WriteString(request.BaseURL)
	sb.WriteString("/rest/api/")
	sb.WriteString(apiVersion)
	sb.WriteString("/")

	switch request.EntityExternalID {
	case GroupExternalID:
		// Groups/picker endpoint doesn't support standard pagination
		sb.Grow(len(e.endpoint))
		sb.WriteString(e.endpoint)

		return sb.String(), nil
	case UserExternalID, GroupMemberExternalID:
		if cursor.CollectionID == nil {
			return "", fmt.Errorf("cursor.CollectionID must not be nil for User entity or GroupMember entity")
		}
		// Preallocate for "group/member", "?groupname=",
		// the escaped CollectionID and an extra character for '&'
		escapedCollection := net_url.QueryEscape(*cursor.CollectionID)
		sectionLen := len(e.endpoint) + len("?groupname=") + len(escapedCollection) + 1
		sb.Grow(sectionLen)

		sb.WriteString(e.endpoint)
		sb.WriteString("?groupname=")
		sb.WriteString(escapedCollection)
		// Only include inactive users if specifically requested
		if request.IncludeInactiveUsers != nil && *request.IncludeInactiveUsers {
			sb.WriteString("&includeInactiveUsers=true")
		}

		sb.WriteString("&")

	case IssueExternalID:
		// Preallocate for endpoint and "?"
		sb.Grow(len(e.endpoint) + 1)
		sb.WriteString(e.endpoint)
		sb.WriteString("?")

		if request.IssuesJQLFilter != nil {
			escapedFilter := net_url.QueryEscape(*request.IssuesJQLFilter)
			// Preallocate for "jql=" + escaped filter + '&'
			sb.Grow(len("jql=") + len(escapedFilter) + 1)
			sb.WriteString("jql=")
			sb.WriteString(escapedFilter)
			sb.WriteRune('&')
		}

		// Build fields parameter from attributes including child entity attributes
		fieldsParam := BuildJiraFieldsParam(request.Entity)
		fieldsStr := "fields=" + fieldsParam + "&"
		sb.Grow(len(fieldsStr))
		sb.WriteString(fieldsStr)

	default:
		// Preallocate for endpoint and "?"
		sb.Grow(len(e.endpoint) + 1)
		sb.WriteString(e.endpoint)
		sb.WriteString("?")
	}

	// Preallocate for the pagination parameters:
	// "startAt=" + cursorStr + "&maxResults=" + pageSizeStr
	paginationLen := len("startAt=") + len(cursorStr) + len("&maxResults=") + len(pageSizeStr)
	sb.Grow(paginationLen)

	sb.WriteString("startAt=")
	sb.WriteString(cursorStr)
	sb.WriteString("&maxResults=")
	sb.WriteString(pageSizeStr)

	return sb.String(), nil
}

// filterAndPaginateGroups filters groups based on the configured group list and
// applies manual pagination to the filtered results.
func filterAndPaginateGroups(
	allGroups []map[string]any,
	configGroups []string,
	pageSize int64,
	cursor int64,
) (objects []map[string]any, nextCursor *int64) {
	// First, filter the groups if config.Groups is non-empty
	var filteredGroups []map[string]any

	if len(configGroups) > 0 {
		groupFilter := make(map[string]bool)
		for _, name := range configGroups {
			groupFilter[name] = true
		}

		filteredGroups = make([]map[string]any, 0)

		for _, obj := range allGroups {
			if groupName, ok := obj["name"].(string); ok {
				if groupFilter[groupName] {
					filteredGroups = append(filteredGroups, obj)
				}
			}
		}
	} else {
		// If no filter is specified, use all groups
		filteredGroups = allGroups
	}

	// Groups are already sorted by the Jira API:
	// https://docs.atlassian.com/software/jira/docs/api/REST/9.17.0/#api/2/groups-findGroups

	// Apply pagination to the filtered groups
	startIndex := int(cursor)
	endIndex := startIndex + int(pageSize)

	if startIndex >= len(filteredGroups) {
		// Beyond the end of the list, return empty result
		return []map[string]any{}, nil
	}

	if endIndex > len(filteredGroups) {
		endIndex = len(filteredGroups)
	}

	objects = filteredGroups[startIndex:endIndex]

	// Calculate the next cursor
	if endIndex < len(filteredGroups) {
		nextCursorVal := int64(endIndex)
		nextCursor = &nextCursorVal
	} else {
		nextCursor = nil // No more pages
	}

	return objects, nextCursor
}

// removeArrayIndices removes array indices (like [0]) from field names.
// Examples:
//   - customfield_10209[0] → customfield_10209
//   - assignee[0] → assignee
//   - summary → summary (unchanged)
func removeArrayIndices(fieldName string) string {
	if idx := strings.Index(fieldName, "["); idx != -1 {
		return fieldName[:idx]
	}

	return fieldName
}

// extractFieldFromJSONPath extracts the Jira field name from a JSON path.
// Examples:
//   - $.fields.summary → summary
//   - $.fields.issuetype.id → issuetype
//   - $.fields.assignee.key → assignee
//   - $.fields.customfield_10209[0].value → customfield_10209
//   - $.id → id
//   - id → id (handles non-JSON path field names)
func extractFieldFromJSONPath(jsonPath string) string {
	// Handle non-JSON path field names (like "id", "key", "self")
	if !strings.HasPrefix(jsonPath, "$.") {
		return removeArrayIndices(jsonPath)
	}

	// Remove the "$." prefix
	path := strings.TrimPrefix(jsonPath, "$.")

	// Split by dots to get path segments
	segments := strings.Split(path, ".")

	// For paths like "$.id", return "id"
	if len(segments) == 1 {
		return removeArrayIndices(segments[0])
	}

	// For paths like "$.fields.summary", return "summary"
	// For paths like "$.fields.issuetype.id", return "issuetype"
	// For paths like "$.fields.customfield_10209[0].value", return "customfield_10209"
	if len(segments) >= 2 && segments[0] == JiraFieldsPrefix {
		return removeArrayIndices(segments[1])
	}

	// For other cases, return the first segment
	return removeArrayIndices(segments[0])
}

// BuildJiraFieldsParam constructs the 'fields' query parameter for Jira API requests.
// It extracts Jira field names from entity attribute JSON paths and deduplicates them.
// Returns "*navigable" if no attributes are provided (Jira's default for search endpoints).
// The returned string is URL-encoded and ready for use in API requests.
// It respects the Jira conventions for field selection:
// - *all - include all fields
// - *navigable - include just navigable fields (default for search)
// - field1,field2 - include specific fields
// - -field - exclude a field.
func BuildJiraFieldsParam(entity *framework.EntityConfig) string {
	encodedAttrs := ExtractEntityFieldNames("", entity)
	if len(encodedAttrs) == 0 {
		return "*navigable"
	}
	// Convert map to slice
	fields := make([]string, 0, len(encodedAttrs))
	for field := range encodedAttrs {
		fields = append(fields, field)
	}

	// Sort fields for consistent output
	sort.Strings(fields)

	// Join fields with comma and then URL-encode the entire string
	return net_url.QueryEscape(strings.Join(fields, ","))
}

// ExtractEntityFieldNames recursively extracts Jira field names from an entity configuration
// and all its child entities. It processes the entity's attributes and combines them with
// field names from nested child entities using dot notation for prefixes.
// Returns a set of unique field names suitable for Jira API field selection.
func ExtractEntityFieldNames(prefix string, entity *framework.EntityConfig) map[string]struct{} {
	if entity == nil {
		return map[string]struct{}{}
	}

	encodedAttrs := ExtractFieldNamesFromAttributes(prefix, entity.Attributes)

	// Include field attribute names from child entities.
	for _, childEntity := range entity.ChildEntities {
		if prefix != "" {
			prefix += "."
		}

		for attr := range ExtractEntityFieldNames(prefix+childEntity.ExternalId, childEntity) {
			encodedAttrs[attr] = struct{}{}
		}
	}

	return encodedAttrs
}

// ExtractFieldNamesFromAttributes extracts Jira field names from a list of attribute configurations.
// It processes each attribute's ExternalId (which may be a JSON path) and converts it to a
// Jira field name using extractFieldFromJSONPath. Prefixes are applied for nested attributes.
// Returns a deduplicated set of field names.
func ExtractFieldNamesFromAttributes(prefix string, attributes []*framework.AttributeConfig) map[string]struct{} {
	// Use a map to deduplicate field names
	fieldSet := make(map[string]struct{})
	if len(attributes) == 0 {
		return fieldSet
	}

	for _, attribute := range attributes {
		if attribute.ExternalId != "" {
			attr := attribute.ExternalId
			if prefix != "" {
				attr = prefix + "." + attr
			}

			fieldName := extractFieldFromJSONPath(attr)

			if fieldName != "" {
				fieldSet[fieldName] = struct{}{}
			}
		}
	}

	return fieldSet
}
