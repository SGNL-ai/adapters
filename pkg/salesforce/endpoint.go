// Copyright 2025 SGNL.ai, Inc.
package salesforce

import (
	"net/url"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
)

// ConstructEndpoint constructs and returns the endpoint to query the datasource.
func ConstructEndpoint(request *Request) string {
	if request == nil {
		return ""
	}

	var sb strings.Builder

	// URL Format:
	// For the first page: baseURL + "/services/data/v" + apiVersion + "/query?q=SELECT+Id" + encodedAttributes
	// 		+ "+FROM+" + entityExternalID + ["+WHERE+" + filter] + "+ORDER+BY+Id+ASC"
	// For subsequent requests: baseURL + cursor
	sb.WriteString(request.BaseURL)

	if request.Cursor != nil {
		sb.WriteString(*request.Cursor)

		return sb.String()
	}

	escapedAPIVersion := url.QueryEscape(request.APIVersion)
	escapedEntityExternalID := url.QueryEscape(request.EntityExternalID)
	encodedAttributes := encodedAttributes(request.Attributes)

	sb.Grow(len(escapedAPIVersion) + len(encodedAttributes) + len(escapedEntityExternalID) + 63)

	sb.WriteString("/services/data/v")
	sb.WriteString(escapedAPIVersion)
	sb.WriteString("/query?q=SELECT+Id")
	sb.WriteString(encodedAttributes)
	sb.WriteString("+FROM+")
	sb.WriteString(escapedEntityExternalID)

	if request.Filter != nil {
		sb.WriteString("+WHERE+")
		sb.WriteString(url.QueryEscape(*request.Filter))
	}

	sb.WriteString("+ORDER+BY+Id+ASC")

	return sb.String()
}

func encodedAttributes(attributes []*framework.AttributeConfig) string {
	var attributesBuilder strings.Builder
	// Guesstimating initial buffer need, len(attributes) * 6 byte strings
	attributesBuilder.Grow(len(attributes) * 6)

	for _, attribute := range attributes {
		if attribute.ExternalId == "Id" {
			// Id is already added above
			continue
		}

		attributesBuilder.WriteRune(',')
		attributesBuilder.WriteString(url.QueryEscape(attribute.ExternalId))
	}

	return attributesBuilder.String()
}
