import 'package:input/app_state.dart';
import 'package:input/bloc.dart';
import 'package:flutter/material.dart';
import 'package:input/field_widget.dart';

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
  String validationError = '';
  List<Widget> fieldWidgets = [];
  AppState state;

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
        TextField(controller: this.databaseController, decoration: InputDecoration(labelText: "database"), onChanged: databaseChanged),
        TextField(controller: this.queryController, decoration: InputDecoration(labelText: "query"), onChanged: queryChanged, style: TextStyle(fontFamily: 'JetBrains Mono')),
        Padding(
          padding: const EdgeInsets.fromLTRB(0, 8, 0, 8),
          child: TextButton(onPressed: state.validateQuery, child: Text("Validate query")),
        ),
        validationError == '' ? SizedBox(height: 0) : SelectableText(
          '$validationError',
          style: TextStyle(color: Colors.red),
        ),
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
        ElevatedButton(onPressed: state.addField, child: Text("Add field")),
      ],
    );
  }
}
