// Copyright 2026 SGNL.ai, Inc.

// nolint: lll, goconst
package crowdstrike

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestParseGraphQLResponse(t *testing.T) {
	tests := map[string]struct {
		body             []byte
		entityExternalID string
		wantObjects      []map[string]any
		wantNextCursor   string
		wantErr          *framework.Error
	}{
		"entities_single_page": {
			body:             []byte(`{"entities": {"nodes": [{"entityId": "e1"}, {"entityId": "e2"}], "pageInfo": {"hasNextPage": false, "endCursor": ""}}}`),
			entityExternalID: "entity",
			wantObjects: []map[string]any{
				{"entityId": "e1"},
				{"entityId": "e2"},
			},
			wantNextCursor: "",
		},
		"entities_next_page": {
			body:             []byte(`{"entities": {"nodes": [{"entityId": "e1"}], "pageInfo": {"hasNextPage": true, "endCursor": "cursor-1"}}}`),
			entityExternalID: "entity",
			wantObjects: []map[string]any{
				{"entityId": "e1"},
			},
			wantNextCursor: "cursor-1",
		},
		"incidents_read_from_incidents_field": {
			body:             []byte(`{"incidents": {"nodes": [{"incidentId": "i1"}], "pageInfo": {"hasNextPage": false, "endCursor": ""}}}`),
			entityExternalID: Incident,
			wantObjects: []map[string]any{
				{"incidentId": "i1"},
			},
			wantNextCursor: "",
		},
		// An empty collection is legitimate and must still sync as an empty page.
		"empty_nodes": {
			body:             []byte(`{"entities": {"nodes": [], "pageInfo": {"hasNextPage": false, "endCursor": ""}}}`),
			entityExternalID: "entity",
			wantObjects:      []map[string]any{},
			wantNextCursor:   "",
		},
		"missing_nodes_field": {
			body:             []byte(`{"entities": {"pageInfo": {"hasNextPage": false, "endCursor": ""}}}`),
			entityExternalID: "entity",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Nodes field is missing in the datasource response.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		// A null body is reachable: getGraphQLPage marshals a nil respData to "null".
		// This must report an error rather than panic.
		"null_body": {
			body:             []byte(`null`),
			entityExternalID: "entity",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Nodes field is missing in the datasource response.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"invalid_object_structure": {
			body:             []byte(`[{"entityId": "e1"}]`),
			entityExternalID: "entity",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: json: cannot unmarshal array into Go value of type crowdstrike.DatasourceResponse.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := parseGraphQLResponse(tt.body, tt.entityExternalID)

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if gotNextCursor != tt.wantNextCursor {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotNextCursor, tt.wantNextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
