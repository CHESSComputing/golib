#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "reader.h"
#include "content.h"

int main(int argc, char **argv) {
    if (argc < 3) {
        printf("Usage: %s <fileName> <datasetName>\n", argv[0]);
        return 1;
    }

    const char *filename = argv[1];
    const char *dataset = argv[2];

    HDF5Result result;
    memset(&result, 0, sizeof(HDF5Result));  // set all fields to 0 / NULL

    /*
    printf("[DEBUG] result.rank = %d\n", result.rank);
    printf("[DEBUG] result.shape = %p\n", result.shape);
    printf("[DEBUG] result.data = %p\n", result.data);
    printf("[DEBUG] result.metadata.length = %zu\n", result.metadata.length);
    printf("[DEBUG] result.metadata.keys = %p\n", result.metadata.keys);
    printf("[DEBUG] result.metadata.values = %p\n", result.metadata.values);
    */
                                             //
    int status = read_hdf5(filename, dataset, &result);
    if (status != 0) {
        printf("Fail to read %s\n", filename);
        return 1;
    }

    printf("Dataset: %s\n", dataset);
    printf("Rank: %d\n", result.rank);
    printf("Shape: ");
    for (int i = 0; i < result.rank; i++) {
        printf("%d", result.shape[i]);
        if (i < result.rank - 1) printf(" x ");
    }
    printf("\n");

    /*
    printf("Metadata:\n");
    for (int i = 0; i < result.metadata.length; i++) {
        printf("  %s: %s\n", result.metadata.keys[i], result.metadata.values[i]);
    }
    */

    printf("Metadata:\n");
    if (result.metadata.length > 0 && result.metadata.keys && result.metadata.values) {
        for (int i = 0; i < result.metadata.length; i++) {
            if (result.metadata.keys[i] && result.metadata.values[i]) {
                printf("  %s: %s\n", result.metadata.keys[i], result.metadata.values[i]);
            }
        }

        // when we done with metadata we free its resources
        for (int i = 0; i < result.metadata.length; i++) {
            free(result.metadata.keys[i]);
            free(result.metadata.values[i]);
        }
        free(result.metadata.keys);
        free(result.metadata.values);
    } else {
        printf("  [No metadata]\n");
    }


    // now it's time to read actual data
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

    return 0;
}
