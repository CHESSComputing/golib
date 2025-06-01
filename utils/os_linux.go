//go:build linux
// +build linux

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

	switch st := fileInfo.Sys().(type) {
	case *syscall.Stat_t:
		return time.Unix(int64(st.Ctim.Sec), int64(st.Ctim.Nsec)), nil
	}

	fmt.Println("Creation time not available, using modification time as fallback.")
	return fileInfo.ModTime(), nil
}
