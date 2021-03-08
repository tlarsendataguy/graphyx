
import 'dart:convert';

import 'package:input/bloc.dart';
import 'package:input/decode_config.dart';
import 'package:input/neo4j_response.dart';
import 'package:rxdart/rxdart.dart' as rx;
import 'package:http/http.dart' as http;

enum updated {
  ReturnValues,
  Fields,
}

class Configuration {
  Configuration({this.connStr, this.username, this.password, this.query, this.lastValidatedResponse, this.fields});
  String connStr;
  String username;
  String password;
  String query;
  ValidatedResponse lastValidatedResponse;
  List<Field> fields;

  Map toJson(){
    return {
      'ConnStr': connStr,
      'Username': username,
      'Password': password,
      'Query': query,
      'LastValidatedResponse': lastValidatedResponse.toJson(),
      'Fields': fields.map((e) => e.toJson()).toList(),
    };
  }
}

class Field {
  Field({this.name, this.dataType, this.path});
  String name;
  String dataType;
  List<PathElement> path;

  Map toJson(){
    return {
      'Name': name,
      'DataType': dataType,
      'Path': path.map((e) => {'Key': e.key, 'DataType': e.dataType}).toList(),
    };
  }
}

class PathElement {
  PathElement({this.key, this.dataType});
  final String key;
  final String dataType;
}

class AppState extends BlocState {
  AppState(String config) {
    _config = decodeConfig(config);
    _lastValidatedResponse = rx.BehaviorSubject<ValidatedResponse>.seeded(_config.lastValidatedResponse);
    _fields = rx.BehaviorSubject<List<Field>>.seeded(_config.fields);
  }

  Configuration _config;
  String get connStr => _config.connStr;
  set connStr(String value) => _config.connStr = value;
  String get username => _config.username;
  set username(String value) => _config.username = value;
  String get password => _config.password;
  set password(String value) => _config.password = value;
  String get query => _config.query;
  set query(String value) => _config.query = value;

  rx.BehaviorSubject<ValidatedResponse> _lastValidatedResponse;
  Stream get lastValidatedResponse => _lastValidatedResponse.stream;
  rx.BehaviorSubject<List<Field>> _fields;
  Stream get fields => _fields.stream;

  Future validateQuery() async {
    var query = _config.query;
    if (!RegExp("\\sLIMIT\\s").hasMatch(query)) {
      query += " LIMIT 1";
    }

    ValidatedResponse validated;
    try {
      var body = {
        "statements": [
          {
            "statement": query,
            "parameters": {},
          },
        ],
      };
      var response = await http.post(
        '${_config.connStr}/db/neo4j/tx/commit',
        headers: {
          'Accept': 'application/vnd.neo4j.jolt+json-seq;strict=true',
          'Content-Type': 'application/json',
          'Authorization': 'Basic ' + base64Encode(
              utf8.encode('${_config.username}:${_config.password}')),
        },
        body: jsonEncode(body),
      );
      validated = validate(response.body);
    }
    catch (ex) {
      validated = ValidatedResponse(error: 'Unable to connect to the Neo4j database.  Double-check the URL make sure you have a working network connection to the database.');
    }
    _config.lastValidatedResponse = validated;
    _lastValidatedResponse.add(validated);
  }

  void moveField(int from, int to) {
    var field = _config.fields[from];
    if (to > from) {
      _config.fields.insert(to, field);
      _config.fields.removeAt(from);
    } else {
      _config.fields.removeAt(from);
      _config.fields.insert(to, field);
    }
    _fields.add(_config.fields);
  }

  void removeField(int at) {
    _config.fields.removeAt(at);
    _fields.add(_config.fields);
  }

  void addField(){
    _config.fields.add(Field(name: '', dataType: '', path: []));
    _fields.add(_config.fields);
  }

  String getConfig(){
    return json.encode(_config.toJson());
  }

  void dispose() {
    _lastValidatedResponse.close();
    _fields.close();
  }

  Future initialize() {}
}