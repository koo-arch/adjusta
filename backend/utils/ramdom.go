package utils

import (
	"crypto/rand"
	"encoding/base32"
)

func GenerateRamdomString(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return base32.StdEncoding.EncodeToString(b)
}
