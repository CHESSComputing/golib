package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// HttpGet performis HTTP GET request with bearer token
func HttpGet(token, rurl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", rurl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept-Encoding", "")
	client := &http.Client{}
	if srvConfig.Config.Frontend.WebServer.Verbose > 1 {
		dump, err := httputil.DumpRequestOut(req, true)
		log.Println("request", string(dump), err)
	}
	resp, err := client.Do(req)
	if srvConfig.Config.Frontend.WebServer.Verbose > 1 {
		dump, err := httputil.DumpResponse(resp, true)
		log.Println("response", string(dump), err)
	}
	return resp, err
}

// HttpPost performs HTTP POST request with bearer token
func HttpPost(token, rurl, contentType string, buffer *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest("POST", rurl, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept", contentType)
	client := &http.Client{}
	if srvConfig.Config.Frontend.WebServer.Verbose > 1 {
		dump, err := httputil.DumpRequestOut(req, true)
		log.Println("request", string(dump), err)
	}
	resp, err := client.Do(req)
	if srvConfig.Config.Frontend.WebServer.Verbose > 1 {
		dump, err := httputil.DumpResponse(resp, true)
		log.Println("response", string(dump), err)
	}
	return resp, err
}

// HttpPostForm perform HTTP POST form request with bearer token
func HttpPostForm(token, rurl string, formData url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", rurl, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	if srvConfig.Config.Frontend.WebServer.Verbose > 1 {
		dump, err := httputil.DumpRequestOut(req, true)
		log.Println("request", string(dump), err)
	}
	resp, err := client.Do(req)
	if srvConfig.Config.Frontend.WebServer.Verbose > 1 {
		dump, err := httputil.DumpResponse(resp, true)
		log.Println("response", string(dump), err)
	}
	return resp, err
}
