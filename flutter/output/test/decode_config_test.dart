

import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';

class Configuration {
  Configuration({this.connStr, this.username, this.password, this.database, this.exportObject, this.batchSize, this.nodeLabel, this.nodeIdFields, this.nodePropFields, this.relLabel, this.relPropFields, this.relLeftLabel, this.relLeftFields, this.relRightLabel, this.relRightFields});
  String connStr;
  String username;
  String password;
  String database;
  String exportObject;
  int batchSize;
  String nodeLabel;
  List<String> nodeIdFields;
  List<String> nodePropFields;
  String relLabel;
  List<String> relPropFields;
  String relLeftLabel;
  List<Map<String, String>> relLeftFields;
  String relRightLabel;
  List<Map<String, String>> relRightFields;
}

Configuration decodeConfig(String configStr) {
  var decoded = json.decode(configStr);
  return Configuration(
    connStr: decoded['ConnStr'] ?? '',
    username: decoded['Username'] ?? '',
    password: decoded['Password'] ?? '',
    database: decoded['Database'] ?? '',
    exportObject: decoded['ExportObject'] ?? 'Node',
    batchSize: decoded['BatchSize'] ?? 10000,
    nodeLabel: decoded['NodeLabel'] ?? '',
    nodeIdFields: decodeStringList(decoded['NodeIdFields']),
    nodePropFields: decodeStringList(decoded['NodePropFields']),
    relLabel: decoded['RelLabel'] ?? '',
    relPropFields: decodeStringList(decoded['RelPropFields']),
    relLeftLabel: decoded['RelLeftLabel'] ?? '',
    relLeftFields: decodeFieldMapping(decoded['RelLeftFields']),
    relRightLabel: decoded['RelRightLabel'] ?? '',
    relRightFields: decodeFieldMapping(decoded['RelRightFields']),
  );
}

List<String> decodeStringList(dynamic strings) {
  List<String> list = [];
  if (strings == null) {
    return list;
  }
  for (var entry in strings) {
    list.add(entry.toString());
  }
  return list;
}

List<Map<String, String>> decodeFieldMapping(dynamic jsonItem) {
  List<Map<String, String>> fieldMap = [];
  if (jsonItem == null) {
    return fieldMap;
  }
  for (var entry in jsonItem) {
    var entryMap = entry as Map<String, dynamic>;
    for (var entry in entryMap.entries) {
      fieldMap.add({entry.key: entry.value.toString()});
    }
  }
  return fieldMap;
}

void main() {
  test('instantiate normal config',(){
    var configStr = '{"ConnStr":"http://localhost:7474","Username":"user","Password":"password","Database":"neo4j","ExportObject":"Node","BatchSize":10000,"NodeLabel":"SomeNode","NodeIdFields":["NodeId1"],"NodePropFields":["NodeProp1","NodeProp2"],"RelLabel":"SomeRel","RelPropFields":["RelProp1"],"RelLeftLabel":"Left","RelLeftFields":[{"AyxField1":"Neo4jField1"}],"RelRightLabel":"Right","RelRightFields":[{"AyxField2":"Neo4jField2"},{"AyxField3":"Neo4jField3"}]}';
    var decoded = decodeConfig(configStr);
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

    expect(decoded.relLeftFields, equals([{"AyxField1":"Neo4jField1"}]));
    expect(decoded.relRightFields, equals([{"AyxField2":"Neo4jField2"},{"AyxField3":"Neo4jField3"}]));
  });
}