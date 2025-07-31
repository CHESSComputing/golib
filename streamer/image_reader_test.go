package streamer

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

// helper function to generate test jpeg images
func setupTestImages(t *testing.T, dir string) []string {
	t.Helper()
	files := []string{
		filepath.Join(dir, "test1.jpg"),
		filepath.Join(dir, "test2.jpg"),
	}
	for _, f := range files {
		buf := new(bytes.Buffer)
		img := image.NewRGBA(image.Rect(0, 0, 10, 10)) // 10x10 black image
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				img.Set(x, y, color.RGBA{255, 0, 0, 255}) // red pixels
			}
		}
		if err := jpeg.Encode(buf, img, nil); err != nil {
			t.Fatalf("failed to encode JPEG: %v", err)
		}
		if err := os.WriteFile(f, buf.Bytes(), 0644); err != nil {
			t.Fatalf("failed to write test image: %v", err)
		}
	}
	return files
}

// TestImageReader_ReadChunk unit test for image reader ReadChunk functionality
func TestImageReader_ReadChunk(t *testing.T) {
	dir := t.TempDir()
	setupTestImages(t, dir)

	reader, err := NewImageReader(dir)
	if err != nil {
		t.Fatalf("failed to create ImageReader: %v", err)
	}

	chunk, err := reader.ReadChunk(0)
	if err != nil {
		t.Fatalf("failed to read chunk: %v", err)
	}
	if len(chunk.Data) == 0 {
		t.Error("chunk data is empty")
	}
	if chunk.ContentType != "image/jpeg" {
		t.Errorf("expected image/jpeg, got %s", chunk.ContentType)
	}
}
