.PHONY: all build test tools docs gorelease clean

all: build tools

gorelease:
	goreleaser release --snapshot --clean

build:
	./go_build.sh

test:
	touch ~/.foxden.yaml
	./go_test.sh

tools:
	@echo "build tools/gendocs ..."
	mkdir -p bin
	go build -o bin/gendocs ./tools/gendocs/main.go

docs: tools
	mkdir -p docs
	./bin/gendocs -root . -out docs/ENDPOINTS.md

clean:
	rm -rf bin
