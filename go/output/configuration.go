package output

import (
	"fmt"
	"strings"
)

type ConfigStruct struct {
	Label      string
	IdFields   []string
	PropFields []string
}

func mergeNodeClause(builder *strings.Builder, config *ConfigStruct) {
	builder.WriteString(fmt.Sprintf("MERGE (newNode:`%v`{", config.Label))
	for index, id := range config.IdFields {
		if index > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("`%v`:row.`%v`", id, id))
	}
	builder.WriteString("})\n")
}

func createNodeClause(builder *strings.Builder, config *ConfigStruct) {
	builder.WriteString(fmt.Sprintf("CREATE (newNode:`%v`{", config.Label))
	for index, id := range config.PropFields {
		if index > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("`%v`:row.`%v`", id, id))
	}
	builder.WriteString("})")
}

func buildSetProperties(builder *strings.Builder, config *ConfigStruct) {
	for index, prop := range config.PropFields {
		if index > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("newNode.`%v`=row.`%v`", prop, prop))
	}
}

func NodeQuery(config *ConfigStruct) string {
	builder := &strings.Builder{}
	builder.WriteString("UNWIND $batch AS row\n")
	if len(config.IdFields) == 0 {
		createNodeClause(builder, config)
		return builder.String()
	}
	mergeNodeClause(builder, config)
	if len(config.PropFields) == 0 {
		return builder.String()
	}
	builder.WriteString("ON CREATE SET ")
	buildSetProperties(builder, config)
	builder.WriteString("\n")
	builder.WriteString("ON MATCH SET ")
	buildSetProperties(builder, config)

	return builder.String()
}
