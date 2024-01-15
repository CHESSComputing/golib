all: build

gorelease:
	goreleaser release --snapshot --clean

build:
	./go_build.sh

test:
	touch ~/.foxden.yaml
	./go_test.sh
