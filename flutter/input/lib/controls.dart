import 'dart:convert';
import 'package:input/app_state.dart';
import 'package:input/bloc.dart';
import 'package:input/configuration.dart' as c;
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:input/field_widget.dart';
import 'package:input/neo4j_response.dart';

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
  List<Widget> fieldWidgets = [];

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
      var body = {
        "statements": [
          {
            "statement": query,
            "parameters": {},
          },
        ],
      };
      var response = await http.post(
          '${this.urlController.text}/db/neo4j/tx/commit',
          headers: {
            'Accept': 'application/vnd.neo4j.jolt+json-seq;strict=true',
            'Content-Type': 'application/json',
            'Authorization': 'Basic ' + base64Encode(utf8.encode('${this.userController.text}:${this.passwordController.text}')),
          },
          body: jsonEncode(body),
      );
      var validated = validate(response.body);
      c.Configuration.LastValidatedResponse = validated;
      if (validated.Error != ''){
        setState((){
          this.validationError = validated.Error;
        });
        return;
      }

      setState((){
        this.validationError = '';
      });
    }
    catch (ex) {
      setState((){
        this.validationError = 'Unable to connect to the Neo4j database.  Double-check the URL make sure you have a working network connection to the database.';
      });
    }
  }

  void moveField(int moveFrom, int moveTo) {
    var field = c.Configuration.Fields[moveFrom];
    if (moveTo > moveFrom) {
      c.Configuration.Fields.insert(moveTo, field);
      c.Configuration.Fields.removeAt(moveFrom);
    } else {
      c.Configuration.Fields.removeAt(moveFrom);
      c.Configuration.Fields.insert(moveTo, field);
    }
    setState(generateFieldWidgets);
  }

  void generateFieldWidgets(){
    var fields = c.Configuration.Fields;
    List<Widget> children = [];
    if (fields != null){
      var indexes = List<int>.generate(c.Configuration.Fields.length, (e)=>e);
      children = indexes.map((e) => FieldWidget(e)).toList();
    }
    fieldWidgets = children;
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
        StreamBuilder(
          stream: BlocProvider.of<AppState>(context).fields,
          builder: (_, __) {
            generateFieldWidgets();
            return Expanded(
              child: ReorderableListView(
                onReorder: moveField,
                children: fieldWidgets,
              ),
            );
          },
        ),
        ElevatedButton(onPressed: (){
          var fields = c.Configuration.Fields;
          if (fields == null){
            c.Configuration.Fields = [c.FieldContainer(Field: c.FieldData(Path: []))];
          } else {
            c.Configuration.Fields.add(c.FieldContainer(Field: c.FieldData(Path: [])));
          }
          setState(generateFieldWidgets);
        }, child: Text("Add field"),)
      ],
    );
  }
}
