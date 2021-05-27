package delete

import "strings"

func GenerateDeleteNodes(label string, idFields []string) string {
	builder := &strings.Builder{}
	builder.WriteString("UNWIND $batch AS row\n")
	builder.WriteString("MATCH (d:`")
	builder.WriteString(escapeName(label))
	builder.WriteString("`{")
	for index, idField := range idFields {
		if index > 0 {
			builder.WriteByte(',')
		}
		escaped := escapeName(idField)
		builder.WriteByte('`')
		builder.WriteString(escaped)
		builder.WriteString("`:row.`")
		builder.WriteString(escaped)
		builder.WriteByte('`')
	}
	builder.WriteString("}) DETACH DELETE d")
	return builder.String()
}

func escapeName(name string) string {
	return strings.Replace(name, "`", "``", -1)
}
