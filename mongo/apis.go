package mongo

// MongoDB represents MongoDB interface with MongoDB backend
type MongoDB struct {
}

// InitDB initializes MongoDB backend
func (d *MongoDB) InitDB(uri string) {
	InitMongoDB(uri)
}

// Insert inserts records into provided database/collection
func (d *MongoDB) Insert(dbname, collname string, records []map[string]any) {
	Insert(dbname, collname, records)
}

// Upsert inserts records into provided database/collection and attribute
func (d *MongoDB) Upsert(dbname, collname, attr string, records []map[string]any) error {
	return Upsert(dbname, collname, attr, records)
}

// Get fetches data from underlying database/collection
func (d *MongoDB) Get(dbname, collname string, spec map[string]any, idx, limit int) []map[string]any {
	return Get(dbname, collname, spec, idx, limit)
}

// Update updates data into given database/collection
func (d *MongoDB) Update(dbname, collname string, spec, newdata map[string]any) error {
	return Update(dbname, collname, spec, newdata)
}

// Count returns total number of records within database/collection and given spec
func (d *MongoDB) Count(dbname, collname string, spec map[string]any) int {
	return Count(dbname, collname, spec)
}

// Remove deletes records in given database/collection using given spec
func (d *MongoDB) Remove(dbname, collname string, spec map[string]any) error {
	return Remove(dbname, collname, spec)
}

// Distinct returns distinct collection of records
func (d *MongoDB) Distinct(dbname, collname, field string) ([]any, error) {
	return Distinct(dbname, collname, field)
}

// InsertRecord inserts single record into given database/collection
func (d *MongoDB) InsertRecord(dbname, collname string, rec map[string]any) error {
	return InsertRecord(dbname, collname, rec)
}

// GetSorted returns sorted records from given database/collection using provided spec, sorted keys, order and limits
func (d *MongoDB) GetSorted(dbname, collname string, spec map[string]any, skeys []string, sortOrder, idx, limit int) []map[string]any {
	return GetSorted(dbname, collname, spec, skeys, sortOrder, idx, limit)
}
