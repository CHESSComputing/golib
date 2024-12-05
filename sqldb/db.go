package sqldb

import (
	"database/sql"
	"log"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// InitDB initializes database with according to server configuration
func InitDB(dbtype, dburi string) (*sql.DB, error) {
	// load web section schema file
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	db, dberr := dbOpen(dbtype, dburi)
	if dberr != nil {
		log.Printf("unable to open dbtype=%s dburi=%s, error %v", dbtype, dburi, dberr)
		return nil, dberr
	}
	dberr = db.Ping()
	if dberr != nil {
		log.Println("DB ping error", dberr)
		return nil, dberr
	}
	db.SetMaxOpenConns(srvConfig.Config.DataBookkeeping.MaxDBConnections)
	db.SetMaxIdleConns(srvConfig.Config.DataBookkeeping.MaxIdleConnections)
	// Disables connection pool for sqlite3. This enables some concurrency with sqlite3 databases
	// See https://stackoverflow.com/questions/57683132/turning-off-connection-pool-for-go-http-client
	// and https://sqlite.org/wal.html
	// This only will apply to sqlite3 databases
	if dbtype == "sqlite3" {
		db.Exec("PRAGMA journal_mode=WAL;")
	}
	return db, nil
}
