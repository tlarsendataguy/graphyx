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

/*
func TestDeleteRelationship(t *testing.T) {
	relType := `IS_RELATED`
	relFields := []string{`Prop`}
	leftNodeLabel := `Customer`
	leftNodeAlteryxFields := []string{`LeftKey`}
	leftNodeNeo4jFields := []string{`Key`}
	rightNodeLabel := `Customer`
	rightNodeAlteryxFields := []string{`RightKey`}
	rightNodeNeo4jFields := []string{`Key`}

	query := delete.GenerateDeleteRelationships(label, idFields)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (d:`Customer`{`Key`:row.`Key`}) DETACH DELETE d"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}
*/
