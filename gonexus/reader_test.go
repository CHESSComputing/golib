package gonexus

import (
	"testing"
)

func TestReadHDF5(t *testing.T) {
	meta, data, err := ReadHDF5("sample.h5", "mydata")
	if err != nil {
		t.Fatalf("Failed to read HDF5: %v", err)
	}

	if len(meta) == 0 {
		t.Errorf("Expected metadata, got none")
	}

	if len(data) != 100 {
		t.Errorf("Expected 100 data points, got %d", len(data))
	}
}
