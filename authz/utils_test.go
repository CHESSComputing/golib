package auth

import (
	"bytes"
	"fmt"
	"testing"
)

// TestRandomString
func TestRandomString(t *testing.T) {
	size := 10
	seed := int64(1)
	s1 := RandomString(size, seed)
	s2 := RandomString(size, seed)
	if s1 != s2 {
		t.Error("RandomString generator is not persistent")
	}
	fmt.Println("random strings:", s1, len(s1), size)
	fmt.Println("random strings:", s2, len(s2), size)
	if len(s1) != size {
		t.Error("RandomString string size test failure")
	}
}

// TestRandomBytes
func TestRandomBytes(t *testing.T) {
	size := 10
	seed := int64(1)
	b1 := RandomBytes(size, seed)
	b2 := RandomBytes(size, seed)
	if !bytes.Equal(b1, b2) {
		t.Error("RandomBytes generator is not persistent")
	}
	if len(b1) != size {
		t.Error("RandomBytes string size test failure")
	}
}

// TestReadSecret
func TestReadSecret(t *testing.T) {
	s1 := ReadSecret("")
	s2 := ReadSecret("test")
	if s1 == s2 {
		t.Error("secrets are the same")
	}
	if s2 != "test" {
		t.Error("fail to read given secret")
	}
}
