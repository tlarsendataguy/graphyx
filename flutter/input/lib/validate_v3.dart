import 'dart:convert';

import 'package:input/return_value_types.dart';
import 'package:input/validated_response.dart';

ValidatedResponse validateV3Response(String response) {
  var decoded = jsonDecode(response);

  List<ReturnValue> returnValues = [];
  var firstResult = decoded['results'][0];
  var firstRow = firstResult['data'][0];
  var columnIndex = 0;
  for (var column in firstResult['columns']) {
    var meta = firstRow['meta'][columnIndex];
    var data = firstRow['row'][columnIndex];
    columnIndex++;
    var dataType = getDataType(meta == null ? data.runtimeType.toString() : meta['type'].toString());

    if (dataType == rUnknown) {
      continue;
    }
    returnValues.add(ReturnValue(name: column, dataType: dataType));
  }

  return ValidatedResponse(error: "", returnValues: returnValues);
}

String getDataType(String typeFromJson) {
  switch (typeFromJson) {
    case 'String':
      return rString;
    case 'node':
      return rNode;
    case 'date':
      return rDate;
    case 'int':
      return rInteger;
    case '_InternalLinkedHashMap<String, dynamic>':
      return rMap;
    default:
      print(typeFromJson);
      return rUnknown;
  }
}