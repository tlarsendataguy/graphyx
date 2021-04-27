import 'package:flutter/material.dart';
import 'package:output/bloc.dart';
import 'package:output/configuration.dart';
import 'package:output/decode_config.dart';

void main() {
  var appState = decodeConfig(configuration);
  registerSaveConfigCallback(appState.getConfig);
  runApp(BlocProvider<Configuration>(
    child: MyApp(),
    bloc: appState,
  ));
}

class MyApp extends StatelessWidget {
  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Neo4j Output',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        accentColor: Colors.green,
      ),
      home: Scaffold(
        body: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Center(
            child: Controls(),
          ),
        )
      ),
    );
  }
}

class Controls extends StatefulWidget {
  Controls({Key key}) : super(key: key);

  _ControlsState createState() => _ControlsState();
}

class _ControlsState extends State<Controls> {
  Configuration config;
  TextEditingController urlController;
  TextEditingController usernameController;
  TextEditingController passwordController;

  void urlChanged(String value) => config.connStr = value;
  void usernameChanged (String value) => config.username = value;
  void passwordChanged (String value) => config.password = value;

  void initState() {
    config = BlocProvider.of<Configuration>(context);
    urlController = TextEditingController(text: config.connStr);
    usernameController = TextEditingController(text: config.username);
    passwordController = TextEditingController(text: config.password);
    super.initState();
  }

  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      mainAxisAlignment: MainAxisAlignment.center,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        TextField(controller: this.urlController, decoration: InputDecoration(labelText: "url"), onChanged: urlChanged),
        TextField(controller: this.usernameController, decoration: InputDecoration(labelText: "username"), onChanged: usernameChanged),
        TextField(controller: this.passwordController, decoration: InputDecoration(labelText: "password"), onChanged: passwordChanged),
      ],
    );
  }
}
