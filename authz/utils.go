package auth

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
)

// RandomBytes generates random bytes from given size and seed
func RandomBytes(size int, seed int64) []byte {
	src := rand.NewSource(seed)
	r := rand.New(src)
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[i] = byte(r.Intn(256))
	}
	return bytes
}

// RandomString generates random string using given seed and size
func RandomString(size int, seed int64) string {
	if size == 0 {
		size = 16
	}
	var iseed int64
	if seed == 0 {
		binary.Read(cryptorand.Reader, binary.LittleEndian, &iseed)
	} else {
		iseed = seed
	}
	return fmt.Sprintf("%x", RandomBytes(size, iseed))[:size]
}

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
		return RandomString(16, 123456789)
	}
	return r
}
