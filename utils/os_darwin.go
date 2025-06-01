//go:build darwin
// +build darwin

package utils

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func DirCreationDate(dir string) (time.Time, error) {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return time.Time{}, err
	}
	if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
		return time.Unix(int64(stat.Birthtimespec.Sec), int64(stat.Birthtimespec.Nsec)), nil
	}
	fmt.Println("Using modification time as fallback")
	return fileInfo.ModTime(), nil
}

