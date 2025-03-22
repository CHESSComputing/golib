package datacite

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
	utils "github.com/CHESSComputing/golib/utils"
)

/*
 * Datacite documentation:
 * https://support.datacite.org/reference/post_dois
 * https://support.datacite.org/docs/api-create-dois
 * https://schema.datacite.org/meta/kernel-4/
 */

// helper function to provide DOI creators
func creators() []Creator {
	return []Creator{
		Creator{
			Name:        "FOXDEN",
			Affiliation: []string{"Cornell University"},
		},
	}
}

// helper function to provide DOI Types
func types() Types {
	return Types{
		RIS:                 "FOXDEN",
		Bibtex:              "misc",
		SchemaOrg:           "dataset",
		ResourceTypeGeneral: "dataset",
	}
}

// helper function to publish foxden metadata in FOXDEN DOI service
func publishFoxdenRecord(record map[string]any) ([]RelatedIdentifier, error) {
	// publish given record in DOIService and obtain its URL
	schema, url, err := utils.Publish2DOIService(record)
	if err != nil {
		log.Println("ERROR: fail to obtain DOIService url, error", err)
		return []RelatedIdentifier{}, err
	}
	if url == "" {
		log.Println("ERROR: empty DOIService url")
		return []RelatedIdentifier{}, errors.New("fail to obtain DOIService url")
	}
	out := []RelatedIdentifier{
		RelatedIdentifier{
			SchemaUri:             schema,
			RelationType:          "FOXDEN metadata",
			RelatedIdentifier:     url,
			RelatedIdentifierType: "URL",
			RelatedMetadataScheme: "foxden",
		},
	}
	return out, nil
}

// DataCiteMetadata provides datacite metadata record for given did and FOXDEN record
func DataCiteMetadata(did, description string, record map[string]any, publish bool) ([]byte, error) {
	foxdenMeta, err := publishFoxdenRecord(record)
	if err != nil {
		log.Println("ERROR: fail to publish foxden record", err)
		return []byte{}, fmt.Errorf("failed to publish foxden record into DOIService: %v", err)
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
		Creators:           creators(),
		Publisher:          "Cornell University",
		PublicationYear:    time.Now().Year(),
		Descriptions:       []string{description},
		Types:              types(),
		RelatedIdentifiers: foxdenMeta,
		URL:                srvConfig.Config.Services.DOIServiceURL,
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
