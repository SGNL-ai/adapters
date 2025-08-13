// Copyright 2025 SGNL.ai, Inc.
package crowdstrike

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	graphql "github.com/machinebox/graphql"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// The REST APIs of CrowdStrike are two level.
// A list endpoint is used to list entity IDs.
// A get endpoint is used to get detailed metadata of a specific entity.
// If a REST API has only one endpoint, it is considered as a get endpoint.
type EndpointInfo struct {
	ListEndpoint string
	GetEndpoint  string
}

var (
	EntityExternalIDToEndpoint = map[string]EndpointInfo{
		Device: {
			ListEndpoint: "devices/queries/devices-scroll/v1", // This is implemented over HTTP GET by CRWD
			GetEndpoint:  "devices/entities/devices/v2",       // This is implemented over HTTP POST by CRWD
		},
		EndpointIncident: {
			ListEndpoint: "incidents/queries/incidents/v1",      // This is implemented over HTTP GET by CRWD
			GetEndpoint:  "incidents/entities/incidents/GET/v1", // This is implemented over HTTP POST by CRWD
		},
		Detect: {
			ListEndpoint: "detects/queries/detects/v1",        // This is implemented over HTTP GET by CRWD
			GetEndpoint:  "detects/entities/summaries/GET/v1", // This is implemented over HTTP POST by CRWD
		},
		Alerts: {
			ListEndpoint: "alerts/queries/alerts/v2",  // This is implemented over HTTP GET by CRWD
			GetEndpoint:  "alerts/entities/alerts/v2", // This is implemented over HTTP POST by CRWD
		},
	}
)

// buildGQLRequest builds a GraphQL query based on the request.
func buildGQLRequest(request *Request) (*graphql.Request, *framework.Error) {
	var gqlReq *graphql.Request

	var pageInfo *PageInfo

	if request.GraphQLCursor != nil && request.GraphQLCursor.Cursor != nil {
		var err *framework.Error
		pageInfo, err = DecodePageInfo(request.GraphQLCursor.Cursor)

		if err != nil {
			return nil, err
		}
	}

	builder, builderErr := GetQueryBuilder(request, pageInfo)
	if builderErr != nil {
		return nil, builderErr
	}

	if builder == nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Unsupported Query for provided entity ID: %s", request.EntityExternalID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	query, err := builder.Build(request)
	if err != nil {
		return nil, err
	}

	gqlReq = graphql.NewRequest(query)
	gqlReq.Header.Set("Cache-Control", "no-cache")
	gqlReq.Header.Add("Authorization", request.Token)

	return gqlReq, nil
}

func DecodePageInfo(cursor *string) (*PageInfo, *framework.Error) {
	b, err := base64.StdEncoding.DecodeString(*cursor)
	if err != nil {
		return nil, &framework.Error{
			Message: "Cursor.Cursor base64 decoding failed.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	var pageInfo PageInfo

	err = json.Unmarshal(b, &pageInfo)
	if err != nil {
		return nil, &framework.Error{
			Message: "PageInfo unmarshalling failed.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return &pageInfo, nil
}

func ConstructRESTEndpoint(request *Request, path string) (*string, *framework.Error) {
	if request == nil {
		return nil, &framework.Error{
			Message: "Request is nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	restEntity := ValidRESTEntityExternalIDs[request.EntityExternalID]

	if len(path) == 0 {
		return nil, &framework.Error{
			Message: "The path to fetch the entity from is nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	var (
		offsetStr     string
		offset        int64
		err           error
		escapedFilter string

		cursor = request.RESTCursor
	)

	if request.Filter != nil {
		escapedFilter = url.QueryEscape(*request.Filter)
	}

	// Some REST APIs use integer cursor, some use string cursor.
	// Handle the endpoint generation accordingly.
	if cursor != nil && cursor.Cursor != nil {
		// validate integer cursor
		if restEntity.UseIntCursor {
			offset, err = strconv.ParseInt(*cursor.Cursor, 10, 64)
			if err != nil {
				return nil, &framework.Error{
					Message: fmt.Sprintf("Expected a numeric cursor for entity: %s.", request.EntityExternalID),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_AUTHENTICATION_FAILED,
				}
			}

			if offset < 0 {
				return nil, &framework.Error{
					Message: "Cursor must be greater than 0.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				}
			}
		}

		offsetStr = *cursor.Cursor
	}

	params := fmt.Sprintf("limit=%d", request.PageSize)

	if len(offsetStr) > 0 {
		params += fmt.Sprintf("&offset=%v", offsetStr)
	}

	if len(escapedFilter) > 0 {
		params += fmt.Sprintf("&filter=%v", escapedFilter)
	}

	endpoint := fmt.Sprintf("%s/%s?%s", request.BaseURL, path, params)

	return &endpoint, nil
}
