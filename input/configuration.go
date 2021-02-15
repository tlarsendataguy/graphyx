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

const source string = `Neo4j Input`

func CreateOutgoingObjects(fields []Field) (OutgoingObjects, error) {
	editor := sdk.EditingRecordInfo{}
	for _, field := range fields {
		switch field.DataType {
		case `Integer`:
			editor.AddInt64Field(field.Name, source)
		case `Float`:
			editor.AddDoubleField(field.Name, source)
		case `Boolean`:
			editor.AddBoolField(field.Name, source)
		case `Date`:
			editor.AddDateField(field.Name, source)
		case `DateTime`:
			editor.AddDateTimeField(field.Name, source)
		case `String`:
			editor.AddV_StringField(field.Name, source, 2147483648)
		default:
			return OutgoingObjects{}, fmt.Errorf(`field %v is invalid type %v`, field.Name, field.DataType)
		}
	}
	outInfo := editor.GenerateOutgoingRecordInfo()
	return OutgoingObjects{RecordInfo: outInfo}, nil
}
