package structs

import (
	"reflect"
	"testing"
)

// MARK: Collection Helpers

func Test_Contains(t *testing.T) {
	collection := []string{"something", "else", "any", "thing"}

	key := "any"
	if !contains(collection, key) {
		t.Errorf(`expected %v to be in collection`, key)
	}

	keys := []string{"test", "art", "think"}
	for _, key := range keys {
		ok := contains(collection, key)
		if ok {
			t.Errorf(`expected %v to not be in collection`, key)
		}
	}
}

// MARK: Reflection Helpers

func Test_PointerElement(t *testing.T) {
	var value *string = stringPointer("something")
	_, err := pointerElement(reflect.ValueOf(value))

	if err != nil {
		t.Errorf(`expected error to be nil, but got %v`, err)
	}
}
func Test_PointerElement_WhenNil(t *testing.T) {
	var value *string
	_, err := pointerElement(reflect.ValueOf(value))

	if err == nil {
		t.Errorf(`expected an error but got nil`)
	}
}
