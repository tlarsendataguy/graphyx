

class Config {
  Config({this.connStr, this.username, this.password, this.query, this.fields});

  String connStr;
  String username;
  String password;
  String query;
  List<Field> fields;

  Map<String, dynamic> toMap() {
    return {
      'ConnStr': connStr,
      'Username': username,
      'Password': password,
      'Query': query,
      'Fields': fieldsToJson(fields),
    };
  }

  String toString() {
    return '{connStr="$connStr", username="$username", query="$query", path=${fields.toString()}}';
  }
}

List<Map<String,dynamic>> fieldsToJson(List<Field> fields) {
  List<Map<String,dynamic>> mapList = [];
  for (var field in fields) {
    mapList.add(field.toMap());
  }
  return mapList;
}

List<Field> jsonToFields(List<Map<String, dynamic>> json) {
  List<Field> fields = [];
  for (var jsonField in json) {
    var name = jsonField['Name'];
    var dataType = jsonField['DataType'];
    var elements = _jsonToPath(jsonField['Path']);
    if (name == null || dataType == null || elements == null){
      continue;
    }
    if (!(name is String && dataType is String && elements is List<Element>)) {
      continue;
    }
    fields.add(Field(name: name, dataType: dataType, path: elements));
  }
  return fields;
}

class Field {
  Field({this.name, this.dataType, this.path});

  String name;
  String dataType;
  List<Element> path;

  Map<String, dynamic> toMap() {
    return {
      'Name': name,
      'DataType': dataType,
      'Path': _pathToJson(path),
    };
  }

  String toString() {
    return '{name="$name", dataType="$dataType", path=${path.toString()}}';
  }
}

List<Map<String, dynamic>> _pathToJson(List<Element> elements) {
  List<Map<String, dynamic>> path = [];
  for (var element in elements) {
    path.add(element.toMap());
  }
  return path;
}

List<Element> _jsonToPath(dynamic json) {
  if (json == null) {
    return null;
  }
  if (!(json is List<Map<String, dynamic>>)) {
    return null;
  }
  List<Element> elements = [];
  for (var jsonElement in json) {
    var key = jsonElement['Key'];
    var dataType = jsonElement['DataType'];
    if (key == null || dataType == null) {
      continue;
    }
    if (!(key is String && dataType is String)) {
      continue;
    }
    elements.add(Element(key: key, dataType: dataType));
  }
  return elements;
}

class Element {
  Element({this.key, this.dataType});

  String key;
  String dataType;

  Map<String, dynamic> toMap() {
    return {'Key': key, 'DataType': dataType};
  }

  String toString() {
    return '{key="$key", dataType="$dataType"}';
  }
}
