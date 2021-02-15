package input

import (
	"github.com/tlarsen7572/goalteryx/sdk"
)

type Neo4jInput struct {
}

func (i *Neo4jInput) Init(provider sdk.Provider) {
	panic("implement me")
}

func (i *Neo4jInput) OnInputConnectionOpened(connection sdk.InputConnection) {
	panic("implement me")
}

func (i *Neo4jInput) OnRecordPacket(connection sdk.InputConnection) {
	panic("implement me")
}

func (i *Neo4jInput) OnComplete() {
	panic("implement me")
}
