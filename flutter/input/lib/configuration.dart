@JS('customToolConfig')
library custom_tool_config;

import "package:js/js.dart";

external Config get Configuration;

@JS()
@anonymous
class Config {
  external factory Config({String ConnStr, String Username, String Password, String Query, List<Field> Fields});

  external String get ConnStr;
  external set ConnStr(String value);
  external String get Username;
  external set Username(String value);
  external String get Password;
  external set Password(String value);
  external String get Query;
  external set Query(String value);
  external List<Field> get Fields;
  external set Fields(List<Field> value);
  external ValidatedResponse get LastValidatedResponse;
  external set LastValidatedResponse(ValidatedResponse value);
}

@JS()
@anonymous
class Field {
  external factory Field({String Name, String DataType, List<Element> Path});

  external String get Name;
  external set Name(String value);
  external String get DataType;
  external set DataType(String value);
  external List<Element> get Path;
}

@JS()
@anonymous
class Element {
  external factory Element({String Key, String DataType});

  external String get Key;
  external String get DataType;
}

String configToString(Config config) {
  var fieldsStr = 'null';
  if (config.Fields != null) {
    fieldsStr = '[';
    var needsDelimiter = false;
    for (var field in config.Fields) {
      if (needsDelimiter) {
        fieldsStr += ',';
      }
      fieldsStr += '{Name="${field.Name}, DataType="${field.DataType}"}';
      needsDelimiter = true;
    }
    fieldsStr += ']';
  }

  return '{ConnStr="${config.ConnStr}", Username="${config.Username}", Query="${config.Query}", Fields=$fieldsStr}';
}

@JS()
@anonymous
class ValidatedResponse {
  external factory ValidatedResponse({ReturnValues, Error});
  external List<ReturnValue> get ReturnValues;
  external String get Error;
}

@JS()
@anonymous
class ReturnValue {
  external factory ReturnValue({String Name, String DataType});

  external String get Name;
  external String get DataType;
}
