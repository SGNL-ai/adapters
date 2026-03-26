// Copyright 2026 SGNL.ai, Inc.

// Generic JSON:API relationship resolver for Rootly API responses.
// Replaces {id, type} stubs in relationships with full included objects.
package rootly

// processIncludes takes the data and included arrays from a Rootly API response,
// and merges the included resources into the main objects based on relationships.
func processIncludes(data []map[string]any, included []map[string]any) []map[string]any {
	if len(included) == 0 {
		return data
	}

	lookup := buildIncludedLookup(included)

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
				relData["data"] = resolveStubArray(relDataArr, lookup)
				relationships[relName] = relData

				continue
			}

			// Handle single object relationship.
			if relDataObj, ok := relDataField.(map[string]any); ok {
				if resolved := resolveStub(relDataObj, lookup); resolved != nil {
					relData["data"] = resolved
					relationships[relName] = relData
				}
			}
		}

		obj["relationships"] = relationships
		data[i] = obj
	}

	return data
}

// buildIncludedLookup indexes included objects by "type:id" for O(1) resolution.
func buildIncludedLookup(included []map[string]any) map[string]map[string]any {
	lookup := make(map[string]map[string]any, len(included))

	for _, inc := range included {
		id, ok := inc["id"].(string)
		if !ok {
			continue
		}

		typ, ok := inc["type"].(string)
		if !ok {
			continue
		}

		lookup[typ+":"+id] = inc
	}

	return lookup
}

// resolveStubArray resolves an array of relationship stubs against the lookup.
func resolveStubArray(items []any, lookup map[string]map[string]any) []any {
	resolved := make([]any, 0, len(items))

	for _, item := range items {
		itemMap, ok := item.(map[string]any)
		if !ok {
			resolved = append(resolved, item)

			continue
		}

		if merged := resolveStub(itemMap, lookup); merged != nil {
			resolved = append(resolved, merged)
		} else {
			resolved = append(resolved, itemMap)
		}
	}

	return resolved
}

// resolveStub looks up a single {id, type} stub in the lookup map and returns
// a merged map if found, or nil if the stub cannot be resolved.
func resolveStub(stub map[string]any, lookup map[string]map[string]any) map[string]any {
	id, hasID := stub["id"].(string)
	typ, hasType := stub["type"].(string)

	if !hasID || !hasType {
		return nil
	}

	includedObj, exists := lookup[typ+":"+id]
	if !exists {
		return nil
	}

	return mergeMaps(stub, includedObj)
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
