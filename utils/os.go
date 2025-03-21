package utils

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// DirCreationDate returns the creation time of a directory,
// or modification time as a fallback if creation time is unavailable.
func DirCreationDate(dir string) (time.Time, error) {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return time.Time{}, err
	}

	// Try extracting creation time based on OS
	switch sys := fileInfo.Sys().(type) {
	case *syscall.Stat_t: // Linux & macOS (Unix systems)
		if sys.Ctimespec.Sec > 0 {
			return time.Unix(sys.Ctimespec.Sec, sys.Ctimespec.Nsec), nil
		}
	}

	// Fallback to modification time
	fmt.Println("Creation time not available, using modification time as fallback.")
	return fileInfo.ModTime(), nil
}
