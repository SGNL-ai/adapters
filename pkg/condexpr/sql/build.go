// Copyright 2025 SGNL.ai, Inc.

// Package sql provides utilities for building SQL expressions from conditional expressions.
// This package is designed to work with the goqu library to construct SQL queries dynamically.
//
// Supported Features:
// - Logical Operators:
//   - AND: Combines multiple conditions with a logical AND.
//   - OR: Combines multiple conditions with a logical OR.
//
// - Comparison Operators:
//   - = : Equal to
//   - !=: Not equal to
//   - > : Greater than
//   - < : Less than
//   - >=: Greater than or equal to
//   - <=: Less than or equal to
//   - IN: Checks if a value exists within a list of values
//
// Limitations:
// - Only supports fields that are valid SQL identifiers:
//   - Must contain only alphanumeric characters, `$`, and `_`.
//   - Must be between 1 and 128 characters in length.
//
// - Does not support advanced SQL features such as:
//   - LIKE or ILIKE for pattern matching.
//   - NOT IN or other negated set operations.
//   - Complex expressions involving functions or subqueries.
//   - Requires all conditions to have a non-empty field, operator, and value.
//     Missing any of these will result in an error.
package sql

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/doug-martin/goqu/v9"

	"github.com/sgnl-ai/adapters/pkg/condexpr"
)

var (
	// validSQLIdentifier checks if a string is a valid SQL identifier:
	// - Contains only alphanumeric characters, $ and _.
	// - Length between 1-128 characters.
	validSQLIdentifier = regexp.MustCompile(`^[a-zA-Z0-9$_]{1,128}$`)

	errMissingValue    = errors.New("missing required value")
	errMissingField    = errors.New("missing required field")
	errMissingOperator = errors.New("missing required operator")
)

type ConditionBuilder struct{}

func NewConditionBuilder() *ConditionBuilder {
	return &ConditionBuilder{}
}

func (cb ConditionBuilder) Build(cond condexpr.Condition) (goqu.Expression, error) {
	return condexpr.DefaultBuild(cb, cond)
}

func (cb ConditionBuilder) BuildCompositeAnd(cond condexpr.Condition) (goqu.Expression, error) {
	exprs := make([]goqu.Expression, 0, len(cond.And))

	for _, c := range cond.And {
		expr, err := cb.Build(c)
		if err != nil {
			return nil, fmt.Errorf("failed to build AND condition: %w", err)
		}

		exprs = append(exprs, expr)
	}

	return goqu.And(exprs...), nil
}

func (cb ConditionBuilder) BuildCompositeOr(cond condexpr.Condition) (goqu.Expression, error) {
	exprs := make([]goqu.Expression, 0, len(cond.Or))

	for _, c := range cond.Or {
		expr, err := cb.Build(c)
		if err != nil {
			return nil, fmt.Errorf("failed to build OR condition: %w", err)
		}

		exprs = append(exprs, expr)
	}

	return goqu.Or(exprs...), nil
}

func (cb ConditionBuilder) BuildLeafCondition(cond condexpr.Condition) (goqu.Expression, error) {
	// Validate leaf condition
	if cond.Value == nil {
		return nil, errMissingValue
	}

	if cond.Field == "" {
		return nil, errMissingField
	}

	if cond.Operator == "" {
		return nil, errMissingOperator
	}

	if valid := validSQLIdentifier.MatchString(cond.Field); !valid {
		return nil, fmt.Errorf("field validation failed: unsupported characters found or length is not in range 1-128")
	}

	// Build leaf condition (Field, Operator, Value)
	col := goqu.C(cond.Field)

	switch cond.Operator {
	case "=":
		return col.Eq(cond.Value), nil
	case "!=":
		return col.Neq(cond.Value), nil
	case ">":
		return col.Gt(cond.Value), nil
	case "<":
		return col.Lt(cond.Value), nil
	case ">=":
		return col.Gte(cond.Value), nil
	case "<=":
		return col.Lte(cond.Value), nil
	case "IN":
		return col.In(cond.Value), nil
	default:
		return nil, fmt.Errorf("unsupported operator: %q", cond.Operator)
	}
}
