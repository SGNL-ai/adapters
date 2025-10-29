// Copyright 2025 SGNL.ai, Inc.
package mysql

import (
	"context"
	"fmt"
	"strconv"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from datasources.
type Adapter struct {
	MySQLClient Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		MySQLClient: client,
	}
}

// GetPage is called by SGNL's ingestion service to query a page of objects
// from a datasource.
func (a *Adapter) GetPage(ctx context.Context, request *framework.Request[Config]) framework.Response {
	if err := a.ValidateGetPageRequest(ctx, request); err != nil {
		return framework.NewGetPageResponseError(err)
	}

	return a.RequestPageFromDatasource(ctx, request)
}

// RequestPageFromDatasource requests a page of objects from a datasource.
func (a *Adapter) RequestPageFromDatasource(
	ctx context.Context, request *framework.Request[Config],
) framework.Response {
	req := &Request{
		Username:     request.Auth.Basic.Username,
		Password:     request.Auth.Basic.Password,
		BaseURL:      request.Address,
		PageSize:     request.PageSize,
		EntityConfig: request.Entity,
		Database:     request.Config.Database,
	}

	if request.Cursor != "" {
		cursor, err := strconv.ParseInt(request.Cursor, 10, 64)
		if err == nil {
			req.Cursor = &cursor
		}
	}

	for _, attribute := range request.Entity.Attributes {
		if attribute.UniqueId {
			req.UniqueAttributeExternalID = attribute.ExternalId

			break
		}
	}

	if request.Config.Filters != nil {
		if curFilter, ok := request.Config.Filters[request.Entity.ExternalId]; ok {
			req.Filter = &curFilter
		}
	}

	resp, err := a.MySQLClient.GetPage(ctx, req)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	// An adapter error message is generated if the response status code is not
	// successful (i.e. if not statusCode >= 200 && statusCode < 300).
	if adapterErr := web.HTTPError(resp.StatusCode, resp.RetryAfterHeader); adapterErr != nil {
		return framework.NewGetPageResponseError(adapterErr)
	}

	// The raw JSON objects from the response must be parsed and converted into framework.Objects.
	// Nested attributes are flattened and delimited by the delimiter specified.
	// DateTime values are parsed using the specified DateTimeFormatWithTimeZone.
	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&request.Entity,
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
		nextCursor = strconv.FormatInt(*resp.NextCursor, 10)
	}

	return framework.NewGetPageResponseSuccess(&framework.Page{
		Objects:    parsedObjects,
		NextCursor: nextCursor,
	})
}
