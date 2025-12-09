// Copyright 2025 SGNL.ai, Inc.
package commonutil_test

import (
	"reflect"
	"sort"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/commonutil"
)

func TestCreateChildEntitiesFromList(t *testing.T) {
	// Standard child entity config with id and value attributes
	childConfig := &framework.EntityConfig{
		Attributes: []*framework.AttributeConfig{
			{ExternalId: "id", Type: framework.AttributeTypeString},
			{ExternalId: "value", Type: framework.AttributeTypeString},
		},
	}

	tests := map[string]struct {
		parentID  string
		fieldName string
		values    []string
		want      []any
	}{
		"multiple_values": {
			parentID:  "123",
			fieldName: "Interests__c",
			values:    []string{"Sports", "Music", "Reading"},
			want: []any{
				map[string]any{"id": "123_Interests__c_sports", "value": "Sports"},
				map[string]any{"id": "123_Interests__c_music", "value": "Music"},
				map[string]any{"id": "123_Interests__c_reading", "value": "Reading"},
			},
		},
		"single_value": {
			parentID:  "456",
			fieldName: "Skills",
			values:    []string{"Technology"},
			want: []any{
				map[string]any{"id": "456_Skills_technology", "value": "Technology"},
			},
		},
		"values_with_spaces": {
			parentID:  "abc",
			fieldName: "Locations",
			values:    []string{"New York", "San Francisco"},
			want: []any{
				map[string]any{"id": "abc_Locations_new-york", "value": "New York"},
				map[string]any{"id": "abc_Locations_san-francisco", "value": "San Francisco"},
			},
		},
		"values_with_special_characters": {
			parentID:  "xyz",
			fieldName: "Languages",
			values:    []string{"C++", "C#"},
			want: []any{
				map[string]any{"id": "xyz_Languages_c++", "value": "C++"},
				map[string]any{"id": "xyz_Languages_c#", "value": "C#"},
			},
		},
		"values_with_slashes": {
			parentID:  "def",
			fieldName: "Paths",
			values:    []string{"A/B", "C\\D"},
			want: []any{
				map[string]any{"id": "def_Paths_a-b", "value": "A/B"},
				map[string]any{"id": "def_Paths_c-d", "value": "C\\D"},
			},
		},
		"duplicate_values_case_insensitive": {
			parentID:  "ghi",
			fieldName: "Items",
			values:    []string{"Apple", "APPLE", "banana", "Banana"},
			want: []any{
				map[string]any{"id": "ghi_Items_apple", "value": "Apple"},
				map[string]any{"id": "ghi_Items_banana", "value": "banana"},
			},
		},
		"values_with_whitespace": {
			parentID:  "jkl",
			fieldName: "Names",
			values:    []string{" John ", "Jane", "  Bob  "},
			want: []any{
				map[string]any{"id": "jkl_Names_john", "value": "John"},
				map[string]any{"id": "jkl_Names_jane", "value": "Jane"},
				map[string]any{"id": "jkl_Names_bob", "value": "Bob"},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := commonutil.CreateChildEntitiesFromList(tt.parentID, tt.fieldName, tt.values, childConfig)

			// Sort both got and want by ID for order-independent comparison
			sortByID := func(items []any) {
				sort.Slice(items, func(i, j int) bool {
					// Check first type assertion
					item1, ok := items[i].(map[string]any)
					if !ok {
						t.Fatalf("items[%d] is not map[string]any", i)
					}

					item2, ok := items[j].(map[string]any)
					if !ok {
						t.Fatalf("items[%d] is not map[string]any", j)
					}

					// Check second type assertion
					id1, ok := item1["id"].(string)
					if !ok {
						t.Fatalf("items[%d][\"id\"] is not string", i)
					}

					id2, ok := item2["id"].(string)
					if !ok {
						t.Fatalf("items[%d][\"id\"] is not string", j)
					}

					return id1 < id2
				})
			}

			sortByID(got)
			sortByID(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateChildEntitiesFromList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUniqueIDValue(t *testing.T) {
	tests := map[string]struct {
		obj          map[string]any
		entityConfig *framework.EntityConfig
		wantID       string
		wantOK       bool
	}{
		"valid_unique_id": {
			obj: map[string]any{"Id": "123", "Name": "Test"},
			entityConfig: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "Id", UniqueId: true},
				},
			},
			wantID: "123",
			wantOK: true,
		},
		"no_unique_id_configured": {
			obj: map[string]any{"Id": "123", "Name": "Test"},
			entityConfig: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "Id", UniqueId: false},
				},
			},
			wantID: "",
			wantOK: false,
		},
		"unique_id_field_missing": {
			obj: map[string]any{"Name": "Test"},
			entityConfig: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "Id", UniqueId: true},
				},
			},
			wantID: "",
			wantOK: false,
		},
		"unique_id_value_is_nil": {
			obj: map[string]any{"Id": nil, "Name": "Test"},
			entityConfig: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "Id", UniqueId: true},
				},
			},
			wantID: "",
			wantOK: false,
		},
		"unique_id_value_not_string": {
			obj: map[string]any{"Id": 123, "Name": "Test"},
			entityConfig: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "Id", UniqueId: true},
				},
			},
			wantID: "",
			wantOK: false,
		},
		"unique_id_value_empty_string": {
			obj: map[string]any{"Id": "", "Name": "Test"},
			entityConfig: &framework.EntityConfig{
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "Id", UniqueId: true},
				},
			},
			wantID: "",
			wantOK: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotID, gotOK := commonutil.GetUniqueIDValue(tt.obj, tt.entityConfig)

			if gotID != tt.wantID || gotOK != tt.wantOK {
				t.Errorf("GetUniqueIDValue() = (%v, %v), want (%v, %v)", gotID, gotOK, tt.wantID, tt.wantOK)
			}
		})
	}
}

func TestCreateChildEntitiesFromDelimitedString(t *testing.T) {
	parentConfig := &framework.EntityConfig{
		Attributes: []*framework.AttributeConfig{
			{ExternalId: "Id", UniqueId: true, Type: framework.AttributeTypeString},
		},
	}

	childConfig := &framework.EntityConfig{
		ExternalId: "Interests",
		Attributes: []*framework.AttributeConfig{
			{ExternalId: "id", Type: framework.AttributeTypeString},
			{ExternalId: "value", Type: framework.AttributeTypeString},
		},
	}

	tests := map[string]struct {
		objects              []map[string]any
		parentEntityConfig   *framework.EntityConfig
		childEntities        []*framework.EntityConfig
		delimiter            string
		wantTransformedCount int
		validateFirst        func(t *testing.T, obj map[string]any)
	}{
		"semicolon_delimited_with_duplicates": {
			objects: []map[string]any{
				{"Id": "123", "Name": "John", "Interests": "Sports;Music;Sports"},
			},
			parentEntityConfig:   parentConfig,
			childEntities:        []*framework.EntityConfig{childConfig},
			delimiter:            ";",
			wantTransformedCount: 1,
			validateFirst: func(t *testing.T, obj map[string]any) {
				interests, ok := obj["Interests"].([]any)
				if !ok {
					t.Fatal("Interests should be []any")
				}
				if len(interests) != 2 {
					t.Errorf("Expected 2 deduplicated interests, got %d", len(interests))
				}
			},
		},
		"nil_value": {
			objects: []map[string]any{
				{"Id": "456", "Name": "Jane", "Interests": nil},
			},
			parentEntityConfig:   parentConfig,
			childEntities:        []*framework.EntityConfig{childConfig},
			delimiter:            ";",
			wantTransformedCount: 1,
			validateFirst: func(t *testing.T, obj map[string]any) {
				interests, ok := obj["Interests"].([]any)
				if !ok {
					t.Fatal("Interests should be []any")
				}
				if len(interests) != 0 {
					t.Errorf("Expected empty array for nil, got %d items", len(interests))
				}
			},
		},
		"empty_string": {
			objects: []map[string]any{
				{"Id": "789", "Name": "Bob", "Interests": ""},
			},
			parentEntityConfig:   parentConfig,
			childEntities:        []*framework.EntityConfig{childConfig},
			delimiter:            ";",
			wantTransformedCount: 1,
			validateFirst: func(t *testing.T, obj map[string]any) {
				interests, ok := obj["Interests"].([]any)
				if !ok {
					t.Fatal("Interests should be []any")
				}
				if len(interests) != 0 {
					t.Errorf("Expected empty array for empty string, got %d items", len(interests))
				}
			},
		},
		"no_child_entities": {
			objects: []map[string]any{
				{"Id": "123", "Name": "John"},
			},
			parentEntityConfig:   parentConfig,
			childEntities:        []*framework.EntityConfig{},
			delimiter:            ";",
			wantTransformedCount: 1,
			validateFirst: func(t *testing.T, obj map[string]any) {
				if _, exists := obj["Interests"]; exists {
					t.Error("Interests field should not exist when no child entities configured")
				}
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := commonutil.CreateChildEntitiesFromDelimitedString(
				tt.objects,
				tt.parentEntityConfig,
				tt.childEntities,
				tt.delimiter,
			)

			if len(got) != tt.wantTransformedCount {
				t.Fatalf("Expected %d transformed objects, got %d", tt.wantTransformedCount, len(got))
			}

			if tt.validateFirst != nil && len(got) > 0 {
				tt.validateFirst(t, got[0])
			}
		})
	}
}
