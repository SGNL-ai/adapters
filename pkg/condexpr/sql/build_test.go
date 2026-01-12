// Copyright 2025 SGNL.ai, Inc.

package sql

import (
	"errors"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/sgnl-ai/adapters/pkg/condexpr"
	"github.com/stretchr/testify/assert"
)

// nolint: lll
func TestBuild(t *testing.T) {
	tests := map[string]struct {
		condition condexpr.Condition
		wantExpr  goqu.Expression
		wantError error
	}{
		"simple_eq_condition": {
			condition: condexpr.Condition{
				Field:    "status",
				Operator: "=",
				Value:    "active",
			},
			wantExpr: goqu.C("status").Eq("active"),
		},
		"simple_gt_condition": {
			condition: condexpr.Condition{
				Field:    "age",
				Operator: ">",
				Value:    30,
			},
			wantExpr: goqu.C("age").Gt(30),
		},
		"and_condition_w_valid_subconditions": {
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
		"or_condition_w_valid_subconditions": {
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
		"nested_and_or_conditions": {
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
		"in_op_with_valid_values": {
			condition: condexpr.Condition{
				Field:    "status",
				Operator: "IN",
				Value:    []string{"active", "inactive"},
			},
			wantExpr: goqu.C("status").In([]string{"active", "inactive"}),
		},
		"in_op_with_single_value": {
			condition: condexpr.Condition{
				Field:    "status",
				Operator: "IN",
				Value:    "active",
			},
			wantExpr: goqu.C("status").In("active"),
		},
		"unsupported_op": {
			condition: condexpr.Condition{
				Field:    "status",
				Operator: "LIKE",
				Value:    "%active%",
			},
			wantError: errors.New(`unsupported operator: "LIKE"`),
		},
		"missing_field": {
			condition: condexpr.Condition{
				Operator: "=",
				Value:    "active",
			},
			wantError: errors.New("missing required field"),
		},
		"missing_op": {
			condition: condexpr.Condition{
				Field: "status",
				Value: "active",
			},
			wantError: errors.New("missing required operator"),
		},
		"missing_value": {
			condition: condexpr.Condition{
				Field:    "status",
				Operator: "=",
			},
			wantError: errors.New("missing required value"),
		},
		"invalid_op": {
			condition: condexpr.Condition{
				Field:    "status",
				Operator: "INVALID",
				Value:    "active",
			},
			wantError: errors.New(`unsupported operator: "INVALID"`),
		},
		"and_condition_with_invalid_subconditions": {
			condition: condexpr.Condition{
				And: []condexpr.Condition{
					{Field: "age", Operator: ">", Value: 30},
					{Operator: "=", Value: "John"},
				},
			},
			wantError: errors.New("failed to build AND condition at index 1: missing required field"),
		},
		"or_condition_with_invalid_subconditions": {
			condition: condexpr.Condition{
				Or: []condexpr.Condition{
					{Field: "age", Operator: "<", Value: 20},
					{Field: "name"},
				},
			},
			wantError: errors.New("failed to build OR condition at index 1: missing required operator"),
		},
		"invalid_chars_in_field": {
			condition: condexpr.Condition{
				Field:    "id; DROP TABLE Users;",
				Operator: "=",
				Value:    "active",
			},
			wantError: errors.New("field validation failed for 'id; DROP TABLE Users;': unsupported characters found or length is not in range 1-128"),
		},
		"is_null_condition": {
			condition: condexpr.Condition{
				Field:    "email",
				Operator: "IS NULL",
			},
			wantExpr: goqu.C("email").IsNull(),
		},
		"is_not_null_condition": {
			condition: condexpr.Condition{
				Field:    "phone",
				Operator: "IS NOT NULL",
			},
			wantExpr: goqu.C("phone").IsNotNull(),
		},
		"is_null_with_value_should_error": {
			condition: condexpr.Condition{
				Field:    "email",
				Operator: "IS NULL",
				Value:    "something",
			},
			wantError: errors.New("value should not be provided for IS NULL operator"),
		},
		"is_not_null_with_value_should_error": {
			condition: condexpr.Condition{
				Field:    "phone",
				Operator: "IS NOT NULL",
				Value:    "something",
			},
			wantError: errors.New("value should not be provided for IS NOT NULL operator"),
		},
		"invalid_condition": {
			condition: condexpr.Condition{
				Field:    "age",
				Operator: ">",
				Value:    21,
				Or: []condexpr.Condition{
					{Field: "age", Operator: ">", Value: 21},
					{Field: "verified", Operator: "=", Value: true},
				},
				And: []condexpr.Condition{
					{Field: "age", Operator: ">", Value: 21},
					{Field: "verified", Operator: "=", Value: true},
				},
			},
			wantError: errors.New("invalid condition: specify exactly one of And, Or, or a valid leaf condition"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			builder := NewConditionBuilder()
			gotExpr, gotErr := builder.Build(tt.condition)

			if tt.wantError != nil {
				assert.EqualError(t, gotErr, tt.wantError.Error())
			} else {
				assert.NoError(t, gotErr)
			}

			assert.Equal(t, tt.wantExpr, gotExpr)
		})
	}
}
