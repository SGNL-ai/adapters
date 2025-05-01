package mysql

import "strings"

func ConstructQuery(request *Request) string {
	if request == nil {
		return ""
	}

	var sb strings.Builder

	sb.Grow(15 + 29 + 37 + len(request.EntityConfig.ExternalId) + len(request.UniqueAttributeExternalID))

	sb.WriteString("SELECT *, CAST(") // len=15
	sb.WriteString(request.UniqueAttributeExternalID)
	sb.WriteString(" as CHAR(50)) as str_id FROM ") // len=29
	sb.WriteString(request.EntityConfig.ExternalId)

	if request.Filter != nil && *request.Filter != "" {
		sb.Grow(7 + len(*request.Filter))

		sb.WriteString(" WHERE ") // len=7
		sb.WriteString(*request.Filter)
	}

	sb.WriteString(" ORDER BY str_id ASC LIMIT ? OFFSET ?") // len=37

	return sb.String()
}
