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

type TransferFunc func(neo4j.Record) error
type GetValueFunc func(neo4j.Record) (interface{}, error)

type OutgoingObjects struct {
	RecordInfo    *sdk.OutgoingRecordInfo
	TransferFuncs []TransferFunc
}

type pathIterator struct {
	elements     []Element
	currentIndex int
}

func (i *pathIterator) NextField() (Element, bool) {
	if i.currentIndex < len(i.elements) {
		element := i.elements[i.currentIndex]
		i.currentIndex++
		return element, true
	}
	return Element{}, false
}

const source string = `Neo4j Input`

func CreateOutgoingObjects(fields []Field) (OutgoingObjects, error) {
	getValueFuncs := make([]GetValueFunc, len(fields))
	transferFuncs := make([]TransferFunc, len(fields))
	editor := sdk.EditingRecordInfo{}
	var err error
	for index, field := range fields {
		field.Name, err = addFieldToEditor(field, &editor)
		if err != nil {
			return OutgoingObjects{}, err
		}
		iterator := &pathIterator{elements: field.Path}
		getValueFuncs[index], err = generateTransferFunc(iterator, field)
		if err != nil {
			return OutgoingObjects{}, err
		}
	}
	outInfo := editor.GenerateOutgoingRecordInfo()
	for index, getValueFunc := range getValueFuncs {
		fieldName := fields[index].Name
		fieldType := fields[index].DataType
		switch fieldType {
		case `Integer`:
			transferFuncs[index] = func(record neo4j.Record) error {
				value, getErr := getValueFunc(record)
				if getErr != nil {
					return getErr
				}
				if value == nil {
					outInfo.IntFields[fieldName].SetNull()
					return nil
				}
				intValue, ok := value.(int)
				if !ok {
					return fmt.Errorf(`value %v is not an integer for field %v`, value, fieldName)
				}
				outInfo.IntFields[fieldName].SetInt(intValue)
				return nil
			}
		case `Float`, `Boolean`, `Date`, `DateTime`, `String`:
			continue
		default:
			return OutgoingObjects{}, fmt.Errorf(`invalid field type '%v' for field '%v'`, fieldType, fieldName)
		}
	}
	return OutgoingObjects{RecordInfo: outInfo, TransferFuncs: transferFuncs}, nil
}

func addFieldToEditor(field Field, editor *sdk.EditingRecordInfo) (string, error) {
	switch field.DataType {
	case `Integer`:
		return editor.AddInt64Field(field.Name, source), nil
	case `Float`:
		return editor.AddDoubleField(field.Name, source), nil
	case `Boolean`:
		return editor.AddBoolField(field.Name, source), nil
	case `Date`:
		return editor.AddDateField(field.Name, source), nil
	case `DateTime`:
		return editor.AddDateTimeField(field.Name, source), nil
	case `String`:
		return editor.AddV_StringField(field.Name, source, 2147483648), nil
	default:
		return ``, fmt.Errorf(`field %v is invalid type %v`, field.Name, field.DataType)
	}
}

func generateTransferFunc(iterator *pathIterator, field Field) (GetValueFunc, error) {
	element, isValid := iterator.NextField()
	if !isValid {
		return nil, fmt.Errorf(`no path was provided for field '%v'`, field.Name)
	}
	switch element.DataType {
	case `Integer`:
		return func(record neo4j.Record) (interface{}, error) {
			value, exists := record.Get(element.Key)
			if !exists {
				return nil, nil
			}
			return value, nil
		}, nil
	case `Float`, `Boolean`, `Date`, `DateTime`, `String`:
		return nil, nil
	default:
		return nil, fmt.Errorf(`invalid field type '%v' for field '%v'`, field.DataType, field.Name)
	}
}
