package docdb

import (
	"errors"
	"fmt"
	"log"

	srvConfig "github.com/CHESSComputing/golib/config"
	embed "github.com/CHESSComputing/golib/embed/badger"
	mongo "github.com/CHESSComputing/golib/mongo"
)

// DocDB represents generic interface for document-oriented database
type DocDB interface {
	InitDB(uri string)
	Insert(dbname, collname string, records []map[string]any)
	Upsert(dbname, collname, attr string, records []map[string]any) error
	Get(dbname, collname string, spec map[string]any, idx, limit int) []map[string]any
	GetProjection(dbname, collname string, spec map[string]any, projection map[string]int, idx, limit int) []map[string]any
	Update(dbname, collname string, spec, newdata map[string]any) error
	Count(dbname, collname string, spec map[string]any) int
	Remove(dbname, collname string, spec map[string]any) error
	Distinct(dbname, collname, field string) ([]any, error)
	InsertRecord(dbname, collname string, rec map[string]any) error
	GetSorted(dbname, collname string, spec map[string]any, skeys []string, sortOrder, idx, limit int) []map[string]any
}

// InitializeDocDB initializes either mongo or embed database based on server configuration
func InitializeDocDB(uri string) (DocDB, error) {
	var docDB DocDB
	var err error
	dbType := "mongo"
	if srvConfig.Config.Embed.DocDb != "" {
		dbType = "embed"
	}
	switch dbType {
	case "mongo":
		log.Println("Initializing DocDB with MongoDB backend")
		docDB = &mongo.MongoDB{}
	case "embed":
		log.Printf("Initializing DocDB with embed DB backend %s", srvConfig.Config.Embed.DocDb)
		docDB = &embed.EmbedDB{}
	default:
		err = errors.New(fmt.Sprintf("Unsupported database type: %s", dbType))
	}
	if docDB != nil {
		docDB.InitDB(uri)
	}
	return docDB, err
}
