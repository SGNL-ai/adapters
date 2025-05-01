// Copyright 2025 SGNL.ai, Inc.
package mysql_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/condexpr"
	mysql "github.com/sgnl-ai/adapters/pkg/my-sql"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

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
			wantQuery: "SELECT *, CAST(`id` AS CHAR(50)) AS `str_id` FROM `users` ORDER BY `str_id` ASC",
		},
		{
			name:         "nil_request",
			inputRequest: nil,
			wantErr:      errors.New("nil request provided"),
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
