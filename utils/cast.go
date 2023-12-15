package utils

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// CastString function to check and cast interface{} to string data-type
func CastString(val interface{}) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	}
	msg := fmt.Sprintf("wrong data type for %v type %T", val, val)
	return "", errors.New(msg)
}

// CastInt function to check and cast interface{} to int data-type
func CastInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	}
	msg := fmt.Sprintf("wrong data type for %v type %T", val, val)
	return 0, errors.New(msg)
}

// CastInt64 function to check and cast interface{} to int64 data-type
func CastInt64(val interface{}) (int64, error) {
	switch v := val.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	}
	msg := fmt.Sprintf("wrong data type for %v type %T", val, val)
	return 0, errors.New(msg)
}

// CastFloat function to check and cast interface{} to int64 data-type
func CastFloat(val interface{}) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	}
	msg := fmt.Sprintf("wrong data type for %v type %T", val, val)
	return 0, errors.New(msg)
}

// ConvertFloat converts string representation of float scientific number to string int
func ConvertFloat(val string) string {
	if strings.Contains(val, "e+") || strings.Contains(val, "E+") {
		// we got float number, should be converted to int
		v, e := strconv.ParseFloat(val, 64)
		if e != nil {
			log.Println("unable to convert", val, " to float, error", e)
			return val
		}
		return strings.Split(fmt.Sprintf("%f", v), ".")[0]
	}
	return val
}
