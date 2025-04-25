package gonexus

/*
#cgo LDFLAGS: -L. -lnexus -lhdf5
#include <stdlib.h>
#include "content.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

// HDF5MetaData represents content of single dataset within HDF5 file
type HDF5MetaData struct {
	Name  string
	DType string
	Shape []int
	Size  int64
}

// Content provides content of the given HDF5 file and return list of HDF5MetaData structs
func Content(fileName string) ([]HDF5MetaData, error) {
	cFile := C.CString(fileName)
	defer C.free(unsafe.Pointer(cFile))

	result := C.get_hdf5_content(cFile)
	if result.count == 0 {
		return nil, errors.New("no datasets found")
	}

	var datasets []HDF5MetaData
	cDatasets := (*[1 << 30]C.HDF5MetaData)(unsafe.Pointer(result.datasets))[:result.count:result.count]
	for _, cdata := range cDatasets {
		shape := (*[1 << 10]C.int)(unsafe.Pointer(cdata.shape))[:cdata.ndim:cdata.ndim]
		dims := make([]int, cdata.ndim)
		for i := range dims {
			dims[i] = int(shape[i])
		}

		data := HDF5MetaData{
			Name:  C.GoString(cdata.name),
			DType: C.GoString(cdata.dtype),
			Shape: dims,
			Size:  int64(cdata.size),
		}
		datasets = append(datasets, data)
	}

	C.free_hdf5_content(result)
	return datasets, nil
}
