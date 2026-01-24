package auth

import "crypto/sha256"

func Encode(password string) []byte {
	hashed := sha256.Sum256([]byte(password))
	return hashed[:]
}
