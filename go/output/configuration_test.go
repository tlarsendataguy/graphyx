package output_test

import (
	"github.com/tlarsendataguy/graphyx/output"
	"testing"
)

func TestGenerateNodeQuery(t *testing.T) {
	config := &output.NodeConfig{
		Label:      `TestLabel`,
		IdFields:   []string{`id1`, `id2`},
		PropFields: []string{`prop1`, `prop2`},
	}
	query, _ := output.NodeQuery(config)
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
	query, _ := output.NodeQuery(config)
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
	query, _ := output.NodeQuery(config)
	expected := "UNWIND $batch AS row\n" +
		"CREATE (newNode:`TestLabel`{`prop1`:row.`prop1`,`prop2`:row.`prop2`})"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}
}

func TestEscapeBackquoteOnNodes(t *testing.T) {
	config := &output.NodeConfig{
		Label:      `TestLabel`,
		IdFields:   []string{"id`1"},
		PropFields: []string{"prop`1"},
	}
	query, _ := output.NodeQuery(config)
	expected := "UNWIND $batch AS row\n" +
		"MERGE (newNode:`TestLabel`{`id``1`:row.`id``1`})\n" +
		"ON CREATE SET newNode.`prop``1`=row.`prop``1`\n" +
		"ON MATCH SET newNode.`prop``1`=row.`prop``1`"
	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}

	config.IdFields = nil
	query, _ = output.NodeQuery(config)
	expected = "UNWIND $batch AS row\n" +
		"CREATE (newNode:`TestLabel`{`prop``1`:row.`prop``1`})"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}

	config.Label = "Test`Label"
	query, _ = output.NodeQuery(config)
	expected = "UNWIND $batch AS row\n" +
		"CREATE (newNode:`Test``Label`{`prop``1`:row.`prop``1`})"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}

	config.IdFields = []string{"id`1"}
	query, _ = output.NodeQuery(config)
	expected = "UNWIND $batch AS row\n" +
		"MERGE (newNode:`Test``Label`{`id``1`:row.`id``1`})\n" +
		"ON CREATE SET newNode.`prop``1`=row.`prop``1`\n" +
		"ON MATCH SET newNode.`prop``1`=row.`prop``1`"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}
}

func TestGenerateRelationshipQuery(t *testing.T) {
	config := &output.RelationshipConfig{
		LeftLabel:          `TestLabel`,
		RightLabel:         `TestLabel`,
		LeftAlteryxFields:  []string{`left1`, `left2`},
		LeftNeo4jFields:    []string{`id1`, `id2`},
		RightAlteryxFields: []string{`right1`, `right2`},
		RightNeo4jFields:   []string{`id1`, `id2`},
		Label:              `TestRel`,
		PropFields:         []string{`prop1`, `prop2`},
	}
	query, _ := output.RelationshipQuery(config)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (left:`TestLabel`{`id1`:row.`left1`,`id2`:row.`left2`})\n" +
		"MATCH (right:`TestLabel`{`id1`:row.`right1`,`id2`:row.`right2`})\n" +
		"MERGE (left)-[newRel:`TestRel`]->(right)\n" +
		"ON CREATE SET newRel.`prop1`=row.`prop1`,newRel.`prop2`=row.`prop2`\n" +
		"ON MATCH SET newRel.`prop1`=row.`prop1`,newRel.`prop2`=row.`prop2`"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}
}

func TestRelationshipQueryWithoutProperties(t *testing.T) {
	config := &output.RelationshipConfig{
		LeftLabel:          `TestLabel`,
		RightLabel:         `TestLabel`,
		LeftAlteryxFields:  []string{`left1`, `left2`},
		LeftNeo4jFields:    []string{`id1`, `id2`},
		RightAlteryxFields: []string{`right1`, `right2`},
		RightNeo4jFields:   []string{`id1`, `id2`},
		Label:              `TestRel`,
		PropFields:         nil,
	}
	query, _ := output.RelationshipQuery(config)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (left:`TestLabel`{`id1`:row.`left1`,`id2`:row.`left2`})\n" +
		"MATCH (right:`TestLabel`{`id1`:row.`right1`,`id2`:row.`right2`})\n" +
		"MERGE (left)-[newRel:`TestRel`]->(right)\n"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}
}

func TestNodeQueryWithoutLabel(t *testing.T) {
	config := &output.NodeConfig{
		Label:      ``,
		IdFields:   []string{`id1`, `id2`},
		PropFields: []string{`prop1`, `prop2`},
	}
	query, err := output.NodeQuery(config)
	if query != `` {
		t.Fatalf(`expected '' but got '%v'`, query)
	}
	if err == nil {
		t.Fatalf(`expected error but got nil`)
	}
	t.Logf(`%v`, err.Error())
}

func TestRelationshipQueryWithoutLabel(t *testing.T) {
	config := &output.RelationshipConfig{
		LeftLabel:          `TestLabel`,
		RightLabel:         `TestLabel`,
		LeftAlteryxFields:  []string{`left1`, `left2`},
		LeftNeo4jFields:    []string{`id1`, `id2`},
		RightAlteryxFields: []string{`right1`, `right2`},
		RightNeo4jFields:   []string{`id1`, `id2`},
		Label:              ``,
		PropFields:         nil,
	}
	query, err := output.RelationshipQuery(config)
	if query != `` {
		t.Fatalf(`expected '' but got '%v'`, query)
	}
	if err == nil {
		t.Fatalf(`expected error but got nil`)
	}
	t.Logf(`%v`, err.Error())
}

func TestEscapeLabelsInRelationshipQueries(t *testing.T) {
	config := &output.RelationshipConfig{
		LeftLabel:          "Test`Label",
		RightLabel:         "Test`Label",
		LeftAlteryxFields:  []string{`left1`, `left2`},
		LeftNeo4jFields:    []string{`id1`, `id2`},
		RightAlteryxFields: []string{`right1`, `right2`},
		RightNeo4jFields:   []string{`id1`, `id2`},
		Label:              "Test`Rel",
		PropFields:         nil,
	}
	query, _ := output.RelationshipQuery(config)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (left:`Test``Label`{`id1`:row.`left1`,`id2`:row.`left2`})\n" +
		"MATCH (right:`Test``Label`{`id1`:row.`right1`,`id2`:row.`right2`})\n" +
		"MERGE (left)-[newRel:`Test``Rel`]->(right)\n"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}
}

func TestEscapeFieldsInRelationshipQueries(t *testing.T) {
	config := &output.RelationshipConfig{
		LeftLabel:          "TestLabel",
		RightLabel:         "TestLabel",
		LeftAlteryxFields:  []string{"left`1", `left2`},
		LeftNeo4jFields:    []string{"id`1", `id2`},
		RightAlteryxFields: []string{"right`1", `right2`},
		RightNeo4jFields:   []string{"id`1", `id2`},
		Label:              "TestRel",
		PropFields:         nil,
	}
	query, _ := output.RelationshipQuery(config)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (left:`TestLabel`{`id``1`:row.`left``1`,`id2`:row.`left2`})\n" +
		"MATCH (right:`TestLabel`{`id``1`:row.`right``1`,`id2`:row.`right2`})\n" +
		"MERGE (left)-[newRel:`TestRel`]->(right)\n"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}
}

func TestRelationshipQueryWithoutLeftLabel(t *testing.T) {
	config := &output.RelationshipConfig{
		LeftLabel:          ``,
		RightLabel:         `TestLabel`,
		LeftAlteryxFields:  []string{`left1`, `left2`},
		LeftNeo4jFields:    []string{`id1`, `id2`},
		RightAlteryxFields: []string{`right1`, `right2`},
		RightNeo4jFields:   []string{`id1`, `id2`},
		Label:              `TestRel`,
		PropFields:         nil,
	}
	query, err := output.RelationshipQuery(config)
	if query != `` {
		t.Fatalf(`expected '' but got '%v'`, query)
	}
	if err == nil {
		t.Fatalf(`expected error but got nil`)
	}
	t.Logf(`%v`, err.Error())
}

func TestRelationshipQueryWithoutRightLabel(t *testing.T) {
	config := &output.RelationshipConfig{
		LeftLabel:          `TestLabel`,
		RightLabel:         ``,
		LeftAlteryxFields:  []string{`left1`, `left2`},
		LeftNeo4jFields:    []string{`id1`, `id2`},
		RightAlteryxFields: []string{`right1`, `right2`},
		RightNeo4jFields:   []string{`id1`, `id2`},
		Label:              `TestRel`,
		PropFields:         nil,
	}
	query, err := output.RelationshipQuery(config)
	if query != `` {
		t.Fatalf(`expected '' but got '%v'`, query)
	}
	if err == nil {
		t.Fatalf(`expected error but got nil`)
	}
	t.Logf(`%v`, err.Error())
}

func TestLeftAlteryxFieldsDoNotMatchNeo4jFields(t *testing.T) {
	config := &output.RelationshipConfig{
		LeftLabel:          `TestLabel`,
		RightLabel:         `TestLabel`,
		LeftAlteryxFields:  []string{`left1`, `left2`},
		LeftNeo4jFields:    []string{`id1`},
		RightAlteryxFields: []string{`right1`, `right2`},
		RightNeo4jFields:   []string{`id1`, `id2`},
		Label:              `TestRel`,
		PropFields:         nil,
	}
	query, err := output.RelationshipQuery(config)
	if query != `` {
		t.Fatalf(`expected '' but got '%v'`, query)
	}
	if err == nil {
		t.Fatalf(`expected error but got nil`)
	}
	t.Logf(`%v`, err.Error())
}

func TestRightAlteryxFieldsDoNotMatchNeo4jFields(t *testing.T) {
	config := &output.RelationshipConfig{
		LeftLabel:          `TestLabel`,
		RightLabel:         `TestLabel`,
		LeftAlteryxFields:  []string{`left1`, `left2`},
		LeftNeo4jFields:    []string{`id1`, `id2`},
		RightAlteryxFields: []string{`right1`, `right2`},
		RightNeo4jFields:   []string{`id1`},
		Label:              `TestRel`,
		PropFields:         nil,
	}
	query, err := output.RelationshipQuery(config)
	if query != `` {
		t.Fatalf(`expected '' but got '%v'`, query)
	}
	if err == nil {
		t.Fatalf(`expected error but got nil`)
	}
	t.Logf(`%v`, err.Error())
}

func TestRelationshipWithIds(t *testing.T) {
	config := &output.RelationshipConfig{
		LeftLabel:          "TestLabel",
		RightLabel:         "TestLabel",
		LeftAlteryxFields:  []string{"left1", `left2`},
		LeftNeo4jFields:    []string{"id1", `id2`},
		RightAlteryxFields: []string{"right1", `right2`},
		RightNeo4jFields:   []string{"id1", `id2`},
		Label:              "TestRel",
		IdFields:           []string{`relId1`, "relId`2"},
	}
	query, _ := output.RelationshipQuery(config)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (left:`TestLabel`{`id1`:row.`left1`,`id2`:row.`left2`})\n" +
		"MATCH (right:`TestLabel`{`id1`:row.`right1`,`id2`:row.`right2`})\n" +
		"MERGE (left)-[newRel:`TestRel` {`relId1`:row.`relId1`,`relId``2`:row.`relId``2`}]->(right)\n"

	if expected != query {
		t.Fatalf("expected\n\n%v\n\nbut got\n\n%v", expected, query)
	}
}
