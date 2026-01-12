// Copyright 2025 SGNL.ai, Inc.
package servicenow

import (
	"context"
	"fmt"
	"strings"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	"github.com/sgnl-ai/adapters/pkg/auth"
	"github.com/sgnl-ai/adapters/pkg/config"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from datasources.
type Adapter struct {
	ServicenowClient Client
}

// NewAdapter instantiates a new Adapter.
func NewAdapter(client Client) framework.Adapter[Config] {
	return &Adapter{
		ServicenowClient: client,
	}
}

// GetPage is called by SGNL's ingestion service to query a page of objects
// from a datasource.
func (a *Adapter) GetPage(ctx context.Context, request *framework.Request[Config]) framework.Response {
	if request.Address == MockServiceNowAddress {
		return MockServiceNowGetPage(ctx, request)
	}

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

	var authorizationHeader string

	switch {
	case request.Auth.Basic != nil:
		authorizationHeader = auth.BasicAuthHeader(request.Auth.Basic.Username, request.Auth.Basic.Password)
	case request.Auth.HTTPAuthorization != "":
		authorizationHeader = request.Auth.HTTPAuthorization
	}

	servicenowReq := &Request{
		BaseURL:               request.Address,
		AuthorizationHeader:   authorizationHeader,
		PageSize:              request.PageSize,
		EntityExternalID:      request.Entity.ExternalId,
		APIVersion:            request.Config.APIVersion,
		Attributes:            request.Entity.Attributes,
		RequestTimeoutSeconds: *commonConfig.RequestTimeoutSeconds,
		CustomURLPath:         request.Config.CustomURLPath,
	}

	var (
		resp *Response
		err  *framework.Error

		// Advanced filters related variables.
		implicitFilters      = map[string][]EntityFilter{}
		relatedFilters       = map[string][]EntityAndRelatedEntityFilter{}
		usingImplicitFilters = false
		usingRelatedFilters  = false
	)

	servicenowReq.Cursor = nil
	if request.Cursor != "" {
		servicenowReq.Cursor = &request.Cursor
	}

	if request.Config.AdvancedFilters != nil && len(request.Config.AdvancedFilters.ScopedObjects) != 0 {
		implicitFilters = ExtractImplicitFilters(*request.Config.AdvancedFilters)
		relatedFilters = ExtractRelatedFilters(*request.Config.AdvancedFilters)

		usingImplicitFilters = len(implicitFilters[request.Entity.ExternalId]) > 0
		usingRelatedFilters = len(relatedFilters[request.Entity.ExternalId]) > 0

		if usingImplicitFilters && usingRelatedFilters {
			return framework.NewGetPageResponseError(&framework.Error{
				Message: fmt.Sprintf("Cannot use both implicit and related filters for entity: %s.", request.Entity.ExternalId),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			})
		}
	}

	if usingImplicitFilters {
		resp, err = a.GetPageUsingImplicitFilters(ctx, implicitFilters[request.Entity.ExternalId], *servicenowReq)
		if err != nil {
			return framework.NewGetPageResponseError(err)
		}
	} else if usingRelatedFilters {
		resp, err = a.GetPageUsingRelatedFilters(ctx, relatedFilters[request.Entity.ExternalId], *servicenowReq)
		if err != nil {
			return framework.NewGetPageResponseError(err)
		}
	} else {
		if request.Config.Filters != nil {
			if filter, found := request.Config.Filters[request.Entity.ExternalId]; found {
				servicenowReq.Filter = &filter
			}
		}

		resp, err = a.ServicenowClient.GetPage(ctx, servicenowReq)
		if err != nil {
			return framework.NewGetPageResponseError(err)
		}
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

		// nolint: lll
		// The below formats are the defaults specified by Servicenow, however users are able to override the
		// global date or time format with a personal preference. TODO [sc-16472].
		// https://docs.servicenow.com/bundle/vancouver-platform-administration/page/administer/time/reference/r_FormatDateAndTimeFields.html
		web.WithDateTimeFormats(
			[]web.DateTimeFormatWithTimeZone{
				{Format: "2006-01-02 15:04:05", HasTimeZone: false},
				{Format: time.DateOnly, HasTimeZone: false},
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

	page := &framework.Page{
		Objects: parsedObjects,
	}

	if resp.NextCursor != nil {
		page.NextCursor = *resp.NextCursor
	}

	return framework.NewGetPageResponseSuccess(page)
}

// GetFilteredEntities retrieves the filtered entities for a given EntityFilter object and cursor.
// For example, if the filter object is:
//
//	{
//	    "scopeEntity": "sys_user_group",
//	    "scopeEntityFilter": "sys_id=0270c251c3200200be647bfaa2d3aea6",
//	    "members": [
//	      {
//	        "memberEntity": "sys_user",
//	        "memberEntityFilter": "user.active=true"
//	      }
//	    ]
//	}
//
// we retrieve the all the member entities (i.e. users).
// If the `members` array is empty, then only the scope entity is retrieved.
// By default, we only retrieve the `sys_id` attribute for each entity. Additional attributes
// can be specified by the `attrs` parameter.
// The CompositeCursor is returned as a separate variable for convenience, but it is identical to
// `Response.NextCursor`.
func (a *Adapter) GetFilteredEntities(
	ctx context.Context,
	baseReq Request,
	filter EntityFilter,
	filterCursor *ImplicitFilterCursor,
	attrs []*framework.AttributeConfig,
) (*Response, *pagination.CompositeCursor[string], *framework.Error) {
	// Make a copy of attrs to avoid modifying the original request attributes.
	additionalAttrs := prefixAndCopyAttributes(attrs, "")

	// The scope entity is equivalent to the parent/collection entity.
	scopeEntityReq := &Request{
		BaseURL:               baseReq.BaseURL,
		AuthorizationHeader:   baseReq.AuthorizationHeader,
		PageSize:              baseReq.PageSize,
		EntityExternalID:      filter.ScopeEntity,
		Filter:                &filter.ScopeEntityFilter,
		APIVersion:            baseReq.APIVersion,
		RequestTimeoutSeconds: baseReq.RequestTimeoutSeconds,
	}

	if filterCursor != nil && filterCursor.Cursor != nil && filterCursor.Cursor.CollectionCursor != nil {
		scopeEntityReq.Cursor = filterCursor.Cursor.CollectionCursor
	}

	// Optimization: Only apply filters if the scope entity is the same as the base request entity.
	// For example, if the request is for users, there's no need to apply user filters to groups.
	if filter.ScopeEntity == baseReq.EntityExternalID {
		scopeEntityReq.Attributes = additionalAttrs
	}

	scopeEntities, err := a.ServicenowClient.GetPage(ctx, scopeEntityReq)
	if err != nil {
		return nil, nil, err
	}

	// If the filter we're evaluating doesn't have any member entities, then we're done.
	// The filtered scope entities are all that need to be retrieved.
	// Populate the cursor for the next page of scope entities, if required.
	if len(filter.Members) == 0 {
		var nextCursor *pagination.CompositeCursor[string]

		if scopeEntities.NextCursor != nil {
			nextCursor = &pagination.CompositeCursor[string]{
				CollectionCursor: scopeEntities.NextCursor,
			}
		}

		return scopeEntities, nextCursor, nil
	}

	// After retrieving the scope entities (i.e. the parent/collection entity), we use these IDs
	// to filter the scope entity's children (i.e. the MemberEntity).
	scopeEntityIDs := extractAttrsFromObjs(uniqueIDAttribute, scopeEntities)

	memberEntityReq := &Request{
		BaseURL:               baseReq.BaseURL,
		AuthorizationHeader:   baseReq.AuthorizationHeader,
		PageSize:              baseReq.PageSize,
		APIVersion:            baseReq.APIVersion,
		RequestTimeoutSeconds: baseReq.RequestTimeoutSeconds,
	}

	if filterCursor != nil && filterCursor.Cursor != nil && filterCursor.Cursor.Cursor != nil {
		memberEntityReq.Cursor = filterCursor.Cursor.Cursor
	}

	member := filter.Members[filterCursor.MemberFilterIndex]
	oldMemberEntityFilter := member.MemberEntityFilter
	newMemberEntityFilter := ""

	switch member.MemberEntity {
	case User:
		// The `group` field is actually a field on the `sys_user_grmember` table, not on the `sys_user` table.
		// The only way to retrieve a Group's users is using the above table.
		newMemberEntityFilter = "groupIN" + strings.Join(scopeEntityIDs, ",")

		// Since we use the `sys_user_grmember` table to retrieve users, all user object attributes
		// must be prefixed with `user.` as demonstrated by the ServiceNow API response.
		additionalAttrs = append(additionalAttrs, prefixAndCopyAttributes(additionalAttrs, "user.")...)

		memberEntityReq.EntityExternalID = GroupMember
	default:
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Member entity %s is not supported.", member.MemberEntity),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	newMemberEntityFilter += "^" + oldMemberEntityFilter
	memberEntityReq.Filter = &newMemberEntityFilter
	memberEntityReq.Attributes = additionalAttrs

	memberEntities, err := a.ServicenowClient.GetPage(ctx, memberEntityReq)
	if err != nil {
		return nil, nil, err
	}

	nextCursor := populateNextCompositeCursorForEntityFilter(scopeEntityReq.Cursor, scopeEntities, memberEntities)

	if nextCursor != nil {
		marshalledCursor, err := pagination.MarshalCursor(nextCursor)
		if err != nil {
			return nil, nil, &framework.Error{
				Message: fmt.Sprintf("Failed to marshal cursor: %v.", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		memberEntities.NextCursor = &marshalledCursor
	}

	return memberEntities, nextCursor, nil
}

// GetPageUsingImplicitFilters retreives a page of entities that are using implicit filters.
func (a *Adapter) GetPageUsingImplicitFilters(
	ctx context.Context,
	implicitFilters []EntityFilter,
	request Request,
) (*Response, *framework.Error) {
	var (
		entityID             = request.EntityExternalID
		err                  *framework.Error
		advancedFilterCursor *AdvancedFilterCursor
		nextCompositeCursor  *pagination.CompositeCursor[string]
	)

	if request.Cursor != nil {
		advancedFilterCursor, err = UnmarshalAdvancedFilterCursor(*request.Cursor)
		if err != nil {
			return nil, err
		}
	}

	// Initialize if we still have an empty cursor.
	if advancedFilterCursor == nil {
		advancedFilterCursor = &AdvancedFilterCursor{ImplicitFilterCursor: &ImplicitFilterCursor{}}
	}

	if advancedFilterCursor.ImplicitFilterCursor == nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Implicit filter cursor is unexpectedly nil for entity: %s.", entityID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	filter := implicitFilters[advancedFilterCursor.ImplicitFilterCursor.EntityFilterIndex]

	if filter.ScopeEntity != Group {
		return nil, &framework.Error{
			Message: fmt.Sprintf("%s is not a supported scope for the current entity: %s.", filter.ScopeEntity, entityID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	filteredEntities, cursor, err := a.GetFilteredEntities(
		ctx,
		request,
		filter,
		advancedFilterCursor.ImplicitFilterCursor,
		request.Attributes,
	)
	if err != nil {
		return nil, err
	}

	nextCompositeCursor = cursor

	// Edge case: User objects are retrieved from the `sys_user_grmember` table which have all user attributes
	// prefixed with `user.`. We need to remove this prefix to match the original request attributes.
	if entityID == User {
		for _, obj := range filteredEntities.Objects {
			for key, value := range obj {
				if !strings.HasPrefix(key, "user.") {
					continue
				}

				newKey := strings.TrimPrefix(key, "user.")
				obj[newKey] = value
				delete(obj, key)
			}
		}
	}

	nextImplicitFilterCursor := PopulateNextImplicitFilterCursor(
		*advancedFilterCursor.ImplicitFilterCursor,
		implicitFilters,
		nextCompositeCursor,
	)

	if nextImplicitFilterCursor != nil {
		nextAdvancedFilterCursor, err := MarshalAdvancedFilterCursor(&AdvancedFilterCursor{
			ImplicitFilterCursor: nextImplicitFilterCursor,
		})
		if err != nil {
			return nil, err
		}

		filteredEntities.NextCursor = &nextAdvancedFilterCursor
	}

	return filteredEntities, nil
}

// GetPageUsingRelatedFilters retrieves a page of entities that are using related filters.
func (a *Adapter) GetPageUsingRelatedFilters(
	ctx context.Context,
	relatedFilters []EntityAndRelatedEntityFilter,
	request Request,
) (*Response, *framework.Error) {
	var (
		entityID             = request.EntityExternalID
		err                  *framework.Error
		advancedFilterCursor *AdvancedFilterCursor
		nextCompositeCursor  *pagination.CompositeCursor[string]
	)

	if request.Cursor != nil {
		advancedFilterCursor, err = UnmarshalAdvancedFilterCursor(*request.Cursor)
		if err != nil {
			return nil, err
		}
	}

	// Initialize if we still have an empty cursor.
	if advancedFilterCursor == nil {
		advancedFilterCursor = &AdvancedFilterCursor{RelatedFilterCursor: &RelatedFilterCursor{}}
	}

	if advancedFilterCursor.RelatedFilterCursor == nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Related filter cursor is unexpectedly nil for entity: %s.", entityID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	filter := relatedFilters[advancedFilterCursor.RelatedFilterCursor.EntityIndex]

	// The entity's filter references a related entity.
	// For example, if entity filter = `assigned_toIN{$.sys_user.sys_id}`, the entity is referencing
	// the user entity. When requesting the related entity, we need to explicitly specify
	// which attributes to request.
	// Therefore, extract the related entity and attribute from the filter.
	relatedEntity, relatedEntityAttribute := ExtractEntityAndAttributeFromJSONPathString(filter.EntityFilter)

	relatedEntities, cursor, err := a.GetFilteredEntities(
		ctx,
		request,
		filter.RelatedEntity,
		&ImplicitFilterCursor{
			Cursor: advancedFilterCursor.RelatedFilterCursor.RelatedEntityCursor,
		},
		[]*framework.AttributeConfig{{ExternalId: relatedEntityAttribute}},
	)
	if err != nil {
		return nil, err
	}

	nextCompositeCursor = cursor

	// Edge case: User objects are retrieved from the `sys_user_grmember` table which have all user attributes
	// prefixed with `user.`. The attribute we need to extract should also be prefixed.
	if relatedEntity == User {
		relatedEntityAttribute = "user." + relatedEntityAttribute
	}

	// Replace `{$.sys_user.sys_id}` with the actual sys_id values as a comma-separated string.
	relatedEntityAttrs := extractAttrsFromObjs(relatedEntityAttribute, relatedEntities)
	updatedFilter := ReplaceEntityAndAttributeInString(filter.EntityFilter, strings.Join(relatedEntityAttrs, ","))

	req := &Request{
		BaseURL:               request.BaseURL,
		AuthorizationHeader:   request.AuthorizationHeader,
		PageSize:              request.PageSize,
		EntityExternalID:      request.EntityExternalID,
		APIVersion:            request.APIVersion,
		Filter:                &updatedFilter,
		RequestTimeoutSeconds: request.RequestTimeoutSeconds,
		Attributes:            request.Attributes,
	}

	if advancedFilterCursor.RelatedFilterCursor.EntityCursor != nil {
		req.Cursor = advancedFilterCursor.RelatedFilterCursor.EntityCursor
	}

	res, err := a.ServicenowClient.GetPage(ctx, req)
	if err != nil {
		return nil, err
	}

	nextRelatedFilterCursor := PopulateNextRelatedFilterCursor(
		*advancedFilterCursor.RelatedFilterCursor,
		res.NextCursor,
		nextCompositeCursor,
		relatedFilters,
	)

	if nextRelatedFilterCursor != nil {
		marshalledCursor, err := MarshalAdvancedFilterCursor(&AdvancedFilterCursor{
			RelatedFilterCursor: nextRelatedFilterCursor,
		})
		if err != nil {
			return nil, err
		}

		res.NextCursor = &marshalledCursor
	}

	return res, nil
}

func extractAttrsFromObjs(attribute string, resp *Response) []string {
	ids := make([]string, 0, len(resp.Objects))

	for _, obj := range resp.Objects {
		rawID, found := obj[attribute]
		if !found {
			continue
		}

		id, ok := rawID.(string)
		if !ok {
			continue
		}

		ids = append(ids, id)
	}

	return ids
}

func prefixAndCopyAttributes(attributes []*framework.AttributeConfig, prefix string) []*framework.AttributeConfig {
	prefixedAttributes := make([]*framework.AttributeConfig, 0, len(attributes))

	for _, attr := range attributes {
		prefixedAttribute := &framework.AttributeConfig{
			ExternalId: attr.ExternalId,
			Type:       attr.Type,
			List:       attr.List,
			UniqueId:   attr.UniqueId,
		}

		if prefix != "" {
			prefixedAttribute.ExternalId = prefix + attr.ExternalId
		}

		prefixedAttributes = append(prefixedAttributes, prefixedAttribute)
	}

	return prefixedAttributes
}

func populateNextCompositeCursorForEntityFilter(
	scopeEntityCurrentCursor *string,
	scopeEntities,
	memberEntities *Response,
) *pagination.CompositeCursor[string] {
	nextCursor := &pagination.CompositeCursor[string]{CollectionCursor: scopeEntityCurrentCursor}

	// If there are more pages for the current member in the (memberEntity, scopeEntity) pair,
	// retrieve more pages of the member entity.
	if memberEntities.NextCursor != nil {
		nextCursor.Cursor = memberEntities.NextCursor

		return nextCursor
	}

	// If there are no more pages for the current member, but there are more pages for the scope entity
	// in the (memberEntity, scopeEntity) pair, retrieve more pages of the scope entity.
	if scopeEntities.NextCursor != nil {
		nextCursor.CollectionCursor = scopeEntities.NextCursor

		return nextCursor
	}

	// Otherwise, we're done syncing.
	return nil
}
