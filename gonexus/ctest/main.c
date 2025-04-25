#include <stdio.h>
#include <stdlib.h>
#include "reader.h"
#include "content.h"

int main() {
    const char *filename = "sample.h5";


    // Second test: content
    HDF5Content content = get_hdf5_content(filename);

    printf("Content of file:\n");
    for (int i = 0; i < content.count; ++i) {
        HDF5MetaData metadata = content.datasets[i];
        printf("  Dataset: %s\n", metadata.name);
        printf("    Type: %s\n", metadata.dtype);
        printf("    Shape: [");
        for (int j = 0; j < metadata.ndim; ++j) {
            printf("%lld", (long long)metadata.shape[j]);
            if (j < metadata.ndim - 1) printf(", ");
        }
        printf("]\n");
        printf("    Size: %zu\n", metadata.size);

        // now let's read data from the content dataset
        HDF5Result result;
        int status = read_hdf5(filename, metadata.name, &result);
        printf("status of read_hdf5 %d\n", status);

        printf("Dataset: %s\n", metadata.name);
        printf("Rank: %d\n", result.rank);
        printf("Shape: ");
        for (int i = 0; i < result.rank; i++) {
            printf("%d", result.shape[i]);
            if (i < result.rank - 1) printf(" x ");
        }
        printf("\n");

        printf("Metadata:\n");
        for (int i = 0; i < result.metadata.length; i++) {
            printf("  %s: %s\n", result.metadata.keys[i], result.metadata.values[i]);
        }

        printf("\nData:\n");
        int size = 1;
        for (int i = 0; i < result.rank; i++) {
            size *= result.shape[i];
        }

        printf("First 10 elements:\n");
        for (int i = 0; i < size && i < 10; i++) {
            printf("%.2f ", result.data[i]);
        }
        printf("\n");

        free(result.data);
        free(result.shape);
        free(result.error);

    }

    return 0;
}

