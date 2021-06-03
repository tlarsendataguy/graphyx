package main_test

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/tlarsen7572/goalteryx/sdk"
	"github.com/tlarsen7572/graphyx/delete"
	"github.com/tlarsen7572/graphyx/input"
	"github.com/tlarsen7572/graphyx/output"
	"strings"
	"testing"
)

func TestInput(t *testing.T) {
	config := `<Configuration>
  <JSON>{"ConnStr":"http://localhost:7474","Username":"test","Password":"test","Database":"neo4j","Query":"MATCH p=()-[r:ACTED_IN]-&gt;() RETURN p","LastValidatedResponse":{"Error":"","ReturnValues":[{"Name":"p","DataType":"Path"}]},"Fields":[{"Name":"Path String","DataType":"String","Path":[{"Key":"p","DataType":"Path"},{"Key":"ToString","DataType":"String"}]},{"Name":"Field1","DataType":"Integer","Path":[{"Key":"p","DataType":"Path"},{"Key":"Nodes","DataType":"List:Node"},{"Key":"First","DataType":"Node"},{"Key":"ID","DataType":"Integer"}]},{"Name":"Field2","DataType":"String","Path":[{"Key":"p","DataType":"Path"},{"Key":"Relationships","DataType":"List:Relationship"},{"Key":"First","DataType":"Relationship"},{"Key":"Type","DataType":"String"}]}]}</JSON>
</Configuration>`
	plugin := &input.Neo4jInput{}
	runner := sdk.RegisterToolTest(plugin, 1, config)
	collector := runner.CaptureOutgoingAnchor(`Output`)
	runner.SimulateLifecycle()
	t.Logf(`%v`, collector.Data)
}

func TestAdHocQuery(t *testing.T) {
	uri := `bolt://localhost:7687`
	username := `test`
	password := `test`
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	query := `MATCH p = (:Person)-[*0..2]-(:Person) RETURN p SKIP 119 LIMIT 1`
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

				str := input.ToString(path)
				t.Log(str)
			}
		}

		if txErr = result.Err(); txErr != nil {
			return nil, txErr
		}

		return result.Consume()
	})
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

func TestEndToEndDelete(t *testing.T) {
	err := deleteTestStuff()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	err = addStuffForDeletion()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	configNodes := `<Configuration>
  <JSON>{"ConnStr":"http://localhost:7474","Username":"test","Password":"test","Database":"neo4j","DeleteObject":"Node","BatchSize":10000,"NodeLabel":"DELETE","NodeIdFields":["Id"]}</JSON>
</Configuration>`
	pluginNodes := &delete.Neo4jDelete{}
	runnerNodes := sdk.RegisterToolTest(pluginNodes, 1, configNodes)
	runnerNodes.ConnectInput(`Input`, `TestNeo4jDeleteNodes.txt`)
	runnerNodes.SimulateLifecycle()

	nodes, err := checkNumberOfItems(`MATCH (n:DELETE) RETURN count(n)`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if nodes != 0 {
		t.Fatalf(`expected 0 but got %v`, nodes)
	}

	err = deleteTestStuff()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
}

func addStuffForDeletion() error {
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

	createQuery := `CREATE (n1:DELETE {Id:1});
CREATE (n2:DELETE {Id:2});
CREATE (n3:DELETE {Id:3});
MATCH (n1:DELETE), (n2:DELETE) WHERE n1.Id=1 AND n2.Id=2
CREATE (n1)-[:Relates_To]->(n2);
MATCH (n2:DELETE), (n3:DELETE) WHERE n2.Id=2 AND n3.Id=3
CREATE (n2)-[:Relates_To]->(n3);`

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queries := strings.Split(createQuery, `;`)
		for _, query := range queries {
			if query == `` {
				continue
			}
			result, txErr := tx.Run(query, nil)
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
		}
		return nil, nil
	})

	return err
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

	deleteNodes := `MATCH (n:TestLabel) DETACH DELETE n;
MATCH (n:DELETE) DETACH DELETE n;`

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queries := strings.Split(deleteNodes, `;`)
		for _, query := range queries {
			if query == `` {
				continue
			}
			result, txErr := tx.Run(query, nil)
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
		}
		return nil, nil
	})

	return err
}
