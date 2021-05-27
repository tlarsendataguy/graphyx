package delete

import (
	"github.com/tlarsen7572/goalteryx/sdk"
)

type Configuration struct {
	ConnStr        string
	Username       string
	Password       string
	Database       string
	DeleteObject   string
	BatchSize      int
	NodeLabel      string
	NodeIdFields   []string
	RelLabel       string
	RelPropFields  []string
	RelLeftLabel   string
	RelLeftFields  []map[string]interface{}
	RelRightLabel  string
	RelRightFields []map[string]interface{}
}

type Neo4jDelete struct {
}

func (d *Neo4jDelete) Init(provider sdk.Provider) {
	panic("implement me")
}

func (d *Neo4jDelete) OnInputConnectionOpened(connection sdk.InputConnection) {
	panic("implement me")
}

func (d *Neo4jDelete) OnRecordPacket(connection sdk.InputConnection) {
	panic("implement me")
}

func (d *Neo4jDelete) OnComplete() {
	panic("implement me")
}
