package validators

import (
	"reflect"
	"testing"
)

type Identifiable struct {
	UUID string `json:"id" validate:"uuid"`
}

type Contact struct {
	IsActive     bool     `json:"is_active"`
	Emails       []string `json:"emails" db:"emails" validate:"email,min=1"`
	PhoneNumbers []string `json:"phones"`
}

type Person struct {
	Identifiable
	Name    string  `json:"name" db:"name" validate:"min=2,max=8"`
	Contact Contact `json:"contact"`
}

type Resource struct {
	Currency         string   `json:"currency" validate:"currency"`
	Price            float64  `json:"price" validate:"min=5,max=30"`
	Group            int      `json:"group" validate:"in=1|3|5|7"`
	Type             string   `json:"type" validate:"in=USED|NEW"`
	Code             string   `json:"code" validate:"eq=5"`
	Quantity         int      `json:"qty" validate:"min=1,max=3"`
	Rating           float64  `json:"rating" validate:"in=0|1|2|4|5"`
	Related          []string `json:"related" validate:"uuid"`
	PublishedAt      string   `json:"published_at" validate:"datetime"`
	PriceAvailableIn []string `json:"price_available_in" validate:"currency"`
}

func Test_Validate(t *testing.T) {
	type S1 struct {
		C []int `json:"c" validate:"currency"`
		D []int `json:"d" validate:"datetime"`
		E []int `json:"e" validate:"email"`
		I []int `json:"i" validate:"in=1|2|3"`
		U []int `json:"u" validate:"uuid"`
	}

	type S2 struct {
		C *string `json:"c" validate:"currency"`
		F string  `json:"f" validate:"min="`
		G *string `json:"g" validate:"uuid"`
		H *string `json:"h" validate:"in=1|2"`
		K *string `json:"k" validate:"email"`
		L *string `json:"l" validate:"datetime"`
		M *int    `json:"m" validate:"min=3"`
	}

	tests := []struct {
		name    string
		model   any
		options ValidationOptions
		want    map[string][]string
	}{
		{
			name:    "person - 1",
			model:   Person{},
			options: ValidationOptions{},
			want: map[string][]string{
				"id":             {"INVALID_FORMAT"},
				"name":           {"INVALID_LENGTH"},
				"contact.emails": {"INVALID_LENGTH"},
			},
		},
		{
			name: "person - 2",
			model: Person{
				Identifiable: Identifiable{UUID: ""},
			},
			options: ValidationOptions{},
			want: map[string][]string{
				"id":             {"INVALID_FORMAT"},
				"name":           {"INVALID_LENGTH"},
				"contact.emails": {"INVALID_LENGTH"},
			},
		},
		{
			name: "person - 3",
			model: Person{
				Identifiable: Identifiable{UUID: "2b852002-f19d-11ec-8ea0-0242ac120002"},
			},
			options: ValidationOptions{},
			want: map[string][]string{
				"name":           {"INVALID_LENGTH"},
				"contact.emails": {"INVALID_LENGTH"},
			},
		},
		{
			name: "person - 4",
			model: Person{
				Identifiable: Identifiable{UUID: "2b852002-f19d-11ec-8ea0-0242ac120002"},
				Name:         "Leonardo",
			},
			options: ValidationOptions{},
			want: map[string][]string{
				"contact.emails": {"INVALID_LENGTH"},
			},
		},
		{
			name: "person - 5",
			model: Person{
				Identifiable: Identifiable{UUID: "2b852002-f19d-11ec-8ea0-0242ac120002"},
				Name:         "Leonardo",
				Contact:      Contact{},
			},
			options: ValidationOptions{},
			want: map[string][]string{
				"contact.emails": {"INVALID_LENGTH"},
			},
		},
		{
			name: "person - 6",
			model: Person{
				Identifiable: Identifiable{UUID: "2b852002-f19d-11ec-8ea0-0242ac120002"},
				Name:         "Leonardo",
				Contact:      Contact{Emails: []string{"email"}},
			},
			options: ValidationOptions{},
			want: map[string][]string{
				"contact.emails[0]": {"INVALID_FORMAT"},
			},
		},
		{
			name: "person - 7",
			model: Person{
				Identifiable: Identifiable{UUID: "2b852002-f19d-11ec-8ea0-0242ac120002"},
				Name:         "Leonardo",
				Contact:      Contact{Emails: []string{"email@example.com"}},
			},
			options: ValidationOptions{},
			want:    map[string][]string{},
		},
		{
			name: "person - 8",
			model: Person{
				Identifiable: Identifiable{UUID: "2b852002-f19d-11ec-8ea0-0242ac120002"},
				Name:         "Leonardo Ribeiro",
				Contact:      Contact{Emails: []string{"email@example.com"}},
			},
			options: ValidationOptions{},
			want: map[string][]string{
				"name": {"INVALID_LENGTH"},
			},
		},
		{
			name: "resource - 1",
			model: Resource{
				Currency:    "BRL",
				Price:       14,
				Group:       7,
				Type:        "NEW",
				Code:        "ABC12",
				Quantity:    2,
				Rating:      5,
				PublishedAt: "2020-01-01T00:00:00+01:00",
			},
			options: ValidationOptions{},
			want:    map[string][]string{},
		},
		{
			name: "resource - 2",
			model: Resource{
				Currency:         "BRL",
				Price:            42,
				Group:            42,
				Type:             "AWESOME",
				Code:             "ABC12",
				Quantity:         2,
				Rating:           42,
				Related:          []string{"abc"},
				PriceAvailableIn: []string{"AUD", "EURO"},
			},
			options: ValidationOptions{},
			want: map[string][]string{
				"price":                 {"INVALID_VALUE"},
				"group":                 {"INVALID_VALUE"},
				"type":                  {"INVALID_VALUE"},
				"rating":                {"INVALID_VALUE"},
				"related[0]":            {"INVALID_FORMAT"},
				"published_at":          {"INVALID_FORMAT"},
				"price_available_in[1]": {"INVALID_VALUE"},
			},
		},
		{
			name: "resource - 3",
			model: Resource{
				Currency:         "BRL",
				Price:            42,
				Group:            42,
				Type:             "AWESOME",
				Code:             "ABC12",
				Quantity:         2,
				Rating:           42,
				Related:          []string{"abc"},
				PriceAvailableIn: []string{"AUD", "EURO"},
				PublishedAt:      "2022-01-01",
			},
			options: ValidationOptions{},
			want: map[string][]string{
				"price":                 {"INVALID_VALUE"},
				"group":                 {"INVALID_VALUE"},
				"type":                  {"INVALID_VALUE"},
				"rating":                {"INVALID_VALUE"},
				"related[0]":            {"INVALID_FORMAT"},
				"published_at":          {"INVALID_FORMAT"},
				"price_available_in[1]": {"INVALID_VALUE"},
			},
		},
		{
			name:    "s1 - 1",
			model:   S1{C: []int{42}, D: []int{42}, E: []int{42}, I: []int{42}, U: []int{42}},
			options: ValidationOptions{},
			want: map[string][]string{
				"c[0]": {"INVALID_TYPE"},
				"d[0]": {"INVALID_TYPE"},
				"e[0]": {"INVALID_TYPE"},
				"i[0]": {"INVALID_VALUE"},
				"u[0]": {"INVALID_TYPE"},
			},
		},
		{
			name:  "s2 - 1",
			model: S2{},
			options: ValidationOptions{
				SkipRules: []string{"currency", "uuid", "datetime", "in"},
			},
			want: map[string][]string{
				"f": {"INVALID_VALUE"},
				"k": {"INVALID_FORMAT"},
				"m": {"INVALID_VALUE"},
			},
		},
		{
			name:    "s2 - 2",
			model:   S2{},
			options: ValidationOptions{},
			want: map[string][]string{
				"c": {"INVALID_VALUE"},
				"f": {"INVALID_VALUE"},
				"g": {"INVALID_FORMAT"},
				"h": {"INVALID_VALUE"},
				"k": {"INVALID_FORMAT"},
				"l": {"INVALID_FORMAT"},
				"m": {"INVALID_VALUE"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Validate(tt.model, tt.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ValidatePayload(t *testing.T) {
	type Person struct {
		UUID    string  `json:"id" validate:"uuid" jsonschema:"required"`
		Name    string  `json:"name" db:"name" validate:"min=2,max=8"`
		Contact Contact `json:"contact"`
	}

	type args struct {
		data    []byte
		model   any
		options ValidationOptions
	}

	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "person - 1",
			args: args{
				data:    []byte(`{"name": "", "contact": {"emails": []}}`),
				model:   &Person{},
				options: ValidationOptions{},
			},
			want: map[string][]string{
				"id":             {"REQUIRED_ATTRIBUTE_MISSING"},
				"name":           {"INVALID_LENGTH"},
				"contact.emails": {"INVALID_LENGTH"},
			},
		},
		{
			name: "person - 2",
			args: args{
				data:    []byte(`{"id": "1108129d-1d98-4a21-837a-ae6319f64c73", "name": 1, "contact": {"emails": ["}}`),
				model:   &Person{},
				options: ValidationOptions{},
			},

			want: map[string][]string{
				"_": {"INVALID_PAYLOAD"},
			},
		},
		{
			name: "person - 3",
			args: args{
				data:    []byte(`{"id": "2b852002-f19d-11ec-8ea0-0242ac120002", "name": 1, "contact": {"emails": ["leo", "leo@example.org"]}}`),
				model:   &Person{},
				options: ValidationOptions{},
			},
			want: map[string][]string{
				"name":              {"INVALID_TYPE"},
				"contact.emails[0]": {"INVALID_FORMAT"},
			},
		},
		{
			name: "person - 4",
			args: args{
				data:    []byte(`{"id": "2b852002-f19d-11ec-8ea0-0242ac120002", "name": "Leonardo", "contact": {"emails": ["leo@example.org"]}}`),
				model:   &Person{},
				options: ValidationOptions{},
			},
			want: map[string][]string{},
		},
		{
			name: "resource - 1",
			args: args{
				model:   &Resource{},
				data:    []byte(`{"currency": "BRL", "price": 14, "group": 7, "type": "NEW", "code": "ABC12", "qty": 2, "rating": 5, "related": ["123", "145"], "published_at": "2020-01-01T00:00:00+01:00", "id": "some-id"}`),
				options: ValidationOptions{},
			},
			want: map[string][]string{
				"id":         {"ADDITIONAL_PROPERTY"},
				"related[0]": {"INVALID_FORMAT"},
				"related[1]": {"INVALID_FORMAT"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePayload(tt.args.data, tt.args.model, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidatePayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

// MARK: Rules

func Test_IsIn(t *testing.T) {
	type args struct {
		value          reflect.Value
		acceptedValues []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "numeric - 1",
			args: args{
				value:          reflect.ValueOf(2),
				acceptedValues: []string{"2", "3", "5", "7"},
			},
			want: true,
		},
		{
			name: "numeric - 2",
			args: args{
				value:          reflect.ValueOf(7),
				acceptedValues: []string{"2", "3", "5", "7"},
			},
			want: true,
		},
		{
			name: "numeric - 3",
			args: args{
				value:          reflect.ValueOf(0),
				acceptedValues: []string{"2", "3", "5", "7"},
			},
			want: false,
		},
		{
			name: "numeric - 4",
			args: args{
				value:          reflect.ValueOf(0),
				acceptedValues: []string{"2a", "B"},
			},
			want: false,
		},
		{
			name: "numeric - 5",
			args: args{
				value:          reflect.ValueOf(10.2),
				acceptedValues: []string{"2a", "B"},
			},
			want: false,
		},
		{
			name: "numeric - 6",
			args: args{
				value:          reflect.ValueOf([]int{}),
				acceptedValues: []string{"2", "B"},
			},
			want: false,
		},
		{
			name: "string - 1",
			args: args{
				value:          reflect.ValueOf("GUEST"),
				acceptedValues: []string{"ADMIN", "GUEST"},
			},
			want: true,
		},
		{
			name: "string - 2",
			args: args{
				value:          reflect.ValueOf("WRITER"),
				acceptedValues: []string{"ADMIN", "GUEST"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIn(tt.args.value, tt.args.acceptedValues); got != tt.want {
				t.Errorf("IsIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsUUID(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{
			name: "test - 1",
			arg:  "",
			want: false,
		},
		{
			name: "test - 2",
			arg:  "someone-is-cool",
			want: false,
		},
		{
			name: "test - 3",
			arg:  " ",
			want: false,
		},
		{
			name: "test - 4",
			arg:  "00000000-0000-0000-0000-000000000000",
			want: true,
		},
		{
			name: "test - 5",
			arg:  "21f2fa1d-c662-4669-abba-095a2416f5b9",
			want: true,
		},
		{
			name: "test - 6",
			arg:  "097e3981-55a2-4f83-a5dd-405e49a80314",
			want: true,
		},
		{
			name: "test - 7",
			arg:  "2b852002-f19d-11ec-8ea0-0242ac120002",
			want: true,
		},
		{
			name: "test - 8",
			arg:  "2g852002-f19d-11ec-8ea0-0242ac120002",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUUID(tt.arg); got != tt.want {
				t.Errorf("IsUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_PassesRegex(t *testing.T) {
	type args struct {
		pattern string
		str     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "numeric - 0",
			want: false,
			args: args{
				pattern: `\d+`,
				str:     "leonardo",
			},
		},
		{
			name: "numeric - 1",
			want: true,
			args: args{
				pattern: `\d+`,
				str:     "299",
			},
		},
		{
			name: "numeric - 2",
			want: true,
			args: args{
				pattern: `\d{3}$`,
				str:     "299",
			},
		},
		{
			name: "numeric - 3",
			want: false,
			args: args{
				pattern: `^\d{2}$`,
				str:     "299",
			},
		},
		{
			name: "US phone number - 1",
			want: false,
			args: args{
				pattern: `^\d{3}-\d{3}-\d{4}$`,
				str:     "55-5-5555-5555",
			},
		},
		{
			name: "US phone number - 2",
			want: true,
			args: args{
				pattern: `^\d{3}-\d{3}-\d{4}$`,
				str:     "555-555-5555",
			},
		},
		{
			name: "invalid - 1",
			want: false,
			args: args{
				pattern: `/\`,
				str:     "#FBA",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PassesRegex(tt.args.pattern, tt.args.str); got != tt.want {
				t.Errorf("PassesRegex() = %v, want %v", got, tt.want)
			}
		})
	}
}
