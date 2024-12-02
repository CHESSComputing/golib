package mongo

import (
	"testing"
)

// TestMongoInsert
func TestMongoInsert(t *testing.T) {
	// our db attributes
	dbname := "chess"
	collname := "test"
	InitMongoDB("mongodb://localhost:8230")

	// remove all records in test collection
	Remove(dbname, collname, map[string]any{})

	// insert one record
	var records []map[string]any
	dataset := "/a/b/c"
	rec := map[string]any{"dataset": dataset}
	records = append(records, rec)
	Insert(dbname, collname, records)

	// look-up one record
	spec := map[string]any{"dataset": dataset}
	idx := 0
	limit := 1
	records = Get(dbname, collname, spec, idx, limit)
	if len(records) != 1 {
		t.Errorf("unable to find records using spec '%s', records %+v", spec, records)
	}

	// modify our record
	rec = map[string]any{"dataset": dataset, "test": 1}
	records = []map[string]any{}
	records = append(records, rec)
	err := Upsert(dbname, collname, "dataset", records)
	if err != nil {
		t.Error(err)
	}
	spec = map[string]any{"test": 1}
	records = Get(dbname, collname, spec, idx, limit)
	if len(records) != 1 {
		t.Errorf("unable to find records using spec '%s', records %+v", spec, records)
	}
}
