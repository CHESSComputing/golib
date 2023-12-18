package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// Token represents response from OAuth server call
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `token_type`
	Expires     int64  `expires_in`
}

// HttpRequest manage http requests
type HttpRequest struct {
	Token   string
	Scope   string
	Expires time.Time
	Verbose int
}

// NewHttpRequest initilizes and returns new HttpRequest object
func NewHttpRequest(scope string, verbose int) *HttpRequest {
	return &HttpRequest{Scope: scope, Verbose: verbose}
}

// GetToken obtains token from OAuth server
func (h *HttpRequest) GetToken() {
	if h.Token == "" || h.Expires.Before(time.Now()) {
		// make a call to Authz service to obtain access token
		rurl := fmt.Sprintf(
			"%s/oauth/token?client_id=%s&response&client_secret=%s&grant_type=client_credentials&scope=%s",
			srvConfig.Config.Services.AuthzURL,
			srvConfig.Config.Authz.ClientID,
			srvConfig.Config.Authz.ClientSecret,
			h.Scope)
		resp, err := h.Get(rurl)
		if err != nil {
			log.Println("ERROR", err)
			return
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		var response Token
		err = json.Unmarshal(data, &response)
		if err != nil {
			log.Println("ERROR", err)
			return
		}
		if h.Verbose > 0 {
			log.Printf("INFO: Authz response %+v, error %v", response, err)
		}
		h.Token = response.AccessToken
		h.Expires = time.Now().Add(time.Duration(response.Expires) * time.Second)
	}
}

// Get performis HTTP GET request
func (h *HttpRequest) Get(rurl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", rurl, nil)
	if err != nil {
		return nil, err
	}
	if h.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.Token))
	}
	req.Header.Add("Accept-Encoding", "")
	client := &http.Client{}
	if h.Verbose > 1 {
		dump, err := httputil.DumpRequestOut(req, true)
		log.Println("request", string(dump), err)
	}
	resp, err := client.Do(req)
	if h.Verbose > 1 {
		dump, err := httputil.DumpResponse(resp, true)
		log.Println("response", string(dump), err)
	}
	return resp, err
}

// Post performs HTTP POST request
func (h *HttpRequest) Post(rurl, contentType string, buffer *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest("POST", rurl, buffer)
	if err != nil {
		return nil, err
	}
	if h.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.Token))
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept", contentType)
	client := &http.Client{}
	if h.Verbose > 1 {
		dump, err := httputil.DumpRequestOut(req, true)
		log.Println("request", string(dump), err)
	}
	resp, err := client.Do(req)
	if h.Verbose > 1 {
		dump, err := httputil.DumpResponse(resp, true)
		log.Println("response", string(dump), err)
	}
	return resp, err
}

// PostForm perform HTTP POST form request with bearer token
func (h *HttpRequest) PostForm(rurl string, formData url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", rurl, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	if h.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.Token))
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	if h.Verbose > 1 {
		dump, err := httputil.DumpRequestOut(req, true)
		log.Println("request", string(dump), err)
	}
	resp, err := client.Do(req)
	if h.Verbose > 1 {
		dump, err := httputil.DumpResponse(resp, true)
		log.Println("response", string(dump), err)
	}
	return resp, err
}

// helper function to return full path of given file name wrt to current location
func fullPath(fname string) string {
	if !strings.HasPrefix(fname, "/") {
		// we got relative path (e.g. server_test.json)
		if wdir, err := os.Getwd(); err == nil {
			fname = filepath.Join(wdir, fname)
		}
	}
	return fname
}

// SchemaFileName obtains schema file name from schema name
func SchemaFileName(sname string) string {
	var fname string
	for _, f := range srvConfig.Config.CHESSMetaData.SchemaFiles {
		if strings.Contains(f, sname) {
			fname = f
			break
		}
	}
	return fullPath(fname)
}

// SchemaName extracts schema name from schema file name
func SchemaName(fname string) string {
	arr := strings.Split(fname, "/")
	return strings.Split(arr[len(arr)-1], ".")[0]
}
