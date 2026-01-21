// Copyright 2026 SGNL.ai, Inc.

package hashicorp_test

import (
	"testing"

	hashicorp_adapter "github.com/sgnl-ai/adapters/pkg/hashicorp"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestConstructEndpoint(t *testing.T) {
	tests := []struct {
		name          string
		request       *hashicorp_adapter.Request
		expectedURL   string
		expectedError bool
	}{
		{
			name: "basic_endpoint_construction",
			request: &hashicorp_adapter.Request{
				BaseURL:          "https://boundary.example.com",
				EntityExternalID: "roles",
				PageSize:         100,
				EntityConfig: map[string]hashicorp_adapter.EntityConfig{
					"roles": {
						ScopeID: "global",
					},
				},
			},
			expectedURL: "https://boundary.example.com/v1/roles?page_size=100&recursive=true&scope_id=global",
		},
		{
			name: "endpoint_with_filter",
			request: &hashicorp_adapter.Request{
				BaseURL:          "https://boundary.example.com",
				EntityExternalID: "hosts",
				PageSize:         50,
				EntityConfig: map[string]hashicorp_adapter.EntityConfig{
					"hosts": {
						ScopeID: "org_123",
						Filter:  "name eq \"test\"",
					},
				},
			},
			expectedURL: "https://boundary.example.com/v1/hosts?" +
				"filter=name+eq+%22test%22&page_size=50&recursive=true&scope_id=org_123",
		},
		{
			name: "endpoint_with_cursor",
			request: &hashicorp_adapter.Request{
				BaseURL:          "https://boundary.example.com",
				EntityExternalID: "users",
				PageSize:         25,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("next_page_token"),
				},
				EntityConfig: map[string]hashicorp_adapter.EntityConfig{
					"users": {
						ScopeID: "proj_456",
					},
				},
			},
			expectedURL: "https://boundary.example.com/v1/users?" +
				"list_token=next_page_token&page_size=25&recursive=true&scope_id=proj_456",
		},
		{
			name: "endpoint_with_additional_params",
			request: &hashicorp_adapter.Request{
				BaseURL:          "https://boundary.example.com",
				EntityExternalID: "groups",
				PageSize:         75,
				EntityConfig: map[string]hashicorp_adapter.EntityConfig{
					"groups": {
						ScopeID: "global",
					},
				},
				AdditionalParams: map[string]string{
					"include": "members",
					"sort":    "name",
				},
			},
			expectedURL: "https://boundary.example.com/v1/groups?" +
				"include=members&page_size=75&recursive=true&scope_id=global&sort=name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := hashicorp_adapter.ConstructEndpoint(tt.request)
			assert.Equal(t, tt.expectedURL, endpoint)
		})
	}
}
