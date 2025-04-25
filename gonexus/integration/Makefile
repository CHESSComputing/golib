# Define the URL of the test file
TEST_FILE_URL = https://github.com/CHESSComputing/golib/releases/download/v0.7.4/test.nxs
TEST_FILE = test.nxs

.PHONY: all fetch-data pytest ctest gotest

all: fetch-data pytest ctest gotest clean

fetch-data:
	@echo "--- Fetching test.nxs ---"
	@if [ ! -f "$(TEST_FILE)" ]; then \
		echo "Downloading test.nxs..."; \
		curl -s -L -o "$(TEST_FILE)" "$(TEST_FILE_URL)"; \
	else \
		echo "test.nxs already exists."; \
	fi

pytest:
	@echo "--- Perform Python test ---"
	cd pytest; $(MAKE) clean; $(MAKE); cd -

ctest:
	@echo "--- Perform C test ---"
	cd ctest; $(MAKE) clean; $(MAKE); cd -

gotest:
	@echo "--- Perform Go test ---"
	cd gotest; $(MAKE) clean; $(MAKE); cd -

clean:
	rm $(TEST_FILE)
	@if [ -f "pytest/$(H5FILE)" ]; then \
		echo "Removing pytest/$(H5FILE)"; \
		rm "$(H5FILE)"; \
	fi
