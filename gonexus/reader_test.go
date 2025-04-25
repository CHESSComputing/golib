package gonexus

import (
	"testing"
)

func TestReadHDF5(t *testing.T) {
	filename := "sample.h5"
	dataset := "mydata"
	h5data, err := ReadHDF5(filename, dataset)
	if err != nil {
		t.Fatalf("Failed to read HDF5: %v", err)
	}
	if h5data == nil {
		t.Fatalf("unable to get data from %s", filename)
	}

	if len(h5data.Data) != 100 {
		t.Errorf("Expected 100 data points, got %d", len(h5data.Data))
	}
}
