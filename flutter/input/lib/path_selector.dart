import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:input/app_state.dart';
import 'package:input/bloc.dart';
import 'package:input/dropdown.dart';
import 'package:input/field_state.dart';
import 'package:input/neo4j_response.dart';

class PathSelector extends StatelessWidget {
  PathSelector(this.field);
  final Field field;

  Widget build(BuildContext context) {
    var fieldState = BlocProvider.of<FieldState>(context);
    return StreamBuilder<List<PathElement>>(
      stream: fieldState.pathChanged,
      builder: (_,AsyncSnapshot<List<PathElement>> value){
        if (!value.hasData || value.data.length == 0) {
          return ChooseReturnValue(field);
        }
        var lastContainer = value.data.last;
        switch (lastContainer.dataType){
          case 'Path':
            return SelectPathChild(field);
          case 'List:Node':
          case 'List:Path':
          case 'List:Relationship':
          case 'List:Integer':
          case 'List:Float':
          case 'List:DateTime':
          case 'List:Boolean':
          case 'List:String':
            var itemType = lastContainer.dataType.split(':')[1];
            return SelectListChild(field, itemType);
          case 'Node':
            return SelectNodeChild(field);
          case 'Relationship':
            return SelectRelationshipChild(field);
          case 'Map':
            return SelectMapChild(field);
          case 'Integer':
          case 'Float':
          case 'DateTime':
          case 'Boolean':
          case 'String':
            return SizedBox(height: 0);
          default:
            return Text("The current path ends in an invalid data type");
        }
      },
    );
  }
}

class ChooseReturnValue extends StatelessWidget {
  ChooseReturnValue(this.field);
  final Field field;

  Widget build(BuildContext context) {
    var appState = BlocProvider.of<AppState>(context);
    return StreamBuilder<ValidatedResponse>(
      stream: appState.lastValidatedResponse,
      builder: (_, AsyncSnapshot<ValidatedResponse> response){
        List<DropdownMenuItem<ReturnValue>> widgets;
        if (response.hasData && response.data.error == ''){
          widgets = response.data.returnValues.map<DropdownMenuItem<ReturnValue>>((e)=>DropdownMenuItem<ReturnValue>(child: Text('${e.name}:${e.dataType}'), value: e)).toList();
        } else {
          widgets = [];
        }
        return DropDown<ReturnValue>(items: widgets, onChanged: (e){
          var fieldState = BlocProvider.of<FieldState>(context);
          fieldState.addElementToPath(PathElement(key: e.name, dataType: e.dataType));
        });
      },
    );
  }
}

class SelectData {
  SelectData(this.name, this.dataType);
  final String name;
  final String dataType;
}

class SelectPathChild extends StatelessWidget {
  SelectPathChild(this.field);
  final Field field;

  Widget build(BuildContext context) {
    return DropDown<SelectData>(
      items: [
        DropdownMenuItem<SelectData>(child: Text("Nodes"), value: SelectData("Nodes", "List:Node")),
        DropdownMenuItem<SelectData>(child: Text("Relationships"), value: SelectData("Relationships", "List:Relationship")),
      ],
      onChanged: (e){
        var fieldState = BlocProvider.of<FieldState>(context);
        fieldState.addElementToPath(PathElement(key: e.name, dataType: e.dataType));
      },
    );
  }
}

class IndexDialog extends StatefulWidget {
  State<StatefulWidget> createState() => _IndexDialogState();
}

class _IndexDialogState extends State<IndexDialog> {
  TextEditingController _controller;
  FocusNode _node;

  initState(){
    _controller = TextEditingController(text: '');
    _node = FocusNode();
    _node.requestFocus();
    super.initState();
  }

  void cancel(){
    Navigator.of(context).pop(null);
  }

  void submit(){
    if (_controller.text == '') {
      return;
    }
    var index = int.parse(_controller.text);
    Navigator.of(context).pop(index);
  }

  Widget build(BuildContext context) {
    return Dialog(
      child: Container(
        padding: EdgeInsets.all(8.0),
        width: 200,
        height: 120,
        child: Column(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            TextField(
              decoration: InputDecoration(labelText: "index"),
              controller: _controller,
              focusNode: _node,
              onSubmitted: (_)=>submit(),
              inputFormatters: [
                FilteringTextInputFormatter.allow(RegExp(r'[0-9]')),
              ],
            ),
            Row(
              children: [
                Expanded(
                  child: TextButton(
                    onPressed: cancel,
                    child: Text("Cancel"),
                  ),
                ),
                Expanded(
                  child: ElevatedButton(
                    onPressed: submit,
                    child: Text("Submit"),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class SelectListChild extends StatelessWidget {
  SelectListChild(this.field, this.itemType);
  final Field field;
  final String itemType;

  Widget build(BuildContext context) {
    return DropDown<SelectData>(
      items: [
        DropdownMenuItem<SelectData>(child: Text("First"), value: SelectData("First", itemType)),
        DropdownMenuItem<SelectData>(child: Text("Last"), value: SelectData("Last", itemType)),
        DropdownMenuItem<SelectData>(child: Text("Index"), value: SelectData("Index", itemType)),
        DropdownMenuItem<SelectData>(child: Text("Count"), value: SelectData("Count", "Integer")),
      ],
      onChanged: (e) async {
        if (e.name == 'Index') {
          var index = await showDialog<int>(
            context: context,
            builder: (context) => IndexDialog(),
          );
          if (index == null){
            return;
          }
          var fieldState = BlocProvider.of<FieldState>(context);
          fieldState.addElementToPath(PathElement(key: "Index:$index", dataType: e.dataType));
          return;
        }
        var fieldState = BlocProvider.of<FieldState>(context);
        fieldState.addElementToPath(PathElement(key: e.name, dataType: e.dataType));
      },
    );
  }
}

class SelectNodeChild extends StatelessWidget {
  SelectNodeChild(this.field);
  final Field field;

  Widget build(BuildContext context) {
    return DropDown<SelectData>(
      items: [
        DropdownMenuItem<SelectData>(child: Text("ID"), value: SelectData("ID", 'Integer')),
        DropdownMenuItem<SelectData>(child: Text("Labels"), value: SelectData("Labels", 'List:String')),
        DropdownMenuItem<SelectData>(child: Text("Properties"), value: SelectData("Properties", 'Map')),
      ],
      onChanged: (e){
        var fieldState = BlocProvider.of<FieldState>(context);
        fieldState.addElementToPath(PathElement(key: e.name, dataType: e.dataType));
      },
    );
  }
}

class SelectRelationshipChild extends StatelessWidget {
  SelectRelationshipChild(this.field);
  final Field field;

  Widget build(BuildContext context) {
    return DropDown<SelectData>(
      items: [
        DropdownMenuItem<SelectData>(child: Text("ID"), value: SelectData("ID", 'Integer')),
        DropdownMenuItem<SelectData>(child: Text("StartId"), value: SelectData("StartId", 'Integer')),
        DropdownMenuItem<SelectData>(child: Text("EndId"), value: SelectData("EndId", 'Integer')),
        DropdownMenuItem<SelectData>(child: Text("Type"), value: SelectData("Type", 'String')),
        DropdownMenuItem<SelectData>(child: Text("Properties"), value: SelectData("Properties", 'Map')),
      ],
      onChanged: (e){
        var fieldState = BlocProvider.of<FieldState>(context);
        fieldState.addElementToPath(PathElement(key: e.name, dataType: e.dataType));
      },
    );
  }
}

class SelectMapChild extends StatefulWidget {
  SelectMapChild(this.field);
  final Field field;

  State<StatefulWidget> createState() => _SelectMapChildState();
}

class _SelectMapChildState extends State<SelectMapChild>{

  TextEditingController _property = TextEditingController(text: '');
  String _selectedType;

  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.end,
      children: [
        SizedBox(
          width: 200,
          child: TextField(
            decoration: InputDecoration(labelText: 'Property name'),
            controller: _property,
          ),
        ),
        Padding(
          padding: const EdgeInsets.fromLTRB(8, 0, 0, 0),
          child: DropDown<String>(
            hint: Text("Select property type"),
            value: _selectedType,
            items: [
              DropdownMenuItem<String>(child: Text("Boolean"), value: "Boolean"),
              DropdownMenuItem<String>(child: Text("DateTime"), value: "DateTime"),
              DropdownMenuItem<String>(child: Text("Float"), value: "Float"),
              DropdownMenuItem<String>(child: Text("Integer"), value: "Integer"),
              DropdownMenuItem<String>(child: Text("String"), value: "String"),
            ],
            onChanged: (e){
              setState(() =>_selectedType = e);
            },
          ),
        ),
        IconButton(
          icon: Icon(Icons.chevron_right),
          onPressed: (){
            var fieldState = BlocProvider.of<FieldState>(context);
            fieldState.addElementToPath(PathElement(key: _property.text, dataType: _selectedType));
          },
        ),
      ],
    );
  }
}