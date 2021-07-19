package bolt_url

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type BoltInfo struct {
	BoltUrlV4 string `json:"bolt_direct"`
	BoltUrlV3 string `json:"bolt"`
}

func GetBoltUrl(httpEndpoint string) (string, error) {
	request, err := http.NewRequest(`GET`, httpEndpoint, nil)
	if err != nil {
		return ``, err
	}
	request.Header.Add(`Accept`, `application/json`)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return ``, err
	}
	if response.StatusCode != 200 {
		return ``, fmt.Errorf(`error connecting to Neo4j: %v`, response.Status)
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ``, err
	}
	return parseResponse(responseBytes)
}

func parseResponse(response []byte) (string, error) {
	var bolt BoltInfo
	err := json.Unmarshal(response, &bolt)
	if err != nil {
		return ``, err
	}
	if bolt.BoltUrlV4 != `` {
		return bolt.BoltUrlV4, nil
	}
	return bolt.BoltUrlV3, nil
}
