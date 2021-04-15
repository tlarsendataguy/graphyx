package output_test

import (
	"fmt"
	"strings"
	"testing"
)

type ConfigStruct struct {
	Label      string
	IdFields   []string
	PropFields []string
}

func buildMergeClause(builder *strings.Builder, config *ConfigStruct) {
	builder.WriteString(fmt.Sprintf("MERGE (newNode:`%v`{", config.Label))
	for index, id := range config.IdFields {
		if index > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("`%v`:row.`%v`", id, id))
	}
	builder.WriteString("})\n")
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
	buildMergeClause(builder, config)
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

func TestGenerateNodeQuery(t *testing.T) {
	config := &ConfigStruct{
		Label:      `TestLabel`,
		IdFields:   []string{`id1`, `id2`},
		PropFields: []string{`prop1`, `prop2`},
	}
	query := NodeQuery(config)
	expected := "UNWIND $batch AS row\n" +
		"MERGE (newNode:`TestLabel`{`id1`:row.`id1`,`id2`:row.`id2`})\n" +
		"ON CREATE SET newNode.`prop1`=row.`prop1`,newNode.`prop2`=row.`prop2`\n" +
		"ON MATCH SET newNode.`prop1`=row.`prop1`,newNode.`prop2`=row.`prop2`"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}
}

func TestGenerateNodeQueryWithNoProperties(t *testing.T) {
	config := &ConfigStruct{
		Label:      `TestLabel`,
		IdFields:   []string{`id1`, `id2`},
		PropFields: nil,
	}
	query := NodeQuery(config)
	expected := "UNWIND $batch AS row\n" +
		"MERGE (newNode:`TestLabel`{`id1`:row.`id1`,`id2`:row.`id2`})\n"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}
}
