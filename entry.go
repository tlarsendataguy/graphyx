package graphyx

import "C"
import (
	"github.com/tlarsen7572/goalteryx/sdk"
	"github.com/tlarsen7572/graphyx/input"
	"github.com/tlarsen7572/graphyx/output"
	"unsafe"
)

//export Neo4jInput
func Neo4jInput(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &input.Neo4jInput{}
	return C.long(sdk.RegisterTool(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}

//export Neo4jOutput
func Neo4jOutput(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &output.Neo4jOutput{}
	return C.long(sdk.RegisterTool(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}
