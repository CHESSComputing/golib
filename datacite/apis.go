package datacite

import (
	"bytes"
	"encoding/json"
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
	return Creator{
		Name:        "FOXDEN",
		Affiliation: "Cornell University",
	}
}

func persistentUrl() string {
	return srvConfig.Config.DOI.Datacite.LandingPageUrl
}

func prefix() string {
	if srvConfig.Config.DOI.Datacite.Url == "https://api.test.datacite.org/dois" {
		return "10.5438"
	}
	return ""
}

func types() []Types {
	t := Types{
		ResourceTypeGeneral: "dataset",
	}
	var out []Types
	out = append(out, t)
	return out
}

// Publish provides publication of did into datacite
func Publish(did, description string, record any) (string, string, error) {
	var doi, doiLink string
	var err error

	title := Title{Title: fmt.Sprintf("FOXDEN did=%s", did)}
	attrs := Attributes{
		Titles:          []Title{title},
		Prefix:          prefix(),
		Creators:        creators(),
		Publisher:       "Cornell University",
		PublicationYear: time.Time.Year(),
		Descriptions:    []string{description},
		Types:           types(),
		MetaData:        record,
		Url:             persistentUrl(),
	}

	// Set the DOI creation endpoint
	url := fmt.Sprintf("%s/dois", srvConfig.Config.DOI.Datacite.Url)

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(attrs)
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
		req.Header.Set("Authorization", fmt.Println("Bearer %s", srvConfig.Config.DOI.Datacite.AccessToken))
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

	return doi, doiLink, err
}
