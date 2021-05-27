package delete_test

import "testing"
import "github.com/tlarsen7572/graphyx/delete"

func TestDeleteNode(t *testing.T) {
	props := delete.DeleteNodesProperties{
		Label:    `Customer`,
		IdFields: []string{`Key`},
	}

	query := delete.GenerateDeleteNodes(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (d:`Customer`{`Key`:row.`Key`}) DETACH DELETE d"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteNodeUsingMultipleProperties(t *testing.T) {
	props := delete.DeleteNodesProperties{
		Label:    `Customer`,
		IdFields: []string{`Key1`, `Key2`},
	}

	query := delete.GenerateDeleteNodes(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (d:`Customer`{`Key1`:row.`Key1`,`Key2`:row.`Key2`}) DETACH DELETE d"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteNodeWithBackticks(t *testing.T) {
	props := delete.DeleteNodesProperties{
		Label:    "Cust`omer",
		IdFields: []string{"Ke`y"},
	}

	query := delete.GenerateDeleteNodes(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (d:`Cust``omer`{`Ke``y`:row.`Ke``y`}) DETACH DELETE d"

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

	query := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Customer` {`Key`:row.`LeftKey`})-[r:`IS_RELATED`]-(:`Customer` {`Key`:row.`RightKey`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}
