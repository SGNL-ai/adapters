// Copyright 2025 SGNL.ai, Inc.
package azuread

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/extractor"
)

// The following are pulled from the set of relationships on each resource type
// that specify `supports $expand`. This is not an exhaustive set and only contains
// a subset of supported relationships that have anticipated use cases.
// Additional relationships can be added here as required.
var (
	expandableRels = map[string]map[string]struct{}{
		// https://learn.microsoft.com/en-us/graph/api/resources/user?view=graph-rest-1.0#relationships
		User: {
			"manager":           {},
			"directReports":     {},
			"ownedDevices":      {},
			"registeredDevices": {},
		},
		// https://learn.microsoft.com/en-us/graph/api/resources/group?view=graph-rest-1.0#relationships
		Group: {
			"appRoleAssignments": {},
			"owners":             {},
		},
		// https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-1.0#relationships
		Application: {
			"owners": {},
		},
		// https://learn.microsoft.com/en-us/graph/api/resources/device?view=graph-rest-1.0#relationships
		Device: {
			"memberOf":           {},
			"registeredOwners":   {},
			"registeredUsers":    {},
			"transitiveMemberOf": {},
		},
	}
)

// ConstructEndpoint constructs and returns the endpoint to query the datasource.
func ConstructEndpoint(request *Request) (string, *framework.Error) {
	if request == nil {
		return "", nil
	}

	var endpoint string

	// The cursor contains the pageNumber if the request is for one of the following entities. A skipToken otherwise.
	// - Role
	// - RoleAssignmentScheduleRequest
	// - GroupAssignmentScheduleRequest
	if request.EntityExternalID != Role &&
		request.EntityExternalID != RoleAssignmentScheduleRequest &&
		request.EntityExternalID != GroupAssignmentScheduleRequest {
		// [!GroupMembers] This is the cursor to the next page of objects.
		// [GroupMembers] This is the cursor to the next page of Members.
		if request.Cursor != nil && request.Cursor.Cursor != nil {
			endpoint = *request.Cursor.Cursor
		}
	}

	if endpoint != "" {
		return endpoint, nil
	}

	formattedPageSize := strconv.FormatInt(request.PageSize, 10)

	var sb strings.Builder

	// URL Format:
	// [User]        baseURL + "/" + apiVersion + "/users" + formAttributeParams(...)
	// [Group]       baseURL + "/" + apiVersion + "/groups" + formAttributeParams(...)
	// [Application] baseURL + "/" + apiVersion + "/applications" + formAttributeParams(...)
	// [Device]      baseURL + "/" + apiVersion + "/devices" + formAttributeParams(...)
	// [Role]        baseURL + "/" + apiVersion + "/directoryRoles" + formAttributeParams(...)
	// [GroupMember] baseURL + "/" + apiVersion + "/groups/" + groupID + "/members?$select=id&$top=" + pageSize
	//        + ["&$filter=" + filter]
	// [RoleMember]  baseURL + "/" + apiVersion + "/users" + userID
	// 			+ "/transitiveMemberOf/microsoft.graph.directoryRole" + "?$select=id&$top=" + pageSize
	// [RoleAssignment] baseURL + "/" + apiVersion + "/roleManagement/directory/roleAssignments"
	//                  + formAttributeParams(...)
	// [RoleAssignmentScheduleRequest] baseURL + "/" + apiVersion
	// 					+ "/roleManagement/directory/roleAssignmentScheduleRequests" + formAttributeParams(...)
	// [GroupAssignmentScheduleRequest] baseURL + "/" + apiVersion
	// 					+ "/identityGovernance/privilegedAccess/group/assignmentScheduleRequests"
	// 					+ formAttributeParams(...)

	sb.Grow(12 + len(request.BaseURL) + len(request.APIVersion) + len(formattedPageSize))

	sb.WriteString(request.BaseURL)
	sb.WriteString("/")
	sb.WriteString(request.APIVersion)

	switch request.EntityExternalID {
	case User:
		if request.UseAdvancedFilters && request.AdvancedFilterMemberExternalID != nil {
			endpoint, err := createMemberEntityEndpoint(request)
			if err != nil {
				return "", err
			}

			sb.WriteString(endpoint)
		} else {
			sb.WriteString("/users")
		}
	case Group:
		if request.UseAdvancedFilters && request.AdvancedFilterMemberExternalID != nil {
			endpoint, err := createMemberEntityEndpoint(request)
			if err != nil {
				return "", err
			}

			sb.WriteString(endpoint)
		} else {
			sb.WriteString("/groups")
		}
	case Role:
		sb.WriteString("/directoryRoles")
	case Application:
		sb.WriteString("/applications")
	case Device:
		sb.WriteString("/devices")
	case RoleAssignment:
		sb.WriteString("/roleManagement/directory/roleAssignments")
	case RoleAssignmentScheduleRequest:
		sb.WriteString("/roleManagement/directory/roleAssignmentScheduleRequests")
	case GroupAssignmentScheduleRequest:
		sb.WriteString("/identityGovernance/privilegedAccess/group/assignmentScheduleRequests")
	case GroupMember:
		if request.Cursor == nil || request.Cursor.CollectionID == nil {
			return "", &framework.Error{
				Message: "Unable to construct group member endpoint without valid cursor.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		pageSizeStr := strconv.FormatInt(request.PageSize, 10)

		sb.Grow(33 + len(*request.Cursor.CollectionID) + len(pageSizeStr))
		sb.WriteString("/groups/")
		sb.WriteString(*request.Cursor.CollectionID)
		sb.WriteString("/members")
		// Validation logic takes care of ensuring that the advanced filter member external ID is valid.
		// The checks added here are simply added defense.
		if request.UseAdvancedFilters && request.AdvancedFilterMemberExternalID != nil {
			suffixes, found := memberEntityToEndpointSuffix[GroupMember]
			if !found {
				return "", &framework.Error{
					Message: "Unable to construct group member endpoint without valid member entity external ID.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			suffix, found := suffixes[*request.AdvancedFilterMemberExternalID]
			if !found {
				return "", &framework.Error{
					Message: "Provided advanced filter member external ID is invalid.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}
			}

			sb.WriteString(suffix)
		}

		sb.WriteString("?$select=id&$top=")
		sb.WriteString(pageSizeStr)

		if request.Filter != nil {
			escapedFilter := url.QueryEscape(*request.Filter)

			sb.Grow(9 + len(escapedFilter))
			sb.WriteString("&$filter=")
			sb.WriteString(escapedFilter)
		}

		// Additional query parameters added on advanced filter requests.
		if request.UseAdvancedFilters {
			sb.WriteString("&$count=true")
		}
	case RoleMember:
		if request.Cursor == nil || request.Cursor.CollectionID == nil {
			return "", &framework.Error{
				Message: "Unable to construct role member endpoint without valid cursor.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		pageSizeStr := strconv.FormatInt(request.PageSize, 10)

		sb.Grow(73 + len(*request.Cursor.CollectionID) + len(pageSizeStr))
		sb.WriteString("/users/")
		sb.WriteString(*request.Cursor.CollectionID)
		sb.WriteString("/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=")
		sb.WriteString(pageSizeStr)

		if request.Filter != nil {
			escapedFilter := url.QueryEscape(*request.Filter)

			sb.Grow(9 + len(escapedFilter))
			sb.WriteString("&$filter=")
			sb.WriteString(escapedFilter)
		}
	default:
		return "", &framework.Error{
			Message: "Provided entity external ID is invalid.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// [!EntityMembers] For all entities other than group, role members, call `formAttributeParams(...)` to construct
	// query params.
	if request.EntityExternalID != GroupMember && request.EntityExternalID != RoleMember {
		selectParams, err := formAttributeParams(
			request.EntityExternalID,
			request.PageSize,
			request.Skip,
			request.Filter,
			request.UseAdvancedFilters,
			"id",
			expandableRels[request.EntityExternalID],
			request.Attributes...,
		)
		if err != nil {
			return "", err
		}

		sb.WriteString(selectParams)
	}

	endpoint = sb.String()

	return endpoint, nil
}

// nolint: lll
func formAttributeParams(
	entityExternalID string,
	pageSize int64,
	offset int64,
	filter *string,
	useAdvancedFilters bool,
	defaultAttribute string,
	relIsExpandable map[string]struct{},
	attributes ...*framework.AttributeConfig,
) (string, *framework.Error) {
	var sb strings.Builder

	escapedDefaultAttribute := url.QueryEscape(defaultAttribute)

	sb.Grow(9 + len(escapedDefaultAttribute))
	sb.WriteString("?$select=")
	sb.WriteString(escapedDefaultAttribute)

	expandRels := map[string][]*framework.AttributeConfig{}
	complexAttrs := make(map[string]*framework.AttributeConfig)

	for idx, attribute := range attributes {
		switch {
		// The defaultAttribute is added by default to all requests, so don't re-add.
		case attribute.ExternalId == defaultAttribute:
			continue

		case strings.HasPrefix(attribute.ExternalId, "$"):
			attributes, err := extractor.AttributesFromJSONPath(attribute.ExternalId)
			if err != nil {
				return "", &framework.Error{
					Message: fmt.Sprintf(
						"Provided entity attribute external id contains unsupported JSON path expression: %q.",
						attribute.ExternalId,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}
			}

			switch len(attributes) {
			case 0:
				return "", &framework.Error{
					Message: fmt.Sprintf(
						"Unable to extract any attributes from JSON path expression in provided attribute external id: %q.",
						attribute.ExternalId,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}
			case 1:
				// If we only find a single attribute, treat it as a child.
				if idx > 0 || defaultAttribute != "" {
					sb.WriteRune(',')
				}

				sb.WriteString(url.QueryEscape(attributes[0]))

				continue
			case 2:
				parentExternalID, childExternalID := attributes[0], attributes[1]

				if _, found := relIsExpandable[parentExternalID]; found {
					// Expand the parent relationship and select the children.
					expandRels[parentExternalID] = append(
						expandRels[parentExternalID],
						&framework.AttributeConfig{
							ExternalId: childExternalID,
						},
					)

					continue
				}

				return "", &framework.Error{
					Message: fmt.Sprintf(
						"Unsupported parent attribute provided for the current entity type: %q.",
						parentExternalID,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}
			default:
				return "", &framework.Error{
					Message: fmt.Sprintf(
						"Too many attributes extracted from JSON path expression in provided attribute external id. Found: %d. Maximum supported: 2.",
						len(attributes),
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}
			}

		// TODO [sc-14078]: Remove support for `__` delimiter
		// In AzureAD, complex attributes need to either select the
		// parent attribute directly, or expand the parent attribute
		// and select the child attribute.
		case strings.Contains(attribute.ExternalId, "__"):
			attrNames := strings.SplitN(attribute.ExternalId, "__", 2)
			parentExternalID := attrNames[0]
			childExternalID := attrNames[1]

			// Return an error for invalid parent or child attributes.
			if parentExternalID == "" || childExternalID == "" {
				return "", &framework.Error{
					Message: fmt.Sprintf(
						"Provided entity attribute list contains the following unsupported attribute: %s.",
						attribute.ExternalId,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}
			}

			if _, found := relIsExpandable[parentExternalID]; found {
				// Expand the parent relationship and select the children.
				expandRels[parentExternalID] = append(expandRels[parentExternalID],
					&framework.AttributeConfig{
						ExternalId: childExternalID,
					})

				continue
			}

			// Select the parent attribute directly, if not already present.
			if _, found := complexAttrs[parentExternalID]; found {
				continue
			}

			attribute = &framework.AttributeConfig{
				ExternalId: parentExternalID,
			}
			complexAttrs[parentExternalID] = attribute

			if idx > 0 || defaultAttribute != "" {
				sb.WriteRune(',')
			}

			sb.WriteString(url.QueryEscape(attribute.ExternalId))

			continue

		default:
			if idx > 0 || defaultAttribute != "" {
				sb.WriteRune(',')
			}

			sb.WriteString(url.QueryEscape(attribute.ExternalId))
		}
	}

	// TODO [sc-38290]: $expand is not supported for the groups/{groupId}/members endpoint.
	if len(expandRels) > 0 && !(useAdvancedFilters && (entityExternalID == User || entityExternalID == Group)) {
		sb.WriteString("&$expand=")

		var parentAttrExternalIDs []string

		for key := range expandRels {
			parentAttrExternalIDs = append(parentAttrExternalIDs, key)
		}

		sort.Strings(parentAttrExternalIDs)

		for idx, parentAttrExternalID := range parentAttrExternalIDs {
			if idx > 0 {
				sb.WriteRune(',')
			}

			escapedParentAttrExternalID := url.QueryEscape(parentAttrExternalID)

			sb.Grow(len(escapedParentAttrExternalID) + 9)
			sb.WriteString(escapedParentAttrExternalID)
			sb.WriteString("($select=")

			childAttrs := expandRels[parentAttrExternalID]
			for idx, childAttr := range childAttrs {
				if idx > 0 {
					sb.WriteRune(',')
				}

				sb.WriteString(url.QueryEscape(childAttr.ExternalId))
			}

			sb.WriteRune(')')
		}
	}

	// In case of a Role, the API returns an error if pageSize ($top) is used
	//
	// Ref: https://learn.microsoft.com/en-us/graph/paging?tabs=http#how-paging-works:~:text=Different%20APIs%20might%20behave,might%20return%20an%20error.
	if entityExternalID != "Role" {
		pageSizeStr := strconv.FormatInt(pageSize, 10)
		sb.Grow(6 + len(pageSizeStr))
		sb.WriteString("&$top=")
		sb.WriteString(pageSizeStr)
	}

	if isPIMEntity(entityExternalID) {
		offsetStr := strconv.FormatInt(offset, 10)
		sb.Grow(7 + len(offsetStr))
		sb.WriteString("&$skip=")
		sb.WriteString(offsetStr)
	}

	if filter != nil {
		escapedFilter := url.QueryEscape(*filter)

		sb.Grow(9 + len(escapedFilter))
		sb.WriteString("&$filter=")
		sb.WriteString(escapedFilter)
	}

	if useAdvancedFilters {
		sb.WriteString("&$count=true")
	}

	return sb.String(), nil
}

func isPIMEntity(entityExternalID string) bool {
	return entityExternalID == RoleAssignmentScheduleRequest || entityExternalID == GroupAssignmentScheduleRequest
}

func createMemberEntityEndpoint(request *Request) (string, *framework.Error) {
	if request.Cursor == nil || request.Cursor.CollectionID == nil {
		return "", &framework.Error{
			Message: "Unable to construct group member endpoint without valid cursor.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	var sb strings.Builder

	sb.Grow(33 + len(*request.Cursor.CollectionID))
	sb.WriteString("/groups/")
	sb.WriteString(*request.Cursor.CollectionID)
	sb.WriteString("/members")
	// Validation logic takes care of ensuring that the advanced filter member external ID is valid.
	// The checks added here are simply added defense.
	if request.UseAdvancedFilters && request.AdvancedFilterMemberExternalID != nil {
		suffixes, found := memberEntityToEndpointSuffix[GroupMember]
		if !found {
			return "", &framework.Error{
				Message: "Unable to construct group member endpoint without valid member entity external ID.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		suffix, found := suffixes[*request.AdvancedFilterMemberExternalID]
		if !found {
			return "", &framework.Error{
				Message: "Provided advanced filter member external ID is invalid.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			}
		}

		sb.WriteString(suffix)
	}

	return sb.String(), nil
}
