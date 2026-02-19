// Copyright 2026 SGNL.ai, Inc.

// Package condexpr defines a simple DSL for building structured, nested filter expressions.
// It supports logical operators (AND, OR) and standard comparison operators (=, !=, >, <, etc.).
// The expressions can be evaluated directly or translated into other languages such as SQL.
//
// This package is designed for use across services that require programmatic filtering logic.
// It is not tied to a specific database or schema and can be extended to support additional
// operators or target languages.
package condexpr

import "errors"

// Condition represents a filter expression that can be either a leaf condition
// or a logical composite condition.
type Condition struct {
	// Field is the name of the field to compare. Only used for leaf conditions.
	Field string `json:"field,omitempty"`

	// Operator is the comparison operator. Only used for leaf conditions.
	Operator string `json:"op,omitempty"`

	// Value is the value to compare against. Only used for leaf conditions.
	// The type should match the field type being compared.
	Value any `json:"value,omitempty"`

	// And combines multiple conditions with logical AND.
	// All conditions must be true for the AND to be true.
	And []Condition `json:"and,omitempty"`

	// Or combines multiple conditions with logical OR.
	// At least one condition must be true for the OR to be true.
	Or []Condition `json:"or,omitempty"`
}

// ConditionBuilder is an interface that can be implemented for specific use cases (e.g. to add support for
// building a SQL expression, etc).
type ConditionBuilder[T any] interface {
	Build(cond Condition) (T, error)
	BuildCompositeAnd(cond Condition) (T, error)
	BuildCompositeOr(cond Condition) (T, error)
	BuildLeafCondition(cond Condition) (T, error)
}

// DefaultBuild provides generic default processing and validation functionality for types implementing
// the ConditionBuilder.Build function.
//
// It processes a given Condition and determines whether it represents a composite AND condition,
// a composite OR condition, or a valid leaf condition. Based on the type, it delegates the construction
// to the appropriate method of the ConditionBuilder.
//
// The function also enforces that the Condition must specify exactly one of the following:
// - An AND condition (non-empty `And` field).
// - An OR condition (non-empty `Or` field).
// - A valid leaf condition (non-empty `Field`, `Operator`, or `Value` fields).
func DefaultBuild[T any, CB ConditionBuilder[T]](cb CB, cond Condition) (out T, err error) {
	// Validate that the condition specifies only one field: And, Or, or a valid leaf condition
	isAnd := len(cond.And) > 0
	isOr := len(cond.Or) > 0
	isLeaf := cond.Field != "" || cond.Operator != "" || cond.Value != nil

	if (isAnd && isOr) || (isAnd && isLeaf) || (isOr && isLeaf) || (!isAnd && !isOr && !isLeaf) {
		err = errors.New("invalid condition: specify exactly one of And, Or, or a valid leaf condition")

		return
	}

	if isAnd {
		return cb.BuildCompositeAnd(cond)
	}

	if isOr {
		return cb.BuildCompositeOr(cond)
	}

	return cb.BuildLeafCondition(cond)
}
