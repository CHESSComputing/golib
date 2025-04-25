IDIR=../..
H5FILE=test.nxs
EXE=read_nexus

all: setup build run clean

setup:
	ln -s ../$(H5FILE)
	ln -s $(IDIR)/libnexus.so

build:
	CGO_LDFLAGS="-L../../" go build -o $(EXE)

run:
	@echo
	@echo "--- show 10 elements of saxs_cake/data/chi_deg dataset"
	./$(EXE) -nexusfile=$(H5FILE) -dataset="saxs_cake/data/chi_deg" -nElements=10
	@echo
	@echo "--- show 30 elements of saxs_integrated/data/I dataset"
	./$(EXE) -nexusfile=$(H5FILE) -dataset="saxs_integrated/data/I" -nElements=30

# target to clean up our area
clean:
	rm -f *.o *.so $(EXE) $(H5FILE)
