// Copyright 2025 SGNL.ai, Inc.
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

	// Build a lookup map for included objects by ID and type
	includedMap := make(map[string]map[string]any)
	for _, inc := range included {
		id, ok := inc["id"].(string)
		if !ok {
			continue
		}

		typ, ok := inc["type"].(string)
		if !ok {
			continue
		}

		// Use "type:id" as the key for the map
		key := fmt.Sprintf("%s:%s", typ, id)
		includedMap[key] = inc
	}

	// Process each data object
	for i, obj := range data {
		// Look for relationships field in the object
		relationships, ok := obj["relationships"].(map[string]any)
		if !ok {
			continue
		}

		// Go through each relationship
		for relName, relValue := range relationships {
			relData, ok := relValue.(map[string]any)
			if !ok {
				continue
			}

			// Check if the relationship has a data field (either single object or array)
			relDataField, ok := relData["data"]
			if !ok {
				continue
			}

			// Handle array of relationships
			if relDataArr, ok := relDataField.([]any); ok {
				mergedItems := make([]any, 0, len(relDataArr))

				for _, relItem := range relDataArr {
					if relItemMap, ok := relItem.(map[string]any); ok {
						relID, hasID := relItemMap["id"].(string)
						relType, hasType := relItemMap["type"].(string)

						if hasID && hasType {
							// Look up the included object
							key := fmt.Sprintf("%s:%s", relType, relID)
							if includedObj, exists := includedMap[key]; exists {
								// Merge the relationship object with the included object
								mergedItem := mergeMaps(relItemMap, includedObj)
								mergedItems = append(mergedItems, mergedItem)
							} else {
								// Keep the original relationship reference if not found
								mergedItems = append(mergedItems, relItemMap)
							}
						} else {
							// Keep the original relationship if it doesn't have id/type
							mergedItems = append(mergedItems, relItemMap)
						}
					}
				}

				// Replace the relationship data with the expanded version
				relData["data"] = mergedItems
				relationships[relName] = relData
			} else if relDataObj, ok := relDataField.(map[string]any); ok {
				// Handle single object relationship
				relID, hasID := relDataObj["id"].(string)
				relType, hasType := relDataObj["type"].(string)

				if hasID && hasType {
					// Look up the included object
					key := fmt.Sprintf("%s:%s", relType, relID)
					if includedObj, exists := includedMap[key]; exists {
						// Merge the relationship object with the included object
						mergedItem := mergeMaps(relDataObj, includedObj)
						relData["data"] = mergedItem
						relationships[relName] = relData
					}
				}
			}
		}

		// Update the relationships in the object
		obj["relationships"] = relationships
		data[i] = obj
	}

	return data
}

// mergeMaps creates a new map that contains all key/values from both maps,
// with the second map's values taking precedence in case of conflicts.
func mergeMaps(m1, m2 map[string]any) map[string]any {
	result := make(map[string]any)

	// Copy m1 into result
	for k, v := range m1 {
		result[k] = v
	}

	// Copy m2 into result, potentially overwriting values from m1
	for k, v := range m2 {
		result[k] = v
	}

	return result
}
