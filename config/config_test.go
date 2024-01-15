package config

import (
	"testing"
)

// TestConfig
func TestConfig(t *testing.T) {
	_, err := ParseConfig("bla")
	if err == nil {
		t.Error("Fail to parse non-existing config")
	}
}
