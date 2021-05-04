import 'dart:convert';
import 'package:output/bloc.dart';

class AyxToNeo4jMap {
  AyxToNeo4jMap(this.ayxField, this.neo4jField);
  String ayxField;
  String neo4jField;

  Map toJson(){
    return {ayxField: neo4jField};
  }
}
class Configuration extends BlocState {
  Configuration({this.connStr, this.username, this.password, this.database, this.exportObject, this.batchSize, this.nodeLabel, this.nodeIdFields, this.nodePropFields, this.relLabel, this.relPropFields, this.relLeftLabel, this.relLeftFields, this.relRightLabel, this.relRightFields, this.incomingFields});

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
  List<AyxToNeo4jMap> relLeftFields;
  String relRightLabel;
  List<AyxToNeo4jMap> relRightFields;

  List<String> incomingFields;

  void dispose() {}
  Future initialize() async {}

  String getConfig(){
    return json.encode(toJson());
  }

  Map toJson() {
    return {
      "ConnStr": connStr,
      "Username": username,
      "Password": password,
      "Database": database,
      "ExportObject": exportObject,
      "BatchSize": batchSize,
      "NodeLabel": nodeLabel,
      "NodeIdFields": nodeIdFields,
      "NodePropFields": nodePropFields,
      "RelLabel": relLabel,
      "RelPropFields": relPropFields,
      "RelLeftLabel": relLeftLabel,
      "RelLeftFields": relLeftFields.map<Map>((e) => e.toJson()).toList(),
      "RelRightLabel": relRightLabel,
      "RelRightFields": relRightFields.map<Map>((e) => e.toJson()).toList(),
    };
  }
}

Configuration decodeConfig(String configStr, List<String> incomingFields) {
  if (configStr == '' || configStr == null) {
    return Configuration(
      connStr: '',
      username: '',
      password: '',
      database: '',
      exportObject: 'Node',
      batchSize: 10000,
      nodeLabel: '',
      nodeIdFields: [],
      nodePropFields: [],
      relLabel: '',
      relPropFields: [],
      relLeftLabel: '',
      relLeftFields: [],
      relRightLabel: '',
      relRightFields: [],
      incomingFields: incomingFields,
    );
  }
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
    incomingFields: incomingFields,
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

List<AyxToNeo4jMap> decodeFieldMapping(dynamic jsonItem) {
  List<AyxToNeo4jMap> fieldMap = [];
  if (jsonItem == null) {
    return fieldMap;
  }
  for (var entry in jsonItem) {
    var entryMap = entry as Map<String, dynamic>;
    for (var entry in entryMap.entries) {
      fieldMap.add(AyxToNeo4jMap(entry.key, entry.value.toString()));
      break;
    }
  }
  return fieldMap;
}
