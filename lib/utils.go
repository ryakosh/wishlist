package lib

import (
	"crypto/rand"
	"encoding/base64"
)

// GenSafeRandomBytes is a utility function used for cryptographically
// secure random bytes generation
func GenSafeRandomBytes(n int) ([]byte, int, error) {
	bytes := make([]byte, n)

	r, err := rand.Read(bytes)

	return bytes, r, err
}

// GenRandCode is a utility function used for cryptographically
// secure random code generation encoded in unpadded base64 format
func GenRandCode(n int) (string, error) {
	rands, _, err := GenSafeRandomBytes(n)
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(rands), nil
}
