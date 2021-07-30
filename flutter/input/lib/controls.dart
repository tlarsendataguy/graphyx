import 'package:input/app_state.dart';
import 'package:input/bloc.dart';
import 'package:flutter/material.dart';
import 'package:input/connection_controls.dart';
import 'package:input/field_widget.dart';
import 'package:input/validated_response.dart';

class Controls extends StatelessWidget {
  Controls({Key key}) : super(key: key);

  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: <Widget>[
        ConnectionControls(),
        Expanded(
          child: QueryControls(),
        )
      ],
    );
  }
}

class QueryControls extends StatefulWidget {
  createState() => _QueryControlsState();
}

class _QueryControlsState extends State<QueryControls> {
  TextEditingController queryController;
  List<Widget> fieldWidgets = [];
  AppState state;
  bool isValidating = false;

  void initState() {
    state = BlocProvider.of<AppState>(context);
    queryController = TextEditingController(text: state.query);
    super.initState();
  }

  void queryChanged(value) {
    state.query = value;
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
    return Card(
      elevation: 12,
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
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
        ),
      ),
    );
  }
}