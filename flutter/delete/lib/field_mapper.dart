import 'package:flutter/material.dart';
import 'package:delete/bloc.dart';
import 'package:delete/decode_config.dart';

class FieldMapper extends StatefulWidget {
  FieldMapper({this.source, this.label="", this.height=200});
  final List<AyxToNeo4jMap> source;
  final double height;
  final String label;

  State<StatefulWidget> createState() => _FieldMapperState();
}

class _FieldMapperState extends State<FieldMapper> {

  TextEditingController _fieldNameController;
  FocusNode _fieldNameNode;

  List<PopupMenuItem<String>> getPopUp(BuildContext context) {
    List<PopupMenuItem<String>> incomingFields = [];
    var config = BlocProvider.of<Configuration>(context);
    for (var field in config.incomingFields) {
      incomingFields.add(PopupMenuItem<String>(value: field, child: Text(field, overflow: TextOverflow.ellipsis)));
    }
    return incomingFields;
  }

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
    _fieldNameController = TextEditingController(text:"");
    _fieldNameNode = FocusNode();
    super.initState();
  }

  Widget build(BuildContext context) {
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
          child: Text(map().ayxField, overflow: TextOverflow.ellipsis),
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