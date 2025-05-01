package mysql

import (
	"errors"

	"github.com/doug-martin/goqu/v9"
	condexprsql "github.com/sgnl-ai/adapters/pkg/condexpr/sql"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

func ConstructQuery(request *Request) (string, []any, error) {
	if request == nil {
		return "", nil, errors.New("nil request provided")
	}

	dialect := goqu.Dialect("mysql")

	expr := dialect.Select("*", goqu.Cast(goqu.I(request.UniqueAttributeExternalID), "CHAR(50)").As("str_id")).From(request.EntityConfig.ExternalId)

	if request.Filter != nil {
		whereExpr, err := condexprsql.Build(*request.Filter)
		if err != nil {
			return "", nil, err
		}

		expr = expr.Where(whereExpr)
	}

	expr = expr.Order(goqu.I("str_id").Asc()).Limit(uint(request.PageSize))

	if request.Cursor != nil {
		expr = expr.Offset(uint(*request.Cursor))
	}

	return expr.ToSQL()
}
