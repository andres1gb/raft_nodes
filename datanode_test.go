package raftnodes

import (
	"reflect"
	"strconv"
	"testing"
)

type TestData struct {
	name string
	age  int
}

func (d *TestData) String() string {
	return d.name + " has " + strconv.Itoa(d.age) + " years"
}

func Test_newDataNode(t *testing.T) {
	tests := []struct {
		name string
		want *dataNode[TestData]
	}{
		{
			name: "Create data node",
			want: newDataNode[TestData](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDataNode[TestData](); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDataNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
