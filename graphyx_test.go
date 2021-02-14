package graphyx

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"testing"
)

func TestConnection(t *testing.T) {
	t.Skip()
	uri := `bolt://localhost:7687`
	username := `test`
	password := `test`
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) { c.Encrypted = false })
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	defer driver.Close()

	session, err := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	defer session.Close()

	query := `MATCH (people:Person)-[relatedTo]-(:Movie {title: "Cloud Atlas"}) RETURN *`
	//query := `MATCH (people:Person)-[relatedTo]-(:Movie {title: "Cloud Atlas"}) RETURN people.name, Type(relatedTo), relatedTo`
	//query := `CALL db.relationshipTypes()`
	//query := `MATCH p=()-[r:ACTED_IN]->() RETURN p`
	result, err := session.Run(query, map[string]interface{}{})
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	records := []neo4j.Path{}
	for result.Next() {
		record := result.Record()
		if p, ok := record.Get(`p`); ok {
			path := p.(neo4j.Path)
			records = append(records, path)
		}
	}
	t.Logf(`%v`, records)
}
