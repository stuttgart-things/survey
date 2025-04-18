package survey

import (
	"reflect"
	"testing"
)

func TestConvertToType(t *testing.T) {
	tests := []struct {
		name  string
		value string
		typ   string
		want  interface{}
	}{
		{
			name:  "String to int",
			value: "42",
			typ:   "int",
			want:  42,
		},
		{
			name:  "Invalid string to int",
			value: "notanumber",
			typ:   "int",
			want:  0,
		},
		{
			name:  "String to boolean (true)",
			value: "true",
			typ:   "boolean",
			want:  true,
		},
		{
			name:  "String to boolean (false)",
			value: "false",
			typ:   "boolean",
			want:  false,
		},
		{
			name:  "String to boolean (Yes)",
			value: "Yes",
			typ:   "boolean",
			want:  true,
		},
		{
			name:  "String remains string",
			value: "hello",
			typ:   "string",
			want:  "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToType(tt.value, tt.typ); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPow10(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{"0", 0, 1},
		{"1", 1, 10},
		{"2", 2, 100},
		{"3", 3, 1000},
		{"negative", -1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pow10(tt.n); got != tt.want {
				t.Errorf("pow10() = %v, want %v", got, tt.want)
			}
		})
	}
}
