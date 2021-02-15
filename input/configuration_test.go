package input_test

import (
	"github.com/tlarsen7572/graphyx/input"
	"reflect"
	"testing"
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
