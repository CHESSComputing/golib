#!/usr/bin/env bash

set -e

for d in $(go list ./... | grep -v vendor); do
    echo "Building $d"
    bdir=`echo $d | awk '{z=split($0,a,"/"); print a[z]}'`
    echo "cd $bdir"
    cd $bdir
    go build
    cd -
done
