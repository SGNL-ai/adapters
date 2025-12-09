// Copyright 2025 SGNL.ai, Inc.
package commonutil

import "fmt"

// CreateChildEntitiesFromValues creates an array of child entity objects from a list of unique values.
// Each child entity will have:
// - id: formatted as "parentID_value"
// - value: the actual value
//
// Parameters:
// - parentID: The ID of the parent object (e.g., "003Hu000020yLuHIAU")
// - values: A list of unique values to transform into child entities
//
// Returns:
// - []any: An array of child entity objects suitable for the framework
//
// Example:
//
//	values := []string{"Sports", "Music", "Reading"}
//	childEntities := CreateChildEntitiesFromValues("123", values)
//	// Returns: []any{
//	//   map[string]any{"id": "123_Sports", "value": "Sports"},
//	//   map[string]any{"id": "123_Music", "value": "Music"},
//	//   map[string]any{"id": "123_Reading", "value": "Reading"}
//	// }
func CreateChildEntitiesFromValues(parentID string, values []string) []any {
	childObjects := make([]any, 0, len(values))

	for _, val := range values {
		childObjects = append(childObjects, map[string]any{
			"id":    fmt.Sprintf("%s_%s", parentID, val),
			"value": val,
		})
	}

	return childObjects
}
