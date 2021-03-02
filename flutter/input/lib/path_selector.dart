import 'package:flutter/material.dart';
import 'package:input/app_state.dart';
import 'package:input/bloc.dart';
import 'package:input/configuration.dart' as c;
import 'package:input/field_state.dart';

class PathSelector extends StatelessWidget {
  PathSelector(this.field);
  final c.Field field;

  Widget build(BuildContext context) {
    var fieldState = BlocProvider.of<FieldState>(context);
    return StreamBuilder(
      stream: fieldState.pathChanged,
      builder: (_,__){
        if (field.Path.length == 0) {
          return ChooseReturnValue(field);
        }
        var lastElement = field.Path.last;
        switch (lastElement.DataType){
          case 'Path':
            return SelectPathChild(field);
          default:
            return Text("hello world");
        }
      },
    );
  }
}

class ChooseReturnValue extends StatelessWidget {
  ChooseReturnValue(this.field);
  final c.Field field;

  Widget build(BuildContext context) {
    var appState = BlocProvider.of<AppState>(context);
    return StreamBuilder(
      stream: appState.returnValues,
      builder: (_, __){
        List<DropdownMenuItem<c.ReturnValue>> widgets;
        if (c.Configuration.LastValidatedResponse == null || c.Configuration.LastValidatedResponse.Error != ''){
          widgets = [];
        } else {
          widgets = c.Configuration.LastValidatedResponse.ReturnValues.map<DropdownMenuItem<c.ReturnValue>>((e)=>DropdownMenuItem<c.ReturnValue>(child: Text('${e.Name}:${e.DataType}'), value: e)).toList();
        }
        return DropdownButton<c.ReturnValue>(items: widgets, onChanged: (e){
          var fieldState = BlocProvider.of<FieldState>(context);
          fieldState.addElementToPath(c.Element(Key: e.Name, DataType: e.DataType));
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
  final c.Field field;

  Widget build(BuildContext context) {
    return DropdownButton<SelectData>(
      items: [
        DropdownMenuItem<SelectData>(child: Text("Nodes"), value: SelectData("Nodes", "List:Node")),
        DropdownMenuItem<SelectData>(child: Text("Relationships"), value: SelectData("Relationships", "List:Relationship")),
      ],
      onChanged: (e){
        var fieldState = BlocProvider.of<FieldState>(context);
        fieldState.addElementToPath(c.Element(Key: e.name, DataType: e.dataType));
      },
    );
  }
}