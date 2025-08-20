#!/usr/bin/env bash

set -e
echo "mode: atomic" > coverage.txt

skipServices="zenodo schema datacite doi s3 gonexus gonexus/integration/gotest"

for d in $(go list ./... | grep -v vendor); do
    echo "Testing $d"
    skip="false"
    for s in $skipServices; do
        if [ $d == "github.com/CHESSComputing/golib/$s" ]; then
            skip="true"
        fi
    done
    if [ "$skip" == "true" ]; then
        echo "Skipping $d, not test files required..."
        continue
    fi
    go test -v $d
    echo "Coverage $d"
    if [ "`hostname`" == "vkarm" ]; then
      go test -race -coverprofile=profile.out -covermode=atomic "$d"
      if [ -f profile.out ]; then
          cat profile.out | grep -v "mode: atomic" >> coverage.txt
          rm profile.out
      fi
    fi
done
if [ "`hostname`" == "vkarm" ]; then
  echo "Run the following command to see coverage:"
  echo "go tool cover -html=coverage.txt"
fi
