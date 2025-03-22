package doi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"text/template"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// Provider represents generic DOI interface
type Provider interface {
	Init()
	Publish(did, description string, record map[string]any, publish bool) (string, string, error)
}

// DOIData represents structure of public DOI attributes which will be written to DOI record
type DOIData struct {
	PI             string
	Facility       string
	Affiliation    string
	StaffScientist string
	Beamline       any
}

// Default template string
const defaultTmpl = `<html><body>
PI: {{.PI}}
<br/>
Facility: {{.Facility}}
<br/>
Beamline: {{.Beamline}}
<br/>
Affiliation: {{.Affiliation}}
<br/>
StaffScientist: {{.StaffScientist}}
</body></html>`

// CreateEntry creates DOI entry for DOIService
func CreateEntry(doi string, rec map[string]any, writeMeta bool) error {
	doiDir := srvConfig.Config.DOI.DocumentDir
	if doiDir == "" {
		return errors.New("no DOI.DocumentDir configuration found")
	}
	// Create directory with 0755 permissions
	err := os.MkdirAll(fmt.Sprintf("%s/%s", doiDir, doi), 0755)
	if err != nil {
		log.Println("Error creating directory:", err)
		return err
	}
	if writeMeta {
		fname := fmt.Sprintf("%s/%s/metadata.json", doiDir, doi)
		file, err := os.Create(fname)
		if err != nil {
			log.Println("unable to create file", err)
			return err
		}
		defer file.Close()
		data, err := json.Marshal(rec)
		if err != nil {
			return err
		}
		file.Write(data)
	}
	doiData := DOIData{}
	if val, ok := rec["beamline"]; ok {
		doiData.Beamline = val
	}
	if val, ok := rec["pi"]; ok {
		doiData.PI = val.(string)
	}
	if val, ok := rec["affiliation"]; ok {
		doiData.Affiliation = val.(string)
	}
	if val, ok := rec["staff_scientist"]; ok {
		doiData.StaffScientist = val.(string)
	}
	result, err := RenderTemplate(srvConfig.Config.DOI.TemplateFile, doiData)
	if err != nil {
		return err
	}
	fname := fmt.Sprintf("%s/%s/index.html", doiDir, doi)
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write([]byte(result))
	return nil
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
