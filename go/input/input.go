package input

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/tlarsen7572/goalteryx/sdk"
	"github.com/tlarsen7572/graphyx/bolt_url"
)

type Neo4jInput struct {
	provider   sdk.Provider
	output     sdk.OutputAnchor
	outObjects OutgoingObjects
	config     Configuration
}

func (i *Neo4jInput) Init(provider sdk.Provider) {
	var err error
	i.provider = provider
	i.output = provider.GetOutputAnchor(`Output`)
	i.config, err = DecodeConfig(provider.ToolConfig())
	if err != nil {
		i.provider.Io().Error(err.Error())
		return
	}
	i.outObjects, err = CreateOutgoingObjects(i.config.Fields)
	if err != nil {
		i.provider.Io().Error(err.Error())
	}
}

func (i *Neo4jInput) OnInputConnectionOpened(_ sdk.InputConnection) {
	panic("should never be called for input tools")
}

func (i *Neo4jInput) OnRecordPacket(_ sdk.InputConnection) {
	panic("should never be called for input tools")
}

func (i *Neo4jInput) OnComplete() {
	i.output.Open(i.outObjects.RecordInfo)
	if i.provider.Environment().UpdateOnly() {
		return
	}

	boltUrl, err := bolt_url.GetBoltUrl(i.config.ConnStr)
	if err != nil {
		i.provider.Io().Error(fmt.Sprintf(`error connecting to Neo4j: %v`, err.Error()))
	}

	driver, err := neo4j.NewDriver(boltUrl, neo4j.BasicAuth(i.config.Username, i.config.Password, ""))
	if err != nil {
		i.provider.Io().Error(fmt.Sprintf(`expected no error but got: %v`, err.Error()))
		return
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	_, err = session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, txErr := tx.Run(i.config.Query, nil)
		if txErr != nil {
			return nil, txErr
		}
		for result.Next() {
			record := result.Record()
			for _, transferFunc := range i.outObjects.TransferFuncs {
				err = transferFunc(record)
				if err != nil {
					i.provider.Io().Error(err.Error())
					return nil, err
				}
			}
			i.output.Write()
		}

		if txErr = result.Err(); txErr != nil {
			return nil, txErr
		}

		return result.Consume()
	})
}
