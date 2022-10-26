package structs

import (
	"reflect"
	"testing"
)

// MARK: Collection Helpers

func Test_Contains(t *testing.T) {
	collection := []string{"something", "else", "any", "thing"}

	key := "any"
	if !Contains(collection, key) {
		t.Errorf(`expected %v to be in collection`, key)
	}

	keys := []string{"test", "art", "think"}
	for _, key := range keys {
		ok := Contains(collection, key)
		if ok {
			t.Errorf(`expected %v to not be in collection`, key)
		}
	}
}

func Test_Map(t *testing.T) {
	type args struct {
		collection    []int
		transformFunc func(int, int) int
	}

	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "multiply",
			args: args{
				collection: []int{3, 5, 7},
				transformFunc: func(i int, n int) int {
					return n * n
				},
			},
			want: []int{9, 25, 49},
		},
		{
			name: "add",
			args: args{
				collection: []int{3, 5, 7},
				transformFunc: func(i int, n int) int {
					return n + 1
				},
			},
			want: []int{4, 6, 8},
		},
		{
			name: "subtract",
			args: args{
				collection: []int{3, 5, 7},
				transformFunc: func(i int, n int) int {
					return n - 1
				},
			},
			want: []int{2, 4, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.collection, tt.args.transformFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

// MARK: Reflection Helpers

func Test_PointerElement(t *testing.T) {
	var value *string = stringPointer("something")
	_, err := PointerElement(reflect.ValueOf(value))

	if err != nil {
		t.Errorf(`expected error to be nil, but got %v`, err)
	}
}

func Test_PointerElement_WhenNil(t *testing.T) {
	var value *string
	_, err := PointerElement(reflect.ValueOf(value))

	if err == nil {
		t.Errorf(`expected an error but got nil`)
	}
}
