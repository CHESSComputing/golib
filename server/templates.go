package server

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	utils "github.com/CHESSComputing/golib/utils"
)

// global variables
var _header, _footer, _footerEmpty string

// Header helper function to define our header
func Header(fsys fs.FS, base string) string {
	if _header == "" {
		tmpl := MakeTmpl(fsys, "Header")
		tmpl["Base"] = base
		_header = TmplPage(fsys, "header.tmpl", tmpl)
	}
	return _header
}

// Footer helper function to define our footer
func Footer(fsys fs.FS, base string) string {
	if _footer == "" {
		tmpl := MakeTmpl(fsys, "Footer")
		tmpl["Base"] = base
		_footer = TmplPage(fsys, "footer.tmpl", tmpl)
	}
	return _footer
}

// FooterEmpty helper function to define our footer
func FooterEmpty(fsys fs.FS, base string) string {
	if _footerEmpty == "" {
		tmpl := MakeTmpl(fsys, "Footer")
		tmpl["Base"] = base
		_footerEmpty = TmplPage(fsys, "footer_empty.tmpl", tmpl)
	}
	return _footerEmpty
}

// Base helper function to handle base path of URL requests
func Base(base, api string) string {
	return utils.BasePath(base, api)
}

// consume list of templates and release their full path counterparts
func fileNames(tdir string, filenames ...string) []string {
	flist := []string{}
	for _, fname := range filenames {
		flist = append(flist, filepath.Join(tdir, fname))
	}
	return flist
}

// ParseTmpl parses template with given data
func ParseTmpl(tdir, tmpl string, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	filenames := fileNames(tdir, tmpl)
	t := template.Must(template.ParseFiles(filenames...))
	err := t.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}

// TmplRecord represent template record
type TmplRecord map[string]interface{}

// GetString converts given value for provided key to string data-type
func (t TmplRecord) GetString(key string) string {
	if v, ok := t[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// GetInt converts given value for provided key to int data-type
func (t TmplRecord) GetInt(key string) int {
	if v, ok := t[key]; ok {
		if val, err := strconv.Atoi(fmt.Sprintf("%v", v)); err == nil {
			return val
		} else {
			log.Println("ERROR:", err)
		}
	}
	return 0
}

// GetError returns error string
func (t TmplRecord) GetError() string {
	if v, ok := t["Error"]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// GetBytes returns bytes object for given key
func (t TmplRecord) GetBytes(key string) []byte {
	if data, ok := t[key]; ok {
		return data.([]byte)
	}
	return []byte{}
}

// GetElapsedTime returns elapsed time
func (t TmplRecord) GetElapsedTime() string {
	if val, ok := t["StartTime"]; ok {
		startTime := time.Unix(val.(int64), 0)
		return time.Since(startTime).String()
	}
	return ""
}

// Templates structure
type Templates struct {
	html string
}

// Tmpl method for ServerTemplates structure
func (q Templates) Tmpl(fsys fs.FS, tfile string, tmplData map[string]interface{}) string {
	if q.html != "" {
		return q.html
	}

	// get template from embed.FS
	filenames := []string{"static/templates/" + tfile}
	t := template.Must(template.New(tfile).ParseFS(fsys, filenames...))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, tmplData)
	if err != nil {
		log.Println("ERROR: template.Tmpl", err)
		return ""
	}
	q.html = buf.String()
	return q.html
}

// Tmpl method for ServerTemplates structure
func (q Templates) TextTmpl(fsys fs.FS, tfile string, tmplData map[string]interface{}) string {
	if q.html != "" {
		return q.html
	}

	// get template from embed.FS
	filenames := []string{"static/templates/" + tfile}
	t := template.Must(template.New(tfile).ParseFS(fsys, filenames...))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, tmplData)
	if err != nil {
		log.Println("ERROR: template.Tmpl", err)
		return ""
	}
	q.html = buf.String()
	return q.html
}

// TmplPage parses given template and return HTML page
func TmplPage(fsys fs.FS, tmpl string, tmplData TmplRecord) string {
	if tmplData == nil {
		tmplData = make(TmplRecord)
	}
	var templates Templates
	page := templates.Tmpl(fsys, tmpl, tmplData)
	return page
}
