package output

import (
	"github.com/tlarsen7572/goalteryx/sdk"
)

type Neo4jOutput struct {
}

func (o *Neo4jOutput) Init(provider sdk.Provider) {
	panic("implement me")
}

func (o *Neo4jOutput) OnInputConnectionOpened(connection sdk.InputConnection) {
	panic("implement me")
}

func (o *Neo4jOutput) OnRecordPacket(connection sdk.InputConnection) {
	panic("implement me")
}

func (o *Neo4jOutput) OnComplete() {
	panic("implement me")
}
