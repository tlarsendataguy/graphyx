import 'dart:typed_data';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:input/app_state.dart';
import 'package:input/bloc.dart';
import 'package:input/configuration.dart';
import 'package:input/controls.dart';
import 'package:input/material_icons.dart';
import 'package:input/mono_font_file.dart';

Future<ByteData> fontFileToByteData(List<int> file) async {
  return ByteData.sublistView(Uint8List.fromList(file));
}

void main() async {
  while (true) {
    if (configurationLoaded) {
      break;
    }
    await Future.delayed(const Duration(milliseconds: 100));
  }
  var appState = AppState(configuration);
  registerSaveConfigCallback(appState.getConfig);

  var monoLoader = FontLoader("JetBrains Mono");
  monoLoader.addFont(fontFileToByteData(monoFontFile));
  var monoFuture = monoLoader.load();

  var materialLoader = FontLoader("MaterialIcons");
  materialLoader.addFont(fontFileToByteData(materialIcons));
  var materialFuture = materialLoader.load();

  await monoFuture;
  await materialFuture;

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
