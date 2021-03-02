import 'package:flutter/material.dart';
import 'package:input/bloc.dart';
import 'package:input/configuration.dart' as c;
import 'package:input/field_state.dart';

class PathChips extends StatelessWidget {
  PathChips(this.field);
  final c.Field field;

  Widget build(BuildContext context) {
    var fieldState = BlocProvider.of<FieldState>(context);
    return StreamBuilder(
      stream: fieldState.pathChanged,
      builder: (_, __){
        List<Widget> widgets;
        if (field.Path == null) {
          widgets = [];
        } else {
          widgets = field.Path.map((e) => Chip(label: Text('${e.Key} (${e.DataType})'))).toList();
        }
        return Wrap(
          spacing: 4,
          runSpacing: 4,
          children: widgets,
        );
      },
    );
  }
}