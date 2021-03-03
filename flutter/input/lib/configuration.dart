@JS('customToolConfig')
library custom_tool_config;

import "package:js/js.dart";

external Config get Configuration;

@JS()
@anonymous
class Config {
  external factory Config({String ConnStr, String Username, String Password, String Query, List<FieldData> Fields});

  external String get ConnStr;
  external set ConnStr(String value);
  external String get Username;
  external set Username(String value);
  external String get Password;
  external set Password(String value);
  external String get Query;
  external set Query(String value);
  external List<FieldContainer> get Fields;
  external set Fields(List<FieldContainer> value);
  external ValidatedResponse get LastValidatedResponse;
  external set LastValidatedResponse(ValidatedResponse value);
}

@JS()
@anonymous
class FieldContainer {
  external factory FieldContainer({FieldData Field});

  external FieldData get Field;
}

@JS()
@anonymous
class FieldData {
  external factory FieldData({String Name, String DataType, List<ElementData> Path});

  external String get Name;
  external set Name(String value);
  external String get DataType;
  external set DataType(String value);
  external List<ElementContainer> get Path;
  external set Path(List<ElementContainer> value);
}

@JS()
@anonymous
class ElementContainer {
  external factory ElementContainer({ElementData Element});
  external ElementData get Element;
}

@JS()
@anonymous
class ElementData {
  external factory ElementData({String Key, String DataType});

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
      fieldsStr += '{Name="${field.Field.Name}, DataType="${field.Field.DataType}"}';
      needsDelimiter = true;
    }
    fieldsStr += ']';
  }

  return '{ConnStr="${config.ConnStr}", Username="${config.Username}", Query="${config.Query}", Fields=$fieldsStr}';
}

@JS()
@anonymous
class ValidatedResponse {
  external factory ValidatedResponse({List<ReturnValueContainer> ReturnValues, String Error});
  external List<ReturnValueContainer> get ReturnValues;
  external String get Error;
}

@JS()
@anonymous
class ReturnValueContainer {
  external factory ReturnValueContainer({ReturnValueData ReturnValue});
  external ReturnValueData get ReturnValue;
}

@JS()
@anonymous
class ReturnValueData {
  external factory ReturnValueData({String Name, String DataType});

  external String get Name;
  external String get DataType;
}
