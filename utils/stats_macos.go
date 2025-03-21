//go:build darwin
// +build darwin

package utils

import (
	"syscall"
	"time"
)

// GetCreationTime returns the creation time for macOS
func GetCreationTime(st syscall.Stat_t) time.Time {
	return time.Unix(int64(st.Ctimespec.Sec), int64(st.Ctimespec.Nsec))
}
