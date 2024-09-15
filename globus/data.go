package globus

type GlobusTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type GlobusEndpointResponse struct {
	// Define the structure based on the response from the API
	Endpoints []struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"DATA"`
}

type GlobusFileListResponse struct {
	// Define the structure based on the response from the API
	Files []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"DATA"`
}
