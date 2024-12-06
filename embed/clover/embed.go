package embed

import (
	"log"
	"os"

	clover "github.com/ostafen/clover/v2"
	cloverD "github.com/ostafen/clover/v2/document"
	cloverQ "github.com/ostafen/clover/v2/query"
)

var db *clover.DB

// InitDB initializes document-oriented db connection object
func InitDB(dbDir string) {
	var err error

	// check if dbDir exist
	_, err = os.Stat(dbDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dbDir, os.ModePerm) // Create all parent directories if necessary
		if err != nil {
			log.Fatal(err)
		}
	}

	// Initialize Clover database
	db, err = clover.Open(dbDir)
	if err != nil {
		log.Fatalf("Failed to open Clover database: %v", err)
	}

	log.Printf("Clover database initialized at %s", dbDir)
}

// Insert records into document-oriented db
func Insert(dbname, collname string, records []map[string]any) {
	Upsert(dbname, collname, "", records)
}

// Upsert records into document-oriented db
func Upsert(dbname, collname, attr string, records []map[string]any) error {
	if err := db.CreateCollection(collname); err != nil && err != clover.ErrCollectionExist {
		log.Fatalf("Failed to create collection: %v", err)
	}
	for _, record := range records {
		// Search for an existing record with the same attribute value
		query := cloverQ.NewQuery(collname).Where(cloverQ.Field(attr).Eq(record[attr]))
		existingDocs, err := db.FindAll(query)
		if err != nil {
			return err
		}

		if len(existingDocs) > 0 {
			// Update the existing record
			doc := existingDocs[0]
			for k, v := range record {
				doc.Set(k, v)
			}
			updater := func(_ *cloverD.Document) *cloverD.Document {
				return doc
			}
			if err := db.UpdateById(collname, doc.ObjectId(), updater); err != nil {
				return err
			}
		} else {
			// Insert a new record
			doc := cloverD.NewDocumentOf(record)
			if _, err := db.InsertOne(collname, doc); err != nil {
				return err
			}
		}
	}
	return nil
}

// Get records from document-oriented db
func Get(dbname, collname string, spec map[string]any, idx, limit int) []map[string]any {
	var results []map[string]any
	if err := db.CreateCollection(collname); err != nil && err != clover.ErrCollectionExist {
		log.Fatalf("Failed to create collection: %v", err)
	}

	// Build query based on the spec
	query := cloverQ.NewQuery(collname)
	for k, v := range spec {
		query = query.Where(cloverQ.Field(k).Eq(v))
	}

	// Set offset and limit
	if idx > 0 {
		query = query.Skip(idx)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Execute the query
	docs, err := db.FindAll(query)
	if err != nil {
		log.Printf("Failed to query documents: %v", err)
		return results
	}

	// Convert documents to map[string]any
	for _, doc := range docs {
		results = append(results, doc.AsMap())
	}
	return results
}

// Update inplace for given spec
func Update(dbname, collname string, spec, newdata map[string]any) {
	if err := db.CreateCollection(collname); err != nil && err != clover.ErrCollectionExist {
		log.Fatalf("Failed to create collection: %v", err)
	}
	// Build query based on the spec
	query := cloverQ.NewQuery(collname)
	for k, v := range spec {
		query = query.Where(cloverQ.Field(k).Eq(v))
	}

	// update document
	err := db.Update(query, newdata)
	if err != nil {
		log.Printf("Failed to update document: %v", err)
	}
}

// Count gets number of records from document-oriented db
func Count(dbname, collname string, spec map[string]any) int {
	if err := db.CreateCollection(collname); err != nil && err != clover.ErrCollectionExist {
		log.Fatalf("Failed to create collection: %v", err)
	}
	// Build query based on the spec
	query := cloverQ.NewQuery(collname)
	for k, v := range spec {
		query = query.Where(cloverQ.Field(k).Eq(v))
	}

	// Count documents
	docs, err := db.FindAll(query)
	if err != nil {
		log.Printf("Failed to count documents: %v", err)
		return 0
	}
	return len(docs)
}

// Remove records from document-oriented db
func Remove(dbname, collname string, spec map[string]any) error {
	if err := db.CreateCollection(collname); err != nil && err != clover.ErrCollectionExist {
		log.Fatalf("Failed to create collection: %v", err)
	}
	// Build query based on the spec
	query := cloverQ.NewQuery(collname)
	for k, v := range spec {
		query = query.Where(cloverQ.Field(k).Eq(v))
	}

	// Execute the query and delete documents
	err := db.Delete(query)
	if err != nil {
		log.Printf("Failed to remove documents: %v", err)
		return err
	}
	return nil
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
