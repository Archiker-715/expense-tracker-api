package auth

import "crypto/sha256"

func Encode(password string) [32]byte {
	return sha256.Sum256([]byte(password))
}
