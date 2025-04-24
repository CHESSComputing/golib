package gonexus

/*
#cgo CFLAGS: -I/opt/local/include
#cgo LDFLAGS: -L/opt/local/lib -lreader -lhdf5
#include <stdlib.h>
#include "reader.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

// ReadHDF5 function reads a given dataset from a file path and returns
// medata (as a map), data array (as floats64) and error
func ReadHDF5(filePath, dataset string) (map[string]string, []float64, error) {
	cfile := C.CString(filePath)
	cdset := C.CString(dataset)
	defer C.free(unsafe.Pointer(cfile))
	defer C.free(unsafe.Pointer(cdset))

	result := C.read_hdf5(cfile, cdset)
	defer C.free_result(result)

	if result.error != nil {
		return nil, nil, errors.New(C.GoString(result.error))
	}

	count := int(result.metadata.count)

	// Convert C **char to Go []string
	keysPtr := unsafe.Pointer(result.metadata.keys)
	keys := (*[1 << 30]*C.char)(keysPtr)[:count:count]

	valsPtr := unsafe.Pointer(result.metadata.values)
	vals := (*[1 << 30]*C.char)(valsPtr)[:count:count]

	meta := make(map[string]string)
	for i := 0; i < count; i++ {
		key := C.GoString(keys[i])
		val := C.GoString(vals[i])
		meta[key] = val
	}

	// Convert C array to Go slice
	dataLen := int(result.dataset.length)
	data := (*[1 << 30]float64)(unsafe.Pointer(result.dataset.data))[:dataLen:dataLen]

	return meta, data, nil
}
