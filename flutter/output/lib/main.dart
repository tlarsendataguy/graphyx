import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:output/bloc.dart';
import 'package:output/configuration.dart';
import 'package:output/decode_config.dart';

void main() {
  var appState = decodeConfig(configuration);
  registerSaveConfigCallback(appState.getConfig);
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
  TextEditingController urlController;
  TextEditingController usernameController;
  TextEditingController passwordController;
  TextEditingController databaseController;
  TextEditingController batchSizeController;

  void urlChanged(String value) => config.connStr = value;
  void usernameChanged (String value) => config.username = value;
  void passwordChanged (String value) => config.password = value;
  void databaseChanged (String value) => config.database = value;
  void batchSizeChanged (String value) {
    var intValue = int.tryParse(value);
    if (intValue == null) {
      return;
    }
    config.batchSize = intValue;
  }

  void initState() {
    config = BlocProvider.of<Configuration>(context);
    urlController = TextEditingController(text: config.connStr);
    usernameController = TextEditingController(text: config.username);
    passwordController = TextEditingController(text: config.password);
    databaseController = TextEditingController(text: config.database);
    batchSizeController = TextEditingController(text: config.batchSize.toString());
    super.initState();
  }

  Widget build(BuildContext context) {
    return ListView(
      children: [
        TextField(controller: urlController, decoration: InputDecoration(labelText: "url"), onChanged: urlChanged),
        TextField(controller: usernameController, decoration: InputDecoration(labelText: "username"), onChanged: usernameChanged),
        TextField(controller: passwordController, decoration: InputDecoration(labelText: "password"), onChanged: passwordChanged),
        TextField(controller: databaseController, decoration: InputDecoration(labelText: "database"), onChanged: databaseChanged),
        TextField(controller: batchSizeController, decoration: InputDecoration(labelText: "batch  size"), onChanged: batchSizeChanged, inputFormatters: [FilteringTextInputFormatter.allow(RegExp(r'[0-9]'))]),
        ExportObjectSelector(),
      ],
    );
  }
}

class ExportObjectSelector extends StatefulWidget {
  ExportObjectSelector();

  createState() => _ExportObjectSelectorState();
}

class _ExportObjectSelectorState extends State<ExportObjectSelector> {
  Configuration config;

  void exportObjectChanged (String value) {
    config.exportObject = value;
    setState(() {});
  }

  initState() {
    config = BlocProvider.of<Configuration>(context);
    super.initState();
  }

  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        SizedBox(height: 20),
        Text("export object:", textScaleFactor: 0.9),
        DropdownButton<String>(
          hint: Text("export object"),
          items: [
            DropdownMenuItem(child: Text("Node"), value: "Node"),
            DropdownMenuItem(child: Text("Relationship"), value: "Relationship"),
          ],
          value: config.exportObject,
          onChanged: exportObjectChanged,
        ),
        NodeOrRelationshipConfig(config.exportObject),
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
    return Text("Invalid export object");
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
    return Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        TextField(controller: nodeLabelController, decoration: InputDecoration(labelText: "node label"), onChanged: nodeLabelChanged),
      ],
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

  void relLabelChanged(String value) {
    config.relLabel = value;
  }

  void relLeftLabelChanged(String value) {
    config.relLeftLabel = value;
  }

  void relRightLabelChanged(String value) {
    config.relRightLabel = value;
  }

  initState(){
    config = BlocProvider.of<Configuration>(context);
    relLabelController = TextEditingController(text: config.relLabel);
    relLeftLabelController = TextEditingController(text: config.relLeftLabel);
    relRightLabelController = TextEditingController(text: config.relRightLabel);
    super.initState();
  }
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        TextField(controller: relLabelController, decoration: InputDecoration(labelText: "relationship label"), onChanged: relLabelChanged),
        TextField(controller: relLeftLabelController, decoration: InputDecoration(labelText: "left node label"), onChanged: relLeftLabelChanged),
        TextField(controller: relRightLabelController, decoration: InputDecoration(labelText: "right node label"), onChanged: relRightLabelChanged),
      ],
    );
  }
}
