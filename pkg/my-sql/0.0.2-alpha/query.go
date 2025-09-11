// Copyright 2025 SGNL.ai, Inc.

package mysql

import (
	"errors"
	"math"

	"github.com/doug-martin/goqu/v9"
	"github.com/sgnl-ai/adapters/pkg/condexpr"
	condexprsql "github.com/sgnl-ai/adapters/pkg/condexpr/sql"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql" // goqu MySQL Dialect required for constructing correct queries.
)

func ConstructQuery(request *Request) (string, []any, error) {
	if request == nil {
		return "", nil, errors.New("nil request provided")
	}

	dialect := goqu.Dialect("mysql")

	expr := dialect.Select(
		"*",
		goqu.Cast(goqu.I(request.UniqueAttributeExternalID), "CHAR(50)").As("str_id"),
	).From(request.EntityConfig.ExternalId).Prepared(true)

	// Create a shared condition object for pagination and request filtering.
	var cond *condexpr.Condition

	if request.Cursor != nil && *request.Cursor != "" {
		cond = &condexpr.Condition{
			// We filter on the string casted id value to ensure that we're paginating
			// on the same field we're sorting on.
			Field:    "str_id",
			Operator: ">",
			Value:    request.Cursor,
		}
	}

	if request.Filter != nil {
		// If we already have a condition already set from pagination, create a new composite condition
		// with both.
		if cond != nil {
			cond = &condexpr.Condition{
				And: []condexpr.Condition{
					*cond,
					*request.Filter,
				},
			}
		} else {
			cond = request.Filter
		}
	}

	if cond != nil {
		builder := condexprsql.NewConditionBuilder()

		whereExpr, err := builder.Build(*cond)
		if err != nil {
			return "", nil, err
		}

		expr = expr.Where(whereExpr)
	}

	if request.PageSize < 0 {
		return "", nil, errors.New("invalid negative pageSize provided")
	}

	// MaxUint will either be equal to MaxUint32 or MaxUint64, depending on the system.
	// For consistency between systems, we'll assert that the cursor is less than MaxUint32.
	if request.PageSize > math.MaxUint32 {
		return "", nil, errors.New("pageSize value exceeds maximum allowed value")
	}

	expr = expr.Order(goqu.I("str_id").Asc()).Limit(uint(request.PageSize))

	return expr.ToSQL()
}
