package embed

import (
	srvConfig "github.com/CHESSComputing/golib/config"
)

// EmbedDB represent embedded database
type EmbedDB struct {
}

// InitDB initialize embedded database
func (d *EmbedDB) InitDB(uri string) {
	InitDB(srvConfig.Config.Embed.DocDb)
}

// Insert inserts records into provided database/collection
func (d *EmbedDB) Insert(dbname, collname string, records []map[string]any) {
	Insert(dbname, collname, records)
}

// Upsert inserts records into provided database/collection and attribute
func (d *EmbedDB) Upsert(dbname, collname, attr string, records []map[string]any) error {
	return Upsert(dbname, collname, attr, records)
}

// Get fetches data from underlying database/collection
func (d *EmbedDB) Get(dbname, collname string, spec map[string]any, idx, limit int) []map[string]any {
	return Get(dbname, collname, spec, idx, limit)
}

// Update updates data into given database/collection
func (d *EmbedDB) Update(dbname, collname string, spec, newdata map[string]any) {
	Update(dbname, collname, spec, newdata)
}

// Count returns total number of records within database/collection and given spec
func (d *EmbedDB) Count(dbname, collname string, spec map[string]any) int {
	return Count(dbname, collname, spec)
}

// Remove deletes records in given database/collection using given spec
func (d *EmbedDB) Remove(dbname, collname string, spec map[string]any) error {
	return Remove(dbname, collname, spec)
}

// Distinct returns distinct collection of records
func (d *EmbedDB) Distinct(dbname, collname, field string) ([]any, error) {
	return Distinct(dbname, collname, field)
}

// InsertRecord inserts single record into given database/collection
func (d *EmbedDB) InsertRecord(dbname, collname string, rec map[string]any) error {
	return InsertRecord(dbname, collname, rec)
}

// GetSorted returns sorted records from given database/collection using provided spec, sorted keys, order and limits
func (d *EmbedDB) GetSorted(dbname, collname string, spec map[string]any, skeys []string, sortOrder, idx, limit int) []map[string]any {
	return GetSorted(dbname, collname, spec, skeys, sortOrder, idx, limit)
}
