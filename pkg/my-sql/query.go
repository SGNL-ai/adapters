package mysql

import "strings"

func ConstructQuery(request *Request) string {
	if request == nil {
		return ""
	}

	var sb strings.Builder

	sb.Grow(45 + 37 + len(request.EntityConfig.ExternalId))

	sb.WriteString("SELECT *, CAST(? as CHAR(50)) as str_id FROM ") // len=45
	sb.WriteString(request.EntityConfig.ExternalId)

	if request.Filter != nil && *request.Filter != "" {
		sb.Grow(7 + len(*request.Filter))

		sb.WriteString(" WHERE ") // len=7
		sb.WriteString(*request.Filter)
	}

	sb.WriteString(" ORDER BY str_id ASC LIMIT ? OFFSET ?") // len=37

	return sb.String()
}
