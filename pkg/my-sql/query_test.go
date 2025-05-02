// Copyright 2025 SGNL.ai, Inc.
package mysql_test

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/condexpr"
	mysql "github.com/sgnl-ai/adapters/pkg/my-sql"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

// nolint: lll
func TestConstructQuery(t *testing.T) {
	tests := []struct {
		name         string
		inputRequest *mysql.Request

		wantQuery string
		wantAttrs []any
		wantErr   error
	}{
		{
			name: "simple",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
				},
				Filter:                    nil,
				UniqueAttributeExternalID: "id",
				PageSize:                  100,
				Cursor:                    testutil.GenPtr(int64(500)),
			},
			wantQuery: "SELECT *, CAST(`id` AS CHAR(50)) AS `str_id` FROM `users` ORDER BY `str_id` ASC LIMIT 100 OFFSET 500",
			wantAttrs: []any{},
		},
		{
			name: "simple_with_filter",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "groups",
				},
				Filter: testutil.GenPtr(condexpr.Condition{
					Field:    "status",
					Operator: "=",
					Value:    "active",
				}),
				UniqueAttributeExternalID: "groupId",
				PageSize:                  100,
			},
			wantQuery: "SELECT *, CAST(`groupId` AS CHAR(50)) AS `str_id` FROM `groups` WHERE (`status` = 'active') ORDER BY `str_id` ASC LIMIT 100",
			wantAttrs: []any{},
		},
		{
			name: "simple_with_complex_filter",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
				},
				Filter: testutil.GenPtr(condexpr.Condition{
					Or: []condexpr.Condition{
						{
							And: []condexpr.Condition{
								{
									Field:    "age",
									Operator: ">",
									Value:    18,
								},
								{
									Field:    "country",
									Operator: "=",
									Value:    "USA",
								},
							},
						},
						{
							Field:    "verified",
							Operator: "=",
							Value:    true,
						},
					},
				}),
				UniqueAttributeExternalID: "id",
			},
			wantQuery: "SELECT *, CAST(`id` AS CHAR(50)) AS `str_id` FROM `users` WHERE (((`age` > 18) AND (`country` = 'USA')) OR (`verified` IS TRUE)) ORDER BY `str_id` ASC",
			wantAttrs: []any{},
		},
		{
			name: "query_empty_filter",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
				},
				Filter:                    testutil.GenPtr(condexpr.Condition{}),
				UniqueAttributeExternalID: "id",
			},
			wantErr: errors.New("invalid condition: specify exactly one of And, Or, or a valid leaf condition"),
		},
		{
			name: "invalid_condition_structure",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "groups",
				},
				Filter: testutil.GenPtr(condexpr.Condition{
					Field:    "status",
					Operator: "=",
					Value:    "active",
					And: []condexpr.Condition{
						{
							Field:    "status",
							Operator: "=",
							Value:    "active",
						},
					},
				}),
				UniqueAttributeExternalID: "groupId",
				PageSize:                  100,
			},
			wantErr: errors.New("invalid condition: specify exactly one of And, Or, or a valid leaf condition"),
		},
		{
			name:         "nil_request",
			inputRequest: nil,
			wantErr:      errors.New("nil request provided"),
		},
		{
			name: "valid_large_cursor",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
				},
				Filter:                    nil,
				UniqueAttributeExternalID: "id",
				PageSize:                  100,
				Cursor:                    testutil.GenPtr(int64(math.MaxUint32)),
			},
			wantQuery: "SELECT *, CAST(`id` AS CHAR(50)) AS `str_id` FROM `users` ORDER BY `str_id` ASC LIMIT 100 OFFSET 4294967295",
			wantAttrs: []any{},
		},
		{
			name: "invalid_cursor_too_large",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
				},
				Filter:                    nil,
				UniqueAttributeExternalID: "id",
				PageSize:                  100,
				Cursor:                    testutil.GenPtr(int64(math.MaxUint32 + 1)),
			},
			wantErr: errors.New("cursor value exceeds maximum allowed value"),
		},
		{
			name: "invalid_cursor_negative",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
				},
				Filter:                    nil,
				UniqueAttributeExternalID: "id",
				PageSize:                  100,
				Cursor:                    testutil.GenPtr(int64(-1)),
			},
			wantErr: errors.New("invalid negative cursor provided"),
		},
		{
			name: "invalid_page_size_too_large",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
				},
				Filter:                    nil,
				UniqueAttributeExternalID: "id",
				PageSize:                  int64(math.MaxUint32 + 1),
			},
			wantErr: errors.New("pageSize value exceeds maximum allowed value"),
		},
		{
			name: "invalid_cursor_negative",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
				},
				Filter:                    nil,
				UniqueAttributeExternalID: "id",
				PageSize:                  int64(-1),
			},
			wantErr: errors.New("invalid negative pageSize provided"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotAttrs, gotErr := mysql.ConstructQuery(tt.inputRequest)

			assert.Equal(t, tt.wantQuery, gotQuery)
			assert.Equal(t, tt.wantAttrs, gotAttrs)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
