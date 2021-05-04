package bolt_url

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type BoltInfo struct {
	BoltUrl string `json:"bolt_direct"`
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
	var bolt BoltInfo
	err = json.Unmarshal(responseBytes, &bolt)
	if err != nil {
		return ``, err
	}
	return bolt.BoltUrl, nil
}