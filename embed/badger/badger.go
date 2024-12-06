package embed

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dgraph-io/badger/v4"
)

var db *badger.DB

// InitDB initializes the Badger database
func InitDB(dbDir string) error {
	// Ensure the directory exists
	err := ensureDirExists(dbDir)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Open the Badger database
	opts := badger.DefaultOptions(dbDir).WithLogger(nil)
	db, err = badger.Open(opts)
	if err != nil {
		return fmt.Errorf("failed to open BadgerDB: %v", err)
	}
	return nil
}

// ensureDirExists checks if a directory exists, and creates it if it doesn't
func ensureDirExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

// Upsert records into BadgerDB
func Upsert(dbname, collname, attr string, records []map[string]any) error {
	return upsert(collname, records)
}

func upsert(collname string, records []map[string]interface{}) error {
	return db.Update(func(txn *badger.Txn) error {
		for _, record := range records {
			key := fmt.Sprintf("%s:%v", collname, record["id"])
			val, err := json.Marshal(record)
			if err != nil {
				return fmt.Errorf("failed to marshal record: %v", err)
			}
			err = txn.Set([]byte(key), val)
			if err != nil {
				return fmt.Errorf("failed to upsert record: %v", err)
			}
		}
		return nil
	})
}

// Get records from BadgerDB
func Get(dbname, collname string, spec map[string]any, idx, limit int) []map[string]any {
	var out []map[string]any
	results, err := get(collname, spec)
	if err != nil {
		log.Println("ERROR: ", err)
		return out
	}
	for i := idx; i < limit; i++ {
		if i < len(results) {
			out = append(out, results[i])
		}
	}
	return out
}
func get(collname string, spec map[string]interface{}) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(collname + ":")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var record map[string]interface{}
				err := json.Unmarshal(val, &record)
				if err != nil {
					return fmt.Errorf("failed to unmarshal record: %v", err)
				}
				match := true
				for k, v := range spec {
					if record[k] != v {
						match = false
						break
					}
				}
				if match {
					results = append(results, record)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("error reading value: %v", err)
			}
		}
		return nil
	})
	return results, err
}

// Update records in BadgerDB
func Update(dbname, collname string, spec, newdata map[string]any) {
	err := update(collname, spec, newdata)
	if err != nil {
		log.Println("ERROR:", err)
	}
}
func update(collname string, spec map[string]interface{}, newdata map[string]interface{}) error {
	return db.Update(func(txn *badger.Txn) error {
		results, err := get(collname, spec)
		if err != nil {
			return fmt.Errorf("failed to get records for update: %v", err)
		}

		for _, record := range results {
			for k, v := range newdata {
				record[k] = v
			}
			key := fmt.Sprintf("%s:%v", collname, record["id"])
			val, err := json.Marshal(record)
			if err != nil {
				return fmt.Errorf("failed to marshal updated record: %v", err)
			}
			err = txn.Set([]byte(key), val)
			if err != nil {
				return fmt.Errorf("failed to update record: %v", err)
			}
		}
		return nil
	})
}

// Count records in BadgerDB
func Count(dbname, collname string, spec map[string]any) int {
	nres, err := count(collname, spec)
	if err != nil {
		log.Println("ERROR:", err)
	}
	return nres
}

func count(collname string, spec map[string]interface{}) (int, error) {
	results, err := get(collname, spec)
	if err != nil {
		return 0, err
	}
	return len(results), nil
}

// Remove records from BadgerDB
func Remove(dbname, collname string, spec map[string]any) error {
	return remove(collname, spec)
}

func remove(collname string, spec map[string]interface{}) error {
	return db.Update(func(txn *badger.Txn) error {
		results, err := get(collname, spec)
		if err != nil {
			return fmt.Errorf("failed to get records for deletion: %v", err)
		}

		for _, record := range results {
			key := fmt.Sprintf("%s:%v", collname, record["id"])
			err := txn.Delete([]byte(key))
			if err != nil {
				return fmt.Errorf("failed to delete record: %v", err)
			}
		}
		return nil
	})
}

// Distinct gets number records from document-oriented db
func Distinct(dbname, collname, field string) ([]any, error) {
	spec := make(map[string]any)
	results, err := get(collname, spec)
	var out []any
	// loop over records and check if record contains the field key
	for _, rec := range results {
		if _, ok := rec[field]; ok {
			out = append(out, rec)
		}
	}
	return out, err
}

// InsertRecord insert record with given spec to document-oriented db
func InsertRecord(dbname, collname string, rec map[string]any) error {
	var records []map[string]any
	records = append(records, rec)
	return upsert(collname, records)
}

// GetSorted fetches records from document-oriented db sorted by given key with specific order
func GetSorted(dbname, collname string, spec map[string]any, skeys []string, sortOrder, idx, limit int) []map[string]any {
	var out []map[string]any
	results, err := get(collname, spec)
	if err != nil {
		log.Println("ERROR: ", err)
		return out
	}
	// TODO: implement how to properly sort records
	for i := idx; i < limit; i++ {
		if i < len(results) {
			out = append(out, results[i])
		}
	}
	return out
}
