//go:build linux
// +build linux

package utils

import (
	"syscall"
	"time"
)

// GetCreationTime returns the creation time for Linux
func GetCreationTime(st syscall.Stat_t) time.Time {
	return time.Unix(int64(st.Ctim.Sec), int64(st.Ctim.Nsec))
}
