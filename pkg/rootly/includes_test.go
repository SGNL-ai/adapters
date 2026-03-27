// Copyright 2026 SGNL.ai, Inc.

package rootly

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveIncludedRelationships_GivenArrayRelationshipStubs_WhenResolved_ThenReplacedWithFullObjects(t *testing.T) {
	tests := []struct {
		name     string
		data     []map[string]any
		included []map[string]any
		check    func(t *testing.T, result []map[string]any)
	}{
		{
			name: "resolves_array_relationship_stubs",
			data: []map[string]any{
				{
					"id":   "incident-1",
					"type": "incidents",
					"relationships": map[string]any{
						"roles": map[string]any{
							"data": []any{
								map[string]any{"id": "role-1", "type": "incident_role_assignments"},
								map[string]any{"id": "role-2", "type": "incident_role_assignments"},
							},
						},
					},
				},
			},
			included: []map[string]any{
				{
					"id":   "role-1",
					"type": "incident_role_assignments",
					"attributes": map[string]any{
						"user": map[string]any{
							"data": map[string]any{
								"attributes": map[string]any{
									"email": "alice@example.com",
								},
							},
						},
					},
				},
				{
					"id":   "role-2",
					"type": "incident_role_assignments",
					"attributes": map[string]any{
						"incident_role": map[string]any{
							"data": map[string]any{
								"attributes": map[string]any{
									"name": "Communications Lead",
								},
							},
						},
					},
				},
			},
			check: func(t *testing.T, result []map[string]any) {
				require.Len(t, result, 1)

				rels, ok := result[0]["relationships"].(map[string]any)
				require.True(t, ok)

				roles, ok := rels["roles"].(map[string]any)
				require.True(t, ok)

				data, ok := roles["data"].([]any)
				require.True(t, ok)
				require.Len(t, data, 2)

				// First role should have full attributes
				role1, ok := data[0].(map[string]any)
				require.True(t, ok)

				attrs1, ok := role1["attributes"].(map[string]any)
				require.True(t, ok)

				user, ok := attrs1["user"].(map[string]any)
				require.True(t, ok)

				userData, ok := user["data"].(map[string]any)
				require.True(t, ok)

				userAttrs, ok := userData["attributes"].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, "alice@example.com", userAttrs["email"])

				// Second role should also be resolved
				role2, ok := data[1].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, "role-2", role2["id"])

				attrs2, ok := role2["attributes"].(map[string]any)
				require.True(t, ok)
				assert.NotNil(t, attrs2["incident_role"])
			},
		},
		{
			name: "resolves_single_object_relationship_stub",
			data: []map[string]any{
				{
					"id":   "incident-1",
					"type": "incidents",
					"relationships": map[string]any{
						"severity": map[string]any{
							"data": map[string]any{"id": "sev-1", "type": "severities"},
						},
					},
				},
			},
			included: []map[string]any{
				{
					"id":   "sev-1",
					"type": "severities",
					"attributes": map[string]any{
						"name": "Critical",
						"slug": "critical",
					},
				},
			},
			check: func(t *testing.T, result []map[string]any) {
				require.Len(t, result, 1)

				rels, ok := result[0]["relationships"].(map[string]any)
				require.True(t, ok)

				severity, ok := rels["severity"].(map[string]any)
				require.True(t, ok)

				data, ok := severity["data"].(map[string]any)
				require.True(t, ok)

				attrs, ok := data["attributes"].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, "Critical", attrs["name"])
				assert.Equal(t, "critical", attrs["slug"])
			},
		},
		{
			name: "preserves_unmatched_stubs",
			data: []map[string]any{
				{
					"id":   "incident-1",
					"type": "incidents",
					"relationships": map[string]any{
						"roles": map[string]any{
							"data": []any{
								map[string]any{"id": "role-unknown", "type": "incident_role_assignments"},
							},
						},
					},
				},
			},
			included: []map[string]any{},
			check: func(t *testing.T, result []map[string]any) {
				require.Len(t, result, 1)

				rels, ok := result[0]["relationships"].(map[string]any)
				require.True(t, ok)

				roles, ok := rels["roles"].(map[string]any)
				require.True(t, ok)

				data, ok := roles["data"].([]any)
				require.True(t, ok)
				require.Len(t, data, 1)

				stub, ok := data[0].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, "role-unknown", stub["id"])

				_, hasAttrs := stub["attributes"]
				assert.False(t, hasAttrs, "unmatched stub should not have attributes")
			},
		},
		{
			name: "handles_missing_relationships_gracefully",
			data: []map[string]any{
				{
					"id":   "incident-1",
					"type": "incidents",
					"attributes": map[string]any{
						"title": "Some incident",
					},
				},
			},
			included: []map[string]any{
				{
					"id":   "role-1",
					"type": "incident_role_assignments",
					"attributes": map[string]any{
						"user": nil,
					},
				},
			},
			check: func(t *testing.T, result []map[string]any) {
				require.Len(t, result, 1)

				_, hasRels := result[0]["relationships"]
				assert.False(t, hasRels, "should not add relationships when none existed")

				attrs, ok := result[0]["attributes"].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, "Some incident", attrs["title"])
			},
		},
		{
			name: "handles_empty_included_array",
			data: []map[string]any{
				{
					"id":   "incident-1",
					"type": "incidents",
					"relationships": map[string]any{
						"roles": map[string]any{
							"data": []any{
								map[string]any{"id": "role-1", "type": "incident_role_assignments"},
							},
						},
					},
				},
			},
			included: []map[string]any{},
			check: func(t *testing.T, result []map[string]any) {
				require.Len(t, result, 1)

				// Data should be returned unchanged
				rels, ok := result[0]["relationships"].(map[string]any)
				require.True(t, ok)

				roles, ok := rels["roles"].(map[string]any)
				require.True(t, ok)

				data, ok := roles["data"].([]any)
				require.True(t, ok)
				require.Len(t, data, 1)
			},
		},
		{
			name: "handles_nil_data",
			data: nil,
			included: []map[string]any{
				{"id": "role-1", "type": "roles"},
			},
			check: func(t *testing.T, result []map[string]any) {
				assert.Nil(t, result)
			},
		},
		{
			name: "handles_empty_data",
			data: []map[string]any{},
			included: []map[string]any{
				{"id": "role-1", "type": "roles"},
			},
			check: func(t *testing.T, result []map[string]any) {
				assert.Empty(t, result)
			},
		},
		{
			name: "handles_stubs_without_id_or_type",
			data: []map[string]any{
				{
					"id":   "incident-1",
					"type": "incidents",
					"relationships": map[string]any{
						"roles": map[string]any{
							"data": []any{
								map[string]any{"id": "role-1"},
								map[string]any{"type": "incident_role_assignments"},
								map[string]any{"other": "value"},
							},
						},
					},
				},
			},
			included: []map[string]any{},
			check: func(t *testing.T, result []map[string]any) {
				require.Len(t, result, 1)

				rels, ok := result[0]["relationships"].(map[string]any)
				require.True(t, ok)
				roles, ok := rels["roles"].(map[string]any)
				require.True(t, ok)
				data, ok := roles["data"].([]any)
				require.True(t, ok)
				require.Len(t, data, 3, "all stubs should be preserved even without id/type")
			},
		},
		{
			name: "handles_empty_data_array_in_relationship",
			data: []map[string]any{
				{
					"id":   "incident-1",
					"type": "incidents",
					"relationships": map[string]any{
						"roles": map[string]any{
							"data": []any{},
						},
					},
				},
			},
			included: []map[string]any{
				{
					"id":   "role-1",
					"type": "incident_role_assignments",
					"attributes": map[string]any{
						"user": nil,
					},
				},
			},
			check: func(t *testing.T, result []map[string]any) {
				require.Len(t, result, 1)

				rels, ok := result[0]["relationships"].(map[string]any)
				require.True(t, ok)

				roles, ok := rels["roles"].(map[string]any)
				require.True(t, ok)

				data, ok := roles["data"].([]any)
				require.True(t, ok)
				assert.Empty(t, data, "empty roles data should remain empty")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := resolveIncludedRelationships(tt.data, tt.included)

			// Assert
			tt.check(t, result)
		})
	}
}

func TestMergeMaps_GivenTwoMaps_WhenMerged_ThenSecondTakesPrecedence(t *testing.T) {
	tests := []struct {
		name     string
		m1       map[string]any
		m2       map[string]any
		expected map[string]any
	}{
		{
			name:     "non_overlapping_keys",
			m1:       map[string]any{"a": 1},
			m2:       map[string]any{"b": 2},
			expected: map[string]any{"a": 1, "b": 2},
		},
		{
			name:     "overlapping_keys_m2_wins",
			m1:       map[string]any{"id": "stub-id", "type": "roles"},
			m2:       map[string]any{"id": "full-id", "type": "roles", "attributes": map[string]any{"name": "Commander"}},
			expected: map[string]any{"id": "full-id", "type": "roles", "attributes": map[string]any{"name": "Commander"}},
		},
		{
			name:     "empty_m1",
			m1:       map[string]any{},
			m2:       map[string]any{"a": 1},
			expected: map[string]any{"a": 1},
		},
		{
			name:     "empty_m2",
			m1:       map[string]any{"a": 1},
			m2:       map[string]any{},
			expected: map[string]any{"a": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := mergeMaps(tt.m1, tt.m2)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResolveIncludedRelationships_GivenRolesAndFormFields_WhenBothStagesRun_ThenBothResolved(t *testing.T) {
	// Arrange: an incident with role relationships AND form field selections in included
	data := []map[string]any{
		{
			"id":   "incident-1",
			"type": "incidents",
			"relationships": map[string]any{
				"roles": map[string]any{
					"data": []any{
						map[string]any{"id": "role-1", "type": "incident_role_assignments"},
					},
				},
			},
		},
	}

	included := []map[string]any{
		{
			"id":   "role-1",
			"type": "incident_role_assignments",
			"attributes": map[string]any{
				"incident_role": map[string]any{
					"data": map[string]any{
						"attributes": map[string]any{
							"name": "Incident Commander",
							"slug": "incident-commander",
						},
					},
				},
				"user": map[string]any{
					"data": map[string]any{
						"attributes": map[string]any{
							"email": "alice@example.com",
						},
					},
				},
			},
		},
		{
			"id":   "field-sel-1",
			"type": "incident_form_field_selections",
			"attributes": map[string]any{
				"incident_id":   "incident-1",
				"form_field_id": "field-1",
				"selected_users": []any{
					map[string]any{"id": "user-1", "name": "Bob"},
				},
			},
		},
	}

	// Act: Stage 1 — resolve relationship stubs
	resolved := resolveIncludedRelationships(data, included)

	// Act: Stage 2 — enrich with form-field data
	enriched := EnrichAllIncidentData(resolved, included)

	// Assert: role relationships are resolved
	require.Len(t, enriched, 1)

	rels, ok := enriched[0]["relationships"].(map[string]any)
	require.True(t, ok)

	roles, ok := rels["roles"].(map[string]any)
	require.True(t, ok)

	roleData, ok := roles["data"].([]any)
	require.True(t, ok)
	require.Len(t, roleData, 1)

	role1, ok := roleData[0].(map[string]any)
	require.True(t, ok)

	attrs, ok := role1["attributes"].(map[string]any)
	require.True(t, ok)

	user, ok := attrs["user"].(map[string]any)
	require.True(t, ok)

	userData, ok := user["data"].(map[string]any)
	require.True(t, ok)

	userAttrs, ok := userData["attributes"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "alice@example.com", userAttrs["email"])

	// Assert: form-field enrichment also worked
	users, ok := enriched[0]["all_selected_users"].([]any)
	require.True(t, ok, "all_selected_users should be present from Stage 2 enrichment")
	require.Len(t, users, 1)

	user1, ok := users[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Bob", user1["name"])
}

func TestResolveIncludedRelationships_GivenOriginalData_WhenResolved_ThenOriginalStubsNotMutated(t *testing.T) {
	// Arrange: create an incident with a relationship stub.
	originalStub := map[string]any{"id": "role-1", "type": "incident_role_assignments"}

	data := []map[string]any{
		{
			"id":   "incident-1",
			"type": "incidents",
			"relationships": map[string]any{
				"roles": map[string]any{
					"data": []any{originalStub},
				},
			},
		},
	}

	included := []map[string]any{
		{
			"id":   "role-1",
			"type": "incident_role_assignments",
			"attributes": map[string]any{
				"incident_role": map[string]any{
					"data": map[string]any{
						"attributes": map[string]any{"name": "Commander"},
					},
				},
			},
		},
	}

	// Act
	_ = resolveIncludedRelationships(data, included)

	// Assert: the original stub map should NOT have been mutated with "attributes".
	_, hasAttrs := originalStub["attributes"]
	assert.False(t, hasAttrs, "resolveIncludedRelationships should not mutate the original stub — mergeMaps creates a new map")
}
