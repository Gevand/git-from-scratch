package integration_test_lib

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz0987654321"

var seededRand *rand.Rand = rand.New(
	rand.NewSource((int64)(time.Now().UnixNano())))

func GenerateRandomString() string {
	length := 20
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
