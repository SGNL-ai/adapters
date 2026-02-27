// Copyright 2026 SGNL.ai, Inc.

package pagerduty

import (
	"context"
	"fmt"
	"time"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	"github.com/sgnl-ai/adapters/pkg/config"
	"github.com/sgnl-ai/adapters/pkg/pagination"

	framework "github.com/sgnl-ai/adapter-framework"
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from datasources.
type Adapter struct {
	PagerDutyClient Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		PagerDutyClient: client,
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
// It calls the PagerDuty datasource client internally to make the datasource request, parses the response,
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

	pagerDutyReq := &Request{
		BaseURL:               request.Address,
		Token:                 request.Auth.HTTPAuthorization,
		PageSize:              request.PageSize,
		EntityExternalID:      request.Entity.ExternalId,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
	}

	if request.Config != nil {
		pagerDutyReq.AdditionalQueryParameters = ParseMap(request.Config.AdditionalQueryParameters)
	}

	// Unmarshal the current cursor.
	cursor, err := pagination.UnmarshalCursor[int64](request.Cursor)
	if err != nil {
		return framework.NewGetPageResponseError(err)
	}

	pagerDutyReq.Cursor = cursor

	res, err := a.PagerDutyClient.GetPage(ctx, pagerDutyReq)
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

		// PagerDuty API dates are represented using ISO 8601.
		// https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTU1-types#datetime.
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
				Message: fmt.Sprintf("Failed to convert PagerDuty response objects: %v.", parserErr),
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

// ParseMap parses a map[string]map[string]any to a map[string]map[string][]string.
// It assumes the input map has already been validated to ensure that the inner values are either []any or string.
// A single string is converted to a []string with a single element.
// A []any is converted to a []string with each element converted to a string.
func ParseMap(inputMap map[string]map[string]any) map[string]map[string][]string {
	if inputMap == nil {
		return nil
	}

	outputMap := make(map[string]map[string][]string, len(inputMap))

	for key, innerMap := range inputMap {
		outputMap[key] = make(map[string][]string, len(innerMap))

		for innerKey, innerValue := range innerMap {
			// The input map has already been validated to ensure that the innerValue is a string or []string.
			var convertedInnerValue []string

			switch v := innerValue.(type) {
			case string:
				convertedInnerValue = []string{v}
			case []any:
				convertedInnerValue = make([]string, 0, len(v))

				for _, e := range v {
					if s, ok := e.(string); ok {
						convertedInnerValue = append(convertedInnerValue, s)
					}
				}
			}

			outputMap[key][innerKey] = convertedInnerValue
		}
	}

	return outputMap
}
