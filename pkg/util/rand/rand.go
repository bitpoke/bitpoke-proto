package rand

import (
	crand "crypto/rand"
	"encoding/base64"
	"math/rand"
)

// GenerateRandomBytes returns securely generated random bytes.
func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := crand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		panic(err)
	}
	return b
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(n int) string {
	return base64.URLEncoding.EncodeToString(GenerateRandomBytes(n))
}

// GenerateRandomInt returns generated random int
func GenerateRandomInt(startAt, endAt int) int {
	return rand.Intn(endAt-startAt) + startAt
}

// GenerateRandomPort returns generated random port
func GenerateRandomPort() int {
	// Private ports are between 49152 and 65535
	firstPort := 49152
	lastPort := 65535
	return GenerateRandomInt(firstPort, lastPort)
}
