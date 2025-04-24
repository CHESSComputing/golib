H5FILE=sample.h5

all: libreader.so build-go

libreader.so: reader.c reader.h
	gcc -c -fPIC -I/opt/local/include reader.c -o reader.o
	gcc -o libreader.so reader.o -shared -L/opt/local/lib -lhdf5

build-go:
	CGO_CFLAGS="-I/opt/local/include" CGO_LDFLAGS="-L/opt/local/lib" go build

test-go:
	CGO_CFLAGS="-I/opt/local/include" CGO_LDFLAGS="-L/opt/local/lib -L./" go test -v

main: ctest/main.c reader.h libreader.so
	gcc -I. -I/opt/local/include ctest/main.c -L. -lreader -L/opt/local/lib -lhdf5 -o main

run-main: main h5file
	LD_LIBRARY_PATH=. ./main

h5file:
	python3 generate_h5.py

test: all h5file test-go run-main

clean:
	rm -f *.o *.so main sample.h5
	go clean -testcache
