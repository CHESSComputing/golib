package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"time"
)

// CodeVerifier generates a random code verifier of 43-128 characters
func CodeVerifier() string {
	rand.Seed(time.Now().UnixNano())
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := rand.Intn(86) + 43
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// CodeChallenge generates code challenge from the code verifier (SHA256 + base64 URL encoding)
func CodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	encoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash[:])
	return encoded
}
