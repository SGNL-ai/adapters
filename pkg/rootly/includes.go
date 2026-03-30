// Copyright 2026 SGNL.ai, Inc.

// Generic JSON:API relationship resolver for Rootly API responses.
// Replaces {id, type} stubs in relationships with full included objects.
package rootly

import "maps"

// resolveIncludedRelationships replaces JSON:API relationship stubs with their
// full objects from the included array.
//
// A Rootly API response follows the JSON:API spec where relationships contain
// bare {id, type} stubs, and the full objects live in a separate "included" array:
//
//	Before: "roles": {"data": [{"id": "123", "type": "role"}]}
//	After:  "roles": {"data": [{"id": "123", "type": "role", "attributes": {"name": "Commander"}}]}
func resolveIncludedRelationships(data []map[string]any, included []map[string]any) []map[string]any {
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

			switch v := relDataField.(type) {
			case []any:
				relData["data"] = resolveStubArray(v, lookup)
				relationships[relName] = relData

			case map[string]any:
				if resolved := resolveStub(v, lookup); resolved != nil {
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

	maps.Copy(result, m1)
	maps.Copy(result, m2)

	return result
}
