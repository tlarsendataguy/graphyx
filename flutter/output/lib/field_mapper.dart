import 'package:flutter/material.dart';
import 'package:output/bloc.dart';
import 'package:output/decode_config.dart';

class FieldMapper extends StatefulWidget {
  FieldMapper({this.source, this.label="", this.height=200});
  final List<AyxToNeo4jMap> source;
  final double height;
  final String label;

  State<StatefulWidget> createState() => _FieldMapperState();
}

class _FieldMapperState extends State<FieldMapper> {

  TextEditingController _fieldNameController;
  List<PopupMenuItem<String>> incomingFields;
  FocusNode _fieldNameNode;

  List<PopupMenuItem<String>> getPopUp(BuildContext context) => incomingFields;

  void saveField(String field) {
    setState((){
      widget.source.add(AyxToNeo4jMap(field, field));
      _fieldNameController.clear();
      _fieldNameNode.requestFocus();
    });
  }

  Function generateFieldRemover(int index) {
    return () {
      setState((){
        widget.source.removeAt(index);
      });
    };
  }

  initState() {
    incomingFields = [];
    var config = BlocProvider.of<Configuration>(context);
    for (var field in config.incomingFields) {
      incomingFields.add(PopupMenuItem<String>(value: field, child: Text(field)));
    }
    _fieldNameController = TextEditingController(text:"");
    _fieldNameNode = FocusNode();
    super.initState();
  }

  Widget build(BuildContext context) {
    var config = BlocProvider.of<Configuration>(context);
    return SizedBox(
      height: widget.height,
      child: Column(
        children: [
          TextField(
            focusNode: _fieldNameNode,
            controller: _fieldNameController,
            decoration: InputDecoration(
              labelText: widget.label,
              suffixIcon: PopupMenuButton<String>(
                icon: Icon(Icons.arrow_drop_down),
                itemBuilder: getPopUp,
                onSelected: saveField,
              ),
            ),
            onSubmitted: saveField,
            autofillHints: config.incomingFields,
          ),
          Expanded(
            child: ListView.builder(
              itemBuilder: (context, index) {
                return Row(
                  children: [
                    Expanded(child: FieldMap(widget.source, index)),
                    IconButton(icon: Icon(Icons.delete), onPressed: generateFieldRemover(index)),
                  ],
                );
              },
              itemCount: widget.source.length,
            ),
          ),
        ],
      ),
    );
  }
}

class FieldMap extends StatefulWidget {
  FieldMap(this.source, this.index);
  final List<AyxToNeo4jMap> source;
  final int index;

  State<StatefulWidget> createState() => _FieldMapState();
}

class _FieldMapState extends State<FieldMap> {

  Configuration config;
  TextEditingController _neo4jNameController;

  AyxToNeo4jMap map() {
    return widget.source[widget.index];
  }

  void neo4jNameChanged(String value) {
    map().neo4jField = value;
  }

  initState(){
    config = BlocProvider.of<Configuration>(context);
    _neo4jNameController = TextEditingController(text: map().neo4jField);
    super.initState();
  }

  Widget build(BuildContext context) {
    return Row(
      children: [
        Expanded(
          child: Text(map().ayxField),
        ),
        Expanded(
          child: TextField(
            controller: _neo4jNameController,
            decoration: InputDecoration(labelText: "Neo4j property"),
            onChanged: neo4jNameChanged,
          ),
        )
      ],
    );
  }

}