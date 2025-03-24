package doi

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
	sqldb "github.com/CHESSComputing/golib/sqldb"
)

// global variables
var _db *sql.DB

// Provider represents generic DOI interface
type Provider interface {
	Init()
	Publish(did, description string, record map[string]any, publish bool) (string, string, error)
}

// DOIData represents structure of public DOI attributes which will be written to DOI record
type DOIData struct {
	Doi            string
	Did            string
	Description    string
	AccessMetadata bool
	Published      int64
}

// Init function for this module
func Init() {
	if _db == nil {
		dbtype, dburi, dbowner := sqldb.ParseDBFile(srvConfig.Config.DOI.DBFile)
		log.Printf("InitDB: type=%s owner=%s", dbtype, dbowner)
		db, err := sqldb.InitDB(dbtype, dburi)
		if err != nil {
			log.Printf("ERROR: unable to initialize database, dbtype=%v, dburi=%v, error=%v", dbtype, dburi, err)
		}
		_db = db
	}
}

// CreateEntry creates DOI entry for DOIService
func CreateEntry(doi string, rec map[string]any, description string, accessMetadata bool) error {
	Init()
	doiData := DOIData{Doi: doi, Published: time.Now().Unix()}
	if val, ok := rec["did"]; ok {
		doiData.Did = val.(string)
	}
	if description != "" {
		doiData.Description = description
	} else {
		if val, ok := rec["description"]; ok {
			doiData.Description = val.(string)
		}
	}
	doiData.AccessMetadata = accessMetadata
	err := InsertData(doiData)
	return err
}

// helper function to insert data into DOI database
func InsertData(data DOIData) error {
	Init()
	tx, err := _db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `INSERT INTO dois (doi,did,description,metadata,published) VALUES (?,?,?,?,?)`
	_, err = tx.Exec(query, data.Doi, data.Did, data.Description, data.AccessMetadata, data.Published)
	if err != nil {
		log.Printf("Could not insert record to dois table; error: %v", err)
		return err
	}
	err = tx.Commit()
	return err
}

// GetData fetches records from the database based on the given ID
func GetData(doi string) ([]DOIData, error) {
	Init()
	var query string
	if strings.Contains(doi, "%") {
		query = `SELECT doi, did, description, metadata, published FROM dois WHERE doi LIKE ?`
	} else {
		query = `SELECT doi, did, description, metadata, published FROM dois WHERE doi = ?`
	}
	rows, err := _db.Query(query, doi)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %v", err)
	}
	defer rows.Close()

	var results []DOIData
	for rows.Next() {
		var d DOIData
		if err := rows.Scan(&d.Doi, &d.Did, &d.Description, &d.AccessMetadata, &d.Published); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		// if public metadata we retrieve its record from MetaData service
		results = append(results, d)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return results, nil
}
