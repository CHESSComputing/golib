#ifndef READER_H
#define READER_H

#include <stddef.h>

typedef struct {
    char **keys;
    char **values;
    size_t length;
} Metadata;

typedef struct {
    double *data;
    int *shape;
    int rank;
    size_t total_size;
    Metadata metadata;
    char *error;
} HDF5Result;

int read_hdf5(const char *filename, const char *dataset_path, HDF5Result *result);
void free_hdf5_result(HDF5Result *result);

#endif

