package datacite

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
)

/*
 * Datacite documentation:
 * https://support.datacite.org/reference/post_dois
 * https://support.datacite.org/docs/api-create-dois
 * https://schema.datacite.org/meta/kernel-4/
 */

// helper function to provide DOI creators
func creatorsInfo() []Creator {
	nameIdent := NameIdentifier{
		AffiliationIdentifier:       "https://ror.org/05bnh6r87",
		AffiliationIdentifierScheme: "ROR",
		SchemeUri:                   "https://ror.org/",
	}
	return []Creator{
		Creator{
			Name:            "FOXDEN",
			NameType:        "Organizational",
			NameIdentifiers: []NameIdentifier{nameIdent},
		},
	}
}

// helper function to provide DOI Types, return pointer since our Attributes.Types is a pointer
// type (to ensure it will be ommitted if nil)
func typesInfo() *Types {
	return &Types{
		ResourceType:        "FOXDEN",
		ResourceTypeGeneral: "Dataset",
	}
}

// helper function to provide DOI Publisher, return pointer since our Attributes.Publisher is a pointer
// type (to ensure it will be ommitted if nil)
func publisherInfo() *Publisher {
	return &Publisher{
		Name:                      "DataCite",
		PublisherIdentifier:       "https://ror.org/04wxnsj81",
		PublisherIdentifierScheme: "ROR",
		SchemeUri:                 "https://ror.org/",
		Lang:                      "en",
	}
}

func descriptionInfo(d string) Description {
	return Description{
		Description:     d,
		DescriptionType: "Other",
		Lang:            "en",
	}
}

// DataciteMetadata provides datacite metadata record for given did and FOXDEN record
func DataciteMetadata(doi, did, description string, record map[string]any, publish bool) ([]byte, error) {
	url := srvConfig.Config.Services.DOIServiceURL
	if doi != "" {
		url += "/doi" // http://DOIServiceURL/doi/<10.xxx/...>
		url = filepath.Join(url, doi)
	}
	relatedIds := []RelatedIdentifier{
		RelatedIdentifier{
			RelationType:          "HasMetadata",
			RelatedIdentifier:     url,
			RelatedTypeGeneral:    "Dataset",
			RelatedIdentifierType: "URL",
		},
	}
	event := ""
	if publish {
		event = "publish"
	}

	title := Title{Title: fmt.Sprintf("FOXDEN did=%s", did)}
	attrs := Attributes{
		Event:              event,
		Titles:             []Title{title},
		Prefix:             srvConfig.Config.DOI.Datacite.Prefix,
		Creators:           creatorsInfo(),
		Publisher:          publisherInfo(),
		PublicationYear:    time.Now().Year(),
		Descriptions:       []Description{descriptionInfo(description)},
		Types:              typesInfo(),
		RelatedIdentifiers: relatedIds,
		URL:                url,
	}

	// Convert payload to JSON
	payload := RequestPayload{
		Data: RequestData{
			Type:       "dois",
			Attributes: attrs,
		},
	}
	payloadBytes, err := json.MarshalIndent(payload, "", "   ")
	return payloadBytes, err
}
