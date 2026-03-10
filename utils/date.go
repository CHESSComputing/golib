package utils

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// ConvertTimestamp converts given input into human readable timestamp
func ConvertTimestamp(val any) string {
	switch v := val.(type) {

	// ---- Signed integers ----
	case int:
		return formatUnix(int64(v))
	case int8:
		return formatUnix(int64(v))
	case int16:
		return formatUnix(int64(v))
	case int32:
		return formatUnix(int64(v))
	case int64:
		return formatUnix(v)

	// ---- Unsigned integers ----
	case uint:
		return formatUnix(int64(v))
	case uint8:
		return formatUnix(int64(v))
	case uint16:
		return formatUnix(int64(v))
	case uint32:
		return formatUnix(int64(v))
	case uint64:
		if v > math.MaxInt64 {
			return fmt.Sprintf("%v", v)
		}
		return formatUnix(int64(v))

	// ---- Floats ----
	case float32:
		return formatUnix(int64(v))
	case float64:
		return formatUnix(int64(v))

	// ---- String ----
	case string:
		// Try parsing as integer timestamp first
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return formatUnix(i)
		}
		// Otherwise return original string
		return v
	}

	return fmt.Sprintf("%v", val)
}

// formatUnix detects seconds/ms/us/ns and formats
func formatUnix(ts int64) string {
	switch {
	case ts > 1e18: // nanoseconds
		return time.Unix(0, ts).UTC().Format(time.RFC3339)
	case ts > 1e15: // microseconds
		return time.Unix(0, ts*1e3).UTC().Format(time.RFC3339)
	case ts > 1e12: // milliseconds
		return time.UnixMilli(ts).UTC().Format(time.RFC3339)
	default: // seconds
		return time.Unix(ts, 0).UTC().Format(time.RFC3339)
	}
}

// RFC3339ToEpoch converts RFC3339 time string into seconds since epoch
func RFC3339ToEpoch(s string) (int64, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return 0, fmt.Errorf("[golib.utils.RFC3339ToEpoch] time.Parse error: %w", err)
	}
	return t.Unix(), nil
}
