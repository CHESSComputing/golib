package embed

import (
	"fmt"
	"os"
	"testing"
)

const testDBPath = "test_clover_db"

func TestInitDB(t *testing.T) {
	// Initialize the test DB
	InitDB(testDBPath)
	defer os.RemoveAll(testDBPath) // Cleanup after test

	if db == nil {
		t.Fatal("Database initialization failed")
	}
}

func TestUpsert(t *testing.T) {
	InitDB(testDBPath)
	defer os.RemoveAll(testDBPath)

	collname := "users"
	records := []map[string]any{
		{"name": "Alice", "age": 30},
		{"name": "Bob", "age": 25},
	}

	err := Upsert(testDBPath, collname, "name", records)
	if err != nil {
		t.Fatalf("Upsert failed: %v", err)
	}

	// Check if the records were inserted
	results := Get(testDBPath, collname, map[string]any{"name": "Alice"}, 0, 0)
	if len(results) != 1 || results[0]["name"] != "Alice" {
		t.Fatalf("Upsert failed to insert record: %+v", results)
	}

	// Upsert (update) an existing record
	err = Upsert(testDBPath, collname, "name", []map[string]any{{"name": "Alice", "age": 35}})
	if err != nil {
		t.Fatalf("Upsert failed: %v", err)
	}

	results = Get(testDBPath, collname, map[string]any{"name": "Alice"}, 0, 0)
	age := fmt.Sprintf("%v", results[0]["age"])
	if len(results) != 1 || age != "35" {
		t.Fatalf("Upsert failed to update record: results %+v expect age=0", results)
	}
}

func TestGet(t *testing.T) {
	InitDB(testDBPath)
	defer os.RemoveAll(testDBPath)

	collname := "users"
	records := []map[string]any{
		{"name": "Charlie", "age": 40},
		{"name": "Daisy", "age": 20},
	}
	_ = Upsert(testDBPath, collname, "name", records)

	// Query all records
	results := Get(testDBPath, collname, map[string]any{}, 0, 0)
	if len(results) != 2 {
		t.Fatalf("Get failed to retrieve all records: %+v", results)
	}

	// Query specific record
	results = Get(testDBPath, collname, map[string]any{"name": "Charlie"}, 0, 0)
	if len(results) != 1 || results[0]["name"] != "Charlie" {
		t.Fatalf("Get failed to retrieve specific record: %+v", results)
	}
}

func TestUpdate(t *testing.T) {
	InitDB(testDBPath)
	defer os.RemoveAll(testDBPath)

	collname := "users"
	records := []map[string]any{
		{"name": "Eve", "age": 50},
	}
	_ = Upsert(testDBPath, collname, "name", records)

	// Update record
	Update(testDBPath, collname, map[string]any{"name": "Eve"}, map[string]any{"age": 55})

	results := Get(testDBPath, collname, map[string]any{"name": "Eve"}, 0, 0)
	if len(results) != 1 {
		t.Fatalf("Update failed: wrong number of records %d expect 1", len(results))
	}
	age := fmt.Sprintf("%v", results[0]["age"])
	if age != "55" {
		t.Fatalf("Update failed: obtain age=%s, expect age=55", age)
	}
}

func TestCount(t *testing.T) {
	InitDB(testDBPath)
	defer os.RemoveAll(testDBPath)

	collname := "users"
	records := []map[string]any{
		{"name": "Frank", "age": 60},
		{"name": "Grace", "age": 70},
	}
	_ = Upsert(testDBPath, collname, "name", records)

	count := Count(testDBPath, collname, map[string]any{})
	if count != 2 {
		t.Fatalf("Count failed: expected 2, got %d", count)
	}

	count = Count(testDBPath, collname, map[string]any{"name": "Grace"})
	if count != 1 {
		t.Fatalf("Count failed: expected 1, got %d", count)
	}
}

func TestRemove(t *testing.T) {
	InitDB(testDBPath)
	defer os.RemoveAll(testDBPath)

	collname := "users"
	records := []map[string]any{
		{"name": "Hank", "age": 30},
		{"name": "Ivy", "age": 25},
	}
	_ = Upsert(testDBPath, collname, "name", records)

	// Remove one record
	err := Remove(testDBPath, collname, map[string]any{"name": "Hank"})
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	results := Get(testDBPath, collname, map[string]any{}, 0, 0)
	if len(results) != 1 || results[0]["name"] != "Ivy" {
		t.Fatalf("Remove failed to delete record: %+v", results)
	}

	// Remove all records
	err = Remove(testDBPath, collname, map[string]any{})
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	results = Get(testDBPath, collname, map[string]any{}, 0, 0)
	if len(results) != 0 {
		t.Fatalf("Remove failed to delete all records: %+v", results)
	}
}
