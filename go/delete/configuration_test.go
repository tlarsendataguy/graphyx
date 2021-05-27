package delete_test

import "testing"
import "github.com/tlarsen7572/graphyx/delete"

func TestDeleteNode(t *testing.T) {
	label := `Customer`
	idFields := []string{`Key`}

	query := delete.GenerateDeleteNodes(label, idFields)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (d:`Customer`{`Key`:row.`Key`}) DETACH DELETE d"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteNodeUsingMultipleProperties(t *testing.T) {
	label := `Customer`
	idFields := []string{`Key1`, `Key2`}

	query := delete.GenerateDeleteNodes(label, idFields)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (d:`Customer`{`Key1`:row.`Key1`,`Key2`:row.`Key2`}) DETACH DELETE d"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteNodeWithBackticks(t *testing.T) {
	label := "Cust`omer"
	idFields := []string{"Ke`y"}

	query := delete.GenerateDeleteNodes(label, idFields)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (d:`Cust``omer`{`Ke``y`:row.`Ke``y`}) DETACH DELETE d"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}
