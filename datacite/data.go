package datacite

/*
 * This module provides all metadata definitions for Datacite
 * References:
 * - https://support.datacite.org/docs/api-create-dois
 * - https://datacite-metadata-schema.readthedocs.io/en/4.6/properties/overview/
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
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

// RelatedIdentifier represents related identifier meta-data
type RelatedIdentifier struct {
	RelationType          string `json:"relationType"`
	RelatedTypeGeneral    string `json:"relatedTypeGeneral"`
	RelatedIdentifier     string `json:"relatedIdentifier"`
	RelatedIdentifierType string `json:"relatedIdentifierType"`
}

// Publisher represents publisher info
type Publisher struct {
	Name                      string `json:"name"`
	PublisherIdentifier       string `json:"publisherIdentifier"`
	PublisherIdentifierScheme string `json:"publisherIdentifier"`
	SchemeUri                 string `json:"schemeUri"`
	Lang                      string `json:"lang"`
}

// Description represents description info
type Description struct {
	Description     string `json:"description"`
	DescriptionType string `json:"descriptionType"`
	Lang            string `json:"lang"`
}

// Attributes represent attributes
type Attributes struct {
	Doi                string              `json:"doi"`
	Prefix             string              `json:"prefix"`
	Event              string              `json:"event"`
	Creators           []Creator           `json:"creators"`
	Titles             []Title             `json:"titles"`
	Publisher          Publisher           `json:"publisher"`
	PublicationYear    int                 `json:"publicationYear"`
	Types              Types               `json:"types"`
	RelatedIdentifiers []RelatedIdentifier `json:"relatedIdentifiers"`
	Descriptions       []Description       `json:"descriptions"`
	URL                string              `json:"url"`
}

// Creator represents creator struct
type Creator struct {
	Name            string           `json:"name"`
	NameType        string           `json:"nameType"`
	NameIdentifiers []NameIdentifier `json:"nameIdentifiers"`
}

// NameIdentifier represents name identifier info
type NameIdentifier struct {
	AffiliationIdentifier       string `json:"affiliationIdentifier"`
	AffiliationIdentifierScheme string `json:"affiliationIdentifierScheme"`
	SchemeUri                   string `json:"schemeUri"`
}

// Title represents title
type Title struct {
	Title string `json:"title"`
}

// Types represents types
type Types struct {
	ResourceType        string `json:"resourceType"`
	ResourceTypeGeneral string `json:"resourceTypeGeneral"`
}
