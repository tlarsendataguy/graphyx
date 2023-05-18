import 'dart:async';
import 'dart:convert';
import 'dart:js_util';
import 'package:input/bloc.dart';
import 'package:input/configuration.dart';
import 'package:input/decode_config.dart';
import 'package:input/validated_response.dart';
import 'package:rxdart/rxdart.dart' as rx;

enum updated {
  ReturnValues,
  Fields,
}

class Configuration {
  Configuration({this.connStr, this.username, this.password, this.database, this.urlCollapsed, this.query, this.lastValidatedResponse, this.fields});
  String connStr;
  String username;
  String password;
  String database;
  bool   urlCollapsed;
  String query;
  ValidatedResponse lastValidatedResponse;
  List<Field> fields;

  Map toJson(){
    return {
      'ConnStr': connStr,
      'Username': username,
      'Password': password,
      'Query': query,
      'Database': database,
      'UrlCollapsed': urlCollapsed,
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
  String _decrypted;
  bool _decrypting = false;
  String get connStr => _config.connStr;
  set connStr(String value) => _config.connStr = value;
  String get username => _config.username;
  set username(String value) => _config.username = value;

  Future<String> getPassword() async {
    _decrypting = true;
    var event = json.encode({"Event": "Decrypt", "text": _config.password, "callback": "decryptPasswordCallback"});
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

  set password(String value) {
    _decrypted = value;
    var event = json.encode({"Event": "Encrypt", "text": value, "encryptionMode": "", "callback": "encryptPasswordCallback"});
    JsEvent(event);
  }

  String get query => _config.query;
  set query(String value) => _config.query = value;
  String get database => _config.database;
  set database(String value) => _config.database = value;
  bool get urlCollapsed => _config.urlCollapsed;
  set urlCollapsed(bool value) => _config.urlCollapsed = value;

  rx.BehaviorSubject<ValidatedResponse> _lastValidatedResponse;
  Stream get lastValidatedResponse => _lastValidatedResponse.stream;
  rx.BehaviorSubject<List<Field>> _fields;
  Stream get fields => _fields.stream;

  void useDecryptedPassword(String value) {
    _decrypted = value;
    _decrypting = false;
  }

  void useEncryptedPassword(String value) {
    _config.password = value;
  }

  Future validateQuery() async {
    var query = _config.query;
    if (!RegExp("\\sLIMIT\\s").hasMatch(query)) {
      query += " LIMIT 1";
    }

    var uri = _config.connStr;
    var password = await getPassword();
    var db = _config.database == '' ? 'neo4j' : _config.database;
    var user = _config.username;

    var result = await promiseToFuture(returnValueMetadata(uri, user, password, db, query));

    ValidatedResponse validated;
    var error = result.Error;
    var jsonFields = result.Fields;
    if (error != '') {
      validated = ValidatedResponse(error: 'Error: $error', returnValues: []);
    } else {
      var fields = <ReturnValue>[];
      for (var raw in jsonFields) {
        fields.add(ReturnValue(name: raw.Name ?? '', dataType: raw.DataType ?? 'Unknown'));
      }
      validated = ValidatedResponse(error: '', returnValues: fields);
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