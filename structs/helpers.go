package structs

import (
	"errors"
	"reflect"
)

// Check if the element is contained within the given collection.
//
// Example:
//
//	contains([]string{"hello", "world", "!"}, "world") // -> true
func contains[T comparable](collection []T, element T) bool {
	for _, item := range collection {
		if item == element {
			return true
		}
	}

	return false
}

// MARK: - Reflection Helpers

func pointerElement(rv reflect.Value) (reflect.Value, error) {
	el := rv

	for el.Kind() == reflect.Pointer {
		if el.IsNil() {
			return el, errors.New("nil pointer")
		}

		el = el.Elem()
	}

	return el, nil
}

func stringPointer(v string) *string {
	return &v
}

func boolPointer(v bool) *bool {
	return &v
}
