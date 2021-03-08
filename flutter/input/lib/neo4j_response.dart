import 'dart:convert';

class ValidatedResponse {
  ValidatedResponse({this.error, this.returnValues});
  final String error;
  final List<ReturnValue> returnValues;

  Map toJson(){
    return {
      'Error': error,
      'ReturnValues': returnValues.map((e) => {'Name': e.name, 'DataType': e.dataType}).toList(),
    };
  }
}

class ReturnValue {
  ReturnValue({this.name, this.dataType});
  final String name;
  final String dataType;
}

ValidatedResponse validate(String response) {
  try {
    var decoded = jsonDecode(response);
    if (decoded['errors'] != null) {
      return ValidatedResponse(error: decoded['errors'][0]['message']);
    }
  } catch (_) {}

  var lines = response.split('\n');
  if (lines.length < 2) {
    return ValidatedResponse(error: 'A response with an unexpected format was returned.  Response was:\n$response');
  }

  var header = jsonDecode(lines[0]);

  if (header['error'] != null) {
    return ValidatedResponse(error: header['error']['errors'][0]['message']['U']);
  }

  var data = jsonDecode(lines[1]);
  if (data['data'] == null) {
    return ValidatedResponse(error: 'The query was successful but no records were returned.  No metadata is available to generate output fields.');
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
    returnValues.add(ReturnValue(name: field, dataType: fieldType));
    index++;
  }
  return ValidatedResponse(returnValues: returnValues, error: '');
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
