package streamer

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockReader struct {
	called bool
}

func (m *mockReader) ReadChunk(_ int) (*Chunk, error) {
	if m.called {
		return nil, io.EOF
	}
	m.called = true
	return &Chunk{
		ContentType: "image/jpeg",
		Data:        []byte{0xFF, 0xD8, 0xFF, 0xD9},
	}, nil
}

func (m *mockReader) Reset() error {
	m.called = false
	return nil
}

// TestGinBinaryStreamHandler provides unit test for GinBinaryStreamHandler
func TestGinBinaryStreamHandler(t *testing.T) {
	handler := GinBinaryStreamHandler(&mockReader{})

	router := setupTestRouter("/test", handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "image/jpeg" {
		t.Errorf("expected image/jpeg content-type, got %s", ct)
	}
}

func setupTestRouter(path string, handler gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.GET(path, handler)
	return r
}
