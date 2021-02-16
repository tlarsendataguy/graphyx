package input_test

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/tlarsen7572/graphyx/input"
	"reflect"
	"testing"
	"time"
)

func TestBasicConfig(t *testing.T) {
	config := `<Configuration>
	<ConnStr>bolt://localhost:7687</ConnStr>
	<Username>user</Username>
	<Password>password</Password>
	<Query>MATCH p=()-[r:ACTED_IN]->() RETURN p</Query>
	<Fields>
		<Field Name="Field1" DataType="Integer">
			<Path>
				<Element DataType="Path" Key="p" />
				<Element DataType="List:Node" Key="Nodes" />
				<Element DataType="Node" Key="0" />
				<Element DataType="Integer" Key="ID" />
			</Path>
		</Field>
		<Field Name="Field2" DataType="Integer">
			<Path>
				<Element DataType="Path" Key="p" />
				<Element DataType="List:Relationship" Key="Relationships" />
				<Element DataType="Relationship" Key="0" />
				<Element DataType="Integer" Key="ID" />
			</Path>
		</Field>
	</Fields>
</Configuration>`

	decoded, err := input.DecodeConfig(config)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	expected := input.Configuration{
		ConnStr:  `bolt://localhost:7687`,
		Username: `user`,
		Password: `password`,
		Query:    `MATCH p=()-[r:ACTED_IN]->() RETURN p`,
		Fields: []input.Field{
			{
				Name:     `Field1`,
				DataType: `Integer`,
				Path: []input.Element{
					{Key: `p`, DataType: `Path`},
					{Key: `Nodes`, DataType: `List:Node`},
					{Key: `0`, DataType: `Node`},
					{Key: `ID`, DataType: `Integer`},
				},
			},
			{
				Name:     `Field2`,
				DataType: `Integer`,
				Path: []input.Element{
					{Key: `p`, DataType: `Path`},
					{Key: `Relationships`, DataType: `List:Relationship`},
					{Key: `0`, DataType: `Relationship`},
					{Key: `ID`, DataType: `Integer`},
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

func NewMockRecord(keys []string, values []interface{}) *MockRecord {
	return &MockRecord{keys: keys, values: values}
}

type MockRecord struct {
	keys   []string
	values []interface{}
}

func (r *MockRecord) Keys() []string {
	return r.keys
}

func (r *MockRecord) Values() []interface{} {
	return r.values
}

func (r *MockRecord) Get(key string) (interface{}, bool) {
	for index, foundKey := range r.keys {
		if foundKey == key {
			return r.values[index], true
		}
	}
	return nil, false
}

func (r *MockRecord) GetByIndex(index int) interface{} {
	return r.values[index]
}

type MockNode struct {
	id     int64
	labels []string
	props  map[string]interface{}
}

func (n *MockNode) Id() int64 {
	return n.id
}

func (n *MockNode) Labels() []string {
	return n.labels
}

func (n *MockNode) Props() map[string]interface{} {
	return n.props
}

type MockRelationship struct {
	id      int64
	startId int64
	endId   int64
	relType string
	props   map[string]interface{}
}

func (r *MockRelationship) Id() int64 {
	return r.id
}

func (r *MockRelationship) StartId() int64 {
	return r.startId
}

func (r *MockRelationship) EndId() int64 {
	return r.endId
}

func (r *MockRelationship) Type() string {
	return r.relType
}

func (r *MockRelationship) Props() map[string]interface{} {
	return r.props
}

type MockPath struct {
	nodes []neo4j.Node
	rels  []neo4j.Relationship
}

func (p *MockPath) Nodes() []neo4j.Node {
	return p.nodes
}

func (p *MockPath) Relationships() []neo4j.Relationship {
	return p.rels
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
		&MockNode{id: 23},
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
		&MockNode{labels: []string{
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
		&MockNode{labels: []string{
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
		&MockNode{labels: []string{}},
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
		&MockNode{labels: []string{
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
		&MockNode{labels: []string{}},
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
		&MockNode{labels: []string{
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
		&MockNode{labels: []string{
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
		&MockNode{props: map[string]interface{}{
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
		&MockNode{props: map[string]interface{}{}},
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
		&MockNode{props: map[string]interface{}{
			`Something`: []string{
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
		&MockNode{props: map[string]interface{}{
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
		&MockRelationship{id: 452},
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
		&MockRelationship{startId: 452},
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
		&MockRelationship{endId: 452},
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
		&MockRelationship{relType: `hello world`},
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
		&MockRelationship{props: map[string]interface{}{
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
		&MockPath{
			nodes: []neo4j.Node{
				&MockNode{id: 234},
				&MockNode{id: 534},
				&MockNode{id: 9349},
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
		&MockPath{
			nodes: []neo4j.Node{
				&MockNode{id: 234},
				&MockNode{id: 534},
				&MockNode{id: 9349},
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
		&MockPath{
			nodes: []neo4j.Node{
				&MockNode{id: 234},
				&MockNode{id: 534},
				&MockNode{id: 9349},
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
		&MockPath{
			nodes: []neo4j.Node{
				&MockNode{id: 234},
				&MockNode{id: 534},
				&MockNode{id: 9349},
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
