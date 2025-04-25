package gonexus

import (
	"fmt"
	"testing"
)

func TestContent(t *testing.T) {
	content, err := Content("sample.h5")
	if err != nil {
		t.Fatalf("Failed to read HDF5: %v", err)
	}
	// [{Name:mydata DType:float Shape:[100] Size:100}]
	if len(content) != 1 {
		t.Fatal("wrong size of fetched content")
	}
	for _, c := range content {
		if c.Name != "mydata" {
			t.Fatal("wrong name of dataset")
		}
		if c.Size != 100 {
			t.Fatal("wrong size of dataset")
		}
	}
	fmt.Printf("content: %+v\n", content)
}
