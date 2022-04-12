package main

import "C"
import (
	"github.com/tlarsendataguy/goalteryx/sdk"
	"github.com/tlarsendataguy/graphyx/delete"
	"github.com/tlarsendataguy/graphyx/input"
	"github.com/tlarsendataguy/graphyx/output"
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

//export Neo4jDelete
func Neo4jDelete(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &delete.Neo4jDelete{}
	return C.long(sdk.RegisterTool(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}

func main() {}
