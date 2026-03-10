// Copyright 2026 SGNL.ai, Inc.

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
					id1, _ := items[i].(map[string]any)["id"].(string)
					id2, _ := items[j].(map[string]any)["id"].(string)

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

func TestCreateChildEntitiesFromStringArray(t *testing.T) {
	parentConfig := &framework.EntityConfig{
		Attributes: []*framework.AttributeConfig{
			{ExternalId: "incidentNumber", UniqueId: true, Type: framework.AttributeTypeString},
		},
	}

	pagedUsersConfig := &framework.EntityConfig{
		ExternalId: "pagedUsers",
		Attributes: []*framework.AttributeConfig{
			{ExternalId: "id", Type: framework.AttributeTypeString},
			{ExternalId: "value", Type: framework.AttributeTypeString},
		},
	}

	pagedTeamsConfig := &framework.EntityConfig{
		ExternalId: "pagedTeams",
		Attributes: []*framework.AttributeConfig{
			{ExternalId: "id", Type: framework.AttributeTypeString},
			{ExternalId: "value", Type: framework.AttributeTypeString},
		},
	}

	tests := map[string]struct {
		objects      []map[string]any
		entityConfig *framework.EntityConfig
		want         []map[string]any
	}{
		"string_array_to_child_entities": {
			objects: []map[string]any{
				{
					"incidentNumber": "1",
					"pagedUsers":     []any{"alice", "bob"},
				},
			},
			entityConfig: &framework.EntityConfig{
				Attributes:    parentConfig.Attributes,
				ChildEntities: []*framework.EntityConfig{pagedUsersConfig},
			},
			want: []map[string]any{
				{
					"incidentNumber": "1",
					"pagedUsers": []any{
						map[string]any{"id": "1_pagedUsers_alice", "value": "alice"},
						map[string]any{"id": "1_pagedUsers_bob", "value": "bob"},
					},
				},
			},
		},
		"nil_array_value": {
			objects: []map[string]any{
				{
					"incidentNumber": "2",
					"pagedUsers":     nil,
				},
			},
			entityConfig: &framework.EntityConfig{
				Attributes:    parentConfig.Attributes,
				ChildEntities: []*framework.EntityConfig{pagedUsersConfig},
			},
			want: []map[string]any{
				{
					"incidentNumber": "2",
					"pagedUsers":     []any{},
				},
			},
		},
		"empty_array": {
			objects: []map[string]any{
				{
					"incidentNumber": "3",
					"pagedUsers":     []any{},
				},
			},
			entityConfig: &framework.EntityConfig{
				Attributes:    parentConfig.Attributes,
				ChildEntities: []*framework.EntityConfig{pagedUsersConfig},
			},
			want: []map[string]any{
				{
					"incidentNumber": "3",
					"pagedUsers":     []any{},
				},
			},
		},
		"no_child_entities_configured": {
			objects: []map[string]any{
				{
					"incidentNumber": "4",
					"pagedUsers":     []any{"alice"},
				},
			},
			entityConfig: &framework.EntityConfig{
				Attributes: parentConfig.Attributes,
			},
			want: []map[string]any{
				{
					"incidentNumber": "4",
					"pagedUsers":     []any{"alice"},
				},
			},
		},
		"multiple_child_entity_fields": {
			objects: []map[string]any{
				{
					"incidentNumber": "5",
					"pagedUsers":     []any{"alice", "bob"},
					"pagedTeams":     []any{"team-alpha"},
				},
			},
			entityConfig: &framework.EntityConfig{
				Attributes:    parentConfig.Attributes,
				ChildEntities: []*framework.EntityConfig{pagedUsersConfig, pagedTeamsConfig},
			},
			want: []map[string]any{
				{
					"incidentNumber": "5",
					"pagedUsers": []any{
						map[string]any{"id": "5_pagedUsers_alice", "value": "alice"},
						map[string]any{"id": "5_pagedUsers_bob", "value": "bob"},
					},
					"pagedTeams": []any{
						map[string]any{"id": "5_pagedTeams_team-alpha", "value": "team-alpha"},
					},
				},
			},
		},
		"field_not_present_in_object": {
			objects: []map[string]any{
				{
					"incidentNumber": "6",
				},
			},
			entityConfig: &framework.EntityConfig{
				Attributes:    parentConfig.Attributes,
				ChildEntities: []*framework.EntityConfig{pagedUsersConfig},
			},
			want: []map[string]any{
				{
					"incidentNumber": "6",
				},
			},
		},
		"non_string_array_elements_unchanged": {
			objects: []map[string]any{
				{
					"incidentNumber": "8",
					"pagedUsers":     []any{map[string]any{"name": "alice"}, map[string]any{"name": "bob"}},
				},
			},
			entityConfig: &framework.EntityConfig{
				Attributes:    parentConfig.Attributes,
				ChildEntities: []*framework.EntityConfig{pagedUsersConfig},
			},
			want: []map[string]any{
				{
					"incidentNumber": "8",
					"pagedUsers":     []any{map[string]any{"name": "alice"}, map[string]any{"name": "bob"}},
				},
			},
		},
		"duplicate_values_deduplicated": {
			objects: []map[string]any{
				{
					"incidentNumber": "7",
					"pagedUsers":     []any{"alice", "Alice", "bob"},
				},
			},
			entityConfig: &framework.EntityConfig{
				Attributes:    parentConfig.Attributes,
				ChildEntities: []*framework.EntityConfig{pagedUsersConfig},
			},
			want: []map[string]any{
				{
					"incidentNumber": "7",
					"pagedUsers": []any{
						map[string]any{"id": "7_pagedUsers_alice", "value": "alice"},
						map[string]any{"id": "7_pagedUsers_bob", "value": "bob"},
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := commonutil.CreateChildEntitiesFromStringArray(
				tt.objects,
				tt.entityConfig,
			)

			// Sort child entity arrays by ID for order-independent comparison.
			// Only sort arrays whose elements are maps (child entity objects).
			sortByID := func(items []any) {
				if len(items) == 0 {
					return
				}

				if _, ok := items[0].(map[string]any); !ok {
					return
				}

				sort.Slice(items, func(i, j int) bool {
					id1, _ := items[i].(map[string]any)["id"].(string)
					id2, _ := items[j].(map[string]any)["id"].(string)

					return id1 < id2
				})
			}

			for i := range got {
				for _, value := range got[i] {
					if childArray, ok := value.([]any); ok {
						sortByID(childArray)
					}
				}
			}

			for i := range tt.want {
				for _, value := range tt.want[i] {
					if childArray, ok := value.([]any); ok {
						sortByID(childArray)
					}
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateChildEntitiesFromStringArray() = %v, want %v", got, tt.want)
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

	skillsConfig := &framework.EntityConfig{
		ExternalId: "Skills",
		Attributes: []*framework.AttributeConfig{
			{ExternalId: "id", Type: framework.AttributeTypeString},
			{ExternalId: "value", Type: framework.AttributeTypeString},
		},
	}

	tests := map[string]struct {
		objects      []map[string]any
		entityConfig *framework.EntityConfig
		delimiter    string
		want         []map[string]any
	}{
		"semicolon_delimited_with_duplicates": {
			objects: []map[string]any{
				{"Id": "123", "Name": "John", "Interests": "Sports;Music;Sports"},
			},
			entityConfig: &framework.EntityConfig{
				Attributes:    parentConfig.Attributes,
				ChildEntities: []*framework.EntityConfig{childConfig},
			},
			delimiter: ";",
			want: []map[string]any{
				{
					"Id":   "123",
					"Name": "John",
					"Interests": []any{
						map[string]any{"id": "123_Interests_sports", "value": "Sports"},
						map[string]any{"id": "123_Interests_music", "value": "Music"},
					},
				},
			},
		},
		"nil_value": {
			objects: []map[string]any{
				{"Id": "456", "Name": "Jane", "Interests": nil},
			},
			entityConfig: &framework.EntityConfig{
				Attributes:    parentConfig.Attributes,
				ChildEntities: []*framework.EntityConfig{childConfig},
			},
			delimiter: ";",
			want: []map[string]any{
				{
					"Id":        "456",
					"Name":      "Jane",
					"Interests": []any{},
				},
			},
		},
		"empty_string": {
			objects: []map[string]any{
				{"Id": "789", "Name": "Bob", "Interests": ""},
			},
			entityConfig: &framework.EntityConfig{
				Attributes:    parentConfig.Attributes,
				ChildEntities: []*framework.EntityConfig{childConfig},
			},
			delimiter: ";",
			want: []map[string]any{
				{
					"Id":        "789",
					"Name":      "Bob",
					"Interests": []any{},
				},
			},
		},
		"no_child_entities": {
			objects: []map[string]any{
				{"Id": "123", "Name": "John"},
			},
			entityConfig: &framework.EntityConfig{
				Attributes:    parentConfig.Attributes,
				ChildEntities: []*framework.EntityConfig{},
			},
			delimiter: ";",
			want: []map[string]any{
				{"Id": "123", "Name": "John"},
			},
		},
		"single_value_without_delimiter": {
			objects: []map[string]any{
				{"Id": "999", "Name": "Alice", "Interests": "Technology"},
			},
			entityConfig: &framework.EntityConfig{
				Attributes:    parentConfig.Attributes,
				ChildEntities: []*framework.EntityConfig{childConfig},
			},
			delimiter: ";",
			want: []map[string]any{
				{
					"Id":   "999",
					"Name": "Alice",
					"Interests": []any{
						map[string]any{"id": "999_Interests_technology", "value": "Technology"},
					},
				},
			},
		},
		"multiple_child_entities": {
			objects: []map[string]any{
				{"Id": "888", "Name": "Charlie", "Interests": "Sports;Music", "Skills": "Go;Python"},
			},
			entityConfig: &framework.EntityConfig{
				Attributes:    parentConfig.Attributes,
				ChildEntities: []*framework.EntityConfig{childConfig, skillsConfig},
			},
			delimiter: ";",
			want: []map[string]any{
				{
					"Id":   "888",
					"Name": "Charlie",
					"Interests": []any{
						map[string]any{"id": "888_Interests_sports", "value": "Sports"},
						map[string]any{"id": "888_Interests_music", "value": "Music"},
					},
					"Skills": []any{
						map[string]any{"id": "888_Skills_go", "value": "Go"},
						map[string]any{"id": "888_Skills_python", "value": "Python"},
					},
				},
			},
		},
		"non_string_field_value": {
			objects: []map[string]any{
				{"Id": "777", "Name": "David", "Interests": 12345},
			},
			entityConfig: &framework.EntityConfig{
				Attributes:    parentConfig.Attributes,
				ChildEntities: []*framework.EntityConfig{childConfig},
			},
			delimiter: ";",
			want: []map[string]any{
				{
					"Id":        "777",
					"Name":      "David",
					"Interests": 12345,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := commonutil.CreateChildEntitiesFromDelimitedString(
				tt.objects,
				tt.entityConfig,
				tt.delimiter,
			)

			// Sort child entity arrays by ID for order-independent comparison
			sortByID := func(items []any) {
				sort.Slice(items, func(i, j int) bool {
					id1, _ := items[i].(map[string]any)["id"].(string)
					id2, _ := items[j].(map[string]any)["id"].(string)

					return id1 < id2
				})
			}

			// Sort both got and want child entity arrays
			for i := range got {
				for _, value := range got[i] {
					if childArray, ok := value.([]any); ok {
						sortByID(childArray)
					}
				}
			}

			for i := range tt.want {
				for _, value := range tt.want[i] {
					if childArray, ok := value.([]any); ok {
						sortByID(childArray)
					}
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateChildEntitiesFromDelimitedString() = %v, want %v", got, tt.want)
			}
		})
	}
}
