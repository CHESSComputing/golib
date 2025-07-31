package streamer

import (
	"io"
	"os"
	"path/filepath"
)

// NPYReader provides structure for numpy files
type NPYReader struct {
	files []string
	index int
}

// NewNPYReader provides numpy reader
func NewNPYReader(dir string) (*NPYReader, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".npy" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &NPYReader{files: files, index: 0}, nil
}

// ReadChunk provides ability to read specific chunk of numpy data file
func (r *NPYReader) ReadChunk(_ int) (*Chunk, error) {
	if r.index >= len(r.files) {
		return nil, io.EOF
	}
	data, err := os.ReadFile(r.files[r.index])
	if err != nil {
		return nil, err
	}
	r.index++
	return &Chunk{ContentType: "application/octet-stream", Data: data}, nil
}

// Reset resets numpy reader index
func (r *NPYReader) Reset() error {
	r.index = 0
	return nil
}
