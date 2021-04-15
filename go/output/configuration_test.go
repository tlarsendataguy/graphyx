package output_test

import (
	"github.com/tlarsen7572/graphyx/output"
	"testing"
)

func TestGenerateNodeQuery(t *testing.T) {
	config := &output.NodeConfig{
		Label:      `TestLabel`,
		IdFields:   []string{`id1`, `id2`},
		PropFields: []string{`prop1`, `prop2`},
	}
	query := output.NodeQuery(config)
	expected := "UNWIND $batch AS row\n" +
		"MERGE (newNode:`TestLabel`{`id1`:row.`id1`,`id2`:row.`id2`})\n" +
		"ON CREATE SET newNode.`prop1`=row.`prop1`,newNode.`prop2`=row.`prop2`\n" +
		"ON MATCH SET newNode.`prop1`=row.`prop1`,newNode.`prop2`=row.`prop2`"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}
}

func TestGenerateNodeQueryWithNoProperties(t *testing.T) {
	config := &output.NodeConfig{
		Label:      `TestLabel`,
		IdFields:   []string{`id1`, `id2`},
		PropFields: nil,
	}
	query := output.NodeQuery(config)
	expected := "UNWIND $batch AS row\n" +
		"MERGE (newNode:`TestLabel`{`id1`:row.`id1`,`id2`:row.`id2`})\n"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}
}

func TestGenerateNodesWithNoIds(t *testing.T) {
	config := &output.NodeConfig{
		Label:      `TestLabel`,
		IdFields:   nil,
		PropFields: []string{`prop1`, `prop2`},
	}
	query := output.NodeQuery(config)
	expected := "UNWIND $batch AS row\n" +
		"CREATE (newNode:`TestLabel`{`prop1`:row.`prop1`,`prop2`:row.`prop2`})"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}
}

func TestEscapeBackspace(t *testing.T) {
	config := &output.NodeConfig{
		Label:      `TestLabel`,
		IdFields:   []string{"id`1"},
		PropFields: []string{"prop`1"},
	}
	query := output.NodeQuery(config)
	expected := "UNWIND $batch AS row\n" +
		"MERGE (newNode:`TestLabel`{`id``1`:row.`id``1`})\n" +
		"ON CREATE SET newNode.`prop``1`=row.`prop``1`\n" +
		"ON MATCH SET newNode.`prop``1`=row.`prop``1`"
	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}

	config.IdFields = nil
	query = output.NodeQuery(config)
	expected = "UNWIND $batch AS row\n" +
		"CREATE (newNode:`TestLabel`{`prop``1`:row.`prop``1`})"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}
}
