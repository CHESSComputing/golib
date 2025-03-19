package datacite

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
)

/*
 * Datacite documentation:
 * https://support.datacite.org/reference/post_dois
 * https://support.datacite.org/docs/api-create-dois
 * https://schema.datacite.org/meta/kernel-4/
 */

func creators() []Creator {
	return []Creator{
		Creator{
			Name:        "FOXDEN",
			Affiliation: []string{"Cornell University"},
		},
	}
}

func types() Types {
	return Types{
		RIS:                 "FOXDEN",
		Bibtex:              "misc",
		SchemaOrg:           "dataset",
		ResourceTypeGeneral: "dataset",
	}
}

func publish2DOIService(rec any) (string, string, error) {
	var schema, url, did string
	var err error
	record := rec.(map[string]any)
	if val, ok := record["schema"]; ok {
		schema = val.(string)
	} else {
		err = errors.New("unable to look-up schema in FOXDEN record")
	}
	if val, ok := record["did"]; ok {
		did = val.(string)
	} else {
		err = errors.New("unable to look-up did in FOXDEN record")
	}
	url = fmt.Sprintf("%s/meta?did=%s", srvConfig.Config.Services.FrontendURL, did)
	return schema, url, err
}

// helper function to publish foxden metadata in FOXDEN DOI service
func publishFoxdenRecord(record any) ([]RelatedIdentifier, error) {
	// publish given record in DOIService and obtain its URL
	schema, url, err := publish2DOIService(record)
	if err != nil {
		return []RelatedIdentifier{}, err
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

// Publish provides publication of did into datacite
func Publish(did, description string, record any) (string, string, error) {
	foxdenMeta, err := publishFoxdenRecord(record)
	if err != nil {
		return "", "", fmt.Errorf("failed to publish foxden record into DOIService: %v", err)
	}

	title := Title{Title: fmt.Sprintf("FOXDEN did=%s", did)}
	attrs := Attributes{
		Titles:             []Title{title},
		Prefix:             srvConfig.Config.DOI.Datacite.Prefix,
		Creators:           creators(),
		Publisher:          "Cornell University",
		PublicationYear:    time.Now().Year(),
		Descriptions:       []string{description},
		Types:              types(),
		RelatedIdentifiers: foxdenMeta,
		URL:                srvConfig.Config.DOI.Datacite.LandingPageUrl,
	}

	// Set the DOI creation endpoint
	url := fmt.Sprintf("%s/dois", srvConfig.Config.DOI.Datacite.Url)

	// Convert payload to JSON
	payload := RequestPayload{
		Data: RequestData{
			Type:       "dois",
			Attributes: attrs,
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal metadata payload: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return "", "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set authentication and headers
	if srvConfig.Config.DOI.Datacite.Username != "" && srvConfig.Config.DOI.Datacite.Password != "" {
		req.SetBasicAuth(srvConfig.Config.DOI.Datacite.Username, srvConfig.Config.DOI.Datacite.Password)
	}
	if srvConfig.Config.DOI.Datacite.AccessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", srvConfig.Config.DOI.Datacite.AccessToken))
	}
	req.Header.Set("Content-Type", "application/vnd.api+json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("failed to create DOI: %s", respBody)
	}

	// Parse the response to extract the DOI
	var response map[string]any
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse response JSON: %v", err)
	}

	// Extract the DOI from the response
	doi, ok := response["data"].(map[string]any)["attributes"].(map[string]any)["doi"].(string)
	if !ok {
		return "", "", fmt.Errorf("failed to extract DOI from response")
	}
	doiLink := fmt.Sprintf("%s/dois/%s", srvConfig.Config.DOI.Datacite.Url, doi)

	return doi, doiLink, err
}
