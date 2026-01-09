// Copyright 2026 SGNL.ai, Inc.
package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
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

// Datasource directly implements a Client interface to allow querying an external datasource.
type Datasource struct {
	Client *http.Client
}

type Entity struct {
	// The RequiredAttributes array specifies attributes necessary to create the uniqueID in the response.
	// In OrganizationUser, $.node.id field is the userId. To create the uniqueId, both userId and orgId are needed.
	// orgId is required for establishing relationships between OrganizationUser and Organization.
	RequiredAttributes []string
	// UniqueExternalIDAttribute is the unique attribute of the entity.
	UniqueExternalIDAttribute string
	// ParsePath is the path to the entity in the GraphQL response. This will be used to unmarshal the response.
	// The length of the ParsePath should always be a multiple of 3. Refer to the DataToObjects function for more details.
	ParsePath []string
	// MemberOf is the external ID of the collection entity that the entity belongs to.
	MemberOf *string
	// CollectionEntityConfig contains the required collection attributes to request for the entity.
	// This matches the default entities to support the test case matching.
	CollectionEntityConfig *framework.EntityConfig
	// CollectionAttribute is the attribute of the collection that needs to be used as the collectionID
	CollectionAttribute string
	// isRestAPI is a boolean that indicates whether the entity is retrieved using the GitHub REST APIs.
	isRestAPI bool
}

type ContainerLayers struct {
	Enterprise   *EnterpriseInfo
	Organization *OrganizationInfo
	Repository   *RepositoryInfo
	Label        *LabelInfo
	Issue        *IssueInfo
	PullRequest  *PullRequestInfo
}

// This entity represents the pagination metadata for each layer of the GraphQL response.
type PageInfo struct {
	HasNextPage bool `json:"hasNextPage"`

	// EndCursor represents the next cursor information.
	EndCursor *string `json:"endCursor"`

	// OrganizationOffset represents the offset within request.Organizations slice.
	OrganizationOffset int `json:"organizationOffset"`

	// InnerPageInfo represents the paging details for the next nested level. (1 layer deeper)
	InnerPageInfo *PageInfo
}

// This is the generic representation of a single layer of GitHub GraphQL response.
type EntitiesInfo struct {
	PageInfo *PageInfo       `json:"pageInfo"`
	Edges    json.RawMessage `json:"edges"`
	Nodes    json.RawMessage `json:"nodes"`
}

type DatasourceResponse struct {
	Data   *Data       `json:"data"`
	Errors []ErrorInfo `json:"errors"`
}

type ErrorInfo struct {
	Message string `json:"message"`
}

type Data struct {
	Enterprise   *EnterpriseInfo   `json:"enterprise"`
	Organization *OrganizationInfo `json:"organization"`
}

type OrganizationDatasourceResponse struct {
	Data   *OrganizationData `json:"data"`
	Errors []ErrorInfo       `json:"errors"`
}

type OrganizationData struct {
	Organization *map[string]any `json:"organization"`
}

type EnterpriseInfo struct {
	ID            *string       `json:"id"`
	Organizations *EntitiesInfo `json:"organizations"`
}

type OrganizationInfo struct {
	ID           *string       `json:"id"`
	Users        *EntitiesInfo `json:"membersWithRole"`
	Repositories *EntitiesInfo `json:"repositories"`
	Teams        *EntitiesInfo `json:"teams"`
}

type RepositoryInfo struct {
	ID            *string       `json:"id"`
	Collaborators *EntitiesInfo `json:"collaborators"`
	Labels        *EntitiesInfo `json:"labels"`
	Issues        *EntitiesInfo `json:"issues"`
	PullRequests  *EntitiesInfo `json:"pullRequests"`
}

type IssueInfo struct {
	ID           *string       `json:"id"`
	Assignees    *EntitiesInfo `json:"assignees"`
	Participants *EntitiesInfo `json:"participants"`
	Labels       *EntitiesInfo `json:"labels"`
}

type LabelInfo struct {
	ID           *string       `json:"id"`
	Issues       *EntitiesInfo `json:"issues"`
	PullRequests *EntitiesInfo `json:"pullRequests"`
}

type PullRequestInfo struct {
	ID           *string       `json:"id"`
	ChangedFiles *EntitiesInfo `json:"files"`
	Reviews      *EntitiesInfo `json:"latestOpinionatedReviews"`
	Commits      *EntitiesInfo `json:"commits"`
	Assignees    *EntitiesInfo `json:"assignees"`
	Participants *EntitiesInfo `json:"participants"`
}

const (
	Organization           string = "Organization"
	OrganizationUser       string = "OrganizationUser"
	User                   string = "User"
	OVDE                   string = "$.node.organizationVerifiedDomainEmails"
	Team                   string = "Team"
	TeamMember             string = "$.members.edges"
	TeamRepository         string = "$.repositories.edges"
	Repository             string = "Repository"
	RepositoryCollaborator string = "$.collaborators.edges"
	Collaborator           string = "Collaborator"
	Label                  string = "Label"
	IssueLabel             string = "IssueLabel"
	PullRequestLabel       string = "PullRequestLabel"
	Issue                  string = "Issue"
	IssueAssignee          string = "IssueAssignee"
	IssueParticipant       string = "IssueParticipant"
	PullRequest            string = "PullRequest"
	PullRequestChangedFile string = "PullRequestChangedFile"
	PullRequestReview      string = "PullRequestReview"
	PullRequestCommit      string = "PullRequestCommit"
	PullRequestAssignee    string = "PullRequestAssignee"
	PullRequestParticipant string = "PullRequestParticipant"
	SecretScanningAlert    string = "SecretScanningAlert"
)

var (
	// ValidEntityExternalIDs is a set of valid external IDs of entities that can be queried.
	// The map value is the Entity struct which contains the unique ID attribute.
	ValidEntityExternalIDs = map[string]Entity{
		Organization: {
			UniqueExternalIDAttribute: "id",
			CollectionEntityConfig:    PopulateOrganizationCollectionConfig(),
			ParsePath:                 []string{"Enterprise", "Organizations", "Nodes"},
		},
		OrganizationUser: {
			UniqueExternalIDAttribute: "uniqueId",
			RequiredAttributes:        []string{"$.node.id", "orgId"},
			CollectionAttribute:       "login",
			MemberOf: func() *string {
				s := Organization

				return &s
			}(),
			ParsePath: []string{"Organization", "Users", "Edges"},
		},
		User: {
			UniqueExternalIDAttribute: "id",
			ParsePath:                 []string{"Enterprise", "Organizations", "Nodes", "Organization", "Users", "Nodes"},
		},
		Team: {
			UniqueExternalIDAttribute: "id",
			ParsePath:                 []string{"Enterprise", "Organizations", "Nodes", "Organization", "Teams", "Nodes"},
		},
		// TeamMember is a child entity of Team.
		TeamMember: {
			UniqueExternalIDAttribute: "$.node.id",
		},
		// TeamRepository is a child entity of Team.
		TeamRepository: {
			UniqueExternalIDAttribute: "$.node.id",
		},
		Repository: {
			UniqueExternalIDAttribute: "id",
			ParsePath:                 []string{"Enterprise", "Organizations", "Nodes", "Organization", "Repositories", "Nodes"},
		},
		// RepositoryCollaborator is a child entity of Repository.
		RepositoryCollaborator: {
			UniqueExternalIDAttribute: "id",
		},
		Collaborator: {
			UniqueExternalIDAttribute: "id",
			ParsePath: []string{"Enterprise", "Organizations", "Nodes", "Organization",
				"Repositories", "Nodes", "Repository", "Collaborators", "Nodes"},
		},
		Label: {
			UniqueExternalIDAttribute: "id",
			ParsePath: []string{"Enterprise", "Organizations", "Nodes", "Organization",
				"Repositories", "Nodes", "Repository", "Labels", "Nodes"},
		},
		IssueLabel: {
			UniqueExternalIDAttribute: "uniqueId",
			RequiredAttributes:        []string{"labelId", "id"},
			ParsePath: []string{"Enterprise", "Organizations", "Nodes", "Organization",
				"Repositories", "Nodes", "Repository", "Labels", "Nodes", "Label", "Issues", "Nodes"},
		},
		PullRequestLabel: {
			UniqueExternalIDAttribute: "uniqueId",
			RequiredAttributes:        []string{"labelId", "id"},
			ParsePath: []string{"Enterprise", "Organizations", "Nodes", "Organization",
				"Repositories", "Nodes", "Repository", "Labels", "Nodes", "Label", "PullRequests", "Nodes"},
		},
		Issue: {
			UniqueExternalIDAttribute: "id",
			ParsePath: []string{"Enterprise", "Organizations", "Nodes", "Organization",
				"Repositories", "Nodes", "Repository", "Issues", "Nodes"},
		},
		IssueAssignee: {
			UniqueExternalIDAttribute: "uniqueId",
			RequiredAttributes:        []string{"issueId", "id"},
			ParsePath: []string{"Enterprise", "Organizations", "Nodes", "Organization",
				"Repositories", "Nodes", "Repository", "Issues", "Nodes", "Issue", "Assignees", "Nodes"},
		},
		IssueParticipant: {
			UniqueExternalIDAttribute: "uniqueId",
			RequiredAttributes:        []string{"issueId", "id"},
			ParsePath: []string{"Enterprise", "Organizations", "Nodes", "Organization",
				"Repositories", "Nodes", "Repository", "Issues", "Nodes", "Issue", "Participants", "Nodes"},
		},
		PullRequest: {
			UniqueExternalIDAttribute: "id",
			ParsePath: []string{"Enterprise", "Organizations", "Nodes", "Organization",
				"Repositories", "Nodes", "Repository", "PullRequests", "Nodes"},
		},
		PullRequestChangedFile: {
			UniqueExternalIDAttribute: "uniqueId",
			RequiredAttributes:        []string{"pullRequestId", "path"},
			ParsePath: []string{"Enterprise", "Organizations", "Nodes", "Organization",
				"Repositories", "Nodes", "Repository", "PullRequests", "Nodes", "PullRequest", "ChangedFiles", "Nodes"},
		},
		PullRequestReview: {
			UniqueExternalIDAttribute: "id",
			ParsePath: []string{"Enterprise", "Organizations", "Nodes", "Organization",
				"Repositories", "Nodes", "Repository", "PullRequests", "Nodes", "PullRequest", "Reviews", "Nodes"},
		},
		PullRequestCommit: {
			UniqueExternalIDAttribute: "id",
			ParsePath: []string{"Enterprise", "Organizations", "Nodes", "Organization",
				"Repositories", "Nodes", "Repository", "PullRequests", "Nodes", "PullRequest", "Commits", "Nodes"},
		},
		PullRequestAssignee: {
			UniqueExternalIDAttribute: "uniqueId",
			RequiredAttributes:        []string{"pullRequestId", "id"},
			ParsePath: []string{"Enterprise", "Organizations", "Nodes", "Organization",
				"Repositories", "Nodes", "Repository", "PullRequests", "Nodes", "PullRequest", "Assignees", "Nodes"},
		},
		PullRequestParticipant: {
			UniqueExternalIDAttribute: "uniqueId",
			RequiredAttributes:        []string{"pullRequestId", "id"},
			ParsePath: []string{"Enterprise", "Organizations", "Nodes", "Organization",
				"Repositories", "Nodes", "Repository", "PullRequests", "Nodes", "PullRequest", "Participants", "Nodes"},
		},
		SecretScanningAlert: {
			UniqueExternalIDAttribute: "number",
			isRestAPI:                 true,
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

	MemberOf := ValidEntityExternalIDs[request.EntityExternalID].MemberOf
	if MemberOf != nil && request.EnterpriseSlug != nil {
		collectionReq := &Request{
			BaseURL:               request.BaseURL,
			EnterpriseSlug:        request.EnterpriseSlug,
			APIVersion:            request.APIVersion,
			IsEnterpriseCloud:     request.IsEnterpriseCloud,
			Token:                 request.Token,
			PageSize:              1,
			EntityExternalID:      *MemberOf,
			RequestTimeoutSeconds: request.RequestTimeoutSeconds,
			EntityConfig:          ValidEntityExternalIDs[*MemberOf].CollectionEntityConfig,
			Organizations:         request.Organizations,
		}

		// If the CollectionCursor is set, use that as the Cursor
		// for the next call to `GetPage`
		if request.Cursor != nil && request.Cursor.CollectionCursor != nil {
			collectionReq.Cursor = &pagination.CompositeCursor[string]{
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
			collectionReq,
			ValidEntityExternalIDs[request.EntityExternalID].CollectionAttribute,
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

	// The `CollectionID` is not set in the cursor for an OrganizationUser entity, when
	// a list of Organizations is passed to the request. This fails validation.
	// This is a temporary fix to pass the validation.
	// This could be refactored so that the `CollectionID` has the `OrgLogin` when
	// a list of Organizations is passed to the request.
	// Created sc-37853 to handle the refactor work.
	if request.EntityExternalID == OrganizationUser && request.EnterpriseSlug != nil {
		if validationErr := pagination.ValidateCompositeCursor(
			request.Cursor,
			request.EntityExternalID,
			ValidEntityExternalIDs[request.EntityExternalID].MemberOf != nil,
		); validationErr != nil {
			return nil, validationErr
		}
	}

	reqInfo, reqErr := PopulateRequestInfo(request)
	if reqErr != nil {
		return nil, reqErr
	}

	// Timeout API calls that take longer than the configured timeout.
	apiCtx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	// Note: An empty query in the HTTP request body defaults to http.[NoBody]
	req, err := http.NewRequestWithContext(apiCtx, reqInfo.HTTPMethod,
		reqInfo.Endpoint, strings.NewReader(reqInfo.Query))
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to create request to datasource: %v.", reqErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	req.Header.Add("Authorization", request.Token)
	req.Header.Set("Content-Type", "application/json")

	logger.Info("Sending request to datasource", fields.RequestURL(reqInfo.Endpoint))

	res, err := d.Client.Do(req)
	if err != nil {
		logger.Error("Request to datasource failed",
			fields.RequestURL(reqInfo.Endpoint),
			fields.SGNLEventTypeError(),
			zap.Error(err),
		)

		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute GitHub request: %v.", err),
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
			fields.RequestURL(reqInfo.Endpoint),
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
			Message: fmt.Sprintf("Failed to read GitHub response body: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	var frameworkErr *framework.Error

	if ValidEntityExternalIDs[request.EntityExternalID].isRestAPI {
		response.Objects, response.NextCursor, frameworkErr = ParseRESTResponse(
			body,
			res.Header.Values("Link"),
			reqInfo.OrganizationOffset,
			len(request.Organizations),
		)
	} else {
		response.Objects, response.NextCursor, frameworkErr = ParseGraphQLResponse(
			body,
			request.EntityExternalID,
			request.Cursor,
			len(request.Organizations))
	}

	if frameworkErr != nil {
		return nil, frameworkErr
	}

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(response.StatusCode),
		fields.ResponseObjectCount(len(response.Objects)),
		fields.ResponseNextCursor(response.NextCursor),
	)

	return response, nil
}

func ParseRESTResponse(
	body []byte,
	links []string,
	currentOrganizationOffset int,
	numberOfOrgs int,
) (
	objects []map[string]any,
	nextCursor *pagination.CompositeCursor[string],
	err *framework.Error,
) {
	// We expect the link header to always be present on REST responses. If the header is not found, an error is returned.
	// The link header is parsed to find the next link. If the next link is found, the next cursor is set to the next link.
	// If there is no link associated with rel="next", the next cursor is nil.
	if len(links) == 0 {
		return nil, nil, &framework.Error{
			Message: "Failed to parse the datasource response: Link header is empty or not found.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// The GitHub REST Response is a JSON array of objects.
	if unmarshalErr := json.Unmarshal(body, &objects); unmarshalErr != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	nextCursor = pagination.GetNextCursorFromLinkHeader(links)

	// If the request contains a list of organizations, we need to paginate through the organizations.
	// For that purpose, we need to modify next cursor to set `CollectionId` to the organization offset.
	if numberOfOrgs > 0 {
		// If the next cursor is nil, we need to increment the organization offset to get the next organization,
		// if there are more organizations to paginate through.
		if nextCursor == nil && currentOrganizationOffset+1 < numberOfOrgs {
			nextOrganizationOffset := currentOrganizationOffset + 1
			organizationOffsetAsCollectionID := strconv.Itoa(nextOrganizationOffset)
			nextCursor = &pagination.CompositeCursor[string]{
				CollectionID: &organizationOffsetAsCollectionID,
			}
		} else if nextCursor != nil {
			// If the next cursor is not nil, we need to set the organization offset as the collection ID.
			// There are more pages to paginate through for the current organization.
			organizationOffsetAsCollectionID := strconv.Itoa(currentOrganizationOffset)
			nextCursor.CollectionID = &organizationOffsetAsCollectionID
		}
	}

	return objects, nextCursor, err
}

// ParseGraphQLResponseForOrganization parses the GraphQL response for a singular organization.
// This function is invoked when a list of organizations are passed to the request and data for just
// one organization is expected.
// This is a sample response for a singular organization:
//
//	{
//		"data": {
//		  	"organization": {
//				"id": "O_kgDOCzkBcw",
//				"login": "dh-test-org-2",
//				"email": null,
//				"url": "https://github.com/dh-test-org-2"
//		  	}
//		}
//	}
//
// The above response is not handled by the ParseGraphQLResponse function
// and hence this function is used to parse the response.
func ParseGraphQLResponseForOrganization(
	body []byte,
	currentCursor *pagination.CompositeCursor[string],
	orgCount int,
) (
	objects []map[string]any,
	nextCursor *pagination.CompositeCursor[string],
	err *framework.Error,
) {
	var (
		response        *OrganizationDatasourceResponse
		currentPageInfo *PageInfo
		nextPageInfo    *PageInfo
	)

	if unmarshalErr := json.Unmarshal(body, &response); unmarshalErr != nil || response == nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if response.Errors != nil {
		return nil, nil, ParseErrors(response.Errors)
	}

	if response.Data == nil {
		return nil, nil, &framework.Error{
			Message: "Failed to unmarshal the datasource response: Data not found.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// nolint: prealloc
	objects = make([]map[string]any, 0)

	objects = append(objects, *response.Data.Organization)

	// If there is a current cursor, we need to decode the PageInfo from it.
	if currentCursor != nil && currentCursor.Cursor != nil {
		var err *framework.Error

		currentPageInfo, err = DecodePageInfo(currentCursor.Cursor)

		if err != nil {
			return nil, nil, err
		}
	}

	// If there are more organizations to paginate through, we need to update the PageInfo.
	// The only field in `PageInfo` that needs to be updated is the `OrganizationOffset`.
	if orgCount > 0 {
		if currentCursor == nil {
			nextOffset := 1
			if nextOffset < orgCount {
				nextPageInfo = &PageInfo{
					OrganizationOffset: 0,
				}
			}
		} else if currentPageInfo != nil &&
			currentPageInfo.OrganizationOffset+1 < orgCount {
			nextPageInfo = &PageInfo{
				OrganizationOffset: currentPageInfo.OrganizationOffset + 1,
			}
		}
	}

	if nextPageInfo != nil {
		nextCursor = &pagination.CompositeCursor[string]{}

		cursorBytes, marshalErr := json.Marshal(nextPageInfo)
		if marshalErr != nil {
			return nil, nil, &framework.Error{
				Message: fmt.Sprintf("Failed to create updated cursor: %v.", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		encodedCursor := base64.StdEncoding.EncodeToString(cursorBytes)
		nextCursor.Cursor = &encodedCursor
	}

	return objects, nextCursor, nil
}

func ParseGraphQLResponse(
	body []byte,
	externalID string,
	currentCursor *pagination.CompositeCursor[string],
	orgCount int,
) (
	objects []map[string]any,
	nextCursor *pagination.CompositeCursor[string],
	err *framework.Error,
) {
	if externalID == Organization && orgCount > 0 {
		return ParseGraphQLResponseForOrganization(body, currentCursor, orgCount)
	}

	var response *DatasourceResponse

	var pageInfo *PageInfo

	var currentPageInfo *PageInfo

	if unmarshalErr := json.Unmarshal(body, &response); unmarshalErr != nil || response == nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if response.Errors != nil {
		return nil, nil, ParseErrors(response.Errors)
	}

	if response.Data == nil {
		return nil, nil, &framework.Error{
			Message: "Failed to unmarshal the datasource response: Data not found.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	objects, pageInfo, err = DataToObjects(response.Data, externalID, orgCount == 0)
	if err != nil {
		return nil, nil, err
	}

	// If there is a current cursor, we need to decode the PageInfo from it.
	if currentCursor != nil && currentCursor.Cursor != nil {
		var err *framework.Error

		currentPageInfo, err = DecodePageInfo(currentCursor.Cursor)

		if err != nil {
			return nil, nil, err
		}
	}

	// Update the PageInfo with the new PageInfo from the response
	isPageInfoUpdated, nextPageInfo := UpdatePageInfo(currentPageInfo, pageInfo)

	// If the request does not have an enterprise slug, we are given a list of organizations.
	// If the cursor is nil (indicating that its the first page) we set the organization offset to 0.
	// If there are no more pages for the current organization, we increment the organization offset
	// after checking if there are more organizations to paginate through.
	// Increment the organization offset to get the next organization.
	if orgCount > 0 {
		if nextPageInfo == nil {
			nextPageInfo = &PageInfo{}
		}

		if currentCursor == nil {
			nextPageInfo.OrganizationOffset = 0

			// There are no more pages for the current organization since the first request has returned
			// page size number of objects.
			if !pageInfo.HasNextPage && nextPageInfo.OrganizationOffset+1 < orgCount {
				nextPageInfo.OrganizationOffset = 1

				// There are more organizations to paginate through.
				// We need to reset the pageInfo to get the next organization.
				isPageInfoUpdated = true
			}
		} else if !pageInfo.HasNextPage &&
			currentPageInfo != nil &&
			currentPageInfo.OrganizationOffset+1 < orgCount {
			nextPageInfo.OrganizationOffset = currentPageInfo.OrganizationOffset + 1
			// There are more organizations to paginate through.
			// We need to reset the pageInfo to get the next organization.
			isPageInfoUpdated = true
		}
	}

	// If there is a collectionCursor in play, it should be preserved for the next call in the sync.
	if currentCursor != nil && currentCursor.CollectionCursor != nil {
		nextCursor = &pagination.CompositeCursor[string]{
			CollectionCursor: currentCursor.CollectionCursor,
		}
	}

	if isPageInfoUpdated {
		// We don't want to overwrite the earlier collection cursor if it exists
		if nextCursor == nil {
			nextCursor = &pagination.CompositeCursor[string]{}
		}

		cursorBytes, marshalErr := json.Marshal(nextPageInfo)
		if marshalErr != nil {
			return nil, nil, &framework.Error{
				Message: fmt.Sprintf("Failed to create updated cursor: %v.", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		encodedCursor := base64.StdEncoding.EncodeToString(cursorBytes)
		nextCursor.Cursor = &encodedCursor

		// We only want to set the CollectionID if we know that nextCursorInfo != nil
		// If nextCursorInfo is nil, the current collection's member objects are finished.
		if currentCursor != nil {
			nextCursor.CollectionID = currentCursor.CollectionID
		}
	}

	return objects, nextCursor, nil
}

func ParseErrors(
	errors []ErrorInfo,
) *framework.Error {
	if len(errors) == 0 {
		return &framework.Error{
			Message: "Failed to get the datasource response. Unexpected error format: Errors array is empty.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	var errorMessages = make([]string, 0, len(errors))

	for _, err := range errors {
		if err.Message == "" {
			return &framework.Error{
				Message: "Failed to get the datasource response. Unexpected error format: message is missing.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		errorMessages = append(errorMessages, err.Message)
	}

	return &framework.Error{
		Message: fmt.Sprintf("Failed to get the datasource response: %v.", errorMessages),
		Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
	}
}

/*
 *  DataToObjects will attempt to convert the entities in the response data to objects.
 *  It also populates the PageInfoLayers to create the next cursor.
 *  The conversion is done using each entity's ParsePath. The ParsePath is a list of strings
 *  that represent the path to the entity in the GraphQL response. The ParsePath is also always a
 *  multiple of three because each 'phase' of the parsing can be divided into three steps:
 *  1. Get the container entity.
 *   a. If this is the first phase, we will get the Enterprise or Organization entity from the 'data'.
 *   b. If this is not the first phase, we will use the json saved from the last phase.
 *      results to retrieve the next container entity.
 *  2. Get the first entry from the container and update the PageInfo.
 *  3. Get the raw entity data from the first entry and save it for the next phase.
 *
 *  Ex. ParsePath for User: ["Enterprise", "Organizations", "Nodes", "Organization", "Users", "Nodes"]
 *
 *  Phase 1: ["Enterprise", "Organizations", "Nodes"]
 *  1. We start by getting the EnterpriseInfo struct from the 'Data' struct through ProcessContainer().
 *  2. We then pass this EnterpriseInfo struct to ProcessEntries() to get the
 *     Organizations container, which is an EntitiesInfo struct. The PageInfo field from the EntitiesInfo
 *     is then added to the return page struct 'pageInfo'.
 *  3. We then get the raw data of the Organizations struct through the "Nodes" field in ProcessEntityData().
 *
 *  Phase 2: ["Organization", "Users", "Nodes"]
 *  1. We use the raw Organizations array data from the last phase to get the first Organization
 *     using ProcessContainer().
 *  2. We then pass this OrganizationInfo struct to ProcessEntries() to get the Users container,
 *     which is an EntitiesInfo struct. The PageInfo field from the EntitiesInfo
 *     is then added to the return page struct 'pageInfo'.
 *  3. We then get the raw data of the Users struct through the "Nodes" field in ProcessEntityData().
 *
 *  The raw entity data is then converted to objects using ConvertEntitiesToObjects().
 *  The objects are then injected with common fields using InjectCommonFields().
 */
func DataToObjects(
	data *Data,
	externalID string,
	hasEnterpriseSlug bool,
) (
	objects []map[string]any,
	pageInfo *PageInfo,
	err *framework.Error,
) {
	parsePath := ValidEntityExternalIDs[externalID].ParsePath

	// If the ParsePath does not have the Enterprise slug, we need to remove the first three elements.
	// This is because the Enterprise slug and its subsequent members are not present in the response.
	// This is a ParsePath value for a Repository:
	// `[]string{"Enterprise", "Organizations", "Nodes", "Organization", "Repositories", "Nodes"}`
	// If a list of Organizations is passed, the response contains `"Organization", "Repositories", "Nodes"`.
	// This is true for all entities except the OrganizationUser entity.
	// For the OrganizationUser entity, the ParsePath value remains the same.
	if !hasEnterpriseSlug && externalID != OrganizationUser {
		if len(parsePath) < 3 {
			return nil, nil, &framework.Error{
				Message: "Failed to parse the datasource response: ParsePath is invalid.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		parsePath = parsePath[3:]
	}

	var collections ContainerLayers

	// We store the raw entity data to be used in the next phase of the parsing.
	// At the end of the parse path, this raw entity data will be converted to objects.
	var rawEntityData json.RawMessage

	// This is intialized as an empty slice to prepare for an early return.
	// Early returns happen in cases where intermediate layers can return with no objects.
	// Ex. In the Collaborators sync, it's possible that organization has no repositories,
	// however there may be more organizations to paginate through.
	objects = make([]map[string]any, 0)

	for i := 0; i < len(parsePath); i += 3 {
		containerName, entriesName, rawEntityFieldName := parsePath[i], parsePath[i+1], parsePath[i+2]

		container, err := ProcessContainer(i == 0, containerName, data, &collections, rawEntityData)
		if err != nil {
			return nil, nil, err
		}

		if container == nil {
			return objects, pageInfo, nil
		}

		firstContainerEntry, err := ProcessEntries(container, containerName, entriesName)
		if err != nil {
			return nil, nil, err
		}

		pageInfo, err = AddPageInfoLayerToLeaf(pageInfo, firstContainerEntry.PageInfo)
		if err != nil {
			return nil, nil, err
		}

		rawEntityData, err = ProcessEntityData(rawEntityFieldName, firstContainerEntry)
		if err != nil {
			return nil, nil, err
		}
	}

	objects, err = ConvertEntitiesToObjects(rawEntityData)
	if err != nil {
		return nil, nil, err
	}

	if externalID == OrganizationUser {
		ConvertOVDEAttribute(&objects)
	}

	injectionErr := InjectCommonFields(&objects, &collections, externalID)
	if injectionErr != nil {
		return nil, nil, injectionErr
	}

	return objects, pageInfo, nil
}

func ProcessEntityData(
	rawEntityFieldName string,
	firstContainerEntry *EntitiesInfo,
) (
	json.RawMessage,
	*framework.Error,
) {
	if rawEntityFieldName == "Nodes" {
		return firstContainerEntry.Nodes, nil
	} else if rawEntityFieldName == "Edges" {
		return firstContainerEntry.Edges, nil
	}

	return nil, &framework.Error{
		Message: fmt.Sprintf("Failed to parse the datasource response: %s not found.", rawEntityFieldName),
		Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
	}
}

// ProcessEntries processes the entries of a container using reflection and validates the following assumption:
// 'container' is expected to a be a struct pointer with a field named 'entriesName' of type *EntitiesInfo.
func ProcessEntries(
	container any,
	containerName string,
	entriesName string,
) (
	*EntitiesInfo,
	*framework.Error,
) {
	// Get the reflect.Value of the container and validate that it is a non-nil pointer to a valid memory address.
	containerValue := reflect.ValueOf(container)
	if containerValue.Kind() != reflect.Ptr && !containerValue.IsValid() {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to parse the datasource response: %s not found.", containerName),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Get the element of the container value and validate that it is a non-nil struct type.
	containerElem := containerValue.Elem()
	if !containerElem.IsValid() || containerElem.Kind() != reflect.Struct {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to parse the datasource response: %s not found.", containerName),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Get the 'entriesName' field of the container struct and validate that the field exists.
	entriesValue := containerElem.FieldByName(entriesName)
	if !entriesValue.IsValid() {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to parse the datasource response: %s not found.", entriesName),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Convert the reflected entries field to an *EntitiesInfo and validate that the conversion is successful.
	entriesInfo, ok := entriesValue.Interface().(*EntitiesInfo)
	if !ok || entriesInfo == nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to parse the datasource response: %s not found.", entriesName),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return entriesInfo, nil
}

/*
 *  ProcessContainer retrieves the collection from the data and returns it as an interface.
 *  The 'isStart' parameter is used to determine if the collection is the first layer of the response.
 *  There are two cases:
 *  1. If we are the start of the parsing, we will directly use the data struct to get the Enterprise or Organization.
 *    a. If the collection is nil or the casting to the appropriate collection struct fails, an error is returned.
 *  2. If we are not the start of the parsing, we will use the rawEntityData to get the next collection.
 *    a. If GetCollectionEntityInterface returns an error, it will be returned.
 *    b. If the collection is nil, we will return nil to indicate that the sub-entities should not be parsed.
 *    c. If the collection is nil, and the containerName is Organization, an error is returned.
 */
func ProcessContainer(
	isStart bool,
	containerName string,
	data *Data,
	collections *ContainerLayers,
	rawEntityData json.RawMessage,
) (
	collection any,
	err *framework.Error,
) {
	var ok bool

	if isStart {
		if containerName == "Enterprise" {
			collection = data.Enterprise
			collections.Enterprise, ok = collection.(*EnterpriseInfo)
		} else if containerName == "Organization" {
			collection = data.Organization
			collections.Organization, ok = collection.(*OrganizationInfo)
		} else {
			return nil, &framework.Error{
				Message: fmt.Sprintf("Failed to parse the datasource response: %s not found.", containerName),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		if collection == nil || !ok {
			return nil, &framework.Error{
				Message: fmt.Sprintf("Failed to parse the datasource response: %s not found.", containerName),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}
	} else {
		collection, err = GetCollectionEntityInterface(rawEntityData, containerName, collections)
		if err != nil {
			return nil, err
		}

		if collection == nil {
			// The organization collection must always have a length of 1.
			if containerName == "Organization" {
				return nil, &framework.Error{
					Message: "Failed to parse the datasource response: Organization collection is length 0, expected 1.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			return nil, nil
		}
	}

	return collection, nil
}

// GetCollectionEntityInterface returns uses the entityName to retrieve the appropriate
// collection struct and returns it as an interface.
func GetCollectionEntityInterface(
	entities json.RawMessage,
	entityName string,
	collections *ContainerLayers,
) (
	any,
	*framework.Error,
) {
	switch entityName {
	case Organization:
		return ParseAndAssignCollection[OrganizationInfo](entities, entityName, &collections.Organization)
	case Repository:
		return ParseAndAssignCollection[RepositoryInfo](entities, entityName, &collections.Repository)
	case Label:
		return ParseAndAssignCollection[LabelInfo](entities, entityName, &collections.Label)
	case Issue:
		return ParseAndAssignCollection[IssueInfo](entities, entityName, &collections.Issue)
	case PullRequest:
		return ParseAndAssignCollection[PullRequestInfo](entities, entityName, &collections.PullRequest)
	default:
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to get the collection entity interface: %s not found.", entityName),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}
}

// ParseAndAssignCollection is a helper function for GetCollectionEntityInterface function that
// retrieves the collection from the ParseContainerEntity function and assigns it to the appropriate collection field.
func ParseAndAssignCollection[T any](
	entities json.RawMessage,
	entityName string,
	collectionField **T,
) (
	any,
	*framework.Error,
) {
	collection, err := ParseContainerEntity[T](entities, entityName)
	if err != nil {
		return nil, err
	}

	if collection == nil {
		return nil, nil
	}

	*collectionField = collection

	return collection, nil
}

// Helper function to add new layer as a leaf of the deepest layer of the parentPageInfo.
// The leafPageInfo will become the new deepest layer in the hierarchy.
func AddPageInfoLayerToLeaf(parentPageInfo *PageInfo, leafPageInfo *PageInfo) (*PageInfo, *framework.Error) {
	if leafPageInfo == nil {
		return nil, &framework.Error{
			Message: "Failed to validate LeafPageInfo: PageInfo is nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if parentPageInfo == nil {
		return leafPageInfo, nil
	}

	if parentPageInfo.InnerPageInfo != nil {
		leafPageInfo, _ = AddPageInfoLayerToLeaf(parentPageInfo.InnerPageInfo, leafPageInfo)
	}

	parentPageInfo.InnerPageInfo = leafPageInfo

	return parentPageInfo, nil
}

/*
 *	PageInfo Update Rules: (Multiple Layer GraphQL Pagination)
 *	1. We will initialize our new PageInfo struct by using the currentPageInfo.
 *		a. This number of layers in currentPageInfo will always be less than or equal to the
 *		   number of layers in the newPageInfo. This is because newPageInfo will always contain the
 *		   maximum number of layers for each entity type. However, currentPageInfo could be nil at
 *		   the start of the sync and also gets rid of all the deepest layers that have nil endCursors.
 *	2. We will iterate through the layers starting from the deepest layer and work our way out.
 *		a. If the new PageInfo is nil, we will return false indicating that we've reached the deepest layer.
 *		b. If the current PageInfo is nil, we will initialize it with the new PageInfo to match
 *		   the layer of the newPageInfo.
 *	3. If the recursive call returns false, this indicates that all deeper layers have nil endCursors and
 *	   therefore we can set the innerPageInfo reference to nil.
 *	4. If the current layer has a next page, we will update the endCursor for the current layer and
 *	   return true to indicate that an update has been made at this layer.
 *	4. If the current layer does not have a next page, we will return false and nil which will set the
 *	   innerPageInfo reference for the outer layer to nil.
 *	5. If we iterate through all layers and none of them have a next page, false is returned indicating no
 *	   updates have been made. This return parameter is used to determine if the nextCursor should be set to nil.
 */
func UpdatePageInfo(currentPageInfo, newPageInfo *PageInfo) (bool, *PageInfo) {
	// Return false if the new PageInfo is nil, indicating no updates are needed.
	if newPageInfo == nil {
		return false, nil
	}

	// If the current PageInfo is nil, initialize it with the new PageInfo.
	if currentPageInfo == nil {
		currentPageInfo = &PageInfo{}
	}

	// Recursively update the inner (deeper) PageInfo layers first.
	innerPageWasUpdated, innerPageInfo := UpdatePageInfo(currentPageInfo.InnerPageInfo, newPageInfo.InnerPageInfo)

	// Update the old PageInfo's reference to the inner layer.
	currentPageInfo.InnerPageInfo = innerPageInfo

	// If an inner page was updated, we should not update any of the outer layers.
	if innerPageWasUpdated {
		return true, currentPageInfo
	}

	// If an inner page was not updated, we know that all deeper layers have nil endCursors and can be set to nil.
	currentPageInfo.InnerPageInfo = nil

	// Update the EndCursor of the current layer if the new PageInfo has a next page.
	if newPageInfo.HasNextPage && newPageInfo.EndCursor != nil {
		currentPageInfo.EndCursor = newPageInfo.EndCursor
		// Return true as the PageInfo layer was updated.
		return true, currentPageInfo
	}

	// Return false as there's no next page indicating the end of updates in this layer and all deeper layers.
	return false, nil
}

/*   ParseContainerEntity will attempt to parse the entities json into an array of T.
 *   Case 1: If the unmarshaling fails, an error is returned.
 *   Case 2: If the container is empty, nil is returned to indicate that the sub-entities should
 *           not be parsed. Empty containers are expected for certain entity types. The caller can determine
 *           if an empty container is an error for type T.
 *   Case 3: The container should never have a length greater than 1, as it is expected to
 *           be a single entity. The function will return an error if the container length is not 1.
 *   Case 4: The function will return the first element of the container as a pointer to T.
 */
func ParseContainerEntity[T any](entities json.RawMessage, entityName string) (*T, *framework.Error) {
	var container []T

	unmarshalErr := json.Unmarshal(entities, &container)

	if unmarshalErr != nil || container == nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response. %s not found: %v.", entityName, unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// If the container is empty, we return nil to indicate that the sub-entities should not be parsed.
	// This case will rawEntityData in an empty list of objects being returned for this specific page of the sync.
	if len(container) == 0 {
		return nil, nil
	}

	// We expect a page size of 1, because this function is used to parse sub-entities of a container.
	// All container entities should be retrieved one at a time.
	if len(container) != 1 {
		return nil, &framework.Error{
			Message: fmt.Sprintf("%s container length is: %d, expected: 1.", entityName, len(container)),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Container page size should always be set to 1, so we want to inspect the sub-entities of the first element.
	return &container[0], nil
}

// ParseEntities will validate that the entity is non-nil.
func ParseEntities[T any](entity T, entityName string) (T, *framework.Error) {
	if reflect.ValueOf(entity).IsNil() {
		return entity, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %s not found.", entityName),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return entity, nil
}

func ConvertEntitiesToObjects(entities json.RawMessage) (objects []map[string]any, unmarshalErr *framework.Error) {
	if unmarshalErr := json.Unmarshal(entities, &objects); unmarshalErr != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return objects, nil
}

// OVDE is originally a list of strings, so we need to convert
// it to a list of JSON objects to allow it to be parsed as a child entity.
func ConvertOVDEAttribute(objects *[]map[string]any) {
	for i, object := range *objects {
		if node, ok := object["node"].(map[string]any); ok {
			if emails, ok := node["organizationVerifiedDomainEmails"].([]any); ok {
				emailObjects := make([]any, 0, len(emails))

				for _, email := range emails {
					if emailStr, ok := email.(string); ok {
						emailObjects = append(emailObjects, map[string]any{"email": emailStr})
					}
				}

				node["organizationVerifiedDomainEmails"] = emailObjects
				(*objects)[i]["node"] = node
			}
		}
	}
}

// InjectCommonFields adds common fields to objects for building relationships based on externalID.
// It appends 'id' attributes for specific entities and handles post-processing tasks like creating uniqueId.
// The function ensures that required fields are present and forms relationships between objects.
func InjectCommonFields(
	objects *[]map[string]any,
	container *ContainerLayers,
	externalID string,
) *framework.Error {
	if objects == nil {
		return &framework.Error{
			Message: "Failed to inject common fields: objects is nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	for i := range *objects {
		if container.Enterprise != nil && container.Enterprise.ID != nil {
			(*objects)[i]["enterpriseId"] = *container.Enterprise.ID
		}

		// This switch statement appends 'id' attributes for specific entities requiring it.
		switch externalID {
		// Teams and Repositories can use the OrgId field to build relationships with the Organization.
		// OrganizationUser needs OrgId to build relationships between Organizations and Users.
		case Team, Repository, OrganizationUser:
			if container.Organization == nil || container.Organization.ID == nil {
				return &framework.Error{
					Message: fmt.Sprintf("Organization is nil or orgID is missing for the %s entity.", externalID),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			(*objects)[i]["orgId"] = *container.Organization.ID
		case Label, Issue:
			if container.Repository == nil || container.Repository.ID == nil {
				return &framework.Error{
					Message: fmt.Sprintf("Repository is nil or repositoryID is missing for the %s entity.", externalID),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			(*objects)[i]["repositoryId"] = *container.Repository.ID
		case IssueLabel, PullRequestLabel:
			if container.Label == nil || container.Label.ID == nil {
				return &framework.Error{
					Message: fmt.Sprintf("Label is nil or labelID is missing for the %s entity.", externalID),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			(*objects)[i]["labelId"] = *container.Label.ID
		case IssueAssignee, IssueParticipant:
			if container.Issue == nil || container.Issue.ID == nil {
				return &framework.Error{
					Message: fmt.Sprintf("Issue is nil or issueID is missing for the %s entity.", externalID),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			(*objects)[i]["issueId"] = *container.Issue.ID
		case PullRequestChangedFile, PullRequestReview, PullRequestCommit, PullRequestAssignee, PullRequestParticipant:
			if container.PullRequest == nil || container.PullRequest.ID == nil {
				return &framework.Error{
					Message: fmt.Sprintf("PullRequest is nil or pullRequestID is missing for the %s entity.", externalID),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			(*objects)[i]["pullRequestId"] = *container.PullRequest.ID
		}

		// This switch statement handles post-processing tasks like creating the uniqueId field.
		switch externalID {
		case OrganizationUser:
			InjectUniqueID(objects, i, externalID, "orgId", "$.node.id")
		case IssueAssignee, IssueParticipant:
			InjectUniqueID(objects, i, externalID, "issueId", "id")
		case IssueLabel, PullRequestLabel:
			InjectUniqueID(objects, i, externalID, "labelId", "id")
		case PullRequestChangedFile:
			InjectUniqueID(objects, i, externalID, "pullRequestId", "path")
		case PullRequestAssignee, PullRequestParticipant:
			InjectUniqueID(objects, i, externalID, "pullRequestId", "id")
		}
	}

	return nil
}

// InjectUniqueID adds the uniqueId field to objects in the format of 'attrA-attrB'.
func InjectUniqueID(
	objects *[]map[string]any,
	index int,
	externalID string,
	attrA string,
	attrB string,
) *framework.Error {
	attrAPath, attrBPath := GetAttributePath(attrA), GetAttributePath(attrB)

	attrAValue, err := GetValueFromPath((*objects)[index], attrAPath)
	if err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("Failed to get %s's value for attribute %s: %v.", externalID, attrA, err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	attrBValue, err := GetValueFromPath((*objects)[index], attrBPath)
	if err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("Failed to get %s's value for attribute %s: %v.", externalID, attrB, err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	(*objects)[index]["uniqueId"] = fmt.Sprintf("%s-%s", attrAValue, attrBValue)

	return nil
}

// GetValueFromPath retrieves the value of a nested attribute from an object.
// This function assume that the nested attribute is a string, and throws an error if it is not.
func GetValueFromPath(object map[string]any, path []string) (string, error) {
	var ok bool

	var value any = object

	for _, key := range path {
		object, ok = value.(map[string]any)
		if !ok {
			return "", fmt.Errorf("expected map[string]any, got %T", value)
		}

		value, ok = object[key]
		if !ok {
			return "", fmt.Errorf("key %s does not exist", key)
		}
	}

	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %T", value)
	}

	return strValue, nil
}
