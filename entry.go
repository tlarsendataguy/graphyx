package graphyx

import "C"
import (
	"github.com/tlarsen7572/goalteryx/sdk"
	"unsafe"
)

func Neo4jInput(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &neo4jInput{}
	return C.long(sdk.RegisterTool(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}

func Neo4jOutput(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &neo4jOutput{}
	return C.long(sdk.RegisterTool(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}
