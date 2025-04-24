H5FILE=sample.h5

# Define variables for virtual environment
VENV_DIR = .venv

# Use the python executable from the venv once created
VENV_PYTHON = $(VENV_DIR)/bin/python
VENV_PIP = $(VENV_DIR)/bin/pip

# Define the python script name
PYTHON_SCRIPT = generate_h5.py

# Add .PHONY for targets that don't represent actual files or should always run
# Keep h5file phony if you want 'make h5file' to *always* run the script
# Or, if OUTPUT_H5_FILE is defined, make that the target instead (see commented example below)
.PHONY: h5file clean_venv venv_setup

all: libreader.so build-go

libreader.so: reader.c reader.h
	gcc -c -fPIC -I/opt/local/include reader.c -o reader.o
	gcc -o libreader.so reader.o -shared -L/opt/local/lib -lhdf5

build-go:
	CGO_CFLAGS="-I/opt/local/include" CGO_LDFLAGS="-L/opt/local/lib" go build

test-go:
	CGO_CFLAGS="-I/opt/local/include" CGO_LDFLAGS="-L/opt/local/lib -L./" go test -v

main: ctest/main.c reader.h libreader.so
	gcc -I. -I/opt/local/include ctest/main.c -L. -lreader -L/opt/local/lib -lhdf5 -o main

run-main: main h5file
	LD_LIBRARY_PATH=. ./main

# h5file:
#     python3 generate_h5.py

test: all h5file test-go run-main

# target to generate the HDF5 file
h5file: venv_setup $(PYTHON_SCRIPT)
	@echo "--- Generating HDF5 file using $(PYTHON_SCRIPT) ---"
	# Execute the python script using the virtual environment's python
	$(VENV_PYTHON) $(PYTHON_SCRIPT)
	@echo "--- HDF5 file generation complete ---"

# Target to set up the virtual environment and install dependencies
venv_setup:
	@echo "--- Checking/Setting up Python virtual environment in $(VENV_DIR) ---"
	# Create venv directory if it doesn't exist using python3's venv module
	test -d $(VENV_DIR) || python3 -m venv $(VENV_DIR)
	# Install required packages using pip from the virtual environment
	# Pip is generally idempotent, so running this is usually fine even if already installed
	$(VENV_PIP) install --upgrade pip # Good practice to upgrade pip first
	$(VENV_PIP) install h5py numpy
	@echo "--- Virtual environment ready. ---"

# Optional: Target to clean up the virtual environment
clean_venv:
	@echo "--- Removing Python virtual environment $(VENV_DIR) ---"
	rm -rf $(VENV_DIR)

# --- Alternative if you want the target to be the actual output file ---
# --- This avoids re-running if the script hasn't changed AND the file exists ---
# .PHONY: clean_venv venv_setup # output file is NOT phony
# $(OUTPUT_H5_FILE): venv_setup $(PYTHON_SCRIPT)
#	@echo "--- Generating $(OUTPUT_H5_FILE) using $(PYTHON_SCRIPT) ---"
#	$(VENV_PYTHON) $(PYTHON_SCRIPT) # Script must create $(OUTPUT_H5_FILE)
#	@echo "--- $(OUTPUT_H5_FILE) generation complete ---"
# clean: # Example clean target including the generated file
#	rm -f $(OUTPUT_H5_FILE)
# clean_venv:
#	rm -rf $(VENV_DIR)

clean:
	rm -f *.o *.so main sample.h5
	go clean -testcache
