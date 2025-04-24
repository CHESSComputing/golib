#ifndef READER_H
#define READER_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
    char **keys;
    char **values;
    int count;
} MetaData;

typedef struct {
    double *data;
    int length;
} DataArray;

typedef struct {
    MetaData metadata;
    DataArray dataset;
    char *error;
} HDF5Result;

HDF5Result read_hdf5(const char *filename, const char *dataset_name);
void free_result(HDF5Result result);

#ifdef __cplusplus
}
#endif

#endif // READER_H

