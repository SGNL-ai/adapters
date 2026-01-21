// Copyright 2026 SGNL.ai, Inc.

package hashicorp

import (
	"context"
	"fmt"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	"github.com/sgnl-ai/adapters/pkg/config"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/validation"
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from datasources. It handles pagination and authentication for HashiCorp API requests.
type Adapter struct {
	HashicorpClient Client
	SSRFValidator   validation.SSRFValidator
}

// NewAdapter creates a new instance of the HashiCorp adapter with the provided client.
// It implements the framework.Adapter interface with Config type parameter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		HashicorpClient: client,
		SSRFValidator:   validation.NewDefaultSSRFValidator(),
	}
}

// GetPage is called by SGNL's ingestion service to query a page of objects
// from a datasource. It validates the request and delegates the actual data fetching
// to RequestPageFromDatasource.
func (a *Adapter) GetPage(ctx context.Context, request *framework.Request[Config]) framework.Response {
	if err := a.ValidateGetPageRequest(ctx, request); err != nil {
		return framework.NewGetPageResponseError(err)
	}

	return a.RequestPageFromDatasource(ctx, request)
}

// RequestPageFromDatasource handles the actual data fetching from the HashiCorp API.
// It processes the request configuration, handles pagination, and converts the response
// into the expected format for SGNL's ingestion service.
func (a *Adapter) RequestPageFromDatasource(
	ctx context.Context, request *framework.Request[Config],
) framework.Response {
	var commonConfig *config.CommonConfig
	if request.Config != nil {
		commonConfig = request.Config.CommonConfig
	}

	commonConfig = config.SetMissingCommonConfigDefaults(commonConfig)

	request.Address = strings.TrimSuffix(request.Address, "/")
	if !strings.HasPrefix(request.Address, "https://") {
		request.Address = "https://" + request.Address
	}

	cursor, err := pagination.UnmarshalCursor[string](request.Cursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	hashicorpReq := &Request{
		BaseURL: request.Address,
		Auth: Auth{
			Username:     request.Auth.Basic.Username,
			Password:     request.Auth.Basic.Password,
			AuthMethodID: request.Config.AuthMethodID,
		},
		PageSize:              request.PageSize,
		EntityExternalID:      request.Entity.ExternalId,
		Attributes:            request.Entity.Attributes,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
		Cursor:                cursor,
		EntityConfig:          request.Config.EntityConfig,
	}

	resp, err := a.HashicorpClient.GetPage(ctx, hashicorpReq)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	// An adapter error message is generated if the response status code is not
	// successful (i.e. if not statusCode >= 200 && statusCode < 300).
	if adapterErr := web.HTTPError(resp.StatusCode, resp.RetryAfterHeader); adapterErr != nil {
		return framework.NewGetPageResponseError(adapterErr)
	}

	// Convert the JSON response into the expected SGNL adapter object format.
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

	page := &framework.Page{
		Objects: parsedObjects,
	}

	if resp.NextCursor != nil {
		page.NextCursor = *resp.NextCursor
	}

	return framework.NewGetPageResponseSuccess(page)
}
