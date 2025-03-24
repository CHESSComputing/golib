package sqldb

import (
	"database/sql"
	"log"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// InitDB initializes database with according to server configuration
func InitDB(dbtype, dburi string) (*sql.DB, error) {
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
	if srvConfig.Config.DataBookkeeping.MaxDBConnections == 0 {
		db.SetMaxOpenConns(100) // Allow up to 100 open connections
	} else {
		db.SetMaxOpenConns(srvConfig.Config.DataBookkeeping.MaxDBConnections)
	}
	if srvConfig.Config.DataBookkeeping.MaxIdleConnections == 0 {
		db.SetMaxIdleConns(50) // Keep up to 50 idle connections ready
	} else {
		db.SetMaxIdleConns(srvConfig.Config.DataBookkeeping.MaxIdleConnections)
	}
	db.SetConnMaxLifetime(5 * time.Minute) // Recycle connections after 5 minutes
	// Disables connection pool for sqlite3. This enables some concurrency with sqlite3 databases
	// See https://stackoverflow.com/questions/57683132/turning-off-connection-pool-for-go-http-client
	// and https://sqlite.org/wal.html
	// This only will apply to sqlite3 databases
	if dbtype == "sqlite3" || dbtype == "sqlite" {
		db.Exec("PRAGMA journal_mode=WAL;")
	}
	return db, nil
}
