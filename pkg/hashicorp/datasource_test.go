// Copyright 2026 SGNL.ai, Inc.

// nolint: lll, goconst
package hashicorp_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	hashicorp_adapter "github.com/sgnl-ai/adapters/pkg/hashicorp"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestParseResponse(t *testing.T) {
	tests := map[string]struct {
		body           []byte
		wantObjects    []map[string]interface{}
		wantNextCursor *string
		wantErr        *framework.Error
	}{
		"single_page": {
			body: []byte(`{"items": [{"id": "h_1234"}, {"id": "h_5678"}], "response_type": "complete", "list_token": ""}`),
			wantObjects: []map[string]interface{}{
				{"id": "h_1234"},
				{"id": "h_5678"},
			},
			wantNextCursor: nil,
		},
		"next_page": {
			body: []byte(`{"items": [{"id": "h_1234"}], "response_type": "delta", "list_token": "next-token"}`),
			wantObjects: []map[string]interface{}{
				{"id": "h_1234"},
			},
			wantNextCursor: testutil.GenPtr("next-token"),
		},
		"invalid_object_structure": {
			body: []byte(`[{"id": "h_1234"}, {"id": "h_5678"}]`),
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: json: cannot unmarshal array into Go value of type hashicorp.DatasourceResponse.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"missing_items_field": {
			body: []byte(`{"response_type": "complete", "list_token": ""}`),
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Items field is missing in the datasource response.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := hashicorp_adapter.ParseResponse(tt.body)

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !reflect.DeepEqual(gotNextCursor, tt.wantNextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotNextCursor, tt.wantNextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
