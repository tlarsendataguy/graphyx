import 'package:flutter_test/flutter_test.dart';
import 'package:input/validate_v3.dart';

main(){
  test("normal response", (){
    var response = """{"results":[{"columns":["n.title","n","date()","ID(n)","properties(n)"],"data":[{"row":["The Matrix",{"tagline":"Welcome to the Real World","title":"The Matrix","released":1999},"2021-08-03",0,{"tagline":"Welcome to the Real World","title":"The Matrix","released":1999}],"meta":[null,{"id":0,"type":"node","deleted":false},{"type":"date"},null,null,null,null]}]}],"errors":[]}""";

    var validated = validateV3Response(response);
    expect(validated, isNotNull);
    expect(validated.error, equals(''));
    expect(validated.returnValues.length, equals(5));
    expect(validated.returnValues[0].dataType, equals('String'));
    expect(validated.returnValues[0].name, equals('n.title'));
    expect(validated.returnValues[1].dataType, equals('Node'));
    expect(validated.returnValues[1].name, equals('n'));
    expect(validated.returnValues[2].dataType, equals('Date'));
    expect(validated.returnValues[2].name, equals('date()'));
    expect(validated.returnValues[3].dataType, equals('Integer'));
    expect(validated.returnValues[3].name, equals('ID(n)'));
    expect(validated.returnValues[4].dataType, equals('Map'));
    expect(validated.returnValues[4].name, equals('properties(n)'));
  });
}