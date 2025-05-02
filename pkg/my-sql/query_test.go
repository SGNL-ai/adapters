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
		{
			name: "filter_all_types_anded",
			inputRequest: &mysql.Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "groups",
				},
				Filter: testutil.GenPtr(condexpr.Condition{
					And: []condexpr.Condition{
						{
							Field:    "age",
							Operator: ">",
							Value:    21,
						},
						{
							Field:    "age",
							Operator: "<",
							Value:    100,
						},
						{
							Field:    "balance",
							Operator: ">=",
							Value:    1.234,
						},
						{
							Field:    "balance",
							Operator: "<=",
							Value:    123.4,
						},
						{
							Field:    "status",
							Operator: "=",
							Value:    "active",
						},
						{
							Field:    "riskScore",
							Operator: "!=",
							Value:    "HIGH",
						},
						{
							Field:    "verified",
							Operator: "=",
							Value:    true,
						},
						{
							Field:    "enabled",
							Operator: "!=",
							Value:    false,
						},
						{
							Field:    "countryCode",
							Operator: "IN",
							Value:    []string{"US", "CA", "MX"},
						},
						{
							Field:    "assignedCases",
							Operator: "IN",
							Value:    []int{1, 2, 3},
						},
						{
							Field:    "status",
							Operator: "IN",
							Value:    "enabled",
						},
					},
				}),
				UniqueAttributeExternalID: "groupId",
				PageSize:                  100,
			},
			wantQuery: "SELECT *, CAST(`groupId` AS CHAR(50)) AS `str_id` FROM `groups` WHERE ((`age` > 21) AND (`age` < 100) AND (`balance` >= 1.234) AND (`balance` <= 123.4) AND (`status` = 'active') AND (`riskScore` != 'HIGH') AND (`verified` IS TRUE) AND (`enabled` IS NOT FALSE) AND (`countryCode` IN ('US', 'CA', 'MX')) AND (`assignedCases` IN (1, 2, 3)) AND (`status` IN ('enabled'))) ORDER BY `str_id` ASC LIMIT 100",
			wantAttrs: []any{},
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
