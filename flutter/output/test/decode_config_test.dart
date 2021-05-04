import 'package:flutter_test/flutter_test.dart';
import 'package:output/decode_config.dart';

void main() {
  test('instantiate normal config',(){
    var configStr = '{"ConnStr":"http://localhost:7474","Username":"user","Password":"password","Database":"neo4j","ExportObject":"Node","BatchSize":10000,"NodeLabel":"SomeNode","NodeIdFields":["NodeId1"],"NodePropFields":["NodeProp1","NodeProp2"],"RelLabel":"SomeRel","RelPropFields":["RelProp1"],"RelLeftLabel":"Left","RelLeftFields":[{"AyxField1":"Neo4jField1"}],"RelRightLabel":"Right","RelRightFields":[{"AyxField2":"Neo4jField2"},{"AyxField3":"Neo4jField3"}]}';
    var decoded = decodeConfig(configStr, []);

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
    expect(decoded.relLeftFields[0].key, equals('AyxField1'));
    expect(decoded.relLeftFields[0].value, equals('Neo4jField1'));
    expect(decoded.relRightFields.length, equals(2));
    expect(decoded.relRightFields[0].key, equals('AyxField2'));
    expect(decoded.relRightFields[0].value, equals('Neo4jField2'));
    expect(decoded.relRightFields[1].key, equals('AyxField3'));
    expect(decoded.relRightFields[1].value, equals('Neo4jField3'));
  });

  test('instantiate empty config',(){
    var configStr = '';
    var decoded = decodeConfig(configStr, []);

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
}