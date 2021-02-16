package graphyx

import (
	"github.com/tlarsen7572/goalteryx/sdk"
	"github.com/tlarsen7572/graphyx/input"
	"testing"
)

func TestInput(t *testing.T) {
	config := `<Configuration>
	<ConnStr>bolt://localhost:7687</ConnStr>
	<Username>neo4j</Username>
	<Password>neo4j</Password>
	<Query>MATCH p=()-[r:ACTED_IN]->() RETURN p</Query>
	<Fields>
		<Field Name="Field1" DataType="Integer">
			<Path>
				<Element DataType="Path" Key="p" />
				<Element DataType="List:Node" Key="Nodes" />
				<Element DataType="Node" Key="First" />
				<Element DataType="Integer" Key="ID" />
			</Path>
		</Field>
		<Field Name="Field2" DataType="Integer">
			<Path>
				<Element DataType="Path" Key="p" />
				<Element DataType="List:Relationship" Key="Relationships" />
				<Element DataType="Relationship" Key="First" />
				<Element DataType="Integer" Key="ID" />
			</Path>
		</Field>
	</Fields>
</Configuration>`
	plugin := &input.Neo4jInput{}
	runner := sdk.RegisterToolTest(plugin, 1, config)
	collector := runner.CaptureOutgoingAnchor(`Output`)
	runner.SimulateLifecycle()
	t.Logf(`%v`, collector.Data)
}

/*
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
*/
