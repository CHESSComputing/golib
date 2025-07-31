package streamer

import (
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/tiff"
)

// ImageChunk defines structure for image chunks
type ImageChunk struct {
	Data     []byte
	MIMEType string
	Name     string
	Width    int
	Height   int
	Format   string
}

// ImageReader defines image reader
type ImageReader struct {
	files []string
	index int
}

// NewImageReader creates new image reader
func NewImageReader(dir string) (*ImageReader, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".tif" || ext == ".tiff" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ImageReader{files: files}, nil
}

// ReadChunk provides read chunk capabilities of image reader
func (r *ImageReader) ReadChunk(idx int) (*Chunk, error) {
	if idx >= len(r.files) {
		return nil, io.EOF
	}

	path := r.files[idx]
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	r.index = idx + 1

	// Determine MIME type
	ext := strings.ToLower(filepath.Ext(path))
	mime := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".tif":  "image/tiff",
		".tiff": "image/tiff",
	}[ext]

	/*
		// Extract metadata
		img, format, err := image.DecodeConfig(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}

		meta := map[string]any{
			"name":   filepath.Base(path),
			"width":  img.Width,
			"height": img.Height,
			"format": format,
		}

		// Embed metadata in header
		buf := bytes.NewBuffer(nil)
		buf.WriteString("-----BEGIN-METADATA-----\n")
		for k, v := range meta {
			strVal := fmt.Sprintf("%v", v)
			strVal = strings.ReplaceAll(strVal, "\n", "")
			strVal = strings.ReplaceAll(strVal, "\r", "")
			strVal = strings.TrimSpace(strVal)
			buf.WriteString(k + ": " + strVal + "\n")

		}
		buf.WriteString("-----END-METADATA-----\n")
		buf.Write(data)
	*/

	return &Chunk{
		ContentType: mime,
		Data:        data,
	}, nil
}

// Reset provides reset of image reader
func (r *ImageReader) Reset() error {
	r.index = 0
	return nil
}
