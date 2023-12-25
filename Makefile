all: build

gorelease:
	goreleaser release --snapshot --clean

build:
	./go_build.sh

test:
	./go_test.sh
