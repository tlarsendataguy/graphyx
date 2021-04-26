import 'dart:convert';
import 'package:flutter_test/flutter_test.dart';
import 'package:input/decode_config.dart';

void main(){
  test('instantiate normal config', (){
    var configStr = '{"ConnStr": "http://localhost:7474", "Username": "user", "Password": "password", "Database": "neo4j", "Query": "MATCH p=()-[r:ACTED_IN]->() RETURN p", "LastValidatedResponse": {"Error": "", "ReturnValues": [{"Name": "p", "DataType": "Path"}]}, "Fields": [{"Name": "Field1", "DataType": "Path", "Path": [{"Key": "p", "DataType": "Path"}]}]}';
    var decoded = decodeConfig(configStr);
    expect(decoded, isNotNull);
    expect(decoded.connStr, equals('http://localhost:7474'));
    expect(decoded.username, equals('user'));
    expect(decoded.password, equals('password'));
    expect(decoded.database, equals('neo4j'));
    expect(decoded.query, equals('MATCH p=()-[r:ACTED_IN]->() RETURN p'));
    expect(decoded.lastValidatedResponse.error, equals(''));
    expect(decoded.lastValidatedResponse.returnValues.length, equals(1));
    expect(decoded.lastValidatedResponse.returnValues[0].name, equals('p'));
    expect(decoded.lastValidatedResponse.returnValues[0].dataType, equals('Path'));
    expect(decoded.fields.length, equals(1));
    expect(decoded.fields[0].name, equals('Field1'));
    expect(decoded.fields[0].dataType, equals('Path'));
    expect(decoded.fields[0].path.length, equals(1));
    expect(decoded.fields[0].path[0].key, equals('p'));
    expect(decoded.fields[0].path[0].dataType, equals('Path'));
  });

  test('instantiate empty config', (){
    var configStr = '{}';
    var decoded = decodeConfig(configStr);
    expect(decoded.connStr, equals(''));
    expect(decoded.username, equals(''));
    expect(decoded.password, equals(''));
    expect(decoded.query, equals(''));
    expect(decoded.database, equals(''));
    expect(decoded.lastValidatedResponse.error, equals(''));
    expect(decoded.lastValidatedResponse.returnValues.length, equals(0));
    expect(decoded.fields.length, equals(0));
  });

  test('instantiate null config', (){
    var configStr;
    var decoded = decodeConfig(configStr);
    expect(decoded, isNotNull);
  });

  test('instantiate default config', (){
    var config = {
      "ConnStr": "http://localhost:7474",
      "Username": "test",
      "Password": "test",
      "Database": "neo4j",
      "Query": "MATCH p=()-[r:ACTED_IN]->() RETURN p",
      "Fields": [],
      "LastValidatedResponse": {
        "Error": "",
        "ReturnValues": [
          {"Name": "p", "DataType": "Path"}
        ]
      }
    };
    var decoded = decodeConfig(json.encode(config));
    expect(decoded.lastValidatedResponse.returnValues.length, equals(1));
    print(json.encode(decoded.toJson()));
  });
}