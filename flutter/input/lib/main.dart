import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:input/app_state.dart';
import 'package:input/bloc.dart';
import 'package:input/configuration.dart' as c;

void main() {
  var appState = AppState();
  runApp(BlocProvider<AppState>(
    child: Neo4jInputApp(),
    bloc: appState,
  ));
}

class Neo4jInputApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Neo4j Input',
      theme: ThemeData(
        primarySwatch: Colors.indigo,
        accentColor: Colors.green,
      ),
      home: Scaffold(
        body: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Center(
            child: Controls(),
          ),
        ),
      ),
    );
  }
}

class Controls extends StatefulWidget {
  Controls({Key key}) : super(key: key);

  @override
  _ControlsState createState() => _ControlsState();
}

class _ControlsState extends State<Controls> {
  String connStr;
  TextEditingController urlController;
  TextEditingController userController;
  TextEditingController passwordController;
  TextEditingController queryController;
  String validationError = '';

  void initState(){
    urlController = TextEditingController(text: c.Configuration.ConnStr);
    userController = TextEditingController(text: c.Configuration.Username);
    passwordController = TextEditingController(text: c.Configuration.Password);
    queryController = TextEditingController(text: c.Configuration.Query);
    super.initState();
  }

  void urlChanged(value) {
    c.Configuration.ConnStr = value;
  }

  void usernameChanged(value) {
    c.Configuration.Username = value;
  }

  void passwordChanged(value) {
    c.Configuration.Password = value;
  }

  void queryChanged(value) {
    c.Configuration.Query = value;
  }

  Future _validateQuery() async {
    var query = this.queryController.text;
    if (!RegExp("\\sLIMIT\\s").hasMatch(query)) {
      query += " LIMIT 1";
    }

    try {
      var response = await http.post(
          '${this.urlController.text}/db/neo4j/tx/commit',
          headers: {
            'Accept': 'application/vnd.neo4j.jolt+json-seq;strict=true',
            'Content-Type': 'application/json',
            'Authorization': 'Basic ' + base64Encode(utf8.encode('${this.userController.text}:${this.passwordController.text}')),
          },
          body: '''{
  "statements": [
    {
      "statement": "$query",
      "parameters": {}
    }
  ]
}'''

      );

      try {
        var errors = jsonDecode(response.body);
        var msg = errors['errors'][0]['message'];
        setState((){
          this.validationError = msg;
        });
        return;
      } catch (_) {}

      var events = response.body.split('\n');
      print(events);
      var header = jsonDecode(events[0]);
      var data = jsonDecode(events[1]);
      if (data['data'] == null) {
        setState((){
          this.validationError = 'The query was successful but no records were returned.  No metadata is available to generate output fields.';
        });
        return;
      }

      var fields = header['header']['fields'];
      var dataTypes = data['data'];
      List<String> metaData = [];
      var index = 0;
      for (var field in fields) {
        var dataType = List.from(dataTypes[index].keys)[0];
        var fieldType = 'Unknown';
        switch (dataType) {
          case '..':
            fieldType = 'Path';
            break;
          case '()':
            fieldType = 'Node';
            break;
          case '->':
          case '<-':
            fieldType = 'Relationship';
            break;
          case '[]':
            var firstItem = dataTypes[index][dataType][0];
            var firstItemType = List.from(firstItem.keys)[0];
            fieldType = 'List:$firstItemType';
            break;
          case '{}':
            fieldType = 'Map';
            break;
          case '?':
            fieldType = 'Boolean';
            break;
          case 'Z':
            fieldType = 'Integer';
            break;
          case 'R':
            fieldType = 'Float';
            break;
          case 'U':
            fieldType = 'String';
            break;
          case 'T':
            fieldType = 'Date';
            break;
          default:
            break;
        }
        metaData.add('${field.toString()}:$fieldType');
        index++;
      }
      print('$metaData');

      setState((){
        this.validationError = '';
      });
    }
    catch (ex) {
      setState((){
        //this.validationError = 'Unable to connect to the Neo4j database.  Double-check the URL and credentials and make sure you have a working network connection to the database.';
        this.validationError = ex.toString();
      });
    }
  }

  void _getConfig() {
    setState((){
      this.validationError = c.configToString(c.Configuration);
    });
  }

  void _setConfig() {
    setState((){
      var newField = c.Field(Name: 'Hello World', DataType: 'Integer', Path: [
      c.Element(Key: 'p', DataType: 'Integer'),
      ]);
      if (c.Configuration.Fields == null) {
        c.Configuration.Fields = [newField];
      } else {
        c.Configuration.Fields.add(c.Field(Name: 'Hello World', DataType: 'Integer', Path: [
          c.Element(Key: 'p', DataType: 'Integer'),
        ]));
      }
      this.validationError = c.configToString(c.Configuration);
    });
  }

  @override
  Widget build(BuildContext context) {

    return Column(
      mainAxisSize: MainAxisSize.min,
      mainAxisAlignment: MainAxisAlignment.center,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: <Widget>[
        TextField(controller: this.urlController, decoration: InputDecoration(labelText: "url"), onChanged: urlChanged),
        TextField(controller: this.userController, decoration: InputDecoration(labelText: "username"), onChanged: usernameChanged),
        TextField(controller: this.passwordController, decoration: InputDecoration(labelText: "password"), onChanged: passwordChanged),
        TextField(controller: this.queryController, decoration: InputDecoration(labelText: "query"), onChanged: queryChanged, style: TextStyle(fontFamily: 'JetBrains')),
        Padding(
          padding: const EdgeInsets.fromLTRB(0, 8, 0, 8),
          child: TextButton(onPressed: _validateQuery, child: Text("Validate query")),
        ),
        validationError == '' ? SizedBox(height: 0) : SelectableText(
          '$validationError',
          style: TextStyle(color: Colors.red),
        ),
        Expanded(
          child: ReorderableListView(
            onReorder: (value1, value2){
              print("value1=$value1, value2=$value2");
            },
            children: [
              Text('1', key: Key('1')),
              Text('2', key: Key('2')),
              Text('3', key: Key('3')),
            ],
          ),
        ),
        ElevatedButton(onPressed: (){}, child: Text("Add field"),)
      ],
    );
  }
}
