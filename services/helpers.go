package services

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
	"strings"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
	utils "github.com/CHESSComputing/golib/utils"
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
	Headers map[string][]string
}

// NewHttpRequest initilizes and returns new HttpRequest object
func NewHttpRequest(scope string, verbose int) *HttpRequest {
	return &HttpRequest{Scope: scope, Verbose: verbose}
}

// Response returns service status record
func Response(srv string, httpCode, srvCode int, err error) ServiceResponse {
	status := "error"
	if err == nil {
		status = "ok"
	}
	if status == "error" {
		log.Printf("ERROR: http code %d srv code %d error %v\n %v", httpCode, srvCode, err, utils.Stack())
	}
	var strError string
	if err != nil {
		strError = err.Error()
	}
	return ServiceResponse{
		HttpCode:  httpCode,
		Service:   srv,
		Status:    status,
		Error:     strError,
		SrvCode:   srvCode,
		Timestamp: time.Now().String(),
	}
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
		if os.Getenv("FOXDEN_DEBUG") != "" || h.Verbose > 2 {
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
	req.Header.Add("Accept", "application/json")
	//     req.Header.Add("Accept-Encoding", "")
	if h.Headers != nil {
		for key, values := range h.Headers {
			for _, val := range values {
				req.Header.Add(key, val)
			}
		}
	}
	// check if we are given Zenodo request
	if strings.Contains(rurl, "zenodo.org") {
		if vals, ok := h.Headers["ZenodoAccessToken"]; ok {
			if len(vals) == 1 && vals[0] != h.Token && vals[0] != "" {
				// we overwrite access bearer value with provided Zenodo access token
				req.Header.Del("Authorization")
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", vals[0]))
			}
		}
	}
	if os.Getenv("FOXDEN_DEBUG") != "" || h.Verbose > 2 {
		log.Printf("HTTP GET request to %s", rurl)
	}
	client := &http.Client{}
	if os.Getenv("FOXDEN_DEBUG") != "" || h.Verbose > 2 {
		dump, err := httputil.DumpRequestOut(req, true)
		log.Println("HttpRequest: GET request", string(dump), err)
	}
	resp, err := client.Do(req)
	if os.Getenv("FOXDEN_DEBUG") != "" || h.Verbose > 2 {
		dump, err := httputil.DumpResponse(resp, true)
		log.Println("HttpRequest: GET response", string(dump), err)
	}
	return resp, err
}

// Post performs HTTP POST request
func (h *HttpRequest) Post(rurl, contentType string, buffer *bytes.Buffer) (*http.Response, error) {
	return h.Request("POST", rurl, contentType, buffer)
}

// Put performs HTTP PUT request
func (h *HttpRequest) Put(rurl, contentType string, buffer *bytes.Buffer) (*http.Response, error) {
	return h.Request("PUT", rurl, contentType, buffer)
}

// Delete performs HTTP PUT request
func (h *HttpRequest) Delete(rurl, contentType string, buffer *bytes.Buffer) (*http.Response, error) {
	return h.Request("DELETE", rurl, contentType, buffer)
}

// Request performs HTTP request for given method
func (h *HttpRequest) Request(method, rurl, contentType string, buffer *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest(method, rurl, buffer)
	if err != nil {
		return nil, err
	}
	if h.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.Token))
	}
	// get all HTTP headers from HttpRequest one
	if h.Headers != nil {
		for key, values := range h.Headers {
			for _, val := range values {
				req.Header.Add(key, val)
			}
		}
	}
	// setup default Content-Type and Accept if they are not set above
	if _, ok := req.Header["Content-Type"]; !ok {
		req.Header.Add("Content-Type", contentType)
	}
	if _, ok := req.Header["Accept"]; !ok {
		// we re-use contentType for accept header as it is usually application/json
		// if accept and content-type are different then it is better to setup them
		// through explicit map of HttpRequest.Headers
		req.Header.Add("Accept", contentType)
	}
	// check if we are given Zenodo request
	if strings.Contains(rurl, "zenodo.org") {
		if vals, ok := h.Headers["ZenodoAccessToken"]; ok {
			if len(vals) == 1 && vals[0] != h.Token && vals[0] != "" {
				// we overwrite access bearer value with provided Zenodo access token
				req.Header.Del("Authorization")
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", vals[0]))
			}
		}
	}
	if os.Getenv("FOXDEN_DEBUG") != "" || h.Verbose > 2 {
		log.Printf("HTTP %s request to %s", method, rurl)
	}
	client := &http.Client{}
	if os.Getenv("FOXDEN_DEBUG") != "" || h.Verbose > 2 {
		dump, err := httputil.DumpRequestOut(req, true)
		log.Printf("HttpRequest: %s request %s, error %v", method, string(dump), err)
	}
	resp, err := client.Do(req)
	if os.Getenv("FOXDEN_DEBUG") != "" || h.Verbose > 2 {
		dump, err := httputil.DumpResponse(resp, true)
		log.Printf("HttpRequest: method %s response %s, error %v", method, string(dump), err)
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
	if os.Getenv("FOXDEN_DEBUG") != "" || h.Verbose > 2 {
		log.Printf("HTTP POST request to %s", rurl)
	}
	client := &http.Client{}
	if os.Getenv("FOXDEN_DEBUG") != "" || h.Verbose > 2 {
		dump, err := httputil.DumpRequestOut(req, true)
		log.Println("HttpRequest: POST form request", string(dump), err)
	}
	resp, err := client.Do(req)
	if os.Getenv("FOXDEN_DEBUG") != "" || h.Verbose > 2 {
		dump, err := httputil.DumpResponse(resp, true)
		log.Println("HttpRequest: POST form response", string(dump), err)
	}
	return resp, err
}
