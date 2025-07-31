package streamer

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

// helper function to create NPY file
func createTestNPYFile(t *testing.T, dir string, name string, data []float64) string {
	t.Helper()

	// Construct a minimal NPY header for a 1D float64 array
	header := []byte("\x93NUMPY")
	header = append(header, 0x01, 0x00) // Version 1.0
	headerDict := "{'descr': '<f8', 'fortran_order': False, 'shape': (4,), }"
	headerPadding := 16 - ((len(headerDict) + 10) % 16)
	headerLen := uint16(len(headerDict) + headerPadding)

	buf := new(bytes.Buffer)
	buf.Write(header)
	binary.Write(buf, binary.LittleEndian, headerLen)
	buf.WriteString(headerDict)
	buf.Write(bytes.Repeat([]byte(" "), headerPadding))

	// Write data
	for _, val := range data {
		binary.Write(buf, binary.LittleEndian, val)
	}

	outPath := filepath.Join(dir, name)
	err := os.WriteFile(outPath, buf.Bytes(), 0644)
	if err != nil {
		t.Fatalf("failed to write npy file: %v", err)
	}
	return outPath
}

// TestNPYReader_ReadChunk provides unit test for NPYReader
func TestNPYReader_ReadChunk(t *testing.T) {
	dir := t.TempDir()
	createTestNPYFile(t, dir, "test.npy", []float64{1.1, 2.2, 3.3, 4.4})

	reader, err := NewNPYReader(dir)
	if err != nil {
		t.Fatalf("failed to create NPYReader: %v", err)
	}

	chunk, err := reader.ReadChunk(0)
	if err != nil {
		t.Fatalf("ReadChunk failed: %v", err)
	}

	if len(chunk.Data) == 0 {
		t.Error("expected non-empty data from NPYReader")
	}
	if chunk.ContentType != "application/octet-stream" {
		t.Errorf("expected content-type 'application/octet-stream', got %s", chunk.ContentType)
	}
}
