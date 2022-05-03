import 'package:flutter/material.dart';
import 'package:input/app_state.dart';
import 'package:input/bloc.dart';
import 'package:input/configuration.dart';
import 'package:input/controls.dart';

void main() async {
  while (true) {
    if (configurationLoaded) {
      break;
    }
    await Future.delayed(const Duration(milliseconds: 100));
  }
  var appState = AppState(configuration);
  registerSaveConfigCallback(appState.getConfig);
  registerDecryptCallback(appState.useDecryptedPassword);
  registerEncryptCallback(appState.useEncryptedPassword);

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
