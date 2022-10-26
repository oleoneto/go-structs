package structs

import (
	"fmt"
	"reflect"
	"strings"
)

// This type abstracts both `reflect.Value` and `reflect.StructField` types.
type StructAttribute struct {
	Value        reflect.Value
	Field        reflect.StructField
	Parents      []StructAttribute
	Children     []StructAttribute
	ListPosition int
	isPrimitive  bool
}

type StructAttributes []StructAttribute

// Returns the name of the field properly scoped under its parents.
//
// Usage:
//
// Imagine you have the following StructAttribute:
//	sa := StructAttribute{
//		Parents: []StructAttribute{parentA, listB}
//		Field: reflect.StructField{
// 			Name:    "attribute_name",
// 			...
// 		}
//	}
//
//	sa.FullName() // -> "parentA.listB[i].attribute_name"
func (sa *StructAttribute) FullName() (name string) {
	if len(sa.Parents) == 0 {
		return GetJSONTagValue(sa.Field)
	}

	scope := sa.Parents[len(sa.Parents)-1].FullName()

	// Adds the array notation to the slice/array field
	if sa.ListPosition >= 0 {
		scope = strings.Join([]string{scope, fmt.Sprint("[", sa.ListPosition, "]")}, "")
	}

	if sa.isPrimitive {
		return scope
	}

	fullName := strings.Join([]string{scope, GetJSONTagValue(sa.Field)}, ".")

	// Ensures field name is never prefixed by a dot (.)
	return strings.TrimSuffix(strings.TrimPrefix(fullName, "."), ".")
}

func (sa *StructAttribute) SkipsPastLastChild() int {
	if len(sa.Children) == 0 {
		return 0
	}

	n := 1 + len(sa.Children)
	for _, child := range sa.Children {
		n += 1 + child.SkipsPastLastChild()
	}

	return n
}
