package utils

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

// Verbose controls verbosity level of function printouts
var Verbose int

// ReadToken function to either read file content or return given string
func ReadToken(r string) string {
	if _, err := os.Stat(r); err == nil {
		b, e := os.ReadFile(r)
		if e != nil {
			log.Fatalf("Unable to read data from file: %s, error: %s", r, e)
		}
		return strings.Replace(string(b), "\n", "", -1)
	}
	return r
}

// HttpClient provides cert/token aware HTTP client
func HttpClient() *http.Client {
	timeout := time.Duration(TIMEOUT) * time.Second
	if TIMEOUT > 0 {
		return &http.Client{Timeout: time.Duration(timeout)}
	}
	return &http.Client{}
}

// HttpGet performs HTTP GET request with bearer token
func HttpGet(rurl string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", rurl, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	if Verbose > 1 {
		dump, err := httputil.DumpRequestOut(req, true)
		log.Println("request", string(dump), err)
	}
	resp, err := client.Do(req)
	if Verbose > 1 {
		dump, err := httputil.DumpResponse(resp, true)
		log.Println("response", string(dump), err)
	}
	return resp, err
}

// HttpPost performs HTTP POST request with bearer token
func HttpPost(rurl string, headers map[string]string, buffer *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest("POST", rurl, buffer)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	if Verbose > 1 {
		dump, err := httputil.DumpRequestOut(req, true)
		log.Println("request", string(dump), err)
	}
	resp, err := client.Do(req)
	if Verbose > 1 {
		dump, err := httputil.DumpResponse(resp, true)
		log.Println("response", string(dump), err)
	}
	return resp, err
}

// HttpPostForm performs HTTP POST form request with bearer token
func HttpPostForm(rurl string, headers map[string]string, formData url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", rurl, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	if Verbose > 1 {
		dump, err := httputil.DumpRequestOut(req, true)
		log.Println("request", string(dump), err)
	}
	resp, err := client.Do(req)
	if Verbose > 1 {
		dump, err := httputil.DumpResponse(resp, true)
		log.Println("response", string(dump), err)
	}
	return resp, err
}
