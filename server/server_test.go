package server

import (
	"strings"
	"testing"
)

// TestLogName
func TestLogName(t *testing.T) {
	lname := logName("test")
	if !strings.Contains(lname, "_%Y%m%d") {
		t.Error("Invalid log name", lname)
	}
}
