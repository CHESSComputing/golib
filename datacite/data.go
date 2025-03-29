package datacite

/*
 * This module provides all metadata definitions for Datacite
 * References:
 * - https://support.datacite.org/docs/api-create-dois
 * - https://datacite-metadata-schema.readthedocs.io/en/4.6/properties/overview/
 */

// RequestPayload represents request payload
type RequestPayload struct {
	Data RequestData `json:"data,omitempty"`
}

// RequestData represents request data
type RequestData struct {
	Type       string     `json:"type,omitempty"`
	Attributes Attributes `json:"attributes,omitempty"`
}

// ResponsePayload represents request payload
type ResponsePayload struct {
	Data ResponseData `json:"data,omitempty"`
}

// ResponseData represents response payload
// NOTE: the Attributes of response is slightly different from Request Attributes,
// for instance Publisher is a string instead of Publisher struct in Request Attributes.
type ResponseData struct {
	ID         string             `json:"id,omitempty"`
	Type       string             `json:"type,omitempty"`
	Attributes ResponseAttributes `json:"attributes,omitempty"`
}

// RelatedIdentifier represents related identifier meta-data
type RelatedIdentifier struct {
	RelationType          string `json:"relationType,omitempty"`
	RelatedTypeGeneral    string `json:"relatedTypeGeneral,omitempty"`
	RelatedIdentifier     string `json:"relatedIdentifier,omitempty"`
	RelatedIdentifierType string `json:"relatedIdentifierType,omitempty"`
}

// Publisher represents publisher info
type Publisher struct {
	Name                      string `json:"name,omitempty"`
	PublisherIdentifier       string `json:"publisherIdentifier,omitempty"`
	PublisherIdentifierScheme string `json:"publisherIdentifier,omitempty"`
	SchemeUri                 string `json:"schemeUri,omitempty"`
	Lang                      string `json:"lang,omitempty"`
}

// Description represents description info
type Description struct {
	Description     string `json:"description,omitempty"`
	DescriptionType string `json:"descriptionType,omitempty"`
	Lang            string `json:"lang,omitempty"`
}

// Attributes represent attributes
// NOTE1: we use pointers to Publisher and Types structs to ensure that they will not be included if nil
// NOTE2: we don't need points for lists since they will be properly omitted
type Attributes struct {
	Doi                string              `json:"doi,omitempty"`
	Prefix             string              `json:"prefix,omitempty"`
	Event              string              `json:"event,omitempty"`
	Creators           []Creator           `json:"creators,omitempty"`
	Titles             []Title             `json:"titles,omitempty"`
	Publisher          *Publisher          `json:"publisher,omitempty"` // omitempty ensures it's not included if nil
	PublicationYear    int                 `json:"publicationYear,omitempty"`
	Types              *Types              `json:"types,omitempty"` // omitempty ensures it's not included if nil
	RelatedIdentifiers []RelatedIdentifier `json:"relatedIdentifiers,omitempty"`
	Descriptions       []Description       `json:"descriptions,omitempty"`
	URL                string              `json:"url,omitempty"`
}

// ResponseAttributes represent attributes
type ResponseAttributes struct {
	Doi                string              `json:"doi,omitempty"`
	Prefix             string              `json:"prefix,omitempty"`
	Event              string              `json:"event,omitempty"`
	Creators           []Creator           `json:"creators,omitempty"`
	Titles             []Title             `json:"titles,omitempty"`
	Publisher          string              `json:"publisher,omitempty"`
	PublicationYear    int                 `json:"publicationYear,omitempty"`
	Types              Types               `json:"types,omitempty"`
	State              string              `json:"state,omitempty"`
	RelatedIdentifiers []RelatedIdentifier `json:"relatedIdentifiers,omitempty"`
	Descriptions       []Description       `json:"descriptions,omitempty"`
	URL                string              `json:"url,omitempty"`
}

// Creator represents creator struct
type Creator struct {
	Name            string           `json:"name,omitempty"`
	NameType        string           `json:"nameType,omitempty"`
	NameIdentifiers []NameIdentifier `json:"nameIdentifiers,omitempty"`
}

// NameIdentifier represents name identifier info
type NameIdentifier struct {
	AffiliationIdentifier       string `json:"affiliationIdentifier,omitempty"`
	AffiliationIdentifierScheme string `json:"affiliationIdentifierScheme,omitempty"`
	SchemeUri                   string `json:"schemeUri,omitempty"`
}

// Title represents title
type Title struct {
	Title string `json:"title,omitempty"`
}

// Types represents types
type Types struct {
	ResourceType        string `json:"resourceType,omitempty"`
	ResourceTypeGeneral string `json:"resourceTypeGeneral,omitempty"`
}
