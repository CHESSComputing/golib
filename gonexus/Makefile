H5FILE=sample.h5

# Define variables for virtual environment
VENV_DIR = .venv

# Use the python executable from the venv once created
VENV_PYTHON = $(VENV_DIR)/bin/python
VENV_PIP = $(VENV_DIR)/bin/pip

# Define the python script name
PYTHON_SCRIPT = generate_h5.py

# Add .PHONY for targets that don't represent actual files or should always run
.PHONY: h5file clean_venv venv_setup

# all targets to execute under `make`
all: libreader.so build-go

# target to build reader C library
libreader.so: reader.c reader.h
	gcc -c -fPIC -I/opt/local/include reader.c -o reader.o
	gcc -o libreader.so reader.o -shared -L/opt/local/lib -lhdf5

# target to build go code with C library
build-go:
	CGO_CFLAGS="-I/opt/local/include" CGO_LDFLAGS="-L/opt/local/lib" go build

# target to test go code
test-go:
	CGO_CFLAGS="-I/opt/local/include" CGO_LDFLAGS="-L/opt/local/lib -L./" go test -v

# target to build C main executable (located at ctest)
main: ctest/main.c reader.h libreader.so
	gcc -I. -I/opt/local/include ctest/main.c -L. -lreader -L/opt/local/lib -lhdf5 -o main

# target to run main C executable with our h5 file
run-main: main h5file
	LD_LIBRARY_PATH=. ./main

# target to run Go tests
test: all h5file test-go run-main

# target to generate the HDF5 file
h5file: venv_setup $(PYTHON_SCRIPT)
	@echo "--- Generating HDF5 file using $(PYTHON_SCRIPT) ---"
	$(VENV_PYTHON) $(PYTHON_SCRIPT)
	@echo "--- HDF5 file generation complete ---"

# target to set up the virtual environment and install dependencies
venv_setup:
	@echo "--- Checking/Setting up Python virtual environment in $(VENV_DIR) ---"
	test -d $(VENV_DIR) || python3 -m venv $(VENV_DIR)
	$(VENV_PIP) install --upgrade pip # Good practice to upgrade pip first
	$(VENV_PIP) install h5py numpy
	@echo "--- Virtual environment ready. ---"

# target to clean up the virtual environment
clean_venv:
	@echo "--- Removing Python virtual environment $(VENV_DIR) ---"
	rm -rf $(VENV_DIR)

# target to clean up our area
clean: clean_venv
	rm -f *.o *.so main sample.h5
	go clean -testcache
