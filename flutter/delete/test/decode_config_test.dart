import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:delete/decode_config.dart';

void main() {
  test('instantiate normal config',(){
    var configStr = '{"ConnStr":"http://localhost:7474","Username":"user","Password":"password","Database":"neo4j","DeleteObject":"Node","BatchSize":10000,"NodeLabel":"DELETE","NodeIdFields":["Id"],"RelType":"Relates_To","RelFields":["RelProp1"],"RelLeftLabel":"Left","RelLeftFields":[{"LeftId":"Id"}],"RelRightLabel":"Right","RelRightFields":[{"RightId":"Id"}]}';
    var decoded = decodeConfig(configStr, ()=>[]);

    expect(decoded, isNotNull);
    expect(decoded.connStr, equals('http://localhost:7474'));
    expect(decoded.username, equals('user'));
    expect(decoded.password, equals("password"));
    expect(decoded.database, equals("neo4j"));
    expect(decoded.deleteObject, equals("Node"));
    expect(decoded.batchSize, equals(10000));
    expect(decoded.nodeLabel, equals("DELETE"));
    expect(decoded.relType, equals("Relates_To"));
    expect(decoded.relLeftLabel, equals("Left"));
    expect(decoded.relRightLabel, equals("Right"));

    expect(decoded.nodeIdFields, equals(["Id"]));
    expect(decoded.relFields, equals(["RelProp1"]));

    expect(decoded.relLeftFields.length, equals(1));
    expect(decoded.relLeftFields[0].ayxField, equals('LeftId'));
    expect(decoded.relLeftFields[0].neo4jField, equals('Id'));
    expect(decoded.relRightFields[0].ayxField, equals('RightId'));
    expect(decoded.relRightFields[0].neo4jField, equals('Id'));
  });

  test('instantiate empty config',(){
    var configStr = '';
    var decoded = decodeConfig(configStr, ()=>[]);

    expect(decoded, isNotNull);
    expect(decoded.connStr, equals(''));
    expect(decoded.username, equals(''));
    expect(decoded.password, equals(""));
    expect(decoded.database, equals(""));
    expect(decoded.deleteObject, equals("Node"));
    expect(decoded.batchSize, equals(10000));
    expect(decoded.nodeLabel, equals(""));
    expect(decoded.relType, equals(""));
    expect(decoded.relLeftLabel, equals(""));
    expect(decoded.relRightLabel, equals(""));

    expect(decoded.nodeIdFields, equals([]));
    expect(decoded.relFields, equals([]));

    expect(decoded.relLeftFields, equals([]));
    expect(decoded.relRightFields, equals([]));
  });

  test("convert config to json", (){
    var config = Configuration(
      connStr: 'http://localhost:7474',
      username: 'user',
      password: 'password',
      database: 'neo4j',
      deleteObject: 'Node',
      batchSize: 10000,
      nodeLabel: 'TestNode',
      nodeIdFields: ['Node1'],
      relType: 'TestRel',
      relFields: ['Rel1'],
      relLeftLabel: 'TestLeftNode',
      relLeftFields: [AyxToNeo4jMap('Ayx1', 'Neo4j1')],
      relRightLabel: 'TestRightNode',
      relRightFields: [AyxToNeo4jMap('Ayx2', 'Neo4j2')]
    );

    var jsonObj = config.toJson();
    var jsonString = json.encode(jsonObj);
    print(jsonString);
  });
}