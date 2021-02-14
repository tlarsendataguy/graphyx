package graphyx_test

import (
	"github.com/tlarsen7572/graphyx"
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
	</Fields>
</Configuration>`

	decoded, err := graphyx.DecodeConfig(config)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	expected := graphyx.Configuration{
		ConnStr:  `bolt://localhost:7687`,
		Username: `user`,
		Password: `password`,
		Query:    `MATCH p=()-[r:ACTED_IN]->() RETURN p`,
		Fields: []graphyx.Field{
			{
				Name:     `Field1`,
				DataType: `Integer`,
				Path: []graphyx.Element{
					{Key: `p`, DataType: `Path`},
					{Key: `Nodes`, DataType: `List:Node`},
					{Key: `0`, DataType: `Node`},
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
