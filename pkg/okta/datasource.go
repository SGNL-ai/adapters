// Copyright 2025 SGNL.ai, Inc.
package okta

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

type DatasourceResponse = []map[string]any

const (
	Users        string = "User"
	Groups       string = "Group"
	GroupMembers string = "GroupMember"
	Applications string = "Application"
)

var (
	// ValidEntityExternalIDs is a set of valid external IDs of entities that can be queried.
	ValidEntityExternalIDs = map[string]struct{}{
		Users:        {},
		Groups:       {},
		GroupMembers: {},
		Applications: {},
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

	// [GroupMembers] For group members, we need to set the `CollectionID` (GroupID) and
	// `CollectionCursor` (GroupCursor).
	if request.EntityExternalID == GroupMembers {
		groupRequest := &Request{
			Token:                 request.Token,
			APIVersion:            request.APIVersion,
			BaseURL:               request.BaseURL,
			EntityExternalID:      Groups,
			PageSize:              1,
			RequestTimeoutSeconds: request.RequestTimeoutSeconds,
		}

		// If the CollectionCursor (GroupCursor) is set, use that as the Cursor
		// for the next call to `GetPage`
		if request.Cursor != nil && request.Cursor.CollectionCursor != nil {
			groupRequest.Cursor = &pagination.CompositeCursor[string]{
				Cursor: request.Cursor.CollectionCursor,
			}
		}

		if request.Cursor == nil {
			request.Cursor = &pagination.CompositeCursor[string]{}
		}

		isEmptyLastPage, cursorErr := pagination.UpdateNextCursorFromCollectionAPI(
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
			groupRequest,
			"id",
		)
		if cursorErr != nil {
			return nil, cursorErr
		}

		// [GroupMembers] If `ConstructEndpoint` hits a page with no `CollectionID` (GroupID) and no
		// `CollectionCursor` (GroupCursor) we should complete the sync at this point.
		if isEmptyLastPage {
			return &Response{
				StatusCode: http.StatusOK,
			}, nil
		}
	}

	// [User/Groups] This verifies that `CollectionID` and `CollectionCursor` are not set.
	// [GroupMembers] This verifies that `CollectionID` is set.
	validationErr := pagination.ValidateCompositeCursor(
		request.Cursor,
		request.EntityExternalID,
		// Send a bool indicating if the entity is a member of a collection.
		request.EntityExternalID == GroupMembers,
	)
	if validationErr != nil {
		return nil, validationErr
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
	req.Header.Add("Content-Type", "application/json;okta-response=omitCredentials,omitCredentialsLinks")

	logger.Info("Sending HTTP request to datasource", fields.RequestURL(endpoint))

	res, err := d.Client.Do(req)
	if err != nil {
		logger.Error("HTTP request to datasource failed",
			fields.RequestURL(endpoint),
			fields.SGNLEventTypeError(),
			zap.Error(err),
		)

		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute Okta request: %v.", err),
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
			fields.RequestURL(endpoint),
			fields.ResponseStatusCode(res.StatusCode),
			fields.ResponseRetryAfterHeader(response.RetryAfterHeader),
			fields.ResponseBody(res.Body),
			fields.SGNLEventTypeError(),
		)

		return response, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to read Okta response: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	objects, frameworkErr := ParseResponse(body)
	if frameworkErr != nil {
		return nil, frameworkErr
	}

	response.NextCursor = pagination.GetNextCursorFromLinkHeader(res.Header.Values("link"))

	// [GroupMembers] Set `id`, `userId` and `groupId`.
	if request.EntityExternalID == GroupMembers {
		for idx, member := range objects {
			memberID, ok := member[uniqueIDAttribute].(string)
			if !ok {
				return nil, &framework.Error{
					Message: fmt.Sprintf(
						"Failed to parse %s field in Okta GroupMember response as string.",
						uniqueIDAttribute,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			objects[idx]["id"] = fmt.Sprintf("%s-%s", memberID, *request.Cursor.CollectionID)
			objects[idx]["userId"] = memberID
			objects[idx]["groupId"] = *request.Cursor.CollectionID
		}

		if response.NextCursor != nil && response.NextCursor.Cursor != nil {
			request.Cursor.Cursor = response.NextCursor.Cursor
		} else {
			request.Cursor.Cursor = nil
		}

		// If we have a next cursor for Groups or Group Members, encode the cursor for the next page.
		// Otherwise, don't set a cursor as this sync is complete.
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

func ParseResponse(body []byte) (objects []map[string]any, err *framework.Error) {
	var data DatasourceResponse

	if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil || data == nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return data, nil
}
