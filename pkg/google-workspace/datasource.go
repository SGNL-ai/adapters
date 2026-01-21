// Copyright 2026 SGNL.ai, Inc.

package googleworkspace

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type DatasourceResponse struct {
	DatasourceResponseTemplate
	Users   []map[string]interface{} `json:"users"`
	Groups  []map[string]interface{} `json:"groups"`
	Members []map[string]interface{} `json:"members"`
}

// Google Workspace API response template.
type DatasourceResponseTemplate struct {
	Kind          string        `json:"kind"`
	Etag          string        `json:"etag"`
	NextPageToken *string       `json:"nextPageToken"`
	Error         *ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  []ErrorInfo `json:"errors"`
	Status  string      `json:"status"`
	Details []ErrorInfo `json:"details"`
}

type ErrorInfo struct {
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

// Entity contains entity specific information, such as the entity's unique ID attribute and the
// endpoint path to query that entity.
type Entity struct {
	// Path is the endpoint to query the entity.
	Path string

	// OrderByAttribute is the attribute to order the results by.
	// This is only supported for certain endpoints.
	// User: https://developers.google.com/admin-sdk/directory/reference/rest/v1/users/list#orderby
	// Group: https://developers.google.com/admin-sdk/directory/reference/rest/v1/groups/list#orderby
	OrderByAttribute string

	// MaxPageSize is the maximum page size allowed for each entity's endpoint.
	MaxPageSize int64

	// The RequiredAttributes array specifies attributes necessary to create the uniqueID in the response.
	// In Member, the id field is the member object's id. To create the uniqueId, both id and groupId are needed.
	// groupId is required for establishing relationships between Groups and Members (members can be groups or users).
	RequiredAttributes []string

	// UniqueIDAttribute is the attribute that uniquely identifies an entity.
	UniqueIDAttribute string
}

const (
	User   = "User"
	Group  = "Group"
	Member = "Member"
)

var (
	// ValidEntityExternalIDs is a set of valid external IDs of entities that can be queried.
	// The map value is the Entity struct which contains the unique ID attribute.
	ValidEntityExternalIDs = map[string]Entity{
		User: {
			// Example URI: /admin/directory/{{APIVersion}}/users
			Path:              "/admin/directory/%s/users",
			OrderByAttribute:  "EMAIL",
			MaxPageSize:       500,
			UniqueIDAttribute: "id",
		},
		Group: {
			// Example URI: /admin/directory/{{APIVersion}}/groups
			Path:              "/admin/directory/%s/groups",
			OrderByAttribute:  "EMAIL",
			MaxPageSize:       1000,
			UniqueIDAttribute: "id",
		},
		// In Google Workspace, a member of a group can be either a user or another group.
		Member: {
			// Example URI: /admin/directory/{{APIVersion}}/groups/{{groupId}}/members
			Path:               "/admin/directory/%s/groups/%s/members",
			MaxPageSize:        1000,
			UniqueIDAttribute:  "uniqueId",
			RequiredAttributes: []string{"id", "groupId"},
		},
	}
)

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

	if request.EntityExternalID == Member {
		groupReq := &Request{
			BaseURL:               request.BaseURL,
			APIVersion:            request.APIVersion,
			Domain:                request.Domain,
			Customer:              request.Customer,
			Filters:               request.Filters,
			Token:                 request.Token,
			PageSize:              1,
			EntityExternalID:      Group,
			RequestTimeoutSeconds: request.RequestTimeoutSeconds,
			Ordered:               false,
		}

		// If the CollectionCursor is set, use that as the Cursor
		// for the next call to `GetPage`
		if request.Cursor != nil && request.Cursor.CollectionCursor != nil {
			groupReq.Cursor = &pagination.CompositeCursor[string]{
				Cursor: request.Cursor.CollectionCursor,
			}
		}

		if request.Cursor == nil {
			request.Cursor = &pagination.CompositeCursor[string]{}
		}

		isLastPageEmpty, cursorErr := pagination.UpdateNextCursorFromCollectionAPI(
			ctx,
			request.Cursor,
			func(ctx context.Context, request *Request) (
				int, string, []map[string]any, *pagination.CompositeCursor[string], *framework.Error,
			) {
				resp, err := d.GetPage(ctx, request)
				if err != nil {
					return 0, "", nil, nil, err
				}

				return resp.StatusCode, resp.RetryAfterHeader, resp.Objects, resp.NextCursor, nil
			},
			groupReq,
			ValidEntityExternalIDs[Group].UniqueIDAttribute,
		)

		if cursorErr != nil {
			return nil, cursorErr
		}

		// If we run out of collection IDs, the sync is complete.
		if isLastPageEmpty {
			return &Response{
				StatusCode: http.StatusOK,
			}, nil
		}
	}

	if validationErr := pagination.ValidateCompositeCursor(
		request.Cursor,
		request.EntityExternalID,
		request.EntityExternalID == Member,
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
			Message: fmt.Sprintf("Failed to execute Google Workspace request: %v.", err),
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
			Message: fmt.Sprintf("Failed to read Google Workspace response body: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	objects, nextCursor, frameworkErr := ParseResponse(body, request)
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

func ParseResponse(body []byte, request *Request) (
	objects []map[string]any,
	nextCursor *pagination.CompositeCursor[string],
	err *framework.Error,
) {
	var response *DatasourceResponse

	if unmarshalErr := json.Unmarshal(body, &response); unmarshalErr != nil || response == nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if response.Error != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Google Workspace API returned an error: %v.", response.Error.Message),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Certain endpoints like the Group endpoint may return an empty response resulting in
	// objects to be set to nil. This is the expected behavior and is handled correctly downstream
	// in the framework provided object parsing functions.
	switch request.EntityExternalID {
	case User:
		objects = response.Users
	case Group:
		objects = response.Groups
	case Member:
		objects = response.Members

		// Post-processing for the Member entity.
		for idx, member := range objects {
			memberID, ok := member["id"].(string)
			if !ok {
				return nil, nil, &framework.Error{
					Message: fmt.Sprintf("Failed to parse 'id' field in Member response as string, actual value: %v.", member["id"]),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			objects[idx]["uniqueId"] = fmt.Sprintf("%s-%s", *request.Cursor.CollectionID, memberID)
			objects[idx]["groupId"] = *request.Cursor.CollectionID
		}
	default:
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Entity ID %v is not supported.", request.EntityExternalID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	return objects, ParseNextCursor(response.NextPageToken, request.Cursor), nil
}

func ParseNextCursor(
	nextPageToken *string,
	oldCursor *pagination.CompositeCursor[string],
) (nextCursor *pagination.CompositeCursor[string]) {
	// If there is a nextPageToken in the response, let's set the nextCursor.
	if nextPageToken != nil {
		nextCursor = &pagination.CompositeCursor[string]{
			Cursor: nextPageToken,
		}

		// [MemberEntity] If there was a cursor in the request, we need to copy the collection ID to the next cursor.
		// [!MemberEntity] This will not do anything.
		if oldCursor != nil {
			nextCursor.CollectionID = oldCursor.CollectionID
		}
	}

	// [MemberEntity] If there was a cursor in the request, we need to copy the collection cursor to the next cursor.
	// This should happen regardless of whether there is a nextPageToken in the response.
	// [!MemberEntity] This will not do anything.
	if oldCursor != nil && oldCursor.CollectionCursor != nil {
		if nextCursor == nil {
			nextCursor = &pagination.CompositeCursor[string]{}
		}

		nextCursor.CollectionCursor = oldCursor.CollectionCursor
	}

	return nextCursor
}
