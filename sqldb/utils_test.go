package sqldb

import (
	"os"
	"testing"
)

// TestParseDBFile tests the ParseDBFile function.
func TestParseDBFile(t *testing.T) {
	// Create a temporary test file.
	tempFile := "testdbfile.txt"
	testData := "dbtype dburi dbowner\n"
	err := os.WriteFile(tempFile, []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tempFile) // Clean up after the test.

	// Call the function and get the results.
	dbtype, dburi, dbowner := ParseDBFile(tempFile)

	// Validate the results.
	if dbtype != "dbtype" {
		t.Errorf("Expected dbtype to be 'dbtype', got '%s'", dbtype)
	}
	if dburi != "dburi" {
		t.Errorf("Expected dburi to be 'dburi', got '%s'", dburi)
	}
	if dbowner != "dbowner" {
		t.Errorf("Expected dbowner to be 'dbowner', got '%s'", dbowner)
	}
}
