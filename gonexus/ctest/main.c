#include <stdio.h>
#include <stdlib.h>
#include "reader.h"
#include "content.h"

int main() {
/*
 * int main(int argc, char **argv) {
    if (argc < 2) {
        printf("Usage: %s <file.h5>\n", argv[0]);
        return 1;
    }

    const char *filename = argv[1];
*/
    const char *filename = "sample.h5";
/*    const char *dataset = "mydata";*/


    // Second test: content
    HDF5Content content = get_hdf5_content(filename);

    printf("Content of file:\n");
    for (int i = 0; i < content.count; ++i) {
        HDF5Data data = content.datasets[i];
        printf("  Dataset: %s\n", data.name);
        printf("    Type: %s\n", data.dtype);
        printf("    Shape: [");
        for (int j = 0; j < data.ndim; ++j) {
            printf("%lld", (long long)data.shape[j]);
            if (j < data.ndim - 1) printf(", ");
        }
        printf("]\n");
        printf("    Size: %zu\n", data.size);

        // now let's read data from the content dataset
        HDF5Result result = read_hdf5(filename, data.name);

        printf("Metadata:\n");
        for (int i = 0; i < result.metadata.count; i++) {
            printf("  %s: %s\n", result.metadata.keys[i], result.metadata.values[i]);
        }

        printf("\nData:\n");
        for (int i = 0; i < result.dataset.length; i++) {
            printf("  %f\n", result.dataset.data[i]);
        }

        free_result(result);

    }

    return 0;
}

