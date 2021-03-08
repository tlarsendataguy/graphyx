import 'package:flutter/material.dart';
import 'package:input/app_state.dart';
import 'package:input/bloc.dart';
import 'package:input/field_state.dart';
import 'package:input/images.dart';
import 'package:input/path_chips.dart';
import 'package:input/path_selector.dart';

class FieldWidget extends StatefulWidget {
  FieldWidget(this.index, this.field) {
    this.key = ObjectKey(field);
  }
  int index;
  Field field;
  Key key;
  State<StatefulWidget> createState() {
    return _FieldWidgetState();
  }
}

class _FieldWidgetState extends State<FieldWidget> {
  TextEditingController _name;
  AppState _appState;

  initState(){
    _appState = BlocProvider.of<AppState>(context);
    _name = TextEditingController(text: widget.field.name);
  }

  void _deleteField(){
    _appState.removeField(widget.index);
  }

  Widget build(BuildContext context) {
    return BlocProvider<FieldState>(
      bloc: FieldState(widget.field),
      child: Padding(
        padding: const EdgeInsets.fromLTRB(0, 4, 0, 4),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                IconButton(icon: Image.asset(delete), onPressed: _deleteField),
                SizedBox(
                  width: 200,
                  child: TextField(controller: _name, onChanged: (value){
                    widget.field.name = value;
                  }),
                ),
                Expanded(child: PathChips(widget.field, widget.index)),
                ReorderableDragStartListener(
                  index: widget.index,
                  child: Image.asset(dragHandle),
                )
              ],
            ),
            PathSelector(widget.field),
          ],
        ),
      ),
    );
  }
}
