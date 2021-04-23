package output

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/tlarsen7572/goalteryx/sdk"
)

type Configuration struct {
	ConnStr        string
	Username       string
	Password       string
	Database       string
	ExportObject   string
	BatchSize      int
	NodeLabel      string
	NodeIdFields   []string
	NodePropFields []string
	RelLabel       string
	RelPropFields  []string
	RelLeftLabel   string
	RelLeftFields  []map[string]interface{}
	RelRightLabel  string
	RelRightFields []map[string]interface{}
}

type CopyData func(sdk.Record, map[string]interface{})

type Neo4jOutput struct {
	query            string
	config           Configuration
	provider         sdk.Provider
	copier           []CopyData
	outputFields     []string
	batch            []map[string]interface{}
	currentBatchSize int
	driver           neo4j.Driver
	session          neo4j.Session
}

func (o *Neo4jOutput) Init(provider sdk.Provider) {
	var err error
	o.provider = provider
	o.config, err = decodeConfig(provider.ToolConfig())
	if err != nil {
		provider.Io().Error(err.Error())
		return
	}
	o.batch = make([]map[string]interface{}, o.config.BatchSize)
	if o.config.ExportObject == `Node` {
		o.generateNodeQuery()
		o.outputFields = append(o.config.NodeIdFields, o.config.NodePropFields...)
	}
	if o.config.ExportObject == `Relationship` {
		o.generateRelationshipQuery()
		for _, field := range o.config.RelLeftFields {
			for ayxField := range field {
				o.outputFields = append(o.outputFields, ayxField)
			}
		}
		for _, field := range o.config.RelRightFields {
			for ayxField := range field {
				o.outputFields = append(o.outputFields, ayxField)
			}
		}
		o.outputFields = append(o.outputFields, o.config.RelPropFields...)
	}
	outputFieldLen := len(o.outputFields)
	for index := range o.batch {
		o.batch[index] = make(map[string]interface{}, outputFieldLen)
	}
}

func (o *Neo4jOutput) findFieldAndGenerateCopier(field string, incomingInfo sdk.IncomingRecordInfo) bool {
	var copier CopyData
	for _, incomingField := range incomingInfo.Fields() {
		if field == incomingField.Name {
			switch incomingField.Type {
			case `Byte`, `Int16`, `Int32`, `Int64`:
				intField, _ := incomingInfo.GetIntField(field)
				getInt := intField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) {
					value, isNull := getInt(copyFrom)
					if isNull {
						copyTo[field] = nil
						return
					}
					copyTo[field] = value
				}
				o.copier = append(o.copier, copier)
				return true
			case `String`, `WString`, `V_String`, `V_WString`:
				stringField, _ := incomingInfo.GetStringField(field)
				getString := stringField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) {
					value, isNull := getString(copyFrom)
					if isNull {
						copyTo[field] = nil
						return
					}
					copyTo[field] = value
				}
				o.copier = append(o.copier, copier)
				return true
			case `Bool`:
				boolField, _ := incomingInfo.GetBoolField(field)
				getBool := boolField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) {
					value, isNull := getBool(copyFrom)
					if isNull {
						copyTo[field] = nil
						return
					}
					copyTo[field] = value
				}
				o.copier = append(o.copier, copier)
				return true
			case `Date`, `DateTime`:
				timeField, _ := incomingInfo.GetTimeField(field)
				getTime := timeField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) {
					value, isNull := getTime(copyFrom)
					if isNull {
						copyTo[field] = nil
						return
					}
					copyTo[field] = value
				}
				o.copier = append(o.copier, copier)
				return true
			case `Float`, `Double`, `FixedDecimal`:
				floatField, _ := incomingInfo.GetFloatField(field)
				getFloat := floatField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) {
					value, isNull := getFloat(copyFrom)
					if isNull {
						copyTo[field] = nil
						return
					}
					copyTo[field] = value
				}
				o.copier = append(o.copier, copier)
				return true
			}
		}
	}
	return false
}

func (o *Neo4jOutput) OnInputConnectionOpened(connection sdk.InputConnection) {
	incomingInfo := connection.Metadata()
	for _, field := range o.outputFields {
		ok := o.findFieldAndGenerateCopier(field, incomingInfo)
		if !ok {
			o.provider.Io().Error(fmt.Sprintf(`field %v was not contained in the record`, field))
			return
		}
	}
}

func (o *Neo4jOutput) OnRecordPacket(connection sdk.InputConnection) {
	packet := connection.Read()
	for packet.Next() {
		if o.currentBatchSize >= o.config.BatchSize {
			o.currentBatchSize = 0
		}
		copyFrom := packet.Record()
		copyTo := o.batch[o.currentBatchSize]
		for _, copyData := range o.copier {
			copyData(copyFrom, copyTo)
		}
		o.currentBatchSize++
	}
}

func (o *Neo4jOutput) OnComplete() {
	return
}

func (o *Neo4jOutput) Config() Configuration {
	return o.config
}

func (o *Neo4jOutput) Query() string {
	return o.query
}

func (o *Neo4jOutput) generateNodeQuery() {
	var err error
	nodeConfig := &NodeConfig{
		Label:      o.config.NodeLabel,
		IdFields:   o.config.NodeIdFields,
		PropFields: o.config.NodePropFields,
	}
	o.query, err = NodeQuery(nodeConfig)
	if err != nil {
		o.provider.Io().Error(err.Error())
		return
	}
}

func (o *Neo4jOutput) generateRelationshipQuery() {
	var err error
	var leftAlteryxFields []string
	var leftNeo4jFields []string
	var rightAlteryxFields []string
	var rightNeo4jFields []string

	leftAlteryxFields, leftNeo4jFields, err = fieldsToAyxAndNeo4jLists(o.config.RelLeftFields)
	if err != nil {
		o.provider.Io().Error(err.Error())
		return
	}
	rightAlteryxFields, rightNeo4jFields, err = fieldsToAyxAndNeo4jLists(o.config.RelRightFields)
	if err != nil {
		o.provider.Io().Error(err.Error())
		return
	}

	relConfig := &RelationshipConfig{
		LeftLabel:          o.config.RelLeftLabel,
		LeftAlteryxFields:  leftAlteryxFields,
		LeftNeo4jFields:    leftNeo4jFields,
		RightLabel:         o.config.RelRightLabel,
		RightAlteryxFields: rightAlteryxFields,
		RightNeo4jFields:   rightNeo4jFields,
		Label:              o.config.RelLabel,
		PropFields:         o.config.RelPropFields,
	}
	o.query, err = RelationshipQuery(relConfig)
	if err != nil {
		o.provider.Io().Error(err.Error())
		return
	}
}

func (o *Neo4jOutput) Batch() []map[string]interface{} {
	return o.batch
}

func (o *Neo4jOutput) OutputFields() []string {
	return o.outputFields
}

func (o *Neo4jOutput) CurrentRecords() []map[string]interface{} {
	return o.batch[:o.currentBatchSize]
}

func fieldsToAyxAndNeo4jLists(fields []map[string]interface{}) ([]string, []string, error) {
	var alteryxFields []string
	var neo4jFields []string
	for _, field := range fields {
		for alteryxField, neo4jField := range field {
			alteryxFields = append(alteryxFields, alteryxField)
			switch f := neo4jField.(type) {
			case string:
				neo4jFields = append(neo4jFields, f)
			default:
				return nil, nil, fmt.Errorf(`the Neo4j field mapping for '%v' is not a string; the tool configuration is not formatted properly`, alteryxField)
			}
		}
	}
	return alteryxFields, neo4jFields, nil
}
