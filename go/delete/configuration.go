package delete

import "strings"

func GenerateDeleteNodes(label string, idFields []string) string {
	builder := &strings.Builder{}
	builder.WriteString("UNWIND $batch AS row\n")
	builder.WriteString("MATCH (d:`Customer`{`Key`:row.`Key`}) DETACH DELETE d")
	return builder.String()
}
