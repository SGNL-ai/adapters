// Copyright 2025 SGNL.ai, Inc.
package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	framework "github.com/sgnl-ai/adapter-framework"
	mysql "github.com/sgnl-ai/adapters/pkg/my-sql"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestConstructQuery(t *testing.T) {
	tests := []struct {
		name         string
		inputRequest *mysql.Request
		wantQuery    string
	}{
		{
			name: "simple",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
				},
				Filter: nil,
			},
			wantQuery: "SELECT *, CAST(? as CHAR(50)) as str_id FROM users ORDER BY str_id ASC LIMIT ? OFFSET ?",
		},
		{
			name: "simple_with_filter",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "groups",
				},
				Filter: testutil.GenPtr("status = 'active'"),
			},
			wantQuery: "SELECT *, CAST(? as CHAR(50)) as str_id FROM groups WHERE status = 'active' ORDER BY str_id ASC LIMIT ? OFFSET ?",
		},
		{
			name: "simple_with_complex_filter",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
				},
				Filter: testutil.GenPtr("(age > 18 AND country = 'USA') OR verified = TRUE"),
			},
			wantQuery: "SELECT *, CAST(? as CHAR(50)) as str_id FROM users WHERE (age > 18 AND country = 'USA') OR verified = TRUE ORDER BY str_id ASC LIMIT ? OFFSET ?",
		},
		{
			name: "query_empty_filter",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
				},
				Filter: testutil.GenPtr(""),
			},
			wantQuery: "SELECT *, CAST(? as CHAR(50)) as str_id FROM users ORDER BY str_id ASC LIMIT ? OFFSET ?",
		},
		{
			name:         "nil_request",
			inputRequest: nil,
			wantQuery:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery := mysql.ConstructQuery(tt.inputRequest)
			assert.Equal(t, tt.wantQuery, gotQuery)
		})
	}
}
