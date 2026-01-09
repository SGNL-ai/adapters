// Copyright 2026 SGNL.ai, Inc.
package azuread

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
	AzureADClient Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		AzureADClient: client,
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

	var (
		curFilter, parentFilter *string
		cursor                  *pagination.CompositeCursor[string]

		advancedFilterCursor           = AdvancedFilterCursor{}
		advancedFilters                = []EntityFilter{}
		useAdvancedFilters             = false
		advancedFilterMemberExternalID *string
	)

	// Identify if advanced filters are to be applied.
	if request.Config.AdvancedFilters != nil && len(request.Config.AdvancedFilters.ScopedObjects) != 0 {
		if _, found := request.Config.AdvancedFilters.ScopedObjects[request.Entity.ExternalId]; found {
			useAdvancedFilters = true
			advancedFilters = request.Config.AdvancedFilters.ScopedObjects[request.Entity.ExternalId]
		}

		implicitFilters := ExtractImplicitFilters(*request.Config.AdvancedFilters)

		if len(implicitFilters[request.Entity.ExternalId]) > 0 {
			// Implicit generated filters cannot be applied with standard filters.
			if _, found := request.Config.Filters[request.Entity.ExternalId]; found {
				return framework.NewGetPageResponseError(
					&framework.Error{
						Message: fmt.Sprintf(
							"Implicit filters generated for entity `%s` are not allowed with standard filters.",
							request.Entity.ExternalId,
						),
						Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
					},
				)
			}

			useAdvancedFilters = true
			advancedFilters = implicitFilters[request.Entity.ExternalId]
		}
	}

	// nolint: nestif
	if useAdvancedFilters {
		parsedAdvancedFilterCursor, err := UnmarshalAdvancedFilterCursor(request.Cursor)
		if err != nil {
			return framework.NewGetPageResponseError(err)
		}

		advancedFilterCursor = *parsedAdvancedFilterCursor
		cursor = advancedFilterCursor.Cursor

		if validationErr := validateAdvancedFilterCursor(
			advancedFilterCursor, advancedFilters, request.Entity.ExternalId,
		); validationErr != nil {
			return framework.NewGetPageResponseError(validationErr)
		}

		parentAdvancedFilterConfig := advancedFilters[advancedFilterCursor.EntityFilterIndex]

		if len(parentAdvancedFilterConfig.Members) > 0 {
			memberFilters := parentAdvancedFilterConfig.Members[advancedFilterCursor.MemberFilterIndex]

			// curFilter is the filter to apply to the member e.g. groupmembers (user/group)
			if len(memberFilters.MemberEntityFilter) > 0 {
				curFilter = &memberFilters.MemberEntityFilter
			}

			// parentFilter is the filter to apply to the parent e.g. a group containing the groupmembers (user/group)
			if len(parentAdvancedFilterConfig.ScopeEntityFilter) > 0 {
				parentFilter = &parentAdvancedFilterConfig.ScopeEntityFilter
			}

			advancedFilterMemberExternalID = &memberFilters.MemberEntity
		} else {
			curFilter = &parentAdvancedFilterConfig.ScopeEntityFilter
		}
	} else {
		// Unmarshal the current cursor.
		parsedCursor, err := pagination.UnmarshalCursor[string](request.Cursor)
		if err != nil {
			return framework.NewGetPageResponseError(err)
		}

		cursor = parsedCursor

		if request.Config.Filters != nil {
			if filter, found := request.Config.Filters[request.Entity.ExternalId]; found {
				curFilter = &filter
			}
		}

		// Set the parent filter if ApplyFiltersToMembers is enabled and the current entity is a member entity.
		if request.Config.ApplyFiltersToMembers {
			parentEntityExternalID := ValidEntityExternalIDs[request.Entity.ExternalId].memberOf
			if parentEntityExternalID != nil && request.Config.Filters != nil {
				if filter, found := request.Config.Filters[*parentEntityExternalID]; found {
					parentFilter = &filter
				}
			}
		}
	}

	azureadReq := &Request{
		BaseURL:                        request.Address,
		Token:                          request.Auth.HTTPAuthorization,
		PageSize:                       request.PageSize,
		EntityExternalID:               request.Entity.ExternalId,
		APIVersion:                     request.Config.APIVersion,
		Attributes:                     request.Entity.Attributes,
		Cursor:                         cursor,
		Filter:                         curFilter,
		ParentFilter:                   parentFilter,
		RequestTimeoutSeconds:          *commonConfig.RequestTimeoutSeconds,
		UseAdvancedFilters:             useAdvancedFilters,
		AdvancedFilterMemberExternalID: advancedFilterMemberExternalID,
	}

	resp, err := a.AzureADClient.GetPage(ctx, azureadReq)
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

		// TODO [sc-14078]: Remove support for complex attribute names.
		web.WithComplexAttributeNameDelimiter("__"),

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

	// Marshal the next cursor. Wrap the cursor with a AdvancedFilterCursor if applicable
	var nextCursorStr string

	if useAdvancedFilters {
		nextAdvancedFilterCursor := populateNextAdvancedFilterCursor(advancedFilterCursor, advancedFilters, resp.NextCursor)
		if nextCursorStr, err = MarshalAdvancedFilterCursor(nextAdvancedFilterCursor); err != nil {
			return framework.NewGetPageResponseError(err)
		}
	} else {
		if nextCursorStr, err = pagination.MarshalCursor(resp.NextCursor); err != nil {
			return framework.NewGetPageResponseError(err)
		}
	}

	return framework.NewGetPageResponseSuccess(&framework.Page{
		Objects:    parsedObjects,
		NextCursor: nextCursorStr,
	})
}
