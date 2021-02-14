package graphyx

import (
	"github.com/tlarsen7572/goalteryx/sdk"
)

type neo4jInput struct {
}

func (i *neo4jInput) Init(provider sdk.Provider) {
	panic("implement me")
}

func (i *neo4jInput) OnInputConnectionOpened(connection sdk.InputConnection) {
	panic("implement me")
}

func (i *neo4jInput) OnRecordPacket(connection sdk.InputConnection) {
	panic("implement me")
}

func (i *neo4jInput) OnComplete() {
	panic("implement me")
}
