package output

import (
	"errors"
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

func NodeQuery(config *NodeConfig) (string, error) {
	if config.Label == `` {
		return ``, errors.New(`label cannot be blank`)
	}
	builder := &strings.Builder{}
	builder.WriteString("UNWIND $batch AS row\n")
	if len(config.IdFields) == 0 {
		createNodeClause(builder, config)
		return builder.String(), nil
	}
	mergeNodeClause(builder, config)
	if len(config.PropFields) == 0 {
		return builder.String(), nil
	}
	onCreateSetQuery(builder, config.PropFields, `newNode`)

	return builder.String(), nil
}

func RelationshipQuery(config *RelationshipConfig) (string, error) {
	if config.Label == `` {
		return ``, errors.New(`label cannot be blank`)
	}
	if config.LeftLabel == `` {
		return ``, errors.New(`left node label cannot be blank`)
	}
	if config.RightLabel == `` {
		return ``, errors.New(`right node label cannot be blank`)
	}
	builder := &strings.Builder{}
	builder.WriteString("UNWIND $batch AS row\n")
	matchNode(builder, escapeName(config.LeftLabel), config.LeftAlteryxFields, config.LeftNeo4jFields, `left`)
	matchNode(builder, escapeName(config.RightLabel), config.RightAlteryxFields, config.RightNeo4jFields, `right`)
	builder.WriteString(fmt.Sprintf("MERGE (left)-[newRel:`%v`]->(right)\n", escapeName(config.Label)))
	if len(config.PropFields) == 0 {
		return builder.String(), nil
	}
	onCreateSetQuery(builder, config.PropFields, `newRel`)
	return builder.String(), nil
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
	label := escapeName(config.Label)
	builder.WriteString(fmt.Sprintf("MERGE (newNode:`%v`{", label))
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
	label := escapeName(config.Label)
	builder.WriteString(fmt.Sprintf("CREATE (newNode:`%v`{", label))
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
