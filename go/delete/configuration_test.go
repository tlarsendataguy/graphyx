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
