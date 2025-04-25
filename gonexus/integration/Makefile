# Define the URL of the test file
TEST_FILE_URL = https://github.com/yourusername/yourrepo/releases/download/v0.7.4/test.nxs
TEST_FILE = test.nxs

.PHONY: all fetch-data pytest ctest gotest

all: fetch-data pytest ctest gotest

fetch-data:
	@echo "--- Fetching test.nxs ---"
	@if [ ! -f "$(TEST_FILE)" ]; then \
		echo "Downloading test.nxs..."; \
		curl -L -o "$(TEST_FILE)" "$(TEST_FILE_URL)"; \
	else \
		echo "test.nxs already exists."; \
	fi

pytest:
	@echo "--- Perform Python test ---"
	$(MAKE) -C pytest

ctest:
	@echo "--- Perform C test ---"
	$(MAKE) -C ctest

gotest:
	@echo "--- Perform Go test ---"
	$(MAKE) -C gotest
