package input_test

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/tlarsen7572/graphyx/input"
	"testing"
	"time"
)

func TestNodeToString(t *testing.T) {
	node := neo4j.Node{
		Id:     2,
		Labels: []string{`Something`},
		Props: map[string]interface{}{
			"Prop1": 2,
			"Prop2": time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC),
		},
	}
	actual := input.ToString(node)
	expected := `(:Something {"Prop1":2,"Prop2":"2020-01-02T03:04:05.000000006Z"})`
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}
}

func TestNodeToStringNoProperties(t *testing.T) {
	node := neo4j.Node{
		Id:     2,
		Labels: []string{`Something`},
		Props:  map[string]interface{}{},
	}
	actual := input.ToString(node)
	expected := `(:Something)`
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}
}

func TestEmptyNodeToString(t *testing.T) {
	node := neo4j.Node{
		Id:     2,
		Labels: []string{},
		Props:  map[string]interface{}{},
	}
	actual := input.ToString(node)
	expected := `()`
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}
}

func TestNodeToStringNoLabels(t *testing.T) {
	node := neo4j.Node{
		Id:     2,
		Labels: []string{},
		Props: map[string]interface{}{
			"Prop1": 2,
			"Prop2": "Hello world",
		},
	}
	actual := input.ToString(node)
	expected := `( {"Prop1":2,"Prop2":"Hello world"})`
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}
}

func TestRelationshipToString(t *testing.T) {
	rel := neo4j.Relationship{
		StartId: 10,
		EndId:   11,
		Id:      2,
		Type:    `Something`,
		Props: map[string]interface{}{
			"Prop1": 2,
			"Prop2": time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC),
		},
	}
	actual := input.ToString(rel)
	expected := `[:Something {"Prop1":2,"Prop2":"2020-01-02T03:04:05.000000006Z"}]`
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}
}

func TestRelationshipToStringNoProperties(t *testing.T) {
	rel := neo4j.Relationship{
		StartId: 10,
		EndId:   11,
		Id:      2,
		Type:    `Something`,
		Props:   map[string]interface{}{},
	}
	actual := input.ToString(rel)
	expected := `[:Something]`
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}
}

func TestEmptyRelationshipToString(t *testing.T) {
	rel := neo4j.Relationship{
		StartId: 10,
		EndId:   11,
		Id:      2,
		Type:    ``,
		Props:   map[string]interface{}{},
	}
	actual := input.ToString(rel)
	expected := `[]`
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}
}

func TestRelationshipToStringNoType(t *testing.T) {
	rel := neo4j.Relationship{
		StartId: 10,
		EndId:   11,
		Id:      2,
		Type:    ``,
		Props: map[string]interface{}{
			"Prop1": 2,
			"Prop2": time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC),
		},
	}
	actual := input.ToString(rel)
	expected := `[ {"Prop1":2,"Prop2":"2020-01-02T03:04:05.000000006Z"}]`
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}
}

func TestLeftToRightPathToString(t *testing.T) {
	path := neo4j.Path{
		Nodes: []neo4j.Node{
			{Id: 1, Labels: []string{`A`}, Props: map[string]interface{}{"Key": 1}},
			{Id: 2, Labels: []string{`B`}, Props: map[string]interface{}{"Key": 2}},
			{Id: 3, Labels: []string{`C`}, Props: map[string]interface{}{"Key": 3}},
		},
		Relationships: []neo4j.Relationship{
			{Id: 4, StartId: 1, EndId: 2, Type: `A_to_B`},
			{Id: 5, StartId: 2, EndId: 3, Type: `B_to_C`},
		},
	}
	actual := input.ToString(path)
	expected := `(:A {"Key":1})-[:A_to_B]->(:B {"Key":2})-[:B_to_C]->(:C {"Key":3})`
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}
}

func TestEmptyPathToString(t *testing.T) {
	path := neo4j.Path{
		Nodes:         []neo4j.Node{},
		Relationships: []neo4j.Relationship{},
	}
	actual := input.ToString(path)
	expected := ``
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}
}

func TestPathToStringOneNode(t *testing.T) {
	path := neo4j.Path{
		Nodes: []neo4j.Node{
			{Id: 1, Labels: []string{`A`}, Props: map[string]interface{}{"Key": 1}},
		},
		Relationships: []neo4j.Relationship{},
	}
	actual := input.ToString(path)
	expected := `(:A {"Key":1})`
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}
}

func TestPrimitivesToString(t *testing.T) {
	expected := `hello world`
	actual := input.ToString(expected)
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}

	expected = `1`
	actual = input.ToString(1)
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}

	expected = `1.2`
	actual = input.ToString(1.2)
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}
}

func TestMultipleRelationshipsFromOneNode(t *testing.T) {
	path := neo4j.Path{
		Nodes: []neo4j.Node{
			{Id: 5, Labels: []string{`Person`}, Props: map[string]interface{}{"Key": 5}},
			{Id: 119, Labels: []string{`Movie`}, Props: map[string]interface{}{"Key": 119}},
		},
		Relationships: []neo4j.Relationship{
			{Id: 4, StartId: 5, EndId: 119, Type: `DIRECTED`},
			{Id: 5, StartId: 5, EndId: 119, Type: `WROTE`},
		},
	}
	actual := input.ToString(path)
	expected := `(:Person {"Key":5})-[:DIRECTED]->(:Movie {"Key":119})<-[:WROTE]-(:Person {"Key":5})`
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}
}

func TestMergingPath(t *testing.T) {
	path := neo4j.Path{
		Nodes: []neo4j.Node{
			{Id: 5, Labels: []string{`Person`}, Props: map[string]interface{}{"Key": 5}},
			{Id: 119, Labels: []string{`Movie`}, Props: map[string]interface{}{"Key": 119}},
			{Id: 200, Labels: []string{`Person`}, Props: map[string]interface{}{"Key": 200}},
		},
		Relationships: []neo4j.Relationship{
			{Id: 4, StartId: 5, EndId: 119, Type: `DIRECTED`},
			{Id: 5, StartId: 200, EndId: 119, Type: `WROTE`},
		},
	}
	actual := input.ToString(path)
	expected := `(:Person {"Key":5})-[:DIRECTED]->(:Movie {"Key":119})<-[:WROTE]-(:Person {"Key":200})`
	if actual != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, actual)
	}
}
