package main

/*
#cgo CFLAGS: -I/opt/local/include -I/usr/local/include
#cgo LDFLAGS: -L/opt/local/lib -L/usr/local/lib -lhdf5
#include "hdf5.h"

// Disable automatic HDF5 error printing
void suppress_hdf5_errors() {
    H5Eset_auto(H5E_DEFAULT, NULL, NULL);
}
*/
import "C"
import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/CHESSComputing/golib/gonexus"
)

func main() {
	var nexusfile string
	flag.StringVar(&nexusfile, "nexusfile", "", "nexus file name")
	var dataset string
	flag.StringVar(&dataset, "dataset", "", "dataset name")
	var nElements int
	flag.IntVar(&nElements, "nElements", 0, "number of elements to show")
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	C.suppress_hdf5_errors()

	// print content of nexus file
	content, err := gonexus.Content(nexusfile)
	if err == nil {
		fmt.Println("METADATA:")
		for _, c := range content {
			fmt.Println("metadata: %+v", c)
		}
	}
	// read actual data
	hdf5data, err := gonexus.ReadHDF5(nexusfile, dataset)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	if nElements > 0 && nElements < hdf5data.Size {
		fmt.Printf("DATA: %v\n", hdf5data.Data[:nElements])
	} else {
		fmt.Printf("DATA: %v\n", hdf5data.Data)
	}
}
