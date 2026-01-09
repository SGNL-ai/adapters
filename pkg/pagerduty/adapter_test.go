// Copyright 2026 SGNL.ai, Inc.
package pagerduty_test

import (
	"reflect"
	"testing"

	"github.com/sgnl-ai/adapters/pkg/pagerduty"
)

func TestParseMap(t *testing.T) {
	// All input maps have values that are either string or []any, as ParseMap() assumes that condition is already met.
	tests := map[string]struct {
		inputMap map[string]map[string]any
		wantMap  map[string]map[string][]string
	}{
		"nil_map": {
			inputMap: nil,
			wantMap:  nil,
		},
		"valid_map": {
			inputMap: map[string]map[string]any{
				"key1": {
					"key2": "value1",
				},
				"key3": {
					// This is []any to simulate raw JSON responses.
					"key4": []any{"value2", "value3"},
				},
			},
			wantMap: map[string]map[string][]string{
				"key1": {
					"key2": []string{"value1"},
				},
				"key3": {
					"key4": []string{"value2", "value3"},
				},
			},
		},
		"empty_values": {
			inputMap: map[string]map[string]any{
				"key1": {
					"key2": "",
				},
				"key3": {
					// This is []any to simulate raw JSON responses.
					"key4": []any{},
				},
			},
			wantMap: map[string]map[string][]string{
				"key1": {
					"key2": []string{""},
				},
				"key3": {
					"key4": []string{},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotMap := pagerduty.ParseMap(tt.inputMap)

			if !reflect.DeepEqual(gotMap, tt.wantMap) {
				t.Errorf("gotMap: %v, wantMap: %v", gotMap, tt.wantMap)
			}
		})
	}
}
