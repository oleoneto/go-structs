package structs

import (
	"reflect"
	"testing"
)

type Identifiable struct {
	UUID string `json:"id"`
}

type Author struct {
	Id string `json:"id"`
}

type Person struct {
	Name         *string  `json:"name" db:"name"`
	Emails       []string `json:"emails" db:"emails"`
	IsActive     *bool
	PhoneNumbers []string `json:"phones"`
}

func Test_GetAllAttributes(t *testing.T) {
	type Article struct {
		Title   string   `json:"title"`
		Authors []Author `json:"authors"`
	}

	type Page struct {
		Identifiable
		PageID   string    `json:"page_id"`
		Articles []Article `json:"articles"`
	}

	type File struct {
		Dir   *string   `json:"dir"`
		Paths *[]string `json:"paths"`
	}

	type SSD struct {
		Owner Person `json:"owner"`
		Files []File `json:"files"`
	}

	var strvalue string

	type Expectation struct {
		Name       string
		Model      any
		Attributes []string
	}

	examples := []Expectation{
		{
			Name: "SSD - 1",
			Model: SSD{
				Owner: Person{},
			},
			Attributes: []string{
				"owner",
				"owner.name",
				"owner.emails",
				"owner.IsActive",
				"owner.phones",
				"files",
			},
		},
		{
			Name: "Person - 1",
			Model: Person{
				Name:         stringPointer("Leonardo Ribeiro"),
				Emails:       []string{"leo@example.com", "lribeiro@example.org"},
				PhoneNumbers: []string{},
				IsActive:     boolPointer(false),
			},
			Attributes: []string{
				"name",
				"emails",
				"emails[0]",
				"emails[1]",
				"IsActive",
				"phones",
			},
		},
		{
			Name: "Person - 2",
			Model: Person{
				Name:         stringPointer("Leonardo Ribeiro"),
				Emails:       []string{},
				PhoneNumbers: []string{"555.555.5555", "111.111.1111"},
				IsActive:     boolPointer(true),
			},
			Attributes: []string{
				"name",
				"emails",
				"IsActive",
				"phones",
				"phones[0]",
				"phones[1]",
			},
		},
		{
			Name: "Person - 3",
			Model: Person{
				Name:         stringPointer("Leonardo Ribeiro"),
				Emails:       []string{"leo@example.com"},
				PhoneNumbers: []string{"555.555.5555", "111.111.1111"},
			},
			Attributes: []string{
				"name",
				"emails",
				"emails[0]",
				"IsActive",
				"phones",
				"phones[0]",
				"phones[1]",
			},
		},
		{
			Name: "Page",
			Model: Page{
				PageID: "PAGE_ID",
				Articles: []Article{
					{
						Title: "Primeiro",
						Authors: []Author{
							{Id: "P1"},
							{Id: "P2"},
							{Id: "P3"},
						},
					},
					{
						Title: "Segundo",
						Authors: []Author{
							{Id: "ZZaa5599"},
							{Id: "Zq"},
						},
					},
				},
			},
			Attributes: []string{
				"id",
				"page_id",
				"articles",
				"articles[0].title",
				"articles[0].authors",
				"articles[0].authors[0].id",
				"articles[0].authors[1].id",
				"articles[0].authors[2].id",
				"articles[1].title",
				"articles[1].authors",
				"articles[1].authors[0].id",
				"articles[1].authors[1].id",
			},
		},
		{
			Name:       "Numeric",
			Model:      4,
			Attributes: []string{},
		},
		{
			Name:       "Pointer to String",
			Model:      stringPointer("something"),
			Attributes: []string{},
		},
		{
			Name:       "Pointer to nil",
			Model:      nil,
			Attributes: []string{},
		},
		{
			Name:       "Pointer to nil struct",
			Model:      &struct{}{},
			Attributes: []string{},
		},
		{
			Name:       "Struct literal - 1",
			Model:      struct{}{},
			Attributes: []string{},
		},
		{
			Name: "Struct literal - 2",
			Model: struct {
				Id    string `json:"id"`
				Notes []struct {
					Title string `json:"title"`
				} `json:"notes"`
			}{},
			Attributes: []string{
				"id",
				"notes",
			},
		},
		{
			Name: "Struct literal - 3",
			Model: struct {
				Id    string `json:"id"`
				Notes []struct {
					Title string `json:"title"`
				} `json:"notes"`
			}{
				Id: "uuid",
				Notes: []struct {
					Title string "json:\"title\""
				}{
					{
						Title: "Note 1",
					},
				},
			},
			Attributes: []string{
				"id",
				"notes",
				"notes[0].title",
			},
		},
		{
			Name: "Pointer to struct literal - 1",
			Model: &struct {
				Id    string `json:"id"`
				Notes []struct {
					Title string `json:"title"`
				} `json:"notes"`
			}{
				Id: "uuid",
				Notes: []struct {
					Title string "json:\"title\""
				}{
					{
						Title: "Note 1",
					},
				},
			},
			Attributes: []string{
				"id",
				"notes",
				"notes[0].title",
			},
		},
		{
			Name: "Files - 1",
			Model: File{
				Dir: nil,
				Paths: &[]string{
					"/home/users/someone/downloads",
					"/home/users/someone/music",
					"/home/users/someone/videos",
				},
			},
			Attributes: []string{
				"dir",
				"paths",
				"paths[0]",
				"paths[1]",
				"paths[2]",
			},
		},
		{
			Name: "Files - 2",
			Model: File{
				Dir:   nil,
				Paths: &[]string{},
			},
			Attributes: []string{
				"dir",
				"paths",
			},
		},
		{
			Name: "Files - 3",
			Model: File{
				Dir:   &strvalue,
				Paths: &[]string{},
			},
			Attributes: []string{
				"dir",
				"paths",
			},
		},
		{
			Name: "Files - 4",
			Model: File{
				Dir:   &strvalue,
				Paths: &[]string{strvalue},
			},
			Attributes: []string{
				"dir",
				"paths",
				"paths[0]",
			},
		},
	}

	for _, example := range examples {
		ignoredFields := []string{}

		t.Run(example.Name, func(t *testing.T) {
			values := GetAttributes(reflect.ValueOf(example.Model), []string{}, ignoredFields...)

			if len(values) != len(example.Attributes) {
				t.Errorf(`expected exactly %v values, but got %v`, len(example.Attributes), len(values))
				return
			}

			for i, field := range values {
				if field.FullName() != example.Attributes[i] {
					t.Errorf(`expected %v to be returned, but got %v`, example.Attributes[i], field.FullName())
					return
				}
			}
		})
	}
}

func Test_GetAttributesWithSpecificTags(t *testing.T) {
	type Expectation struct {
		Name       string
		Model      any
		Attributes []string
	}

	examples := []Expectation{
		{
			Name: "Person - 1",
			Model: Person{
				Name:         stringPointer("Leonardo Ribeiro"),
				Emails:       []string{"leo@example.com", "lribeiro@example.org"},
				PhoneNumbers: []string{},
				IsActive:     boolPointer(false),
			},
			Attributes: []string{
				"name",
				"emails",
				"emails[0]",
				"emails[1]",
			},
		},
		{
			Name: "Person - 2",
			Model: Person{
				Name:         stringPointer("Leonardo Ribeiro"),
				Emails:       []string{},
				PhoneNumbers: []string{"555.555.5555", "111.111.1111"},
				IsActive:     boolPointer(true),
			},
			Attributes: []string{
				"name",
				"emails",
			},
		},
		{
			Name: "Person - 3",
			Model: Person{
				Name:         stringPointer("Leonardo Ribeiro"),
				Emails:       []string{"leo@example.com"},
				PhoneNumbers: []string{"555.555.5555", "111.111.1111"},
			},
			Attributes: []string{
				"name",
				"emails",
				"emails[0]",
			},
		},
	}

	for _, example := range examples {
		ignoredFields := []string{}

		t.Run(example.Name, func(t *testing.T) {
			values := GetAttributes(reflect.ValueOf(example.Model), []string{"db"}, ignoredFields...)

			if len(values) != len(example.Attributes) {
				t.Errorf(`expected exactly %v values, but got %v`, len(example.Attributes), len(values))
				return
			}

			for i, field := range values {
				if field.FullName() != example.Attributes[i] {
					t.Errorf(`expected %v to be returned, but got %v`, example.Attributes[i], field.FullName())
					return
				}
			}
		})
	}
}

func Test_GetTagValues(t *testing.T) {
	var field reflect.StructField = reflect.StructField{
		Tag: `json:"id,omitempty" db:"_id"`,
	}

	jExpectation := "id"
	jValues := GetTagValues(field, "json")
	if jValues[0] != jExpectation {
		t.Errorf(`expected values %v, but got %v`, jExpectation, jValues)
	}

	dExpectation := "_id"
	dValues := GetTagValues(field, "db")
	if dValues[0] != dExpectation {
		t.Errorf(`expected values %v, but got %v`, dExpectation, dValues)
	}

	eValues := GetTagValues(field, "quadrado")
	if len(eValues) != 0 {
		t.Errorf(`expected values to be empty, but got %v`, eValues)
	}
}

func Test_GetTags(t *testing.T) {
	var field reflect.StructField = reflect.StructField{
		Tag: `json:"id,omitempty" db:"_id"`,
	}

	tags := GetTags(field)

	if len(tags["json"]) != 2 {
		t.Errorf(`expected different values for "json" tag, but got %v`, tags["json"])
	}

	if len(tags["db"]) != 1 {
		t.Errorf(`expected different value "db" tag, but got %v`, tags["db"])
	}
}

func Test_TagContainsValues(t *testing.T) {
	type Expectation struct {
		Name   string
		Field  reflect.StructField
		Tag    string
		Values []string
		Result bool
	}

	tests := []Expectation{
		{
			Name:   "Nullable field - 1",
			Field:  reflect.StructField{Name: "id", Tag: `json:"id,omitempty"`},
			Tag:    "json",
			Values: []string{"omitempty"},
			Result: true,
		},
		{
			Name:   "Nullable field - 2",
			Field:  reflect.StructField{Name: "id", Tag: `json:"id,omitempty"`},
			Tag:    "json",
			Values: []string{"arroz"},
			Result: false,
		},
		{
			Name:   "Email - 1",
			Field:  reflect.StructField{Name: "email", Tag: `validate:"email,contains(@dock.tech)"`},
			Tag:    "json",
			Values: []string{"omitempty"},
			Result: false,
		},
		{
			Name:   "Email - 2",
			Field:  reflect.StructField{Name: "email", Tag: `validate:"email,contains(@dock.tech)"`},
			Tag:    "validate",
			Values: []string{"omitempty"},
			Result: false,
		},
		{
			Name:   "Email - 3",
			Field:  reflect.StructField{Name: "email", Tag: `validate:"email,contains(@dock.tech)"`},
			Tag:    "validate",
			Values: []string{"email"},
			Result: true,
		},
		{
			Name:   "Email - 4",
			Field:  reflect.StructField{Name: "email", Tag: `validate:"email,matches(@dock.tech)"`},
			Tag:    "validate",
			Values: []string{"@dock.tech"},
			Result: false,
		},
		{
			Name:   "Email - 4.1",
			Field:  reflect.StructField{Name: "email", Tag: `validate:"email,matches(@dock.tech)"`},
			Tag:    "validate",
			Values: []string{"@dock.tech"},
			Result: false,
		},
		{
			Name:   "Email - 5",
			Field:  reflect.StructField{Name: "email", Tag: `validate:"email,matches(@dock.tech)"`},
			Tag:    "validate",
			Values: []string{"matches(@dock.tech)"},
			Result: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			res := TagConstainsValues(test.Field, test.Tag, test.Values)

			if res != test.Result {
				t.Errorf(`expected result to be %v, but got %v`, test.Result, res)
			}
		})
	}
}

func Test_MatchingFields(t *testing.T) {
	type Expectation struct {
		Name   string
		Model  any
		Tag    string
		Values []string
		Fields []string
	}

	type Person struct {
		Name           string   `json:"name,omitempty" orm:"pk=name,noupdate" check:"uuid"`
		PrimaryEmail   string   `json:"email1" check:"email,primary"`
		SecondaryEmail []string `json:"email2" check:"email,backup"`
	}

	tests := []Expectation{
		{
			Name:   "person - 1",
			Model:  Person{},
			Tag:    "json",
			Values: []string{"omitempty"},
			Fields: []string{"name"},
		},
		{
			Name:   "person - 2",
			Model:  Person{},
			Tag:    "check",
			Values: []string{"email"},
			Fields: []string{"email1", "email2"},
		},
		{
			Name:   "person - 3",
			Model:  Person{},
			Tag:    "check",
			Values: []string{"primary"},
			Fields: []string{"email1"},
		},
		{
			Name:   "person - 4",
			Model:  Person{},
			Tag:    "check",
			Values: []string{"backup"},
			Fields: []string{"email2"},
		},
		{
			Name:   "person - 5",
			Model:  Person{},
			Tag:    "json",
			Values: []string{"email"},
			Fields: []string{},
		},
		{
			Name:   "person - 5.1",
			Model:  Person{},
			Tag:    "json",
			Values: []string{"email"},
			Fields: []string{},
		},
		{
			Name:   "person - 6",
			Model:  Person{},
			Tag:    "json",
			Values: []string{"email1"},
			Fields: []string{"email1"},
		},
		{
			Name:   "person - 7",
			Model:  Person{},
			Tag:    "orm",
			Values: []string{"email"},
			Fields: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			fields := MatchingFields(test.Model, test.Tag, test.Values)

			if len(fields) != len(test.Fields) {
				t.Errorf(`expected %d fields, but got %d`, len(test.Fields), len(fields))
				return
			}

			for index, field := range fields {
				if field != test.Fields[index] {
					t.Errorf(`expected field to be %v, but got %v instead`, test.Fields[index], field)
				}
			}
		})
	}
}
