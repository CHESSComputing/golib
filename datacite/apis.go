package datacite

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// Publish provides publication of did into datacite
func Publish(did, description string, record map[string]any, publish bool, verbose int) (string, string, error) {
	if srvConfig.Config.Services.DOIServiceURL == "" {
		return "", "", errors.New("Missing DOIService url in FOXDEN configuration")
	}

	// Set the DOI creation endpoint
	url := fmt.Sprintf("%s/dois", srvConfig.Config.DOI.Datacite.Url)
	payloadBytes, err := DataciteMetadata("", did, description, record, publish)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal metadata payload: %v", err)
	}
	if verbose > 1 {
		log.Println("Publish\n", string(payloadBytes))
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

// GetRecord fetches DOI record for given doi
func GetRecord(doi string, verbose int) (ResponsePayload, error) {
	var record ResponsePayload
	rurl := fmt.Sprintf("%s/dois/%s", srvConfig.Config.DOI.Datacite.Url, doi)
	req, err := http.NewRequest("GET", rurl, nil)
	if err != nil {
		log.Println("ERROR: unable to create PUT request", err)
		return record, err
	}

	// Set authentication and headers
	if srvConfig.Config.DOI.Datacite.Username != "" && srvConfig.Config.DOI.Datacite.Password != "" {
		req.SetBasicAuth(srvConfig.Config.DOI.Datacite.Username, srvConfig.Config.DOI.Datacite.Password)
	}
	if srvConfig.Config.DOI.Datacite.AccessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", srvConfig.Config.DOI.Datacite.AccessToken))
	}
	req.Header.Set("Accept", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR: unable to make HTTP request", err)
		return record, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if verbose > 2 {
		log.Println("receive", string(body), err)
	}
	if err != nil {
		log.Println("ERROR: unable to read HTTP response", err)
		return record, err
	}
	err = json.Unmarshal(body, &record)
	return record, err
}

// MakePublic implements logic of publishing draft DOI
// curl -X PUT -H "Content-Type: application/vnd.api+json" --user YOUR_REPOSITORY_ID:YOUR_PASSWORD -d @my_doi_update.json https://api.test.datacite.org/dois/:id
func MakePublic(doi string, verbose int) error {
	// first we should check if we already has Public DOI
	if rec, err := GetRecord(doi, verbose); err == nil {
		if rec.Data.Attributes.State == "findable" {
			log.Printf("WARNING: our doi record %s is already findable", doi)
			return nil
		}
	}

	// update record
	rurl := fmt.Sprintf("%s/%s", srvConfig.Config.Services.DOIServiceURL, doi)

	rid := RelatedIdentifier{
		RelationType:          "HasMetadata",
		RelatedIdentifier:     rurl,
		RelatedTypeGeneral:    "Dataset",
		RelatedIdentifierType: "URL",
	}

	attrs := Attributes{
		RelatedIdentifiers: []RelatedIdentifier{rid},
		Types:              &Types{ResourceType: "FOXDEN", ResourceTypeGeneral: "Dataset"},
		Event:              "publish",
		URL:                rurl,
	}

	// Convert payload to JSON
	payload := RequestPayload{
		Data: RequestData{
			Type:       "dois",
			Attributes: attrs,
		},
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if verbose > 1 {
		log.Printf("update DOI record %+v", string(data))
	}
	if err != nil {
		log.Println("ERROR: unable to create JSON payload", err)
		return err
	}
	// make HTTP PUT request to update our record
	url := fmt.Sprintf("%s/dois/%s", srvConfig.Config.Datacite.Url, doi)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		log.Println("ERROR: unable to create PUT request", err)
		return err
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
		log.Println("ERROR: unable to make HTTP request", err)
		return err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERROR: unable to read HTTP response", err)
		return err
	}
	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("unable to update DOI document, status %v response %+v", resp.StatusCode, string(body))
		log.Println("ERROR:", msg)
		return errors.New(msg)
	}
	return nil

}
