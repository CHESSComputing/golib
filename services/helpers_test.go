package services

import (
	"testing"
)

// TestHTTPResponse
func TestHTTPResponse(t *testing.T) {
	httpReq := NewHttpRequest("scope", 0)
	resp, err := httpReq.Get("http://localhost:12345")
	if err == nil {
		t.Errorf("Fail to process http request, err=%v\n", err)
	}
	if resp != nil {
		t.Errorf("HTTP response is not nil, resp=%v\n", resp)
	}
}
