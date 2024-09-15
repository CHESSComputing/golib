package globus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// Token API obtains Globus access token
// example through curl
// curl -X POST https://auth.globus.org/v2/oauth2/token --header "Content-Type: application/x-www-form-urlencoded" --data-urlencode "grant_type=client_credentials" --data-urlencode "scope=$scope" --data-urlencode "client_id=$clientid" --data-urlencode "client_secret=$secret"
func Token(scope string) (string, error) {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	data := fmt.Sprintf("grant_type=client_credentials&scope=%s&client_id=%s&client_secret=%s",
		scope,
		srvConfig.Config.Globus.ClientID,
		srvConfig.Config.Globus.ClientSecret)
	req, err := http.NewRequest("POST", srvConfig.Config.Globus.AuthURL, strings.NewReader(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get token, status code: %d", resp.StatusCode)
	}

	var tokenResponse GlobusTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

// Search API provides search API to globus search pattern
// example of curl
// curl -H "Accept: application/json" -H "Authorization: Bearer $t" "https://transfer.api.globus.org/v0.10/endpoint_search?filter_fulltext=CHESS"
func Search(token string, pattern string) {
	url := fmt.Sprintf("%s/endpoint_search?filter_fulltext=%s", srvConfig.Config.Globus.TransferURL, pattern)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		return
	}

	var searchResponse GlobusEndpointResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	for _, endpoint := range searchResponse.Endpoints {
		fmt.Printf("Found endpoint: ID=%s, Name=%s\n", endpoint.ID, endpoint.DisplayName)
	}
}

// Ls API provides listing files within Globus
func Ls(token, endpointID, path string) {
	url := fmt.Sprintf("%s/operation/endpoint/%s/ls?path=%s", srvConfig.Config.Globus.TransferURL, endpointID, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		return
	}

	var fileList GlobusFileListResponse
	err = json.NewDecoder(resp.Body).Decode(&fileList)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	for _, file := range fileList.Files {
		fmt.Printf("File: %s (Type: %s)\n", file.Name, file.Type)
	}
}

// Mkdir provides API to create Globus directory
func Mkdir(token, endpointID, path string) error {
	url := fmt.Sprintf("%s/operation/endpoint/%s/mkdir", srvConfig.Config.Globus.TransferURL, endpointID)
	payload := fmt.Sprintf(`{"path": "%s"}`, path)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create directory, status code: %d", resp.StatusCode)
	}

	fmt.Println("Directory created successfully")
	return nil
}

// Upload API
func Upload(token, endpointID, localFile, remotePath string) error {
	// You would typically initiate a transfer task to upload files using the transfer API
	// Upload logic would involve creating a transfer task to move the file from one endpoint to another
	// e.g., from your local machine endpoint to a remote Globus endpoint
	fmt.Println("Upload logic to be implemented based on transfer tasks")
	return nil
}

// Download API
func Download(token string, endpointID string, remotePath string) ([]byte, error) {
	// Similar to upload, download logic would involve transfer tasks via the Globus Transfer API
	// This typically requires setting up the transfer between endpoints
	fmt.Println("Download logic to be implemented based on transfer tasks")
	return nil, nil
}

// SharedLink provides shared link to Globus data
func SharedLink(token, endpointID, path string) (string, error) {
	url := fmt.Sprintf("%s/endpoint/%s/share", srvConfig.Config.Globus.TransferURL, endpointID)
	payload := fmt.Sprintf(`{"path": "%s"}`, path)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create shared link, status code: %d", resp.StatusCode)
	}

	var linkResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&linkResponse)
	if err != nil {
		return "", err
	}

	return linkResponse["link_url"].(string), nil
}
