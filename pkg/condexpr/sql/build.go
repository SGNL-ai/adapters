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
//   - IS NULL: Checks if a field is null (no value required)
//   - IS NOT NULL: Checks if a field is not null (no value required)
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
//   - Requires all conditions to have a non-empty field and operator.
//     For most operators, a value is also required, except for IS NULL and IS NOT NULL.
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

	for idx, c := range cond.And {
		expr, err := cb.Build(c)
		if err != nil {
			return nil, fmt.Errorf("failed to build AND condition at index %d: %w", idx, err)
		}

		exprs = append(exprs, expr)
	}

	return goqu.And(exprs...), nil
}

func (cb ConditionBuilder) BuildCompositeOr(cond condexpr.Condition) (goqu.Expression, error) {
	exprs := make([]goqu.Expression, 0, len(cond.Or))

	for idx, c := range cond.Or {
		expr, err := cb.Build(c)
		if err != nil {
			return nil, fmt.Errorf("failed to build OR condition at index %d: %w", idx, err)
		}

		exprs = append(exprs, expr)
	}

	return goqu.Or(exprs...), nil
}

func (cb ConditionBuilder) BuildLeafCondition(cond condexpr.Condition) (goqu.Expression, error) {
	// Validate leaf condition
	if cond.Field == "" {
		return nil, errMissingField
	}

	if cond.Operator == "" {
		return nil, errMissingOperator
	}

	// For null operators, value should not be provided
	if cond.Operator == "IS NULL" || cond.Operator == "IS NOT NULL" {
		if cond.Value != nil {
			return nil, fmt.Errorf("value should not be provided for %s operator", cond.Operator)
		}
	} else {
		// For all other operators, value is required
		if cond.Value == nil {
			return nil, errMissingValue
		}
	}

	if valid := validSQLIdentifier.MatchString(cond.Field); !valid {
		return nil, fmt.Errorf(
			"field validation failed for '%s': unsupported characters found or length is not in range 1-128",
			cond.Field,
		)
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
	case "IS NULL":
		return col.IsNull(), nil
	case "IS NOT NULL":
		return col.IsNotNull(), nil
	default:
		return nil, fmt.Errorf("unsupported operator: %q", cond.Operator)
	}
}
