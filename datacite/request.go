package datacite

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var DATACITE_REPOSITORY_ID, DATACITE_PASSWORD string

func DOIRequest(url, payload string) error {
	// Define the API endpoint
	//     url := "https://api.test.datacite.org/dois"

	// Read the JSON file
	jsonData, err := ioutil.ReadFile(payload)
	if err != nil {
		log.Printf("ERROR: Failed to read JSON file: %v", err)
		return err
	}

	// Create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request: %v", err)
		return err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/vnd.api+json")

	// Add basic auth
	req.SetBasicAuth(DATACITE_REPOSITORY_ID, DATACITE_PASSWORD)

	// Send the request using the http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make the request: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read the response: %v", err)
		return err
	}

	// Print the status code and response body
	fmt.Printf("Response status code: %d\n", resp.StatusCode)
	fmt.Printf("Response body: %s\n", string(body))
	return nil
}
