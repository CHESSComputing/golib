package docdb

import (
	srvConfig "github.com/CHESSComputing/golib/config"
	embed "github.com/CHESSComputing/golib/embed/badger"
	"github.com/CHESSComputing/golib/mongo"
)

// InitDocDB initializes document-oriented db connection object
func InitDocDB(uri string) {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	if srvConfig.Config.Embed.DocDb != "" {
		embed.InitDB(uri)
		return
	}
	mongo.InitMongoDB(uri)
}

// Upsert records into document-oriented db
func Upsert(dbname, collname, attr string, records []map[string]any) error {
	if srvConfig.Config.Embed.DocDb != "" {
		return embed.Upsert(dbname, collname, attr, records)
	}
	return mongo.Upsert(dbname, collname, attr, records)
}

// Get records from document-oriented db
func Get(dbname, collname string, spec map[string]any, idx, limit int) []map[string]any {
	if srvConfig.Config.Embed.DocDb != "" {
		return embed.Get(dbname, collname, spec, idx, limit)
	}
	return mongo.Get(dbname, collname, spec, idx, limit)
}

// Update inplace for given spec
func Update(dbname, collname string, spec, newdata map[string]any) {
	if srvConfig.Config.Embed.DocDb != "" {
		embed.Update(dbname, collname, spec, newdata)
		return
	}
	mongo.Update(dbname, collname, spec, newdata)
}

// Count gets number records from document-oriented db
func Count(dbname, collname string, spec map[string]any) int {
	if srvConfig.Config.Embed.DocDb != "" {
		return embed.Count(dbname, collname, spec)
	}
	return mongo.Count(dbname, collname, spec)
}

// Remove records from document-oriented db
func Remove(dbname, collname string, spec map[string]any) error {
	if srvConfig.Config.Embed.DocDb != "" {
		embed.Remove(dbname, collname, spec)
	}
	return mongo.Remove(dbname, collname, spec)
}

// Distinct gets number records from document-oriented db
func Distinct(dbname, collname, field string) ([]any, error) {
	if srvConfig.Config.Embed.DocDb != "" {
		embed.Distinct(dbname, collname, field)
	}
	return mongo.Distinct(dbname, collname, field)
}

// InsertRecord insert record with given spec to document-oriented db
func InsertRecord(dbname, collname string, rec map[string]any) error {
	if srvConfig.Config.Embed.DocDb != "" {
		embed.InsertRecord(dbname, collname, rec)
	}
	return mongo.InsertRecord(dbname, collname, rec)
}

// GetSorted fetches records from document-oriented db sorted by given key with specific order
func GetSorted(dbname, collname string, spec map[string]any, skeys []string, sortOrder, idx, limit int) []map[string]any {
	if srvConfig.Config.Embed.DocDb != "" {
		embed.GetSorted(dbname, collname, spec, skeys, sortOrder, idx, limit)
	}
	return mongo.GetSorted(dbname, collname, spec, skeys, sortOrder, idx, limit)
}
