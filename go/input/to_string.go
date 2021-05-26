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

type RelationshipEndType int

const (
	relInvalid RelationshipEndType = 0
	relStartId RelationshipEndType = 1
	relEndId   RelationshipEndType = 2
)

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
		if nodeCount == 1 && len(v.Relationships) == 0 {
			recursiveToString(v.Nodes[0], builder)
			return
		}
		var nodeMap = make(map[int64]string, nodeCount)
		for _, node := range v.Nodes {
			nodeMap[node.Id] = ToString(node)
		}

		priorEnd := relInvalid
		var err error

		for index, rel := range v.Relationships {
			if index > 0 {
				priorEnd, err = processNotFirstRelationship(builder, v, index, rel, priorEnd, nodeMap)
				if err != nil {
					return
				}
				continue
			}
			priorEnd, err = processFirstRelationship(builder, v, rel, nodeMap)
			if err != nil {
				return
			}
		}
	case string:
		builder.WriteString(v)
	default:
		builder.WriteString(fmt.Sprintf(`%v`, value))
	}
}

func processNotFirstRelationship(builder *strings.Builder, v neo4j.Path, index int, rel neo4j.Relationship, priorEnd RelationshipEndType, nodeMap map[int64]string) (RelationshipEndType, error) {
	priorRel := v.Relationships[index-1]
	priorEndId := priorRel.EndId
	if priorEnd == relStartId {
		priorEndId = priorRel.StartId
	}
	if priorEndId == rel.StartId {
		node, ok := nodeMap[rel.EndId]
		if !ok {
			return relInvalid, fmt.Errorf(`could not find node with ID %v`, rel.EndId)
		}
		writeLeftToRight(builder, rel, node)
		return relEndId, nil
	}
	if priorEndId == rel.EndId {
		node, ok := nodeMap[rel.StartId]
		if !ok {
			return relInvalid, fmt.Errorf(`could not find node with ID %v`, rel.StartId)
		}
		writeRightToLeft(builder, rel, node)
		return relStartId, nil
	}
	node, ok := nodeMap[rel.StartId]
	if !ok {
		return relInvalid, fmt.Errorf(`could not find node with ID %v`, rel.StartId)
	}
	writeSeparator(builder, node)

	node, ok = nodeMap[rel.EndId]
	if !ok {
		return relInvalid, fmt.Errorf(`could not find node with ID %v`, rel.EndId)
	}
	writeLeftToRight(builder, rel, node)
	return relEndId, nil
}

func processFirstRelationship(builder *strings.Builder, v neo4j.Path, rel neo4j.Relationship, nodeMap map[int64]string) (RelationshipEndType, error) {
	if len(v.Relationships) > 1 {
		relationship2 := v.Relationships[1]
		if relationship2.StartId == rel.StartId {
			node, ok := nodeMap[rel.EndId]
			if !ok {
				return relInvalid, fmt.Errorf(`could not find node with ID %v`, rel.EndId)
			}
			builder.WriteString(node)

			node, ok = nodeMap[rel.StartId]
			if !ok {
				return relInvalid, fmt.Errorf(`could not find node with ID %v`, rel.StartId)
			}
			writeRightToLeft(builder, rel, node)
			return relStartId, nil
		}
	}
	node, ok := nodeMap[rel.StartId]
	if !ok {
		return relInvalid, fmt.Errorf(`could not find node with ID %v`, rel.StartId)
	}
	builder.WriteString(node)

	node, ok = nodeMap[rel.EndId]
	if !ok {
		return relInvalid, fmt.Errorf(`could not find node with ID %v`, rel.EndId)
	}
	writeLeftToRight(builder, rel, node)
	return relEndId, nil
}

func writeLeftToRight(builder *strings.Builder, rel neo4j.Relationship, node string) {
	builder.WriteByte('-')
	recursiveToString(rel, builder)
	builder.WriteByte('-')
	builder.WriteByte('>')
	builder.WriteString(node)
}

func writeRightToLeft(builder *strings.Builder, rel neo4j.Relationship, node string) {
	builder.WriteByte('<')
	builder.WriteByte('-')
	recursiveToString(rel, builder)
	builder.WriteByte('-')
	builder.WriteString(node)
}

func writeSeparator(builder *strings.Builder, node string) {
	builder.WriteByte(' ')
	builder.WriteByte('|')
	builder.WriteByte(' ')
	builder.WriteString(node)
}
