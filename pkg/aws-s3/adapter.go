// Copyright 2026 SGNL.ai, Inc.

package awss3

import (
	"context"
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	"github.com/sgnl-ai/adapters/pkg/config"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from datasources.
type Adapter struct {
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

	// Unmarshal the current cursor.
	cursor, err := UnmarshalS3Cursor(request.Cursor)
	if err != nil {
		// Fallback: try to unmarshal using the old pagination.CompositeCursor format.
		// This handles cursors created before headers were cached in the cursor.
		legacyCursor, legacyErr := pagination.UnmarshalCursor[int64](request.Cursor)
		if legacyErr != nil || legacyCursor == nil {
			// Both formats failed, return the original error.
			return framework.NewGetPageResponseError(err)
		}

		// Successfully parsed old format - convert to S3Cursor with nil headers.
		// The datasource will fetch headers from S3 when headers are nil.
		cursor = &S3Cursor{
			Cursor:  legacyCursor.Cursor,
			Headers: nil,
		}
	}

	if request.Config.FileType == nil {
		request.Config.FileType = &DefaultFileType
	}

	awsReq := &Request{
		Auth: Auth{
			AccessKey: request.Auth.Basic.Username,
			SecretKey: request.Auth.Basic.Password,
			Region:    request.Config.Region,
		},
		Bucket:                request.Config.Bucket,
		PathPrefix:            request.Config.Prefix,
		FileType:              *request.Config.FileType,
		PageSize:              request.PageSize,
		EntityExternalID:      request.Entity.ExternalId,
		Cursor:                cursor,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
		AttributeConfig:       request.Entity.Attributes,
	}

	resp, err := a.Client.GetPage(ctx, awsReq)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	// An adapter error message is generated if the response status code is not
	// successful (i.e. if not statusCode >= 200 && statusCode < 300).
	if adapterErr := web.HTTPError(resp.StatusCode, resp.RetryAfterHeader); adapterErr != nil {
		return framework.NewGetPageResponseError(adapterErr)
	}

	parsedObjects, parserErr := web.ConvertJSONObjectList(
		&request.Entity,
		resp.Objects,
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
	nextCursor, err := MarshalS3Cursor(resp.NextCursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	return framework.NewGetPageResponseSuccess(&framework.Page{
		Objects:    parsedObjects,
		NextCursor: nextCursor,
	})
}
