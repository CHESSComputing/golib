package embed

import (
	"fmt"
	"testing"
)

func setupBadgerDB(t *testing.T) string {
	tempDir := t.TempDir()
	err := InitDB(tempDir)
	if err != nil {
		t.Fatalf("Failed to initialize BadgerDB: %v", err)
	}
	return tempDir
}

func teardownBadgerDB() {
	if db != nil {
		db.Close()
	}
}

func TestUpsertAndGet(t *testing.T) {
	defer teardownBadgerDB()
	setupBadgerDB(t)

	records := []map[string]interface{}{
		{"id": 1, "name": "Alice", "age": 25},
		{"id": 2, "name": "Bob", "age": 30},
	}

	err := upsert("users", records)
	if err != nil {
		t.Fatalf("Failed to upsert records: %v", err)
	}

	results, err := get("users", map[string]interface{}{"name": "Alice"})
	if err != nil {
		t.Fatalf("Failed to get records: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(results))
	}

	if results[0]["name"] != "Alice" {
		t.Fatalf("Expected name to be Alice, got %s", results[0]["name"])
	}
}

func TestUpdate(t *testing.T) {
	defer teardownBadgerDB()
	setupBadgerDB(t)

	records := []map[string]interface{}{
		{"id": 1, "name": "Alice", "age": 25},
	}

	err := upsert("users", records)
	if err != nil {
		t.Fatalf("Failed to upsert records: %v", err)
	}

	err = update("users", map[string]interface{}{"name": "Alice"}, map[string]interface{}{"age": 26})
	if err != nil {
		t.Fatalf("Failed to update records: %v", err)
	}

	results, err := get("users", map[string]interface{}{"name": "Alice"})
	if err != nil {
		t.Fatalf("Failed to get records: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(results))
	}

	age := fmt.Sprintf("%v", results[0]["age"])
	if age != "26" {
		t.Fatalf("Expected age to be 26, got %v", age)
	}
}

func TestCount(t *testing.T) {
	defer teardownBadgerDB()
	setupBadgerDB(t)

	records := []map[string]interface{}{
		{"id": 1, "name": "Alice", "age": 25},
		{"id": 2, "name": "Bob", "age": 30},
	}

	err := upsert("users", records)
	if err != nil {
		t.Fatalf("Failed to upsert records: %v", err)
	}

	count, err := count("users", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to count records: %v", err)
	}

	if count != 2 {
		t.Fatalf("Expected count to be 2, got %d", count)
	}
}

func TestRemove(t *testing.T) {
	defer teardownBadgerDB()
	setupBadgerDB(t)

	records := []map[string]interface{}{
		{"id": 1, "name": "Alice", "age": 25},
	}

	err := upsert("users", records)
	if err != nil {
		t.Fatalf("Failed to upsert records: %v", err)
	}

	err = remove("users", map[string]interface{}{"name": "Alice"})
	if err != nil {
		t.Fatalf("Failed to remove records: %v", err)
	}

	results, err := get("users", map[string]interface{}{"name": "Alice"})
	if err != nil {
		t.Fatalf("Failed to get records: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("Expected 0 records, got %d", len(results))
	}
}
