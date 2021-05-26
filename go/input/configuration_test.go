package input_test

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/tlarsen7572/graphyx/input"
	"reflect"
	"testing"
	"time"
)

func TestBasicConfig(t *testing.T) {
	config := `<Configuration>
  <JSON>{"ConnStr":"http://localhost:7474","Username":"user","Password":"password","Database":"neo4j","Query":"MATCH p=()-[r:ACTED_IN]-&gt;() RETURN p","LastValidatedResponse":{"Error":"","ReturnValues":[{"Name":"p","DataType":"Path"}]},"Fields":[{"Name":"Field1","DataType":"Integer","Path":[{"Key":"p","DataType":"Path"},{"Key":"Nodes","DataType":"List:Node"},{"Key":"First","DataType":"Node"},{"Key":"ID","DataType":"Integer"}]},{"Name":"Field2","DataType":"String","Path":[{"Key":"p","DataType":"Path"},{"Key":"Relationships","DataType":"List:Relationship"},{"Key":"First","DataType":"Relationship"},{"Key":"Type","DataType":"String"}]}]}</JSON>
</Configuration>`

	decoded, err := input.DecodeConfig(config)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	expected := input.Configuration{
		ConnStr:  `http://localhost:7474`,
		Username: `user`,
		Password: `password`,
		Database: `neo4j`,
		Query:    `MATCH p=()-[r:ACTED_IN]->() RETURN p`,
		Fields: []input.Field{
			{
				Name:     `Field1`,
				DataType: `Integer`,
				Path: []input.Element{
					{Key: `p`, DataType: `Path`},
					{Key: `Nodes`, DataType: `List:Node`},
					{Key: `First`, DataType: `Node`},
					{Key: `ID`, DataType: `Integer`},
				},
			},
			{
				Name:     `Field2`,
				DataType: `String`,
				Path: []input.Element{
					{Key: `p`, DataType: `Path`},
					{Key: `Relationships`, DataType: `List:Relationship`},
					{Key: `First`, DataType: `Relationship`},
					{Key: `Type`, DataType: `String`},
				},
			},
		},
	}
	if !reflect.DeepEqual(expected, decoded) {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, decoded)
	}
	t.Logf(`%v`, decoded)
}

func TestOutgoingRecordInfoFromConfig(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Integer`},
			},
		},
		{
			Name:     `Field2`,
			DataType: `Float`,
			Path: []input.Element{
				{Key: `value`, DataType: `Float`},
			},
		},
		{
			Name:     `Field3`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `String`},
			},
		},
		{
			Name:     `Field4`,
			DataType: `Boolean`,
			Path: []input.Element{
				{Key: `value`, DataType: `Boolean`},
			},
		},
		{
			Name:     `Field5`,
			DataType: `Date`,
			Path: []input.Element{
				{Key: `value`, DataType: `Date`},
			},
		},
		{
			Name:     `Field5`,
			DataType: `DateTime`,
			Path: []input.Element{
				{Key: `value`, DataType: `DateTime`},
			},
		},
	}
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	info := outgoingStuff.RecordInfo
	if count := len(info.IntFields); count != 1 {
		t.Fatalf(`expected 1 int field but got %v`, count)
	}
	if count := len(info.FloatFields); count != 1 {
		t.Fatalf(`expected 1 float field but got %v`, count)
	}
	if count := len(info.StringFields); count != 1 {
		t.Fatalf(`expected 1 string field but got %v`, count)
	}
	if count := len(info.BoolFields); count != 1 {
		t.Fatalf(`expected 1 bool field but got %v`, count)
	}
	if count := len(info.DateTimeFields); count != 2 {
		t.Fatalf(`expected 2 datetime field but got %v`, count)
	}
	t.Logf(`%v`, outgoingStuff)
}

func NewMockRecord(keys []string, values []interface{}) *neo4j.Record {
	return &neo4j.Record{Keys: keys, Values: values}
}

func TestIntegerToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{int64(12345)})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 12345 {
		t.Fatalf(`expected 12345 but got %v`, value)
	}
}

func TestFloatToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Float`,
			Path: []input.Element{
				{Key: `value`, DataType: `Float`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{123.45})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.FloatFields[`Field1`].GetCurrentFloat()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 123.45 {
		t.Fatalf(`expected 123.45 but got %v`, value)
	}
}

func TestBoolToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Boolean`,
			Path: []input.Element{
				{Key: `value`, DataType: `Boolean`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{true})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.BoolFields[`Field1`].GetCurrentBool()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != true {
		t.Fatalf(`expected true but got %v`, value)
	}
}

func TestStringToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{`hello world`})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetCurrentString()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `hello world` {
		t.Fatalf(`expected 'hello world' but got %v`, value)
	}
}

func TestDateToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Date`,
			Path: []input.Element{
				{Key: `value`, DataType: `Date`},
			},
		},
	}
	expectedDate := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	record := NewMockRecord([]string{`value`}, []interface{}{expectedDate})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.DateTimeFields[`Field1`].GetCurrentDateTime()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != expectedDate {
		t.Fatalf(`expected %v but got %v`, expectedDate, value)
	}
}

func TestDateTimeToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `DateTime`,
			Path: []input.Element{
				{Key: `value`, DataType: `DateTime`},
			},
		},
	}
	expectedDate := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	record := NewMockRecord([]string{`value`}, []interface{}{expectedDate})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.DateTimeFields[`Field1`].GetCurrentDateTime()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != expectedDate {
		t.Fatalf(`expected %v but got %v`, expectedDate, value)
	}
}

func TestNodeIdRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `ID`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Id: 23},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 23 {
		t.Fatalf(`expected 23 but got %v`, value)
	}
}

func TestConcatenatedNodeLabelsToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Labels`, DataType: `List:String`},
				{Key: `Concatenate`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Labels: []string{
			`Label1`,
			`Label2`,
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetCurrentString()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `Label1,Label2` {
		t.Fatalf(`expected 'Label1,Label2' but got %v`, value)
	}
}

func TestFirstNodeLabelToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Labels`, DataType: `List:String`},
				{Key: `First`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Labels: []string{
			`Label1`,
			`Label2`,
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetCurrentString()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `Label1` {
		t.Fatalf(`expected 'Label1' but got %v`, value)
	}
}

func TestFirstNodeLabelZeroLabelsToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Labels`, DataType: `List:String`},
				{Key: `First`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Labels: []string{}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetNull()
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestLastNodeLabelToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Labels`, DataType: `List:String`},
				{Key: `Last`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Labels: []string{
			`Label1`,
			`Label2`,
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetCurrentString()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `Label2` {
		t.Fatalf(`expected 'Label2' but got %v`, value)
	}
}

func TestCountNodeLabelToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Labels`, DataType: `List:String`},
				{Key: `Count`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Labels: []string{
			`Label1`,
			`Label2`,
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 2 {
		t.Fatalf(`expected 2 but got %v`, value)
	}
}

func TestLastNodeLabelZeroLabelsToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Labels`, DataType: `List:String`},
				{Key: `Last`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Labels: []string{}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetNull()
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestLabelIndexToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Labels`, DataType: `List:String`},
				{Key: `Index:1`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Labels: []string{
			`Label1`,
			`Label2`,
			`Label3`,
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetCurrentString()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `Label2` {
		t.Fatalf(`expected 'Label2' but got %v`, value)
	}
}

func TestLabelIndexOutOfBoundsToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Labels`, DataType: `List:String`},
				{Key: `Index:1`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Labels: []string{
			`Label1`,
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetNull()
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestNodeStringPropertyToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Properties`, DataType: `Map`},
				{Key: `Something`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Props: map[string]interface{}{
			`Something`: `hello world`,
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetCurrentString()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `hello world` {
		t.Fatalf(`expected 'hello world' but got %v`, value)
	}
}

func TestNodeMissingStringPropertyToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Properties`, DataType: `Map`},
				{Key: `Something`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Props: map[string]interface{}{}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetNull()
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestNodeStringListPropertyToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Properties`, DataType: `Map`},
				{Key: `Something`, DataType: `List:String`},
				{Key: `First`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Props: map[string]interface{}{
			`Something`: []interface{}{
				`hello world`,
				`abcdefg`,
				`wxyz`,
			},
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetCurrentString()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `hello world` {
		t.Fatalf(`expected 'hello world' but got %v`, value)
	}
}

func TestNodeWithWrongDataTypePropertyToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Properties`, DataType: `Map`},
				{Key: `Something`, DataType: `List:String`},
				{Key: `First`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		&neo4j.Node{Props: map[string]interface{}{
			`Something`: `hello world`,
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err == nil {
		t.Fatalf(`expected an error but got none`)
	}
	t.Logf(`error: %v`, err.Error())
}

func TestWrongType(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{`abcdefg`})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err == nil {
		t.Fatalf(`expected an error but got none`)
	}
	t.Logf(`error: %v`, err.Error())
}

func TestRelationshipIdToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Relationship`},
				{Key: `ID`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Relationship{Id: 452},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 452 {
		t.Fatalf(`expected 452 but got %v`, value)
	}
}

func TestRelationshipStartIdToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Relationship`},
				{Key: `StartId`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Relationship{StartId: 452},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 452 {
		t.Fatalf(`expected 452 but got %v`, value)
	}
}

func TestRelationshipEndIdToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Relationship`},
				{Key: `EndId`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Relationship{EndId: 452},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 452 {
		t.Fatalf(`expected 452 but got %v`, value)
	}
}

func TestRelationshipTypeToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `Relationship`},
				{Key: `Type`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Relationship{Type: `hello world`},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetCurrentString()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `hello world` {
		t.Fatalf(`expected 'hello world' but got %v`, value)
	}
}

func TestRelationshipPropertiesToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `Relationship`},
				{Key: `Properties`, DataType: `Map`},
				{Key: `Something`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Relationship{Props: map[string]interface{}{
			`Something`: `hello world`,
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetCurrentString()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `hello world` {
		t.Fatalf(`expected 'hello world' but got %v`, value)
	}
}

func TestRelationshipStringToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `Relationship`},
				{Key: `ToString`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Relationship{Props: map[string]interface{}{
			`Something`: `hello world`,
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetCurrentString()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	expected := `[ {"Something":"hello world"}]`
	if value != expected {
		t.Fatalf(`expected '%v' but got '%v'`, expected, value)
	}
}

func TestPathFirstNodeIdRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Path`},
				{Key: `Nodes`, DataType: `List:Node`},
				{Key: `First`, DataType: `Node`},
				{Key: `ID`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Path{
			Nodes: []neo4j.Node{
				{Id: 234},
				{Id: 534},
				{Id: 9349},
			},
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 234 {
		t.Fatalf(`expected 234 but got %v`, value)
	}
}

func TestPathLastNodeIdRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Path`},
				{Key: `Nodes`, DataType: `List:Node`},
				{Key: `Last`, DataType: `Node`},
				{Key: `ID`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Path{
			Nodes: []neo4j.Node{
				{Id: 234},
				{Id: 534},
				{Id: 9349},
			},
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 9349 {
		t.Fatalf(`expected 9349 but got %v`, value)
	}
}

func TestPathIndexedNodeIdRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Path`},
				{Key: `Nodes`, DataType: `List:Node`},
				{Key: `Index:1`, DataType: `Node`},
				{Key: `ID`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Path{
			Nodes: []neo4j.Node{
				{Id: 234},
				{Id: 534},
				{Id: 9349},
			},
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 534 {
		t.Fatalf(`expected 534 but got %v`, value)
	}
}

func TestPathIndexedNodeOutOfBoundsToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Path`},
				{Key: `Nodes`, DataType: `List:Node`},
				{Key: `Index:3`, DataType: `Node`},
				{Key: `ID`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Path{
			Nodes: []neo4j.Node{
				{Id: 234},
				{Id: 534},
				{Id: 9349},
			},
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetNull()
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestPathCountNodesRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Path`},
				{Key: `Nodes`, DataType: `List:Node`},
				{Key: `Count`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Path{
			Nodes: []neo4j.Node{
				{Id: 234},
				{Id: 534},
				{Id: 9349},
			},
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 3 {
		t.Fatalf(`expected 3 but got %v`, value)
	}
}

func TestPathFirstRelationshipIdToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Path`},
				{Key: `Relationships`, DataType: `List:Relationship`},
				{Key: `First`, DataType: `Relationship`},
				{Key: `ID`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Path{
			Relationships: []neo4j.Relationship{
				{Id: 234},
				{Id: 534},
				{Id: 9349},
			},
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 234 {
		t.Fatalf(`expected 234 but got %v`, value)
	}
}

func TestPathLastRelationshipIdToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Path`},
				{Key: `Relationships`, DataType: `List:Relationship`},
				{Key: `Last`, DataType: `Relationship`},
				{Key: `ID`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Path{
			Relationships: []neo4j.Relationship{
				{Id: 234},
				{Id: 534},
				{Id: 9349},
			},
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 9349 {
		t.Fatalf(`expected 9349 but got %v`, value)
	}
}

func TestPathIndexedRelationshipIdToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Path`},
				{Key: `Relationships`, DataType: `List:Relationship`},
				{Key: `Index:1`, DataType: `Relationship`},
				{Key: `ID`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Path{
			Relationships: []neo4j.Relationship{
				{Id: 234},
				{Id: 534},
				{Id: 9349},
			},
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 534 {
		t.Fatalf(`expected 534 but got %v`, value)
	}
}

func TestPathCountRelationshipsToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Path`},
				{Key: `Relationships`, DataType: `List:Relationship`},
				{Key: `Count`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Path{
			Relationships: []neo4j.Relationship{
				{Id: 234},
				{Id: 534},
				{Id: 9349},
			},
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 3 {
		t.Fatalf(`expected 3 but got %v`, value)
	}
}

func TestFirstStringListToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `List:String`},
				{Key: `First`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		[]interface{}{
			`a`,
			`b`,
			`c`,
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetCurrentString()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `a` {
		t.Fatalf(`expected 'a' but got %v`, value)
	}
}

func TestLastStringListToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `List:String`},
				{Key: `Last`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		[]interface{}{
			`a`,
			`b`,
			`c`,
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetCurrentString()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `c` {
		t.Fatalf(`expected 'c' but got %v`, value)
	}
}

func TestIndexedStringListToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `List:String`},
				{Key: `Index:1`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		[]interface{}{
			`a`,
			`b`,
			`c`,
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetCurrentString()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `b` {
		t.Fatalf(`expected 'b' but got %v`, value)
	}
}

func TestStringListCountRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `List:String`},
				{Key: `Count`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		[]interface{}{
			`a`,
			`b`,
			`c`,
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 3 {
		t.Fatalf(`expected 3 but got %v`, value)
	}
}

func TestConcatenateStringListToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `String`,
			Path: []input.Element{
				{Key: `value`, DataType: `List:String`},
				{Key: `Concatenate`, DataType: `String`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		[]interface{}{
			`a`,
			`b`,
			`c`,
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.StringFields[`Field1`].GetCurrentString()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `a,b,c` {
		t.Fatalf(`expected 'a,b,c' but got %v`, value)
	}
}

func TestFirstIntegerListToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `List:Integer`},
				{Key: `First`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		[]interface{}{
			int64(1),
			int64(2),
			int64(3),
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 1 {
		t.Fatalf(`expected 1 but got %v`, value)
	}
}

func TestLastIntegerListToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `List:Integer`},
				{Key: `Last`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		[]interface{}{
			int64(1),
			int64(2),
			int64(3),
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 3 {
		t.Fatalf(`expected 3 but got %v`, value)
	}
}

func TestIndexedIntegerListToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `List:Integer`},
				{Key: `Index:1`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		[]interface{}{
			int64(1),
			int64(2),
			int64(3),
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 2 {
		t.Fatalf(`expected 2 but got %v`, value)
	}
}

func TestFirstFloatListToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Float`,
			Path: []input.Element{
				{Key: `value`, DataType: `List:Float`},
				{Key: `First`, DataType: `Float`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		[]interface{}{
			1.1,
			2.2,
			3.3,
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.FloatFields[`Field1`].GetCurrentFloat()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 1.1 {
		t.Fatalf(`expected 1.1 but got %v`, value)
	}
}

func TestLastFloatListToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Float`,
			Path: []input.Element{
				{Key: `value`, DataType: `List:Float`},
				{Key: `Last`, DataType: `Float`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		[]interface{}{
			1.1,
			2.2,
			3.3,
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.FloatFields[`Field1`].GetCurrentFloat()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 3.3 {
		t.Fatalf(`expected 3.3 but got %v`, value)
	}
}

func TestCountFloatListToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `List:Float`},
				{Key: `Count`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		[]interface{}{
			1.1,
			2.2,
			3.3,
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 3 {
		t.Fatalf(`expected 3 but got %v`, value)
	}
}

func TestLastBooleanListToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Boolean`,
			Path: []input.Element{
				{Key: `value`, DataType: `List:Boolean`},
				{Key: `Last`, DataType: `Boolean`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		[]interface{}{
			false,
			false,
			true,
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.BoolFields[`Field1`].GetCurrentBool()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != true {
		t.Fatalf(`expected true but got %v`, value)
	}
}

func TestLastDateTimeListToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `DateTime`,
			Path: []input.Element{
				{Key: `value`, DataType: `List:DateTime`},
				{Key: `Last`, DataType: `DateTime`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		[]interface{}{
			neo4j.DateOf(time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)),
			neo4j.LocalDateTimeOf(time.Date(2022, 1, 2, 3, 4, 5, 6, time.UTC)),
			time.Date(2024, 1, 2, 3, 4, 5, 6, time.UTC),
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.DateTimeFields[`Field1`].GetCurrentDateTime()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if expected := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC); value != expected {
		t.Fatalf(`expected 2024-01-02 03:04:05 but got %v`, value)
	}
}

func TestIndexedFloatListToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Float`,
			Path: []input.Element{
				{Key: `value`, DataType: `List:Float`},
				{Key: `Index:1`, DataType: `Float`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		[]interface{}{
			1.1,
			2.2,
			3.3,
		},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.FloatFields[`Field1`].GetCurrentFloat()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 2.2 {
		t.Fatalf(`expected 2.2 but got %v`, value)
	}
}

func TestNodeIntegerListPropertyToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Integer`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Properties`, DataType: `Map`},
				{Key: `Something`, DataType: `List:Integer`},
				{Key: `First`, DataType: `Integer`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Props: map[string]interface{}{
			`Something`: []interface{}{
				int64(1),
				int64(2),
				int64(3),
			},
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.IntFields[`Field1`].GetCurrentInt()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 1 {
		t.Fatalf(`expected 1 but got %v`, value)
	}
}

func TestNodeFloatListPropertyToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Float`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Properties`, DataType: `Map`},
				{Key: `Something`, DataType: `List:Float`},
				{Key: `First`, DataType: `Float`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Props: map[string]interface{}{
			`Something`: []interface{}{
				1.1,
				2.2,
				3.3,
			},
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.FloatFields[`Field1`].GetCurrentFloat()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 1.1 {
		t.Fatalf(`expected 1.1 but got %v`, value)
	}
}

func TestNodeBooleanListPropertyToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `Boolean`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Properties`, DataType: `Map`},
				{Key: `Something`, DataType: `List:Boolean`},
				{Key: `First`, DataType: `Boolean`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Props: map[string]interface{}{
			`Something`: []interface{}{
				true,
				false,
				false,
			},
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.BoolFields[`Field1`].GetCurrentBool()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != true {
		t.Fatalf(`expected true but got %v`, value)
	}
}

func TestNodeDateTimeListPropertyToRecordInfo(t *testing.T) {
	fields := []input.Field{
		{
			Name:     `Field1`,
			DataType: `DateTime`,
			Path: []input.Element{
				{Key: `value`, DataType: `Node`},
				{Key: `Properties`, DataType: `Map`},
				{Key: `Something`, DataType: `List:DateTime`},
				{Key: `First`, DataType: `DateTime`},
			},
		},
	}
	record := NewMockRecord([]string{`value`}, []interface{}{
		neo4j.Node{Props: map[string]interface{}{
			`Something`: []interface{}{
				neo4j.DateOf(time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)),
				neo4j.LocalDateTimeOf(time.Date(2022, 1, 2, 3, 4, 5, 6, time.UTC)),
				time.Date(2024, 1, 2, 3, 4, 5, 6, time.UTC),
			},
		}},
	})
	outgoingStuff, err := input.CreateOutgoingObjects(fields)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if count := len(outgoingStuff.TransferFuncs); count != 1 {
		t.Fatalf(`expected 1 transfer func but got %v`, count)
	}
	err = outgoingStuff.TransferFuncs[0](record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull := outgoingStuff.RecordInfo.DateTimeFields[`Field1`].GetCurrentDateTime()
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if expected := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC); value != expected {
		t.Fatalf(`expected 2020-01-02 but got %v`, value)
	}
}
