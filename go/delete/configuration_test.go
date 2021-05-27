package delete_test

import "testing"
import "github.com/tlarsen7572/graphyx/delete"

func TestDeleteNode(t *testing.T) {
	props := &delete.DeleteNodesProperties{
		Label:    `Customer`,
		IdFields: []string{`Key`},
	}

	query := delete.GenerateDeleteNodes(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (d:`Customer` {`Key`:row.`Key`}) DETACH DELETE d"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteNodeUsingMultipleProperties(t *testing.T) {
	props := &delete.DeleteNodesProperties{
		Label:    `Customer`,
		IdFields: []string{`Key1`, `Key2`},
	}

	query := delete.GenerateDeleteNodes(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (d:`Customer` {`Key1`:row.`Key1`,`Key2`:row.`Key2`}) DETACH DELETE d"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteNodeWithBackticks(t *testing.T) {
	props := &delete.DeleteNodesProperties{
		Label:    "Cust`omer",
		IdFields: []string{"Ke`y"},
	}

	query := delete.GenerateDeleteNodes(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (d:`Cust``omer` {`Ke``y`:row.`Ke``y`}) DETACH DELETE d"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteNodeWithoutIdFields(t *testing.T) {
	props := &delete.DeleteNodesProperties{
		Label: `Customer`,
	}

	query := delete.GenerateDeleteNodes(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (d:`Customer`) DETACH DELETE d"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteNodeWithoutLabel(t *testing.T) {
	props := &delete.DeleteNodesProperties{
		IdFields: []string{`Key`},
	}

	query := delete.GenerateDeleteNodes(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (d {`Key`:row.`Key`}) DETACH DELETE d"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationship(t *testing.T) {
	props := delete.DeleteRelationshipsProperties{
		RelType:                `IS_RELATED`,
		RelFields:              []string{`Prop`},
		LeftNodeLabel:          `Customer`,
		LeftNodeAlteryxFields:  []string{`LeftKey`},
		LeftNodeNeo4jFields:    []string{`Key`},
		RightNodeLabel:         `Customer`,
		RightNodeAlteryxFields: []string{`RightKey`},
		RightNodeNeo4jFields:   []string{`Key`},
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Customer` {`Key`:row.`LeftKey`})-[r:`IS_RELATED` {`Prop`:row.`Prop`}]-(:`Customer` {`Key`:row.`RightKey`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipsWithBackticks(t *testing.T) {
	props := delete.DeleteRelationshipsProperties{
		RelType:                "IS_`RELATED",
		RelFields:              []string{"Pro`p"},
		LeftNodeLabel:          "Cust`omer",
		LeftNodeAlteryxFields:  []string{"LeftKe`y"},
		LeftNodeNeo4jFields:    []string{"Ke`y"},
		RightNodeLabel:         "Cust`omer",
		RightNodeAlteryxFields: []string{"RightKe`y"},
		RightNodeNeo4jFields:   []string{"Ke`y"},
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Cust``omer` {`Ke``y`:row.`LeftKe``y`})-[r:`IS_``RELATED` {`Pro``p`:row.`Pro``p`}]-(:`Cust``omer` {`Ke``y`:row.`RightKe``y`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipWithoutLeftFields(t *testing.T) {
	props := delete.DeleteRelationshipsProperties{
		RelType:                `IS_RELATED`,
		RelFields:              []string{`Prop`},
		LeftNodeLabel:          `Customer`,
		RightNodeLabel:         `Customer`,
		RightNodeAlteryxFields: []string{`RightKey`},
		RightNodeNeo4jFields:   []string{`Key`},
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Customer`)-[r:`IS_RELATED` {`Prop`:row.`Prop`}]-(:`Customer` {`Key`:row.`RightKey`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipWithoutRightFields(t *testing.T) {
	props := delete.DeleteRelationshipsProperties{
		RelType:               `IS_RELATED`,
		RelFields:             []string{`Prop`},
		LeftNodeLabel:         `Customer`,
		LeftNodeAlteryxFields: []string{"LeftKey"},
		LeftNodeNeo4jFields:   []string{"Key"},
		RightNodeLabel:        `Customer`,
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Customer` {`Key`:row.`LeftKey`})-[r:`IS_RELATED` {`Prop`:row.`Prop`}]-(:`Customer`) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipWithoutRelationshipFields(t *testing.T) {
	props := delete.DeleteRelationshipsProperties{
		RelType:                `IS_RELATED`,
		LeftNodeLabel:          `Customer`,
		LeftNodeAlteryxFields:  []string{`LeftKey`},
		LeftNodeNeo4jFields:    []string{`Key`},
		RightNodeLabel:         `Customer`,
		RightNodeAlteryxFields: []string{`RightKey`},
		RightNodeNeo4jFields:   []string{`Key`},
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Customer` {`Key`:row.`LeftKey`})-[r:`IS_RELATED`]-(:`Customer` {`Key`:row.`RightKey`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipWithoutLeftLabel(t *testing.T) {
	props := delete.DeleteRelationshipsProperties{
		RelType:                `IS_RELATED`,
		RelFields:              []string{`Prop`},
		LeftNodeAlteryxFields:  []string{`LeftKey`},
		LeftNodeNeo4jFields:    []string{`Key`},
		RightNodeLabel:         `Customer`,
		RightNodeAlteryxFields: []string{`RightKey`},
		RightNodeNeo4jFields:   []string{`Key`},
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH ( {`Key`:row.`LeftKey`})-[r:`IS_RELATED` {`Prop`:row.`Prop`}]-(:`Customer` {`Key`:row.`RightKey`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipWithoutRelationshipType(t *testing.T) {
	props := delete.DeleteRelationshipsProperties{
		RelFields:              []string{`Prop`},
		LeftNodeLabel:          `Customer`,
		LeftNodeAlteryxFields:  []string{`LeftKey`},
		LeftNodeNeo4jFields:    []string{`Key`},
		RightNodeLabel:         `Customer`,
		RightNodeAlteryxFields: []string{`RightKey`},
		RightNodeNeo4jFields:   []string{`Key`},
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Customer` {`Key`:row.`LeftKey`})-[r {`Prop`:row.`Prop`}]-(:`Customer` {`Key`:row.`RightKey`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipWithoutRightLabel(t *testing.T) {
	props := delete.DeleteRelationshipsProperties{
		RelType:                `IS_RELATED`,
		RelFields:              []string{`Prop`},
		LeftNodeLabel:          `Customer`,
		LeftNodeAlteryxFields:  []string{`LeftKey`},
		LeftNodeNeo4jFields:    []string{`Key`},
		RightNodeAlteryxFields: []string{`RightKey`},
		RightNodeNeo4jFields:   []string{`Key`},
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Customer` {`Key`:row.`LeftKey`})-[r:`IS_RELATED` {`Prop`:row.`Prop`}]-( {`Key`:row.`RightKey`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipMismatchedLeftFields(t *testing.T) {
	props := delete.DeleteRelationshipsProperties{
		RelType:                `IS_RELATED`,
		RelFields:              []string{`Prop`},
		LeftNodeLabel:          `Customer`,
		LeftNodeAlteryxFields:  []string{`LeftKey`},
		LeftNodeNeo4jFields:    []string{`Key1`, `Key2`},
		RightNodeLabel:         `Customer`,
		RightNodeAlteryxFields: []string{`RightKey`},
		RightNodeNeo4jFields:   []string{`Key`},
	}

	query, err := delete.GenerateDeleteRelationships(props)
	if err == nil {
		t.Fatalf(`expected an error but got none`)
	}
	if query != `` {
		t.Fatalf("expected empty string but got\n%v", query)
	}
}

func TestDeleteRelationshipMismatchedRightFields(t *testing.T) {
	props := delete.DeleteRelationshipsProperties{
		RelType:                `IS_RELATED`,
		RelFields:              []string{`Prop`},
		LeftNodeLabel:          `Customer`,
		LeftNodeAlteryxFields:  []string{`LeftKey`},
		LeftNodeNeo4jFields:    []string{`Key`},
		RightNodeLabel:         `Customer`,
		RightNodeAlteryxFields: []string{`RightKey`},
		RightNodeNeo4jFields:   []string{`Key1`, `Key2`},
	}

	query, err := delete.GenerateDeleteRelationships(props)
	if err == nil {
		t.Fatalf(`expected an error but got none`)
	}
	if query != `` {
		t.Fatalf("expected empty string but got\n%v", query)
	}
}
