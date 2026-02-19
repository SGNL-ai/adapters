// Copyright 2026 SGNL.ai, Inc.

package workday

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/extractor"
)

// ConstructEndpoint constructs and returns the endpoint to query the datasource.
func ConstructEndpoint(request *Request) (string, *framework.Error) {
	if request == nil {
		return "", &framework.Error{
			Message: "Request is nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	var sb strings.Builder

	var offset int64

	if request.Cursor != nil && request.Cursor.Cursor != nil {
		if *request.Cursor.Cursor <= 0 {
			return "", &framework.Error{
				Message: "Cursor value must be greater than 0.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		offset = *request.Cursor.Cursor
	}

	sb.Grow(len(request.BaseURL) + len(request.APIVersion) + len(request.OrganizationID) + 15)
	sb.WriteString(request.BaseURL)
	sb.WriteString("/api/wql/")
	sb.WriteString(request.APIVersion)
	sb.WriteString("/")
	sb.WriteString(request.OrganizationID)
	sb.WriteString("/data")

	params := url.Values{}
	params.Add("limit", strconv.FormatInt(request.PageSize, 10))
	params.Add("offset", strconv.FormatInt(offset, 10))

	query, err := BuildQuery(request.EntityConfig, request.Ordered)
	if err != nil {
		return "", err
	}

	params.Add("query", query)

	paramString := params.Encode()

	sb.Grow(len(paramString) + 1)
	sb.WriteString("?")
	sb.WriteString(paramString)

	return sb.String(), nil
}

func BuildQuery(entity *framework.EntityConfig, ordered bool) (string, *framework.Error) {
	var sb strings.Builder

	// Create a map to act as a set for storing unique attributes
	uniqueAttributes := make(map[string]struct{})

	var uniqueAttrExternalID string

	// First pass: add all unique attributes to the set
	for _, attribute := range entity.Attributes {
		columnName, err := AttrExternalIDToColumnName(attribute.ExternalId)
		if err != nil {
			return "", &framework.Error{
				Message: fmt.Sprintf("Error extracting column name from attribute %s: %s", attribute.ExternalId, err.Error()),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		if attribute.UniqueId {
			uniqueAttrExternalID = columnName
		}

		uniqueAttributes[columnName] = struct{}{}
	}

	if uniqueAttrExternalID == "" {
		return "", &framework.Error{
			Message: "No unique attribute found for ordering.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	for _, child := range entity.ChildEntities {
		columnName, err := AttrExternalIDToColumnName(child.ExternalId)
		if err != nil {
			return "", &framework.Error{
				Message: fmt.Sprintf("Error extracting column name from attribute %s: %s", child.ExternalId, err.Error()),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		uniqueAttributes[columnName] = struct{}{}
	}

	// Convert the map keys to a slice
	attrs := make([]string, 0, len(uniqueAttributes))
	for attr := range uniqueAttributes {
		attrs = append(attrs, attr)
	}

	// Sort the slice of attributes
	sort.Strings(attrs)

	joinedAttributes := strings.Join(attrs, ", ")

	// Build the final query
	sb.Grow(len(joinedAttributes) + 30)
	sb.WriteString("SELECT ")
	sb.WriteString(joinedAttributes)
	sb.WriteString(" FROM ")
	sb.WriteString(entity.ExternalId)

	if ordered {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(uniqueAttrExternalID)
		sb.WriteString(" ASC")
	}

	return sb.String(), nil
}

func AttrExternalIDToColumnName(externalID string) (string, error) {
	if strings.HasPrefix(externalID, "$.") {
		columnName, err := extractor.AttributesFromJSONPath(externalID)
		if err != nil {
			return "", err
		}

		if len(columnName) == 0 {
			return "", errors.New("JSON path extraction did not return any attributes")
		}

		return columnName[0], nil
	}

	return externalID, nil
}
