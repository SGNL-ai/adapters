// Package condexpr defines a simple DSL for building structured, nested filter expressions.
// It supports logical operators (AND, OR) and standard comparison operators (=, !=, >, <, etc.).
// The expressions can be evaluated directly or translated into other languages such as SQL.
//
// This package is designed for use across services that require programmatic filtering logic.
// It is not tied to a specific database or schema and can be extended to support additional
// operators or target languages.
//
// Example:
//
//	cond := Condition{
//		And: []Condition{
//			{
//				Field: "age",
//				Operator: ">",
//				Value: 18
//			},
//			{
//				Field: "status",
//				Operator: "=",
//				Value: "active"
//			},
//		},
//	}
package condexpr

type Condition struct {
	Field    string `json:"field,omitempty"`
	Operator string `json:"op,omitempty"`
	Value    any    `json:"value,omitempty"`

	// Nested conditions
	And []Condition `json:"and,omitempty"`
	Or  []Condition `json:"or,omitempty"`
}
