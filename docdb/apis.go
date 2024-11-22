package docdb

import (
	srvConfig "github.com/CHESSComputing/golib/config"
	embed "github.com/CHESSComputing/golib/embed"
	"github.com/CHESSComputing/golib/mongo"
	bson "go.mongodb.org/mongo-driver/bson"
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
func Get(dbname, collname string, spec bson.M, idx, limit int) []map[string]any {
	if srvConfig.Config.Embed.DocDb != "" {
		embed.Get(dbname, collname, spec, idx, limit)
	}
	return mongo.Get(dbname, collname, spec, idx, limit)
}

// Update inplace for given spec
func Update(dbname, collname string, spec, newdata bson.M) {
	if srvConfig.Config.Embed.DocDb != "" {
		embed.Update(dbname, collname, spec, newdata)
		return
	}
	mongo.Update(dbname, collname, spec, newdata)
}

// Count gets number records from document-oriented db
func Count(dbname, collname string, spec bson.M) int {
	if srvConfig.Config.Embed.DocDb != "" {
		return embed.Count(dbname, collname, spec)
	}
	return mongo.Count(dbname, collname, spec)
}

// Remove records from document-oriented db
func Remove(dbname, collname string, spec bson.M) error {
	if srvConfig.Config.Embed.DocDb != "" {
		embed.Remove(dbname, collname, spec)
	}
	return mongo.Remove(dbname, collname, spec)
}
