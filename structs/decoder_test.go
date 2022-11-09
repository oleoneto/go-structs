package structs

import (
	"reflect"
	"testing"
)

func Test_Decode(t *testing.T) {
	type BasePerson struct {
		Name   *string  `json:"name" db:"name"`
		Emails []string `json:"emails" db:"emails"`
	}

	type Person struct {
		Id    string  `json:"id" jsonschema:"required"`
		Name  *string `json:"name" db:"name" jsonschema:"required"`
		Email string  `json:"email" db:"email"`
	}

	type args struct {
		data    []byte
		model   any
		options DecoderOptions
	}

	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "base person - 0",
			args: args{
				model:   &BasePerson{},
				data:    []byte{},
				options: DecoderOptions{Rules: []SchemaValidationRule{}},
			},
			want: map[string][]string{},
		},
		{
			name: "base person - 1",
			args: args{
				model:   &BasePerson{},
				data:    []byte{},
				options: DecoderOptions{Rules: []SchemaValidationRule{INVALID_TYPE}},
			},
			want: map[string][]string{},
		},
		{
			name: "base person - 2",
			args: args{
				model:   &BasePerson{},
				data:    []byte(`{}`),
				options: DecoderOptions{Rules: []SchemaValidationRule{INVALID_TYPE}},
			},
			want: map[string][]string{},
		},
		{
			name: "base person - 3",
			args: args{
				model:   &BasePerson{},
				data:    []byte(`{}`),
				options: DecoderOptions{Rules: []SchemaValidationRule{INVALID_TYPE, REQUIRED_ATTRIBUTE}},
			},
			want: map[string][]string{},
		},
		{
			name: "base person - 4",
			args: args{
				model:   &BasePerson{},
				data:    []byte(`{"extra": 1}`),
				options: DecoderOptions{Rules: []SchemaValidationRule{INVALID_TYPE, REQUIRED_ATTRIBUTE}},
			},
			want: map[string][]string{},
		},
		{
			name: "base person - 5",
			args: args{
				model:   &BasePerson{},
				data:    []byte(`{"extra": 1}`),
				options: DecoderOptions{Rules: []SchemaValidationRule{INVALID_TYPE, REQUIRED_ATTRIBUTE, ADDITIONAL_PROPERTY}},
			},
			want: map[string][]string{
				"extra": {"ADDITIONAL_PROPERTY"},
			},
		},
		{
			name: "person - 6",
			args: args{
				data:    []byte(`{}`),
				model:   &Person{},
				options: DecoderOptions{Rules: []SchemaValidationRule{INVALID_TYPE, REQUIRED_ATTRIBUTE}},
			},
			want: map[string][]string{
				"id":   {"REQUIRED_ATTRIBUTE_MISSING"},
				"name": {"REQUIRED_ATTRIBUTE_MISSING"},
			},
		},
		{
			name: "person - 7",
			args: args{
				data:    []byte(`{"id": 1, "name": 2}`),
				model:   &Person{},
				options: DecoderOptions{Rules: []SchemaValidationRule{INVALID_TYPE, REQUIRED_ATTRIBUTE}},
			},
			want: map[string][]string{
				"id":   {"INVALID_TYPE"},
				"name": {"INVALID_TYPE"},
			},
		},
		{
			name: "person - 8",
			args: args{
				data:    []byte(`{"email": "leonardo@example.com"}`),
				model:   &Person{},
				options: DecoderOptions{Rules: []SchemaValidationRule{INVALID_TYPE, REQUIRED_ATTRIBUTE}},
			},
			want: map[string][]string{
				"id":   {"REQUIRED_ATTRIBUTE_MISSING"},
				"name": {"REQUIRED_ATTRIBUTE_MISSING"},
			},
		},
		{
			name: "invalid payload - 1",
			args: args{
				data:  []byte(`{`),
				model: &Person{},
				options: DecoderOptions{
					Rules: []SchemaValidationRule{INVALID_TYPE, REQUIRED_ATTRIBUTE},
				},
			},
			want: map[string][]string{
				"_": {"INVALID_PAYLOAD"},
			},
		},
		{
			name: "invalid payload - 2",
			args: args{
				data:  []byte(`}`),
				model: &Person{},
				options: DecoderOptions{
					Rules: []SchemaValidationRule{INVALID_TYPE, REQUIRED_ATTRIBUTE},
				},
			},
			want: map[string][]string{
				"_": {"INVALID_PAYLOAD"},
			},
		},
		{
			name: "invalid payload - 3",
			args: args{
				data:  []byte(`-`),
				model: &Person{},
				options: DecoderOptions{
					Rules: []SchemaValidationRule{INVALID_TYPE, REQUIRED_ATTRIBUTE},
				},
			},
			want: map[string][]string{
				"_": {"INVALID_PAYLOAD"},
			},
		},
		{
			name: "before hook - set default values",
			args: args{
				data:  []byte(`{`),
				model: &Person{},
				options: DecoderOptions{
					Rules: []SchemaValidationRule{INVALID_TYPE},
					BeforeHook: func(data []byte, model any) []byte {
						return []byte(`{"id": "", "name": ""}`)
					},
				},
			},
			want: map[string][]string{},
		},
		{
			name: "after hook - set custom errors",
			args: args{
				data:  []byte(`}`),
				model: 0,
				options: DecoderOptions{
					Rules: []SchemaValidationRule{INVALID_TYPE},
					AfterHook: func(m map[string][]string) map[string][]string {
						return map[string][]string{"error": {"CUSTOM_ERROR"}}
					},
				},
			},
			want: map[string][]string{"error": {"CUSTOM_ERROR"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decode(tt.args.data, tt.args.model, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonAttributeName(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "generic - 1",
			arg:  "resource_id is not valid",
			want: "resource_id is not valid",
		},
		{
			name: "required - 1",
			arg:  "(root): resource_id is required",
			want: "resource_id",
		},
		{
			name: "required - 2",
			arg:  "resources.0: id is required",
			want: "resources[0].id",
		},
		{
			name: "invalid type - 1",
			arg:  "resource_id: Invalid type. Expected: string, given: integer",
			want: "resource_id",
		},
		{
			name: "invalid type - 2",
			arg:  "resources.0: Invalid type. Expected: string, given: integer",
			want: "resources[0]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jsonAttributeName(tt.arg); got != tt.want {
				t.Errorf("jsonAttributeName() = %v, want %v", got, tt.want)
			}
		})
	}
}
