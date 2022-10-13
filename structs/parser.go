package structs

import (
	"reflect"
	"strings"

	"github.com/google/uuid"
)

// Fetches all the fields of the given struct instance and returns a flattened list with all of its attributes.
//
// Params:
//
//	- filterTags
// when set, any fields not containing at least one of these tags
// will be ignored. An empty list allows all fields to be included.
//	- ignoredFields
// when set, any fields contained in this list will be ignored.
// Note that the name of the field should be the one defined in the struct.
//
// Usage:
//
// Imagine you have the struct:
// 	type Person struct {
//		Name string `json:"name"`
//		Emails []string `json:"emails"`
// 	}
//
// and create the following instance:
//	person := Person{Name: "Leonardo Ribeiro", Emails: []string{"leo@example.com", "lribeiro@example.org"}}
//
// `GetAttributes(person, ...)` would return a slice containing the following elements:
//
// 	name
//		-> Value: "Leonardo Ribeiro"
//	emails
//		-> Value: ["leo@example.com", "lribeiro@example.org"]
//		-> Children: [emails[0], emails[1]]
//	emails[0]
//		-> Value: "leo@example.com"
//		-> Parents: [emails]
//	emails[1]
//		-> Value: "lribeiro@example.org"
//		-> Parents: [emails]
//
// Each returned attribute will expose its underlying value as well as
// the definitions for its field type as found in the parent struct type.
func GetAttributes(entity reflect.Value, filterTags []string, ignoredFields ...string) (attributes []StructAttribute) {
	currentIndex := 0
	parents := []StructAttribute{}

	return getAttributes(entity, parents, filterTags, ignoredFields, currentIndex)
}

// Get the first value of the `json` tag.
//
// This is equivalent to calling:
//		GetTagValue(sf, "json")
func GetJSONTagValue(sf reflect.StructField) string {
	return GetTagValue(sf, "json")
}

// Get the first value of the given tag.
//
// Usage:
//
// Imagine you have the struct:
// 	type Person struct {
//		Name string `json:"name,omitempty" orm:"pk=name,noupdate,required,pk"`
//		Emails []string `json:"emails"`
// 	}
//
// You can obtain the `orm` tag the following way:
//	GetTagValue(name_sf, "orm") // -> "pk=name"
func GetTagValue(sf reflect.StructField, tagName string) string {
	name := sf.Name

	// Attribute name should come from json tag
	tag := strings.Split(sf.Tag.Get(tagName), ",")

	if len(tag) != 0 && tag[0] != "" {
		name = tag[0]
	}

	return name
}

// Get the full value of the given tag.
//
// Usage:
//
// Imagine you have the struct:
// 	type Person struct {
//		Name string `json:"name,omitempty" orm:"pk=name,noupdate,required,pk"`
//		Emails []string `json:"emails"`
// 	}
//
//	You can obtain the `orm` tag the following way:
//	GetTagValues(name_sf, "orm") // -> "pk=name,noupdate,required,pk"
func GetTagValues(sf reflect.StructField, tagName string) []string {
	r, exists := sf.Tag.Lookup(tagName)

	if exists {
		return strings.Split(r, ",")
	}

	return []string{}
}

// Fetches all the fields of the given struct.
func getAttributes(rv reflect.Value, parents []StructAttribute, filterTags, ignoredFields []string, currentIndex int) (attributes []StructAttribute) {
	if rv.Kind() == reflect.Pointer {
		rv, _ = pointerElement(rv)
	}

	if rv.Kind() != reflect.Struct {
		return attributes
	}

	for position := 0; position < rv.NumField(); position++ {
		// Concrete value type of the field at this position
		value := rv.Field(position)
		value, _ = pointerElement(value)

		// Struct field definition
		rsf := rv.Type().Field(position)

		sa := StructAttribute{
			Value:        value,
			Field:        rsf,
			Parents:      parents,
			ListPosition: currentIndex,
		}

		// Do not include an anonymous field at the top level.
		// Only include its inner fields.
		if sa.Field.Anonymous {
			anonValues := getAttributes(value, parents, filterTags, ignoredFields, currentIndex)
			sa.Children = append(sa.Children, anonValues...)
			attributes = append(attributes, anonValues...)
			continue
		}

		shouldBeIncluded := len(filterTags) == 0
		for _, tag := range filterTags {
			_, shouldBeIncluded = sa.Field.Tag.Lookup(tag)
		}

		if !shouldBeIncluded || contains(ignoredFields, rsf.Name) {
			continue
		}

		// Save field
		attributes = append(attributes, sa)

		// Check if the field needs further processing.
		switch value.Kind() {
		case reflect.Slice, reflect.Array:
			isListOfPrimitives := false
			newParents := append(parents, sa)

			if value.Len() > 0 {
				containsStructs := value.Index(0).Kind() == reflect.Struct

				// Google's UUID is a special case. Should not be considered a primitive type.
				isGoogleUUID := value.Type() == reflect.TypeOf(uuid.New())

				// Primitive types as in int, string, bool, etc
				isListOfPrimitives = !containsStructs && !isGoogleUUID
			}

			// Process each element in slice/array
			for l := 0; l < value.Len(); l++ {
				el := value.Index(l)

				if isListOfPrimitives {
					child := StructAttribute{
						Value:        el,
						Parents:      newParents,
						ListPosition: l,
						isPrimitive:  true,
					}

					// Copy information from parent StructField
					child.Field = reflect.StructField{
						Type:    el.Type(),
						Name:    child.FullName(),
						Tag:     sa.Field.Tag,
						PkgPath: sa.Field.PkgPath,
					}

					attributes[len(attributes)-1].Children = append(sa.Children, child)
					attributes = append(attributes, child)
					continue
				}

				nestedValues := getAttributes(el, newParents, filterTags, ignoredFields, l)
				if len(attributes) != 0 {
					attributes[len(attributes)-1].Children = append(sa.Children, nestedValues...)
				}

				attributes = append(attributes, nestedValues...)
			}
		}
	}

	return attributes
}
