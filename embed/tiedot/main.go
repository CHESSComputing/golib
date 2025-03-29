package embed

import (
	"log"

	tiedodb "github.com/HouzuoGuo/tiedot/db"
)

var db *tiedodb.DB

// InitDB initializes document-oriented db connection object
func InitDB(uri string) {
	dbDir := "./tiedotDB"
	var err error
	db, err = tiedodb.OpenDB(dbDir)
	if err != nil {
		log.Fatal(err)
	}
}

// Insert records into document-oriented db
func Insert(dbname, collname string, records []map[string]any) {
	Upsert(dbname, collname, "", records)
}

// Upsert records into document-oriented db
func Upsert(dbname, collname, attr string, records []map[string]any) error {
	var err error
	return err
}

// Get records from document-oriented db
func Get(dbname, collname string, spec map[string]any, idx, limit int) []map[string]any {
	coll := db.Use(collname)
	var results []map[string]any
	row := make(map[string]any)
	queryRes := make(map[int]struct{})
	if err := tiedodb.EvalQuery(spec, coll, &queryRes); err != nil {
		log.Fatal(err)
	}
	results = append(results, row)
	return results
}

// Update inplace for given spec
func Update(dbname, collname string, spec, newdata map[string]any) error {
	return nil
}

// Count gets number records from document-oriented db
func Count(dbname, collname string, spec map[string]any) int {
	return 0
}

// Remove records from document-oriented db
func Remove(dbname, collname string, spec map[string]any) error {
	var err error
	return err
}

// Distinct gets number records from document-oriented db
func Distinct(dbname, collname, field string) ([]any, error) {
	var out []any
	var err error
	// Not implemented yet
	return out, err
}

// InsertRecord insert record with given spec to document-oriented db
func InsertRecord(dbname, collname string, rec map[string]any) error {
	var records []map[string]any
	records = append(records, rec)
	return Upsert(dbname, collname, "", records)
}

// GetSorted fetches records from document-oriented db sorted by given key with specific order
func GetSorted(dbname, collname string, spec map[string]any, skeys []string, sortOrder, idx, limit int) []map[string]any {
	return Get(dbname, collname, spec, idx, limit)
}
