import h5py
import numpy as np

with h5py.File("sample.h5", "w") as f:
    f.attrs["Creator"] = "Test Suite"
    f.attrs["Version"] = "1.0"
    f.create_dataset("mydata", data=np.arange(100, dtype=np.float64))

