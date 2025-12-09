// Copyright 2025 SGNL.ai, Inc.
package commonutil

import (
	"fmt"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
)

// GetUniqueIDValue extracts the unique ID value from an object based on the entity configuration.
// Returns the unique ID value and true if found, empty string and false otherwise.
//
// Parameters:
// - obj: The object to extract the unique ID from
// - entityConfig: The entity configuration containing attribute metadata
//
// Returns:
// - string: The unique ID value
// - bool: Whether the unique ID was found and is a string
func GetUniqueIDValue(obj map[string]any, entityConfig *framework.EntityConfig) (string, bool) {
	// Find the unique ID attribute in the entity config
	var uniqueIDExternalID string
	for _, attr := range entityConfig.Attributes {
		if attr.UniqueId {
			uniqueIDExternalID = attr.ExternalId
			break
		}
	}

	// If no unique ID attribute is configured, return false
	if uniqueIDExternalID == "" {
		return "", false
	}

	// Extract the value from the object
	value, exists := obj[uniqueIDExternalID]
	if !exists || value == nil {
		return "", false
	}

	// Ensure it's a string
	strValue, ok := value.(string)
	if !ok || strValue == "" {
		return "", false
	}

	return strValue, true
}

// CreateChildEntitiesFromList creates an array of child entity objects from a list of values.
// The function automatically:
// - De-duplicates values (case-insensitive)
// - Trims whitespace
// - Filters out empty values
// - Generates deterministic IDs
// - Uses attribute names from the child entity config
//
// Parameters:
// - parentID: The ID of the parent object (e.g., "003Hu000020yLuHIAU")
// - fieldName: The name of the field (e.g., "Interests__c")
// - values: A list of values to transform into child entities (may contain duplicates)
// - childEntityConfig: The child entity configuration defining the attributes
//
// Returns:
// - []any: An array of child entity objects suitable for the framework
//
// Expected child entity config attributes:
// - One attribute with ExternalId "id" for the unique identifier
// - One attribute with ExternalId "value" for the actual value
//
// Example:
//
//	values := []string{"Sports", "Music", "Sports", " Reading "}
//	config := &framework.EntityConfig{
//	  Attributes: []*framework.AttributeConfig{
//	    {ExternalId: "id", Type: framework.AttributeTypeString},
//	    {ExternalId: "value", Type: framework.AttributeTypeString},
//	  },
//	}
//	childEntities := CreateChildEntitiesFromList("123", "Interests__c", values, config)
//	// Returns: []any{
//	//   map[string]any{"id": "123_Interests__c_sports", "value": "Sports"},
//	//   map[string]any{"id": "123_Interests__c_music", "value": "Music"},
//	//   map[string]any{"id": "123_Interests__c_reading", "value": "Reading"}
//	// }
func CreateChildEntitiesFromList(
	parentID string,
	fieldName string,
	values []string,
	childEntityConfig *framework.EntityConfig,
) []any {
	// Extract attribute names from config
	var idAttr, valueAttr string
	for _, attr := range childEntityConfig.Attributes {
		switch attr.ExternalId {
		case "id":
			idAttr = attr.ExternalId
		case "value":
			valueAttr = attr.ExternalId
		}
	}

	// De-duplicate values using a map (case-insensitive key, original value)
	uniqueValues := make(map[string]string)

	for _, val := range values {
		trimmedVal := strings.TrimSpace(val)
		if trimmedVal == "" {
			continue
		}

		// Use lowercase as key for case-insensitive deduplication
		// But store the first occurrence's original casing
		lowerKey := strings.ToLower(trimmedVal)
		if _, exists := uniqueValues[lowerKey]; !exists {
			uniqueValues[lowerKey] = trimmedVal
		}
	}

	// Create child objects with deterministic IDs
	childObjects := make([]any, 0, len(uniqueValues))

	for lowerKey, originalValue := range uniqueValues {
		childObj := make(map[string]any)

		if idAttr != "" {
			childObj[idAttr] = generateCompositeID(parentID, fieldName, lowerKey)
		}

		if valueAttr != "" {
			childObj[valueAttr] = originalValue
		}

		childObjects = append(childObjects, childObj)
	}

	return childObjects
}

// CreateChildEntitiesFromDelimitedString transforms delimited string fields into child entity arrays.
// This is a generic function that works for any adapter with delimited values.
//
// The function automatically extracts the parent unique ID from the entity configuration,
// so adapters don't need to manually find the unique ID attribute.
//
// Parameters:
// - objects: The raw objects from the datasource
// - parentEntityConfig: The parent entity configuration (used to find unique ID)
// - childEntities: Configuration of which fields should be transformed
// - delimiter: The delimiter used to separate values (e.g., ";", ",", "|")
//
// Returns:
// - []map[string]any: Transformed objects with child entity arrays
//
// Example:
//
//	objects := []map[string]any{
//	  {"Id": "123", "Interests": "Sports;Music;Reading"},
//	}
//	parentConfig := &framework.EntityConfig{
//	  Attributes: []*framework.AttributeConfig{
//	    {ExternalId: "Id", UniqueId: true},
//	  },
//	}
//	childEntities := []*framework.EntityConfig{
//	  {
//	    ExternalId: "Interests",
//	    Attributes: []*framework.AttributeConfig{
//	      {ExternalId: "id"}, {ExternalId: "value"},
//	    },
//	  },
//	}
//	transformed := CreateChildEntitiesFromDelimitedString(objects, parentConfig, childEntities, ";")
func CreateChildEntitiesFromDelimitedString(
	objects []map[string]any,
	parentEntityConfig *framework.EntityConfig,
	childEntities []*framework.EntityConfig,
	delimiter string,
) []map[string]any {
	if len(childEntities) == 0 {
		return objects
	}

	// Map of field name -> child entity config for O(1) lookup
	childEntityFields := make(map[string]*framework.EntityConfig, len(childEntities))
	for _, childEntity := range childEntities {
		childEntityFields[childEntity.ExternalId] = childEntity
	}

	transformedObjects := make([]map[string]any, len(objects))

	for i, obj := range objects {
		transformedObj := make(map[string]any, len(obj))

		// Copy all existing fields
		for key, value := range obj {
			transformedObj[key] = value
		}

		// Extract parent unique ID using the entity config
		parentID, ok := GetUniqueIDValue(obj, parentEntityConfig)
		if !ok {
			// If we can't get the parent ID, keep the object as-is
			transformedObjects[i] = transformedObj
			continue
		}

		// Transform each child entity field
		for fieldName, childConfig := range childEntityFields {
			value, exists := obj[fieldName]

			// If field doesn't exist in the original object, don't add it
			if !exists {
				continue
			}

			// Field exists but is nil - set to empty array for consistency
			if value == nil {
				transformedObj[fieldName] = []any{}

				continue
			}

			// Only transform string values (delimited lists)
			strValue, ok := value.(string)
			if !ok {
				// Field exists but isn't a string - keep the original value (already copied)
				continue
			}

			// Handle empty strings - set to empty array
			if strValue == "" {
				transformedObj[fieldName] = []any{}

				continue
			}

			// Split by delimiter and convert to child entities
			values := strings.Split(strValue, delimiter)
			transformedObj[fieldName] = CreateChildEntitiesFromList(parentID, fieldName, values, childConfig)
		}

		transformedObjects[i] = transformedObj
	}

	return transformedObjects
}

// generateCompositeID generates a deterministic, human-readable composite ID.
// Format: {parentID}_{fieldName}_{cleanValue}
// Example: "003Hu000020yLuHIAU_Interests__c_sports"
func generateCompositeID(parentID string, fieldName string, value string) string {
	// Clean value to make it URL-safe and consistent
	// Value should already be lowercase from the caller
	cleanValue := strings.ReplaceAll(value, " ", "-")
	cleanValue = strings.ReplaceAll(cleanValue, "/", "-")
	cleanValue = strings.ReplaceAll(cleanValue, "\\", "-")

	return fmt.Sprintf("%s_%s_%s", parentID, fieldName, cleanValue)
}
