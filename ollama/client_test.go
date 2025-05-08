package ollama

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendPrompt(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"response":"hello world"}`))
	}))
	defer server.Close()

	client := NewClient(Config{
		Host:   "localhost",
		Port:   80,
		Model:  "llama2",
		Client: server.Client(),
	})
	client.url = server.URL

	ctx := context.Background()
	resp, err := client.SendPrompt(ctx, "hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", resp)
	}
}
