all: build

gorelease:
	goreleaser release --snapshot --clean

build:
	./go_build.sh

test:
	touch $HOME/.foxden.yaml
	./go_test.sh
