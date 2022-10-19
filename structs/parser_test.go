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
