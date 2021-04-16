package output

import (
	"encoding/json"
	"encoding/xml"
)

type xmlJson struct {
	JSON string `xml:",text"`
}

func decodeConfig(configStr string) (Configuration, error) {
	var err error
	var xmlConfig xmlJson
	err = xml.Unmarshal([]byte(configStr), &xmlConfig)
	if err != nil {
		return Configuration{}, err
	}
	var config Configuration
	err = json.Unmarshal([]byte(xmlConfig.JSON), &config)
	if err != nil {
		return Configuration{}, err
	}
	return config, nil
}
