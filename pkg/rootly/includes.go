// Copyright 2026 SGNL.ai, Inc.

// Generic JSON:API relationship resolver for Rootly API responses.
// Replaces {id, type} stubs in relationships with full included objects.
package rootly

import (
	"fmt"
)

// processIncludes takes the data and included arrays from a Rootly API response,
// and merges the included resources into the main objects based on relationships.
func processIncludes(data []map[string]any, included []map[string]any) []map[string]any {
	if len(included) == 0 {
		return data
	}

	// Build a lookup map for included objects by "type:id".
	includedMap := make(map[string]map[string]any, len(included))

	for _, inc := range included {
		id, ok := inc["id"].(string)
		if !ok {
			continue
		}

		typ, ok := inc["type"].(string)
		if !ok {
			continue
		}

		key := fmt.Sprintf("%s:%s", typ, id)
		includedMap[key] = inc
	}

	// Process each data object.
	for i, obj := range data {
		relationships, ok := obj["relationships"].(map[string]any)
		if !ok {
			continue
		}

		// Walk each relationship and resolve stubs.
		for relName, relValue := range relationships {
			relData, ok := relValue.(map[string]any)
			if !ok {
				continue
			}

			relDataField, ok := relData["data"]
			if !ok {
				continue
			}

			// Handle array of relationship references.
			if relDataArr, ok := relDataField.([]any); ok {
				mergedItems := make([]any, 0, len(relDataArr))

				for _, relItem := range relDataArr {
					relItemMap, ok := relItem.(map[string]any)
					if !ok {
						mergedItems = append(mergedItems, relItem)

						continue
					}

					relID, hasID := relItemMap["id"].(string)
					relType, hasType := relItemMap["type"].(string)

					if hasID && hasType {
						key := fmt.Sprintf("%s:%s", relType, relID)
						if includedObj, exists := includedMap[key]; exists {
							mergedItems = append(mergedItems, mergeMaps(relItemMap, includedObj))
						} else {
							mergedItems = append(mergedItems, relItemMap)
						}
					} else {
						mergedItems = append(mergedItems, relItemMap)
					}
				}

				relData["data"] = mergedItems
				relationships[relName] = relData

				continue
			}

			// Handle single object relationship.
			if relDataObj, ok := relDataField.(map[string]any); ok {
				relID, hasID := relDataObj["id"].(string)
				relType, hasType := relDataObj["type"].(string)

				if hasID && hasType {
					key := fmt.Sprintf("%s:%s", relType, relID)
					if includedObj, exists := includedMap[key]; exists {
						relData["data"] = mergeMaps(relDataObj, includedObj)
						relationships[relName] = relData
					}
				}
			}
		}

		obj["relationships"] = relationships
		data[i] = obj
	}

	return data
}

// mergeMaps creates a new map that contains all key/values from both maps,
// with the second map's values taking precedence in case of conflicts.
func mergeMaps(m1, m2 map[string]any) map[string]any {
	result := make(map[string]any, len(m1)+len(m2))

	for k, v := range m1 {
		result[k] = v
	}

	for k, v := range m2 {
		result[k] = v
	}

	return result
}
