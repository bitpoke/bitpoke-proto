package rand

import (
	rand "crypto/rand"
	"encoding/base64"
)

// GenerateRandomBytes returns securely generated random bytes.
func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		panic(err)
	}
	return b
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) string {
	return base64.URLEncoding.EncodeToString(GenerateRandomBytes(s))
}
