#include "reader.h"
#include "hdf5.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

int read_hdf5(const char *filename, const char *dataset_path, HDF5Result *result) {
    hid_t file = -1, dataset = -1, dataspace = -1;
    herr_t status;

    file = H5Fopen(filename, H5F_ACC_RDONLY, H5P_DEFAULT);
    if (file < 0) {
        result->error = strdup("Failed to open file");
        return -1;
    }

    dataset = H5Dopen(file, dataset_path, H5P_DEFAULT);
    if (dataset < 0) {
        result->error = strdup("Failed to open dataset");
        goto cleanup;
    }

    dataspace = H5Dget_space(dataset);
    result->rank = H5Sget_simple_extent_ndims(dataspace);

    hsize_t dims[16];  // arbitrary max rank
    H5Sget_simple_extent_dims(dataspace, dims, NULL);

    result->shape = malloc(result->rank * sizeof(int));
    result->total_size = 1;
    for (int i = 0; i < result->rank; i++) {
        result->shape[i] = dims[i];
        result->total_size *= dims[i];
    }

    result->data = malloc(result->total_size * sizeof(double));
    status = H5Dread(dataset, H5T_NATIVE_DOUBLE, H5S_ALL, H5S_ALL, H5P_DEFAULT, result->data);
    if (status < 0) {
        result->error = strdup("Failed to read dataset");
        goto cleanup;
    }

    result->error = NULL;

cleanup:
    if (dataspace >= 0) H5Sclose(dataspace);
    if (dataset >= 0) H5Dclose(dataset);
    if (file >= 0) H5Fclose(file);
    return result->error ? -1 : 0;
}

void free_hdf5_result(HDF5Result *result) {
    if (result->data) free(result->data);
    if (result->shape) free(result->shape);
    if (result->error) free(result->error);
}

