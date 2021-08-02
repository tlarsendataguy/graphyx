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
    if (meta == null) {
      switch (data.runtimeType.toString()) {
        case 'String':
          returnValues.add(ReturnValue(name: column, dataType: 'String'));
      }
      continue;
    }

    switch (meta['type'].toString()) {
      case 'node':
        returnValues.add(ReturnValue(name: column, dataType: 'Node'));
    }
  }

  return ValidatedResponse(error: "", returnValues: returnValues);
}
