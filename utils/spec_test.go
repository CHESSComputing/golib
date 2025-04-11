package utils

import (
	"testing"
)

// TestNormalizeSpec provide unit test for NormalizeSpec function
func TestNormalizeSpec(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Title:Some Value", "title:Some Value"},
		{"Author : John Doe", "author:John Doe"},
		{"publisher:Penguin", "publisher:Penguin"},
		{`{"Title": "Some Value"}`, `{"title":"Some Value"}`},
		{`{ "Key": "Value", "Other": 123 }`, `{"key":"Value","other":123}`},
		{"invalidinput", "invalidinput"},
		{"NoColonHere", "NoColonHere"},
	}

	for _, tt := range tests {
		result := NormalizeSpec(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeSpec(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}
