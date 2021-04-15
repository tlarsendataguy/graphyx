package output

import (
	"fmt"
	"strings"
)

type NodeConfig struct {
	Label      string
	IdFields   []string
	PropFields []string
}

type RelationshipConfig struct {
	LeftLabel          string
	LeftAlteryxFields  []string
	LeftNeo4jFields    []string
	RightLabel         string
	RightAlteryxFields []string
	RightNeo4jFields   []string
	Label              string
	PropFields         []string
}

func NodeQuery(config *NodeConfig) string {
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
	onCreateSetQuery(builder, config.PropFields, `newNode`)

	return builder.String()
}

func RelationshipQuery(config *RelationshipConfig) string {
	builder := &strings.Builder{}
	builder.WriteString("UNWIND $batch AS row\n")
	matchNode(builder, config.LeftLabel, config.LeftAlteryxFields, config.LeftNeo4jFields, `left`)
	matchNode(builder, config.RightLabel, config.RightAlteryxFields, config.RightNeo4jFields, `right`)
	builder.WriteString(fmt.Sprintf("MERGE (left)-[newRel:`%v`]->(right)\n", config.Label))
	if len(config.PropFields) == 0 {
		return builder.String()
	}
	onCreateSetQuery(builder, config.PropFields, `newRel`)
	return builder.String()
}

func matchNode(builder *strings.Builder, label string, alteryxFields []string, neo4jFields []string, neo4jVariable string) {
	builder.WriteString(fmt.Sprintf("MATCH (%v:`%v`{", neo4jVariable, label))
	for index, neo4jId := range neo4jFields {
		neo4jId = escapeName(neo4jId)
		alteryxId := escapeName(alteryxFields[index])
		if index > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("`%v`:row.`%v`", neo4jId, alteryxId))
	}
	builder.WriteString("})\n")
}

func mergeNodeClause(builder *strings.Builder, config *NodeConfig) {
	builder.WriteString(fmt.Sprintf("MERGE (newNode:`%v`{", config.Label))
	for index, id := range config.IdFields {
		id = escapeName(id)
		if index > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("`%v`:row.`%v`", id, id))
	}
	builder.WriteString("})\n")
}

func createNodeClause(builder *strings.Builder, config *NodeConfig) {
	builder.WriteString(fmt.Sprintf("CREATE (newNode:`%v`{", config.Label))
	for index, id := range config.PropFields {
		id = escapeName(id)
		if index > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("`%v`:row.`%v`", id, id))
	}
	builder.WriteString("})")
}

func onCreateSetQuery(builder *strings.Builder, props []string, neo4jVariable string) {
	builder.WriteString("ON CREATE SET ")
	buildSetProperties(builder, props, neo4jVariable)
	builder.WriteString("\n")
	builder.WriteString("ON MATCH SET ")
	buildSetProperties(builder, props, neo4jVariable)
}

func buildSetProperties(builder *strings.Builder, props []string, neo4jVariable string) {
	for index, prop := range props {
		prop = escapeName(prop)
		if index > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("%v.`%v`=row.`%v`", neo4jVariable, prop, prop))
	}
}

func escapeName(name string) string {
	return strings.Replace(name, "`", "``", -1)
}
