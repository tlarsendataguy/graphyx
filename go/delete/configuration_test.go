package delete_test

import "testing"
import "github.com/tlarsendataguy/graphyx/delete"

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
	props := &delete.DeleteRelationshipsProperties{
		RelType:         `IS_RELATED`,
		RelFields:       []string{`Prop`},
		LeftNodeLabel:   `Customer`,
		LeftNodeFields:  []map[string]interface{}{{`LeftKey`: `Key`}},
		RightNodeLabel:  `Customer`,
		RightNodeFields: []map[string]interface{}{{`RightKey`: `Key`}},
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Customer` {`Key`:row.`LeftKey`})-[r:`IS_RELATED` {`Prop`:row.`Prop`}]-(:`Customer` {`Key`:row.`RightKey`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipsWithBackticks(t *testing.T) {
	props := &delete.DeleteRelationshipsProperties{
		RelType:         "IS_`RELATED",
		RelFields:       []string{"Pro`p"},
		LeftNodeLabel:   "Cust`omer",
		LeftNodeFields:  []map[string]interface{}{{"LeftKe`y": "Ke`y"}},
		RightNodeLabel:  "Cust`omer",
		RightNodeFields: []map[string]interface{}{{"RightKe`y": "Ke`y"}},
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Cust``omer` {`Ke``y`:row.`LeftKe``y`})-[r:`IS_``RELATED` {`Pro``p`:row.`Pro``p`}]-(:`Cust``omer` {`Ke``y`:row.`RightKe``y`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipWithoutLeftFields(t *testing.T) {
	props := &delete.DeleteRelationshipsProperties{
		RelType:         `IS_RELATED`,
		RelFields:       []string{`Prop`},
		LeftNodeLabel:   `Customer`,
		RightNodeLabel:  `Customer`,
		RightNodeFields: []map[string]interface{}{{"RightKey": "Key"}},
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Customer`)-[r:`IS_RELATED` {`Prop`:row.`Prop`}]-(:`Customer` {`Key`:row.`RightKey`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipWithoutRightFields(t *testing.T) {
	props := &delete.DeleteRelationshipsProperties{
		RelType:        `IS_RELATED`,
		RelFields:      []string{`Prop`},
		LeftNodeLabel:  `Customer`,
		LeftNodeFields: []map[string]interface{}{{"LeftKey": "Key"}},
		RightNodeLabel: `Customer`,
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Customer` {`Key`:row.`LeftKey`})-[r:`IS_RELATED` {`Prop`:row.`Prop`}]-(:`Customer`) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipWithoutRelationshipFields(t *testing.T) {
	props := &delete.DeleteRelationshipsProperties{
		RelType:         `IS_RELATED`,
		LeftNodeLabel:   `Customer`,
		LeftNodeFields:  []map[string]interface{}{{"LeftKey": "Key"}},
		RightNodeLabel:  `Customer`,
		RightNodeFields: []map[string]interface{}{{"RightKey": "Key"}},
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Customer` {`Key`:row.`LeftKey`})-[r:`IS_RELATED`]-(:`Customer` {`Key`:row.`RightKey`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipWithoutLeftLabel(t *testing.T) {
	props := &delete.DeleteRelationshipsProperties{
		RelType:         `IS_RELATED`,
		RelFields:       []string{`Prop`},
		LeftNodeFields:  []map[string]interface{}{{"LeftKey": "Key"}},
		RightNodeLabel:  `Customer`,
		RightNodeFields: []map[string]interface{}{{"RightKey": "Key"}},
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH ( {`Key`:row.`LeftKey`})-[r:`IS_RELATED` {`Prop`:row.`Prop`}]-(:`Customer` {`Key`:row.`RightKey`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipWithoutRelationshipType(t *testing.T) {
	props := &delete.DeleteRelationshipsProperties{
		RelFields:       []string{`Prop`},
		LeftNodeLabel:   `Customer`,
		LeftNodeFields:  []map[string]interface{}{{"LeftKey": "Key"}},
		RightNodeLabel:  `Customer`,
		RightNodeFields: []map[string]interface{}{{"RightKey": "Key"}},
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Customer` {`Key`:row.`LeftKey`})-[r {`Prop`:row.`Prop`}]-(:`Customer` {`Key`:row.`RightKey`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipWithoutRightLabel(t *testing.T) {
	props := &delete.DeleteRelationshipsProperties{
		RelType:         `IS_RELATED`,
		RelFields:       []string{`Prop`},
		LeftNodeLabel:   `Customer`,
		LeftNodeFields:  []map[string]interface{}{{"LeftKey": "Key"}},
		RightNodeFields: []map[string]interface{}{{"RightKey": "Key"}},
	}

	query, _ := delete.GenerateDeleteRelationships(props)
	expected := "UNWIND $batch AS row\n" +
		"MATCH (:`Customer` {`Key`:row.`LeftKey`})-[r:`IS_RELATED` {`Prop`:row.`Prop`}]-( {`Key`:row.`RightKey`}) DELETE r"

	if query != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, query)
	}
}

func TestDeleteRelationshipMismatchedLeftFields(t *testing.T) {
	props := &delete.DeleteRelationshipsProperties{
		RelType:         `IS_RELATED`,
		RelFields:       []string{`Prop`},
		LeftNodeLabel:   `Customer`,
		LeftNodeFields:  []map[string]interface{}{{"LeftKey": "Key"}, {"LeftKey2": 12345}},
		RightNodeLabel:  `Customer`,
		RightNodeFields: []map[string]interface{}{{"RightKey": "Key"}},
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
	props := &delete.DeleteRelationshipsProperties{
		RelType:         `IS_RELATED`,
		RelFields:       []string{`Prop`},
		LeftNodeLabel:   `Customer`,
		LeftNodeFields:  []map[string]interface{}{{"LeftKey": "Key"}},
		RightNodeLabel:  `Customer`,
		RightNodeFields: []map[string]interface{}{{"RightKey": "Key"}, {"RightKey2": 12345}},
	}

	query, err := delete.GenerateDeleteRelationships(props)
	if err == nil {
		t.Fatalf(`expected an error but got none`)
	}
	if query != `` {
		t.Fatalf("expected empty string but got\n%v", query)
	}
}
