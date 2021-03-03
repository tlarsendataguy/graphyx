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
    expect(validatedResponse.Error, isNot(''));
    print(validatedResponse.Error);
  });

  test("normal response",(){
    var response = '''{"header":{"fields":["p"]}}
{"data":[{"..":[{"()":[8,["Person"],{"born":{"Z":"1978"},"name":{"U":"Emil Eifrem"}}]},{"->":[7,8,"ACTED_IN",0,{"roles":[{"U":"Emil"}]}]},{"()":[0,["Movie"],{"tagline":{"U":"Welcome to the Real World"},"title":{"U":"The Matrix"},"released":{"Z":"1999"}}]}]}]}
{"summary":{}}
{"info":{}}''';

    var validatedResponse = validate(response);
    expect(validatedResponse.Error, equals(''));
    expect(validatedResponse.ReturnValues.length, equals(1));
    expect(validatedResponse.ReturnValues[0].ReturnValue.Name, equals('p'));
    expect(validatedResponse.ReturnValues[0].ReturnValue.DataType, equals('Path'));
  });

  test("multiple return values with a list",(){
    var response = '''{"header":{"fields":["p","nodes(p)"]}}
{"data":[{"..":[{"()":[8,["Person"],{"born":{"Z":"1978"},"name":{"U":"Emil Eifrem"}}]},{"->":[7,8,"ACTED_IN",0,{"roles":[{"U":"Emil"}]}]},{"()":[0,["Movie"],{"tagline":{"U":"Welcome to the Real World"},"title":{"U":"The Matrix"},"released":{"Z":"1999"}}]}]},{"[]":[{"()":[8,["Person"],{"born":{"Z":"1978"},"name":{"U":"Emil Eifrem"}}]},{"()":[0,["Movie"],{"tagline":{"U":"Welcome to the Real World"},"title":{"U":"The Matrix"},"released":{"Z":"1999"}}]}]}]}
{"summary":{}}
{"info":{}}''';

    var validatedResponse = validate(response);
    expect(validatedResponse.Error, equals(''));
    expect(validatedResponse.ReturnValues.length, equals(2));
    expect(validatedResponse.ReturnValues[0].ReturnValue.Name, equals('p'));
    expect(validatedResponse.ReturnValues[0].ReturnValue.DataType, equals('Path'));
    expect(validatedResponse.ReturnValues[1].ReturnValue.Name, equals('nodes(p)'));
    expect(validatedResponse.ReturnValues[1].ReturnValue.DataType, equals('List:Node'));
  });
}