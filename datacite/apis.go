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
func Publish(did, description string, record map[string]any, publish bool) (string, string, error) {
	if srvConfig.Config.Services.DOIServiceURL == "" {
		return "", "", errors.New("Missing DOIService url in FOXDEN configuration")
	}

	// Set the DOI creation endpoint
	url := fmt.Sprintf("%s/dois", srvConfig.Config.DOI.Datacite.Url)

	payloadBytes, err := DataCiteMetadata(did, description, record, publish)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal metadata payload: %v", err)
	}
	log.Println("Publish\n", string(payloadBytes))

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
