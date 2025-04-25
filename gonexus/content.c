#include "content.h"
#include "hdf5.h"
#include <stdlib.h>
#include <string.h>

#define MAX_DATASETS 128
#define MAX_NAME_LEN 256

static char* copy_string(const char* src) {
    char* dst = malloc(strlen(src) + 1);
    strcpy(dst, src);
    return dst;
}

static char* dtype_to_str(H5T_class_t class_id) {
    switch (class_id) {
        case H5T_INTEGER: return copy_string("int");
        case H5T_FLOAT:   return copy_string("float");
        case H5T_STRING:  return copy_string("string");
        default:          return copy_string("unknown");
    }
}

static herr_t dataset_info_cb(hid_t loc_id, const char *name, const H5L_info_t *info, void *op_data) {
    HDF5Content* content = (HDF5Content*)op_data;
    if (content->count >= MAX_DATASETS) return 1;

    hid_t dset_id = H5Dopen2(loc_id, name, H5P_DEFAULT);
    if (dset_id < 0) return 0;

    HDF5MetaData* item = &content->datasets[content->count];
    item->name = copy_string(name);

    hid_t space = H5Dget_space(dset_id);
    item->ndim = H5Sget_simple_extent_ndims(space);
    item->shape = malloc(item->ndim * sizeof(int));
    hsize_t dims[item->ndim];
    H5Sget_simple_extent_dims(space, dims, NULL);
    item->size = 1;
    for (int i = 0; i < item->ndim; i++) {
        item->shape[i] = (int)dims[i];
        item->size *= dims[i];
    }
    H5Sclose(space);

    hid_t dtype = H5Dget_type(dset_id);
    item->dtype = dtype_to_str(H5Tget_class(dtype));
    H5Tclose(dtype);

    H5Dclose(dset_id);
    content->count++;
    return 0;
}

HDF5Content get_hdf5_content(const char* filename) {
    HDF5Content content;
    content.datasets = malloc(MAX_DATASETS * sizeof(HDF5MetaData));
    content.count = 0;

    hid_t file_id = H5Fopen(filename, H5F_ACC_RDONLY, H5P_DEFAULT);
    if (file_id < 0) return content;

    H5Literate(file_id, H5_INDEX_NAME, H5_ITER_NATIVE, NULL, dataset_info_cb, &content);
    H5Fclose(file_id);

    return content;
}

void free_hdf5_content(HDF5Content content) {
    for (int i = 0; i < content.count; i++) {
        free(content.datasets[i].name);
        free(content.datasets[i].dtype);
        free(content.datasets[i].shape);
    }
    free(content.datasets);
}

