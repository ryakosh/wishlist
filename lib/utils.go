package lib

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

const (
	LError = "ERROR"
	LPanic = "PANIC"
	LFatal = "FATAL"
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

func LogError(lev, sum string, reason error) {
	if lev == LError {
		log.Printf("%s: %s\n\treason: %s", lev, sum, reason)
	} else if lev == LPanic {
		log.Panicf("%s: %s\n\treason: %s", lev, sum, reason)
	} else if lev == LFatal {
		log.Fatalf("%s: %s\n\treason: %s", lev, sum, reason)
	}
}
