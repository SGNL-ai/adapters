// Copyright 2026 SGNL.ai, Inc.
package db2

import (
	"context"
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from DB2 datasources.
type Adapter struct {
	DB2Client Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		DB2Client: client,
	}
}

// GetPage is called by SGNL's ingestion service to query a page of objects
// from a datasource.
func (a *Adapter) GetPage(ctx context.Context, request *framework.Request[Config]) framework.Response {
	req, err := NewRequestFromConfig(request)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	return a.RequestPageFromDatasource(ctx, req, &request.Entity)
}

// RequestPageFromDatasource requests a page of objects from a datasource.
func (a *Adapter) RequestPageFromDatasource(
	ctx context.Context, req *Request, entity *framework.EntityConfig,
) framework.Response {
	resp, adapterErr := a.DB2Client.GetPage(ctx, req)
	if adapterErr != nil {
		return framework.NewGetPageResponseError(adapterErr)
	}

	// Check if the response status code indicates an error.
	if httpErr := web.HTTPError(resp.StatusCode, ""); httpErr != nil {
		return framework.NewGetPageResponseError(httpErr)
	}

	// The raw JSON objects from the response must be parsed and converted into framework.Objects.
	// Nested attributes are flattened and delimited by the delimiter specified.
	// DateTime values are parsed using the specified DateTimeFormatWithTimeZone.
	parsedObjects, parserErr := web.ConvertJSONObjectList(
		entity,
		resp.Objects,
		web.WithJSONPathAttributeNames(),
	)
	if parserErr != nil {
		return framework.NewGetPageResponseError(
			&framework.Error{
				Message: fmt.Sprintf("Failed to convert datasource response objects: %v.", parserErr),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		)
	}

	var nextCursor string
	if resp.NextCursor != nil {
		nextCursor = *resp.NextCursor
	}

	return framework.NewGetPageResponseSuccess(&framework.Page{
		Objects:    parsedObjects,
		NextCursor: nextCursor,
	})
}
