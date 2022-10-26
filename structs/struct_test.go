package structs

import (
	"reflect"
	"testing"
)

func Test_StructAttribute_SkipsPastLastChild(t *testing.T) {
	type fields struct {
		Value        reflect.Value
		Field        reflect.StructField
		Parents      []StructAttribute
		Children     []StructAttribute
		ListPosition int
		isPrimitive  bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "test - 1",
			fields: fields{
				Children: []StructAttribute{},
			},
			want: 0,
		},
		{
			name: "test - 2",
			fields: fields{
				Children: []StructAttribute{{}},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := &StructAttribute{
				Value:        tt.fields.Value,
				Field:        tt.fields.Field,
				Parents:      tt.fields.Parents,
				Children:     tt.fields.Children,
				ListPosition: tt.fields.ListPosition,
				isPrimitive:  tt.fields.isPrimitive,
			}
			if got := sa.SkipsPastLastChild(); got != tt.want {
				t.Errorf("StructAttribute.SkipsPastLastChild() = %v, want %v", got, tt.want)
			}
		})
	}
}
