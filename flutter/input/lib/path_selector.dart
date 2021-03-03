import 'package:flutter/material.dart';
import 'package:input/app_state.dart';
import 'package:input/bloc.dart';
import 'package:input/configuration.dart' as c;
import 'package:input/field_state.dart';

class PathSelector extends StatelessWidget {
  PathSelector(this.field);
  final c.FieldData field;

  Widget build(BuildContext context) {
    var fieldState = BlocProvider.of<FieldState>(context);
    return StreamBuilder(
      stream: fieldState.pathChanged,
      builder: (_,__){
        if (field.Path.length == 0) {
          return ChooseReturnValue(field);
        }
        var lastContainer = field.Path.last;
        switch (lastContainer.Element.DataType){
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
            var itemType = lastContainer.Element.DataType.split(':')[1];
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
  final c.FieldData field;

  Widget build(BuildContext context) {
    var appState = BlocProvider.of<AppState>(context);
    return StreamBuilder(
      stream: appState.returnValues,
      builder: (_, __){
        List<DropdownMenuItem<c.ReturnValueData>> widgets;
        if (c.Configuration.LastValidatedResponse == null || c.Configuration.LastValidatedResponse.Error != ''){
          widgets = [];
        } else {
          widgets = c.Configuration.LastValidatedResponse.ReturnValues.map<DropdownMenuItem<c.ReturnValueData>>((e)=>DropdownMenuItem<c.ReturnValueData>(child: Text('${e.ReturnValue.Name}:${e.ReturnValue.DataType}'), value: e.ReturnValue)).toList();
        }
        return DropdownButton<c.ReturnValueData>(items: widgets, onChanged: (e){
          var fieldState = BlocProvider.of<FieldState>(context);
          fieldState.addElementToPath(c.ElementData(Key: e.Name, DataType: e.DataType));
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
  final c.FieldData field;

  Widget build(BuildContext context) {
    return DropdownButton<SelectData>(
      items: [
        DropdownMenuItem<SelectData>(child: Text("Nodes"), value: SelectData("Nodes", "List:Node")),
        DropdownMenuItem<SelectData>(child: Text("Relationships"), value: SelectData("Relationships", "List:Relationship")),
      ],
      onChanged: (e){
        var fieldState = BlocProvider.of<FieldState>(context);
        fieldState.addElementToPath(c.ElementData(Key: e.name, DataType: e.dataType));
      },
    );
  }
}

class SelectListChild extends StatelessWidget {
  SelectListChild(this.field, this.itemType);
  final c.FieldData field;
  final String itemType;

  Widget build(BuildContext context) {
    return DropdownButton<SelectData>(
      items: [
        DropdownMenuItem<SelectData>(child: Text("First"), value: SelectData("First", itemType)),
        DropdownMenuItem<SelectData>(child: Text("Last"), value: SelectData("Last", itemType)),
      ],
      onChanged: (e){
        var fieldState = BlocProvider.of<FieldState>(context);
        fieldState.addElementToPath(c.ElementData(Key: e.name, DataType: e.dataType));
      },
    );
  }
}

class SelectNodeChild extends StatelessWidget {
  SelectNodeChild(this.field);
  final c.FieldData field;

  Widget build(BuildContext context) {
    return DropdownButton<SelectData>(
      items: [
        DropdownMenuItem<SelectData>(child: Text("ID"), value: SelectData("ID", 'Integer')),
        DropdownMenuItem<SelectData>(child: Text("Labels"), value: SelectData("Labels", 'List:String')),
        DropdownMenuItem<SelectData>(child: Text("Properties"), value: SelectData("Properties", 'Map')),
      ],
      onChanged: (e){
        var fieldState = BlocProvider.of<FieldState>(context);
        fieldState.addElementToPath(c.ElementData(Key: e.name, DataType: e.dataType));
      },
    );
  }
}

class SelectRelationshipChild extends StatelessWidget {
  SelectRelationshipChild(this.field);
  final c.FieldData field;

  Widget build(BuildContext context) {
    return DropdownButton<SelectData>(
      items: [
        DropdownMenuItem<SelectData>(child: Text("ID"), value: SelectData("ID", 'Integer')),
        DropdownMenuItem<SelectData>(child: Text("StartId"), value: SelectData("StartId", 'Integer')),
        DropdownMenuItem<SelectData>(child: Text("EndId"), value: SelectData("EndId", 'Integer')),
        DropdownMenuItem<SelectData>(child: Text("Type"), value: SelectData("Type", 'String')),
        DropdownMenuItem<SelectData>(child: Text("Properties"), value: SelectData("Properties", 'Map')),
      ],
      onChanged: (e){
        var fieldState = BlocProvider.of<FieldState>(context);
        fieldState.addElementToPath(c.ElementData(Key: e.name, DataType: e.dataType));
      },
    );
  }
}

class SelectMapChild extends StatefulWidget {
  SelectMapChild(this.field);
  final c.FieldData field;

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
          child: DropdownButton<String>(
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
            fieldState.addElementToPath(c.ElementData(Key: _property.text, DataType: _selectedType));
          },
        ),
      ],
    );
  }
}