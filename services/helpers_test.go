package services

import (
	"testing"
)

// TestHTTPResponse
func TestHTTPResponse(t *testing.T) {
	httpReq := NewHttpRequest("scope", 0)
	resp, err := httpReq.Get("http://localhost:1234")
	if err == nil {
		t.Error("Fail to process http request")
	}
	if resp != nil {
		t.Error("HTTP response is not nil")
	}
}
