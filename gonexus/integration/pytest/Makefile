H5FILE=test.nxs

# Define variables for virtual environment
VENV_DIR = .venv

# Use the python executable from the venv once created
VENV_PYTHON = $(VENV_DIR)/bin/python
VENV_PIP = $(VENV_DIR)/bin/pip

# Define the python script name
PYTHON_SCRIPT = read_nexus.py

# Add .PHONY for targets that don't represent actual files or should always run
.PHONY: h5file clean_venv venv_setup

all: setup run clean

setup:
	ln -s ../$(H5FILE)

run: venv_setup
	@echo "-- run $(PYTHON_SCRIPT)"
	$(VENV_PYTHON) $(PYTHON_SCRIPT) $(H5FILE)

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
	rm $(H5FILE)
