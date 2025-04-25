IDIR=../..
H5FILE=test.nxs
EXE=main

all: setup build run clean

setup:
	ln -s ../$(H5FILE)
	ln -s $(IDIR)/libnexus.so

build:
	gcc -I$(IDIR) -I/opt/local/include -I/usr/local/include ./main.c -L$(IDIR) -lnexus -L/opt/local/lib -L/usr/local/lib -lhdf5 -o $(EXE)

run:
	@echo
	@echo "--- show 10 elements of saxs_cake/data/chi_deg dataset"
	./$(EXE) $(H5FILE) "saxs_cake/data/chi_deg"
	@echo
	@echo "--- show 30 elements of saxs_cake/data/I dataset"
	./$(EXE) $(H5FILE) "saxs_cake/data/I"

# target to clean up our area
clean:
	rm -f *.o *.so $(EXE) $(H5FILE)
