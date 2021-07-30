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
