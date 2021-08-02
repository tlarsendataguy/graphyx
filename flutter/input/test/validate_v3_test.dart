import 'package:flutter_test/flutter_test.dart';
import 'package:input/validate_v3.dart';

main(){
  test("normal response of string and node", (){
    var response = """{"results":[{"columns":["n.title","n"],"data":[{"row":["The Matrix",{"tagline":"Welcome to the Real World","title":"The Matrix","released":1999}],"meta":[null,{"id":0,"type":"node","deleted":false}]}]}],"errors":[]}""";

    var validated = validateV3Response(response);
    expect(validated, isNotNull);
    expect(validated.error, equals(''));
    expect(validated.returnValues.length, equals(2));
    expect(validated.returnValues[0].dataType, equals('String'));
    expect(validated.returnValues[0].name, equals('n.title'));
    expect(validated.returnValues[1].dataType, equals('Node'));
    expect(validated.returnValues[1].name, equals('n'));
  });
}