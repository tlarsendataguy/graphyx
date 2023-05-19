@JS()
library custom_tool_config;

import 'dart:html';

import "package:js/js.dart";
import 'package:js/js_util.dart';

typedef String GenerateConfig();

void registerSaveConfigCallback(GenerateConfig generator) {
  setProperty(window, 'getCustomToolConfig', allowInterop(generator));
}

typedef void PasswordFunc(String value);

void registerEncryptCallback(PasswordFunc f) {
  setProperty(window, 'encryptPasswordCallback', allowInterop(f));
}

void registerDecryptCallback(PasswordFunc f) {
  setProperty(window, 'decryptPasswordCallback', allowInterop(f));
}

@JS('customToolConfig')
external String get configuration;

@JS('customToolConfig')
external set configuration(String value);

@JS('customToolConfigLoaded')
external bool get configurationLoaded;

@JS('Alteryx.JsEvent')
external void JsEvent(String eventStr);

@JS()
class FieldMetadata {
  external String get Name;
  external String get DataType;
}

@JS()
class ReturnValueMetadata {
  external String get Error;
  external List<FieldMetadata> get Fields;
}

@JS('returnValueMetadata')
external ReturnValueMetadata returnValueMetadata(String uri, String user, String password, String database, String query);


