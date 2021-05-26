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
		if nodeCount == 1 && len(v.Relationships) == 0 {
			recursiveToString(v.Nodes[0], builder)
			return
		}
		var nodeMap = make(map[int64]string, nodeCount)
		for _, node := range v.Nodes {
			nodeMap[node.Id] = ToString(node)
		}

		firstRel := v.Relationships[0]
		firstNode, firstNodeOk := nodeMap[firstRel.StartId]
		if !firstNodeOk {
			return
		}
		builder.WriteString(firstNode)

		for index, rel := range v.Relationships {
			if index > 0 {
				priorRel := v.Relationships[index-1]
				if priorRel.EndId == rel.StartId {
					node, ok := nodeMap[rel.EndId]
					if !ok {
						return
					}
					builder.WriteByte('-')
					recursiveToString(rel, builder)
					builder.WriteByte('-')
					builder.WriteByte('>')
					builder.WriteString(node)
					continue
				}
				if priorRel.EndId == rel.EndId {
					node, ok := nodeMap[rel.StartId]
					if !ok {
						return
					}
					builder.WriteByte('<')
					builder.WriteByte('-')
					recursiveToString(rel, builder)
					builder.WriteByte('-')
					builder.WriteString(node)
					continue
				}
				node, ok := nodeMap[rel.StartId]
				if !ok {
					return
				}
				builder.WriteByte(' ')
				builder.WriteByte('|')
				builder.WriteByte(' ')
				builder.WriteString(node)

				node, ok = nodeMap[rel.EndId]
				if !ok {
					return
				}
				builder.WriteByte('-')
				recursiveToString(rel, builder)
				builder.WriteByte('-')
				builder.WriteByte('>')
				builder.WriteString(node)
				continue
			}

			node, ok := nodeMap[rel.EndId]
			if !ok {
				return
			}
			builder.WriteByte('-')
			recursiveToString(rel, builder)
			builder.WriteByte('-')
			builder.WriteByte('>')
			builder.WriteString(node)
		}
	case string:
		builder.WriteString(v)
	default:
		builder.WriteString(fmt.Sprintf(`%v`, value))
	}
}
