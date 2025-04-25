#!/usr/bin/env python3

import h5py
import sys
import numpy as np # Used for checking dtypes and potentially handling data

# Define max elements to print fully to avoid overwhelming the console
MAX_ELEMENTS_TO_PRINT = 100

def print_nexus_item(name, obj):
    """
    Callback function for h5py's visititems.
    Prints information about the group or dataset object.
    """
    indent_level = name.count('/')
    indent = "  " * indent_level

    print(f"{indent}{name}:") # Print the full path of the item

    item_type = ""
    nx_class = obj.attrs.get('NX_class', None) # Check for NeXus class attribute

    # --- Identify Type (Group or Dataset) ---
    if isinstance(obj, h5py.Group):
        item_type = "Group"
        if nx_class:
            print(f"{indent}  Type: {item_type} (NX_class: {nx_class.decode() if isinstance(nx_class, bytes) else nx_class})")
        else:
            print(f"{indent}  Type: {item_type}")

    elif isinstance(obj, h5py.Dataset):
        item_type = "Dataset"
        if nx_class:
            print(f"{indent}  Type: {item_type} (NX_class: {nx_class.decode() if isinstance(nx_class, bytes) else nx_class})")
        else:
            print(f"{indent}  Type: {item_type}")
        # Print dataset-specific info
        print(f"{indent}  Shape: {obj.shape}")
        print(f"{indent}  Data Type: {obj.dtype}")
        print(f"{indent}  Total Elements: {obj.size}")

        # --- Print Data (with limits) ---
        try:
            if obj.size == 0:
                 print(f"{indent}  Data: (empty dataset)")
            elif obj.size <= MAX_ELEMENTS_TO_PRINT:
                 # Read and print all data if it's small
                 # Check if it's a string type that needs decoding
                 if h5py.check_string_dtype(obj.dtype):
                      # Use .asstr() context manager for robust string reading
                      with obj.astype(str) as str_data:
                           data_read = str_data[()]
                 else:
                      data_read = obj[()] # Read numerical/other data

                 # Format nicely, especially for multi-dimensional arrays
                 data_str = np.array2string(np.array(data_read), precision=4, suppress_small=True, separator=', ')
                 print(f"{indent}  Data:\n{indent}    {data_str}")
            else:
                 # Print only a slice if the dataset is large
                 print(f"{indent}  Data: (Dataset too large to print fully, showing slice)")
                 # Create a slice object for the first few elements (e.g., first 5 of first dim)
                 slice_indices = tuple([slice(0, 5)] + [slice(None)] * (obj.ndim - 1))
                 if obj.ndim == 0: # Scalar dataset but size > MAX_ELEMENTS? Unlikely but handle
                    slice_indices = ()

                 if h5py.check_string_dtype(obj.dtype):
                    with obj.astype(str) as str_data:
                         data_slice = str_data[slice_indices]
                 else:
                     data_slice = obj[slice_indices]

                 data_str = np.array2string(np.array(data_slice), precision=4, suppress_small=True, separator=', ')
                 print(f"{indent}    {data_str}...")

        except Exception as e:
            print(f"{indent}  Data: Error reading data - {e}")

    else:
        print(f"{indent}  Type: Unknown/Link?") # Could be an external link etc.

    # --- Print Attributes (Metadata) ---
    if obj.attrs:
        print(f"{indent}  Attributes:")
        for key, value in obj.attrs.items():
            # Decode byte strings if necessary for cleaner printing
            if isinstance(value, bytes):
                try:
                    value_str = value.decode('utf-8', errors='replace')
                except: # Fallback if decoding fails
                    value_str = str(value)
            else:
                value_str = str(value)

            # Add quotes around string attributes
            if isinstance(value, (str, bytes)):
                 value_str = f'"{value_str}"'

            print(f"{indent}    '{key}': {value_str}")

    print("-" * (len(indent) + 40)) # Separator


# --- Main Script ---
if __name__ == "__main__":
    if len(sys.argv) != 2:
        print(f"Usage: python {sys.argv[0]} <nexus_file.h5>")
        sys.exit(1)

    filepath = sys.argv[1]

    try:
        # Open the HDF5/NeXus file in read-only mode ('r')
        with h5py.File(filepath, 'r') as f:
            print(f"--- Reading NeXus File: {filepath} ---")
            # visititems traverses the hierarchy and calls the function for each item
            f.visititems(print_nexus_item)
            print(f"--- End of file contents ---")

    except FileNotFoundError:
        print(f"Error: File not found at '{filepath}'")
        sys.exit(1)
    except OSError as e:
        print(f"Error: Could not open or read file '{filepath}'. Is it a valid HDF5/NeXus file?")
        print(f"       ({e})")
        sys.exit(1)
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
        sys.exit(1)
