// Copyright 2025 SGNL.ai, Inc.
package rootly

import (
	"reflect"
	"testing"
)

func TestIncludedItemProcessor_FindMatching(t *testing.T) {
	tests := []struct {
		name       string
		incidentID string
		included   []map[string]any
		want       int // Number of matching items
	}{
		{
			name:       "single matching item",
			incidentID: "incident-1",
			included: []map[string]any{
				{
					"id":   "item-1",
					"type": "custom_field_selection",
					"attributes": map[string]any{
						"incident_id":   "incident-1",
						"form_field_id": "field-1",
					},
				},
			},
			want: 1,
		},
		{
			name:       "multiple matching items",
			incidentID: "incident-1",
			included: []map[string]any{
				{
					"attributes": map[string]any{
						"incident_id":   "incident-1",
						"form_field_id": "field-1",
					},
				},
				{
					"attributes": map[string]any{
						"incident_id":   "incident-1",
						"form_field_id": "field-2",
					},
				},
				{
					"attributes": map[string]any{
						"incident_id": "incident-2",
					},
				},
			},
			want: 2,
		},
		{
			name:       "no matching items",
			incidentID: "incident-1",
			included: []map[string]any{
				{
					"attributes": map[string]any{
						"incident_id": "incident-2",
					},
				},
			},
			want: 0,
		},
		{
			name:       "empty included",
			incidentID: "incident-1",
			included:   []map[string]any{},
			want:       0,
		},
		{
			name:       "missing attributes",
			incidentID: "incident-1",
			included: []map[string]any{
				{
					"id":   "item-1",
					"type": "test",
				},
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			processor := NewIncludedItemProcessor(tt.incidentID, tt.included)

			// Act
			got := processor.FindMatching()

			// Assert
			if len(got) != tt.want {
				t.Errorf("FindMatching() returned %d items, want %d", len(got), tt.want)
			}
		})
	}
}

func TestIncludedItemProcessor_ExtractEntities(t *testing.T) {
	// Arrange
	const testIncidentID = "incident-1"

	included := []map[string]any{
		{
			"attributes": map[string]any{
				"incident_id":   testIncidentID,
				"form_field_id": "field-1",
				"selected_users": []any{
					map[string]any{"id": "user-1", "name": "Alice"},
					map[string]any{"id": "user-2", "name": "Bob"},
				},
				"selected_groups": []any{
					map[string]any{"id": "group-1", "name": "Engineering"},
				},
				"selected_services": []any{
					map[string]any{"id": "service-1", "name": "API Service"},
					map[string]any{"id": "service-2", "name": "Web Service"},
				},
				"selected_options": []any{
					map[string]any{"id": "option-1", "value": "high"},
				},
				"selected_functionalities": []any{
					map[string]any{"id": "func-1", "name": "Authentication"},
				},
				"selected_catalog_entities": []any{
					map[string]any{"id": "entity-1", "name": "User Database"},
				},
			},
		},
	}

	processor := NewIncludedItemProcessor(testIncidentID, included)

	// Act
	allEntities := processor.ExtractEntities()

	// Assert
	// Should have all 6 entity types
	expectedTypes := []string{
		"selected_users",
		"selected_groups",
		"selected_services",
		"selected_options",
		"selected_functionalities",
		"selected_catalog_entities",
	}

	if len(allEntities) != len(expectedTypes) {
		t.Errorf("ExtractEntities() returned %d entity types, want %d", len(allEntities), len(expectedTypes))
	}

	// Verify each entity type exists and has correct count
	expectations := map[string]int{
		"selected_users":            2,
		"selected_groups":           1,
		"selected_services":         2,
		"selected_options":          1,
		"selected_functionalities":  1,
		"selected_catalog_entities": 1,
	}

	for entityType, expectedCount := range expectations {
		entities, ok := allEntities[entityType]
		if !ok {
			t.Errorf("ExtractEntities() missing entity type %s", entityType)

			continue
		}

		if len(entities) != expectedCount {
			t.Errorf("ExtractEntities() %s has %d items, want %d", entityType, len(entities), expectedCount)
		}

		// Verify field_id was added
		for _, entityAny := range entities {
			entity, ok := entityAny.(map[string]any)
			if !ok {
				t.Errorf("ExtractEntities() %s entity is not a map", entityType)

				continue
			}
			if fieldID, ok := entity["field_id"].(string); !ok || fieldID != "field-1" {
				t.Errorf("ExtractEntities() %s entity missing correct field_id", entityType)
			}
		}
	}
}

func TestIncludedItemProcessor_ProcessAndExpand(t *testing.T) {
	// Arrange
	incidentID := "incident-1"
	included := []map[string]any{
		{
			"id":   "item-1",
			"type": "custom_field_selection",
			"attributes": map[string]any{
				"incident_id":   "incident-1",
				"form_field_id": "field-1",
				"selected_users": []any{
					map[string]any{"id": "user-1", "name": "Alice"},
				},
				"selected_groups": []any{
					map[string]any{"id": "group-1", "name": "Engineering"},
				},
			},
		},
	}

	processor := NewIncludedItemProcessor(incidentID, included)

	// Act
	expanded := processor.ProcessAndExpand()

	// Assert
	// Should have: 1 user + 1 group + 1 original item = 3 items
	if len(expanded) != 3 {
		t.Errorf("ProcessAndExpand() returned %d items, want 3", len(expanded))
	}

	// Verify entity_type is added
	foundUser := false
	foundGroup := false

	for _, item := range expanded {
		if entityType, ok := item["entity_type"].(string); ok {
			if entityType == "selected_users" {
				foundUser = true
			}
			if entityType == "selected_groups" {
				foundGroup = true
			}
		}
	}

	if !foundUser {
		t.Error("ProcessAndExpand() did not add entity_type for users")
	}
	if !foundGroup {
		t.Error("ProcessAndExpand() did not add entity_type for groups")
	}
}

func TestEnrichIncidentData(t *testing.T) {
	tests := []struct {
		name       string
		dataObject map[string]any
		included   []map[string]any
		checkFunc  func(t *testing.T, result map[string]any)
	}{
		{
			name: "enriches incident with included items",
			dataObject: map[string]any{
				"id":   "incident-1",
				"type": "incidents",
				"attributes": map[string]any{
					"title":  "Test Incident",
					"status": "open",
				},
			},
			included: []map[string]any{
				{
					"attributes": map[string]any{
						"incident_id":   "incident-1",
						"form_field_id": "field-1",
						"selected_users": []any{
							map[string]any{"id": "user-1", "name": "Alice"},
						},
					},
				},
			},
			checkFunc: func(t *testing.T, result map[string]any) {
				// Check that included field exists
				if _, ok := result["included"]; !ok {
					t.Error("EnrichIncidentData() missing 'included' field")
				}

				// Check that all_selected_users exists
				if users, ok := result["all_selected_users"].([]any); !ok {
					t.Error("EnrichIncidentData() missing 'all_selected_users' field")
				} else if len(users) != 1 {
					t.Errorf("EnrichIncidentData() all_selected_users has %d items, want 1", len(users))
				}

				// Check original fields are preserved
				if result["id"] != "incident-1" {
					t.Error("EnrichIncidentData() did not preserve original id")
				}
			},
		},
		{
			name: "enriches incident with all entity types",
			dataObject: map[string]any{
				"id":   "incident-1",
				"type": "incidents",
			},
			included: []map[string]any{
				{
					"attributes": map[string]any{
						"incident_id":               "incident-1",
						"form_field_id":             "field-1",
						"selected_users":            []any{map[string]any{"id": "user-1"}},
						"selected_groups":           []any{map[string]any{"id": "group-1"}},
						"selected_services":         []any{map[string]any{"id": "service-1"}},
						"selected_options":          []any{map[string]any{"id": "option-1"}},
						"selected_functionalities":  []any{map[string]any{"id": "func-1"}},
						"selected_catalog_entities": []any{map[string]any{"id": "entity-1"}},
					},
				},
			},
			checkFunc: func(t *testing.T, result map[string]any) {
				// Verify all entity types are enriched with "all_" prefix
				expectedFields := []string{
					"all_selected_users",
					"all_selected_groups",
					"all_selected_services",
					"all_selected_options",
					"all_selected_functionalities",
					"all_selected_catalog_entities",
				}

				for _, field := range expectedFields {
					if entities, ok := result[field].([]any); !ok {
						t.Errorf("EnrichIncidentData() missing '%s' field", field)
					} else if len(entities) != 1 {
						t.Errorf("EnrichIncidentData() %s has %d items, want 1", field, len(entities))
					}
				}
			},
		},
		{
			name: "handles incident with no ID",
			dataObject: map[string]any{
				"type": "incidents",
			},
			included: []map[string]any{},
			checkFunc: func(t *testing.T, result map[string]any) {
				// When included is empty, should not add included field
				if _, ok := result["included"]; ok {
					t.Error("EnrichIncidentData() should not add included field when empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange is in test table (tt.dataObject, tt.included)

			// Act
			result := EnrichIncidentData(tt.dataObject, tt.included)

			// Assert
			tt.checkFunc(t, result)
		})
	}
}

func TestEnrichAllIncidentData(t *testing.T) {
	// Arrange
	data := []map[string]any{
		{
			"id":   "incident-1",
			"type": "incidents",
		},
		{
			"id":   "incident-2",
			"type": "incidents",
		},
	}

	included := []map[string]any{
		{
			"attributes": map[string]any{
				"incident_id": "incident-1",
			},
		},
	}

	// Act
	result := EnrichAllIncidentData(data, included)

	// Assert
	if len(result) != 2 {
		t.Errorf("EnrichAllIncidentData() returned %d items, want 2", len(result))
	}

	// incident-1 should have included field (has matching items), incident-2 should not
	if _, ok := result[0]["included"]; !ok {
		t.Errorf("EnrichAllIncidentData() item 0 (incident-1) should have 'included' field")
	}

	// incident-2 should not have included field (no matching items)
	if _, ok := result[1]["included"]; ok {
		t.Errorf("EnrichAllIncidentData() item 1 (incident-2) should not have 'included' field when no matching items")
	}
}

func TestToMapSlice(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  int // Number of maps in result
	}{
		{
			name: "[]any with maps",
			input: []any{
				map[string]any{"id": "1"},
				map[string]any{"id": "2"},
			},
			want: 2,
		},
		{
			name: "[]interface{} with maps",
			input: []interface{}{
				map[string]interface{}{"id": "1"},
				map[string]interface{}{"id": "2"},
			},
			want: 2,
		},
		{
			name:  "[]map[string]any",
			input: []map[string]any{{"id": "1"}, {"id": "2"}},
			want:  2,
		},
		{
			name:  "empty slice",
			input: []any{},
			want:  0,
		},
		{
			name:  "nil value",
			input: nil,
			want:  0,
		},
		{
			name: "mixed types - filters non-maps",
			input: []any{
				map[string]any{"id": "1"},
				"not a map",
				map[string]any{"id": "2"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange is in test table (tt.input)

			// Act
			got := toMapSlice(tt.input)

			// Assert
			if len(got) != tt.want {
				t.Errorf("toMapSlice() returned %d items, want %d", len(got), tt.want)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantNil bool
	}{
		{
			name:    "map[string]any",
			input:   map[string]any{"id": "1"},
			wantNil: false,
		},
		{
			name:    "map[string]interface{}",
			input:   map[string]interface{}{"id": "1"},
			wantNil: false,
		},
		{
			name:    "string - not a map",
			input:   "not a map",
			wantNil: true,
		},
		{
			name:    "nil",
			input:   nil,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange is in test table (tt.input)

			// Act
			got := toMap(tt.input)

			// Assert
			if (got == nil) != tt.wantNil {
				t.Errorf("toMap() returned nil=%v, want nil=%v", got == nil, tt.wantNil)
			}
		})
	}
}

func TestCopyMap(t *testing.T) {
	// Arrange
	original := map[string]any{
		"id":   "1",
		"name": "Test",
		"nested": map[string]any{
			"value": 42,
		},
	}

	// Act
	copied := copyMap(original)

	// Assert
	// Check that all keys are copied
	if !reflect.DeepEqual(original, copied) {
		t.Error("copyMap() did not create equal map")
	}

	// Check that it's a different map (not same reference)
	copied["id"] = "2"
	if original["id"] == "2" {
		t.Error("copyMap() did not create a new map (shallow copy issue)")
	}

	// Note: This is a shallow copy, so nested maps are still shared
	// This is intentional for performance reasons
}

func TestGetAttributes(t *testing.T) {
	tests := []struct {
		name    string
		item    map[string]any
		wantOk  bool
		wantLen int
	}{
		{
			name: "valid attributes",
			item: map[string]any{
				"attributes": map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
			},
			wantOk:  true,
			wantLen: 2,
		},
		{
			name:    "missing attributes",
			item:    map[string]any{"id": "1"},
			wantOk:  false,
			wantLen: 0,
		},
		{
			name: "attributes is not a map",
			item: map[string]any{
				"attributes": "not a map",
			},
			wantOk:  false,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange is in test table (tt.item)

			// Act
			got, ok := getAttributes(tt.item)

			// Assert
			if ok != tt.wantOk {
				t.Errorf("getAttributes() ok = %v, want %v", ok, tt.wantOk)
			}

			if ok && len(got) != tt.wantLen {
				t.Errorf("getAttributes() returned map with %d items, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestExtractArrayField(t *testing.T) {
	tests := []struct {
		name string
		obj  map[string]any
		key  string
		want int
	}{
		{
			name: "extract existing array",
			obj: map[string]any{
				"users": []any{
					map[string]any{"id": "1"},
					map[string]any{"id": "2"},
				},
			},
			key:  "users",
			want: 2,
		},
		{
			name: "missing key",
			obj: map[string]any{
				"other": "value",
			},
			key:  "users",
			want: 0,
		},
		{
			name: "empty array",
			obj: map[string]any{
				"users": []any{},
			},
			key:  "users",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange is in test table (tt.obj, tt.key)

			// Act
			got := extractArrayField(tt.obj, tt.key)

			// Assert
			if len(got) != tt.want {
				t.Errorf("extractArrayField() returned %d items, want %d", len(got), tt.want)
			}
		})
	}
}
