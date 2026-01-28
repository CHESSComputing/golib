package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// SplitStr2List splits on comma or whitespace and trims elements.
func SplitStr2List(val string) []string {
	val = strings.TrimSpace(val)
	if val == "" {
		return nil
	}

	var parts []string
	if strings.Contains(val, ",") {
		parts = strings.Split(val, ",")
	} else {
		parts = strings.Fields(val)
	}

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// Convert2dtype converts string value to proper data-type.
func Convert2dtype(val string, dtype string) (any, error) {
	val = strings.TrimSpace(val)

	switch dtype {

	case "string", "str":
		return val, nil

	// ---- integers ----
	case "int":
		i, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		return i, nil

	case "int8":
		i, err := strconv.ParseInt(val, 10, 8)
		if err != nil {
			return nil, err
		}
		return int8(i), nil

	case "int16":
		i, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			return nil, err
		}
		return int16(i), nil

	case "int32":
		i, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return nil, err
		}
		return int32(i), nil

	case "int64":
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		return i, nil

	// ---- floats ----
	case "float", "float64":
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, err
		}
		return f, nil

	case "float32":
		f, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return nil, err
		}
		return float32(f), nil

	// ---- bool ----
	case "bool":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
		return b, nil

	// ---- lists ----
	case "list_str":
		return SplitStr2List(val), nil

	case "list_int":
		parts := SplitStr2List(val)
		out := make([]int, 0, len(parts))
		for _, p := range parts {
			i, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			out = append(out, i)
		}
		return out, nil

	case "list_float":
		parts := SplitStr2List(val)
		out := make([]float64, 0, len(parts))
		for _, p := range parts {
			f, err := strconv.ParseFloat(p, 64)
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
	max := 0
	for _, vals := range input {
		if len(vals) > max {
			max = len(vals)
		}
	}

	records := make([]map[string]string, max)
	for i := 0; i < max; i++ {
		records[i] = make(map[string]string)
	}

	for key, vals := range input {
		for i, v := range vals {
			records[i][key] = v
		}
	}

	return records
}
