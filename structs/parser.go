package structs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

const (
	// The literal name of the validation tag as it'll appear in the struct.
	//
	// Example:
	//
	//	type Resource struct {
	//		Name string `json:"name" validate:"min=6"`
	//	}
	VALIDATION_TAG_KEYWORD string = "validate"
)

var (
	// Tag attributes that should be excluded
	NON_INHERITABLE_TAG_ATTRIBUTES = []string{"max", "min"}
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
// You can obtain the `orm` tag the following way:
//	GetTagValues(name_sf, "orm") // -> "pk=name,noupdate,required,pk"
func GetTagValues(sf reflect.StructField, tagName string) []string {
	r, exists := sf.Tag.Lookup(tagName)

	if exists {
		return strings.Split(r, ",")
	}

	return []string{}
}

// Get each of the attributes of the given tag.
func GetTag(sf reflect.StructField, tagName string) map[string]string {
	values := make(map[string]string, 0)

	if r, exists := sf.Tag.Lookup(tagName); exists {
		rule := strings.Split(r, ",")

		for _, rl := range rule {
			t := strings.SplitN(rl, "=", 2)

			if len(t) == 1 {
				values[t[0]] = ""
				continue
			}

			values[t[0]] = t[1]
		}
	}

	return values
}

// Get all the tags of the given struct field.
//
// Usage:
//
// Imagine you have the struct:
//	type Person struct {
//		Name string `json:"name,omitempty" orm:"pk=name"`
//		Emails []string `json:"emails"`
//	}
//
// You can get all the tags set on the name field:
// 	GetTags(name_sf) // -> {json: [name, omitempty], orm: [pk=name]}
func GetTags(sf reflect.StructField) map[string][]string {
	tags := make(map[string][]string)

	for _, tag := range strings.Split(string(sf.Tag), " ") {
		t := strings.Split(tag, ":")
		name := t[0]
		values := t[1]

		tags[name] = strings.Split(values, ",")
	}

	return tags
}

// Returns whether or not a struct field contains the provided values in the specified tag.
//
// Usage:
//
// Imagine you have the struct:
// 	type Person struct {
//		Name           string   `json:"name,omitempty" orm:"pk=name,noupdate,required,pk" validate:"uuid"`
//		PrimaryEmail   string   `json:"email1" validate:"email"`
//		SecondaryEmail []string `json:"email2" validate:"email"`
// 	}
//
// Does the field `Name` have the value `email` in the `validate` tag?
//	TagConstainsValues(name_sf, "validate", []string{"email"}) // -> false
//
// Does the field `PrimaryEmail` has the value `email` in the `validate` tag?
//	TagConstainsValues(primary_email_sf, "validate", []string{"email"}) // -> true
//
// IMPORTANT:
//
// Values will only match if they are the same. If you pass only a substring, the method will return false.
//
// For example:
//	TagConstainsValues(primary_email_sf, "json", []string{"email"}) // -> false
func TagConstainsValues(field reflect.StructField, tag string, values []string) bool {
	if tag, ok := field.Tag.Lookup(tag); ok {
		for _, value := range values {
			if Contains(strings.Split(tag, ","), value) {
				return true
			}
		}
	}

	return false
}

// Get a list of all the struct fields that contain the provided values in the specified tag.
//
// Usage:
//
// Imagine you have the struct:
//	type Person struct {
//		Name           string   `json:"name,omitempty" orm:"pk=name,noupdate" validate:"uuid"`
//		PrimaryEmail   string   `json:"email1" validate:"email"`
//		SecondaryEmail []string `json:"email2" validate:"email"`
//	}
//
// You can get all the fields that include the value `email` in the `validate` tag:
//	MatchingFields(Person{}, "validate", []string{"email"}) // -> [email1, email2]
func MatchingFields(v any, tag string, requiredKeywords []string) (result []string) {
	rv := reflect.ValueOf(v)
	parents := []string{}
	return matchingFields(rv, parents, tag, requiredKeywords)
}

func SetValuesFromMap(entity any, values map[string]any) {
	rv := reflect.ValueOf(entity)
	attrs := GetAttributes(rv, []string{})

	for _, attr := range attrs {
		if v, ok := values[attr.FullName()]; ok {
			if sf := rv.Elem().FieldByName(attr.Field.Name); sf.CanSet() {
				value := reflect.ValueOf(v)

				switch sf.Type().Kind() {
				case reflect.Array, reflect.Slice:
					if value.Kind() != sf.Type().Kind() {
						delete(values, attr.FullName())
					}
				case reflect.Pointer:
					if value.Kind() != sf.Type().Elem().Kind() {
						delete(values, attr.FullName())
					}
				case reflect.Struct:
				}
			}
		}
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(values)
	json.NewDecoder(buf).Decode(entity)
}

func SetValuesFromBytes(entity any, data []byte) {
	values := map[string]any{}
	_ = json.Unmarshal(data, &values)
	SetValuesFromMap(entity, values)
}

// -------------------------------------------------------
// -------------------------------------------------------
// -------------------------------------------------------

// Fetches all the fields of the given struct.
func getAttributes(rv reflect.Value, parents []StructAttribute, filterTags, ignoredFields []string, currentIndex int) (attributes []StructAttribute) {
	if rv.Kind() == reflect.Pointer {
		rv, _ = PointerElement(rv)
	}

	if rv.Kind() != reflect.Struct {
		return attributes
	}

	for position := 0; position < rv.NumField(); position++ {
		// Concrete value type of the field at this position
		value := rv.Field(position)
		value, _ = PointerElement(value)

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

		if !shouldBeIncluded || Contains(ignoredFields, rsf.Name) {
			continue
		}

		// Save field
		attributes = append(attributes, sa)

		// Check if the field needs further processing.
		switch value.Kind() {
		case reflect.Struct:
			nestedAttributes := getAttributes(value, append(parents, sa), filterTags, ignoredFields, -1)
			attributes = append(attributes, nestedAttributes...)
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

					// Exclude some predefined validation tag attributes
					childTag := RemoveValuesFromTag(VALIDATION_TAG_KEYWORD, NON_INHERITABLE_TAG_ATTRIBUTES, sa.Field)
					child.Field.Tag = reflect.StructTag(childTag)

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

func RemoveValuesFromTag(tag string, removeList []string, field reflect.StructField) string {
	result := string(field.Tag)

	t := GetTag(field, tag)
	for _, i := range removeList {
		if v, ok := t[i]; ok {
			pattern := regexp.MustCompile(fmt.Sprintf(`,?%v=?%v?,?`, i, v))
			result = pattern.ReplaceAllString(result, "")
		}
	}

	return result
}

func matchingFields(rv reflect.Value, parents []string, tag string, requiredKeywords []string) (fields []string) {
	if rv.Kind() == reflect.Pointer {
		rv, _ = PointerElement(rv)
	}

	if rv.Kind() != reflect.Struct {
		return fields
	}

	for position := 0; position < rv.NumField(); position++ {
		f := rv.Type().Field(position)
		value := rv.Field(position)

		prefix := strings.Join(parents, ".")
		fieldName := strings.TrimPrefix(strings.Join([]string{prefix, GetTagValue(f, "json")}, "."), ".")
		if TagConstainsValues(f, tag, requiredKeywords) {
			fields = append(fields, fieldName)
		}

		switch value.Kind() {
		case reflect.Array, reflect.Slice:
			newParents := append(parents, fieldName)

			t := reflect.New(value.Type().Elem())
			fields = append(fields, matchingFields(t, newParents, tag, requiredKeywords)...)
		}
	}

	return fields
}
