package mysql

import (
	"errors"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/sgnl-ai/adapters/pkg/condexpr"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	tests := []struct {
		name      string
		condition condexpr.Condition
		wantExpr  goqu.Expression
		wantError error
	}{
		{
			name: "Simple equality condition",
			condition: condexpr.Condition{
				Field:    "status",
				Operator: "=",
				Value:    "active",
			},
			wantExpr: goqu.C("status").Eq("active"),
		},
		{
			name: "Simple greater than condition",
			condition: condexpr.Condition{
				Field:    "age",
				Operator: ">",
				Value:    30,
			},
			wantExpr: goqu.C("age").Gt(30),
		},
		{
			name: "AND condition with valid sub-conditions",
			condition: condexpr.Condition{
				And: []condexpr.Condition{
					{Field: "age", Operator: ">", Value: 30},
					{Field: "status", Operator: "=", Value: "active"},
				},
			},
			wantExpr: goqu.And(
				goqu.C("age").Gt(30),
				goqu.C("status").Eq("active"),
			),
		},
		{
			name: "OR condition with valid sub-conditions",
			condition: condexpr.Condition{
				Or: []condexpr.Condition{
					{Field: "age", Operator: "<", Value: 20},
					{Field: "status", Operator: "=", Value: "inactive"},
				},
			},
			wantExpr: goqu.Or(
				goqu.C("age").Lt(20),
				goqu.C("status").Eq("inactive"),
			),
		},
		{
			name: "Nested AND/OR condition",
			condition: condexpr.Condition{
				And: []condexpr.Condition{
					{
						Or: []condexpr.Condition{
							{Field: "age", Operator: ">", Value: 30},
							{Field: "status", Operator: "=", Value: "active"},
						},
					},
					{Field: "country", Operator: "=", Value: "US"},
				},
			},
			wantExpr: goqu.And(
				goqu.Or(
					goqu.C("age").Gt(30),
					goqu.C("status").Eq("active"),
				),
				goqu.C("country").Eq("US"),
			),
		},
		{
			name: "IN operator with valid values",
			condition: condexpr.Condition{
				Field:    "status",
				Operator: "IN",
				Value:    []string{"active", "inactive"},
			},
			wantExpr: goqu.C("status").In([]string{"active", "inactive"}),
		},
		{
			name: "Unsupported operator",
			condition: condexpr.Condition{
				Field:    "status",
				Operator: "LIKE",
				Value:    "%active%",
			},
			wantError: errors.New("unsupported operator: LIKE"),
		},
		{
			name: "Missing field",
			condition: condexpr.Condition{
				Operator: "=",
				Value:    "active",
			},
			wantError: errors.New("missing required field"),
		},
		{
			name: "Missing operator",
			condition: condexpr.Condition{
				Field: "status",
				Value: "active",
			},
			wantError: errors.New("missing required operator"),
		},
		{
			name: "Missing value",
			condition: condexpr.Condition{
				Field:    "status",
				Operator: "=",
			},
			wantError: errors.New("missing required value"),
		},
		{
			name: "Invalid operator",
			condition: condexpr.Condition{
				Field:    "status",
				Operator: "INVALID",
				Value:    "active",
			},
			wantError: errors.New("unsupported operator: INVALID"),
		},
		{
			name: "AND condition with invalid sub-condition",
			condition: condexpr.Condition{
				And: []condexpr.Condition{
					{Field: "age", Operator: ">", Value: 30},
					{Operator: "=", Value: "John"},
				},
			},
			wantError: errors.New("missing required field"),
		},
		{
			name: "OR condition with invalid sub-condition",
			condition: condexpr.Condition{
				Or: []condexpr.Condition{
					{Field: "age", Operator: "<", Value: 20},
					{Field: "name"},
				},
			},
			wantError: errors.New("missing required value"),
		},
		{
			name: "Invalid chars in field",
			condition: condexpr.Condition{
				Field:    "id; DROP TABLE Users;",
				Operator: "=",
				Value:    "active",
			},
			wantError: errors.New("field validation failed: unsupported characters found or length is not in range 1-128"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExpr, gotErr := Build(tt.condition)

			assert.Equal(t, tt.wantError, gotErr)
			assert.Equal(t, tt.wantExpr, gotExpr)
		})
	}
}
