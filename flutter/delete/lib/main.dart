import 'dart:typed_data';

import 'package:delete/connection_controls.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:delete/bloc.dart';
import 'package:delete/configuration.dart';
import 'package:delete/decode_config.dart';
import 'package:delete/field_mapper.dart';
import 'package:delete/field_selector.dart';
import 'package:delete/material_icons.dart';

Future<ByteData> fontFileToByteData(List<int> file) async {
  return ByteData.sublistView(Uint8List.fromList(file));
}

List<String> lazyLoadIncomingFields(){
  List<String> incomingFields = [];
  for (var field in getIncomingFields) {
    if (field.strType == 'Blob' || field.strType == 'SpatialObj') continue;
    incomingFields.add(field.strName);
  }
  return incomingFields;
}

void main() async {
  while (true) {
    if (configurationLoaded) {
      break;
    }
    await Future.delayed(const Duration(milliseconds: 100));
  }
  var appState = decodeConfig(configuration, lazyLoadIncomingFields);
  registerSaveConfigCallback(appState.getConfig);
  registerDecryptCallback(appState.callbackDecryptedPassword);
  registerEncryptCallback(appState.callbackEncryptedPassword);

  var materialLoader = FontLoader("MaterialIcons");
  materialLoader.addFont(fontFileToByteData(materialIcons));
  await materialLoader.load();

  runApp(BlocProvider<Configuration>(
    child: MyApp(),
    bloc: appState,
  ));
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Neo4j Output',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        accentColor: Colors.green,
      ),
      home: Scaffold(
        body: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Center(
            child: Controls(),
          ),
        )
      ),
    );
  }
}

class Controls extends StatefulWidget {
  Controls({Key key}) : super(key: key);

  _ControlsState createState() => _ControlsState();
}

class _ControlsState extends State<Controls> {
  Configuration config;
  TextEditingController batchSizeController;

  void batchSizeChanged (String value) {
    var intValue = int.tryParse(value);
    if (intValue == null) {
      return;
    }
    config.batchSize = intValue;
  }

  void initState() {
    config = BlocProvider.of<Configuration>(context);
    batchSizeController = TextEditingController(text: config.batchSize.toString());
    super.initState();
  }

  Widget build(BuildContext context) {
    return ListView(
      children: [
        ConnectionControls(),
        Card(
          elevation: 12,
          child: Padding(
            padding: const EdgeInsets.all(8),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                TextField(controller: batchSizeController, decoration: InputDecoration(labelText: "batch  size"), onChanged: batchSizeChanged, inputFormatters: [FilteringTextInputFormatter.allow(RegExp(r'[0-9]'))]),
                ExportObjectSelector(()=>setState((){})),
              ],
            ),
          ),
        ),
        NodeOrRelationshipConfig(config.deleteObject),
      ],
    );
  }
}

class ExportObjectSelector extends StatefulWidget {
  ExportObjectSelector(this.onChanged);
  final VoidCallback onChanged;

  createState() => _ExportObjectSelectorState();
}

class _ExportObjectSelectorState extends State<ExportObjectSelector> {
  Configuration config;

  void exportObjectChanged (String value) {
    config.deleteObject = value;
    setState(widget.onChanged);
  }

  initState() {
    config = BlocProvider.of<Configuration>(context);
    super.initState();
  }

  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        SizedBox(height: 20),
        Text("delete object:", textScaleFactor: 0.9),
        DropdownButton<String>(
          hint: Text("delete object"),
          items: [
            DropdownMenuItem(child: Text("Node"), value: "Node"),
            DropdownMenuItem(child: Text("Relationship"), value: "Relationship"),
          ],
          value: config.deleteObject,
          onChanged: exportObjectChanged,
        ),
      ],
    );
  }
}

class NodeOrRelationshipConfig extends StatelessWidget {
  NodeOrRelationshipConfig(this.exportObject);
  final String exportObject;

  Widget build(BuildContext context) {
    if (exportObject == 'Node') {
      return NodeConfig();
    }
    if (exportObject == 'Relationship') {
      return RelationshipConfig();
    }
    return Text("Invalid delete object");
  }
}

class NodeConfig extends StatefulWidget {
  State<StatefulWidget> createState() => _NodeConfigState();
}

class _NodeConfigState extends State<NodeConfig> {
  Configuration config;
  TextEditingController nodeLabelController;

  void nodeLabelChanged(String value) {
    config.nodeLabel = value;
  }

  initState(){
    config = BlocProvider.of<Configuration>(context);
    nodeLabelController = TextEditingController(text: config.nodeLabel);
    super.initState();
  }

  Widget build(BuildContext context) {
    return Card(
      elevation: 12,
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            TextField(controller: nodeLabelController, decoration: InputDecoration(labelText: "node label"), onChanged: nodeLabelChanged),
            FieldSelector(source: config.nodeIdFields, label: "node ID fields"),
          ],
        ),
      ),
    );
  }
}

class RelationshipConfig extends StatefulWidget {
  State<StatefulWidget> createState() => _RelationshipConfigState();
}

class _RelationshipConfigState extends State<RelationshipConfig> {
  Configuration config;
  TextEditingController relLabelController;
  TextEditingController relLeftLabelController;
  TextEditingController relRightLabelController;

  void relTypeChanged(String value) {
    config.relType = value;
  }

  void relLeftLabelChanged(String value) {
    config.relLeftLabel = value;
  }

  void relRightLabelChanged(String value) {
    config.relRightLabel = value;
  }

  initState(){
    config = BlocProvider.of<Configuration>(context);
    relLabelController = TextEditingController(text: config.relType);
    relLeftLabelController = TextEditingController(text: config.relLeftLabel);
    relRightLabelController = TextEditingController(text: config.relRightLabel);
    super.initState();
  }
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Card(
          elevation: 12,
          child: Padding(
            padding: const EdgeInsets.all(8),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                TextField(controller: relLabelController, decoration: InputDecoration(labelText: "relationship type"), onChanged: relTypeChanged),
                FieldSelector(source: config.relFields, label: "match the following properties of the relationship"),
              ],
            ),
          ),
        ),
        Card(
          elevation: 12,
          child: Padding(
            padding: const EdgeInsets.all(8),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                TextField(controller: relLeftLabelController, decoration: InputDecoration(labelText: "left node label"), onChanged: relLeftLabelChanged),
                FieldMapper(source: config.relLeftFields, label: "match the following properties of the left node"),
              ],
            ),
          ),
        ),
        Card(
          elevation: 12,
          child: Padding(
            padding: const EdgeInsets.all(8),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                TextField(controller: relRightLabelController, decoration: InputDecoration(labelText: "right node label"), onChanged: relRightLabelChanged),
                FieldMapper(source: config.relRightFields, label: "match the following properties of the right node"),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
