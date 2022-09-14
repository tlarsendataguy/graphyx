import 'package:flutter/material.dart';
import 'package:delete/bloc.dart';
import 'package:delete/decode_config.dart';

class FieldSelector extends StatefulWidget {
  FieldSelector({this.source, this.label="", this.height=200});
  final List<String> source;
  final double height;
  final String label;

  State<StatefulWidget> createState() => _FieldSelectorState();
}

class _FieldSelectorState extends State<FieldSelector> {

  TextEditingController _fieldNameController;
  FocusNode _fieldNameNode;

  List<PopupMenuItem<String>> getPopUp(BuildContext context) {
    var config = BlocProvider.of<Configuration>(context);
    List<PopupMenuItem<String>> incomingFields = [];
    for (var field in config.incomingFields) {
      incomingFields.add(PopupMenuItem<String>(value: field, child: Text(field, overflow: TextOverflow.ellipsis)));
    }
    return incomingFields;
  }

  void saveField(String field) {
    setState((){
      widget.source.add(field);
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
                    Expanded(child: Text(widget.source[index], overflow: TextOverflow.ellipsis)),
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
