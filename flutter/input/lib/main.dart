import 'package:flutter/material.dart';
import 'package:input/app_state.dart';
import 'package:input/bloc.dart';
import 'package:input/controls.dart';

void main() {
  var appState = AppState();
  runApp(BlocProvider<AppState>(
    child: Neo4jInputApp(),
    bloc: appState,
  ));
}

class Neo4jInputApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Neo4j Input',
      theme: ThemeData(
        primarySwatch: Colors.indigo,
        accentColor: Colors.green,
      ),
      home: Scaffold(
        body: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Center(
            child: Controls(),
          ),
        ),
      ),
    );
  }
}
