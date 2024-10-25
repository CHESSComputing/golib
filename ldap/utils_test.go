package ldap

import (
	"testing"
)

// TestRemoveSuffix provides unit test for removeSuffix function
func TestRemoveSuffix(t *testing.T) {
	tests := []struct {
		input    string
		suffix   string
		expected string
	}{
		{"bla-m-m", "-m", "bla-m"},
		{"hello-world", "-world", "hello"},
		{"example", "-m", "example"}, // suffix not present
		{"test-suffix-suffix", "-suffix", "test-suffix"},
		{"-suffix", "-suffix", ""},                      // entire string is suffix
		{"no-suffix-here", "-suffix", "no-suffix-here"}, // suffix not present
	}

	for _, test := range tests {
		result := removeSuffix(test.input, test.suffix)
		if result != test.expected {
			t.Errorf("removeSuffix(%q, %q) = %q; want %q", test.input, test.suffix, result, test.expected)
		}
	}
}
