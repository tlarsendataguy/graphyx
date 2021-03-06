package main_test

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/tlarsen7572/goalteryx/sdk"
	"github.com/tlarsen7572/graphyx/input"
	"testing"
)

func TestInput(t *testing.T) {
	config := `<Configuration>
  <ConnStr>http://localhost:7474</ConnStr>
  <Username>user</Username>
  <Password>password</Password>
  <Query>MATCH p=()-[r:ACTED_IN]-&gt;() RETURN p</Query>
  <Fields>
    <Field>
      <Path>
        <Element>
          <Key>p</Key>
          <DataType>Path</DataType>
        </Element>
      </Path>
      <Path>
        <Element>
          <Key>Nodes</Key>
          <DataType>List:Node</DataType>
        </Element>
      </Path>
      <Path>
        <Element>
          <Key>First</Key>
          <DataType>Node</DataType>
        </Element>
      </Path>
      <Path>
        <Element>
          <Key>ID</Key>
          <DataType>Integer</DataType>
        </Element>
      </Path>
      <Name>Field1</Name>
      <DataType>Integer</DataType>
    </Field>
  </Fields>
  <Fields>
    <Field>
      <Path>
        <Element>
          <Key>p</Key>
          <DataType>Path</DataType>
        </Element>
      </Path>
      <Path>
        <Element>
          <Key>Relationships</Key>
          <DataType>List:Relationship</DataType>
        </Element>
      </Path>
      <Path>
        <Element>
          <Key>First</Key>
          <DataType>Relationship</DataType>
        </Element>
      </Path>
      <Path>
        <Element>
          <Key>Type</Key>
          <DataType>String</DataType>
        </Element>
      </Path>
      <Name>Field2</Name>
      <DataType>String</DataType>
    </Field>
  </Fields>
  <LastValidatedResponse>
    <ReturnValues>
      <ReturnValue>
        <Name>p</Name>
        <DataType>Path</DataType>
      </ReturnValue>
    </ReturnValues>
    <Error>
    </Error>
  </LastValidatedResponse>
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

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
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
