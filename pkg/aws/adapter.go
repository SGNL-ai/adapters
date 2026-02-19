// Copyright 2026 SGNL.ai, Inc.

package aws

import (
	"context"
	"fmt"
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
	cursor, err := pagination.UnmarshalCursor[string](request.Cursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	var accountIDRequested bool

	if request.Entity.Attributes != nil {
		for _, attr := range request.Entity.Attributes {
			if attr.ExternalId == AccountID {
				accountIDRequested = true

				break
			}
		}
	}

	awsReq := &Request{
		Auth: Auth{
			AccessKey: request.Auth.Basic.Username,
			SecretKey: request.Auth.Basic.Password,
			Region:    request.Config.Region,
		},
		MaxItems:              int32(request.PageSize),
		EntityExternalID:      request.Entity.ExternalId,
		Cursor:                cursor,
		AccountIDRequested:    accountIDRequested,
		EntityConfig:          request.Config.EntityConfig,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
		ResourceAccountRoles:  request.Config.ResourceAccountRoles,
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

		web.WithJSONPathAttributeNames(),
		web.WithDateTimeFormats(
			[]web.DateTimeFormatWithTimeZone{
				// The AWS SDK specifies ISO 8601 for all timestamps. RFC 3339, a subset of ISO 8601, is used in the API responses.
				// We'll use RFC 3339's predefined layout since Go doesn't natively support ISO 8601. However, because RFC 3339
				// compliance isn't guaranteed, we may need to add more formats as necessary.
				//
				// [Example SDK Usage]: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/iam@v1.33.1/types#RoleLastUsed
				// [RFC 3339]: https://datatracker.ietf.org/doc/html/rfc3339
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
