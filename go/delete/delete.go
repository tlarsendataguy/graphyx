package delete

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/tlarsendataguy/goalteryx/sdk"
	"github.com/tlarsendataguy/graphyx/bolt_url"
	"github.com/tlarsendataguy/graphyx/util"
)

type xmlConfig struct {
	JSON string `xml:",text"`
}

type Configuration struct {
	ConnStr        string
	Username       string
	Password       string
	Database       string
	DeleteObject   string
	BatchSize      int
	NodeLabel      string
	NodeIdFields   []string
	RelType        string
	RelFields      []string
	RelLeftLabel   string
	RelLeftFields  []map[string]interface{}
	RelRightLabel  string
	RelRightFields []map[string]interface{}
}

type Neo4jDelete struct {
	provider         sdk.Provider
	config           Configuration
	doExport         bool
	query            string
	copiers          []util.CopyData
	requiredFields   []string
	driver           neo4j.Driver
	session          neo4j.Session
	batch            []map[string]interface{}
	currentBatchSize int
}

func (d *Neo4jDelete) Init(provider sdk.Provider) {
	d.provider = provider
	var rawConfig xmlConfig
	err := xml.Unmarshal([]byte(provider.ToolConfig()), &rawConfig)
	if err != nil {
		d.error(fmt.Sprintf(`error parsing XML configuration: %v`, err.Error()))
		return
	}
	err = json.Unmarshal([]byte(rawConfig.JSON), &d.config)
	if err != nil {
		d.error(fmt.Sprintf(`error parsing JSON configuration: %v`, err.Error()))
		return
	}

	switch d.config.DeleteObject {
	case `Node`:
		d.query = GenerateDeleteNodes(&DeleteNodesProperties{
			Label:    d.config.NodeLabel,
			IdFields: d.config.NodeIdFields,
		})
		d.requiredFields = d.config.NodeIdFields
	case `Relationship`:
		d.query, err = GenerateDeleteRelationships(&DeleteRelationshipsProperties{
			RelType:         d.config.RelType,
			RelFields:       d.config.RelFields,
			LeftNodeLabel:   d.config.RelLeftLabel,
			LeftNodeFields:  d.config.RelLeftFields,
			RightNodeLabel:  d.config.RelRightLabel,
			RightNodeFields: d.config.RelRightFields,
		})
		if err != nil {
			d.error(err.Error())
			return
		}
		for _, field := range d.config.RelFields {
			d.requiredFields = append(d.requiredFields, field)
		}
		for _, fieldLists := range [][]map[string]interface{}{d.config.RelLeftFields, d.config.RelRightFields} {
			for _, fieldList := range fieldLists {
				for key := range fieldList {
					d.requiredFields = append(d.requiredFields, key)
				}
			}
		}
	default:
		d.error(`the DeleteObject property is not valid, expected either 'Node' or 'Relationship'`)
		return
	}

	d.batch = make([]map[string]interface{}, d.config.BatchSize)
	numFields := len(d.requiredFields)
	for index := range d.batch {
		d.batch[index] = make(map[string]interface{}, numFields)
	}

	if !d.provider.Environment().UpdateOnly() {
		d.doExport = true
	}
}

func (d *Neo4jDelete) OnInputConnectionOpened(connection sdk.InputConnection) {
	if !d.doExport {
		return
	}

	url, err := bolt_url.GetBoltUrl(d.config.ConnStr)
	if err != nil {
		d.error(err.Error())
		return
	}

	var copier util.CopyData
	incomingInfo := connection.Metadata()
	for _, field := range d.requiredFields {
		copier, err = util.FindFieldAndGenerateCopier(field, incomingInfo)
		if err != nil {
			d.error(fmt.Sprintf(`field %v was not contained in the record`, field))
			return
		}
		d.copiers = append(d.copiers, copier)
	}

	username, password := util.GetCredentials(d.config.Username, d.config.Password, d.provider)
	d.driver, err = neo4j.NewDriver(url, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		d.error(err.Error())
		return
	}
	err = d.driver.VerifyConnectivity()
	if err != nil {
		d.error(err.Error())
		return
	}
	d.session = d.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: d.config.Database})
}

func (d *Neo4jDelete) OnRecordPacket(connection sdk.InputConnection) {
	if !d.doExport {
		return
	}

	packet := connection.Read()
	for packet.Next() {
		if d.currentBatchSize >= d.config.BatchSize {
			d.sendBatch()
		}
		copyFrom := packet.Record()
		copyTo := d.batch[d.currentBatchSize]
		for _, copyData := range d.copiers {
			copyData(copyFrom, copyTo)
		}
		d.currentBatchSize++
	}
	d.provider.Io().UpdateProgress(connection.Progress())
}

func (d *Neo4jDelete) OnComplete() {
	if d.currentBatchSize > 0 && d.doExport {
		d.sendBatch()
	}
	if d.session != nil {
		d.session.Close()
	}
	if d.driver != nil {
		d.driver.Close()
	}
	d.provider.Io().UpdateProgress(1.0)
}

func (d *Neo4jDelete) sendBatch() {
	_, err := d.session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(d.query, map[string]interface{}{`batch`: d.batch[:d.currentBatchSize]})
	})
	if err != nil {
		d.error(err.Error())
		return
	}
	d.currentBatchSize = 0
}

func (d *Neo4jDelete) error(msg string) {
	d.provider.Io().Error(msg)
	d.doExport = false
}
