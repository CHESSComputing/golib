//go:build windows
// +build windows

package utils

import (
	"fmt"
	"os"
	"time"
)

func DirCreationDate(dir string) (time.Time, error) {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return time.Time{}, err
	}

	fmt.Println("Creation time not available on Windows, using modification time.")
	return fileInfo.ModTime(), nil
}

