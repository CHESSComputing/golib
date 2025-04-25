package gonexus

/*
#cgo CFLAGS: -I/opt/local/include
#cgo LDFLAGS: -L/opt/local/lib -lnexus -lhdf5
#include <stdlib.h>
#include "reader.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

// HDF5Data represents HDF5 data
type HDF5Data struct {
	Data  []float64
	Shape []int
	Rank  int
	Size  int
}

// ReadHDF5 function reads a given dataset from a file path and returns
// medata (as a map), data array (as floats64) and error
func ReadHDF5(filename, dataset string) (*HDF5Data, error) {
	cfile := C.CString(filename)
	cdataset := C.CString(dataset)
	defer C.free(unsafe.Pointer(cfile))
	defer C.free(unsafe.Pointer(cdataset))

	var result C.HDF5Result
	if C.read_hdf5(cfile, cdataset, &result) != 0 {
		defer C.free_hdf5_result(&result)
		return nil, errors.New(C.GoString(result.error))
	}

	defer C.free_hdf5_result(&result)

	// Convert C.int* to Go []int
	shapeSlice := unsafe.Slice(result.shape, result.rank)
	shape := make([]int, result.rank)
	size := 1
	for i := 0; i < int(result.rank); i++ {
		dim := int(shapeSlice[i])
		shape[i] = dim
		size *= dim
	}

	data := make([]float64, size)
	cData := (*[1 << 30]C.double)(unsafe.Pointer(result.data))[:size:size]
	for i := range data {
		data[i] = float64(cData[i])
	}

	/*
		// Convert data buffer to Go []float64
		dataSlice := unsafe.Slice(result.data, size)
		data := make([]float64, size)
		copy(data, dataSlice)
	*/

	return &HDF5Data{
		Data:  data,
		Shape: shape,
		Rank:  int(result.rank),
		Size:  size,
	}, nil
}
