package output

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/tlarsendataguy/goalteryx/sdk"
	"github.com/tlarsendataguy/graphyx/bolt_url"
	"github.com/tlarsendataguy/graphyx/util"
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
	RelIdFields    []string
	RelPropFields  []string
	RelLeftLabel   string
	RelLeftFields  []map[string]interface{}
	RelRightLabel  string
	RelRightFields []map[string]interface{}
}

type Neo4jOutput struct {
	query            string
	config           Configuration
	provider         sdk.Provider
	copier           []util.CopyData
	outputFields     []string
	batch            []map[string]interface{}
	currentBatchSize int
	driver           neo4j.Driver
	session          neo4j.Session
	doExport         bool
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
		o.outputFields = append(o.outputFields, o.config.RelIdFields...)
	}
	outputFieldLen := len(o.outputFields)
	for index := range o.batch {
		o.batch[index] = make(map[string]interface{}, outputFieldLen)
	}

	if !o.provider.Environment().UpdateOnly() {
		o.doExport = true
	}
}

func (o *Neo4jOutput) OnInputConnectionOpened(connection sdk.InputConnection) {
	if !o.doExport {
		return
	}

	url, err := bolt_url.GetBoltUrl(o.config.ConnStr)
	if err != nil {
		o.error(err.Error())
		return
	}

	var copier util.CopyData
	incomingInfo := connection.Metadata()
	for _, field := range o.outputFields {
		copier, err = util.FindFieldAndGenerateCopier(field, incomingInfo)
		if err != nil {
			o.error(fmt.Sprintf(`field %v was not contained in the record`, field))
			return
		}
		o.copier = append(o.copier, copier)
	}

	username, password := util.GetCredentials(o.config.Username, o.config.Password, o.provider)
	o.driver, err = neo4j.NewDriver(url, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		o.error(err.Error())
		return
	}
	o.session = o.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: o.config.Database})
	err = o.driver.VerifyConnectivity()
	if err != nil {
		o.error(err.Error())
	}
}

func (o *Neo4jOutput) OnRecordPacket(connection sdk.InputConnection) {
	if !o.doExport {
		return
	}

	packet := connection.Read()
	for packet.Next() {
		if o.currentBatchSize >= o.config.BatchSize {
			o.sendBatch()
		}
		copyFrom := packet.Record()
		copyTo := o.batch[o.currentBatchSize]
		for _, copyData := range o.copier {
			err := copyData(copyFrom, copyTo)
			if err != nil {
				o.provider.Io().Error(err.Error())
			}
		}
		o.currentBatchSize++
	}
	o.provider.Io().UpdateProgress(connection.Progress())
}

func (o *Neo4jOutput) sendBatch() {
	_, err := o.session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(o.query, map[string]interface{}{`batch`: o.batch[:o.currentBatchSize]})
	})
	if err != nil {
		o.error(err.Error())
		return
	}
	o.currentBatchSize = 0
}

func (o *Neo4jOutput) OnComplete() {
	if o.currentBatchSize > 0 && o.doExport {
		o.sendBatch()
	}
	if o.session != nil {
		o.session.Close()
	}
	if o.driver != nil {
		o.driver.Close()
	}
	o.provider.Io().UpdateProgress(1.0)
}

func (o *Neo4jOutput) error(msg string) {
	o.doExport = false
	o.provider.Io().Error(msg)
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
		IdFields:           o.config.RelIdFields,
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
