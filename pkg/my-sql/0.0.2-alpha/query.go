// Copyright 2026 SGNL.ai, Inc.

package mysql

import (
	"errors"
	"math"

	"github.com/doug-martin/goqu/v9"
	condexprsql "github.com/sgnl-ai/adapters/pkg/condexpr/sql"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql" // goqu MySQL Dialect required for constructing correct queries.
)

// TotalRemainingRowsColumn is the column name for the COUNT window function result
// that provides the total count of rows matching the query conditions.
const TotalRemainingRowsColumn = "total_remaining_rows"

func ConstructQuery(request *Request) (string, []any, error) {
	if request == nil {
		return "", nil, errors.New("nil request provided")
	}

	dialect := goqu.Dialect("mysql")

	expr := dialect.Select(
		"*",
		goqu.Cast(goqu.I(request.UniqueAttributeExternalID), "CHAR(50)").As("str_id"),
		goqu.L("COUNT(*) OVER()").As(TotalRemainingRowsColumn),
	).From(request.EntityConfig.ExternalId).Prepared(true)

	if request.Cursor != nil && *request.Cursor != "" {
		expr = expr.Where(goqu.Cast(goqu.I(request.UniqueAttributeExternalID), "CHAR(50)").Gt(*request.Cursor))
	}

	if request.Filter != nil {
		builder := condexprsql.NewConditionBuilder()

		whereExpr, err := builder.Build(*request.Filter)
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
