package zenodo

import "errors"

// Some info are reverse engineered from the followin resources:
// https://felipecrp.com/2021/01/01/uploading-to-zenodo-through-api.html
// Zenodo REST API: https://developers.zenodo.org/
// snd, some discussion about zenodo APIs:
// https://github.com/zenodo/zenodo/issues/2168

// Creator represents creator record
type Creator struct {
	Name        string `json:"name"`
	Affiliation string `json:"affiliation"`
}

// MetaDataRecord represents meta-data record
type MetaDataRecord struct {
	PublicationType string    `json:"publication_type"`
	UploadType      string    `json:"upload_type"`
	Description     string    `json:"description"`
	Keywords        []string  `json:"keywords"`
	Title           string    `json:"title"`
	Licences        []string  `json:"licenses,omitempty"`
	Version         string    `json:"version,omitempty"`
	Publisher       string    `json:"publisher,omitempty"`
	Contributors    []Creator `json:"contributors,omitempty"`
	Creators        []Creator `json:"creators"`
}

// Validate provides validation of meta data record
func (m *MetaDataRecord) Validate() error {
	if m.PublicationType == "" {
		return errors.New("missing publication type, e.g. article")
	} else if m.UploadType == "" {
		return errors.New("missing upload type, e.g. publication")
	} else if m.Description == "" {
		return errors.New("missing description")
	} else if m.Title == "" {
		return errors.New("missing title")
	} else if len(m.Creators) == 0 {
		return errors.New("missing creators, e.g. [{\"name\":\"First Last\", \"affiliation\": \"Zenodo\"}]")
	}
	return nil
}

// Error represents individual zenodo error struct
type Error struct {
	Field    string
	Messages []string
}

// Reponse represents zenodo response
type Response struct {
	Status  int
	Message string
	Error   []Error
}

// CreateResponse represents output of /create API
type CreateResponse struct {
	Id        int64    `json:"id"`
	MetaData  MetaData `json:"metadata"`
	Created   string   `json:"created"`
	Modified  string   `json:"modified"`
	Owner     int      `json:"owner"`
	RecordId  int64    `json:"record_id"`
	State     string   `json:"state"`
	Submitted bool     `json:"submitted"`
	Title     string   `json:"title"`
	Links     Links    `json:"links"`
}

// PrereserveDoi represents PrereserveDoi struct
type PrereserveDoi struct {
	Doi   string `json:"doi"`
	RecId int64  `json:"recid"`
}

// MetaData represents meta-data struct
type MetaData struct {
	AccessRight   string        `json:"access_right"`
	PrereserveDoi PrereserveDoi `json:"prereserve_doi"`
}

// Links contains zenodo links
type Links struct {
	Bucket             string `json:"bucket"`
	Discard            string `json:"discard"`
	Edit               string `json:"edit"`
	Files              string `json:"files"`
	Html               string `json:"html"`
	LatestDraft        string `json:"latest_draft"`
	LasestDraftHtml    string `json:"latest_draft_html"`
	NewVersion         string `json:"newversion"`
	Publish            string `json:"publish"`
	RegisterConceptDoi string `json:"registerconceptdoi"`
	Self               string `json:"self"`
}

// AddResponse represents output of /create API
type AddResponse struct {
	Created  string `json:"created"`
	Modified string `json:"modified"`
	Size     int64  `json:"size"`
	Key      string `json:"key"`
	MimeType string `json:"mimetype"`
	Checksum string `json:"checksum"`
	Owner    int    `json:"owner"`
	RecordId int64  `json:"record_id"`
	Links    Links  `json:"links"`
}

// File describes common file record
type File struct {
	Name string `json:"name"`
	File string `json:"file"`
}

// DoiRecord represents doi record
type DoiRecord struct {
	Id     int64  `json:"id"`
	Doi    string `json:"doi"`
	DoiUrl string `json:"doi_url"`
	Files  []File `json:"files,omitempty"`
	Links  Links  `json:"links"`
}
