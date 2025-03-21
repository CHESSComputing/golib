package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// CountEntries returns the total number of files and directories in a given directory.
func CountEntries(dirPath string) (int, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return 0, err
	}
	defer dir.Close()

	files, err := dir.Readdirnames(-1) // Read all entries
	if err != nil {
		return 0, err
	}

	return len(files), nil
}

// DirExists checks if a given directory exists
func DirExists(path string) bool {
	info, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	stat, err := os.Stat(info)
	return err == nil && stat.IsDir()
}

// FindMatches searches for all files and directories matching a given pattern within dir
func FindMatches(dir, pattern string) ([]string, error) {
	var matches []string

	// Walk through the directory recursively
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file/directory name matches the pattern
		matched, err := filepath.Match(pattern, info.Name())
		if err != nil {
			return err
		}

		if matched {
			matches = append(matches, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return matches, nil
}

// FindMatchingDirectories searches for directories matching a given pattern within dir.
func FindMatchingDirectories(dir, pattern string) ([]string, error) {
	var matches []string

	dirs, err := FindAllDirectories(dir)
	if err != nil {
		return matches, err
	}
	for _, d := range dirs {
		if strings.Contains(d, pattern) {
			matches = append(matches, d)
		}
	}
	return matches, nil
}

// FindAllDirectories returns a list of final directories (directories that do not contain other directories)
func FindAllDirectories(rootDir string) ([]string, error) {
	var finalDirs []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if it's not a directory
		if !info.IsDir() {
			return nil
		}

		// Check if the directory has any subdirectories
		hasSubDir := false
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				hasSubDir = true
				break
			}
		}

		// If no subdirectories found, add to finalDirs
		if !hasSubDir {
			relPath, err := filepath.Rel(rootDir, path)
			if err != nil {
				return err
			}
			finalDirs = append(finalDirs, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return finalDirs, nil
}
