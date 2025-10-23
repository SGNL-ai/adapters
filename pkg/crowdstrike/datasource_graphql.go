// Copyright 2025 SGNL.ai, Inc.
package crowdstrike

// We ingest CrowdStrike data using both GraphQL and REST APIs.
// This file contains the functions and structs that are used to interact with the CrowdStrike GraphQL APIs.

import (
	"context"
	"encoding/json"
	"fmt"

	graphql "github.com/machinebox/graphql"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
	// InnerPageInfo represents the paging details for the next nested level. (1 layer deeper)
	InnerPageInfo *PageInfo
}

type DatasourceResponse struct {
	Entities  ResponseItems `json:"entities"`
	Incidents ResponseItems `json:"incidents"`
}

type ResponseItems struct {
	Nodes    []map[string]interface{} `json:"nodes"`
	PageInfo PageInfo                 `json:"pageInfo"`
}

func (d *Datasource) getGraphQLPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	logger := zaplogger.FromContext(ctx).With(
		fields.RequestEntityExternalID(request.EntityExternalID),
		fields.RequestPageSize(request.PageSize),
	)

	logger.Info("Starting datasource request")

	url := fmt.Sprintf("%s/identity-protection/combined/graphql/%s", request.BaseURL, request.Config.APIVersion)

	logger.Info("Sending HTTP request to datasource", fields.URL(url))

	client := graphql.NewClient(url, graphql.WithHTTPClient(d.Client))

	gqlReq, err := buildGQLRequest(request)
	if err != nil {
		return nil, err
	}

	var respData interface{}
	if err := client.Run(ctx, gqlReq, &respData); err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to run the query: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	bodyBytes, marshalErr := json.Marshal(respData)
	if marshalErr != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to marshal the response: %v.", marshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	response := &Response{
		StatusCode: 200,
	}

	objects, nextCursor, parseErr := parseGraphQLResponse(bodyBytes, request.EntityExternalID)

	if parseErr != nil {
		return nil, parseErr
	}

	response.Objects = objects

	if nextCursor != "" {
		response.NextGraphQLCursor = &pagination.CompositeCursor[string]{Cursor: &nextCursor}
	}

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(response.StatusCode),
		fields.ResponseObjectCount(len(response.Objects)),
		fields.ResponseNextCursor(response.NextGraphQLCursor),
	)

	return response, nil
}

func parseGraphQLResponse(
	body []byte,
	entityExternalID string,
) ([]map[string]any, string, *framework.Error) {
	var data *DatasourceResponse

	unmarshalErr := json.Unmarshal(body, &data)
	if unmarshalErr != nil {
		return nil, "", &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	items := data.Entities
	if entityExternalID == Incident {
		items = data.Incidents
	}

	if items.PageInfo.HasNextPage {
		return items.Nodes, items.PageInfo.EndCursor, nil
	}

	return items.Nodes, "", nil
}
