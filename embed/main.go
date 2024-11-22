package embed

import (
	"log"

	tiedodb "github.com/HouzuoGuo/tiedot/db"
	bson "go.mongodb.org/mongo-driver/bson"
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

// Upsert records into document-oriented db
func Upsert(dbname, collname, attr string, records []map[string]any) error {
	var err error
	return err
}

// Get records from document-oriented db
func Get(dbname, collname string, spec bson.M, idx, limit int) []map[string]any {
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
func Update(dbname, collname string, spec, newdata bson.M) {
}

// Count gets number records from document-oriented db
func Count(dbname, collname string, spec bson.M) int {
	return 0
}

// Remove records from document-oriented db
func Remove(dbname, collname string, spec bson.M) error {
	var err error
	return err
}
