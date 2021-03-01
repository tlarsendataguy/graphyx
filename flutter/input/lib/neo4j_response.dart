import 'dart:convert';
import 'package:input/configuration.dart';

ValidatedResponse validate(String response) {
  try {
    var decoded = jsonDecode(response);
    if (decoded['errors'] != null) {
      return ValidatedResponse(Error: decoded['errors'][0]['message']);
    }
  } catch (_) {}

  var lines = response.split('\n');
  if (lines.length < 2) {
    return ValidatedResponse(Error: 'A response with an unexpected format was returned.  Response was:\n$response');
  }

  var header = jsonDecode(lines[0]);

  if (header['error'] != null) {
    return ValidatedResponse(Error: header['error']['errors'][0]['message']['U']);
  }

  var data = jsonDecode(lines[1]);
  if (data['data'] == null) {
    return ValidatedResponse(Error: 'The query was successful but no records were returned.  No metadata is available to generate output fields.');
  }

  var fields = header['header']['fields'];
  var dataTypes = data['data'];
  List<ReturnValue> returnValues = [];
  var index = 0;
  for (var field in fields) {
    var dataType = List.from(dataTypes[index].keys)[0];
    String fieldType;
    switch (dataType) {
      case '[]':
        var firstItem = dataTypes[index][dataType][0];
        var firstItemType = List.from(firstItem.keys)[0];
        fieldType = 'List:${decodeNonListDataType(firstItemType)}';
        break;
      default:
        fieldType = decodeNonListDataType(dataType);
    }
    returnValues.add(ReturnValue(Name: field, DataType: fieldType));
    index++;
  }
  return ValidatedResponse(ReturnValues: returnValues, Error: '');
}

String decodeNonListDataType(String dataType) {
  switch (dataType) {
    case '..':
      return 'Path';
    case '()':
      return 'Node';
    case '->':
    case '<-':
      return 'Relationship';
    case '{}':
      return 'Map';
    case '?':
      return 'Boolean';
    case 'Z':
      return 'Integer';
    case 'R':
      return 'Float';
    case 'U':
      return 'String';
    case 'T':
      return 'Date';
    default:
      return 'Unknown';
  }
}
