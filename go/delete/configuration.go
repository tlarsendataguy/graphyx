package delete

import "strings"

type DeleteNodesProperties struct {
	Label    string
	IdFields []string
}

func GenerateDeleteNodes(props DeleteNodesProperties) string {
	builder := &strings.Builder{}
	builder.WriteString("UNWIND $batch AS row\n")
	builder.WriteString("MATCH (d:`")
	builder.WriteString(escapeName(props.Label))
	builder.WriteString("`{")
	for index, idField := range props.IdFields {
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
