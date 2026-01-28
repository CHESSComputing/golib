package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// Convert2dtype converts string value to proper data-type.
func Convert2dtype(val string, dtype string) (any, error) {
	switch dtype {

	case "string", "str":
		return val, nil

	case "int", "int8", "int32", "int64":
		i, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		return i, nil

	case "float", "float32", "float64":
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, err
		}
		return f, nil

	case "bool":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
		return b, nil

	// arrays of primitives (optional but future-proof)
	case "list_str":
		var out []string
		if strings.Contains(val, ",") {
			for _, v := range strings.Split(val, ",") {
				out = append(out, strings.Trim(v, " "))
			}
		} else if strings.Contains(val, " ") {
			for _, v := range strings.Split(val, " ") {
				out = append(out, strings.Trim(v, " "))
			}
		}
		return out, nil

	case "list_int":
		out := make([]int, 0, len(val))
		sep := " "
		if strings.Contains(val, ",") {
			sep = ","
		}
		for _, v := range strings.Split(val, sep) {
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			out = append(out, i)
		}
		return out, nil

	case "list_float":
		out := make([]float64, 0, len(val))
		sep := " "
		if strings.Contains(val, ",") {
			sep = ","
		}
		for _, v := range strings.Split(val, sep) {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}
			out = append(out, f)
		}
		return out, nil
	}

	return nil, fmt.Errorf("unsupported data type %q", dtype)
}

// Convert2records converts input map of parsed web form values to list of records
func Convert2records(input map[string][]string) []map[string]string {
	// determine number of records
	max := 0
	for _, vals := range input {
		if len(vals) > max {
			max = len(vals)
		}
	}

	// allocate result slice
	records := make([]map[string]string, max)
	for i := 0; i < max; i++ {
		records[i] = make(map[string]string)
	}

	// populate records
	for key, vals := range input {
		for i, v := range vals {
			records[i][key] = v
		}
	}

	return records
}
