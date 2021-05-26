package input

import (
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"strings"
)

func ToString(value interface{}) string {
	builder := &strings.Builder{}
	recursiveToString(value, builder)
	return builder.String()
}

func recursiveToString(value interface{}, builder *strings.Builder) {
	switch v := value.(type) {
	case neo4j.Node:
		builder.WriteByte('(')
		for _, label := range v.Labels {
			builder.WriteByte(':')
			builder.WriteString(label)
		}
		if len(v.Props) > 0 {
			builder.WriteByte(' ')
			jsonBytes, _ := json.Marshal(v.Props)
			builder.Write(jsonBytes)
		}
		builder.WriteByte(')')
	case neo4j.Relationship:
		builder.WriteByte('[')
		if v.Type != `` {
			builder.WriteByte(':')
			builder.WriteString(v.Type)
		}
		if len(v.Props) > 0 {
			builder.WriteByte(' ')
			jsonBytes, _ := json.Marshal(v.Props)
			builder.Write(jsonBytes)
		}
		builder.WriteByte(']')
	case neo4j.Path:
		nodeCount := len(v.Nodes)
		if nodeCount == 0 {
			return
		}
		recursiveToString(v.Nodes[0], builder)
		if nodeCount == 1 {
			return
		}

		for index, rel := range v.Relationships {
			builder.WriteByte('-')
			recursiveToString(rel, builder)
			builder.WriteByte('-')
			recursiveToString(v.Nodes[index+1], builder)
		}
	case string:
		builder.WriteString(v)
	default:
		builder.WriteString(fmt.Sprintf(`%v`, value))
	}
}
