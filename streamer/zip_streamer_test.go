package streamer

import (
	"archive/zip"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestStreamAsZip provides unit test for zip streamer
func TestStreamAsZip(t *testing.T) {
	dir := t.TempDir()
	setupTestImages(t, dir)

	reader, err := NewImageReader(dir)
	if err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	router.GET("/zip", func(c *gin.Context) {
		StreamAsZip(c, reader, "bundle.zip")
	})

	req := httptest.NewRequest("GET", "/zip", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Check zip content
	r := bytes.NewReader(w.Body.Bytes())
	zr, err := zip.NewReader(r, int64(len(w.Body.Bytes())))
	if err != nil {
		t.Fatalf("invalid zip: %v", err)
	}
	if len(zr.File) == 0 {
		t.Error("zip file is empty")
	}
}
