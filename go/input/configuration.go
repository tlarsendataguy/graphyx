package input

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/tlarsen7572/goalteryx/sdk"
	"strconv"
	"strings"
	"time"
)

type XmlJson struct {
	JSON string `xml:",text"`
}

type Configuration struct {
	ConnStr  string
	Username string
	Password string
	Query    string
	Database string
	Fields   []Field
}

type Field struct {
	Name     string
	DataType string
	Path     []Element
}

type Element struct {
	Key      string
	DataType string
}

func DecodeConfig(config string) (Configuration, error) {
	xmlDecoded := XmlJson{}
	err := xml.Unmarshal([]byte(config), &xmlDecoded)
	if err != nil {
		return Configuration{}, err
	}
	decoded := Configuration{}
	err = json.Unmarshal([]byte(xmlDecoded.JSON), &decoded)
	return decoded, err
}

type TransferFunc func(*neo4j.Record) error
type GetValueFunc func(*neo4j.Record) (interface{}, error)

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
		return editor.AddV_WStringField(field.Name, source, 1073741823), nil
	default:
		return ``, fmt.Errorf(`field %v is invalid type %v`, field.Name, field.DataType)
	}
}

func integerTransferFunc(fieldName string, info *sdk.OutgoingRecordInfo, getValueFunc GetValueFunc) TransferFunc {
	return func(record *neo4j.Record) error {
		value, getErr := getValueFunc(record)
		if getErr != nil {
			return getErr
		}
		if value == nil {
			info.IntFields[fieldName].SetNull()
			return nil
		}
		intValue, ok := value.(int64)
		if !ok {
			return fmt.Errorf(`value %v is not an integer for field %v`, value, fieldName)
		}
		info.IntFields[fieldName].SetInt(int(intValue))
		return nil
	}
}

func floatTransferFunc(fieldName string, info *sdk.OutgoingRecordInfo, getValueFunc GetValueFunc) TransferFunc {
	return func(record *neo4j.Record) error {
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
	return func(record *neo4j.Record) error {
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
	return func(record *neo4j.Record) error {
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
	return func(record *neo4j.Record) error {
		value, getErr := getValueFunc(record)
		if getErr != nil {
			return getErr
		}
		if value == nil {
			info.DateTimeFields[fieldName].SetNull()
			return nil
		}
		var dateTimeValue time.Time
		switch typed := value.(type) {
		case neo4j.Time:
			dateTimeValue = typed.Time()
		case neo4j.LocalTime:
			dateTimeValue = typed.Time()
		case neo4j.LocalDateTime:
			dateTimeValue = typed.Time()
		case neo4j.Date:
			dateTimeValue = typed.Time()
		case time.Time:
			dateTimeValue = typed
		default:
			return fmt.Errorf(`value %v is not a datetime for field %v`, value, fieldName)
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
		return func(record *neo4j.Record) (interface{}, error) {
			value, exists := record.Get(element.Key)
			if !exists {
				return nil, nil
			}
			return value, nil
		}, nil
	case `List:String`, `List:Integer`, `List:Float`, `List:Boolean`, `List:Date`, `List:DateTime`:
		extractListFunc := func(record *neo4j.Record) ([]interface{}, error) {
			value, exists := record.Get(element.Key)
			if !exists {
				return nil, nil
			}
			valueList, ok := value.([]interface{})
			if !ok {
				return nil, fmt.Errorf(`path key %v for field %v is not a list, but is %T`, element.Key, field.Name, value)
			}
			return valueList, nil
		}
		return listTransferFunc(iterator, field, extractListFunc)
	case `Node`:
		extractNodeFunc := func(record *neo4j.Record) (neo4j.Node, error) {
			value, exists := record.Get(element.Key)
			if !exists {
				return emptyNode, nil
			}
			nodeValue, ok := value.(neo4j.Node)
			if !ok {
				return emptyNode, fmt.Errorf(`path key %v for field %v is not a Node as expected, but is %T`, element.Key, field.Name, value)
			}
			return nodeValue, nil
		}
		return nodeTransferFunc(iterator, field, extractNodeFunc)
	case `Relationship`:
		extractRelationshipFunc := func(record *neo4j.Record) (neo4j.Relationship, error) {
			value, exists := record.Get(element.Key)
			if !exists {
				return emptyRel, nil
			}
			relValue, ok := value.(neo4j.Relationship)
			if !ok {
				return emptyRel, fmt.Errorf(`path key %v for field %v is not a Relationship as expected, but is %T`, element.Key, field.Name, value)
			}
			return relValue, nil
		}
		return relationshipTransferFunc(iterator, field, extractRelationshipFunc)
	case `Path`:
		extractPathFunc := func(record *neo4j.Record) (neo4j.Path, error) {
			value, exists := record.Get(element.Key)
			if !exists {
				return neo4j.Path{}, nil
			}
			pathValue, ok := value.(neo4j.Path)
			if !ok {
				return neo4j.Path{}, fmt.Errorf(`path key %v for field %v is not a Path as expected, but is %T`, element.Key, field.Name, value)
			}
			return pathValue, nil
		}
		return pathTransferFunc(iterator, field, extractPathFunc)
	default:
		return nil, fmt.Errorf(`invalid data type '%v' for path in field '%v'`, element.DataType, field.Name)
	}
}

type extractNode func(record *neo4j.Record) (neo4j.Node, error)

func nodeTransferFunc(iterator *pathIterator, field Field, nodeExtractor extractNode) (GetValueFunc, error) {
	element, ok := iterator.NextField()
	if !ok {
		return nil, fmt.Errorf(`the path for field %v ends in a Node and not in a property data type`, field.Name)
	}
	switch element.Key {
	case `ID`:
		return func(record *neo4j.Record) (interface{}, error) {
			node, err := nodeExtractor(record)
			if err != nil {
				return nil, err
			}
			if node.Id < 0 {
				return nil, nil
			}
			return node.Id, nil
		}, nil
	case `Labels`:
		nodeFunc := func(record *neo4j.Record) ([]string, error) {
			node, err := nodeExtractor(record)
			if err != nil {
				return nil, err
			}
			return node.Labels, nil
		}
		return labelsTransferFunc(iterator, field, nodeFunc)
	case `Properties`:
		nodeFunc := func(record *neo4j.Record) (map[string]interface{}, error) {
			node, err := nodeExtractor(record)
			if err != nil {
				return nil, err
			}
			return node.Props, nil
		}
		return mapTransferFunc(iterator, field, nodeFunc)
	case `ToString`:
		return func(record *neo4j.Record) (interface{}, error) {
			node, err := nodeExtractor(record)
			if err != nil {
				return nil, err
			}
			str := ToString(node)
			return str, nil
		}, nil
	default:
		return nil, fmt.Errorf(`field %v has an invalid key '%v' for Node`, field.Name, element.Key)
	}
}

type extractRelationship func(record *neo4j.Record) (neo4j.Relationship, error)

func relationshipTransferFunc(iterator *pathIterator, field Field, relExtractor extractRelationship) (GetValueFunc, error) {
	element, ok := iterator.NextField()
	if !ok {
		return nil, fmt.Errorf(`the path for field %v ends in a Relationship and not in a property data type`, field.Name)
	}
	switch element.Key {
	case `ID`:
		return func(record *neo4j.Record) (interface{}, error) {
			relationship, err := relExtractor(record)
			if err != nil {
				return nil, err
			}
			if relationship.Id < 0 {
				return nil, nil
			}
			return relationship.Id, nil
		}, nil
	case `StartId`:
		return func(record *neo4j.Record) (interface{}, error) {
			relationship, err := relExtractor(record)
			if err != nil {
				return nil, err
			}
			if relationship.Id < 0 {
				return nil, nil
			}
			return relationship.StartId, nil
		}, nil
	case `EndId`:
		return func(record *neo4j.Record) (interface{}, error) {
			relationship, err := relExtractor(record)
			if err != nil {
				return nil, err
			}
			if relationship.Id < 0 {
				return nil, nil
			}
			return relationship.EndId, nil
		}, nil
	case `Type`:
		return func(record *neo4j.Record) (interface{}, error) {
			relationship, err := relExtractor(record)
			if err != nil {
				return nil, err
			}
			if relationship.Id < 0 {
				return nil, nil
			}
			return relationship.Type, nil
		}, nil
	case `Properties`:
		nodeFunc := func(record *neo4j.Record) (map[string]interface{}, error) {
			relationship, err := relExtractor(record)
			if err != nil {
				return nil, err
			}
			return relationship.Props, nil
		}
		return mapTransferFunc(iterator, field, nodeFunc)
	case `ToString`:
		return func(record *neo4j.Record) (interface{}, error) {
			relationship, err := relExtractor(record)
			if err != nil {
				return nil, err
			}
			str := ToString(relationship)
			return str, nil
		}, nil
	default:
		return nil, fmt.Errorf(`field %v has an invalid key '%v' for Relationship`, field.Name, element.Key)
	}
}

type extractPath func(record *neo4j.Record) (neo4j.Path, error)

func pathTransferFunc(iterator *pathIterator, field Field, extract extractPath) (GetValueFunc, error) {
	element, ok := iterator.NextField()
	if !ok {
		return nil, fmt.Errorf(`the path for field %v ends in a Path and not in a property data type`, field.Name)
	}

	switch element.Key {
	case `Nodes`:
		nodesFunc := func(record *neo4j.Record) ([]neo4j.Node, error) {
			extractedPath, err := extract(record)
			if err != nil {
				return nil, err
			}
			return extractedPath.Nodes, nil
		}
		return nodeListTransferFunc(iterator, field, nodesFunc)
	case `Relationships`:
		relsFunc := func(record *neo4j.Record) ([]neo4j.Relationship, error) {
			extractedPath, err := extract(record)
			if err != nil {
				return nil, err
			}
			return extractedPath.Relationships, nil
		}
		return relListTransferFunc(iterator, field, relsFunc)
	case `ToString`:
		return func(record *neo4j.Record) (interface{}, error) {
			extractedPath, err := extract(record)
			if err != nil {
				return nil, err
			}
			str := ToString(extractedPath)
			return str, nil
		}, nil
	default:
		return nil, fmt.Errorf(`field %v has an invalid key '%v' for Path`, field.Name, element.Key)
	}
}

type extractRelList func(record *neo4j.Record) ([]neo4j.Relationship, error)

var emptyRel = neo4j.Relationship{Id: -1}

func relListTransferFunc(iterator *pathIterator, field Field, extractList extractRelList) (GetValueFunc, error) {
	element, ok := iterator.NextField()
	if !ok {
		return nil, fmt.Errorf(`the path for field %v ends in a list of Relationships and not in a property data type`, field.Name)
	}
	switch element.Key {
	case `First`:
		relFunc := func(record *neo4j.Record) (neo4j.Relationship, error) {
			list, err := extractList(record)
			if err != nil {
				return emptyRel, err
			}
			if len(list) == 0 {
				return emptyRel, nil
			}
			return list[0], nil
		}
		return relationshipTransferFunc(iterator, field, relFunc)
	case `Last`:
		relFunc := func(record *neo4j.Record) (neo4j.Relationship, error) {
			list, err := extractList(record)
			if err != nil {
				return emptyRel, err
			}
			if len(list) == 0 {
				return emptyRel, nil
			}
			return list[len(list)-1], nil
		}
		return relationshipTransferFunc(iterator, field, relFunc)
	case `Count`:
		return func(record *neo4j.Record) (interface{}, error) {
			list, err := extractList(record)
			if err != nil {
				return 0, err
			}
			return int64(len(list)), nil
		}, nil
	default:
		if len(element.Key) < 7 || element.Key[:6] != `Index:` {
			return nil, fmt.Errorf(`field %v has an invalid key '%v' for List:Relationship`, field.Name, element.Key)
		}
		index, err := strconv.Atoi(element.Key[6:])
		if err != nil {
			return nil, fmt.Errorf(`field %v does not have a valid index in key '%v'`, field.Name, element.Key)
		}
		relFunc := func(record *neo4j.Record) (neo4j.Relationship, error) {
			list, getErr := extractList(record)
			if getErr != nil {
				return emptyRel, getErr
			}
			if len(list) <= index {
				return emptyRel, nil
			}
			return list[index], nil
		}
		return relationshipTransferFunc(iterator, field, relFunc)
	}
}

type extractNodeList func(record *neo4j.Record) ([]neo4j.Node, error)

var emptyNode = neo4j.Node{Id: -1}

func nodeListTransferFunc(iterator *pathIterator, field Field, extractList extractNodeList) (GetValueFunc, error) {
	element, ok := iterator.NextField()
	if !ok {
		return nil, fmt.Errorf(`the path for field %v ends in a list of Nodes and not in a property data type`, field.Name)
	}
	switch element.Key {
	case `First`:
		nodeFunc := func(record *neo4j.Record) (neo4j.Node, error) {
			list, err := extractList(record)
			if err != nil {
				return emptyNode, err
			}
			if len(list) == 0 {
				return emptyNode, nil
			}
			return list[0], nil
		}
		return nodeTransferFunc(iterator, field, nodeFunc)
	case `Last`:
		nodeFunc := func(record *neo4j.Record) (neo4j.Node, error) {
			list, err := extractList(record)
			if err != nil {
				return emptyNode, err
			}
			if len(list) == 0 {
				return emptyNode, nil
			}
			return list[len(list)-1], nil
		}
		return nodeTransferFunc(iterator, field, nodeFunc)
	case `Count`:
		return func(record *neo4j.Record) (interface{}, error) {
			list, err := extractList(record)
			if err != nil {
				return 0, err
			}
			return int64(len(list)), nil
		}, nil
	default:
		if len(element.Key) < 7 || element.Key[:6] != `Index:` {
			return nil, fmt.Errorf(`field %v has an invalid key '%v' for List:Node`, field.Name, element.Key)
		}
		index, err := strconv.Atoi(element.Key[6:])
		if err != nil {
			return nil, fmt.Errorf(`field %v does not have a valid index in key '%v'`, field.Name, element.Key)
		}
		nodeFunc := func(record *neo4j.Record) (neo4j.Node, error) {
			list, getErr := extractList(record)
			if getErr != nil {
				return emptyNode, getErr
			}
			if len(list) <= index {
				return emptyNode, nil
			}
			return list[index], nil
		}
		return nodeTransferFunc(iterator, field, nodeFunc)
	}
}

type extractLabels func(record *neo4j.Record) ([]string, error)

func labelsTransferFunc(iterator *pathIterator, field Field, extract extractLabels) (GetValueFunc, error) {
	element, ok := iterator.NextField()
	if !ok {
		return nil, fmt.Errorf(`the path for field %v ends in a list of strings and not in a property data type`, field.Name)
	}
	switch element.Key {
	case `Concatenate`:
		return func(record *neo4j.Record) (interface{}, error) {
			list, err := extract(record)
			if err != nil {
				return nil, err
			}
			return strings.Join(list, `,`), nil
		}, nil
	case `First`:
		return func(record *neo4j.Record) (interface{}, error) {
			list, err := extract(record)
			if err != nil {
				return nil, err
			}
			if len(list) == 0 {
				return nil, nil
			}
			return list[0], nil
		}, nil
	case `Last`:
		return func(record *neo4j.Record) (interface{}, error) {
			list, err := extract(record)
			if err != nil {
				return nil, err
			}
			if len(list) == 0 {
				return nil, nil
			}
			return list[len(list)-1], nil
		}, nil
	case `Count`:
		return func(record *neo4j.Record) (interface{}, error) {
			list, err := extract(record)
			if err != nil {
				return 0, err
			}
			return int64(len(list)), nil
		}, nil
	default:
		if len(element.Key) < 7 || element.Key[:6] != `Index:` {
			return nil, fmt.Errorf(`field %v has an invalid key '%v' for List:String`, field.Name, element.Key)
		}
		index, err := strconv.Atoi(element.Key[6:])
		if err != nil {
			return nil, fmt.Errorf(`field %v does not have a valid index in key '%v'`, field.Name, element.Key)
		}
		return func(record *neo4j.Record) (interface{}, error) {
			list, getErr := extract(record)
			if getErr != nil {
				return nil, getErr
			}
			if len(list) <= index {
				return nil, nil
			}
			return list[index], nil
		}, nil
	}
}

type extractList func(record *neo4j.Record) ([]interface{}, error)

func listTransferFunc(iterator *pathIterator, field Field, extract extractList) (GetValueFunc, error) {
	element, ok := iterator.NextField()
	if !ok {
		return nil, fmt.Errorf(`the path for field %v ends in a list of strings and not in a property data type`, field.Name)
	}

	switch element.Key {
	case `Concatenate`:
		return func(record *neo4j.Record) (interface{}, error) {
			list, err := extract(record)
			if err != nil {
				return nil, err
			}
			if len(list) == 0 {
				return ``, nil
			}
			var builder strings.Builder
			builder.WriteString(list[0].(string))
			for _, value := range list[1:] {
				builder.WriteByte(',')
				builder.WriteString(value.(string))
			}
			return builder.String(), nil
		}, nil
	case `First`:
		return func(record *neo4j.Record) (interface{}, error) {
			list, err := extract(record)
			if err != nil {
				return nil, err
			}
			if len(list) == 0 {
				return nil, nil
			}
			return list[0], nil
		}, nil
	case `Last`:
		return func(record *neo4j.Record) (interface{}, error) {
			list, err := extract(record)
			if err != nil {
				return nil, err
			}
			if len(list) == 0 {
				return nil, nil
			}
			return list[len(list)-1], nil
		}, nil
	case `Count`:
		return func(record *neo4j.Record) (interface{}, error) {
			list, err := extract(record)
			if err != nil {
				return 0, err
			}
			return int64(len(list)), nil
		}, nil
	default:
		if len(element.Key) < 7 || element.Key[:6] != `Index:` {
			return nil, fmt.Errorf(`field %v has an invalid key '%v' for list`, field.Name, element.Key)
		}
		index, err := strconv.Atoi(element.Key[6:])
		if err != nil {
			return nil, fmt.Errorf(`field %v does not have a valid index in key '%v'`, field.Name, element.Key)
		}
		return func(record *neo4j.Record) (interface{}, error) {
			list, getErr := extract(record)
			if getErr != nil {
				return nil, getErr
			}
			if len(list) <= index {
				return nil, nil
			}
			return list[index], nil
		}, nil
	}
}

type extractMap func(record *neo4j.Record) (map[string]interface{}, error)

func mapTransferFunc(iterator *pathIterator, field Field, extract extractMap) (GetValueFunc, error) {
	element, ok := iterator.NextField()
	if !ok {
		return nil, fmt.Errorf(`the path for field %v ends in a list of strings and not in a property data type`, field.Name)
	}

	switch element.DataType {
	case `String`, `Integer`, `Boolean`, `Float`, `Date`, `DateTime`:
		return func(record *neo4j.Record) (interface{}, error) {
			extractedMap, err := extract(record)
			if err != nil {
				return nil, err
			}
			value, hasKey := extractedMap[element.Key]
			if !hasKey {
				return nil, nil
			}
			return value, nil
		}, nil
	case `List:String`, `List:Integer`, `List:Boolean`, `List:Float`, `List:Date`, `List:DateTime`:
		listFunc := func(record *neo4j.Record) ([]interface{}, error) {
			extractedMap, err := extract(record)
			if err != nil {
				return nil, err
			}
			value, hasKey := extractedMap[element.Key]
			if !hasKey {
				return nil, nil
			}
			listValue, convertOk := value.([]interface{})
			if !convertOk {
				return nil, fmt.Errorf(`map value with key '%v' on field %v is not a list; it is %T`, element.Key, field.Name, value)
			}
			return listValue, nil
		}
		return listTransferFunc(iterator, field, listFunc)
	default:
		return nil, fmt.Errorf(`field %v has an invalid data type '%v' for Map`, field.Name, element.DataType)
	}
}
