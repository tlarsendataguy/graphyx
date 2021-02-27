import 'package:flutter_test/flutter_test.dart';
import 'package:input/fields_conv.dart';

Config testConfig = Config(
  connStr: 'bolt://localhost:7687',
  username: 'test',
  password: 'test',
  query: 'MATCH p=()-[r:ACTED_IN]->() RETURN p',
  fields: [
    Field(name: 'Field1', dataType: 'Integer', path: [
      Element(key: 'p', dataType: 'List:Integer'),
      Element(key: 'First', dataType: 'Integer'),
    ]),
    Field(name: 'Field2', dataType: 'String', path: [
      Element(key: 's', dataType: 'String'),
    ]),
  ],
);

Map<String, dynamic> testMap = {
  'ConnStr': 'bolt://localhost:7687',
  'Username': 'test',
  'Password': 'test',
  'Query': 'MATCH p=()-[r:ACTED_IN]->() RETURN p',
  'Fields': [
    {
      'Name': 'Field1',
      'DataType': 'Integer',
      'Path': [
        {'Key': 'p', 'DataType': 'List:Integer'},
        {'Key': 'First', 'DataType': 'Integer'},
      ],
    },
    {
      'Name': 'Field2',
      'DataType': 'String',
      'Path': [
        {'Key': 's', 'DataType': 'String'},
      ],
    },
  ]
};

void main(){
  test("convert config to json-acceptable map", (){
    var result = testConfig.toMap();
    expect(result, equals(testMap));
  });

  test("convert json map to config", (){
    var result = configFromMap(testMap);
    print('$result');

    var actualFields = result.fields;
    var expectedFields = testConfig.fields;
    expect(actualFields.length, equals(expectedFields.length));
    expect(actualFields.length, equals(2));
    expect(actualFields[0].name, equals(expectedFields[0].name));
    expect(actualFields[0].path.length, equals(expectedFields[0].path.length));
    expect(actualFields[0].path.length, equals(2));
  });

  test("convert empty json map to config", (){
    var result = configFromMap({});
    expect(result.query, equals(''));
    expect(result.connStr, equals(''));
    expect(result.username, equals(''));
    expect(result.password, equals(''));
    expect(result.fields, equals([]));
  });
}