package input

import "testing"

func TestGetBoltUrl(t *testing.T) {
	var httpEndpoint = `http://localhost:7474`
	boltUrl, err := getBoltUrl(httpEndpoint)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	expected := `bolt://localhost:7687`
	if boltUrl != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, boltUrl)
	}
}
