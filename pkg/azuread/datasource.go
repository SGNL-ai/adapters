// Copyright 2025 SGNL.ai, Inc.
package azuread

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
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
)

// NewClient returns a Client to query the datasource.
func NewClient(client *http.Client) Client {
	return &Datasource{
		Client: client,
	}
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	// [MemberEntities] For member entities, we need to set the `CollectionID` and `CollectionCursor`.
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
				resp, err := d.GetPage(ctx, memberReq)
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

	if request.UseAdvancedFilters {
		req.Header.Add("ConsistencyLevel", "eventual")
	}

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
		// TEMP: Log the response body for debugging purposes.
		body, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			slog.Error(
				"Failed to read error response body",
				"error", readErr,
			)
		} else {
			slog.Error(
				"Azure AD API error",
				slog.Int("status", res.StatusCode),
				slog.String("response", string(body)),
			)
		}
		// END TEMP.

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

	return response, nil
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
