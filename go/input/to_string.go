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

		processor := PathProcessor{
			Builder: builder,
			Path:    v,
			NodeMap: nodeMap,
		}

		for processor.Next() {
		}
	case string:
		builder.WriteString(v)
	default:
		builder.WriteString(fmt.Sprintf(`%v`, value))
	}
}

type RelationshipEndType int

const (
	relInvalid RelationshipEndType = 0
	relStartId RelationshipEndType = 1
	relEndId   RelationshipEndType = 2
)

type PathProcessor struct {
	Builder                 *strings.Builder
	Path                    neo4j.Path
	NodeMap                 map[int64]string
	currentRelIndex         int
	currentRel              neo4j.Relationship
	lastRelationshipEndType RelationshipEndType
}

func (p *PathProcessor) Next() bool {
	if p.currentRelIndex == 0 {
		return p.processFirstRelationship()
	}
	return p.processNotFirstRelationship()
}

func (p *PathProcessor) success(endType RelationshipEndType) bool {
	p.lastRelationshipEndType = endType
	nextIndex := p.currentRelIndex + 1
	if nextIndex < len(p.Path.Relationships) {
		p.currentRelIndex = nextIndex
		p.currentRel = p.Path.Relationships[nextIndex]
		return true
	}
	return false
}

func (p *PathProcessor) processNotFirstRelationship() bool {
	priorRel := p.Path.Relationships[p.currentRelIndex-1]
	priorEndId := priorRel.EndId
	if p.lastRelationshipEndType == relStartId {
		priorEndId = priorRel.StartId
	}
	if priorEndId == p.currentRel.StartId {
		node, ok := p.NodeMap[p.currentRel.EndId]
		if !ok {
			return false
		}
		p.writeLeftToRight(node)

		return p.success(relEndId)
	}
	if priorEndId == p.currentRel.EndId {
		node, ok := p.NodeMap[p.currentRel.StartId]
		if !ok {
			return false
		}
		p.writeRightToLeft(node)
		return p.success(relStartId)
	}
	node, ok := p.NodeMap[p.currentRel.StartId]
	if !ok {
		return false
	}
	p.writeSeparator(node)

	node, ok = p.NodeMap[p.currentRel.EndId]
	if !ok {
		return false
	}
	p.writeLeftToRight(node)
	return p.success(relEndId)
}

func (p *PathProcessor) processFirstRelationship() bool {
	p.currentRel = p.Path.Relationships[0]
	if len(p.Path.Relationships) > 1 {
		relationship2 := p.Path.Relationships[1]
		if relationship2.StartId == p.currentRel.StartId {
			node, ok := p.NodeMap[p.currentRel.EndId]
			if !ok {
				return false
			}
			p.Builder.WriteString(node)

			node, ok = p.NodeMap[p.currentRel.StartId]
			if !ok {
				return false
			}
			p.writeRightToLeft(node)
			return p.success(relStartId)
		}
	}
	node, ok := p.NodeMap[p.currentRel.StartId]
	if !ok {
		return false
	}
	p.Builder.WriteString(node)

	node, ok = p.NodeMap[p.currentRel.EndId]
	if !ok {
		return false
	}
	p.writeLeftToRight(node)
	return p.success(relEndId)
}

func (p *PathProcessor) writeLeftToRight(node string) {
	p.Builder.WriteByte('-')
	recursiveToString(p.currentRel, p.Builder)
	p.Builder.WriteByte('-')
	p.Builder.WriteByte('>')
	p.Builder.WriteString(node)
}

func (p *PathProcessor) writeRightToLeft(node string) {
	p.Builder.WriteByte('<')
	p.Builder.WriteByte('-')
	recursiveToString(p.currentRel, p.Builder)
	p.Builder.WriteByte('-')
	p.Builder.WriteString(node)
}

func (p *PathProcessor) writeSeparator(node string) {
	p.Builder.WriteByte(' ')
	p.Builder.WriteByte('|')
	p.Builder.WriteByte(' ')
	p.Builder.WriteString(node)
}
