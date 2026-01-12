// Copyright 2025 SGNL.ai, Inc.

package rootly

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

// TestResponseProcessingIntegration tests the full response processing pipeline
// using real Rootly response structure with all entity types.
func TestResponseProcessingIntegration(t *testing.T) {
	// Arrange
	// Load test data from file
	data, err := os.ReadFile("testdata/incident_response_with_included.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	// Parse the response
	var response DatasourceResponse
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("Failed to unmarshal test data: %v", err)
	}

	// Act
	// Process the data using our enrichment function
	processedData := EnrichAllIncidentData(response.Data, response.Included)

	// Assert
	// Verify we have the expected number of incidents
	if len(processedData) != 1 {
		t.Errorf("Expected 1 processed incident, got %d", len(processedData))
	}

	incident := processedData[0]

	// Verify original fields are preserved
	if incident["id"] != "test-incident-1" {
		t.Errorf("Expected incident id 'test-incident-1', got %v", incident["id"])
	}

	if incident["type"] != "incidents" {
		t.Errorf("Expected incident type 'incidents', got %v", incident["type"])
	}

	// Verify all_selected_* fields are added for all 6 entity types
	expectedEntityFields := map[string]int{
		"all_selected_users":            2, // 2 users
		"all_selected_groups":           2, // 2 groups
		"all_selected_services":         2, // 2 services
		"all_selected_options":          1, // 1 option
		"all_selected_functionalities":  2, // 2 functionalities
		"all_selected_catalog_entities": 1, // 1 catalog entity
	}

	for fieldName, expectedCount := range expectedEntityFields {
		entitiesAny, ok := incident[fieldName].([]any)
		if !ok {
			t.Errorf("Expected incident to have '%s' field as []any", fieldName)

			continue
		}

		if len(entitiesAny) != expectedCount {
			t.Errorf("Expected %s to have %d items, got %d", fieldName, expectedCount, len(entitiesAny))
		}

		// Verify each entity has field_id added
		for i, entityAny := range entitiesAny {
			entity, ok := entityAny.(map[string]any)
			if !ok {
				t.Errorf("%s[%d] is not a map", fieldName, i)

				continue
			}
			if fieldID, ok := entity["field_id"].(string); !ok || fieldID == "" {
				t.Errorf("%s[%d] missing field_id", fieldName, i)
			}
		}
	}

	// Verify included field exists and contains expanded items
	includedItemsAny, ok := incident["included"].([]any)
	if !ok {
		t.Fatal("Expected incident to have 'included' field as []any")
	}

	if len(includedItemsAny) == 0 {
		t.Error("Expected included items to be non-empty")
	}

	// Count expanded items by entity_type
	entityTypeCounts := make(map[string]int)
	for _, itemAny := range includedItemsAny {
		item, ok := itemAny.(map[string]any)
		if !ok {
			continue
		}

		if entityType, ok := item["entity_type"].(string); ok {
			entityTypeCounts[entityType]++
		}
	}

	// Verify we have expanded items for each entity type
	expectedEntityTypes := []string{
		"selected_users",
		"selected_groups",
		"selected_services",
		"selected_options",
		"selected_functionalities",
		"selected_catalog_entities",
	}

	for _, entityType := range expectedEntityTypes {
		count := entityTypeCounts[entityType]
		if count == 0 {
			t.Errorf("Expected to find expanded items with entity_type '%s', but found none", entityType)
		}
	}

	// Now test JSONPath queries on the enriched data
	t.Run("JSONPath queries", func(t *testing.T) {
		// Convert processed data to JSON for JSONPath queries
		processedJSON, err := json.Marshal(processedData)
		if err != nil {
			t.Fatalf("Failed to marshal processed data: %v", err)
		}

		// Parse JSON into interface{} for ojg/jp
		var jsonData any
		if err := json.Unmarshal(processedJSON, &jsonData); err != nil {
			t.Fatalf("Failed to unmarshal for JSONPath: %v", err)
		}

		// Test 1: Filter users by field_id using JSONPath filter expression
		t.Run("filter users by field_id", func(t *testing.T) {
			// JSONPath: $[0].all_selected_users[?(@.field_id=='field-users')]
			x, err := jp.ParseString("$[0].all_selected_users[?(@.field_id=='field-users')]")
			if err != nil {
				t.Fatalf("Failed to parse JSONPath: %v", err)
			}

			result := x.Get(jsonData)
			if len(result) == 0 {
				t.Error("Expected to find users with field_id=='field-users'")

				return
			}

			// Should find 2 users
			if len(result) != 2 {
				t.Errorf("Expected 2 users with field_id=='field-users', got %d", len(result))
			}

			// Check first user
			firstUser, ok := result[0].(map[string]any)
			if !ok {
				t.Error("Expected first result to be a map")

				return
			}

			if firstUser["id"] != "user-10" {
				t.Errorf("Expected user id 'user-10', got %v", firstUser["id"])
			}
		})

		// Test 2: Filter services by field_id
		t.Run("filter services by field_id", func(t *testing.T) {
			// JSONPath: $[0].all_selected_services[?(@.field_id=='field-services')]
			x, err := jp.ParseString("$[0].all_selected_services[?(@.field_id=='field-services')]")
			if err != nil {
				t.Fatalf("Failed to parse JSONPath: %v", err)
			}

			result := x.Get(jsonData)
			if len(result) == 0 {
				t.Error("Expected to find services with field_id=='field-services'")

				return
			}

			// Should find 2 services
			if len(result) != 2 {
				t.Errorf("Expected 2 services with field_id=='field-services', got %d", len(result))
			}

			firstService, ok := result[0].(map[string]any)
			if !ok {
				t.Error("Expected first result to be a map")

				return
			}

			if firstService["name"] != "API Service" {
				t.Errorf("Expected service name 'API Service', got %v", firstService["name"])
			}
		})

		// Test 3: Count entities using JSONPath
		t.Run("count entities", func(t *testing.T) {
			tests := []struct {
				name     string
				path     string
				expected int
			}{
				{"count users", "$[0].all_selected_users", 2},
				{"count groups", "$[0].all_selected_groups", 2},
				{"count services", "$[0].all_selected_services", 2},
				{"count options", "$[0].all_selected_options", 1},
				{"count functionalities", "$[0].all_selected_functionalities", 2},
				{"count catalog entities", "$[0].all_selected_catalog_entities", 1},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					x, err := jp.ParseString(tt.path)
					if err != nil {
						t.Fatalf("Failed to parse JSONPath: %v", err)
					}

					result := x.Get(jsonData)
					if len(result) == 0 {
						t.Errorf("%s: expected to find array", tt.name)

						return
					}

					arr, ok := result[0].([]any)
					if !ok {
						t.Errorf("%s: expected array", tt.name)

						return
					}

					if len(arr) != tt.expected {
						t.Errorf("%s: expected %d, got %d", tt.name, tt.expected, len(arr))
					}
				})
			}
		})

		// Test 4: Filter included by entity_type
		t.Run("filter included by entity_type", func(t *testing.T) {
			// JSONPath: $[0].included[?(@.entity_type=='selected_users')]
			x, err := jp.ParseString("$[0].included[?(@.entity_type=='selected_users')]")
			if err != nil {
				t.Fatalf("Failed to parse JSONPath: %v", err)
			}

			result := x.Get(jsonData)
			if len(result) == 0 {
				t.Error("Expected to find included items with entity_type=='selected_users'")

				return
			}

			// Should find 2 user items
			if len(result) < 1 {
				t.Errorf("Expected at least 1 user in included, got %d", len(result))
			}
		})

		// Test 5: Query nested attributes
		t.Run("query nested attributes", func(t *testing.T) {
			x, err := jp.ParseString("$[0].attributes.title")
			if err != nil {
				t.Fatalf("Failed to parse JSONPath: %v", err)
			}

			result := x.Get(jsonData)
			if len(result) == 0 {
				t.Error("Expected incident title to exist")

				return
			}

			if result[0] != "[TEST] Test Incident" {
				t.Errorf("Expected title '[TEST] Test Incident', got %v", result[0])
			}
		})
	})
}

// TestResponseProcessingWithNoIncluded verifies behavior when included array is empty.
func TestResponseProcessingWithNoIncluded(t *testing.T) {
	// Arrange
	data := []map[string]any{
		{
			"id":   "incident-1",
			"type": "incidents",
			"attributes": map[string]any{
				"title": "Test Incident",
			},
		},
	}

	included := []map[string]any{}

	// Act
	processedData := EnrichAllIncidentData(data, included)

	// Assert
	if len(processedData) != 1 {
		t.Errorf("Expected 1 incident, got %d", len(processedData))
	}

	incident := processedData[0]

	// Should not have included field when empty
	if _, hasIncluded := incident["included"]; hasIncluded {
		t.Error("Expected no 'included' field when included array is empty")
	}

	// Should not have any all_selected_* fields
	entityFields := []string{
		"all_selected_users",
		"all_selected_groups",
		"all_selected_services",
		"all_selected_options",
		"all_selected_functionalities",
		"all_selected_catalog_entities",
	}

	for _, field := range entityFields {
		if _, hasField := incident[field]; hasField {
			t.Errorf("Expected no '%s' field when included array is empty", field)
		}
	}
}

// TestResponseProcessingWithMultipleIncidents verifies processing of multiple incidents.
func TestResponseProcessingWithMultipleIncidents(t *testing.T) {
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
				"incident_id":   "incident-1",
				"form_field_id": "field-1",
				"selected_users": []any{
					map[string]any{"id": "user-1", "name": "User 1"},
				},
			},
		},
		{
			"attributes": map[string]any{
				"incident_id":   "incident-2",
				"form_field_id": "field-2",
				"selected_groups": []any{
					map[string]any{"id": "group-1", "name": "Group 1"},
				},
			},
		},
	}

	// Act
	processedData := EnrichAllIncidentData(data, included)

	// Assert
	if len(processedData) != 2 {
		t.Errorf("Expected 2 incidents, got %d", len(processedData))
	}

	// Incident 1 should have users
	incident1 := processedData[0]
	if users, ok := incident1["all_selected_users"].([]any); !ok || len(users) != 1 {
		t.Error("Incident 1 should have 1 user")
	}

	if _, hasGroups := incident1["all_selected_groups"]; hasGroups {
		t.Error("Incident 1 should not have groups")
	}

	// Incident 2 should have groups
	incident2 := processedData[1]
	if groups, ok := incident2["all_selected_groups"].([]any); !ok || len(groups) != 1 {
		t.Error("Incident 2 should have 1 group")
	}

	if _, hasUsers := incident2["all_selected_users"]; hasUsers {
		t.Error("Incident 2 should not have users")
	}
}

// TestJSONPathQueryWithSpecificFieldID tests the exact query pattern from the user's request.
// Query pattern: $.all_selected_users[?(@.field_id=="0c696034-4c87-4a29-b7a1-8a5524a443ca")].
func TestJSONPathQueryWithSpecificFieldID(t *testing.T) {
	// Arrange
	// Create test data with specific UUID field_id
	const specificFieldID = "0c696034-4c87-4a29-b7a1-8a5524a443ca"

	data := []map[string]any{
		{
			"id":   "incident-1",
			"type": "incidents",
			"attributes": map[string]any{
				"title": "Test Incident",
			},
		},
	}

	included := []map[string]any{
		{
			"attributes": map[string]any{
				"incident_id":   "incident-1",
				"form_field_id": specificFieldID,
				"selected_users": []any{
					map[string]any{"id": "user-1", "name": "Alice", "email": "alice@example.com"},
					map[string]any{"id": "user-2", "name": "Bob", "email": "bob@example.com"},
				},
			},
		},
		{
			"attributes": map[string]any{
				"incident_id":   "incident-1",
				"form_field_id": "different-field-id",
				"selected_users": []any{
					map[string]any{"id": "user-3", "name": "Charlie", "email": "charlie@example.com"},
				},
			},
		},
	}

	// Act
	processedData := EnrichAllIncidentData(data, included)

	// Convert to JSON for JSONPath queries
	processedJSON, err := json.Marshal(processedData)
	if err != nil {
		t.Fatalf("Failed to marshal processed data: %v", err)
	}

	// Parse JSON into interface{} for ojg/jp
	var jsonData any
	if err := json.Unmarshal(processedJSON, &jsonData); err != nil {
		t.Fatalf("Failed to unmarshal for JSONPath: %v", err)
	}

	// Assert
	// Test the exact query pattern: $.all_selected_users[?(@.field_id=="0c696034-4c87-4a29-b7a1-8a5524a443ca")]
	t.Run("filter by specific UUID field_id", func(t *testing.T) {
		// Build JSONPath query with the specific UUID
		query := "$[0].all_selected_users[?(@.field_id=='" + specificFieldID + "')]"

		x, err := jp.ParseString(query)
		if err != nil {
			t.Fatalf("Failed to parse JSONPath: %v", err)
		}

		result := x.Get(jsonData)

		if len(result) == 0 {
			t.Errorf("Expected to find users with field_id=='%s'", specificFieldID)
			t.Logf("Processed data: %s", oj.JSON(jsonData))

			return
		}

		// Should find 2 users (Alice and Bob)
		if len(result) != 2 {
			t.Errorf("Expected 2 users with field_id=='%s', got %d", specificFieldID, len(result))
		}

		// Verify the first user
		firstUser, ok := result[0].(map[string]any)
		if !ok {
			t.Error("Expected first result to be a map")

			return
		}

		if firstUser["id"] != "user-1" {
			t.Errorf("Expected user id 'user-1', got %v", firstUser["id"])
		}

		if firstUser["name"] != "Alice" {
			t.Errorf("Expected name 'Alice', got %v", firstUser["name"])
		}

		// Verify field_id is present
		if firstUser["field_id"] != specificFieldID {
			t.Errorf("Expected field_id '%s', got %v", specificFieldID, firstUser["field_id"])
		}
	})

	// Test filtering by different field_id
	t.Run("filter by different field_id", func(t *testing.T) {
		query := "$[0].all_selected_users[?(@.field_id=='different-field-id')]"

		x, err := jp.ParseString(query)
		if err != nil {
			t.Fatalf("Failed to parse JSONPath: %v", err)
		}

		result := x.Get(jsonData)

		if len(result) == 0 {
			t.Error("Expected to find user with field_id=='different-field-id'")

			return
		}

		// Should find 1 user (Charlie)
		if len(result) != 1 {
			t.Errorf("Expected 1 user with field_id=='different-field-id', got %d", len(result))
		}

		user, ok := result[0].(map[string]any)
		if !ok {
			t.Error("Expected result to be a map")

			return
		}

		if user["name"] != "Charlie" {
			t.Errorf("Expected name 'Charlie', got %v", user["name"])
		}
	})

	// Test query on all entity types with field_id
	t.Run("query all entity types by field_id", func(t *testing.T) {
		entityTypes := []string{
			"all_selected_users",
			"all_selected_groups",
			"all_selected_services",
			"all_selected_options",
			"all_selected_functionalities",
			"all_selected_catalog_entities",
		}

		for _, entityType := range entityTypes {
			t.Run(entityType, func(t *testing.T) {
				// Query: $[0].<entityType>[?(@.field_id)]
				query := "$[0]." + entityType + "[?(@.field_id)]"

				x, err := jp.ParseString(query)
				if err != nil {
					t.Fatalf("Failed to parse JSONPath: %v", err)
				}

				result := x.Get(jsonData)

				// Verify that entities have field_id
				for i, item := range result {
					itemMap, ok := item.(map[string]any)
					if !ok {
						t.Errorf("Item %d is not a map", i)

						continue
					}

					if itemMap["field_id"] == nil || itemMap["field_id"] == "" {
						t.Errorf("Item %d missing field_id", i)
					}
				}
			})
		}
	})
}
