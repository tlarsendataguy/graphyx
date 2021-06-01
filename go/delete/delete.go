package delete

import (
	"encoding/json"
	"github.com/tlarsen7572/goalteryx/sdk"
)

type Configuration struct {
	ConnStr             string
	Username            string
	Password            string
	Database            string
	DeleteObject        string
	BatchSize           int
	NodeLabel           string
	NodeIdFields        []string
	RelType             string
	RelFields           []string
	RelLeftLabel        string
	RelLeftAyxFields    []string
	RelLeftNeo4jFields  []string
	RelRightLabel       string
	RelRightAyxFields   []string
	RelRightNeo4jFields []string
}

type Neo4jDelete struct {
	provider sdk.Provider
	config   Configuration
	doExport bool
	query    string
}

func (d *Neo4jDelete) Init(provider sdk.Provider) {
	d.provider = provider
	var err error
	err = json.Unmarshal([]byte(provider.ToolConfig()), &d.config)
	if err != nil {
		d.provider.Io().Error(err.Error())
		return
	}
	switch d.config.DeleteObject {
	case `Node`:
		d.query = GenerateDeleteNodes(&DeleteNodesProperties{
			Label:    d.config.NodeLabel,
			IdFields: d.config.NodeIdFields,
		})
	case `Relationship`:
		d.query, err = GenerateDeleteRelationships(&DeleteRelationshipsProperties{
			RelType:                d.config.RelType,
			RelFields:              d.config.RelFields,
			LeftNodeLabel:          d.config.RelLeftLabel,
			LeftNodeAlteryxFields:  d.config.RelLeftAyxFields,
			LeftNodeNeo4jFields:    d.config.RelLeftNeo4jFields,
			RightNodeLabel:         d.config.RelRightLabel,
			RightNodeAlteryxFields: d.config.RelRightAyxFields,
			RightNodeNeo4jFields:   d.config.RelRightNeo4jFields,
		})
		if err != nil {
			d.provider.Io().Error(err.Error())
		}
	default:
		d.provider.Io().Error(`the DeleteObject property is not valid, expected either 'Node' or 'Relationship'`)
	}
	d.doExport = true
}

func (d *Neo4jDelete) OnInputConnectionOpened(connection sdk.InputConnection) {
	panic("implement me")
}

func (d *Neo4jDelete) OnRecordPacket(connection sdk.InputConnection) {
	panic("implement me")
}

func (d *Neo4jDelete) OnComplete() {}
