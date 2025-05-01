package mysql

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/doug-martin/goqu/v9"

	"github.com/sgnl-ai/adapters/pkg/condexpr"
)

var validSQLIdentifier = regexp.MustCompile(`^[a-zA-Z0-9$_]{1,128}$`)

// Build recursively builds a goqu expression for the provided condexpr.Condition.
func Build(cond condexpr.Condition) (goqu.Expression, error) {
	// Handle AND condition
	if len(cond.And) > 0 {
		andExprs := []goqu.Expression{}
		for _, c := range cond.And {
			expr, err := Build(c)
			if err != nil {
				return nil, err
			}
			andExprs = append(andExprs, expr)
		}
		return goqu.And(andExprs...), nil
	}

	// Handle OR condition
	if len(cond.Or) > 0 {
		orExprs := []goqu.Expression{}
		for _, c := range cond.Or {
			expr, err := Build(c)
			if err != nil {
				return nil, err
			}
			orExprs = append(orExprs, expr)
		}
		return goqu.Or(orExprs...), nil
	}

	if cond.Value == nil {
		return nil, errors.New("missing required value")
	}

	if cond.Field == "" {
		return nil, errors.New("missing required field")
	}

	if cond.Operator == "" {
		return nil, errors.New("missing required operator")
	}

	if valid := validSQLIdentifier.MatchString(cond.Field); !valid {
		return nil, fmt.Errorf("field validation failed: unsupported characters found or length is not in range 1-128")
	}

	// Base condition (Field, Operator, Value)
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
		return nil, fmt.Errorf("unsupported operator: %s", cond.Operator)
	}
}
