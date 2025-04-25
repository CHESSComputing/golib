#ifndef CONTENT_H
#define CONTENT_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
    char* name;
    char* dtype;
    int* shape;
    int ndim;
    long size;
} HDF5MetaData;

typedef struct {
    HDF5MetaData* datasets;
    int count;
} HDF5Content;

HDF5Content get_hdf5_content(const char* filename);
void free_hdf5_content(HDF5Content content);

#ifdef __cplusplus
}
#endif

#endif // CONTENT_H

