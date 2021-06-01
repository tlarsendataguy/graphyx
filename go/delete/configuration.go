package delete

import (
	"errors"
	"strings"
)

type DeleteNodesProperties struct {
	Label    string
	IdFields []string
}

func (p *DeleteNodesProperties) escape() {
	p.Label = escapeName(p.Label)
	for index, idField := range p.IdFields {
		p.IdFields[index] = escapeName(idField)
	}
}

func GenerateDeleteNodes(props *DeleteNodesProperties) string {
	props.escape()
	builder := &strings.Builder{}
	builder.WriteString("UNWIND $batch AS row\n")
	builder.WriteString("MATCH (d")
	if props.Label != `` {
		writeLabel(builder, props.Label)
	}
	if len(props.IdFields) > 0 {
		writeProperties(builder, props.IdFields, props.IdFields)
	}
	builder.WriteString(") DETACH DELETE d")
	return builder.String()
}

type DeleteRelationshipsProperties struct {
	RelType                string
	RelFields              []string
	LeftNodeLabel          string
	LeftNodeAlteryxFields  []string
	LeftNodeNeo4jFields    []string
	RightNodeLabel         string
	RightNodeAlteryxFields []string
	RightNodeNeo4jFields   []string
}

func (p *DeleteRelationshipsProperties) escape() {
	p.RelType = escapeName(p.RelType)
	p.LeftNodeLabel = escapeName(p.LeftNodeLabel)
	p.RightNodeLabel = escapeName(p.RightNodeLabel)
	for _, strList := range [][]string{p.RelFields, p.LeftNodeAlteryxFields, p.LeftNodeNeo4jFields, p.RightNodeAlteryxFields, p.RightNodeNeo4jFields} {
		for index, item := range strList {
			strList[index] = escapeName(item)
		}
	}
}

func GenerateDeleteRelationships(props *DeleteRelationshipsProperties) (string, error) {
	if len(props.LeftNodeAlteryxFields) != len(props.LeftNodeNeo4jFields) {
		return ``, errors.New(`the number of left node Alteryx fields does not match the number of left node Neo4j fields`)
	}
	if len(props.RightNodeAlteryxFields) != len(props.RightNodeNeo4jFields) {
		return ``, errors.New(`the number of right node Alteryx fields does not match the number of right node Neo4j fields`)
	}

	props.escape()
	builder := &strings.Builder{}
	builder.WriteString("UNWIND $batch AS row\n")
	builder.WriteString("MATCH (")
	if props.LeftNodeLabel != `` {
		writeLabel(builder, props.LeftNodeLabel)
	}
	if len(props.LeftNodeNeo4jFields) > 0 {
		writeProperties(builder, props.LeftNodeNeo4jFields, props.LeftNodeAlteryxFields)
	}
	builder.WriteString(")-[r")
	if props.RelType != `` {
		writeLabel(builder, props.RelType)
	}
	if len(props.RelFields) > 0 {
		writeProperties(builder, props.RelFields, props.RelFields)
	}
	builder.WriteString("]-(")
	if props.RightNodeLabel != `` {
		writeLabel(builder, props.RightNodeLabel)
	}
	if len(props.RightNodeNeo4jFields) > 0 {
		writeProperties(builder, props.RightNodeNeo4jFields, props.RightNodeAlteryxFields)
	}
	builder.WriteString(") DELETE r")
	return builder.String(), nil
}

func escapeName(name string) string {
	return strings.Replace(name, "`", "``", -1)
}

func writeLabel(builder *strings.Builder, label string) {
	builder.WriteString(":`")
	builder.WriteString(label)
	builder.WriteByte('`')
}

func writeProperties(builder *strings.Builder, neo4jFields []string, alteryxFields []string) {
	builder.WriteString(" {")
	for index, neo4jKey := range neo4jFields {
		if index > 0 {
			builder.WriteByte(',')
		}
		builder.WriteByte('`')
		builder.WriteString(neo4jKey)
		builder.WriteString("`:row.`")
		builder.WriteString(alteryxFields[index])
		builder.WriteByte('`')
	}
	builder.WriteByte('}')
}
