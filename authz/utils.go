package auth

import (
	"log"
	"os"
)

// ReadSecret provides unified way to read secret either from provided file
// or a string, and fall back to a default value if string is empty
func ReadSecret(r string) string {
	if _, err := os.Stat(r); err == nil {
		b, e := os.ReadFile(r)
		if e != nil {
			log.Fatalf("Unable to read data from file: %s, error: %s", r, e)
		}
		return string(b)
	}
	if r == "" {
		return "lkdsjflkjsdoiweuior"
	}
	return r
}
