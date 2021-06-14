import 'dart:convert';
import 'package:output/bloc.dart';
import 'package:output/configuration.dart';

class AyxToNeo4jMap {
  AyxToNeo4jMap(this.ayxField, this.neo4jField);
  String ayxField;
  String neo4jField;

  Map toJson(){
    return {ayxField: neo4jField};
  }
}

typedef List<String> LazyFieldLoader();

class Configuration extends BlocState {
  Configuration({this.connStr, this.username, this.password, this.database, this.urlCollapsed, this.exportObject, this.batchSize, this.nodeLabel, this.nodeIdFields, this.nodePropFields, this.relLabel, this.relPropFields, this.relLeftLabel, this.relLeftFields, this.relRightLabel, this.relRightFields, this.loadFields});

  String connStr;
  String username;
  String password;
  String database;
  String exportObject;
  bool urlCollapsed;
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
  LazyFieldLoader loadFields;

  bool _decrypting = false;
  String _decrypted;

  List<String> get incomingFields {
    if (_incomingFields == null) {
      _incomingFields = loadFields();
    }
    return _incomingFields;
  }

  List<String> _incomingFields;

  void encryptPassword(String value) {
    _decrypted = value;
    var event = json.encode({"Event": "Encrypt", "text": value, "encryptionMode": "", "callback": "encryptPasswordCallback"});
    JsEvent(event);
  }

  Future<String> decryptPassword() async {
    _decrypting = true;
    var event = json.encode({"Event": "Decrypt", "text": password, "callback": "decryptPasswordCallback"});
    JsEvent(event);
    while (true) {
      if (_decrypting) {
        await Future.delayed(Duration(milliseconds: 10));
        continue;
      }
      break;
    }
    return _decrypted;
  }

  void callbackDecryptedPassword(String value) {
    _decrypted = value;
    _decrypting = false;
  }

  void callbackEncryptedPassword(String value) {
    password = value;
  }

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
      "UrlCollapsed": urlCollapsed,
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

Configuration decodeConfig(String configStr, LazyFieldLoader incomingFields) {
  if (configStr == '' || configStr == null) {
    return Configuration(
      connStr: '',
      username: '',
      password: '',
      database: '',
      urlCollapsed: false,
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
      loadFields: incomingFields,
    );
  }
  var decoded = json.decode(configStr);
  return Configuration(
    connStr: decoded['ConnStr'] ?? '',
    username: decoded['Username'] ?? '',
    password: decoded['Password'] ?? '',
    database: decoded['Database'] ?? '',
    urlCollapsed: decoded['UrlCollapsed'] ?? false,
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
    loadFields: incomingFields,
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
