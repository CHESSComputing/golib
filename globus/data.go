package globus

// GlobusTokenResponse represents globus token response
type GlobusTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// GlobusEndpointResponse represents globus endpoint response
type GlobusEndpointResponse struct {
	// Define the structure based on the response from the API
	Endpoints []struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
		Owner       string `json:"owner_string"`
		Description string `json:"description"`
	} `json:"DATA"`
}

// GlobusFileListResponse represents globus file list response
type GlobusFileListResponse struct {
	// Define the structure based on the response from the API
	Files []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"DATA"`
}

// GlobusSearchResponse represents globus search response
type GlobusSearchResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Owner       string `json:"owner_string"`
	Description string `json:"description"`
}
