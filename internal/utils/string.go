package utils

import (
	"math/rand"
	"time"
)

func RandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	b := make([]byte, 7)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}

	return string(b)
}
