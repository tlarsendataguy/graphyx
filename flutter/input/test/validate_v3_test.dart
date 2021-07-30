import 'package:flutter_test/flutter_test.dart';
import 'package:input/validate_v3.dart';

main(){
  test("normal response of string and node", (){
    var response = """{"results":[{"columns":["n.title","n"],"data":[{"row":["The Matrix",{"tagline":"Welcome to the Real World","title":"The Matrix","released":1999}],"meta":[null,{"id":0,"type":"node","deleted":false}]}]}],"errors":[]}""";

    var validated = validateV3Response(response);
    expect(validated, isNotNull);
  });
}