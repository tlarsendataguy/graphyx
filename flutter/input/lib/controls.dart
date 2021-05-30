import 'dart:math';

import 'package:input/app_state.dart';
import 'package:input/bloc.dart';
import 'package:flutter/material.dart';
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
  TextEditingController databaseController;
  List<Widget> fieldWidgets = [];
  AppState state;
  bool isValidating = false;

  void initState(){
    state = BlocProvider.of<AppState>(context);
    urlController = TextEditingController(text: state.connStr);
    userController = TextEditingController(text: state.username);
    passwordController = TextEditingController(text: state.password);
    databaseController = TextEditingController(text: state.database);
    queryController = TextEditingController(text: state.query);
    super.initState();
  }

  void urlChanged(value) {
    state.connStr = value;
  }

  void usernameChanged(value) {
    state.username = value;
  }

  void passwordChanged(value) {
    state.password = value;
  }

  void queryChanged(value) {
    state.query = value;
  }

  void databaseChanged(value) {
    state.database = value;
  }

  void generateFieldWidgets(List<Field> fields){
    if (fields == null) {
      fieldWidgets = [];
      return;
    }
    List<Widget> children = [];
    var index = 0;
    for (var field in fields) {
      children.add(FieldWidget(index, field));
      index++;
    }
    fieldWidgets = children;
  }

  Future validateQuery() async {
    setState(()=>isValidating=true);
    await state.validateQuery();
    setState(()=>isValidating=false);
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: <Widget>[
        TextField(controller: this.urlController, decoration: InputDecoration(labelText: "url"), onChanged: urlChanged, autocorrect: false),
        TextField(controller: this.userController, decoration: InputDecoration(labelText: "username"), onChanged: usernameChanged, autocorrect: false),
        TextField(controller: this.passwordController, decoration: InputDecoration(labelText: "password"), onChanged: passwordChanged, autocorrect: false),
        TextField(controller: this.databaseController, decoration: InputDecoration(labelText: "database"), onChanged: databaseChanged, autocorrect: false),
        TextField(controller: this.queryController, decoration: InputDecoration(labelText: "query"), onChanged: queryChanged, style: TextStyle(fontFamily: 'JetBrains Mono'), minLines: 1, maxLines: 10, autocorrect: false),
        Padding(
          padding: const EdgeInsets.fromLTRB(0, 8, 0, 8),
          child: SizedBox(
            height: 40,
            child: TextButton(
              onPressed: validateQuery,
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text("Validate query"),
                  isValidating ? CircularProgressIndicator(strokeWidth: 2) : SizedBox(width: 0),
                ],
              ),
            ),
          ),
        ),
        StreamBuilder<ValidatedResponse>(
          stream: state.lastValidatedResponse,
          builder: (_, AsyncSnapshot<ValidatedResponse> value){
            if (value.hasData && value.data.error != '') {
              return SelectableText(
                '${value.data.error}',
                style: TextStyle(color: Colors.red)
              );
            }
            return SizedBox(height: 0);
          },
        ),
        ElevatedButton(onPressed: state.addField, child: Text("Add field")),
        StreamBuilder<List<Field>>(
          stream: BlocProvider.of<AppState>(context).fields,
          builder: (_, AsyncSnapshot<List<Field>> value) {
            generateFieldWidgets(value.data);
            return Expanded(
              child: ReorderableListView(
                onReorder: state.moveField,
                children: fieldWidgets,
                buildDefaultDragHandles: false,
              ),
            );
          },
        ),
      ],
    );
  }
}
