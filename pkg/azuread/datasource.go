// Copyright 2025 SGNL.ai, Inc.
package azuread

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Datasource directly implements a Client interface to allow querying an external datasource.
type Datasource struct {
	Client *http.Client
}

// Graph API response format example:
// User Response: https://learn.microsoft.com/en-us/graph/api/user-list?view=graph-rest-1.0&tabs=http#examples
// Paging: https://learn.microsoft.com/en-us/graph/paging?tabs=http
type DatasourceResponse struct {
	Values   []map[string]any `json:"value"`
	NextLink *string          `json:"@odata.nextLink"`
}

type EntityInfo struct {
	memberOf *string
}

const (
	User           string = "User"
	Group          string = "Group"
	GroupMember    string = "GroupMember"
	Application    string = "Application"
	Device         string = "Device"
	Role           string = "Role"
	RoleMember     string = "RoleMember"
	RoleAssignment string = "RoleAssignment"

	// Use a combination of $top and $skip to paginate the response for these two PIM entities.
	RoleAssignmentScheduleRequest  string = "RoleAssignmentScheduleRequest"
	GroupAssignmentScheduleRequest string = "GroupAssignmentScheduleRequest"
)

var (
	// ValidEntityExternalIDs is a set of valid external IDs of entities that can be queried.
	ValidEntityExternalIDs = map[string]EntityInfo{
		User:  {},
		Group: {},
		GroupMember: {memberOf: func() *string {
			s := Group // Entity containing the group member data

			return &s
		}()},
		Application: {},
		Device:      {},
		Role:        {},
		RoleMember: {memberOf: func() *string {
			s := User // Entity containing the role member data

			return &s
		}()},
		RoleAssignment:                 {},
		RoleAssignmentScheduleRequest:  {},
		GroupAssignmentScheduleRequest: {},
	}

	// Advanced query operators that require the `ConsistencyLevel: eventual` header.
	// These need to be matched as whole words/operators, not substrings.
	advancedQueryOperators = map[string]struct{}{
		"endswith":   {},
		"contains":   {},
		"startswith": {},
	}

	// Regex patterns for advanced query operators that need word boundary matching.
	neOperatorRegex  = regexp.MustCompile(`\bne\b`)  // 'ne' as whole word
	notOperatorRegex = regexp.MustCompile(`\bnot\b`) // 'not' as whole word
)

// NewClient returns a Client to query the datasource.
func NewClient(client *http.Client) Client {
	return &Datasource{
		Client: client,
	}
}

// deepCopyCursor creates a deep copy of a CompositeCursor.
func deepCopyCursor(cursor *pagination.CompositeCursor[string]) *pagination.CompositeCursor[string] {
	if cursor == nil {
		return nil
	}

	result := &pagination.CompositeCursor[string]{}

	if cursor.Cursor != nil {
		cursorVal := *cursor.Cursor
		result.Cursor = &cursorVal
	}

	if cursor.CollectionID != nil {
		collectionIDVal := *cursor.CollectionID
		result.CollectionID = &collectionIDVal
	}

	if cursor.CollectionCursor != nil {
		collectionCursorVal := *cursor.CollectionCursor
		result.CollectionCursor = &collectionCursorVal
	}

	return result
}

// GetPage fetches a page of data, accumulating members from multiple parent groups if needed to satisfy page size.
func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	// Check if this is a member entity that needs accumulation
	parentEntityExternalID := getParentEntityExternalID(request)

	// For non-member entities, just use the base implementation
	if parentEntityExternalID == nil {
		return d.getPageBase(ctx, request)
	}

	// For member entities, accumulate members from multiple parent groups to satisfy page size
	accumulatedObjects := make([]map[string]any, 0)
	currentRequest := *request

	var nextCursor *pagination.CompositeCursor[string]

	for int64(len(accumulatedObjects)) < request.PageSize {
		response, err := d.getPageBase(ctx, &currentRequest)
		if err != nil {
			return nil, err
		}

		if response.StatusCode != http.StatusOK {
			return response, nil
		}

		if len(response.Objects) > 0 {
			if int64(len(response.Objects)+len(accumulatedObjects)) > request.PageSize {
				// Received more data than asked for.
				break
			}

			accumulatedObjects = append(accumulatedObjects, response.Objects...)
		}

		// Check if we have more data to fetch
		if response.NextCursor == nil {
			nextCursor = nil

			break // No more data available
		}

		// Save the cursor
		nextCursor = deepCopyCursor(response.NextCursor)

		// Update the cursor
		currentRequest.Cursor = response.NextCursor
	}

	// Build final response with accumulated objects
	return &Response{
		StatusCode: http.StatusOK,
		Objects:    accumulatedObjects,
		NextCursor: nextCursor,
	}, nil
}

func (d *Datasource) getPageBase(ctx context.Context, request *Request) (*Response, *framework.Error) {
	logger := zaplogger.FromContext(ctx).With(
		fields.RequestEntityExternalID(request.EntityExternalID),
		fields.RequestPageSize(request.PageSize),
	)

	logger.Info("Starting datasource request")

	// [MemberEntities] For member entities, we need to set the `CollectionID` and `CollectionCursor`.
	parentEntityExternalID := getParentEntityExternalID(request)

	if parentEntityExternalID != nil {
		memberReq := &Request{
			Token:                 request.Token,
			APIVersion:            request.APIVersion,
			BaseURL:               request.BaseURL,
			EntityExternalID:      *parentEntityExternalID,
			PageSize:              1,
			Filter:                request.ParentFilter,
			RequestTimeoutSeconds: request.RequestTimeoutSeconds,
		}

		// If the CollectionCursor is set, use that as the Cursor for the next call to `GetPage`.
		if request.Cursor != nil && request.Cursor.CollectionCursor != nil {
			memberReq.Cursor = &pagination.CompositeCursor[string]{
				Cursor: request.Cursor.CollectionCursor,
			}
		}

		if request.Cursor == nil {
			request.Cursor = &pagination.CompositeCursor[string]{}
		}

		isEmptyLastPage, cursorErr := pagination.UpdateNextCursorFromCollectionAPI(
			ctx,
			request.Cursor,
			func(ctx context.Context, _ *Request) (
				int, string, []map[string]any, *pagination.CompositeCursor[string], *framework.Error,
			) {
				resp, err := d.getPageBase(ctx, memberReq)
				if err != nil {
					return 0, "", nil, nil, err
				}

				return resp.StatusCode, resp.RetryAfterHeader, resp.Objects, resp.NextCursor, nil
			},
			memberReq,
			"id",
		)
		if cursorErr != nil {
			return nil, cursorErr
		}

		// If `ConstructEndpoint` hits a page with no `CollectionID` and no
		// `CollectionCursor` we should complete the sync at this point.
		if isEmptyLastPage {
			return &Response{
				StatusCode: http.StatusOK,
			}, nil
		}
	}

	// [!MemberEntities] This verifies that `CollectionID` and `CollectionCursor` are not set.
	// [MemberEntities] This verifies that `CollectionID` is set.
	validationErr := pagination.ValidateCompositeCursor(
		request.Cursor,
		request.EntityExternalID,
		// Send a bool indicating if the entity is a member of a collection.
		parentEntityExternalID != nil,
	)
	if validationErr != nil {
		return nil, validationErr
	}

	if isPIMEntity(request.EntityExternalID) {
		if request.Cursor != nil {
			offset, err := strconv.ParseInt(*request.Cursor.Cursor, 10, 64)
			if err != nil {
				return nil, &framework.Error{
					Message: fmt.Sprintf("Expected a numeric cursor for PIM entities: %v.", *request.Cursor.CollectionCursor),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			request.Skip = offset
		}
	}

	endpoint, endpointErr := ConstructEndpoint(request)
	if endpointErr != nil {
		return nil, endpointErr
	}

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

	req.Header.Add("Authorization", request.Token)

	// Enhanced advanced query detection - check if we need ConsistencyLevel: eventual
	if IsAdvancedQuery(request, endpoint) {
		req.Header.Add("ConsistencyLevel", "eventual")
	}

	logger.Info("Sending HTTP request to datasource", fields.URL(endpoint))

	res, err := d.Client.Do(req)
	if err != nil {
		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute Azure AD request: %v.", err),
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
		logger.Error("Datasource request failed",
			fields.ResponseStatusCode(response.StatusCode),
			fields.ResponseRetryAfterHeader(response.RetryAfterHeader),
			fields.ResponseBody(res.Body),
		)

		return response, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to read Azure AD response: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	objects, nextLink, frameworkErr := ParseResponse(body)
	if frameworkErr != nil {
		return nil, frameworkErr
	}

	// [Roles] No pagination support from the server side for this entity.
	if request.EntityExternalID == Role {
		objects, nextLink, frameworkErr = pagination.PaginateObjects(objects, request.PageSize, request.Cursor)
		if frameworkErr != nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf("Failed to paginate Directory Roles response - %v.", frameworkErr.Message),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}
	}

	// [RoleAssignmentScheduleRequest, GroupAssignmentScheduleRequest] Use $top and $skip to paginate the response.
	if isPIMEntity(request.EntityExternalID) {
		// Last page
		if int64(len(objects)) < request.PageSize {
			nextLink = nil
		} else {
			nextOffsetStr := strconv.FormatInt(request.Skip+request.PageSize, 10)
			nextLink = &nextOffsetStr
		}
	}

	if nextLink != nil {
		response.NextCursor = &pagination.CompositeCursor[string]{
			Cursor: nextLink,
		}
	}

	// [MemberEntities] Set `id`, `memberId` and `memberType`.
	if parentEntityExternalID != nil {
		for idx, member := range objects {
			memberID, ok := member[uniqueIDAttribute].(string)
			if !ok {
				return nil, &framework.Error{
					Message: fmt.Sprintf(
						"Failed to parse %s field in Azure AD Member response as string.",
						uniqueIDAttribute,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			switch request.EntityExternalID {
			case GroupMember:
				objects[idx]["id"] = fmt.Sprintf("%s-%s", memberID, *request.Cursor.CollectionID)
				objects[idx]["memberId"] = memberID
				objects[idx]["groupId"] = *request.Cursor.CollectionID
			case RoleMember:
				// To show relation between roles and users (user is a member of a role)
				objects[idx]["id"] = fmt.Sprintf("%s-%s", memberID, *request.Cursor.CollectionID)
				objects[idx]["memberId"] = *request.Cursor.CollectionID // userId
				objects[idx]["roleId"] = memberID                       // roleId
			}
		}

		if response.NextCursor != nil && response.NextCursor.Cursor != nil {
			request.Cursor.Cursor = response.NextCursor.Cursor
		} else {
			request.Cursor.Cursor = nil
		}

		// If we have a next cursor for either the base collection (Groups, Roles, etc.) or members (Group/Role/etc. Members),
		// encode the cursor for the next page. Otherwise, don't set a cursor as this sync is complete.
		if request.Cursor.Cursor != nil || request.Cursor.CollectionCursor != nil {
			response.NextCursor = request.Cursor
		}
	}

	response.Objects = objects

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(response.StatusCode),
		fields.ResponseObjectCount(len(response.Objects)),
		fields.ResponseNextCursor(response.NextCursor),
	)

	return response, nil
}

func getParentEntityExternalID(request *Request) *string {
	parentEntityExternalID := ValidEntityExternalIDs[request.EntityExternalID].memberOf

	if request.UseAdvancedFilters {
		switch request.EntityExternalID {
		case User:
			parentEntityExternalID = func() *string {
				s := Group

				return &s
			}()
		case Group:
			if request.AdvancedFilterMemberExternalID != nil {
				parentEntityExternalID = func() *string {
					s := Group

					return &s
				}()
			}
		}
	}

	return parentEntityExternalID
}

func ParseResponse(body []byte) (objects []map[string]any, nextLink *string, err *framework.Error) {
	var data *DatasourceResponse

	if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil || data == nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return data.Values, data.NextLink, nil
}

// nolint: lll
// IsAdvancedQuery determines if the request requires advanced query capabilities
// based on Microsoft Graph advanced query documentation.
// If any of the conditions are met, it returns true and a header "ConsistencyLevel: eventual"
// should be added to the request.
// See: https://learn.microsoft.com/en-us/graph/aad-advanced-queries?tabs=http#microsoft-entra-id-directory-objects-that-support-advanced-query-capabilities
// for more information.
// The ODATA logical operator NOT also needs the ConsistencyLevel header.
// See: https://learn.microsoft.com/en-us/graph/api/group-list?view=graph-rest-1.0&tabs=http
// The function checks for the following conditions:
// 1. If the request already has UseAdvancedFilters set to true, it returns true.
// 2. If the endpoint contains $count, $search, or $orderby query parameters.
// 3. If the endpoint contains $filter with advanced operators like endsWith, contains, startsWith.
// 4. If the endpoint contains $filter with 'ne' or 'not' operators as whole words.
// 5. If the endpoint contains $filter with 'NOT' operator as a whole word.
func IsAdvancedQuery(request *Request, endpoint string) bool {
	// If UseAdvancedFilters is already set, respect that.
	if request.UseAdvancedFilters {
		return true
	}

	// Check if endpoint contains $count, $search, or $orderby (all require advanced query).
	// Also check for $count as URL segment (e.g., ~/groups/$count).
	if strings.Contains(endpoint, "$count=true") ||
		strings.Contains(endpoint, "$search=") ||
		strings.Contains(endpoint, "$orderby=") ||
		strings.Contains(endpoint, "/$count") {
		return true
	}

	// Check if endpoint contains $filter with advanced operators.
	if strings.Contains(endpoint, "$filter=") {
		// URL decode the endpoint for proper pattern matching
		decodedEndpoint := endpoint
		if decoded, err := url.QueryUnescape(endpoint); err == nil {
			decodedEndpoint = decoded
		}

		endpointLower := strings.ToLower(decodedEndpoint)

		// Check for advanced function operators (can be substring matches).
		for operator := range advancedQueryOperators {
			if strings.Contains(endpointLower, operator) {
				return true
			}
		}

		// Check for 'ne' and 'not' operators using word boundary regex on decoded endpoint.
		decodedEndpointLower := strings.ToLower(decodedEndpoint)
		if neOperatorRegex.MatchString(decodedEndpointLower) || notOperatorRegex.MatchString(decodedEndpointLower) {
			return true
		}
	}

	return false
}
