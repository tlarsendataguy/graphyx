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

type Neo4jOutput struct {
	query    string
	params   []map[string]interface{}
	config   Configuration
	provider sdk.Provider
}

func (o *Neo4jOutput) Init(provider sdk.Provider) {
	var err error
	o.provider = provider
	o.config, err = decodeConfig(provider.ToolConfig())
	if err != nil {
		provider.Io().Error(err.Error())
		return
	}
	if o.config.ExportObject == `Node` {
		o.generateNodeQuery()
		return
	}
	if o.config.ExportObject == `Relationship` {
		o.generateRelationshipQuery()
		return
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
