import 'package:flutter_test/flutter_test.dart';
import 'package:input/fields_conv.dart';

void main(){
  test("convert list of fields to json-acceptable list", (){
    List<Field> fields = [
      Field(name: 'Field1', dataType: 'Integer', path: [
        Element(key: 'p', dataType: 'List:Integer'),
        Element(key: 'First', dataType: 'Integer'),
      ]),
      Field(name: 'Field2', dataType: 'String', path: [
        Element(key: 's', dataType: 'String'),
      ]),
    ];

    var result = fieldsToJson(fields);
    var expected = [
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
    ];
    expect(result, equals(expected));
  });

  test("convert json list of fields to class objects", (){
    var json = [
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
    ];
    var result = jsonToFields(json);
    List<Field> expected = [
      Field(name: 'Field1', dataType: 'Integer', path: [
        Element(key: 'p', dataType: 'List:Integer'),
        Element(key: 'First', dataType: 'Integer'),
      ]),
      Field(name: 'Field2', dataType: 'String', path: [
        Element(key: 's', dataType: 'String'),
      ]),
    ];
    print('$result');
    expect(result.length, equals(expected.length));
    expect(result.length, equals(2));
    expect(result[0].name, equals(expected[0].name));
    expect(result[0].path.length, equals(result[0].path.length));
    expect(result[0].path.length, equals(2));
  });
}