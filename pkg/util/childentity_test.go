// Copyright 2025 SGNL.ai, Inc.
package util_test

import (
	"reflect"
	"testing"

	"github.com/sgnl-ai/adapters/pkg/util"
)

func TestCreateChildEntitiesFromValues(t *testing.T) {
	tests := map[string]struct {
		parentID string
		values   []string
		want     []any
	}{
		"multiple_values": {
			parentID: "123",
			values:   []string{"Sports", "Music", "Reading"},
			want: []any{
				map[string]any{"id": "123_Sports", "value": "Sports"},
				map[string]any{"id": "123_Music", "value": "Music"},
				map[string]any{"id": "123_Reading", "value": "Reading"},
			},
		},
		"single_value": {
			parentID: "456",
			values:   []string{"Technology"},
			want: []any{
				map[string]any{"id": "456_Technology", "value": "Technology"},
			},
		},
		"empty_values": {
			parentID: "789",
			values:   []string{},
			want:     []any{},
		},
		"empty_parent_id": {
			parentID: "",
			values:   []string{"Value1", "Value2"},
			want: []any{
				map[string]any{"id": "_Value1", "value": "Value1"},
				map[string]any{"id": "_Value2", "value": "Value2"},
			},
		},
		"values_with_spaces": {
			parentID: "abc",
			values:   []string{"New York", "San Francisco"},
			want: []any{
				map[string]any{"id": "abc_New York", "value": "New York"},
				map[string]any{"id": "abc_San Francisco", "value": "San Francisco"},
			},
		},
		"values_with_special_characters": {
			parentID: "xyz",
			values:   []string{"C++", "C#"},
			want: []any{
				map[string]any{"id": "xyz_C++", "value": "C++"},
				map[string]any{"id": "xyz_C#", "value": "C#"},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := util.CreateChildEntitiesFromValues(tt.parentID, tt.values)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateChildEntitiesFromValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
