import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:input/configuration.dart' as c;

void main() {
  runApp(Neo4jInputApp());
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
          child: SingleChildScrollView(
            child: Center(
              child: Controls(),
            ),
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
  TextEditingController urlController = TextEditingController(text: "http://localhost:7474");
  TextEditingController userController = TextEditingController(text: "test");
  TextEditingController passwordController = TextEditingController(text: "test");
  TextEditingController queryController = TextEditingController(text: "MATCH p=()-[r:ACTED_IN]->() RETURN p");
  String response;

  Future _connect() async {
    var response = await http.get(urlController.text);
    setState((){
      this.response = response.body;
    });
  }

  Future _query() async {
    var query = this.queryController.text;
    if (!RegExp("\\sLIMIT\\s").hasMatch(query)) {
      query += " LIMIT 1";
    }

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
    setState((){
      this.response = response.body;
    });
  }

  void _getConfig() {
    setState((){
      this.response = c.configToString(c.Configuration);
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
      this.response = c.configToString(c.Configuration);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: <Widget>[
        TextField(controller: this.urlController, decoration: InputDecoration(labelText: "url")),
        TextField(controller: this.userController, decoration: InputDecoration(labelText: "username")),
        TextField(controller: this.passwordController, decoration: InputDecoration(labelText: "password")),
        TextField(controller: this.queryController, decoration: InputDecoration(labelText: "query")),
        TextButton(onPressed: _connect, child: Text("Test Connection")),
        TextButton(onPressed: _query, child: Text("Run Query")),
        TextButton(onPressed: _getConfig, child: Text("Get Config")),
        TextButton(onPressed: _setConfig, child: Text("Set Config")),
        SelectableText(
          '$response',
        ),
      ],
    );
  }
}
