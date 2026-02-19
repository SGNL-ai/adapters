// Copyright 2026 SGNL.ai, Inc.

package crowdstrike

import (
	"context"
	"fmt"
	"strings"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	"github.com/sgnl-ai/adapters/pkg/config"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from datasources.
type Adapter struct {
	// Client provides access to the datasource.
	Client Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		Client: client,
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
	var commonConfig *config.CommonConfig
	if request.Config != nil {
		commonConfig = request.Config.CommonConfig
	}

	commonConfig = config.SetMissingCommonConfigDefaults(commonConfig)

	trimmedAddress := strings.TrimSpace(request.Address)
	sanitizedAddress := strings.ToLower(trimmedAddress)

	if !strings.HasPrefix(sanitizedAddress, "https://") {
		request.Address = "https://" + trimmedAddress
	}

	_, isGraphQLEntity := ValidGraphQLEntityExternalIDs[request.Entity.ExternalId]
	_, isRESTEntity := ValidRESTEntityExternalIDs[request.Entity.ExternalId]

	req := &Request{
		BaseURL:               request.Address,
		Token:                 request.Auth.HTTPAuthorization,
		PageSize:              request.PageSize,
		EntityExternalID:      request.Entity.ExternalId,
		EntityConfig:          &request.Entity,
		Ordered:               request.Ordered,
		Config:                request.Config,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
	}

	// Unmarshal the current cursor.
	if isGraphQLEntity {
		graphQLCursor, unmarshalErr := pagination.UnmarshalCursor[string](request.Cursor)
		if unmarshalErr != nil {
			return framework.NewGetPageResponseError(unmarshalErr)
		}

		req.GraphQLCursor = graphQLCursor
	}

	if isRESTEntity {
		restCursor, unmarshalErr := pagination.UnmarshalCursor[string](request.Cursor)
		if unmarshalErr != nil {
			return framework.NewGetPageResponseError(unmarshalErr)
		}

		req.RESTCursor = restCursor
	}

	if request.Config.Filters != nil {
		if curFilter, ok := request.Config.Filters[request.Entity.ExternalId]; ok && isRESTEntity {
			req.Filter = &curFilter
		}
	}

	resp, err := a.Client.GetPage(ctx, req)
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
		web.WithDateTimeFormats(
			[]web.DateTimeFormatWithTimeZone{
				// The CrowdStrike Document specifies ISO 8601 for all timestamps. RFC 3339, a subset of ISO 8601, is used in the
				// API responses. We'll use RFC 3339's predefined layout since Go doesn't natively support ISO 8601. However,
				// because RFC 3339 compliance isn't guaranteed, we may need to add more formats as necessary.
				{Format: time.RFC3339, HasTimeZone: true},
				{Format: time.RFC3339Nano, HasTimeZone: true},
				{Format: "2006-01-02T15:04:05.000Z0700", HasTimeZone: true},
				{Format: "2006-01-02T15:04:05Z", HasTimeZone: false},
				{Format: "2006-01-02", HasTimeZone: false},
			}...,
		),
	)
	if parserErr != nil {
		return framework.NewGetPageResponseError(
			&framework.Error{
				Message: fmt.Sprintf("Failed to convert datasource response objects: %v.", parserErr),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		)
	}

	page := framework.Page{
		Objects: parsedObjects,
	}

	// Marshal the next cursor.
	if isGraphQLEntity {
		nextCursor, err := pagination.MarshalCursor(resp.NextGraphQLCursor)
		if err != nil {
			return framework.NewGetPageResponseError(err)
		}

		page.NextCursor = nextCursor
	}

	if isRESTEntity {
		nextCursor, err := pagination.MarshalCursor(resp.NextRESTCursor)
		if err != nil {
			return framework.NewGetPageResponseError(err)
		}

		page.NextCursor = nextCursor
	}

	return framework.NewGetPageResponseSuccess(&page)
}
