
import 'dart:convert';

import 'package:input/app_state.dart';
import 'package:input/validated_response.dart';

Configuration decodeConfig(String configStr) {
  if (configStr == null) {
    return Configuration(
      connStr: '',
      username: '',
      password: '',
      query: '',
      database: '',
      urlCollapsed: false,
      lastValidatedResponse: ValidatedResponse(error: '', returnValues: []),
      fields: [],
    );
  }
  var decoded = json.decode(configStr);
  return Configuration(
    connStr: decoded['ConnStr'] ?? '',
    username: decoded['Username'] ?? '',
    password: decoded['Password'] ?? '',
    database: decoded['Database'] ?? '',
    urlCollapsed: decoded['UrlCollapsed'] ?? false,
    query: decoded['Query'] ?? '',
    lastValidatedResponse: _decodeValidatedResponse(decoded['LastValidatedResponse']),
    fields: _decodeFields(decoded['Fields']),
  );
}

ValidatedResponse _decodeValidatedResponse(dynamic responseObj) {
  if (responseObj == null) {
    return ValidatedResponse(error: '', returnValues: []);
  }
  return ValidatedResponse(
    error: responseObj['Error'] ?? '',
    returnValues: _decodeReturnValues(responseObj['ReturnValues']),
  );
}

List<ReturnValue> _decodeReturnValues(dynamic returnValuesArray) {
  if (returnValuesArray == null) {
    return [];
  }
  List<ReturnValue> returnValues = [];
  for (var returnValue in returnValuesArray) {
    if (returnValue == null) {
      continue;
    }
    returnValues.add(ReturnValue(
      name: returnValue['Name'] ?? '',
      dataType: returnValue['DataType'] ?? '',
    ));
  }
  return returnValues;
}

List<Field> _decodeFields(dynamic fieldsArray) {
  if (fieldsArray == null) {
    return [];
  }
  List<Field> fields = [];
  for (var field in fieldsArray) {
    if (field == null) {
      continue;
    }
    fields.add(Field(
      name: field['Name'] ?? '',
      dataType: field['DataType'] ?? '',
      path: _decodePath(field['Path']),
    ));
  }
  return fields;
}

List<PathElement> _decodePath(dynamic elementArray) {
  if (elementArray == null) {
    return [];
  }
  List<PathElement> path = [];
  for (var element in elementArray) {
    if (element == null) {
      continue;
    }
    path.add(PathElement(
      key: element['Key'] ?? '',
      dataType: element['DataType'] ?? '',
    ));
  }
  return path;
}