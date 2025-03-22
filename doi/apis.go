package doi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
	sqldb "github.com/CHESSComputing/golib/sqldb"
)

// global variables
var dbtype, dburi, dbowner string

// Provider represents generic DOI interface
type Provider interface {
	Init()
	Publish(did, description string, record map[string]any, publish bool) (string, string, error)
}

// DOIData represents structure of public DOI attributes which will be written to DOI record
type DOIData struct {
	Doi         string
	Did         string
	Description string
	Metadata    string
	Published   int64
}

// Default template string
const defaultTmpl = `<html><body>
DOI: {{.DOI}}
<br/>
DID: {{.DID}}
<br/>
Description: {{.Description}}
<br/>
Metadata: {{.Metadata}}
<br/>
Published: {{.Published}}
</body></html>`

// CreateEntry creates DOI entry for DOIService
func CreateEntry(doi string, rec map[string]any, description string, writeMeta bool) error {
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
	if writeMeta {
		data, err := json.MarshalIndent(rec, "", "   ")
		if err != nil {
			return err
		}
		doiData.Metadata = string(data)
	}
	err := insertData(doiData)
	return err
}

// helper function to insert data into DOI database
func insertData(data DOIData) error {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	if dbtype == "" || dburi == "" || dbowner == "" {
		log.Printf("InitDB: type=%s owner=%s", dbtype, dbowner)
		dbtype, dburi, dbowner = sqldb.ParseDBFile(srvConfig.Config.DOI.DBFile)
	}
	db, err := sqldb.InitDB(dbtype, dburi)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `INSERT INTO dois (doi,did,description,metadata,published) VALUES (?,?,?,?,?)`
	_, err = tx.Exec(query, data.Doi, data.Did, data.Description, data.Metadata, data.Published)
	if err != nil {
		log.Printf("Could not insert record to dois table; error: %v", err)
		return err
	}
	err = tx.Commit()
	return err
}

// RenderTemplate processes a template from a file if provided, otherwise, it uses a default template.
func RenderTemplate(fileName string, data DOIData) (string, error) {
	var tmplContent string

	// If a file name is provided, read the template from the file
	if fileName != "" {
		content, err := os.ReadFile(fileName)
		if err != nil {
			return "", fmt.Errorf("failed to read template file: %v", err)
		}
		tmplContent = string(content)
	} else {
		tmplContent = defaultTmpl
	}

	// Parse the template content
	t, err := template.New("template").Parse(tmplContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	var result bytes.Buffer
	err = t.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return result.String(), nil
}
