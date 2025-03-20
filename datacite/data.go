package datacite

import "time"

/*
 * https://support.datacite.org/docs/api-create-dois
 */

// RequestPayload represents request payload
type RequestPayload struct {
	Data RequestData `json:"data"`
}

// RequestData represents request data
type RequestData struct {
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

// ResponseData represents response payload
type ResposeData struct {
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	Attributes Attributes    `json:"attributes"`
	Relations  Relationships `json:"relationships"`
}

// RelatedIdentifier represents related identifier meta-data
type RelatedIdentifier struct {
	SchemaUri             string `json:"schemaUri"`
	RelationType          string `json:"relationType"`
	RelatedIdentifier     string `json:"relatedIdentifier"`
	RelatedIdentifierType string `json:"relatedIdentifierType"`
	RelatedMetadataScheme string `json:"relatedMetadataScheme"`
}

// Attributes represent attributes
type Attributes struct {
	Doi                  string              `json:"doi"`
	Prefix               string              `json:"prefix"`
	Event                string              `json:"event"`
	Suffix               string              `json:"suffix"`
	Identifiers          []string            `json:"identifiers"`
	AlternateIdentifiers []string            `json:"alternateIdentifiers"`
	Creators             []Creator           `json:"creators"`
	Titles               []Title             `json:"titles"`
	Publisher            string              `json:"publisher"`
	Container            Container           `json:"container"`
	PublicationYear      int                 `json:"publicationYear"`
	Subjects             []string            `json:"subjects"`
	Contributors         []Contributor       `json:"contributors"`
	Dates                []Date              `json:"dates"`
	Language             interface{}         `json:"language"`
	Types                Types               `json:"types"`
	RelatedIdentifiers   []RelatedIdentifier `json:"relatedIdentifiers"`
	RelatedItems         []string            `json:"relatedItems"`
	Sizes                []string            `json:"sizes"`
	Formats              []string            `json:"formats"`
	Version              interface{}         `json:"version"`
	RightsList           []string            `json:"rightsList"`
	Descriptions         []string            `json:"descriptions"`
	GeoLocations         []string            `json:"geoLocations"`
	FundingReferences    []string            `json:"fundingReferences"`
	XML                  string              `json:"xml"`
	URL                  string              `json:"url"`
	ContentURL           interface{}         `json:"contentUrl"`
	MetadataVersion      int                 `json:"metadataVersion"`
	SchemaVersion        interface{}         `json:"schemaVersion"`
	Source               string              `json:"source"`
	IsActive             bool                `json:"isActive"`
	State                string              `json:"state"`
	Reason               interface{}         `json:"reason"`
	LandingPage          interface{}         `json:"landingPage"`
	ViewCount            int                 `json:"viewCount"`
	ViewsOverTime        []string            `json:"viewsOverTime"`
	DownloadCount        int                 `json:"downloadCount"`
	DownloadsOverTime    []string            `json:"downloadsOverTime"`
	ReferenceCount       int                 `json:"referenceCount"`
	CitationCount        int                 `json:"citationCount"`
	CitationsOverTime    []string            `json:"citationsOverTime"`
	PartCount            int                 `json:"partCount"`
	PartOfCount          int                 `json:"partOfCount"`
	VersionCount         int                 `json:"versionCount"`
	VersionOfCount       int                 `json:"versionOfCount"`
	Created              time.Time           `json:"created"`
	Registered           time.Time           `json:"registered"`
	Published            string              `json:"published"`
	Updated              time.Time           `json:"updated"`

	MetaData any `json:"metaData"`
}

// Creator represents creator struct
type Creator struct {
	Name            string   `json:"name"`
	Affiliation     []string `json:"affiliation"`
	NameIdentifiers []string `json:"nameIdentifiers"`
}

// Title represents title
type Title struct {
	Title string `json:"title"`
}

// Container represents container
type Container struct {
}

// Affiliation represents affiliation
type Affiliation struct {
	AffiliationIdentifier       string `json:"affiliationIdentifier"`
	AffiliationIdentifierScheme string `json:"affiliationIdentifierScheme"`
	Name                        string `json:"name"`
	SchemeUri                   string `json:"schemeUri"`
}

// Contributor represents contributor
type Contributor struct {
	Name        string
	NameType    string
	GivenName   string
	FamilyName  string
	Affiliation []Affiliation
}

// Date represents date
type Date struct {
	Date            string `json:"dataType"`
	DataType        string `json:"dataType"`
	DateInformation string `json:"dateInformation"`
}

// Types represents types
type Types struct {
	SchemaOrg           string `json:"schemaOrg"`
	Citeproc            string `json:"citeproc"`
	Bibtex              string `json:"bibtex"`
	RIS                 string `json:"ris"`
	ResourceTypeGeneral string `json:"resourceTypeGeneral"`
}

// Relationships represents relationships
type Relationships struct {
	Client     RelationData `json:"client"`
	Provider   RelationData `json:"provider"`
	Media      RelationData `json:"media"`
	References RelationData `json:"references"`
	Citations  RelationData `json:"citations"`
	Parts      RelationData `json:"parts"`
	PartOf     RelationData `json:"partOf"`
	Versions   RelationData `json:"versions"`
	VersionOf  RelationData `json:"versionOf"`
}

// RelationData represents relation data
type RelationData struct {
	Data RelationInfo `json:"data"`
}

// RelationInfo represents relation info
type RelationInfo struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
