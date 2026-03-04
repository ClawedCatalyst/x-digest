package application

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// RandURLSafe returns a URL-safe random string of n bytes (base64 encoded).
func RandURLSafe(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func pkceChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}
