package utils

import (
	"reflect"
	"testing"
)

// TestSplitStr2List provides unit test for SplitStr2List function
func TestSplitStr2List(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "empty string",
			in:   "",
			want: nil,
		},
		{
			name: "whitespace only",
			in:   "   ",
			want: nil,
		},
		{
			name: "single value",
			in:   "foo",
			want: []string{"foo"},
		},
		{
			name: "space separated",
			in:   "foo bar baz",
			want: []string{"foo", "bar", "baz"},
		},
		{
			name: "multiple spaces",
			in:   "foo   bar    baz",
			want: []string{"foo", "bar", "baz"},
		},
		{
			name: "comma separated",
			in:   "foo,bar,baz",
			want: []string{"foo", "bar", "baz"},
		},
		{
			name: "comma with spaces",
			in:   "foo, bar ,  baz ",
			want: []string{"foo", "bar", "baz"},
		},
		{
			name: "single comma value",
			in:   "foo,",
			want: []string{"foo", ""},
		},
		{
			name: "numbers",
			in:   "1, 2, 3",
			want: []string{"1", "2", "3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitStr2List(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("splitList(%q) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

// TestConvert2dtype provides unit test for Convert2dtype function
func TestConvert2dtype(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		dtype   string
		want    any
		wantErr bool
	}{
		{"string", "hello", "string", "hello", false},
		{"int", "42", "int", 42, false},
		{"int8", "127", "int8", int8(127), false},
		{"int8 overflow", "128", "int8", nil, true},
		{"int16", "32000", "int16", int16(32000), false},
		{"int32", "123456", "int32", int32(123456), false},
		{"int64", "9223372036854775807", "int64", int64(9223372036854775807), false},

		{"float64", "3.14", "float64", 3.14, false},
		{"float32", "3.14", "float32", float32(3.14), false},

		{"bool true", "true", "bool", true, false},
		{"bool false", "false", "bool", false, false},

		{"list_str comma", "a,b,c", "list_str", []string{"a", "b", "c"}, false},
		{"list_str space", "a b c", "list_str", []string{"a", "b", "c"}, false},

		{"list_int", "1,2,3", "list_int", []int{1, 2, 3}, false},
		{"list_float", "1.1 2.2 3.3", "list_float", []float64{1.1, 2.2, 3.3}, false},

		{"unsupported", "x", "uuid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert2dtype(tt.val, tt.dtype)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v (%T), want %#v (%T)", got, got, tt.want, tt.want)
			}
		})
	}
}

// TestConvert2records provides unit test for Convert2records function
func TestConvert2records(t *testing.T) {
	input := map[string][]string{
		"name":  {"alice", "bob"},
		"age":   {"30", "40"},
		"email": {"a@example.com"},
	}

	want := []map[string]string{
		{
			"name":  "alice",
			"age":   "30",
			"email": "a@example.com",
		},
		{
			"name": "bob",
			"age":  "40",
		},
	}

	got := Convert2records(input)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}
