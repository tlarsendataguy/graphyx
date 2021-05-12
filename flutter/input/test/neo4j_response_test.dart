import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:input/neo4j_response.dart';

void main(){
  test("error response",(){
    var response = '''{
  "errors" : [ {
    "code" : "Neo.ClientError.Security.Unauthorized",
    "message" : "Invalid username or password."
  } ]
}''';

    var validatedResponse = validate(response);
    expect(validatedResponse, isNotNull);
    expect(validatedResponse.error, isNot(''));
    print(validatedResponse.error);
  });

  test("normal response",(){
    var response = '''{"header":{"fields":["p"]}}
{"data":[{"..":[{"()":[8,["Person"],{"born":{"Z":"1978"},"name":{"U":"Emil Eifrem"}}]},{"->":[7,8,"ACTED_IN",0,{"roles":[{"U":"Emil"}]}]},{"()":[0,["Movie"],{"tagline":{"U":"Welcome to the Real World"},"title":{"U":"The Matrix"},"released":{"Z":"1999"}}]}]}]}
{"summary":{}}
{"info":{}}''';

    var validatedResponse = validate(response);
    expect(validatedResponse.error, equals(''));
    expect(validatedResponse.returnValues.length, equals(1));
    expect(validatedResponse.returnValues[0].name, equals('p'));
    expect(validatedResponse.returnValues[0].dataType, equals('Path'));
    print(json.encode(validatedResponse.toJson()));
  });

  test("multiple return values with a list",(){
    var response = '''{"header":{"fields":["p","nodes(p)"]}}
{"data":[{"..":[{"()":[8,["Person"],{"born":{"Z":"1978"},"name":{"U":"Emil Eifrem"}}]},{"->":[7,8,"ACTED_IN",0,{"roles":[{"U":"Emil"}]}]},{"()":[0,["Movie"],{"tagline":{"U":"Welcome to the Real World"},"title":{"U":"The Matrix"},"released":{"Z":"1999"}}]}]},{"[]":[{"()":[8,["Person"],{"born":{"Z":"1978"},"name":{"U":"Emil Eifrem"}}]},{"()":[0,["Movie"],{"tagline":{"U":"Welcome to the Real World"},"title":{"U":"The Matrix"},"released":{"Z":"1999"}}]}]}]}
{"summary":{}}
{"info":{}}''';

    var validatedResponse = validate(response);
    expect(validatedResponse.error, equals(''));
    expect(validatedResponse.returnValues.length, equals(2));
    expect(validatedResponse.returnValues[0].name, equals('p'));
    expect(validatedResponse.returnValues[0].dataType, equals('Path'));
    expect(validatedResponse.returnValues[1].name, equals('nodes(p)'));
    expect(validatedResponse.returnValues[1].dataType, equals('List:Node'));
    print(json.encode(validatedResponse.toJson()));
  });

  test("invalid json",(){
    var response = '''invalid json''';
    var validatedResponse = validate(response);
    expect(validatedResponse, isNotNull);
    print(json.encode(validatedResponse.toJson()));
  });


  test("normal response with record separators",(){
    var response = '''\u001E{"header":{"fields":["p"]}}
\u001E{"data":[{"..":[{"()":[8,["Person"],{"born":{"Z":"1978"},"name":{"U":"Emil Eifrem"}}]},{"->":[7,8,"ACTED_IN",0,{"roles":[{"U":"Emil"}]}]},{"()":[0,["Movie"],{"tagline":{"U":"Welcome to the Real World"},"title":{"U":"The Matrix"},"released":{"Z":"1999"}}]}]}]}
\u001E{"summary":{}}
\u001E{"info":{}}''';

    var validatedResponse = validate(response);
    expect(validatedResponse.error, equals(''));
    expect(validatedResponse.returnValues.length, equals(1));
    expect(validatedResponse.returnValues[0].name, equals('p'));
    expect(validatedResponse.returnValues[0].dataType, equals('Path'));
    print(json.encode(validatedResponse.toJson()));
  });

}