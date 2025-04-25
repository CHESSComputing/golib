#!/usr/bin/env bash

set -e

for d in $(go list ./... | grep -v vendor); do
    echo "Building $d"
    if [ "$d" == "github.com/CHESSComputing/golib/gonexus/integration/gotest" ]; then
        continue
    fi
    bdir=`echo $d | awk '{z=split($0,a,"/"); print a[z]}'`
    if [ "$bdir" == "badger" ]; then
        bdir="embed/badger"
    fi
    if [ "$bdir" == "clover" ]; then
        bdir="embed/clover"
    fi
    if [ "$bdir" == "tiedot" ]; then
        bdir="embed/tiedot"
    fi
    if [ "$bdir" == "gonexus" ]; then
        continue
    fi
    echo "cd $bdir"
    cd $bdir
    go build
    cd - 2>&1 1>& /dev/null
done
