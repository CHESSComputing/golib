#include "reader.h"
#include "hdf5.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

HDF5Result read_hdf5(const char *filename, const char *dataset_name) {
    HDF5Result result = {0};

    hid_t file = H5Fopen(filename, H5F_ACC_RDONLY, H5P_DEFAULT);
    if (file < 0) {
        result.error = strdup("Failed to open file");
        return result;
    }

    hid_t dset = H5Dopen2(file, dataset_name, H5P_DEFAULT);
    if (dset < 0) {
        H5Fclose(file);
        result.error = strdup("Failed to open dataset");
        return result;
    }

    hid_t space = H5Dget_space(dset);
    hsize_t dims[1];
    H5Sget_simple_extent_dims(space, dims, NULL);

    result.dataset.length = dims[0];
    result.dataset.data = (double *)malloc(dims[0] * sizeof(double));
    H5Dread(dset, H5T_NATIVE_DOUBLE, H5S_ALL, H5S_ALL, H5P_DEFAULT, result.dataset.data);

    // Metadata: simple implementation, returns fixed dummy keys/values
    result.metadata.count = 1;
    result.metadata.keys = malloc(sizeof(char*));
    result.metadata.values = malloc(sizeof(char*));
    result.metadata.keys[0] = strdup("dummy_key");
    result.metadata.values[0] = strdup("dummy_value");

    H5Sclose(space);
    H5Dclose(dset);
    H5Fclose(file);

    return result;
}

void free_result(HDF5Result result) {
    for (int i = 0; i < result.metadata.count; i++) {
        free(result.metadata.keys[i]);
        free(result.metadata.values[i]);
    }
    free(result.metadata.keys);
    free(result.metadata.values);
    free(result.dataset.data);
    free(result.error);
}

