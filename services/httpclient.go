package services

import (
	"bytes"
	"net/http"
	"net/url"
)

// HTTPClient defines generic interface of HTTP client
type HTTPClient interface {
	SetToken(token string)
	GetToken()
	Get(rurl string) (*http.Response, error)
	Post(rurl, contentType string, buffer *bytes.Buffer) (*http.Response, error)
	Put(rurl, contentType string, buffer *bytes.Buffer) (*http.Response, error)
	Delete(rurl, contentType string, buffer *bytes.Buffer) (*http.Response, error)
	Request(method, rurl, contentType string, buffer *bytes.Buffer) (*http.Response, error)
	PostForm(rurl string, formData url.Values) (*http.Response, error)
}

var _ HTTPClient = (*HttpRequest)(nil) // compile-time check
