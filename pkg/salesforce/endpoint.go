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
		// Extract the actual field name from JSONPath or use as-is
		fieldName := extractFieldName(attribute.ExternalId)

		if fieldName == "Id" {
			// Id is already added above
			continue
		}

		attributesBuilder.WriteRune(',')
		attributesBuilder.WriteString(url.QueryEscape(fieldName))
	}

	return attributesBuilder.String()
}

// extractFieldName extracts the field name from a JSON path or attribute name.
// Examples:
//   - $.CustomField__c → CustomField__c
//   - $.Account.Name → Account.Name
//   - $.Emails[*].value → Emails
//   - $.Contacts[?(@.primary==true)].Email → Contacts
//   - Name → Name (handles non-JSON path field names)
func extractFieldName(attributeName string) string {
	// Handle non-JSON path field names (like "Id", "Name", etc.)
	if !strings.HasPrefix(attributeName, "$.") {
		return removeArraySyntax(attributeName)
	}

	path := strings.TrimPrefix(attributeName, "$.")

	hasArraySyntax := strings.Contains(path, "[")

	path = removeArraySyntax(path)

	if hasArraySyntax {
		if dotIdx := strings.Index(path, "."); dotIdx != -1 {
			return path[:dotIdx]
		}
	}

	return path
}

// removeArraySyntax removes array access syntax from field names.
// Examples:
//   - Emails[*] → Emails
//   - Contacts[0] → Contacts
//   - Emails[?(@.primary==true)] → Emails
//   - CustomField__c → CustomField__c (unchanged)
func removeArraySyntax(fieldName string) string {
	if idx := strings.Index(fieldName, "["); idx != -1 {
		return fieldName[:idx]
	}

	return fieldName
}
