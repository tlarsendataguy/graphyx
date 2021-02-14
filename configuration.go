package graphyx

import "encoding/xml"

type Configuration struct {
	ConnStr  string
	Username string
	Password string
	Query    string
	Fields   []Field `xml:"Fields>Field"`
}

type Field struct {
	Name     string    `xml:",attr"`
	DataType string    `xml:",attr"`
	Path     []Element `xml:"Path>Element"`
}

type Element struct {
	Key      string `xml:",attr"`
	DataType string `xml:",attr"`
}

func DecodeConfig(config string) (Configuration, error) {
	decoded := Configuration{}
	err := xml.Unmarshal([]byte(config), &decoded)
	return decoded, err
}
