package delete

import (
	"fmt"
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
	LeftNodeFields         []map[string]interface{}
	RightNodeLabel         string
	RightNodeFields        []map[string]interface{}
	leftNodeAlteryxFields  []string
	leftNodeNeo4jFields    []string
	rightNodeAlteryxFields []string
	rightNodeNeo4jFields   []string
}

func (p *DeleteRelationshipsProperties) escape() error {
	p.RelType = escapeName(p.RelType)
	p.LeftNodeLabel = escapeName(p.LeftNodeLabel)
	p.RightNodeLabel = escapeName(p.RightNodeLabel)

	for _, maps := range p.LeftNodeFields {
		for key, value := range maps {
			p.leftNodeAlteryxFields = append(p.leftNodeAlteryxFields, escapeName(key))
			valueStr, ok := value.(string)
			if !ok {
				return fmt.Errorf(`the Neo4j field mapping for Alteryx field '%v' in the left fields list is not a string`, key)
			}
			p.leftNodeNeo4jFields = append(p.leftNodeNeo4jFields, escapeName(valueStr))
		}
	}

	for _, maps := range p.RightNodeFields {
		for key, value := range maps {
			p.rightNodeAlteryxFields = append(p.rightNodeAlteryxFields, escapeName(key))
			valueStr, ok := value.(string)
			if !ok {
				return fmt.Errorf(`the Neo4j field mapping for Alteryx field '%v' in the right fields list is not a string`, key)
			}
			p.rightNodeNeo4jFields = append(p.rightNodeNeo4jFields, escapeName(valueStr))
		}
	}

	for index, item := range p.RelFields {
		p.RelFields[index] = escapeName(item)
	}
	return nil
}

func GenerateDeleteRelationships(props *DeleteRelationshipsProperties) (string, error) {
	err := props.escape()
	if err != nil {
		return ``, err
	}

	builder := &strings.Builder{}
	builder.WriteString("UNWIND $batch AS row\n")
	builder.WriteString("MATCH (")
	if props.LeftNodeLabel != `` {
		writeLabel(builder, props.LeftNodeLabel)
	}
	if len(props.leftNodeNeo4jFields) > 0 {
		writeProperties(builder, props.leftNodeNeo4jFields, props.leftNodeAlteryxFields)
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
	if len(props.rightNodeNeo4jFields) > 0 {
		writeProperties(builder, props.rightNodeNeo4jFields, props.rightNodeAlteryxFields)
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
