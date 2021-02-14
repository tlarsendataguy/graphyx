package graphyx

import (
	"github.com/tlarsen7572/goalteryx/sdk"
)

type neo4jOutput struct {
}

func (o *neo4jOutput) Init(provider sdk.Provider) {
	panic("implement me")
}

func (o *neo4jOutput) OnInputConnectionOpened(connection sdk.InputConnection) {
	panic("implement me")
}

func (o *neo4jOutput) OnRecordPacket(connection sdk.InputConnection) {
	panic("implement me")
}

func (o *neo4jOutput) OnComplete() {
	panic("implement me")
}
