import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:io';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: MyHomePage(title: 'Flutter Demo Home Page'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  MyHomePage({Key key, this.title}) : super(key: key);
  final String title;

  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
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
      query += " LIMIT 25";
    }

    var response = await http.post(
      '${this.urlController.text}/db/neo4j/tx/commit',
      headers: {
        'Accept': 'application/vnd.neo4j.jolt+json-seq;strict=true',
        'Content-Type': 'application/json',
        'X-Stream': 'true',
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

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: SingleChildScrollView(
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: <Widget>[
              TextField(controller: this.urlController, decoration: InputDecoration(hintText: "url")),
              TextField(controller: this.userController, decoration: InputDecoration(hintText: "username")),
              TextField(controller: this.passwordController, decoration: InputDecoration(hintText: "password")),
              TextField(controller: this.queryController, decoration: InputDecoration(hintText: "query")),
              TextButton(onPressed: _connect, child: Text("Test Connection")),
              TextButton(onPressed: _query, child: Text("Run Query")),
              SelectableText(
                '$response',
              ),
            ],
          ),
        ),
      ),
    );
  }
}
