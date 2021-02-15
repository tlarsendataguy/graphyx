package input

import (
	"encoding/xml"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/tlarsen7572/goalteryx/sdk"
)

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

type OutgoingObjects struct {
	RecordInfo    *sdk.OutgoingRecordInfo
	TransferFuncs []func(neo4j.Record)
}

func CreateOutgoingObjects(fields []Field) (OutgoingObjects, error) {
	editor := sdk.EditingRecordInfo{}
	for _, field := range fields {
		switch field.DataType {
		case `Integer`:
			editor.AddInt64Field(field.Name, `Neo4j Input`)
		default:
			return OutgoingObjects{}, fmt.Errorf(`field %v is invalid type %v`, field.Name, field.DataType)
		}
	}
	outInfo := editor.GenerateOutgoingRecordInfo()
	return OutgoingObjects{RecordInfo: outInfo}, nil
}
