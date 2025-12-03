// Copyright 2025 SGNL.ai, Inc.

// ABOUTME: Response processing utilities for Rootly JSON:API responses.
// ABOUTME: Handles enrichment of incident data with included relationships.
package rootly

import "fmt"

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
		attrs, ok := getAttributes(item)
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
// Returns a map where keys are entity types (e.g., "selected_users") and values are the flattened entities.
func (p *IncludedItemProcessor) ExtractEntities() map[string][]map[string]any {
	result := make(map[string][]map[string]any)

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
func (p *IncludedItemProcessor) extractEntitiesWithFieldID(entityKey string) []map[string]any {
	var result []map[string]any

	for _, item := range p.included {
		attrs, ok := getAttributes(item)
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
func (p *IncludedItemProcessor) ProcessAndExpand() []map[string]any {
	expanded := make([]map[string]any, 0, len(p.included))

	for _, item := range p.included {
		attrs, ok := getAttributes(item)
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
				expandedItem := copyMap(entity)
				expandedItem["entity_type"] = entityKey

				if formFieldID != "" {
					expandedItem["form_field_id"] = formFieldID
				}

				expanded = append(expanded, expandedItem)
			}
		}

		// Also include the original item with form_field_id flattened
		modifiedItem := copyMap(item)
		if formFieldID != "" {
			modifiedItem["form_field_id"] = formFieldID
		}

		expanded = append(expanded, modifiedItem)
	}

	return expanded
}

// EnrichIncidentData enriches a single incident data object with its included items.
func EnrichIncidentData(dataObject map[string]any, included []map[string]any) map[string]any {
	// Create a copy of the original object
	enriched := copyMap(dataObject)

	// Only process included items if there are any
	if len(included) == 0 {
		return enriched
	}

	// Get the incident ID
	incidentID, ok := dataObject["id"]
	if !ok {
		// If no id found, nothing to enrich
		return enriched
	}

	incidentIDStr := fmt.Sprint(incidentID)

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

// EnrichAllIncidentData enriches all incident data objects with their included items.
func EnrichAllIncidentData(data []map[string]any, included []map[string]any) []map[string]any {
	result := make([]map[string]any, len(data))

	for i, dataObject := range data {
		result[i] = EnrichIncidentData(dataObject, included)
	}

	return result
}

// Helper functions

// getAttributes safely extracts the attributes field from a map.
func getAttributes(item map[string]any) (map[string]any, bool) {
	attrs, ok := item["attributes"].(map[string]any)

	return attrs, ok
}

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
// Handles []any and nested conversions.
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

// copyMap creates a shallow copy of a map[string]any.
func copyMap(original map[string]any) map[string]any {
	copied := make(map[string]any, len(original))
	for key, value := range original {
		copied[key] = value
	}

	return copied
}
