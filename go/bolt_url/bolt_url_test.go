package bolt_url

import "testing"

func TestGetBoltUrlFromNeo4jVersion4EndToEnd(t *testing.T) {
	var httpEndpoint = `http://localhost:7474`
	boltUrl, err := GetBoltUrl(httpEndpoint)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	expected := `bolt://localhost:7687`
	if boltUrl != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, boltUrl)
	}
}

func TestParseV4Response(t *testing.T) {
	response := []byte(`{
  "bolt_routing" : "neo4j://localhost:7687",
  "transaction" : "http://localhost:7474/db/{databaseName}/tx",
  "bolt_direct" : "bolt://localhost:7687",
  "neo4j_version" : "4.2.1",
  "neo4j_edition" : "enterprise"
}`)

	parsed, err := parseResponse(response)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if parsed != `bolt://localhost:7687` {
		t.Fatalf(`expected 'bolt://localhost:7687' but got '%v'`, parsed)
	}
}

func TestParseV3Response(t *testing.T) {
	response := []byte(`{
  "data" : "http://localhost:7474/db/data/",
  "management" : "http://localhost:7474/db/manage/",
  "bolt" : "bolt://localhost:7687"
}`)

	parsed, err := parseResponse(response)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if parsed != `bolt://localhost:7687` {
		t.Fatalf(`expected 'bolt://localhost:7687' but got '%v'`, parsed)
	}
}
