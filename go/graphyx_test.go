package main_test

import (
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/tlarsen7572/goalteryx/sdk"
	"github.com/tlarsen7572/graphyx/input"
	"github.com/tlarsen7572/graphyx/output"
	"testing"
)

func TestInput(t *testing.T) {
	config := `<Configuration>
  <JSON>{"ConnStr":"http://localhost:7474","Username":"test","Password":"test","Query":"MATCH p=()-[r:ACTED_IN]-&gt;() RETURN p","LastValidatedResponse":{"Error":"","ReturnValues":[{"Name":"p","DataType":"Path"}]},"Fields":[{"Name":"Field1","DataType":"Integer","Path":[{"Key":"p","DataType":"Path"},{"Key":"Nodes","DataType":"List:Node"},{"Key":"First","DataType":"Node"},{"Key":"ID","DataType":"Integer"}]},{"Name":"Field2","DataType":"String","Path":[{"Key":"p","DataType":"Path"},{"Key":"Relationships","DataType":"List:Relationship"},{"Key":"First","DataType":"Relationship"},{"Key":"Type","DataType":"String"}]}]}</JSON>
</Configuration>`
	plugin := &input.Neo4jInput{}
	runner := sdk.RegisterToolTest(plugin, 1, config)
	collector := runner.CaptureOutgoingAnchor(`Output`)
	runner.SimulateLifecycle()
	t.Logf(`%v`, collector.Data)
}

func TestConnection(t *testing.T) {
	uri := `bolt://localhost:7687`
	username := `test`
	password := `test`
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: `Movie Database`})
	defer session.Close()

	//query := `MATCH (people:Person)-[relatedTo]-(:Movie {title: "Cloud Atlas"}) RETURN *`
	//query := `MATCH (people:Person)-[relatedTo]-(:Movie {title: "Cloud Atlas"}) RETURN people.name, Type(relatedTo), relatedTo`
	//query := `CALL db.relationshipTypes()`
	query := `MATCH p=()-[r:ACTED_IN]->() RETURN p`

	_, err = session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, txErr := tx.Run(query, nil)
		if txErr != nil {
			return nil, txErr
		}
		for result.Next() {
			record := result.Record()
			if p, ok := record.Get(`p`); ok {
				path := p.(neo4j.Path)
				t.Logf(`%v`, path)
			}
		}

		if txErr = result.Err(); txErr != nil {
			return nil, txErr
		}

		return result.Consume()
	})

}

func TestOutput(t *testing.T) {
	uri := `bolt://localhost:7687`
	username := `test`
	password := `test`
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: `Movie Database`})
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		tx.Run(``, nil)
		return nil, nil
	})
}

func TestOutputToolNodes(t *testing.T) {
	config := `<Configuration>
  <JSON>{"ConnStr":"http://localhost:7474","Username":"test","Password":"test","Database":"Movie Database","ExportObject":"Node","NodeLabel":"TestLabel","NodeIdFields":["ID"],"NodePropFields":["Value"],"RelLabel":"","RelPropFields":[],"RelLeftLabel":"","RelLeftFields":[],"RelRightLabel":"","RelRightFields":[]}</JSON>
</Configuration>`
	plugin := &output.Neo4jOutput{}
	runner := sdk.RegisterToolTest(plugin, 1, config)
	runner.ConnectInput(`Input`, `TestNeo4jOutputNodes.txt`)
	runner.SimulateLifecycle()
	configBytes, err := json.Marshal(plugin.Config())
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	t.Logf(string(configBytes))
	t.Logf(plugin.Query())
}

func TestOutputToolRelationships(t *testing.T) {
	config := `<Configuration>
  <JSON>{"ConnStr":"http://localhost:7474","Username":"test","Password":"test","Database":"Movie Database","ExportObject":"Relationship","NodeLabel":"","NodeIdFields":[],"NodePropFields":[],"RelLabel":"TestRel","RelPropFields":["Value"],"RelLeftLabel":"TestLabel","RelLeftFields":[{"LeftID":"ID"}],"RelRightLabel":"TestLabel","RelRightFields":[{"RightID":"ID"}]}</JSON>
</Configuration>`
	plugin := &output.Neo4jOutput{}
	runner := sdk.RegisterToolTest(plugin, 1, config)
	runner.ConnectInput(`Input`, `TestNeo4jOutputNodes.txt`)
	runner.SimulateLifecycle()
	configBytes, err := json.Marshal(plugin.Config())
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	t.Logf(string(configBytes))
	t.Logf(plugin.Query())
}
