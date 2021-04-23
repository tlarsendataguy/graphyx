package main_test

import (
	"fmt"
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

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: `neo4j`})
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

func TestBatch(t *testing.T) {
	err := deleteTestStuff()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	config := `<Configuration>
  <JSON>{"ConnStr":"http://localhost:7474","Username":"test","Password":"test","Database":"neo4j","ExportObject":"Node","BatchSize":10000,"NodeLabel":"TestLabel","NodeIdFields":["ID"],"NodePropFields":["Value"],"RelLabel":"","RelPropFields":[],"RelLeftLabel":"","RelLeftFields":[],"RelRightLabel":"","RelRightFields":[]}</JSON>
</Configuration>`
	plugin := &output.Neo4jOutput{}
	runner := sdk.RegisterToolTest(plugin, 1, config)
	runner.ConnectInput(`Input`, `TestNeo4jOutputNodes.txt`)
	runner.SimulateLifecycle()
	records, err := checkNumberOfItems(`MATCH (n:TestLabel) RETURN count(n)`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if records != 3 {
		t.Fatalf(`expected 3 records but got %v`, records)
	}
}

func TestDataTypesCopyCorrectly(t *testing.T) {
	err := deleteTestStuff()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	config := `<Configuration>
  <JSON>{"ConnStr":"http://localhost:7474","Username":"test","Password":"test","Database":"neo4j","ExportObject":"Node","BatchSize":10000,"NodeLabel":"TestLabel","NodeIdFields":["ID"],"NodePropFields":["ByteField","Int16Field","Int32Field","StringField","WStringField","V_StringField","V_WStringField","DateField","DateTimeField","FloatField","DoubleField","FixedDecimalField","BoolField"],"RelLabel":"","RelPropFields":[],"RelLeftLabel":"","RelLeftFields":[],"RelRightLabel":"","RelRightFields":[]}</JSON>
</Configuration>`
	plugin := &output.Neo4jOutput{}
	runner := sdk.RegisterToolTest(plugin, 1, config)
	runner.ConnectInput(`Input`, `TestNeo4jOutputTypes.txt`)
	runner.SimulateLifecycle()
	records, err := checkNumberOfItems(`MATCH (n:TestLabel) RETURN count(n)`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if records != 3 {
		t.Fatalf(`expected 3 records but got %v`, records)
	}
	properties, err := checkNumberOfItems(`MATCH (n:TestLabel {ID:1}) UNWIND keys(n) as k RETURN count(k)`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if properties != 14 {
		t.Fatalf(`expected 14 properties but got %v`, records)
	}
}

func TestDoNotRunIfUpdateOnly(t *testing.T) {
	err := deleteTestStuff()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	config := `<Configuration>
  <JSON>{"ConnStr":"http://localhost:7474","Username":"test","Password":"test","Database":"neo4j","ExportObject":"Node","BatchSize":10000,"NodeLabel":"TestLabel","NodeIdFields":["ID"],"NodePropFields":["ByteField","Int16Field","Int32Field","StringField","WStringField","V_StringField","V_WStringField","DateField","DateTimeField","FloatField","DoubleField","FixedDecimalField","BoolField"],"RelLabel":"","RelPropFields":[],"RelLeftLabel":"","RelLeftFields":[],"RelRightLabel":"","RelRightFields":[]}</JSON>
</Configuration>`
	plugin := &output.Neo4jOutput{}
	runner := sdk.RegisterToolTest(plugin, 1, config, sdk.UpdateOnly(true))
	runner.ConnectInput(`Input`, `TestNeo4jOutputTypes.txt`)
	runner.SimulateLifecycle()
	records, err := checkNumberOfItems(`MATCH (n:TestLabel) RETURN count(n)`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if records != 0 {
		t.Fatalf(`expected 0 records but got %v`, records)
	}
}

func TestEndToEnd(t *testing.T) {
	err := deleteTestStuff()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	configNodes := `<Configuration>
  <JSON>{"ConnStr":"http://localhost:7474","Username":"test","Password":"test","Database":"neo4j","ExportObject":"Node","BatchSize":10000,"NodeLabel":"TestLabel","NodeIdFields":["ID"],"NodePropFields":["Value"],"RelLabel":"","RelPropFields":[],"RelLeftLabel":"","RelLeftFields":[],"RelRightLabel":"","RelRightFields":[]}</JSON>
</Configuration>`
	pluginNodes := &output.Neo4jOutput{}
	runnerNodes := sdk.RegisterToolTest(pluginNodes, 1, configNodes)
	runnerNodes.ConnectInput(`Input`, `TestNeo4jOutputNodes.txt`)
	runnerNodes.SimulateLifecycle()

	configRelationships := `<Configuration>
  <JSON>{"ConnStr":"http://localhost:7474","Username":"test","Password":"test","Database":"neo4j","ExportObject":"Relationship","BatchSize":10000,"NodeLabel":"","NodeIdFields":[],"NodePropFields":[],"RelLabel":"TestRel","RelPropFields":["Value"],"RelLeftLabel":"TestLabel","RelLeftFields":[{"LeftID":"ID"}],"RelRightLabel":"TestLabel","RelRightFields":[{"RightID":"ID"}]}</JSON>
</Configuration>`
	pluginRelationships := &output.Neo4jOutput{}
	runnerRelationships := sdk.RegisterToolTest(pluginRelationships, 100, configRelationships)
	runnerRelationships.ConnectInput(`Input`, `TestNeo4jOutputRelationships.txt`)
	runnerRelationships.SimulateLifecycle()

	relationships, err := checkNumberOfItems(`MATCH (:TestLabel)-[r:TestRel]->(:TestLabel) RETURN count(r)`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if relationships != 3 {
		t.Fatalf(`expected 3 but got %v`, relationships)
	}

	err = deleteTestStuff()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
}

func deleteTestStuff() error {
	uri := `bolt://localhost:7687`
	database := `neo4j`
	username := `test`
	password := `test`
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return err
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: database})
	defer session.Close()

	deleteRelationships := `MATCH ()-[r:TestRel]-() DELETE r`
	deleteNodes := `MATCH (n:TestLabel) DELETE n`

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, txErr := tx.Run(deleteRelationships, nil)
		if txErr != nil {
			return nil, txErr
		}
		if txErr = result.Err(); txErr != nil {
			return nil, txErr
		}
		_, txErr = result.Consume()
		if txErr != nil {
			return nil, txErr
		}

		result, txErr = tx.Run(deleteNodes, nil)
		if txErr != nil {
			return nil, txErr
		}
		if txErr = result.Err(); txErr != nil {
			return nil, txErr
		}
		return result.Consume()
	})

	return err
}

func checkNumberOfItems(query string) (int, error) {
	uri := `bolt://localhost:7687`
	database := `neo4j`
	username := `test`
	password := `test`
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return 0, err
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: database})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, txErr := tx.Run(query, nil)
		if txErr != nil {
			return nil, txErr
		}
		if txErr = result.Err(); txErr != nil {
			return nil, txErr
		}
		hasRecord := result.Next()
		if !hasRecord {
			return nil, fmt.Errorf(`no record was returned`)
		}
		record := result.Record()
		return record.Values[0], nil
	})

	return int(result.(int64)), err
}
