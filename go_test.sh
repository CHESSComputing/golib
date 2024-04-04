#!/usr/bin/env bash

set -e
echo "mode: atomic" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    echo "Testing $d"
    go test -v $d
    if [ $d == "github.com/CHESSComputing/golib/zenodo" ]; then
        continue
    fi
    if [ $d == "github.com/CHESSComputing/golib/ldap" ]; then
        continue
    fi
    echo "Coverage $d"
    go test -race -coverprofile=profile.out -covermode=atomic "$d"
    if [ -f profile.out ]; then
        cat profile.out | grep -v "mode: atomic" >> coverage.txt
        rm profile.out
    fi
done
echo "Run the following command to see coverage:"
echo "go tool cover -html=coverage.txt"
