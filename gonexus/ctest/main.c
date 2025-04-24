#include <stdio.h>
#include "reader.h"

int main() {
    const char *filename = "sample.h5";
    const char *dataset = "mydata";

    HDF5Result result = read_hdf5(filename, dataset);

    printf("Metadata:\n");
    for (int i = 0; i < result.metadata.count; i++) {
        printf("  %s: %s\n", result.metadata.keys[i], result.metadata.values[i]);
    }

    printf("\nData:\n");
    for (int i = 0; i < result.dataset.length; i++) {
/*        printf("  %f\n", result.dataset.data.values[i]);*/
        printf("  %f\n", result.dataset.data[i]);
    }

    free_result(result);
    return 0;
}

