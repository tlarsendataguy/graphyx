package output

import (
	"fmt"
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
		return
	}
	if o.config.ExportObject == `Relationship` {
		o.generateRelationshipQuery()
		for _, field := range o.config.RelLeftFields {
			for key := range field {
				o.outputFields = append(o.outputFields, key)
			}
		}
		for _, field := range o.config.RelRightFields {
			for key := range field {
				o.outputFields = append(o.outputFields, key)
			}
		}
		o.outputFields = append(o.outputFields, o.config.RelPropFields...)
		return
	}
	outputFieldLen := len(o.outputFields)
	for index := range o.batch {
		o.batch[index] = make(map[string]interface{}, outputFieldLen)
	}
}

func (o *Neo4jOutput) OnInputConnectionOpened(connection sdk.InputConnection) {
	return
}

func (o *Neo4jOutput) OnRecordPacket(connection sdk.InputConnection) {
	return
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
