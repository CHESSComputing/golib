### gonexus
gonexus is a Go library for reading NeXus (HDF5) scientific data files
using CGO bindings to HDF5 C libraries.
It provides Go-friendly access to metadata and dataset contents for
scientific computing, machine learning, or data analysis.

Features
- Read NeXus (HDF5) files in Go
- Extract metadata as Go maps
- Read datasets (1D, 2D, 3D) into Go slices
- Fast C-backed access using CGO

Requirements: Go 1.18+, HDF5 development libraries (libhdf5 + headers), C compiler (GCC/Clang)

---

### Installation
```
# Clone the repository:
git clone https://github.com/CHESSComputing/golib.git
cd golib/gonexus

# Build the C helper library:
make
```
This creates libnexus.so library

---

### Environment Setup
Before running Go code, point your system to the shared library:
```
# on Linux
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/path/to/gonexus

# On macOS:
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:/path/to/gonexus
```

---

### Example Program
For concrete examples of how to use this library please refer to
[integration](integration) area which contains examples written
in C, Go, and Python.
