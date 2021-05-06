import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:output/decode_config.dart';

void main() {
  test('instantiate normal config',(){
    var configStr = '{"ConnStr":"http://localhost:7474","Username":"user","Password":"password","Database":"neo4j","ExportObject":"Node","BatchSize":10000,"NodeLabel":"SomeNode","NodeIdFields":["NodeId1"],"NodePropFields":["NodeProp1","NodeProp2"],"RelLabel":"SomeRel","RelPropFields":["RelProp1"],"RelLeftLabel":"Left","RelLeftFields":[{"AyxField1":"Neo4jField1"}],"RelRightLabel":"Right","RelRightFields":[{"AyxField2":"Neo4jField2"},{"AyxField3":"Neo4jField3"}]}';
    var decoded = decodeConfig(configStr, ()=>[]);

    expect(decoded, isNotNull);
    expect(decoded.connStr, equals('http://localhost:7474'));
    expect(decoded.username, equals('user'));
    expect(decoded.password, equals("password"));
    expect(decoded.database, equals("neo4j"));
    expect(decoded.exportObject, equals("Node"));
    expect(decoded.batchSize, equals(10000));
    expect(decoded.nodeLabel, equals("SomeNode"));
    expect(decoded.relLabel, equals("SomeRel"));
    expect(decoded.relLeftLabel, equals("Left"));
    expect(decoded.relRightLabel, equals("Right"));

    expect(decoded.nodeIdFields, equals(["NodeId1"]));
    expect(decoded.nodePropFields, equals(["NodeProp1","NodeProp2"]));
    expect(decoded.relPropFields, equals(["RelProp1"]));

    expect(decoded.relLeftFields.length, equals(1));
    expect(decoded.relLeftFields[0].ayxField, equals('AyxField1'));
    expect(decoded.relLeftFields[0].neo4jField, equals('Neo4jField1'));
    expect(decoded.relRightFields.length, equals(2));
    expect(decoded.relRightFields[0].ayxField, equals('AyxField2'));
    expect(decoded.relRightFields[0].neo4jField, equals('Neo4jField2'));
    expect(decoded.relRightFields[1].ayxField, equals('AyxField3'));
    expect(decoded.relRightFields[1].neo4jField, equals('Neo4jField3'));
  });

  test('instantiate empty config',(){
    var configStr = '';
    var decoded = decodeConfig(configStr, ()=>[]);

    expect(decoded, isNotNull);
    expect(decoded.connStr, equals(''));
    expect(decoded.username, equals(''));
    expect(decoded.password, equals(""));
    expect(decoded.database, equals(""));
    expect(decoded.exportObject, equals("Node"));
    expect(decoded.batchSize, equals(10000));
    expect(decoded.nodeLabel, equals(""));
    expect(decoded.relLabel, equals(""));
    expect(decoded.relLeftLabel, equals(""));
    expect(decoded.relRightLabel, equals(""));

    expect(decoded.nodeIdFields, equals([]));
    expect(decoded.nodePropFields, equals([]));
    expect(decoded.relPropFields, equals([]));

    expect(decoded.relLeftFields, equals([]));
    expect(decoded.relRightFields, equals([]));
  });

  test("convert config to json", (){
    var config = Configuration(
      connStr: 'http://localhost:7474',
      username: 'user',
      password: 'password',
      database: 'neo4j',
      exportObject: 'Node',
      batchSize: 10000,
      nodeLabel: 'TestNode',
      nodeIdFields: ['Node1'],
      nodePropFields: ['Node2'],
      relLabel: 'TestRel',
      relPropFields: ['Rel1'],
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