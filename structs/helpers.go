package structs

import (
	"errors"
	"reflect"
)

// Check if the element is contained within the given collection.
//
// Example:
//
//	Contains([]string{"hello", "world", "!"}, "world") // -> true
func Contains[T comparable](collection []T, element T) bool {
	for _, item := range collection {
		if item == element {
			return true
		}
	}

	return false
}

// Applies a `transformer` function to every element in a list.
//
// Usage:
//
//	Map([]int{3, 4}, func(index, n int) int { return n * n }) // -> [9, 16]
func Map[A any, B any](collection []A, transformFunc func(int, A) B) []B {
	result := make([]B, len(collection))

	for index, item := range collection {
		result[index] = transformFunc(index, item)
	}

	return result
}

// Filters a collection and returns only the elements
// that match the provided `inclusionTest`.
//
// Example:
//
// Filter and return all even numbers:
//
// 	Filter([]int{16, 9, 25}, func(i, n int) bool { return n%2 == 0 }) // -> [16]
func Filter[T any](collection []T, inclusionTest func(int, T) bool) []T {
	result := make([]T, 0)

	for index, item := range collection {
		if inclusionTest(index, item) {
			result = append(result, item)
		}
	}

	return result
}

// MARK: - Reflection Helpers

func PointerElement(rv reflect.Value) (reflect.Value, error) {
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
