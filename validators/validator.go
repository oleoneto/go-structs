package validators

import (
	"errors"
	"net/mail"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/oleoneto/go-structs/structs"
	"golang.org/x/text/currency"
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

	// Use if field must have a valid currency code as value (only works on strings).
	//
	// If the field is a slice or an array of strings, the slice/array type itself
	// won't be validated, but each of its contained elements will be validated individually.
	//
	// Examples:
	//
	//	Currency    string    `validate:"currency"`
	//	Currencies  []string  `validate:"currency"`
	CURRENCY string = "currency"

	// Use if field must have a valid datetime value (only works on strings).
	//
	// If the field is a slice or an array of strings, the slice/array type itself
	// won't be validated, but each of its contained elements will be validated individually.
	//
	// Examples:
	//
	//	Date   string    `validate:"datetime"`
	//	Dates  []string  `validate:"datetime"`
	DATETIME string = "datetime"

	// Use if field must contain an email address (only works on strings).
	//
	// If the field is a slice or an array of strings, the slice/array type itself
	// won't be validated, but each of its contained elements will be validated individually.
	//
	// Examples:
	//
	//	Email  string   `format:"email"`
	//	Emails []string `format:"email"`
	EMAIL string = "email"

	// Use if string must have exactly 'eq' number of characters
	// or if integer must be exactly equal to this value.
	//
	// Examples:
	//
	//	Cards   []Card  `validate:"eq=2"`
	EQUAL string = "eq"

	// Use if field must be equal to one of the provided options.
	//
	// If the field is an array or a slice, each of its contained elements will be validated individually.
	//
	// Examples:
	//
	//	Role   string   `validate:"in=ADMIN|GUEST|SUPER USER"`
	//	Roles  []string `validate:"in=ADMIN|GUEST|SUPER USER"`
	//	Level  int      `validate:"in=1|5|20"`
	//	Levels []int    `validate:"in=1|5|20"`
	IN string = "in"

	// Use if string must have at least 'min' number of characters
	// or if integer must be greater than or equal to this value.
	//
	// Examples:
	//
	//	Name   string   `validate:"max=5"`
	//	Roles  []string `validate:"max=1"`
	//	Age    int      `validate:"max=18"`
	MAX string = "max"

	// Use if string must have at least 'min' number of characters
	// or if integer must be greater than or equal to this value.
	//
	// Examples:
	//
	//	Name   string   `validate:"min=5"`
	//	Roles  []string `validate:"min=1"`
	//	Age    int      `validate:"min=18"`
	MIN string = "min"

	// Use if field must contain a value that matches the specified regular expression.
	//
	// If the field is a slice or an array, the slice/array type itself
	// won't be validated, but each of its contained elements will be validated individually.
	//
	// Examples:
	//
	//	Name   string   `json:"name"   format:"regex([aA-zZ]+)"`
	//	Phone  string   `json:"phone"  format:"regex(\d{3}.\d{3}.\d{4})"`
	//	Phones []string `json:"phones" format:"regex(\d{3}.\d{3}.\d{4})"`
	REGEX string = "regex"

	// Use if field must contain a URL (only works on strings).
	//
	// If the field is a slice or an array of strings, the slice/array type itself
	// won't be validated, but each of its contained elements will be validated individually.
	//
	// Examples:
	//
	//	Website  string   `format:"url"`
	//	Websites []string `format:"url"`
	URL string = "url"

	// Use if field must contain a UUID-formated string (only works on strings).
	//
	// If the field is a slice or an array of strings, the slice/array type itself
	// won't be validated, but each of its contained elements will be validated individually.
	//
	// Examples:
	//
	//	Id        string   `format:"uuid"`
	//	Accounts  []string `format:"uuid"`
	UUID string = "uuid"
)

var Errors = map[string]string{
	"immutable": "IMMUTABLE_VALUE",
	"format":    "INVALID_FORMAT",
	"length":    "INVALID_LENGTH",
	"type":      "INVALID_TYPE",
	"value":     "INVALID_VALUE",
}

type ValidationOptions struct {
	Ignore    []string
	SkipRules []string
}

// Validates a struct and its attributes and returns a list of validation errors.
//
// Usage:
//
//	type Resource struct {
//		Id string `json:"id" validate:"uuid"`
//	}
//
//	r := Resource{Id: "abc"}
//	errs := ValidateAttribute(r) // -> {id: ["INVALID_FORMAT"]}
func Validate(model any, options ValidationOptions) map[string][]string {
	validations := make(map[string][]string)

	attributes := structs.GetAttributes(
		reflect.ValueOf(model),
		[]string{},
		options.Ignore...,
	)

	for pos := 0; pos < len(attributes); pos++ {
		attr := attributes[pos]
		errs := ValidateAttribute(attr, options)

		if len(errs) != 0 {
			validations[attr.FullName()] = errs

			switch attr.Value.Kind() {
			case reflect.Slice, reflect.Array:
				pos += attr.SkipsPastLastChild()
			}
		}
	}

	return validations
}

// Validates a struct attribute and returns a list of validation errors.
//
// Usage:
//
//	type Resource struct {
//		Id string `json:"name" validate:"uuid"`
//	}
//
//	r := Resource{Name: "abc"}
//	errs := ValidateAttribute(r["name"]) // -> ["INVALID_FORMAT"]
func ValidateAttribute(attribute structs.StructAttribute, options ValidationOptions) []string {
	validations := []string{}

	FORMAT_ERROR := []string{Errors["format"]}
	TYPE_ERROR := []string{Errors["type"]}
	VALUE_ERROR := []string{Errors["value"]}

	rules := structs.GetTagValues(attribute.Field, VALIDATION_TAG_KEYWORD)
	for _, validationRule := range rules {
		// The full validation ruleType. i.e min=20, required, nullable
		ruleType := validationRule

		// If the full validation rule contains a value, like min=20, this will be set to 20.
		var ruleValue string

		// This will split the rule and its value if one exits.
		// For example, `min=20` will become (min, 20)
		indexOfAsignment := strings.IndexByte(validationRule, '=')
		if indexOfAsignment != -1 {
			ruleType = validationRule[:indexOfAsignment]
			ruleValue = validationRule[indexOfAsignment+1:]
		}

		// Skip this rule
		if structs.Contains(options.SkipRules, ruleType) {
			continue
		}

		switch ruleType {
		case CURRENCY:
			f, err := structs.PointerElement(attribute.Value)
			if err != nil {
				return VALUE_ERROR
			}

			switch f.Kind() {
			case reflect.Array, reflect.Slice:
				// Assume children will be validated individually
				continue
			case reflect.String:
				if _, err := currency.ParseISO(f.String()); err != nil {
					return VALUE_ERROR
				}
			default:
				return TYPE_ERROR
			}
		case DATETIME:
			f, err := structs.PointerElement(attribute.Value)
			if err != nil {
				return FORMAT_ERROR
			}

			switch f.Kind() {
			case reflect.Array, reflect.Slice:
				// Assume that children will be validated individually
				continue
			case reflect.String:
				if f.Kind() == reflect.String {
					if _, err := time.Parse(time.RFC3339, f.String()); err != nil {
						return FORMAT_ERROR
					}

					continue
				}
			default:
				return TYPE_ERROR
			}
		case EMAIL:
			f, err := structs.PointerElement(attribute.Value)
			if err != nil {
				return FORMAT_ERROR
			}

			switch f.Kind() {
			case reflect.Array, reflect.Slice:
				// Assume that children will be validated individually
				continue
			case reflect.String:
				if _, err := mail.ParseAddress(f.String()); err != nil {
					return FORMAT_ERROR
				}
			default:
				return TYPE_ERROR
			}
		case EQUAL, MAX, MIN:
			length, err := parsedLengthAttribute(ruleValue)
			if err != nil {
				return VALUE_ERROR
			}

			f, err := structs.PointerElement(attribute.Value)
			if err != nil {
				return VALUE_ERROR
			}

			if !IsValidLength(f, length, ruleType) {
				var defaultError string

				switch f.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
					return VALUE_ERROR
				default:
					defaultError = Errors["length"]
				}

				return append(validations, defaultError)
			}
		case IN:
			f, err := structs.PointerElement(attribute.Value)
			if err != nil {
				return VALUE_ERROR
			}

			switch f.Kind() {
			case reflect.Array, reflect.Slice:
				// Assume that children will be validated individually
				continue
			default:
				acceptedValues := strings.Split(ruleValue, "|")
				if !IsIn(f, acceptedValues) {
					return VALUE_ERROR
				}
			}
		case UUID:
			f, err := structs.PointerElement(attribute.Value)
			if err != nil {
				return FORMAT_ERROR
			}

			switch f.Kind() {
			case reflect.Array, reflect.Slice:
				// Assume that children will be validated individually
				continue
			case reflect.String:
				if !IsUUID(f.String()) && len(validations) == 0 {
					return FORMAT_ERROR
				}
			default:
				return TYPE_ERROR
			}
		}
	}

	return validations
}

// Decodes and validates the provided payload.
//
// Usage:
//
//	type Resource struct {
//		Id   string `json:"id" validate:"uuid" jsonschema:"required"`
//		Name string `json:"name" validate:"min=3" jsonschema:"required"`
//	}
//
//	var r Resource
//	errs := ValidatePayload([]byte(`{"id": null}`), &r)
// /*
// {
// id: ["INVALID_TYPE"],
// name: ["REQUIRED_ATTRIBUTE_MISSING"]
// }
// */
func ValidatePayload(data []byte, model any, options ValidationOptions) map[string][]string {
	decoderErrors := structs.Decode(
		data,
		model,
		structs.DecoderOptions{
			Rules: []structs.SchemaValidationRule{
				structs.ADDITIONAL_PROPERTY,
				structs.INVALID_TYPE,
				structs.REQUIRED_ATTRIBUTE,
			},
		},
	)

	// NOTE: no need to go any further because the payload is invalid.
	if _, ok := decoderErrors["_"]; ok {
		return decoderErrors
	}

	validations := Validate(model, options)

	for k, v := range decoderErrors {
		validations[k] = v
	}

	return validations
}

// Returns `true` if value is one of the accepted values.
//
// Usage:
// 		IsIn(reflect.ValueOf("John"), []string{"Mario", "Luigi"}) // -> false
func IsIn(value reflect.Value, acceptedValues []string) bool {
	switch value.Kind() {
	case reflect.String:
		return structs.Contains(acceptedValues, value.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		found := false
		for _, v := range acceptedValues {
			vf, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return false
			}

			if vf == value.Int() {
				found = true
				break
			}
		}

		return found
	case reflect.Float32, reflect.Float64:
		found := false
		for _, v := range acceptedValues {
			vf, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return false
			}

			if vf == value.Float() {
				found = true
				break
			}
		}

		return found
	}

	return false
}

// Returns `true` of value is a UUID-formatted string.
//
// Usage:
//	IsUUID("something") // -> false
//	IsUUID("2bf99c42-4777-4796-9131-6cbc13d951c8") // -> true
func IsUUID(value string) bool {
	pattern := regexp.MustCompile(`^(?:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$`)
	return pattern.MatchString(value)
}

func IsValidLength(v reflect.Value, length float64, rule string) bool {
	var value float64 = -42

	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		value = float64(v.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = float64(v.Int())
	case reflect.Float32, reflect.Float64:
		value = v.Float()
	}

	if rule == MIN {
		return value >= length
	} else if rule == MAX {
		return value <= length
	}

	return value == length
}

// Returns `true` if the str is a valid value for the provided regular expression pattern.
//
// Usage:
//
//	PassesRegex(`\d+`, "23")       // -> true
//	PassesRegex(`\d+`, "leonardo") // -> false
func PassesRegex(pattern, str string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	return re.MatchString(str)
}

func parsedLengthAttribute(value string) (length float64, err error) {
	if value == "" {
		return length, errors.New("required length attribute")
	}

	return strconv.ParseFloat(value, 64)
}
