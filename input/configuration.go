package input

import (
	"encoding/xml"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/tlarsen7572/goalteryx/sdk"
	"time"
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
			transferFuncs[index] = integerTransferFunc(fieldName, outInfo, getValueFunc)
		case `Float`:
			transferFuncs[index] = floatTransferFunc(fieldName, outInfo, getValueFunc)
		case `Boolean`:
			transferFuncs[index] = boolTransferFunc(fieldName, outInfo, getValueFunc)
		case `String`:
			transferFuncs[index] = stringTransferFunc(fieldName, outInfo, getValueFunc)
		case `Date`, `DateTime`:
			transferFuncs[index] = dateTimeTransferFunc(fieldName, outInfo, getValueFunc)
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

func integerTransferFunc(fieldName string, info *sdk.OutgoingRecordInfo, getValueFunc GetValueFunc) TransferFunc {
	return func(record neo4j.Record) error {
		value, getErr := getValueFunc(record)
		if getErr != nil {
			return getErr
		}
		if value == nil {
			info.IntFields[fieldName].SetNull()
			return nil
		}
		intValue, ok := value.(int)
		if !ok {
			return fmt.Errorf(`value %v is not an integer for field %v`, value, fieldName)
		}
		info.IntFields[fieldName].SetInt(intValue)
		return nil
	}
}

func floatTransferFunc(fieldName string, info *sdk.OutgoingRecordInfo, getValueFunc GetValueFunc) TransferFunc {
	return func(record neo4j.Record) error {
		value, getErr := getValueFunc(record)
		if getErr != nil {
			return getErr
		}
		if value == nil {
			info.FloatFields[fieldName].SetNull()
			return nil
		}
		floatValue, ok := value.(float64)
		if !ok {
			return fmt.Errorf(`value %v is not a float for field %v`, value, fieldName)
		}
		info.FloatFields[fieldName].SetFloat(floatValue)
		return nil
	}
}

func boolTransferFunc(fieldName string, info *sdk.OutgoingRecordInfo, getValueFunc GetValueFunc) TransferFunc {
	return func(record neo4j.Record) error {
		value, getErr := getValueFunc(record)
		if getErr != nil {
			return getErr
		}
		if value == nil {
			info.BoolFields[fieldName].SetNull()
			return nil
		}
		boolValue, ok := value.(bool)
		if !ok {
			return fmt.Errorf(`value %v is not a boolean for field %v`, value, fieldName)
		}
		info.BoolFields[fieldName].SetBool(boolValue)
		return nil
	}
}

func stringTransferFunc(fieldName string, info *sdk.OutgoingRecordInfo, getValueFunc GetValueFunc) TransferFunc {
	return func(record neo4j.Record) error {
		value, getErr := getValueFunc(record)
		if getErr != nil {
			return getErr
		}
		if value == nil {
			info.StringFields[fieldName].SetNull()
			return nil
		}
		stringValue, ok := value.(string)
		if !ok {
			return fmt.Errorf(`value %v is not a string for field %v`, value, fieldName)
		}
		info.StringFields[fieldName].SetString(stringValue)
		return nil
	}
}

func dateTimeTransferFunc(fieldName string, info *sdk.OutgoingRecordInfo, getValueFunc GetValueFunc) TransferFunc {
	return func(record neo4j.Record) error {
		value, getErr := getValueFunc(record)
		if getErr != nil {
			return getErr
		}
		if value == nil {
			info.DateTimeFields[fieldName].SetNull()
			return nil
		}
		dateTimeValue, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf(`value %v is not a string for field %v`, value, fieldName)
		}
		info.DateTimeFields[fieldName].SetDateTime(dateTimeValue)
		return nil
	}
}

func generateTransferFunc(iterator *pathIterator, field Field) (GetValueFunc, error) {
	element, isValid := iterator.NextField()
	if !isValid {
		return nil, fmt.Errorf(`no path was provided for field '%v'`, field.Name)
	}
	switch element.DataType {
	case `Integer`, `Float`, `Boolean`, `String`, `Date`, `DateTime`:
		return func(record neo4j.Record) (interface{}, error) {
			value, exists := record.Get(element.Key)
			if !exists {
				return nil, nil
			}
			return value, nil
		}, nil
	default:
		return nil, fmt.Errorf(`invalid field type '%v' for field '%v'`, field.DataType, field.Name)
	}
}
