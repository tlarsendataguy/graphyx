import 'dart:convert';

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

    if (dataType == '') {
      continue;
    }
    returnValues.add(ReturnValue(name: column, dataType: dataType));
  }

  return ValidatedResponse(error: "", returnValues: returnValues);
}

String getDataType(String typeFromJson) {
  switch (typeFromJson) {
    case 'String':
      return 'String';
    case 'node':
      return 'Node';
    case 'date':
      return 'Date';
    case 'int':
      return 'Integer';
    default:
      return '';
  }
}