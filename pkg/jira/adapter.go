// Copyright 2025 SGNL.ai, Inc.
package jira

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

var DefaultAssetBaseURL = "https://api.atlassian.com/jsm/assets"

// Adapter implements the framework.Adapter interface to query pages of objects
// from datasources.
type Adapter struct {
	JiraClient Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		JiraClient: client,
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
// It calls the Jira datasource client internally to make the datasource request, parses the response,
// and handles any errors.
// It also handles parsing the current cursor and generating the next cursor.
func (a *Adapter) RequestPageFromDatasource(
	ctx context.Context,
	request *framework.Request[Config],
) framework.Response {
	var commonConfig *config.CommonConfig
	if request.Config != nil {
		commonConfig = request.Config.CommonConfig
	}

	commonConfig = config.SetMissingCommonConfigDefaults(commonConfig)

	sanitizedAddress := strings.TrimSpace(strings.ToLower(request.Address))
	if !strings.HasPrefix(sanitizedAddress, "https://") {
		request.Address = "https://" + request.Address
	}

	jiraReq := &Request{
		BaseURL:               request.Address,
		Username:              request.Auth.Basic.Username,
		Password:              request.Auth.Basic.Password,
		PageSize:              request.PageSize,
		EntityExternalID:      request.Entity.ExternalId,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
	}

	if request.Config != nil {
		// TODO: Remove this after fully deprecating the legacy Issue endpoint.
		if request.Config.EnhancedIssueSearch && request.Entity.ExternalId == Issue {
			request.Entity.ExternalId = EnhancedIssue
			jiraReq.EntityExternalID = EnhancedIssue
		}

		switch request.Entity.ExternalId {
		case Issue, EnhancedIssue:
			jiraReq.IssuesJQLFilter = request.Config.IssuesJQLFilter
		case Object:
			jiraReq.ObjectsQLQuery = request.Config.ObjectsQLQuery

			if request.Config.AssetBaseURL == nil {
				jiraReq.AssetBaseURL = &DefaultAssetBaseURL
			} else {
				jiraReq.AssetBaseURL = request.Config.AssetBaseURL
			}
		}
	}

	// Unmarshal the current cursor.
	cursor, err := pagination.UnmarshalCursor[string](request.Cursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	jiraReq.Cursor = cursor

	res, err := a.JiraClient.GetPage(ctx, jiraReq)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	// An adapter error message is generated if the response status code is not
	// successful (i.e. if not statusCode >= 200 && statusCode < 300).
	adapterErr := web.HTTPError(res.StatusCode, res.RetryAfterHeader)
	if adapterErr != nil {
		return framework.NewGetPageResponseError(adapterErr)
	}

	// The raw JSON objects from the response must be parsed and converted into framework.Objects.
	// Nested attributes are flattened and delimited by the delimiter specified.
	// DateTime values are parsed using the specified DateTimeFormatWithTimeZone.
	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&request.Entity,
		res.Objects,

		// TODO [sc-14078]: Remove support for complex attribute names.
		web.WithComplexAttributeNameDelimiter("__"),

		web.WithJSONPathAttributeNames(),

		web.WithDateTimeFormats(
			[]web.DateTimeFormatWithTimeZone{
				{Format: "2006-01-02T15:04:05.000Z0700", HasTimeZone: true},
				{Format: time.DateOnly, HasTimeZone: false},
			}...,
		),
		web.WithLocalTimeZoneOffset(commonConfig.LocalTimeZoneOffset),
	)
	if parserErr != nil {
		return framework.NewGetPageResponseError(
			&framework.Error{
				Message: fmt.Sprintf("Failed to convert Jira response objects: %v.", parserErr),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		)
	}

	// Marshal the next cursor.
	nextCursor, err := pagination.MarshalCursor(res.NextCursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	return framework.NewGetPageResponseSuccess(&framework.Page{
		Objects:    parsedObjects,
		NextCursor: nextCursor,
	})
}
