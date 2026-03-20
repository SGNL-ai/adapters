// Copyright 2026 SGNL.ai, Inc.

// ABOUTME: Response processing utilities for Rootly JSON:API responses.
// ABOUTME: Handles enrichment of incident data with included relationships.
package rootly

import (
	"fmt"
	"maps"
)

// IncludedItemProcessor processes included items and enriches incident data.
// It handles the extraction and flattening of nested relationships from Rootly's JSON:API format.
type IncludedItemProcessor struct {
	incidentID string
	included   []map[string]any
}

// selectedEntityTypes defines all the entity types that should be extracted from included items.
var selectedEntityTypes = []string{
	"selected_groups",
	"selected_options",
	"selected_services",
	"selected_functionalities",
	"selected_catalog_entities",
	"selected_users",
}

// NewIncludedItemProcessor creates a new processor for the given incident ID.
func NewIncludedItemProcessor(incidentID string, included []map[string]any) *IncludedItemProcessor {
	return &IncludedItemProcessor{
		incidentID: incidentID,
		included:   included,
	}
}

// FindMatching returns all included items that belong to this incident.
// It matches based on the incident_id field in the attributes.
func (p *IncludedItemProcessor) FindMatching() []map[string]any {
	var matching []map[string]any

	for _, item := range p.included {
		attrs, ok := item["attributes"].(map[string]any)
		if !ok {
			continue
		}

		if incidentID, ok := attrs["incident_id"].(string); ok && incidentID == p.incidentID {
			matching = append(matching, item)
		}
	}

	return matching
}

// ExtractEntities extracts and flattens entities of a specific type from all matching included items.
// Each entity will have a field_id property added for filtering.
// Returns a map where keys are entity types (e.g., "selected_users") and values are slices of
// flattened entities ([]any), to ensure compatibility with JSONPath filter expressions.
func (p *IncludedItemProcessor) ExtractEntities() map[string][]any {
	result := make(map[string][]any)

	for _, entityType := range selectedEntityTypes {
		entities := p.extractEntitiesWithFieldID(entityType)
		if len(entities) > 0 {
			result[entityType] = entities
		}
	}

	return result
}

// extractEntitiesWithFieldID is a helper that extracts entities of a given type
// and adds the form_field_id as field_id for easy JSONPath filtering.
// Returns []any to ensure compatibility with JSONPath filter expressions.
func (p *IncludedItemProcessor) extractEntitiesWithFieldID(entityKey string) []any {
	var result []any

	for _, item := range p.included {
		attrs, ok := item["attributes"].(map[string]any)
		if !ok {
			continue
		}

		// Only process items that belong to this incident
		incidentID, _ := attrs["incident_id"].(string)
		if incidentID != p.incidentID {
			continue
		}

		formFieldID, _ := attrs["form_field_id"].(string)
		entities := extractArrayField(attrs, entityKey)

		for _, entity := range entities {
			// Add form_field_id as field_id for easier JSONPath access
			if formFieldID != "" {
				entity["field_id"] = formFieldID
			}

			result = append(result, entity)
		}
	}

	return result
}

// ProcessAndExpand processes matching included items and expands nested arrays.
// It flattens all selected_* entity types into individual items with metadata.
// Returns []any to ensure compatibility with JSONPath filter expressions.
func (p *IncludedItemProcessor) ProcessAndExpand() []any {
	expanded := make([]any, 0, len(p.included))

	for _, item := range p.included {
		attrs, ok := item["attributes"].(map[string]any)
		if !ok {
			continue
		}

		// Only process items that belong to this incident
		incidentID, _ := attrs["incident_id"].(string)
		if incidentID != p.incidentID {
			continue
		}

		formFieldID, _ := attrs["form_field_id"].(string)

		// Process all selected entity types
		for _, entityKey := range selectedEntityTypes {
			entities := extractArrayField(attrs, entityKey)

			for _, entity := range entities {
				expandedItem := make(map[string]any, len(entity)+2)
				maps.Copy(expandedItem, entity)
				expandedItem["entity_type"] = entityKey

				if formFieldID != "" {
					expandedItem["form_field_id"] = formFieldID
				}

				expanded = append(expanded, expandedItem)
			}
		}

		// Also include the original item with form_field_id flattened
		modifiedItem := make(map[string]any, len(item))
		maps.Copy(modifiedItem, item)

		if formFieldID != "" {
			modifiedItem["form_field_id"] = formFieldID
		}

		expanded = append(expanded, modifiedItem)
	}

	return expanded
}

// EnrichIncidentData enriches a single incident data object with its included items.
func EnrichIncidentData(dataObject map[string]any, included []map[string]any) map[string]any {
	if len(included) == 0 {
		return dataObject
	}

	return enrichIncidentDataWithLookup(dataObject, included, buildIncludedLookup(included))
}

// enrichIncidentDataWithLookup is the internal implementation of EnrichIncidentData that accepts
// a pre-built included lookup to avoid redundant map construction when processing multiple incidents.
func enrichIncidentDataWithLookup(
	dataObject map[string]any, included []map[string]any, includedLookup map[string]map[string]any,
) map[string]any {
	// Only process included items if there are any
	if len(included) == 0 {
		return dataObject
	}

	// Get the incident ID
	incidentID, ok := dataObject["id"]
	if !ok {
		// If no id found, nothing to enrich
		return dataObject
	}

	incidentIDStr := fmt.Sprint(incidentID)

	// Create a copy of the original object
	enriched := make(map[string]any, len(dataObject))
	maps.Copy(enriched, dataObject)

	// Resolve relationship stubs against included objects.
	// This replaces bare {id, type} stubs in relationships with the full included objects,
	// enabling JSONPath traversal into nested attributes (e.g., role assignments).
	resolveRelationshipIncludes(enriched, includedLookup)

	// Create processor for this incident
	processor := NewIncludedItemProcessor(incidentIDStr, included)

	// Process and expand included items
	expandedItems := processor.ProcessAndExpand()
	if len(expandedItems) > 0 {
		enriched["included"] = expandedItems
	}

	// Extract all flattened entity types
	allEntities := processor.ExtractEntities()
	for entityType, entities := range allEntities {
		// Use "all_" prefix for flattened entity arrays (e.g., "all_selected_users")
		enriched["all_"+entityType] = entities
	}

	return enriched
}

// resolveRelationshipIncludes replaces relationship stubs with full included objects.
// It walks the data object's "relationships" map and for each relationship that has a
// "data" field (array or single object), replaces {id, type} stubs with the corresponding
// full object from the included lookup. This enables JSONPath traversal into nested attributes
// such as $.relationships.roles.data[*].attributes.user.data.attributes.email.
// The relationships map and each relationship wrapper map are copied to avoid mutating the original data object.
func resolveRelationshipIncludes(dataObject map[string]any, includedMap map[string]map[string]any) {
	originalRelationships, ok := dataObject["relationships"].(map[string]any)
	if !ok {
		return
	}

	if len(includedMap) == 0 {
		return
	}

	// Copy the relationships map and each wrapper map to avoid mutating the original data object.
	relationships := make(map[string]any, len(originalRelationships))

	for relName, relValue := range originalRelationships {
		relData, ok := relValue.(map[string]any)
		if !ok {
			relationships[relName] = relValue

			continue
		}

		// Copy the relationship wrapper map (contains "data" key).
		relDataCopy := make(map[string]any, len(relData))
		maps.Copy(relDataCopy, relData)
		relationships[relName] = relDataCopy
	}

	// Walk each relationship and resolve stubs.
	for relName, relValue := range relationships {
		relData, ok := relValue.(map[string]any)
		if !ok {
			continue
		}

		dataField, ok := relData["data"]
		if !ok {
			continue
		}

		// Handle array of relationship references.
		if dataArr, ok := dataField.([]any); ok {
			resolved := resolveRelationshipArray(dataArr, includedMap)
			relData["data"] = resolved
			relationships[relName] = relData

			continue
		}

		// Handle single relationship reference.
		if dataObj, ok := dataField.(map[string]any); ok {
			if resolved := resolveRelationshipObject(dataObj, includedMap); resolved != nil {
				relData["data"] = resolved
				relationships[relName] = relData
			}
		}
	}

	dataObject["relationships"] = relationships
}

// buildIncludedLookup builds a lookup map from included objects keyed by "type:id".
func buildIncludedLookup(included []map[string]any) map[string]map[string]any {
	lookup := make(map[string]map[string]any, len(included))

	for _, inc := range included {
		id, hasID := inc["id"].(string)
		typ, hasType := inc["type"].(string)

		if hasID && hasType {
			key := fmt.Sprintf("%s:%s", typ, id)
			lookup[key] = inc
		}
	}

	return lookup
}

// resolveRelationshipArray resolves an array of relationship stubs against the included lookup.
func resolveRelationshipArray(dataArr []any, includedMap map[string]map[string]any) []any {
	resolved := make([]any, 0, len(dataArr))

	for _, item := range dataArr {
		itemMap, ok := item.(map[string]any)
		if !ok {
			resolved = append(resolved, item)

			continue
		}

		if merged := resolveRelationshipObject(itemMap, includedMap); merged != nil {
			resolved = append(resolved, merged)
		} else {
			resolved = append(resolved, itemMap)
		}
	}

	return resolved
}

// resolveRelationshipObject resolves a single relationship stub against the included lookup.
// Returns a new merged map if a match is found, or nil if no match exists.
func resolveRelationshipObject(stub map[string]any, includedMap map[string]map[string]any) map[string]any {
	relID, hasID := stub["id"].(string)
	relType, hasType := stub["type"].(string)

	if !hasID || !hasType {
		return nil
	}

	key := fmt.Sprintf("%s:%s", relType, relID)

	includedObj, exists := includedMap[key]
	if !exists {
		return nil
	}

	// Merge: start with the stub, overlay with the full included object.
	merged := make(map[string]any, len(includedObj))
	maps.Copy(merged, stub)
	maps.Copy(merged, includedObj)

	return merged
}

// EnrichAllIncidentData enriches all incident data objects with their included items.
func EnrichAllIncidentData(data []map[string]any, included []map[string]any) []map[string]any {
	// Only process included items if there are any
	if len(included) == 0 {
		return data
	}

	// Build the included lookup once for all incidents to avoid redundant map construction.
	includedLookup := buildIncludedLookup(included)

	result := make([]map[string]any, len(data))

	for i, dataObject := range data {
		result[i] = enrichIncidentDataWithLookup(dataObject, included, includedLookup)
	}

	return result
}

// Helper functions

// extractArrayField extracts an array field from a map and converts it to []map[string]any.
// Handles both []any and []interface{} types automatically.
func extractArrayField(obj map[string]any, key string) []map[string]any {
	value, ok := obj[key]
	if !ok {
		return nil
	}

	return toMapSlice(value)
}

// toMapSlice converts various slice types to []map[string]any.
// Handles []any, nested conversions, and single map objects.
func toMapSlice(value any) []map[string]any {
	var result []map[string]any

	switch v := value.(type) {
	case []any:
		for _, item := range v {
			if m := toMap(item); m != nil {
				result = append(result, m)
			}
		}
	case []map[string]any:
		return v
	case map[string]any:
		// Handle single object case (e.g., selected_users as object with id/value)
		// Convert single object to a slice with one element
		return []map[string]any{v}
	}

	return result
}

// toMap converts an interface{} to map[string]any.
func toMap(value any) map[string]any {
	switch v := value.(type) {
	case map[string]any:
		return v
	default:
		return nil
	}
}
