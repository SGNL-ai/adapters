// Copyright 2025 SGNL.ai, Inc.
package github

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
	GithubClient Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		GithubClient: client,
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

	// Unmarshal the current cursor.
	cursor, err := pagination.UnmarshalCursor[string](request.Cursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	githubReq := &Request{
		BaseURL:               request.Address,
		Token:                 request.Auth.HTTPAuthorization,
		EnterpriseSlug:        request.Config.EnterpriseSlug,
		APIVersion:            request.Config.APIVersion,
		IsEnterpriseCloud:     request.Config.IsEnterpriseCloud,
		EntityConfig:          &request.Entity,
		EntityExternalID:      request.Entity.ExternalId,
		Cursor:                cursor,
		PageSize:              request.PageSize,
		Organizations:         request.Config.Organizations,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
	}

	resp, err := a.GithubClient.GetPage(ctx, githubReq)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	// An adapter error message is generated if the response status code is not
	// successful (i.e. if not statusCode >= 200 && statusCode < 300).
	if adapterErr := web.HTTPError(resp.StatusCode, resp.RetryAfterHeader); adapterErr != nil {
		return framework.NewGetPageResponseError(adapterErr)
	}

	// GitHub GraphQL API dates are represented using ISO 8601.
	// https://docs.github.com/en/enterprise-cloud@latest/graphql/reference/scalars#datetime.
	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&request.Entity,
		resp.Objects,
		web.WithJSONPathAttributeNames(),
		web.WithDateTimeFormats(
			[]web.DateTimeFormatWithTimeZone{
				// While the API technically specifies ISO 8601, RFC 3339 is a profile (subset) of ISO 8601 and it
				// appears datetimes in API response are RFC 3339 compliant, so we'll be using the RFC 3339 predefined
				// layout since golang does not have built-in support for ISO 8601. However, this cannot be guaranteed
				// so additional formats should be added here as necessary.
				// https://datatracker.ietf.org/doc/html/rfc3339
				{Format: time.RFC3339, HasTimeZone: true},
			}...,
		),
		web.WithLocalTimeZoneOffset(commonConfig.LocalTimeZoneOffset),
	)
	if parserErr != nil {
		return framework.NewGetPageResponseError(
			&framework.Error{
				Message: fmt.Sprintf("Failed to convert datasource response objects: %v.", parserErr),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		)
	}

	// Marshal the next cursor.
	nextCursor, err := pagination.MarshalCursor(resp.NextCursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	return framework.NewGetPageResponseSuccess(&framework.Page{
		Objects:    parsedObjects,
		NextCursor: nextCursor,
	})
}
