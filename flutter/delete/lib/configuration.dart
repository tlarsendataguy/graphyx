@JS()
library custom_tool_config;

import 'dart:html';

import "package:js/js.dart";
import 'package:js/js_util.dart';

typedef String GenerateConfig();

void registerSaveConfigCallback(GenerateConfig generator) {
  setProperty(window, 'getCustomToolConfig', allowInterop(generator));
}

@JS('customToolConfig')
external String get configuration;

@JS('customToolConfig')
external set configuration(String value);

@JS('customToolConfigLoaded')
external bool get configurationLoaded;

@JS('incomingFields')
external List<FieldInfo> get getIncomingFields;

@JS()
@anonymous
class FieldInfo {
  external String get strName;
  external String get strType;
}